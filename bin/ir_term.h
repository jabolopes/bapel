#pragma once

#include "ir_type.h"

namespace ir {

struct MatchArm {
  std::string tag;
  std::string arg;
  std::shared_ptr<IrTerm> body;
  std::optional<int> index;
};

struct LabelValue {
  std::string label;
  std::shared_ptr<IrTerm> value;
};

struct AppTermData {
  std::shared_ptr<IrTerm> fun;
  std::shared_ptr<IrTerm> arg;
};

struct AppTypeTermData {
  std::shared_ptr<IrTerm> fun;
  IrType arg;
};

struct AssignTermData {
  std::shared_ptr<IrTerm> arg;
  std::shared_ptr<IrTerm> ret;
};

struct BlockTermData {
  std::vector<IrTerm> terms;
};

struct ConstTermData {
  IrLiteral literal;
};

struct InjectionTermData {
  IrType variant_type;
  std::string tag;
  std::shared_ptr<IrTerm> value;
  std::optional<int> tag_index;
};

struct LambdaTermData {
  FunctionArg arg;
  std::shared_ptr<IrTerm> body;
};

struct LetTermData {
  std::string var;
  std::optional<IrType> var_type;
  std::shared_ptr<IrTerm> value;
};

struct MatchTermData {
  std::shared_ptr<IrTerm> term;
  std::vector<MatchArm> arms;
};

struct ProjectionTermData {
  std::shared_ptr<IrTerm> term;
  std::string label;
  std::optional<IrType> reduced_type;
};

struct ReturnTermData {
  std::shared_ptr<IrTerm> expr;
};

struct SetTermData {
  std::shared_ptr<IrTerm> term;
  std::vector<LabelValue> values;
  std::optional<IrType> reduced_type;
};

struct StructTermData {
  std::vector<LabelValue> values;
};

struct TupleTermData {
  std::vector<IrTerm> elems;
};

struct TypeAbsTermData {
  TypeParam type_param;
  std::shared_ptr<IrTerm> body;
};

struct VarTermData {
  std::string id;
};

struct IrTerm {
  IrTermCase case_val = IrTermCase::VarTerm;
  std::shared_ptr<AppTermData> app_term;
  std::shared_ptr<AppTypeTermData> app_type;
  std::shared_ptr<AssignTermData> assign;
  std::shared_ptr<BlockTermData> block;
  std::shared_ptr<ConstTermData> const_data;
  std::shared_ptr<InjectionTermData> injection;
  std::shared_ptr<LambdaTermData> lambda;
  std::shared_ptr<LetTermData> let_data;
  std::shared_ptr<MatchTermData> match_data;
  std::shared_ptr<ProjectionTermData> projection;
  std::shared_ptr<ReturnTermData> return_data;
  std::shared_ptr<SetTermData> set_data;
  std::shared_ptr<StructTermData> struct_data;
  std::shared_ptr<TupleTermData> tuple_data;
  std::shared_ptr<TypeAbsTermData> type_abs;
  std::shared_ptr<VarTermData> var_data;

  std::optional<IrType> type;
  Pos pos;

  bool is(IrTermCase c) const { return case_val == c; }

  std::pair<IrTerm, std::vector<IrType>> app_types() const {
    std::vector<IrType> types;
    IrTerm cur = *this;
    while (cur.is(IrTermCase::AppTypeTerm) && cur.app_type) {
      types.insert(types.begin(), cur.app_type->arg);
      if (cur.app_type->fun) {
        cur = *cur.app_type->fun;
      } else {
        break;
      }
    }
    return {cur, types};
  }

  std::tuple<IrTerm, std::vector<IrType>, IrTerm> app_args() const {
    if (!is(IrTermCase::AppTermTerm) || !app_term) {
      return {*this, {}, *this};
    }
    IrTerm fun_term = app_term->fun ? *app_term->fun : IrTerm{};
    IrTerm arg_term = app_term->arg ? *app_term->arg : IrTerm{};
    auto [base_fun, types] = fun_term.app_types();
    return {base_fun, types, arg_term};
  }

  std::tuple<std::vector<TypeParam>, std::vector<FunctionArg>, IrTerm> to_function() const {
    std::vector<TypeParam> tvars;
    std::vector<FunctionArg> args;
    IrTerm cur = *this;

    while (cur.is(IrTermCase::TypeAbsTerm) && cur.type_abs) {
      tvars.push_back(cur.type_abs->type_param);
      if (cur.type_abs->body) {
        cur = *cur.type_abs->body;
      } else {
        break;
      }
    }

    while (cur.is(IrTermCase::LambdaTerm) && cur.lambda) {
      args.push_back(cur.lambda->arg);
      if (cur.lambda->body) {
        cur = *cur.lambda->body;
      } else {
        break;
      }
    }

    return {tvars, args, cur};
  }

  std::string to_string() const {
    switch (case_val) {
      case IrTermCase::AppTermTerm:
        return "(" + (app_term && app_term->fun ? app_term->fun->to_string() : "") + " " +
               (app_term && app_term->arg ? app_term->arg->to_string() : "") + ")";
      case IrTermCase::AppTypeTerm:
        return "(" + (app_type && app_type->fun ? app_type->fun->to_string() : "") + " [" +
               (app_type ? app_type->arg.to_string() : "") + "])";
      case IrTermCase::AssignTerm:
        return (assign && assign->ret ? assign->ret->to_string() : "") + " = " +
               (assign && assign->arg ? assign->arg->to_string() : "");
      case IrTermCase::BlockTerm: {
        std::string s = "{ ";
        if (block) {
          for (size_t i = 0; i < block->terms.size(); ++i) {
            if (i > 0) s += "; ";
            s += block->terms[i].to_string();
          }
        }
        s += " }";
        return s;
      }
      case IrTermCase::ConstTerm:
        return const_data ? const_data->literal.to_string() : "const";
      case IrTermCase::InjectionTerm:
        return (injection ? injection->tag : "") + " " + (injection && injection->value ? injection->value->to_string() : "");
      case IrTermCase::LambdaTerm:
        return "(\\" + (lambda ? lambda->arg.to_string() : "") + " -> " +
               (lambda && lambda->body ? lambda->body->to_string() : "") + ")";
      case IrTermCase::LetTerm:
        return "let " + (let_data ? let_data->var : "") + " = " +
               (let_data && let_data->value ? let_data->value->to_string() : "");
      case IrTermCase::MatchTerm:
        return "match " + (match_data && match_data->term ? match_data->term->to_string() : "");
      case IrTermCase::ProjectionTerm:
        return (projection && projection->term ? projection->term->to_string() : "") + "." +
               (projection ? projection->label : "");
      case IrTermCase::ReturnTerm:
        return "return " + (return_data && return_data->expr ? return_data->expr->to_string() : "");
      case IrTermCase::SetTerm:
        return "set " + (set_data && set_data->term ? set_data->term->to_string() : "");
      case IrTermCase::StructTerm:
        return "struct";
      case IrTermCase::TupleTerm: {
        std::string s = "(";
        if (tuple_data) {
          for (size_t i = 0; i < tuple_data->elems.size(); ++i) {
            if (i > 0) s += ", ";
            s += tuple_data->elems[i].to_string();
          }
        }
        s += ")";
        return s;
      }
      case IrTermCase::TypeAbsTerm:
        return "/\\" + (type_abs ? type_abs->type_param.to_string() : "") + " -> " +
               (type_abs && type_abs->body ? type_abs->body->to_string() : "");
      case IrTermCase::VarTerm:
        return var_data ? var_data->id : "";
    }
    return "";
  }

  std::string to_json() const {
    std::stringstream ss;
    ss << "{\"Case\":" << static_cast<int>(case_val);
    if (type.has_value()) {
      ss << ",\"Type\":" << type->to_json();
    }
    ss << ",\"Pos\":" << pos.to_json();
    switch (case_val) {
      case IrTermCase::AppTermTerm:
        if (app_term) {
          ss << ",\"AppTerm\":{\"Fun\":" << (app_term->fun ? app_term->fun->to_json() : "null")
             << ",\"Arg\":" << (app_term->arg ? app_term->arg->to_json() : "null") << "}";
        }
        break;
      case IrTermCase::AppTypeTerm:
        if (app_type) {
          ss << ",\"AppType\":{\"Fun\":" << (app_type->fun ? app_type->fun->to_json() : "null")
             << ",\"Arg\":" << app_type->arg.to_json() << "}";
        }
        break;
      case IrTermCase::AssignTerm:
        if (assign) {
          ss << ",\"Assign\":{\"Arg\":" << (assign->arg ? assign->arg->to_json() : "null")
             << ",\"Ret\":" << (assign->ret ? assign->ret->to_json() : "null") << "}";
        }
        break;
      case IrTermCase::BlockTerm:
        if (block) {
          ss << ",\"Block\":{\"Terms\":[";
          Interleave(block->terms, [&]() { ss << ","; }, [&](int, const IrTerm& t) {
            ss << t.to_json();
          });
          ss << "]}";
        }
        break;
      case IrTermCase::ConstTerm:
        if (const_data) {
          ss << ",\"Const\":" << const_data->literal.to_json();
        }
        break;
      case IrTermCase::InjectionTerm:
        if (injection) {
          ss << ",\"Injection\":{\"VariantType\":" << injection->variant_type.to_json()
             << ",\"Tag\":\"" << json_escape(injection->tag) << "\""
             << ",\"Value\":" << (injection->value ? injection->value->to_json() : "null");
          if (injection->tag_index.has_value()) {
            ss << ",\"TagIndex\":" << injection->tag_index.value();
          }
          ss << "}";
        }
        break;
      case IrTermCase::LambdaTerm:
        if (lambda) {
          ss << ",\"Lambda\":{\"Arg\":" << lambda->arg.to_json()
             << ",\"Body\":" << (lambda->body ? lambda->body->to_json() : "null") << "}";
        }
        break;
      case IrTermCase::LetTerm:
        if (let_data) {
          ss << ",\"Let\":{\"Var\":\"" << json_escape(let_data->var) << "\"";
          if (let_data->var_type.has_value()) {
            ss << ",\"VarType\":" << let_data->var_type->to_json();
          }
          ss << ",\"Value\":" << (let_data->value ? let_data->value->to_json() : "null") << "}";
        }
        break;
      case IrTermCase::MatchTerm:
        if (match_data) {
          ss << ",\"Match\":{\"Term\":" << (match_data->term ? match_data->term->to_json() : "null")
             << ",\"Arms\":[";
          Interleave(match_data->arms, [&]() { ss << ","; }, [&](int, const MatchArm& arm) {
            ss << "{\"Tag\":\"" << json_escape(arm.tag) << "\",\"Arg\":\"" << json_escape(arm.arg) << "\""
               << ",\"Body\":" << (arm.body ? arm.body->to_json() : "null");
            if (arm.index.has_value()) {
              ss << ",\"Index\":" << arm.index.value();
            }
            ss << "}";
          });
          ss << "]}";
        }
        break;
      case IrTermCase::ProjectionTerm:
        if (projection) {
          ss << ",\"Projection\":{\"Term\":" << (projection->term ? projection->term->to_json() : "null")
             << ",\"Label\":\"" << json_escape(projection->label) << "\"";
          if (projection->reduced_type.has_value()) {
            ss << ",\"ReducedType\":" << projection->reduced_type->to_json();
          }
          ss << "}";
        }
        break;
      case IrTermCase::ReturnTerm:
        if (return_data) {
          ss << ",\"Return\":{\"Expr\":" << (return_data->expr ? return_data->expr->to_json() : "null") << "}";
        }
        break;
      case IrTermCase::SetTerm:
        if (set_data) {
          ss << ",\"Set\":{\"Term\":" << (set_data->term ? set_data->term->to_json() : "null")
             << ",\"Values\":[";
          Interleave(set_data->values, [&]() { ss << ","; }, [&](int, const LabelValue& lv) {
            ss << "{\"Label\":\"" << json_escape(lv.label) << "\",\"Value\":" << (lv.value ? lv.value->to_json() : "null") << "}";
          });
          ss << "]";
          if (set_data->reduced_type.has_value()) {
            ss << ",\"ReducedType\":" << set_data->reduced_type->to_json();
          }
          ss << "}";
        }
        break;
      case IrTermCase::StructTerm:
        if (struct_data) {
          ss << ",\"Struct\":{\"Values\":[";
          Interleave(struct_data->values, [&]() { ss << ","; }, [&](int, const LabelValue& lv) {
            ss << "{\"Label\":\"" << json_escape(lv.label) << "\",\"Value\":" << (lv.value ? lv.value->to_json() : "null") << "}";
          });
          ss << "]}";
        }
        break;
      case IrTermCase::TupleTerm:
        if (tuple_data) {
          ss << ",\"Tuple\":{\"Elems\":[";
          Interleave(tuple_data->elems, [&]() { ss << ","; }, [&](int, const IrTerm& elem) {
            ss << elem.to_json();
          });
          ss << "]}";
        }
        break;
      case IrTermCase::TypeAbsTerm:
        if (type_abs) {
          ss << ",\"TypeAbs\":{\"TypeParam\":" << type_abs->type_param.to_json()
             << ",\"Body\":" << (type_abs->body ? type_abs->body->to_json() : "null") << "}";
        }
        break;
      case IrTermCase::VarTerm:
        if (var_data) {
          ss << ",\"Var\":{\"ID\":\"" << json_escape(var_data->id) << "\"}";
        }
        break;
    }
    ss << "}";
    return ss.str();
  }
};

inline IrTerm new_var_term(const std::string& id) {
  IrTerm t;
  t.case_val = IrTermCase::VarTerm;
  t.var_data = std::make_shared<VarTermData>();
  t.var_data->id = id;
  return t;
}

inline IrTerm new_const_int_term(int64_t v) {
  IrTerm t;
  t.case_val = IrTermCase::ConstTerm;
  t.const_data = std::make_shared<ConstTermData>();
  t.const_data->literal.int_val = v;
  return t;
}

inline IrTerm new_const_float_term(int64_t integer, int64_t decimal) {
  IrTerm t;
  t.case_val = IrTermCase::ConstTerm;
  t.const_data = std::make_shared<ConstTermData>();
  t.const_data->literal.float_val = FloatLit{integer, decimal};
  return t;
}

inline IrTerm new_const_str_term(const std::string& s) {
  IrTerm t;
  t.case_val = IrTermCase::ConstTerm;
  t.const_data = std::make_shared<ConstTermData>();
  t.const_data->literal.str_val = s;
  return t;
}

inline IrTerm new_const_rune_term(const std::string& r) {
  IrTerm t;
  t.case_val = IrTermCase::ConstTerm;
  t.const_data = std::make_shared<ConstTermData>();
  t.const_data->literal.rune_val = r;
  return t;
}

inline IrTerm new_let_term(std::string var, std::optional<IrType> var_type, IrTerm value) {
  IrTerm t;
  t.case_val = IrTermCase::LetTerm;
  t.let_data = std::make_shared<LetTermData>();
  t.let_data->var = std::move(var);
  t.let_data->var_type = std::move(var_type);
  t.let_data->value = std::make_shared<IrTerm>(std::move(value));
  return t;
}

inline IrTerm new_assign_term(IrTerm ret, IrTerm arg) {
  IrTerm t;
  t.case_val = IrTermCase::AssignTerm;
  t.assign = std::make_shared<AssignTermData>();
  t.assign->ret = std::make_shared<IrTerm>(std::move(ret));
  t.assign->arg = std::make_shared<IrTerm>(std::move(arg));
  return t;
}

inline IrTerm new_block_term(std::vector<IrTerm> terms) {
  IrTerm t;
  t.case_val = IrTermCase::BlockTerm;
  t.block = std::make_shared<BlockTermData>();
  t.block->terms = std::move(terms);
  return t;
}

inline IrTerm new_return_term(IrTerm expr) {
  IrTerm t;
  t.case_val = IrTermCase::ReturnTerm;
  t.return_data = std::make_shared<ReturnTermData>();
  t.return_data->expr = std::make_shared<IrTerm>(std::move(expr));
  return t;
}

inline IrTerm new_tuple_term(std::vector<IrTerm> elems) {
  if (elems.size() == 1) {
    return elems[0];
  }
  IrTerm t;
  t.case_val = IrTermCase::TupleTerm;
  t.tuple_data = std::make_shared<TupleTermData>();
  t.tuple_data->elems = std::move(elems);
  return t;
}

inline IrTerm new_struct_term(std::vector<LabelValue> values) {
  IrTerm t;
  t.case_val = IrTermCase::StructTerm;
  t.struct_data = std::make_shared<StructTermData>();
  t.struct_data->values = std::move(values);
  return t;
}

inline IrTerm new_injection_term(IrType vtype, std::string tag, IrTerm value, std::optional<int> tag_idx = std::nullopt) {
  IrTerm t;
  t.case_val = IrTermCase::InjectionTerm;
  t.injection = std::make_shared<InjectionTermData>();
  t.injection->variant_type = std::move(vtype);
  t.injection->tag = std::move(tag);
  t.injection->value = std::make_shared<IrTerm>(std::move(value));
  t.injection->tag_index = tag_idx;
  return t;
}

inline IrTerm new_app_term(IrTerm fun, IrTerm arg) {
  IrTerm t;
  t.case_val = IrTermCase::AppTermTerm;
  t.app_term = std::make_shared<AppTermData>();
  t.app_term->fun = std::make_shared<IrTerm>(std::move(fun));
  t.app_term->arg = std::make_shared<IrTerm>(std::move(arg));
  return t;
}

inline IrTerm new_app_type_term(IrTerm fun, IrType arg) {
  IrTerm t;
  t.case_val = IrTermCase::AppTypeTerm;
  t.app_type = std::make_shared<AppTypeTermData>();
  t.app_type->fun = std::make_shared<IrTerm>(std::move(fun));
  t.app_type->arg = std::move(arg);
  return t;
}

inline IrTerm new_projection_term(IrTerm term, std::string label) {
  IrTerm t;
  t.case_val = IrTermCase::ProjectionTerm;
  t.projection = std::make_shared<ProjectionTermData>();
  t.projection->term = std::make_shared<IrTerm>(std::move(term));
  t.projection->label = std::move(label);
  return t;
}

} // namespace ir
