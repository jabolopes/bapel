// Code generated from cpp_parser/bapel.g4 by ANTLR 4.9.2. DO NOT EDIT.

package parser // bapel

import "github.com/antlr/antlr4/runtime/Go/antlr"

type BasebapelVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BasebapelVisitor) VisitBaseSourceFile(ctx *BaseSourceFileContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitImplSourceFile(ctx *ImplSourceFileContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitModuleHeader(ctx *ModuleHeaderContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitImplementsHeader(ctx *ImplementsHeaderContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitWorkspace(ctx *WorkspaceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitPackagesSection(ctx *PackagesSectionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitPackageRule(ctx *PackageRuleContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitImportsSection(ctx *ImportsSectionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitImplsSection(ctx *ImplsSectionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitFlagsSection(ctx *FlagsSectionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitModuleID(ctx *ModuleIDContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitFilename(ctx *FilenameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitSources(ctx *SourcesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitSource(ctx *SourceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitTraitDecl(ctx *TraitDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitTraitMethod(ctx *TraitMethodContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitImplBlock(ctx *ImplBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitDeclNoExport(ctx *DeclNoExportContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitFunctionNoExport(ctx *FunctionNoExportContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitFunctionArgs(ctx *FunctionArgsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitArg(ctx *ArgContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitDecl(ctx *DeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitUnexportedDecl(ctx *UnexportedDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitDeclNoTerm(ctx *DeclNoTermContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitTermDecl(ctx *TermDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitTypeDecl(ctx *TypeDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitTypeAbstraction(ctx *TypeAbstractionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitTvar(ctx *TvarContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitType_(ctx *Type_Context) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitForallType(ctx *ForallTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitFunctionType(ctx *FunctionTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitPtrType(ctx *PtrTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitAppType(ctx *AppTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitPrimaryType(ctx *PrimaryTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitArrayType(ctx *ArrayTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitStructType(ctx *StructTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitFields(ctx *FieldsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitField(ctx *FieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitTupleType(ctx *TupleTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitTupleTypeArgs(ctx *TupleTypeArgsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitVariantType(ctx *VariantTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitTags(ctx *TagsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitTag(ctx *TagContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitExpression(ctx *ExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitExpressionWithoutBlock(ctx *ExpressionWithoutBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitExpressionWithBlock(ctx *ExpressionWithBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitAssignTerm(ctx *AssignTermContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitReturnTerm(ctx *ReturnTermContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitIfTerm(ctx *IfTermContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitForTerm(ctx *ForTermContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitLambdaTerm(ctx *LambdaTermContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitMatchTerm(ctx *MatchTermContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitMatchArms(ctx *MatchArmsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitMatchArm(ctx *MatchArmContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitSetTerm(ctx *SetTermContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitBlockExpr(ctx *BlockExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitBlockStatements(ctx *BlockStatementsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitStatements(ctx *StatementsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitStatement(ctx *StatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitLetStatement(ctx *LetStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitExpressionStatement(ctx *ExpressionStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitOperatorExpr(ctx *OperatorExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitLogicalOrExpr(ctx *LogicalOrExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitLogicalAndExpr(ctx *LogicalAndExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitEqualityExpr(ctx *EqualityExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitComparisonExpr(ctx *ComparisonExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitAdditiveExpr(ctx *AdditiveExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitMultiplicativeExpr(ctx *MultiplicativeExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitUnaryExpr(ctx *UnaryExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitApplicativeExpr(ctx *ApplicativeExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitTypeApplicativeExpr(ctx *TypeApplicativeExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitTypeApplicativeArgs(ctx *TypeApplicativeArgsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitBasePrimaryExpr(ctx *BasePrimaryExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitPrimaryExpr(ctx *PrimaryExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitProjectionExpr(ctx *ProjectionExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitDerefExpr(ctx *DerefExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitInjectionExpr(ctx *InjectionExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitStructExpr(ctx *StructExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitLabelValues(ctx *LabelValuesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitLabelValue(ctx *LabelValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitTupleExpr(ctx *TupleExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitTupleExprArgs(ctx *TupleExprArgsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitVarExpr(ctx *VarExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitId(ctx *IdContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasebapelVisitor) VisitIdTokens(ctx *IdTokensContext) interface{} {
	return v.VisitChildren(ctx)
}
