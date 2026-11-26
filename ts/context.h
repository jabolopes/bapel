#pragma once

#include "bin/ir_base.h"
#include "bin/ir_decl.h"
#include "bin/ir_function.h"
#include "bin/ir_type.h"
#include "ts/bind.h"
#include "ts/list.h"

#include <algorithm>
#include <cstdint>
#include <functional>
#include <map>
#include <memory>
#include <optional>
#include <set>
#include <sstream>
#include <stdexcept>
#include <string>
#include <tuple>
#include <utility>
#include <vector>

namespace ts {

inline std::string base_type_name(const ir::IrType& typ) {
  if (typ.data.is(ir::IrTypeData::Tag::name_t)) {
    return typ.data.as_name_t();
  }
  if (typ.data.is(ir::IrTypeData::Tag::app_t) && typ.data.as_app_t().fun) {
    return base_type_name(*typ.data.as_app_t().fun);
  }
  return "";
}

inline bool match_type(
    const ir::IrType& pattern,
    const ir::IrType& target,
    const std::set<std::string>& vars,
    std::map<std::string, ir::IrType>& subs) {
  if (pattern.data.is(ir::IrTypeData::Tag::var_t) && vars.count(pattern.data.as_var_t())) {
    const auto& pv = pattern.data.as_var_t();
    auto it = subs.find(pv);
    if (it != subs.end()) {
      return ir::equals_type(it->second, target);
    }
    subs[pv] = target;
    return true;
  }

  if (pattern.data.tag() != target.data.tag()) {
    return false;
  }

  switch (pattern.data.tag()) {
    case ir::IrTypeData::Tag::name_t:
      return pattern.data.as_name_t() == target.data.as_name_t();

    case ir::IrTypeData::Tag::mu_t: {
      const auto& pm = pattern.data.as_mu_t();
      const auto& tm = target.data.as_mu_t();
      if (!pm.body || !tm.body) {
        return false;
      }
      if (!ir::equals_kind(pm.kind, tm.kind)) {
        return false;
      }
      ir::IrType p_body = ir::substitute_type(
          *pm.body, ir::new_var_type(pm.var), ir::new_var_type("__mu"));
      ir::IrType t_body = ir::substitute_type(
          *tm.body, ir::new_var_type(tm.var), ir::new_var_type("__mu"));
      return match_type(p_body, t_body, vars, subs);
    }

    case ir::IrTypeData::Tag::app_t: {
      const auto& pa = pattern.data.as_app_t();
      const auto& ta = target.data.as_app_t();
      if (!pa.fun || !ta.fun || !pa.arg || !ta.arg) {
        return false;
      }
      if (!match_type(*pa.fun, *ta.fun, vars, subs)) {
        return false;
      }
      return match_type(*pa.arg, *ta.arg, vars, subs);
    }

    case ir::IrTypeData::Tag::tuple_t: {
      auto p_elems = ir::type_elems(pattern);
      auto t_elems = ir::type_elems(target);
      if (p_elems.size() != t_elems.size()) {
        return false;
      }
      for (size_t i = 0; i < p_elems.size(); ++i) {
        if (!match_type(p_elems[i], t_elems[i], vars, subs)) {
          return false;
        }
      }
      return true;
    }

    case ir::IrTypeData::Tag::struct_t: {
      auto p_fields = ir::type_fields(pattern);
      auto t_fields = ir::type_fields(target);
      if (p_fields.size() != t_fields.size()) {
        return false;
      }
      for (size_t i = 0; i < p_fields.size(); ++i) {
        if (p_fields[i].id != t_fields[i].id || !p_fields[i].typ || !t_fields[i].typ) {
          return false;
        }
        if (!match_type(*p_fields[i].typ, *t_fields[i].typ, vars, subs)) {
          return false;
        }
      }
      return true;
    }

    case ir::IrTypeData::Tag::fun_t: {
      const auto& pf = pattern.data.as_fun_t();
      const auto& tf = target.data.as_fun_t();
      if (!pf.arg || !tf.arg || !pf.ret || !tf.ret) {
        return false;
      }
      if (!match_type(*pf.arg, *tf.arg, vars, subs)) {
        return false;
      }
      return match_type(*pf.ret, *tf.ret, vars, subs);
    }

    default:
      return ir::equals_type(pattern, target);
  }
}

class Context {
 public:
  Context() : list_(), wellformed_size_(0) {}
  explicit Context(List<Binding> list, size_t wellformed_size = 0)
      : list_(std::move(list)), wellformed_size_(wellformed_size) {}

  bool empty() const { return list_.empty(); }
  size_t size() const { return list_.size(); }
  const List<Binding>& list() const { return list_; }

  std::string to_string() const {
    std::ostringstream ss;
    if (!list_.empty()) {
      auto binds = list_.collect();
      ss << binds[0].to_string();
      for (size_t i = 1; i < binds.size(); ++i) {
        ss << ", " << binds[i].to_string();
      }
    }
    return ss.str();
  }

  template <typename Predicate>
  const Binding* lookup_bind_ptr(Predicate&& predicate) const {
    return list_.find_if(std::forward<Predicate>(predicate));
  }

  template <typename Predicate>
  bool lookup_bind(Predicate&& predicate, Binding& out_bind) const {
    const Binding* b = lookup_bind_ptr(std::forward<Predicate>(predicate));
    if (b) {
      out_bind = *b;
      return true;
    }
    return false;
  }

  const Binding* lookup_alias_bind_ptr(const std::string& name) const {
    return lookup_bind_ptr([&](const Binding& b) {
      return b.is_alias() && b.alias && b.alias->name == name;
    });
  }

  bool lookup_alias_bind(const std::string& name, Binding& out_bind) const {
    const Binding* b = lookup_alias_bind_ptr(name);
    if (b) {
      out_bind = *b;
      return true;
    }
    return false;
  }

  bool contains_alias_bind(const std::string& name) const {
    return lookup_alias_bind_ptr(name) != nullptr;
  }

  const Binding* lookup_const_bind_ptr(const std::string& name) const {
    return lookup_bind_ptr([&](const Binding& b) {
      return b.is_const() && b.const_data && b.const_data->name == name;
    });
  }

  bool lookup_const_bind(const std::string& name, Binding& out_bind) const {
    const Binding* b = lookup_const_bind_ptr(name);
    if (b) {
      out_bind = *b;
      return true;
    }
    return false;
  }

  bool contains_const_bind(const std::string& name) const {
    return lookup_const_bind_ptr(name) != nullptr;
  }

  const Binding* lookup_term_decl_or_def_bind_ptr(const std::string& name) const {
    return lookup_bind_ptr([&](const Binding& b) {
      return (b.is_term_decl() && b.term_decl && b.term_decl->name == name) ||
             (b.is_term_def() && b.term_def && b.term_def->name == name);
    });
  }

  bool lookup_term_decl_or_def_bind(const std::string& name, Binding& out_bind) const {
    const Binding* b = lookup_term_decl_or_def_bind_ptr(name);
    if (b) {
      out_bind = *b;
      return true;
    }
    return false;
  }

  const Binding* lookup_term_decl_bind_in_scope_ptr(const std::string& name) const {
    const Binding* b = lookup_bind_ptr([&](const Binding& b) {
      return b.is_scope() || (b.is_term_decl() && b.term_decl && b.term_decl->name == name);
    });
    if (!b || b->is_scope()) {
      return nullptr;
    }
    return b;
  }

  bool contains_term_decl_bind_in_scope(const std::string& name) const {
    return lookup_term_decl_bind_in_scope_ptr(name) != nullptr;
  }

  const Binding* lookup_term_def_bind_in_scope_ptr(const std::string& name) const {
    const Binding* b = lookup_bind_ptr([&](const Binding& b) {
      return b.is_scope() || (b.is_term_def() && b.term_def && b.term_def->name == name);
    });
    if (!b || b->is_scope()) {
      return nullptr;
    }
    return b;
  }

  bool contains_term_def_bind_in_scope(const std::string& name) const {
    return lookup_term_def_bind_in_scope_ptr(name) != nullptr;
  }

  const Binding* lookup_scope_bind_ptr() const {
    return lookup_bind_ptr([](const Binding& b) {
      return b.is_scope();
    });
  }

  const Binding* lookup_type_param_bind_ptr(const std::string& tvar) const {
    return lookup_bind_ptr([&](const Binding& b) {
      return b.is_type_param() && b.type_param && b.type_param->name == tvar;
    });
  }

  bool lookup_type_param_bind(const std::string& tvar, Binding& out_bind) const {
    const Binding* b = lookup_type_param_bind_ptr(tvar);
    if (b) {
      out_bind = *b;
      return true;
    }
    return false;
  }

  bool contains_type_param_bind(const std::string& tvar) const {
    return lookup_type_param_bind_ptr(tvar) != nullptr;
  }

  const Binding* lookup_type_param_bind_in_scope_ptr(const std::string& tvar) const {
    const Binding* b = lookup_bind_ptr([&](const Binding& b) {
      return b.is_scope() || (b.is_type_param() && b.type_param && b.type_param->name == tvar);
    });
    if (!b || b->is_scope()) {
      return nullptr;
    }
    return b;
  }

  bool contains_type_param_bind_in_scope(const std::string& tvar) const {
    return lookup_type_param_bind_in_scope_ptr(tvar) != nullptr;
  }

  const Binding* lookup_trait_bind_ptr(const std::string& name) const {
    return lookup_bind_ptr([&](const Binding& b) {
      return b.is_trait() && b.trait && b.trait->name == name;
    });
  }

  bool lookup_trait_bind(const std::string& name, Binding& out_bind) const {
    const Binding* b = lookup_trait_bind_ptr(name);
    if (b) {
      out_bind = *b;
      return true;
    }
    return false;
  }

  bool contains_trait_bind(const std::string& name) const {
    return lookup_trait_bind_ptr(name) != nullptr;
  }

  bool lookup_trait_impl_bind(const ir::IrType& trait_type, const ir::IrType& type_name, Binding& out_bind) const {
    return lookup_bind([&](const Binding& b) {
      if (!b.is_trait_impl() || !b.trait_impl) return false;
      const auto& impl = *b.trait_impl;
      std::set<std::string> vars;
      for (const auto& tp : impl.type_params) {
        vars.insert(tp.var);
      }
      std::map<std::string, ir::IrType> subs;
      if (!match_type(impl.trait_type, trait_type, vars, subs)) {
        return false;
      }
      return match_type(impl.type_name, type_name, vars, subs);
    }, out_bind);
  }

  std::optional<std::string> lookup_mu_alias_name(const ir::IrType& mu_type) const {
    if (!mu_type.data.is(ir::IrTypeData::Tag::mu_t)) {
      return std::nullopt;
    }
    const Binding* alias_b = lookup_bind_ptr([&](const Binding& b) {
      if (!b.is_alias() || !b.alias || !b.alias->type.data.is(ir::IrTypeData::Tag::mu_t)) {
        return false;
      }
      std::set<std::string> empty_vars;
      std::map<std::string, ir::IrType> empty_subs;
      return match_type(b.alias->type, mu_type, empty_vars, empty_subs);
    });
    if (alias_b && alias_b->alias) {
      return alias_b->alias->name;
    }
    return std::nullopt;
  }

  bool lookup_method(
      const ir::IrType& receiver_type,
      const std::string& method_name,
      std::string& out_name,
      ir::IrType& out_type) const {
    auto extract_term_type = [&](const Binding* bind, const std::string& name) -> bool {
      if (!bind) return false;
      out_name = name;
      if (bind->is_term_decl() && bind->term_decl) {
        out_type = bind->term_decl->type;
      } else if (bind->is_term_def() && bind->term_def) {
        out_type = bind->term_def->type;
      }
      return true;
    };

    ir::IrType value_type = receiver_type;
    while (value_type.data.is(ir::IrTypeData::Tag::app_t) && base_type_name(value_type) == "Ptr") {
      if (value_type.data.as_app_t().arg) {
        value_type = *value_type.data.as_app_t().arg;
      } else {
        break;
      }
    }

    if (auto mu_alias = lookup_mu_alias_name(value_type)) {
      value_type = ir::new_name_type(*mu_alias);
    }

    // 1. Inherent methods on value_type: Type::method
    std::string base_name = base_type_name(value_type);
    if (!base_name.empty()) {
      std::string inherent_name = base_name + "::" + method_name;
      if (extract_term_type(lookup_term_decl_or_def_bind_ptr(inherent_name), inherent_name)) {
        return true;
      }
    }

    // 2. Trait methods on bounded type variables: 't: Size -> Size::method
    if (value_type.data.is(ir::IrTypeData::Tag::var_t)) {
      const Binding* bind = lookup_type_param_bind_ptr(value_type.data.as_var_t());
      if (bind && bind->type_param) {
        for (const auto& bound : bind->type_param->bounds) {
          std::string trait_name = base_type_name(bound);
          if (trait_name.empty()) continue;
          std::string trait_method_name = trait_name + "::" + method_name;
          if (extract_term_type(lookup_term_decl_or_def_bind_ptr(trait_method_name), trait_method_name)) {
            return true;
          }
        }
      }
    }

    // 3. Trait methods across all in-scope trait implementations for value_type
    const Binding* matched_method_binding = nullptr;
    std::string matched_method_name;

    const Binding* matched_impl = list_.find_if([&](const Binding& bind) {
      if (!bind.is_trait_impl() || !bind.trait_impl) return false;
      const auto& impl = *bind.trait_impl;
      std::set<std::string> vars;
      for (const auto& tp : impl.type_params) {
        vars.insert(tp.var);
      }
      std::map<std::string, ir::IrType> subs;
      if (!match_type(impl.type_name, value_type, vars, subs)) {
        return false;
      }

      std::string trait_name = base_type_name(impl.trait_type);
      if (trait_name.empty()) return false;

      std::string trait_method_name = trait_name + "::" + method_name;
      const Binding* m_bind = lookup_term_decl_or_def_bind_ptr(trait_method_name);
      if (m_bind) {
        matched_method_binding = m_bind;
        matched_method_name = std::move(trait_method_name);
        return true;
      }
      return false;
    });

    if (matched_impl && extract_term_type(matched_method_binding, matched_method_name)) {
      return true;
    }

    return false;
  }

  Context add_tvar(std::string name, ir::IrKind kind = ir::new_type_kind(), std::vector<ir::IrType> bounds = {}) const {
    return add_bind(new_tvar_bind(std::move(name), std::move(kind), std::move(bounds)));
  }

  // Scope and symbol management
  Context enter_scope() const {
    const Binding* b = lookup_scope_bind_ptr();
    if (b && b->scope) {
      return add_bind(new_scope_bind(b->scope->level + 1));
    }
    return add_bind(new_scope_bind(1));
  }

  Context enter_function(
      const std::vector<ir::TypeParam>& type_params,
      const std::vector<ir::FunctionArg>& args) const {
    Context c = enter_scope();
    for (const auto& tp : type_params) {
      c = c.add_bind(new_type_param_bind(tp));
    }
    for (const auto& arg : args) {
      c = c.add_bind(new_term_def_bind(arg.id, arg.typ));
    }
    return c;
  }

  Context add_symbol(const ir::IrDecl& decl) const {
    switch (decl.data.tag()) {
      case ir::IrDeclData::Tag::term_decl: {
        const auto& td = decl.data.as_term_decl();
        return add_bind(new_term_decl_bind(td.id, td.typ));
      }

      case ir::IrDeclData::Tag::alias_decl: {
        const auto& ad = decl.data.as_alias_decl();
        return add_bind(new_alias_bind(ad.id, ad.typ));
      }

      case ir::IrDeclData::Tag::name_decl: {
        const auto& nd = decl.data.as_name_decl();
        return add_bind(new_const_bind(nd.id, nd.kind));
      }

      case ir::IrDeclData::Tag::trait_decl: {
        const auto& td = decl.data.as_trait_decl();
        Context c = add_bind(new_trait_bind(td.id, td.type_params, td.methods));
        for (const auto& m : td.methods) {
          std::vector<ir::IrType> args;
          args.reserve(m.args.size());
          for (const auto& arg : m.args) {
            args.push_back(arg.typ);
          }
          ir::IrType ret = m.ret_type;
          ir::IrType method_type = ir::new_function_type(ir::new_tuple_type(std::move(args)), std::move(ret));
          for (auto it = td.type_params.rbegin(); it != td.type_params.rend(); ++it) {
            method_type = ir::new_forall_type(*it, std::move(method_type));
          }
          method_type = ir::new_forall_type(ir::TypeParam{"Self", ir::new_type_kind(), {}}, std::move(method_type));
          c = c.add_bind(new_term_decl_bind(td.id + "::" + m.id, std::move(method_type)));
        }
        return c;
      }
      case ir::IrDeclData::Tag::impl_decl: {
        const auto& impl = decl.data.as_impl_decl();
        if (impl.case_val.is(ir::ImplCase::Tag::inherent_impl)) {
          std::string base_name = base_type_name(impl.type_name);
          if (base_name.empty()) {
            throw std::runtime_error("invalid type for inherent impl: " + ir::type_to_string(impl.type_name));
          }
          Context ctx = *this;
          for (const auto& m : impl.methods) {
            std::string name = base_name + "::" + m.id;
            std::vector<ir::IrType> args;
            for (const auto& a : m.args) {
              args.push_back(a.typ);
            }
            ir::IrType ret = m.ret_type;
            ir::IrType method_type = ir::new_function_type(ir::new_tuple_type(std::move(args)), std::move(ret));
            std::vector<ir::TypeParam> tvars = impl.type_params;
            tvars.insert(tvars.end(), m.type_params.begin(), m.type_params.end());
            for (auto it = tvars.rbegin(); it != tvars.rend(); ++it) {
              method_type = ir::new_forall_type(*it, std::move(method_type));
            }
            ctx = ctx.add_bind(new_term_decl_bind(name, std::move(method_type)));
          }
          return ctx;
        }
        return add_bind(new_trait_impl_bind(impl.type_params, impl.trait_type, impl.type_name));
      }
    }
    return *this;
  }

  Context add_bind(Binding bind) const {
    bool fast_wf = bind.is_trait_impl();
    List<Binding> new_list = list_.add(std::move(bind));
    if (fast_wf && wellformed_size_ == list_.size()) {
      return Context(std::move(new_list), new_list.size());
    }
    Context new_ctx(std::move(new_list), wellformed_size_);
    std::string err;
    if (!new_ctx.is_wellformed(err)) {
      throw std::runtime_error("Context is not wellformed: " + err);
    }
    new_ctx.wellformed_size_ = new_ctx.list_.size();
    return new_ctx;
  }

  static int64_t& global_tvar_gen() {
    static int64_t tvar_gen = 0;
    return tvar_gen;
  }

  static int64_t& global_evar_gen() {
    static int64_t evar_gen = 0;
    return evar_gen;
  }

  static void reset_fresh_var_generators() {
    global_tvar_gen() = 0;
    global_evar_gen() = 0;
  }

  // Fresh variable generation
  ir::IrType gen_fresh_var_type() const {
    std::vector<bool> short_name_used(26, false);
    list_.for_each([&](const Binding& b) {
      if (b.is_type_param() && b.type_param) {
        if (b.type_param->name.size() == 1) {
          char c = b.type_param->name[0];
          if (c >= 'a' && c <= 'z') {
            short_name_used[c - 'a'] = true;
          }
        }
      }
    });
    for (size_t i = 0; i < 26; ++i) {
      if (!short_name_used[i]) {
        return ir::new_var_type(std::string(1, static_cast<char>('a' + i)));
      }
    }

    return ir::new_var_type("ꞇ" + std::to_string(global_tvar_gen()++));
  }

  ir::IrType gen_fresh_exist_var() const {
    return ir::new_exist_var_type(global_evar_gen()++);
  }

  std::tuple<Context, ir::TypeParam, ir::IrType> add_fresh_type(const ir::IrType& typ) const {
    if (!typ.data.is(ir::IrTypeData::Tag::forall_t) && !typ.data.is(ir::IrTypeData::Tag::lambda_t) && !typ.data.is(ir::IrTypeData::Tag::mu_t)) {
      return {*this, ir::TypeParam{}, typ};
    }

    if (typ.data.is(ir::IrTypeData::Tag::forall_t)) {
      const auto& fa = typ.data.as_forall_t();
      ir::IrType fresh_var = gen_fresh_var_type();
      std::string fresh_name = fresh_var.data.as_var_t();
      std::vector<ir::IrType> bounds;
      for (const auto& b : fa.param.bounds) {
        bounds.push_back(ir::substitute_type(b, ir::new_var_type(fa.param.var), fresh_var));
      }
      ir::TypeParam tp{fresh_name, fa.param.kind, std::move(bounds)};
      ir::IrType body = fa.body ? ir::substitute_type(*fa.body, ir::new_var_type(fa.param.var), fresh_var) : ir::IrType{};
      Context new_ctx = add_bind(new_type_param_bind(tp));
      return {new_ctx, tp, body};
    }

    if (typ.data.is(ir::IrTypeData::Tag::lambda_t)) {
      const auto& l = typ.data.as_lambda_t();
      ir::IrType fresh_var = gen_fresh_var_type();
      std::string fresh_name = fresh_var.data.as_var_t();
      ir::TypeParam tp{fresh_name, l.kind, {}};
      ir::IrType body = l.body ? ir::substitute_type(*l.body, ir::new_var_type(l.var), fresh_var) : ir::IrType{};
      Context new_ctx = add_bind(new_type_param_bind(tp));
      return {new_ctx, tp, body};
    }

    if (typ.data.is(ir::IrTypeData::Tag::mu_t)) {
      const auto& m = typ.data.as_mu_t();
      ir::IrType fresh_var = gen_fresh_var_type();
      std::string fresh_name = fresh_var.data.as_var_t();
      ir::TypeParam tp{fresh_name, m.kind, {}};
      ir::IrType body = m.body ? ir::substitute_type(*m.body, ir::new_var_type(m.var), fresh_var) : ir::IrType{};
      Context new_ctx = add_bind(new_type_param_bind(tp));
      return {new_ctx, tp, body};
    }

    return {*this, ir::TypeParam{}, typ};
  }

  // Well-formedness checks
  bool is_wellformed(std::string& err) const {
    if (empty()) return true;
    if (wellformed_size_ == list_.size()) return true;

    const Binding& bind = list_.front();
    Context rest(list_.remove(), wellformed_size_);
    if (!rest.is_wellformed(err)) {
      return false;
    }

    switch (bind.case_val) {
      case BindCase::AliasBind:
        if (bind.alias) {
          if (rest.contains_alias_bind(bind.alias->name) ||
              rest.contains_const_bind(bind.alias->name) ||
              rest.contains_trait_bind(bind.alias->name)) {
            err = "type/trait \"" + bind.alias->name + "\" is already defined";
            return false;
          }
        }
        break;

      case BindCase::ConstBind:
        if (bind.const_data) {
          if (rest.contains_alias_bind(bind.const_data->name) ||
              rest.contains_const_bind(bind.const_data->name) ||
              rest.contains_trait_bind(bind.const_data->name)) {
            err = "type/trait \"" + bind.const_data->name + "\" is already defined";
            return false;
          }
        }
        break;

      case BindCase::ScopeBind:
        if (bind.scope) {
          const Binding* prev_scope = rest.lookup_scope_bind_ptr();
          int want_level = 1;
          if (prev_scope && prev_scope->scope) {
            want_level = prev_scope->scope->level + 1;
          }
          if (bind.scope->level != want_level) {
            err = "expected scope " + std::to_string(want_level) + "; got " + bind.scope->to_string();
            return false;
          }
        }
        break;

      case BindCase::TermDeclBind:
        if (bind.term_decl) {
          if (rest.contains_term_decl_bind_in_scope(bind.term_decl->name)) {
            err = "term \"" + bind.term_decl->name + "\" is already declared";
            return false;
          }
          if (rest.contains_term_def_bind_in_scope(bind.term_decl->name)) {
            err = "term \"" + bind.term_decl->name + "\" is already defined";
            return false;
          }
        }
        break;

      case BindCase::TermDefBind:
        if (bind.term_def) {
          if (rest.contains_term_def_bind_in_scope(bind.term_def->name)) {
            err = "term \"" + bind.term_def->name + "\" is already defined";
            return false;
          }
        }
        break;

      case BindCase::TypeParamBind:
        if (bind.type_param) {
          if (rest.contains_type_param_bind_in_scope(bind.type_param->name)) {
            err = "type parameter \"" + bind.type_param->name + "\" is already declared";
            return false;
          }
        }
        break;

      case BindCase::TraitBind:
        if (bind.trait) {
          if (rest.contains_alias_bind(bind.trait->name) ||
              rest.contains_const_bind(bind.trait->name) ||
              rest.contains_trait_bind(bind.trait->name)) {
            err = "type/trait \"" + bind.trait->name + "\" is already defined";
            return false;
          }
        }
        break;

      case BindCase::TraitImplBind:
        break;
    }

    return true;
  }

 private:
  List<Binding> list_;
  size_t wellformed_size_ = 0;
};

inline void reset_fresh_var_generators() {
  Context::reset_fresh_var_generators();
}

} // namespace ts
