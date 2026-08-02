#pragma once

#include "ast/ast_decl.h"
#include "ast/ast_expr.h"
#include "ast/ast_source_file.h"
#include "bin/ir_base.h"
#include "bin/ir_decl.h"
#include "bin/ir_function.h"
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

} // namespace ast
