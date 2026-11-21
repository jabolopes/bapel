#pragma once

#include "comp/desugar.h"
#include "bin/ir_base.h"
#include "bin/ir_decl.h"
#include "bin/ir_function.h"
#include "bin/ir_typed_module.h"
#include "bin/ir_module.h"
#include "comp/rewrite_recursive_types.h"
#include "comp/normalizer.h"
#include "comp/querier.h"
#include "comp/resolver.h"
#include "ts/bind.h"
#include "ts/context.h"
#include "ts/inferencer.h"
#include "ts/typecheck.h"
#include "ts/wellformed.h"

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

class ModuleChecker {
 public:
  ModuleChecker(TypecheckOptions options, ts::Context base_context)
      : options_(options), base_context_(std::move(base_context)) {}

  bool check_module(ir::IrModule* mod, ir::IrTypedModule& out_typed, std::string& err) {
    // Step 1: Index cached imports by ModuleID.
    std::map<std::string, std::vector<ir::IrDecl>> decls_by_mod;
    for (const auto& decl : mod->import_decls) {
      decls_by_mod[decl.pos.filename].push_back(decl);
    }

    std::map<std::string, std::vector<ir::ModuleID>> deps_by_mod;
    for (const auto& imp : mod->imports) {
      deps_by_mod[imp.module_id.name] = imp.deps;
    }

    auto reachable_imports = [&](const std::vector<ir::IrImport>& direct) {
      std::vector<std::string> order;
      std::set<std::string> visited;
      std::vector<std::string> stack;
      for (const auto& imp : direct) {
        if (visited.insert(imp.module_id.name).second) {
          stack.push_back(imp.module_id.name);
        }
      }
      while (!stack.empty()) {
        std::string cur = stack.back();
        stack.pop_back();
        order.push_back(cur);
        auto it = deps_by_mod.find(cur);
        if (it != deps_by_mod.end()) {
          for (const auto& dep : it->second) {
            if (visited.insert(dep.name).second) {
              stack.push_back(dep.name);
            }
          }
        }
      }
      return order;
    };

    // Step 2: Collect & topologically sort module-local declarations.
    std::vector<ir::IrDecl> merged_decls;
    for (const auto& file : mod->files) {
      bool is_ffi = (file.case_val == ir::IrFileCase::Ffi);
      for (const auto& decl : file.decls) {
        merged_decls.push_back(decl);
        if (decl.is(ir::IrDeclCase::NameDecl) && decl.name) {
          local_types_.insert(decl.name->id);
        } else if (decl.is(ir::IrDeclCase::AliasDecl) && decl.alias) {
          local_types_.insert(decl.alias->id);
        } else if (!is_ffi && decl.is(ir::IrDeclCase::TermDecl) && decl.term) {
          auto& sym = symbols_[decl.term->id];
          if (sym.declared) {
            err = decl.pos.to_string() + ": symbol \"" + decl.term->id + "\" already declared";
            return false;
          }
          sym.decl = decl;
          sym.declared = true;
        }
      }
    }

    std::vector<ir::IrDecl> sorted_decls;
    if (!ir::topo_sort_decls(merged_decls, sorted_decls, err)) {
      return false;
    }

    // Step 3: Build per-file checking contexts (file_ctx[i]).
    std::vector<ts::Context> file_ctxs(mod->files.size());
    for (size_t i = 0; i < mod->files.size(); ++i) {
      auto& file = mod->files[i];
      if (file.case_val == ir::IrFileCase::Ffi) {
        continue;
      }

      ts::Context ctx = base_context_;
      auto mods = reachable_imports(file.imports);

      std::vector<ir::IrDecl> file_import_decls;
      for (const auto& m : mods) {
        auto it = decls_by_mod.find(m);
        if (it != decls_by_mod.end()) {
          file_import_decls.insert(file_import_decls.end(), it->second.begin(), it->second.end());
        }
      }
      if (!finalize_decls(file_import_decls, err)) {
        return false;
      }

      for (const auto& decl : file_import_decls) {
        if (decl.is(ir::IrDeclCase::ImplDecl) && decl.impl_decl) {
          if (!add_impl_decl_to_ctx(ctx, *decl.impl_decl, err)) return false;
        } else {
          if (!add_symbol_to_ctx(ctx, decl, err)) return false;
        }
      }

      try {
        ctx = ctx.enter_scope();
      } catch (const std::exception& e) {
        err = e.what();
        return false;
      }

      for (const auto& decl : sorted_decls) {
        if (!add_symbol_to_ctx(ctx, decl, err)) return false;
      }

      file_ctxs[i] = std::move(ctx);
    }

    // Step 3b: Reduce each module-local trait_impl in its defining file's context,
    // then register it into every non-FFI file's context.
    for (size_t j = 0; j < mod->files.size(); ++j) {
      auto& def_file = mod->files[j];
      if (def_file.case_val == ir::IrFileCase::Ffi) {
        continue;
      }
      for (auto& ti : def_file.trait_impls) {
        std::vector<ts::Binding> bindings;
        ir::IrDecl impl_decl = ti.decl();
        if (!build_impl_decl_bindings(file_ctxs[j], *impl_decl.impl_decl, bindings, err)) {
          return false;
        }
        for (size_t i = 0; i < mod->files.size(); ++i) {
          if (mod->files[i].case_val == ir::IrFileCase::Ffi) continue;
          for (const auto& b : bindings) {
            try {
              file_ctxs[i] = file_ctxs[i].add_bind(b);
            } catch (const std::exception& e) {
              err = e.what();
              return false;
            }
          }
        }
      }
    }

    // Step 3c: Validate declaration well-formedness in each file's own context.
    for (size_t i = 0; i < mod->files.size(); ++i) {
      auto& file = mod->files[i];
      if (file.case_val == ir::IrFileCase::Ffi) continue;
      for (const auto& decl : file.decls) {
        if (!ts::is_wellformed_decl(file_ctxs[i], decl, err)) {
          err = decl.pos.to_string() + ": " + err;
          return false;
        }
      }
    }

    // Step 4: Infer & verify bodies per file (in-place).
    for (size_t i = 0; i < mod->files.size(); ++i) {
      auto& file = mod->files[i];
      if (file.case_val == ir::IrFileCase::Ffi) {
        continue;
      }

      ts::Context& ctx = file_ctxs[i];
      for (auto& fn : file.functions) {
        ir::IrDecl decl = fn.decl();
        auto& sym = symbols_[decl.id()];
        sym.decl = decl;
        if (sym.defined) {
          err = decl.pos.to_string() + ": symbol \"" + decl.id() + "\" already defined";
          return false;
        }
        sym.defined = true;
        if (!check_function_in_ctx(ctx, &fn, err)) return false;
      }

      for (auto& ti : file.trait_impls) {
        if (!check_trait_impl_in_ctx(ctx, &ti, err)) return false;
      }
    }

    // Step 5: Check undefined terms across the module.
    if (!options_.skip_undefined_term_checks) {
      for (const auto& pair : symbols_) {
        if (pair.second.declared && !pair.second.defined) {
          err = pair.second.decl.pos.to_string() + ": symbol \"" + pair.first + "\" is declared but not defined";
          return false;
        }
      }
    }

    // Step 6: Flatten into IrTypedModule.
    out_typed = ir::IrTypedModule{};
    out_typed.module_id = mod->module_id;
    out_typed.imports = mod->imports;
    out_typed.impls = mod->impls;
    out_typed.flags = mod->flags;
    out_typed.import_decls = mod->import_decls;

    for (const auto& file : mod->files) {
      if (file.case_val == ir::IrFileCase::Ffi) {
        out_typed.ffi_decls.insert(out_typed.ffi_decls.end(), file.decls.begin(), file.decls.end());
      } else {
        out_typed.decls.insert(out_typed.decls.end(), file.decls.begin(), file.decls.end());
        out_typed.trait_impls.insert(out_typed.trait_impls.end(), file.trait_impls.begin(), file.trait_impls.end());
        out_typed.functions.insert(out_typed.functions.end(), file.functions.begin(), file.functions.end());
      }
    }

    out_typed.ffi_decls = ir::clean_decls(std::move(out_typed.ffi_decls));
    if (!finalize_decls(out_typed.decls, err)) {
      return false;
    }
    out_typed.trait_impls = ir::clean_trait_impls(std::move(out_typed.trait_impls));

    return true;
  }

 private:
  bool add_symbol_to_ctx(ts::Context& ctx, const ir::IrDecl& decl, std::string& err) {
    try {
      ctx = ctx.add_symbol(decl);
      return true;
    } catch (const std::exception& e) {
      err = e.what();
      return false;
    }
  }

  bool check_function_in_ctx(ts::Context& ctx, ir::IrFunction* function, std::string& err) {
    ts::Typechecker typechecker(ctx);
    ts::Context updated_ctx;

    if (options_.skip_term_typechecker) {
      if (!typechecker.infer_function(function, updated_ctx)) {
        err = typechecker.error();
        return false;
      }
      ctx = updated_ctx;
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
      ctx = out_ctx;
    }
    return true;
  }

  bool build_impl_decl_bindings(const ts::Context& def_ctx, const ir::ImplDeclData& impl, std::vector<ts::Binding>& out_bindings, std::string& err) {
    if (impl.case_val == ir::ImplCase::InherentImpl) {
      std::string base_name = ts::base_type_name(impl.type_name);
      if (base_name.empty()) {
        err = "invalid type for inherent impl: " + impl.type_name.to_string();
        return false;
      }
      for (const auto& m : impl.methods) {
        std::string name = base_name + "::" + m.id;
        std::vector<ir::IrType> args;
        for (const auto& a : m.args) {
          ir::IrType t = ir::substitute_type(a.type, ir::new_name_type("Self"), impl.type_name);
          args.push_back(std::move(t));
        }
        ir::IrType ret = ir::substitute_type(m.ret_type, ir::new_name_type("Self"), impl.type_name);
        ir::IrType method_type = ir::new_function_type(ir::new_tuple_type(std::move(args)), std::move(ret));

        std::vector<ir::TypeParam> tvars = impl.type_params;
        tvars.insert(tvars.end(), m.type_params.begin(), m.type_params.end());
        for (auto it = tvars.rbegin(); it != tvars.rend(); ++it) {
          method_type = ir::new_forall_type(*it, std::move(method_type));
        }
        out_bindings.push_back(ts::new_term_decl_bind(name, method_type));
      }
      return true;
    }

    ts::Context local_ctx = def_ctx;
    try {
      for (const auto& tp : impl.type_params) {
        local_ctx = local_ctx.add_bind(ts::new_type_param_bind(tp));
      }
      ts::Typechecker tc(local_ctx);
      ir::IrType reduced_type = tc.reduce_type(impl.type_name);
      ir::IrType reduced_trait_type = tc.reduce_type(impl.trait_type);
      out_bindings.push_back(ts::new_trait_impl_bind(impl.type_params, reduced_trait_type, reduced_type));
      return true;
    } catch (const std::exception& e) {
      err = e.what();
      return false;
    }
  }

  bool add_impl_decl_to_ctx(ts::Context& ctx, const ir::ImplDeclData& impl, std::string& err) {
    std::vector<ts::Binding> bindings;
    if (!build_impl_decl_bindings(ctx, impl, bindings, err)) {
      return false;
    }
    for (const auto& b : bindings) {
      try {
        ctx = ctx.add_bind(b);
      } catch (const std::exception& e) {
        err = e.what();
        return false;
      }
    }
    return true;
  }

  bool check_trait_impl_in_ctx(const ts::Context& ctx, ir::IrTraitImpl* impl, std::string& err) {
    std::string base_type = ts::base_type_name(impl->type_name);
    if (base_type.empty()) {
      err = impl->pos.to_string() + ": invalid type in impl: " + impl->type_name.to_string();
      return false;
    }
    if (local_types_.find(base_type) == local_types_.end() && base_type != "Ptr" && base_type != "StringView") {
      err = impl->pos.to_string() + ": coherence violation: cannot implement trait for foreign type \"" + base_type + "\"";
      return false;
    }

    ts::Context method_context = ctx;
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
        if (!options_.skip_term_typechecker) {
          if (!tc.typecheck_function(&m, out_ctx)) {
            err = tc.error();
            return false;
          }
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
    if (!ctx.lookup_trait_bind(trait_name, bind) || !bind.trait) {
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
      if (!options_.skip_term_typechecker) {
        if (!tc.typecheck_function(impl_method, out_ctx)) {
          err = tc.error();
          return false;
        }
      }
    }

    return true;
  }

  TypecheckOptions options_;
  ts::Context base_context_;
  std::map<std::string, SymbolStatus> symbols_;
  std::set<std::string> local_types_;
};

inline bool typecheck_module(TypecheckOptions options, ir::IrModule* mod, ir::IrTypedModule& out_typed, std::string& err) {
  ts::Context ctx = options.skip_default_context ? ts::Context{} : new_default_context();
  ModuleChecker checker(options, std::move(ctx));
  return checker.check_module(mod, out_typed, err);
}

inline bool typecheck_module_from_file(const Querier& querier, TypecheckOptions options, const std::string& base_filename, ir::IrTypedModule& out_typed, std::string& err) {
  ir::IrModule mod;
  if (!normalize_module(querier, base_filename, mod, err)) {
    return false;
  }
  return typecheck_module(options, &mod, out_typed, err);
}

} // namespace comp
