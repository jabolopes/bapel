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
    auto it = imported_modules_.find(module_id.name);
    if (it != imported_modules_.end()) {
      if (it->second.visiting) {
        err = "import cycle with module \"" + module_id.name + "\"";
        return false;
      }
      return true;
    }

    imported_modules_[module_id.name] = ImportedModuleState{true, false};

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
    unit.imported_trait_impls.insert(unit.imported_trait_impls.end(), module_query.trait_impls.begin(), module_query.trait_impls.end());

    imported_modules_[module_id.name] = ImportedModuleState{false, true};
    return true;
  }

  Querier querier_;
  std::map<std::string, ImportedModuleState> imported_modules_;
};

inline bool resolve_source_file(Querier querier, ir::IrUnit& unit, std::string& err) {
  Resolver resolver(std::move(querier));
  return resolver.resolve_source_unit(unit, err);
}

} // namespace comp
