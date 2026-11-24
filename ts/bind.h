#pragma once

#include "bin/ir_base.h"
#include "bin/ir_decl.h"
#include "bin/ir_function.h"
#include "bin/ir_type.h"

#include <cstdint>
#include <memory>
#include <sstream>
#include <string>
#include <vector>

namespace ts {

enum class BindCase {
  AliasBind = 0,
  ConstBind = 1,
  ScopeBind = 2,
  TermDeclBind = 3,
  TermDefBind = 4,
  TypeParamBind = 5,
  TraitBind = 6,
  TraitImplBind = 7
};

struct AliasBindData {
  std::string name;
  ir::IrType type;

  std::string to_string() const {
    return "type " + name + " = " + type.to_string();
  }

  bool operator==(const AliasBindData& other) const {
    return name == other.name && ir::equals_type(type, other.type);
  }
};

struct ConstBindData {
  std::string name;
  ir::IrKind kind;

  std::string to_string() const {
    return "type " + name + " :: " + ir::kind_to_string(kind);
  }

  bool operator==(const ConstBindData& other) const {
    return name == other.name && ir::equals_kind(kind, other.kind);
  }
};

struct ScopeBindData {
  int level = 0;

  std::string to_string() const {
    return "scope " + std::to_string(level);
  }

  bool operator==(const ScopeBindData& other) const {
    return level == other.level;
  }
};

struct TermDeclBindData {
  std::string name;
  ir::IrType type;

  std::string to_string() const {
    return name + ": " + type.to_string();
  }

  bool operator==(const TermDeclBindData& other) const {
    return name == other.name && ir::equals_type(type, other.type);
  }
};

struct TermDefBindData {
  std::string name;
  ir::IrType type;

  std::string to_string() const {
    return "let " + name + ": " + type.to_string();
  }

  bool operator==(const TermDefBindData& other) const {
    return name == other.name && ir::equals_type(type, other.type);
  }
};

struct TypeParamBindData {
  std::string name;
  ir::IrKind kind;
  std::vector<ir::IrType> bounds;

  std::string to_string() const {
    std::string s = "type '" + name;
    if (!bounds.empty()) {
      s += ": ";
      ir::Interleave(bounds, [&]() { s += " + "; }, [&](int, const ir::IrType& b) {
        s += b.to_string();
      });
    }
    return s;
  }

  bool operator==(const TypeParamBindData& other) const {
    if (name != other.name || !ir::equals_kind(kind, other.kind)) return false;
    if (bounds.size() != other.bounds.size()) return false;
    for (size_t i = 0; i < bounds.size(); ++i) {
      if (!ir::equals_type(bounds[i], other.bounds[i])) return false;
    }
    return true;
  }
};

struct TraitBindData {
  std::string name;
  std::vector<ir::TypeParam> type_params;
  std::vector<ir::IrSignature> methods;

  std::string to_string() const {
    return "trait " + name;
  }

  bool operator==(const TraitBindData& other) const {
    if (name != other.name) return false;
    if (type_params.size() != other.type_params.size()) return false;
    for (size_t i = 0; i < type_params.size(); ++i) {
      if (!(type_params[i] == other.type_params[i])) return false;
    }
    if (methods.size() != other.methods.size()) return false;
    for (size_t i = 0; i < methods.size(); ++i) {
      if (methods[i].id != other.methods[i].id) return false;
      if (!ir::equals_type(methods[i].ret_type, other.methods[i].ret_type)) return false;
      if (methods[i].args.size() != other.methods[i].args.size()) return false;
      for (size_t a = 0; a < methods[i].args.size(); ++a) {
        if (methods[i].args[a].id != other.methods[i].args[a].id ||
            !ir::equals_type(methods[i].args[a].type, other.methods[i].args[a].type)) {
          return false;
        }
      }
    }
    return true;
  }
};

struct TraitImplBindData {
  std::vector<ir::TypeParam> type_params;
  ir::IrType trait_type;
  ir::IrType type_name;
  std::vector<ir::IrFunction> methods;

  std::string to_string() const {
    return "impl " + trait_type.to_string() + " for " + type_name.to_string();
  }

  bool operator==(const TraitImplBindData& other) const {
    if (type_params.size() != other.type_params.size()) return false;
    for (size_t i = 0; i < type_params.size(); ++i) {
      if (!(type_params[i] == other.type_params[i])) return false;
    }
    return ir::equals_type(trait_type, other.trait_type) &&
           ir::equals_type(type_name, other.type_name);
  }
};

struct Binding {
  BindCase case_val = BindCase::ScopeBind;
  std::shared_ptr<AliasBindData> alias;
  std::shared_ptr<ConstBindData> const_data;
  std::shared_ptr<ScopeBindData> scope;
  std::shared_ptr<TermDeclBindData> term_decl;
  std::shared_ptr<TermDefBindData> term_def;
  std::shared_ptr<TypeParamBindData> type_param;
  std::shared_ptr<TraitBindData> trait;
  std::shared_ptr<TraitImplBindData> trait_impl;

  bool is(BindCase c) const { return case_val == c; }
  bool is_alias() const { return case_val == BindCase::AliasBind; }
  bool is_const() const { return case_val == BindCase::ConstBind; }
  bool is_scope() const { return case_val == BindCase::ScopeBind; }
  bool is_term_decl() const { return case_val == BindCase::TermDeclBind; }
  bool is_term_def() const { return case_val == BindCase::TermDefBind; }
  bool is_type_param() const { return case_val == BindCase::TypeParamBind; }
  bool is_tvar() const { return case_val == BindCase::TypeParamBind; }
  bool is_trait() const { return case_val == BindCase::TraitBind; }
  bool is_trait_impl() const { return case_val == BindCase::TraitImplBind; }
  bool is_impl() const { return case_val == BindCase::TraitImplBind; }

  std::string to_string() const {
    switch (case_val) {
      case BindCase::AliasBind:
        return alias ? alias->to_string() : "";
      case BindCase::ConstBind:
        return const_data ? const_data->to_string() : "";
      case BindCase::ScopeBind:
        return scope ? scope->to_string() : "";
      case BindCase::TermDeclBind:
        return term_decl ? term_decl->to_string() : "";
      case BindCase::TermDefBind:
        return term_def ? term_def->to_string() : "";
      case BindCase::TypeParamBind:
        return type_param ? type_param->to_string() : "";
      case BindCase::TraitBind:
        return trait ? trait->to_string() : "";
      case BindCase::TraitImplBind:
        return trait_impl ? trait_impl->to_string() : "";
    }
    return "";
  }

  bool operator==(const Binding& other) const {
    if (case_val != other.case_val) return false;
    switch (case_val) {
      case BindCase::AliasBind:
        if (!alias || !other.alias) return alias == other.alias;
        return *alias == *other.alias;
      case BindCase::ConstBind:
        if (!const_data || !other.const_data) return const_data == other.const_data;
        return *const_data == *other.const_data;
      case BindCase::ScopeBind:
        if (!scope || !other.scope) return scope == other.scope;
        return *scope == *other.scope;
      case BindCase::TermDeclBind:
        if (!term_decl || !other.term_decl) return term_decl == other.term_decl;
        return *term_decl == *other.term_decl;
      case BindCase::TermDefBind:
        if (!term_def || !other.term_def) return term_def == other.term_def;
        return *term_def == *other.term_def;
      case BindCase::TypeParamBind:
        if (!type_param || !other.type_param) return type_param == other.type_param;
        return *type_param == *other.type_param;
      case BindCase::TraitBind:
        if (!trait || !other.trait) return trait == other.trait;
        return *trait == *other.trait;
      case BindCase::TraitImplBind:
        if (!trait_impl || !other.trait_impl) return trait_impl == other.trait_impl;
        return *trait_impl == *other.trait_impl;
    }
    return false;
  }

  bool operator!=(const Binding& other) const {
    return !(*this == other);
  }
};

inline Binding new_alias_bind(std::string name, ir::IrType type) {
  Binding b;
  b.case_val = BindCase::AliasBind;
  b.alias = std::make_shared<AliasBindData>(AliasBindData{std::move(name), std::move(type)});
  return b;
}

inline Binding new_const_bind(std::string name, ir::IrKind kind) {
  Binding b;
  b.case_val = BindCase::ConstBind;
  b.const_data = std::make_shared<ConstBindData>(ConstBindData{std::move(name), std::move(kind)});
  return b;
}

inline Binding new_scope_bind(int level) {
  Binding b;
  b.case_val = BindCase::ScopeBind;
  b.scope = std::make_shared<ScopeBindData>(ScopeBindData{level});
  return b;
}

inline Binding new_term_decl_bind(std::string name, ir::IrType type) {
  Binding b;
  b.case_val = BindCase::TermDeclBind;
  b.term_decl = std::make_shared<TermDeclBindData>(TermDeclBindData{std::move(name), std::move(type)});
  return b;
}

inline Binding new_term_def_bind(std::string name, ir::IrType type) {
  Binding b;
  b.case_val = BindCase::TermDefBind;
  b.term_def = std::make_shared<TermDefBindData>(TermDefBindData{std::move(name), std::move(type)});
  return b;
}

inline Binding new_type_param_bind(std::string name, ir::IrKind kind = ir::new_type_kind(), std::vector<ir::IrType> bounds = {}) {
  Binding b;
  b.case_val = BindCase::TypeParamBind;
  b.type_param = std::make_shared<TypeParamBindData>(TypeParamBindData{std::move(name), std::move(kind), std::move(bounds)});
  return b;
}

inline Binding new_type_param_bind(ir::TypeParam tp) {
  return new_type_param_bind(std::move(tp.var), std::move(tp.kind), std::move(tp.bounds));
}

inline Binding new_tvar_bind(std::string name, ir::IrKind kind = ir::new_type_kind(), std::vector<ir::IrType> bounds = {}) {
  return new_type_param_bind(std::move(name), std::move(kind), std::move(bounds));
}

inline Binding new_tvar_bind(ir::TypeParam tp) {
  return new_type_param_bind(std::move(tp));
}

inline Binding new_trait_bind(std::string name, std::vector<ir::TypeParam> type_params, std::vector<ir::IrSignature> methods) {
  Binding b;
  b.case_val = BindCase::TraitBind;
  b.trait = std::make_shared<TraitBindData>(TraitBindData{std::move(name), std::move(type_params), std::move(methods)});
  return b;
}

inline Binding new_trait_impl_bind(std::vector<ir::TypeParam> type_params, ir::IrType trait_type, ir::IrType type_name, std::vector<ir::IrFunction> methods = {}) {
  Binding b;
  b.case_val = BindCase::TraitImplBind;
  b.trait_impl = std::make_shared<TraitImplBindData>(TraitImplBindData{std::move(type_params), std::move(trait_type), std::move(type_name), std::move(methods)});
  return b;
}

} // namespace ts
