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

class TypeReducer {
 public:
  ir::IrType reduce(const Context& ctx, const ir::IrType& typ) {
    std::string err;
    if (!is_wellformed_type(ctx, typ, err)) {
      throw std::runtime_error("failed to reduce " + ir::type_to_string(typ) + " because type is not wellformed: " + err);
    }
    ir::IrType reduced = reduce_impl(ctx, typ);
    reduced.pos = typ.pos;
    return reduced;
  }

 private:
  ir::IrType reduce_impl(const Context& ctx, const ir::IrType& typ) {
    switch (typ.data.tag()) {
      case ir::IrTypeData::Tag::app_t: {
        const auto& a = typ.data.as_app_t();
        ir::IrType fun = a.fun ? reduce(ctx, *a.fun) : ir::IrType{};
        ir::IrType arg = a.arg ? reduce(ctx, *a.arg) : ir::IrType{};

        if (fun.data.is(ir::IrTypeData::Tag::lambda_t) && fun.data.as_lambda_t().body) {
          const auto& l = fun.data.as_lambda_t();
          return reduce(ctx, ir::substitute_type(*l.body, ir::new_var_type(l.var), arg));
        }
        return ir::new_app_type(std::move(fun), std::move(arg));
      }

      case ir::IrTypeData::Tag::forall_t: {
        auto [new_ctx, tp, body_type] = ctx.add_fresh_type(typ);
        std::vector<ir::IrType> bounds;
        bounds.reserve(tp.bounds.size());
        for (const auto& b : tp.bounds) {
          bounds.push_back(reduce(ctx, b));
        }
        tp.bounds = std::move(bounds);
        return ir::new_forall_type(std::move(tp), reduce(new_ctx, body_type));
      }

      case ir::IrTypeData::Tag::lambda_t: {
        auto [new_ctx, tp, body_type] = ctx.add_fresh_type(typ);
        return ir::new_lambda_type(tp.var, tp.kind, reduce(new_ctx, body_type));
      }

      case ir::IrTypeData::Tag::mu_t: {
        if (!typ.data.as_mu_t().body) return typ;
        auto [new_ctx, tp, body_type] = ctx.add_fresh_type(typ);
        return ir::new_mu_type(tp.var, tp.kind, reduce(new_ctx, body_type));
      }

      case ir::IrTypeData::Tag::name_t: {
        const auto& name = typ.data.as_name_t();
        if (ctx.contains_alias_bind(name)) {
          Binding bind;
          if (ctx.lookup_alias_bind(name, bind) && bind.alias) {
            return reduce(ctx, bind.alias->type);
          }
        }
        return typ;
      }

      default:
        return ir::map_type(typ, [&](const ir::IrType& child) {
          return reduce(ctx, child);
        });
    }
  }
};

inline ir::IrType reduce_type(const Context& ctx, const ir::IrType& typ) {
  TypeReducer reducer;
  return reducer.reduce(ctx, typ);
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
