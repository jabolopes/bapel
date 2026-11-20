#pragma once

#include "bin/ir_base.h"
#include "bin/ir_decl.h"
#include "bin/ir_function.h"
#include "bin/ir_unit.h"
#include "comp/querier.h"
#include "ts/context.h"

#include <map>
#include <set>
#include <string>
#include <vector>

namespace comp {

struct ImportedModuleState {
  bool visiting = false;
  bool completed = false;
};

inline bool finalize_decls(std::vector<ir::IrDecl>& decls, std::string& err) {
  decls = ir::clean_decls(std::move(decls));
  return ir::topo_sort_decls(decls, decls, err);
}

class Resolver {
 public:
  explicit Resolver(Querier querier) : querier_(std::move(querier)) {}

  bool resolve_module(ir::IrModule& mod, std::string& err) {
    imported_modules_.clear();

    // 1. Resolve all impls declared in the base file and append them to mod.files
    // in the exact order they appear in the base file's impls section.
    std::vector<ir::Filename> rel_impls;
    for (const auto& imp : mod.impls) {
      rel_impls.push_back(imp.relative_filename);
    }
    mod.impls.clear();

    ir::Filename base_filename = mod.files.empty() ? ir::Filename{} : mod.files[0].filename;
    for (const auto& rel_impl : rel_impls) {
      ir::Filename impl_fn = querier_.impl_source_filename(base_filename, rel_impl);
      mod.impls.push_back(ir::new_impl(impl_fn));

      const std::string& path = impl_fn.value;
      bool is_bpl = (path.size() >= 4 && path.substr(path.size() - 4) == ".bpl");
      if (is_bpl) {
        auto res = parser::parse_source_file_from_file(path);
        if (!res.ok) {
          err = res.error().empty() ? ("failed to parse implementation file: " + path) : res.error();
          return false;
        }
        ast::DesugaredFile impl_desugared = ast::desugar_source_file(res.value);
        if (impl_desugared.module_id.name != mod.module_id.name) {
          err = "implementation file \"" + path + "\" declares module \"" + impl_desugared.module_id.name +
                "\", expected \"" + mod.module_id.name + "\"";
          return false;
        }
        if (!impl_desugared.impls.empty()) {
          err = "implementation file \"" + path + "\" cannot declare an impls section (only the base file may)";
          return false;
        }
        if (!impl_desugared.flags.empty()) {
          err = "implementation file \"" + path + "\" cannot declare a flags section (only the base file may)";
          return false;
        }
        impl_desugared.file.case_val = ir::IrFileCase::Impl;
        impl_desugared.file.filename = impl_fn;
        mod.files.push_back(std::move(impl_desugared.file));
      } else {
        SourceFileQuery sf_query;
        if (!query_source_file(path, sf_query, err)) {
          return false;
        }
        ir::IrFile f;
        f.case_val = ir::IrFileCase::Ffi;
        f.filename = impl_fn;
        f.decls = std::move(sf_query.decls);
        mod.files.push_back(std::move(f));
      }
    }

    // 2. For each file, add implicit function decls, clean per-file imports/decls/trait_impls,
    // and resolve all unique imported modules.
    for (auto& file : mod.files) {
      for (const auto& fn : file.functions) {
        file.decls.push_back(fn.decl());
      }
      file.imports = ir::clean_imports(std::move(file.imports));
      file.trait_impls = ir::clean_trait_impls(std::move(file.trait_impls));
      if (!finalize_decls(file.decls, err)) {
        return false;
      }

      for (const auto& imp : file.imports) {
        if (!resolve_module_import(imp.module_id, mod, err)) {
          return false;
        }
      }
    }

    std::sort(mod.imports.begin(), mod.imports.end(), [](const ir::IrResolvedImport& a, const ir::IrResolvedImport& b) {
      return a.module_id.name < b.module_id.name;
    });

    // 3. Module-wide trait coherence check:
    // Any trait implementation `impl Trait for Type` must be defined in the same module
    // as the `Type` definition.
    std::set<std::string> local_types;
    for (const auto& file : mod.files) {
      for (const auto& decl : file.decls) {
        if (decl.is(ir::IrDeclCase::NameDecl) && decl.name) {
          local_types.insert(decl.name->id);
        } else if (decl.is(ir::IrDeclCase::AliasDecl) && decl.alias) {
          local_types.insert(decl.alias->id);
        }
      }
    }

    for (const auto& file : mod.files) {
      for (const auto& impl : file.trait_impls) {
        if (impl.case_val == ir::ImplCase::TraitImpl) {
          std::string base_name = ts::base_type_name(impl.type_name);
          if (!base_name.empty() && base_name != "Ptr" && base_name != "StringView" &&
              local_types.find(base_name) == local_types.end()) {
            err = impl.pos.to_string() + ":\n trait coherence violation: 'impl " +
                  impl.trait_type.to_string() + " for " + impl.type_name.to_string() +
                  "' must be defined in the same module as " + base_name;
            return false;
          }
        }
      }
    }

    mod.imported_trait_impls = ir::clean_trait_impls(std::move(mod.imported_trait_impls));
    if (!finalize_decls(mod.import_decls, err)) {
      return false;
    }

    return true;
  }

 private:
  bool resolve_module_import(const ir::ModuleID& module_id, ir::IrModule& mod, std::string& err) {
    std::string norm_name = module_id.to_filename();
    auto it = imported_modules_.find(norm_name);
    if (it != imported_modules_.end()) {
      if (it->second.visiting) {
        err = "import cycle with module \"" + module_id.name + "\"";
        return false;
      }
      return true;
    }

    imported_modules_[norm_name] = ImportedModuleState{true, false};

    ModuleQuery module_query;
    if (!querier_.query_module_exports(module_id, module_query, err)) {
      return false;
    }

    ir::Filename imp_fn = querier_.base_source_filename(module_id);
    mod.imports.push_back(ir::IrResolvedImport{module_id, imp_fn, module_query.imports});

    for (const auto& dep : module_query.imports) {
      if (!resolve_module_import(dep, mod, err)) {
        return false;
      }
    }

    mod.import_decls.insert(mod.import_decls.end(), module_query.decls.begin(), module_query.decls.end());
    mod.import_decls = ir::clean_decls(std::move(mod.import_decls));
    mod.imported_trait_impls.insert(mod.imported_trait_impls.end(), module_query.trait_impls.begin(), module_query.trait_impls.end());
    mod.imported_trait_impls = ir::clean_trait_impls(std::move(mod.imported_trait_impls));

    imported_modules_[norm_name] = ImportedModuleState{false, true};
    return true;
  }

  Querier querier_;
  std::map<std::string, ImportedModuleState> imported_modules_;
};

inline bool resolve_module(Querier querier, ir::IrModule& mod, std::string& err) {
  Resolver resolver(std::move(querier));
  return resolver.resolve_module(mod, err);
}

} // namespace comp
