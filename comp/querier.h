#pragma once

#include "ast/ast_decl.h"
#include "ast/ast_desugar.h"
#include "ast/ast_source_file.h"
#include "bin/ir_base.h"
#include "bin/ir_decl.h"
#include "bin/ir_parser.h"
#include "bin/ir_unit.h"
#include "comp/module_finder.h"

#include <algorithm>
#include <cstdio>
#include <fstream>
#include <map>
#include <sstream>
#include <string>
#include <vector>

namespace comp {

inline std::string find_parser_binary() {
  if (std::ifstream("./bootstrap/parser").good()) return "./bootstrap/parser";
  if (std::ifstream("../bootstrap/parser").good()) return "../bootstrap/parser";
  return "bootstrap/parser";
}

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

inline void replace_all_str(std::string& str, const std::string& from, const std::string& to) {
  size_t start_pos = 0;
  while ((start_pos = str.find(from, start_pos)) != std::string::npos) {
    str.replace(start_pos, from.length(), to);
    start_pos += to.length();
  }
}

inline std::string escape_shell_arg(const std::string& str) {
  std::string res = "'";
  for (char c : str) {
    if (c == '\'') {
      res += "'\\''";
    } else {
      res += c;
    }
  }
  res += "'";
  return res;
}

inline bool query_source_file(const std::string& filename, SourceFileQuery& out_query, std::string& err) {
  if (filename.size() >= 2 && filename.substr(filename.size() - 2) == ".h") {
    std::ifstream ifs(filename);
    if (!ifs) {
      err = "failed to open header file: " + filename;
      return false;
    }
    std::string line;
    while (std::getline(ifs, line)) {
      auto idx = line.find("@bpl:");
      if (idx == std::string::npos) continue;
      std::string decl_text = line.substr(idx + 5);
      size_t first = decl_text.find_first_not_of(" \t\r\n");
      if (first == std::string::npos) continue;
      decl_text = decl_text.substr(first);
      if (decl_text.rfind("export ", 0) == 0) {
        decl_text = "pub " + decl_text.substr(7);
      }
      replace_all_str(decl_text, "∗", "*");
      replace_all_str(decl_text, ":: * -> * -> *", "['a, 'b]");
      replace_all_str(decl_text, ":: * -> *", "['a]");
      
      std::string parser_bin = find_parser_binary();
      std::string cmd = "printf '%s\\n' " + escape_shell_arg(decl_text) + " | " + parser_bin + " --symbol=Decl --format=json 2>/dev/null";
      FILE* fp = popen(cmd.c_str(), "r");
      if (!fp) continue;
      char buf[2048];
      std::string json_str;
      while (fgets(buf, sizeof(buf), fp)) {
        json_str += buf;
      }
      pclose(fp);
      if (!json_str.empty()) {
        auto j = ir::JsonParser::parse(json_str);
        ir::IrDecl d = ir::deserialize_decl(j);
        out_query.decls.push_back(std::move(d));
      }
    }
    return true;
  }

  // Parse .bpl file via parser
  std::string parser_bin = find_parser_binary();
  std::string cmd = parser_bin + " --symbol=SourceFile --format=json " + filename + " 2>/dev/null";
  FILE* fp = popen(cmd.c_str(), "r");
  if (!fp) {
    err = "failed to parse source file: " + filename;
    return false;
  }
  char buf[4096];
  std::string json_str;
  while (fgets(buf, sizeof(buf), fp)) {
    json_str += buf;
  }
  int status = pclose(fp);
  if (status != 0 || json_str.empty()) {
    err = "failed to parse source file: " + filename;
    return false;
  }

  auto j = ir::JsonParser::parse(json_str);
  ast::SourceFile sf = ast::deserialize_ast_source_file(j);
  out_query.imports = sf.imports.ids;
  out_query.impls = sf.impls.filenames;
  out_query.flags = sf.flags.filenames;

  for (const auto& s : sf.body) {
    switch (s.case_val) {
      case ast::SourceCase::DeclSource:
        if (s.decl_data) out_query.decls.push_back(*s.decl_data);
        break;
      case ast::SourceCase::FunctionSource:
        if (s.function_data) out_query.decls.push_back(s.function_data->decl());
        break;
      case ast::SourceCase::TraitSource:
        if (s.trait_data) out_query.decls.push_back(s.trait_data->decl());
        break;
      case ast::SourceCase::ImplSource:
        if (s.impl_data) out_query.trait_impls.push_back(ast::desugar_impl(*s.impl_data));
        break;
    }
  }
  return true;
}

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

    ir::Filename base_fn = finder_.base_source_filename(module_id);
    if (base_fn.value.empty()) {
      err = "could not find module " + module_id.name;
      return false;
    }

    SourceFileQuery base_query;
    if (!query_source_file(base_fn.value, base_query, err)) {
      return false;
    }

    out_query.imports = base_query.imports;
    out_query.impls = base_query.impls;
    out_query.flags = base_query.flags;
    out_query.decls = base_query.decls;
    out_query.trait_impls = base_query.trait_impls;

    for (const auto& rel_impl : base_query.impls) {
      ir::Filename impl_fn = finder_.impl_source_filename(base_fn, rel_impl);
      SourceFileQuery impl_query;
      if (query_source_file(impl_fn.value, impl_query, err)) {
        out_query.decls.insert(out_query.decls.end(), impl_query.decls.begin(), impl_query.decls.end());
        out_query.trait_impls.insert(out_query.trait_impls.end(), impl_query.trait_impls.begin(), impl_query.trait_impls.end());
      }
    }

    for (auto& d : out_query.decls) {
      d.pos.filename = module_id.name;
    }

    cached_modules_[module_id.name] = out_query;
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
