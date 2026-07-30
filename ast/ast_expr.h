#pragma once

#include "ast_pos.h"
#include "bin/ir_type.h"
#include <memory>
#include <optional>
#include <sstream>
#include <string>
#include <vector>

namespace ast {

enum class ExprCase {
  AppTermExpr = 0,
  AppTypeExpr = 1,
  AssignExpr = 2,
  BlockExpr = 3,
  ConstExpr = 4,
  InjectionExpr = 5,
  LambdaExpr = 6,
  LetExpr = 7,
  MatchExpr = 8,
  ProjectionExpr = 9,
  ReturnExpr = 10,
  SetExpr = 11,
  StructExpr = 12,
  TupleExpr = 13,
  TypeAbsExpr = 14,
  VarExpr = 15,
  ForExpr = 16,
};

struct Expr;

struct MatchArm {
  std::string tag;
  std::string arg;
  std::shared_ptr<Expr> body;
  std::optional<int> index;

  std::string to_string(bool with_pos = false) const;
  std::string to_json() const;
};

struct LabelValue {
  std::string label;
  std::shared_ptr<Expr> value;

  std::string to_string(bool with_pos = false) const;
  std::string to_json() const;
};

struct AppTermData {
  std::shared_ptr<Expr> fun;
  std::shared_ptr<Expr> arg;
};

struct AppTypeData {
  std::shared_ptr<Expr> fun;
  ir::IrType arg;
};

struct AssignData {
  std::shared_ptr<Expr> arg;
  std::shared_ptr<Expr> ret;
};

struct BlockData {
  std::vector<Expr> exprs;
};

struct ConstData {
  ir::IrLiteral literal;
};

struct InjectionData {
  ir::IrType variant_type;
  std::string tag;
  std::shared_ptr<Expr> expr;
  std::optional<int> tag_index;
};

struct LambdaData {
  ir::FunctionArg arg;
  std::shared_ptr<Expr> body;
};

struct LetData {
  std::string var;
  std::optional<ir::IrType> var_type;
  std::shared_ptr<Expr> expr;
};

struct MatchData {
  std::shared_ptr<Expr> expr;
  std::vector<MatchArm> arms;
};

struct ProjectionData {
  std::shared_ptr<Expr> expr;
  std::string label;
};

struct ReturnData {
  std::shared_ptr<Expr> expr;
};

struct SetData {
  std::shared_ptr<Expr> expr;
  std::vector<LabelValue> values;
};

struct StructData {
  std::vector<LabelValue> values;
};

struct TupleData {
  std::vector<Expr> elems;
};

struct TypeAbsData {
  ir::TypeParam arg;
  std::shared_ptr<Expr> body;
};

struct VarData {
  std::string id;
};

struct ForData {
  std::shared_ptr<Expr> condition;
  std::shared_ptr<Expr> body;
};

struct Expr {
  ExprCase case_val = ExprCase::VarExpr;
  Pos pos;

  std::shared_ptr<AppTermData> app_term_data;
  std::shared_ptr<AppTypeData> app_type_data;
  std::shared_ptr<AssignData> assign_data;
  std::shared_ptr<BlockData> block_data;
  std::shared_ptr<ConstData> const_data;
  std::shared_ptr<InjectionData> injection_data;
  std::shared_ptr<LambdaData> lambda_data;
  std::shared_ptr<LetData> let_data;
  std::shared_ptr<MatchData> match_data;
  std::shared_ptr<ProjectionData> projection_data;
  std::shared_ptr<ReturnData> return_data;
  std::shared_ptr<SetData> set_data;
  std::shared_ptr<StructData> struct_data;
  std::shared_ptr<TupleData> tuple_data;
  std::shared_ptr<TypeAbsData> type_abs_data;
  std::shared_ptr<VarData> var_data;
  std::shared_ptr<ForData> for_data;

  bool is(ExprCase c) const { return case_val == c; }

  std::string to_string(bool with_pos = false) const {
    std::stringstream ss;
    switch (case_val) {
      case ExprCase::AppTermExpr: {
        if (!app_term_data) return "";
        std::string fun_str = app_term_data->fun ? app_term_data->fun->to_string(with_pos) : "";
        std::string arg_str = app_term_data->arg ? app_term_data->arg->to_string(with_pos) : "";
        bool fun_paren = app_term_data->fun && app_term_data->fun->is(ExprCase::LambdaExpr);
        bool arg_paren = app_term_data->arg && (app_term_data->arg->is(ExprCase::AppTermExpr) || app_term_data->arg->is(ExprCase::LambdaExpr));
        if (fun_paren) ss << "(";
        ss << fun_str;
        if (fun_paren) ss << ")";
        ss << " ";
        if (arg_paren) ss << "(";
        ss << arg_str;
        if (arg_paren) ss << ")";
        break;
      }
      case ExprCase::AppTypeExpr: {
        if (!app_type_data) return "";
        std::string fun_str = app_type_data->fun ? app_type_data->fun->to_string(with_pos) : "";
        ss << fun_str << " [" << app_type_data->arg.to_string() << "]";
        break;
      }
      case ExprCase::AssignExpr: {
        if (!assign_data) return "";
        std::string ret_str = assign_data->ret ? assign_data->ret->to_string(with_pos) : "";
        std::string arg_str = assign_data->arg ? assign_data->arg->to_string(with_pos) : "";
        ss << ret_str << " <- " << arg_str;
        break;
      }
      case ExprCase::BlockExpr: {
        if (!block_data) return "{}";
        ss << "{\n";
        for (const auto& e : block_data->exprs) {
          ss << "  " << e.to_string(with_pos) << "\n";
        }
        ss << "}";
        break;
      }
      case ExprCase::ConstExpr: {
        if (!const_data) return "";
        ss << const_data->literal.to_string();
        break;
      }
      case ExprCase::InjectionExpr: {
        if (!injection_data) return "";
        ss << "variant{" << injection_data->variant_type.to_string() << " " << injection_data->tag << " = "
           << (injection_data->expr ? injection_data->expr->to_string(with_pos) : "") << "}";
        break;
      }
      case ExprCase::LambdaExpr: {
        if (!lambda_data) return "";
        ss << "\\(" << lambda_data->arg.to_string() << ") -> " << (lambda_data->body ? lambda_data->body->to_string(with_pos) : "");
        break;
      }
      case ExprCase::LetExpr: {
        if (!let_data) return "";
        if (let_data->var_type.has_value()) {
          ss << "let " << let_data->var << ": " << let_data->var_type->to_string() << " = "
             << (let_data->expr ? let_data->expr->to_string(with_pos) : "");
        } else {
          ss << "let " << let_data->var << " = " << (let_data->expr ? let_data->expr->to_string(with_pos) : "");
        }
        break;
      }
      case ExprCase::MatchExpr: {
        if (!match_data) return "";
        ss << "case " << (match_data->expr ? match_data->expr->to_string(with_pos) : "") << " {";
        if (match_data->arms.size() == 1) {
          ss << " " << match_data->arms[0].to_string(with_pos) << " }";
        } else {
          ss << "\n";
          for (const auto& a : match_data->arms) {
            ss << "    " << a.to_string(with_pos) << "\n";
          }
          ss << "}";
        }
        break;
      }
      case ExprCase::ProjectionExpr: {
        if (!projection_data) return "";
        ss << (projection_data->expr ? projection_data->expr->to_string(with_pos) : "") << "." << projection_data->label;
        break;
      }
      case ExprCase::ReturnExpr: {
        if (!return_data) return "";
        ss << "return " << (return_data->expr ? return_data->expr->to_string(with_pos) : "");
        break;
      }
      case ExprCase::SetExpr: {
        if (!set_data) return "";
        ss << "set " << (set_data->expr ? set_data->expr->to_string(with_pos) : "") << " {";
        ir::Interleave(set_data->values, [&]() { ss << ", "; }, [&](int, const LabelValue& lv) {
          ss << lv.to_string(with_pos);
        });
        ss << "}";
        break;
      }
      case ExprCase::StructExpr: {
        if (!struct_data) return "struct{}";
        ss << "struct{";
        ir::Interleave(struct_data->values, [&]() { ss << ", "; }, [&](int, const LabelValue& lv) {
          ss << lv.to_string(with_pos);
        });
        ss << "}";
        break;
      }
      case ExprCase::TupleExpr: {
        if (!tuple_data) return "()";
        ss << "(";
        ir::Interleave(tuple_data->elems, [&]() { ss << ", "; }, [&](int, const Expr& e) {
          ss << e.to_string(with_pos);
        });
        ss << ")";
        break;
      }
      case ExprCase::TypeAbsExpr: {
        if (!type_abs_data) return "";
        ss << "\xce\x9b" << type_abs_data->arg.to_string() << ". "
           << (type_abs_data->body ? type_abs_data->body->to_string(with_pos) : "");
        break;
      }
      case ExprCase::VarExpr: {
        if (!var_data) return "";
        ss << var_data->id;
        break;
      }
      case ExprCase::ForExpr: {
        if (!for_data) return "";
        ss << "for " << (for_data->condition ? for_data->condition->to_string(with_pos) : "")
           << " " << (for_data->body ? for_data->body->to_string(with_pos) : "");
        break;
      }
    }
    return ss.str();
  }

  std::string to_json() const {
    std::stringstream ss;
    ss << "{\"Case\":" << static_cast<int>(case_val);
    switch (case_val) {
      case ExprCase::AppTermExpr:
        if (app_term_data) {
          ss << ",\"AppTerm\":{\"Fun\":" << (app_term_data->fun ? app_term_data->fun->to_json() : "null")
             << ",\"Arg\":" << (app_term_data->arg ? app_term_data->arg->to_json() : "null") << "}";
        }
        break;
      case ExprCase::AppTypeExpr:
        if (app_type_data) {
          ss << ",\"AppType\":{\"Fun\":" << (app_type_data->fun ? app_type_data->fun->to_json() : "null")
             << ",\"Arg\":" << app_type_data->arg.to_json() << "}";
        }
        break;
      case ExprCase::AssignExpr:
        if (assign_data) {
          ss << ",\"Assign\":{\"Arg\":" << (assign_data->arg ? assign_data->arg->to_json() : "null")
             << ",\"Ret\":" << (assign_data->ret ? assign_data->ret->to_json() : "null") << "}";
        }
        break;
      case ExprCase::BlockExpr:
        if (block_data) {
          ss << ",\"Block\":{\"Exprs\":[";
          ir::Interleave(block_data->exprs, [&]() { ss << ","; }, [&](int, const Expr& e) {
            ss << e.to_json();
          });
          ss << "]}";
        }
        break;
      case ExprCase::ConstExpr:
        if (const_data) {
          ss << ",\"Const\":" << const_data->literal.to_json();
        }
        break;
      case ExprCase::InjectionExpr:
        if (injection_data) {
          ss << ",\"Injection\":{\"VariantType\":" << injection_data->variant_type.to_json()
             << ",\"Tag\":\"" << ir::json_escape(injection_data->tag) << "\""
             << ",\"Expr\":" << (injection_data->expr ? injection_data->expr->to_json() : "null");
          if (injection_data->tag_index.has_value()) {
            ss << ",\"TagIndex\":" << *injection_data->tag_index;
          } else {
            ss << ",\"TagIndex\":null";
          }
          ss << "}";
        }
        break;
      case ExprCase::LambdaExpr:
        if (lambda_data) {
          ss << ",\"Lambda\":{\"Arg\":" << lambda_data->arg.to_json()
             << ",\"Body\":" << (lambda_data->body ? lambda_data->body->to_json() : "null") << "}";
        }
        break;
      case ExprCase::LetExpr:
        if (let_data) {
          ss << ",\"Let\":{\"Var\":\"" << ir::json_escape(let_data->var) << "\"";
          if (let_data->var_type.has_value()) {
            ss << ",\"VarType\":" << let_data->var_type->to_json();
          } else {
            ss << ",\"VarType\":null";
          }
          ss << ",\"Expr\":" << (let_data->expr ? let_data->expr->to_json() : "null") << "}";
        }
        break;
      case ExprCase::MatchExpr:
        if (match_data) {
          ss << ",\"Match\":{\"Expr\":" << (match_data->expr ? match_data->expr->to_json() : "null")
             << ",\"Arms\":[";
          ir::Interleave(match_data->arms, [&]() { ss << ","; }, [&](int, const MatchArm& a) {
            ss << a.to_json();
          });
          ss << "]}";
        }
        break;
      case ExprCase::ProjectionExpr:
        if (projection_data) {
          ss << ",\"Projection\":{\"Expr\":" << (projection_data->expr ? projection_data->expr->to_json() : "null")
             << ",\"Label\":\"" << ir::json_escape(projection_data->label) << "\"}";
        }
        break;
      case ExprCase::ReturnExpr:
        if (return_data) {
          ss << ",\"Return\":{\"Expr\":" << (return_data->expr ? return_data->expr->to_json() : "null") << "}";
        }
        break;
      case ExprCase::SetExpr:
        if (set_data) {
          ss << ",\"Set\":{\"Expr\":" << (set_data->expr ? set_data->expr->to_json() : "null")
             << ",\"Values\":[";
          ir::Interleave(set_data->values, [&]() { ss << ","; }, [&](int, const LabelValue& lv) {
            ss << lv.to_json();
          });
          ss << "]}";
        }
        break;
      case ExprCase::StructExpr:
        if (struct_data) {
          ss << ",\"Struct\":{\"Values\":[";
          ir::Interleave(struct_data->values, [&]() { ss << ","; }, [&](int, const LabelValue& lv) {
            ss << lv.to_json();
          });
          ss << "]}";
        }
        break;
      case ExprCase::TupleExpr:
        if (tuple_data) {
          ss << ",\"Tuple\":{\"Elems\":[";
          ir::Interleave(tuple_data->elems, [&]() { ss << ","; }, [&](int, const Expr& e) {
            ss << e.to_json();
          });
          ss << "]}";
        }
        break;
      case ExprCase::TypeAbsExpr:
        if (type_abs_data) {
          ss << ",\"TypeAbs\":{\"Arg\":" << type_abs_data->arg.to_json()
             << ",\"Body\":" << (type_abs_data->body ? type_abs_data->body->to_json() : "null") << "}";
        }
        break;
      case ExprCase::VarExpr:
        if (var_data) {
          ss << ",\"Var\":{\"ID\":\"" << ir::json_escape(var_data->id) << "\"}";
        }
        break;
      case ExprCase::ForExpr:
        if (for_data) {
          ss << ",\"For\":{\"Condition\":" << (for_data->condition ? for_data->condition->to_json() : "null")
             << ",\"Body\":" << (for_data->body ? for_data->body->to_json() : "null") << "}";
        }
        break;
    }
    ss << ",\"Pos\":" << pos.to_json() << "}";
    return ss.str();
  }
};

inline std::string MatchArm::to_string(bool with_pos) const {
  std::stringstream ss;
  ss << tag << " " << arg << " -> " << (body ? body->to_string(with_pos) : "");
  return ss.str();
}

inline std::string MatchArm::to_json() const {
  std::stringstream ss;
  ss << "{\"Tag\":\"" << ir::json_escape(tag) << "\",\"Arg\":\"" << ir::json_escape(arg) << "\""
     << ",\"Body\":" << (body ? body->to_json() : "null");
  if (index.has_value()) {
    ss << ",\"Index\":" << *index;
  } else {
    ss << ",\"Index\":null";
  }
  ss << "}";
  return ss.str();
}

inline std::string LabelValue::to_string(bool with_pos) const {
  std::stringstream ss;
  ss << label << " = " << (value ? value->to_string(with_pos) : "");
  return ss.str();
}

inline std::string LabelValue::to_json() const {
  std::stringstream ss;
  ss << "{\"Label\":\"" << ir::json_escape(label) << "\",\"Value\":" << (value ? value->to_json() : "null") << "}";
  return ss.str();
}

// Helper Constructors
inline Expr new_var_expr(Pos pos, std::string id) {
  Expr e;
  e.case_val = ExprCase::VarExpr;
  e.pos = pos;
  e.var_data = std::make_shared<VarData>(VarData{std::move(id)});
  return e;
}

inline Expr new_const_expr(Pos pos, ir::IrLiteral lit) {
  Expr e;
  e.case_val = ExprCase::ConstExpr;
  e.pos = pos;
  e.const_data = std::make_shared<ConstData>(ConstData{std::move(lit)});
  return e;
}

inline Expr new_app_term_expr(Pos pos, Expr fun, Expr arg) {
  Expr e;
  e.case_val = ExprCase::AppTermExpr;
  e.pos = pos;
  e.app_term_data = std::make_shared<AppTermData>(AppTermData{
      std::make_shared<Expr>(std::move(fun)),
      std::make_shared<Expr>(std::move(arg)),
  });
  return e;
}

inline Expr new_app_type_expr(Pos pos, Expr fun, ir::IrType arg) {
  Expr e;
  e.case_val = ExprCase::AppTypeExpr;
  e.pos = pos;
  e.app_type_data = std::make_shared<AppTypeData>(AppTypeData{
      std::make_shared<Expr>(std::move(fun)),
      std::move(arg),
  });
  return e;
}

inline Expr new_assign_expr(Pos pos, Expr arg, Expr ret) {
  Expr e;
  e.case_val = ExprCase::AssignExpr;
  e.pos = pos;
  e.assign_data = std::make_shared<AssignData>(AssignData{
      std::make_shared<Expr>(std::move(arg)),
      std::make_shared<Expr>(std::move(ret)),
  });
  return e;
}

inline Expr new_block_expr(Pos pos, std::vector<Expr> exprs) {
  Expr e;
  e.case_val = ExprCase::BlockExpr;
  e.pos = pos;
  e.block_data = std::make_shared<BlockData>(BlockData{std::move(exprs)});
  return e;
}

inline Expr new_let_expr(Pos pos, std::string var, std::optional<ir::IrType> var_type, Expr expr) {
  Expr e;
  e.case_val = ExprCase::LetExpr;
  e.pos = pos;
  e.let_data = std::make_shared<LetData>(LetData{
      std::move(var),
      std::move(var_type),
      std::make_shared<Expr>(std::move(expr)),
  });
  return e;
}

inline Expr new_match_expr(Pos pos, Expr expr, std::vector<MatchArm> arms) {
  Expr e;
  e.case_val = ExprCase::MatchExpr;
  e.pos = pos;
  e.match_data = std::make_shared<MatchData>(MatchData{
      std::make_shared<Expr>(std::move(expr)),
      std::move(arms),
  });
  return e;
}

inline Expr new_projection_expr(Pos pos, Expr expr, std::string label) {
  Expr e;
  e.case_val = ExprCase::ProjectionExpr;
  e.pos = pos;
  e.projection_data = std::make_shared<ProjectionData>(ProjectionData{
      std::make_shared<Expr>(std::move(expr)),
      std::move(label),
  });
  return e;
}

inline Expr new_return_expr(Pos pos, Expr expr) {
  Expr e;
  e.case_val = ExprCase::ReturnExpr;
  e.pos = pos;
  e.return_data = std::make_shared<ReturnData>(ReturnData{
      std::make_shared<Expr>(std::move(expr)),
  });
  return e;
}

inline Expr new_set_expr(Pos pos, Expr expr, std::vector<LabelValue> values) {
  Expr e;
  e.case_val = ExprCase::SetExpr;
  e.pos = pos;
  e.set_data = std::make_shared<SetData>(SetData{
      std::make_shared<Expr>(std::move(expr)),
      std::move(values),
  });
  return e;
}

inline Expr new_struct_expr(Pos pos, std::vector<LabelValue> values) {
  Expr e;
  e.case_val = ExprCase::StructExpr;
  e.pos = pos;
  e.struct_data = std::make_shared<StructData>(StructData{std::move(values)});
  return e;
}

inline Expr new_tuple_expr(Pos pos, std::vector<Expr> elems) {
  Expr e;
  e.case_val = ExprCase::TupleExpr;
  e.pos = pos;
  e.tuple_data = std::make_shared<TupleData>(TupleData{std::move(elems)});
  return e;
}

inline Expr new_lambda_expr(Pos pos, ir::FunctionArg arg, Expr body) {
  Expr e;
  e.case_val = ExprCase::LambdaExpr;
  e.pos = pos;
  e.lambda_data = std::make_shared<LambdaData>(LambdaData{
      std::move(arg),
      std::make_shared<Expr>(std::move(body)),
  });
  return e;
}

inline Expr new_for_expr(Pos pos, Expr condition, Expr body) {
  Expr e;
  e.case_val = ExprCase::ForExpr;
  e.pos = pos;
  e.for_data = std::make_shared<ForData>(ForData{
      std::make_shared<Expr>(std::move(condition)),
      std::make_shared<Expr>(std::move(body)),
  });
  return e;
}

inline Expr new_injection_expr(Pos pos, ir::IrType variant_type, std::string tag, Expr expr, std::optional<int> tag_index = std::nullopt) {
  Expr e;
  e.case_val = ExprCase::InjectionExpr;
  e.pos = pos;
  e.injection_data = std::make_shared<InjectionData>(InjectionData{
      std::move(variant_type),
      std::move(tag),
      std::make_shared<Expr>(std::move(expr)),
      tag_index,
  });
  return e;
}

inline Expr new_type_abs_expr(Pos pos, ir::TypeParam arg, Expr body) {
  Expr e;
  e.case_val = ExprCase::TypeAbsExpr;
  e.pos = pos;
  e.type_abs_data = std::make_shared<TypeAbsData>(TypeAbsData{
      std::move(arg),
      std::make_shared<Expr>(std::move(body)),
  });
  return e;
}

} // namespace ast
