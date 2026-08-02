#pragma once

#include "antlr4-runtime.h"
#include "ast/ast.h"
#include "cpp_parser/ast_builder.h"
#include "cpp_parser/error_listener.h"
#include "cpp_parser/generated/bapelLexer.h"
#include "cpp_parser/generated/bapelParser.h"

#include <fstream>
#include <sstream>
#include <string>
#include <utility>
#include <vector>

namespace parser {

template <typename T>
struct ParseResult {
  bool ok = false;
  T value{};
  std::vector<std::string> errors;

  std::string error() const {
    return errors.empty() ? "" : errors.front();
  }

  explicit operator bool() const { return ok; }
};

inline ParseResult<ast::SourceFile> parse_source_file(const std::string& input_code,
                                                      const std::string& filename = "<inline>") {
  ParseResult<ast::SourceFile> result;
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
    result.ok = false;
    result.errors = error_listener.errors();
    return result;
  }

  ast::AstBuilder builder(filename);
  ast::SourceFile sf = tree->accept(&builder).as<ast::SourceFile>();
  sf.header.filename = ir::new_filename(filename, ir::Pos{});

  std::vector<std::string> val_errors;
  if (!ast::validate_source_file(sf, val_errors)) {
    result.ok = false;
    result.errors = std::move(val_errors);
    return result;
  }

  result.ok = true;
  result.value = std::move(sf);
  return result;
}

inline ParseResult<ast::SourceFile> parse_source_file_from_file(const std::string& filepath) {
  std::ifstream file(filepath);
  if (!file.is_open()) {
    ParseResult<ast::SourceFile> result;
    result.ok = false;
    result.errors = {"failed to open file: " + filepath};
    return result;
  }
  std::stringstream buffer;
  buffer << file.rdbuf();
  return parse_source_file(buffer.str(), filepath);
}

inline ParseResult<ast::Workspace> parse_workspace(const std::string& input_code,
                                                   const std::string& filename = "<inline>") {
  ParseResult<ast::Workspace> result;
  antlr4::ANTLRInputStream stream(input_code);
  bapelLexer lexer(&stream);
  ast::BapelErrorListener error_listener(filename);
  lexer.removeErrorListeners();
  lexer.addErrorListener(&error_listener);

  antlr4::CommonTokenStream tokens(&lexer);
  bapelParser parser(&tokens);
  parser.removeErrorListeners();
  parser.addErrorListener(&error_listener);

  auto* tree = parser.workspace();
  if (error_listener.has_errors()) {
    result.ok = false;
    result.errors = error_listener.errors();
    return result;
  }

  ast::AstBuilder builder(filename);
  ast::Workspace ws = tree->accept(&builder).as<ast::Workspace>();

  std::vector<std::string> val_errors;
  if (!ast::validate_workspace(ws, val_errors)) {
    result.ok = false;
    result.errors = std::move(val_errors);
    return result;
  }

  result.ok = true;
  result.value = std::move(ws);
  return result;
}

inline ParseResult<ast::Workspace> parse_workspace_from_file(const std::string& filepath) {
  std::ifstream file(filepath);
  if (!file.is_open()) {
    ParseResult<ast::Workspace> result;
    result.ok = false;
    result.errors = {"failed to open file: " + filepath};
    return result;
  }
  std::stringstream buffer;
  buffer << file.rdbuf();
  return parse_workspace(buffer.str(), filepath);
}

inline ParseResult<ir::IrDecl> parse_decl(const std::string& input_code,
                                          const std::string& filename = "<inline>") {
  ParseResult<ir::IrDecl> result;
  antlr4::ANTLRInputStream stream(input_code);
  bapelLexer lexer(&stream);
  ast::BapelErrorListener error_listener(filename);
  lexer.removeErrorListeners();
  lexer.addErrorListener(&error_listener);

  antlr4::CommonTokenStream tokens(&lexer);
  bapelParser parser(&tokens);
  parser.removeErrorListeners();
  parser.addErrorListener(&error_listener);

  auto* tree = parser.decl();
  if (error_listener.has_errors()) {
    result.ok = false;
    result.errors = error_listener.errors();
    return result;
  }

  ast::AstBuilder builder(filename);
  ir::IrDecl decl = tree->accept(&builder).as<ir::IrDecl>();

  result.ok = true;
  result.value = std::move(decl);
  return result;
}

} // namespace parser
