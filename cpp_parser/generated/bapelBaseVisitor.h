
// Generated from cpp_parser/bapel.g4 by ANTLR 4.9.2

#pragma once


#include "antlr4-runtime.h"
#include "bapelVisitor.h"


/**
 * This class provides an empty implementation of bapelVisitor, which can be
 * extended to create a visitor which only needs to handle a subset of the available methods.
 */
class  bapelBaseVisitor : public bapelVisitor {
public:

  virtual antlrcpp::Any visitBaseSourceFile(bapelParser::BaseSourceFileContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitImplSourceFile(bapelParser::ImplSourceFileContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitModuleHeader(bapelParser::ModuleHeaderContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitImplementsHeader(bapelParser::ImplementsHeaderContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitWorkspace(bapelParser::WorkspaceContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitPackagesSection(bapelParser::PackagesSectionContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitPackageRule(bapelParser::PackageRuleContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitImportsSection(bapelParser::ImportsSectionContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitImplsSection(bapelParser::ImplsSectionContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitFlagsSection(bapelParser::FlagsSectionContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitModuleID(bapelParser::ModuleIDContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitFilename(bapelParser::FilenameContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitSources(bapelParser::SourcesContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitSource(bapelParser::SourceContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitTraitDecl(bapelParser::TraitDeclContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitTraitMethod(bapelParser::TraitMethodContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitTraitImpl(bapelParser::TraitImplContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitInherentImpl(bapelParser::InherentImplContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitDeclNoExport(bapelParser::DeclNoExportContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitFunctionNoExport(bapelParser::FunctionNoExportContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitFunctionArgs(bapelParser::FunctionArgsContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitArg(bapelParser::ArgContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitDecl(bapelParser::DeclContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitUnexportedDecl(bapelParser::UnexportedDeclContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitDeclNoTerm(bapelParser::DeclNoTermContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitTermDecl(bapelParser::TermDeclContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitTypeDecl(bapelParser::TypeDeclContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitTypeAbstraction(bapelParser::TypeAbstractionContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitBoundedTvar(bapelParser::BoundedTvarContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitTvar(bapelParser::TvarContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitTraitBound(bapelParser::TraitBoundContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitType_(bapelParser::Type_Context *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitForallType(bapelParser::ForallTypeContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitFunctionType(bapelParser::FunctionTypeContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitPtrType(bapelParser::PtrTypeContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitAppType(bapelParser::AppTypeContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitPrimaryType(bapelParser::PrimaryTypeContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitArrayType(bapelParser::ArrayTypeContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitStructType(bapelParser::StructTypeContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitFields(bapelParser::FieldsContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitField(bapelParser::FieldContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitTupleType(bapelParser::TupleTypeContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitTupleTypeArgs(bapelParser::TupleTypeArgsContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitVariantType(bapelParser::VariantTypeContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitTags(bapelParser::TagsContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitTag(bapelParser::TagContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitExpression(bapelParser::ExpressionContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitExpressionWithoutBlock(bapelParser::ExpressionWithoutBlockContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitExpressionWithBlock(bapelParser::ExpressionWithBlockContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitAssignTerm(bapelParser::AssignTermContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitReturnTerm(bapelParser::ReturnTermContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitIfTerm(bapelParser::IfTermContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitForTerm(bapelParser::ForTermContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitLambdaTerm(bapelParser::LambdaTermContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitMatchTerm(bapelParser::MatchTermContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitMatchArms(bapelParser::MatchArmsContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitMatchArm(bapelParser::MatchArmContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitSetTerm(bapelParser::SetTermContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitBlockExpr(bapelParser::BlockExprContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitBlockStatements(bapelParser::BlockStatementsContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitStatements(bapelParser::StatementsContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitStatement(bapelParser::StatementContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitLetStatement(bapelParser::LetStatementContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitExpressionStatement(bapelParser::ExpressionStatementContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitOperatorExpr(bapelParser::OperatorExprContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitLogicalOrExpr(bapelParser::LogicalOrExprContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitLogicalAndExpr(bapelParser::LogicalAndExprContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitEqualityExpr(bapelParser::EqualityExprContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitComparisonExpr(bapelParser::ComparisonExprContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitAdditiveExpr(bapelParser::AdditiveExprContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitMultiplicativeExpr(bapelParser::MultiplicativeExprContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitUnaryExpr(bapelParser::UnaryExprContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitApplicativeExpr(bapelParser::ApplicativeExprContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitTypeApplicativeExpr(bapelParser::TypeApplicativeExprContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitTypeApplicativeArgs(bapelParser::TypeApplicativeArgsContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitBasePrimaryExpr(bapelParser::BasePrimaryExprContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitPrimaryExpr(bapelParser::PrimaryExprContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitProjectionExpr(bapelParser::ProjectionExprContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitDerefExpr(bapelParser::DerefExprContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitInjectionExpr(bapelParser::InjectionExprContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitStructExpr(bapelParser::StructExprContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitLabelValues(bapelParser::LabelValuesContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitLabelValue(bapelParser::LabelValueContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitTupleExpr(bapelParser::TupleExprContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitTupleExprArgs(bapelParser::TupleExprArgsContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitVarExpr(bapelParser::VarExprContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitId(bapelParser::IdContext *ctx) override {
    return visitChildren(ctx);
  }

  virtual antlrcpp::Any visitIdTokens(bapelParser::IdTokensContext *ctx) override {
    return visitChildren(ctx);
  }


};

