#pragma once

#include "bin/ir_base.h"
#include "bin/ir_decl.h"
#include "bin/ir_type.h"
#include "ts/bind.h"
#include "ts/context.h"
#include "ts/rename_vars.h"

#include <algorithm>
#include <string>
#include <vector>

namespace ts {

inline bool contains_duplicates(std::vector<std::string> ids) {
  std::sort(ids.begin(), ids.end());
  for (size_t i = 0; i + 1 < ids.size(); ++i) {
    if (ids[i] == ids[i + 1]) {
      return true;
    }
  }
  return false;
}

inline bool is_wellformed_type(const Context& c, const ir::IrType& t, std::string& err) {
  switch (t.case_val) {
    case ir::IrTypeCase::AppType: {
      if (!t.app || !t.app->fun || !t.app->arg) {
        err = "invalid application type";
        return false;
      }
      if (!is_wellformed_type(c, *t.app->fun, err)) {
        return false;
      }
      return is_wellformed_type(c, *t.app->arg, err);
    }

    case ir::IrTypeCase::ArrayType: {
      if (!t.array || !t.array->elem_type) {
        err = "invalid array type";
        return false;
      }
      return is_wellformed_type(c, *t.array->elem_type, err);
    }

    case ir::IrTypeCase::ExistVarType:
      return true;

    case ir::IrTypeCase::ForallType: {
      auto [new_ctx, tp, body_type] = c.add_fresh_type(t);
      return is_wellformed_type(new_ctx, body_type, err);
    }

    case ir::IrTypeCase::FunType: {
      if (!t.fun || !t.fun->arg || !t.fun->ret) {
        err = "invalid function type";
        return false;
      }
      if (!is_wellformed_type(c, *t.fun->arg, err)) {
        return false;
      }
      return is_wellformed_type(c, *t.fun->ret, err);
    }

    case ir::IrTypeCase::LambdaType: {
      auto [new_ctx, tp, body_type] = c.add_fresh_type(t);
      return is_wellformed_type(new_ctx, body_type, err);
    }

    case ir::IrTypeCase::MuType: {
      if (!t.mu || !t.mu->body) {
        err = "invalid mu type";
        return false;
      }
      if (!t.mu->kind.is_type_kind()) {
        err = "recursive mu type must have kind *";
        return false;
      }
      if (t.mu->body->is(ir::IrTypeCase::VarType) && t.mu->body->var == t.mu->var) {
        err = "non-contractive recursive type: " + t.to_string();
        return false;
      }
      auto [new_ctx, tp, body_type] = c.add_fresh_type(t);
      return is_wellformed_type(new_ctx, body_type, err);
    }

    case ir::IrTypeCase::NameType:
      if (c.contains_alias_bind(t.name) ||
          c.contains_const_bind(t.name) ||
          c.contains_trait_bind(t.name)) {
        return true;
      }
      err = "type \"" + t.name + "\" is undefined";
      return false;

    case ir::IrTypeCase::StructType: {
      if (contains_duplicates(t.field_ids())) {
        err = "struct type " + t.to_string() + " contains duplicate fields";
        return false;
      }
      for (const auto& field : t.fields()) {
        if (field.type) {
          if (!is_wellformed_type(c, *field.type, err)) {
            return false;
          }
        }
      }
      return true;
    }

    case ir::IrTypeCase::TupleType: {
      for (const auto& elem : t.elems()) {
        if (!is_wellformed_type(c, elem, err)) {
          return false;
        }
      }
      return true;
    }

    case ir::IrTypeCase::VariantType: {
      if (contains_duplicates(t.tag_ids())) {
        err = "variant type " + t.to_string() + " contains duplicate tags";
        return false;
      }
      for (const auto& tag : t.tags()) {
        if (tag.type) {
          if (!is_wellformed_type(c, *tag.type, err)) {
            return false;
          }
        }
      }
      return true;
    }

    case ir::IrTypeCase::VarType:
      if (c.contains_type_param_bind(t.var)) {
        return true;
      }
      err = "\"" + t.to_string() + "\" is undefined";
      return false;
  }

  err = "unhandled type case in is_wellformed_type";
  return false;
}

inline bool is_wellformed_decl(const Context& c, const ir::IrDecl& decl, std::string& err) {
  switch (decl.tag()) {
    case ir::IrDeclData::Tag::alias_decl:
      return is_wellformed_type(c, decl.data.as_alias_decl().type, err);

    case ir::IrDeclData::Tag::name_decl:
      return true;

    case ir::IrDeclData::Tag::term_decl:
      return is_wellformed_type(c, decl.data.as_term_decl().type, err);

    case ir::IrDeclData::Tag::trait_decl: {
      const auto& td = decl.data.as_trait_decl();

      std::vector<std::string> tp_ids;
      tp_ids.reserve(td.type_params.size());
      for (const auto& tp : td.type_params) {
        tp_ids.push_back(tp.var);
      }
      if (contains_duplicates(tp_ids)) {
        err = "trait \"" + td.id + "\" contains duplicate type parameters";
        return false;
      }

      std::vector<std::string> method_ids;
      method_ids.reserve(td.methods.size());
      for (const auto& m : td.methods) {
        method_ids.push_back(m.id);
      }
      if (contains_duplicates(method_ids)) {
        err = "trait \"" + td.id + "\" contains duplicate methods";
        return false;
      }

      Context trait_ctx = c.enter_scope();
      trait_ctx = trait_ctx.add_tvar("Self");
      for (const auto& tp : td.type_params) {
        trait_ctx = trait_ctx.add_bind(new_type_param_bind(tp));
      }

      for (const auto& m : td.methods) {
        for (const auto& arg : m.args) {
          if (!is_wellformed_type(trait_ctx, arg.type, err)) {
            return false;
          }
        }
        if (!is_wellformed_type(trait_ctx, m.ret_type, err)) {
          return false;
        }
      }
      return true;
    }

    case ir::IrDeclData::Tag::impl_decl:
      return true;
  }

  err = "unhandled decl case in is_wellformed_decl";
  return false;
}

} // namespace ts
