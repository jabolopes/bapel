#pragma once

#include "ir.h"
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

} // namespace ir
