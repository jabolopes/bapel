#pragma once

#include "bin/ir_base.h"
#include "bin/ir_type.h"
#include "ts/bind.h"
#include "ts/context.h"
#include "ts/rename_vars.h"

#include <stdexcept>
#include <string>
#include <utility>

namespace ts {

inline bool infer_kind(const Context& context, const ir::IrType& typ, ir::IrKind& out_kind, std::string& out_err);

inline bool infer_kind_apply(
    const Context& context,
    const ir::IrType& fun,
    const ir::IrType& arg,
    ir::IrKind& out_kind,
    std::string& out_err) {
  ir::IrKind fun_kind;
  if (!infer_kind(context, fun, fun_kind, out_err)) {
    return false;
  }

  ir::IrKind arg_kind;
  if (!infer_kind(context, arg, arg_kind, out_err)) {
    return false;
  }

  if (!fun_kind.is_arrow_kind()) {
    out_err = "expected arrow kind (" + ir::new_arrow_kind(ir::new_type_kind(), ir::new_type_kind()).to_string() + ") in type application";
    return false;
  }

  if (!fun_kind.left || !ir::equals_kind(*fun_kind.left, arg_kind)) {
    std::string expected_arg = fun_kind.left ? fun_kind.left->to_string() : "∗";
    out_err = "expected argument in type application (" + arg_kind.to_string() + ") to match function argument (" + expected_arg + ")";
    return false;
  }

  out_kind = fun_kind.right ? *fun_kind.right : ir::new_type_kind();
  return true;
}

inline bool infer_kind(const Context& context, const ir::IrType& typ, ir::IrKind& out_kind, std::string& out_err) {
  switch (typ.case_val) {
    case ir::IrTypeCase::AppType: {
      if (!typ.app || !typ.app->fun || !typ.app->arg) {
        out_err = "invalid AppType in infer_kind";
        return false;
      }
      return infer_kind_apply(context, *typ.app->fun, *typ.app->arg, out_kind, out_err);
    }

    case ir::IrTypeCase::ArrayType:
    case ir::IrTypeCase::ExistVarType:
      out_kind = ir::new_type_kind();
      return true;

    case ir::IrTypeCase::ForallType: {
      auto [new_context, tp, body_type] = context.add_fresh_type(typ);
      ir::IrKind body_kind;
      if (!infer_kind(new_context, body_type, body_kind, out_err)) {
        return false;
      }
      if (!body_kind.is_type_kind()) {
        out_err = "expected type " + typ.to_string() + " to have kind " + ir::new_type_kind().to_string() + ", but got kind " + body_kind.to_string();
        return false;
      }
      out_kind = body_kind;
      return true;
    }

    case ir::IrTypeCase::FunType:
      out_kind = ir::new_type_kind();
      return true;

    case ir::IrTypeCase::LambdaType: {
      auto [new_context, tp, body_type] = context.add_fresh_type(typ);
      ir::IrKind ret_kind;
      if (!infer_kind(new_context, body_type, ret_kind, out_err)) {
        return false;
      }
      ir::IrKind lambda_kind = typ.lambda ? typ.lambda->kind : ir::new_type_kind();
      out_kind = ir::new_arrow_kind(std::move(lambda_kind), std::move(ret_kind));
      return true;
    }

    case ir::IrTypeCase::NameType: {
      if (context.contains_alias_bind(typ.name)) {
        Binding bind;
        if (context.lookup_alias_bind(typ.name, bind) && bind.alias) {
          return infer_kind(context, bind.alias->type, out_kind, out_err);
        }
      }
      if (context.contains_const_bind(typ.name)) {
        Binding bind;
        if (context.lookup_const_bind(typ.name, bind) && bind.const_data) {
          out_kind = bind.const_data->kind;
          return true;
        }
      }
      out_err = "type \"" + typ.name + "\" is undefined";
      return false;
    }

    case ir::IrTypeCase::StructType:
    case ir::IrTypeCase::TupleType:
    case ir::IrTypeCase::VariantType:
      out_kind = ir::new_type_kind();
      return true;

    case ir::IrTypeCase::VarType: {
      Binding bind;
      if (context.lookup_type_param_bind(typ.var, bind) && bind.type_param) {
        out_kind = bind.type_param->kind;
        return true;
      }
      out_err = "type parameter \"" + typ.var + "\" is undefined";
      return false;
    }
  }

  out_err = "unhandled type case in infer_kind";
  return false;
}

inline ir::IrKind infer_kind(const Context& context, const ir::IrType& typ) {
  ir::IrKind kind;
  std::string err;
  if (!infer_kind(context, typ, kind, err)) {
    throw std::runtime_error(err + "\n  inferring kind for type " + typ.to_string());
  }
  return kind;
}

} // namespace ts
