#pragma once

#include "bin/ir_base.h"
#include "bin/ir_type.h"
#include "ts/bind.h"
#include "ts/context.h"
#include "ts/predicative.h"
#include "ts/wellformed.h"

#include <stdexcept>
#include <string>
#include <vector>

namespace ts {

inline bool satisfies_bound(const Context& context, const ir::IrType& self_type, const ir::IrType& bound, std::string& err) {
  if (self_type.is(ir::IrTypeCase::VarType)) {
    Binding bind;
    if (context.lookup_type_param_bind(self_type.var, bind) && bind.type_param) {
      for (const auto& b : bind.type_param->bounds) {
        if (ir::equals_type(b, bound)) {
          return true;
        }
      }
    }
  }

  Binding impl_bind;
  if (context.lookup_trait_impl_bind(bound, self_type, impl_bind)) {
    return true;
  }

  err = "type " + self_type.to_string() + " does not implement trait " + bound.to_string();
  return false;
}

inline bool satisfies_bounds(const Context& context, const ir::IrType& self_type, const std::vector<ir::IrType>& bounds, std::string& err) {
  for (const auto& b : bounds) {
    if (!satisfies_bound(context, self_type, b, err)) {
      return false;
    }
  }
  return true;
}

inline bool subtype(Context& context, const ir::IrType& left, const ir::IrType& right, std::string& err);

inline bool subtype_impl(Context& context, ir::IrType left, ir::IrType right, std::string& err) {
  left = reduce_and_predicate_type(context, left);
  right = reduce_and_predicate_type(context, right);

  if (left.is(ir::IrTypeCase::AppType) && right.is(ir::IrTypeCase::AppType)) {
    if (!subtype(context, *left.app->fun, *right.app->fun, err)) {
      err = "mismatch in function types: " + err;
      return false;
    }
    if (!subtype(context, *left.app->arg, *right.app->arg, err)) {
      err = "mismatch in argument types: " + err;
      return false;
    }
    return true;
  }

  if (left.is(ir::IrTypeCase::ArrayType) && right.is(ir::IrTypeCase::ArrayType)) {
    if (!subtype(context, *left.array->elem_type, *right.array->elem_type, err)) {
      err = "mismatch in array element types: " + err;
      return false;
    }
    if (left.array->size != right.array->size) {
      err = "expected array with " + std::to_string(left.array->size) + " elements; got " + std::to_string(right.array->size) + " elements";
      return false;
    }
    return true;
  }

  if (left.is(ir::IrTypeCase::ForallType) && right.is(ir::IrTypeCase::ForallType)) {
    Context orig_context = context;

    auto [new_ctx, tvar, right_body_type] = context.add_fresh_type(right);
    context = new_ctx;

    std::vector<ir::IrType> left_bounds;
    left_bounds.reserve(left.forall->type_param.bounds.size());
    for (const auto& b : left.forall->type_param.bounds) {
      left_bounds.push_back(ir::substitute_type(b, ir::new_var_type(left.forall->type_param.var), ir::new_var_type(tvar.var)));
    }

    for (const auto& lb : left_bounds) {
      bool found = false;
      for (const auto& rb : tvar.bounds) {
        if (ir::equals_type(lb, rb)) {
          found = true;
          break;
        }
      }
      if (!found) {
        context = orig_context;
        err = "bounds mismatch: right bounds do not satisfy left bounds";
        return false;
      }
    }

    ir::IrType left_body_type = ir::substitute_type(*left.forall->type, ir::new_var_type(left.forall->type_param.var), ir::new_var_type(tvar.var));
    bool ok = subtype(context, left_body_type, right_body_type, err);
    context = orig_context;
    return ok;
  }

  if (left.is(ir::IrTypeCase::FunType) && right.is(ir::IrTypeCase::FunType)) {
    // Contravariant argument
    if (!subtype(context, *right.fun->arg, *left.fun->arg, err)) {
      return false;
    }
    // Covariant return
    return subtype(context, *left.fun->ret, *right.fun->ret, err);
  }

  if (left.is(ir::IrTypeCase::StructType) && right.is(ir::IrTypeCase::StructType)) {
    auto lf = left.fields();
    auto rf = right.fields();
    if (lf.size() != rf.size()) {
      err = "expected " + std::to_string(left.fields().size()) + " fields; got " + std::to_string(right.fields().size());
      return false;
    }
    for (size_t i = 0; i < lf.size(); ++i) {
      if (lf[i].id != rf[i].id) {
        err = "field name mismatch: " + lf[i].id + " != " + rf[i].id;
        return false;
      }
      if (!subtype(context, *lf[i].type, *rf[i].type, err)) {
        return false;
      }
    }
    return true;
  }

  if (left.is(ir::IrTypeCase::TupleType) && right.is(ir::IrTypeCase::TupleType)) {
    auto le = left.elems();
    auto re = right.elems();
    if (le.size() != re.size()) {
      err = "expected " + std::to_string(le.size()) + " elements; got " + std::to_string(re.size());
      return false;
    }
    for (size_t i = 0; i < le.size(); ++i) {
      if (!subtype(context, le[i], re[i], err)) {
        return false;
      }
    }
    return true;
  }

  if (left.is(ir::IrTypeCase::VariantType) && right.is(ir::IrTypeCase::VariantType)) {
    auto lt = left.tags();
    auto rt = right.tags();
    if (lt.size() != rt.size()) {
      err = "expected " + std::to_string(lt.size()) + " tags; got " + std::to_string(rt.size());
      return false;
    }
    for (size_t i = 0; i < lt.size(); ++i) {
      if (lt[i].id != rt[i].id) {
        err = "tag name mismatch: " + lt[i].id + " != " + rt[i].id;
        return false;
      }
      if (!subtype(context, *lt[i].type, *rt[i].type, err)) {
        return false;
      }
    }
    return true;
  }

  if (left.is(ir::IrTypeCase::VarType) && right.is(ir::IrTypeCase::VarType) && left.var == right.var) {
    return true;
  }

  if (left.is(ir::IrTypeCase::NameType) && context.contains_const_bind(left.name) &&
      right.is(ir::IrTypeCase::NameType) && context.contains_const_bind(right.name) &&
      left.name == right.name) {
    return true;
  }

  if (left.is(ir::IrTypeCase::NameType) && context.contains_alias_bind(left.name)) {
    Binding bind;
    if (context.lookup_alias_bind(left.name, bind) && bind.alias) {
      return subtype(context, bind.alias->type, right, err);
    }
  }

  if (right.is(ir::IrTypeCase::NameType) && context.contains_alias_bind(right.name)) {
    Binding bind;
    if (context.lookup_alias_bind(right.name, bind) && bind.alias) {
      return subtype(context, left, bind.alias->type, err);
    }
  }

  err = "expected type " + left.to_string() + "; got " + right.to_string();
  return false;
}

inline bool subtype(Context& context, const ir::IrType& left, const ir::IrType& right, std::string& err) {
  if (!is_wellformed_type(context, left, err)) {
    return false;
  }
  if (!is_wellformed_type(context, right, err)) {
    return false;
  }
  if (!subtype_impl(context, left, right, err)) {
    err += "\n  subtyping " + left.to_string() + " and " + right.to_string();
    return false;
  }
  return true;
}

} // namespace ts
