#pragma once

#include "ir_base.h"
#include <sstream>
#include <stdexcept>

namespace ir {

struct StructField {
  std::string id;
  std::shared_ptr<IrType> type;

  std::string to_string() const;
  std::string to_json() const;
};

struct VariantTag {
  std::string id;
  std::shared_ptr<IrType> type;

  std::string to_string() const;
  std::string to_json() const;
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
  std::string to_json() const;

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
      return var;
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
  std::string to_json() const {
    return "{\"ID\":\"" + json_escape(id) + "\",\"Type\":" + type.to_json() + "}";
  }
};

inline std::string StructField::to_json() const {
  return "{\"ID\":\"" + json_escape(id) + "\",\"Type\":" + (type ? type->to_json() : "null") + "}";
}

inline std::string VariantTag::to_json() const {
  return "{\"ID\":\"" + json_escape(id) + "\",\"Type\":" + (type ? type->to_json() : "null") + "}";
}

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

inline std::string TypeParam::to_json() const {
  std::stringstream ss;
  ss << "{\"Var\":\"" << json_escape(var) << "\",\"Kind\":" << kind.to_json() << ",\"Bounds\":[";
  Interleave(bounds, [&]() { ss << ","; }, [&](int, const IrType& b) {
    ss << b.to_json();
  });
  ss << "]}";
  return ss.str();
}

inline std::string IrType::to_json() const {
  std::stringstream ss;
  ss << "{\"Case\":" << static_cast<int>(case_val);
  switch (case_val) {
    case IrTypeCase::AppType:
      if (app) {
        ss << ",\"App\":{\"Fun\":" << (app->fun ? app->fun->to_json() : "null")
           << ",\"Arg\":" << (app->arg ? app->arg->to_json() : "null") << "}";
      }
      break;
    case IrTypeCase::ArrayType:
      if (array) {
        ss << ",\"Array\":{\"ElemType\":" << (array->elem_type ? array->elem_type->to_json() : "null")
           << ",\"Size\":" << array->size << "}";
      }
      break;
    case IrTypeCase::ExistVarType:
      ss << ",\"ExistVar\":" << exist_var;
      break;
    case IrTypeCase::ForallType:
      if (forall) {
        ss << ",\"Forall\":{\"Var\":\"" << json_escape(forall->type_param.var) << "\""
           << ",\"Kind\":" << forall->type_param.kind.to_json()
           << ",\"Bounds\":[";
        Interleave(forall->type_param.bounds, [&]() { ss << ","; }, [&](int, const IrType& b) {
          ss << b.to_json();
        });
        ss << "],\"Type\":" << (forall->type ? forall->type->to_json() : "null") << "}";
      }
      break;
    case IrTypeCase::FunType:
      if (fun) {
        ss << ",\"Fun\":{\"Arg\":" << (fun->arg ? fun->arg->to_json() : "null")
           << ",\"Ret\":" << (fun->ret ? fun->ret->to_json() : "null") << "}";
      }
      break;
    case IrTypeCase::LambdaType:
      if (lambda) {
        ss << ",\"Lambda\":{\"Var\":\"" << json_escape(lambda->var) << "\""
           << ",\"Kind\":" << lambda->kind.to_json()
           << ",\"Type\":" << (lambda->type ? lambda->type->to_json() : "null") << "}";
      }
      break;
    case IrTypeCase::NameType:
      ss << ",\"Name\":\"" << json_escape(name) << "\"";
      break;
    case IrTypeCase::StructType:
      if (struct_data) {
        ss << ",\"Struct\":{\"Fields\":[";
        Interleave(struct_data->fields, [&]() { ss << ","; }, [&](int, const StructField& f) {
          ss << f.to_json();
        });
        ss << "]}";
      }
      break;
    case IrTypeCase::TupleType:
      if (tuple_data) {
        ss << ",\"Tuple\":{\"Elems\":[";
        Interleave(tuple_data->elems, [&]() { ss << ","; }, [&](int, const IrType& t) {
          ss << t.to_json();
        });
        ss << "]}";
      }
      break;
    case IrTypeCase::VariantType:
      if (variant_data) {
        ss << ",\"Variant\":{\"Tags\":[";
        Interleave(variant_data->tags, [&]() { ss << ","; }, [&](int, const VariantTag& t) {
          ss << t.to_json();
        });
        ss << "]}";
      }
      break;
    case IrTypeCase::VarType:
      ss << ",\"Var\":\"" << json_escape(var) << "\"";
      break;
  }
  ss << ",\"Pos\":" << pos.to_json() << "}";
  return ss.str();
}

} // namespace ir
