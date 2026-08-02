#pragma once

#include "bin/ir_base.h"
#include "bin/ir_decl.h"
#include "bin/ir_unit.h"
#include "comp/module_finder.h"

#include <algorithm>
#include <fstream>
#include <map>
#include <sstream>
#include <string>
#include <vector>

namespace comp {

struct SourceFileQuery {
  std::vector<ir::ModuleID> imports;
  std::vector<ir::Filename> impls;
  std::vector<ir::Filename> flags;
  std::vector<ir::IrDecl> decls;
  std::vector<ir::IrTraitImpl> trait_impls;
};

struct ModuleQuery {
  std::vector<ir::ModuleID> imports;
  std::vector<ir::Filename> impls;
  std::vector<ir::Filename> flags;
  std::vector<ir::IrDecl> decls;
  std::vector<ir::IrTraitImpl> trait_impls;
};

class Querier {
 public:
  Querier() = default;
  explicit Querier(ModuleFinder finder) : finder_(std::move(finder)) {}

  ir::Filename base_source_filename(const ir::ModuleID& module_id) const {
    return finder_.base_source_filename(module_id);
  }

  ir::Filename impl_source_filename(const ir::Filename& base_filename, const ir::Filename& relative_impl_filename) const {
    return finder_.impl_source_filename(base_filename, relative_impl_filename);
  }

  bool query_module(const ir::ModuleID& module_id, ModuleQuery& out_query, std::string& err) const {
    auto it = cached_modules_.find(module_id.name);
    if (it != cached_modules_.end()) {
      out_query = it->second;
      return true;
    }
    out_query = ModuleQuery{};
    return true;
  }

  bool query_module_exports(const ir::ModuleID& module_id, ModuleQuery& out_query, std::string& err) const {
    if (!query_module(module_id, out_query, err)) {
      return false;
    }
    std::vector<ir::IrDecl> exported;
    for (const auto& d : out_query.decls) {
      if (d.export_flag) {
        exported.push_back(d);
      }
    }
    out_query.decls = std::move(exported);
    return true;
  }

  void register_module(const std::string& name, ModuleQuery query) {
    cached_modules_[name] = std::move(query);
  }

  const ModuleFinder& finder() const { return finder_; }

 private:
  ModuleFinder finder_;
  mutable std::map<std::string, ModuleQuery> cached_modules_;
};

} // namespace comp
