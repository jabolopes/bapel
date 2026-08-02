#include "tests/test_util.h"
#include "cpp_parser/parser.h"

#include <string>

TEST(ParserTest, GoldenFiles) {
  auto matches = tests::glob("tests/testdata/parse/in/*.in");
  ASSERT_FALSE(matches.empty());

  for (const auto& inFile : matches) {
    ctx.run(inFile, [&](tests::TestContext& sub_ctx) {
      auto directives = tests::TestDirectives::parse_from_file(inFile);
      if (!directives.should_run_stage("parse")) return;

      std::string wantFile = tests::replace_string(tests::replace_extension(inFile, ".bpl"), "/in/", "/parsed/");
      auto res = parser::parse_source_file_from_file(inFile);

      if (directives.expects_error("parse")) {
        if (res.ok) {
          sub_ctx.add_error("expected parse error for " + inFile + " but parser succeeded", __FILE__, __LINE__);
          return;
        }
      } else {
        if (!res.ok) {
          sub_ctx.add_error("parser failed on " + inFile + ":\n" + res.error(), __FILE__, __LINE__);
          return;
        }
      }

      std::string got = res.ok ? res.value.to_string(true /* with_pos */) : res.error();
      // Ensure trailing newline
      if (got.empty() || got.back() != '\n') {
        got += "\n";
      }

      std::string diff_str;
      if (!tests::diff_out_regen(got, wantFile, diff_str)) {
        sub_ctx.add_error("in test " + inFile + ":\n" + diff_str, __FILE__, __LINE__);
      }
    });
  }
}

