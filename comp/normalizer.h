#pragma once

#include "ast/ast_desugar.h"
#include "bin/ir_base.h"
#include "bin/ir_decl.h"
#include "bin/ir_function.h"
#include "bin/ir_unit.h"
#include "comp/desugar.h"
#include "comp/querier.h"
#include "comp/resolver.h"
#include "cpp_parser/parser.h"

#include <string>

namespace comp {

inline bool normalize_module(Querier querier, const std::string& base_filename, ir::IrModule& out_module, std::string& err) {
  auto res = parser::parse_source_file_from_file(base_filename);
  if (!res.ok) {
    err = res.error().empty() ? ("failed to parse source file: " + base_filename) : res.error();
    return false;
  }
  ir::IrUnit base_unit = ast::desugar_source_file(res.value);

  out_module = ir::IrModule{};
  out_module.module_id = base_unit.module_id;
  for (const auto& f : res.value.flags.filenames) {
    out_module.flags.push_back(f.value);
  }
  out_module.impls = std::move(base_unit.impls);

  ir::IrFile base_file;
  base_file.case_val = ir::IrFileCase::Base;
  base_file.filename = base_unit.filename;
  base_file.imports = std::move(base_unit.imports);
  base_file.decls = std::move(base_unit.decls);
  base_file.trait_impls = std::move(base_unit.trait_impls);
  base_file.functions = std::move(base_unit.functions);
  out_module.files.push_back(std::move(base_file));

  if (!resolve_module(std::move(querier), out_module, err)) {
    return false;
  }

  if (!desugar_self_recursive_module(out_module, err)) {
    return false;
  }

  return true;
}

} // namespace comp
