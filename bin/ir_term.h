#pragma once

#include "ir_type.h"

namespace ir {

struct MatchArm {
  std::string tag;
  std::string arg;
  std::shared_ptr<IrTerm> body;
  std::optional<int> index;

  bool is_wildcard() const { return tag == "_"; }
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

  BlockTermData();
  explicit BlockTermData(std::vector<IrTerm> t);
  BlockTermData(const BlockTermData&);
  BlockTermData(BlockTermData&&) noexcept;
  BlockTermData& operator=(const BlockTermData&);
  BlockTermData& operator=(BlockTermData&&) noexcept;
  ~BlockTermData();
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

  TupleTermData();
  explicit TupleTermData(std::vector<IrTerm> e);
  TupleTermData(const TupleTermData&);
  TupleTermData(TupleTermData&&) noexcept;
  TupleTermData& operator=(const TupleTermData&);
  TupleTermData& operator=(TupleTermData&&) noexcept;
  ~TupleTermData();
};

struct TypeAbsTermData {
  TypeParam type_param;
  std::shared_ptr<IrTerm> body;
};

struct VarTermData {
  std::string id;
};

BAPEL_VARIANT(IrTermData,
  (app_term, (AppTermData)),
  (app_type, (AppTypeTermData)),
  (assign, (AssignTermData)),
  (block, (BlockTermData)),
  (const_t, (ConstTermData)),
  (injection, (InjectionTermData)),
  (lambda, (LambdaTermData)),
  (let_t, (LetTermData)),
  (match_t, (MatchTermData)),
  (projection, (ProjectionTermData)),
  (return_t, (ReturnTermData)),
  (set_t, (SetTermData)),
  (struct_t, (StructTermData)),
  (tuple_t, (TupleTermData)),
  (type_abs, (TypeAbsTermData)),
  (var_t, (VarTermData))
);

struct IrTerm {
  IrTermData data{IrTermData::mk_var_t(VarTermData{})};
  std::optional<IrType> type;
  Pos pos;

  IrTermData::Tag tag() const { return data.tag(); }
  bool is(IrTermData::Tag t) const { return data.is(t); }

  std::pair<IrTerm, std::vector<IrType>> app_types() const {
    std::vector<IrType> types;
    IrTerm cur = *this;
    while (cur.is(IrTermData::Tag::app_type)) {
      const auto& at = cur.data.as_app_type();
      types.insert(types.begin(), at.arg);
      if (at.fun) {
        cur = *at.fun;
      } else {
        break;
      }
    }
    return {cur, types};
  }

  std::tuple<IrTerm, std::vector<IrType>, IrTerm> app_args() const {
    if (!is(IrTermData::Tag::app_term)) {
      return {*this, {}, *this};
    }
    const auto& at = data.as_app_term();
    IrTerm fun_term = at.fun ? *at.fun : IrTerm{};
    IrTerm arg_term = at.arg ? *at.arg : IrTerm{};
    auto [base_fun, types] = fun_term.app_types();
    return {base_fun, types, arg_term};
  }

  std::tuple<std::vector<TypeParam>, std::vector<FunctionArg>, IrTerm> to_function() const {
    std::vector<TypeParam> tvars;
    std::vector<FunctionArg> args;
    IrTerm cur = *this;

    while (cur.is(IrTermData::Tag::type_abs)) {
      const auto& ta = cur.data.as_type_abs();
      tvars.push_back(ta.type_param);
      if (ta.body) {
        cur = *ta.body;
      } else {
        break;
      }
    }

    while (cur.is(IrTermData::Tag::lambda)) {
      const auto& l = cur.data.as_lambda();
      args.push_back(l.arg);
      if (l.body) {
        cur = *l.body;
      } else {
        break;
      }
    }

    return {tvars, args, cur};
  }

  std::string to_string() const {
    std::string inner;
    switch (data.tag()) {
      case IrTermData::Tag::app_term: {
        const auto& at = data.as_app_term();
        std::string fun_l = (at.fun && at.fun->is(IrTermData::Tag::lambda)) ? "(" : "";
        std::string fun_r = (at.fun && at.fun->is(IrTermData::Tag::lambda)) ? ")" : "";
        std::string arg_l = "";
        std::string arg_r = "";
        if (at.arg &&
            (at.arg->is(IrTermData::Tag::app_term) || at.arg->is(IrTermData::Tag::lambda))) {
          arg_l = "(";
          arg_r = ")";
        }
        inner = fun_l + (at.fun ? at.fun->to_string() : "") + fun_r + " " +
                arg_l + (at.arg ? at.arg->to_string() : "") + arg_r;
        break;
      }
      case IrTermData::Tag::app_type: {
        const auto& at = data.as_app_type();
        inner = (at.fun ? at.fun->to_string() : "") + " [" + at.arg.to_string() + "]";
        break;
      }
      case IrTermData::Tag::assign: {
        const auto& a = data.as_assign();
        inner = (a.ret ? a.ret->to_string() : "") + " <- " + (a.arg ? a.arg->to_string() : "");
        break;
      }
      case IrTermData::Tag::block: {
        std::string s = "{\n";
        for (const auto& term : data.as_block().terms) {
          std::string ts = term.to_string();
          size_t start = 0;
          while (start < ts.size()) {
            size_t end = ts.find('\n', start);
            if (end == std::string::npos) {
              s += "  " + ts.substr(start) + "\n";
              break;
            } else {
              s += "  " + ts.substr(start, end - start + 1);
              start = end + 1;
            }
          }
        }
        s += "}";
        inner = s;
        break;
      }
      case IrTermData::Tag::const_t:
        inner = data.as_const_t().literal.to_string();
        break;
      case IrTermData::Tag::injection: {
        const auto& inj = data.as_injection();
        bool type_needs_parens = false;
        switch (inj.variant_type.tag()) {
          case IrTypeData::Tag::app_t:
          case IrTypeData::Tag::forall_t:
          case IrTypeData::Tag::fun_t:
          case IrTypeData::Tag::lambda_t:
            type_needs_parens = true;
            break;
          default:
            break;
        }
        std::string lp = type_needs_parens ? "(" : "";
        std::string rp = type_needs_parens ? ")" : "";
        inner = "variant{" + lp + inj.variant_type.to_string() + rp + " " +
                inj.tag + " = " + (inj.value ? inj.value->to_string() : "") + "}";
        break;
      }
      case IrTermData::Tag::lambda: {
        const auto& l = data.as_lambda();
        inner = "\\(" + l.arg.id + ": " + l.arg.type.to_string() + ") -> " +
                (l.body ? l.body->to_string() : "");
        break;
      }
      case IrTermData::Tag::let_t: {
        const auto& lt = data.as_let_t();
        if (lt.var_type.has_value()) {
          inner = "let " + lt.var + ": " + lt.var_type->to_string() + " = " +
                  (lt.value ? lt.value->to_string() : "");
        } else {
          inner = "let " + lt.var + " = " + (lt.value ? lt.value->to_string() : "");
        }
        break;
      }
      case IrTermData::Tag::match_t: {
        const auto& m = data.as_match_t();
        auto format_arm = [](const MatchArm& arm) {
          if (arm.is_wildcard()) {
            return "_ -> " + (arm.body ? arm.body->to_string() : "");
          }
          return arm.tag + " " + arm.arg + " -> " + (arm.body ? arm.body->to_string() : "");
        };
        std::string s = "case " + (m.term ? m.term->to_string() : "") + " {";
        if (m.arms.size() == 1) {
          s += " " + format_arm(m.arms[0]) + " }";
        } else if (m.arms.size() > 1) {
          s += "\n";
          for (size_t i = 0; i < m.arms.size(); ++i) {
            if (i > 0) s += "\n";
            s += "    " + format_arm(m.arms[i]);
          }
          s += "}";
        } else {
          s += "}";
        }
        inner = s;
        break;
      }
      case IrTermData::Tag::projection: {
        const auto& p = data.as_projection();
        inner = (p.term ? p.term->to_string() : "") + "." + p.label;
        break;
      }
      case IrTermData::Tag::return_t: {
        const auto& r = data.as_return_t();
        inner = "return " + (r.expr ? r.expr->to_string() : "");
        break;
      }
      case IrTermData::Tag::set_t: {
        const auto& st = data.as_set_t();
        std::string s = "set " + (st.term ? st.term->to_string() : "") + " {";
        for (size_t i = 0; i < st.values.size(); ++i) {
          if (i > 0) s += ", ";
          s += st.values[i].label + " = " +
               (st.values[i].value ? st.values[i].value->to_string() : "");
        }
        s += "}";
        inner = s;
        break;
      }
      case IrTermData::Tag::struct_t: {
        const auto& st = data.as_struct_t();
        std::string s = "struct{";
        for (size_t i = 0; i < st.values.size(); ++i) {
          if (i > 0) s += ", ";
          s += st.values[i].label + " = " +
               (st.values[i].value ? st.values[i].value->to_string() : "");
        }
        s += "}";
        inner = s;
        break;
      }
      case IrTermData::Tag::tuple_t: {
        const auto& tt = data.as_tuple_t();
        std::string s = "(";
        for (size_t i = 0; i < tt.elems.size(); ++i) {
          if (i > 0) s += ", ";
          s += tt.elems[i].to_string();
        }
        s += ")";
        inner = s;
        break;
      }
      case IrTermData::Tag::type_abs: {
        const auto& ta = data.as_type_abs();
        inner = "Λ" + ta.type_param.to_string() + ". " +
                (ta.body ? ta.body->to_string() : "");
        break;
      }
      case IrTermData::Tag::var_t:
        inner = data.as_var_t().id;
        break;
    }

    if (!type.has_value()) {
      return inner;
    }

    bool term_needs_parens = false;
    switch (data.tag()) {
      case IrTermData::Tag::app_term:
      case IrTermData::Tag::app_type:
      case IrTermData::Tag::assign:
      case IrTermData::Tag::injection:
      case IrTermData::Tag::lambda:
      case IrTermData::Tag::let_t:
      case IrTermData::Tag::match_t:
      case IrTermData::Tag::projection:
      case IrTermData::Tag::return_t:
      case IrTermData::Tag::type_abs:
        term_needs_parens = true;
        break;
      default:
        break;
    }

    bool type_needs_parens = false;
    switch (type->tag()) {
      case IrTypeData::Tag::app_t:
      case IrTypeData::Tag::forall_t:
      case IrTypeData::Tag::fun_t:
      case IrTypeData::Tag::lambda_t:
        type_needs_parens = true;
        break;
      default:
        break;
    }

    std::string res;
    if (term_needs_parens) res += "(";
    res += inner;
    if (term_needs_parens) res += ")";
    res += ":";
    if (type_needs_parens) res += "(";
    res += type->to_string();
    if (type_needs_parens) res += ")";
    return res;
  }
};

inline BlockTermData::BlockTermData() = default;
inline BlockTermData::BlockTermData(std::vector<IrTerm> t) : terms(std::move(t)) {}
inline BlockTermData::BlockTermData(const BlockTermData&) = default;
inline BlockTermData::BlockTermData(BlockTermData&&) noexcept = default;
inline BlockTermData& BlockTermData::operator=(const BlockTermData&) = default;
inline BlockTermData& BlockTermData::operator=(BlockTermData&&) noexcept = default;
inline BlockTermData::~BlockTermData() = default;

inline TupleTermData::TupleTermData() = default;
inline TupleTermData::TupleTermData(std::vector<IrTerm> e) : elems(std::move(e)) {}
inline TupleTermData::TupleTermData(const TupleTermData&) = default;
inline TupleTermData::TupleTermData(TupleTermData&&) noexcept = default;
inline TupleTermData& TupleTermData::operator=(const TupleTermData&) = default;
inline TupleTermData& TupleTermData::operator=(TupleTermData&&) noexcept = default;
inline TupleTermData::~TupleTermData() = default;

inline IrTerm new_var_term(const std::string& id) {
  return IrTerm{IrTermData::mk_var_t(VarTermData{id}), std::nullopt, Pos{}};
}

inline IrTerm new_const_term(IrLiteral literal) {
  return IrTerm{IrTermData::mk_const_t(ConstTermData{std::move(literal)}), std::nullopt, Pos{}};
}

inline IrTerm new_let_term(std::string var, std::optional<IrType> var_type, IrTerm value) {
  return IrTerm{
      IrTermData::mk_let_t(LetTermData{
          std::move(var),
          std::move(var_type),
          std::make_shared<IrTerm>(std::move(value)),
      }),
      std::nullopt,
      Pos{},
  };
}

inline IrTerm new_assign_term(IrTerm ret, IrTerm arg) {
  return IrTerm{
      IrTermData::mk_assign(AssignTermData{
          std::make_shared<IrTerm>(std::move(arg)),
          std::make_shared<IrTerm>(std::move(ret)),
      }),
      std::nullopt,
      Pos{},
  };
}

inline IrTerm new_block_term(std::vector<IrTerm> terms) {
  return IrTerm{IrTermData::mk_block(BlockTermData{std::move(terms)}), std::nullopt, Pos{}};
}

inline IrTerm new_return_term(IrTerm expr) {
  return IrTerm{
      IrTermData::mk_return_t(ReturnTermData{std::make_shared<IrTerm>(std::move(expr))}),
      std::nullopt,
      Pos{},
  };
}

inline IrTerm new_tuple_term(std::vector<IrTerm> elems) {
  if (elems.size() == 1) {
    return elems[0];
  }
  return IrTerm{IrTermData::mk_tuple_t(TupleTermData{std::move(elems)}), std::nullopt, Pos{}};
}

inline IrTerm new_struct_term(std::vector<LabelValue> values) {
  return IrTerm{IrTermData::mk_struct_t(StructTermData{std::move(values)}), std::nullopt, Pos{}};
}

inline IrTerm new_injection_term(IrType vtype, std::string tag, IrTerm value, std::optional<int> tag_idx = std::nullopt) {
  return IrTerm{
      IrTermData::mk_injection(InjectionTermData{
          std::move(vtype),
          std::move(tag),
          std::make_shared<IrTerm>(std::move(value)),
          tag_idx,
      }),
      std::nullopt,
      Pos{},
  };
}

inline IrTerm new_app_term(IrTerm fun, IrTerm arg) {
  return IrTerm{
      IrTermData::mk_app_term(AppTermData{
          std::make_shared<IrTerm>(std::move(fun)),
          std::make_shared<IrTerm>(std::move(arg)),
      }),
      std::nullopt,
      Pos{},
  };
}

inline IrTerm new_app_type_term(IrTerm fun, IrType arg) {
  return IrTerm{
      IrTermData::mk_app_type(AppTypeTermData{
          std::make_shared<IrTerm>(std::move(fun)),
          std::move(arg),
      }),
      std::nullopt,
      Pos{},
  };
}

inline IrTerm new_projection_term(IrTerm term, std::string label) {
  return IrTerm{
      IrTermData::mk_projection(ProjectionTermData{
          std::make_shared<IrTerm>(std::move(term)),
          std::move(label),
          std::nullopt,
      }),
      std::nullopt,
      Pos{},
  };
}

inline IrTerm new_lambda_term(FunctionArg arg, IrTerm body) {
  return IrTerm{
      IrTermData::mk_lambda(LambdaTermData{
          std::move(arg),
          std::make_shared<IrTerm>(std::move(body)),
      }),
      std::nullopt,
      Pos{},
  };
}

inline IrTerm new_match_term(IrTerm term, std::vector<MatchArm> arms) {
  return IrTerm{
      IrTermData::mk_match_t(MatchTermData{
          std::make_shared<IrTerm>(std::move(term)),
          std::move(arms),
      }),
      std::nullopt,
      Pos{},
  };
}

inline IrTerm new_set_term(IrTerm term, std::vector<LabelValue> values) {
  return IrTerm{
      IrTermData::mk_set_t(SetTermData{
          std::make_shared<IrTerm>(std::move(term)),
          std::move(values),
          std::nullopt,
      }),
      std::nullopt,
      Pos{},
  };
}

inline IrTerm new_type_abs_term(TypeParam tp, IrTerm body) {
  return IrTerm{
      IrTermData::mk_type_abs(TypeAbsTermData{
          std::move(tp),
          std::make_shared<IrTerm>(std::move(body)),
      }),
      std::nullopt,
      Pos{},
  };
}

} // namespace ir
