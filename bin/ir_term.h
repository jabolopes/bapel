#pragma once

#include "ir_type.h"

namespace ir {

using IrTermData = ::IrTermData<::IrTerm>;
using MatchArm = ::MatchArm<::IrTerm>;
using LabelValue = ::LabelValue<::IrTerm>;

inline bool is_wildcard_arm(const MatchArm& arm) {
  return arm.tag == "_";
}

inline std::pair<IrTerm, std::vector<IrType>> term_app_types(const IrTerm& term) {
  std::vector<IrType> types;
  IrTerm cur = term;
  while (cur.data.is(IrTermData::Tag::app_type)) {
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

inline std::tuple<IrTerm, std::vector<IrType>, IrTerm> term_app_args(const IrTerm& term) {
  if (!term.data.is(IrTermData::Tag::app_term)) {
    return {term, {}, term};
  }
  const auto& at = term.data.as_app_term();
  IrTerm fun_term = at.fun ? *at.fun : IrTerm{};
  IrTerm arg_term = at.arg ? *at.arg : IrTerm{};
  auto [base_fun, types] = term_app_types(fun_term);
  return {base_fun, types, arg_term};
}

inline std::tuple<std::vector<TypeParam>, std::vector<FunctionArg>, IrTerm> term_to_function(const IrTerm& term) {
  std::vector<TypeParam> tvars;
  std::vector<FunctionArg> args;
  IrTerm cur = term;

  while (cur.data.is(IrTermData::Tag::type_abs)) {
    const auto& ta = cur.data.as_type_abs();
    tvars.push_back(ta.param);
    if (ta.body) {
      cur = *ta.body;
    } else {
      break;
    }
  }

  while (cur.data.is(IrTermData::Tag::lambda)) {
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

inline std::string term_to_string(const IrTerm& term) {
  std::string inner;
  switch (term.data.tag()) {
    case IrTermData::Tag::app_term: {
      const auto& at = term.data.as_app_term();
      std::string fun_l = (at.fun && at.fun->data.is(IrTermData::Tag::lambda)) ? "(" : "";
      std::string fun_r = (at.fun && at.fun->data.is(IrTermData::Tag::lambda)) ? ")" : "";
      std::string arg_l = "";
      std::string arg_r = "";
      if (at.arg &&
          (at.arg->data.is(IrTermData::Tag::app_term) || at.arg->data.is(IrTermData::Tag::lambda))) {
        arg_l = "(";
        arg_r = ")";
      }
      inner = fun_l + (at.fun ? term_to_string(*at.fun) : "") + fun_r + " " +
              arg_l + (at.arg ? term_to_string(*at.arg) : "") + arg_r;
      break;
    }
    case IrTermData::Tag::app_type: {
      const auto& at = term.data.as_app_type();
      inner = (at.fun ? term_to_string(*at.fun) : "") + " [" + type_to_string(at.arg) + "]";
      break;
    }
    case IrTermData::Tag::assign: {
      const auto& a = term.data.as_assign();
      inner = (a.ret ? term_to_string(*a.ret) : "") + " <- " + (a.arg ? term_to_string(*a.arg) : "");
      break;
    }
    case IrTermData::Tag::block: {
      std::string s = "{\n";
      for (const auto& t : term.data.as_block()) {
        std::string ts = term_to_string(t);
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
      inner = literal_to_string(term.data.as_const_t());
      break;
    case IrTermData::Tag::injection: {
      const auto& inj = term.data.as_injection();
      bool type_needs_parens = false;
      switch (inj.variant_type.data.tag()) {
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
      inner = "variant{" + lp + type_to_string(inj.variant_type) + rp + " " +
              inj.tag + " = " + (inj.value ? term_to_string(*inj.value) : "") + "}";
      break;
    }
    case IrTermData::Tag::lambda: {
      const auto& l = term.data.as_lambda();
      inner = "\\(" + l.arg.id + ": " + type_to_string(l.arg.typ) + ") -> " +
              (l.body ? term_to_string(*l.body) : "");
      break;
    }
    case IrTermData::Tag::let_t: {
      const auto& lt = term.data.as_let_t();
      if (lt.var_type.has_value()) {
        inner = "let " + lt.var + ": " + type_to_string(*lt.var_type) + " = " +
                (lt.value ? term_to_string(*lt.value) : "");
      } else {
        inner = "let " + lt.var + " = " + (lt.value ? term_to_string(*lt.value) : "");
      }
      break;
    }
    case IrTermData::Tag::match_t: {
      const auto& m = term.data.as_match_t();
      auto format_arm = [](const MatchArm& arm) {
        if (is_wildcard_arm(arm)) {
          return "_ -> " + (arm.body ? term_to_string(*arm.body) : "");
        }
        return arm.tag + " " + arm.arg + " -> " + (arm.body ? term_to_string(*arm.body) : "");
      };
      std::string s = "case " + (m.term ? term_to_string(*m.term) : "") + " {";
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
      const auto& p = term.data.as_projection();
      inner = (p.term ? term_to_string(*p.term) : "") + "." + p.label;
      break;
    }
    case IrTermData::Tag::return_t: {
      const auto& r = term.data.as_return_t();
      inner = "return " + (r ? term_to_string(*r) : "");
      break;
    }
    case IrTermData::Tag::set_t: {
      const auto& st = term.data.as_set_t();
      std::string s = "set " + (st.term ? term_to_string(*st.term) : "") + " {";
      for (size_t i = 0; i < st.values.size(); ++i) {
        if (i > 0) s += ", ";
        s += st.values[i].label + " = " +
             (st.values[i].value ? term_to_string(*st.values[i].value) : "");
      }
      s += "}";
      inner = s;
      break;
    }
    case IrTermData::Tag::struct_t: {
      const auto& vals = term.data.as_struct_t();
      std::string s = "struct{";
      for (size_t i = 0; i < vals.size(); ++i) {
        if (i > 0) s += ", ";
        s += vals[i].label + " = " +
             (vals[i].value ? term_to_string(*vals[i].value) : "");
      }
      s += "}";
      inner = s;
      break;
    }
    case IrTermData::Tag::tuple_t: {
      const auto& elems = term.data.as_tuple_t();
      std::string s = "(";
      for (size_t i = 0; i < elems.size(); ++i) {
        if (i > 0) s += ", ";
        s += term_to_string(elems[i]);
      }
      s += ")";
      inner = s;
      break;
    }
    case IrTermData::Tag::type_abs: {
      const auto& ta = term.data.as_type_abs();
      inner = "Λ" + type_param_to_string(ta.param) + ". " +
              (ta.body ? term_to_string(*ta.body) : "");
      break;
    }
    case IrTermData::Tag::var_t:
      inner = term.data.as_var_t();
      break;
  }

  if (!term.typ.has_value()) {
    return inner;
  }

  bool term_needs_parens = false;
  switch (term.data.tag()) {
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
  switch (term.typ->data.tag()) {
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
  res += type_to_string(*term.typ);
  if (type_needs_parens) res += ")";
  return res;
}

inline IrTerm new_var_term(const std::string& id) {
  return IrTerm{IrTermData::mk_var_t(id), std::nullopt, Pos{}};
}

inline IrTerm new_const_term(IrLiteral literal) {
  return IrTerm{IrTermData::mk_const_t(std::move(literal)), std::nullopt, Pos{}};
}

inline IrTerm new_let_term(std::string var, std::optional<IrType> var_type, IrTerm value) {
  return IrTerm{
      IrTermData::mk_let_t({
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
      IrTermData::mk_assign({
          std::make_shared<IrTerm>(std::move(arg)),
          std::make_shared<IrTerm>(std::move(ret)),
      }),
      std::nullopt,
      Pos{},
  };
}

inline IrTerm new_block_term(std::vector<IrTerm> terms) {
  return IrTerm{IrTermData::mk_block(std::move(terms)), std::nullopt, Pos{}};
}

inline IrTerm new_return_term(IrTerm expr) {
  return IrTerm{
      IrTermData::mk_return_t(std::make_shared<IrTerm>(std::move(expr))),
      std::nullopt,
      Pos{},
  };
}

inline IrTerm new_tuple_term(std::vector<IrTerm> elems) {
  if (elems.size() == 1) {
    return elems[0];
  }
  return IrTerm{IrTermData::mk_tuple_t(std::move(elems)), std::nullopt, Pos{}};
}

inline IrTerm new_struct_term(std::vector<LabelValue> values) {
  return IrTerm{IrTermData::mk_struct_t(std::move(values)), std::nullopt, Pos{}};
}

inline IrTerm new_injection_term(IrType vtype, std::string tag, IrTerm value, std::optional<int64_t> tag_idx = std::nullopt) {
  return IrTerm{
      IrTermData::mk_injection({
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
      IrTermData::mk_app_term({
          std::make_shared<IrTerm>(std::move(fun)),
          std::make_shared<IrTerm>(std::move(arg)),
      }),
      std::nullopt,
      Pos{},
  };
}

inline IrTerm new_app_type_term(IrTerm fun, IrType arg) {
  return IrTerm{
      IrTermData::mk_app_type({
          std::make_shared<IrTerm>(std::move(fun)),
          std::move(arg),
      }),
      std::nullopt,
      Pos{},
  };
}

inline IrTerm new_projection_term(IrTerm term, std::string label) {
  return IrTerm{
      IrTermData::mk_projection({
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
      IrTermData::mk_lambda({
          std::move(arg),
          std::make_shared<IrTerm>(std::move(body)),
      }),
      std::nullopt,
      Pos{},
  };
}

inline IrTerm new_match_term(IrTerm term, std::vector<MatchArm> arms) {
  return IrTerm{
      IrTermData::mk_match_t({
          std::make_shared<IrTerm>(std::move(term)),
          std::move(arms),
      }),
      std::nullopt,
      Pos{},
  };
}

inline IrTerm new_set_term(IrTerm term, std::vector<LabelValue> values) {
  return IrTerm{
      IrTermData::mk_set_t({
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
      IrTermData::mk_type_abs({
          std::move(tp),
          std::make_shared<IrTerm>(std::move(body)),
      }),
      std::nullopt,
      Pos{},
  };
}

} // namespace ir
