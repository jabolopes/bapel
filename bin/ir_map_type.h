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
  switch (typ.tag()) {
    case IrTypeData::Tag::app_t: {
      const auto& a = typ.data.as_app_t();
      IrType fun, arg;
      if (a.fun && !fn(*a.fun, fun, err)) return false;
      if (a.arg && !fn(*a.arg, arg, err)) return false;
      out_type = new_app_type(std::move(fun), std::move(arg));
      out_type.pos = typ.pos;
      return true;
    }
    case IrTypeData::Tag::array_t: {
      const auto& arr = typ.data.as_array_t();
      if (!arr.elem_type) { out_type = typ; return true; }
      IrType elem;
      if (!fn(*arr.elem_type, elem, err)) return false;
      out_type = new_array_type(std::move(elem), arr.size);
      out_type.pos = typ.pos;
      return true;
    }
    case IrTypeData::Tag::fun_t: {
      const auto& f = typ.data.as_fun_t();
      IrType arg, ret;
      if (f.arg && !fn(*f.arg, arg, err)) return false;
      if (f.ret && !fn(*f.ret, ret, err)) return false;
      out_type = new_function_type(std::move(arg), std::move(ret));
      out_type.pos = typ.pos;
      return true;
    }
    case IrTypeData::Tag::forall_t: {
      const auto& fa = typ.data.as_forall_t();
      TypeParam tp = fa.type_param;
      for (auto& b : tp.bounds) {
        IrType nb;
        if (!fn(b, nb, err)) return false;
        b = std::move(nb);
      }
      IrType body;
      if (fa.type && !fn(*fa.type, body, err)) return false;
      out_type = new_forall_type(std::move(tp), std::move(body));
      out_type.pos = typ.pos;
      return true;
    }
    case IrTypeData::Tag::lambda_t: {
      const auto& l = typ.data.as_lambda_t();
      if (!l.type) { out_type = typ; return true; }
      IrType body;
      if (!fn(*l.type, body, err)) return false;
      out_type = new_lambda_type(l.var, l.kind, std::move(body));
      out_type.pos = typ.pos;
      return true;
    }
    case IrTypeData::Tag::mu_t: {
      const auto& m = typ.data.as_mu_t();
      if (!m.body) { out_type = typ; return true; }
      IrType body;
      if (!fn(*m.body, body, err)) return false;
      out_type = new_mu_type(m.var, m.kind, std::move(body));
      out_type.pos = typ.pos;
      return true;
    }
    case IrTypeData::Tag::struct_t: {
      const auto& s = typ.data.as_struct_t();
      std::vector<StructField> fields;
      fields.reserve(s.fields.size());
      for (const auto& f : s.fields) {
        IrType ft;
        if (f.type && !fn(*f.type, ft, err)) return false;
        fields.push_back(StructField{f.id, std::make_shared<IrType>(std::move(ft))});
      }
      out_type = new_struct_type(std::move(fields));
      out_type.pos = typ.pos;
      return true;
    }
    case IrTypeData::Tag::tuple_t: {
      const auto& tup = typ.data.as_tuple_t();
      std::vector<IrType> elems;
      elems.reserve(tup.elems.size());
      for (const auto& e : tup.elems) {
        IrType ne;
        if (!fn(e, ne, err)) return false;
        elems.push_back(std::move(ne));
      }
      out_type = new_tuple_type(std::move(elems));
      out_type.pos = typ.pos;
      return true;
    }
    case IrTypeData::Tag::variant_t: {
      const auto& v = typ.data.as_variant_t();
      std::vector<VariantTag> tags;
      tags.reserve(v.tags.size());
      for (const auto& tag : v.tags) {
        IrType tt;
        if (tag.type && !fn(*tag.type, tt, err)) return false;
        tags.push_back(VariantTag{tag.id, std::make_shared<IrType>(std::move(tt))});
      }
      out_type = new_variant_type(std::move(tags));
      out_type.pos = typ.pos;
      return true;
    }
    case IrTypeData::Tag::name_t:
    case IrTypeData::Tag::var_t:
    case IrTypeData::Tag::exist_var_t:
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
