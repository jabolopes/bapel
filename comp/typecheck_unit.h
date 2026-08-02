#pragma once

#include "ast/ast_desugar.h"
#include "bin/ir_base.h"
#include "bin/ir_decl.h"
#include "bin/ir_function.h"
#include "bin/ir_unit.h"
#include "comp/desugar.h"
#include "comp/normalizer.h"
#include "comp/querier.h"
#include "comp/resolver.h"
#include "ts/bind.h"
#include "ts/context.h"
#include "ts/inferencer.h"
#include "ts/typecheck.h"

#include <cstdio>
#include <fstream>
#include <map>
#include <memory>
#include <set>
#include <sstream>
#include <string>
#include <vector>

namespace comp {

inline ts::Context new_default_context() {
  ts::Context ctx;

  // Fundamental types
  ctx = ctx.add_bind(ts::new_const_bind("bool", ir::new_type_kind()));
  ctx = ctx.add_bind(ts::new_const_bind("i8", ir::new_type_kind()));
  ctx = ctx.add_bind(ts::new_const_bind("i16", ir::new_type_kind()));
  ctx = ctx.add_bind(ts::new_const_bind("i32", ir::new_type_kind()));
  ctx = ctx.add_bind(ts::new_const_bind("i64", ir::new_type_kind()));
  ctx = ctx.add_bind(ts::new_const_bind("f32", ir::new_type_kind()));
  ctx = ctx.add_bind(ts::new_const_bind("f64", ir::new_type_kind()));
  ctx = ctx.add_bind(ts::new_const_bind("void", ir::new_type_kind()));

  // Fundamental terms
  ctx = ctx.add_bind(ts::new_term_decl_bind("true", ir::new_name_type("bool")));
  ctx = ctx.add_bind(ts::new_term_decl_bind("false", ir::new_name_type("bool")));

  // Operators
  std::vector<std::string> ops = {"||", "&&", "!=", "==", ">", ">=", "<", "<=", "+", "-", "*", "/", "!"};
  for (const auto& op : ops) {
    ctx = ctx.add_bind(ts::new_term_decl_bind(op, ir::operator_type(op)));
  }

  // ifthen
  ctx = ctx.add_bind(ts::new_term_decl_bind(
      "ifthen",
      ir::new_forall_type(
          ir::TypeParam{"a", ir::new_type_kind(), {}},
          ir::new_function_type(
              ir::new_tuple_type({ir::new_name_type("bool"), ir::new_var_type("a")}),
              ir::new_var_type("a")))));

  // ifelse
  ctx = ctx.add_bind(ts::new_term_decl_bind(
      "ifelse",
      ir::new_forall_type(
          ir::TypeParam{"a", ir::new_type_kind(), {}},
          ir::new_function_type(
              ir::new_tuple_type({ir::new_name_type("bool"), ir::new_var_type("a"), ir::new_var_type("a")}),
              ir::new_var_type("a")))));

  // core::for
  ctx = ctx.add_bind(ts::new_term_decl_bind(
      "core::for",
      ir::new_forall_type(
          ir::TypeParam{"a", ir::new_type_kind(), {}},
          ir::new_function_type(
              ir::new_tuple_type({
                  ir::new_name_type("bool"),
                  ir::new_function_type(ir::new_tuple_type({}), ir::new_var_type("a")),
              }),
              ir::new_tuple_type({})))));

  return ctx;
}

struct TypecheckOptions {
  bool skip_default_context = false;
  bool skip_term_typechecker = false;
  bool skip_undefined_term_checks = false;
};

struct SymbolStatus {
  ir::IrDecl decl;
  bool declared = false;
  bool defined = false;
};

class UnitChecker {
 public:
  UnitChecker(TypecheckOptions options, ts::Context context)
      : options_(options), context_(std::move(context)) {}

  bool check_unit(ir::IrUnit* unit, std::string& err) {
    for (const auto& decl : unit->import_decls) {
      if (!add_symbol(decl, err)) return false;
    }

    for (auto& ti : unit->imported_trait_impls) {
      if (!add_trait_impl(&ti, err)) return false;
    }

    try {
      context_ = context_.enter_scope();
    } catch (const std::exception& e) {
      err = e.what();
      return false;
    }

    std::set<std::string> local_decl_ids;
    std::vector<ir::IrDecl> merged_decls;
    for (const auto& decl : unit->decls) {
      local_decl_ids.insert(decl.id());
      merged_decls.push_back(decl);
      if (decl.is(ir::IrDeclCase::NameDecl) && decl.name) {
        local_types_.insert(decl.name->id);
      } else if (decl.is(ir::IrDeclCase::AliasDecl) && decl.alias) {
        local_types_.insert(decl.alias->id);
      }
    }
    for (const auto& decl : unit->impl_decls) {
      merged_decls.push_back(decl);
      if (decl.is(ir::IrDeclCase::NameDecl) && decl.name) {
        local_types_.insert(decl.name->id);
      } else if (decl.is(ir::IrDeclCase::AliasDecl) && decl.alias) {
        local_types_.insert(decl.alias->id);
      }
    }

    std::vector<ir::IrDecl> sorted_decls;
    if (!ir::topo_sort_decls(merged_decls, sorted_decls, err)) {
      return false;
    }

    for (const auto& decl : sorted_decls) {
      if (local_decl_ids.find(decl.id()) != local_decl_ids.end()) {
        if (!add_decl(decl, err)) return false;
      } else {
        if (!add_symbol(decl, err)) return false;
      }
    }

    for (auto& ti : unit->trait_impls) {
      if (!add_trait_impl(&ti, err)) return false;
    }

    for (auto& fn : unit->functions) {
      if (!add_function(&fn, err)) return false;
    }

    for (auto& ti : unit->trait_impls) {
      if (!check_trait_impl(&ti, err)) return false;
    }

    if (!options_.skip_undefined_term_checks) {
      for (const auto& pair : symbols_) {
        if (pair.second.declared && !pair.second.defined) {
          err = pair.second.decl.pos.to_string() + ": symbol \"" + pair.first + "\" is declared but not defined";
          return false;
        }
      }
    }

    return true;
  }

  const ts::Context& context() const { return context_; }

 private:
  bool add_symbol(const ir::IrDecl& decl, std::string& err) {
    try {
      context_ = context_.add_symbol(decl);
      return true;
    } catch (const std::exception& e) {
      err = e.what();
      return false;
    }
  }

  bool add_decl(const ir::IrDecl& decl, std::string& err) {
    if (decl.is(ir::IrDeclCase::TermDecl) && decl.term) {
      auto& sym = symbols_[decl.term->id];
      sym.decl = decl;
      if (sym.declared) {
        err = decl.pos.to_string() + ": symbol \"" + decl.term->id + "\" already declared";
        return false;
      }
      sym.declared = true;
    }
    return add_symbol(decl, err);
  }

  bool check_function(ir::IrFunction* function, std::string& err) {
    ts::Typechecker typechecker(context_);
    ts::Context updated_ctx;

    if (options_.skip_term_typechecker) {
      if (!typechecker.infer_function(function, updated_ctx)) {
        err = typechecker.error();
        return false;
      }
      context_ = updated_ctx;
    } else {
      if (!typechecker.infer_function(function, updated_ctx)) {
        err = typechecker.error();
        return false;
      }
      ts::Context out_ctx;
      if (!typechecker.typecheck_function(function, out_ctx)) {
        err = typechecker.error();
        return false;
      }
      context_ = out_ctx;
    }
    return true;
  }

  bool add_function(ir::IrFunction* function, std::string& err) {
    ir::IrDecl decl = function->decl();
    auto& sym = symbols_[decl.id()];
    sym.decl = decl;
    if (sym.defined) {
      err = decl.pos.to_string() + ": symbol \"" + decl.id() + "\" already defined";
      return false;
    }
    sym.defined = true;
    return check_function(function, err);
  }

  bool add_trait_impl(ir::IrTraitImpl* impl, std::string& err) {
    ts::Context ctx = context_;
    for (const auto& tp : impl->type_params) {
      try {
        ctx = ctx.add_bind(ts::new_type_param_bind(tp));
      } catch (const std::exception& e) {
        err = e.what();
        return false;
      }
    }

    ts::Typechecker tc(ctx);
    ir::IrType reduced_type = tc.reduce_type(impl->type_name);

    if (impl->case_val == ir::ImplCase::InherentImpl) {
      std::string base_name = ts::base_type_name(reduced_type);
      if (base_name.empty()) {
        err = "invalid type for inherent impl: " + impl->type_name.to_string();
        return false;
      }
      for (const auto& m : impl->methods) {
        std::string name = base_name + "::" + m.id;
        std::vector<ir::IrType> args;
        for (const auto& a : m.args) {
          ir::IrType t = ir::substitute_type(a.type, ir::new_name_type("Self"), reduced_type);
          args.push_back(std::move(t));
        }
        ir::IrType ret = ir::substitute_type(m.ret_type, ir::new_name_type("Self"), reduced_type);
        ir::IrType method_type = ir::new_function_type(ir::new_tuple_type(std::move(args)), std::move(ret));

        std::vector<ir::TypeParam> tvars = impl->type_params;
        tvars.insert(tvars.end(), m.type_params.begin(), m.type_params.end());
        for (auto it = tvars.rbegin(); it != tvars.rend(); ++it) {
          method_type = ir::new_forall_type(*it, std::move(method_type));
        }

        try {
          context_ = context_.add_bind(ts::new_term_decl_bind(name, method_type));
        } catch (const std::exception& e) {
          err = e.what();
          return false;
        }
      }
      return true;
    }

    ir::IrType reduced_trait_type = tc.reduce_type(impl->trait_type);
    try {
      context_ = context_.add_bind(ts::new_trait_impl_bind(impl->type_params, reduced_trait_type, reduced_type));
      return true;
    } catch (const std::exception& e) {
      err = e.what();
      return false;
    }
  }

  bool check_trait_impl(ir::IrTraitImpl* impl, std::string& err) {
    std::string base_type = ts::base_type_name(impl->type_name);
    if (base_type.empty()) {
      err = impl->pos.to_string() + ": invalid type in impl: " + impl->type_name.to_string();
      return false;
    }
    if (local_types_.find(base_type) == local_types_.end() && base_type != "Ptr" && base_type != "StringView") {
      err = impl->pos.to_string() + ": coherence violation: cannot implement trait for foreign type \"" + base_type + "\"";
      return false;
    }

    ts::Context method_context = context_;
    for (const auto& tp : impl->type_params) {
      try {
        method_context = method_context.add_bind(ts::new_type_param_bind(tp));
      } catch (const std::exception& e) {
        err = e.what();
        return false;
      }
    }
    try {
      method_context = method_context.add_bind(ts::new_alias_bind("Self", impl->type_name));
    } catch (const std::exception& e) {
      err = e.what();
      return false;
    }

    if (impl->case_val == ir::ImplCase::InherentImpl) {
      ts::Typechecker tc(method_context);
      for (auto& m : impl->methods) {
        ts::Context out_ctx;
        if (!tc.infer_function(&m, out_ctx)) {
          err = tc.error();
          return false;
        }
        if (!tc.typecheck_function(&m, out_ctx)) {
          err = tc.error();
          return false;
        }
      }
      return true;
    }

    std::string trait_name;
    if (impl->trait_type.is(ir::IrTypeCase::NameType)) {
      trait_name = impl->trait_type.name;
    } else if (impl->trait_type.is(ir::IrTypeCase::AppType)) {
      trait_name = ts::base_type_name(impl->trait_type);
    }

    ts::Binding bind;
    if (!context_.lookup_trait_bind(trait_name, bind) || !bind.trait) {
      err = "trait \"" + trait_name + "\" not found";
      return false;
    }
    const auto& trait = *bind.trait;

    std::map<std::string, ir::IrFunction*> implemented;
    for (auto& m : impl->methods) {
      implemented[m.id] = &m;
    }

    for (const auto& tm : trait.methods) {
      auto it = implemented.find(tm.id);
      if (it == implemented.end()) {
        err = "method \"" + tm.id + "\" of trait " + impl->trait_type.to_string() + " is not implemented for " + impl->type_name.to_string();
        return false;
      }

      auto* impl_method = it->second;
      auto trait_args = impl->trait_type.app_args();

      auto substitute = [&](ir::IrType t) -> ir::IrType {
        ir::IrType res = ir::substitute_type(t, ir::new_name_type("Self"), impl->type_name);
        for (size_t i = 0; i < trait.type_params.size() && i < trait_args.size(); ++i) {
          res = ir::substitute_type(res, ir::new_var_type(trait.type_params[i].var), trait_args[i]);
        }
        return res;
      };

      std::vector<ir::IrType> expected_args;
      for (const auto& a : tm.args) {
        expected_args.push_back(substitute(a.type));
      }
      ir::IrType expected_ret = substitute(tm.ret_type);
      ir::IrType expected_type = ir::new_function_type(ir::new_tuple_type(std::move(expected_args)), std::move(expected_ret));

      std::vector<ir::IrType> actual_args;
      for (const auto& a : impl_method->args) {
        actual_args.push_back(a.type);
      }
      ir::IrType actual_type = ir::new_function_type(ir::new_tuple_type(std::move(actual_args)), impl_method->ret_type);

      ts::Typechecker tc(method_context);
      if (!tc.subtype(expected_type, actual_type)) {
        err = "method \"" + impl_method->id + "\" has type " + actual_type.to_string() + " that does not match trait signature " + expected_type.to_string();
        return false;
      }
      if (!tc.subtype(actual_type, expected_type)) {
        err = "method \"" + impl_method->id + "\" has type " + actual_type.to_string() + " that does not match trait signature " + expected_type.to_string();
        return false;
      }

      ts::Context out_ctx;
      if (!tc.infer_function(impl_method, out_ctx)) {
        err = tc.error();
        return false;
      }
      if (!tc.typecheck_function(impl_method, out_ctx)) {
        err = tc.error();
        return false;
      }
    }

    return true;
  }

  TypecheckOptions options_;
  ts::Context context_;
  std::map<std::string, SymbolStatus> symbols_;
  std::set<std::string> local_types_;
};

inline bool typecheck_unit(TypecheckOptions options, ir::IrUnit* unit, std::string& err) {
  ts::Context ctx = options.skip_default_context ? ts::Context{} : new_default_context();
  UnitChecker checker(options, std::move(ctx));
  return checker.check_unit(unit, err);
}

inline bool typecheck_source_file(Querier querier, TypecheckOptions options, const std::string& input_filename, ir::IrUnit& out_unit, std::string& err) {
  if (!normalize_source_file(std::move(querier), input_filename, out_unit, err)) {
    return false;
  }
  return typecheck_unit(options, &out_unit, err);
}

} // namespace comp
