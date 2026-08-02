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
    switch (typ.case_val) {
      case ir::IrTypeCase::AppType: {
        if (!typ.app) return typ;
        ir::IrType fun = typ.app->fun ? predicate_impl(*typ.app->fun) : ir::IrType{};
        ir::IrType arg = typ.app->arg ? predicate_impl(*typ.app->arg) : ir::IrType{};
        return ir::new_app_type(std::move(fun), std::move(arg));
      }

      case ir::IrTypeCase::ArrayType: {
        if (!typ.array) return typ;
        ir::IrType elem = typ.array->elem_type ? predicate_impl(*typ.array->elem_type) : ir::IrType{};
        return ir::new_array_type(std::move(elem), typ.array->size);
      }

      case ir::IrTypeCase::ExistVarType:
        return typ;

      case ir::IrTypeCase::ForallType: {
        auto [new_ctx, tp, body_type] = context_.add_fresh_type(typ);
        context_ = new_ctx;
        tvars_.push_back(tp);
        return predicate_impl(body_type);
      }

      case ir::IrTypeCase::FunType: {
        if (!typ.fun) return typ;
        ir::IrType arg = typ.fun->arg ? predicate_impl(*typ.fun->arg) : ir::IrType{};
        ir::IrType ret = typ.fun->ret ? predicate_impl(*typ.fun->ret) : ir::IrType{};
        return ir::new_function_type(std::move(arg), std::move(ret));
      }

      case ir::IrTypeCase::LambdaType: {
        auto [new_ctx, tp, body_type] = context_.add_fresh_type(typ);
        context_ = new_ctx;
        tvars_.push_back(tp);
        return predicate_impl(body_type);
      }

      case ir::IrTypeCase::NameType:
        return typ;

      case ir::IrTypeCase::StructType: {
        if (!typ.struct_data) return typ;
        std::vector<ir::StructField> fields;
        fields.reserve(typ.struct_data->fields.size());
        for (const auto& f : typ.struct_data->fields) {
          ir::IrType ft = f.type ? predicate_impl(*f.type) : ir::IrType{};
          fields.push_back(ir::StructField{f.id, std::make_shared<ir::IrType>(std::move(ft))});
        }
        return ir::new_struct_type(std::move(fields));
      }

      case ir::IrTypeCase::TupleType: {
        if (!typ.tuple_data) return typ;
        std::vector<ir::IrType> elems;
        elems.reserve(typ.tuple_data->elems.size());
        for (const auto& e : typ.tuple_data->elems) {
          elems.push_back(predicate_impl(e));
        }
        return ir::new_tuple_type(std::move(elems));
      }

      case ir::IrTypeCase::VariantType: {
        if (!typ.variant_data) return typ;
        std::vector<ir::VariantTag> tags;
        tags.reserve(typ.variant_data->tags.size());
        for (const auto& tag : typ.variant_data->tags) {
          ir::IrType tt = tag.type ? predicate_impl(*tag.type) : ir::IrType{};
          tags.push_back(ir::VariantTag{tag.id, std::make_shared<ir::IrType>(std::move(tt))});
        }
        return ir::new_variant_type(std::move(tags));
      }

      case ir::IrTypeCase::VarType:
        return typ;
    }
    return typ;
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
