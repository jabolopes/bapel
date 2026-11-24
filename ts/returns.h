#pragma once

#include "bin/ir_function.h"
#include "bin/ir_term.h"
#include "bin/ir_type.h"

#include <stdexcept>
#include <string>
#include <vector>

namespace ts {

inline void last_terms_impl(ir::IrTerm* term, std::vector<ir::IrTerm*>& last) {
  if (!term) return;
  switch (term->tag()) {
    case ir::IrTermData::Tag::app_term:
    case ir::IrTermData::Tag::app_type:
    case ir::IrTermData::Tag::assign:
    case ir::IrTermData::Tag::const_t:
    case ir::IrTermData::Tag::injection:
    case ir::IrTermData::Tag::lambda:
    case ir::IrTermData::Tag::let_t:
    case ir::IrTermData::Tag::match_t:
    case ir::IrTermData::Tag::projection:
    case ir::IrTermData::Tag::return_t:
    case ir::IrTermData::Tag::set_t:
    case ir::IrTermData::Tag::struct_t:
    case ir::IrTermData::Tag::tuple_t:
    case ir::IrTermData::Tag::type_abs:
    case ir::IrTermData::Tag::var_t:
      last.push_back(term);
      break;

    case ir::IrTermData::Tag::block:
      if (!term->data.as_block().terms.empty()) {
        last_terms_impl(&term->data.as_block().terms.back(), last);
      }
      break;
  }
}

inline std::vector<ir::IrTerm*> last_terms(ir::IrTerm* term) {
  std::vector<ir::IrTerm*> last;
  last_terms_impl(term, last);
  return last;
}

} // namespace ts
