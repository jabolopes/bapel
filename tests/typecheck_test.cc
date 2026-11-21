#include "tests/test_util.h"
#include "comp/typecheck_module.h"
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

      ir::IrTypedModule mod;
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
    ir::IrTypedModule mod;
    std::string err;
    bool ok = comp::typecheck_module_from_file(querier, options, ok_base, mod, err);
    ASSERT_TRUE(ok);
    ASSERT_EQ(mod.functions.size(), static_cast<size_t>(2));
    ASSERT_EQ(mod.trait_impls.size(), static_cast<size_t>(2));
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
    ir::IrTypedModule mod;
    std::string err;
    bool ok = comp::typecheck_module_from_file(querier, options, bad_base, mod, err);
    ASSERT_FALSE(ok);
    ASSERT_TRUE(err.find("String") != std::string::npos);
  }

  // 3. Whole-module typechecking of bapel/compiler/ir.bpl (6 files).
  {
    ir::IrTypedModule mod;
    std::string err;
    bool ok = comp::typecheck_module_from_file(querier, options, "bapel/compiler/ir.bpl", mod, err);
    ASSERT_TRUE(ok);
    ASSERT_EQ(mod.impls.size(), static_cast<size_t>(5));
  }

  // 4. Imported ImplDecl export filtering: pub type impls are exported (signatures only),
  //    private type impls are NOT exported.
  std::string exp_mod = tmp_dir + "/exp_mod.bpl";
  std::string imp_mod = tmp_dir + "/imp_mod.bpl";
  write_file(exp_mod,
      "module exp_mod\n"
      "pub type PubT = struct { v: i64 }\n"
      "type PrivT = struct { v: i64 }\n"
      "impl PubT {\n"
      "  fn val(self: Self) -> i64 { self.v }\n"
      "}\n"
      "impl PrivT {\n"
      "  fn secret(self: Self) -> i64 { self.v }\n"
      "}\n");
  write_file(imp_mod,
      "module imp_mod\n"
      "imports {\n"
      "  exp_mod\n"
      "}\n"
      "pub fn read_pub(x: PubT) -> i64 { x.val () }\n");

  {
    comp::ModuleFinder local_finder({{"exp_mod", tmp_dir}}, {{"", "."}});
    comp::Querier local_querier(local_finder);
    ir::IrTypedModule mod;
    std::string err;
    bool ok = comp::typecheck_module_from_file(local_querier, options, imp_mod, mod, err);
    ASSERT_TRUE(ok);
    std::string printed = mod.to_bpl_string(false);
    ASSERT_TRUE(printed.find("impl PubT") != std::string::npos);
    ASSERT_TRUE(printed.find("PrivT") == std::string::npos);
  }
}
