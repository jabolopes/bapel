#pragma once

#include "bin/ir_base.h"
#include "bin/ir_type.h"
#include "ts/bind.h"
#include "ts/context.h"
#include "ts/predicative.h"
#include "ts/wellformed.h"

#include <map>
#include <memory>
#include <stdexcept>
#include <string>
#include <vector>

namespace ts {

inline bool same_tags(const ir::IrType& left, const ir::IrType& right) {
  auto lt = left.tags();
  auto rt = right.tags();
  if (lt.size() != rt.size()) return false;
  for (size_t i = 0; i < lt.size(); ++i) {
    if (lt[i].id != rt[i].id) return false;
  }
  return true;
}

inline bool same_fields(const ir::IrType& left, const ir::IrType& right) {
  auto lf = left.fields();
  auto rf = right.fields();
  if (lf.size() != rf.size()) return false;
  for (size_t i = 0; i < lf.size(); ++i) {
    if (lf[i].id != rf[i].id) return false;
  }
  return true;
}

struct ExistVarData {
  std::shared_ptr<ir::IrType> solution;
};

class UnificationState {
 public:
  std::map<int64_t, ExistVarData> exist_vars;

  ir::IrType new_evar(const Context& ctx) {
    ir::IrType evar = ctx.gen_fresh_exist_var();
    exist_vars[evar.exist_var] = ExistVarData{nullptr};
    return evar;
  }

  bool is_exist_var_unassigned(const ir::IrType& tvar) const {
    if (!tvar.is(ir::IrTypeCase::ExistVarType)) return false;
    auto it = exist_vars.find(tvar.exist_var);
    return it != exist_vars.end() && it->second.solution == nullptr;
  }

  bool is_exist_var_assigned(const ir::IrType& tvar) const {
    if (!tvar.is(ir::IrTypeCase::ExistVarType)) return false;
    auto it = exist_vars.find(tvar.exist_var);
    return it != exist_vars.end() && it->second.solution != nullptr;
  }

  ir::IrType exist_var_solution(const ir::IrType& evar) {
    if (!evar.is(ir::IrTypeCase::ExistVarType)) {
      throw std::runtime_error("expected existential variable; got " + evar.to_string());
    }
    auto& data = exist_vars[evar.exist_var];
    if (!data.solution) {
      return evar;
    }
    if (is_exist_var_assigned(*data.solution)) {
      ir::IrType typ = exist_var_solution(*data.solution);
      data.solution = std::make_shared<ir::IrType>(typ);
      return typ;
    }
    return *data.solution;
  }

  bool can_assign(const ir::IrType& evar, const ir::IrType& typ) const {
    if (typ.is(ir::IrTypeCase::ForallType)) {
      return false;
    }
    if (typ.is(ir::IrTypeCase::ExistVarType)) {
      return typ.exist_var < evar.exist_var;
    }
    return true;
  }

  void solve_exist_var(const ir::IrType& evar, const ir::IrType& typ) {
    if (!evar.is(ir::IrTypeCase::ExistVarType)) {
      throw std::runtime_error("expected existential variable; got " + evar.to_string());
    }
    auto& data = exist_vars[evar.exist_var];
    if (data.solution != nullptr) {
      throw std::runtime_error(evar.to_string() + " is already solved with solution " + data.solution->to_string());
    }
    data.solution = std::make_shared<ir::IrType>(typ);
  }

  ir::IrType solve_type(const ir::IrType& typ) {
    switch (typ.case_val) {
      case ir::IrTypeCase::AppType: {
        if (!typ.app) return typ;
        ir::IrType fun = typ.app->fun ? solve_type(*typ.app->fun) : ir::IrType{};
        ir::IrType arg = typ.app->arg ? solve_type(*typ.app->arg) : ir::IrType{};
        return ir::new_app_type(std::move(fun), std::move(arg));
      }
      case ir::IrTypeCase::ArrayType: {
        if (!typ.array) return typ;
        ir::IrType elem = typ.array->elem_type ? solve_type(*typ.array->elem_type) : ir::IrType{};
        return ir::new_array_type(std::move(elem), typ.array->size);
      }
      case ir::IrTypeCase::ExistVarType: {
        if (is_exist_var_assigned(typ)) {
          return solve_type(exist_var_solution(typ));
        }
        return typ;
      }
      case ir::IrTypeCase::ForallType: {
        if (!typ.forall) return typ;
        ir::IrType body = typ.forall->type ? solve_type(*typ.forall->type) : ir::IrType{};
        return ir::new_forall_type(typ.forall->type_param, std::move(body));
      }
      case ir::IrTypeCase::FunType: {
        if (!typ.fun) return typ;
        ir::IrType arg = typ.fun->arg ? solve_type(*typ.fun->arg) : ir::IrType{};
        ir::IrType ret = typ.fun->ret ? solve_type(*typ.fun->ret) : ir::IrType{};
        return ir::new_function_type(std::move(arg), std::move(ret));
      }
      case ir::IrTypeCase::LambdaType: {
        if (!typ.lambda) return typ;
        ir::IrType body = typ.lambda->type ? solve_type(*typ.lambda->type) : ir::IrType{};
        return ir::new_lambda_type(typ.lambda->var, typ.lambda->kind, std::move(body));
      }
      case ir::IrTypeCase::NameType:
        return typ;
      case ir::IrTypeCase::StructType: {
        if (!typ.struct_data) return typ;
        std::vector<ir::StructField> fields;
        fields.reserve(typ.struct_data->fields.size());
        for (const auto& f : typ.struct_data->fields) {
          ir::IrType ft = f.type ? solve_type(*f.type) : ir::IrType{};
          fields.push_back(ir::StructField{f.id, std::make_shared<ir::IrType>(std::move(ft))});
        }
        return ir::new_struct_type(std::move(fields));
      }
      case ir::IrTypeCase::TupleType: {
        if (!typ.tuple_data) return typ;
        std::vector<ir::IrType> elems;
        elems.reserve(typ.tuple_data->elems.size());
        for (const auto& e : typ.tuple_data->elems) {
          elems.push_back(solve_type(e));
        }
        return ir::new_tuple_type(std::move(elems));
      }
      case ir::IrTypeCase::VariantType: {
        if (!typ.variant_data) return typ;
        std::vector<ir::VariantTag> tags;
        tags.reserve(typ.variant_data->tags.size());
        for (const auto& tag : typ.variant_data->tags) {
          ir::IrType tt = tag.type ? solve_type(*tag.type) : ir::IrType{};
          tags.push_back(ir::VariantTag{tag.id, std::make_shared<ir::IrType>(std::move(tt))});
        }
        return ir::new_variant_type(std::move(tags));
      }
      case ir::IrTypeCase::VarType:
        return typ;
    }
    return typ;
  }
};

inline bool unify_impl(Context& context, UnificationState& state, ir::IrType left, ir::IrType right, std::string& err) {
  if (state.is_exist_var_assigned(left)) {
    left = state.exist_var_solution(left);
    return unify_impl(context, state, left, right, err);
  }

  if (state.is_exist_var_assigned(right)) {
    right = state.exist_var_solution(right);
    return unify_impl(context, state, left, right, err);
  }

  if (state.is_exist_var_unassigned(left) && state.can_assign(left, right)) {
    state.solve_exist_var(left, right);
    return true;
  }

  if (state.is_exist_var_unassigned(right) && state.can_assign(right, left)) {
    state.solve_exist_var(right, left);
    return true;
  }

  if (left.is(ir::IrTypeCase::AppType) && right.is(ir::IrTypeCase::AppType)) {
    if (!unify_impl(context, state, *left.app->fun, *right.app->fun, err)) {
      err = "mismatch in function types: " + err;
      return false;
    }
    if (!unify_impl(context, state, *left.app->arg, *right.app->arg, err)) {
      err = "mismatch in argument types: " + err;
      return false;
    }
    return true;
  }

  if (left.is(ir::IrTypeCase::ArrayType) && right.is(ir::IrTypeCase::ArrayType)) {
    if (!unify_impl(context, state, *left.array->elem_type, *right.array->elem_type, err)) {
      err = "mismatch in array element types: " + err;
      return false;
    }
    if (left.array->size != right.array->size) {
      err = "expected array with " + std::to_string(left.array->size) + " elements; got " + std::to_string(right.array->size) + " elements";
      return false;
    }
    return true;
  }

  if (right.is(ir::IrTypeCase::ForallType)) {
    auto [new_ctx, tp, body] = context.add_fresh_type(right);
    context = new_ctx;
    return unify_impl(context, state, left, body, err);
  }

  if (left.is(ir::IrTypeCase::ForallType)) {
    ir::IrType evar = state.new_evar(context);
    ir::IrType new_body = ir::substitute_type(*left.forall->type, ir::new_var_type(left.forall->type_param.var), evar);
    return unify_impl(context, state, new_body, right, err);
  }

  if (left.is(ir::IrTypeCase::FunType) && right.is(ir::IrTypeCase::FunType)) {
    if (!unify_impl(context, state, *left.fun->arg, *right.fun->arg, err)) {
      return false;
    }
    return unify_impl(context, state, *left.fun->ret, *right.fun->ret, err);
  }

  if (left.is(ir::IrTypeCase::NameType) && context.contains_const_bind(left.name) &&
      right.is(ir::IrTypeCase::NameType) && context.contains_const_bind(right.name) &&
      left.name == right.name) {
    return true;
  }

  if (left.is(ir::IrTypeCase::NameType) && context.contains_alias_bind(left.name)) {
    Binding bind;
    if (context.lookup_alias_bind(left.name, bind) && bind.alias) {
      return unify_impl(context, state, bind.alias->type, right, err);
    }
  }

  if (right.is(ir::IrTypeCase::NameType) && context.contains_alias_bind(right.name)) {
    Binding bind;
    if (context.lookup_alias_bind(right.name, bind) && bind.alias) {
      return unify_impl(context, state, left, bind.alias->type, err);
    }
  }

  if (left.is(ir::IrTypeCase::StructType) && right.is(ir::IrTypeCase::StructType) && same_fields(left, right)) {
    auto lf = left.fields();
    auto rf = right.fields();
    for (size_t i = 0; i < lf.size(); ++i) {
      if (!unify_impl(context, state, *lf[i].type, *rf[i].type, err)) {
        return false;
      }
    }
    return true;
  }

  if (left.is(ir::IrTypeCase::TupleType) && right.is(ir::IrTypeCase::TupleType) && left.elems().size() == right.elems().size()) {
    auto le = left.elems();
    auto re = right.elems();
    for (size_t i = 0; i < le.size(); ++i) {
      if (!unify_impl(context, state, le[i], re[i], err)) {
        return false;
      }
    }
    return true;
  }

  if (left.is(ir::IrTypeCase::VariantType) && right.is(ir::IrTypeCase::VariantType) && same_tags(left, right)) {
    auto lt = left.tags();
    auto rt = right.tags();
    for (size_t i = 0; i < lt.size(); ++i) {
      if (!unify_impl(context, state, *lt[i].type, *rt[i].type, err)) {
        return false;
      }
    }
    return true;
  }

  if (left.is(ir::IrTypeCase::VarType) && right.is(ir::IrTypeCase::VarType) && left.var == right.var) {
    return true;
  }

  err = "expected type " + left.to_string() + "; got " + right.to_string();
  return false;
}

inline bool unify(Context& context, UnificationState& state, const ir::IrType& left, const ir::IrType& right, std::string& err) {
  if (!is_wellformed_type(context, left, err)) {
    return false;
  }
  if (!is_wellformed_type(context, right, err)) {
    return false;
  }
  if (!unify_impl(context, state, left, right, err)) {
    err += "\n  unifying " + left.to_string() + " and " + right.to_string();
    return false;
  }
  return true;
}

} // namespace ts
