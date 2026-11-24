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
    if (typ.data.is(ir::IrTypeData::Tag::name_t) && typ.data.as_name_t() == alias_id_) {
      if (params_.empty()) {
        out_type = ir::new_var_type(self_var_);
        out_type.pos = typ.pos;
        return true;
      }
      err = ir::pos_to_string(pos_) + ": non-regular recursive type \"" + alias_id_ + "\": recursive occurrences must match the declared type parameters";
      return false;
    }
    if (typ.data.is(ir::IrTypeData::Tag::app_t)) {
      ir::IrType head = ir::type_app_fun(typ);
      if (head.data.is(ir::IrTypeData::Tag::name_t) && head.data.as_name_t() == alias_id_) {
        auto args = ir::type_app_args(typ);
        bool matches = (args.size() == params_.size());
        for (size_t i = 0; matches && i < args.size(); ++i) {
          if (!args[i].data.is(ir::IrTypeData::Tag::var_t) || args[i].data.as_var_t() != params_[i].first) {
            matches = false;
          }
        }
        if (matches) {
          out_type = ir::new_var_type(self_var_);
          out_type.pos = typ.pos;
          return true;
        }
        err = ir::pos_to_string(pos_) + ": non-regular recursive type \"" + alias_id_ + "\": recursive occurrences must match the declared type parameters";
        return false;
      }
    }
    return ir::map_type(typ, [&](const ir::IrType& child, ir::IrType& out_child, std::string& e) {
      return rewrite(child, out_child, e);
    }, out_type, err);
  }

 private:
  const std::string& alias_id_;
  const std::vector<std::pair<std::string, ir::IrKind>>& params_;
  const std::string& self_var_;
  const ir::Pos& pos_;
};

inline bool rewrite_self_recursive_decls(std::vector<ir::IrDecl>& decls, std::string& err) {
  for (auto& d : decls) {
    if (!d.is(ir::IrDeclData::Tag::alias_decl)) continue;
    auto& alias = d.data.as_alias_decl();
    const std::string& alias_id = alias.id;

    std::vector<std::pair<std::string, ir::IrKind>> params;
    const ir::IrType* cur = &alias.type;
    while (cur->data.is(ir::IrTypeData::Tag::lambda_t) && cur->data.as_lambda_t().body) {
      const auto& l = cur->data.as_lambda_t();
      params.push_back({l.var, l.kind});
      cur = l.body.get();
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
    alias.type = std::move(result);
  }
  return true;
}

inline bool rewrite_self_recursive_module(ir::IrModule& mod, std::string& err) {
  if (!rewrite_self_recursive_decls(mod.import_decls, err)) return false;
  for (auto& file : mod.files) {
    if (!rewrite_self_recursive_decls(file.decls, err)) return false;
  }
  return true;
}

} // namespace comp
