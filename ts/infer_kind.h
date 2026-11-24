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

  if (!ir::is_arrow_kind(fun_kind)) {
    out_err = "expected arrow kind (" + ir::kind_to_string(ir::new_arrow_kind(ir::new_type_kind(), ir::new_type_kind())) + ") in type application";
    return false;
  }

  const auto& ak = fun_kind.as_arrow_kind();
  if (!ak.left || !ir::equals_kind(*ak.left, arg_kind)) {
    std::string expected_arg = ak.left ? ir::kind_to_string(*ak.left) : "∗";
    out_err = "expected argument in type application (" + ir::kind_to_string(arg_kind) + ") to match function argument (" + expected_arg + ")";
    return false;
  }

  out_kind = ak.right ? *ak.right : ir::new_type_kind();
  return true;
}

inline bool infer_kind(const Context& context, const ir::IrType& typ, ir::IrKind& out_kind, std::string& out_err) {
  switch (typ.data.tag()) {
    case ir::IrTypeData::Tag::app_t: {
      const auto& a = typ.data.as_app_t();
      if (!a.fun || !a.arg) {
        out_err = "invalid AppType in infer_kind";
        return false;
      }
      return infer_kind_apply(context, *a.fun, *a.arg, out_kind, out_err);
    }

    case ir::IrTypeData::Tag::array_t:
    case ir::IrTypeData::Tag::exist_var_t:
      out_kind = ir::new_type_kind();
      return true;

    case ir::IrTypeData::Tag::forall_t: {
      auto [new_context, tp, body_type] = context.add_fresh_type(typ);
      ir::IrKind body_kind;
      if (!infer_kind(new_context, body_type, body_kind, out_err)) {
        return false;
      }
      if (!ir::is_type_kind(body_kind)) {
        out_err = "expected type " + ir::type_to_string(typ) + " to have kind " + ir::kind_to_string(ir::new_type_kind()) + ", but got kind " + ir::kind_to_string(body_kind);
        return false;
      }
      out_kind = body_kind;
      return true;
    }

    case ir::IrTypeData::Tag::fun_t:
      out_kind = ir::new_type_kind();
      return true;

    case ir::IrTypeData::Tag::lambda_t: {
      auto [new_context, tp, body_type] = context.add_fresh_type(typ);
      ir::IrKind ret_kind;
      if (!infer_kind(new_context, body_type, ret_kind, out_err)) {
        return false;
      }
      ir::IrKind lambda_kind = typ.data.as_lambda_t().kind;
      out_kind = ir::new_arrow_kind(std::move(lambda_kind), std::move(ret_kind));
      return true;
    }

    case ir::IrTypeData::Tag::mu_t: {
      const auto& m = typ.data.as_mu_t();
      if (!m.body) {
        out_err = "invalid MuType in infer_kind";
        return false;
      }
      if (!ir::is_type_kind(m.kind)) {
        out_err = "recursive mu type must have kind *";
        return false;
      }
      auto [new_context, tp, body_type] = context.add_fresh_type(typ);
      ir::IrKind body_kind;
      if (!infer_kind(new_context, body_type, body_kind, out_err)) {
        return false;
      }
      if (!ir::is_type_kind(body_kind)) {
        out_err = "expected mu type body " + ir::type_to_string(typ) + " to have kind " + ir::kind_to_string(ir::new_type_kind()) + ", but got kind " + ir::kind_to_string(body_kind);
        return false;
      }
      out_kind = body_kind;
      return true;
    }

    case ir::IrTypeData::Tag::name_t: {
      const auto& name = typ.data.as_name_t();
      if (context.contains_alias_bind(name)) {
        Binding bind;
        if (context.lookup_alias_bind(name, bind) && bind.alias) {
          return infer_kind(context, bind.alias->type, out_kind, out_err);
        }
      }
      if (context.contains_const_bind(name)) {
        Binding bind;
        if (context.lookup_const_bind(name, bind) && bind.const_data) {
          out_kind = bind.const_data->kind;
          return true;
        }
      }
      out_err = "type \"" + name + "\" is undefined";
      return false;
    }

    case ir::IrTypeData::Tag::struct_t:
    case ir::IrTypeData::Tag::tuple_t:
    case ir::IrTypeData::Tag::variant_t:
      out_kind = ir::new_type_kind();
      return true;

    case ir::IrTypeData::Tag::var_t: {
      const auto& var = typ.data.as_var_t();
      Binding bind;
      if (context.lookup_type_param_bind(var, bind) && bind.type_param) {
        out_kind = bind.type_param->kind;
        return true;
      }
      out_err = "type parameter \"" + var + "\" is undefined";
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
    throw std::runtime_error(err + "\n  inferring kind for type " + ir::type_to_string(typ));
  }
  return kind;
}

} // namespace ts
