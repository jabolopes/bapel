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

struct IrUnit {
  IrUnitCase case_val = IrUnitCase::BaseUnit;
  ModuleID module_id;
  Filename filename;
  std::vector<IrImport> imports;
  std::vector<IrImpl> impls;
  std::vector<IrDecl> import_decls;
  std::vector<IrDecl> impl_decls;
  std::vector<IrDecl> decls;
  std::vector<IrFunction> functions;
  std::vector<IrTraitImpl> trait_impls;
  std::vector<IrTraitImpl> imported_trait_impls;

  std::string to_json() const {
    std::stringstream ss;
    ss << "{\"Case\":" << static_cast<int>(case_val)
       << ",\"ModuleID\":" << module_id.to_json()
       << ",\"Filename\":" << filename.to_json()
       << ",\"Imports\":[";
    Interleave(imports, [&]() { ss << ","; }, [&](int, const IrImport& imp) {
      ss << "{\"ModuleID\":" << imp.module_id.to_json() << "}";
    });
    ss << "],\"Impls\":[";
    Interleave(impls, [&]() { ss << ","; }, [&](int, const IrImpl& imp) {
      ss << "{\"RelativeFilename\":" << imp.relative_filename.to_json() << "}";
    });
    ss << "],\"ImportDecls\":[";
    Interleave(import_decls, [&]() { ss << ","; }, [&](int, const IrDecl& d) {
      ss << d.to_json();
    });
    ss << "],\"ImplDecls\":[";
    Interleave(impl_decls, [&]() { ss << ","; }, [&](int, const IrDecl& d) {
      ss << d.to_json();
    });
    ss << "],\"Decls\":[";
    Interleave(decls, [&]() { ss << ","; }, [&](int, const IrDecl& d) {
      ss << d.to_json();
    });
    ss << "],\"Functions\":[";
    Interleave(functions, [&]() { ss << ","; }, [&](int, const IrFunction& f) {
      ss << f.to_json();
    });
    ss << "],\"TraitImpls\":[";
    Interleave(trait_impls, [&]() { ss << ","; }, [&](int, const IrTraitImpl& ti) {
      ss << ti.to_json();
    });
    ss << "],\"ImportedTraitImpls\":[";
    Interleave(imported_trait_impls, [&]() { ss << ","; }, [&](int, const IrTraitImpl& ti) {
      ss << ti.to_json();
    });
    ss << "]}";
    return ss.str();
  }

  std::string to_string() const {
    std::stringstream ss;
    ss << "MODULE " << module_id.to_string() << "\n";
    ss << (case_val == IrUnitCase::BaseUnit ? "CASE base\n" : "CASE impl\n");
    for (const auto& imp : imports) {
      ss << "IMPORT " << imp.module_id.to_string() << "\n";
    }
    for (const auto& impl : impls) {
      ss << "IMPL " << impl.relative_filename.value << "\n";
    }
    for (const auto& d : decls) {
      ss << "DECL " << d.to_string() << "\n";
    }
    for (const auto& ti : trait_impls) {
      ss << "TRAIT_IMPL " << ti.to_string() << "\n";
    }
    for (const auto& f : functions) {
      ss << "FUNC " << f.to_string() << "\n";
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
  switch (type.case_val) {
    case IrTypeCase::NameType:
      out.insert(type.name);
      break;
    case IrTypeCase::AppType:
      if (type.app) {
        if (type.app->fun) collect_name_types(*type.app->fun, out);
        if (type.app->arg) collect_name_types(*type.app->arg, out);
      }
      break;
    case IrTypeCase::FunType:
      if (type.fun) {
        if (type.fun->arg) collect_name_types(*type.fun->arg, out);
        if (type.fun->ret) collect_name_types(*type.fun->ret, out);
      }
      break;
    case IrTypeCase::ForallType:
      if (type.forall && type.forall->type) {
        collect_name_types(*type.forall->type, out);
      }
      break;
    case IrTypeCase::LambdaType:
      if (type.lambda && type.lambda->type) {
        collect_name_types(*type.lambda->type, out);
      }
      break;
    case IrTypeCase::StructType:
      for (const auto& f : type.fields()) {
        if (f.type) collect_name_types(*f.type, out);
      }
      break;
    case IrTypeCase::TupleType:
      for (const auto& elem : type.elems()) {
        collect_name_types(elem, out);
      }
      break;
    case IrTypeCase::VariantType:
      for (const auto& tag : type.tags()) {
        if (tag.type) collect_name_types(*tag.type, out);
      }
      break;
    case IrTypeCase::ArrayType:
      if (type.array && type.array->elem_type) {
        collect_name_types(*type.array->elem_type, out);
      }
      break;
    case IrTypeCase::VarType:
    case IrTypeCase::ExistVarType:
      break;
  }
}

inline bool topo_sort_alias_decls(const std::vector<IrDecl>& decls, std::vector<IrDecl>& out, std::string& err) {
  std::map<std::string, size_t> name_to_id;
  for (size_t i = 0; i < decls.size(); ++i) {
    if (decls[i].alias) {
      name_to_id[decls[i].alias->id] = i;
    }
  }

  std::vector<std::vector<size_t>> adj(decls.size());
  std::vector<int> in_degree(decls.size(), 0);

  for (size_t i = 0; i < decls.size(); ++i) {
    if (!decls[i].alias) continue;
    std::set<std::string> names;
    collect_name_types(decls[i].alias->type, names);
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
  std::vector<IrDecl> name_decls, alias_decls, trait_decls, term_decls;
  for (const auto& d : decls) {
    switch (d.case_val) {
      case IrDeclCase::NameDecl:
        name_decls.push_back(d);
        break;
      case IrDeclCase::AliasDecl:
        alias_decls.push_back(d);
        break;
      case IrDeclCase::TraitDecl:
        trait_decls.push_back(d);
        break;
      case IrDeclCase::TermDecl:
        term_decls.push_back(d);
        break;
    }
  }

  auto comp = [](const IrDecl& a, const IrDecl& b) {
    return a.id() < b.id();
  };
  std::sort(name_decls.begin(), name_decls.end(), comp);
  std::sort(trait_decls.begin(), trait_decls.end(), comp);
  std::sort(term_decls.begin(), term_decls.end(), comp);

  std::vector<IrDecl> sorted_aliases;
  if (!topo_sort_alias_decls(alias_decls, sorted_aliases, err)) {
    return false;
  }

  out.clear();
  out.insert(out.end(), name_decls.begin(), name_decls.end());
  out.insert(out.end(), sorted_aliases.begin(), sorted_aliases.end());
  out.insert(out.end(), trait_decls.begin(), trait_decls.end());
  out.insert(out.end(), term_decls.begin(), term_decls.end());
  return true;
}

} // namespace ir
