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

  std::string to_string() const;
  std::string to_json() const;
};

struct TraitDeclData {
  std::string id;
  std::vector<TypeParam> type_params;
  std::vector<IrSignature> methods;
};

struct IrDecl {
  IrDeclCase case_val = IrDeclCase::TermDecl;
  std::shared_ptr<TermDeclData> term;
  std::shared_ptr<AliasDeclData> alias;
  std::shared_ptr<NameDeclData> name;
  std::shared_ptr<TraitDeclData> trait;
  bool export_flag = false;
  Pos pos;

  bool is(IrDeclCase c) const { return case_val == c; }
  std::string to_string(bool with_pos = false) const;
  std::string to_json() const;

  std::string id() const {
    switch (case_val) {
      case IrDeclCase::TermDecl:
        return term ? term->id : "";
      case IrDeclCase::AliasDecl:
        return alias ? alias->id : "";
      case IrDeclCase::NameDecl:
        return name ? name->id : "";
      case IrDeclCase::TraitDecl:
        return trait ? trait->id : "";
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

  std::string to_string(bool with_pos = false) const;
  std::string to_json() const;
};

inline IrDecl new_term_decl(std::string id, IrType type, bool export_flag = false) {
  IrDecl d;
  d.case_val = IrDeclCase::TermDecl;
  d.term = std::make_shared<TermDeclData>();
  d.term->id = std::move(id);
  d.term->type = std::move(type);
  d.export_flag = export_flag;
  return d;
}

inline IrDecl new_alias_decl(std::string id, IrKind kind, IrType type, bool export_flag = false) {
  IrDecl d;
  d.case_val = IrDeclCase::AliasDecl;
  d.alias = std::make_shared<AliasDeclData>();
  d.alias->id = std::move(id);
  d.alias->kind = std::move(kind);
  d.alias->type = std::move(type);
  d.export_flag = export_flag;
  return d;
}

inline IrDecl new_name_decl(std::string id, IrKind kind, bool export_flag = false) {
  IrDecl d;
  d.case_val = IrDeclCase::NameDecl;
  d.name = std::make_shared<NameDeclData>();
  d.name->id = std::move(id);
  d.name->kind = std::move(kind);
  d.export_flag = export_flag;
  return d;
}

inline IrDecl new_trait_decl(std::string id, std::vector<TypeParam> type_params, std::vector<IrSignature> methods, bool export_flag = false) {
  IrDecl d;
  d.case_val = IrDeclCase::TraitDecl;
  d.trait = std::make_shared<TraitDeclData>();
  d.trait->id = std::move(id);
  d.trait->type_params = std::move(type_params);
  d.trait->methods = std::move(methods);
  d.export_flag = export_flag;
  return d;
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
  ss << "fn " << id << "(";
  Interleave(args, [&]() { ss << ", "; }, [&](int, const FunctionArg& a) {
    ss << a.to_string();
  });
  ss << ") -> " << ret_type.to_string();
  return ss.str();
}

inline std::string IrSignature::to_json() const {
  std::stringstream ss;
  ss << "{\"ID\":\"" << json_escape(id) << "\",\"Args\":[";
  Interleave(args, [&]() { ss << ","; }, [&](int, const FunctionArg& a) {
    ss << a.to_json();
  });
  ss << "],\"RetType\":" << ret_type.to_json() << "}";
  return ss.str();
}

inline std::string IrDecl::to_string(bool with_pos) const {
  std::stringstream ss;
  if (with_pos && !pos.filename.empty()) {
    ss << pos.to_string(true) << " ";
  }
  if (export_flag) {
    ss << "export ";
  }
  switch (case_val) {
    case IrDeclCase::TermDecl:
      if (term) {
        ss << term->id << ": " << term->type.to_string();
      }
      break;
    case IrDeclCase::AliasDecl:
      if (alias) {
        if (alias->kind.is_type_kind()) {
          ss << "type " << alias->id << " = " << alias->type.to_string();
        } else {
          ss << "type " << alias->id << " :: " << alias->kind.to_string() << " = " << alias->type.to_string();
        }
      }
      break;
    case IrDeclCase::NameDecl:
      if (name) {
        if (name->kind.is_type_kind()) {
          ss << "type " << name->id;
        } else {
          ss << "type " << name->id << " :: " << name->kind.to_string();
        }
      }
      break;
    case IrDeclCase::TraitDecl:
      if (trait) {
        ss << "trait " << trait->id;
        if (!trait->type_params.empty()) {
          ss << " [";
          Interleave(trait->type_params, [&]() { ss << ", "; }, [&](int, const TypeParam& tp) {
            ss << "'" << tp.var;
          });
          ss << "]";
        }
        ss << " {\n";
        for (const auto& m : trait->methods) {
          ss << "  " << m.to_string() << "\n";
        }
        ss << "}";
      }
      break;
  }
  return ss.str();
}

inline std::string IrDecl::to_json() const {
  std::stringstream ss;
  ss << "{\"Case\":" << static_cast<int>(case_val);
  switch (case_val) {
    case IrDeclCase::TermDecl:
      if (term) {
        ss << ",\"Term\":{\"ID\":\"" << json_escape(term->id) << "\",\"Type\":" << term->type.to_json() << "}";
      }
      break;
    case IrDeclCase::AliasDecl:
      if (alias) {
        ss << ",\"Alias\":{\"ID\":\"" << json_escape(alias->id) << "\",\"Kind\":" << alias->kind.to_json()
           << ",\"Type\":" << alias->type.to_json() << "}";
      }
      break;
    case IrDeclCase::NameDecl:
      if (name) {
        ss << ",\"Name\":{\"ID\":\"" << json_escape(name->id) << "\",\"Kind\":" << name->kind.to_json() << "}";
      }
      break;
    case IrDeclCase::TraitDecl:
      if (trait) {
        ss << ",\"Trait\":{\"ID\":\"" << json_escape(trait->id) << "\",\"TypeParams\":[";
        Interleave(trait->type_params, [&]() { ss << ","; }, [&](int, const TypeParam& tp) {
          ss << tp.to_json();
        });
        ss << "],\"Methods\":[";
        Interleave(trait->methods, [&]() { ss << ","; }, [&](int, const IrSignature& m) {
          ss << m.to_json();
        });
        ss << "]}";
      }
      break;
  }
  ss << ",\"Export\":" << (export_flag ? "true" : "false")
     << ",\"Pos\":" << pos.to_json() << "}";
  return ss.str();
}

inline std::string IrTraitImpl::to_string(bool with_pos) const {
  std::stringstream ss;
  if (with_pos && !pos.filename.empty()) {
    ss << pos.to_string(true);
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

inline std::string IrTraitImpl::to_json() const {
  std::stringstream ss;
  ss << "{\"Case\":" << static_cast<int>(case_val)
     << ",\"TypeParams\":[";
  Interleave(type_params, [&]() { ss << ","; }, [&](int, const TypeParam& tp) {
    ss << tp.to_json();
  });
  ss << "],\"TraitType\":" << trait_type.to_json()
     << ",\"TypeName\":" << type_name.to_json()
     << ",\"Methods\":[";
  Interleave(methods, [&]() { ss << ","; }, [&](int, const IrFunction& m) {
    ss << m.to_json();
  });
  ss << "],\"Pos\":" << pos.to_json() << "}";
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
