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
  if (typ.is(ir::IrTypeCase::NameType)) {
    return typ.name;
  }
  if (typ.is(ir::IrTypeCase::AppType) && typ.app && typ.app->fun) {
    return base_type_name(*typ.app->fun);
  }
  return "";
}

inline bool match_type(
    const ir::IrType& pattern,
    const ir::IrType& target,
    const std::set<std::string>& vars,
    std::map<std::string, ir::IrType>& subs) {
  if (pattern.is(ir::IrTypeCase::VarType) && vars.count(pattern.var)) {
    auto it = subs.find(pattern.var);
    if (it != subs.end()) {
      return ir::equals_type(it->second, target);
    }
    subs[pattern.var] = target;
    return true;
  }

  if (pattern.case_val != target.case_val) {
    return false;
  }

  switch (pattern.case_val) {
    case ir::IrTypeCase::NameType:
      return pattern.name == target.name;

    case ir::IrTypeCase::AppType: {
      if (!pattern.app || !target.app || !pattern.app->fun || !target.app->fun ||
          !pattern.app->arg || !target.app->arg) {
        return false;
      }
      if (!match_type(*pattern.app->fun, *target.app->fun, vars, subs)) {
        return false;
      }
      return match_type(*pattern.app->arg, *target.app->arg, vars, subs);
    }

    case ir::IrTypeCase::TupleType: {
      auto p_elems = pattern.elems();
      auto t_elems = target.elems();
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

    case ir::IrTypeCase::FunType: {
      if (!pattern.fun || !target.fun || !pattern.fun->arg || !target.fun->arg ||
          !pattern.fun->ret || !target.fun->ret) {
        return false;
      }
      if (!match_type(*pattern.fun->arg, *target.fun->arg, vars, subs)) {
        return false;
      }
      return match_type(*pattern.fun->ret, *target.fun->ret, vars, subs);
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

  std::pair<Binding, Context> pop() const {
    if (list_.empty()) {
      throw std::runtime_error("Context is empty");
    }
    Binding b = list_.front();
    List<Binding> rem = list_.remove();
    size_t wf_size = std::min(wellformed_size_, rem.size());
    return {std::move(b), Context(std::move(rem), wf_size)};
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

  Binding get_alias_bind(const std::string& name) const {
    Binding b;
    if (!lookup_alias_bind(name, b)) {
      throw std::runtime_error("type \"" + name + "\" is undefined");
    }
    return b;
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

  Binding get_const_bind(const std::string& name) const {
    Binding b;
    if (!lookup_const_bind(name, b)) {
      throw std::runtime_error("type \"" + name + "\" is undefined");
    }
    return b;
  }

  const Binding* lookup_term_decl_or_def_bind_ptr(const std::string& name) const {
    return lookup_bind_ptr([&](const Binding& b) {
      return (b.is_term_decl() && b.term_decl && b.term_decl->name == name) ||
             (b.is_term_def() && b.term_def && b.term_def->name == name) ||
             (b.is_term_var() && b.term_var && b.term_var->name == name);
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

  Binding get_term_decl_or_def_bind(const std::string& name) const {
    Binding b;
    if (!lookup_term_decl_or_def_bind(name, b)) {
      throw std::runtime_error("term \"" + name + "\" is undefined");
    }
    return b;
  }

  bool lookup_term_var(const std::string& name, ir::IrType& out_type) const {
    const Binding* b = lookup_term_decl_or_def_bind_ptr(name);
    if (b) {
      if (b->is_term_decl() && b->term_decl) {
        out_type = b->term_decl->type;
        return true;
      }
      if (b->is_term_def() && b->term_def) {
        out_type = b->term_def->type;
        return true;
      }
      if (b->is_term_var() && b->term_var) {
        out_type = b->term_var->type;
        return true;
      }
    }
    return false;
  }

  bool lookup_decl(const std::string& id, ir::IrDecl& out_decl) const {
    const Binding* b = lookup_bind_ptr([&](const Binding& b) {
      return b.is_decl() && b.decl && b.decl->decl.id() == id;
    });
    if (b) {
      out_decl = b->decl->decl;
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

  bool lookup_term_decl_bind_in_scope(const std::string& name, Binding& out_bind) const {
    const Binding* b = lookup_term_decl_bind_in_scope_ptr(name);
    if (b) {
      out_bind = *b;
      return true;
    }
    return false;
  }

  bool contains_term_decl_bind_in_scope(const std::string& name) const {
    return lookup_term_decl_bind_in_scope_ptr(name) != nullptr;
  }

  const Binding* lookup_term_def_bind_in_scope_ptr(const std::string& name) const {
    const Binding* b = lookup_bind_ptr([&](const Binding& b) {
      return b.is_scope() || (b.is_term_def() && b.term_def && b.term_def->name == name) ||
             (b.is_term_var() && b.term_var && b.term_var->name == name && b.term_var->is_def);
    });
    if (!b || b->is_scope()) {
      return nullptr;
    }
    return b;
  }

  bool lookup_term_def_bind_in_scope(const std::string& name, Binding& out_bind) const {
    const Binding* b = lookup_term_def_bind_in_scope_ptr(name);
    if (b) {
      out_bind = *b;
      return true;
    }
    return false;
  }

  bool contains_term_def_bind_in_scope(const std::string& name) const {
    return lookup_term_def_bind_in_scope_ptr(name) != nullptr;
  }

  const Binding* lookup_scope_bind_ptr() const {
    return lookup_bind_ptr([](const Binding& b) {
      return b.is_scope();
    });
  }

  bool lookup_scope_bind(Binding& out_bind) const {
    const Binding* b = lookup_scope_bind_ptr();
    if (b) {
      out_bind = *b;
      return true;
    }
    return false;
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

  Binding get_type_param_bind(const std::string& tvar) const {
    Binding b;
    if (!lookup_type_param_bind(tvar, b)) {
      throw std::runtime_error("type parameter \"" + tvar + "\" is undefined");
    }
    return b;
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

  bool lookup_type_param_bind_in_scope(const std::string& tvar, Binding& out_bind) const {
    const Binding* b = lookup_type_param_bind_in_scope_ptr(tvar);
    if (b) {
      out_bind = *b;
      return true;
    }
    return false;
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

  Binding get_trait_bind(const std::string& name) const {
    Binding b;
    if (!lookup_trait_bind(name, b)) {
      throw std::runtime_error("trait \"" + name + "\" is undefined");
    }
    return b;
  }

  bool lookup_type_or_trait_bind(const std::string& name, Binding& out_bind) const {
    return lookup_bind([&](const Binding& b) {
      return (b.is_alias() && b.alias && b.alias->name == name) ||
             (b.is_const() && b.const_data && b.const_data->name == name) ||
             (b.is_trait() && b.trait && b.trait->name == name);
    }, out_bind);
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

  bool contains_trait_impl(const ir::IrType& trait_type, const ir::IrType& type_name) const {
    Binding b;
    return lookup_trait_impl_bind(trait_type, type_name, b);
  }

  bool lookup_method(
      const ir::IrType& receiver_type,
      const std::string& method_name,
      std::string& out_name,
      ir::IrType& out_type) const {
    ir::IrType value_type = receiver_type;
    while (value_type.is(ir::IrTypeCase::AppType) && base_type_name(value_type) == "Ptr") {
      if (value_type.app && value_type.app->arg) {
        value_type = *value_type.app->arg;
      } else {
        break;
      }
    }

    // 1. Inherent methods on value_type: Type::method
    std::string base_name = base_type_name(value_type);
    if (!base_name.empty()) {
      std::string inherent_name = base_name + "::" + method_name;
      const Binding* bind = lookup_term_decl_or_def_bind_ptr(inherent_name);
      if (bind) {
        out_name = inherent_name;
        if (bind->is_term_decl() && bind->term_decl) {
          out_type = bind->term_decl->type;
        } else if (bind->is_term_def() && bind->term_def) {
          out_type = bind->term_def->type;
        } else if (bind->is_term_var() && bind->term_var) {
          out_type = bind->term_var->type;
        }
        return true;
      }
    }

    // 2. Trait methods on bounded type variables: 't: Size -> Size::method
    if (value_type.is(ir::IrTypeCase::VarType)) {
      const Binding* bind = lookup_type_param_bind_ptr(value_type.var);
      if (bind && bind->type_param) {
        for (const auto& bound : bind->type_param->bounds) {
          std::string trait_name = base_type_name(bound);
          if (trait_name.empty()) continue;
          std::string trait_method_name = trait_name + "::" + method_name;
          const Binding* m_bind = lookup_term_decl_or_def_bind_ptr(trait_method_name);
          if (m_bind) {
            out_name = trait_method_name;
            if (m_bind->is_term_decl() && m_bind->term_decl) {
              out_type = m_bind->term_decl->type;
            } else if (m_bind->is_term_def() && m_bind->term_def) {
              out_type = m_bind->term_def->type;
            } else if (m_bind->is_term_var() && m_bind->term_var) {
              out_type = m_bind->term_var->type;
            }
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

    if (matched_impl && matched_method_binding) {
      out_name = std::move(matched_method_name);
      if (matched_method_binding->is_term_decl() && matched_method_binding->term_decl) {
        out_type = matched_method_binding->term_decl->type;
      } else if (matched_method_binding->is_term_def() && matched_method_binding->term_def) {
        out_type = matched_method_binding->term_def->type;
      } else if (matched_method_binding->is_term_var() && matched_method_binding->term_var) {
        out_type = matched_method_binding->term_var->type;
      }
      return true;
    }

    return false;
  }

  // Dunfield-Krishnaswami algorithmic context manipulation
  Context add_tvar(std::string name, ir::IrKind kind = ir::new_type_kind(), std::vector<ir::IrType> bounds = {}) const {
    return add_bind(new_tvar_bind(std::move(name), std::move(kind), std::move(bounds)));
  }

  Context add_tvar(const ir::TypeParam& tp) const {
    return add_bind(new_type_param_bind(tp));
  }

  Context add_exist_var(int64_t id) const {
    return add_bind(new_exist_var_bind(id));
  }

  Context add_solved_exist_var(int64_t id, ir::IrType solution) const {
    return add_bind(new_solved_exist_var_bind(id, std::move(solution)));
  }

  Context add_term_var(std::string name, ir::IrType type, bool is_def = false) const {
    return add_bind(new_term_var_bind(std::move(name), std::move(type), is_def));
  }

  Context add_marker(int64_t id) const {
    return add_bind(new_marker_bind(id));
  }

  Context drop_marker(int64_t id) const {
    List<Binding> curr = list_;
    while (!curr.empty()) {
      const Binding& b = curr.front();
      if (b.is_marker() && b.marker && b.marker->id == id) {
        List<Binding> rem = curr.remove();
        size_t wf_size = std::min(wellformed_size_, rem.size());
        return Context(std::move(rem), wf_size);
      }
      curr = curr.remove();
    }
    return *this;
  }

  std::pair<Context, std::vector<Binding>> split_at_marker(int64_t id) const {
    std::vector<Binding> after;
    List<Binding> curr = list_;
    while (!curr.empty()) {
      const Binding& b = curr.front();
      if (b.is_marker() && b.marker && b.marker->id == id) {
        List<Binding> rem = curr.remove();
        size_t wf_size = std::min(wellformed_size_, rem.size());
        std::reverse(after.begin(), after.end());
        return {Context(std::move(rem), wf_size), std::move(after)};
      }
      after.push_back(b);
      curr = curr.remove();
    }
    std::reverse(after.begin(), after.end());
    return {*this, std::move(after)};
  }

  std::pair<Context, std::vector<Binding>> split_at_exist_var(int64_t id) const {
    std::vector<Binding> after;
    List<Binding> curr = list_;
    while (!curr.empty()) {
      const Binding& b = curr.front();
      if ((b.is_exist_var() && b.exist_var && b.exist_var->id == id) ||
          (b.is_solved_exist_var() && b.solved_exist_var && b.solved_exist_var->id == id)) {
        List<Binding> rem = curr.remove();
        size_t wf_size = std::min(wellformed_size_, rem.size());
        std::reverse(after.begin(), after.end());
        return {Context(std::move(rem), wf_size), std::move(after)};
      }
      after.push_back(b);
      curr = curr.remove();
    }
    std::reverse(after.begin(), after.end());
    return {*this, std::move(after)};
  }

  bool contains_var(const std::string& name) const {
    Binding b;
    return lookup_type_param_bind(name, b) || lookup_term_decl_or_def_bind(name, b);
  }

  const Binding* lookup_exist_var_ptr(int64_t id) const {
    return lookup_bind_ptr([&](const Binding& bind) {
      return (bind.is_exist_var() && bind.exist_var && bind.exist_var->id == id) ||
             (bind.is_solved_exist_var() && bind.solved_exist_var && bind.solved_exist_var->id == id);
    });
  }

  bool contains_exist_var(int64_t id) const {
    return lookup_exist_var_ptr(id) != nullptr;
  }

  bool lookup_exist_var(int64_t id, bool& out_solved, ir::IrType& out_solution) const {
    const Binding* b = lookup_exist_var_ptr(id);
    if (!b) return false;
    if (b->is_solved_exist_var() && b->solved_exist_var) {
      out_solved = true;
      out_solution = b->solved_exist_var->solution;
    } else {
      out_solved = false;
      out_solution = {};
    }
    return true;
  }

  std::map<int64_t, ir::IrType> collect_exist_var_solutions() const {
    std::map<int64_t, ir::IrType> solutions;
    list_.for_each([&](const Binding& b) {
      if (b.is_solved_exist_var() && b.solved_exist_var) {
        solutions.try_emplace(b.solved_exist_var->id, b.solved_exist_var->solution);
      }
    });
    return solutions;
  }

  ir::IrType substitute_exist_vars(const ir::IrType& typ) const {
    auto solutions = collect_exist_var_solutions();
    return apply_exist_substitutions(typ, solutions);
  }

  Context substitute_exist_vars_in_context() const {
    auto solutions = collect_exist_var_solutions();
    if (solutions.empty()) return *this;

    List<Binding> new_list;
    auto binds = list_.collect();
    for (const auto& b : binds) {
      if (b.is_solved_exist_var() && b.solved_exist_var) {
        ir::IrType new_sol = apply_exist_substitutions(b.solved_exist_var->solution, solutions);
        new_list = new_list.add(new_solved_exist_var_bind(b.solved_exist_var->id, std::move(new_sol)));
      } else if (b.is_term_decl() && b.term_decl) {
        ir::IrType new_t = apply_exist_substitutions(b.term_decl->type, solutions);
        new_list = new_list.add(new_term_decl_bind(b.term_decl->name, std::move(new_t)));
      } else if (b.is_term_def() && b.term_def) {
        ir::IrType new_t = apply_exist_substitutions(b.term_def->type, solutions);
        new_list = new_list.add(new_term_def_bind(b.term_def->name, std::move(new_t)));
      } else if (b.is_term_var() && b.term_var) {
        ir::IrType new_t = apply_exist_substitutions(b.term_var->type, solutions);
        new_list = new_list.add(new_term_var_bind(b.term_var->name, std::move(new_t), b.term_var->is_def));
      } else if (b.is_alias() && b.alias) {
        ir::IrType new_t = apply_exist_substitutions(b.alias->type, solutions);
        new_list = new_list.add(new_alias_bind(b.alias->name, std::move(new_t)));
      } else {
        new_list = new_list.add(b);
      }
    }
    return Context(std::move(new_list), wellformed_size_);
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
      c = c.add_bind(new_term_def_bind(arg.id, arg.type));
    }
    return c;
  }

  Context add_symbol(const ir::IrDecl& decl) const {
    switch (decl.case_val) {
      case ir::IrDeclCase::TermDecl:
        if (!decl.term) return *this;
        return add_bind(new_term_decl_bind(decl.term->id, decl.term->type));

      case ir::IrDeclCase::AliasDecl:
        if (!decl.alias) return *this;
        return add_bind(new_alias_bind(decl.alias->id, decl.alias->type));

      case ir::IrDeclCase::NameDecl:
        if (!decl.name) return *this;
        return add_bind(new_const_bind(decl.name->id, decl.name->kind));

      case ir::IrDeclCase::TraitDecl: {
        if (!decl.trait) return *this;
        Context c = add_bind(new_trait_bind(decl.trait->id, decl.trait->type_params, decl.trait->methods));
        for (const auto& m : decl.trait->methods) {
          std::vector<ir::IrType> args;
          args.reserve(m.args.size());
          for (const auto& arg : m.args) {
            args.push_back(ir::substitute_type(arg.type, ir::new_name_type("Self"), ir::new_var_type("Self")));
          }
          ir::IrType ret = ir::substitute_type(m.ret_type, ir::new_name_type("Self"), ir::new_var_type("Self"));
          ir::IrType method_type = ir::new_function_type(ir::new_tuple_type(std::move(args)), std::move(ret));
          for (auto it = decl.trait->type_params.rbegin(); it != decl.trait->type_params.rend(); ++it) {
            method_type = ir::new_forall_type(*it, std::move(method_type));
          }
          method_type = ir::new_forall_type(ir::TypeParam{"Self", ir::new_type_kind(), {}}, std::move(method_type));
          c = c.add_bind(new_term_decl_bind(decl.trait->id + "::" + m.id, std::move(method_type)));
        }
        return c;
      }
    }
    return *this;
  }

  Context add_bind(Binding bind) const {
    bool fast_wf = (bind.is_term_var() || bind.is_exist_var() || bind.is_solved_exist_var() ||
                    bind.is_marker() || bind.is_decl() || bind.is_trait_impl());
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
    if (!typ.is(ir::IrTypeCase::ForallType) && !typ.is(ir::IrTypeCase::LambdaType)) {
      return {*this, ir::TypeParam{}, typ};
    }

    if (typ.is(ir::IrTypeCase::ForallType) && typ.forall) {
      ir::IrType fresh_var = gen_fresh_var_type();
      std::string fresh_name = fresh_var.var;
      std::vector<ir::IrType> bounds;
      for (const auto& b : typ.forall->type_param.bounds) {
        bounds.push_back(ir::substitute_type(b, ir::new_var_type(typ.forall->type_param.var), fresh_var));
      }
      ir::TypeParam tp{fresh_name, typ.forall->type_param.kind, std::move(bounds)};
      ir::IrType body = typ.forall->type ? ir::substitute_type(*typ.forall->type, ir::new_var_type(typ.forall->type_param.var), fresh_var) : ir::IrType{};
      Context new_ctx = add_bind(new_type_param_bind(tp));
      return {new_ctx, tp, body};
    }

    if (typ.is(ir::IrTypeCase::LambdaType) && typ.lambda) {
      ir::IrType fresh_var = gen_fresh_var_type();
      std::string fresh_name = fresh_var.var;
      ir::TypeParam tp{fresh_name, typ.lambda->kind, {}};
      ir::IrType body = typ.lambda->type ? ir::substitute_type(*typ.lambda->type, ir::new_var_type(typ.lambda->var), fresh_var) : ir::IrType{};
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
      case BindCase::TermVarBind:
      case BindCase::ExistVarBind:
      case BindCase::SolvedExistVarBind:
      case BindCase::MarkerBind:
      case BindCase::DeclBind:
        break;
    }

    return true;
  }

 private:
  static ir::IrType apply_exist_substitutions(
      const ir::IrType& typ,
      const std::map<int64_t, ir::IrType>& solutions) {
    if (typ.is(ir::IrTypeCase::ExistVarType)) {
      auto it = solutions.find(typ.exist_var);
      if (it != solutions.end()) {
        return apply_exist_substitutions(it->second, solutions);
      }
      return typ;
    }

    switch (typ.case_val) {
      case ir::IrTypeCase::AppType: {
        if (!typ.app) return typ;
        ir::IrType fun = typ.app->fun ? apply_exist_substitutions(*typ.app->fun, solutions) : ir::IrType{};
        ir::IrType arg = typ.app->arg ? apply_exist_substitutions(*typ.app->arg, solutions) : ir::IrType{};
        return ir::new_app_type(std::move(fun), std::move(arg));
      }
      case ir::IrTypeCase::ArrayType: {
        if (!typ.array) return typ;
        ir::IrType elem = typ.array->elem_type ? apply_exist_substitutions(*typ.array->elem_type, solutions) : ir::IrType{};
        return ir::new_array_type(std::move(elem), typ.array->size);
      }
      case ir::IrTypeCase::ForallType: {
        if (!typ.forall) return typ;
        std::vector<ir::IrType> bounds;
        for (const auto& b : typ.forall->type_param.bounds) {
          bounds.push_back(apply_exist_substitutions(b, solutions));
        }
        ir::TypeParam tp{typ.forall->type_param.var, typ.forall->type_param.kind, std::move(bounds)};
        ir::IrType body = typ.forall->type ? apply_exist_substitutions(*typ.forall->type, solutions) : ir::IrType{};
        return ir::new_forall_type(std::move(tp), std::move(body));
      }
      case ir::IrTypeCase::FunType: {
        if (!typ.fun) return typ;
        ir::IrType arg = typ.fun->arg ? apply_exist_substitutions(*typ.fun->arg, solutions) : ir::IrType{};
        ir::IrType ret = typ.fun->ret ? apply_exist_substitutions(*typ.fun->ret, solutions) : ir::IrType{};
        return ir::new_function_type(std::move(arg), std::move(ret));
      }
      case ir::IrTypeCase::LambdaType: {
        if (!typ.lambda) return typ;
        ir::IrType body = typ.lambda->type ? apply_exist_substitutions(*typ.lambda->type, solutions) : ir::IrType{};
        return ir::new_lambda_type(typ.lambda->var, typ.lambda->kind, std::move(body));
      }
      case ir::IrTypeCase::StructType: {
        if (!typ.struct_data) return typ;
        std::vector<ir::StructField> fields;
        for (const auto& f : typ.struct_data->fields) {
          ir::IrType ft = f.type ? apply_exist_substitutions(*f.type, solutions) : ir::IrType{};
          fields.push_back(ir::StructField{f.id, std::make_shared<ir::IrType>(std::move(ft))});
        }
        return ir::new_struct_type(std::move(fields));
      }
      case ir::IrTypeCase::TupleType: {
        if (!typ.tuple_data) return typ;
        std::vector<ir::IrType> elems;
        for (const auto& e : typ.tuple_data->elems) {
          elems.push_back(apply_exist_substitutions(e, solutions));
        }
        return ir::new_tuple_type(std::move(elems));
      }
      case ir::IrTypeCase::VariantType: {
        if (!typ.variant_data) return typ;
        std::vector<ir::VariantTag> tags;
        for (const auto& tag : typ.variant_data->tags) {
          ir::IrType tt = tag.type ? apply_exist_substitutions(*tag.type, solutions) : ir::IrType{};
          tags.push_back(ir::VariantTag{tag.id, std::make_shared<ir::IrType>(std::move(tt))});
        }
        return ir::new_variant_type(std::move(tags));
      }
      default:
        return typ;
    }
  }

  List<Binding> list_;
  size_t wellformed_size_ = 0;
};

inline Context new_context() {
  return Context();
}

inline void reset_fresh_var_generators() {
  Context::reset_fresh_var_generators();
}

} // namespace ts
