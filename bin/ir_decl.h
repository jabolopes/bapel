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

} // namespace ir
