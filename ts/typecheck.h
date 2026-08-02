#pragma once

#include "bin/ir_base.h"
#include "bin/ir_decl.h"
#include "bin/ir_function.h"
#include "bin/ir_term.h"
#include "bin/ir_type.h"
#include "ts/bind.h"
#include "ts/context.h"
#include "ts/infer_kind.h"
#include "ts/inferencer.h"
#include "ts/predicative.h"
#include "ts/reduce_type.h"
#include "ts/returns.h"
#include "ts/subtype.h"
#include "ts/wellformed.h"

#include <memory>
#include <stdexcept>
#include <string>
#include <vector>

namespace ts {

class Typechecker {
 public:
  explicit Typechecker(Context context) : context_(std::move(context)) {}

  ir::IrType reduce_type(const ir::IrType& typ) {
    return ts::reduce_type(context_, typ);
  }

  ir::IrType reduce_and_predicate_type(const ir::IrType& typ) {
    return ts::reduce_and_predicate_type(context_, typ);
  }

  bool subtype(const ir::IrType& left, const ir::IrType& right) {
    return ts::subtype(context_, left, right, err_);
  }

  bool infer_function(ir::IrFunction* function, Context& out_context) {
    Inferencer inferencer(context_);
    bool ok = inferencer.infer_function(function, out_context);
    if (!ok) {
      err_ = function->pos.to_string() + ":\n" + inferencer.error();
    }
    return ok;
  }

  bool typecheck_function(ir::IrFunction* function, Context& out_context) {
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
    bool body_ok = typecheck_term(&function->body);
    expect_returns_ = expect_returns_.remove();
    if (!body_ok) {
      context_ = orig_context;
      return false;
    }

    auto last = last_terms(&function->body);
    for (auto* term : last) {
      if (term->is(ir::IrTermCase::ReturnTerm)) {
        err_ = term->pos.to_string() + ":\n redundant 'return' statement as the last term of a function";
        context_ = orig_context;
        return false;
      }
      if (term->type) {
        if (!subtype(function->ret_type, *term->type)) {
          err_ = term->pos.to_string() + ":\n" + err_;
          context_ = orig_context;
          return false;
        }
      }
    }

    if (function->body.is(ir::IrTermCase::BlockTerm) && last.empty()) {
      err_ = function->body.pos.to_string() + ":\nexpected non-empty function block";
      context_ = orig_context;
      return false;
    }

    context_ = out_context;
    return true;
  }

  bool typecheck_term(ir::IrTerm* term) {
    if (!term) return true;
    if (!typecheck_impl(term)) {
      err_ = term->pos.to_string() + ":\n" + err_;
      return false;
    }
    return true;
  }

  const Context& context() const { return context_; }
  const std::string& error() const { return err_; }

 private:
  bool resolve_trait_method_app(ir::IrTerm* term, std::string& trait_name, std::vector<ir::IrType>& type_args) {
    type_args.clear();
    ir::IrTerm* curr = term;
    while (curr->is(ir::IrTermCase::AppTypeTerm) && curr->app_type) {
      type_args.push_back(curr->app_type->arg);
      curr = curr->app_type->fun.get();
    }
    std::reverse(type_args.begin(), type_args.end());

    if (curr->is(ir::IrTermCase::VarTerm) && curr->var_data) {
      size_t pos = curr->var_data->id.find("::");
      if (pos != std::string::npos) {
        trait_name = curr->var_data->id.substr(0, pos);
        return context_.contains_trait_bind(trait_name);
      }
    }
    return false;
  }

  bool typecheck_app_term(ir::IrTerm* term) {
    auto& c = *term->app_term;
    if (!typecheck_term(c.fun.get())) return false;
    if (!typecheck_term(c.arg.get())) return false;

    if (!c.fun->type || !c.fun->type->is(ir::IrTypeCase::FunType) || !c.fun->type->fun) {
      err_ = "expected term to have function type";
      return false;
    }

    const auto& fun_type = *c.fun->type->fun;
    if (c.arg->type) {
      if (!subtype(*c.arg->type, *fun_type.arg)) return false;
    }

    if (term->type) {
      if (!subtype(*fun_type.ret, *term->type)) return false;
    }
    return true;
  }

  bool typecheck_app_type(ir::IrTerm* term) {
    auto& c = *term->app_type;
    if (!typecheck_term(c.fun.get())) return false;

    if (!c.fun->type || !c.fun->type->is(ir::IrTypeCase::ForallType) || !c.fun->type->forall) {
      err_ = "expected term to have forall type";
      return false;
    }

    const auto& forall_type = *c.fun->type->forall;
    if (!is_wellformed_type(context_, c.arg, err_)) {
      return false;
    }

    ir::IrKind arg_kind;
    if (!infer_kind(context_, c.arg, arg_kind, err_)) {
      return false;
    }

    if (!ir::equals_kind(forall_type.type_param.kind, arg_kind)) {
      err_ = "expected argument in type application to match forall type's kind";
      return false;
    }

    if (!satisfies_bounds(context_, c.arg, forall_type.type_param.bounds, err_)) {
      return false;
    }

    std::string trait_name;
    std::vector<ir::IrType> type_args;
    if (resolve_trait_method_app(term, trait_name, type_args)) {
      Binding trait_bind;
      if (context_.lookup_trait_bind(trait_name, trait_bind) && trait_bind.trait) {
        size_t expected_args = trait_bind.trait->type_params.size() + 1;
        if (type_args.size() == expected_args) {
          ir::IrType self_type = reduce_type(type_args[0]);
          ir::IrType trait_type = ir::new_name_type(trait_name);
          for (size_t i = 1; i < type_args.size(); ++i) {
            trait_type = ir::new_app_type(trait_type, reduce_type(type_args[i]));
          }
          if (!satisfies_bound(context_, self_type, trait_type, err_)) {
            return false;
          }
        }
      }
    }

    ir::IrType typ = ir::substitute_type(*forall_type.type, ir::new_var_type(forall_type.type_param.var), c.arg);
    if (term->type) {
      if (!subtype(typ, *term->type)) return false;
    }
    return true;
  }

  bool typecheck_block(ir::IrTerm* term) {
    auto& c = *term->block;
    Context orig_context = context_;

    try {
      context_ = context_.enter_scope();
    } catch (const std::exception& e) {
      err_ = e.what();
      return false;
    }

    for (auto& t : c.terms) {
      if (!typecheck_term(&t)) {
        context_ = orig_context;
        return false;
      }
    }

    if (!c.terms.empty() && c.terms.back().type && term->type) {
      if (!subtype(*c.terms.back().type, *term->type)) {
        context_ = orig_context;
        return false;
      }
    }

    context_ = orig_context;
    return true;
  }

  bool typecheck_const(ir::IrTerm* term) {
    auto& c = *term->const_data;
    if (c.literal.is_rune()) {
      ir::IrType typ = ir::new_name_type("i8");
      if (term->type) {
        if (!subtype(typ, *term->type)) return false;
      }
    } else if (c.literal.is_str()) {
      ir::IrType typ = ir::new_name_type("StringView");
      if (term->type) {
        if (!subtype(typ, *term->type)) return false;
      }
    } else {
      if (term->type) {
        ir::IrKind kind;
        if (!infer_kind(context_, *term->type, kind, err_)) return false;
        if (!kind.is_type_kind()) {
          err_ = "expected type to have kind *";
          return false;
        }
      }
    }
    return true;
  }

  bool typecheck_injection(ir::IrTerm* term) {
    auto& c = *term->injection;
    ir::IrType variant_type = reduce_and_predicate_type(c.variant_type);
    if (!variant_type.is(ir::IrTypeCase::VariantType)) {
      err_ = "expected variant type";
      return false;
    }

    ir::IrKind kind;
    if (!infer_kind(context_, variant_type, kind, err_)) return false;
    if (!kind.is_type_kind()) {
      err_ = "expected variant type to have kind *";
      return false;
    }

    bool found = false;
    ir::IrType tag_type;
    for (const auto& tag : variant_type.tags()) {
      if (tag.id == c.tag) {
        tag_type = tag.type ? *tag.type : ir::IrType{};
        found = true;
        break;
      }
    }
    if (!found) {
      err_ = "tag " + c.tag + " not found in variant";
      return false;
    }

    if (!typecheck_term(c.value.get())) return false;

    if (c.value->type) {
      if (!subtype(*c.value->type, tag_type)) return false;
    }

    if (term->type) {
      if (!subtype(variant_type, *term->type)) return false;
    }
    return true;
  }

  bool typecheck_lambda(ir::IrTerm* term) {
    auto& c = *term->lambda;
    ir::IrKind arg_kind;
    if (!infer_kind(context_, c.arg.type, arg_kind, err_)) return false;
    if (!arg_kind.is_type_kind()) {
      err_ = "expected lambda argument to have kind *";
      return false;
    }

    Context orig_context = context_;
    try {
      context_ = context_.enter_function({}, {c.arg});
    } catch (const std::exception& e) {
      err_ = e.what();
      return false;
    }

    ir::IrType lambda_ret_type;
    if (term->type && term->type->is(ir::IrTypeCase::FunType) && term->type->fun && term->type->fun->ret) {
      lambda_ret_type = *term->type->fun->ret;
    } else if (c.body && c.body->type) {
      lambda_ret_type = *c.body->type;
    }

    expect_returns_ = expect_returns_.add(lambda_ret_type);
    bool body_ok = typecheck_term(c.body.get());
    expect_returns_ = expect_returns_.remove();
    if (!body_ok) {
      context_ = orig_context;
      return false;
    }

    auto last = last_terms(c.body.get());
    for (auto* t : last) {
      if (t->is(ir::IrTermCase::ReturnTerm)) {
        err_ = "redundant 'return' statement as the last term of a function";
        context_ = orig_context;
        return false;
      }
      if (t->type && c.body->type) {
        if (!subtype(*c.body->type, *t->type)) {
          context_ = orig_context;
          return false;
        }
      }
    }

    if (c.body->type && term->type) {
      ir::IrType typ = ir::new_function_type(c.arg.type, *c.body->type);
      if (!subtype(typ, *term->type)) {
        context_ = orig_context;
        return false;
      }
    }

    context_ = orig_context;
    return true;
  }

  bool typecheck_let(ir::IrTerm* term) {
    auto& c = *term->let_data;
    ir::IrType var_type = c.var_type ? *c.var_type : (c.value->type ? *c.value->type : ir::IrType{});

    if (c.var_type) {
      if (!is_wellformed_type(context_, *c.var_type, err_)) return false;
      ir::IrKind kind;
      if (!infer_kind(context_, *c.var_type, kind, err_)) return false;
      if (!kind.is_type_kind()) {
        err_ = "expected let variable to have kind *";
        return false;
      }
    }

    try {
      context_ = context_.add_bind(new_term_def_bind(c.var, var_type));
    } catch (const std::exception& e) {
      err_ = e.what();
      return false;
    }

    if (!typecheck_term(c.value.get())) return false;

    if (c.value->type && c.var_type) {
      if (!subtype(*c.value->type, *c.var_type)) return false;
    }

    if (term->type && c.value->type) {
      if (!subtype(*c.value->type, *term->type)) return false;
    }
    return true;
  }

  bool typecheck_match(ir::IrTerm* term) {
    auto& c = *term->match_data;
    if (!typecheck_term(c.term.get())) return false;

    if (!c.term->type) return true;

    ir::IrType variant_type = reduce_and_predicate_type(*c.term->type);
    if (!variant_type.is(ir::IrTypeCase::VariantType)) {
      err_ = "expected variant type in match";
      return false;
    }

    std::optional<ir::IrType> match_type;
    for (auto& arm : c.arms) {
      ir::IrType tag_type;
      bool found = false;
      for (const auto& tag : variant_type.tags()) {
        if (tag.id == arm.tag) {
          tag_type = tag.type ? *tag.type : ir::IrType{};
          found = true;
          break;
        }
      }
      if (!found) {
        err_ = "tag \"" + arm.tag + "\" is not a valid tag of variant type";
        return false;
      }

      Context orig_context = context_;
      try {
        context_ = context_.add_bind(new_term_def_bind(arm.arg, tag_type));
      } catch (const std::exception& e) {
        err_ = e.what();
        return false;
      }

      if (!typecheck_term(arm.body.get())) {
        context_ = orig_context;
        return false;
      }

      if (!match_type && arm.body->type) {
        match_type = *arm.body->type;
      } else if (match_type && arm.body->type) {
        if (!subtype(*arm.body->type, *match_type)) {
          context_ = orig_context;
          return false;
        }
      }

      context_ = orig_context;
    }

    if (match_type && term->type) {
      if (!subtype(*match_type, *term->type)) return false;
    }
    return true;
  }

  bool typecheck_projection(ir::IrTerm* term) {
    auto& c = *term->projection;
    if (!typecheck_term(c.term.get())) return false;
    return true;
  }

  bool typecheck_return(ir::IrTerm* term) {
    auto& c = *term->return_data;
    if (expect_returns_.empty()) {
      err_ = "return statement must be inside a function or lambda";
      return false;
    }
    if (!typecheck_term(c.expr.get())) return false;
    if (c.expr->type) {
      if (!subtype(expect_returns_.front(), *c.expr->type)) {
        return false;
      }
    }
    if (term->type && c.expr->type) {
      if (!subtype(*c.expr->type, *term->type)) return false;
    }
    return true;
  }

  bool typecheck_set(ir::IrTerm* term) {
    auto& c = *term->set_data;
    if (!typecheck_term(c.term.get())) return false;
    for (auto& lv : c.values) {
      if (!typecheck_term(lv.value.get())) return false;
    }
    return true;
  }

  bool typecheck_struct(ir::IrTerm* term) {
    auto& c = *term->struct_data;
    for (auto& fv : c.values) {
      if (!typecheck_term(fv.value.get())) return false;
    }
    return true;
  }

  bool typecheck_tuple(ir::IrTerm* term) {
    auto& c = *term->tuple_data;
    for (auto& elem : c.elems) {
      if (!typecheck_term(&elem)) return false;
    }
    return true;
  }

  bool typecheck_type_abs(ir::IrTerm* term) {
    auto& c = *term->type_abs;
    Context orig_context = context_;
    try {
      context_ = context_.enter_function({c.type_param}, {});
    } catch (const std::exception& e) {
      err_ = e.what();
      return false;
    }

    if (!typecheck_term(c.body.get())) {
      context_ = orig_context;
      return false;
    }

    if (c.body->type && term->type) {
      ir::IrType typ = ir::new_forall_type(c.type_param, *c.body->type);
      if (!subtype(typ, *term->type)) {
        context_ = orig_context;
        return false;
      }
    }

    context_ = orig_context;
    return true;
  }

  bool typecheck_var(ir::IrTerm* term) {
    auto& c = *term->var_data;
    Binding bind;
    if (context_.lookup_term_decl_or_def_bind(c.id, bind)) {
      ir::IrType var_type = bind.is_term_decl() && bind.term_decl ? bind.term_decl->type : (bind.is_term_def() && bind.term_def ? bind.term_def->type : ir::IrType{});
      if (term->type) {
        return subtype(var_type, *term->type);
      }
      return true;
    }
    err_ = "variable \"" + c.id + "\" is undefined";
    return false;
  }

  bool typecheck_impl(ir::IrTerm* term) {
    switch (term->case_val) {
      case ir::IrTermCase::AppTermTerm:
        return typecheck_app_term(term);
      case ir::IrTermCase::AppTypeTerm:
        return typecheck_app_type(term);
      case ir::IrTermCase::AssignTerm:
        return true;
      case ir::IrTermCase::BlockTerm:
        return typecheck_block(term);
      case ir::IrTermCase::ConstTerm:
        return typecheck_const(term);
      case ir::IrTermCase::InjectionTerm:
        return typecheck_injection(term);
      case ir::IrTermCase::LambdaTerm:
        return typecheck_lambda(term);
      case ir::IrTermCase::LetTerm:
        return typecheck_let(term);
      case ir::IrTermCase::MatchTerm:
        return typecheck_match(term);
      case ir::IrTermCase::ProjectionTerm:
        return typecheck_projection(term);
      case ir::IrTermCase::ReturnTerm:
        return typecheck_return(term);
      case ir::IrTermCase::SetTerm:
        return typecheck_set(term);
      case ir::IrTermCase::StructTerm:
        return typecheck_struct(term);
      case ir::IrTermCase::TupleTerm:
        return typecheck_tuple(term);
      case ir::IrTermCase::TypeAbsTerm:
        return typecheck_type_abs(term);
      case ir::IrTermCase::VarTerm:
        return typecheck_var(term);
    }
    err_ = "unhandled term case in typecheck";
    return false;
  }

  Context context_;
  List<ir::IrType> expect_returns_;
  std::string err_;
};

} // namespace ts
