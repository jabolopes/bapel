#pragma once

#include "bin/ir_base.h"
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

inline bool wellformed_type(const Context& c, const ir::IrType& t, std::string& err) {
  return is_wellformed_type(c, t, err);
}

inline bool is_wellformed_context(const Context& ctx, std::string& err) {
  return ctx.is_wellformed(err);
}

} // namespace ts
