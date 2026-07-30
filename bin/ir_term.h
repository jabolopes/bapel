#pragma once

#include "ir_type.h"

namespace ir {

struct FunctionArg {
  std::string id;
  IrType type;

  std::string to_string() const {
    return id + ": " + type.to_string();
  }
};

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

} // namespace ir
