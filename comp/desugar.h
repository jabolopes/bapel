#pragma once

#include "bin/ir_base.h"
#include "bin/ir_decl.h"
#include "bin/ir_type.h"
#include "bin/ir_module.h"

#include <memory>
#include <set>
#include <string>
#include <utility>
#include <vector>

namespace comp {

class SelfRecursiveRewriter {
 public:
  SelfRecursiveRewriter(
      const std::string& alias_id,
      const std::vector<std::pair<std::string, ir::IrKind>>& params,
      const std::string& self_var,
      const ir::Pos& pos)
      : alias_id_(alias_id), params_(params), self_var_(self_var), pos_(pos) {}

  bool rewrite(const ir::IrType& typ, ir::IrType& out_type, std::string& err) const {
    switch (typ.case_val) {
      case ir::IrTypeCase::NameType: {
        if (typ.name == alias_id_) {
          if (params_.empty()) {
            out_type = ir::new_var_type(self_var_);
            out_type.pos = typ.pos;
            return true;
          }
          err = pos_.to_string() + ": non-regular recursive type \"" + alias_id_ + "\": recursive occurrences must match the declared type parameters";
          return false;
        }
        out_type = typ;
        return true;
      }
      case ir::IrTypeCase::AppType: {
        if (!typ.app) {
          out_type = typ;
          return true;
        }
        ir::IrType head = typ.app_fun();
        if (head.is(ir::IrTypeCase::NameType) && head.name == alias_id_) {
          auto args = typ.app_args();
          bool matches = (args.size() == params_.size());
          for (size_t i = 0; matches && i < args.size(); ++i) {
            if (!args[i].is(ir::IrTypeCase::VarType) || args[i].var != params_[i].first) {
              matches = false;
            }
          }
          if (matches) {
            out_type = ir::new_var_type(self_var_);
            out_type.pos = typ.pos;
            return true;
          }
          err = pos_.to_string() + ": non-regular recursive type \"" + alias_id_ + "\": recursive occurrences must match the declared type parameters";
          return false;
        }
        ir::IrType fun, arg;
        if (typ.app->fun && !rewrite(*typ.app->fun, fun, err)) return false;
        if (typ.app->arg && !rewrite(*typ.app->arg, arg, err)) return false;
        out_type = ir::new_app_type(std::move(fun), std::move(arg));
        out_type.pos = typ.pos;
        return true;
      }
      case ir::IrTypeCase::ArrayType: {
        if (!typ.array || !typ.array->elem_type) { out_type = typ; return true; }
        ir::IrType elem;
        if (!rewrite(*typ.array->elem_type, elem, err)) return false;
        out_type = ir::new_array_type(std::move(elem), typ.array->size);
        out_type.pos = typ.pos;
        return true;
      }
      case ir::IrTypeCase::FunType: {
        if (!typ.fun) { out_type = typ; return true; }
        ir::IrType arg, ret;
        if (typ.fun->arg && !rewrite(*typ.fun->arg, arg, err)) return false;
        if (typ.fun->ret && !rewrite(*typ.fun->ret, ret, err)) return false;
        out_type = ir::new_function_type(std::move(arg), std::move(ret));
        out_type.pos = typ.pos;
        return true;
      }
      case ir::IrTypeCase::ForallType: {
        if (!typ.forall) { out_type = typ; return true; }
        ir::TypeParam tp = typ.forall->type_param;
        for (auto& b : tp.bounds) {
          ir::IrType nb;
          if (!rewrite(b, nb, err)) return false;
          b = std::move(nb);
        }
        ir::IrType body;
        if (typ.forall->type && !rewrite(*typ.forall->type, body, err)) return false;
        out_type = ir::new_forall_type(std::move(tp), std::move(body));
        out_type.pos = typ.pos;
        return true;
      }
      case ir::IrTypeCase::LambdaType: {
        if (!typ.lambda || !typ.lambda->type) { out_type = typ; return true; }
        ir::IrType body;
        if (!rewrite(*typ.lambda->type, body, err)) return false;
        out_type = ir::new_lambda_type(typ.lambda->var, typ.lambda->kind, std::move(body));
        out_type.pos = typ.pos;
        return true;
      }
      case ir::IrTypeCase::MuType: {
        if (!typ.mu || !typ.mu->body) { out_type = typ; return true; }
        ir::IrType body;
        if (!rewrite(*typ.mu->body, body, err)) return false;
        out_type = ir::new_mu_type(typ.mu->var, typ.mu->kind, std::move(body));
        out_type.pos = typ.pos;
        return true;
      }
      case ir::IrTypeCase::StructType: {
        if (!typ.struct_data) { out_type = typ; return true; }
        std::vector<ir::StructField> fields;
        fields.reserve(typ.struct_data->fields.size());
        for (const auto& f : typ.struct_data->fields) {
          ir::IrType ft;
          if (f.type && !rewrite(*f.type, ft, err)) return false;
          fields.push_back(ir::StructField{f.id, std::make_shared<ir::IrType>(std::move(ft))});
        }
        out_type = ir::new_struct_type(std::move(fields));
        out_type.pos = typ.pos;
        return true;
      }
      case ir::IrTypeCase::TupleType: {
        if (!typ.tuple_data) { out_type = typ; return true; }
        std::vector<ir::IrType> elems;
        elems.reserve(typ.tuple_data->elems.size());
        for (const auto& e : typ.tuple_data->elems) {
          ir::IrType ne;
          if (!rewrite(e, ne, err)) return false;
          elems.push_back(std::move(ne));
        }
        out_type = ir::new_tuple_type(std::move(elems));
        out_type.pos = typ.pos;
        return true;
      }
      case ir::IrTypeCase::VariantType: {
        if (!typ.variant_data) { out_type = typ; return true; }
        std::vector<ir::VariantTag> tags;
        tags.reserve(typ.variant_data->tags.size());
        for (const auto& tag : typ.variant_data->tags) {
          ir::IrType tt;
          if (tag.type && !rewrite(*tag.type, tt, err)) return false;
          tags.push_back(ir::VariantTag{tag.id, std::make_shared<ir::IrType>(std::move(tt))});
        }
        out_type = ir::new_variant_type(std::move(tags));
        out_type.pos = typ.pos;
        return true;
      }
      case ir::IrTypeCase::VarType:
      case ir::IrTypeCase::ExistVarType:
        out_type = typ;
        return true;
    }
    out_type = typ;
    return true;
  }

 private:
  const std::string& alias_id_;
  const std::vector<std::pair<std::string, ir::IrKind>>& params_;
  const std::string& self_var_;
  const ir::Pos& pos_;
};

inline bool desugar_self_recursive_decls(std::vector<ir::IrDecl>& decls, std::string& err) {
  for (auto& d : decls) {
    if (!d.is(ir::IrDeclCase::AliasDecl) || !d.alias) continue;
    const std::string& alias_id = d.alias->id;

    std::vector<std::pair<std::string, ir::IrKind>> params;
    const ir::IrType* cur = &d.alias->type;
    while (cur->is(ir::IrTypeCase::LambdaType) && cur->lambda && cur->lambda->type) {
      params.push_back({cur->lambda->var, cur->lambda->kind});
      cur = cur->lambda->type.get();
    }

    std::set<std::string> names;
    ir::collect_name_types(*cur, names);
    if (names.find(alias_id) == names.end()) {
      continue;
    }

    std::vector<std::string> bound_v, free_v;
    ir::collect_free_vars(*cur, bound_v, free_v);
    std::set<std::string> used(free_v.begin(), free_v.end());
    for (const auto& p : params) used.insert(p.first);
    std::string self_var = "self";
    int idx = 0;
    while (used.find(self_var) != used.end()) {
      self_var = "self" + std::to_string(idx++);
    }

    SelfRecursiveRewriter rewriter(alias_id, params, self_var, d.pos);
    ir::IrType rewritten;
    if (!rewriter.rewrite(*cur, rewritten, err)) {
      return false;
    }

    ir::IrType result = ir::new_mu_type(self_var, ir::new_type_kind(), std::move(rewritten));
    for (auto it = params.rbegin(); it != params.rend(); ++it) {
      result = ir::new_lambda_type(it->first, it->second, std::move(result));
    }
    d.alias->type = std::move(result);
  }
  return true;
}

inline bool desugar_self_recursive_module(ir::IrModule& mod, std::string& err) {
  if (!desugar_self_recursive_decls(mod.import_decls, err)) return false;
  for (auto& file : mod.files) {
    if (!desugar_self_recursive_decls(file.decls, err)) return false;
  }
  return true;
}

} // namespace comp
