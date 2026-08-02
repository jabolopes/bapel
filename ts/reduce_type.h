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
      throw std::runtime_error("failed to reduce " + typ.to_string() + " because type is not wellformed: " + err);
    }
    ir::IrType reduced = reduce_impl(ctx, typ);
    reduced.pos = typ.pos;
    return reduced;
  }

 private:
  ir::IrType reduce_impl(const Context& ctx, const ir::IrType& typ) {
    switch (typ.case_val) {
      case ir::IrTypeCase::AppType: {
        if (!typ.app) return typ;
        ir::IrType fun = typ.app->fun ? reduce(ctx, *typ.app->fun) : ir::IrType{};
        ir::IrType arg = typ.app->arg ? reduce(ctx, *typ.app->arg) : ir::IrType{};

        if (fun.is(ir::IrTypeCase::LambdaType) && fun.lambda && fun.lambda->type) {
          return reduce(ctx, ir::substitute_type(*fun.lambda->type, ir::new_var_type(fun.lambda->var), arg));
        }
        return ir::new_app_type(std::move(fun), std::move(arg));
      }

      case ir::IrTypeCase::ArrayType: {
        if (!typ.array) return typ;
        ir::IrType elem = typ.array->elem_type ? reduce(ctx, *typ.array->elem_type) : ir::IrType{};
        return ir::new_array_type(std::move(elem), typ.array->size);
      }

      case ir::IrTypeCase::ForallType: {
        auto [new_ctx, tp, body_type] = ctx.add_fresh_type(typ);
        std::vector<ir::IrType> bounds;
        bounds.reserve(tp.bounds.size());
        for (const auto& b : tp.bounds) {
          bounds.push_back(reduce(ctx, b));
        }
        tp.bounds = std::move(bounds);
        return ir::new_forall_type(std::move(tp), reduce(new_ctx, body_type));
      }

      case ir::IrTypeCase::ExistVarType:
        return typ;

      case ir::IrTypeCase::FunType: {
        if (!typ.fun) return typ;
        ir::IrType arg = typ.fun->arg ? reduce(ctx, *typ.fun->arg) : ir::IrType{};
        ir::IrType ret = typ.fun->ret ? reduce(ctx, *typ.fun->ret) : ir::IrType{};
        return ir::new_function_type(std::move(arg), std::move(ret));
      }

      case ir::IrTypeCase::LambdaType: {
        auto [new_ctx, tp, body_type] = ctx.add_fresh_type(typ);
        return ir::new_lambda_type(tp.var, tp.kind, reduce(new_ctx, body_type));
      }

      case ir::IrTypeCase::NameType: {
        if (ctx.contains_alias_bind(typ.name)) {
          Binding bind;
          if (ctx.lookup_alias_bind(typ.name, bind) && bind.alias) {
            return reduce(ctx, bind.alias->type);
          }
        }
        return typ;
      }

      case ir::IrTypeCase::StructType: {
        if (!typ.struct_data) return typ;
        std::vector<ir::StructField> fields;
        fields.reserve(typ.struct_data->fields.size());
        for (const auto& f : typ.struct_data->fields) {
          ir::IrType ft = f.type ? reduce(ctx, *f.type) : ir::IrType{};
          fields.push_back(ir::StructField{f.id, std::make_shared<ir::IrType>(std::move(ft))});
        }
        return ir::new_struct_type(std::move(fields));
      }

      case ir::IrTypeCase::TupleType: {
        if (!typ.tuple_data) return typ;
        std::vector<ir::IrType> elems;
        elems.reserve(typ.tuple_data->elems.size());
        for (const auto& e : typ.tuple_data->elems) {
          elems.push_back(reduce(ctx, e));
        }
        return ir::new_tuple_type(std::move(elems));
      }

      case ir::IrTypeCase::VariantType: {
        if (!typ.variant_data) return typ;
        std::vector<ir::VariantTag> tags;
        tags.reserve(typ.variant_data->tags.size());
        for (const auto& tag : typ.variant_data->tags) {
          ir::IrType tt = tag.type ? reduce(ctx, *tag.type) : ir::IrType{};
          tags.push_back(ir::VariantTag{tag.id, std::make_shared<ir::IrType>(std::move(tt))});
        }
        return ir::new_variant_type(std::move(tags));
      }

      case ir::IrTypeCase::VarType:
        return typ;
    }
    return typ;
  }
};

inline ir::IrType reduce_type(const Context& ctx, const ir::IrType& typ) {
  TypeReducer reducer;
  return reducer.reduce(ctx, typ);
}

} // namespace ts
