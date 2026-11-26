#pragma once

#include "ir_module.h"

#include <sstream>
#include <string>
#include <vector>

namespace ir {

using ::IrTypedModule;

inline std::string typed_module_to_bpl_string(const IrTypedModule& m, bool with_pos = true) {
  std::stringstream ss;
  ss << "module " << m.module_id.name << "\n";

  if (!m.mod_imports.empty()) {
    ss << "\nimports {\n";
    for (const auto& imp : m.mod_imports) {
      ss << "  " << imp.module_id.name << " \"" << imp.resolved_filename.value << "\"\n";
    }
    ss << "}\n";
  }

  if (!m.mod_impls.empty()) {
    ss << "\nimpls {\n";
    for (const auto& impl : m.mod_impls) {
      ss << "  \"" << impl.relative_filename.value << "\"\n";
    }
    ss << "}\n";
  }

  if (!m.mod_flags.empty()) {
    ss << "\nflags {\n";
    for (const auto& flag : m.mod_flags) {
      ss << "  \"" << flag << "\"\n";
    }
    ss << "}\n";
  }

  if (!m.import_decls.empty()) {
    ss << "\nimportDecls {\n";
    for (const auto& d : m.import_decls) {
      ss << "  " << decl_to_string(d, with_pos) << "\n";
    }
    ss << "}\n";
  }

  if (!m.ffi_decls.empty()) {
    ss << "\nffiDecls {\n";
    for (const auto& d : m.ffi_decls) {
      ss << "  " << decl_to_string(d, with_pos) << "\n";
    }
    ss << "}\n";
  }

  if (!m.decls.empty()) {
    ss << "\ndecls {\n";
    for (const auto& d : m.decls) {
      ss << "  " << decl_to_string(d, with_pos) << "\n";
    }
    ss << "}\n";
  }

  if (!m.trait_impls.empty()) {
    ss << "\ntraitImpls {\n";
    for (const auto& ti : m.trait_impls) {
      ss << "  " << trait_impl_to_string(ti, with_pos) << "\n";
    }
    ss << "}\n";
  }

  if (!m.functions.empty()) {
    ss << "\nfunctions {\n";
    for (const auto& f : m.functions) {
      ss << "  " << function_to_string(f, with_pos) << "\n";
    }
    ss << "}\n";
  }

  return ss.str();
}

} // namespace ir
