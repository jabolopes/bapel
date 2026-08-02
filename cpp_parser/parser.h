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

namespace bapel_parser {

static inline void replace_all(std::string& str, const std::string& from, const std::string& to) {
  size_t start_pos = 0;
  while ((start_pos = str.find(from, start_pos)) != std::string::npos) {
    str.replace(start_pos, from.length(), to);
    start_pos += to.length();
  }
}

// @bpl: pub bapel_parser::run: Vector String -> (i64, String)
inline std::tuple<int64_t, std::string> run(const std::vector<std::string>& args) {
  std::string symbol = "SourceFile";
  std::string format = "json";
  bool with_pos = false;
  std::string filename;
  std::string input_path;

  for (size_t i = 0; i < args.size(); ++i) {
    const std::string& arg = args[i];
    if (arg == "--with-pos" || arg == "-with-pos") {
      with_pos = true;
    } else if (arg.rfind("--symbol=", 0) == 0) {
      symbol = arg.substr(9);
    } else if (arg.rfind("-symbol=", 0) == 0) {
      symbol = arg.substr(8);
    } else if (arg == "--symbol" || arg == "-symbol") {
      if (i + 1 < args.size()) symbol = args[++i];
    } else if (arg.rfind("--format=", 0) == 0) {
      format = arg.substr(9);
    } else if (arg.rfind("-format=", 0) == 0) {
      format = arg.substr(8);
    } else if (arg == "--format" || arg == "-format") {
      if (i + 1 < args.size()) format = args[++i];
    } else if (arg.rfind("--filename=", 0) == 0) {
      filename = arg.substr(11);
    } else if (arg.rfind("-filename=", 0) == 0) {
      filename = arg.substr(10);
    } else if (arg == "--filename" || arg == "-filename") {
      if (i + 1 < args.size()) filename = args[++i];
    } else if (arg == "--workspace" || arg == "-workspace") {
      symbol = "Workspace";
    } else if (arg.rfind("-", 0) != 0) {
      if (input_path.empty()) {
        input_path = arg;
      } else {
        return {1, "Expected at most one argument\n"};
      }
    }
  }

  std::string input_code;
  if (!input_path.empty()) {
    std::ifstream file(input_path);
    if (file.is_open()) {
      std::stringstream buffer;
      buffer << file.rdbuf();
      input_code = buffer.str();
      if (filename.empty()) {
        filename = input_path;
      }
    } else {
      input_code = input_path;
      if (filename.empty()) {
        filename = "<inline>";
      }
    }
  } else {
    filename = "<inline>";
  }

  std::stringstream out;
  if (symbol == "SourceFile") {
    if ((filename.size() >= 3 && filename.substr(filename.size() - 3) == ".cc") ||
        (filename.size() >= 4 && filename.substr(filename.size() - 4) == ".cpp") ||
        (filename.size() >= 4 && filename.substr(filename.size() - 4) == ".cxx") ||
        (filename.size() >= 2 && filename.substr(filename.size() - 2) == ".c")) {
      if (format == "json") {
        ast::SourceFile sf;
        out << sf.to_json() << "\n";
      }
      return {0, out.str()};
    }
    auto res = parser::parse_source_file(input_code, filename);
    if (!res.ok) {
      std::string err_out;
      for (const auto& err : res.errors) {
        err_out += err + "\n";
      }
      return {1, err_out};
    }
    const auto& sf = res.value;

    if (format == "bpl") {
      out << sf.to_string(with_pos) << "\n";
      return {0, out.str()};
    }

    if (format == "flat") {
      for (const auto& imp : sf.imports.ids) {
        out << "IMPORT " << imp.name << "\n";
      }
      for (const auto& impl : sf.impls.filenames) {
        out << "IMPL " << impl.value << "\n";
      }
      for (const auto& flag : sf.flags.filenames) {
        out << "FLAG " << flag.value << "\n";
      }
      for (const auto& source : sf.body) {
        if (source.case_val == ast::SourceCase::DeclSource && source.decl_data) {
          std::string s = source.decl_data->to_string();
          replace_all(s, "\n", "\\n");
          out << "DECL " << s << "\n";
        } else if (source.case_val == ast::SourceCase::FunctionSource && source.function_data) {
          std::string s = source.function_data->decl().to_string();
          replace_all(s, "\n", "\\n");
          out << "DECL " << s << "\n";
          std::string fn_str = source.function_data->to_string(false);
          replace_all(fn_str, "\n", "\\n");
          out << "FUNC " << fn_str << "\n";
        } else if (source.case_val == ast::SourceCase::TraitSource && source.trait_data) {
          std::string s = source.trait_data->decl().to_string();
          replace_all(s, "\n", "\\n");
          out << "DECL " << s << "\n";
        } else if (source.case_val == ast::SourceCase::ImplSource && source.impl_data) {
          std::string s = source.impl_data->to_string(false);
          replace_all(s, "\n", "\\n");
          out << "TRAIT_IMPL " << s << "\n";
        }
      }
      return {0, out.str()};
    }

    out << sf.to_json() << "\n";
    return {0, out.str()};
  } else if (symbol == "Workspace") {
    auto res = parser::parse_workspace(input_code, filename);
    if (!res.ok) {
      std::string err_out;
      for (const auto& err : res.errors) {
        err_out += err + "\n";
      }
      return {1, err_out};
    }
    const auto& ws = res.value;

    if (format == "flat") {
      for (const auto& pkg : ws.packages.packages) {
        if (pkg.case_val == ast::PackageCase::PrefixPackage) {
          out << "PREFIX " << pkg.module_id.name << " " << pkg.filename.value << "\n";
        } else if (pkg.case_val == ast::PackageCase::ModulePackage) {
          out << "MODULE " << pkg.module_id.name << " " << pkg.filename.value << "\n";
        }
      }
      return {0, out.str()};
    }

    out << ws.to_json() << "\n";
    return {0, out.str()};
  } else if (symbol == "Decl") {
    auto res = parser::parse_decl(input_code, filename);
    if (!res.ok) {
      std::string err_out;
      for (const auto& err : res.errors) {
        err_out += err + "\n";
      }
      return {1, err_out};
    }
    const auto& decl = res.value;

    if (format == "flat") {
      std::string s = decl.to_string();
      replace_all(s, "\n", "\\n");
      out << "DECL " << s << "\n";
      return {0, out.str()};
    }
    out << decl.to_json() << "\n";
    return {0, out.str()};
  }

  return {1, "Unsupported symbol \"" + symbol + "\"\n"};
}

} // namespace bapel_parser
