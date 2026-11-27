#pragma once

#include "bin/ir_base.h"
#include "bin/ir_type.h"
#include "ts/bind.h"
#include "ts/context.h"
#include "ts/predicative.h"
#include "ts/wellformed.h"

#include <stdexcept>
#include <string>
#include <utility>
#include <vector>

namespace ts {

inline bool satisfies_bound(const Context& context, const ir::IrType& self_type, const ir::IrType& bound, std::string& err) {
  if (self_type.data.is(ir::IrTypeData::Tag::var_t)) {
    Binding bind;
    if (context.lookup_type_param_bind(self_type.data.as_var_t(), bind) && bind.type_param) {
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

  // If self_type is a MuType (from unfolding a recursive alias's self-reference),
  // resolve it back to its defining alias name and check nominal impls.
  if (auto mu_alias = context.lookup_mu_alias_name(self_type)) {
    ir::IrType alias_type = ir::new_name_type(*mu_alias);
    if (context.lookup_trait_impl_bind(bound, alias_type, impl_bind)) {
      return true;
    }
  }

  err = "type " + ir::type_to_string(self_type) + " does not implement trait " + ir::type_to_string(bound);
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

class Subtyper {
 public:
  bool check_impl(Context& context, ir::IrType left, ir::IrType right, std::string& err) {
    if (ir::equals_type(left, right)) {
      return true;
    }

    left = reduce_whnf(context, left);
    right = reduce_whnf(context, right);

    if (ir::equals_type(left, right)) {
      return true;
    }

    if (is_assumed(left, right)) {
      return true;
    }

    if (left.data.is(ir::IrTypeData::Tag::mu_t) && right.data.is(ir::IrTypeData::Tag::mu_t)) {
      assumed_.push_back({left, right});
      Context orig_context = context;
      auto [new_ctx, tvar, right_body] = context.add_fresh_type(right);
      const auto& lm = left.data.as_mu_t();
      ir::IrType left_body = ir::substitute_type(*lm.body, ir::new_var_type(lm.var), ir::new_var_type(tvar.var));
      bool ok = check_impl(new_ctx, left_body, right_body, err);
      context = orig_context;
      assumed_.pop_back();
      return ok;
    }

    if (left.data.is(ir::IrTypeData::Tag::mu_t)) {
      assumed_.push_back({left, right});
      bool ok = check_impl(context, unfold_mu_type(left), right, err);
      assumed_.pop_back();
      return ok;
    }

    if (right.data.is(ir::IrTypeData::Tag::mu_t)) {
      assumed_.push_back({left, right});
      bool ok = check_impl(context, left, unfold_mu_type(right), err);
      assumed_.pop_back();
      return ok;
    }

    if (left.data.is(ir::IrTypeData::Tag::app_t) && right.data.is(ir::IrTypeData::Tag::app_t)) {
      const auto& la = left.data.as_app_t();
      const auto& ra = right.data.as_app_t();
      if (!check_impl(context, *la.fun, *ra.fun, err)) {
        err = "mismatch in function types: " + err;
        return false;
      }
      if (!check_impl(context, *la.arg, *ra.arg, err)) {
        err = "mismatch in argument types: " + err;
        return false;
      }
      return true;
    }

    if (left.data.is(ir::IrTypeData::Tag::array_t) && right.data.is(ir::IrTypeData::Tag::array_t)) {
      const auto& la = left.data.as_array_t();
      const auto& ra = right.data.as_array_t();
      if (!check_impl(context, *la.elem, *ra.elem, err)) {
        err = "mismatch in array element types: " + err;
        return false;
      }
      if (la.size != ra.size) {
        err = "expected array with " + std::to_string(la.size) + " elements; got " + std::to_string(ra.size) + " elements";
        return false;
      }
      return true;
    }

    if (left.data.is(ir::IrTypeData::Tag::forall_t) && right.data.is(ir::IrTypeData::Tag::forall_t)) {
      Context orig_context = context;

      auto [new_ctx, tvar, right_body_type] = context.add_fresh_type(right);
      context = new_ctx;

      const auto& lf = left.data.as_forall_t();
      std::vector<ir::IrType> left_bounds;
      left_bounds.reserve(lf.param.bounds.size());
      for (const auto& b : lf.param.bounds) {
        left_bounds.push_back(ir::substitute_type(b, ir::new_var_type(lf.param.var), ir::new_var_type(tvar.var)));
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

      ir::IrType left_body_type = ir::substitute_type(*lf.body, ir::new_var_type(lf.param.var), ir::new_var_type(tvar.var));
      bool ok = check_impl(context, left_body_type, right_body_type, err);
      context = orig_context;
      return ok;
    }

    if (left.data.is(ir::IrTypeData::Tag::lambda_t) && right.data.is(ir::IrTypeData::Tag::lambda_t)) {
      Context orig_context = context;
      auto [new_ctx, tvar, right_body] = context.add_fresh_type(right);
      const auto& ll = left.data.as_lambda_t();
      ir::IrType left_body = ir::substitute_type(*ll.body, ir::new_var_type(ll.var), ir::new_var_type(tvar.var));
      bool ok = check_impl(new_ctx, left_body, right_body, err);
      context = orig_context;
      return ok;
    }

    if (left.data.is(ir::IrTypeData::Tag::fun_t) && right.data.is(ir::IrTypeData::Tag::fun_t)) {
      const auto& lf = left.data.as_fun_t();
      const auto& rf = right.data.as_fun_t();
      // Contravariant argument
      if (!check_impl(context, *rf.arg, *lf.arg, err)) {
        return false;
      }
      // Covariant return
      return check_impl(context, *lf.ret, *rf.ret, err);
    }

    if (left.data.is(ir::IrTypeData::Tag::struct_t) && right.data.is(ir::IrTypeData::Tag::struct_t)) {
      auto lf = ir::type_fields(left);
      auto rf = ir::type_fields(right);
      if (lf.size() != rf.size()) {
        err = "expected " + std::to_string(lf.size()) + " fields; got " + std::to_string(rf.size());
        return false;
      }
      for (size_t i = 0; i < lf.size(); ++i) {
        if (lf[i].id != rf[i].id) {
          err = "field name mismatch: " + lf[i].id + " != " + rf[i].id;
          return false;
        }
        if (!check_impl(context, *lf[i].typ, *rf[i].typ, err)) {
          return false;
        }
      }
      return true;
    }

    if (left.data.is(ir::IrTypeData::Tag::tuple_t) && right.data.is(ir::IrTypeData::Tag::tuple_t)) {
      auto le = ir::type_elems(left);
      auto re = ir::type_elems(right);
      if (le.size() != re.size()) {
        err = "expected " + std::to_string(le.size()) + " elements; got " + std::to_string(re.size());
        return false;
      }
      for (size_t i = 0; i < le.size(); ++i) {
        if (!check_impl(context, le[i], re[i], err)) {
          return false;
        }
      }
      return true;
    }

    if (left.data.is(ir::IrTypeData::Tag::variant_t) && right.data.is(ir::IrTypeData::Tag::variant_t)) {
      auto lt = ir::type_tags(left);
      auto rt = ir::type_tags(right);
      if (lt.size() != rt.size()) {
        err = "expected " + std::to_string(lt.size()) + " tags; got " + std::to_string(rt.size());
        return false;
      }
      for (size_t i = 0; i < lt.size(); ++i) {
        if (lt[i].id != rt[i].id) {
          err = "tag name mismatch: " + lt[i].id + " != " + rt[i].id;
          return false;
        }
        if (!check_impl(context, *lt[i].typ, *rt[i].typ, err)) {
          return false;
        }
      }
      return true;
    }

    if (left.data.is(ir::IrTypeData::Tag::var_t) && right.data.is(ir::IrTypeData::Tag::var_t) &&
        left.data.as_var_t() == right.data.as_var_t()) {
      return true;
    }

    if (left.data.is(ir::IrTypeData::Tag::name_t) && context.contains_const_bind(left.data.as_name_t()) &&
        right.data.is(ir::IrTypeData::Tag::name_t) && context.contains_const_bind(right.data.as_name_t()) &&
        left.data.as_name_t() == right.data.as_name_t()) {
      return true;
    }

    if ((left.data.is(ir::IrTypeData::Tag::name_t) && context.contains_alias_bind(left.data.as_name_t())) ||
        (right.data.is(ir::IrTypeData::Tag::name_t) && context.contains_alias_bind(right.data.as_name_t()))) {
      throw std::logic_error("internal error: alias NameType was not expanded by reduce_and_predicate_type");
    }

    err = "expected type " + ir::type_to_string(left) + "; got " + ir::type_to_string(right);
    return false;
  }

 private:
  bool is_assumed(const ir::IrType& left, const ir::IrType& right) const {
    for (const auto& p : assumed_) {
      if (ir::equals_type(p.first, left) && ir::equals_type(p.second, right)) {
        return true;
      }
    }
    return false;
  }

  std::vector<std::pair<ir::IrType, ir::IrType>> assumed_;
};

inline bool subtype(Context& context, const ir::IrType& left, const ir::IrType& right, std::string& err) {
  if (!is_wellformed_type(context, left, err)) {
    return false;
  }
  if (!is_wellformed_type(context, right, err)) {
    return false;
  }
  Subtyper subtyper;
  if (!subtyper.check_impl(context, left, right, err)) {
    err += "\n  subtyping " + ir::type_to_string(left) + " and " + ir::type_to_string(right);
    return false;
  }
  return true;
}

} // namespace ts
