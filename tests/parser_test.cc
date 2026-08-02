#include "tests/test_util.h"
#include "antlr4-runtime.h"
#include "ast/ast.h"
#include "cpp_parser/ast_builder.h"
#include "cpp_parser/error_listener.h"
#include "cpp_parser/generated/bapelLexer.h"
#include "cpp_parser/generated/bapelParser.h"

#include <fstream>
#include <memory>
#include <sstream>
#include <string>
#include <vector>

namespace {

struct ProcessResult {
  int exit_code = 0;
  std::string output;
};

ProcessResult run_parser(const std::string& input_file, bool with_pos) {
  std::ifstream file(input_file);
  if (!file.is_open()) {
    return {-1, "failed to open file: " + input_file};
  }

  std::stringstream buffer;
  buffer << file.rdbuf();
  std::string input_code = buffer.str();
  std::string filename = input_file;

  antlr4::ANTLRInputStream stream(input_code);
  bapelLexer lexer(&stream);
  ast::BapelErrorListener error_listener(filename);
  lexer.removeErrorListeners();
  lexer.addErrorListener(&error_listener);

  antlr4::CommonTokenStream tokens(&lexer);
  bapelParser parser(&tokens);
  parser.removeErrorListeners();
  parser.addErrorListener(&error_listener);

  auto* tree = parser.sourceFile();

  if (error_listener.has_errors()) {
    return {1, error_listener.errors()[0]};
  }

  ast::AstBuilder builder(filename);
  ast::SourceFile sf = tree->accept(&builder).as<ast::SourceFile>();
  sf.header.filename = ir::new_filename(filename, ir::Pos{});
  std::vector<std::string> val_errors;
  if (!ast::validate_source_file(sf, val_errors)) {
    std::string err_str;
    for (const auto& err : val_errors) {
      if (!err_str.empty()) err_str += "\n";
      err_str += err;
    }
    return {1, err_str};
  }

  return {0, sf.to_string(with_pos)};
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
