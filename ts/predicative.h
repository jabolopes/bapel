#pragma once

#include "bin/ir_base.h"
#include "bin/ir_type.h"
#include "ts/bind.h"
#include "ts/context.h"
#include "ts/reduce_type.h"

#include <stdexcept>
#include <string>
#include <utility>
#include <vector>

namespace ts {

class TypePredicator {
 public:
  explicit TypePredicator(Context ctx) : context_(std::move(ctx)) {}

  std::pair<ir::IrType, std::vector<ir::TypeParam>> predicate(const ir::IrType& typ) {
    tvars_.clear();
    ir::IrType res = predicate_impl(typ);
    res.pos = typ.pos;
    return {std::move(res), tvars_};
  }

  const Context& context() const { return context_; }

 private:
  ir::IrType predicate_impl(const ir::IrType& typ) {
    switch (typ.tag()) {
      case ir::IrTypeData::Tag::forall_t:
      case ir::IrTypeData::Tag::lambda_t: {
        auto [new_ctx, tp, body_type] = context_.add_fresh_type(typ);
        context_ = new_ctx;
        tvars_.push_back(tp);
        return predicate_impl(body_type);
      }
      case ir::IrTypeData::Tag::mu_t: {
        if (!typ.data.as_mu_t().body) return typ;
        Context orig_ctx = context_;
        auto [new_ctx, tp, body_type] = context_.add_fresh_type(typ);
        context_ = new_ctx;
        ir::IrType pred_body = predicate_impl(body_type);
        context_ = orig_ctx;
        return ir::new_mu_type(tp.var, tp.kind, std::move(pred_body));
      }
      default:
        return ir::map_type(typ, [&](const ir::IrType& child) {
          return predicate_impl(child);
        });
    }
  }

  Context context_;
  std::vector<ir::TypeParam> tvars_;
};

inline ir::IrType predicate_type(const Context& ctx, const ir::IrType& typ) {
  TypePredicator predicator(ctx);
  auto [body, tvars] = predicator.predicate(typ);
  return ir::forall_vars(tvars, std::move(body));
}

inline ir::IrType reduce_and_predicate_type(const Context& ctx, const ir::IrType& typ) {
  return predicate_type(ctx, reduce_type(ctx, typ));
}

} // namespace ts
