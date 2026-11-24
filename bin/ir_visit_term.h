#pragma once

#include "bin/ir_term.h"

namespace ir {

// Applies `fn(IrTerm&)` to every immediate child IrTerm of `term`.
template <typename F>
inline void visit_term(IrTerm& term, F&& fn) {
  switch (term.tag()) {
    case IrTermData::Tag::app_term: {
      auto& c = term.data.as_app_term();
      if (c.fun) fn(*c.fun);
      if (c.arg) fn(*c.arg);
      break;
    }
    case IrTermData::Tag::app_type:
      if (term.data.as_app_type().fun) {
        fn(*term.data.as_app_type().fun);
      }
      break;
    case IrTermData::Tag::assign: {
      auto& c = term.data.as_assign();
      if (c.ret) fn(*c.ret);
      if (c.arg) fn(*c.arg);
      break;
    }
    case IrTermData::Tag::block:
      for (auto& t : term.data.as_block().terms) {
        fn(t);
      }
      break;
    case IrTermData::Tag::const_t:
      break;
    case IrTermData::Tag::injection:
      if (term.data.as_injection().value) {
        fn(*term.data.as_injection().value);
      }
      break;
    case IrTermData::Tag::lambda:
      if (term.data.as_lambda().body) {
        fn(*term.data.as_lambda().body);
      }
      break;
    case IrTermData::Tag::let_t:
      if (term.data.as_let_t().value) {
        fn(*term.data.as_let_t().value);
      }
      break;
    case IrTermData::Tag::match_t: {
      auto& c = term.data.as_match_t();
      if (c.term) fn(*c.term);
      for (auto& arm : c.arms) {
        if (arm.body) fn(*arm.body);
      }
      break;
    }
    case IrTermData::Tag::projection:
      if (term.data.as_projection().term) {
        fn(*term.data.as_projection().term);
      }
      break;
    case IrTermData::Tag::return_t:
      if (term.data.as_return_t().expr) {
        fn(*term.data.as_return_t().expr);
      }
      break;
    case IrTermData::Tag::set_t: {
      auto& c = term.data.as_set_t();
      if (c.term) fn(*c.term);
      for (auto& lv : c.values) {
        if (lv.value) fn(*lv.value);
      }
      break;
    }
    case IrTermData::Tag::struct_t:
      for (auto& lv : term.data.as_struct_t().values) {
        if (lv.value) fn(*lv.value);
      }
      break;
    case IrTermData::Tag::tuple_t:
      for (auto& elem : term.data.as_tuple_t().elems) {
        fn(elem);
      }
      break;
    case IrTermData::Tag::type_abs:
      if (term.data.as_type_abs().body) {
        fn(*term.data.as_type_abs().body);
      }
      break;
    case IrTermData::Tag::var_t:
      break;
  }
}

// Recursively applies `fn(const IrType&) -> IrType` to every IrType embedded in `term`.
template <typename F>
inline void rewrite_types_in_term(IrTerm& term, F&& fn) {
  if (term.type.has_value()) {
    term.type = fn(*term.type);
  }
  switch (term.tag()) {
    case IrTermData::Tag::app_type:
      term.data.as_app_type().arg = fn(term.data.as_app_type().arg);
      break;
    case IrTermData::Tag::injection:
      term.data.as_injection().variant_type = fn(term.data.as_injection().variant_type);
      break;
    case IrTermData::Tag::lambda:
      term.data.as_lambda().arg.typ = fn(term.data.as_lambda().arg.typ);
      break;
    case IrTermData::Tag::let_t:
      if (term.data.as_let_t().var_type.has_value()) {
        term.data.as_let_t().var_type = fn(*term.data.as_let_t().var_type);
      }
      break;
    case IrTermData::Tag::projection:
      if (term.data.as_projection().reduced_type.has_value()) {
        term.data.as_projection().reduced_type = fn(*term.data.as_projection().reduced_type);
      }
      break;
    case IrTermData::Tag::set_t:
      if (term.data.as_set_t().reduced_type.has_value()) {
        term.data.as_set_t().reduced_type = fn(*term.data.as_set_t().reduced_type);
      }
      break;
    case IrTermData::Tag::type_abs:
      for (auto& b : term.data.as_type_abs().type_param.bounds) {
        b = fn(b);
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
