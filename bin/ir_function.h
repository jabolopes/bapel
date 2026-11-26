#pragma once

#include "ir_term.h"

namespace ir {

using ::IrFunction;

inline std::string function_to_string(const IrFunction& f, bool with_pos = false) {
  std::stringstream ss;
  if (with_pos && !f.pos.filename.empty()) {
    ss << pos_to_string(f.pos, true);
  }
  if (f.export_flag) {
    ss << "pub ";
  }
  ss << "fn " << f.id;
  if (!f.type_params.empty()) {
    ss << "[";
    Interleave(f.type_params, [&]() { ss << ", "; }, [&](int, const TypeParam& tp) {
      ss << "'" << tp.var << " " << kind_to_string(tp.kind);
    });
    ss << "]";
  }
  ss << "(";
  Interleave(f.args, [&]() { ss << ", "; }, [&](int, const FunctionArg& a) {
    ss << function_arg_to_string(a);
  });
  ss << ") -> " << type_to_string(f.ret_type) << " " << term_to_string(f.body);
  return ss.str();
}

inline IrFunction new_function(
    bool export_flag,
    std::string id,
    std::vector<TypeParam> type_params,
    std::vector<FunctionArg> args,
    IrType ret_type,
    IrTerm body) {
  IrFunction f{};
  f.id = std::move(id);
  f.type_params = std::move(type_params);
  f.args = std::move(args);
  f.ret_type = std::move(ret_type);
  f.body = std::move(body);
  f.export_flag = export_flag;
  return f;
}

} // namespace ir
