#pragma once

#include "ir_base.h"
#include <set>
#include <sstream>
#include <stdexcept>

namespace ir {

struct StructField {
  std::string id;
  std::shared_ptr<IrType> type;

  std::string to_string() const;
};

struct VariantTag {
  std::string id;
  std::shared_ptr<IrType> type;

  std::string to_string() const;
};

struct AppTypeData {
  std::shared_ptr<IrType> fun;
  std::shared_ptr<IrType> arg;
};

struct ArrayTypeData {
  std::shared_ptr<IrType> elem_type;
  int64_t size = 0;
};

struct ForallTypeData {
  TypeParam type_param;
  std::shared_ptr<IrType> type;
};

struct FunctionTypeData {
  std::shared_ptr<IrType> arg;
  std::shared_ptr<IrType> ret;
};

struct LambdaTypeData {
  std::string var;
  IrKind kind;
  std::shared_ptr<IrType> type;
};

struct MuTypeData {
  std::string var;
  IrKind kind;
  std::shared_ptr<IrType> body;
};

struct StructTypeData {
  std::vector<StructField> fields;
};

struct TupleTypeData {
  std::vector<IrType> elems;

  TupleTypeData();
  explicit TupleTypeData(std::vector<IrType> e);
  TupleTypeData(const TupleTypeData&);
  TupleTypeData(TupleTypeData&&) noexcept;
  TupleTypeData& operator=(const TupleTypeData&);
  TupleTypeData& operator=(TupleTypeData&&) noexcept;
  ~TupleTypeData();
};

struct VariantTypeData {
  std::vector<VariantTag> tags;
};

BAPEL_VARIANT(IrTypeData,
  (app_t, (AppTypeData)),
  (array_t, (ArrayTypeData)),
  (exist_var_t, (int64_t)),
  (forall_t, (ForallTypeData)),
  (fun_t, (FunctionTypeData)),
  (lambda_t, (LambdaTypeData)),
  (mu_t, (MuTypeData)),
  (name_t, (std::string)),
  (struct_t, (StructTypeData)),
  (tuple_t, (TupleTypeData)),
  (variant_t, (VariantTypeData)),
  (var_t, (std::string))
);

struct IrType {
  IrTypeData data{IrTypeData::mk_name_t("")};
  Pos pos;

  IrTypeData::Tag tag() const { return data.tag(); }
  bool is(IrTypeData::Tag t) const { return data.is(t); }

  IrType app_fun() const {
    if (is(IrTypeData::Tag::app_t)) {
      const auto& a = data.as_app_t();
      return a.fun ? a.fun->app_fun() : *this;
    }
    return *this;
  }

  std::vector<IrType> app_args() const {
    std::vector<IrType> res;
    collect_app_args(res);
    return res;
  }

  void collect_app_args(std::vector<IrType>& acc) const {
    if (is(IrTypeData::Tag::app_t)) {
      const auto& a = data.as_app_t();
      if (a.fun) {
        a.fun->collect_app_args(acc);
      }
      if (a.arg) {
        acc.push_back(*a.arg);
      }
    }
  }

  std::vector<StructField> fields() const {
    if (is(IrTypeData::Tag::struct_t)) {
      return data.as_struct_t().fields;
    }
    return {};
  }

  std::vector<IrType> elems() const;

  std::vector<VariantTag> tags() const {
    if (is(IrTypeData::Tag::variant_t)) {
      return data.as_variant_t().tags;
    }
    return {};
  }

  std::vector<std::string> field_ids() const {
    std::vector<std::string> ids;
    for (const auto& f : fields()) {
      ids.push_back(f.id);
    }
    return ids;
  }

  std::vector<std::string> tag_ids() const {
    std::vector<std::string> ids;
    for (const auto& t : tags()) {
      ids.push_back(t.id);
    }
    return ids;
  }

  std::vector<std::string> lambda_vars() const {
    std::vector<std::string> res;
    if (is(IrTypeData::Tag::lambda_t)) {
      const auto& l = data.as_lambda_t();
      res.push_back(l.var);
      if (l.type && l.type->is(IrTypeData::Tag::lambda_t)) {
        auto next = l.type->lambda_vars();
        res.insert(res.end(), next.begin(), next.end());
      }
    }
    return res;
  }

  IrType lambda_body() const {
    if (is(IrTypeData::Tag::lambda_t)) {
      const auto& l = data.as_lambda_t();
      if (l.type && l.type->is(IrTypeData::Tag::lambda_t)) {
        return l.type->lambda_body();
      }
      return l.type ? *l.type : *this;
    }
    return *this;
  }

  std::vector<TypeParam> forall_type_params() const {
    std::vector<TypeParam> res;
    if (is(IrTypeData::Tag::forall_t)) {
      const auto& fa = data.as_forall_t();
      res.push_back(fa.type_param);
      if (fa.type && fa.type->is(IrTypeData::Tag::forall_t)) {
        auto next = fa.type->forall_type_params();
        res.insert(res.end(), next.begin(), next.end());
      }
    }
    return res;
  }

  std::vector<std::string> forall_vars() const {
    std::vector<std::string> res;
    if (is(IrTypeData::Tag::forall_t)) {
      const auto& fa = data.as_forall_t();
      res.push_back(fa.type_param.var);
      if (fa.type && fa.type->is(IrTypeData::Tag::forall_t)) {
        auto next = fa.type->forall_vars();
        res.insert(res.end(), next.begin(), next.end());
      }
    }
    return res;
  }

  IrType forall_body() const {
    if (is(IrTypeData::Tag::forall_t)) {
      const auto& fa = data.as_forall_t();
      if (fa.type && fa.type->is(IrTypeData::Tag::forall_t)) {
        return fa.type->forall_body();
      }
      return fa.type ? *fa.type : *this;
    }
    return *this;
  }

  std::string trait_name() const {
    if (is(IrTypeData::Tag::name_t)) {
      return data.as_name_t();
    }
    if (is(IrTypeData::Tag::app_t)) {
      return app_fun().trait_name();
    }
    return "";
  }

  bool field_by_label(const std::string& label, int& out_index, StructField& out_field) const {
    if (is(IrTypeData::Tag::struct_t)) {
      const auto& fields = data.as_struct_t().fields;
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

  bool elem_by_label(const std::string& label, int& out_index, IrType& out_elem) const {
    if (is(IrTypeData::Tag::tuple_t)) {
      const auto& elems = data.as_tuple_t().elems;
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

  bool tag_by_label(const std::string& label, int& out_index) const {
    if (is(IrTypeData::Tag::variant_t)) {
      const auto& tags = data.as_variant_t().tags;
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

  std::string to_string() const;
};

inline std::string StructField::to_string() const {
  return id + ": " + (type ? type->to_string() : "()");
}

inline TypeParam::TypeParam() = default;
inline TypeParam::TypeParam(std::string v, IrKind k)
    : var(std::move(v)), kind(std::move(k)) {}
inline TypeParam::TypeParam(std::string v, IrKind k, std::vector<IrType> b)
    : var(std::move(v)), kind(std::move(k)), bounds(std::move(b)) {}
inline TypeParam::TypeParam(const TypeParam&) = default;
inline TypeParam::TypeParam(TypeParam&&) noexcept = default;
inline TypeParam& TypeParam::operator=(const TypeParam&) = default;
inline TypeParam& TypeParam::operator=(TypeParam&&) noexcept = default;
inline TypeParam::~TypeParam() = default;

inline TupleTypeData::TupleTypeData() = default;
inline TupleTypeData::TupleTypeData(std::vector<IrType> e) : elems(std::move(e)) {}
inline TupleTypeData::TupleTypeData(const TupleTypeData&) = default;
inline TupleTypeData::TupleTypeData(TupleTypeData&&) noexcept = default;
inline TupleTypeData& TupleTypeData::operator=(const TupleTypeData&) = default;
inline TupleTypeData& TupleTypeData::operator=(TupleTypeData&&) noexcept = default;
inline TupleTypeData::~TupleTypeData() = default;

inline std::vector<IrType> IrType::elems() const {
  if (is(IrTypeData::Tag::tuple_t)) {
    return data.as_tuple_t().elems;
  }
  return {};
}

inline std::string VariantTag::to_string() const {
  return id + " " + (type ? type->to_string() : "()");
}

inline std::string IrType::to_string() const {
  switch (data.tag()) {
    case IrTypeData::Tag::app_t: {
      const auto& a = data.as_app_t();
      std::string lparen = "", rparen = "";
      if (a.arg && (a.arg->is(IrTypeData::Tag::app_t) || a.arg->is(IrTypeData::Tag::forall_t) ||
                    a.arg->is(IrTypeData::Tag::fun_t) || a.arg->is(IrTypeData::Tag::lambda_t))) {
        lparen = "("; rparen = ")";
      }
      return (a.fun ? a.fun->to_string() : "") + " " + lparen + (a.arg ? a.arg->to_string() : "") + rparen;
    }
    case IrTypeData::Tag::array_t: {
      const auto& arr = data.as_array_t();
      return "[" + (arr.elem_type ? arr.elem_type->to_string() : "") + ", " + std::to_string(arr.size) + "]";
    }
    case IrTypeData::Tag::exist_var_t:
      return "^" + std::to_string(data.as_exist_var_t());
    case IrTypeData::Tag::forall_t: {
      std::ostringstream ss;
      ss << "forall [";
      auto tps = forall_type_params();
      for (size_t i = 0; i < tps.size(); ++i) {
        if (i > 0) ss << ", ";
        ss << "'" << tps[i].var;
        if (!tps[i].bounds.empty()) {
          ss << ": ";
          for (size_t b = 0; b < tps[i].bounds.size(); ++b) {
            if (b > 0) ss << " + ";
            ss << tps[i].bounds[b].to_string();
          }
        }
      }
      ss << "] " << forall_body().to_string();
      return ss.str();
    }
    case IrTypeData::Tag::fun_t: {
      const auto& f = data.as_fun_t();
      return (f.arg ? f.arg->to_string() : "") + " -> " + (f.ret ? f.ret->to_string() : "");
    }
    case IrTypeData::Tag::lambda_t: {
      const auto& l = data.as_lambda_t();
      return "fun (" + l.var + ") (" + (l.type ? l.type->to_string() : "") + ")";
    }
    case IrTypeData::Tag::mu_t: {
      const auto& m = data.as_mu_t();
      return "mu '" + m.var + ". " + (m.body ? m.body->to_string() : "");
    }
    case IrTypeData::Tag::name_t:
      return data.as_name_t();
    case IrTypeData::Tag::struct_t: {
      const auto& fields = data.as_struct_t().fields;
      std::string s = "struct{";
      for (size_t i = 0; i < fields.size(); ++i) {
        if (i > 0) s += ", ";
        s += fields[i].to_string();
      }
      return s + "}";
    }
    case IrTypeData::Tag::tuple_t: {
      const auto& elems = data.as_tuple_t().elems;
      std::string s = "(";
      for (size_t i = 0; i < elems.size(); ++i) {
        if (i > 0) s += ", ";
        s += elems[i].to_string();
      }
      return s + ")";
    }
    case IrTypeData::Tag::variant_t: {
      const auto& tags = data.as_variant_t().tags;
      std::string s = "variant{";
      for (size_t i = 0; i < tags.size(); ++i) {
        if (i > 0) s += ", ";
        s += tags[i].to_string();
      }
      return s + "}";
    }
    case IrTypeData::Tag::var_t: {
      const auto& v = data.as_var_t();
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
      IrTypeData::mk_app_t(AppTypeData{
          std::make_shared<IrType>(std::move(fun)),
          std::make_shared<IrType>(std::move(arg)),
      }),
      Pos{},
  };
}

inline IrType new_array_type(IrType elem, int64_t size) {
  return IrType{
      IrTypeData::mk_array_t(ArrayTypeData{
          std::make_shared<IrType>(std::move(elem)),
          size,
      }),
      Pos{},
  };
}

inline IrType new_function_type(IrType arg, IrType ret) {
  return IrType{
      IrTypeData::mk_fun_t(FunctionTypeData{
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
  return IrType{IrTypeData::mk_tuple_t(TupleTypeData{std::move(elems)}), Pos{}};
}

inline IrType new_struct_type(std::vector<StructField> fields) {
  return IrType{IrTypeData::mk_struct_t(StructTypeData{std::move(fields)}), Pos{}};
}

inline IrType new_variant_type(std::vector<VariantTag> tags) {
  return IrType{IrTypeData::mk_variant_t(VariantTypeData{std::move(tags)}), Pos{}};
}

inline IrType new_forall_type(TypeParam tp, IrType body) {
  return IrType{
      IrTypeData::mk_forall_t(ForallTypeData{
          std::move(tp),
          std::make_shared<IrType>(std::move(body)),
      }),
      Pos{},
  };
}

inline IrType new_lambda_type(std::string var, IrKind kind, IrType body) {
  return IrType{
      IrTypeData::mk_lambda_t(LambdaTypeData{
          std::move(var),
          std::move(kind),
          std::make_shared<IrType>(std::move(body)),
      }),
      Pos{},
  };
}

inline IrType new_mu_type(std::string var, IrKind kind, IrType body) {
  return IrType{
      IrTypeData::mk_mu_t(MuTypeData{
          std::move(var),
          std::move(kind),
          std::make_shared<IrType>(std::move(body)),
      }),
      Pos{},
  };
}

inline IrType forall_vars(const std::vector<TypeParam>& params, IrType body) {
  for (auto it = params.rbegin(); it != params.rend(); ++it) {
    body = new_forall_type(*it, std::move(body));
  }
  return body;
}

inline void collect_free_vars(const IrType& typ, std::vector<std::string>& bound_vars, std::vector<std::string>& free_vars) {
  switch (typ.tag()) {
    case IrTypeData::Tag::app_t: {
      const auto& a = typ.data.as_app_t();
      if (a.fun) collect_free_vars(*a.fun, bound_vars, free_vars);
      if (a.arg) collect_free_vars(*a.arg, bound_vars, free_vars);
      break;
    }
    case IrTypeData::Tag::array_t: {
      const auto& arr = typ.data.as_array_t();
      if (arr.elem_type) collect_free_vars(*arr.elem_type, bound_vars, free_vars);
      break;
    }
    case IrTypeData::Tag::forall_t: {
      const auto& fa = typ.data.as_forall_t();
      for (const auto& b : fa.type_param.bounds) {
        collect_free_vars(b, bound_vars, free_vars);
      }
      bound_vars.push_back(fa.type_param.var);
      if (fa.type) collect_free_vars(*fa.type, bound_vars, free_vars);
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
      if (l.type) collect_free_vars(*l.type, bound_vars, free_vars);
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
      for (const auto& f : typ.data.as_struct_t().fields) {
        if (f.type) collect_free_vars(*f.type, bound_vars, free_vars);
      }
      break;
    case IrTypeData::Tag::tuple_t:
      for (const auto& elem : typ.data.as_tuple_t().elems) {
        collect_free_vars(elem, bound_vars, free_vars);
      }
      break;
    case IrTypeData::Tag::variant_t:
      for (const auto& tag : typ.data.as_variant_t().tags) {
        if (tag.type) collect_free_vars(*tag.type, bound_vars, free_vars);
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

struct FunctionArg {
  std::string id;
  IrType type;

  std::string to_string() const {
    return id + ": " + type.to_string();
  }
};

inline std::string TypeParam::to_string() const {
  std::string s = "'" + var;
  if (!bounds.empty()) {
    s += ": ";
    Interleave(bounds, [&]() { s += " + "; }, [&](int, const IrType& b) {
      s += b.to_string();
    });
  }
  return s;
}

inline bool equals_struct_field(const StructField& f1, const StructField& f2);
inline bool equals_variant_tag(const VariantTag& t1, const VariantTag& t2);

inline bool equals_type(const IrType& t1, const IrType& t2) {
  if (t1.tag() != t2.tag()) return false;
  switch (t1.tag()) {
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
      if (!a1.elem_type || !a2.elem_type) return a1.elem_type == a2.elem_type;
      return equals_type(*a1.elem_type, *a2.elem_type);
    }
    case IrTypeData::Tag::exist_var_t:
      return t1.data.as_exist_var_t() == t2.data.as_exist_var_t();
    case IrTypeData::Tag::forall_t: {
      const auto& f1 = t1.data.as_forall_t();
      const auto& f2 = t2.data.as_forall_t();
      if (f1.type_param.var != f2.type_param.var) return false;
      if (!equals_kind(f1.type_param.kind, f2.type_param.kind)) return false;
      if (f1.type_param.bounds.size() != f2.type_param.bounds.size()) return false;
      for (size_t i = 0; i < f1.type_param.bounds.size(); ++i) {
        if (!equals_type(f1.type_param.bounds[i], f2.type_param.bounds[i])) return false;
      }
      if (!f1.type || !f2.type) return f1.type == f2.type;
      return equals_type(*f1.type, *f2.type);
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
      if (!l1.type || !l2.type) return l1.type == l2.type;
      return equals_type(*l1.type, *l2.type);
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
      const auto& s1 = t1.data.as_struct_t().fields;
      const auto& s2 = t2.data.as_struct_t().fields;
      if (s1.size() != s2.size()) return false;
      for (size_t i = 0; i < s1.size(); ++i) {
        if (!equals_struct_field(s1[i], s2[i])) return false;
      }
      return true;
    }
    case IrTypeData::Tag::tuple_t: {
      const auto& e1 = t1.data.as_tuple_t().elems;
      const auto& e2 = t2.data.as_tuple_t().elems;
      if (e1.size() != e2.size()) return false;
      for (size_t i = 0; i < e1.size(); ++i) {
        if (!equals_type(e1[i], e2[i])) return false;
      }
      return true;
    }
    case IrTypeData::Tag::variant_t: {
      const auto& v1 = t1.data.as_variant_t().tags;
      const auto& v2 = t2.data.as_variant_t().tags;
      if (v1.size() != v2.size()) return false;
      for (size_t i = 0; i < v1.size(); ++i) {
        if (!equals_variant_tag(v1[i], v2[i])) return false;
      }
      return true;
    }
    case IrTypeData::Tag::var_t:
      return t1.data.as_var_t() == t2.data.as_var_t();
  }
  return false;
}

inline bool equals_struct_field(const StructField& f1, const StructField& f2) {
  if (f1.id != f2.id) return false;
  if (!f1.type || !f2.type) return f1.type == f2.type;
  return equals_type(*f1.type, *f2.type);
}

inline bool equals_variant_tag(const VariantTag& t1, const VariantTag& t2) {
  if (t1.id != t2.id) return false;
  if (!t1.type || !t2.type) return t1.type == t2.type;
  return equals_type(*t1.type, *t2.type);
}

inline bool operator==(const IrType& t1, const IrType& t2) {
  return equals_type(t1, t2);
}

inline bool operator!=(const IrType& t1, const IrType& t2) {
  return !(t1 == t2);
}

inline bool operator==(const StructField& f1, const StructField& f2) {
  return equals_struct_field(f1, f2);
}

inline bool operator!=(const StructField& f1, const StructField& f2) {
  return !(f1 == f2);
}

inline bool operator==(const VariantTag& t1, const VariantTag& t2) {
  return equals_variant_tag(t1, t2);
}

inline bool operator!=(const VariantTag& t1, const VariantTag& t2) {
  return !(t1 == t2);
}

inline bool operator==(const TypeParam& tp1, const TypeParam& tp2) {
  if (tp1.var != tp2.var) return false;
  if (!equals_kind(tp1.kind, tp2.kind)) return false;
  if (tp1.bounds.size() != tp2.bounds.size()) return false;
  for (size_t i = 0; i < tp1.bounds.size(); ++i) {
    if (!equals_type(tp1.bounds[i], tp2.bounds[i])) return false;
  }
  return true;
}

inline bool operator!=(const TypeParam& tp1, const TypeParam& tp2) {
  return !(tp1 == tp2);
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

#include "ir_map_type.h"
#include "ir_substitute_type.h"
