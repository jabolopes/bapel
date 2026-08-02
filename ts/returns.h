#pragma once

#include "bin/ir_function.h"
#include "bin/ir_term.h"
#include "bin/ir_type.h"

#include <stdexcept>
#include <string>
#include <vector>

namespace ts {

inline void all_blocks_impl(const ir::IrTerm& term, std::vector<ir::IrTerm>& blocks) {
  switch (term.case_val) {
    case ir::IrTermCase::AppTermTerm:
    case ir::IrTermCase::AppTypeTerm:
    case ir::IrTermCase::AssignTerm:
    case ir::IrTermCase::ConstTerm:
    case ir::IrTermCase::InjectionTerm:
    case ir::IrTermCase::LambdaTerm:
    case ir::IrTermCase::LetTerm:
    case ir::IrTermCase::MatchTerm:
    case ir::IrTermCase::ProjectionTerm:
    case ir::IrTermCase::ReturnTerm:
    case ir::IrTermCase::SetTerm:
    case ir::IrTermCase::StructTerm:
    case ir::IrTermCase::TupleTerm:
    case ir::IrTermCase::TypeAbsTerm:
    case ir::IrTermCase::VarTerm:
      break;

    case ir::IrTermCase::BlockTerm:
      if (term.block) {
        blocks.push_back(term);
        for (const auto& t : term.block->terms) {
          all_blocks_impl(t, blocks);
        }
      }
      break;
  }
}

inline std::vector<ir::IrTerm> all_returns(const ir::IrTerm& term) {
  std::vector<ir::IrTerm> blocks;
  all_blocks_impl(term, blocks);

  std::vector<ir::IrTerm> returns;
  for (const auto& block : blocks) {
    if (block.is(ir::IrTermCase::BlockTerm) && block.block) {
      for (const auto& t : block.block->terms) {
        if (t.is(ir::IrTermCase::ReturnTerm)) {
          returns.push_back(t);
        }
      }
    }
  }
  return returns;
}

inline void last_terms_impl(ir::IrTerm* term, std::vector<ir::IrTerm*>& last) {
  if (!term) return;
  switch (term->case_val) {
    case ir::IrTermCase::AppTermTerm:
    case ir::IrTermCase::AppTypeTerm:
    case ir::IrTermCase::AssignTerm:
    case ir::IrTermCase::ConstTerm:
    case ir::IrTermCase::InjectionTerm:
    case ir::IrTermCase::LambdaTerm:
    case ir::IrTermCase::LetTerm:
    case ir::IrTermCase::MatchTerm:
    case ir::IrTermCase::ProjectionTerm:
    case ir::IrTermCase::ReturnTerm:
    case ir::IrTermCase::SetTerm:
    case ir::IrTermCase::StructTerm:
    case ir::IrTermCase::TupleTerm:
    case ir::IrTermCase::TypeAbsTerm:
    case ir::IrTermCase::VarTerm:
      last.push_back(term);
      break;

    case ir::IrTermCase::BlockTerm:
      if (term->block && !term->block->terms.empty()) {
        last_terms_impl(&term->block->terms.back(), last);
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
