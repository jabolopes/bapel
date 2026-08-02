#pragma once

#include "ast/ast_desugar.h"
#include "bin/ir_base.h"
#include "bin/ir_decl.h"
#include "bin/ir_function.h"
#include "bin/ir_parser.h"
#include "bin/ir_unit.h"
#include "comp/desugar.h"
#include "comp/querier.h"
#include "comp/resolver.h"
#include "cpp_parser/parser.h"

#include <fstream>
#include <string>

namespace comp {

inline bool normalize_source_file(Querier querier, const std::string& input_filename, ir::IrUnit& out_unit, std::string& err) {
  // If input is .json, parse directly
  if (input_filename.size() >= 5 && input_filename.substr(input_filename.size() - 5) == ".json") {
    std::ifstream ifs(input_filename);
    if (!ifs) {
      err = "failed to open file: " + input_filename;
      return false;
    }
    std::string content((std::istreambuf_iterator<char>(ifs)), std::istreambuf_iterator<char>());
    out_unit = ir::parse_ir_unit_from_json(content);
    return true;
  }

  // Parse .bpl / .in file directly via parser
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
