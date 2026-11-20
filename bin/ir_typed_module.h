#pragma once

#include "ir_unit.h"

#include <sstream>
#include <string>
#include <vector>

namespace ir {

struct IrTypedModule {
  ModuleID module_id;
  std::vector<IrResolvedImport> imports;
  std::vector<IrImpl> impls;
  std::vector<std::string> flags;
  std::vector<IrDecl> import_decls;
  std::vector<IrDecl> ffi_decls;
  std::vector<IrDecl> decls;
  std::vector<IrTraitImpl> trait_impls;
  std::vector<IrFunction> functions;

  std::string to_bpl_string(bool with_pos = true) const {
    std::stringstream ss;
    ss << "module " << module_id.to_string(false) << "\n";

    if (!imports.empty()) {
      ss << "\nimports {\n";
      for (const auto& imp : imports) {
        ss << "  " << imp.module_id.to_string(false) << " \"" << imp.resolved_filename.value << "\"\n";
      }
      ss << "}\n";
    }

    if (!impls.empty()) {
      ss << "\nimpls {\n";
      for (const auto& impl : impls) {
        ss << "  \"" << impl.relative_filename.value << "\"\n";
      }
      ss << "}\n";
    }

    if (!flags.empty()) {
      ss << "\nflags {\n";
      for (const auto& flag : flags) {
        ss << "  \"" << flag << "\"\n";
      }
      ss << "}\n";
    }

    if (!import_decls.empty()) {
      ss << "\nimportDecls {\n";
      for (const auto& d : import_decls) {
        ss << "  " << d.to_string(with_pos) << "\n";
      }
      ss << "}\n";
    }

    if (!ffi_decls.empty()) {
      ss << "\nffiDecls {\n";
      for (const auto& d : ffi_decls) {
        ss << "  " << d.to_string(with_pos) << "\n";
      }
      ss << "}\n";
    }

    if (!decls.empty()) {
      ss << "\ndecls {\n";
      for (const auto& d : decls) {
        ss << "  " << d.to_string(with_pos) << "\n";
      }
      ss << "}\n";
    }

    if (!trait_impls.empty()) {
      ss << "\ntraitImpls {\n";
      for (const auto& ti : trait_impls) {
        ss << "  " << ti.to_string(with_pos) << "\n";
      }
      ss << "}\n";
    }

    if (!functions.empty()) {
      ss << "\nfunctions {\n";
      for (const auto& f : functions) {
        ss << "  " << f.to_string(with_pos) << "\n";
      }
      ss << "}\n";
    }

    return ss.str();
  }
};

} // namespace ir
