#pragma once

#include "ir_function.h"

namespace ir {

using ::IrSignature;
using ::IrDeclData;
using ::IrDecl;
using ::IrTraitImpl;

inline std::string signature_to_string(const IrSignature& sig) {
  std::stringstream ss;
  ss << "fn " << sig.id;
  if (!sig.type_params.empty()) {
    ss << " [";
    Interleave(sig.type_params, [&]() { ss << ", "; }, [&](int, const TypeParam& tp) {
      ss << "'" << tp.var;
    });
    ss << "]";
  }
  ss << "(";
  Interleave(sig.args, [&]() { ss << ", "; }, [&](int, const FunctionArg& a) {
    ss << function_arg_to_string(a);
  });
  ss << ") -> " << type_to_string(sig.ret_type);
  return ss.str();
}

inline std::string decl_id(const IrDecl& d) {
  switch (d.data.tag()) {
    case IrDeclData::Tag::term_decl:
      return d.data.as_term_decl().id;
    case IrDeclData::Tag::alias_decl:
      return d.data.as_alias_decl().id;
    case IrDeclData::Tag::name_decl:
      return d.data.as_name_decl().id;
    case IrDeclData::Tag::trait_decl:
      return d.data.as_trait_decl().id;
    case IrDeclData::Tag::impl_decl: {
      const auto& impl = d.data.as_impl_decl();
      if (impl.case_val.is(ImplCase::Tag::inherent_impl)) {
        return "impl:" + type_to_string(impl.type_name);
      }
      return "impl:" + type_to_string(impl.trait_type) + ":" + type_to_string(impl.type_name);
    }
  }
  return "";
}

inline IrDecl new_term_decl(std::string id, IrType type, bool export_flag = false) {
  return IrDecl{
      IrDeclData::mk_term_decl({std::move(id), std::move(type)}),
      export_flag,
      Pos{},
  };
}

inline IrDecl new_alias_decl(std::string id, IrKind kind, IrType type, bool export_flag = false) {
  return IrDecl{
      IrDeclData::mk_alias_decl({std::move(id), std::move(kind), std::move(type)}),
      export_flag,
      Pos{},
  };
}

inline IrDecl new_name_decl(std::string id, IrKind kind, bool export_flag = false) {
  return IrDecl{
      IrDeclData::mk_name_decl({std::move(id), std::move(kind)}),
      export_flag,
      Pos{},
  };
}

inline IrDecl new_trait_decl(std::string id, std::vector<TypeParam> type_params, std::vector<IrSignature> methods, bool export_flag = false) {
  return IrDecl{
      IrDeclData::mk_trait_decl({std::move(id), std::move(type_params), std::move(methods)}),
      export_flag,
      Pos{},
  };
}

inline IrTraitImpl new_trait_impl(std::vector<TypeParam> type_params, IrType trait_type, IrType type_name, std::vector<IrFunction> methods = {}, Pos pos = {}) {
  return IrTraitImpl{
      ImplCase::mk_trait_impl(),
      std::move(type_params),
      std::move(trait_type),
      std::move(type_name),
      std::move(methods),
      std::move(pos),
  };
}

inline IrTraitImpl new_inherent_impl(std::vector<TypeParam> type_params, IrType type_name, std::vector<IrFunction> methods = {}, Pos pos = {}) {
  return IrTraitImpl{
      ImplCase::mk_inherent_impl(),
      std::move(type_params),
      IrType{},
      std::move(type_name),
      std::move(methods),
      std::move(pos),
  };
}

inline std::string decl_to_string(const IrDecl& d, bool with_pos = false) {
  std::stringstream ss;
  if (with_pos && !d.pos.filename.empty()) {
    ss << pos_to_string(d.pos, true) << " ";
  }
  if (d.export_flag) {
    ss << "pub ";
  }
  switch (d.data.tag()) {
    case IrDeclData::Tag::term_decl: {
      const auto& term = d.data.as_term_decl();
      ss << term.id << ": " << type_to_string(term.typ);
      break;
    }
    case IrDeclData::Tag::alias_decl: {
      const auto& alias = d.data.as_alias_decl();
      if (is_type_kind(alias.kind)) {
        ss << "type " << alias.id << " = " << type_to_string(alias.typ);
      } else {
        ss << "type " << alias.id << " :: " << kind_to_string(alias.kind) << " = " << type_to_string(alias.typ);
      }
      break;
    }
    case IrDeclData::Tag::name_decl: {
      const auto& name = d.data.as_name_decl();
      if (is_type_kind(name.kind)) {
        ss << "type " << name.id;
      } else {
        ss << "type " << name.id << " :: " << kind_to_string(name.kind);
      }
      break;
    }
    case IrDeclData::Tag::trait_decl: {
      const auto& trait = d.data.as_trait_decl();
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
        ss << "  " << signature_to_string(m) << "\n";
      }
      ss << "}";
      break;
    }
    case IrDeclData::Tag::impl_decl: {
      const auto& impl = d.data.as_impl_decl();
      ss << "impl";
      if (!impl.type_params.empty()) {
        ss << " [";
        Interleave(impl.type_params, [&]() { ss << ", "; }, [&](int, const TypeParam& tp) {
          ss << "'" << tp.var;
        });
        ss << "]";
      }
      if (impl.case_val.is(ImplCase::Tag::inherent_impl)) {
        ss << " " << type_to_string(impl.type_name) << " {\n";
      } else {
        ss << " " << type_to_string(impl.trait_type) << " for " << type_to_string(impl.type_name) << " {\n";
      }
      for (const auto& m : impl.methods) {
        ss << "  " << signature_to_string(m) << "\n";
      }
      ss << "}";
      break;
    }
  }
  return ss.str();
}

inline IrDecl trait_impl_decl(const IrTraitImpl& ti) {
  std::vector<IrSignature> sigs;
  sigs.reserve(ti.methods.size());
  for (const auto& m : ti.methods) {
    sigs.push_back(IrSignature{m.id, m.type_params, m.args, m.ret_type});
  }
  return IrDecl{
      IrDeclData::mk_impl_decl({ti.case_val, ti.type_params, ti.trait_type, ti.type_name, std::move(sigs)}),
      false,
      ti.pos,
  };
}

inline std::string trait_impl_to_string(const IrTraitImpl& ti, bool with_pos = false) {
  std::stringstream ss;
  if (with_pos && !ti.pos.filename.empty()) {
    ss << pos_to_string(ti.pos, true);
  }
  ss << "impl";
  if (!ti.type_params.empty()) {
    ss << " [";
    Interleave(ti.type_params, [&]() { ss << ", "; }, [&](int, const TypeParam& tp) {
      ss << "'" << tp.var;
    });
    ss << "]";
  }
  if (ti.case_val.is(ImplCase::Tag::inherent_impl)) {
    ss << " " << type_to_string(ti.type_name) << " {\n";
  } else {
    ss << " " << type_to_string(ti.trait_type) << " for " << type_to_string(ti.type_name) << " {\n";
  }
  for (const auto& m : ti.methods) {
    ss << "  " << function_to_string(m, false) << "\n";
  }
  ss << "}";
  return ss.str();
}

inline IrDecl function_decl(const IrFunction& f) {
  std::vector<IrType> arg_types;
  arg_types.reserve(f.args.size());
  for (const auto& a : f.args) {
    arg_types.push_back(a.typ);
  }
  IrType t = new_function_type(new_tuple_type(std::move(arg_types)), f.ret_type);
  for (auto it = f.type_params.rbegin(); it != f.type_params.rend(); ++it) {
    t = new_forall_type(*it, std::move(t));
  }
  IrDecl d = new_term_decl(f.id, std::move(t), f.export_flag);
  d.pos = f.pos;
  return d;
}

} // namespace ir
