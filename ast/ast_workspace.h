#pragma once

#include "ast_pos.h"
#include <sstream>
#include <string>
#include <vector>

namespace ast {

enum class PackageCase {
  ModulePackage = 0,
  PrefixPackage = 1,
};

struct Package {
  PackageCase case_val = PackageCase::ModulePackage;
  ir::ModuleID module_id; // Also used for Prefix
  ir::Filename filename;
  Pos pos;

  bool is(PackageCase c) const { return case_val == c; }

  std::string to_string(bool with_pos = false) const {
    std::stringstream ss;
    if (with_pos) ss << pos.to_string() << " ";
    if (case_val == PackageCase::ModulePackage) {
      ss << "module " << module_id.to_string(with_pos);
    } else {
      ss << "prefix " << module_id.to_string(with_pos);
    }
    ss << " in \"" << filename.value << "\"";
    return ss.str();
  }

  std::string to_json() const {
    std::stringstream ss;
    ss << "{\"Case\":" << static_cast<int>(case_val);
    if (case_val == PackageCase::ModulePackage) {
      ss << ",\"Module\":{\"ModuleID\":" << module_id.to_json() << "}";
    } else {
      ss << ",\"Prefix\":{\"Prefix\":" << module_id.to_json() << "}";
    }
    ss << ",\"Filename\":" << filename.to_json()
       << ",\"Pos\":" << pos.to_json() << "}";
    return ss.str();
  }
};

struct Packages {
  std::vector<Package> packages;
  Pos pos;

  std::string to_string(bool with_pos = false) const {
    std::stringstream ss;
    if (with_pos) ss << pos.to_string() << "\n";
    ss << "packages {\n";
    for (const auto& pkg : packages) {
      ss << "  " << pkg.to_string(with_pos) << "\n";
    }
    ss << "}\n";
    return ss.str();
  }

  std::string to_json() const {
    std::stringstream ss;
    ss << "{\"Packages\":[";
    ir::Interleave(packages, [&]() { ss << ","; }, [&](int, const Package& p) {
      ss << p.to_json();
    });
    ss << "],\"Pos\":" << pos.to_json() << "}";
    return ss.str();
  }
};

struct Workspace {
  Packages packages;

  std::string to_string(bool with_pos = false) const {
    std::stringstream ss;
    ss << "workspace {\n  " << packages.to_string(with_pos) << "}";
    return ss.str();
  }

  std::string to_json() const {
    std::stringstream ss;
    ss << "{\"Packages\":" << packages.to_json() << "}";
    return ss.str();
  }
};

inline bool validate_workspace(const Workspace& ws, std::vector<std::string>& out_errors) {
  bool ok = true;
  for (const auto& pkg : ws.packages.packages) {
    if (pkg.module_id.name.empty()) {
      out_errors.push_back("package rule has empty module or prefix name");
      ok = false;
    }
    if (pkg.filename.value.empty()) {
      out_errors.push_back("package rule has empty filename");
      ok = false;
    }
  }
  return ok;
}

inline Package new_module_package(ir::ModuleID module_id, ir::Filename filename, Pos pos) {
  Package p;
  p.case_val = PackageCase::ModulePackage;
  p.module_id = std::move(module_id);
  p.filename = std::move(filename);
  p.pos = pos;
  return p;
}

inline Package new_prefix_package(ir::ModuleID prefix, ir::Filename filename, Pos pos) {
  Package p;
  p.case_val = PackageCase::PrefixPackage;
  p.module_id = std::move(prefix);
  p.filename = std::move(filename);
  p.pos = pos;
  return p;
}

} // namespace ast
