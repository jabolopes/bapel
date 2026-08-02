#pragma once

#include "bin/ir_base.h"
#include "bin/ir_decl.h"
#include "bin/ir_function.h"
#include "bin/ir_unit.h"
#include "comp/desugar.h"
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

class Resolver {
 public:
  explicit Resolver(Querier querier) : querier_(std::move(querier)) {}

  bool resolve_source_unit(ir::IrUnit& unit, std::string& err) {
    imported_modules_.clear();

    for (const auto& imp : unit.imports) {
      if (!resolve_import(imp.module_id, unit, err)) {
        return false;
      }
    }

    if (unit.case_val == ir::IrUnitCase::BaseUnit) {
      std::vector<ir::Filename> rel_impls;
      for (const auto& imp : unit.impls) {
        rel_impls.push_back(imp.relative_filename);
      }
      unit.impls.clear();
      if (!resolve_impls(unit.filename, rel_impls, unit, err)) {
        return false;
      }
    } else {
      if (!resolve_impl_source_file_impls(unit, err)) {
        return false;
      }
    }

    for (auto& fn : unit.functions) {
      fn = desugar_function(std::move(fn));
      unit.decls.push_back(fn.decl());
    }

    // Trait coherence check (Same-file / same-unit rule per AGENTS.md):
    // Any trait implementation `impl Trait for Type` must be defined in the same unit
    // as the `Type` definition.
    std::set<std::string> local_types;
    for (const auto& decl : unit.decls) {
      if (decl.is(ir::IrDeclCase::NameDecl) && decl.name) {
        local_types.insert(decl.name->id);
      } else if (decl.is(ir::IrDeclCase::AliasDecl) && decl.alias) {
        local_types.insert(decl.alias->id);
      }
    }

    for (const auto& impl : unit.trait_impls) {
      if (impl.case_val == ir::ImplCase::TraitImpl) {
        std::string base_name = ts::base_type_name(impl.type_name);
        // Primitive types and tuples cannot have manual impl blocks;
        // user-defined types must be declared in the local unit.
        if (!base_name.empty() && base_name != "Ptr" && base_name != "StringView" && local_types.find(base_name) == local_types.end()) {
          // If the type is not defined locally, check if it's an imported type
          // Trait coherence rule forbids implementing external traits for external types.
          err = impl.pos.to_string() + ":\n trait coherence violation: 'impl " +
                impl.trait_type.to_string() + " for " + impl.type_name.to_string() +
                "' must be defined in the same source file as " + base_name;
          return false;
        }
      }
    }

    unit.imports = ir::clean_imports(std::move(unit.imports));
    unit.import_decls = ir::clean_decls(std::move(unit.import_decls));
    unit.impl_decls = ir::clean_decls(std::move(unit.impl_decls));
    unit.decls = ir::clean_decls(std::move(unit.decls));
    unit.imported_trait_impls = ir::clean_trait_impls(std::move(unit.imported_trait_impls));
    unit.trait_impls = ir::clean_trait_impls(std::move(unit.trait_impls));

    if (!ir::topo_sort_decls(unit.import_decls, unit.import_decls, err)) {
      return false;
    }
    if (!ir::topo_sort_decls(unit.impl_decls, unit.impl_decls, err)) {
      return false;
    }
    if (!ir::topo_sort_decls(unit.decls, unit.decls, err)) {
      return false;
    }

    return true;
  }

 private:
  bool resolve_import(const ir::ModuleID& module_id, ir::IrUnit& unit, std::string& err) {
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

    for (const auto& dep : module_query.imports) {
      if (!resolve_import(dep, unit, err)) {
        return false;
      }
    }

    unit.import_decls.insert(unit.import_decls.end(), module_query.decls.begin(), module_query.decls.end());
    unit.import_decls = ir::clean_decls(std::move(unit.import_decls));
    unit.imported_trait_impls.insert(unit.imported_trait_impls.end(), module_query.trait_impls.begin(), module_query.trait_impls.end());
    unit.imported_trait_impls = ir::clean_trait_impls(std::move(unit.imported_trait_impls));

    imported_modules_[norm_name] = ImportedModuleState{false, true};
    return true;
  }

  bool resolve_impl(const ir::Filename& impl_filename, ir::IrUnit& unit, std::string& err) {
    SourceFileQuery sf_query;
    if (!query_source_file(impl_filename.value, sf_query, err)) {
      return false;
    }

    for (const auto& dep : sf_query.imports) {
      if (!resolve_import(dep, unit, err)) {
        return false;
      }
    }

    unit.impls.push_back(ir::new_impl(impl_filename));
    unit.impl_decls.insert(unit.impl_decls.end(), sf_query.decls.begin(), sf_query.decls.end());
    unit.imported_trait_impls.insert(unit.imported_trait_impls.end(), sf_query.trait_impls.begin(), sf_query.trait_impls.end());
    return true;
  }

  bool resolve_impls(const ir::Filename& base_filename, const std::vector<ir::Filename>& relative_impls, ir::IrUnit& unit, std::string& err) {
    for (const auto& rel_impl : relative_impls) {
      ir::Filename impl_fn = querier_.impl_source_filename(base_filename, rel_impl);
      if (!resolve_impl(impl_fn, unit, err)) {
        return false;
      }
    }
    return true;
  }

  bool resolve_base_decls(const ir::Filename& base_filename, ir::IrUnit& unit, std::string& err) {
    SourceFileQuery sf_query;
    if (!query_source_file(base_filename.value, sf_query, err)) {
      return false;
    }

    for (const auto& dep : sf_query.imports) {
      if (!resolve_import(dep, unit, err)) {
        return false;
      }
    }

    unit.impl_decls.insert(unit.impl_decls.end(), sf_query.decls.begin(), sf_query.decls.end());
    unit.imported_trait_impls.insert(unit.imported_trait_impls.end(), sf_query.trait_impls.begin(), sf_query.trait_impls.end());
    return true;
  }

  bool resolve_impl_source_file_impls(ir::IrUnit& unit, std::string& err) {
    ir::Filename base_filename = querier_.base_source_filename(unit.module_id);
    if (!resolve_base_decls(base_filename, unit, err)) {
      return false;
    }

    ModuleQuery mod_query;
    if (!querier_.query_module(unit.module_id, mod_query, err)) {
      return false;
    }

    std::string base_name = unit.filename.value;
    auto last_slash = base_name.rfind('/');
    if (last_slash != std::string::npos) {
      base_name = base_name.substr(last_slash + 1);
    }

    int index = -1;
    for (size_t i = 0; i < mod_query.impls.size(); ++i) {
      if (mod_query.impls[i].value == base_name) {
        index = static_cast<int>(i);
        break;
      }
    }

    if (index == -1) {
      err = "implementation file \"" + unit.filename.value + "\" belongs to module \"" + unit.module_id.name + "\" but is not in impls section";
      return false;
    }

    std::vector<ir::Filename> above_impls(mod_query.impls.begin(), mod_query.impls.begin() + index);
    return resolve_impls(base_filename, above_impls, unit, err);
  }

  Querier querier_;
  std::map<std::string, ImportedModuleState> imported_modules_;
};

inline bool resolve_source_file(Querier querier, ir::IrUnit& unit, std::string& err) {
  Resolver resolver(std::move(querier));
  return resolver.resolve_source_unit(unit, err);
}

} // namespace comp
