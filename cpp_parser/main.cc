#include "antlr4-runtime.h"
#include "ast/ast.h"
#include "ast_builder.h"
#include "error_listener.h"
#include "generated/bapelLexer.h"
#include "generated/bapelParser.h"
#include <fstream>
#include <iostream>
#include <sstream>
#include <string>
#include <vector>

static void replace_all(std::string& str, const std::string& from, const std::string& to) {
  size_t start_pos = 0;
  while ((start_pos = str.find(from, start_pos)) != std::string::npos) {
    str.replace(start_pos, from.length(), to);
    start_pos += to.length();
  }
}

int main(int argc, char* argv[]) {
  std::string symbol = "SourceFile";
  std::string format = "json";
  std::string filename;
  std::string input_path;

  for (int i = 1; i < argc; ++i) {
    std::string arg = argv[i];
    if (arg.rfind("--symbol=", 0) == 0) {
      symbol = arg.substr(9);
    } else if (arg.rfind("-symbol=", 0) == 0) {
      symbol = arg.substr(8);
    } else if (arg == "--symbol" || arg == "-symbol") {
      if (i + 1 < argc) symbol = argv[++i];
    } else if (arg.rfind("--format=", 0) == 0) {
      format = arg.substr(9);
    } else if (arg.rfind("-format=", 0) == 0) {
      format = arg.substr(8);
    } else if (arg == "--format" || arg == "-format") {
      if (i + 1 < argc) format = argv[++i];
    } else if (arg.rfind("--filename=", 0) == 0) {
      filename = arg.substr(11);
    } else if (arg.rfind("-filename=", 0) == 0) {
      filename = arg.substr(10);
    } else if (arg == "--filename" || arg == "-filename") {
      if (i + 1 < argc) filename = argv[++i];
    } else if (arg == "--workspace" || arg == "-workspace") {
      symbol = "Workspace";
    } else if (arg.rfind("-", 0) != 0) {
      if (input_path.empty()) {
        input_path = arg;
      } else {
        std::cerr << "Expected at most one argument\n";
        return 1;
      }
    }
  }

  std::string input_code;
  if (!input_path.empty()) {
    std::ifstream file(input_path);
    if (!file.is_open()) {
      std::cerr << "Failed to open input file: " << input_path << "\n";
      return 1;
    }
    std::stringstream buffer;
    buffer << file.rdbuf();
    input_code = buffer.str();
    if (filename.empty()) {
      filename = input_path;
    }
  } else {
    std::stringstream buffer;
    buffer << std::cin.rdbuf();
    input_code = buffer.str();
    if (filename.empty()) {
      filename = "<stdin>";
    }
  }

  antlr4::ANTLRInputStream stream(input_code);
  bapelLexer lexer(&stream);
  ast::BapelErrorListener error_listener(filename);
  lexer.removeErrorListeners();
  lexer.addErrorListener(&error_listener);

  antlr4::CommonTokenStream tokens(&lexer);
  bapelParser parser(&tokens);
  parser.removeErrorListeners();
  parser.addErrorListener(&error_listener);

  antlr4::tree::ParseTree* tree = nullptr;
  if (symbol == "SourceFile") {
    tree = parser.sourceFile();
  } else if (symbol == "Workspace") {
    tree = parser.workspace();
  } else if (symbol == "Decl") {
    tree = parser.decl();
  } else {
    std::cerr << "Unsupported symbol \"" << symbol << "\"\n";
    return 1;
  }

  if (error_listener.has_errors()) {
    std::cerr << error_listener.errors()[0] << "\n";
    return 1;
  }

  ast::AstBuilder builder(filename);

  if (symbol == "SourceFile") {
    ast::SourceFile sf = tree->accept(&builder).as<ast::SourceFile>();
    sf.header.filename = ir::new_filename(filename, ir::Pos{});
    std::vector<std::string> val_errors;
    if (!ast::validate_source_file(sf, val_errors)) {
      for (const auto& err : val_errors) {
        std::cerr << err << "\n";
      }
      return 1;
    }

    if (format == "flat") {
      for (const auto& imp : sf.imports.ids) {
        std::cout << "IMPORT " << imp.name << "\n";
      }
      for (const auto& impl : sf.impls.filenames) {
        std::cout << "IMPL " << impl.value << "\n";
      }
      for (const auto& flag : sf.flags.filenames) {
        std::cout << "FLAG " << flag.value << "\n";
      }
      for (const auto& source : sf.body) {
        if (source.case_val == ast::SourceCase::DeclSource && source.decl_data) {
          std::string s = source.decl_data->to_string();
          replace_all(s, "\n", "\\n");
          std::cout << "DECL " << s << "\n";
        } else if (source.case_val == ast::SourceCase::FunctionSource && source.function_data) {
          std::string s = source.function_data->decl().to_string();
          replace_all(s, "\n", "\\n");
          std::cout << "DECL " << s << "\n";
          std::string fn_str = source.function_data->to_string(false);
          replace_all(fn_str, "\n", "\\n");
          std::cout << "FUNC " << fn_str << "\n";
        } else if (source.case_val == ast::SourceCase::TraitSource && source.trait_data) {
          std::string s = source.trait_data->decl().to_string();
          replace_all(s, "\n", "\\n");
          std::cout << "DECL " << s << "\n";
        } else if (source.case_val == ast::SourceCase::ImplSource && source.impl_data) {
          std::string s = source.impl_data->to_string(false);
          replace_all(s, "\n", "\\n");
          std::cout << "TRAIT_IMPL " << s << "\n";
        }
      }
      return 0;
    }

    std::cout << sf.to_json() << "\n";
    return 0;
  } else if (symbol == "Workspace") {
    ast::Workspace ws = tree->accept(&builder).as<ast::Workspace>();
    std::vector<std::string> val_errors;
    if (!ast::validate_workspace(ws, val_errors)) {
      for (const auto& err : val_errors) {
        std::cerr << err << "\n";
      }
      return 1;
    }

    if (format == "flat") {
      for (const auto& pkg : ws.packages.packages) {
        if (pkg.case_val == ast::PackageCase::PrefixPackage) {
          std::cout << "PREFIX " << pkg.module_id.name << " " << pkg.filename.value << "\n";
        } else if (pkg.case_val == ast::PackageCase::ModulePackage) {
          std::cout << "MODULE " << pkg.module_id.name << " " << pkg.filename.value << "\n";
        }
      }
      return 0;
    }

    std::cout << ws.to_json() << "\n";
    return 0;
  } else if (symbol == "Decl") {
    ir::IrDecl decl = tree->accept(&builder).as<ir::IrDecl>();
    if (format == "flat") {
      std::string s = decl.to_string();
      replace_all(s, "\n", "\\n");
      std::cout << "DECL " << s << "\n";
      return 0;
    }
    std::cout << decl.to_json() << "\n";
    return 0;
  }

  return 0;
}
