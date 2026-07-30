#pragma once

#include "antlr4-runtime.h"
#include "ast/ast.h"
#include "generated/bapelBaseVisitor.h"
#include "generated/bapelParser.h"
#include <cinttypes>
#include <climits>
#include <cstdio>
#include <sstream>
#include <string>
#include <vector>

namespace ast {

struct ID {
  std::string value;
  Pos pos;
};

inline ID new_id(std::string value, Pos pos) {
  return ID{std::move(value), std::move(pos)};
}

inline Pos pos_from_token(const std::string& filename, antlr4::Token* token) {
  if (token == nullptr) return Pos{filename, 0, 0};
  int64_t line = static_cast<int64_t>(token->getLine());
  return ir::new_line_pos(filename, line);
}

inline Pos pos_from_context(const std::string& filename, antlr4::ParserRuleContext* ctx) {
  if (ctx == nullptr || ctx->getStart() == nullptr) return Pos{filename, 0, 0};
  int64_t start_line = static_cast<int64_t>(ctx->getStart()->getLine());
  int64_t stop_line = ctx->getStop() != nullptr ? static_cast<int64_t>(ctx->getStop()->getLine()) : start_line;
  return ir::new_range_pos(filename, start_line, stop_line);
}

inline Pos make_pos(const Pos& p1, const Pos& p2) {
  std::string fn = !p1.filename.empty() ? p1.filename : p2.filename;
  int64_t start = p1.begin_line_num > 0 ? p1.begin_line_num : p2.begin_line_num;
  int64_t stop = p2.end_line_num > 0 ? p2.end_line_num : p1.end_line_num;
  return ir::new_range_pos(fn, start, stop);
}

inline int64_t parse_number(const std::string& arg) {
  int64_t val = 0;
  if (arg.rfind("0x", 0) == 0 || arg.rfind("0X", 0) == 0) {
    sscanf(arg.c_str(), "0x%" SCNx64, &val);
  } else {
    sscanf(arg.c_str(), "%" SCNd64, &val);
  }
  return val;
}

inline bool parse_float(const std::string& arg, int64_t& integer, int64_t& decimal) {
  size_t dot = arg.find('.');
  if (dot == std::string::npos) return false;
  integer = parse_number(arg.substr(0, dot));
  decimal = parse_number(arg.substr(dot + 1));
  return true;
}

inline ir::IrLiteral new_number_literal(const Pos& pos, const std::string& text) {
  if (text.find('.') != std::string::npos) {
    int64_t integer = 0, decimal = 0;
    parse_float(text, integer, decimal);
    return ir::new_float_literal(pos, integer, decimal);
  }
  int64_t val = parse_number(text);
  return ir::new_int_literal(pos, val);
}

inline ir::IrLiteral parse_rune_literal(const Pos& pos, std::string text) {
  if (text.size() >= 2 && text.front() == '\'' && text.back() == '\'') {
    text = text.substr(1, text.size() - 2);
  }
  return ir::new_rune_literal(pos, text);
}

inline ir::IrLiteral parse_string_literal(const Pos& pos, std::string text) {
  if (text.size() >= 2 && text.front() == '"' && text.back() == '"') {
    text = text.substr(1, text.size() - 2);
  } else if (text.size() >= 2 && text.front() == '`' && text.back() == '`') {
    text = text.substr(1, text.size() - 2);
  }
  return ir::new_str_literal(pos, text);
}

inline Expr call(Pos pos, Expr id, const std::vector<ir::IrType>& types, std::vector<Expr> args) {
  Expr expr = std::move(id);
  for (const auto& typ : types) {
    expr = new_app_type_expr(pos, std::move(expr), typ);
  }
  if (args.empty()) {
    return expr;
  }
  return new_app_term_expr(pos, std::move(expr), new_tuple_expr(pos, std::move(args)));
}

inline Expr lambda_dsl(Pos pos, const std::vector<ir::TypeParam>& tvars, const std::vector<ir::FunctionArg>& args, Expr body, size_t tvar_idx = 0, size_t arg_idx = 0) {
  if (tvar_idx < tvars.size()) {
    return new_type_abs_expr(pos, tvars[tvar_idx], lambda_dsl(pos, tvars, args, std::move(body), tvar_idx + 1, arg_idx));
  }
  if (arg_idx < args.size()) {
    return new_lambda_expr(pos, args[arg_idx], lambda_dsl(pos, tvars, args, std::move(body), tvar_idx, arg_idx + 1));
  }
  return body;
}

inline Expr new_bin_op_expr(Expr op_expr, const std::vector<ir::IrType>& type_args, Expr t1, Expr t2) {
  Pos pos = make_pos(t1.pos, t2.pos);
  std::vector<Expr> args;
  args.push_back(std::move(t1));
  args.push_back(std::move(t2));
  return call(pos, std::move(op_expr), type_args, std::move(args));
}

inline Expr new_unary_op_expr(Expr op_expr, const std::vector<ir::IrType>& type_args, Expr operand) {
  Pos pos = make_pos(op_expr.pos, operand.pos);
  if (op_expr.is(ExprCase::VarExpr) && op_expr.var_data && op_expr.var_data->id == "-") {
    std::vector<Expr> args;
    args.push_back(new_const_expr(pos, ir::new_int_literal(pos, 0)));
    args.push_back(std::move(operand));
    return call(pos, std::move(op_expr), type_args, std::move(args));
  }
  std::vector<Expr> args;
  args.push_back(std::move(operand));
  return call(pos, std::move(op_expr), type_args, std::move(args));
}

class AstBuilder : public bapelBaseVisitor {
public:
  explicit AstBuilder(std::string filename) : filename_(std::move(filename)) {}

  antlrcpp::Any visitBaseSourceFile(bapelParser::BaseSourceFileContext *ctx) override {
    SourceFile sf;
    sf.header = ctx->moduleHeader()->accept(this).as<SourceFileHeader>();
    if (ctx->importsSection()) {
      sf.imports = ctx->importsSection()->accept(this).as<Imports>();
    }
    if (ctx->implsSection()) {
      sf.impls = ctx->implsSection()->accept(this).as<Impls>();
    }
    if (ctx->flagsSection()) {
      sf.flags = ctx->flagsSection()->accept(this).as<Flags>();
    }
    if (ctx->sources()) {
      sf.body = ctx->sources()->accept(this).as<std::vector<Source>>();
    }
    return sf;
  }

  antlrcpp::Any visitImplSourceFile(bapelParser::ImplSourceFileContext *ctx) override {
    SourceFile sf;
    sf.header = ctx->implementsHeader()->accept(this).as<SourceFileHeader>();
    if (ctx->importsSection()) {
      sf.imports = ctx->importsSection()->accept(this).as<Imports>();
    }
    if (ctx->implsSection()) {
      sf.impls = ctx->implsSection()->accept(this).as<Impls>();
    }
    if (ctx->flagsSection()) {
      sf.flags = ctx->flagsSection()->accept(this).as<Flags>();
    }
    if (ctx->sources()) {
      sf.body = ctx->sources()->accept(this).as<std::vector<Source>>();
    }
    return sf;
  }

  antlrcpp::Any visitModuleHeader(bapelParser::ModuleHeaderContext *ctx) override {
    ir::ModuleID mid = ctx->moduleID()->accept(this).as<ir::ModuleID>();
    SourceFileHeader h;
    h.case_val = SourceFileCase::BaseSourceFile;
    h.module_id = std::move(mid);
    return h;
  }

  antlrcpp::Any visitImplementsHeader(bapelParser::ImplementsHeaderContext *ctx) override {
    ir::ModuleID mid = ctx->moduleID()->accept(this).as<ir::ModuleID>();
    SourceFileHeader h;
    h.case_val = SourceFileCase::ImplSourceFile;
    h.module_id = std::move(mid);
    return h;
  }

  antlrcpp::Any visitWorkspace(bapelParser::WorkspaceContext *ctx) override {
    Packages pkgs = ctx->packagesSection()->accept(this).as<Packages>();
    Workspace ws;
    ws.packages = std::move(pkgs);
    return ws;
  }

  antlrcpp::Any visitPackagesSection(bapelParser::PackagesSectionContext *ctx) override {
    std::vector<Package> pkgs;
    for (auto* rule : ctx->packageRule()) {
      pkgs.push_back(rule->accept(this).as<Package>());
    }
    Packages p;
    p.packages = std::move(pkgs);
    p.pos = pos_from_context(filename_, ctx);
    return p;
  }

  antlrcpp::Any visitPackageRule(bapelParser::PackageRuleContext *ctx) override {
    ir::ModuleID mid = ctx->moduleID()->accept(this).as<ir::ModuleID>();
    ir::Filename fn = ctx->filename()->accept(this).as<ir::Filename>();
    Pos pos = pos_from_context(filename_, ctx);
    if (ctx->getStart()->getText() == "prefix") {
      return new_prefix_package(std::move(mid), std::move(fn), pos);
    }
    return new_module_package(std::move(mid), std::move(fn), pos);
  }

  antlrcpp::Any visitImportsSection(bapelParser::ImportsSectionContext *ctx) override {
    std::vector<ir::ModuleID> ids;
    for (auto* mid : ctx->moduleID()) {
      ids.push_back(mid->accept(this).as<ir::ModuleID>());
    }
    Imports imp;
    imp.ids = std::move(ids);
    imp.pos = pos_from_context(filename_, ctx);
    return imp;
  }

  antlrcpp::Any visitImplsSection(bapelParser::ImplsSectionContext *ctx) override {
    std::vector<ir::Filename> fns;
    for (auto* fn : ctx->filename()) {
      fns.push_back(fn->accept(this).as<ir::Filename>());
    }
    Impls im;
    im.filenames = std::move(fns);
    im.pos = pos_from_context(filename_, ctx);
    return im;
  }

  antlrcpp::Any visitFlagsSection(bapelParser::FlagsSectionContext *ctx) override {
    std::vector<ir::Filename> fns;
    for (auto* fn : ctx->filename()) {
      fns.push_back(fn->accept(this).as<ir::Filename>());
    }
    Flags fl;
    fl.filenames = std::move(fns);
    fl.pos = pos_from_context(filename_, ctx);
    return fl;
  }

  antlrcpp::Any visitModuleID(bapelParser::ModuleIDContext *ctx) override {
    std::string name;
    auto identifiers = ctx->IDENTIFIER();
    for (size_t i = 0; i < identifiers.size(); ++i) {
      if (i > 0) name += ir::ModuleIDSeparator;
      name += identifiers[i]->getText();
    }
    return ir::new_module_id(name, pos_from_context(filename_, ctx));
  }

  antlrcpp::Any visitFilename(bapelParser::FilenameContext *ctx) override {
    std::string text = ctx->STRING_LITERAL()->getText();
    if (text.size() >= 2 && text.front() == '"' && text.back() == '"') {
      text = text.substr(1, text.size() - 2);
    }
    return ir::new_filename(text, pos_from_context(filename_, ctx));
  }

  antlrcpp::Any visitSources(bapelParser::SourcesContext *ctx) override {
    std::vector<Source> sources;
    for (auto* src : ctx->source()) {
      sources.push_back(src->accept(this).as<Source>());
    }
    return sources;
  }

  antlrcpp::Any visitSource(bapelParser::SourceContext *ctx) override {
    if (ctx->declNoExport()) {
      return ctx->declNoExport()->accept(this).as<Source>();
    }
    if (ctx->functionNoExport()) {
      Source fn_src = ctx->functionNoExport()->accept(this).as<Source>();
      if (ctx->getStart()->getText() == "pub" && fn_src.function_data) {
        fn_src.function_data->export_flag = true;
      }
      return fn_src;
    }
    if (ctx->traitDecl()) {
      Trait trait = ctx->traitDecl()->accept(this).as<Trait>();
      if (ctx->getStart()->getText() == "pub") {
        trait.export_flag = true;
      }
      return new_trait_source(std::move(trait));
    }
    if (ctx->implBlock()) {
      return ctx->implBlock()->accept(this).as<Source>();
    }
    return Source{};
  }

  antlrcpp::Any visitDeclNoExport(bapelParser::DeclNoExportContext *ctx) override {
    ir::IrDecl decl;
    if (ctx->declNoTerm()) {
      decl = ctx->declNoTerm()->accept(this).as<ir::IrDecl>();
    } else {
      decl = ctx->termDecl()->accept(this).as<ir::IrDecl>();
      if (ctx->getStart()->getText() == "pub") {
        decl.export_flag = true;
      }
    }
    return new_decl_source(std::move(decl));
  }

  antlrcpp::Any visitFunctionNoExport(bapelParser::FunctionNoExportContext *ctx) override {
    ID id = ctx->id()->accept(this).as<ID>();
    std::vector<ir::TypeParam> tvars;
    if (ctx->typeAbstraction()) {
      tvars = ctx->typeAbstraction()->accept(this).as<std::vector<ir::TypeParam>>();
    }
    std::vector<ir::FunctionArg> fun_args = ctx->functionArgs()->accept(this).as<std::vector<ir::FunctionArg>>();
    ir::IrType ret_type = ctx->type_()->accept(this).as<ir::IrType>();
    Expr body = ctx->blockExpr()->accept(this).as<Expr>();

    Function fn;
    fn.export_flag = false;
    fn.id = std::move(id.value);
    fn.type_params = std::move(tvars);
    fn.args = std::move(fun_args);
    fn.ret_type = std::move(ret_type);
    fn.body = std::move(body);
    fn.pos = pos_from_context(filename_, ctx);
    return new_function_source(std::move(fn));
  }

  antlrcpp::Any visitFunctionArgs(bapelParser::FunctionArgsContext *ctx) override {
    std::vector<ir::FunctionArg> args;
    for (auto* arg : ctx->arg()) {
      args.push_back(arg->accept(this).as<ir::FunctionArg>());
    }
    return args;
  }

  antlrcpp::Any visitArg(bapelParser::ArgContext *ctx) override {
    std::string id_text = ctx->IDENTIFIER()->getText();
    ir::IrType typ = ctx->type_()->accept(this).as<ir::IrType>();
    return ir::FunctionArg{std::move(id_text), std::move(typ)};
  }

  antlrcpp::Any visitDecl(bapelParser::DeclContext *ctx) override {
    ir::IrDecl decl = ctx->unexportedDecl()->accept(this).as<ir::IrDecl>();
    if (ctx->getStart()->getText() == "pub") {
      decl.export_flag = true;
    }
    return decl;
  }

  antlrcpp::Any visitUnexportedDecl(bapelParser::UnexportedDeclContext *ctx) override {
    if (ctx->termDecl()) {
      return ctx->termDecl()->accept(this).as<ir::IrDecl>();
    }
    return ctx->typeDecl()->accept(this).as<ir::IrDecl>();
  }

  antlrcpp::Any visitDeclNoTerm(bapelParser::DeclNoTermContext *ctx) override {
    ir::IrDecl decl = ctx->typeDecl()->accept(this).as<ir::IrDecl>();
    if (ctx->getStart()->getText() == "pub") {
      decl.export_flag = true;
    }
    return decl;
  }

  antlrcpp::Any visitTermDecl(bapelParser::TermDeclContext *ctx) override {
    ID id = ctx->id()->accept(this).as<ID>();
    ir::IrType typ = ctx->type_()->accept(this).as<ir::IrType>();
    ir::IrType qtyp = ir::quantify_type(typ);
    qtyp.pos = typ.pos;
    ir::IrDecl decl = ir::new_term_decl(id.value, std::move(qtyp), false);
    decl.pos = make_pos(id.pos, typ.pos);
    return decl;
  }

  antlrcpp::Any visitTypeDecl(bapelParser::TypeDeclContext *ctx) override {
    ID id = ctx->id()->accept(this).as<ID>();
    std::vector<ir::TypeParam> tvars;
    if (ctx->typeAbstraction()) {
      tvars = ctx->typeAbstraction()->accept(this).as<std::vector<ir::TypeParam>>();
    }
    ir::IrKind kind = ir::new_type_kind();
    for (size_t i = 0; i < tvars.size(); ++i) {
      kind = ir::new_arrow_kind(ir::new_type_kind(), std::move(kind));
    }
    if (ctx->type_()) {
      ir::IrType typ = ctx->type_()->accept(this).as<ir::IrType>();
      ir::IrDecl decl = ir::new_alias_decl(id.value, kind, ir::lambda_vars(tvars, std::move(typ)), false);
      decl.pos = make_pos(id.pos, typ.pos);
      return decl;
    } else {
      ir::IrDecl decl = ir::new_name_decl(id.value, kind, false);
      decl.pos = id.pos;
      return decl;
    }
  }

  antlrcpp::Any visitTypeAbstraction(bapelParser::TypeAbstractionContext *ctx) override {
    std::vector<ir::TypeParam> tvars;
    for (auto* tvar : ctx->boundedTvar()) {
      tvars.push_back(tvar->accept(this).as<ir::TypeParam>());
    }
    return tvars;
  }

  antlrcpp::Any visitBoundedTvar(bapelParser::BoundedTvarContext *ctx) override {
    ir::TypeParam tvar = ctx->tvar()->accept(this).as<ir::TypeParam>();
    if (ctx->traitBound()) {
      tvar.bounds = ctx->traitBound()->accept(this).as<std::vector<ir::IrType>>();
    }
    return tvar;
  }

  antlrcpp::Any visitTraitBound(bapelParser::TraitBoundContext *ctx) override {
    std::vector<ir::IrType> bounds;
    for (auto* typ : ctx->type_()) {
      bounds.push_back(typ->accept(this).as<ir::IrType>());
    }
    return bounds;
  }

  antlrcpp::Any visitTvar(bapelParser::TvarContext *ctx) override {
    std::string id_text = ctx->IDENTIFIER()->getText();
    return ir::TypeParam{std::move(id_text), ir::new_type_kind(), {}};
  }

  antlrcpp::Any visitType_(bapelParser::Type_Context *ctx) override {
    return ctx->forallType()->accept(this);
  }

  antlrcpp::Any visitForallType(bapelParser::ForallTypeContext *ctx) override {
    if (ctx->typeAbstraction()) {
      std::vector<ir::TypeParam> tvars = ctx->typeAbstraction()->accept(this).as<std::vector<ir::TypeParam>>();
      ir::IrType sub_type = ctx->functionType()->accept(this).as<ir::IrType>();
      ir::IrType forall = ir::forall_vars(tvars, std::move(sub_type));
      forall.pos = pos_from_context(filename_, ctx);
      return forall;
    }
    return ctx->functionType()->accept(this);
  }

  antlrcpp::Any visitFunctionType(bapelParser::FunctionTypeContext *ctx) override {
    if (ctx->functionType()) {
      ir::IrType arg = ctx->ptrType()->accept(this).as<ir::IrType>();
      ir::IrType ret = ctx->functionType()->accept(this).as<ir::IrType>();
      Pos pos = make_pos(arg.pos, ret.pos);
      ir::IrType typ = ir::new_function_type(std::move(arg), std::move(ret));
      typ.pos = pos;
      return typ;
    }
    return ctx->ptrType()->accept(this);
  }

  antlrcpp::Any visitAppType(bapelParser::AppTypeContext *ctx) override {
    if (ctx->appType()) {
      ir::IrType fun = ctx->appType()->accept(this).as<ir::IrType>();
      ir::IrType arg = ctx->primaryType()->accept(this).as<ir::IrType>();
      Pos pos = make_pos(fun.pos, arg.pos);
      ir::IrType typ = ir::new_app_type(std::move(fun), std::move(arg));
      typ.pos = pos;
      return typ;
    }
    return ctx->primaryType()->accept(this);
  }

  antlrcpp::Any visitPtrType(bapelParser::PtrTypeContext *ctx) override {
    if (ctx->AMP()) {
      ID id = new_id("Ptr", pos_from_token(filename_, ctx->AMP()->getSymbol()));
      ir::IrType name_typ = ir::new_name_type(id.value);
      name_typ.pos = id.pos;
      ir::IrType typ = ctx->ptrType()->accept(this).as<ir::IrType>();
      Pos pos = make_pos(id.pos, typ.pos);
      ir::IrType res = ir::new_app_type(std::move(name_typ), std::move(typ));
      res.pos = pos;
      return res;
    }
    return ctx->appType()->accept(this);
  }

  antlrcpp::Any visitPrimaryType(bapelParser::PrimaryTypeContext *ctx) override {
    if (ctx->arrayType()) return ctx->arrayType()->accept(this);
    if (ctx->structType()) return ctx->structType()->accept(this);
    if (ctx->tupleType()) return ctx->tupleType()->accept(this);
    if (ctx->variantType()) return ctx->variantType()->accept(this);
    if (ctx->SINGLE_QUOTE()) {
      antlr4::Token* id_token = ctx->IDENTIFIER()->getSymbol();
      ID id = new_id(id_token->getText(), pos_from_token(filename_, id_token));
      ir::IrType typ = ir::new_var_type(id.value);
      typ.pos = id.pos;
      return typ;
    }
    if (ctx->id()) {
      ID id = ctx->id()->accept(this).as<ID>();
      ir::IrType typ = ir::new_name_type(id.value);
      typ.pos = id.pos;
      return typ;
    }
    return ctx->type_()->accept(this);
  }

  antlrcpp::Any visitArrayType(bapelParser::ArrayTypeContext *ctx) override {
    ir::IrType elem_type = ctx->type_()->accept(this).as<ir::IrType>();
    int64_t length = LLONG_MAX;
    if (ctx->INT_LITERAL()) {
      length = parse_number(ctx->INT_LITERAL()->getText());
      if (ctx->MINUS()) {
        length = -length;
      }
    }
    ir::IrType typ = ir::new_array_type(std::move(elem_type), length);
    typ.pos = pos_from_context(filename_, ctx);
    return typ;
  }

  antlrcpp::Any visitStructType(bapelParser::StructTypeContext *ctx) override {
    std::vector<ir::StructField> fields;
    if (ctx->fields()) {
      fields = ctx->fields()->accept(this).as<std::vector<ir::StructField>>();
    }
    ir::IrType typ = ir::new_struct_type(std::move(fields));
    typ.pos = pos_from_context(filename_, ctx);
    return typ;
  }

  antlrcpp::Any visitFields(bapelParser::FieldsContext *ctx) override {
    std::vector<ir::StructField> fields;
    for (auto* field : ctx->field()) {
      fields.push_back(field->accept(this).as<ir::StructField>());
    }
    return fields;
  }

  antlrcpp::Any visitField(bapelParser::FieldContext *ctx) override {
    ID id = ctx->id()->accept(this).as<ID>();
    ir::IrType typ = ctx->type_()->accept(this).as<ir::IrType>();
    return ir::StructField{std::move(id.value), std::make_shared<ir::IrType>(std::move(typ))};
  }

  antlrcpp::Any visitTupleType(bapelParser::TupleTypeContext *ctx) override {
    std::vector<ir::IrType> elems;
    if (ctx->tupleTypeArgs()) {
      elems = ctx->tupleTypeArgs()->accept(this).as<std::vector<ir::IrType>>();
    }
    ir::IrType typ = ir::new_tuple_type(std::move(elems));
    typ.pos = pos_from_context(filename_, ctx);
    return typ;
  }

  antlrcpp::Any visitTupleTypeArgs(bapelParser::TupleTypeArgsContext *ctx) override {
    std::vector<ir::IrType> elems;
    for (auto* typ : ctx->type_()) {
      elems.push_back(typ->accept(this).as<ir::IrType>());
    }
    return elems;
  }

  antlrcpp::Any visitVariantType(bapelParser::VariantTypeContext *ctx) override {
    std::vector<ir::VariantTag> tags;
    if (ctx->tags()) {
      tags = ctx->tags()->accept(this).as<std::vector<ir::VariantTag>>();
    }
    ir::IrType typ = ir::new_variant_type(std::move(tags));
    typ.pos = pos_from_context(filename_, ctx);
    return typ;
  }

  antlrcpp::Any visitTags(bapelParser::TagsContext *ctx) override {
    std::vector<ir::VariantTag> tags;
    for (auto* tag : ctx->tag()) {
      tags.push_back(tag->accept(this).as<ir::VariantTag>());
    }
    return tags;
  }

  antlrcpp::Any visitTag(bapelParser::TagContext *ctx) override {
    ID id = ctx->id()->accept(this).as<ID>();
    ir::IrType typ = ctx->type_()->accept(this).as<ir::IrType>();
    return ir::VariantTag{std::move(id.value), std::make_shared<ir::IrType>(std::move(typ))};
  }

  antlrcpp::Any visitExpression(bapelParser::ExpressionContext *ctx) override {
    if (ctx->expressionWithoutBlock()) {
      return ctx->expressionWithoutBlock()->accept(this);
    }
    return ctx->expressionWithBlock()->accept(this);
  }

  antlrcpp::Any visitExpressionWithoutBlock(bapelParser::ExpressionWithoutBlockContext *ctx) override {
    if (ctx->assignTerm()) return ctx->assignTerm()->accept(this);
    if (ctx->operatorExpr()) return ctx->operatorExpr()->accept(this);
    return ctx->returnTerm()->accept(this);
  }

  antlrcpp::Any visitExpressionWithBlock(bapelParser::ExpressionWithBlockContext *ctx) override {
    if (ctx->blockExpr()) return ctx->blockExpr()->accept(this);
    if (ctx->ifTerm()) return ctx->ifTerm()->accept(this);
    if (ctx->forTerm()) return ctx->forTerm()->accept(this);
    if (ctx->lambdaTerm()) return ctx->lambdaTerm()->accept(this);
    if (ctx->matchTerm()) return ctx->matchTerm()->accept(this);
    return ctx->setTerm()->accept(this);
  }

  antlrcpp::Any visitAssignTerm(bapelParser::AssignTermContext *ctx) override {
    Expr ret;
    if (ctx->id()) {
      ID id = ctx->id()->accept(this).as<ID>();
      ret = new_var_expr(id.pos, id.value);
    } else {
      ret = ctx->tupleExpr()->accept(this).as<Expr>();
    }
    Expr arg = ctx->expression()->accept(this).as<Expr>();
    Pos pos = make_pos(arg.pos, ret.pos);
    return new_assign_expr(pos, std::move(arg), std::move(ret));
  }

  antlrcpp::Any visitReturnTerm(bapelParser::ReturnTermContext *ctx) override {
    Expr expr = ctx->expressionWithoutBlock()->accept(this).as<Expr>();
    return new_return_expr(pos_from_context(filename_, ctx), std::move(expr));
  }

  antlrcpp::Any visitIfTerm(bapelParser::IfTermContext *ctx) override {
    Expr condition = ctx->expressionWithoutBlock()->accept(this).as<Expr>();
    Expr then_branch = ctx->blockExpr(0)->accept(this).as<Expr>();
    Pos pos = pos_from_context(filename_, ctx);
    if (ctx->blockExpr().size() > 1) {
      Expr else_branch = ctx->blockExpr(1)->accept(this).as<Expr>();
      std::vector<Expr> args = {std::move(condition), std::move(then_branch), std::move(else_branch)};
      return new_app_term_expr(pos, new_var_expr(pos, "ifelse"), new_tuple_expr(pos, std::move(args)));
    } else if (ctx->ifTerm()) {
      Expr else_branch = ctx->ifTerm()->accept(this).as<Expr>();
      std::vector<Expr> args = {std::move(condition), std::move(then_branch), std::move(else_branch)};
      return new_app_term_expr(pos, new_var_expr(pos, "ifelse"), new_tuple_expr(pos, std::move(args)));
    }
    std::vector<Expr> args = {std::move(condition), std::move(then_branch)};
    return new_app_term_expr(pos, new_var_expr(pos, "ifthen"), new_tuple_expr(pos, std::move(args)));
  }

  antlrcpp::Any visitForTerm(bapelParser::ForTermContext *ctx) override {
    Expr condition = ctx->expressionWithoutBlock()->accept(this).as<Expr>();
    Expr body = ctx->blockExpr()->accept(this).as<Expr>();
    return new_for_expr(pos_from_context(filename_, ctx), std::move(condition), std::move(body));
  }

  antlrcpp::Any visitLambdaTerm(bapelParser::LambdaTermContext *ctx) override {
    std::vector<ir::TypeParam> tvars;
    if (ctx->typeAbstraction()) {
      tvars = ctx->typeAbstraction()->accept(this).as<std::vector<ir::TypeParam>>();
    }
    std::vector<ir::FunctionArg> fun_args = ctx->functionArgs()->accept(this).as<std::vector<ir::FunctionArg>>();
    Expr body = ctx->blockExpr()->accept(this).as<Expr>();
    return lambda_dsl(pos_from_context(filename_, ctx), tvars, fun_args, std::move(body));
  }

  antlrcpp::Any visitMatchTerm(bapelParser::MatchTermContext *ctx) override {
    Expr expr = ctx->expression()->accept(this).as<Expr>();
    std::vector<MatchArm> arms = ctx->matchArms()->accept(this).as<std::vector<MatchArm>>();
    return new_match_expr(pos_from_context(filename_, ctx), std::move(expr), std::move(arms));
  }

  antlrcpp::Any visitMatchArms(bapelParser::MatchArmsContext *ctx) override {
    std::vector<MatchArm> arms;
    for (auto* arm : ctx->matchArm()) {
      arms.push_back(arm->accept(this).as<MatchArm>());
    }
    return arms;
  }

  antlrcpp::Any visitMatchArm(bapelParser::MatchArmContext *ctx) override {
    ID tag = ctx->id()->accept(this).as<ID>();
    std::string arg = ctx->IDENTIFIER()->getText();
    Expr body = ctx->expression()->accept(this).as<Expr>();
    MatchArm arm;
    arm.tag = std::move(tag.value);
    arm.arg = std::move(arg);
    arm.body = std::make_shared<Expr>(std::move(body));
    return arm;
  }

  antlrcpp::Any visitSetTerm(bapelParser::SetTermContext *ctx) override {
    Expr expr = ctx->expression()->accept(this).as<Expr>();
    std::vector<LabelValue> values = ctx->labelValues()->accept(this).as<std::vector<LabelValue>>();
    return new_set_expr(pos_from_context(filename_, ctx), std::move(expr), std::move(values));
  }

  antlrcpp::Any visitBlockExpr(bapelParser::BlockExprContext *ctx) override {
    std::vector<Expr> exprs = ctx->blockStatements()->accept(this).as<std::vector<Expr>>();
    return new_block_expr(pos_from_context(filename_, ctx), std::move(exprs));
  }

  antlrcpp::Any visitBlockStatements(bapelParser::BlockStatementsContext *ctx) override {
    std::vector<Expr> exprs;
    if (ctx->statements()) {
      exprs = ctx->statements()->accept(this).as<std::vector<Expr>>();
      if (ctx->expressionWithoutBlock()) {
        exprs.push_back(ctx->expressionWithoutBlock()->accept(this).as<Expr>());
      }
      return exprs;
    }
    exprs.push_back(ctx->expressionWithoutBlock()->accept(this).as<Expr>());
    return exprs;
  }

  antlrcpp::Any visitStatements(bapelParser::StatementsContext *ctx) override {
    std::vector<Expr> exprs;
    for (auto* stmt : ctx->statement()) {
      exprs.push_back(stmt->accept(this).as<Expr>());
    }
    return exprs;
  }

  antlrcpp::Any visitStatement(bapelParser::StatementContext *ctx) override {
    if (ctx->letStatement()) return ctx->letStatement()->accept(this);
    return ctx->expressionStatement()->accept(this);
  }

  antlrcpp::Any visitLetStatement(bapelParser::LetStatementContext *ctx) override {
    ID id = ctx->id()->accept(this).as<ID>();
    Expr value = ctx->expression()->accept(this).as<Expr>();
    std::optional<ir::IrType> var_type = std::nullopt;
    if (ctx->type_()) {
      var_type = ctx->type_()->accept(this).as<ir::IrType>();
    }
    Pos pos = make_pos(id.pos, value.pos);
    return new_let_expr(pos, id.value, std::move(var_type), std::move(value));
  }

  antlrcpp::Any visitExpressionStatement(bapelParser::ExpressionStatementContext *ctx) override {
    if (ctx->expressionWithoutBlock()) return ctx->expressionWithoutBlock()->accept(this);
    return ctx->expressionWithBlock()->accept(this);
  }

  antlrcpp::Any visitOperatorExpr(bapelParser::OperatorExprContext *ctx) override {
    return ctx->logicalOrExpr()->accept(this);
  }

  antlrcpp::Any visitLogicalOrExpr(bapelParser::LogicalOrExprContext *ctx) override {
    if (ctx->logicalOrExpr()) {
      Expr left = ctx->logicalOrExpr()->accept(this).as<Expr>();
      Expr right = ctx->logicalAndExpr()->accept(this).as<Expr>();
      antlr4::Token* op_tok = ctx->OR()->getSymbol();
      Expr op_expr = new_var_expr(pos_from_token(filename_, op_tok), op_tok->getText());
      return new_bin_op_expr(std::move(op_expr), {}, std::move(left), std::move(right));
    }
    return ctx->logicalAndExpr()->accept(this);
  }

  antlrcpp::Any visitLogicalAndExpr(bapelParser::LogicalAndExprContext *ctx) override {
    if (ctx->logicalAndExpr()) {
      Expr left = ctx->logicalAndExpr()->accept(this).as<Expr>();
      Expr right = ctx->equalityExpr()->accept(this).as<Expr>();
      antlr4::Token* op_tok = ctx->AND()->getSymbol();
      Expr op_expr = new_var_expr(pos_from_token(filename_, op_tok), op_tok->getText());
      return new_bin_op_expr(std::move(op_expr), {}, std::move(left), std::move(right));
    }
    return ctx->equalityExpr()->accept(this);
  }

  antlrcpp::Any visitEqualityExpr(bapelParser::EqualityExprContext *ctx) override {
    if (ctx->equalityExpr()) {
      Expr left = ctx->equalityExpr()->accept(this).as<Expr>();
      Expr right = ctx->comparisonExpr()->accept(this).as<Expr>();
      antlr4::Token* op_tok = ctx->NE() ? ctx->NE()->getSymbol() : ctx->EQ()->getSymbol();
      Expr op_expr = new_var_expr(pos_from_token(filename_, op_tok), op_tok->getText());
      std::vector<ir::IrType> type_args;
      if (ctx->typeApplicativeArgs()) {
        type_args = ctx->typeApplicativeArgs()->accept(this).as<std::vector<ir::IrType>>();
      }
      return new_bin_op_expr(std::move(op_expr), type_args, std::move(left), std::move(right));
    }
    return ctx->comparisonExpr()->accept(this);
  }

  antlrcpp::Any visitComparisonExpr(bapelParser::ComparisonExprContext *ctx) override {
    if (ctx->comparisonExpr()) {
      Expr left = ctx->comparisonExpr()->accept(this).as<Expr>();
      Expr right = ctx->additiveExpr()->accept(this).as<Expr>();
      antlr4::Token* op_tok = nullptr;
      if (ctx->GT()) op_tok = ctx->GT()->getSymbol();
      else if (ctx->GE()) op_tok = ctx->GE()->getSymbol();
      else if (ctx->LT()) op_tok = ctx->LT()->getSymbol();
      else op_tok = ctx->LE()->getSymbol();
      Expr op_expr = new_var_expr(pos_from_token(filename_, op_tok), op_tok->getText());
      std::vector<ir::IrType> type_args;
      if (ctx->typeApplicativeArgs()) {
        type_args = ctx->typeApplicativeArgs()->accept(this).as<std::vector<ir::IrType>>();
      }
      return new_bin_op_expr(std::move(op_expr), type_args, std::move(left), std::move(right));
    }
    return ctx->additiveExpr()->accept(this);
  }

  antlrcpp::Any visitAdditiveExpr(bapelParser::AdditiveExprContext *ctx) override {
    if (ctx->additiveExpr()) {
      Expr left = ctx->additiveExpr()->accept(this).as<Expr>();
      Expr right = ctx->multiplicativeExpr()->accept(this).as<Expr>();
      antlr4::Token* op_tok = ctx->PLUS() ? ctx->PLUS()->getSymbol() : ctx->MINUS()->getSymbol();
      Expr op_expr = new_var_expr(pos_from_token(filename_, op_tok), op_tok->getText());
      std::vector<ir::IrType> type_args;
      if (ctx->typeApplicativeArgs()) {
        type_args = ctx->typeApplicativeArgs()->accept(this).as<std::vector<ir::IrType>>();
      }
      return new_bin_op_expr(std::move(op_expr), type_args, std::move(left), std::move(right));
    }
    return ctx->multiplicativeExpr()->accept(this);
  }

  antlrcpp::Any visitMultiplicativeExpr(bapelParser::MultiplicativeExprContext *ctx) override {
    if (ctx->multiplicativeExpr()) {
      Expr left = ctx->multiplicativeExpr()->accept(this).as<Expr>();
      Expr right = ctx->unaryExpr()->accept(this).as<Expr>();
      antlr4::Token* op_tok = ctx->MUL() ? ctx->MUL()->getSymbol() : ctx->DIV()->getSymbol();
      Expr op_expr = new_var_expr(pos_from_token(filename_, op_tok), op_tok->getText());
      std::vector<ir::IrType> type_args;
      if (ctx->typeApplicativeArgs()) {
        type_args = ctx->typeApplicativeArgs()->accept(this).as<std::vector<ir::IrType>>();
      }
      return new_bin_op_expr(std::move(op_expr), type_args, std::move(left), std::move(right));
    }
    return ctx->unaryExpr()->accept(this);
  }

  antlrcpp::Any visitUnaryExpr(bapelParser::UnaryExprContext *ctx) override {
    if (ctx->unaryExpr()) {
      Expr operand = ctx->unaryExpr()->accept(this).as<Expr>();
      antlr4::Token* op_tok = ctx->NOT() ? ctx->NOT()->getSymbol() : ctx->MINUS()->getSymbol();
      Expr op_expr = new_var_expr(pos_from_token(filename_, op_tok), op_tok->getText());
      std::vector<ir::IrType> type_args;
      if (ctx->typeApplicativeArgs()) {
        type_args = ctx->typeApplicativeArgs()->accept(this).as<std::vector<ir::IrType>>();
      }
      return new_unary_op_expr(std::move(op_expr), type_args, std::move(operand));
    }
    return ctx->applicativeExpr()->accept(this);
  }

  antlrcpp::Any visitApplicativeExpr(bapelParser::ApplicativeExprContext *ctx) override {
    if (ctx->applicativeExpr()) {
      Expr fun = ctx->applicativeExpr()->accept(this).as<Expr>();
      Expr arg = ctx->basePrimaryExpr()->accept(this).as<Expr>();
      Pos pos = make_pos(fun.pos, arg.pos);
      return new_app_term_expr(pos, std::move(fun), std::move(arg));
    }
    return ctx->typeApplicativeExpr()->accept(this);
  }

  antlrcpp::Any visitTypeApplicativeExpr(bapelParser::TypeApplicativeExprContext *ctx) override {
    Expr primary = ctx->primaryExpr()->accept(this).as<Expr>();
    if (ctx->typeApplicativeArgs()) {
      std::vector<ir::IrType> type_args = ctx->typeApplicativeArgs()->accept(this).as<std::vector<ir::IrType>>();
      Pos pos = primary.pos;
      for (auto& typ : type_args) {
        primary = new_app_type_expr(pos, std::move(primary), std::move(typ));
      }
    }
    return primary;
  }

  antlrcpp::Any visitTypeApplicativeArgs(bapelParser::TypeApplicativeArgsContext *ctx) override {
    if (ctx->tupleTypeArgs()) {
      return ctx->tupleTypeArgs()->accept(this);
    }
    return std::vector<ir::IrType>{ctx->type_()->accept(this).as<ir::IrType>()};
  }

  antlrcpp::Any visitPrimaryExpr(bapelParser::PrimaryExprContext *ctx) override {
    if (ctx->MUL()) {
      Pos mul_pos = pos_from_token(filename_, ctx->MUL()->getSymbol());
      Expr id = new_var_expr(mul_pos, "Ptr::get");
      Expr expr = ctx->primaryExpr()->accept(this).as<Expr>();
      Pos pos = make_pos(id.pos, expr.pos);
      std::vector<Expr> args;
      args.push_back(std::move(expr));
      return call(pos, std::move(id), {}, std::move(args));
    }
    if (ctx->basePrimaryExpr()) {
      return ctx->basePrimaryExpr()->accept(this);
    }
    return Expr{};
  }

  antlrcpp::Any visitBasePrimaryExpr(bapelParser::BasePrimaryExprContext *ctx) override {
    if (ctx->AMP()) {
      Pos amp_pos = pos_from_token(filename_, ctx->AMP()->getSymbol());
      Expr id = new_var_expr(amp_pos, "Ptr::mk");
      Expr target = ctx->projectionExpr()->accept(this).as<Expr>();
      Pos pos = make_pos(id.pos, target.pos);
      std::vector<Expr> args;
      args.push_back(std::move(target));
      return call(pos, std::move(id), {}, std::move(args));
    }
    if (ctx->projectionExpr()) {
      return ctx->projectionExpr()->accept(this);
    }
    if (ctx->INT_LITERAL()) {
      antlr4::Token* tok = ctx->INT_LITERAL()->getSymbol();
      Pos pos = pos_from_token(filename_, tok);
      return new_const_expr(pos, new_number_literal(pos, tok->getText()));
    }
    if (ctx->FLOAT_LITERAL()) {
      antlr4::Token* tok = ctx->FLOAT_LITERAL()->getSymbol();
      Pos pos = pos_from_token(filename_, tok);
      return new_const_expr(pos, new_number_literal(pos, tok->getText()));
    }
    return Expr{};
  }

  antlrcpp::Any visitProjectionExpr(bapelParser::ProjectionExprContext *ctx) override {
    if (ctx->projectionExpr()) {
      Expr expr = ctx->projectionExpr()->accept(this).as<Expr>();
      std::string label;
      Pos label_pos;
      if (ctx->INT_LITERAL()) {
        antlr4::Token* tok = ctx->INT_LITERAL()->getSymbol();
        label = tok->getText();
        label_pos = pos_from_token(filename_, tok);
      } else {
        antlr4::Token* tok = ctx->IDENTIFIER()->getSymbol();
        label = tok->getText();
        label_pos = pos_from_token(filename_, tok);
      }
      Pos pos = make_pos(expr.pos, label_pos);
      return new_projection_expr(pos, std::move(expr), std::move(label));
    }
    return ctx->derefExpr()->accept(this);
  }

  antlrcpp::Any visitDerefExpr(bapelParser::DerefExprContext *ctx) override {
    if (ctx->injectionExpr()) return ctx->injectionExpr()->accept(this);
    if (ctx->RUNE_LITERAL()) {
      antlr4::Token* tok = ctx->RUNE_LITERAL()->getSymbol();
      Pos pos = pos_from_token(filename_, tok);
      return new_const_expr(pos, parse_rune_literal(pos, tok->getText()));
    }
    if (ctx->STRING_LITERAL()) {
      antlr4::Token* tok = ctx->STRING_LITERAL()->getSymbol();
      Pos pos = pos_from_token(filename_, tok);
      return new_const_expr(pos, parse_string_literal(pos, tok->getText()));
    }
    if (ctx->structExpr()) return ctx->structExpr()->accept(this);
    if (ctx->tupleExpr()) return ctx->tupleExpr()->accept(this);
    if (ctx->varExpr()) return ctx->varExpr()->accept(this);
    return ctx->expression()->accept(this);
  }

  antlrcpp::Any visitInjectionExpr(bapelParser::InjectionExprContext *ctx) override {
    ir::IrType variant_type = ctx->type_()->accept(this).as<ir::IrType>();
    LabelValue lv = ctx->labelValue()->accept(this).as<LabelValue>();
    return new_injection_expr(pos_from_context(filename_, ctx), std::move(variant_type), std::move(lv.label), *lv.value);
  }

  antlrcpp::Any visitStructExpr(bapelParser::StructExprContext *ctx) override {
    std::vector<LabelValue> values;
    if (ctx->labelValues()) {
      values = ctx->labelValues()->accept(this).as<std::vector<LabelValue>>();
    }
    return new_struct_expr(pos_from_context(filename_, ctx), std::move(values));
  }

  antlrcpp::Any visitLabelValues(bapelParser::LabelValuesContext *ctx) override {
    std::vector<LabelValue> values;
    for (auto* lv : ctx->labelValue()) {
      values.push_back(lv->accept(this).as<LabelValue>());
    }
    return values;
  }

  antlrcpp::Any visitLabelValue(bapelParser::LabelValueContext *ctx) override {
    Expr value = ctx->expression()->accept(this).as<Expr>();
    if (ctx->id()) {
      ID id = ctx->id()->accept(this).as<ID>();
      LabelValue lv;
      lv.label = std::move(id.value);
      lv.value = std::make_shared<Expr>(std::move(value));
      return lv;
    }
    std::string label_val = ctx->INT_LITERAL()->getText();
    LabelValue lv;
    lv.label = std::move(label_val);
    lv.value = std::make_shared<Expr>(std::move(value));
    return lv;
  }

  antlrcpp::Any visitTupleExpr(bapelParser::TupleExprContext *ctx) override {
    std::vector<Expr> elems;
    if (ctx->tupleExprArgs()) {
      elems = ctx->tupleExprArgs()->accept(this).as<std::vector<Expr>>();
    }
    return new_tuple_expr(pos_from_context(filename_, ctx), std::move(elems));
  }

  antlrcpp::Any visitTupleExprArgs(bapelParser::TupleExprArgsContext *ctx) override {
    std::vector<Expr> elems;
    for (auto* expr : ctx->expression()) {
      elems.push_back(expr->accept(this).as<Expr>());
    }
    return elems;
  }

  antlrcpp::Any visitVarExpr(bapelParser::VarExprContext *ctx) override {
    ID id = ctx->id()->accept(this).as<ID>();
    return new_var_expr(id.pos, std::move(id.value));
  }

  antlrcpp::Any visitId(bapelParser::IdContext *ctx) override {
    if (ctx->idTokens()) {
      return ctx->idTokens()->accept(this);
    }
    auto* op_node = dynamic_cast<antlr4::tree::TerminalNode*>(ctx->children[1]);
    antlr4::Token* op_tok = op_node ? op_node->getSymbol() : ctx->getStart();
    return new_id(op_tok->getText(), pos_from_token(filename_, op_tok));
  }

  antlrcpp::Any visitIdTokens(bapelParser::IdTokensContext *ctx) override {
    return new_id(ctx->getText(), pos_from_context(filename_, ctx));
  }

  antlrcpp::Any visitTraitDecl(bapelParser::TraitDeclContext *ctx) override {
    ID id = ctx->id()->accept(this).as<ID>();
    std::vector<ir::TypeParam> tvars;
    if (ctx->typeAbstraction()) {
      tvars = ctx->typeAbstraction()->accept(this).as<std::vector<ir::TypeParam>>();
    }
    std::vector<Signature> methods;
    for (auto* m : ctx->traitMethod()) {
      methods.push_back(m->accept(this).as<Signature>());
    }
    Trait t;
    t.export_flag = false;
    t.id = std::move(id.value);
    t.type_params = std::move(tvars);
    t.methods = std::move(methods);
    t.pos = pos_from_context(filename_, ctx);
    return t;
  }

  antlrcpp::Any visitTraitMethod(bapelParser::TraitMethodContext *ctx) override {
    ID id = ctx->id()->accept(this).as<ID>();
    std::vector<ir::FunctionArg> fun_args = ctx->functionArgs()->accept(this).as<std::vector<ir::FunctionArg>>();
    ir::IrType ret_type = ctx->type_()->accept(this).as<ir::IrType>();
    Signature sig;
    sig.id = std::move(id.value);
    sig.args = std::move(fun_args);
    sig.ret_type = std::move(ret_type);
    sig.pos = pos_from_context(filename_, ctx);
    return sig;
  }

  antlrcpp::Any visitTraitImpl(bapelParser::TraitImplContext *ctx) override {
    std::vector<ir::TypeParam> tvars;
    if (ctx->typeAbstraction()) {
      tvars = ctx->typeAbstraction()->accept(this).as<std::vector<ir::TypeParam>>();
    }
    ir::IrType trait_type = ctx->type_(0)->accept(this).as<ir::IrType>();
    ir::IrType target_type = ctx->type_(1)->accept(this).as<ir::IrType>();
    std::vector<Function> methods;
    for (auto* fn_ctx : ctx->functionNoExport()) {
      Source fn_src = fn_ctx->accept(this).as<Source>();
      if (fn_src.function_data) {
        methods.push_back(*fn_src.function_data);
      }
    }
    Impl impl;
    impl.case_val = ImplCase::TraitImpl;
    impl.type_params = std::move(tvars);
    impl.trait_type = std::move(trait_type);
    impl.type_name = std::move(target_type);
    impl.methods = std::move(methods);
    impl.pos = pos_from_context(filename_, ctx);
    return new_impl_source(std::move(impl));
  }

  antlrcpp::Any visitInherentImpl(bapelParser::InherentImplContext *ctx) override {
    std::vector<ir::TypeParam> tvars;
    if (ctx->typeAbstraction()) {
      tvars = ctx->typeAbstraction()->accept(this).as<std::vector<ir::TypeParam>>();
    }
    ir::IrType target_type = ctx->type_()->accept(this).as<ir::IrType>();
    std::vector<Function> methods;
    for (auto* fn_ctx : ctx->functionNoExport()) {
      Source fn_src = fn_ctx->accept(this).as<Source>();
      if (fn_src.function_data) {
        methods.push_back(*fn_src.function_data);
      }
    }
    Impl impl;
    impl.case_val = ImplCase::InherentImpl;
    impl.type_params = std::move(tvars);
    impl.type_name = std::move(target_type);
    impl.methods = std::move(methods);
    impl.pos = pos_from_context(filename_, ctx);
    return new_impl_source(std::move(impl));
  }

private:
  std::string filename_;
};

} // namespace ast
