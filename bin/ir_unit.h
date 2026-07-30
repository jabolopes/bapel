#pragma once

#include "ir_decl.h"

namespace ir {

struct IrImport {
  ModuleID module_id;
};

struct IrImpl {
  Filename relative_filename;
};

struct IrUnit {
  IrUnitCase case_val = IrUnitCase::BaseUnit;
  ModuleID module_id;
  Filename filename;
  std::vector<IrImport> imports;
  std::vector<IrImpl> impls;
  std::vector<IrDecl> import_decls;
  std::vector<IrDecl> impl_decls;
  std::vector<IrDecl> decls;
  std::vector<IrFunction> functions;
  std::vector<IrTraitImpl> trait_impls;
  std::vector<IrTraitImpl> imported_trait_impls;
};

} // namespace ir
