#include "tests/test_util.h"
#include "comp/typecheck_unit.h"
#include "comp/querier.h"
#include "comp/module_finder.h"

#include <filesystem>
#include <fstream>
#include <string>
#include <vector>

TEST(TypecheckTest, GoldenFiles) {
  auto matches = tests::glob("tests/testdata/in/*.in");
  ASSERT_FALSE(matches.empty());

  comp::ModuleFinder finder({}, {{"", "."}});
  comp::Querier querier(finder);

  for (const auto& inFile : matches) {
    ctx.run(inFile, [&](tests::TestContext& sub_ctx) {
      auto directives = tests::TestDirectives::parse_from_file(inFile);
      if (!directives.should_run_stage("typecheck")) return;

      std::string wantFile = tests::replace_string(tests::replace_extension(inFile, ".out"), "tests/testdata/in/", "tests/testdata/typecheck/");

      comp::TypecheckOptions options = directives.typecheck_options;
      options.skip_term_typechecker = false;

      ir::IrModule mod;
      std::string err;
      bool ok = comp::typecheck_module_from_file(querier, options, inFile, mod, err);

      if (directives.expects_error("typecheck")) {
        if (ok) {
          sub_ctx.add_error("expected typecheck error for " + inFile + " but typechecking succeeded", __FILE__, __LINE__);
          return;
        }
      } else {
        if (!ok) {
          sub_ctx.add_error("typecheck failed on " + inFile + ":\n" + err, __FILE__, __LINE__);
          return;
        }
      }

      std::string got = ok ? mod.to_bpl_string(true) : (err + "\n");

      std::string diff_str;
      if (!tests::diff_out_regen(got, wantFile, diff_str)) {
        sub_ctx.add_error("in test " + inFile + ":\n" + diff_str, __FILE__, __LINE__);
      }
    });
  }
}

TEST(TypecheckTest, MultiFileModule) {
  namespace fs = std::filesystem;
  std::string tmp_dir = "/tmp/bapel_multi_file_test";
  fs::create_directories(tmp_dir);

  auto write_file = [&](const std::string& path, const std::string& content) {
    std::ofstream ofs(path);
    ofs << content;
  };

  comp::ModuleFinder finder({}, {{"", "."}});
  comp::Querier querier(finder);
  comp::TypecheckOptions options;

  // 1. Equal-weight type, fn, trait, and impl blocks across base and impl files,
  //    plus per-file import scoping (impl imports bapel.stl, base does not).
  std::string ok_base = tmp_dir + "/ok_mod.bpl";
  std::string ok_impl = tmp_dir + "/ok_mod_impl.bpl";
  write_file(ok_base,
      "module ok_mod\n"
      "impls {\n"
      "  \"ok_mod_impl.bpl\"\n"
      "}\n"
      "pub type Boxed = struct { val: i64 }\n"
      "pub fn use_both(b: Boxed, s: Label) -> i64 {\n"
      "  b.get () + s.tag ()\n"
      "}\n");
  write_file(ok_impl,
      "implements ok_mod\n"
      "imports {\n"
      "  bapel.core\n"
      "  bapel.stl\n"
      "}\n"
      "pub type Label = struct { id: i64 }\n"
      "pub fn make_str() -> String { String::from_view \"hi\" }\n"
      "pub trait Tagged {\n"
      "  fn tag(self: Self) -> i64\n"
      "}\n"
      "impl Boxed {\n"
      "  fn get(self: Self) -> i64 { self.val }\n"
      "}\n"
      "impl Tagged for Label {\n"
      "  fn tag(self: Self) -> i64 { self.id }\n"
      "}\n");

  {
    ir::IrModule mod;
    std::string err;
    bool ok = comp::typecheck_module_from_file(querier, options, ok_base, mod, err);
    ASSERT_TRUE(ok);
    ASSERT_EQ(mod.files.size(), static_cast<size_t>(2));
    ASSERT_TRUE(mod.files[0].case_val == ir::IrFileCase::Base);
    ASSERT_TRUE(mod.files[1].case_val == ir::IrFileCase::Impl);
  }

  // 2. Per-file import scoping: base file must NOT see String imported only by impl file.
  std::string bad_base = tmp_dir + "/bad_scope.bpl";
  std::string bad_impl = tmp_dir + "/bad_scope_impl.bpl";
  write_file(bad_base,
      "module bad_scope\n"
      "impls {\n"
      "  \"bad_scope_impl.bpl\"\n"
      "}\n"
      "pub type Leak = struct { s: String }\n");
  write_file(bad_impl,
      "implements bad_scope\n"
      "imports {\n"
      "  bapel.stl\n"
      "}\n"
      "pub type Ok = struct { s: String }\n");

  {
    ir::IrModule mod;
    std::string err;
    bool ok = comp::typecheck_module_from_file(querier, options, bad_base, mod, err);
    ASSERT_FALSE(ok);
    ASSERT_TRUE(err.find("String") != std::string::npos);
  }

  // 3. Whole-module typechecking of bapel/compiler/ir.bpl (6 files).
  {
    ir::IrModule mod;
    std::string err;
    bool ok = comp::typecheck_module_from_file(querier, options, "bapel/compiler/ir.bpl", mod, err);
    ASSERT_TRUE(ok);
    ASSERT_EQ(mod.files.size(), static_cast<size_t>(6));
  }
}
