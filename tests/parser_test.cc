#include "tests/test_util.h"

#include <cstdio>
#include <memory>
#include <string>
#include <vector>

namespace {

struct ProcessResult {
  int exit_code = 0;
  std::string output;
};

ProcessResult run_parser(const std::string& input_file, bool with_pos) {
  std::string parser_bin = "./bootstrap/parser";
  if (!tests::fs::exists(parser_bin)) {
    parser_bin = "bootstrap/parser";
  }

  std::string cmd = parser_bin + " --symbol=SourceFile --format=bpl" +
                    (with_pos ? " --with-pos " : " ") + input_file + " 2>&1";

  FILE* pipe = popen(cmd.c_str(), "r");
  if (!pipe) {
    return {-1, "failed to popen: " + cmd};
  }

  char buffer[4096];
  std::string result;
  while (fgets(buffer, sizeof(buffer), pipe) != nullptr) {
    result += buffer;
  }

  int status = pclose(pipe);
  int exit_code = WIFEXITED(status) ? WEXITSTATUS(status) : -1;
  return {exit_code, result};
}

} // namespace

TEST(ParserTest, GoldenFiles) {
  auto matches = tests::glob("tests/testdata/parse/in/*.in");
  ASSERT_FALSE(matches.empty());

  for (const auto& inFile : matches) {
    ctx.run(inFile, [&](tests::TestContext& sub_ctx) {
      std::string wantFile = tests::replace_string(tests::replace_extension(inFile, ".bpl"), "/in/", "/parsed/");
      bool wantErr = (tests::path_base(inFile).rfind("bad_", 0) == 0);

      auto res = run_parser(inFile, true /* with_pos */);

      if (wantErr) {
        if (res.exit_code == 0) {
          sub_ctx.add_error("expected parse error for " + inFile + " but parser succeeded", __FILE__, __LINE__);
          return;
        }
      } else {
        if (res.exit_code != 0) {
          sub_ctx.add_error("parser failed on " + inFile + ":\n" + res.output, __FILE__, __LINE__);
          return;
        }
      }

      std::string got = res.output;
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
