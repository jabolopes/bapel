#pragma once

#include "ir_type.h"
#include <set>
#include <string>
#include <vector>

namespace ir {

inline bool has_free_var(const std::vector<TypeParam>& free_vars, const std::string& name) {
  for (const auto& fv : free_vars) {
    if (fv.var == name) return true;
  }
  return false;
}

inline std::string pick_fresh_binder_var(
    const IrType& t,
    const std::string& binder_var,
    const IrType& source,
    const std::vector<TypeParam>& target_free) {
  std::vector<std::string> bound_v, all_free_v;
  collect_free_vars(t, bound_v, all_free_v);
  for (const auto& tf : target_free) all_free_v.push_back(tf.var);
  if (source.is(IrTypeData::Tag::var_t)) all_free_v.push_back(source.data.as_var_t());
  std::set<std::string> used(all_free_v.begin(), all_free_v.end());
  used.insert(binder_var);

  for (char c = 'a'; c <= 'z'; ++c) {
    std::string s(1, c);
    if (used.find(s) == used.end()) return s;
  }
  int64_t idx = 0;
  while (true) {
    std::string s = "t" + std::to_string(idx++);
    if (used.find(s) == used.end()) return s;
  }
}

inline IrType substitute_type(const IrType& t, const IrType& source, const IrType& target) {
  if (equals_type(t, source)) {
    return target;
  }

  if (t.is(IrTypeData::Tag::forall_t)) {
    const auto& fa = t.data.as_forall_t();
    const std::string& var = fa.type_param.var;
    if (source.is(IrTypeData::Tag::var_t) && source.data.as_var_t() == var) {
      std::vector<IrType> bounds;
      bounds.reserve(fa.type_param.bounds.size());
      for (const auto& b : fa.type_param.bounds) {
        bounds.push_back(substitute_type(b, source, target));
      }
      TypeParam tp{var, fa.type_param.kind, std::move(bounds)};
      return new_forall_type(std::move(tp), fa.type ? *fa.type : IrType{});
    }
    std::vector<TypeParam> target_free = get_free_type_vars(target);
    if (has_free_var(target_free, var)) {
      std::string fresh = pick_fresh_binder_var(t, var, source, target_free);
      std::vector<IrType> final_bounds;
      final_bounds.reserve(fa.type_param.bounds.size());
      for (const auto& b : fa.type_param.bounds) {
        IrType rb = substitute_type(b, new_var_type(var), new_var_type(fresh));
        final_bounds.push_back(substitute_type(rb, source, target));
      }
      IrType renamed_body = fa.type ? substitute_type(*fa.type, new_var_type(var), new_var_type(fresh)) : IrType{};
      TypeParam tp{fresh, fa.type_param.kind, std::move(final_bounds)};
      return new_forall_type(std::move(tp), substitute_type(renamed_body, source, target));
    }
  } else if (t.is(IrTypeData::Tag::lambda_t)) {
    const auto& l = t.data.as_lambda_t();
    const std::string& var = l.var;
    if (source.is(IrTypeData::Tag::var_t) && source.data.as_var_t() == var) return t;
    std::vector<TypeParam> target_free = get_free_type_vars(target);
    if (has_free_var(target_free, var)) {
      std::string fresh = pick_fresh_binder_var(t, var, source, target_free);
      IrType renamed_body = l.type ? substitute_type(*l.type, new_var_type(var), new_var_type(fresh)) : IrType{};
      return new_lambda_type(fresh, l.kind, substitute_type(renamed_body, source, target));
    }
  } else if (t.is(IrTypeData::Tag::mu_t)) {
    const auto& m = t.data.as_mu_t();
    const std::string& var = m.var;
    if (source.is(IrTypeData::Tag::var_t) && source.data.as_var_t() == var) return t;
    std::vector<TypeParam> target_free = get_free_type_vars(target);
    if (has_free_var(target_free, var)) {
      std::string fresh = pick_fresh_binder_var(t, var, source, target_free);
      IrType renamed_body = m.body ? substitute_type(*m.body, new_var_type(var), new_var_type(fresh)) : IrType{};
      return new_mu_type(fresh, m.kind, substitute_type(renamed_body, source, target));
    }
  }

  return map_type(t, [&](const IrType& child) {
    return substitute_type(child, source, target);
  });
}

} // namespace ir
