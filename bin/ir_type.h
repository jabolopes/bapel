#pragma once

#include "ir_base.h"
#include <set>
#include <sstream>
#include <stdexcept>

namespace ir {

using IrTypeData = ::IrTypeData<::IrType>;
using StructField = ::NamedType<::IrType>;
using VariantTag = ::NamedType<::IrType>;
using ::FunctionArg;

inline std::string type_to_string(const IrType& typ);

inline IrType type_app_fun(const IrType& typ) {
  if (typ.data.is(IrTypeData::Tag::app_t)) {
    const auto& a = typ.data.as_app_t();
    return a.fun ? type_app_fun(*a.fun) : typ;
  }
  return typ;
}

inline void collect_app_args(const IrType& typ, std::vector<IrType>& acc) {
  if (typ.data.is(IrTypeData::Tag::app_t)) {
    const auto& a = typ.data.as_app_t();
    if (a.fun) {
      collect_app_args(*a.fun, acc);
    }
    if (a.arg) {
      acc.push_back(*a.arg);
    }
  }
}

inline std::vector<IrType> type_app_args(const IrType& typ) {
  std::vector<IrType> res;
  collect_app_args(typ, res);
  return res;
}

inline std::vector<StructField> type_fields(const IrType& typ) {
  if (typ.data.is(IrTypeData::Tag::struct_t)) {
    return typ.data.as_struct_t();
  }
  return {};
}

inline std::vector<IrType> type_elems(const IrType& typ) {
  if (typ.data.is(IrTypeData::Tag::tuple_t)) {
    return typ.data.as_tuple_t();
  }
  return {};
}

inline std::vector<VariantTag> type_tags(const IrType& typ) {
  if (typ.data.is(IrTypeData::Tag::variant_t)) {
    return typ.data.as_variant_t();
  }
  return {};
}

inline std::vector<std::string> type_field_ids(const IrType& typ) {
  std::vector<std::string> ids;
  for (const auto& f : type_fields(typ)) {
    ids.push_back(f.id);
  }
  return ids;
}

inline std::vector<std::string> type_tag_ids(const IrType& typ) {
  std::vector<std::string> ids;
  for (const auto& t : type_tags(typ)) {
    ids.push_back(t.id);
  }
  return ids;
}

inline std::vector<std::string> type_lambda_vars(const IrType& typ) {
  std::vector<std::string> res;
  if (typ.data.is(IrTypeData::Tag::lambda_t)) {
    const auto& l = typ.data.as_lambda_t();
    res.push_back(l.var);
    if (l.body && l.body->data.is(IrTypeData::Tag::lambda_t)) {
      auto next = type_lambda_vars(*l.body);
      res.insert(res.end(), next.begin(), next.end());
    }
  }
  return res;
}

inline IrType type_lambda_body(const IrType& typ) {
  if (typ.data.is(IrTypeData::Tag::lambda_t)) {
    const auto& l = typ.data.as_lambda_t();
    if (l.body && l.body->data.is(IrTypeData::Tag::lambda_t)) {
      return type_lambda_body(*l.body);
    }
    return l.body ? *l.body : typ;
  }
  return typ;
}

inline std::vector<TypeParam> type_forall_type_params(const IrType& typ) {
  std::vector<TypeParam> res;
  if (typ.data.is(IrTypeData::Tag::forall_t)) {
    const auto& fa = typ.data.as_forall_t();
    res.push_back(fa.param);
    if (fa.body && fa.body->data.is(IrTypeData::Tag::forall_t)) {
      auto next = type_forall_type_params(*fa.body);
      res.insert(res.end(), next.begin(), next.end());
    }
  }
  return res;
}

inline std::vector<std::string> type_forall_vars(const IrType& typ) {
  std::vector<std::string> res;
  if (typ.data.is(IrTypeData::Tag::forall_t)) {
    const auto& fa = typ.data.as_forall_t();
    res.push_back(fa.param.var);
    if (fa.body && fa.body->data.is(IrTypeData::Tag::forall_t)) {
      auto next = type_forall_vars(*fa.body);
      res.insert(res.end(), next.begin(), next.end());
    }
  }
  return res;
}

inline IrType type_forall_body(const IrType& typ) {
  if (typ.data.is(IrTypeData::Tag::forall_t)) {
    const auto& fa = typ.data.as_forall_t();
    if (fa.body && fa.body->data.is(IrTypeData::Tag::forall_t)) {
      return type_forall_body(*fa.body);
    }
    return fa.body ? *fa.body : typ;
  }
  return typ;
}

inline std::string type_trait_name(const IrType& typ) {
  if (typ.data.is(IrTypeData::Tag::name_t)) {
    return typ.data.as_name_t();
  }
  if (typ.data.is(IrTypeData::Tag::app_t)) {
    return type_trait_name(type_app_fun(typ));
  }
  return "";
}

inline bool type_field_by_label(const IrType& typ, const std::string& label, int& out_index, StructField& out_field) {
  if (typ.data.is(IrTypeData::Tag::struct_t)) {
    const auto& fields = typ.data.as_struct_t();
    try {
      int idx = std::stoi(label);
      if (idx >= 0 && idx < static_cast<int>(fields.size())) {
        out_index = idx;
        out_field = fields[idx];
        return true;
      }
    } catch (...) {}

    for (size_t i = 0; i < fields.size(); ++i) {
      if (fields[i].id == label) {
        out_index = static_cast<int>(i);
        out_field = fields[i];
        return true;
      }
    }
  }
  return false;
}

inline bool type_elem_by_label(const IrType& typ, const std::string& label, int& out_index, IrType& out_elem) {
  if (typ.data.is(IrTypeData::Tag::tuple_t)) {
    const auto& elems = typ.data.as_tuple_t();
    try {
      int idx = std::stoi(label);
      if (idx >= 0 && idx < static_cast<int>(elems.size())) {
        out_index = idx;
        out_elem = elems[idx];
        return true;
      }
    } catch (...) {}
  }
  return false;
}

inline bool type_tag_by_label(const IrType& typ, const std::string& label, int& out_index) {
  if (typ.data.is(IrTypeData::Tag::variant_t)) {
    const auto& tags = typ.data.as_variant_t();
    try {
      int idx = std::stoi(label);
      if (idx >= 0 && idx < static_cast<int>(tags.size())) {
        out_index = idx;
        return true;
      }
    } catch (...) {}

    for (size_t i = 0; i < tags.size(); ++i) {
      if (tags[i].id == label) {
        out_index = static_cast<int>(i);
        return true;
      }
    }
  }
  return false;
}

inline std::string struct_field_to_string(const StructField& f) {
  return f.id + ": " + (f.typ ? type_to_string(*f.typ) : "()");
}

inline std::string variant_tag_to_string(const VariantTag& t) {
  return t.id + " " + (t.typ ? type_to_string(*t.typ) : "()");
}

inline std::string function_arg_to_string(const FunctionArg& a) {
  return a.id + ": " + type_to_string(a.typ);
}

inline std::string type_param_to_string(const TypeParam& tp) {
  std::string s = "'" + tp.var;
  if (!tp.bounds.empty()) {
    s += ": ";
    Interleave(tp.bounds, [&]() { s += " + "; }, [&](int, const IrType& b) {
      s += type_to_string(b);
    });
  }
  return s;
}

inline std::string type_to_string(const IrType& typ) {
  switch (typ.data.tag()) {
    case IrTypeData::Tag::app_t: {
      const auto& a = typ.data.as_app_t();
      std::string lparen = "", rparen = "";
      if (a.arg && (a.arg->data.is(IrTypeData::Tag::app_t) || a.arg->data.is(IrTypeData::Tag::forall_t) ||
                    a.arg->data.is(IrTypeData::Tag::fun_t) || a.arg->data.is(IrTypeData::Tag::lambda_t))) {
        lparen = "("; rparen = ")";
      }
      return (a.fun ? type_to_string(*a.fun) : "") + " " + lparen + (a.arg ? type_to_string(*a.arg) : "") + rparen;
    }
    case IrTypeData::Tag::array_t: {
      const auto& arr = typ.data.as_array_t();
      return "[" + (arr.elem ? type_to_string(*arr.elem) : "") + ", " + std::to_string(arr.size) + "]";
    }
    case IrTypeData::Tag::exist_var_t:
      return "^" + std::to_string(typ.data.as_exist_var_t());
    case IrTypeData::Tag::forall_t: {
      std::ostringstream ss;
      ss << "forall [";
      auto tps = type_forall_type_params(typ);
      for (size_t i = 0; i < tps.size(); ++i) {
        if (i > 0) ss << ", ";
        ss << "'" << tps[i].var;
        if (!tps[i].bounds.empty()) {
          ss << ": ";
          for (size_t b = 0; b < tps[i].bounds.size(); ++b) {
            if (b > 0) ss << " + ";
            ss << type_to_string(tps[i].bounds[b]);
          }
        }
      }
      ss << "] " << type_to_string(type_forall_body(typ));
      return ss.str();
    }
    case IrTypeData::Tag::fun_t: {
      const auto& fn = typ.data.as_fun_t();
      std::string lparen = "", rparen = "";
      if (fn.arg && (fn.arg->data.is(IrTypeData::Tag::forall_t) || fn.arg->data.is(IrTypeData::Tag::fun_t) ||
                     fn.arg->data.is(IrTypeData::Tag::lambda_t))) {
        lparen = "("; rparen = ")";
      }
      return lparen + (fn.arg ? type_to_string(*fn.arg) : "") + rparen + " -> " +
             (fn.ret ? type_to_string(*fn.ret) : "");
    }
    case IrTypeData::Tag::lambda_t: {
      const auto& l = typ.data.as_lambda_t();
      return "fun (" + l.var + ") (" + (l.body ? type_to_string(*l.body) : "") + ")";
    }
    case IrTypeData::Tag::mu_t: {
      const auto& m = typ.data.as_mu_t();
      return "mu '" + m.var + ". " + (m.body ? type_to_string(*m.body) : "");
    }
    case IrTypeData::Tag::name_t:
      return typ.data.as_name_t();
    case IrTypeData::Tag::struct_t: {
      const auto& fields = typ.data.as_struct_t();
      std::ostringstream ss;
      ss << "struct{";
      for (size_t i = 0; i < fields.size(); ++i) {
        if (i > 0) ss << ", ";
        ss << struct_field_to_string(fields[i]);
      }
      ss << "}";
      return ss.str();
    }
    case IrTypeData::Tag::tuple_t: {
      const auto& elems = typ.data.as_tuple_t();
      std::ostringstream ss;
      ss << "(";
      for (size_t i = 0; i < elems.size(); ++i) {
        if (i > 0) ss << ", ";
        ss << type_to_string(elems[i]);
      }
      ss << ")";
      return ss.str();
    }
    case IrTypeData::Tag::variant_t: {
      const auto& tags = typ.data.as_variant_t();
      std::ostringstream ss;
      ss << "variant{";
      for (size_t i = 0; i < tags.size(); ++i) {
        if (i > 0) ss << ", ";
        ss << variant_tag_to_string(tags[i]);
      }
      ss << "}";
      return ss.str();
    }
    case IrTypeData::Tag::var_t: {
      const auto& v = typ.data.as_var_t();
      return v.empty() || v[0] == '\'' ? v : ("'" + v);
    }
  }
  return "";
}

inline IrType new_name_type(const std::string& name) {
  return IrType{IrTypeData::mk_name_t(name), Pos{}};
}

inline IrType new_var_type(const std::string& var) {
  return IrType{IrTypeData::mk_var_t(var), Pos{}};
}

inline IrType new_exist_var_type(int64_t evar) {
  return IrType{IrTypeData::mk_exist_var_t(evar), Pos{}};
}

inline IrType new_app_type(IrType fun, IrType arg) {
  return IrType{
      IrTypeData::mk_app_t({
          std::make_shared<IrType>(std::move(fun)),
          std::make_shared<IrType>(std::move(arg)),
      }),
      Pos{},
  };
}

inline IrType new_array_type(IrType elem, int64_t size) {
  return IrType{
      IrTypeData::mk_array_t({
          std::make_shared<IrType>(std::move(elem)),
          size,
      }),
      Pos{},
  };
}

inline IrType new_function_type(IrType arg, IrType ret) {
  return IrType{
      IrTypeData::mk_fun_t({
          std::make_shared<IrType>(std::move(arg)),
          std::make_shared<IrType>(std::move(ret)),
      }),
      Pos{},
  };
}

inline IrType new_tuple_type(std::vector<IrType> elems) {
  if (elems.size() == 1) {
    return elems[0];
  }
  return IrType{IrTypeData::mk_tuple_t(std::move(elems)), Pos{}};
}

inline IrType new_struct_type(std::vector<StructField> fields) {
  return IrType{IrTypeData::mk_struct_t(std::move(fields)), Pos{}};
}

inline IrType new_variant_type(std::vector<VariantTag> tags) {
  return IrType{IrTypeData::mk_variant_t(std::move(tags)), Pos{}};
}

inline IrType new_forall_type(TypeParam tp, IrType body) {
  return IrType{
      IrTypeData::mk_forall_t({
          std::move(tp),
          std::make_shared<IrType>(std::move(body)),
      }),
      Pos{},
  };
}

inline IrType new_lambda_type(std::string var, IrKind kind, IrType body) {
  return IrType{
      IrTypeData::mk_lambda_t({
          std::move(var),
          std::move(kind),
          std::make_shared<IrType>(std::move(body)),
      }),
      Pos{},
  };
}

inline IrType new_mu_type(std::string var, IrKind kind, IrType body) {
  return IrType{
      IrTypeData::mk_mu_t({
          std::move(var),
          std::move(kind),
          std::make_shared<IrType>(std::move(body)),
      }),
      Pos{},
  };
}

inline IrType apply_types(IrType fun, const std::vector<IrType>& args) {
  for (const auto& arg : args) {
    fun = new_app_type(std::move(fun), arg);
  }
  return fun;
}

inline IrType forall_vars(const std::vector<TypeParam>& params, IrType body) {
  for (auto it = params.rbegin(); it != params.rend(); ++it) {
    body = new_forall_type(*it, std::move(body));
  }
  return body;
}

inline void collect_free_vars(const IrType& typ, std::vector<std::string>& bound_vars, std::vector<std::string>& free_vars) {
  switch (typ.data.tag()) {
    case IrTypeData::Tag::app_t: {
      const auto& a = typ.data.as_app_t();
      if (a.fun) collect_free_vars(*a.fun, bound_vars, free_vars);
      if (a.arg) collect_free_vars(*a.arg, bound_vars, free_vars);
      break;
    }
    case IrTypeData::Tag::array_t: {
      const auto& arr = typ.data.as_array_t();
      if (arr.elem) collect_free_vars(*arr.elem, bound_vars, free_vars);
      break;
    }
    case IrTypeData::Tag::forall_t: {
      const auto& fa = typ.data.as_forall_t();
      for (const auto& b : fa.param.bounds) {
        collect_free_vars(b, bound_vars, free_vars);
      }
      bound_vars.push_back(fa.param.var);
      if (fa.body) collect_free_vars(*fa.body, bound_vars, free_vars);
      bound_vars.pop_back();
      break;
    }
    case IrTypeData::Tag::fun_t: {
      const auto& f = typ.data.as_fun_t();
      if (f.arg) collect_free_vars(*f.arg, bound_vars, free_vars);
      if (f.ret) collect_free_vars(*f.ret, bound_vars, free_vars);
      break;
    }
    case IrTypeData::Tag::lambda_t: {
      const auto& l = typ.data.as_lambda_t();
      bound_vars.push_back(l.var);
      if (l.body) collect_free_vars(*l.body, bound_vars, free_vars);
      bound_vars.pop_back();
      break;
    }
    case IrTypeData::Tag::mu_t: {
      const auto& m = typ.data.as_mu_t();
      bound_vars.push_back(m.var);
      if (m.body) collect_free_vars(*m.body, bound_vars, free_vars);
      bound_vars.pop_back();
      break;
    }
    case IrTypeData::Tag::struct_t:
      for (const auto& f : typ.data.as_struct_t()) {
        if (f.typ) collect_free_vars(*f.typ, bound_vars, free_vars);
      }
      break;
    case IrTypeData::Tag::tuple_t:
      for (const auto& elem : typ.data.as_tuple_t()) {
        collect_free_vars(elem, bound_vars, free_vars);
      }
      break;
    case IrTypeData::Tag::variant_t:
      for (const auto& tag : typ.data.as_variant_t()) {
        if (tag.typ) collect_free_vars(*tag.typ, bound_vars, free_vars);
      }
      break;
    case IrTypeData::Tag::var_t: {
      const auto& v = typ.data.as_var_t();
      if (std::find(bound_vars.begin(), bound_vars.end(), v) == bound_vars.end()) {
        if (std::find(free_vars.begin(), free_vars.end(), v) == free_vars.end()) {
          free_vars.push_back(v);
        }
      }
      break;
    }
    default:
      break;
  }
}

inline std::vector<TypeParam> get_free_type_vars(const IrType& typ) {
  std::vector<std::string> bound_vars;
  std::vector<std::string> free_vars;
  collect_free_vars(typ, bound_vars, free_vars);
  std::sort(free_vars.begin(), free_vars.end());
  std::vector<TypeParam> result;
  result.reserve(free_vars.size());
  for (auto& v : free_vars) {
    result.push_back(TypeParam{std::move(v), new_type_kind(), {}});
  }
  return result;
}

inline IrType quantify_type(IrType typ) {
  auto free_vars = get_free_type_vars(typ);
  return forall_vars(free_vars, std::move(typ));
}

inline IrType lambda_vars(const std::vector<TypeParam>& params, IrType body) {
  for (auto it = params.rbegin(); it != params.rend(); ++it) {
    body = new_lambda_type(it->var, it->kind, std::move(body));
  }
  return body;
}

inline bool equals_named_type(const StructField& f1, const StructField& f2);

inline bool equals_type(const IrType& t1, const IrType& t2) {
  if (t1.data.tag() != t2.data.tag()) return false;
  switch (t1.data.tag()) {
    case IrTypeData::Tag::app_t: {
      const auto& a1 = t1.data.as_app_t();
      const auto& a2 = t2.data.as_app_t();
      if (!a1.fun || !a2.fun || !a1.arg || !a2.arg) return false;
      return equals_type(*a1.fun, *a2.fun) && equals_type(*a1.arg, *a2.arg);
    }
    case IrTypeData::Tag::array_t: {
      const auto& a1 = t1.data.as_array_t();
      const auto& a2 = t2.data.as_array_t();
      if (a1.size != a2.size) return false;
      if (!a1.elem || !a2.elem) return a1.elem == a2.elem;
      return equals_type(*a1.elem, *a2.elem);
    }
    case IrTypeData::Tag::exist_var_t:
      return t1.data.as_exist_var_t() == t2.data.as_exist_var_t();
    case IrTypeData::Tag::forall_t: {
      const auto& f1 = t1.data.as_forall_t();
      const auto& f2 = t2.data.as_forall_t();
      if (f1.param.var != f2.param.var) return false;
      if (!equals_kind(f1.param.kind, f2.param.kind)) return false;
      if (f1.param.bounds.size() != f2.param.bounds.size()) return false;
      for (size_t i = 0; i < f1.param.bounds.size(); ++i) {
        if (!equals_type(f1.param.bounds[i], f2.param.bounds[i])) return false;
      }
      if (!f1.body || !f2.body) return f1.body == f2.body;
      return equals_type(*f1.body, *f2.body);
    }
    case IrTypeData::Tag::fun_t: {
      const auto& f1 = t1.data.as_fun_t();
      const auto& f2 = t2.data.as_fun_t();
      if (!f1.arg || !f2.arg || !f1.ret || !f2.ret) return false;
      return equals_type(*f1.arg, *f2.arg) && equals_type(*f1.ret, *f2.ret);
    }
    case IrTypeData::Tag::lambda_t: {
      const auto& l1 = t1.data.as_lambda_t();
      const auto& l2 = t2.data.as_lambda_t();
      if (l1.var != l2.var) return false;
      if (!equals_kind(l1.kind, l2.kind)) return false;
      if (!l1.body || !l2.body) return l1.body == l2.body;
      return equals_type(*l1.body, *l2.body);
    }
    case IrTypeData::Tag::mu_t: {
      const auto& m1 = t1.data.as_mu_t();
      const auto& m2 = t2.data.as_mu_t();
      if (m1.var != m2.var) return false;
      if (!equals_kind(m1.kind, m2.kind)) return false;
      if (!m1.body || !m2.body) return m1.body == m2.body;
      return equals_type(*m1.body, *m2.body);
    }
    case IrTypeData::Tag::name_t:
      return t1.data.as_name_t() == t2.data.as_name_t();
    case IrTypeData::Tag::struct_t: {
      const auto& s1 = t1.data.as_struct_t();
      const auto& s2 = t2.data.as_struct_t();
      if (s1.size() != s2.size()) return false;
      for (size_t i = 0; i < s1.size(); ++i) {
        if (!equals_named_type(s1[i], s2[i])) return false;
      }
      return true;
    }
    case IrTypeData::Tag::tuple_t: {
      const auto& e1 = t1.data.as_tuple_t();
      const auto& e2 = t2.data.as_tuple_t();
      if (e1.size() != e2.size()) return false;
      for (size_t i = 0; i < e1.size(); ++i) {
        if (!equals_type(e1[i], e2[i])) return false;
      }
      return true;
    }
    case IrTypeData::Tag::variant_t: {
      const auto& v1 = t1.data.as_variant_t();
      const auto& v2 = t2.data.as_variant_t();
      if (v1.size() != v2.size()) return false;
      for (size_t i = 0; i < v1.size(); ++i) {
        if (!equals_named_type(v1[i], v2[i])) return false;
      }
      return true;
    }
    case IrTypeData::Tag::var_t:
      return t1.data.as_var_t() == t2.data.as_var_t();
  }
  return false;
}

inline bool equals_named_type(const StructField& f1, const StructField& f2) {
  if (f1.id != f2.id) return false;
  if (!f1.typ || !f2.typ) return f1.typ == f2.typ;
  return equals_type(*f1.typ, *f2.typ);
}

inline IrType operator_type(const std::string& id) {
  auto comparison = new_forall_type(
      TypeParam{"a", new_type_kind(), {}},
      new_function_type(new_tuple_type({new_var_type("a"), new_var_type("a")}), new_name_type("bool")));

  auto additive = new_forall_type(
      TypeParam{"a", new_type_kind(), {}},
      new_function_type(new_tuple_type({new_var_type("a"), new_var_type("a")}), new_var_type("a")));

  auto logical_unary = new_function_type(new_name_type("bool"), new_name_type("bool"));

  auto logical_binary = new_function_type(new_tuple_type({new_name_type("bool"), new_name_type("bool")}), new_name_type("bool"));

  if (id == "||" || id == "&&") return logical_binary;
  if (id == "!=" || id == "==" || id == ">" || id == ">=" || id == "<" || id == "<=") return comparison;
  if (id == "+" || id == "-" || id == "*" || id == "/") return additive;
  if (id == "!") return logical_unary;
  return new_tuple_type({});
}

} // namespace ir

inline bool operator==(const ::IrType& t1, const ::IrType& t2) {
  return ir::equals_type(t1, t2);
}

inline bool operator!=(const ::IrType& t1, const ::IrType& t2) {
  return !(t1 == t2);
}

inline bool operator==(const ::NamedType<::IrType>& f1, const ::NamedType<::IrType>& f2) {
  return ir::equals_named_type(f1, f2);
}

inline bool operator!=(const ::NamedType<::IrType>& f1, const ::NamedType<::IrType>& f2) {
  return !(f1 == f2);
}

inline bool operator==(const ::TypeParam<::IrType>& tp1, const ::TypeParam<::IrType>& tp2) {
  if (tp1.var != tp2.var) return false;
  if (!ir::equals_kind(tp1.kind, tp2.kind)) return false;
  if (tp1.bounds.size() != tp2.bounds.size()) return false;
  for (size_t i = 0; i < tp1.bounds.size(); ++i) {
    if (!ir::equals_type(tp1.bounds[i], tp2.bounds[i])) return false;
  }
  return true;
}

inline bool operator!=(const ::TypeParam<::IrType>& tp1, const ::TypeParam<::IrType>& tp2) {
  return !(tp1 == tp2);
}

#include "ir_map_type.h"
#include "ir_substitute_type.h"
