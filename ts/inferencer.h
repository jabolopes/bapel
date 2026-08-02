#pragma once

#include "bin/ir_base.h"
#include "bin/ir_decl.h"
#include "bin/ir_function.h"
#include "bin/ir_term.h"
#include "bin/ir_type.h"
#include "ts/bind.h"
#include "ts/context.h"
#include "ts/list.h"
#include "ts/predicative.h"
#include "ts/reduce_type.h"
#include "ts/returns.h"
#include "ts/solver.h"
#include "ts/unify.h"

#include <memory>
#include <stdexcept>
#include <string>
#include <vector>

namespace ts {

class Inferencer {
 public:
  explicit Inferencer(Context context) : context_(std::move(context)) {}

  ir::IrType new_evar() {
    return state_.new_evar(context_);
  }

  ir::IrType predicate_type(const ir::IrType& typ) {
    return ts::predicate_type(context_, typ);
  }

  ir::IrType reduce_and_predicate_type(const ir::IrType& typ) {
    return ts::reduce_and_predicate_type(context_, typ);
  }

  bool unify(const ir::IrType& left, const ir::IrType& right) {
    return ts::unify(context_, state_, left, right, err_);
  }

  bool solve_term(ir::IrTerm* term) {
    return ts::solve_term(context_, state_, term, err_);
  }

  ir::IrType solve_type(const ir::IrType& typ) {
    return state_.solve_type(typ);
  }

  const Context& context() const { return context_; }
  const std::string& error() const { return err_; }

  bool infer_function(ir::IrFunction* function, Context& out_context) {
    Context orig_context = context_;

    ir::IrDecl decl = function->decl();
    try {
      context_ = context_.add_bind(new_term_def_bind(decl.term->id, decl.term->type));
    } catch (const std::exception& e) {
      err_ = e.what();
      return false;
    }

    out_context = context_;

    try {
      context_ = context_.enter_function(function->type_params, function->args);
    } catch (const std::exception& e) {
      err_ = e.what();
      context_ = orig_context;
      return false;
    }

    expect_returns_ = expect_returns_.add(function->ret_type);

    if (!infer(&function->body, nullptr, &function->ret_type)) {
      context_ = orig_context;
      return false;
    }

    if (!solve_term(&function->body)) {
      context_ = orig_context;
      return false;
    }

    expect_returns_ = expect_returns_.remove();
    return true;
  }

  bool infer(ir::IrTerm* term, ir::IrTerm* parent_term, const ir::IrType* expect_type) {
    if (!term) return true;

    ir::IrType evar = new_evar();
    if (expect_type) {
      unify(evar, *expect_type);
    }

    if (!infer_impl(evar, term, parent_term, expect_type)) {
      err_ += "\n  inferring " + term->to_string();
      return false;
    }

    if (term->type) {
      term->type = solve_type(*term->type);
      term->type = predicate_type(*term->type);
    }

    return true;
  }

 private:
  bool try_resolve_method_call(ir::IrTerm* term, ir::IrTerm* parent_term, const ir::IrType* expect_type, bool& resolved) {
    resolved = false;
    if (!term->is(ir::IrTermCase::AppTermTerm) || !term->app_term) {
      return true;
    }

    auto& c = *term->app_term;
    if (!c.fun || !c.fun->is(ir::IrTermCase::ProjectionTerm) || !c.fun->projection) {
      return true;
    }

    auto& proj = *c.fun->projection;
    if (!infer(proj.term.get(), c.fun.get(), nullptr)) {
      return false;
    }
    if (!proj.term->type) {
      return true;
    }

    ir::IrType obj_type = reduce_and_predicate_type(*proj.term->type);

    int idx = -1;
    ir::StructField sf;
    ir::IrType elem;
    if (obj_type.is(ir::IrTypeCase::StructType) && obj_type.field_by_label(proj.label, idx, sf)) {
      return true;
    }
    if (obj_type.is(ir::IrTypeCase::TupleType) && obj_type.elem_by_label(proj.label, idx, elem)) {
      return true;
    }
    if (obj_type.is(ir::IrTypeCase::VariantType) && obj_type.tag_by_label(proj.label, idx)) {
      return true;
    }

    std::string method_name;
    ir::IrType method_type;
    if (!context_.lookup_method(obj_type, proj.label, method_name, method_type)) {
      return true;
    }

    ir::IrType ft = method_type;
    while (ft.is(ir::IrTypeCase::ForallType) && ft.forall && ft.forall->type) {
      ft = *ft.forall->type;
    }
    if (!ft.is(ir::IrTypeCase::FunType) || !ft.fun) {
      return true;
    }

    ir::IrType expected_receiver;
    bool is_multi_arg = false;
    if (ft.fun->arg && ft.fun->arg->is(ir::IrTypeCase::TupleType) && !ft.fun->arg->elems().empty()) {
      expected_receiver = ft.fun->arg->elems()[0];
      is_multi_arg = true;
    } else if (ft.fun->arg) {
      expected_receiver = *ft.fun->arg;
    }

    bool expects_ptr = expected_receiver.is(ir::IrTypeCase::AppType) && base_type_name(expected_receiver) == "Ptr";
    bool is_ptr = obj_type.is(ir::IrTypeCase::AppType) && base_type_name(obj_type) == "Ptr";

    ir::IrTerm adjusted_s = *proj.term;
    if (expects_ptr && !is_ptr) {
      adjusted_s = ir::new_app_term(ir::new_var_term("Ptr::mk"), *proj.term);
      adjusted_s.pos = proj.term->pos;
    } else if (!expects_ptr && is_ptr) {
      adjusted_s = ir::new_app_term(ir::new_var_term("Ptr::get"), *proj.term);
      adjusted_s.pos = proj.term->pos;
    }

    ir::IrTerm new_arg;
    if (!is_multi_arg) {
      if (c.arg->is(ir::IrTermCase::TupleTerm) && c.arg->tuple_data && !c.arg->tuple_data->elems.empty()) {
        err_ = "method " + method_name + " takes no arguments, but was called with arguments";
        return false;
      }
      new_arg = adjusted_s;
    } else {
      if (c.arg->is(ir::IrTermCase::TupleTerm) && c.arg->tuple_data) {
        std::vector<ir::IrTerm> new_elems = {adjusted_s};
        new_elems.insert(new_elems.end(), c.arg->tuple_data->elems.begin(), c.arg->tuple_data->elems.end());
        new_arg = ir::new_tuple_term(std::move(new_elems));
        new_arg.pos = c.arg->pos;
      } else {
        new_arg = ir::new_tuple_term({adjusted_s, *c.arg});
        new_arg.pos = c.arg->pos;
      }
    }

    ir::IrTerm new_fun = ir::new_var_term(method_name);
    new_fun.pos = c.fun->pos;
    ir::Pos pos = term->pos;
    *term = ir::new_app_term(std::move(new_fun), std::move(new_arg));
    term->pos = pos;
    resolved = true;
    return infer(term, parent_term, expect_type);
  }

  bool infer_app_term(const ir::IrType& evar, ir::IrTerm* term, ir::IrTerm* parent_term, const ir::IrType* expect_type) {
    bool resolved = false;
    if (!try_resolve_method_call(term, parent_term, expect_type, resolved)) {
      return false;
    }
    if (resolved) {
      return true;
    }

    auto& c = *term->app_term;
    if (!infer(c.fun.get(), term, nullptr)) {
      return false;
    }

    if (!c.fun->type) {
      return infer(c.arg.get(), term, nullptr);
    }

    if (c.fun->type->is(ir::IrTypeCase::FunType) && c.fun->type->fun) {
      ir::IrType arg_type = *c.fun->type->fun->arg;
      ir::IrType ret_type = *c.fun->type->fun->ret;

      if (!infer(c.arg.get(), term, &arg_type)) {
        return false;
      }

      unify(evar, ret_type);
      unify(*c.fun->type, ir::new_function_type(arg_type, evar));
      term->type = ret_type;
      return true;
    }

    if (c.fun->type->is(ir::IrTypeCase::ForallType) && !c.fun->is(ir::IrTermCase::AppTermTerm)) {
      if (!infer(c.arg.get(), term, nullptr)) {
        return false;
      }
      ir::IrTerm app_type = ir::new_app_type_term(*c.fun, new_evar());
      c.fun = std::make_shared<ir::IrTerm>(std::move(app_type));
      return infer(term, parent_term, expect_type);
    }

    if (c.fun->type->is(ir::IrTypeCase::ForallType)) {
      if (!infer(c.arg.get(), term, nullptr)) {
        return false;
      }
      return infer(term, parent_term, expect_type);
    }

    return true;
  }

  bool infer_app_type(const ir::IrType& evar, ir::IrTerm* term, ir::IrTerm* parent_term, const ir::IrType* expect_type) {
    auto& c = *term->app_type;
    if (!infer(c.fun.get(), term, nullptr)) {
      return false;
    }
    if (!c.fun->type || !c.fun->type->is(ir::IrTypeCase::ForallType) || !c.fun->type->forall) {
      return true;
    }

    const auto& forall_type = *c.fun->type->forall;
    ir::IrType typ = ir::substitute_type(*forall_type.type, ir::new_var_type(forall_type.type_param.var), c.arg);
    unify(evar, typ);
    term->type = typ;
    return true;
  }

  bool infer_assign(const ir::IrType& evar, ir::IrTerm* term, ir::IrTerm* parent_term, const ir::IrType* expect_type) {
    auto& c = *term->assign;
    if (!infer(c.ret.get(), term, nullptr)) {
      return false;
    }
    const ir::IrType* ret_t = c.ret->type ? &*c.ret->type : nullptr;
    if (!infer(c.arg.get(), term, ret_t)) {
      return false;
    }
    if (c.arg->type) {
      unify(evar, *c.arg->type);
      term->type = *c.arg->type;
    }
    return true;
  }

  bool infer_block(const ir::IrType& evar, ir::IrTerm* term, ir::IrTerm* parent_term, const ir::IrType* expect_type) {
    auto& c = *term->block;
    Context orig_context = context_;

    try {
      context_ = context_.enter_scope();
    } catch (const std::exception& e) {
      err_ = e.what();
      return false;
    }

    for (size_t i = 0; i < c.terms.size(); ++i) {
      const ir::IrType* item_expect = (i + 1 == c.terms.size()) ? expect_type : nullptr;
      if (!infer(&c.terms[i], term, item_expect)) {
        context_ = orig_context;
        return false;
      }
    }

    if (!c.terms.empty() && c.terms.back().type) {
      unify(evar, *c.terms.back().type);
      term->type = *c.terms.back().type;
    }

    if (!solve_term(term)) {
      context_ = orig_context;
      return false;
    }

    context_ = orig_context;
    return true;
  }

  bool infer_const(const ir::IrType& evar, ir::IrTerm* term, ir::IrTerm* parent_term, const ir::IrType* expect_type) {
    auto& c = *term->const_data;
    if (c.literal.is_rune()) {
      ir::IrType typ = ir::new_name_type("i8");
      unify(evar, typ);
      term->type = typ;
      return true;
    }
    if (c.literal.is_str()) {
      ir::IrType typ = ir::new_name_type("StringView");
      unify(evar, typ);
      term->type = typ;
      return true;
    }

    if ((!parent_term || !parent_term->is(ir::IrTermCase::AppTypeTerm)) && expect_type) {
      *term = ir::new_app_type_term(*term, *expect_type);
      return infer(term, parent_term, expect_type);
    }

    ir::IrType typ = ir::new_forall_type(ir::TypeParam{"a", ir::new_type_kind(), {}}, ir::new_var_type("a"));
    unify(evar, typ);
    term->type = typ;
    return true;
  }

  bool infer_injection(const ir::IrType& evar, ir::IrTerm* term, const ir::IrType* expect_type) {
    auto& c = *term->injection;
    ir::IrType variant_type = reduce_and_predicate_type(c.variant_type);
    if (!variant_type.is(ir::IrTypeCase::VariantType)) {
      err_ = "expected variant type; got " + variant_type.to_string();
      return false;
    }

    auto tags = variant_type.tags();
    ir::IrType tag_type;
    bool found = false;
    for (size_t i = 0; i < tags.size(); ++i) {
      if (tags[i].id == c.tag) {
        c.tag_index = static_cast<int>(i);
        tag_type = tags[i].type ? *tags[i].type : ir::IrType{};
        found = true;
        break;
      }
    }
    if (!found) {
      err_ = "tag " + c.tag + " is not a valid tag in variant type " + variant_type.to_string();
      return false;
    }

    if (!infer(c.value.get(), term, &tag_type)) {
      return false;
    }

    unify(evar, c.variant_type);
    term->type = c.variant_type;
    return true;
  }

  bool infer_lambda(const ir::IrType& evar, ir::IrTerm* term, const ir::IrType* expect_type) {
    auto& c = *term->lambda;
    Context orig_context = context_;

    try {
      context_ = context_.enter_function({}, {c.arg});
    } catch (const std::exception& e) {
      err_ = e.what();
      return false;
    }

    if (!expect_type || !expect_type->is(ir::IrTypeCase::FunType) || !expect_type->fun) {
      ir::IrType ret_evar = new_evar();
      expect_returns_ = expect_returns_.add(ret_evar);

      if (!infer(c.body.get(), term, &ret_evar)) {
        context_ = orig_context;
        return false;
      }

      if (c.body->type) {
        unify(ret_evar, *c.body->type);
      }

      expect_returns_ = expect_returns_.remove();
      ir::IrType typ = ir::new_function_type(c.arg.type, ret_evar);
      unify(evar, typ);
      term->type = typ;
      context_ = orig_context;
      return true;
    }

    ir::IrType expect_body_type = *expect_type->fun->ret;
    expect_returns_ = expect_returns_.add(expect_body_type);

    if (!infer(c.body.get(), term, &expect_body_type)) {
      context_ = orig_context;
      return false;
    }

    expect_returns_ = expect_returns_.remove();
    ir::IrType typ = ir::new_function_type(c.arg.type, expect_body_type);
    unify(evar, typ);
    term->type = typ;
    context_ = orig_context;
    return true;
  }

  bool infer_let(const ir::IrType& evar, ir::IrTerm* term, ir::IrTerm* parent_term, const ir::IrType* expect_type) {
    auto& c = *term->let_data;
    if (!c.var_type) {
      if (!infer(c.value.get(), term, nullptr)) {
        return false;
      }
      if (!c.value->type) {
        err_ = "failed to infer type for let variable " + c.var + "; please add a type annotation";
        return false;
      }

      ir::IrType var_type = predicate_type(*c.value->type);
      try {
        context_ = context_.add_bind(new_term_def_bind(c.var, var_type));
      } catch (const std::exception& e) {
        err_ = e.what();
        return false;
      }

      c.var_type = var_type;
      unify(evar, var_type);
      term->type = var_type;
      return true;
    }

    ir::IrType var_type = *c.var_type;
    try {
      context_ = context_.add_bind(new_term_def_bind(c.var, var_type));
    } catch (const std::exception& e) {
      err_ = e.what();
      return false;
    }

    if (!infer(c.value.get(), term, &var_type)) {
      return false;
    }

    unify(evar, var_type);
    term->type = var_type;
    return true;
  }

  bool infer_match(const ir::IrType& evar, ir::IrTerm* term, const ir::IrType* expect_type) {
    auto& c = *term->match_data;
    if (!infer(c.term.get(), term, nullptr)) {
      return false;
    }

    ir::IrType obj_type;
    if (c.term->type) {
      obj_type = reduce_and_predicate_type(*c.term->type);
    }

    if (!obj_type.is(ir::IrTypeCase::VariantType)) {
      return true;
    }

    auto tags = obj_type.tags();
    std::optional<ir::IrType> match_type;

    for (auto& arm : c.arms) {
      bool found = false;
      ir::IrType tag_type;
      for (size_t i = 0; i < tags.size(); ++i) {
        if (tags[i].id == arm.tag) {
          arm.index = static_cast<int>(i);
          tag_type = tags[i].type ? *tags[i].type : ir::IrType{};
          found = true;
          break;
        }
      }
      if (!found) {
        err_ = "tag " + arm.tag + " is not a valid tag of variant type " + obj_type.to_string();
        return false;
      }

      Context orig_context = context_;
      try {
        context_ = context_.add_bind(new_term_def_bind(arm.arg, tag_type));
      } catch (const std::exception& e) {
        err_ = e.what();
        return false;
      }

      const ir::IrType* expect_arm = match_type ? &*match_type : nullptr;
      if (!infer(arm.body.get(), term, expect_arm)) {
        context_ = orig_context;
        return false;
      }

      if (!match_type && arm.body->type) {
        match_type = *arm.body->type;
      }

      context_ = orig_context;
    }

    if (match_type) {
      unify(evar, *match_type);
      term->type = *match_type;
    }

    return true;
  }

  bool infer_projection(const ir::IrType& evar, ir::IrTerm* term, ir::IrTerm* parent_term, const ir::IrType* expect_type) {
    auto& c = *term->projection;
    if (!infer(c.term.get(), term, nullptr)) {
      return false;
    }

    if (!c.term->type) return true;

    ir::IrType obj_type = reduce_and_predicate_type(*c.term->type);
    c.reduced_type = obj_type;

    int idx = -1;
    ir::StructField sf;
    ir::IrType elem;
    if (obj_type.is(ir::IrTypeCase::StructType) && obj_type.field_by_label(c.label, idx, sf)) {
      ir::IrType ft = sf.type ? *sf.type : ir::IrType{};
      unify(evar, ft);
      term->type = ft;
      return true;
    }
    if (obj_type.is(ir::IrTypeCase::TupleType) && obj_type.elem_by_label(c.label, idx, elem)) {
      unify(evar, elem);
      term->type = elem;
      return true;
    }
    if (obj_type.is(ir::IrTypeCase::VariantType) && obj_type.tag_by_label(c.label, idx)) {
      auto tags = obj_type.tags();
      ir::IrType tt = tags[idx].type ? *tags[idx].type : ir::IrType{};
      unify(evar, tt);
      term->type = tt;
      return true;
    }

    std::string method_name;
    ir::IrType method_type;
    if (context_.lookup_method(obj_type, c.label, method_name, method_type)) {
      ir::IrType ft = method_type;
      while (ft.is(ir::IrTypeCase::ForallType) && ft.forall && ft.forall->type) {
        ft = *ft.forall->type;
      }
      if (ft.is(ir::IrTypeCase::FunType) && ft.fun) {
        bool is_multi_arg = ft.fun->arg && ft.fun->arg->is(ir::IrTypeCase::TupleType) && !ft.fun->arg->elems().empty();
        if (!is_multi_arg) {
          ir::IrType expected_receiver = ft.fun->arg ? *ft.fun->arg : ir::IrType{};
          bool expects_ptr = expected_receiver.is(ir::IrTypeCase::AppType) && base_type_name(expected_receiver) == "Ptr";
          bool is_ptr = obj_type.is(ir::IrTypeCase::AppType) && base_type_name(obj_type) == "Ptr";

          ir::IrTerm adjusted_s = *c.term;
          if (expects_ptr && !is_ptr) {
            adjusted_s = ir::new_app_term(ir::new_var_term("Ptr::mk"), *c.term);
            adjusted_s.pos = c.term->pos;
          } else if (!expects_ptr && is_ptr) {
            adjusted_s = ir::new_app_term(ir::new_var_term("Ptr::get"), *c.term);
            adjusted_s.pos = c.term->pos;
          }

          ir::IrTerm new_arg = adjusted_s;
          ir::IrTerm new_fun = ir::new_var_term(method_name);
          new_fun.pos = term->pos;
          *term = ir::new_app_term(std::move(new_fun), std::move(new_arg));
          return infer(term, parent_term, expect_type);
        }
      }
    }

    err_ = "type " + obj_type.to_string() + " has no field or method named \"" + c.label + "\"";
    return false;
  }

  bool infer_return(const ir::IrType& evar, ir::IrTerm* term, ir::IrTerm* parent_term, const ir::IrType* expect_type) {
    auto& c = *term->return_data;
    if (expect_returns_.empty()) {
      err_ = "return statement must be inside a function or lambda";
      return false;
    }

    ir::IrType ret_type = expect_returns_.front();
    unify(evar, ret_type);

    if (!infer(c.expr.get(), term, &ret_type)) {
      return false;
    }

    if (c.expr->type) {
      unify(evar, *c.expr->type);
      term->type = *c.expr->type;
    }
    return true;
  }

  bool infer_set(const ir::IrType& evar, ir::IrTerm* term, ir::IrTerm* parent_term, const ir::IrType* expect_type) {
    auto& c = *term->set_data;
    if (!infer(c.term.get(), parent_term, expect_type)) {
      return false;
    }

    ir::IrType obj_type;
    if (c.term->type) {
      obj_type = reduce_and_predicate_type(*c.term->type);
      c.reduced_type = obj_type;
    }

    if (obj_type.is(ir::IrTypeCase::StructType)) {
      for (auto& lv : c.values) {
        int idx = -1;
        ir::StructField sf;
        if (!obj_type.field_by_label(lv.label, idx, sf)) {
          err_ = "field " + lv.label + " not found in struct " + obj_type.to_string();
          return false;
        }
        ir::IrType field_type = sf.type ? *sf.type : ir::IrType{};
        if (!infer(lv.value.get(), parent_term, &field_type)) {
          return false;
        }
      }
    } else if (obj_type.is(ir::IrTypeCase::TupleType)) {
      for (auto& lv : c.values) {
        int idx = -1;
        ir::IrType elem;
        if (!obj_type.elem_by_label(lv.label, idx, elem)) {
          err_ = "tuple index " + lv.label + " out of bounds in " + obj_type.to_string();
          return false;
        }
        if (!infer(lv.value.get(), parent_term, &elem)) {
          return false;
        }
      }
    } else {
      for (auto& lv : c.values) {
        if (!infer(lv.value.get(), parent_term, nullptr)) {
          return false;
        }
      }
    }

    if (c.term->type) {
      unify(evar, *c.term->type);
      term->type = *c.term->type;
    }
    return true;
  }

  bool infer_struct(const ir::IrType& evar, ir::IrTerm* term, ir::IrTerm* parent_term, const ir::IrType* expect_type) {
    auto& c = *term->struct_data;
    std::optional<ir::IrType> struct_type;

    if (expect_type) {
      ir::IrType typ = reduce_and_predicate_type(*expect_type);
      if (typ.is(ir::IrTypeCase::StructType)) {
        struct_type = typ;
      }
    }

    std::vector<ir::StructField> result_fields;
    for (auto& value : c.values) {
      const ir::IrType* field_type = nullptr;
      if (struct_type) {
        for (const auto& f : struct_type->fields()) {
          if (f.id == value.label && f.type) {
            field_type = f.type.get();
            break;
          }
        }
      }

      if (!infer(value.value.get(), term, field_type)) {
        return false;
      }

      ir::IrType val_type = value.value->type ? *value.value->type : ir::IrType{};
      result_fields.push_back(ir::StructField{value.label, std::make_shared<ir::IrType>(std::move(val_type))});
    }

    if (struct_type) {
      unify(evar, *struct_type);
      term->type = *struct_type;
    } else {
      ir::IrType typ = ir::new_struct_type(std::move(result_fields));
      unify(evar, typ);
      term->type = typ;
    }
    return true;
  }

  bool infer_tuple(const ir::IrType& evar, ir::IrTerm* term, ir::IrTerm* parent_term, const ir::IrType* expect_type) {
    auto& c = *term->tuple_data;
    std::optional<ir::IrType> tuple_type;

    if (expect_type) {
      ir::IrType typ = reduce_and_predicate_type(*expect_type);
      if (typ.is(ir::IrTypeCase::TupleType)) {
        tuple_type = typ;
      }
    }

    std::vector<ir::IrType> tuple_elems;
    if (tuple_type) {
      tuple_elems = tuple_type->elems();
    }

    std::vector<ir::IrType> result_elems;
    for (size_t i = 0; i < c.elems.size(); ++i) {
      const ir::IrType* elem_type = nullptr;
      if (tuple_type && i < tuple_elems.size()) {
        elem_type = &tuple_elems[i];
      }

      if (!infer(&c.elems[i], term, elem_type)) {
        return false;
      }

      ir::IrType el_type = c.elems[i].type ? *c.elems[i].type : ir::IrType{};
      result_elems.push_back(std::move(el_type));
    }

    if (tuple_type) {
      unify(evar, *tuple_type);
      term->type = *tuple_type;
    } else {
      ir::IrType typ = ir::new_tuple_type(std::move(result_elems));
      unify(evar, typ);
      term->type = typ;
    }
    return true;
  }

  bool infer_type_abs(const ir::IrType& evar, ir::IrTerm* term, ir::IrTerm* parent_term, const ir::IrType* expect_type) {
    auto& c = *term->type_abs;
    Context orig_context = context_;

    try {
      context_ = context_.enter_function({c.type_param}, {});
    } catch (const std::exception& e) {
      err_ = e.what();
      return false;
    }

    if (!infer(c.body.get(), term, nullptr)) {
      context_ = orig_context;
      return false;
    }

    if (c.body->type) {
      ir::IrType typ = ir::new_forall_type(c.type_param, *c.body->type);
      unify(evar, typ);
      term->type = typ;
    }

    context_ = orig_context;
    return true;
  }

  bool infer_var(const ir::IrType& evar, ir::IrTerm* term, ir::IrTerm* parent_term, const ir::IrType* expect_type) {
    auto& c = *term->var_data;
    Binding bind;
    if (context_.lookup_term_decl_or_def_bind(c.id, bind)) {
      if (bind.is_term_decl() && bind.term_decl) {
        unify(evar, bind.term_decl->type);
        term->type = bind.term_decl->type;
        return true;
      }
      if (bind.is_term_def() && bind.term_def) {
        unify(evar, bind.term_def->type);
        term->type = bind.term_def->type;
        return true;
      }
    }
    err_ = "term variable \"" + c.id + "\" is undefined";
    return false;
  }

  bool infer_impl(const ir::IrType& evar, ir::IrTerm* term, ir::IrTerm* parent_term, const ir::IrType* expect_type) {
    switch (term->case_val) {
      case ir::IrTermCase::AppTermTerm:
        return infer_app_term(evar, term, parent_term, expect_type);
      case ir::IrTermCase::AppTypeTerm:
        return infer_app_type(evar, term, parent_term, expect_type);
      case ir::IrTermCase::AssignTerm:
        return infer_assign(evar, term, parent_term, expect_type);
      case ir::IrTermCase::BlockTerm:
        return infer_block(evar, term, parent_term, expect_type);
      case ir::IrTermCase::ConstTerm:
        return infer_const(evar, term, parent_term, expect_type);
      case ir::IrTermCase::InjectionTerm:
        return infer_injection(evar, term, expect_type);
      case ir::IrTermCase::LambdaTerm:
        return infer_lambda(evar, term, expect_type);
      case ir::IrTermCase::LetTerm:
        return infer_let(evar, term, parent_term, expect_type);
      case ir::IrTermCase::MatchTerm:
        return infer_match(evar, term, expect_type);
      case ir::IrTermCase::ProjectionTerm:
        return infer_projection(evar, term, parent_term, expect_type);
      case ir::IrTermCase::ReturnTerm:
        return infer_return(evar, term, parent_term, expect_type);
      case ir::IrTermCase::SetTerm:
        return infer_set(evar, term, parent_term, expect_type);
      case ir::IrTermCase::StructTerm:
        return infer_struct(evar, term, parent_term, expect_type);
      case ir::IrTermCase::TupleTerm:
        return infer_tuple(evar, term, parent_term, expect_type);
      case ir::IrTermCase::TypeAbsTerm:
        return infer_type_abs(evar, term, parent_term, expect_type);
      case ir::IrTermCase::VarTerm:
        return infer_var(evar, term, parent_term, expect_type);
    }
    err_ = "unhandled term case in infer";
    return false;
  }

  Context context_;
  UnificationState state_;
  List<ir::IrType> expect_returns_;
  std::string err_;
};

} // namespace ts
