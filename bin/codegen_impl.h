#pragma once

#include "comp/typecheck_module.h"
#include "comp/session.h"
#include "comp/querier.h"
#include "comp/module_finder.h"
#include <algorithm>
#include <cstdio>
#include <cstdlib>
#include <fstream>
#include <functional>
#include <iostream>
#include <map>
#include <memory>
#include <set>
#include <sstream>
#include <string>
#include <vector>

namespace codegen {

inline std::string to_header_path(const ir::ModuleID& module_id) {
  std::string s = module_id.name;
  for (char& c : s) {
    if (c == '.') c = '/';
  }
  return s + ".h";
}

inline std::string inherent_cpp_name(const std::string& bapel_name) {
  auto pos = bapel_name.rfind(ir::NamespaceSeparator);
  if (pos == std::string::npos) {
    return "inherents::" + bapel_name;
  }
  std::string prefix = bapel_name.substr(0, pos);
  std::string name = bapel_name.substr(pos + ir::NamespaceSeparator.size());
  return prefix + "::inherents::" + name;
}

inline std::string trait_cpp_name(const std::string& bapel_name) {
  auto pos = bapel_name.rfind(ir::NamespaceSeparator);
  if (pos == std::string::npos) {
    return "traits::" + bapel_name;
  }
  std::string prefix = bapel_name.substr(0, pos);
  std::string name = bapel_name.substr(pos + ir::NamespaceSeparator.size());
  return prefix + "::traits::" + name;
}

inline std::string base_type_name(const ir::IrType& typ) {
  if (typ.is(ir::IrTypeData::Tag::name_t)) {
    return typ.data.as_name_t();
  }
  if (typ.is(ir::IrTypeData::Tag::app_t)) {
    return base_type_name(typ.app_fun());
  }
  return "";
}

inline int count_type_vars(const ir::IrKind& kind) {
  if (kind.is_arrow_kind() && kind.as_arrow_kind().right) {
    return 1 + count_type_vars(*kind.as_arrow_kind().right);
  }
  return 0;
}

inline int count_alias_type_vars(const ir::IrType& typ) {
  if (typ.is(ir::IrTypeData::Tag::lambda_t)) {
    return static_cast<int>(typ.lambda_vars().size());
  }
  return 0;
}

inline bool is_cpp_statement(const ir::IrTerm& term) {
  switch (term.tag()) {
    case ir::IrTermData::Tag::assign:
    case ir::IrTermData::Tag::block:
    case ir::IrTermData::Tag::let_t:
    case ir::IrTermData::Tag::match_t:
    case ir::IrTermData::Tag::return_t:
      return true;
    case ir::IrTermData::Tag::app_term: {
      auto [id, types, arg] = term.app_args();
      if (id.is(ir::IrTermData::Tag::var_t)) {
        const std::string& v = id.data.as_var_t().id;
        if (v == "ifthen" || v == "ifelse" || v == "core::for") {
          return true;
        }
      }
      return false;
    }
    default:
      return false;
  }
}

struct AnonymousType {
  ir::IrType name_type;
  ir::IrDecl decl;
};

class CppPrinter {
public:
  std::ostream& out;
  ir::PrinterMode mode = ir::PrinterMode::ModePublicHeader;
  int idgen = 0;
  ir::Position position = ir::Position::TypePosition;
  bool auto_type = false;
  bool last_term = false;
  std::string var_destination;
  std::map<std::string, std::string> mu_self_names_;
  bool in_variant_base_ = false;
  const ir::IrTypedModule* mod = nullptr;

  CppPrinter(std::ostream& output, ir::PrinterMode m, const ir::IrTypedModule* u)
      : out(output), mode(m), mod(u) {}

  std::string gen_id() {
    return "__v_" + std::to_string(idgen++);
  }

  bool is_type_decl(const ir::IrDecl& decl) const {
    return decl.is(ir::IrDeclData::Tag::name_decl) || decl.is(ir::IrDeclData::Tag::alias_decl);
  }

  bool find_decl(const std::string& id, ir::IrDecl& out_decl) const {
    if (!mod) return false;
    for (const auto& d : mod->decls) {
      if (d.id() == id) { out_decl = d; return true; }
    }
    for (const auto& d : mod->import_decls) {
      if (d.id() == id) { out_decl = d; return true; }
    }
    for (const auto& d : mod->ffi_decls) {
      if (d.id() == id) { out_decl = d; return true; }
    }
    return false;
  }

  std::string canonical_mu_hash(const ir::IrType& typ) const {
    if (typ.is(ir::IrTypeData::Tag::mu_t) && typ.data.as_mu_t().body) {
      const auto& m = typ.data.as_mu_t();
      ir::IrType norm = ir::new_mu_type(
          "__mu", m.kind,
          ir::substitute_type(*m.body, ir::new_var_type(m.var), ir::new_var_type("__mu")));
      return hash_type(norm);
    }
    return hash_type(typ);
  }

  ir::IrType normalize_mu_refs(const ir::IrType& typ) const {
    switch (typ.tag()) {
      case ir::IrTypeData::Tag::mu_t: {
        if (mod) {
          std::string h = canonical_mu_hash(typ);
          auto find_in = [&](const std::vector<ir::IrDecl>& list) -> std::string {
            for (const auto& d : list) {
              if (d.is(ir::IrDeclData::Tag::alias_decl) &&
                  d.data.as_alias_decl().type.is(ir::IrTypeData::Tag::mu_t) &&
                  canonical_mu_hash(d.data.as_alias_decl().type) == h) {
                return d.id();
              }
            }
            return "";
          };
          std::string name = find_in(mod->decls);
          if (name.empty()) name = find_in(mod->import_decls);
          if (name.empty()) name = find_in(mod->ffi_decls);
          if (!name.empty()) return ir::new_name_type(name);
        }
        return typ;
      }
      case ir::IrTypeData::Tag::app_t: {
        const auto& a = typ.data.as_app_t();
        if (a.fun && a.arg) {
          return ir::new_app_type(normalize_mu_refs(*a.fun), normalize_mu_refs(*a.arg));
        }
        return typ;
      }
      case ir::IrTypeData::Tag::array_t: {
        const auto& arr = typ.data.as_array_t();
        if (arr.elem_type) {
          return ir::new_array_type(normalize_mu_refs(*arr.elem_type), arr.size);
        }
        return typ;
      }
      case ir::IrTypeData::Tag::fun_t: {
        const auto& f = typ.data.as_fun_t();
        if (f.arg && f.ret) {
          return ir::new_function_type(normalize_mu_refs(*f.arg), normalize_mu_refs(*f.ret));
        }
        return typ;
      }
      case ir::IrTypeData::Tag::struct_t: {
        std::vector<ir::StructField> fields;
        for (const auto& f : typ.fields()) {
          ir::StructField nf = f;
          if (nf.type) nf.type = std::make_shared<ir::IrType>(normalize_mu_refs(*nf.type));
          fields.push_back(std::move(nf));
        }
        return ir::new_struct_type(std::move(fields));
      }
      case ir::IrTypeData::Tag::tuple_t: {
        std::vector<ir::IrType> elems;
        for (const auto& e : typ.elems()) elems.push_back(normalize_mu_refs(e));
        return ir::new_tuple_type(std::move(elems));
      }
      case ir::IrTypeData::Tag::variant_t: {
        std::vector<ir::VariantTag> tags;
        for (const auto& tag : typ.tags()) {
          ir::VariantTag nt = tag;
          if (nt.type) nt.type = std::make_shared<ir::IrType>(normalize_mu_refs(*nt.type));
          tags.push_back(std::move(nt));
        }
        return ir::new_variant_type(std::move(tags));
      }
      default:
        return typ;
    }
  }

  std::string canonical_hash_type(const ir::IrType& typ) const {
    if (typ.is(ir::IrTypeData::Tag::mu_t)) {
      return canonical_mu_hash(typ);
    }
    return hash_type(normalize_mu_refs(typ));
  }

  bool find_decl_for_type(const ir::IrType& typ, ir::IrDecl& out_decl) const {
    if (!mod) return false;
    std::string target_hash = canonical_hash_type(typ);
    auto check_list = [&](const std::vector<ir::IrDecl>& list) -> bool {
      for (const auto& d : list) {
        if (d.is(ir::IrDeclData::Tag::alias_decl) && canonical_hash_type(d.data.as_alias_decl().type) == target_hash) {
          out_decl = d;
          return true;
        }
      }
      return false;
    };
    if (check_list(mod->decls) || check_list(mod->import_decls) || check_list(mod->ffi_decls)) {
      return true;
    }

    // Fallback for MuType whose body had inner aliases (e.g. IrTypeData, Pos)
    // expanded during inference: match by top-level struct fields or variant tags.
    if (typ.is(ir::IrTypeData::Tag::mu_t) && typ.data.as_mu_t().body) {
      const ir::IrType& tbody = *typ.data.as_mu_t().body;
      auto match_mu_shape = [&](const std::vector<ir::IrDecl>& list) -> bool {
        for (const auto& d : list) {
          if (!d.is(ir::IrDeclData::Tag::alias_decl)) continue;
          const auto& ad = d.data.as_alias_decl();
          if (!ad.type.is(ir::IrTypeData::Tag::mu_t) || !ad.type.data.as_mu_t().body) {
            continue;
          }
          const ir::IrType& dbody = *ad.type.data.as_mu_t().body;
          if (tbody.is(ir::IrTypeData::Tag::struct_t) && dbody.is(ir::IrTypeData::Tag::struct_t) &&
              tbody.field_ids() == dbody.field_ids()) {
            out_decl = d;
            return true;
          }
          if (tbody.is(ir::IrTypeData::Tag::variant_t) && dbody.is(ir::IrTypeData::Tag::variant_t) &&
              tbody.tag_ids() == dbody.tag_ids()) {
            out_decl = d;
            return true;
          }
        }
        return false;
      };
      if (match_mu_shape(mod->decls) || match_mu_shape(mod->import_decls) || match_mu_shape(mod->ffi_decls)) {
        return true;
      }
    }
    return false;
  }

  bool find_trait_decl(const std::string& id, ir::IrDecl& out_decl) const {
    if (!mod) return false;
    auto check_list = [&](const std::vector<ir::IrDecl>& list) -> bool {
      for (const auto& d : list) {
        if (d.id() == id && d.is(ir::IrDeclData::Tag::trait_decl)) {
          out_decl = d;
          return true;
        }
      }
      return false;
    };
    return check_list(mod->decls) || check_list(mod->import_decls) || check_list(mod->ffi_decls);
  }

  bool has_function(const std::string& id) const {
    if (!mod) return false;
    for (const auto& f : mod->functions) {
      if (f.id == id) return true;
    }
    return false;
  }

  std::string to_id(const std::string& id) const {
    auto pos = id.find(ir::NamespaceSeparator);
    if (pos != std::string::npos) {
      auto last_pos = id.rfind(ir::NamespaceSeparator);
      std::string type_prefix = id.substr(0, last_pos);
      std::string method_name = id.substr(last_pos + ir::NamespaceSeparator.size());
      ir::IrDecl decl;
      if (find_decl(type_prefix, decl) && is_type_decl(decl)) {
        return "::" + inherent_cpp_name(type_prefix) + "::" + method_name;
      }
      return "::" + id;
    }
    return id;
  }

  std::string hash_type(const ir::IrType& typ) const {
    return std::to_string(std::hash<std::string>{}(typ.to_string()));
  }

  template <typename F>
  void with_bind_position(F&& cb) {
    auto old = position;
    position = ir::Position::BindPosition;
    cb();
    position = old;
  }

  template <typename F>
  void with_auto_type(bool val, F&& cb) {
    auto old = auto_type;
    auto_type = val;
    cb();
    auto_type = old;
  }

  template <typename F>
  void with_last_term(bool val, F&& cb) {
    auto old = last_term;
    last_term = val;
    cb();
    last_term = old;
  }

  template <typename F>
  void with_var_destination(const std::string& dest, F&& cb) {
    auto old = var_destination;
    var_destination = dest;
    with_last_term(true, std::forward<F>(cb));
    var_destination = old;
  }

  template <typename F>
  void handle_last_term(F&& cb) {
    if (!last_term) {
      cb();
      return;
    }
    if (var_destination.empty()) {
      out << "return ";
    } else {
      out << var_destination << " = ";
    }
    with_last_term(false, std::forward<F>(cb));
  }

  template <typename F>
  void print_in_namespace(const std::string& id, F&& cb) {
    auto pos = id.rfind(ir::NamespaceSeparator);
    if (pos == std::string::npos) {
      cb(id);
      return;
    }
    std::string ns = id.substr(0, pos);
    std::string base_id = id.substr(pos + ir::NamespaceSeparator.size());

    out << "namespace ";
    std::vector<std::string> tokens;
    size_t start = 0;
    while ((pos = ns.find(ir::NamespaceSeparator, start)) != std::string::npos) {
      tokens.push_back(ns.substr(start, pos - start));
      start = pos + ir::NamespaceSeparator.size();
    }
    tokens.push_back(ns.substr(start));

    ir::Interleave(tokens, [&]() { out << "::"; }, [&](int, const std::string& t) { out << t; });
    out << " { ";
    cb(base_id);
    out << " }";
  }

  void print_type(const ir::IrType& typ) {
    if (auto_type) {
      out << "auto";
      return;
    }

    switch (typ.tag()) {
      case ir::IrTypeData::Tag::app_t: {
        print_type(typ.app_fun());
        auto args = typ.app_args();
        out << "<";
        ir::Interleave(args, [&]() { out << ", "; }, [&](int, const ir::IrType& a) { print_type(a); });
        out << ">";
        break;
      }
      case ir::IrTypeData::Tag::array_t: {
        const auto& arr = typ.data.as_array_t();
        out << "std::array<";
        if (arr.elem_type) {
          print_type(*arr.elem_type);
        }
        out << ", " << arr.size << ">";
        break;
      }
      case ir::IrTypeData::Tag::exist_var_t:
        out << typ.to_string();
        break;
      case ir::IrTypeData::Tag::forall_t: {
        auto tvars = typ.forall_vars();
        out << "template <";
        ir::Interleave(tvars, [&]() { out << ", "; }, [&](int, const std::string& tv) { out << "typename " << tv; });
        out << "> ";
        print_type(typ.forall_body());
        break;
      }
      case ir::IrTypeData::Tag::fun_t: {
        const auto& fn = typ.data.as_fun_t();
        out << "std::function<";
        with_bind_position([&]() {
          if (fn.ret) print_type(*fn.ret);
        });
        out << "(";
        if (fn.arg) print_type(*fn.arg);
        out << ")>";
        break;
      }
      case ir::IrTypeData::Tag::name_t: {
        const std::string& n = typ.data.as_name_t();
        if (n == "bool") out << "bool";
        else if (n == "i8") out << "int8_t";
        else if (n == "i16") out << "int16_t";
        else if (n == "i32") out << "int32_t";
        else if (n == "i64") out << "int64_t";
        else if (n == "f32") out << "float";
        else if (n == "f64") out << "double";
        else if (n.rfind("__anonym_", 0) == 0) {
          out << n;
        } else {
          ir::IrDecl decl;
          if (n.find(ir::NamespaceSeparator) == std::string::npos && find_decl(n, decl)) {
            out << "::" << to_id(n);
          } else {
            out << to_id(n);
          }
        }
        break;
      }
      case ir::IrTypeData::Tag::struct_t: {
        ir::IrDecl decl;
        if (find_decl_for_type(typ, decl)) {
          if (decl.id().rfind("__anonym_", 0) == 0) {
            out << decl.id();
          } else {
            out << "::" << to_id(decl.id());
          }
          return;
        }
        // Match instantiated generic struct aliases (e.g. StructField<::IrType>).
        if (mod) {
          auto try_generic = [&](const std::vector<ir::IrDecl>& list) -> bool {
            for (const auto& d : list) {
              if (!d.is(ir::IrDeclData::Tag::alias_decl)) continue;
              const auto& ad = d.data.as_alias_decl();
              if (!ad.type.is(ir::IrTypeData::Tag::lambda_t)) {
                continue;
              }
              ir::IrType body = ad.type.lambda_body();
              if (!body.is(ir::IrTypeData::Tag::struct_t)) continue;
              auto lvars = ad.type.lambda_vars();
              std::set<std::string> var_set(lvars.begin(), lvars.end());
              std::map<std::string, ir::IrType> subs;
              if (ts::match_type(body, typ, var_set, subs) && subs.size() == lvars.size()) {
                out << "::" << to_id(d.id()) << "<";
                ir::Interleave(lvars, [&]() { out << ", "; }, [&](int, const std::string& v) {
                  print_type(subs[v]);
                });
                out << ">";
                return true;
              }
            }
            return false;
          };
          if (try_generic(mod->decls) || try_generic(mod->import_decls) || try_generic(mod->ffi_decls)) {
            return;
          }
        }
        out << "struct {";
        auto fields = typ.fields();
        ir::Interleave(fields, [&]() { out << " "; }, [&](int, const ir::StructField& f) {
          if (f.type) print_type(*f.type);
          out << " " << f.id << ";";
        });
        out << "}";
        break;
      }
      case ir::IrTypeData::Tag::tuple_t: {
        auto elems = typ.elems();
        if (position == ir::Position::TypePosition) {
          ir::Interleave(elems, [&]() { out << ", "; }, [&](int, const ir::IrType& e) { print_type(e); });
        } else {
          if (elems.empty()) {
            out << "std::monostate";
          } else if (elems.size() == 1) {
            print_type(elems[0]);
          } else {
            out << "std::tuple<";
            ir::Interleave(elems, [&]() { out << ", "; }, [&](int, const ir::IrType& e) { print_type(e); });
            out << ">";
          }
        }
        break;
      }
      case ir::IrTypeData::Tag::variant_t: {
        if (!in_variant_base_) {
          ir::IrDecl decl;
          if (find_decl_for_type(typ, decl)) {
            out << "::" << to_id(decl.id());
            return;
          }
        }
        bool saved = in_variant_base_;
        in_variant_base_ = false;
        out << "std::variant<";
        auto tags = typ.tags();
        ir::Interleave(tags, [&]() { out << ", "; }, [&](int, const ir::VariantTag& tag) {
          with_bind_position([&]() {
            if (tag.type) print_type(*tag.type);
          });
          out << "/* " << to_id(tag.id) << " */";
        });
        out << ">";
        in_variant_base_ = saved;
        break;
      }
      case ir::IrTypeData::Tag::mu_t: {
        ir::IrDecl decl;
        if (find_decl_for_type(typ, decl)) {
          out << "::" << to_id(decl.id());
          return;
        }
        if (typ.data.as_mu_t().body) {
          print_type(*typ.data.as_mu_t().body);
        }
        break;
      }
      case ir::IrTypeData::Tag::var_t: {
        const std::string& v = typ.data.as_var_t();
        auto it = mu_self_names_.find(v);
        if (it != mu_self_names_.end()) {
          out << it->second;
        } else {
          out << v;
        }
        break;
      }
      default:
        break;
    }
  }

  void print_app_type_term(const ir::IrTerm& term) {
    auto app_res = term.app_types();
    const auto& arg = app_res.first;
    const auto& types = app_res.second;
    if (arg.is(ir::IrTermData::Tag::const_t)) {
      out << "static_cast<";
      with_bind_position([&]() {
        ir::Interleave(types, [&]() { out << ", "; }, [&](int, const ir::IrType& t) { print_type(t); });
      });
      out << ">(";
      print_term(arg);
      out << ")";
      return;
    }

    print_term(arg);
    out << "<";
    with_bind_position([&]() {
      ir::Interleave(types, [&]() { out << ", "; }, [&](int, const ir::IrType& t) { print_type(t); });
    });
    out << ">";
  }

  void print_app_term_term(const ir::IrTerm& term) {
    auto app_res = term.app_args();
    const auto& id = std::get<0>(app_res);
    const auto& types = std::get<1>(app_res);
    const auto& arg = std::get<2>(app_res);

    if (id.is(ir::IrTermData::Tag::var_t)) {
      const std::string& vname = id.data.as_var_t().id;
      if (vname == "ifthen" && arg.is(ir::IrTermData::Tag::tuple_t) && arg.data.as_tuple_t().elems.size() >= 2) {
        const auto& elems = arg.data.as_tuple_t().elems;
        const auto& cond = elems[0];
        const auto& then_term = elems[1];
        out << "if (";
        with_last_term(false, [&]() { print_term(cond); });
        out << ") ";
        print_term(then_term);
        return;
      }
      if (vname == "ifelse" && arg.is(ir::IrTermData::Tag::tuple_t) && arg.data.as_tuple_t().elems.size() >= 3) {
        const auto& elems = arg.data.as_tuple_t().elems;
        const auto& cond = elems[0];
        const auto& then_term = elems[1];
        const auto& else_term = elems[2];
        out << "if (";
        with_last_term(false, [&]() { print_term(cond); });
        out << ") ";
        print_term(then_term);
        out << " else ";
        print_term(else_term);
        return;
      }
      if (vname == "core::for" && arg.is(ir::IrTermData::Tag::tuple_t) && arg.data.as_tuple_t().elems.size() >= 2) {
        const auto& elems = arg.data.as_tuple_t().elems;
        const auto& cond = elems[0];
        auto fn_res = elems[1].to_function();
        const auto& body = std::get<2>(fn_res);
        out << "while (";
        with_last_term(false, [&]() { print_term(cond); });
        out << ") ";
        with_last_term(false, [&]() { print_term(body); });
        return;
      }
      if (ir::is_operator(vname)) {
        if (arg.is(ir::IrTermData::Tag::tuple_t) && !arg.data.as_tuple_t().elems.empty()) {
          const auto& elems = arg.data.as_tuple_t().elems;
          out << "(";
          print_term(elems[0]);
          out << ")";
          print_term(id);
          out << "(";
          std::vector<ir::IrTerm> rest(elems.begin() + 1, elems.end());
          print_term(ir::new_tuple_term(rest));
          out << ")";
        } else {
          print_term(id);
          out << " ";
          print_term(arg);
        }
        return;
      }

      // Trait method resolution
      auto pos = vname.rfind(ir::NamespaceSeparator);
      if (pos != std::string::npos) {
        std::string trait_name = vname.substr(0, pos);
        std::string method_name = vname.substr(pos + ir::NamespaceSeparator.size());
        ir::IrDecl decl;
        if (find_trait_decl(trait_name, decl)) {
          out << "::" << trait_cpp_name(trait_name) << "<";
          with_bind_position([&]() {
            ir::Interleave(types, [&]() { out << ", "; }, [&](int, const ir::IrType& t) { print_type(t); });
          });
          out << ">::" << method_name << "(";
          if (arg.is(ir::IrTermData::Tag::tuple_t)) {
            ir::Interleave(arg.data.as_tuple_t().elems, [&]() { out << ", "; }, [&](int, const ir::IrTerm& t) { print_term(t); });
          } else {
            print_term(arg);
          }
          out << ")";
          return;
        }

        // Inherent method resolution
        if (find_decl(trait_name, decl) && (decl.is(ir::IrDeclData::Tag::name_decl) || decl.is(ir::IrDeclData::Tag::alias_decl))) {
          int arity = 0;
          if (decl.is(ir::IrDeclData::Tag::name_decl)) {
            arity = count_type_vars(decl.data.as_name_decl().kind);
          } else {
            arity = count_alias_type_vars(decl.data.as_alias_decl().type);
          }

          if (static_cast<int>(types.size()) >= arity) {
            std::vector<ir::IrType> type_args(types.begin(), types.begin() + arity);
            std::vector<ir::IrType> method_args(types.begin() + arity, types.end());
            std::string mapped_type = is_type_decl(decl) ? inherent_cpp_name(trait_name) : trait_name;

            out << "::" << mapped_type;
            if (arity > 0) {
              out << "<";
              with_bind_position([&]() {
                ir::Interleave(type_args, [&]() { out << ", "; }, [&](int, const ir::IrType& t) { print_type(t); });
              });
              out << ">";
            }
            out << "::" << method_name;
            if (!method_args.empty()) {
              out << "<";
              with_bind_position([&]() {
                ir::Interleave(method_args, [&]() { out << ", "; }, [&](int, const ir::IrType& t) { print_type(t); });
              });
              out << ">";
            }
            out << "(";
            if (arg.is(ir::IrTermData::Tag::tuple_t)) {
              ir::Interleave(arg.data.as_tuple_t().elems, [&]() { out << ", "; }, [&](int, const ir::IrTerm& t) { print_term(t); });
            } else {
              print_term(arg);
            }
            out << ")";
            return;
          }
        }
      }
    }

    print_term(id);
    if (id.is(ir::IrTermData::Tag::var_t) && !ir::is_operator(id.data.as_var_t().id) && !types.empty()) {
      out << "<";
      with_bind_position([&]() {
        ir::Interleave(types, [&]() { out << ", "; }, [&](int, const ir::IrType& t) { print_type(t); });
      });
      out << ">";
    }

    out << "(";
    if (arg.is(ir::IrTermData::Tag::tuple_t)) {
      ir::Interleave(arg.data.as_tuple_t().elems, [&]() { out << ", "; }, [&](int, const ir::IrTerm& t) { print_term(t); });
    } else {
      print_term(arg);
    }
    out << ")";
  }

  void print_let_term(const ir::IrTerm& term) {
    const auto& c = term.data.as_let_t();
    bool auto_flag = c.value && c.value->is(ir::IrTermData::Tag::type_abs);

    with_auto_type(auto_flag, [&]() {
      with_bind_position([&]() {
        if (c.var_type) {
          print_type(*c.var_type);
          out << " ";
        }
        out << c.var;
      });
    });

    if (c.value && is_cpp_statement(*c.value)) {
      with_var_destination(c.var, [&]() {
        out << ";\n";
        print_term(*c.value);
      });
    } else {
      out << " = ";
      if (c.value) print_term(*c.value);
    }
  }

  void print_match_term(const ir::IrTerm& term) {
    const auto& c = term.data.as_match_t();
    std::string variant_id = gen_id();

    out << "{\n";
    out << "auto " << variant_id << " = ";
    with_last_term(false, [&]() {
      if (c.term) print_term(*c.term);
    });
    out << ";\n";
    out << "switch (" << variant_id << ".index()) {\n";
    for (const auto& arm : c.arms) {
      if (arm.is_wildcard()) {
        out << "default: {\n";
        if (arm.body) print_term(*arm.body);
        out << ";\n";
        out << "}\n";
      } else {
        int idx = arm.index.value_or(0);
        out << "case " << idx << ": {\n";
        out << "auto &" << arm.arg << " = std::get<" << idx << ">(" << variant_id << ");\n";
        if (arm.body) print_term(*arm.body);
        out << ";\n";
        out << "}\n";
      }
    }
    out << "} }";
  }

  void print_projection_term(const ir::IrTerm& term) {
    const auto& c = term.data.as_projection();
    ir::IrType obj_type = term.type.value_or(ir::IrType{});
    if (c.reduced_type) obj_type = *c.reduced_type;
    if (obj_type.is(ir::IrTypeData::Tag::mu_t) && obj_type.data.as_mu_t().body) {
      obj_type = *obj_type.data.as_mu_t().body;
    }

    if (obj_type.is(ir::IrTypeData::Tag::struct_t)) {
      int idx = 0;
      ir::StructField field;
      if (obj_type.field_by_label(c.label, idx, field)) {
        if (c.term) print_term(*c.term);
        out << "." << field.id;
        return;
      }
    }

    int index = 0;
    ir::IrType elem;
    if (obj_type.is(ir::IrTypeData::Tag::tuple_t) && obj_type.elem_by_label(c.label, index, elem)) {
      out << "std::get<" << index << ">(";
      if (c.term) print_term(*c.term);
      out << ")";
      return;
    }

    if (obj_type.is(ir::IrTypeData::Tag::variant_t) && obj_type.tag_by_label(c.label, index)) {
      out << "std::get<" << index << ">(";
      if (c.term) print_term(*c.term);
      out << ")";
      return;
    }

    out << "std::get<" << c.label << ">(";
    if (c.term) print_term(*c.term);
    out << ")";
  }

  void print_set_term(const ir::IrTerm& term) {
    const auto& c = term.data.as_set_t();
    ir::IrType obj_type = term.type.value_or(ir::IrType{});
    if (c.reduced_type) obj_type = *c.reduced_type;
    if (obj_type.is(ir::IrTypeData::Tag::mu_t) && obj_type.data.as_mu_t().body) {
      obj_type = *obj_type.data.as_mu_t().body;
    }

    if (obj_type.is(ir::IrTypeData::Tag::struct_t)) {
      std::string struct_id = gen_id();
      out << "([&, " << struct_id << " = ";
      if (c.term) print_term(*c.term);
      out << "]() mutable {\n";
      for (const auto& lv : c.values) {
        int idx = 0;
        ir::StructField field;
        if (obj_type.field_by_label(lv.label, idx, field)) {
          out << struct_id << "." << field.id << " = ";
          if (lv.value) print_term(*lv.value);
          out << ";\n";
        }
      }
      out << "return " << struct_id << ";\n";
      out << "})()";
      return;
    }

    if (obj_type.is(ir::IrTypeData::Tag::tuple_t)) {
      std::string tuple_id = gen_id();
      out << "([" << tuple_id << " = ";
      if (c.term) print_term(*c.term);
      out << "]() mutable {\n";
      for (const auto& lv : c.values) {
        out << "std::get<" << lv.label << ">(" << tuple_id << ") = ";
        if (lv.value) print_term(*lv.value);
        out << ";\n";
      }
      out << "return " << tuple_id << ";\n";
      out << "})()";
      return;
    }
  }

  void print_lambda_term(const ir::IrTerm& term) {
    auto fn_res = term.to_function();
    const auto& tvars = std::get<0>(fn_res);
    const auto& args = std::get<1>(fn_res);
    const auto& body = std::get<2>(fn_res);
    out << "[&]";

    if (!tvars.empty()) {
      out << "<";
      ir::Interleave(tvars, [&]() { out << ", "; }, [&](int, const ir::TypeParam& tp) {
        out << "typename " << tp.var;
      });
      out << ">";
    }

    bool argless = args.size() == 1 && args[0].type.is(ir::IrTypeData::Tag::tuple_t) && args[0].type.elems().empty();
    out << "(";
    if (!argless) {
      ir::Interleave(args, [&]() { out << ", "; }, [&](int, const ir::FunctionArg& arg) {
        print_type(arg.type);
        out << " " << to_id(arg.id);
      });
    }
    out << ")";

    with_last_term(true, [&]() { print_term(body); });
  }

  void print_term(const ir::IrTerm& term) {
    switch (term.tag()) {
      case ir::IrTermData::Tag::app_type:
        handle_last_term([&]() { print_app_type_term(term); });
        break;
      case ir::IrTermData::Tag::app_term:
        if (is_cpp_statement(term)) {
          print_app_term_term(term);
        } else {
          handle_last_term([&]() { print_app_term_term(term); });
        }
        break;
      case ir::IrTermData::Tag::assign:
        handle_last_term([&]() {
          const auto& a = term.data.as_assign();
          with_bind_position([&]() { if (a.ret) print_term(*a.ret); });
          out << " = ";
          if (a.arg) print_term(*a.arg);
        });
        break;
      case ir::IrTermData::Tag::block: {
        const auto& b = term.data.as_block();
        out << "{\n";
        for (size_t i = 0; i < b.terms.size(); ++i) {
          bool is_last = (i + 1 == b.terms.size());
          if (is_last) {
            print_term(b.terms[i]);
            out << ";}\n";
          } else {
            with_last_term(false, [&]() {
              print_term(b.terms[i]);
              out << ";\n";
            });
          }
        }
        break;
      }
      case ir::IrTermData::Tag::const_t:
        handle_last_term([&]() {
          const auto& lit = term.data.as_const_t().literal;
          switch (lit.data.tag()) {
            case ir::IrLiteralData::Tag::int_lit: out << lit.data.as_int_lit(); break;
            case ir::IrLiteralData::Tag::float_lit: out << lit.data.as_float_lit().integer << "." << lit.data.as_float_lit().decimal; break;
            case ir::IrLiteralData::Tag::rune_lit: out << "'" << lit.data.as_rune_lit() << "'"; break;
            case ir::IrLiteralData::Tag::str_lit: out << "\"" << lit.data.as_str_lit() << "\""; break;
          }
        });
        break;
      case ir::IrTermData::Tag::injection:
        handle_last_term([&]() {
          const auto& inj = term.data.as_injection();
          print_type(inj.variant_type);
          out << "(";
          if (inj.tag_index.has_value()) {
            out << "std::in_place_index<" << *inj.tag_index << ">, ";
          }
          if (inj.value) {
            if (inj.value->is(ir::IrTermData::Tag::struct_t)) {
              ir::IrType tag_type;
              bool found_tag_type = false;
              if (inj.tag_index.has_value()) {
                std::string vname = base_type_name(inj.variant_type);
                ir::IrDecl vdecl;
                if (!vname.empty() && find_decl(vname, vdecl) && vdecl.is(ir::IrDeclData::Tag::alias_decl)) {
                  ir::IrType vbody = vdecl.data.as_alias_decl().type;
                  std::vector<std::string> vparams;
                  while (vbody.is(ir::IrTypeData::Tag::lambda_t) && vbody.data.as_lambda_t().type) {
                    vparams.push_back(vbody.data.as_lambda_t().var);
                    vbody = *vbody.data.as_lambda_t().type;
                  }
                  if (vbody.is(ir::IrTypeData::Tag::mu_t) && vbody.data.as_mu_t().body) {
                    vbody = *vbody.data.as_mu_t().body;
                  }
                  auto tags = vbody.tags();
                  int idx = *inj.tag_index;
                  if (idx >= 0 && idx < static_cast<int>(tags.size()) && tags[idx].type) {
                    tag_type = *tags[idx].type;
                    auto vargs = inj.variant_type.app_args();
                    for (size_t i = 0; i < vparams.size() && i < vargs.size(); ++i) {
                      tag_type = ir::substitute_type(tag_type, ir::new_var_type(vparams[i]), vargs[i]);
                    }
                    found_tag_type = true;
                  }
                }
              }
              if (found_tag_type) {
                with_bind_position([&]() { print_type(tag_type); });
              } else if (inj.value->type.has_value()) {
                with_bind_position([&]() { print_type(*inj.value->type); });
              }
            }
            print_term(*inj.value);
          }
          out << ")";
        });
        break;
      case ir::IrTermData::Tag::lambda:
      case ir::IrTermData::Tag::type_abs:
        handle_last_term([&]() {
          with_last_term(false, [&]() { print_lambda_term(term); });
        });
        break;
      case ir::IrTermData::Tag::let_t:
        with_last_term(false, [&]() { print_let_term(term); });
        break;
      case ir::IrTermData::Tag::match_t:
        print_match_term(term);
        break;
      case ir::IrTermData::Tag::projection:
        handle_last_term([&]() { print_projection_term(term); });
        break;
      case ir::IrTermData::Tag::return_t:
        with_last_term(false, [&]() {
          const auto& r = term.data.as_return_t();
          if (r.expr) {
            out << "return ";
            print_term(*r.expr);
            out << ";";
          }
        });
        break;
      case ir::IrTermData::Tag::tuple_t:
        handle_last_term([&]() {
          const auto& elems = term.data.as_tuple_t().elems;
          if (position == ir::Position::BindPosition) {
            out << "std::tie(";
          } else if (elems.empty()) {
            out << "std::monostate(";
          } else {
            out << "std::make_tuple(";
          }
          ir::Interleave(elems, [&]() { out << ", "; }, [&](int, const ir::IrTerm& t) {
            print_term(t);
          });
          out << ")";
        });
        break;
      case ir::IrTermData::Tag::set_t:
        handle_last_term([&]() { print_set_term(term); });
        break;
      case ir::IrTermData::Tag::struct_t:
        handle_last_term([&]() {
          out << "{";
          ir::Interleave(term.data.as_struct_t().values, [&]() { out << ", "; }, [&](int, const ir::LabelValue& lv) {
            out << "." << lv.label << " = ";
            if (lv.value) print_term(*lv.value);
          });
          out << "}";
        });
        break;
      case ir::IrTermData::Tag::var_t:
        handle_last_term([&]() {
          out << to_id(term.data.as_var_t().id);
        });
        break;
    }
  }

  void print_template_params(const std::vector<ir::TypeParam>& tvars, bool is_definition) {
    if (tvars.empty()) return;
    out << "template <";
    ir::Interleave(tvars, [&]() { out << ", "; }, [&](int, const ir::TypeParam& tp) {
      out << "typename " << tp.var;
    });

    std::vector<std::string> constraints;
    for (const auto& tp : tvars) {
      for (const auto& bound : tp.bounds) {
        constraints.push_back(sfinae_constraint(tp.var, bound));
      }
    }

    if (!constraints.empty()) {
      if (is_definition) {
        out << ", typename";
      } else {
        out << ", typename = std::enable_if_t<";
        for (size_t i = 0; i < constraints.size(); ++i) {
          if (i > 0) out << " && ";
          out << constraints[i];
        }
        out << ">";
      }
    }
    out << "> ";
  }

  std::string sfinae_constraint(const std::string& vk_var, const ir::IrType& bound) {
    std::string trait_name;
    std::vector<ir::IrType> args;
    if (bound.is(ir::IrTypeData::Tag::name_t)) {
      trait_name = bound.data.as_name_t();
    } else if (bound.is(ir::IrTypeData::Tag::app_t)) {
      trait_name = bound.app_fun().trait_name();
      args = bound.app_args();
    }

    std::string cpp_name = "::" + trait_cpp_name(trait_name);
    std::ostringstream ss;
    ss << "(sizeof(" << cpp_name << "<" << vk_var;
    for (const auto& arg : args) {
      ss << ", ";
      std::ostringstream temp_out;
      CppPrinter temp_printer(temp_out, mode, mod);
      temp_printer.print_type(arg);
      ss << temp_out.str();
    }
    ss << ">) > 0)";
    return ss.str();
  }

  void print_alias_decl(const std::string& id, const ir::IrType& typ, const std::vector<std::string>& tvars = {}) {
    switch (typ.tag()) {
      case ir::IrTypeData::Tag::lambda_t: {
        auto lvars = typ.lambda_vars();
        out << "template <";
        ir::Interleave(lvars, [&]() { out << ", "; }, [&](int, const std::string& tv) { out << "typename " << tv; });
        out << "> ";
        print_alias_decl(id, typ.lambda_body(), lvars);
        break;
      }
      case ir::IrTypeData::Tag::mu_t: {
        const auto& m = typ.data.as_mu_t();
        if (!m.body) break;
        std::string self_name = "::" + to_id(id);
        if (!tvars.empty()) {
          self_name += "<";
          for (size_t i = 0; i < tvars.size(); ++i) {
            if (i > 0) self_name += ", ";
            self_name += tvars[i];
          }
          self_name += ">";
        }
        auto prev = mu_self_names_.find(m.var);
        bool had_prev = (prev != mu_self_names_.end());
        std::string prev_val = had_prev ? prev->second : "";
        mu_self_names_[m.var] = self_name;
        print_alias_decl(id, *m.body, tvars);
        if (had_prev) {
          mu_self_names_[m.var] = prev_val;
        } else {
          mu_self_names_.erase(m.var);
        }
        break;
      }
      case ir::IrTypeData::Tag::name_t:
        out << "using " << id << " = ";
        print_type(typ);
        break;
      case ir::IrTypeData::Tag::struct_t:
        out << "struct " << id << " {\n";
        for (const auto& f : typ.fields()) {
          if (f.type) print_type(*f.type);
          out << " " << f.id << ";\n";
        }
        out << "}";
        break;
      case ir::IrTypeData::Tag::tuple_t:
        out << "struct " << id << " : ";
        with_bind_position([&]() { print_type(typ); });
        out << " { " << id << "(const ";
        with_bind_position([&]() { print_type(typ); });
        out << "& arg) : ";
        with_bind_position([&]() { print_type(typ); });
        out << "(arg) {}}";
        break;
      case ir::IrTypeData::Tag::variant_t: {
        out << "BAPEL_VARIANT(" << id;
        auto tags = typ.tags();
        for (const auto& tag : tags) {
          out << ", (" << to_id(tag.id) << ", (";
          with_bind_position([&]() {
            if (tag.type) print_type(*tag.type);
          });
          out << "))";
        }
        out << ")";
        break;
      }
      default:
        break;
    }
  }

  void print_decl(const ir::IrDecl& decl) {
    if (mode == ir::PrinterMode::ModeSource) return;
    if (mode == ir::PrinterMode::ModePublicHeader && !decl.export_flag) return;
    if (mode == ir::PrinterMode::ModePrivateHeader && decl.export_flag) return;

    if (decl.is(ir::IrDeclData::Tag::trait_decl)) {
      print_trait_decl(decl);
      return;
    }
    if (decl.is(ir::IrDeclData::Tag::name_decl)) return;

    if (decl.is(ir::IrDeclData::Tag::alias_decl)) {
      const auto& ad = decl.data.as_alias_decl();
      print_in_namespace(ad.id, [&](const std::string& id) {
        if (ad.type.is(ir::IrTypeData::Tag::lambda_t)) {
          auto tvars = ad.type.lambda_vars();
          out << "template <";
          ir::Interleave(tvars, [&]() { out << ", "; }, [&](int, const std::string& tv) { out << "typename " << tv; });
          out << "> struct " << id;
        } else {
          out << "struct " << id;
        }
        out << ";\n";
      });
      return;
    }

    if (decl.is(ir::IrDeclData::Tag::term_decl)) {
      const auto& td = decl.data.as_term_decl();
      if (has_function(td.id)) return;
      const auto& typ = td.type;
      switch (typ.tag()) {
        case ir::IrTypeData::Tag::app_t:
        case ir::IrTypeData::Tag::array_t:
        case ir::IrTypeData::Tag::exist_var_t:
        case ir::IrTypeData::Tag::name_t:
        case ir::IrTypeData::Tag::tuple_t:
        case ir::IrTypeData::Tag::var_t:
          print_in_namespace(td.id, [&](const std::string& id) {
            with_bind_position([&]() {
              print_type(typ);
              out << " " << id;
            });
          });
          break;
        case ir::IrTypeData::Tag::forall_t:
          print_in_namespace(td.id, [&](const std::string& id) {
            auto tvars = typ.forall_type_params();
            print_template_params(tvars, false);
            print_decl(ir::new_term_decl(id, typ.forall_body(), decl.export_flag));
          });
          break;
        case ir::IrTypeData::Tag::fun_t:
          print_in_namespace(td.id, [&](const std::string& id) {
            const auto& fn = typ.data.as_fun_t();
            with_bind_position([&]() {
              if (fn.ret) print_type(*fn.ret);
            });
            out << " " << id << "(";
            if (fn.arg) print_type(*fn.arg);
            out << ");";
          });
          break;
        default:
          break;
      }
    }
  }

  void print_type_def(const ir::IrDecl& decl) {
    if (mode == ir::PrinterMode::ModeSource) return;
    if (mode == ir::PrinterMode::ModePublicHeader && !decl.export_flag) return;
    if (mode == ir::PrinterMode::ModePrivateHeader && decl.export_flag) return;

    if (decl.is(ir::IrDeclData::Tag::name_decl)) {
      print_type(ir::new_name_type(decl.data.as_name_decl().id));
      out << ";\n";
    } else if (decl.is(ir::IrDeclData::Tag::alias_decl)) {
      const auto& ad = decl.data.as_alias_decl();
      bool is_anonym = ad.id.rfind("__anonym_", 0) == 0;
      if (is_anonym) {
        out << "#ifndef " << ad.id << "_DEFINED\n";
        out << "#define " << ad.id << "_DEFINED\n";
      }
      print_in_namespace(ad.id, [&](const std::string& id) {
        print_alias_decl(id, ad.type);
        out << ";\n";
      });
      if (is_anonym) {
        out << "#endif\n";
      }
    }
  }

  void print_function_signature(const ir::IrFunction& f) {
    print_in_namespace(f.id, [&](const std::string& id) {
      print_template_params(f.type_params, false);
      with_bind_position([&]() { print_type(f.ret_type); });
      out << " " << id << "(";
      ir::Interleave(f.args, [&]() { out << ", "; }, [&](int, const ir::FunctionArg& arg) {
        print_type(arg.type);
        out << " " << arg.id;
      });
      out << ");\n";
    });
  }

  void print_function_full(const ir::IrFunction& f) {
    print_in_namespace(f.id, [&](const std::string& id) {
      print_template_params(f.type_params, true);
      with_bind_position([&]() { print_type(f.ret_type); });
      out << " " << id << "(";
      ir::Interleave(f.args, [&]() { out << ", "; }, [&](int, const ir::FunctionArg& arg) {
        print_type(arg.type);
        out << " " << arg.id;
      });
      out << ")\n";
      with_last_term(true, [&]() { print_term(f.body); });
      out << "\n";
    });
  }

  void print_trait_decl(const ir::IrDecl& decl) {
    if (!decl.is(ir::IrDeclData::Tag::trait_decl)) return;
    const auto& td = decl.data.as_trait_decl();
    print_in_namespace(trait_cpp_name(td.id), [&](const std::string& id) {
      out << "template <typename Self";
      for (const auto& tp : td.type_params) {
        out << ", typename " << tp.var;
      }
      out << ">\nstruct " << id << ";\n";
    });
  }

  void print_trait_impl(const ir::IrTraitImpl& impl) {
    bool exported = false;
    std::string base_name = base_type_name(impl.type_name);
    ir::IrDecl decl;
    if (!base_name.empty() && find_decl(base_name, decl)) {
      exported = decl.export_flag;
    }

    if (mode == ir::PrinterMode::ModeSource) return;
    if (mode == ir::PrinterMode::ModePublicHeader && !exported) return;
    if (mode == ir::PrinterMode::ModePrivateHeader && exported) return;

    if (impl.case_val == ir::ImplCase::InherentImpl) {
      print_in_namespace(inherent_cpp_name(base_name), [&](const std::string& id) {
        if (!impl.type_params.empty()) {
          out << "template <";
          ir::Interleave(impl.type_params, [&]() { out << ", "; }, [&](int, const ir::TypeParam& tp) {
            out << "typename " << tp.var;
          });
          out << ">\n";
        }
        out << "struct " << id << " {\n";
        out << "  " << id << "() = delete;\n";
        for (const auto& m : impl.methods) {
          print_template_params(m.type_params, false);
          out << "  static inline ";
          with_bind_position([&]() { print_type(m.ret_type); });
          out << " " << m.id << "(";
          ir::Interleave(m.args, [&]() { out << ", "; }, [&](int, const ir::FunctionArg& arg) {
            print_type(arg.type);
            out << " " << arg.id;
          });
          out << ") ";
          with_last_term(true, [&]() { print_term(m.body); });
          out << "\n";
        }
        out << "};\n";
      });
      return;
    }

    std::string trait_name = impl.trait_type.trait_name();
    print_in_namespace(trait_cpp_name(trait_name), [&](const std::string& id) {
      if (!impl.type_params.empty()) {
        out << "template <";
        ir::Interleave(impl.type_params, [&]() { out << ", "; }, [&](int, const ir::TypeParam& tp) {
          out << "typename " << tp.var;
        });
        out << ">\n";
      } else {
        out << "template <>\n";
      }
      out << "struct " << id << "<";
      with_bind_position([&]() { print_type(impl.type_name); });
      for (const auto& arg : impl.trait_type.app_args()) {
        out << ", ";
        with_bind_position([&]() { print_type(arg); });
      }
      out << "> {\n";
      for (const auto& m : impl.methods) {
        print_template_params(m.type_params, false);
        out << "  static inline ";
        with_bind_position([&]() { print_type(m.ret_type); });
        out << " " << m.id << "(";
        ir::Interleave(m.args, [&]() { out << ", "; }, [&](int, const ir::FunctionArg& arg) {
          print_type(arg.type);
          out << " " << arg.id;
        });
        out << ") ";
        with_last_term(true, [&]() { print_term(m.body); });
        out << "\n";
      }
      out << "};\n";
    });
  }

  void print_module_top() {
    switch (mode) {
      case ir::PrinterMode::ModePublicHeader:
        out << "#pragma once\n\n";
        out << "#include <array>\n";
        out << "#include <cstdlib>\n";
        out << "#include <cmath>\n";
        out << "#include <functional>\n";
        out << "#include <optional>\n";
        out << "#include <string>\n";
        out << "#include <tuple>\n";
        out << "#include <variant>\n";
        out << "#include <vector>\n\n";
        out << "#include \"bapel/core_variant.h\"\n\n";
        break;
      case ir::PrinterMode::ModePrivateHeader: {
        out << "#pragma once\n\n";
        if (mod) {
          std::string h_path = to_header_path(mod->module_id);
          out << "#include \"" << h_path << "\"\n\n";
        }
        break;
      }
      case ir::PrinterMode::ModeSource: {
        out << "\n";
        if (mod) {
          std::string h_path = to_header_path(mod->module_id);
          auto pos = h_path.rfind(".h");
          std::string priv_path = (pos != std::string::npos ? h_path.substr(0, pos) : h_path) + "_private.h";
          out << "#include \"" << priv_path << "\"\n\n";
        }
        break;
      }
    }
  }

  void print_imports() {
    if (!mod) return;
    if (mode == ir::PrinterMode::ModePublicHeader || mode == ir::PrinterMode::ModeSource) {
      out << "\n";
      for (const auto& imp : mod->imports) {
        out << "#include \"" << to_header_path(imp.module_id) << "\"\n";
      }
      out << "\n";
    }
  }

  void print_impls() {
    if (mode != ir::PrinterMode::ModePublicHeader || !mod) return;
    out << "\n";
    for (const auto& impl : mod->impls) {
      const std::string& fn = impl.relative_filename.value;
      if (fn.size() >= 2 && fn.substr(fn.size() - 2) == ".h") {
        out << "#include \"" << fn << "\"\n";
      }
    }
    out << "\n";
  }

  static bool is_mu_alias(const ir::IrDecl& d) {
    if (!d.is(ir::IrDeclData::Tag::alias_decl)) return false;
    const ir::IrType* cur = &d.data.as_alias_decl().type;
    while (cur->is(ir::IrTypeData::Tag::lambda_t) && cur->data.as_lambda_t().type) {
      cur = cur->data.as_lambda_t().type.get();
    }
    return cur->is(ir::IrTypeData::Tag::mu_t);
  }

  void print_decls(const std::vector<ir::IrDecl>& decls) {
    for (const auto& d : decls) {
      if (is_mu_alias(d)) {
        print_decl(d);
      }
    }
    for (const auto& d : decls) {
      if (d.is(ir::IrDeclData::Tag::alias_decl)) {
        print_type_def(d);
      } else {
        print_decl(d);
      }
    }
  }

  void print_module_content() {
    if (!mod) return;
    print_module_top();
    print_impls();
    print_imports();

    print_decls(mod->decls);

    if (mode == ir::PrinterMode::ModePublicHeader || mode == ir::PrinterMode::ModePrivateHeader) {
      for (const auto& f : mod->functions) {
        bool is_pub = f.export_flag;
        if (mode == ir::PrinterMode::ModePublicHeader && !is_pub) continue;
        if (mode == ir::PrinterMode::ModePrivateHeader && is_pub) continue;
        print_function_signature(f);
      }
      for (const auto& f : mod->functions) {
        bool is_template = !f.type_params.empty();
        bool is_pub = f.export_flag;
        if (mode == ir::PrinterMode::ModePublicHeader && is_pub && is_template) {
          print_function_full(f);
        } else if (mode == ir::PrinterMode::ModePrivateHeader && !is_pub && is_template) {
          print_function_full(f);
        }
      }
    } else if (mode == ir::PrinterMode::ModeSource) {
      for (const auto& f : mod->functions) {
        bool is_template = !f.type_params.empty();
        if (!is_template) {
          print_function_full(f);
        }
      }
      return;
    }

    for (const auto& impl : mod->trait_impls) {
      print_trait_impl(impl);
    }
  }
};

// Anonymous struct extraction & topological sorting (Phase 8.3)

inline uint64_t fnv1a_hash(std::string_view s) {
  uint64_t h = 14695981039346656037ULL;
  for (unsigned char c : s) {
    h ^= c;
    h *= 1099511628211ULL;
  }
  return h;
}

inline ir::IrType gen_name_type(ir::IrType typ, bool export_flag, std::map<std::string, AnonymousType>& anonymous_types) {
  std::vector<std::string> bound_vars;
  std::vector<std::string> raw_free_vars;
  ir::collect_free_vars(typ, bound_vars, raw_free_vars);
  std::set<std::string> free_vars(raw_free_vars.begin(), raw_free_vars.end());

  ir::IrType decl_type = typ;
  ir::IrKind decl_kind = ir::new_type_kind();
  for (auto it = free_vars.rbegin(); it != free_vars.rend(); ++it) {
    decl_type = ir::new_lambda_type(*it, ir::new_type_kind(), std::move(decl_type));
    decl_kind = ir::new_arrow_kind(ir::new_type_kind(), std::move(decl_kind));
  }

  std::string hash = std::to_string(fnv1a_hash(decl_type.to_string()));
  std::string name = "__anonym_" + hash;
  auto it = anonymous_types.find(hash);
  if (it != anonymous_types.end()) {
    if (export_flag && !it->second.decl.export_flag) {
      it->second.decl.export_flag = true;
    }
  } else {
    anonymous_types[hash] = AnonymousType{
        ir::new_name_type(name),
        ir::new_alias_decl(name, std::move(decl_kind), std::move(decl_type), export_flag)
    };
  }

  ir::IrType ref = ir::new_name_type(name);
  for (const auto& fv : free_vars) {
    ref = ir::new_app_type(std::move(ref), ir::new_var_type(fv));
  }
  return ref;
}

inline ir::IrType record_anonymous_types(
    ir::IrType typ,
    bool export_flag,
    std::map<std::string, AnonymousType>& anonymous_types,
    const std::map<std::string, ir::IrType>& mu_subst = {}) {
  ir::IrType mapped = ir::map_type(typ, [&](const ir::IrType& child) {
    return record_anonymous_types(child, export_flag, anonymous_types, mu_subst);
  });
  if (mapped.is(ir::IrTypeData::Tag::struct_t)) {
    for (const auto& [var, target] : mu_subst) {
      mapped = ir::substitute_type(mapped, ir::new_var_type(var), target);
    }
    return gen_name_type(std::move(mapped), export_flag, anonymous_types);
  }
  return mapped;
}

inline void record_anonymous_types_in_alias(ir::IrDecl& d, std::map<std::string, AnonymousType>& anonymous_types) {
  if (!d.is(ir::IrDeclData::Tag::alias_decl)) return;
  auto& ad = d.data.as_alias_decl();

  std::vector<std::pair<std::string, ir::IrKind>> lambdas;
  ir::IrType cur = ad.type;
  while (cur.is(ir::IrTypeData::Tag::lambda_t) && cur.data.as_lambda_t().type) {
    lambdas.push_back({cur.data.as_lambda_t().var, cur.data.as_lambda_t().kind});
    cur = *cur.data.as_lambda_t().type;
  }

  ir::IrType self_type = ir::new_name_type(ad.id);
  for (const auto& [var, kind] : lambdas) {
    self_type = ir::new_app_type(std::move(self_type), ir::new_var_type(var));
  }

  std::map<std::string, ir::IrType> mu_subst;
  bool has_mu = false;
  std::string mu_var;
  ir::IrKind mu_kind;
  if (cur.is(ir::IrTypeData::Tag::mu_t) && cur.data.as_mu_t().body) {
    has_mu = true;
    mu_var = cur.data.as_mu_t().var;
    mu_kind = cur.data.as_mu_t().kind;
    mu_subst[mu_var] = self_type;
    cur = *cur.data.as_mu_t().body;
  }

  if (cur.is(ir::IrTypeData::Tag::struct_t)) {
    std::vector<ir::StructField> fields;
    for (const auto& f : cur.fields()) {
      ir::StructField nf = f;
      if (nf.type) {
        nf.type = std::make_shared<ir::IrType>(record_anonymous_types(*nf.type, d.export_flag, anonymous_types, mu_subst));
      }
      fields.push_back(std::move(nf));
    }
    cur = ir::new_struct_type(std::move(fields));
  } else {
    cur = record_anonymous_types(cur, d.export_flag, anonymous_types, mu_subst);
  }

  if (has_mu) {
    cur = ir::new_mu_type(mu_var, mu_kind, std::move(cur));
  }
  for (auto it = lambdas.rbegin(); it != lambdas.rend(); ++it) {
    cur = ir::new_lambda_type(it->first, it->second, std::move(cur));
  }
  ad.type = std::move(cur);
}

inline void get_free_vars_from_type(const ir::IrType& typ, std::set<std::string>& bound_names, std::vector<ir::IrType>& free_vars) {
  switch (typ.tag()) {
    case ir::IrTypeData::Tag::name_t: {
      const std::string& n = typ.data.as_name_t();
      if (bound_names.find(n) == bound_names.end()) {
        bound_names.insert(n);
        free_vars.push_back(typ);
      }
      return;
    }
    case ir::IrTypeData::Tag::app_t: {
      const auto& a = typ.data.as_app_t();
      if (a.fun) get_free_vars_from_type(*a.fun, bound_names, free_vars);
      if (a.arg) get_free_vars_from_type(*a.arg, bound_names, free_vars);
      return;
    }
    case ir::IrTypeData::Tag::array_t: {
      const auto& arr = typ.data.as_array_t();
      if (arr.elem_type) get_free_vars_from_type(*arr.elem_type, bound_names, free_vars);
      return;
    }
    case ir::IrTypeData::Tag::fun_t: {
      const auto& fn = typ.data.as_fun_t();
      if (fn.arg) get_free_vars_from_type(*fn.arg, bound_names, free_vars);
      if (fn.ret) get_free_vars_from_type(*fn.ret, bound_names, free_vars);
      return;
    }
    case ir::IrTypeData::Tag::struct_t:
      for (const auto& f : typ.fields()) {
        if (f.type) get_free_vars_from_type(*f.type, bound_names, free_vars);
      }
      return;
    case ir::IrTypeData::Tag::tuple_t:
      for (const auto& e : typ.elems()) {
        get_free_vars_from_type(e, bound_names, free_vars);
      }
      return;
    case ir::IrTypeData::Tag::variant_t:
      for (const auto& tag : typ.tags()) {
        if (tag.type) get_free_vars_from_type(*tag.type, bound_names, free_vars);
      }
      return;
    case ir::IrTypeData::Tag::lambda_t:
      if (typ.data.as_lambda_t().type) {
        get_free_vars_from_type(*typ.data.as_lambda_t().type, bound_names, free_vars);
      }
      return;
    case ir::IrTypeData::Tag::mu_t:
      if (typ.data.as_mu_t().body) {
        get_free_vars_from_type(*typ.data.as_mu_t().body, bound_names, free_vars);
      }
      return;
    case ir::IrTypeData::Tag::exist_var_t:
    case ir::IrTypeData::Tag::forall_t:
    case ir::IrTypeData::Tag::var_t:
      return;
  }
}

inline bool is_mu_alias_type(const ir::IrType& t) {
  const ir::IrType* cur = &t;
  while (cur->is(ir::IrTypeData::Tag::lambda_t) && cur->data.as_lambda_t().type) {
    cur = cur->data.as_lambda_t().type.get();
  }
  return cur->is(ir::IrTypeData::Tag::mu_t);
}

inline std::vector<ir::IrDecl> topo_sort_alias_decls(const std::vector<ir::IrDecl>& decls) {
  std::map<std::string, int> nodes_by_name;
  for (size_t i = 0; i < decls.size(); ++i) {
    nodes_by_name[decls[i].is(ir::IrDeclData::Tag::alias_decl) ? decls[i].data.as_alias_decl().id : ""] = static_cast<int>(i);
  }

  std::vector<std::vector<int>> adj(decls.size());
  std::vector<int> in_degree(decls.size(), 0);

  std::vector<std::set<std::string>> decl_refs(decls.size());
  for (size_t i = 0; i < decls.size(); ++i) {
    if (!decls[i].is(ir::IrDeclData::Tag::alias_decl)) continue;
    std::set<std::string> bound;
    std::vector<ir::IrType> free_vars;
    get_free_vars_from_type(decls[i].data.as_alias_decl().type, bound, free_vars);
    for (const auto& fvar : free_vars) {
      if (fvar.is(ir::IrTypeData::Tag::name_t)) {
        decl_refs[i].insert(fvar.data.as_name_t());
      }
    }
  }

  for (size_t i = 0; i < decls.size(); ++i) {
    if (!decls[i].is(ir::IrDeclData::Tag::alias_decl)) continue;
    const auto& ad = decls[i].data.as_alias_decl();
    bool is_anonym = ad.id.rfind("__anonym_", 0) == 0;
    for (const auto& ref_name : decl_refs[i]) {
      auto it = nodes_by_name.find(ref_name);
      if (it != nodes_by_name.end() && it->second != static_cast<int>(i)) {
        if (is_anonym && decls[it->second].is(ir::IrDeclData::Tag::alias_decl) &&
            is_mu_alias_type(decls[it->second].data.as_alias_decl().type) &&
            decl_refs[it->second].count(ad.id) > 0) {
          continue;
        }
        adj[it->second].push_back(static_cast<int>(i));
        in_degree[i]++;
      }
    }
  }

  std::vector<int> zero_queue;
  for (size_t i = 0; i < decls.size(); ++i) {
    if (in_degree[i] == 0) zero_queue.push_back(static_cast<int>(i));
  }

  std::vector<ir::IrDecl> sorted;
  while (!zero_queue.empty()) {
    int u = zero_queue.back();
    zero_queue.pop_back();
    sorted.push_back(decls[u]);
    for (int v : adj[u]) {
      if (--in_degree[v] == 0) {
        zero_queue.push_back(v);
      }
    }
  }

  if (sorted.size() != decls.size()) {
    return decls; // fallback if cycle
  }
  return sorted;
}

inline void topo_sort_decls(std::vector<ir::IrDecl>& decls) {
  std::vector<ir::IrDecl> term_decls, alias_decls, name_decls, trait_decls;
  for (const auto& d : decls) {
    switch (d.tag()) {
      case ir::IrDeclData::Tag::term_decl: term_decls.push_back(d); break;
      case ir::IrDeclData::Tag::alias_decl: alias_decls.push_back(d); break;
      case ir::IrDeclData::Tag::name_decl: name_decls.push_back(d); break;
      case ir::IrDeclData::Tag::trait_decl: trait_decls.push_back(d); break;
      case ir::IrDeclData::Tag::impl_decl: break;
    }
  }

  alias_decls = topo_sort_alias_decls(alias_decls);

  decls.clear();
  decls.insert(decls.end(), name_decls.begin(), name_decls.end());
  decls.insert(decls.end(), alias_decls.begin(), alias_decls.end());
  decls.insert(decls.end(), trait_decls.begin(), trait_decls.end());
  decls.insert(decls.end(), term_decls.begin(), term_decls.end());
}

inline void record_anonymous_types_from_module(ir::IrTypedModule* mod) {
  if (!mod) return;
  std::map<std::string, AnonymousType> anonymous_types;

  auto rewrite_decl = [&](ir::IrDecl& d) {
    if (d.is(ir::IrDeclData::Tag::term_decl)) {
      auto& td = d.data.as_term_decl();
      td.type = record_anonymous_types(td.type, d.export_flag, anonymous_types);
    } else if (d.is(ir::IrDeclData::Tag::alias_decl)) {
      record_anonymous_types_in_alias(d, anonymous_types);
    } else if (d.is(ir::IrDeclData::Tag::trait_decl)) {
      for (auto& m : d.data.as_trait_decl().methods) {
        for (auto& arg : m.args) {
          arg.type = record_anonymous_types(arg.type, d.export_flag, anonymous_types);
        }
        m.ret_type = record_anonymous_types(m.ret_type, d.export_flag, anonymous_types);
      }
    }
  };

  for (auto& d : mod->decls) rewrite_decl(d);
  for (auto& f : mod->functions) {
    for (auto& arg : f.args) {
      arg.type = record_anonymous_types(arg.type, f.export_flag, anonymous_types);
    }
    f.ret_type = record_anonymous_types(f.ret_type, f.export_flag, anonymous_types);
  }
  for (auto& ti : mod->trait_impls) {
    for (auto& m : ti.methods) {
      for (auto& arg : m.args) {
        arg.type = record_anonymous_types(arg.type, m.export_flag, anonymous_types);
      }
      m.ret_type = record_anonymous_types(m.ret_type, m.export_flag, anonymous_types);
    }
  }

  for (const auto& [k, v] : anonymous_types) {
    mod->decls.push_back(v.decl);
  }
  topo_sort_decls(mod->decls);
}

inline void format_file(const std::string& filename) {
  std::string cmd = "clang-format -i " + filename + " 2>/dev/null";
  std::system(cmd.c_str());
}

inline bool print_module_to_file(const ir::IrTypedModule& mod, ir::PrinterMode mode, const std::string& filename) {
  std::ofstream f(filename);
  if (!f.is_open()) return false;
  CppPrinter printer(f, mode, &mod);
  printer.print_module_content();
  f.close();
  format_file(filename);
  return true;
}

inline bool compile_ir_module_direct(ir::IrTypedModule mod, const std::string& output_base) {
  record_anonymous_types_from_module(&mod);

  std::string base = output_base;
  auto dot = base.rfind('.');
  if (dot != std::string::npos) {
    base = base.substr(0, dot);
  }

  if (!print_module_to_file(mod, ir::PrinterMode::ModePublicHeader, base + ".h")) return false;
  if (!print_module_to_file(mod, ir::PrinterMode::ModePrivateHeader, base + "_private.h")) return false;
  if (!print_module_to_file(mod, ir::PrinterMode::ModeSource, base + ".cc")) return false;
  return true;
}

// @bpl: pub codegen::compile_module: (&comp::Session, String, String) -> i64
inline int64_t compile_module(comp::Session* session, const std::string& base_file, const std::string& output_base) {
  comp::TypecheckOptions options;

  ir::IrTypedModule mod;
  std::string err;
  if (!comp::typecheck_module_from_file(session->querier, options, base_file, mod, err)) {
    std::cerr << "Typecheck error in " << base_file << ":\n" << err << "\n";
    return 1;
  }

  return compile_ir_module_direct(std::move(mod), output_base) ? 0 : 1;
}

} // namespace codegen

namespace typechecker {

// @bpl: pub typechecker::run: (&comp::Session, Vector String) -> (i64, String)
inline std::tuple<int64_t, std::string> run(comp::Session* session, const std::vector<std::string>& args) {
  std::string format = "flat";
  std::string input_file;

  for (size_t i = 0; i < args.size(); ++i) {
    const std::string& arg = args[i];
    if (arg.rfind("-format=", 0) == 0) {
      format = arg.substr(8);
    } else if (arg.rfind("--format=", 0) == 0) {
      format = arg.substr(9);
    } else if (arg == "-format" || arg == "--format") {
      if (i + 1 < args.size()) {
        format = args[++i];
      }
    } else if (arg.rfind("-", 0) != 0) {
      if (input_file.empty()) {
        input_file = arg;
      } else {
        return {1, "Usage: typechecker [-format=flat] <input_file>\n"};
      }
    }
  }

  if (input_file.empty()) {
    return {1, "Usage: typechecker [-format=flat] <input_file>\n"};
  }

  comp::TypecheckOptions options;

  ir::IrTypedModule mod;
  std::string err;
  if (!comp::typecheck_module_from_file(session->querier, options, input_file, mod, err)) {
    return {1, "Failed to typecheck \"" + input_file + "\": " + err + "\n"};
  }

  std::stringstream out;
  if (format == "flat") {
    out << "MODULE " << mod.module_id.to_string() << "\n";
    out << "CASE base\n";

    for (const auto& imp : mod.imports) {
      out << "IMPORT " << imp.module_id.to_string() << "\n";
    }
    for (const auto& impl : mod.impls) {
      out << "IMPL " << impl.relative_filename.value << "\n";
    }
    for (const auto& decl : mod.decls) {
      std::string s = decl.to_string();
      std::string escaped_s = s;
      comp::replace_all_str(escaped_s, "\n", "\\n");
      out << "DECL " << escaped_s << "\n";
      std::string export_str = decl.export_flag ? "1" : "0";
      out << "DECL_DEF " << export_str << " " << decl.id() << " " << escaped_s << "\n";
    }
    for (const auto& trait_impl : mod.trait_impls) {
      std::string s = trait_impl.to_string();
      std::string escaped_s = s;
      comp::replace_all_str(escaped_s, "\n", "\\n");
      out << "TRAIT_IMPL " << escaped_s << "\n";
      std::string trait_type_str = trait_impl.trait_type.to_string();
      comp::replace_all_str(trait_type_str, "\n", "\\n");
      std::string type_name_str = trait_impl.type_name.to_string();
      comp::replace_all_str(type_name_str, "\n", "\\n");
      out << "TRAIT_DEF " << trait_type_str << " " << type_name_str << " " << escaped_s << "\n";
    }
    for (const auto& fn : mod.functions) {
      std::string s = fn.to_string();
      std::string escaped_s = s;
      comp::replace_all_str(escaped_s, "\n", "\\n");
      out << "FUNC " << escaped_s << "\n";
      std::string export_str = fn.export_flag ? "1" : "0";
      out << "FUNC_DEF " << export_str << " " << fn.id << " " << escaped_s << "\n";
    }
  }
  return {0, out.str()};
}

} // namespace typechecker
