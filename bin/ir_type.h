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

struct StructTypeData {
  std::vector<StructField> fields;
};

struct TupleTypeData {
  std::vector<IrType> elems;
};

struct VariantTypeData {
  std::vector<VariantTag> tags;
};

struct IrType {
  IrTypeCase case_val = IrTypeCase::NameType;
  std::shared_ptr<AppTypeData> app;
  std::shared_ptr<ArrayTypeData> array;
  int64_t exist_var = 0;
  std::shared_ptr<ForallTypeData> forall;
  std::shared_ptr<FunctionTypeData> fun;
  std::shared_ptr<LambdaTypeData> lambda;
  std::string name;
  std::shared_ptr<StructTypeData> struct_data;
  std::shared_ptr<TupleTypeData> tuple_data;
  std::shared_ptr<VariantTypeData> variant_data;
  std::string var;
  Pos pos;

  bool is(IrTypeCase c) const { return case_val == c; }

  IrType app_fun() const {
    if (case_val == IrTypeCase::AppType && app) {
      return app->fun ? app->fun->app_fun() : *this;
    }
    return *this;
  }

  std::vector<IrType> app_args() const {
    std::vector<IrType> res;
    collect_app_args(res);
    return res;
  }

  void collect_app_args(std::vector<IrType>& acc) const {
    if (case_val == IrTypeCase::AppType && app) {
      if (app->fun) {
        app->fun->collect_app_args(acc);
      }
      if (app->arg) {
        acc.push_back(*app->arg);
      }
    }
  }

  std::vector<StructField> fields() const {
    if (case_val == IrTypeCase::StructType && struct_data) {
      return struct_data->fields;
    }
    return {};
  }

  std::vector<IrType> elems() const {
    if (case_val == IrTypeCase::TupleType && tuple_data) {
      return tuple_data->elems;
    }
    return {};
  }

  std::vector<VariantTag> tags() const {
    if (case_val == IrTypeCase::VariantType && variant_data) {
      return variant_data->tags;
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
    if (case_val == IrTypeCase::LambdaType && lambda) {
      res.push_back(lambda->var);
      if (lambda->type && lambda->type->is(IrTypeCase::LambdaType)) {
        auto next = lambda->type->lambda_vars();
        res.insert(res.end(), next.begin(), next.end());
      }
    }
    return res;
  }

  IrType lambda_body() const {
    if (case_val == IrTypeCase::LambdaType && lambda) {
      if (lambda->type && lambda->type->is(IrTypeCase::LambdaType)) {
        return lambda->type->lambda_body();
      }
      return lambda->type ? *lambda->type : *this;
    }
    return *this;
  }

  std::vector<TypeParam> forall_type_params() const {
    std::vector<TypeParam> res;
    if (case_val == IrTypeCase::ForallType && forall) {
      res.push_back(forall->type_param);
      if (forall->type && forall->type->is(IrTypeCase::ForallType)) {
        auto next = forall->type->forall_type_params();
        res.insert(res.end(), next.begin(), next.end());
      }
    }
    return res;
  }

  std::vector<std::string> forall_vars() const {
    std::vector<std::string> res;
    if (case_val == IrTypeCase::ForallType && forall) {
      res.push_back(forall->type_param.var);
      if (forall->type && forall->type->is(IrTypeCase::ForallType)) {
        auto next = forall->type->forall_vars();
        res.insert(res.end(), next.begin(), next.end());
      }
    }
    return res;
  }

  IrType forall_body() const {
    if (case_val == IrTypeCase::ForallType && forall) {
      if (forall->type && forall->type->is(IrTypeCase::ForallType)) {
        return forall->type->forall_body();
      }
      return forall->type ? *forall->type : *this;
    }
    return *this;
  }

  std::string trait_name() const {
    if (case_val == IrTypeCase::NameType) {
      return name;
    }
    if (case_val == IrTypeCase::AppType && app) {
      return app_fun().trait_name();
    }
    return "";
  }

  bool field_by_label(const std::string& label, int& out_index, StructField& out_field) const {
    if (case_val == IrTypeCase::StructType && struct_data) {
      try {
        int idx = std::stoi(label);
        if (idx >= 0 && idx < static_cast<int>(struct_data->fields.size())) {
          out_index = idx;
          out_field = struct_data->fields[idx];
          return true;
        }
      } catch (...) {}

      for (size_t i = 0; i < struct_data->fields.size(); ++i) {
        if (struct_data->fields[i].id == label) {
          out_index = static_cast<int>(i);
          out_field = struct_data->fields[i];
          return true;
        }
      }
    }
    return false;
  }

  bool elem_by_label(const std::string& label, int& out_index, IrType& out_elem) const {
    if (case_val == IrTypeCase::TupleType && tuple_data) {
      try {
        int idx = std::stoi(label);
        if (idx >= 0 && idx < static_cast<int>(tuple_data->elems.size())) {
          out_index = idx;
          out_elem = tuple_data->elems[idx];
          return true;
        }
      } catch (...) {}
    }
    return false;
  }

  bool tag_by_label(const std::string& label, int& out_index) const {
    if (case_val == IrTypeCase::VariantType && variant_data) {
      try {
        int idx = std::stoi(label);
        if (idx >= 0 && idx < static_cast<int>(variant_data->tags.size())) {
          out_index = idx;
          return true;
        }
      } catch (...) {}

      for (size_t i = 0; i < variant_data->tags.size(); ++i) {
        if (variant_data->tags[i].id == label) {
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

inline std::string VariantTag::to_string() const {
  return id + " " + (type ? type->to_string() : "()");
}

inline std::string IrType::to_string() const {
  switch (case_val) {
    case IrTypeCase::AppType: {
      if (!app) return "";
      std::string lparen = "", rparen = "";
      if (app->arg && (app->arg->is(IrTypeCase::AppType) || app->arg->is(IrTypeCase::ForallType) ||
                       app->arg->is(IrTypeCase::FunType) || app->arg->is(IrTypeCase::LambdaType))) {
        lparen = "("; rparen = ")";
      }
      return (app->fun ? app->fun->to_string() : "") + " " + lparen + (app->arg ? app->arg->to_string() : "") + rparen;
    }
    case IrTypeCase::ArrayType: {
      if (!array) return "[]";
      return "[" + (array->elem_type ? array->elem_type->to_string() : "") + ", " + std::to_string(array->size) + "]";
    }
    case IrTypeCase::ExistVarType:
      return "^" + std::to_string(exist_var);
    case IrTypeCase::ForallType: {
      if (!forall) return "";
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
    case IrTypeCase::FunType:
      if (!fun) return "";
      return (fun->arg ? fun->arg->to_string() : "") + " -> " + (fun->ret ? fun->ret->to_string() : "");
    case IrTypeCase::LambdaType:
      if (!lambda) return "";
      return "fun (" + lambda->var + ") (" + (lambda->type ? lambda->type->to_string() : "") + ")";
    case IrTypeCase::NameType:
      return name;
    case IrTypeCase::StructType: {
      if (!struct_data) return "struct{}";
      std::string s = "struct{";
      for (size_t i = 0; i < struct_data->fields.size(); ++i) {
        if (i > 0) s += ", ";
        s += struct_data->fields[i].to_string();
      }
      return s + "}";
    }
    case IrTypeCase::TupleType: {
      if (!tuple_data) return "()";
      std::string s = "(";
      for (size_t i = 0; i < tuple_data->elems.size(); ++i) {
        if (i > 0) s += ", ";
        s += tuple_data->elems[i].to_string();
      }
      return s + ")";
    }
    case IrTypeCase::VariantType: {
      if (!variant_data) return "variant{}";
      std::string s = "variant{";
      for (size_t i = 0; i < variant_data->tags.size(); ++i) {
        if (i > 0) s += ", ";
        s += variant_data->tags[i].to_string();
      }
      return s + "}";
    }
    case IrTypeCase::VarType:
      return var.empty() || var[0] == '\'' ? var : ("'" + var);
  }
  return "";
}

inline IrType new_name_type(const std::string& name) {
  IrType t;
  t.case_val = IrTypeCase::NameType;
  t.name = name;
  return t;
}

inline IrType new_var_type(const std::string& var) {
  IrType t;
  t.case_val = IrTypeCase::VarType;
  t.var = var;
  return t;
}

inline IrType new_exist_var_type(int64_t evar) {
  IrType t;
  t.case_val = IrTypeCase::ExistVarType;
  t.exist_var = evar;
  return t;
}

inline IrType new_app_type(IrType fun, IrType arg) {
  IrType t;
  t.case_val = IrTypeCase::AppType;
  t.app = std::make_shared<AppTypeData>();
  t.app->fun = std::make_shared<IrType>(std::move(fun));
  t.app->arg = std::make_shared<IrType>(std::move(arg));
  return t;
}

inline IrType new_array_type(IrType elem, int64_t size) {
  IrType t;
  t.case_val = IrTypeCase::ArrayType;
  t.array = std::make_shared<ArrayTypeData>();
  t.array->elem_type = std::make_shared<IrType>(std::move(elem));
  t.array->size = size;
  return t;
}

inline IrType new_function_type(IrType arg, IrType ret) {
  IrType t;
  t.case_val = IrTypeCase::FunType;
  t.fun = std::make_shared<FunctionTypeData>();
  t.fun->arg = std::make_shared<IrType>(std::move(arg));
  t.fun->ret = std::make_shared<IrType>(std::move(ret));
  return t;
}

inline IrType new_tuple_type(std::vector<IrType> elems) {
  if (elems.size() == 1) {
    return elems[0];
  }
  IrType t;
  t.case_val = IrTypeCase::TupleType;
  t.tuple_data = std::make_shared<TupleTypeData>();
  t.tuple_data->elems = std::move(elems);
  return t;
}

inline IrType new_struct_type(std::vector<StructField> fields) {
  IrType t;
  t.case_val = IrTypeCase::StructType;
  t.struct_data = std::make_shared<StructTypeData>();
  t.struct_data->fields = std::move(fields);
  return t;
}

inline IrType new_variant_type(std::vector<VariantTag> tags) {
  IrType t;
  t.case_val = IrTypeCase::VariantType;
  t.variant_data = std::make_shared<VariantTypeData>();
  t.variant_data->tags = std::move(tags);
  return t;
}

inline IrType new_forall_type(TypeParam tp, IrType body) {
  IrType t;
  t.case_val = IrTypeCase::ForallType;
  t.forall = std::make_shared<ForallTypeData>();
  t.forall->type_param = std::move(tp);
  t.forall->type = std::make_shared<IrType>(std::move(body));
  return t;
}

inline IrType new_lambda_type(std::string var, IrKind kind, IrType body) {
  IrType t;
  t.case_val = IrTypeCase::LambdaType;
  t.lambda = std::make_shared<LambdaTypeData>();
  t.lambda->var = std::move(var);
  t.lambda->kind = std::move(kind);
  t.lambda->type = std::make_shared<IrType>(std::move(body));
  return t;
}

inline IrType forall_vars(const std::vector<TypeParam>& params, IrType body) {
  for (auto it = params.rbegin(); it != params.rend(); ++it) {
    body = new_forall_type(*it, std::move(body));
  }
  return body;
}

inline void collect_free_vars(const IrType& typ, std::vector<std::string>& bound_vars, std::vector<std::string>& free_vars) {
  switch (typ.case_val) {
    case IrTypeCase::AppType:
      if (typ.app) {
        if (typ.app->fun) collect_free_vars(*typ.app->fun, bound_vars, free_vars);
        if (typ.app->arg) collect_free_vars(*typ.app->arg, bound_vars, free_vars);
      }
      break;
    case IrTypeCase::ArrayType:
      if (typ.array && typ.array->elem_type) {
        collect_free_vars(*typ.array->elem_type, bound_vars, free_vars);
      }
      break;
    case IrTypeCase::ForallType:
      if (typ.forall) {
        for (const auto& b : typ.forall->type_param.bounds) {
          collect_free_vars(b, bound_vars, free_vars);
        }
        bound_vars.push_back(typ.forall->type_param.var);
        if (typ.forall->type) collect_free_vars(*typ.forall->type, bound_vars, free_vars);
        bound_vars.pop_back();
      }
      break;
    case IrTypeCase::FunType:
      if (typ.fun) {
        if (typ.fun->arg) collect_free_vars(*typ.fun->arg, bound_vars, free_vars);
        if (typ.fun->ret) collect_free_vars(*typ.fun->ret, bound_vars, free_vars);
      }
      break;
    case IrTypeCase::LambdaType:
      if (typ.lambda) {
        bound_vars.push_back(typ.lambda->var);
        if (typ.lambda->type) collect_free_vars(*typ.lambda->type, bound_vars, free_vars);
        bound_vars.pop_back();
      }
      break;
    case IrTypeCase::StructType:
      if (typ.struct_data) {
        for (const auto& f : typ.struct_data->fields) {
          if (f.type) collect_free_vars(*f.type, bound_vars, free_vars);
        }
      }
      break;
    case IrTypeCase::TupleType:
      if (typ.tuple_data) {
        for (const auto& elem : typ.tuple_data->elems) {
          collect_free_vars(elem, bound_vars, free_vars);
        }
      }
      break;
    case IrTypeCase::VariantType:
      if (typ.variant_data) {
        for (const auto& tag : typ.variant_data->tags) {
          if (tag.type) collect_free_vars(*tag.type, bound_vars, free_vars);
        }
      }
      break;
    case IrTypeCase::VarType:
      if (std::find(bound_vars.begin(), bound_vars.end(), typ.var) == bound_vars.end()) {
        if (std::find(free_vars.begin(), free_vars.end(), typ.var) == free_vars.end()) {
          free_vars.push_back(typ.var);
        }
      }
      break;
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
  if (t1.case_val != t2.case_val) return false;
  switch (t1.case_val) {
    case IrTypeCase::AppType:
      if (!t1.app || !t2.app) return t1.app == t2.app;
      if (!t1.app->fun || !t2.app->fun || !t1.app->arg || !t2.app->arg) return false;
      return equals_type(*t1.app->fun, *t2.app->fun) && equals_type(*t1.app->arg, *t2.app->arg);
    case IrTypeCase::ArrayType:
      if (!t1.array || !t2.array) return t1.array == t2.array;
      if (t1.array->size != t2.array->size) return false;
      if (!t1.array->elem_type || !t2.array->elem_type) return t1.array->elem_type == t2.array->elem_type;
      return equals_type(*t1.array->elem_type, *t2.array->elem_type);
    case IrTypeCase::ExistVarType:
      return t1.exist_var == t2.exist_var;
    case IrTypeCase::ForallType: {
      if (!t1.forall || !t2.forall) return t1.forall == t2.forall;
      if (t1.forall->type_param.var != t2.forall->type_param.var) return false;
      if (!equals_kind(t1.forall->type_param.kind, t2.forall->type_param.kind)) return false;
      if (t1.forall->type_param.bounds.size() != t2.forall->type_param.bounds.size()) return false;
      for (size_t i = 0; i < t1.forall->type_param.bounds.size(); ++i) {
        if (!equals_type(t1.forall->type_param.bounds[i], t2.forall->type_param.bounds[i])) return false;
      }
      if (!t1.forall->type || !t2.forall->type) return t1.forall->type == t2.forall->type;
      return equals_type(*t1.forall->type, *t2.forall->type);
    }
    case IrTypeCase::FunType:
      if (!t1.fun || !t2.fun) return t1.fun == t2.fun;
      if (!t1.fun->arg || !t2.fun->arg || !t1.fun->ret || !t2.fun->ret) return false;
      return equals_type(*t1.fun->arg, *t2.fun->arg) && equals_type(*t1.fun->ret, *t2.fun->ret);
    case IrTypeCase::LambdaType:
      if (!t1.lambda || !t2.lambda) return t1.lambda == t2.lambda;
      if (t1.lambda->var != t2.lambda->var) return false;
      if (!equals_kind(t1.lambda->kind, t2.lambda->kind)) return false;
      if (!t1.lambda->type || !t2.lambda->type) return t1.lambda->type == t2.lambda->type;
      return equals_type(*t1.lambda->type, *t2.lambda->type);
    case IrTypeCase::NameType:
      return t1.name == t2.name;
    case IrTypeCase::StructType: {
      if (!t1.struct_data || !t2.struct_data) return t1.struct_data == t2.struct_data;
      if (t1.struct_data->fields.size() != t2.struct_data->fields.size()) return false;
      for (size_t i = 0; i < t1.struct_data->fields.size(); ++i) {
        if (!equals_struct_field(t1.struct_data->fields[i], t2.struct_data->fields[i])) return false;
      }
      return true;
    }
    case IrTypeCase::TupleType: {
      if (!t1.tuple_data || !t2.tuple_data) return t1.tuple_data == t2.tuple_data;
      if (t1.tuple_data->elems.size() != t2.tuple_data->elems.size()) return false;
      for (size_t i = 0; i < t1.tuple_data->elems.size(); ++i) {
        if (!equals_type(t1.tuple_data->elems[i], t2.tuple_data->elems[i])) return false;
      }
      return true;
    }
    case IrTypeCase::VariantType: {
      if (!t1.variant_data || !t2.variant_data) return t1.variant_data == t2.variant_data;
      if (t1.variant_data->tags.size() != t2.variant_data->tags.size()) return false;
      for (size_t i = 0; i < t1.variant_data->tags.size(); ++i) {
        if (!equals_variant_tag(t1.variant_data->tags[i], t2.variant_data->tags[i])) return false;
      }
      return true;
    }
    case IrTypeCase::VarType:
      return t1.var == t2.var;
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

inline IrType substitute_type(const IrType& t, const IrType& source, const IrType& target) {
  if (equals_type(t, source)) {
    return target;
  }
  switch (t.case_val) {
    case IrTypeCase::AppType: {
      if (!t.app) return t;
      IrType fun = t.app->fun ? substitute_type(*t.app->fun, source, target) : IrType{};
      IrType arg = t.app->arg ? substitute_type(*t.app->arg, source, target) : IrType{};
      return new_app_type(std::move(fun), std::move(arg));
    }
    case IrTypeCase::ArrayType: {
      if (!t.array) return t;
      IrType elem = t.array->elem_type ? substitute_type(*t.array->elem_type, source, target) : IrType{};
      return new_array_type(std::move(elem), t.array->size);
    }
    case IrTypeCase::ExistVarType:
      return t;
    case IrTypeCase::ForallType: {
      if (!t.forall) return t;
      if (source.is(IrTypeCase::VarType) && source.var == t.forall->type_param.var) {
        std::vector<IrType> bounds;
        bounds.reserve(t.forall->type_param.bounds.size());
        for (const auto& b : t.forall->type_param.bounds) {
          bounds.push_back(substitute_type(b, source, target));
        }
        TypeParam tp{t.forall->type_param.var, t.forall->type_param.kind, std::move(bounds)};
        return new_forall_type(std::move(tp), t.forall->type ? *t.forall->type : IrType{});
      }
      std::vector<TypeParam> target_free = get_free_type_vars(target);
      bool capture = false;
      for (const auto& tf : target_free) {
        if (tf.var == t.forall->type_param.var) {
          capture = true;
          break;
        }
      }
      if (capture) {
        std::vector<std::string> bound_v, all_free_v;
        collect_free_vars(t, bound_v, all_free_v);
        for (const auto& tf : target_free) all_free_v.push_back(tf.var);
        if (source.is(IrTypeCase::VarType)) all_free_v.push_back(source.var);
        std::set<std::string> used(all_free_v.begin(), all_free_v.end());
        used.insert(t.forall->type_param.var);

        std::string fresh_name;
        for (char c = 'a'; c <= 'z'; ++c) {
          std::string s(1, c);
          if (used.find(s) == used.end()) {
            fresh_name = s;
            break;
          }
        }
        if (fresh_name.empty()) {
          int64_t idx = 0;
          while (true) {
            std::string s = "t" + std::to_string(idx++);
            if (used.find(s) == used.end()) {
              fresh_name = s;
              break;
            }
          }
        }

        std::vector<IrType> renamed_bounds;
        renamed_bounds.reserve(t.forall->type_param.bounds.size());
        for (const auto& b : t.forall->type_param.bounds) {
          renamed_bounds.push_back(substitute_type(b, new_var_type(t.forall->type_param.var), new_var_type(fresh_name)));
        }
        IrType renamed_body = t.forall->type ? substitute_type(*t.forall->type, new_var_type(t.forall->type_param.var), new_var_type(fresh_name)) : IrType{};

        std::vector<IrType> final_bounds;
        final_bounds.reserve(renamed_bounds.size());
        for (const auto& b : renamed_bounds) {
          final_bounds.push_back(substitute_type(b, source, target));
        }
        IrType final_body = substitute_type(renamed_body, source, target);
        TypeParam tp{fresh_name, t.forall->type_param.kind, std::move(final_bounds)};
        return new_forall_type(std::move(tp), std::move(final_body));
      }

      std::vector<IrType> bounds;
      bounds.reserve(t.forall->type_param.bounds.size());
      for (const auto& b : t.forall->type_param.bounds) {
        bounds.push_back(substitute_type(b, source, target));
      }
      TypeParam tp{t.forall->type_param.var, t.forall->type_param.kind, std::move(bounds)};
      IrType body = t.forall->type ? substitute_type(*t.forall->type, source, target) : IrType{};
      return new_forall_type(std::move(tp), std::move(body));
    }
    case IrTypeCase::FunType: {
      if (!t.fun) return t;
      IrType arg = t.fun->arg ? substitute_type(*t.fun->arg, source, target) : IrType{};
      IrType ret = t.fun->ret ? substitute_type(*t.fun->ret, source, target) : IrType{};
      return new_function_type(std::move(arg), std::move(ret));
    }
    case IrTypeCase::LambdaType: {
      if (!t.lambda) return t;
      if (source.is(IrTypeCase::VarType) && source.var == t.lambda->var) {
        return t;
      }
      std::vector<TypeParam> target_free = get_free_type_vars(target);
      bool capture = false;
      for (const auto& tf : target_free) {
        if (tf.var == t.lambda->var) {
          capture = true;
          break;
        }
      }
      if (capture) {
        std::vector<std::string> bound_v, all_free_v;
        collect_free_vars(t, bound_v, all_free_v);
        for (const auto& tf : target_free) all_free_v.push_back(tf.var);
        if (source.is(IrTypeCase::VarType)) all_free_v.push_back(source.var);
        std::set<std::string> used(all_free_v.begin(), all_free_v.end());
        used.insert(t.lambda->var);

        std::string fresh_name;
        for (char c = 'a'; c <= 'z'; ++c) {
          std::string s(1, c);
          if (used.find(s) == used.end()) {
            fresh_name = s;
            break;
          }
        }
        if (fresh_name.empty()) {
          int64_t idx = 0;
          while (true) {
            std::string s = "t" + std::to_string(idx++);
            if (used.find(s) == used.end()) {
              fresh_name = s;
              break;
            }
          }
        }

        IrType renamed_body = t.lambda->type ? substitute_type(*t.lambda->type, new_var_type(t.lambda->var), new_var_type(fresh_name)) : IrType{};
        IrType final_body = substitute_type(renamed_body, source, target);
        return new_lambda_type(fresh_name, t.lambda->kind, std::move(final_body));
      }
      IrType body = t.lambda->type ? substitute_type(*t.lambda->type, source, target) : IrType{};
      return new_lambda_type(t.lambda->var, t.lambda->kind, std::move(body));
    }
    case IrTypeCase::NameType:
      return t;
    case IrTypeCase::StructType: {
      if (!t.struct_data) return t;
      std::vector<StructField> fields;
      fields.reserve(t.struct_data->fields.size());
      for (const auto& f : t.struct_data->fields) {
        IrType ft = f.type ? substitute_type(*f.type, source, target) : IrType{};
        fields.push_back(StructField{f.id, std::make_shared<IrType>(std::move(ft))});
      }
      return new_struct_type(std::move(fields));
    }
    case IrTypeCase::TupleType: {
      if (!t.tuple_data) return t;
      std::vector<IrType> elems;
      elems.reserve(t.tuple_data->elems.size());
      for (const auto& e : t.tuple_data->elems) {
        elems.push_back(substitute_type(e, source, target));
      }
      return new_tuple_type(std::move(elems));
    }
    case IrTypeCase::VariantType: {
      if (!t.variant_data) return t;
      std::vector<VariantTag> tags;
      tags.reserve(t.variant_data->tags.size());
      for (const auto& tag : t.variant_data->tags) {
        IrType tt = tag.type ? substitute_type(*tag.type, source, target) : IrType{};
        tags.push_back(VariantTag{tag.id, std::make_shared<IrType>(std::move(tt))});
      }
      return new_variant_type(std::move(tags));
    }
    case IrTypeCase::VarType:
      return t;
  }
  return t;
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
