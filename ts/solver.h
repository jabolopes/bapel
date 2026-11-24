#pragma once

#include "bin/ir_base.h"
#include "bin/ir_term.h"
#include "bin/ir_type.h"
#include "ts/bind.h"
#include "ts/context.h"
#include "ts/predicative.h"
#include "ts/unify.h"

#include <stdexcept>
#include <string>
#include <vector>

namespace ts {

inline bool solve_term(const Context& context, UnificationState& state, ir::IrTerm* term, std::string& err);

inline bool solve_app_term(const Context& context, UnificationState& state, ir::IrTerm* term, std::string& err) {
  auto& c = term->data.as_app_term();
  if (!solve_term(context, state, c.fun.get(), err)) return false;
  if (!solve_term(context, state, c.arg.get(), err)) return false;
  return true;
}

inline bool solve_app_type(const Context& context, UnificationState& state, ir::IrTerm* term, std::string& err) {
  auto& c = term->data.as_app_type();
  if (!solve_term(context, state, c.fun.get(), err)) return false;
  c.arg = state.solve_type(c.arg);
  return true;
}

inline bool solve_assign(const Context& context, UnificationState& state, ir::IrTerm* term, std::string& err) {
  auto& c = term->data.as_assign();
  if (!solve_term(context, state, c.ret.get(), err)) return false;
  if (!solve_term(context, state, c.arg.get(), err)) return false;
  return true;
}

inline bool solve_block(const Context& context, UnificationState& state, ir::IrTerm* term, std::string& err) {
  auto& c = term->data.as_block();
  for (auto& t : c.terms) {
    if (!solve_term(context, state, &t, err)) return false;
  }
  return true;
}

inline bool solve_injection(const Context& context, UnificationState& state, ir::IrTerm* term, std::string& err) {
  auto& c = term->data.as_injection();
  c.variant_type = state.solve_type(c.variant_type);
  if (c.variant_type.data.is(ir::IrTypeData::Tag::variant_t)) {
    auto tags = ir::type_tags(c.variant_type);
    for (size_t i = 0; i < tags.size(); ++i) {
      if (tags[i].id == c.tag) {
        c.tag_index = static_cast<int>(i);
        break;
      }
    }
  }
  return solve_term(context, state, c.value.get(), err);
}

inline bool solve_lambda(const Context& context, UnificationState& state, ir::IrTerm* term, std::string& err) {
  auto& c = term->data.as_lambda();
  c.arg.typ = state.solve_type(c.arg.typ);
  return solve_term(context, state, c.body.get(), err);
}

inline bool solve_let(const Context& context, UnificationState& state, ir::IrTerm* term, std::string& err) {
  auto& c = term->data.as_let_t();
  if (c.var_type) {
    c.var_type = state.solve_type(*c.var_type);
  }
  return solve_term(context, state, c.value.get(), err);
}

inline bool solve_match(const Context& context, UnificationState& state, ir::IrTerm* term, std::string& err) {
  auto& c = term->data.as_match_t();
  if (!solve_term(context, state, c.term.get(), err)) return false;

  ir::IrType obj_type;
  if (c.term->type) {
    obj_type = reduce_and_predicate_type(context, *c.term->type);
  }

  for (auto& arm : c.arms) {
    if (obj_type.data.is(ir::IrTypeData::Tag::variant_t)) {
      auto tags = ir::type_tags(obj_type);
      for (size_t i = 0; i < tags.size(); ++i) {
        if (tags[i].id == arm.tag) {
          arm.index = static_cast<int>(i);
          break;
        }
      }
    }
    if (!solve_term(context, state, arm.body.get(), err)) return false;
  }
  return true;
}

inline bool solve_projection(const Context& context, UnificationState& state, ir::IrTerm* term, std::string& err) {
  auto& c = term->data.as_projection();
  if (!solve_term(context, state, c.term.get(), err)) return false;
  if (c.term->type) {
    c.reduced_type = reduce_and_predicate_type(context, *c.term->type);
  }
  return true;
}

inline bool solve_return(const Context& context, UnificationState& state, ir::IrTerm* term, std::string& err) {
  auto& c = term->data.as_return_t();
  return solve_term(context, state, c.expr.get(), err);
}

inline bool solve_set(const Context& context, UnificationState& state, ir::IrTerm* term, std::string& err) {
  auto& c = term->data.as_set_t();
  if (!solve_term(context, state, c.term.get(), err)) return false;
  if (c.term->type) {
    c.reduced_type = reduce_and_predicate_type(context, *c.term->type);
  }
  for (auto& lv : c.values) {
    if (!solve_term(context, state, lv.value.get(), err)) return false;
  }
  return true;
}

inline bool solve_struct(const Context& context, UnificationState& state, ir::IrTerm* term, std::string& err) {
  auto& c = term->data.as_struct_t();
  for (auto& fv : c.values) {
    if (!solve_term(context, state, fv.value.get(), err)) return false;
  }
  return true;
}

inline bool solve_tuple(const Context& context, UnificationState& state, ir::IrTerm* term, std::string& err) {
  auto& c = term->data.as_tuple_t();
  for (auto& elem : c.elems) {
    if (!solve_term(context, state, &elem, err)) return false;
  }
  return true;
}

inline bool solve_type_abs(const Context& context, UnificationState& state, ir::IrTerm* term, std::string& err) {
  auto& c = term->data.as_type_abs();
  return solve_term(context, state, c.body.get(), err);
}

// TODO: Refactor solve_term_impl and its per-case helpers to use ir::visit_term /
// ir::rewrite_types_in_term from bin/ir_visit_term.h.
inline bool solve_term_impl(const Context& context, UnificationState& state, ir::IrTerm* term, std::string& err) {
  switch (term->tag()) {
    case ir::IrTermData::Tag::app_term:
      return solve_app_term(context, state, term, err);
    case ir::IrTermData::Tag::app_type:
      return solve_app_type(context, state, term, err);
    case ir::IrTermData::Tag::assign:
      return solve_assign(context, state, term, err);
    case ir::IrTermData::Tag::block:
      return solve_block(context, state, term, err);
    case ir::IrTermData::Tag::const_t:
      return true;
    case ir::IrTermData::Tag::injection:
      return solve_injection(context, state, term, err);
    case ir::IrTermData::Tag::lambda:
      return solve_lambda(context, state, term, err);
    case ir::IrTermData::Tag::let_t:
      return solve_let(context, state, term, err);
    case ir::IrTermData::Tag::match_t:
      return solve_match(context, state, term, err);
    case ir::IrTermData::Tag::projection:
      return solve_projection(context, state, term, err);
    case ir::IrTermData::Tag::return_t:
      return solve_return(context, state, term, err);
    case ir::IrTermData::Tag::set_t:
      return solve_set(context, state, term, err);
    case ir::IrTermData::Tag::struct_t:
      return solve_struct(context, state, term, err);
    case ir::IrTermData::Tag::tuple_t:
      return solve_tuple(context, state, term, err);
    case ir::IrTermData::Tag::type_abs:
      return solve_type_abs(context, state, term, err);
    case ir::IrTermData::Tag::var_t:
      return true;
  }
  return true;
}

inline bool solve_term(const Context& context, UnificationState& state, ir::IrTerm* term, std::string& err) {
  if (!term) return true;
  if (!solve_term_impl(context, state, term, err)) {
    err += "\n  solving " + term->to_string();
    return false;
  }
  if (term->type) {
    term->type = state.solve_type(*term->type);
  }
  return true;
}

} // namespace ts
