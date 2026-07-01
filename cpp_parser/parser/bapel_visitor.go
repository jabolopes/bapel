// Code generated from cpp_parser/bapel.g4 by ANTLR 4.9.2. DO NOT EDIT.

package parser // bapel

import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by bapelParser.
type bapelVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by bapelParser#baseSourceFile.
	VisitBaseSourceFile(ctx *BaseSourceFileContext) interface{}

	// Visit a parse tree produced by bapelParser#implSourceFile.
	VisitImplSourceFile(ctx *ImplSourceFileContext) interface{}

	// Visit a parse tree produced by bapelParser#moduleHeader.
	VisitModuleHeader(ctx *ModuleHeaderContext) interface{}

	// Visit a parse tree produced by bapelParser#implementsHeader.
	VisitImplementsHeader(ctx *ImplementsHeaderContext) interface{}

	// Visit a parse tree produced by bapelParser#workspace.
	VisitWorkspace(ctx *WorkspaceContext) interface{}

	// Visit a parse tree produced by bapelParser#packagesSection.
	VisitPackagesSection(ctx *PackagesSectionContext) interface{}

	// Visit a parse tree produced by bapelParser#packageRule.
	VisitPackageRule(ctx *PackageRuleContext) interface{}

	// Visit a parse tree produced by bapelParser#importsSection.
	VisitImportsSection(ctx *ImportsSectionContext) interface{}

	// Visit a parse tree produced by bapelParser#implsSection.
	VisitImplsSection(ctx *ImplsSectionContext) interface{}

	// Visit a parse tree produced by bapelParser#flagsSection.
	VisitFlagsSection(ctx *FlagsSectionContext) interface{}

	// Visit a parse tree produced by bapelParser#moduleID.
	VisitModuleID(ctx *ModuleIDContext) interface{}

	// Visit a parse tree produced by bapelParser#filename.
	VisitFilename(ctx *FilenameContext) interface{}

	// Visit a parse tree produced by bapelParser#sources.
	VisitSources(ctx *SourcesContext) interface{}

	// Visit a parse tree produced by bapelParser#source.
	VisitSource(ctx *SourceContext) interface{}

	// Visit a parse tree produced by bapelParser#traitDecl.
	VisitTraitDecl(ctx *TraitDeclContext) interface{}

	// Visit a parse tree produced by bapelParser#traitMethod.
	VisitTraitMethod(ctx *TraitMethodContext) interface{}

	// Visit a parse tree produced by bapelParser#traitImpl.
	VisitTraitImpl(ctx *TraitImplContext) interface{}

	// Visit a parse tree produced by bapelParser#inherentImpl.
	VisitInherentImpl(ctx *InherentImplContext) interface{}

	// Visit a parse tree produced by bapelParser#declNoExport.
	VisitDeclNoExport(ctx *DeclNoExportContext) interface{}

	// Visit a parse tree produced by bapelParser#functionNoExport.
	VisitFunctionNoExport(ctx *FunctionNoExportContext) interface{}

	// Visit a parse tree produced by bapelParser#functionArgs.
	VisitFunctionArgs(ctx *FunctionArgsContext) interface{}

	// Visit a parse tree produced by bapelParser#arg.
	VisitArg(ctx *ArgContext) interface{}

	// Visit a parse tree produced by bapelParser#decl.
	VisitDecl(ctx *DeclContext) interface{}

	// Visit a parse tree produced by bapelParser#unexportedDecl.
	VisitUnexportedDecl(ctx *UnexportedDeclContext) interface{}

	// Visit a parse tree produced by bapelParser#declNoTerm.
	VisitDeclNoTerm(ctx *DeclNoTermContext) interface{}

	// Visit a parse tree produced by bapelParser#termDecl.
	VisitTermDecl(ctx *TermDeclContext) interface{}

	// Visit a parse tree produced by bapelParser#typeDecl.
	VisitTypeDecl(ctx *TypeDeclContext) interface{}

	// Visit a parse tree produced by bapelParser#typeAbstraction.
	VisitTypeAbstraction(ctx *TypeAbstractionContext) interface{}

	// Visit a parse tree produced by bapelParser#tvar.
	VisitTvar(ctx *TvarContext) interface{}

	// Visit a parse tree produced by bapelParser#type_.
	VisitType_(ctx *Type_Context) interface{}

	// Visit a parse tree produced by bapelParser#forallType.
	VisitForallType(ctx *ForallTypeContext) interface{}

	// Visit a parse tree produced by bapelParser#functionType.
	VisitFunctionType(ctx *FunctionTypeContext) interface{}

	// Visit a parse tree produced by bapelParser#ptrType.
	VisitPtrType(ctx *PtrTypeContext) interface{}

	// Visit a parse tree produced by bapelParser#appType.
	VisitAppType(ctx *AppTypeContext) interface{}

	// Visit a parse tree produced by bapelParser#primaryType.
	VisitPrimaryType(ctx *PrimaryTypeContext) interface{}

	// Visit a parse tree produced by bapelParser#arrayType.
	VisitArrayType(ctx *ArrayTypeContext) interface{}

	// Visit a parse tree produced by bapelParser#structType.
	VisitStructType(ctx *StructTypeContext) interface{}

	// Visit a parse tree produced by bapelParser#fields.
	VisitFields(ctx *FieldsContext) interface{}

	// Visit a parse tree produced by bapelParser#field.
	VisitField(ctx *FieldContext) interface{}

	// Visit a parse tree produced by bapelParser#tupleType.
	VisitTupleType(ctx *TupleTypeContext) interface{}

	// Visit a parse tree produced by bapelParser#tupleTypeArgs.
	VisitTupleTypeArgs(ctx *TupleTypeArgsContext) interface{}

	// Visit a parse tree produced by bapelParser#variantType.
	VisitVariantType(ctx *VariantTypeContext) interface{}

	// Visit a parse tree produced by bapelParser#tags.
	VisitTags(ctx *TagsContext) interface{}

	// Visit a parse tree produced by bapelParser#tag.
	VisitTag(ctx *TagContext) interface{}

	// Visit a parse tree produced by bapelParser#expression.
	VisitExpression(ctx *ExpressionContext) interface{}

	// Visit a parse tree produced by bapelParser#expressionWithoutBlock.
	VisitExpressionWithoutBlock(ctx *ExpressionWithoutBlockContext) interface{}

	// Visit a parse tree produced by bapelParser#expressionWithBlock.
	VisitExpressionWithBlock(ctx *ExpressionWithBlockContext) interface{}

	// Visit a parse tree produced by bapelParser#assignTerm.
	VisitAssignTerm(ctx *AssignTermContext) interface{}

	// Visit a parse tree produced by bapelParser#returnTerm.
	VisitReturnTerm(ctx *ReturnTermContext) interface{}

	// Visit a parse tree produced by bapelParser#ifTerm.
	VisitIfTerm(ctx *IfTermContext) interface{}

	// Visit a parse tree produced by bapelParser#forTerm.
	VisitForTerm(ctx *ForTermContext) interface{}

	// Visit a parse tree produced by bapelParser#lambdaTerm.
	VisitLambdaTerm(ctx *LambdaTermContext) interface{}

	// Visit a parse tree produced by bapelParser#matchTerm.
	VisitMatchTerm(ctx *MatchTermContext) interface{}

	// Visit a parse tree produced by bapelParser#matchArms.
	VisitMatchArms(ctx *MatchArmsContext) interface{}

	// Visit a parse tree produced by bapelParser#matchArm.
	VisitMatchArm(ctx *MatchArmContext) interface{}

	// Visit a parse tree produced by bapelParser#setTerm.
	VisitSetTerm(ctx *SetTermContext) interface{}

	// Visit a parse tree produced by bapelParser#blockExpr.
	VisitBlockExpr(ctx *BlockExprContext) interface{}

	// Visit a parse tree produced by bapelParser#blockStatements.
	VisitBlockStatements(ctx *BlockStatementsContext) interface{}

	// Visit a parse tree produced by bapelParser#statements.
	VisitStatements(ctx *StatementsContext) interface{}

	// Visit a parse tree produced by bapelParser#statement.
	VisitStatement(ctx *StatementContext) interface{}

	// Visit a parse tree produced by bapelParser#letStatement.
	VisitLetStatement(ctx *LetStatementContext) interface{}

	// Visit a parse tree produced by bapelParser#expressionStatement.
	VisitExpressionStatement(ctx *ExpressionStatementContext) interface{}

	// Visit a parse tree produced by bapelParser#operatorExpr.
	VisitOperatorExpr(ctx *OperatorExprContext) interface{}

	// Visit a parse tree produced by bapelParser#logicalOrExpr.
	VisitLogicalOrExpr(ctx *LogicalOrExprContext) interface{}

	// Visit a parse tree produced by bapelParser#logicalAndExpr.
	VisitLogicalAndExpr(ctx *LogicalAndExprContext) interface{}

	// Visit a parse tree produced by bapelParser#equalityExpr.
	VisitEqualityExpr(ctx *EqualityExprContext) interface{}

	// Visit a parse tree produced by bapelParser#comparisonExpr.
	VisitComparisonExpr(ctx *ComparisonExprContext) interface{}

	// Visit a parse tree produced by bapelParser#additiveExpr.
	VisitAdditiveExpr(ctx *AdditiveExprContext) interface{}

	// Visit a parse tree produced by bapelParser#multiplicativeExpr.
	VisitMultiplicativeExpr(ctx *MultiplicativeExprContext) interface{}

	// Visit a parse tree produced by bapelParser#unaryExpr.
	VisitUnaryExpr(ctx *UnaryExprContext) interface{}

	// Visit a parse tree produced by bapelParser#applicativeExpr.
	VisitApplicativeExpr(ctx *ApplicativeExprContext) interface{}

	// Visit a parse tree produced by bapelParser#typeApplicativeExpr.
	VisitTypeApplicativeExpr(ctx *TypeApplicativeExprContext) interface{}

	// Visit a parse tree produced by bapelParser#typeApplicativeArgs.
	VisitTypeApplicativeArgs(ctx *TypeApplicativeArgsContext) interface{}

	// Visit a parse tree produced by bapelParser#basePrimaryExpr.
	VisitBasePrimaryExpr(ctx *BasePrimaryExprContext) interface{}

	// Visit a parse tree produced by bapelParser#primaryExpr.
	VisitPrimaryExpr(ctx *PrimaryExprContext) interface{}

	// Visit a parse tree produced by bapelParser#projectionExpr.
	VisitProjectionExpr(ctx *ProjectionExprContext) interface{}

	// Visit a parse tree produced by bapelParser#derefExpr.
	VisitDerefExpr(ctx *DerefExprContext) interface{}

	// Visit a parse tree produced by bapelParser#injectionExpr.
	VisitInjectionExpr(ctx *InjectionExprContext) interface{}

	// Visit a parse tree produced by bapelParser#structExpr.
	VisitStructExpr(ctx *StructExprContext) interface{}

	// Visit a parse tree produced by bapelParser#labelValues.
	VisitLabelValues(ctx *LabelValuesContext) interface{}

	// Visit a parse tree produced by bapelParser#labelValue.
	VisitLabelValue(ctx *LabelValueContext) interface{}

	// Visit a parse tree produced by bapelParser#tupleExpr.
	VisitTupleExpr(ctx *TupleExprContext) interface{}

	// Visit a parse tree produced by bapelParser#tupleExprArgs.
	VisitTupleExprArgs(ctx *TupleExprArgsContext) interface{}

	// Visit a parse tree produced by bapelParser#varExpr.
	VisitVarExpr(ctx *VarExprContext) interface{}

	// Visit a parse tree produced by bapelParser#id.
	VisitId(ctx *IdContext) interface{}

	// Visit a parse tree produced by bapelParser#idTokens.
	VisitIdTokens(ctx *IdTokensContext) interface{}
}
