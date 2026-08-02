#pragma once

#include "bin/ir_base.h"
#include "bin/ir_decl.h"
#include "bin/ir_type.h"
#include "ts/bind.h"
#include "ts/context.h"
#include "ts/list.h"

#include <string>
#include <tuple>
#include <vector>

namespace ts {

struct TypeSubstitution {
  ir::IrType source;
  ir::IrType target;
};

class TypeVarRenamer {
 public:
  explicit TypeVarRenamer(Context context, List<TypeSubstitution> substitutions = {})
      : context_(std::move(context)), substitutions_(std::move(substitutions)) {}

  bool lookup_substitution(const std::string& source_tvar, ir::IrType& out_target) const {
    const auto* sub = substitutions_.find_if([&](const TypeSubstitution& s) {
      return s.source.is(ir::IrTypeCase::VarType) && s.source.var == source_tvar;
    });
    if (sub) {
      out_target = sub->target;
      return true;
    }
    return false;
  }

  template <typename Callback>
  ir::IrType with_type_substitution(
      const std::string& source_tvar,
      const ir::IrKind& source_kind,
      const std::vector<ir::IrType>& source_bounds,
      Callback&& callback) {
    Context orig_context = context_;
    List<TypeSubstitution> orig_substitutions = substitutions_;

    ir::IrType target_tvar = context_.gen_fresh_var_type();
    substitutions_ = substitutions_.add(TypeSubstitution{ir::new_var_type(source_tvar), target_tvar});

    std::vector<ir::IrType> renamed_bounds;
    renamed_bounds.reserve(source_bounds.size());
    for (const auto& b : source_bounds) {
      renamed_bounds.push_back(rename(b));
    }

    try {
      context_ = context_.add_bind(new_type_param_bind(target_tvar.var, source_kind, renamed_bounds));
    } catch (const std::exception& e) {
      if (!err_.empty()) err_ += "\n";
      err_ += e.what();
      return ir::IrType{};
    }

    ir::IrType result = callback(target_tvar.var);

    context_ = orig_context;
    substitutions_ = orig_substitutions;
    return result;
  }

  ir::IrType rename(const ir::IrType& typ) {
    if (!err_.empty()) {
      return ir::IrType{};
    }
    return rename_impl(typ);
  }

  const std::string& error() const { return err_; }
  const Context& context() const { return context_; }

 private:
  ir::IrType rename_impl(const ir::IrType& typ) {
    switch (typ.case_val) {
      case ir::IrTypeCase::AppType: {
        if (!typ.app) return typ;
        ir::IrType fun = typ.app->fun ? rename(*typ.app->fun) : ir::IrType{};
        ir::IrType arg = typ.app->arg ? rename(*typ.app->arg) : ir::IrType{};
        return ir::new_app_type(std::move(fun), std::move(arg));
      }
      case ir::IrTypeCase::ArrayType: {
        if (!typ.array) return typ;
        ir::IrType elem = typ.array->elem_type ? rename(*typ.array->elem_type) : ir::IrType{};
        return ir::new_array_type(std::move(elem), typ.array->size);
      }
      case ir::IrTypeCase::ExistVarType:
        return typ;
      case ir::IrTypeCase::ForallType: {
        if (!typ.forall) return typ;
        const auto& tp = typ.forall->type_param;
        return with_type_substitution(tp.var, tp.kind, tp.bounds, [&](const std::string& target_tvar) {
          std::vector<ir::IrType> bounds;
          bounds.reserve(tp.bounds.size());
          for (const auto& b : tp.bounds) {
            bounds.push_back(rename(b));
          }
          ir::TypeParam new_tp{target_tvar, tp.kind, std::move(bounds)};
          ir::IrType body = typ.forall->type ? rename(*typ.forall->type) : ir::IrType{};
          return ir::new_forall_type(std::move(new_tp), std::move(body));
        });
      }
      case ir::IrTypeCase::FunType: {
        if (!typ.fun) return typ;
        ir::IrType arg = typ.fun->arg ? rename(*typ.fun->arg) : ir::IrType{};
        ir::IrType ret = typ.fun->ret ? rename(*typ.fun->ret) : ir::IrType{};
        return ir::new_function_type(std::move(arg), std::move(ret));
      }
      case ir::IrTypeCase::LambdaType: {
        if (!typ.lambda) return typ;
        return with_type_substitution(typ.lambda->var, typ.lambda->kind, {}, [&](const std::string& target_tvar) {
          ir::IrType body = typ.lambda->type ? rename(*typ.lambda->type) : ir::IrType{};
          return ir::new_lambda_type(target_tvar, typ.lambda->kind, std::move(body));
        });
      }
      case ir::IrTypeCase::NameType:
        return typ;
      case ir::IrTypeCase::StructType: {
        if (!typ.struct_data) return typ;
        std::vector<ir::StructField> fields;
        fields.reserve(typ.struct_data->fields.size());
        for (const auto& f : typ.struct_data->fields) {
          ir::IrType ft = f.type ? rename(*f.type) : ir::IrType{};
          fields.push_back(ir::StructField{f.id, std::make_shared<ir::IrType>(std::move(ft))});
        }
        return ir::new_struct_type(std::move(fields));
      }
      case ir::IrTypeCase::TupleType: {
        if (!typ.tuple_data) return typ;
        std::vector<ir::IrType> elems;
        elems.reserve(typ.tuple_data->elems.size());
        for (const auto& e : typ.tuple_data->elems) {
          elems.push_back(rename(e));
        }
        return ir::new_tuple_type(std::move(elems));
      }
      case ir::IrTypeCase::VariantType: {
        if (!typ.variant_data) return typ;
        std::vector<ir::VariantTag> tags;
        tags.reserve(typ.variant_data->tags.size());
        for (const auto& tag : typ.variant_data->tags) {
          ir::IrType tt = tag.type ? rename(*tag.type) : ir::IrType{};
          tags.push_back(ir::VariantTag{tag.id, std::make_shared<ir::IrType>(std::move(tt))});
        }
        return ir::new_variant_type(std::move(tags));
      }
      case ir::IrTypeCase::VarType: {
        ir::IrType target;
        if (lookup_substitution(typ.var, target)) {
          return target;
        }
        return typ;
      }
    }
    return typ;
  }

  Context context_;
  List<TypeSubstitution> substitutions_;
  std::string err_;
};

inline std::tuple<Context, ir::IrType, std::string> rename_type_vars_with_substitutions(
    const Context& context,
    const ir::IrType& typ,
    const std::vector<TypeSubstitution>& substitutions) {
  List<TypeSubstitution> subs;
  for (const auto& s : substitutions) {
    subs = subs.add(s);
  }
  TypeVarRenamer renamer(context, std::move(subs));
  ir::IrType result = renamer.rename(typ);
  return {renamer.context(), std::move(result), renamer.error()};
}

inline std::tuple<Context, ir::IrType, std::string> rename_type_vars(
    const Context& context,
    const ir::IrType& typ) {
  return rename_type_vars_with_substitutions(context, typ, {});
}

} // namespace ts
