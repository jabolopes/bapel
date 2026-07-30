
// Generated from cpp_parser/bapel.g4 by ANTLR 4.9.2

#pragma once


#include "antlr4-runtime.h"
#include "bapelParser.h"



/**
 * This class defines an abstract visitor for a parse tree
 * produced by bapelParser.
 */
class  bapelVisitor : public antlr4::tree::AbstractParseTreeVisitor {
public:

  /**
   * Visit parse trees produced by bapelParser.
   */
    virtual antlrcpp::Any visitBaseSourceFile(bapelParser::BaseSourceFileContext *context) = 0;

    virtual antlrcpp::Any visitImplSourceFile(bapelParser::ImplSourceFileContext *context) = 0;

    virtual antlrcpp::Any visitModuleHeader(bapelParser::ModuleHeaderContext *context) = 0;

    virtual antlrcpp::Any visitImplementsHeader(bapelParser::ImplementsHeaderContext *context) = 0;

    virtual antlrcpp::Any visitWorkspace(bapelParser::WorkspaceContext *context) = 0;

    virtual antlrcpp::Any visitPackagesSection(bapelParser::PackagesSectionContext *context) = 0;

    virtual antlrcpp::Any visitPackageRule(bapelParser::PackageRuleContext *context) = 0;

    virtual antlrcpp::Any visitImportsSection(bapelParser::ImportsSectionContext *context) = 0;

    virtual antlrcpp::Any visitImplsSection(bapelParser::ImplsSectionContext *context) = 0;

    virtual antlrcpp::Any visitFlagsSection(bapelParser::FlagsSectionContext *context) = 0;

    virtual antlrcpp::Any visitModuleID(bapelParser::ModuleIDContext *context) = 0;

    virtual antlrcpp::Any visitFilename(bapelParser::FilenameContext *context) = 0;

    virtual antlrcpp::Any visitSources(bapelParser::SourcesContext *context) = 0;

    virtual antlrcpp::Any visitSource(bapelParser::SourceContext *context) = 0;

    virtual antlrcpp::Any visitTraitDecl(bapelParser::TraitDeclContext *context) = 0;

    virtual antlrcpp::Any visitTraitMethod(bapelParser::TraitMethodContext *context) = 0;

    virtual antlrcpp::Any visitTraitImpl(bapelParser::TraitImplContext *context) = 0;

    virtual antlrcpp::Any visitInherentImpl(bapelParser::InherentImplContext *context) = 0;

    virtual antlrcpp::Any visitDeclNoExport(bapelParser::DeclNoExportContext *context) = 0;

    virtual antlrcpp::Any visitFunctionNoExport(bapelParser::FunctionNoExportContext *context) = 0;

    virtual antlrcpp::Any visitFunctionArgs(bapelParser::FunctionArgsContext *context) = 0;

    virtual antlrcpp::Any visitArg(bapelParser::ArgContext *context) = 0;

    virtual antlrcpp::Any visitDecl(bapelParser::DeclContext *context) = 0;

    virtual antlrcpp::Any visitUnexportedDecl(bapelParser::UnexportedDeclContext *context) = 0;

    virtual antlrcpp::Any visitDeclNoTerm(bapelParser::DeclNoTermContext *context) = 0;

    virtual antlrcpp::Any visitTermDecl(bapelParser::TermDeclContext *context) = 0;

    virtual antlrcpp::Any visitTypeDecl(bapelParser::TypeDeclContext *context) = 0;

    virtual antlrcpp::Any visitTypeAbstraction(bapelParser::TypeAbstractionContext *context) = 0;

    virtual antlrcpp::Any visitBoundedTvar(bapelParser::BoundedTvarContext *context) = 0;

    virtual antlrcpp::Any visitTvar(bapelParser::TvarContext *context) = 0;

    virtual antlrcpp::Any visitTraitBound(bapelParser::TraitBoundContext *context) = 0;

    virtual antlrcpp::Any visitType_(bapelParser::Type_Context *context) = 0;

    virtual antlrcpp::Any visitForallType(bapelParser::ForallTypeContext *context) = 0;

    virtual antlrcpp::Any visitFunctionType(bapelParser::FunctionTypeContext *context) = 0;

    virtual antlrcpp::Any visitPtrType(bapelParser::PtrTypeContext *context) = 0;

    virtual antlrcpp::Any visitAppType(bapelParser::AppTypeContext *context) = 0;

    virtual antlrcpp::Any visitPrimaryType(bapelParser::PrimaryTypeContext *context) = 0;

    virtual antlrcpp::Any visitArrayType(bapelParser::ArrayTypeContext *context) = 0;

    virtual antlrcpp::Any visitStructType(bapelParser::StructTypeContext *context) = 0;

    virtual antlrcpp::Any visitFields(bapelParser::FieldsContext *context) = 0;

    virtual antlrcpp::Any visitField(bapelParser::FieldContext *context) = 0;

    virtual antlrcpp::Any visitTupleType(bapelParser::TupleTypeContext *context) = 0;

    virtual antlrcpp::Any visitTupleTypeArgs(bapelParser::TupleTypeArgsContext *context) = 0;

    virtual antlrcpp::Any visitVariantType(bapelParser::VariantTypeContext *context) = 0;

    virtual antlrcpp::Any visitTags(bapelParser::TagsContext *context) = 0;

    virtual antlrcpp::Any visitTag(bapelParser::TagContext *context) = 0;

    virtual antlrcpp::Any visitExpression(bapelParser::ExpressionContext *context) = 0;

    virtual antlrcpp::Any visitExpressionWithoutBlock(bapelParser::ExpressionWithoutBlockContext *context) = 0;

    virtual antlrcpp::Any visitExpressionWithBlock(bapelParser::ExpressionWithBlockContext *context) = 0;

    virtual antlrcpp::Any visitAssignTerm(bapelParser::AssignTermContext *context) = 0;

    virtual antlrcpp::Any visitReturnTerm(bapelParser::ReturnTermContext *context) = 0;

    virtual antlrcpp::Any visitIfTerm(bapelParser::IfTermContext *context) = 0;

    virtual antlrcpp::Any visitForTerm(bapelParser::ForTermContext *context) = 0;

    virtual antlrcpp::Any visitLambdaTerm(bapelParser::LambdaTermContext *context) = 0;

    virtual antlrcpp::Any visitMatchTerm(bapelParser::MatchTermContext *context) = 0;

    virtual antlrcpp::Any visitMatchArms(bapelParser::MatchArmsContext *context) = 0;

    virtual antlrcpp::Any visitMatchArm(bapelParser::MatchArmContext *context) = 0;

    virtual antlrcpp::Any visitSetTerm(bapelParser::SetTermContext *context) = 0;

    virtual antlrcpp::Any visitBlockExpr(bapelParser::BlockExprContext *context) = 0;

    virtual antlrcpp::Any visitBlockStatements(bapelParser::BlockStatementsContext *context) = 0;

    virtual antlrcpp::Any visitStatements(bapelParser::StatementsContext *context) = 0;

    virtual antlrcpp::Any visitStatement(bapelParser::StatementContext *context) = 0;

    virtual antlrcpp::Any visitLetStatement(bapelParser::LetStatementContext *context) = 0;

    virtual antlrcpp::Any visitExpressionStatement(bapelParser::ExpressionStatementContext *context) = 0;

    virtual antlrcpp::Any visitOperatorExpr(bapelParser::OperatorExprContext *context) = 0;

    virtual antlrcpp::Any visitLogicalOrExpr(bapelParser::LogicalOrExprContext *context) = 0;

    virtual antlrcpp::Any visitLogicalAndExpr(bapelParser::LogicalAndExprContext *context) = 0;

    virtual antlrcpp::Any visitEqualityExpr(bapelParser::EqualityExprContext *context) = 0;

    virtual antlrcpp::Any visitComparisonExpr(bapelParser::ComparisonExprContext *context) = 0;

    virtual antlrcpp::Any visitAdditiveExpr(bapelParser::AdditiveExprContext *context) = 0;

    virtual antlrcpp::Any visitMultiplicativeExpr(bapelParser::MultiplicativeExprContext *context) = 0;

    virtual antlrcpp::Any visitUnaryExpr(bapelParser::UnaryExprContext *context) = 0;

    virtual antlrcpp::Any visitApplicativeExpr(bapelParser::ApplicativeExprContext *context) = 0;

    virtual antlrcpp::Any visitTypeApplicativeExpr(bapelParser::TypeApplicativeExprContext *context) = 0;

    virtual antlrcpp::Any visitTypeApplicativeArgs(bapelParser::TypeApplicativeArgsContext *context) = 0;

    virtual antlrcpp::Any visitBasePrimaryExpr(bapelParser::BasePrimaryExprContext *context) = 0;

    virtual antlrcpp::Any visitPrimaryExpr(bapelParser::PrimaryExprContext *context) = 0;

    virtual antlrcpp::Any visitProjectionExpr(bapelParser::ProjectionExprContext *context) = 0;

    virtual antlrcpp::Any visitDerefExpr(bapelParser::DerefExprContext *context) = 0;

    virtual antlrcpp::Any visitInjectionExpr(bapelParser::InjectionExprContext *context) = 0;

    virtual antlrcpp::Any visitStructExpr(bapelParser::StructExprContext *context) = 0;

    virtual antlrcpp::Any visitLabelValues(bapelParser::LabelValuesContext *context) = 0;

    virtual antlrcpp::Any visitLabelValue(bapelParser::LabelValueContext *context) = 0;

    virtual antlrcpp::Any visitTupleExpr(bapelParser::TupleExprContext *context) = 0;

    virtual antlrcpp::Any visitTupleExprArgs(bapelParser::TupleExprArgsContext *context) = 0;

    virtual antlrcpp::Any visitVarExpr(bapelParser::VarExprContext *context) = 0;

    virtual antlrcpp::Any visitId(bapelParser::IdContext *context) = 0;

    virtual antlrcpp::Any visitIdTokens(bapelParser::IdTokensContext *context) = 0;


};

