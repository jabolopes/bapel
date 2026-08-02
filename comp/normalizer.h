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

inline bool normalize_source_file(Querier querier, const std::string& input_filename, ir::IrUnit& out_unit, std::string& err) {
  // Parse source file directly via parser
  auto res = parser::parse_source_file_from_file(input_filename);
  if (!res.ok) {
    err = res.error().empty() ? ("failed to parse source file: " + input_filename) : res.error();
    return false;
  }
  out_unit = ast::desugar_source_file(res.value);

  if (!resolve_source_file(std::move(querier), out_unit, err)) {
    return false;
  }

  return true;
}

} // namespace comp
