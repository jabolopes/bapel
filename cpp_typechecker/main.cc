#include "bin/ir_base.h"
#include "bin/ir_decl.h"
#include "bin/ir_function.h"
#include "bin/ir_parser.h"
#include "bin/ir_unit.h"
#include "comp/module_finder.h"
#include "comp/querier.h"
#include "comp/resolver.h"
#include "comp/typecheck_unit.h"

#include <iostream>
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
  std::string format = "flat";
  std::string input_file;

  for (int i = 1; i < argc; ++i) {
    std::string arg = argv[i];
    if (arg.rfind("-format=", 0) == 0) {
      format = arg.substr(8);
    } else if (arg.rfind("--format=", 0) == 0) {
      format = arg.substr(9);
    } else if (arg == "-format" || arg == "--format") {
      if (i + 1 < argc) {
        format = argv[++i];
      }
    } else if (arg.rfind("-", 0) != 0) {
      if (input_file.empty()) {
        input_file = arg;
      } else {
        std::cerr << "Usage: typechecker [-format=flat|json|ir] <input_file>\n";
        return 1;
      }
    }
  }

  if (input_file.empty()) {
    std::cerr << "Usage: typechecker [-format=flat|json|ir] <input_file>\n";
    return 1;
  }

  comp::ModuleFinder finder;
  comp::Querier querier(finder);
  comp::TypecheckOptions options;

  ir::IrUnit unit;
  std::string err;
  if (!comp::typecheck_source_file(querier, options, input_file, unit, err)) {
    std::cerr << "Failed to typecheck \"" << input_file << "\": " << err << "\n";
    return 1;
  }

  if (format == "json") {
    std::cout << unit.to_json() << "\n";
  } else if (format == "ir") {
    std::cout << unit.to_string() << "\n";
  } else if (format == "flat") {
    std::cout << "MODULE " << unit.module_id.to_string() << "\n";
    if (unit.case_val == ir::IrUnitCase::BaseUnit) {
      std::cout << "CASE base\n";
    } else {
      std::cout << "CASE impl\n";
    }

    for (const auto& imp : unit.imports) {
      std::cout << "IMPORT " << imp.module_id.to_string() << "\n";
    }
    for (const auto& impl : unit.impls) {
      std::cout << "IMPL " << impl.relative_filename.value << "\n";
    }
    for (const auto& decl : unit.decls) {
      std::string s = decl.to_string();
      std::string escaped_s = s;
      replace_all(escaped_s, "\n", "\\n");
      std::cout << "DECL " << escaped_s << "\n";
      std::string export_str = decl.export_flag ? "1" : "0";
      std::cout << "DECL_DEF " << export_str << " " << decl.id() << " " << escaped_s << "\n";
    }
    for (const auto& trait_impl : unit.trait_impls) {
      std::string s = trait_impl.to_string();
      std::string escaped_s = s;
      replace_all(escaped_s, "\n", "\\n");
      std::cout << "TRAIT_IMPL " << escaped_s << "\n";
      std::string trait_type_str = trait_impl.trait_type.to_string();
      replace_all(trait_type_str, "\n", "\\n");
      std::string type_name_str = trait_impl.type_name.to_string();
      replace_all(type_name_str, "\n", "\\n");
      std::cout << "TRAIT_DEF " << trait_type_str << " " << type_name_str << " " << escaped_s << "\n";
    }
    for (const auto& fn : unit.functions) {
      std::string s = fn.to_string();
      std::string escaped_s = s;
      replace_all(escaped_s, "\n", "\\n");
      std::cout << "FUNC " << escaped_s << "\n";
      std::string export_str = fn.export_flag ? "1" : "0";
      std::string body_str = fn.body.to_string();
      replace_all(body_str, "\n", "\\n");
      std::string ret_type_str = fn.ret_type.to_string();
      replace_all(ret_type_str, "\n", "\\n");
      std::cout << "FUNC_DEF " << export_str << " " << fn.id << " " << ret_type_str << " " << body_str << "\n";
    }
  } else {
    std::cerr << "Unknown format \"" << format << "\" (expected flat, json, or ir)\n";
    return 1;
  }

  return 0;
}
