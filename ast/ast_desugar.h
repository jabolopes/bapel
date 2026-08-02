#pragma once

#include "ast/ast_decl.h"
#include "ast/ast_expr.h"
#include "ast/ast_source_file.h"
#include "bin/ir_base.h"
#include "bin/ir_decl.h"
#include "bin/ir_function.h"
#include "bin/ir_parser.h"
#include "bin/ir_term.h"
#include "bin/ir_type.h"
#include "bin/ir_unit.h"

#include <memory>
#include <string>
#include <vector>

namespace ast {

inline ir::IrTerm desugar_expr(const Expr& expr);

inline ir::IrTerm desugar_expr(const Expr& expr) {
  ir::IrTerm res;
  res.pos = expr.pos;

  switch (expr.case_val) {
    case ExprCase::AppTermExpr:
      if (expr.app_term_data) {
        ir::IrTerm fun = expr.app_term_data->fun ? desugar_expr(*expr.app_term_data->fun) : ir::IrTerm{};
        ir::IrTerm arg = expr.app_term_data->arg ? desugar_expr(*expr.app_term_data->arg) : ir::IrTerm{};
        res = ir::new_app_term(std::move(fun), std::move(arg));
      }
      break;

    case ExprCase::AppTypeExpr:
      if (expr.app_type_data) {
        ir::IrTerm fun = expr.app_type_data->fun ? desugar_expr(*expr.app_type_data->fun) : ir::IrTerm{};
        res = ir::new_app_type_term(std::move(fun), expr.app_type_data->arg);
      }
      break;

    case ExprCase::AssignExpr:
      if (expr.assign_data) {
        ir::IrTerm ret = expr.assign_data->ret ? desugar_expr(*expr.assign_data->ret) : ir::IrTerm{};
        ir::IrTerm arg = expr.assign_data->arg ? desugar_expr(*expr.assign_data->arg) : ir::IrTerm{};
        res = ir::new_assign_term(std::move(ret), std::move(arg));
      }
      break;

    case ExprCase::BlockExpr:
      if (expr.block_data) {
        std::vector<ir::IrTerm> terms;
        terms.reserve(expr.block_data->exprs.size());
        for (const auto& e : expr.block_data->exprs) {
          terms.push_back(desugar_expr(e));
        }
        res = ir::new_block_term(std::move(terms));
      }
      break;

    case ExprCase::ConstExpr:
      if (expr.const_data) {
        res.case_val = ir::IrTermCase::ConstTerm;
        res.const_data = std::make_shared<ir::ConstTermData>();
        res.const_data->literal = expr.const_data->literal;
      }
      break;

    case ExprCase::ForExpr:
      if (expr.for_data) {
        ir::IrTerm body = expr.for_data->body ? desugar_expr(*expr.for_data->body) : ir::IrTerm{};
        ir::IrTerm lambda;
        lambda.case_val = ir::IrTermCase::LambdaTerm;
        lambda.lambda = std::make_shared<ir::LambdaTermData>();
        lambda.lambda->arg = ir::FunctionArg{"_", ir::new_tuple_type({})};
        lambda.lambda->body = std::make_shared<ir::IrTerm>(std::move(body));

        ir::IrTerm cond = expr.for_data->condition ? desugar_expr(*expr.for_data->condition) : ir::IrTerm{};
        ir::IrTerm tuple_arg = ir::new_tuple_term({std::move(cond), std::move(lambda)});
        res = ir::new_app_term(ir::new_var_term("core::for"), std::move(tuple_arg));
      }
      break;

    case ExprCase::InjectionExpr:
      if (expr.injection_data) {
        ir::IrTerm val = expr.injection_data->expr ? desugar_expr(*expr.injection_data->expr) : ir::IrTerm{};
        res = ir::new_injection_term(expr.injection_data->variant_type, expr.injection_data->tag, std::move(val));
      }
      break;

    case ExprCase::LambdaExpr:
      if (expr.lambda_data) {
        res.case_val = ir::IrTermCase::LambdaTerm;
        res.lambda = std::make_shared<ir::LambdaTermData>();
        res.lambda->arg = expr.lambda_data->arg;
        res.lambda->body = std::make_shared<ir::IrTerm>(expr.lambda_data->body ? desugar_expr(*expr.lambda_data->body) : ir::IrTerm{});
      }
      break;

    case ExprCase::LetExpr:
      if (expr.let_data) {
        ir::IrTerm val = expr.let_data->expr ? desugar_expr(*expr.let_data->expr) : ir::IrTerm{};
        res = ir::new_let_term(expr.let_data->var, expr.let_data->var_type, std::move(val));
      }
      break;

    case ExprCase::MatchExpr:
      if (expr.match_data) {
        res.case_val = ir::IrTermCase::MatchTerm;
        res.match_data = std::make_shared<ir::MatchTermData>();
        res.match_data->term = std::make_shared<ir::IrTerm>(expr.match_data->expr ? desugar_expr(*expr.match_data->expr) : ir::IrTerm{});
        for (const auto& arm : expr.match_data->arms) {
          ir::MatchArm ir_arm;
          ir_arm.tag = arm.tag;
          ir_arm.arg = arm.arg;
          ir_arm.index = arm.index;
          ir_arm.body = std::make_shared<ir::IrTerm>(arm.body ? desugar_expr(*arm.body) : ir::IrTerm{});
          res.match_data->arms.push_back(std::move(ir_arm));
        }
      }
      break;

    case ExprCase::ProjectionExpr:
      if (expr.projection_data) {
        ir::IrTerm term = expr.projection_data->expr ? desugar_expr(*expr.projection_data->expr) : ir::IrTerm{};
        res = ir::new_projection_term(std::move(term), expr.projection_data->label);
      }
      break;

    case ExprCase::ReturnExpr:
      if (expr.return_data) {
        ir::IrTerm val = expr.return_data->expr ? desugar_expr(*expr.return_data->expr) : ir::IrTerm{};
        res = ir::new_return_term(std::move(val));
      }
      break;

    case ExprCase::SetExpr:
      if (expr.set_data) {
        res.case_val = ir::IrTermCase::SetTerm;
        res.set_data = std::make_shared<ir::SetTermData>();
        res.set_data->term = std::make_shared<ir::IrTerm>(expr.set_data->expr ? desugar_expr(*expr.set_data->expr) : ir::IrTerm{});
        for (const auto& lv : expr.set_data->values) {
          ir::LabelValue ir_lv;
          ir_lv.label = lv.label;
          ir_lv.value = std::make_shared<ir::IrTerm>(lv.value ? desugar_expr(*lv.value) : ir::IrTerm{});
          res.set_data->values.push_back(std::move(ir_lv));
        }
      }
      break;

    case ExprCase::StructExpr:
      if (expr.struct_data) {
        std::vector<ir::LabelValue> values;
        for (const auto& lv : expr.struct_data->values) {
          ir::LabelValue ir_lv;
          ir_lv.label = lv.label;
          ir_lv.value = std::make_shared<ir::IrTerm>(lv.value ? desugar_expr(*lv.value) : ir::IrTerm{});
          values.push_back(std::move(ir_lv));
        }
        res = ir::new_struct_term(std::move(values));
      }
      break;

    case ExprCase::TupleExpr:
      if (expr.tuple_data) {
        std::vector<ir::IrTerm> elems;
        elems.reserve(expr.tuple_data->elems.size());
        for (const auto& elem : expr.tuple_data->elems) {
          elems.push_back(desugar_expr(elem));
        }
        res = ir::new_tuple_term(std::move(elems));
      }
      break;

    case ExprCase::TypeAbsExpr:
      if (expr.type_abs_data) {
        res.case_val = ir::IrTermCase::TypeAbsTerm;
        res.type_abs = std::make_shared<ir::TypeAbsTermData>();
        res.type_abs->type_param = expr.type_abs_data->arg;
        res.type_abs->body = std::make_shared<ir::IrTerm>(expr.type_abs_data->body ? desugar_expr(*expr.type_abs_data->body) : ir::IrTerm{});
      }
      break;

    case ExprCase::VarExpr:
      if (expr.var_data) {
        res = ir::new_var_term(expr.var_data->id);
      }
      break;
  }

  res.pos = expr.pos;
  return res;
}

inline ir::IrFunction desugar_function(const Function& fn) {
  ir::IrTerm body = desugar_expr(fn.body);
  ir::IrFunction f = ir::new_function(fn.export_flag, fn.id, fn.type_params, fn.args, fn.ret_type, std::move(body));
  f.pos = fn.pos;
  return f;
}

inline ir::IrTraitImpl desugar_impl(const Impl& impl) {
  ir::IrTraitImpl res;
  res.case_val = (impl.case_val == ImplCase::TraitImpl) ? ir::ImplCase::TraitImpl : ir::ImplCase::InherentImpl;
  res.type_params = impl.type_params;
  res.trait_type = impl.trait_type;
  res.type_name = impl.type_name;
  res.pos = impl.pos;
  for (const auto& m : impl.methods) {
    res.methods.push_back(desugar_function(m));
  }
  return res;
}

inline ir::IrUnit desugar_source_file(const SourceFile& sf) {
  ir::IrUnit unit;
  unit.case_val = sf.header.is(SourceFileCase::BaseSourceFile) ? ir::IrUnitCase::BaseUnit : ir::IrUnitCase::ImplUnit;
  unit.module_id = sf.header.module_id;
  unit.filename = sf.header.filename;

  for (const auto& id : sf.imports.ids) {
    unit.imports.push_back(ir::new_import(id));
  }

  for (const auto& fn : sf.impls.filenames) {
    unit.impls.push_back(ir::new_impl(fn));
  }

  for (const auto& s : sf.body) {
    switch (s.case_val) {
      case SourceCase::DeclSource:
        if (s.decl_data) {
          unit.decls.push_back(*s.decl_data);
        }
        break;
      case SourceCase::FunctionSource:
        if (s.function_data) {
          unit.functions.push_back(desugar_function(*s.function_data));
        }
        break;
      case SourceCase::TraitSource:
        if (s.trait_data) {
          unit.decls.push_back(s.trait_data->decl());
        }
        break;
      case SourceCase::ImplSource:
        if (s.impl_data) {
          unit.trait_impls.push_back(desugar_impl(*s.impl_data));
        }
        break;
    }
  }

  return unit;
}

// JSON deserialization of SourceFile from parser output
inline Expr deserialize_ast_expr(const ir::JsonValue& j);

inline ir::IrLiteral deserialize_ast_literal(const ir::JsonValue& j) {
  ir::IrLiteral lit;
  if (!j.is_object()) return lit;
  lit.pos = ir::deserialize_pos(j.get("Pos"));
  if (j.has_field("Int")) {
    lit.case_val = ir::IrLiteralCase::IntLiteral;
    lit.int_val = j.get("Int").as_int();
  } else if (j.has_field("Float")) {
    lit.case_val = ir::IrLiteralCase::FloatLiteral;
    const auto& f = j.get("Float");
    lit.float_val = ir::FloatLit{f.get("Integer").as_int(), f.get("Decimal").as_int()};
  } else if (j.has_field("Rune")) {
    lit.case_val = ir::IrLiteralCase::RuneLiteral;
    lit.rune_val = j.get("Rune").as_string();
  } else if (j.has_field("Str")) {
    lit.case_val = ir::IrLiteralCase::StrLiteral;
    lit.str_val = j.get("Str").as_string();
  }
  return lit;
}

inline MatchArm deserialize_ast_match_arm(const ir::JsonValue& j) {
  MatchArm arm;
  if (!j.is_object()) return arm;
  arm.tag = j.get("Tag").as_string();
  arm.arg = j.get("Arg").as_string();
  arm.body = std::make_shared<Expr>(deserialize_ast_expr(j.get("Body")));
  if (j.has_field("Index")) {
    arm.index = static_cast<int>(j.get("Index").as_int());
  }
  return arm;
}

inline LabelValue deserialize_ast_label_value(const ir::JsonValue& j) {
  LabelValue lv;
  if (!j.is_object()) return lv;
  lv.label = j.get("Label").as_string();
  lv.value = std::make_shared<Expr>(deserialize_ast_expr(j.get("Value")));
  return lv;
}

inline Expr deserialize_ast_expr(const ir::JsonValue& j) {
  Expr e;
  if (!j.is_object()) return e;
  e.case_val = static_cast<ExprCase>(j.get("Case").as_int());
  e.pos = ir::deserialize_pos(j.get("Pos"));

  switch (e.case_val) {
    case ExprCase::AppTermExpr:
      if (j.has_field("AppTerm")) {
        const auto& at = j.get("AppTerm");
        e.app_term_data = std::make_shared<AppTermData>();
        e.app_term_data->fun = std::make_shared<Expr>(deserialize_ast_expr(at.get("Fun")));
        e.app_term_data->arg = std::make_shared<Expr>(deserialize_ast_expr(at.get("Arg")));
      }
      break;

    case ExprCase::AppTypeExpr:
      if (j.has_field("AppType")) {
        const auto& at = j.get("AppType");
        e.app_type_data = std::make_shared<AppTypeData>();
        e.app_type_data->fun = std::make_shared<Expr>(deserialize_ast_expr(at.get("Fun")));
        e.app_type_data->arg = ir::deserialize_type(at.get("Arg"));
      }
      break;

    case ExprCase::AssignExpr:
      if (j.has_field("Assign")) {
        const auto& asg = j.get("Assign");
        e.assign_data = std::make_shared<AssignData>();
        e.assign_data->ret = std::make_shared<Expr>(deserialize_ast_expr(asg.get("Ret")));
        e.assign_data->arg = std::make_shared<Expr>(deserialize_ast_expr(asg.get("Arg")));
      }
      break;

    case ExprCase::BlockExpr:
      if (j.has_field("Block") && j.get("Block").has_field("Exprs")) {
        e.block_data = std::make_shared<BlockData>();
        for (const auto& item : j.get("Block").get("Exprs").as_array()) {
          e.block_data->exprs.push_back(deserialize_ast_expr(item));
        }
      }
      break;

    case ExprCase::ConstExpr:
      if (j.has_field("Const")) {
        e.const_data = std::make_shared<ConstData>();
        e.const_data->literal = deserialize_ast_literal(j.get("Const"));
      }
      break;

    case ExprCase::ForExpr:
      if (j.has_field("For")) {
        const auto& fd = j.get("For");
        e.for_data = std::make_shared<ForData>();
        e.for_data->condition = std::make_shared<Expr>(deserialize_ast_expr(fd.get("Condition")));
        e.for_data->body = std::make_shared<Expr>(deserialize_ast_expr(fd.get("Body")));
      }
      break;

    case ExprCase::InjectionExpr:
      if (j.has_field("Injection")) {
        const auto& inj = j.get("Injection");
        e.injection_data = std::make_shared<InjectionData>();
        e.injection_data->variant_type = ir::deserialize_type(inj.get("VariantType"));
        e.injection_data->tag = inj.get("Tag").as_string();
        e.injection_data->expr = std::make_shared<Expr>(deserialize_ast_expr(inj.get("Expr")));
      }
      break;

    case ExprCase::LambdaExpr:
      if (j.has_field("Lambda")) {
        const auto& l = j.get("Lambda");
        e.lambda_data = std::make_shared<LambdaData>();
        e.lambda_data->arg = ir::deserialize_function_arg(l.get("Arg"));
        e.lambda_data->body = std::make_shared<Expr>(deserialize_ast_expr(l.get("Body")));
      }
      break;

    case ExprCase::LetExpr:
      if (j.has_field("Let")) {
        const auto& lt = j.get("Let");
        e.let_data = std::make_shared<LetData>();
        e.let_data->var = lt.get("Var").as_string();
        if (lt.has_field("VarType")) {
          e.let_data->var_type = ir::deserialize_type(lt.get("VarType"));
        }
        e.let_data->expr = std::make_shared<Expr>(deserialize_ast_expr(lt.get("Expr")));
      }
      break;

    case ExprCase::MatchExpr:
      if (j.has_field("Match")) {
        const auto& m = j.get("Match");
        e.match_data = std::make_shared<MatchData>();
        e.match_data->expr = std::make_shared<Expr>(deserialize_ast_expr(m.get("Expr")));
        if (m.has_field("Arms")) {
          for (const auto& arm_j : m.get("Arms").as_array()) {
            e.match_data->arms.push_back(deserialize_ast_match_arm(arm_j));
          }
        }
      }
      break;

    case ExprCase::ProjectionExpr:
      if (j.has_field("Projection")) {
        const auto& p = j.get("Projection");
        e.projection_data = std::make_shared<ProjectionData>();
        e.projection_data->expr = std::make_shared<Expr>(deserialize_ast_expr(p.get("Expr")));
        e.projection_data->label = p.get("Label").as_string();
      }
      break;

    case ExprCase::ReturnExpr:
      if (j.has_field("Return")) {
        e.return_data = std::make_shared<ReturnData>();
        e.return_data->expr = std::make_shared<Expr>(deserialize_ast_expr(j.get("Return").get("Expr")));
      }
      break;

    case ExprCase::SetExpr:
      if (j.has_field("Set")) {
        const auto& st = j.get("Set");
        e.set_data = std::make_shared<SetData>();
        e.set_data->expr = std::make_shared<Expr>(deserialize_ast_expr(st.get("Expr")));
        if (st.has_field("Values")) {
          for (const auto& lv_j : st.get("Values").as_array()) {
            e.set_data->values.push_back(deserialize_ast_label_value(lv_j));
          }
        }
      }
      break;

    case ExprCase::StructExpr:
      if (j.has_field("Struct") && j.get("Struct").has_field("Values")) {
        e.struct_data = std::make_shared<StructData>();
        for (const auto& lv_j : j.get("Struct").get("Values").as_array()) {
          e.struct_data->values.push_back(deserialize_ast_label_value(lv_j));
        }
      }
      break;

    case ExprCase::TupleExpr:
      if (j.has_field("Tuple") && j.get("Tuple").has_field("Elems")) {
        e.tuple_data = std::make_shared<TupleData>();
        for (const auto& elem_j : j.get("Tuple").get("Elems").as_array()) {
          e.tuple_data->elems.push_back(deserialize_ast_expr(elem_j));
        }
      }
      break;

    case ExprCase::TypeAbsExpr:
      if (j.has_field("TypeAbs")) {
        const auto& ta = j.get("TypeAbs");
        e.type_abs_data = std::make_shared<TypeAbsData>();
        e.type_abs_data->arg = ir::deserialize_type_param(ta.has_field("TypeParam") ? ta.get("TypeParam") : ta.get("Arg"));
        e.type_abs_data->body = std::make_shared<Expr>(deserialize_ast_expr(ta.get("Body")));
      }
      break;

    case ExprCase::VarExpr:
      if (j.has_field("Var")) {
        e.var_data = std::make_shared<VarData>();
        e.var_data->id = j.get("Var").get("ID").as_string();
      }
      break;
  }

  return e;
}

inline Function deserialize_ast_function(const ir::JsonValue& j) {
  Function f;
  if (!j.is_object()) return f;
  f.export_flag = j.get("Export").as_bool();
  f.id = j.get("ID").as_string();
  f.ret_type = ir::deserialize_type(j.get("RetType"));
  f.body = deserialize_ast_expr(j.get("Body"));
  f.pos = ir::deserialize_pos(j.get("Pos"));

  if (j.has_field("TypeParams")) {
    for (const auto& tp_j : j.get("TypeParams").as_array()) {
      f.type_params.push_back(ir::deserialize_type_param(tp_j));
    }
  }
  if (j.has_field("Args")) {
    for (const auto& arg_j : j.get("Args").as_array()) {
      f.args.push_back(ir::deserialize_function_arg(arg_j));
    }
  }
  return f;
}

inline Trait deserialize_ast_trait(const ir::JsonValue& j) {
  Trait tr;
  if (!j.is_object()) return tr;
  tr.export_flag = j.get("Export").as_bool();
  tr.id = j.get("ID").as_string();
  tr.pos = ir::deserialize_pos(j.get("Pos"));

  if (j.has_field("TypeParams")) {
    for (const auto& tp_j : j.get("TypeParams").as_array()) {
      tr.type_params.push_back(ir::deserialize_type_param(tp_j));
    }
  }
  if (j.has_field("Methods")) {
    for (const auto& m_j : j.get("Methods").as_array()) {
      Signature sig;
      sig.id = m_j.get("ID").as_string();
      sig.ret_type = ir::deserialize_type(m_j.get("RetType"));
      sig.pos = ir::deserialize_pos(m_j.get("Pos"));
      if (m_j.has_field("Args")) {
        for (const auto& a_j : m_j.get("Args").as_array()) {
          sig.args.push_back(ir::deserialize_function_arg(a_j));
        }
      }
      tr.methods.push_back(std::move(sig));
    }
  }
  return tr;
}

inline Impl deserialize_ast_impl(const ir::JsonValue& j) {
  Impl impl;
  if (!j.is_object()) return impl;
  impl.case_val = static_cast<ImplCase>(j.get("Case").as_int());
  impl.trait_type = ir::deserialize_type(j.get("TraitType"));
  impl.type_name = ir::deserialize_type(j.get("TypeName"));
  impl.pos = ir::deserialize_pos(j.get("Pos"));

  if (j.has_field("TypeParams")) {
    for (const auto& tp_j : j.get("TypeParams").as_array()) {
      impl.type_params.push_back(ir::deserialize_type_param(tp_j));
    }
  }
  if (j.has_field("Methods")) {
    for (const auto& m_j : j.get("Methods").as_array()) {
      impl.methods.push_back(deserialize_ast_function(m_j));
    }
  }
  return impl;
}

inline SourceFile deserialize_ast_source_file(const ir::JsonValue& j) {
  SourceFile sf;
  if (!j.is_object()) return sf;

  if (j.has_field("Header")) {
    const auto& h = j.get("Header");
    sf.header.case_val = static_cast<SourceFileCase>(h.get("Case").as_int());
    sf.header.module_id = ir::deserialize_module_id(h.get("ModuleID"));
    sf.header.filename = ir::deserialize_filename(h.get("Filename"));
  }

  if (j.has_field("Imports") && j.get("Imports").has_field("IDs")) {
    for (const auto& id_j : j.get("Imports").get("IDs").as_array()) {
      sf.imports.ids.push_back(ir::deserialize_module_id(id_j));
    }
  }

  if (j.has_field("Impls") && j.get("Impls").has_field("Filenames")) {
    for (const auto& fn_j : j.get("Impls").get("Filenames").as_array()) {
      sf.impls.filenames.push_back(ir::deserialize_filename(fn_j));
    }
  }

  if (j.has_field("Flags") && j.get("Flags").has_field("Filenames")) {
    for (const auto& fn_j : j.get("Flags").get("Filenames").as_array()) {
      sf.flags.filenames.push_back(ir::deserialize_filename(fn_j));
    }
  }

  if (j.has_field("Body")) {
    for (const auto& item_j : j.get("Body").as_array()) {
      int item_case = item_j.get("Case").as_int();
      switch (static_cast<SourceCase>(item_case)) {
        case SourceCase::DeclSource:
          if (item_j.has_field("Decl")) {
            Source s;
            s.case_val = SourceCase::DeclSource;
            s.decl_data = std::make_shared<ir::IrDecl>(ir::deserialize_decl(item_j.get("Decl").has_field("Decl") ? item_j.get("Decl").get("Decl") : item_j.get("Decl")));
            sf.body.push_back(std::move(s));
          }
          break;
        case SourceCase::FunctionSource:
          if (item_j.has_field("Function")) {
            Source s;
            s.case_val = SourceCase::FunctionSource;
            s.function_data = std::make_shared<Function>(deserialize_ast_function(item_j.get("Function")));
            sf.body.push_back(std::move(s));
          }
          break;
        case SourceCase::TraitSource:
          if (item_j.has_field("Trait")) {
            Source s;
            s.case_val = SourceCase::TraitSource;
            s.trait_data = std::make_shared<Trait>(deserialize_ast_trait(item_j.get("Trait")));
            sf.body.push_back(std::move(s));
          }
          break;
        case SourceCase::ImplSource:
          if (item_j.has_field("Impl")) {
            Source s;
            s.case_val = SourceCase::ImplSource;
            s.impl_data = std::make_shared<Impl>(deserialize_ast_impl(item_j.get("Impl")));
            sf.body.push_back(std::move(s));
          }
          break;
      }
    }
  }

  return sf;
}

} // namespace ast
