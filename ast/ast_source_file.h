#pragma once

#include "ast_decl.h"
#include "ast_expr.h"
#include "ast_pos.h"
#include <algorithm>
#include <sstream>
#include <string>
#include <vector>

namespace ast {

enum class SourceFileCase {
  BaseSourceFile = 0,
  ImplSourceFile = 1,
};

struct SourceFileHeader {
  SourceFileCase case_val = SourceFileCase::BaseSourceFile;
  ir::ModuleID module_id;
  ir::Filename filename;

  bool is(SourceFileCase c) const { return case_val == c; }

  std::string to_string(bool with_pos = false) const {
    std::stringstream ss;
    if (case_val == SourceFileCase::BaseSourceFile) {
      ss << "module ";
    } else {
      ss << "implements ";
    }
    ss << module_id.to_string(with_pos);
    return ss.str();
  }

  std::string to_json() const {
    std::stringstream ss;
    ss << "{\"Case\":" << static_cast<int>(case_val)
       << ",\"ModuleID\":" << module_id.to_json()
       << ",\"Filename\":" << filename.to_json() << "}";
    return ss.str();
  }
};

struct Imports {
  std::vector<ir::ModuleID> ids;
  Pos pos;

  std::string to_string(bool with_pos = false) const {
    if (ids.empty()) return "";
    std::stringstream ss;
    if (with_pos) ss << pos.to_string() << "\n";
    ss << "imports {\n";
    for (const auto& id : ids) {
      ss << "  " << id.to_string(with_pos) << "\n";
    }
    ss << "}\n";
    return ss.str();
  }

  std::string to_json() const {
    std::stringstream ss;
    ss << "{\"IDs\":[";
    ir::Interleave(ids, [&]() { ss << ","; }, [&](int, const ir::ModuleID& id) {
      ss << id.to_json();
    });
    ss << "],\"Pos\":" << pos.to_json() << "}";
    return ss.str();
  }
};

struct Impls {
  std::vector<ir::Filename> filenames;
  Pos pos;

  std::string to_string(bool with_pos = false) const {
    if (filenames.empty()) return "";
    std::stringstream ss;
    if (with_pos) ss << pos.to_string() << "\n";
    ss << "impls {\n";
    for (const auto& f : filenames) {
      ss << "  " << f.to_string(with_pos) << "\n";
    }
    ss << "}\n";
    return ss.str();
  }

  std::string to_json() const {
    std::stringstream ss;
    ss << "{\"Filenames\":[";
    ir::Interleave(filenames, [&]() { ss << ","; }, [&](int, const ir::Filename& f) {
      ss << f.to_json();
    });
    ss << "],\"Pos\":" << pos.to_json() << "}";
    return ss.str();
  }
};

struct Flags {
  std::vector<ir::Filename> filenames;
  Pos pos;

  std::string to_string(bool with_pos = false) const {
    if (filenames.empty()) return "";
    std::stringstream ss;
    if (with_pos) ss << pos.to_string() << "\n";
    ss << "flags {\n";
    for (const auto& f : filenames) {
      ss << "  " << f.to_string(with_pos) << "\n";
    }
    ss << "}\n";
    return ss.str();
  }

  std::string to_json() const {
    std::stringstream ss;
    ss << "{\"Filenames\":[";
    ir::Interleave(filenames, [&]() { ss << ","; }, [&](int, const ir::Filename& f) {
      ss << f.to_json();
    });
    ss << "],\"Pos\":" << pos.to_json() << "}";
    return ss.str();
  }
};

struct SourceFile {
  SourceFileHeader header;
  Imports imports;
  Impls impls;
  Flags flags;
  std::vector<Source> body;

  std::string to_string(bool with_pos = false) const {
    std::stringstream ss;
    ss << header.to_string(with_pos) << "\n";
    if (!imports.ids.empty()) {
      ss << "\n" << imports.to_string(with_pos);
    }
    if (!impls.filenames.empty()) {
      ss << "\n" << impls.to_string(with_pos);
    }
    if (!flags.filenames.empty()) {
      ss << "\n" << flags.to_string(with_pos);
    }
    for (const auto& s : body) {
      ss << "\n" << s.to_string(with_pos) << "\n";
    }
    return ss.str();
  }

  std::string to_json() const {
    std::stringstream ss;
    ss << "{\"Header\":" << header.to_json()
       << ",\"Imports\":" << imports.to_json()
       << ",\"Impls\":" << impls.to_json()
       << ",\"Flags\":" << flags.to_json()
       << ",\"Body\":[";
    ir::Interleave(body, [&]() { ss << ","; }, [&](int, const Source& s) {
      ss << s.to_json();
    });
    ss << "]}";
    return ss.str();
  }
};

inline bool validate_source_file(const SourceFile& sf, std::vector<std::string>& out_errors) {
  bool ok = true;
  if (sf.header.is(SourceFileCase::ImplSourceFile)) {
    if (!sf.impls.filenames.empty()) {
      out_errors.push_back(
          "implementation file \"" + sf.header.module_id.name +
          "\" has an 'impls' section. The 'impls' section can only be used in base files");
      ok = false;
    }
  }

  std::vector<std::string> seen_impls;
  for (const auto& f : sf.impls.filenames) {
    if (std::find(seen_impls.begin(), seen_impls.end(), f.value) != seen_impls.end()) {
      out_errors.push_back("file \"" + sf.header.module_id.name +
                           "\" has an 'impls' section that contains duplicated implementation files");
      ok = false;
      break;
    }
    seen_impls.push_back(f.value);
  }

  return ok;
}

} // namespace ast
