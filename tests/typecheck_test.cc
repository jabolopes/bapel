#include "tests/test_util.h"
#include "comp/typecheck_unit.h"
#include "comp/querier.h"
#include "comp/module_finder.h"

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

      ir::IrUnit unit;
      std::string err;
      bool ok = comp::typecheck_source_file(querier, options, inFile, unit, err);

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

      std::string got = ok ? unit.to_bpl_string(true) : (err + "\n");

      std::string diff_str;
      if (!tests::diff_out_regen(got, wantFile, diff_str)) {
        sub_ctx.add_error("in test " + inFile + ":\n" + diff_str, __FILE__, __LINE__);
      }
    });
  }
}
