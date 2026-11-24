#pragma once

#include "ir_decl.h"
#include <algorithm>
#include <map>
#include <set>
#include <string>
#include <vector>

namespace ir {

struct IrImport {
  ModuleID module_id;
};

inline IrImport new_import(ModuleID id) {
  return IrImport{std::move(id)};
}

struct IrImpl {
  Filename relative_filename;
};

inline IrImpl new_impl(Filename filename) {
  return IrImpl{std::move(filename)};
}

struct IrResolvedImport {
  ModuleID module_id;
  Filename resolved_filename;
  std::vector<ModuleID> deps;
};

struct IrFile {
  IrFileCase case_val = IrFileCase::Base;
  Filename filename;
  std::vector<IrImport> imports;
  std::vector<IrDecl> decls;
  std::vector<IrTraitImpl> trait_impls;
  std::vector<IrFunction> functions;

  std::string case_string() const {
    switch (case_val) {
      case IrFileCase::Base: return "base";
      case IrFileCase::Impl: return "impl";
      case IrFileCase::Ffi: return "ffi";
    }
    return "base";
  }

  std::string to_bpl_string(bool with_pos = true) const {
    std::stringstream ss;
    ss << "file \"" << filename.value << "\" (" << case_string() << ") {\n";
    if (!imports.empty()) {
      ss << "  imports {\n";
      for (const auto& imp : imports) {
        ss << "    " << imp.module_id.to_string(with_pos) << "\n";
      }
      ss << "  }\n";
    }
    if (!decls.empty()) {
      ss << "  decls {\n";
      for (const auto& d : decls) {
        ss << "    " << d.to_string(with_pos) << "\n";
      }
      ss << "  }\n";
    }
    if (!trait_impls.empty()) {
      ss << "  traitImpls {\n";
      for (const auto& ti : trait_impls) {
        ss << "    " << ti.to_string(with_pos) << "\n";
      }
      ss << "  }\n";
    }
    if (!functions.empty()) {
      ss << "  functions {\n";
      for (const auto& f : functions) {
        ss << "    " << f.to_string(with_pos) << "\n";
      }
      ss << "  }\n";
    }
    ss << "}\n";
    return ss.str();
  }
};

struct IrModule {
  ModuleID module_id;
  std::vector<IrResolvedImport> imports;
  std::vector<IrImpl> impls;
  std::vector<std::string> flags;
  std::vector<IrDecl> import_decls;
  std::vector<IrFile> files;

  std::string to_bpl_string(bool with_pos = true) const {
    std::stringstream ss;
    ss << "module " << module_id.to_string(false) << "\n";

    if (!imports.empty()) {
      ss << "\nimports {\n";
      for (const auto& imp : imports) {
        ss << "  " << imp.module_id.to_string(false);
        if (!imp.resolved_filename.value.empty()) {
          ss << " \"" << imp.resolved_filename.value << "\"";
        }
        ss << "\n";
      }
      ss << "}\n";
    }

    if (!impls.empty()) {
      ss << "\nimpls {\n";
      for (const auto& impl : impls) {
        ss << "  \"" << impl.relative_filename.value << "\"\n";
      }
      ss << "}\n";
    }

    if (!flags.empty()) {
      ss << "\nflags {\n";
      for (const auto& flag : flags) {
        ss << "  \"" << flag << "\"\n";
      }
      ss << "}\n";
    }

    if (!import_decls.empty()) {
      ss << "\nimportDecls {\n";
      for (const auto& d : import_decls) {
        ss << "  " << d.to_string(with_pos) << "\n";
      }
      ss << "}\n";
    }

    for (const auto& f : files) {
      ss << "\n" << f.to_bpl_string(with_pos);
    }

    return ss.str();
  }
};

inline std::vector<IrImport> clean_imports(std::vector<IrImport> imports) {
  std::vector<IrImport> out;
  std::set<std::string> seen;
  for (auto& imp : imports) {
    std::string norm = imp.module_id.to_filename();
    if (seen.insert(norm).second) {
      out.push_back(std::move(imp));
    }
  }
  return out;
}

inline std::vector<IrDecl> clean_decls(std::vector<IrDecl> decls) {
  std::vector<IrDecl> out;
  std::set<std::string> seen;
  for (auto& d : decls) {
    if (seen.insert(d.id()).second) {
      out.push_back(std::move(d));
    }
  }
  return out;
}

inline std::vector<IrTraitImpl> clean_trait_impls(std::vector<IrTraitImpl> trait_impls) {
  std::vector<IrTraitImpl> out;
  std::set<std::string> seen;
  for (auto& ti : trait_impls) {
    std::string key = ti.to_string();
    if (seen.insert(key).second) {
      out.push_back(std::move(ti));
    }
  }
  return out;
}

inline void collect_name_types(const IrType& type, std::set<std::string>& out) {
  switch (type.tag()) {
    case IrTypeData::Tag::name_t:
      out.insert(type.data.as_name_t());
      break;
    case IrTypeData::Tag::app_t: {
      const auto& a = type.data.as_app_t();
      if (a.fun) collect_name_types(*a.fun, out);
      if (a.arg) collect_name_types(*a.arg, out);
      break;
    }
    case IrTypeData::Tag::fun_t: {
      const auto& f = type.data.as_fun_t();
      if (f.arg) collect_name_types(*f.arg, out);
      if (f.ret) collect_name_types(*f.ret, out);
      break;
    }
    case IrTypeData::Tag::forall_t: {
      const auto& fa = type.data.as_forall_t();
      if (fa.type) collect_name_types(*fa.type, out);
      break;
    }
    case IrTypeData::Tag::lambda_t: {
      const auto& l = type.data.as_lambda_t();
      if (l.type) collect_name_types(*l.type, out);
      break;
    }
    case IrTypeData::Tag::mu_t: {
      const auto& m = type.data.as_mu_t();
      if (m.body) collect_name_types(*m.body, out);
      break;
    }
    case IrTypeData::Tag::struct_t:
      for (const auto& f : type.fields()) {
        if (f.type) collect_name_types(*f.type, out);
      }
      break;
    case IrTypeData::Tag::tuple_t:
      for (const auto& elem : type.elems()) {
        collect_name_types(elem, out);
      }
      break;
    case IrTypeData::Tag::variant_t:
      for (const auto& tag : type.tags()) {
        if (tag.type) collect_name_types(*tag.type, out);
      }
      break;
    case IrTypeData::Tag::array_t: {
      const auto& arr = type.data.as_array_t();
      if (arr.elem_type) collect_name_types(*arr.elem_type, out);
      break;
    }
    case IrTypeData::Tag::var_t:
    case IrTypeData::Tag::exist_var_t:
      break;
  }
}

inline bool topo_sort_alias_decls(const std::vector<IrDecl>& decls, std::vector<IrDecl>& out, std::string& err) {
  std::map<std::string, size_t> name_to_id;
  for (size_t i = 0; i < decls.size(); ++i) {
    if (decls[i].is(IrDeclData::Tag::alias_decl)) {
      name_to_id[decls[i].data.as_alias_decl().id] = i;
    }
  }

  std::vector<std::vector<size_t>> adj(decls.size());
  std::vector<int> in_degree(decls.size(), 0);

  for (size_t i = 0; i < decls.size(); ++i) {
    if (!decls[i].is(IrDeclData::Tag::alias_decl)) continue;
    std::set<std::string> names;
    collect_name_types(decls[i].data.as_alias_decl().type, names);
    for (const auto& n : names) {
      auto it = name_to_id.find(n);
      if (it != name_to_id.end() && it->second != i) {
        adj[it->second].push_back(i);
        in_degree[i]++;
      }
    }
  }

  std::vector<size_t> queue;
  for (size_t i = 0; i < decls.size(); ++i) {
    if (in_degree[i] == 0) {
      queue.push_back(i);
    }
  }

  std::vector<IrDecl> sorted;
  while (!queue.empty()) {
    size_t u = queue.back();
    queue.pop_back();
    sorted.push_back(decls[u]);

    for (size_t v : adj[u]) {
      if (--in_degree[v] == 0) {
        queue.push_back(v);
      }
    }
  }

  if (sorted.size() != decls.size()) {
    err = "cyclic dependency between type declarations";
    return false;
  }

  out = std::move(sorted);
  return true;
}

inline bool topo_sort_decls(const std::vector<IrDecl>& decls, std::vector<IrDecl>& out, std::string& err) {
  std::vector<IrDecl> name_decls, alias_decls, trait_decls, term_decls, impl_decls;
  for (const auto& d : decls) {
    switch (d.tag()) {
      case IrDeclData::Tag::name_decl:
        name_decls.push_back(d);
        break;
      case IrDeclData::Tag::alias_decl:
        alias_decls.push_back(d);
        break;
      case IrDeclData::Tag::trait_decl:
        trait_decls.push_back(d);
        break;
      case IrDeclData::Tag::term_decl:
        term_decls.push_back(d);
        break;
      case IrDeclData::Tag::impl_decl:
        impl_decls.push_back(d);
        break;
    }
  }

  auto comp = [](const IrDecl& a, const IrDecl& b) {
    return a.id() < b.id();
  };
  std::sort(name_decls.begin(), name_decls.end(), comp);
  std::sort(trait_decls.begin(), trait_decls.end(), comp);
  std::sort(term_decls.begin(), term_decls.end(), comp);
  std::sort(impl_decls.begin(), impl_decls.end(), comp);

  std::vector<IrDecl> sorted_aliases;
  if (!topo_sort_alias_decls(alias_decls, sorted_aliases, err)) {
    return false;
  }

  out.clear();
  out.insert(out.end(), name_decls.begin(), name_decls.end());
  out.insert(out.end(), sorted_aliases.begin(), sorted_aliases.end());
  out.insert(out.end(), trait_decls.begin(), trait_decls.end());
  out.insert(out.end(), term_decls.begin(), term_decls.end());
  out.insert(out.end(), impl_decls.begin(), impl_decls.end());
  return true;
}

} // namespace ir
