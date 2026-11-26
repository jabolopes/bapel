#pragma once

#include "comp/desugar.h"
#include "bin/ir_base.h"
#include "bin/ir_decl.h"
#include "bin/ir_function.h"
#include "bin/ir_module.h"
#include "comp/rewrite_recursive_types.h"
#include "comp/querier.h"
#include "comp/resolver.h"
#include "cpp_parser/parser.h"

#include <string>

namespace comp {

inline bool normalize_module(const Querier& querier, const std::string& base_filename, ir::IrModule& out_module, std::string& err) {
  auto res = parser::parse_source_file_from_file(base_filename);
  if (!res.ok) {
    err = res.error().empty() ? ("failed to parse source file: " + base_filename) : res.error();
    return false;
  }
  comp::DesugaredFile base = comp::desugar_source_file(res.value);

  out_module = ir::IrModule{};
  out_module.module_id = std::move(base.module_id);
  out_module.mod_flags = std::move(base.flags);
  out_module.mod_impls = std::move(base.impls);
  base.file.case_val = ir::IrFileCase::mk_file_base();
  out_module.files.push_back(std::move(base.file));

  if (!resolve_module(querier, out_module, err)) {
    return false;
  }

  if (!rewrite_self_recursive_module(out_module, err)) {
    return false;
  }

  return true;
}

} // namespace comp
