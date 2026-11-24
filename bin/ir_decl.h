#pragma once

#include "ir_function.h"

namespace ir {

struct TermDeclData {
  std::string id;
  IrType type;
};

struct AliasDeclData {
  std::string id;
  IrKind kind;
  IrType type;
};

struct NameDeclData {
  std::string id;
  IrKind kind;
};

struct IrSignature {
  std::string id;
  std::vector<FunctionArg> args;
  IrType ret_type;
  std::vector<TypeParam> type_params;

  std::string to_string() const;
};

struct TraitDeclData {
  std::string id;
  std::vector<TypeParam> type_params;
  std::vector<IrSignature> methods;
};

struct ImplDeclData {
  ImplCase case_val = ImplCase::TraitImpl;
  std::vector<TypeParam> type_params;
  IrType trait_type;
  IrType type_name;
  std::vector<IrSignature> methods;
};

BAPEL_VARIANT(IrDeclData,
  (term_decl, (TermDeclData)),
  (alias_decl, (AliasDeclData)),
  (name_decl, (NameDeclData)),
  (trait_decl, (TraitDeclData)),
  (impl_decl, (ImplDeclData))
);

struct IrDecl {
  IrDeclData data{IrDeclData::mk_term_decl(TermDeclData{})};
  bool export_flag = false;
  Pos pos;

  IrDeclData::Tag tag() const { return data.tag(); }
  bool is(IrDeclData::Tag t) const { return data.is(t); }
  std::string to_string(bool with_pos = false) const;

  std::string id() const {
    switch (data.tag()) {
      case IrDeclData::Tag::term_decl:
        return data.as_term_decl().id;
      case IrDeclData::Tag::alias_decl:
        return data.as_alias_decl().id;
      case IrDeclData::Tag::name_decl:
        return data.as_name_decl().id;
      case IrDeclData::Tag::trait_decl:
        return data.as_trait_decl().id;
      case IrDeclData::Tag::impl_decl: {
        const auto& impl = data.as_impl_decl();
        if (impl.case_val == ImplCase::InherentImpl) {
          return "impl:" + impl.type_name.to_string();
        }
        return "impl:" + impl.trait_type.to_string() + ":" + impl.type_name.to_string();
      }
    }
    return "";
  }
};

struct IrTraitImpl {
  ImplCase case_val = ImplCase::TraitImpl;
  std::vector<TypeParam> type_params;
  IrType trait_type;
  IrType type_name;
  std::vector<IrFunction> methods;
  Pos pos;

  IrDecl decl() const;
  std::string to_string(bool with_pos = false) const;
};

inline IrDecl new_term_decl(std::string id, IrType type, bool export_flag = false) {
  return IrDecl{
      IrDeclData::mk_term_decl(TermDeclData{std::move(id), std::move(type)}),
      export_flag,
      Pos{},
  };
}

inline IrDecl new_alias_decl(std::string id, IrKind kind, IrType type, bool export_flag = false) {
  return IrDecl{
      IrDeclData::mk_alias_decl(AliasDeclData{std::move(id), std::move(kind), std::move(type)}),
      export_flag,
      Pos{},
  };
}

inline IrDecl new_name_decl(std::string id, IrKind kind, bool export_flag = false) {
  return IrDecl{
      IrDeclData::mk_name_decl(NameDeclData{std::move(id), std::move(kind)}),
      export_flag,
      Pos{},
  };
}

inline IrDecl new_trait_decl(std::string id, std::vector<TypeParam> type_params, std::vector<IrSignature> methods, bool export_flag = false) {
  return IrDecl{
      IrDeclData::mk_trait_decl(TraitDeclData{std::move(id), std::move(type_params), std::move(methods)}),
      export_flag,
      Pos{},
  };
}

inline IrTraitImpl new_trait_impl(std::vector<TypeParam> type_params, IrType trait_type, IrType type_name, std::vector<IrFunction> methods = {}, Pos pos = {}) {
  IrTraitImpl impl;
  impl.case_val = ImplCase::TraitImpl;
  impl.type_params = std::move(type_params);
  impl.trait_type = std::move(trait_type);
  impl.type_name = std::move(type_name);
  impl.methods = std::move(methods);
  impl.pos = std::move(pos);
  return impl;
}

inline IrTraitImpl new_inherent_impl(std::vector<TypeParam> type_params, IrType type_name, std::vector<IrFunction> methods = {}, Pos pos = {}) {
  IrTraitImpl impl;
  impl.case_val = ImplCase::InherentImpl;
  impl.type_params = std::move(type_params);
  impl.type_name = std::move(type_name);
  impl.methods = std::move(methods);
  impl.pos = std::move(pos);
  return impl;
}

inline std::string IrSignature::to_string() const {
  std::stringstream ss;
  ss << "fn " << id;
  if (!type_params.empty()) {
    ss << " [";
    Interleave(type_params, [&]() { ss << ", "; }, [&](int, const TypeParam& tp) {
      ss << "'" << tp.var;
    });
    ss << "]";
  }
  ss << "(";
  Interleave(args, [&]() { ss << ", "; }, [&](int, const FunctionArg& a) {
    ss << a.to_string();
  });
  ss << ") -> " << ret_type.to_string();
  return ss.str();
}

inline std::string IrDecl::to_string(bool with_pos) const {
  std::stringstream ss;
  if (with_pos && !pos.filename.empty()) {
    ss << pos_to_string(pos, true) << " ";
  }
  if (export_flag) {
    ss << "pub ";
  }
  switch (data.tag()) {
    case IrDeclData::Tag::term_decl: {
      const auto& term = data.as_term_decl();
      ss << term.id << ": " << term.type.to_string();
      break;
    }
    case IrDeclData::Tag::alias_decl: {
      const auto& alias = data.as_alias_decl();
      if (is_type_kind(alias.kind)) {
        ss << "type " << alias.id << " = " << alias.type.to_string();
      } else {
        ss << "type " << alias.id << " :: " << kind_to_string(alias.kind) << " = " << alias.type.to_string();
      }
      break;
    }
    case IrDeclData::Tag::name_decl: {
      const auto& name = data.as_name_decl();
      if (is_type_kind(name.kind)) {
        ss << "type " << name.id;
      } else {
        ss << "type " << name.id << " :: " << kind_to_string(name.kind);
      }
      break;
    }
    case IrDeclData::Tag::trait_decl: {
      const auto& trait = data.as_trait_decl();
      ss << "trait " << trait.id;
      if (!trait.type_params.empty()) {
        ss << " [";
        Interleave(trait.type_params, [&]() { ss << ", "; }, [&](int, const TypeParam& tp) {
          ss << "'" << tp.var;
        });
        ss << "]";
      }
      ss << " {\n";
      for (const auto& m : trait.methods) {
        ss << "  " << m.to_string() << "\n";
      }
      ss << "}";
      break;
    }
    case IrDeclData::Tag::impl_decl: {
      const auto& impl = data.as_impl_decl();
      ss << "impl";
      if (!impl.type_params.empty()) {
        ss << " [";
        Interleave(impl.type_params, [&]() { ss << ", "; }, [&](int, const TypeParam& tp) {
          ss << "'" << tp.var;
        });
        ss << "]";
      }
      if (impl.case_val == ImplCase::InherentImpl) {
        ss << " " << impl.type_name.to_string() << " {\n";
      } else {
        ss << " " << impl.trait_type.to_string() << " for " << impl.type_name.to_string() << " {\n";
      }
      for (const auto& m : impl.methods) {
        ss << "  " << m.to_string() << "\n";
      }
      ss << "}";
      break;
    }
  }
  return ss.str();
}

inline IrDecl IrTraitImpl::decl() const {
  ImplDeclData impl_data;
  impl_data.case_val = case_val;
  impl_data.type_params = type_params;
  impl_data.trait_type = trait_type;
  impl_data.type_name = type_name;
  for (const auto& m : methods) {
    impl_data.methods.push_back(IrSignature{m.id, m.args, m.ret_type, m.type_params});
  }
  return IrDecl{IrDeclData::mk_impl_decl(std::move(impl_data)), false, pos};
}

inline std::string IrTraitImpl::to_string(bool with_pos) const {
  std::stringstream ss;
  if (with_pos && !pos.filename.empty()) {
    ss << pos_to_string(pos, true);
  }
  ss << "impl";
  if (!type_params.empty()) {
    ss << " [";
    Interleave(type_params, [&]() { ss << ", "; }, [&](int, const TypeParam& tp) {
      ss << "'" << tp.var;
    });
    ss << "]";
  }
  if (case_val == ImplCase::InherentImpl) {
    ss << " " << type_name.to_string() << " {\n";
  } else {
    ss << " " << trait_type.to_string() << " for " << type_name.to_string() << " {\n";
  }
  for (const auto& m : methods) {
    ss << "  " << m.to_string(false) << "\n";
  }
  ss << "}";
  return ss.str();
}

inline IrDecl IrFunction::decl() const {
  std::vector<IrType> arg_types;
  arg_types.reserve(args.size());
  for (const auto& a : args) {
    arg_types.push_back(a.type);
  }
  IrType t = new_function_type(new_tuple_type(std::move(arg_types)), ret_type);
  for (auto it = type_params.rbegin(); it != type_params.rend(); ++it) {
    t = new_forall_type(*it, std::move(t));
  }
  IrDecl d = new_term_decl(id, std::move(t), export_flag);
  d.pos = pos;
  return d;
}

} // namespace ir
