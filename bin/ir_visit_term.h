#pragma once

#include "bin/ir_term.h"

namespace ir {

// Applies `fn(IrTerm&)` to every immediate child IrTerm of `term`.
template <typename F>
inline void visit_term(IrTerm& term, F&& fn) {
  switch (term.case_val) {
    case IrTermCase::AppTermTerm:
      if (term.app_term) {
        if (term.app_term->fun) fn(*term.app_term->fun);
        if (term.app_term->arg) fn(*term.app_term->arg);
      }
      break;
    case IrTermCase::AppTypeTerm:
      if (term.app_type && term.app_type->fun) {
        fn(*term.app_type->fun);
      }
      break;
    case IrTermCase::AssignTerm:
      if (term.assign) {
        if (term.assign->ret) fn(*term.assign->ret);
        if (term.assign->arg) fn(*term.assign->arg);
      }
      break;
    case IrTermCase::BlockTerm:
      if (term.block) {
        for (auto& t : term.block->terms) {
          fn(t);
        }
      }
      break;
    case IrTermCase::ConstTerm:
      break;
    case IrTermCase::InjectionTerm:
      if (term.injection && term.injection->value) {
        fn(*term.injection->value);
      }
      break;
    case IrTermCase::LambdaTerm:
      if (term.lambda && term.lambda->body) {
        fn(*term.lambda->body);
      }
      break;
    case IrTermCase::LetTerm:
      if (term.let_data && term.let_data->value) {
        fn(*term.let_data->value);
      }
      break;
    case IrTermCase::MatchTerm:
      if (term.match_data) {
        if (term.match_data->term) fn(*term.match_data->term);
        for (auto& arm : term.match_data->arms) {
          if (arm.body) fn(*arm.body);
        }
      }
      break;
    case IrTermCase::ProjectionTerm:
      if (term.projection && term.projection->term) {
        fn(*term.projection->term);
      }
      break;
    case IrTermCase::ReturnTerm:
      if (term.return_data && term.return_data->expr) {
        fn(*term.return_data->expr);
      }
      break;
    case IrTermCase::SetTerm:
      if (term.set_data) {
        if (term.set_data->term) fn(*term.set_data->term);
        for (auto& lv : term.set_data->values) {
          if (lv.value) fn(*lv.value);
        }
      }
      break;
    case IrTermCase::StructTerm:
      if (term.struct_data) {
        for (auto& lv : term.struct_data->values) {
          if (lv.value) fn(*lv.value);
        }
      }
      break;
    case IrTermCase::TupleTerm:
      if (term.tuple_data) {
        for (auto& elem : term.tuple_data->elems) {
          fn(elem);
        }
      }
      break;
    case IrTermCase::TypeAbsTerm:
      if (term.type_abs && term.type_abs->body) {
        fn(*term.type_abs->body);
      }
      break;
    case IrTermCase::VarTerm:
      break;
  }
}

// Recursively applies `fn(const IrType&) -> IrType` to every IrType embedded in `term`.
template <typename F>
inline void rewrite_types_in_term(IrTerm& term, F&& fn) {
  if (term.type.has_value()) {
    term.type = fn(*term.type);
  }
  switch (term.case_val) {
    case IrTermCase::AppTypeTerm:
      if (term.app_type) {
        term.app_type->arg = fn(term.app_type->arg);
      }
      break;
    case IrTermCase::InjectionTerm:
      if (term.injection) {
        term.injection->variant_type = fn(term.injection->variant_type);
      }
      break;
    case IrTermCase::LambdaTerm:
      if (term.lambda) {
        term.lambda->arg.type = fn(term.lambda->arg.type);
      }
      break;
    case IrTermCase::LetTerm:
      if (term.let_data && term.let_data->var_type.has_value()) {
        term.let_data->var_type = fn(*term.let_data->var_type);
      }
      break;
    case IrTermCase::ProjectionTerm:
      if (term.projection && term.projection->reduced_type.has_value()) {
        term.projection->reduced_type = fn(*term.projection->reduced_type);
      }
      break;
    case IrTermCase::SetTerm:
      if (term.set_data && term.set_data->reduced_type.has_value()) {
        term.set_data->reduced_type = fn(*term.set_data->reduced_type);
      }
      break;
    case IrTermCase::TypeAbsTerm:
      if (term.type_abs) {
        for (auto& b : term.type_abs->type_param.bounds) {
          b = fn(b);
        }
      }
      break;
    default:
      break;
  }
  visit_term(term, [&](IrTerm& child) {
    rewrite_types_in_term(child, fn);
  });
}

} // namespace ir
