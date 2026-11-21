#pragma once

#include <memory>
#include <string>
#include <utility>
#include <vector>

namespace ir {

// Applies `fn(const IrType& child, IrType& out_child, std::string& err)` to
// every direct child IrType of `typ`, returning a rebuilt IrType with `pos`
// preserved. Callers control recursion (e.g. intercepting specific cases or
// handling binder shadowing).
template <typename Fn>
inline bool map_type(const IrType& typ, Fn&& fn, IrType& out_type, std::string& err) {
  switch (typ.case_val) {
    case IrTypeCase::AppType: {
      if (!typ.app) { out_type = typ; return true; }
      IrType fun, arg;
      if (typ.app->fun && !fn(*typ.app->fun, fun, err)) return false;
      if (typ.app->arg && !fn(*typ.app->arg, arg, err)) return false;
      out_type = new_app_type(std::move(fun), std::move(arg));
      out_type.pos = typ.pos;
      return true;
    }
    case IrTypeCase::ArrayType: {
      if (!typ.array || !typ.array->elem_type) { out_type = typ; return true; }
      IrType elem;
      if (!fn(*typ.array->elem_type, elem, err)) return false;
      out_type = new_array_type(std::move(elem), typ.array->size);
      out_type.pos = typ.pos;
      return true;
    }
    case IrTypeCase::FunType: {
      if (!typ.fun) { out_type = typ; return true; }
      IrType arg, ret;
      if (typ.fun->arg && !fn(*typ.fun->arg, arg, err)) return false;
      if (typ.fun->ret && !fn(*typ.fun->ret, ret, err)) return false;
      out_type = new_function_type(std::move(arg), std::move(ret));
      out_type.pos = typ.pos;
      return true;
    }
    case IrTypeCase::ForallType: {
      if (!typ.forall) { out_type = typ; return true; }
      TypeParam tp = typ.forall->type_param;
      for (auto& b : tp.bounds) {
        IrType nb;
        if (!fn(b, nb, err)) return false;
        b = std::move(nb);
      }
      IrType body;
      if (typ.forall->type && !fn(*typ.forall->type, body, err)) return false;
      out_type = new_forall_type(std::move(tp), std::move(body));
      out_type.pos = typ.pos;
      return true;
    }
    case IrTypeCase::LambdaType: {
      if (!typ.lambda || !typ.lambda->type) { out_type = typ; return true; }
      IrType body;
      if (!fn(*typ.lambda->type, body, err)) return false;
      out_type = new_lambda_type(typ.lambda->var, typ.lambda->kind, std::move(body));
      out_type.pos = typ.pos;
      return true;
    }
    case IrTypeCase::MuType: {
      if (!typ.mu || !typ.mu->body) { out_type = typ; return true; }
      IrType body;
      if (!fn(*typ.mu->body, body, err)) return false;
      out_type = new_mu_type(typ.mu->var, typ.mu->kind, std::move(body));
      out_type.pos = typ.pos;
      return true;
    }
    case IrTypeCase::StructType: {
      if (!typ.struct_data) { out_type = typ; return true; }
      std::vector<StructField> fields;
      fields.reserve(typ.struct_data->fields.size());
      for (const auto& f : typ.struct_data->fields) {
        IrType ft;
        if (f.type && !fn(*f.type, ft, err)) return false;
        fields.push_back(StructField{f.id, std::make_shared<IrType>(std::move(ft))});
      }
      out_type = new_struct_type(std::move(fields));
      out_type.pos = typ.pos;
      return true;
    }
    case IrTypeCase::TupleType: {
      if (!typ.tuple_data) { out_type = typ; return true; }
      std::vector<IrType> elems;
      elems.reserve(typ.tuple_data->elems.size());
      for (const auto& e : typ.tuple_data->elems) {
        IrType ne;
        if (!fn(e, ne, err)) return false;
        elems.push_back(std::move(ne));
      }
      out_type = new_tuple_type(std::move(elems));
      out_type.pos = typ.pos;
      return true;
    }
    case IrTypeCase::VariantType: {
      if (!typ.variant_data) { out_type = typ; return true; }
      std::vector<VariantTag> tags;
      tags.reserve(typ.variant_data->tags.size());
      for (const auto& tag : typ.variant_data->tags) {
        IrType tt;
        if (tag.type && !fn(*tag.type, tt, err)) return false;
        tags.push_back(VariantTag{tag.id, std::make_shared<IrType>(std::move(tt))});
      }
      out_type = new_variant_type(std::move(tags));
      out_type.pos = typ.pos;
      return true;
    }
    case IrTypeCase::NameType:
    case IrTypeCase::VarType:
    case IrTypeCase::ExistVarType:
      out_type = typ;
      return true;
  }
  out_type = typ;
  return true;
}

// Infallible overload: `fn(const IrType& child) -> IrType`.
template <typename Fn>
inline IrType map_type(const IrType& typ, Fn&& fn) {
  IrType out;
  std::string err;
  map_type(typ, [&](const IrType& child, IrType& out_child, std::string&) {
    out_child = fn(child);
    return true;
  }, out, err);
  return out;
}

} // namespace ir
