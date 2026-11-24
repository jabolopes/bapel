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
#include <utility>
#include <vector>

namespace ts {

inline bool same_tags(const ir::IrType& left, const ir::IrType& right) {
  auto lt = ir::type_tags(left);
  auto rt = ir::type_tags(right);
  if (lt.size() != rt.size()) return false;
  for (size_t i = 0; i < lt.size(); ++i) {
    if (lt[i].id != rt[i].id) return false;
  }
  return true;
}

inline bool same_fields(const ir::IrType& left, const ir::IrType& right) {
  auto lf = ir::type_fields(left);
  auto rf = ir::type_fields(right);
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
    exist_vars[evar.data.as_exist_var_t()] = ExistVarData{nullptr};
    return evar;
  }

  bool is_exist_var_unassigned(const ir::IrType& tvar) const {
    if (!tvar.data.is(ir::IrTypeData::Tag::exist_var_t)) return false;
    auto it = exist_vars.find(tvar.data.as_exist_var_t());
    return it != exist_vars.end() && it->second.solution == nullptr;
  }

  bool is_exist_var_assigned(const ir::IrType& tvar) const {
    if (!tvar.data.is(ir::IrTypeData::Tag::exist_var_t)) return false;
    auto it = exist_vars.find(tvar.data.as_exist_var_t());
    return it != exist_vars.end() && it->second.solution != nullptr;
  }

  ir::IrType exist_var_solution(const ir::IrType& evar) {
    if (!evar.data.is(ir::IrTypeData::Tag::exist_var_t)) {
      throw std::runtime_error("expected existential variable; got " + ir::type_to_string(evar));
    }
    auto& data = exist_vars[evar.data.as_exist_var_t()];
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
    if (typ.data.is(ir::IrTypeData::Tag::forall_t)) {
      return false;
    }
    if (typ.data.is(ir::IrTypeData::Tag::exist_var_t)) {
      return typ.data.as_exist_var_t() < evar.data.as_exist_var_t();
    }
    return true;
  }

  void solve_exist_var(const ir::IrType& evar, const ir::IrType& typ) {
    if (!evar.data.is(ir::IrTypeData::Tag::exist_var_t)) {
      throw std::runtime_error("expected existential variable; got " + ir::type_to_string(evar));
    }
    auto& data = exist_vars[evar.data.as_exist_var_t()];
    if (data.solution != nullptr) {
      throw std::runtime_error(ir::type_to_string(evar) + " is already solved with solution " + ir::type_to_string(*data.solution));
    }
    data.solution = std::make_shared<ir::IrType>(typ);
  }

  ir::IrType solve_type(const ir::IrType& typ) {
    if (typ.data.is(ir::IrTypeData::Tag::exist_var_t)) {
      if (is_exist_var_assigned(typ)) {
        return solve_type(exist_var_solution(typ));
      }
      return typ;
    }
    if (typ.data.is(ir::IrTypeData::Tag::forall_t)) {
      const auto& fa = typ.data.as_forall_t();
      ir::IrType body = fa.body ? solve_type(*fa.body) : ir::IrType{};
      return ir::new_forall_type(fa.param, std::move(body));
    }
    return ir::map_type(typ, [&](const ir::IrType& child) {
      return solve_type(child);
    });
  }
};

class Unifier {
 public:
  bool unify(Context& context, UnificationState& state, ir::IrType left, ir::IrType right, std::string& err) {
    if (state.is_exist_var_assigned(left)) {
      left = state.exist_var_solution(left);
      return unify(context, state, left, right, err);
    }

    if (state.is_exist_var_assigned(right)) {
      right = state.exist_var_solution(right);
      return unify(context, state, left, right, err);
    }

    if (left.data.is(ir::IrTypeData::Tag::exist_var_t) && right.data.is(ir::IrTypeData::Tag::exist_var_t) &&
        left.data.as_exist_var_t() == right.data.as_exist_var_t()) {
      return true;
    }

    if (state.is_exist_var_unassigned(left) && state.can_assign(left, right)) {
      state.solve_exist_var(left, right);
      return true;
    }

    if (state.is_exist_var_unassigned(right) && state.can_assign(right, left)) {
      state.solve_exist_var(right, left);
      return true;
    }

    if (ir::equals_type(left, right)) {
      return true;
    }

    if (is_assumed(left, right)) {
      return true;
    }

    if (left.data.is(ir::IrTypeData::Tag::app_t) && right.data.is(ir::IrTypeData::Tag::app_t)) {
      const auto& la = left.data.as_app_t();
      const auto& ra = right.data.as_app_t();
      if (!unify(context, state, *la.fun, *ra.fun, err)) {
        err = "mismatch in function types: " + err;
        return false;
      }
      if (!unify(context, state, *la.arg, *ra.arg, err)) {
        err = "mismatch in argument types: " + err;
        return false;
      }
      return true;
    }

    if (left.data.is(ir::IrTypeData::Tag::array_t) && right.data.is(ir::IrTypeData::Tag::array_t)) {
      const auto& la = left.data.as_array_t();
      const auto& ra = right.data.as_array_t();
      if (!unify(context, state, *la.elem, *ra.elem, err)) {
        err = "mismatch in array element types: " + err;
        return false;
      }
      if (la.size != ra.size) {
        err = "expected array with " + std::to_string(la.size) + " elements; got " + std::to_string(ra.size) + " elements";
        return false;
      }
      return true;
    }

    if (right.data.is(ir::IrTypeData::Tag::forall_t)) {
      auto [new_ctx, tp, body] = context.add_fresh_type(right);
      context = new_ctx;
      return unify(context, state, left, body, err);
    }

    if (left.data.is(ir::IrTypeData::Tag::forall_t)) {
      const auto& lf = left.data.as_forall_t();
      ir::IrType evar = state.new_evar(context);
      ir::IrType new_body = ir::substitute_type(*lf.body, ir::new_var_type(lf.param.var), evar);
      return unify(context, state, new_body, right, err);
    }

    if (left.data.is(ir::IrTypeData::Tag::lambda_t) && right.data.is(ir::IrTypeData::Tag::lambda_t)) {
      auto [new_ctx, tp, body] = context.add_fresh_type(left);
      const auto& rl = right.data.as_lambda_t();
      ir::IrType right_body = ir::substitute_type(*rl.body, ir::new_var_type(rl.var), ir::new_var_type(tp.var));
      return unify(new_ctx, state, body, right_body, err);
    }

    if (left.data.is(ir::IrTypeData::Tag::mu_t) && right.data.is(ir::IrTypeData::Tag::mu_t)) {
      assumed_.push_back({left, right});
      auto [new_ctx, tp, body] = context.add_fresh_type(left);
      const auto& rm = right.data.as_mu_t();
      ir::IrType right_body = ir::substitute_type(*rm.body, ir::new_var_type(rm.var), ir::new_var_type(tp.var));
      bool ok = unify(new_ctx, state, body, right_body, err);
      assumed_.pop_back();
      return ok;
    }

    if (left.data.is(ir::IrTypeData::Tag::mu_t)) {
      assumed_.push_back({left, right});
      bool ok = unify(context, state, unfold_mu_type(left), right, err);
      assumed_.pop_back();
      return ok;
    }

    if (right.data.is(ir::IrTypeData::Tag::mu_t)) {
      assumed_.push_back({left, right});
      bool ok = unify(context, state, left, unfold_mu_type(right), err);
      assumed_.pop_back();
      return ok;
    }

    if (left.data.is(ir::IrTypeData::Tag::fun_t) && right.data.is(ir::IrTypeData::Tag::fun_t)) {
      const auto& lf = left.data.as_fun_t();
      const auto& rf = right.data.as_fun_t();
      if (!unify(context, state, *lf.arg, *rf.arg, err)) {
        return false;
      }
      return unify(context, state, *lf.ret, *rf.ret, err);
    }

    if (left.data.is(ir::IrTypeData::Tag::name_t) && context.contains_const_bind(left.data.as_name_t()) &&
        right.data.is(ir::IrTypeData::Tag::name_t) && context.contains_const_bind(right.data.as_name_t()) &&
        left.data.as_name_t() == right.data.as_name_t()) {
      return true;
    }

    if (left.data.is(ir::IrTypeData::Tag::name_t) && context.contains_alias_bind(left.data.as_name_t())) {
      Binding bind;
      if (context.lookup_alias_bind(left.data.as_name_t(), bind) && bind.alias) {
        return unify(context, state, bind.alias->type, right, err);
      }
    }

    if (right.data.is(ir::IrTypeData::Tag::name_t) && context.contains_alias_bind(right.data.as_name_t())) {
      Binding bind;
      if (context.lookup_alias_bind(right.data.as_name_t(), bind) && bind.alias) {
        return unify(context, state, left, bind.alias->type, err);
      }
    }

    if (left.data.is(ir::IrTypeData::Tag::struct_t) && right.data.is(ir::IrTypeData::Tag::struct_t) && same_fields(left, right)) {
      auto lf = ir::type_fields(left);
      auto rf = ir::type_fields(right);
      for (size_t i = 0; i < lf.size(); ++i) {
        if (!unify(context, state, *lf[i].typ, *rf[i].typ, err)) {
          return false;
        }
      }
      return true;
    }

    if (left.data.is(ir::IrTypeData::Tag::tuple_t) && right.data.is(ir::IrTypeData::Tag::tuple_t) && ir::type_elems(left).size() == ir::type_elems(right).size()) {
      auto le = ir::type_elems(left);
      auto re = ir::type_elems(right);
      for (size_t i = 0; i < le.size(); ++i) {
        if (!unify(context, state, le[i], re[i], err)) {
          return false;
        }
      }
      return true;
    }

    if (left.data.is(ir::IrTypeData::Tag::variant_t) && right.data.is(ir::IrTypeData::Tag::variant_t) && same_tags(left, right)) {
      auto lt = ir::type_tags(left);
      auto rt = ir::type_tags(right);
      for (size_t i = 0; i < lt.size(); ++i) {
        if (!unify(context, state, *lt[i].typ, *rt[i].typ, err)) {
          return false;
        }
      }
      return true;
    }

    if (left.data.is(ir::IrTypeData::Tag::var_t) && right.data.is(ir::IrTypeData::Tag::var_t) &&
        left.data.as_var_t() == right.data.as_var_t()) {
      return true;
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

inline bool unify(Context& context, UnificationState& state, const ir::IrType& left, const ir::IrType& right, std::string& err) {
  if (!is_wellformed_type(context, left, err)) {
    return false;
  }
  if (!is_wellformed_type(context, right, err)) {
    return false;
  }
  Unifier unifier;
  if (!unifier.unify(context, state, left, right, err)) {
    err += "\n  unifying " + ir::type_to_string(left) + " and " + ir::type_to_string(right);
    return false;
  }
  return true;
}

} // namespace ts
