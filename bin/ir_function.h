#pragma once

#include "ir_term.h"

namespace ir {

struct IrDecl;

struct IrFunction {
  bool export_flag = false;
  std::string id;
  std::vector<TypeParam> type_params;
  std::vector<FunctionArg> args;
  IrType ret_type;
  IrTerm body;
  Pos pos;

  IrDecl decl() const;
  std::string to_string(bool with_pos = false) const;
};

inline std::string IrFunction::to_string(bool with_pos) const {
  std::stringstream ss;
  if (with_pos && !pos.filename.empty()) {
    ss << pos.to_string(true);
  }
  if (export_flag) {
    ss << "export ";
  }
  ss << "fn " << id;
  if (!type_params.empty()) {
    ss << "[";
    Interleave(type_params, [&]() { ss << ", "; }, [&](int, const TypeParam& tp) {
      ss << "'" << tp.var << " " << tp.kind.to_string();
    });
    ss << "]";
  }
  ss << "(";
  Interleave(args, [&]() { ss << ", "; }, [&](int, const FunctionArg& a) {
    ss << a.to_string();
  });
  ss << ") -> " << ret_type.to_string() << " " << body.to_string();
  return ss.str();
}

inline IrFunction new_function(
    bool export_flag,
    std::string id,
    std::vector<TypeParam> type_params,
    std::vector<FunctionArg> args,
    IrType ret_type,
    IrTerm body) {
  IrFunction f;
  f.export_flag = export_flag;
  f.id = std::move(id);
  f.type_params = std::move(type_params);
  f.args = std::move(args);
  f.ret_type = std::move(ret_type);
  f.body = std::move(body);
  return f;
}

} // namespace ir
