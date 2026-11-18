#pragma once

#include "ir_type.h"
#include <set>
#include <string>
#include <vector>

namespace ir {

// TODO: Factor out shared capture-avoiding binder substitution (ForallType, LambdaType, MuType) into a helper.
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
    case IrTypeCase::MuType: {
      if (!t.mu) return t;
      if (source.is(IrTypeCase::VarType) && source.var == t.mu->var) {
        return t;
      }
      std::vector<TypeParam> target_free = get_free_type_vars(target);
      bool capture = false;
      for (const auto& tf : target_free) {
        if (tf.var == t.mu->var) {
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
        used.insert(t.mu->var);

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

        IrType renamed_body = t.mu->body ? substitute_type(*t.mu->body, new_var_type(t.mu->var), new_var_type(fresh_name)) : IrType{};
        IrType final_body = substitute_type(renamed_body, source, target);
        return new_mu_type(fresh_name, t.mu->kind, std::move(final_body));
      }
      IrType body = t.mu->body ? substitute_type(*t.mu->body, source, target) : IrType{};
      return new_mu_type(t.mu->var, t.mu->kind, std::move(body));
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

} // namespace ir
