#pragma once

#include "bin/ir_base.h"
#include "bin/ir_function.h"
#include "bin/ir_term.h"
#include "bin/ir_type.h"

#include <memory>
#include <string>
#include <vector>

namespace comp {

inline ir::IrTerm desugar_term(ir::IrTerm term);

inline ir::IrTerm desugar_if_term(ir::IrTerm cond, ir::IrTerm then_branch, ir::IrTerm else_branch) {
  ir::MatchArm arm_true;
  arm_true.tag = "True";
  arm_true.arg = "_";
  arm_true.body = std::make_shared<ir::IrTerm>(desugar_term(std::move(then_branch)));

  ir::MatchArm arm_false;
  arm_false.tag = "False";
  arm_false.arg = "_";
  arm_false.body = std::make_shared<ir::IrTerm>(desugar_term(std::move(else_branch)));

  ir::IrTerm match_term;
  match_term.case_val = ir::IrTermCase::MatchTerm;
  match_term.match_data = std::make_shared<ir::MatchTermData>();
  match_term.match_data->term = std::make_shared<ir::IrTerm>(desugar_term(std::move(cond)));
  match_term.match_data->arms = {std::move(arm_true), std::move(arm_false)};
  return match_term;
}

inline ir::IrTerm desugar_binop_term(const std::string& op, ir::IrTerm left, ir::IrTerm right) {
  if (op == "&&") {
    ir::IrTerm false_term = ir::new_injection_term(ir::new_name_type("bool"), "False", ir::new_tuple_term({}));
    return desugar_if_term(std::move(left), std::move(right), std::move(false_term));
  }
  if (op == "||") {
    ir::IrTerm true_term = ir::new_injection_term(ir::new_name_type("bool"), "True", ir::new_tuple_term({}));
    return desugar_if_term(std::move(left), std::move(true_term), std::move(right));
  }

  ir::IrTerm fun = ir::new_var_term(op);
  ir::IrTerm arg = ir::new_tuple_term({desugar_term(std::move(left)), desugar_term(std::move(right))});
  return ir::new_app_term(std::move(fun), std::move(arg));
}

inline ir::IrTerm desugar_unop_term(const std::string& op, ir::IrTerm arg) {
  ir::IrTerm fun = ir::new_var_term(op);
  return ir::new_app_term(std::move(fun), desugar_term(std::move(arg)));
}

inline ir::IrTerm desugar_term(ir::IrTerm term) {
  switch (term.case_val) {
    case ir::IrTermCase::AppTermTerm:
      if (term.app_term) {
        if (term.app_term->fun) *term.app_term->fun = desugar_term(std::move(*term.app_term->fun));
        if (term.app_term->arg) *term.app_term->arg = desugar_term(std::move(*term.app_term->arg));
      }
      break;

    case ir::IrTermCase::AppTypeTerm:
      if (term.app_type && term.app_type->fun) {
        *term.app_type->fun = desugar_term(std::move(*term.app_type->fun));
      }
      break;

    case ir::IrTermCase::AssignTerm:
      if (term.assign) {
        if (term.assign->ret) *term.assign->ret = desugar_term(std::move(*term.assign->ret));
        if (term.assign->arg) *term.assign->arg = desugar_term(std::move(*term.assign->arg));
      }
      break;

    case ir::IrTermCase::BlockTerm:
      if (term.block) {
        for (auto& t : term.block->terms) {
          t = desugar_term(std::move(t));
        }
      }
      break;

    case ir::IrTermCase::ConstTerm:
      break;

    case ir::IrTermCase::InjectionTerm:
      if (term.injection && term.injection->value) {
        *term.injection->value = desugar_term(std::move(*term.injection->value));
      }
      break;

    case ir::IrTermCase::LambdaTerm:
      if (term.lambda && term.lambda->body) {
        *term.lambda->body = desugar_term(std::move(*term.lambda->body));
      }
      break;

    case ir::IrTermCase::LetTerm:
      if (term.let_data && term.let_data->value) {
        *term.let_data->value = desugar_term(std::move(*term.let_data->value));
      }
      break;

    case ir::IrTermCase::MatchTerm:
      if (term.match_data) {
        if (term.match_data->term) *term.match_data->term = desugar_term(std::move(*term.match_data->term));
        for (auto& arm : term.match_data->arms) {
          if (arm.body) *arm.body = desugar_term(std::move(*arm.body));
        }
      }
      break;

    case ir::IrTermCase::ProjectionTerm:
      if (term.projection && term.projection->term) {
        *term.projection->term = desugar_term(std::move(*term.projection->term));
      }
      break;

    case ir::IrTermCase::ReturnTerm:
      if (term.return_data && term.return_data->expr) {
        *term.return_data->expr = desugar_term(std::move(*term.return_data->expr));
      }
      break;

    case ir::IrTermCase::SetTerm:
      if (term.set_data) {
        if (term.set_data->term) *term.set_data->term = desugar_term(std::move(*term.set_data->term));
        for (auto& lv : term.set_data->values) {
          if (lv.value) *lv.value = desugar_term(std::move(*lv.value));
        }
      }
      break;

    case ir::IrTermCase::StructTerm:
      if (term.struct_data) {
        for (auto& lv : term.struct_data->values) {
          if (lv.value) *lv.value = desugar_term(std::move(*lv.value));
        }
      }
      break;

    case ir::IrTermCase::TupleTerm:
      if (term.tuple_data) {
        for (auto& elem : term.tuple_data->elems) {
          elem = desugar_term(std::move(elem));
        }
      }
      break;

    case ir::IrTermCase::TypeAbsTerm:
      if (term.type_abs && term.type_abs->body) {
        *term.type_abs->body = desugar_term(std::move(*term.type_abs->body));
      }
      break;

    case ir::IrTermCase::VarTerm:
      break;
  }
  return term;
}

inline ir::IrFunction desugar_function(ir::IrFunction fn) {
  fn.body = desugar_term(std::move(fn.body));
  return fn;
}

} // namespace comp
