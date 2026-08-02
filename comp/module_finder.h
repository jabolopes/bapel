#pragma once

#include "bin/ir_base.h"
#include <filesystem>
#include <map>
#include <string>
#include <utility>

namespace comp {

class ModuleFinder {
 public:
  ModuleFinder() = default;
  ModuleFinder(std::map<std::string, std::string> by_name, std::map<std::string, std::string> by_prefix)
      : modules_by_name_(std::move(by_name)), modules_by_prefix_(std::move(by_prefix)) {}

  bool lookup_module_by_name(const ir::ModuleID& module_id, std::string& out_pkg) const {
    auto it = modules_by_name_.find(module_id.name);
    if (it != modules_by_name_.end()) {
      out_pkg = it->second;
      return true;
    }
    return false;
  }

  bool lookup_module_by_prefix(const ir::ModuleID& module_id, std::string& out_pkg) const {
    std::string name = module_id.name;
    while (true) {
      size_t idx = name.rfind('.');
      if (idx == std::string::npos) {
        name.clear();
      } else {
        name = name.substr(0, idx);
      }

      auto it = modules_by_prefix_.find(name);
      if (it != modules_by_prefix_.end()) {
        out_pkg = it->second;
        return true;
      }

      if (name.empty()) {
        return false;
      }
    }
  }

  ir::Filename base_source_filename(const ir::ModuleID& module_id) const {
    std::string package_name;
    if (!lookup_module_by_name(module_id, package_name)) {
      lookup_module_by_prefix(module_id, package_name);
    }

    std::string to_file = module_id.to_filename();
    std::string full_path = package_name.empty() ? to_file : (package_name + "/" + to_file);
    return ir::new_filename(full_path + ".bpl", module_id.pos);
  }

  ir::Filename impl_source_filename(const ir::Filename& base_filename, const ir::Filename& relative_impl_filename) const {
    size_t idx = base_filename.value.rfind('/');
    std::string dir = (idx == std::string::npos) ? "" : base_filename.value.substr(0, idx);
    std::string full_path = dir.empty() ? relative_impl_filename.value : (dir + "/" + relative_impl_filename.value);
    full_path = std::filesystem::path(full_path).lexically_normal().string();
    return ir::new_filename(full_path, relative_impl_filename.pos);
  }

  void add_module(const std::string& name, const std::string& pkg) {
    modules_by_name_[name] = pkg;
  }

  void add_prefix(const std::string& prefix, const std::string& pkg) {
    modules_by_prefix_[prefix] = pkg;
  }

 private:
  std::map<std::string, std::string> modules_by_name_;
  std::map<std::string, std::string> modules_by_prefix_;
};

} // namespace comp
