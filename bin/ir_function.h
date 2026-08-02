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
  std::string to_json() const;
  std::string to_string() const;
};

inline std::string IrFunction::to_json() const {
  std::stringstream ss;
  ss << "{\"Export\":" << (export_flag ? "true" : "false")
     << ",\"ID\":\"" << json_escape(id) << "\""
     << ",\"TypeParams\":[";
  Interleave(type_params, [&]() { ss << ","; }, [&](int, const TypeParam& tp) {
    ss << tp.to_json();
  });
  ss << "],\"Args\":[";
  Interleave(args, [&]() { ss << ","; }, [&](int, const FunctionArg& a) {
    ss << a.to_json();
  });
  ss << "],\"RetType\":" << ret_type.to_json()
     << ",\"Body\":" << body.to_json()
     << ",\"Pos\":" << pos.to_json() << "}";
  return ss.str();
}

inline std::string IrFunction::to_string() const {
  std::stringstream ss;
  if (!pos.filename.empty()) {
    ss << pos.to_string(true) << " ";
  }
  if (export_flag) {
    ss << "pub ";
  }
  ss << "fn " << id;
  if (!type_params.empty()) {
    ss << " [";
    Interleave(type_params, [&]() { ss << ", "; }, [&](int, const TypeParam& tp) {
      ss << "'" << tp.var;
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
