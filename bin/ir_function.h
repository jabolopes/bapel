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
};

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
