#pragma once

#include "bin/ir_base.h"
#include "bin/ir_type.h"
#include "ts/bind.h"
#include "ts/context.h"
#include "ts/wellformed.h"

#include <stdexcept>
#include <string>
#include <vector>

namespace ts {

inline ir::IrType reduce_whnf(const Context& ctx, const ir::IrType& typ) {
  switch (typ.data.tag()) {
    case ir::IrTypeData::Tag::app_t: {
      const auto& a = typ.data.as_app_t();
      ir::IrType fun = a.fun ? reduce_whnf(ctx, *a.fun) : ir::IrType{};
      if (fun.data.is(ir::IrTypeData::Tag::lambda_t) && fun.data.as_lambda_t().body) {
        const auto& l = fun.data.as_lambda_t();
        ir::IrType arg = a.arg ? *a.arg : ir::IrType{};
        return reduce_whnf(ctx, ir::substitute_type(*l.body, ir::new_var_type(l.var), arg));
      }
      return ir::new_app_type(std::move(fun), a.arg ? *a.arg : ir::IrType{});
    }

    case ir::IrTypeData::Tag::name_t: {
      const auto& name = typ.data.as_name_t();
      if (ctx.contains_alias_bind(name)) {
        Binding bind;
        if (ctx.lookup_alias_bind(name, bind) && bind.alias) {
          return reduce_whnf(ctx, bind.alias->type);
        }
      }
      return typ;
    }

    default:
      return typ;
  }
}

inline ir::IrType reduce_type(const Context& ctx, const ir::IrType& typ) {
  std::string err;
  if (!is_wellformed_type(ctx, typ, err)) {
    throw std::runtime_error("failed to reduce " + ir::type_to_string(typ) + " because type is not wellformed: " + err);
  }
  ir::IrType reduced = reduce_whnf(ctx, typ);
  reduced.pos = typ.pos;
  return reduced;
}

inline ir::IrType unfold_mu_type(const ir::IrType& typ) {
  if (typ.data.is(ir::IrTypeData::Tag::mu_t) && typ.data.as_mu_t().body) {
    const auto& m = typ.data.as_mu_t();
    ir::IrType unfolded = ir::substitute_type(*m.body, ir::new_var_type(m.var), typ);
    unfolded.pos = typ.pos;
    return unfolded;
  }
  return typ;
}

} // namespace ts
