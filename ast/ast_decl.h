#pragma once

#include "ast_expr.h"
#include "ast_pos.h"
#include "bin/ir_decl.h"
#include "bin/ir_function.h"
#include "bin/ir_type.h"
#include <memory>
#include <sstream>
#include <string>
#include <vector>

namespace ast {

struct Signature {
  std::string id;
  std::vector<ir::FunctionArg> args;
  ir::IrType ret_type;
  Pos pos;

  std::string to_string(bool with_pos = false) const {
    std::stringstream ss;
    if (with_pos) ss << pos.to_string() << " ";
    ss << "fn " << id << "(";
    ir::Interleave(args, [&]() { ss << ", "; }, [&](int, const ir::FunctionArg& a) {
      ss << a.to_string();
    });
    ss << ") -> " << ret_type.to_string();
    return ss.str();
  }
};

struct Function {
  bool export_flag = false;
  std::string id;
  std::vector<ir::TypeParam> type_params;
  std::vector<ir::FunctionArg> args;
  ir::IrType ret_type;
  Expr body;
  Pos pos;

  ir::IrDecl decl() const {
    std::vector<ir::IrType> arg_types;
    arg_types.reserve(args.size());
    for (const auto& a : args) {
      arg_types.push_back(a.type);
    }
    ir::IrType fn_type = ir::new_function_type(ir::new_tuple_type(arg_types), ret_type);
    ir::IrType typ = type_params.empty() ? fn_type : ir::forall_vars(type_params, fn_type);
    ir::IrDecl d = ir::new_term_decl(id, typ, export_flag);
    d.pos = pos;
    return d;
  }

  std::string to_string(bool with_pos = false) const {
    std::stringstream ss;
    if (with_pos) ss << pos.to_string() << " ";
    if (export_flag) ss << "pub ";
    ss << "fn " << id;
    if (!type_params.empty()) {
      ss << " [";
      ir::Interleave(type_params, [&]() { ss << ", "; }, [&](int, const ir::TypeParam& tp) {
        ss << "'" << tp.var;
        if (!tp.bounds.empty()) {
          ss << ": ";
          ir::Interleave(tp.bounds, [&]() { ss << " + "; }, [&](int, const ir::IrType& b) {
            ss << b.to_string();
          });
        }
      });
      ss << "]";
    }
    ss << "(";
    ir::Interleave(args, [&]() { ss << ", "; }, [&](int, const ir::FunctionArg& a) {
      ss << a.to_string();
    });
    ss << ") -> " << ret_type.to_string() << " " << body.to_string(with_pos);
    return ss.str();
  }
};

struct Trait {
  bool export_flag = false;
  std::string id;
  std::vector<ir::TypeParam> type_params;
  std::vector<Signature> methods;
  Pos pos;

  ir::IrDecl decl() const {
    std::vector<ir::IrSignature> ir_methods;
    ir_methods.reserve(methods.size());
    for (const auto& m : methods) {
      ir_methods.push_back(ir::IrSignature{m.id, m.args, m.ret_type});
    }
    ir::IrDecl d = ir::new_trait_decl(id, type_params, ir_methods, export_flag);
    d.pos = pos;
    return d;
  }

  std::string to_string(bool with_pos = false) const {
    std::stringstream ss;
    if (with_pos) ss << pos.to_string() << " ";
    if (export_flag) ss << "pub ";
    ss << "trait " << id;
    if (!type_params.empty()) {
      ss << " [";
      ir::Interleave(type_params, [&]() { ss << ", "; }, [&](int, const ir::TypeParam& tv) {
        ss << "'" << tv.var;
      });
      ss << "]";
    }
    ss << " {\n";
    for (const auto& m : methods) {
      ss << "  " << m.to_string(with_pos) << "\n";
    }
    ss << "}";
    return ss.str();
  }
};

enum class ImplCase {
  TraitImpl = 0,
  InherentImpl = 1,
};

struct Impl {
  ImplCase case_val = ImplCase::TraitImpl;
  std::vector<ir::TypeParam> type_params;
  ir::IrType trait_type;
  ir::IrType type_name;
  std::vector<Function> methods;
  Pos pos;

  std::string to_string(bool with_pos = false) const {
    std::stringstream ss;
    if (with_pos) ss << pos.to_string() << " ";
    ss << "impl";
    if (!type_params.empty()) {
      ss << " [";
      ir::Interleave(type_params, [&]() { ss << ", "; }, [&](int, const ir::TypeParam& tv) {
        ss << "'" << tv.var;
      });
      ss << "]";
    }
    if (case_val == ImplCase::InherentImpl) {
      ss << " " << type_name.to_string() << " {\n";
    } else {
      ss << " " << trait_type.to_string() << " for " << type_name.to_string() << " {\n";
    }
    for (const auto& m : methods) {
      ss << "  " << m.to_string(with_pos) << "\n";
    }
    ss << "}";
    return ss.str();
  }
};

enum class SourceCase {
  DeclSource = 0,
  FunctionSource = 1,
  TraitSource = 2,
  ImplSource = 3,
};

struct Source {
  SourceCase case_val = SourceCase::DeclSource;
  std::shared_ptr<ir::IrDecl> decl_data;
  std::shared_ptr<Function> function_data;
  std::shared_ptr<Trait> trait_data;
  std::shared_ptr<Impl> impl_data;

  bool is(SourceCase c) const { return case_val == c; }

  std::string to_string(bool with_pos = false) const {
    switch (case_val) {
      case SourceCase::DeclSource:
        return decl_data ? decl_data->to_string(with_pos) : "";
      case SourceCase::FunctionSource:
        return function_data ? function_data->to_string(with_pos) : "";
      case SourceCase::TraitSource:
        return trait_data ? trait_data->to_string(with_pos) : "";
      case SourceCase::ImplSource:
        return impl_data ? impl_data->to_string(with_pos) : "";
    }
    return "";
  }
};

inline Source new_decl_source(ir::IrDecl decl) {
  Source s;
  s.case_val = SourceCase::DeclSource;
  s.decl_data = std::make_shared<ir::IrDecl>(std::move(decl));
  return s;
}

inline Source new_function_source(Function function) {
  Source s;
  s.case_val = SourceCase::FunctionSource;
  s.function_data = std::make_shared<Function>(std::move(function));
  return s;
}

inline Source new_trait_source(Trait trait) {
  Source s;
  s.case_val = SourceCase::TraitSource;
  s.trait_data = std::make_shared<Trait>(std::move(trait));
  return s;
}

inline Source new_impl_source(Impl impl) {
  Source s;
  s.case_val = SourceCase::ImplSource;
  s.impl_data = std::make_shared<Impl>(std::move(impl));
  return s;
}

} // namespace ast
