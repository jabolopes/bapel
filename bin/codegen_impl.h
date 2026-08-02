#pragma once

#include "ir_parser.h"
#include "comp/typecheck_unit.h"
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
  if (typ.is(ir::IrTypeCase::NameType)) {
    return typ.name;
  }
  if (typ.is(ir::IrTypeCase::AppType)) {
    return base_type_name(typ.app_fun());
  }
  return "";
}

inline int count_type_vars(const ir::IrKind& kind) {
  if (kind.is_arrow_kind() && kind.right) {
    return 1 + count_type_vars(*kind.right);
  }
  return 0;
}

inline int count_alias_type_vars(const ir::IrType& typ) {
  if (typ.is(ir::IrTypeCase::LambdaType)) {
    return static_cast<int>(typ.lambda_vars().size());
  }
  return 0;
}

inline bool is_cpp_statement(const ir::IrTerm& term) {
  switch (term.case_val) {
    case ir::IrTermCase::AssignTerm:
    case ir::IrTermCase::BlockTerm:
    case ir::IrTermCase::LetTerm:
    case ir::IrTermCase::MatchTerm:
    case ir::IrTermCase::ReturnTerm:
      return true;
    case ir::IrTermCase::AppTermTerm: {
      auto [id, types, arg] = term.app_args();
      if (id.is(ir::IrTermCase::VarTerm) && id.var_data) {
        const std::string& v = id.var_data->id;
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
  std::map<std::string, AnonymousType> anonymous_types;
  const ir::IrUnit* unit = nullptr;

  CppPrinter(std::ostream& output, ir::PrinterMode m, const ir::IrUnit* u)
      : out(output), mode(m), unit(u) {}

  std::string gen_id() {
    return "__v_" + std::to_string(idgen++);
  }

  bool is_type_decl(const ir::IrDecl& decl) const {
    return decl.is(ir::IrDeclCase::NameDecl) || decl.is(ir::IrDeclCase::AliasDecl);
  }

  bool find_decl(const std::string& id, ir::IrDecl& out_decl) const {
    if (!unit) return false;
    for (const auto& d : unit->decls) {
      if (d.id() == id) { out_decl = d; return true; }
    }
    for (const auto& d : unit->import_decls) {
      if (d.id() == id) { out_decl = d; return true; }
    }
    for (const auto& d : unit->impl_decls) {
      if (d.id() == id) { out_decl = d; return true; }
    }
    return false;
  }

  bool find_decl_for_type(const ir::IrType& typ, ir::IrDecl& out_decl) const {
    if (!unit) return false;
    std::string target_hash = hash_type(typ);
    auto check_list = [&](const std::vector<ir::IrDecl>& list) -> bool {
      for (const auto& d : list) {
        if (d.is(ir::IrDeclCase::AliasDecl) && d.alias && hash_type(d.alias->type) == target_hash) {
          out_decl = d;
          return true;
        }
      }
      return false;
    };
    return check_list(unit->decls) || check_list(unit->import_decls) || check_list(unit->impl_decls);
  }

  bool find_trait_decl(const std::string& id, ir::IrDecl& out_decl) const {
    if (!unit) return false;
    auto check_list = [&](const std::vector<ir::IrDecl>& list) -> bool {
      for (const auto& d : list) {
        if (d.id() == id && d.is(ir::IrDeclCase::TraitDecl)) {
          out_decl = d;
          return true;
        }
      }
      return false;
    };
    return check_list(unit->decls) || check_list(unit->import_decls) || check_list(unit->impl_decls);
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

    switch (typ.case_val) {
      case ir::IrTypeCase::AppType: {
        print_type(typ.app_fun());
        auto args = typ.app_args();
        out << "<";
        ir::Interleave(args, [&]() { out << ", "; }, [&](int, const ir::IrType& a) { print_type(a); });
        out << ">";
        break;
      }
      case ir::IrTypeCase::ArrayType: {
        out << "std::array<";
        if (typ.array && typ.array->elem_type) {
          print_type(*typ.array->elem_type);
        }
        out << ", " << (typ.array ? typ.array->size : 0) << ">";
        break;
      }
      case ir::IrTypeCase::ExistVarType:
        out << typ.to_string();
        break;
      case ir::IrTypeCase::ForallType: {
        auto tvars = typ.forall_vars();
        out << "template <";
        ir::Interleave(tvars, [&]() { out << ", "; }, [&](int, const std::string& tv) { out << "typename " << tv; });
        out << "> ";
        print_type(typ.forall_body());
        break;
      }
      case ir::IrTypeCase::FunType: {
        out << "std::function<";
        with_bind_position([&]() {
          if (typ.fun && typ.fun->ret) print_type(*typ.fun->ret);
        });
        out << "(";
        if (typ.fun && typ.fun->arg) print_type(*typ.fun->arg);
        out << ")>";
        break;
      }
      case ir::IrTypeCase::NameType: {
        const std::string& n = typ.name;
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
      case ir::IrTypeCase::StructType: {
        std::string h = hash_type(typ);
        auto it = anonymous_types.find(h);
        if (it != anonymous_types.end()) {
          print_type(it->second.name_type);
          return;
        }
        ir::IrDecl decl;
        if (find_decl_for_type(typ, decl)) {
          out << "::" << to_id(decl.id());
          return;
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
      case ir::IrTypeCase::TupleType: {
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
      case ir::IrTypeCase::VariantType: {
        out << "std::variant<";
        auto tags = typ.tags();
        ir::Interleave(tags, [&]() { out << ", "; }, [&](int, const ir::VariantTag& tag) {
          with_bind_position([&]() {
            if (tag.type) print_type(*tag.type);
          });
          out << "/* " << to_id(tag.id) << " */";
        });
        out << ">";
        break;
      }
      case ir::IrTypeCase::VarType:
        out << typ.var;
        break;
      default:
        break;
    }
  }

  void print_app_type_term(const ir::IrTerm& term) {
    auto app_res = term.app_types();
    const auto& arg = app_res.first;
    const auto& types = app_res.second;
    if (arg.is(ir::IrTermCase::ConstTerm)) {
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

    if (id.is(ir::IrTermCase::VarTerm) && id.var_data) {
      const std::string& vname = id.var_data->id;
      if (vname == "ifthen" && arg.is(ir::IrTermCase::TupleTerm) && arg.tuple_data->elems.size() >= 2) {
        const auto& cond = arg.tuple_data->elems[0];
        const auto& then_term = arg.tuple_data->elems[1];
        out << "if (";
        with_last_term(false, [&]() { print_term(cond); });
        out << ") ";
        print_term(then_term);
        return;
      }
      if (vname == "ifelse" && arg.is(ir::IrTermCase::TupleTerm) && arg.tuple_data->elems.size() >= 3) {
        const auto& cond = arg.tuple_data->elems[0];
        const auto& then_term = arg.tuple_data->elems[1];
        const auto& else_term = arg.tuple_data->elems[2];
        out << "if (";
        with_last_term(false, [&]() { print_term(cond); });
        out << ") ";
        print_term(then_term);
        out << " else ";
        print_term(else_term);
        return;
      }
      if (vname == "core::for" && arg.is(ir::IrTermCase::TupleTerm) && arg.tuple_data->elems.size() >= 2) {
        const auto& cond = arg.tuple_data->elems[0];
        auto fn_res = arg.tuple_data->elems[1].to_function();
        const auto& body = std::get<2>(fn_res);
        out << "while (";
        with_last_term(false, [&]() { print_term(cond); });
        out << ") ";
        with_last_term(false, [&]() { print_term(body); });
        return;
      }
      if (ir::is_operator(vname)) {
        if (arg.is(ir::IrTermCase::TupleTerm) && !arg.tuple_data->elems.empty()) {
          out << "(";
          print_term(arg.tuple_data->elems[0]);
          out << ")";
          print_term(id);
          out << "(";
          std::vector<ir::IrTerm> rest(arg.tuple_data->elems.begin() + 1, arg.tuple_data->elems.end());
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
          if (arg.is(ir::IrTermCase::TupleTerm)) {
            ir::Interleave(arg.tuple_data->elems, [&]() { out << ", "; }, [&](int, const ir::IrTerm& t) { print_term(t); });
          } else {
            print_term(arg);
          }
          out << ")";
          return;
        }

        // Inherent method resolution
        if (find_decl(trait_name, decl) && (decl.is(ir::IrDeclCase::NameDecl) || decl.is(ir::IrDeclCase::AliasDecl))) {
          int arity = 0;
          if (decl.is(ir::IrDeclCase::NameDecl) && decl.name) {
            arity = count_type_vars(decl.name->kind);
          } else if (decl.alias) {
            arity = count_alias_type_vars(decl.alias->type);
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
            if (arg.is(ir::IrTermCase::TupleTerm)) {
              ir::Interleave(arg.tuple_data->elems, [&]() { out << ", "; }, [&](int, const ir::IrTerm& t) { print_term(t); });
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
    if (id.is(ir::IrTermCase::VarTerm) && id.var_data && !ir::is_operator(id.var_data->id) && !types.empty()) {
      out << "<";
      with_bind_position([&]() {
        ir::Interleave(types, [&]() { out << ", "; }, [&](int, const ir::IrType& t) { print_type(t); });
      });
      out << ">";
    }

    out << "(";
    if (arg.is(ir::IrTermCase::TupleTerm)) {
      ir::Interleave(arg.tuple_data->elems, [&]() { out << ", "; }, [&](int, const ir::IrTerm& t) { print_term(t); });
    } else {
      print_term(arg);
    }
    out << ")";
  }

  void print_let_term(const ir::IrTerm& term) {
    if (!term.let_data) return;
    const auto& c = *term.let_data;
    bool auto_flag = c.value && c.value->is(ir::IrTermCase::TypeAbsTerm);

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
    if (!term.match_data) return;
    const auto& c = *term.match_data;
    std::string variant_id = gen_id();

    out << "{\n";
    out << "auto " << variant_id << " = ";
    with_last_term(false, [&]() {
      if (c.term) print_term(*c.term);
    });
    out << ";\n";
    out << "switch (" << variant_id << ".index()) {\n";
    for (const auto& arm : c.arms) {
      int idx = arm.index.value_or(0);
      out << "case " << idx << ": {\n";
      out << "auto &" << arm.arg << " = std::get<" << idx << ">(" << variant_id << ");\n";
      if (arm.body) print_term(*arm.body);
      out << ";\n";
      out << "}\n";
    }
    out << "} }";
  }

  void print_projection_term(const ir::IrTerm& term) {
    if (!term.projection) return;
    const auto& c = *term.projection;
    ir::IrType obj_type = term.type.value_or(ir::IrType{});
    if (c.reduced_type) obj_type = *c.reduced_type;

    if (obj_type.is(ir::IrTypeCase::StructType)) {
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
    if (obj_type.is(ir::IrTypeCase::TupleType) && obj_type.elem_by_label(c.label, index, elem)) {
      out << "std::get<" << index << ">(";
      if (c.term) print_term(*c.term);
      out << ")";
      return;
    }

    if (obj_type.is(ir::IrTypeCase::VariantType) && obj_type.tag_by_label(c.label, index)) {
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
    if (!term.set_data) return;
    const auto& c = *term.set_data;
    ir::IrType obj_type = term.type.value_or(ir::IrType{});
    if (c.reduced_type) obj_type = *c.reduced_type;

    if (obj_type.is(ir::IrTypeCase::StructType)) {
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

    if (obj_type.is(ir::IrTypeCase::TupleType)) {
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

    bool argless = args.size() == 1 && args[0].type.is(ir::IrTypeCase::TupleType) && args[0].type.elems().empty();
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
    switch (term.case_val) {
      case ir::IrTermCase::AppTypeTerm:
        handle_last_term([&]() { print_app_type_term(term); });
        break;
      case ir::IrTermCase::AppTermTerm:
        if (is_cpp_statement(term)) {
          print_app_term_term(term);
        } else {
          handle_last_term([&]() { print_app_term_term(term); });
        }
        break;
      case ir::IrTermCase::AssignTerm:
        handle_last_term([&]() {
          if (term.assign) {
            with_bind_position([&]() { if (term.assign->ret) print_term(*term.assign->ret); });
            out << " = ";
            if (term.assign->arg) print_term(*term.assign->arg);
          }
        });
        break;
      case ir::IrTermCase::BlockTerm:
        if (term.block) {
          out << "{\n";
          for (size_t i = 0; i < term.block->terms.size(); ++i) {
            bool is_last = (i + 1 == term.block->terms.size());
            if (is_last) {
              print_term(term.block->terms[i]);
              out << ";}\n";
            } else {
              with_last_term(false, [&]() {
                print_term(term.block->terms[i]);
                out << ";\n";
              });
            }
          }
        }
        break;
      case ir::IrTermCase::ConstTerm:
        handle_last_term([&]() {
          if (term.const_data) {
            const auto& lit = term.const_data->literal;
            if (lit.is_int()) out << *lit.int_val;
            else if (lit.is_float()) out << lit.float_val->integer << "." << lit.float_val->decimal;
            else if (lit.is_rune()) out << "'" << *lit.rune_val << "'";
            else if (lit.is_str()) out << "\"" << *lit.str_val << "\"";
          }
        });
        break;
      case ir::IrTermCase::InjectionTerm:
        handle_last_term([&]() {
          if (term.injection) {
            print_type(term.injection->variant_type);
            out << "(";
            if (term.injection->tag_index.has_value()) {
              out << "std::in_place_index<" << *term.injection->tag_index << ">, ";
            }
            if (term.injection->value) print_term(*term.injection->value);
            out << ")";
          }
        });
        break;
      case ir::IrTermCase::LambdaTerm:
      case ir::IrTermCase::TypeAbsTerm:
        handle_last_term([&]() {
          with_last_term(false, [&]() { print_lambda_term(term); });
        });
        break;
      case ir::IrTermCase::LetTerm:
        with_last_term(false, [&]() { print_let_term(term); });
        break;
      case ir::IrTermCase::MatchTerm:
        print_match_term(term);
        break;
      case ir::IrTermCase::ProjectionTerm:
        handle_last_term([&]() { print_projection_term(term); });
        break;
      case ir::IrTermCase::ReturnTerm:
        with_last_term(false, [&]() {
          if (term.return_data && term.return_data->expr) {
            out << "return ";
            print_term(*term.return_data->expr);
            out << ";";
          }
        });
        break;
      case ir::IrTermCase::TupleTerm:
        handle_last_term([&]() {
          if (term.tuple_data) {
            if (position == ir::Position::BindPosition) {
              out << "std::tie(";
            } else if (term.tuple_data->elems.empty()) {
              out << "std::monostate(";
            } else {
              out << "std::make_tuple(";
            }
            ir::Interleave(term.tuple_data->elems, [&]() { out << ", "; }, [&](int, const ir::IrTerm& t) {
              print_term(t);
            });
            out << ")";
          }
        });
        break;
      case ir::IrTermCase::SetTerm:
        handle_last_term([&]() { print_set_term(term); });
        break;
      case ir::IrTermCase::StructTerm:
        handle_last_term([&]() {
          if (term.struct_data) {
            out << "{";
            ir::Interleave(term.struct_data->values, [&]() { out << ", "; }, [&](int, const ir::LabelValue& lv) {
              out << "." << lv.label << " = ";
              if (lv.value) print_term(*lv.value);
            });
            out << "}";
          }
        });
        break;
      case ir::IrTermCase::VarTerm:
        handle_last_term([&]() {
          if (term.var_data) out << to_id(term.var_data->id);
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
    if (bound.is(ir::IrTypeCase::NameType)) {
      trait_name = bound.name;
    } else if (bound.is(ir::IrTypeCase::AppType)) {
      trait_name = bound.app_fun().trait_name();
      args = bound.app_args();
    }

    std::string cpp_name = "::" + trait_cpp_name(trait_name);
    std::ostringstream ss;
    ss << "(sizeof(" << cpp_name << "<" << vk_var;
    for (const auto& arg : args) {
      ss << ", ";
      std::ostringstream temp_out;
      CppPrinter temp_printer(temp_out, mode, unit);
      temp_printer.print_type(arg);
      ss << temp_out.str();
    }
    ss << ">) > 0)";
    return ss.str();
  }

  void print_alias_decl(const std::string& id, const ir::IrType& typ) {
    switch (typ.case_val) {
      case ir::IrTypeCase::LambdaType: {
        auto tvars = typ.lambda_vars();
        out << "template <";
        ir::Interleave(tvars, [&]() { out << ", "; }, [&](int, const std::string& tv) { out << "typename " << tv; });
        out << "> ";
        print_alias_decl(id, typ.lambda_body());
        break;
      }
      case ir::IrTypeCase::NameType:
        out << "using " << id << " = ";
        print_type(typ);
        break;
      case ir::IrTypeCase::StructType:
        out << "struct " << id << " {\n";
        for (const auto& f : typ.fields()) {
          if (f.type) print_type(*f.type);
          out << " " << f.id << ";\n";
        }
        out << "}";
        break;
      case ir::IrTypeCase::TupleType:
        out << "struct " << id << " : ";
        with_bind_position([&]() { print_type(typ); });
        out << " { " << id << "(const ";
        with_bind_position([&]() { print_type(typ); });
        out << "& arg) : ";
        with_bind_position([&]() { print_type(typ); });
        out << "(arg) {}}";
        break;
      case ir::IrTypeCase::VariantType:
        out << "struct " << id << " : ";
        print_type(typ);
        out << " { using ";
        print_type(typ);
        out << "::variant; }";
        break;
      default:
        break;
    }
  }

  void print_decl(const ir::IrDecl& decl) {
    if (mode == ir::PrinterMode::ModeSource) return;
    if (mode == ir::PrinterMode::ModePublicHeader && !decl.export_flag) return;
    if (mode == ir::PrinterMode::ModePrivateHeader && decl.export_flag) return;

    if (decl.is(ir::IrDeclCase::TraitDecl)) {
      print_trait_decl(decl);
      return;
    }
    if (decl.is(ir::IrDeclCase::NameDecl)) return;

    if (decl.is(ir::IrDeclCase::AliasDecl) && decl.alias) {
      print_in_namespace(decl.alias->id, [&](const std::string& id) {
        if (decl.alias->type.is(ir::IrTypeCase::LambdaType)) {
          auto tvars = decl.alias->type.lambda_vars();
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

    if (decl.is(ir::IrDeclCase::TermDecl) && decl.term) {
      const auto& typ = decl.term->type;
      switch (typ.case_val) {
        case ir::IrTypeCase::AppType:
        case ir::IrTypeCase::ArrayType:
        case ir::IrTypeCase::ExistVarType:
        case ir::IrTypeCase::NameType:
        case ir::IrTypeCase::TupleType:
        case ir::IrTypeCase::VarType:
          print_in_namespace(decl.term->id, [&](const std::string& id) {
            with_bind_position([&]() {
              print_type(typ);
              out << " " << id;
            });
          });
          break;
        case ir::IrTypeCase::ForallType:
          print_in_namespace(decl.term->id, [&](const std::string& id) {
            auto tvars = typ.forall_type_params();
            print_template_params(tvars, false);
            print_decl(ir::new_term_decl(id, typ.forall_body(), decl.export_flag));
          });
          break;
        case ir::IrTypeCase::FunType:
          print_in_namespace(decl.term->id, [&](const std::string& id) {
            with_bind_position([&]() {
              if (typ.fun && typ.fun->ret) print_type(*typ.fun->ret);
            });
            out << " " << id << "(";
            if (typ.fun && typ.fun->arg) print_type(*typ.fun->arg);
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

    if (decl.is(ir::IrDeclCase::NameDecl) && decl.name) {
      print_type(ir::new_name_type(decl.name->id));
      out << ";\n";
    } else if (decl.is(ir::IrDeclCase::AliasDecl) && decl.alias) {
      print_in_namespace(decl.alias->id, [&](const std::string& id) {
        print_alias_decl(id, decl.alias->type);
        out << ";\n";
      });
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

  void print_function(const ir::IrFunction& f) {
    bool is_template = !f.type_params.empty();
    bool is_pub = f.export_flag;

    switch (mode) {
      case ir::PrinterMode::ModePublicHeader:
        if (!is_pub) return;
        if (is_template) print_function_full(f);
        else print_function_signature(f);
        break;
      case ir::PrinterMode::ModePrivateHeader:
        if (is_pub) return;
        if (is_template) print_function_full(f);
        else print_function_signature(f);
        break;
      case ir::PrinterMode::ModeSource:
        if (is_template) return;
        print_function_full(f);
        break;
    }
  }

  void print_trait_decl(const ir::IrDecl& decl) {
    if (!decl.trait) return;
    print_in_namespace(trait_cpp_name(decl.trait->id), [&](const std::string& id) {
      out << "template <typename Self";
      for (const auto& tp : decl.trait->type_params) {
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
        out << "  using Self = ";
        with_bind_position([&]() { print_type(impl.type_name); });
        out << ";\n";
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
      out << "  using Self = ";
      with_bind_position([&]() { print_type(impl.type_name); });
      out << ";\n";
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
        break;
      case ir::PrinterMode::ModePrivateHeader: {
        out << "#pragma once\n\n";
        if (unit) {
          std::string h_path = to_header_path(unit->module_id);
          out << "#include \"" << h_path << "\"\n\n";
        }
        break;
      }
      case ir::PrinterMode::ModeSource: {
        out << "\n";
        if (unit) {
          std::string h_path = to_header_path(unit->module_id);
          auto pos = h_path.rfind(".h");
          std::string priv_path = (pos != std::string::npos ? h_path.substr(0, pos) : h_path) + "_private.h";
          out << "#include \"" << priv_path << "\"\n\n";
        }
        break;
      }
    }
  }

  void print_imports() {
    if (!unit) return;
    if (mode == ir::PrinterMode::ModePublicHeader || mode == ir::PrinterMode::ModeSource) {
      out << "\n";
      for (const auto& imp : unit->imports) {
        out << "#include \"" << to_header_path(imp.module_id) << "\"\n";
      }
      out << "\n";
    }
  }

  void print_impls() {
    if (mode != ir::PrinterMode::ModePublicHeader || !unit) return;
    out << "\n";
    for (const auto& impl : unit->impls) {
      const std::string& fn = impl.relative_filename.value;
      if (fn.size() >= 2 && fn.substr(fn.size() - 2) == ".h") {
        out << "#include \"" << fn << "\"\n";
      }
    }
    out << "\n";
  }

  void print_decls(const std::vector<ir::IrDecl>& decls) {
    for (const auto& d : decls) {
      if (d.is(ir::IrDeclCase::AliasDecl)) {
        print_type_def(d);
      } else {
        print_decl(d);
      }
    }
  }

  void print_unit_content() {
    if (!unit) return;
    print_module_top();
    print_impls();
    print_imports();

    print_decls(unit->decls);
    if (mode == ir::PrinterMode::ModePrivateHeader) {
      print_decls(unit->impl_decls);
    }
    for (const auto& f : unit->functions) {
      print_function(f);
    }
    for (const auto& impl : unit->trait_impls) {
      print_trait_impl(impl);
    }
  }
};

// Anonymous struct extraction & topological sorting (Phase 8.3)

inline ir::IrType gen_name_type(ir::IrType typ, bool export_flag, std::map<std::string, AnonymousType>& anonymous_types) {
  std::string hash = std::to_string(std::hash<std::string>{}(typ.to_string()));
  std::string name = "__anonym_" + hash;
  auto it = anonymous_types.find(hash);
  if (it != anonymous_types.end()) {
    if (export_flag && !it->second.decl.export_flag) {
      it->second.decl.export_flag = true;
    }
  } else {
    anonymous_types[hash] = AnonymousType{
        ir::new_name_type(name),
        ir::new_alias_decl(name, ir::IrKind{}, typ, export_flag)
    };
  }
  return ir::new_name_type(name);
}

inline ir::IrType record_anonymous_types(ir::IrType typ, bool export_flag, std::map<std::string, AnonymousType>& anonymous_types) {
  switch (typ.case_val) {
    case ir::IrTypeCase::AppType: {
      if (typ.app && typ.app->fun && typ.app->arg) {
        return ir::new_app_type(record_anonymous_types(*typ.app->fun, export_flag, anonymous_types),
                                record_anonymous_types(*typ.app->arg, export_flag, anonymous_types));
      }
      return typ;
    }
    case ir::IrTypeCase::ArrayType: {
      if (typ.array && typ.array->elem_type) {
        return ir::new_array_type(record_anonymous_types(*typ.array->elem_type, export_flag, anonymous_types), typ.array->size);
      }
      return typ;
    }
    case ir::IrTypeCase::ForallType: {
      if (typ.forall && typ.forall->type) {
        ir::TypeParam tp = typ.forall->type_param;
        for (auto& b : tp.bounds) {
          b = record_anonymous_types(b, export_flag, anonymous_types);
        }
        return ir::new_forall_type(std::move(tp), record_anonymous_types(*typ.forall->type, export_flag, anonymous_types));
      }
      return typ;
    }
    case ir::IrTypeCase::FunType: {
      if (typ.fun && typ.fun->arg && typ.fun->ret) {
        return ir::new_function_type(record_anonymous_types(*typ.fun->arg, export_flag, anonymous_types),
                                     record_anonymous_types(*typ.fun->ret, export_flag, anonymous_types));
      }
      return typ;
    }
    case ir::IrTypeCase::LambdaType: {
      if (typ.lambda && typ.lambda->type) {
        return ir::new_lambda_type(typ.lambda->var, typ.lambda->kind,
                                   record_anonymous_types(*typ.lambda->type, export_flag, anonymous_types));
      }
      return typ;
    }
    case ir::IrTypeCase::StructType:
      return gen_name_type(typ, export_flag, anonymous_types);
    case ir::IrTypeCase::TupleType: {
      std::vector<ir::IrType> elems;
      for (const auto& e : typ.elems()) {
        elems.push_back(record_anonymous_types(e, export_flag, anonymous_types));
      }
      return ir::new_tuple_type(std::move(elems));
    }
    case ir::IrTypeCase::VariantType: {
      std::vector<ir::VariantTag> tags;
      for (const auto& tag : typ.tags()) {
        ir::VariantTag nt = tag;
        if (nt.type) {
          nt.type = std::make_shared<ir::IrType>(record_anonymous_types(*nt.type, export_flag, anonymous_types));
        }
        tags.push_back(std::move(nt));
      }
      return ir::new_variant_type(std::move(tags));
    }
    default:
      return typ;
  }
}

inline void get_free_vars_from_type(const ir::IrType& typ, std::set<std::string>& bound_names, std::vector<ir::IrType>& free_vars) {
  if (typ.is(ir::IrTypeCase::NameType)) {
    if (bound_names.find(typ.name) == bound_names.end()) {
      bound_names.insert(typ.name);
      free_vars.push_back(typ);
    }
    return;
  }
  if (typ.is(ir::IrTypeCase::AppType) && typ.app) {
    if (typ.app->fun) get_free_vars_from_type(*typ.app->fun, bound_names, free_vars);
    if (typ.app->arg) get_free_vars_from_type(*typ.app->arg, bound_names, free_vars);
  } else if (typ.is(ir::IrTypeCase::ArrayType) && typ.array && typ.array->elem_type) {
    get_free_vars_from_type(*typ.array->elem_type, bound_names, free_vars);
  } else if (typ.is(ir::IrTypeCase::FunType) && typ.fun) {
    if (typ.fun->arg) get_free_vars_from_type(*typ.fun->arg, bound_names, free_vars);
    if (typ.fun->ret) get_free_vars_from_type(*typ.fun->ret, bound_names, free_vars);
  } else if (typ.is(ir::IrTypeCase::StructType)) {
    for (const auto& f : typ.fields()) {
      if (f.type) get_free_vars_from_type(*f.type, bound_names, free_vars);
    }
  } else if (typ.is(ir::IrTypeCase::TupleType)) {
    for (const auto& e : typ.elems()) {
      get_free_vars_from_type(e, bound_names, free_vars);
    }
  } else if (typ.is(ir::IrTypeCase::VariantType)) {
    for (const auto& tag : typ.tags()) {
      if (tag.type) get_free_vars_from_type(*tag.type, bound_names, free_vars);
    }
  }
}

inline std::vector<ir::IrDecl> topo_sort_alias_decls(const std::vector<ir::IrDecl>& decls) {
  std::map<std::string, int> nodes_by_name;
  for (size_t i = 0; i < decls.size(); ++i) {
    nodes_by_name[decls[i].alias ? decls[i].alias->id : ""] = static_cast<int>(i);
  }

  std::vector<std::vector<int>> adj(decls.size());
  std::vector<int> in_degree(decls.size(), 0);

  for (size_t i = 0; i < decls.size(); ++i) {
    if (!decls[i].alias) continue;
    std::set<std::string> bound;
    std::vector<ir::IrType> free_vars;
    get_free_vars_from_type(decls[i].alias->type, bound, free_vars);
    for (const auto& fvar : free_vars) {
      if (fvar.is(ir::IrTypeCase::NameType)) {
        auto it = nodes_by_name.find(fvar.name);
        if (it != nodes_by_name.end() && it->second != static_cast<int>(i)) {
          adj[it->second].push_back(static_cast<int>(i));
          in_degree[i]++;
        }
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
    switch (d.case_val) {
      case ir::IrDeclCase::TermDecl: term_decls.push_back(d); break;
      case ir::IrDeclCase::AliasDecl: alias_decls.push_back(d); break;
      case ir::IrDeclCase::NameDecl: name_decls.push_back(d); break;
      case ir::IrDeclCase::TraitDecl: trait_decls.push_back(d); break;
    }
  }

  alias_decls = topo_sort_alias_decls(alias_decls);

  decls.clear();
  decls.insert(decls.end(), name_decls.begin(), name_decls.end());
  decls.insert(decls.end(), alias_decls.begin(), alias_decls.end());
  decls.insert(decls.end(), trait_decls.begin(), trait_decls.end());
  decls.insert(decls.end(), term_decls.begin(), term_decls.end());
}

inline void record_anonymous_types_from_unit(ir::IrUnit* unit, std::map<std::string, AnonymousType>& anonymous_types) {
  if (!unit) return;
  for (auto& d : unit->decls) {
    if (d.is(ir::IrDeclCase::TermDecl) && d.term) {
      d.term->type = record_anonymous_types(d.term->type, d.export_flag, anonymous_types);
    }
  }
  for (auto& d : unit->impl_decls) {
    if (d.is(ir::IrDeclCase::TermDecl) && d.term) {
      d.term->type = record_anonymous_types(d.term->type, d.export_flag, anonymous_types);
    }
  }
  for (auto& f : unit->functions) {
    f.ret_type = record_anonymous_types(f.ret_type, f.export_flag, anonymous_types);
  }

  for (const auto& [k, v] : anonymous_types) {
    unit->decls.push_back(v.decl);
  }
  topo_sort_decls(unit->decls);
}

inline void format_file(const std::string& filename) {
  std::string cmd = "clang-format -i " + filename + " 2>/dev/null";
  std::system(cmd.c_str());
}

inline bool print_unit_to_file(ir::IrUnit unit, ir::PrinterMode mode, const std::string& filename,
                               const std::map<std::string, AnonymousType>& atypes) {
  std::ofstream f(filename);
  if (!f.is_open()) return false;
  CppPrinter printer(f, mode, &unit);
  printer.anonymous_types = atypes;
  printer.print_unit_content();
  f.close();
  format_file(filename);
  return true;
}

inline bool compile_ir_unit_direct(ir::IrUnit unit, const std::string& output_base) {
  std::map<std::string, AnonymousType> atypes;
  record_anonymous_types_from_unit(&unit, atypes);

  std::string base = output_base;
  auto dot = base.rfind('.');
  if (dot != std::string::npos) {
    base = base.substr(0, dot);
  }

  if (unit.case_val == ir::IrUnitCase::BaseUnit) {
    if (!print_unit_to_file(unit, ir::PrinterMode::ModePublicHeader, base + ".h", atypes)) return false;
    if (!print_unit_to_file(unit, ir::PrinterMode::ModePrivateHeader, base + "_private.h", atypes)) return false;
    if (!print_unit_to_file(unit, ir::PrinterMode::ModeSource, base + ".cc", atypes)) return false;
  } else {
    if (!print_unit_to_file(unit, ir::PrinterMode::ModeSource, base + ".cc", atypes)) return false;
  }
  return true;
}

// @bpl: pub codegen::compile_unit: (String, String) -> i64
inline int64_t compile_unit(const std::string& input_file, const std::string& output_base) {
  comp::ModuleFinder finder({}, {{"", "."}});
  comp::Querier querier(finder);
  comp::TypecheckOptions options;

  ir::IrUnit unit;
  std::string err;
  if (!comp::typecheck_source_file(querier, options, input_file, unit, err)) {
    std::cerr << "Typecheck error in " << input_file << ":\n" << err << "\n";
    return 1;
  }

  return compile_ir_unit_direct(std::move(unit), output_base) ? 0 : 1;
}

} // namespace codegen

namespace typechecker {

// @bpl: pub typechecker::run: Vector String -> (i64, String)
inline std::tuple<int64_t, std::string> run(const std::vector<std::string>& args) {
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
        return {1, "Usage: typechecker [-format=flat|json|ir] <input_file>\n"};
      }
    }
  }

  if (input_file.empty()) {
    return {1, "Usage: typechecker [-format=flat|json|ir] <input_file>\n"};
  }

  comp::ModuleFinder finder({}, {{"", "."}});
  comp::Querier querier(finder);
  comp::TypecheckOptions options;

  ir::IrUnit unit;
  std::string err;
  if (!comp::typecheck_source_file(querier, options, input_file, unit, err)) {
    return {1, "Failed to typecheck \"" + input_file + "\": " + err + "\n"};
  }

  std::stringstream out;
  if (format == "json") {
    out << unit.to_json() << "\n";
  } else if (format == "ir") {
    out << unit.to_string() << "\n";
  } else if (format == "flat") {
    out << "MODULE " << unit.module_id.to_string() << "\n";
    if (unit.case_val == ir::IrUnitCase::BaseUnit) {
      out << "CASE base\n";
    } else {
      out << "CASE impl\n";
    }

    for (const auto& imp : unit.imports) {
      out << "IMPORT " << imp.module_id.to_string() << "\n";
    }
    for (const auto& impl : unit.impls) {
      out << "IMPL " << impl.relative_filename.value << "\n";
    }
    for (const auto& decl : unit.decls) {
      std::string s = decl.to_string();
      std::string escaped_s = s;
      comp::replace_all_str(escaped_s, "\n", "\\n");
      out << "DECL " << escaped_s << "\n";
      std::string export_str = decl.export_flag ? "1" : "0";
      out << "DECL_DEF " << export_str << " " << decl.id() << " " << escaped_s << "\n";
    }
    for (const auto& trait_impl : unit.trait_impls) {
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
    for (const auto& fn : unit.functions) {
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
