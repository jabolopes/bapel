#include "tests/test_util.h"
#include "comp/typecheck_unit.h"
#include "comp/querier.h"
#include "comp/module_finder.h"

#include <string>
#include <vector>

TEST(StlcTest, InferTerm) {
  auto matches = tests::glob("tests/testdata/stlc/*.stlc.in");
  ASSERT_FALSE(matches.empty());

  comp::ModuleFinder finder({}, {{"", "."}});
  comp::Querier querier(finder);

  for (const auto& inFile : matches) {
    ctx.run(inFile, [&](tests::TestContext& sub_ctx) {
      std::string wantFile = tests::replace_extension(inFile, ".out");

      comp::TypecheckOptions options;
      options.skip_default_context = true;
      options.skip_term_typechecker = true;
      options.skip_undefined_term_checks = true;

      ir::IrUnit unit;
      std::string err;
      bool ok = comp::typecheck_source_file(querier, options, inFile, unit, err);

      std::string got;
      if (ok) {
        got = unit.to_bpl_string(true);
      } else {
        got = err + "\n";
      }

      std::string diff_str;
      if (!tests::diff_out_regen(got, wantFile, diff_str)) {
        sub_ctx.add_error("in test " + inFile + ":\n" + diff_str, __FILE__, __LINE__);
      }
    });
  }
}

TEST(StlcTest, TypecheckTerm) {
  auto matches = tests::glob("tests/testdata/stlc/*.stlc.in");
  ASSERT_FALSE(matches.empty());

  comp::ModuleFinder finder({}, {{"", "."}});
  comp::Querier querier(finder);

  for (const auto& inFile : matches) {
    ctx.run(inFile, [&](tests::TestContext& sub_ctx) {
      comp::TypecheckOptions options;
      options.skip_default_context = true;
      options.skip_term_typechecker = false;
      options.skip_undefined_term_checks = true;

      ir::IrUnit unit;
      std::string err;
      bool ok = comp::typecheck_source_file(querier, options, inFile, unit, err);
      if (!ok) {
        sub_ctx.add_error("in test " + inFile + ": " + err, __FILE__, __LINE__);
      }
    });
  }
}
