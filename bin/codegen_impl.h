#pragma once

#include "bapel/compiler/codegen/cpp.h"
#include "comp/typecheck_module.h"
#include "comp/session.h"
#include "comp/querier.h"
#include "comp/module_finder.h"
#include <cstdlib>
#include <fstream>
#include <iostream>
#include <sstream>
#include <string>
#include <tuple>
#include <vector>

namespace codegen {

inline void format_file(const std::string& filename) {
  std::string cmd = "clang-format -i " + filename + " 2>/dev/null";
  std::system(cmd.c_str());
}

inline bool print_module_to_file(const ir::IrTypedModule& mod, ::PrinterMode mode, const std::string& filename) {
  std::ofstream f(filename);
  if (!f.is_open()) return false;
  ir::IrTypedModule mod_copy = mod;
  f << ::print_module_to_string(&mod_copy, mode);
  f.close();
  format_file(filename);
  return true;
}

inline bool compile_ir_module_direct(ir::IrTypedModule mod, const std::string& output_base) {
  mod = ::record_anonymous_types_from_module(&mod);

  std::string base = output_base;
  auto dot = base.rfind('.');
  if (dot != std::string::npos) {
    base = base.substr(0, dot);
  }

  if (!print_module_to_file(mod, ::PrinterMode::mk_mode_public_header(), base + ".h")) return false;
  if (!print_module_to_file(mod, ::PrinterMode::mk_mode_private_header(), base + "_private.h")) return false;
  if (!print_module_to_file(mod, ::PrinterMode::mk_mode_source(), base + ".cc")) return false;
  return true;
}

// @bpl: pub codegen::compile_module: (&comp::Session, String, String) -> i64
inline int64_t compile_module(comp::Session* session, const std::string& base_file, const std::string& output_base) {
  comp::TypecheckOptions options;

  ir::IrTypedModule mod;
  std::string err;
  if (!comp::typecheck_module_from_file(session->querier, options, base_file, mod, err)) {
    std::cerr << "Typecheck error in " << base_file << ":\n" << err << "\n";
    return 1;
  }

  return compile_ir_module_direct(std::move(mod), output_base) ? 0 : 1;
}

} // namespace codegen

namespace typechecker {

// @bpl: pub typechecker::run: (&comp::Session, Vector String) -> (i64, String)
inline std::tuple<int64_t, std::string> run(comp::Session* session, const std::vector<std::string>& args) {
  std::string format = "flat";
  std::string input_file;

  for (size_t i = 0; i < args.size(); ++i) {
    const std::string& arg = args[i];
    if (arg.rfind("-format=", 0) == 0) {
      format = arg.substr(8);
    } else if (arg.rfind("--format=", 0) == 0) {
      format = arg.substr(9);
    } else if (arg == "-format" || arg == "--format") {
      if (i + 1 < args.size()) {
        format = args[++i];
      }
    } else if (arg.rfind("-", 0) != 0) {
      if (input_file.empty()) {
        input_file = arg;
      } else {
        return {1, "Usage: typechecker [-format=flat] <input_file>\n"};
      }
    }
  }

  if (input_file.empty()) {
    return {1, "Usage: typechecker [-format=flat] <input_file>\n"};
  }

  comp::TypecheckOptions options;

  ir::IrTypedModule mod;
  std::string err;
  if (!comp::typecheck_module_from_file(session->querier, options, input_file, mod, err)) {
    return {1, "Failed to typecheck \"" + input_file + "\": " + err + "\n"};
  }

  std::stringstream out;
  if (format == "flat") {
    out << "MODULE " << mod.module_id.name << "\n";
    out << "CASE base\n";

    for (const auto& imp : mod.mod_imports) {
      out << "IMPORT " << imp.module_id.name << "\n";
    }
    for (const auto& impl : mod.mod_impls) {
      out << "IMPL " << impl.relative_filename.value << "\n";
    }
    for (const auto& decl : mod.decls) {
      std::string s = ir::decl_to_string(decl);
      std::string escaped_s = s;
      comp::replace_all_str(escaped_s, "\n", "\\n");
      out << "DECL " << escaped_s << "\n";
      std::string export_str = decl.export_flag ? "1" : "0";
      out << "DECL_DEF " << export_str << " " << ir::decl_id(decl) << " " << escaped_s << "\n";
    }
    for (const auto& trait_impl : mod.trait_impls) {
      std::string s = ir::trait_impl_to_string(trait_impl);
      std::string escaped_s = s;
      comp::replace_all_str(escaped_s, "\n", "\\n");
      out << "TRAIT_IMPL " << escaped_s << "\n";
      std::string trait_type_str = ir::type_to_string(trait_impl.trait_type);
      comp::replace_all_str(trait_type_str, "\n", "\\n");
      std::string type_name_str = ir::type_to_string(trait_impl.type_name);
      comp::replace_all_str(type_name_str, "\n", "\\n");
      out << "TRAIT_DEF " << trait_type_str << " " << type_name_str << " " << escaped_s << "\n";
    }
    for (const auto& fn : mod.functions) {
      std::string s = ir::function_to_string(fn);
      std::string escaped_s = s;
      comp::replace_all_str(escaped_s, "\n", "\\n");
      out << "FUNC " << escaped_s << "\n";
      std::string export_str = fn.export_flag ? "1" : "0";
      out << "FUNC_DEF " << export_str << " " << fn.id << " " << escaped_s << "\n";
    }
  }
  return {0, out.str()};
}

} // namespace typechecker
