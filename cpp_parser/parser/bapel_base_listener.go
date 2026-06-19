// Code generated from cpp_parser/bapel.g4 by ANTLR 4.9.2. DO NOT EDIT.

package parser // bapel

import "github.com/antlr/antlr4/runtime/Go/antlr"

// BasebapelListener is a complete listener for a parse tree produced by bapelParser.
type BasebapelListener struct{}

var _ bapelListener = &BasebapelListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BasebapelListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BasebapelListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BasebapelListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BasebapelListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterBaseSourceFile is called when production baseSourceFile is entered.
func (s *BasebapelListener) EnterBaseSourceFile(ctx *BaseSourceFileContext) {}

// ExitBaseSourceFile is called when production baseSourceFile is exited.
func (s *BasebapelListener) ExitBaseSourceFile(ctx *BaseSourceFileContext) {}

// EnterImplSourceFile is called when production implSourceFile is entered.
func (s *BasebapelListener) EnterImplSourceFile(ctx *ImplSourceFileContext) {}

// ExitImplSourceFile is called when production implSourceFile is exited.
func (s *BasebapelListener) ExitImplSourceFile(ctx *ImplSourceFileContext) {}

// EnterModuleHeader is called when production moduleHeader is entered.
func (s *BasebapelListener) EnterModuleHeader(ctx *ModuleHeaderContext) {}

// ExitModuleHeader is called when production moduleHeader is exited.
func (s *BasebapelListener) ExitModuleHeader(ctx *ModuleHeaderContext) {}

// EnterImplementsHeader is called when production implementsHeader is entered.
func (s *BasebapelListener) EnterImplementsHeader(ctx *ImplementsHeaderContext) {}

// ExitImplementsHeader is called when production implementsHeader is exited.
func (s *BasebapelListener) ExitImplementsHeader(ctx *ImplementsHeaderContext) {}

// EnterWorkspace is called when production workspace is entered.
func (s *BasebapelListener) EnterWorkspace(ctx *WorkspaceContext) {}

// ExitWorkspace is called when production workspace is exited.
func (s *BasebapelListener) ExitWorkspace(ctx *WorkspaceContext) {}

// EnterPackagesSection is called when production packagesSection is entered.
func (s *BasebapelListener) EnterPackagesSection(ctx *PackagesSectionContext) {}

// ExitPackagesSection is called when production packagesSection is exited.
func (s *BasebapelListener) ExitPackagesSection(ctx *PackagesSectionContext) {}

// EnterPackageRule is called when production packageRule is entered.
func (s *BasebapelListener) EnterPackageRule(ctx *PackageRuleContext) {}

// ExitPackageRule is called when production packageRule is exited.
func (s *BasebapelListener) ExitPackageRule(ctx *PackageRuleContext) {}

// EnterImportsSection is called when production importsSection is entered.
func (s *BasebapelListener) EnterImportsSection(ctx *ImportsSectionContext) {}

// ExitImportsSection is called when production importsSection is exited.
func (s *BasebapelListener) ExitImportsSection(ctx *ImportsSectionContext) {}

// EnterImplsSection is called when production implsSection is entered.
func (s *BasebapelListener) EnterImplsSection(ctx *ImplsSectionContext) {}

// ExitImplsSection is called when production implsSection is exited.
func (s *BasebapelListener) ExitImplsSection(ctx *ImplsSectionContext) {}

// EnterFlagsSection is called when production flagsSection is entered.
func (s *BasebapelListener) EnterFlagsSection(ctx *FlagsSectionContext) {}

// ExitFlagsSection is called when production flagsSection is exited.
func (s *BasebapelListener) ExitFlagsSection(ctx *FlagsSectionContext) {}

// EnterModuleID is called when production moduleID is entered.
func (s *BasebapelListener) EnterModuleID(ctx *ModuleIDContext) {}

// ExitModuleID is called when production moduleID is exited.
func (s *BasebapelListener) ExitModuleID(ctx *ModuleIDContext) {}

// EnterFilename is called when production filename is entered.
func (s *BasebapelListener) EnterFilename(ctx *FilenameContext) {}

// ExitFilename is called when production filename is exited.
func (s *BasebapelListener) ExitFilename(ctx *FilenameContext) {}

// EnterSources is called when production sources is entered.
func (s *BasebapelListener) EnterSources(ctx *SourcesContext) {}

// ExitSources is called when production sources is exited.
func (s *BasebapelListener) ExitSources(ctx *SourcesContext) {}

// EnterSource is called when production source is entered.
func (s *BasebapelListener) EnterSource(ctx *SourceContext) {}

// ExitSource is called when production source is exited.
func (s *BasebapelListener) ExitSource(ctx *SourceContext) {}

// EnterDeclNoExport is called when production declNoExport is entered.
func (s *BasebapelListener) EnterDeclNoExport(ctx *DeclNoExportContext) {}

// ExitDeclNoExport is called when production declNoExport is exited.
func (s *BasebapelListener) ExitDeclNoExport(ctx *DeclNoExportContext) {}

// EnterFunctionNoExport is called when production functionNoExport is entered.
func (s *BasebapelListener) EnterFunctionNoExport(ctx *FunctionNoExportContext) {}

// ExitFunctionNoExport is called when production functionNoExport is exited.
func (s *BasebapelListener) ExitFunctionNoExport(ctx *FunctionNoExportContext) {}

// EnterFunctionArgs is called when production functionArgs is entered.
func (s *BasebapelListener) EnterFunctionArgs(ctx *FunctionArgsContext) {}

// ExitFunctionArgs is called when production functionArgs is exited.
func (s *BasebapelListener) ExitFunctionArgs(ctx *FunctionArgsContext) {}

// EnterArg is called when production arg is entered.
func (s *BasebapelListener) EnterArg(ctx *ArgContext) {}

// ExitArg is called when production arg is exited.
func (s *BasebapelListener) ExitArg(ctx *ArgContext) {}

// EnterDecl is called when production decl is entered.
func (s *BasebapelListener) EnterDecl(ctx *DeclContext) {}

// ExitDecl is called when production decl is exited.
func (s *BasebapelListener) ExitDecl(ctx *DeclContext) {}

// EnterUnexportedDecl is called when production unexportedDecl is entered.
func (s *BasebapelListener) EnterUnexportedDecl(ctx *UnexportedDeclContext) {}

// ExitUnexportedDecl is called when production unexportedDecl is exited.
func (s *BasebapelListener) ExitUnexportedDecl(ctx *UnexportedDeclContext) {}

// EnterDeclNoTerm is called when production declNoTerm is entered.
func (s *BasebapelListener) EnterDeclNoTerm(ctx *DeclNoTermContext) {}

// ExitDeclNoTerm is called when production declNoTerm is exited.
func (s *BasebapelListener) ExitDeclNoTerm(ctx *DeclNoTermContext) {}

// EnterTermDecl is called when production termDecl is entered.
func (s *BasebapelListener) EnterTermDecl(ctx *TermDeclContext) {}

// ExitTermDecl is called when production termDecl is exited.
func (s *BasebapelListener) ExitTermDecl(ctx *TermDeclContext) {}

// EnterTypeDecl is called when production typeDecl is entered.
func (s *BasebapelListener) EnterTypeDecl(ctx *TypeDeclContext) {}

// ExitTypeDecl is called when production typeDecl is exited.
func (s *BasebapelListener) ExitTypeDecl(ctx *TypeDeclContext) {}

// EnterTypeAbstraction is called when production typeAbstraction is entered.
func (s *BasebapelListener) EnterTypeAbstraction(ctx *TypeAbstractionContext) {}

// ExitTypeAbstraction is called when production typeAbstraction is exited.
func (s *BasebapelListener) ExitTypeAbstraction(ctx *TypeAbstractionContext) {}

// EnterTvar is called when production tvar is entered.
func (s *BasebapelListener) EnterTvar(ctx *TvarContext) {}

// ExitTvar is called when production tvar is exited.
func (s *BasebapelListener) ExitTvar(ctx *TvarContext) {}

// EnterType_ is called when production type_ is entered.
func (s *BasebapelListener) EnterType_(ctx *Type_Context) {}

// ExitType_ is called when production type_ is exited.
func (s *BasebapelListener) ExitType_(ctx *Type_Context) {}

// EnterForallType is called when production forallType is entered.
func (s *BasebapelListener) EnterForallType(ctx *ForallTypeContext) {}

// ExitForallType is called when production forallType is exited.
func (s *BasebapelListener) ExitForallType(ctx *ForallTypeContext) {}

// EnterFunctionType is called when production functionType is entered.
func (s *BasebapelListener) EnterFunctionType(ctx *FunctionTypeContext) {}

// ExitFunctionType is called when production functionType is exited.
func (s *BasebapelListener) ExitFunctionType(ctx *FunctionTypeContext) {}

// EnterPtrType is called when production ptrType is entered.
func (s *BasebapelListener) EnterPtrType(ctx *PtrTypeContext) {}

// ExitPtrType is called when production ptrType is exited.
func (s *BasebapelListener) ExitPtrType(ctx *PtrTypeContext) {}

// EnterAppType is called when production appType is entered.
func (s *BasebapelListener) EnterAppType(ctx *AppTypeContext) {}

// ExitAppType is called when production appType is exited.
func (s *BasebapelListener) ExitAppType(ctx *AppTypeContext) {}

// EnterPrimaryType is called when production primaryType is entered.
func (s *BasebapelListener) EnterPrimaryType(ctx *PrimaryTypeContext) {}

// ExitPrimaryType is called when production primaryType is exited.
func (s *BasebapelListener) ExitPrimaryType(ctx *PrimaryTypeContext) {}

// EnterArrayType is called when production arrayType is entered.
func (s *BasebapelListener) EnterArrayType(ctx *ArrayTypeContext) {}

// ExitArrayType is called when production arrayType is exited.
func (s *BasebapelListener) ExitArrayType(ctx *ArrayTypeContext) {}

// EnterStructType is called when production structType is entered.
func (s *BasebapelListener) EnterStructType(ctx *StructTypeContext) {}

// ExitStructType is called when production structType is exited.
func (s *BasebapelListener) ExitStructType(ctx *StructTypeContext) {}

// EnterFields is called when production fields is entered.
func (s *BasebapelListener) EnterFields(ctx *FieldsContext) {}

// ExitFields is called when production fields is exited.
func (s *BasebapelListener) ExitFields(ctx *FieldsContext) {}

// EnterField is called when production field is entered.
func (s *BasebapelListener) EnterField(ctx *FieldContext) {}

// ExitField is called when production field is exited.
func (s *BasebapelListener) ExitField(ctx *FieldContext) {}

// EnterTupleType is called when production tupleType is entered.
func (s *BasebapelListener) EnterTupleType(ctx *TupleTypeContext) {}

// ExitTupleType is called when production tupleType is exited.
func (s *BasebapelListener) ExitTupleType(ctx *TupleTypeContext) {}

// EnterTupleTypeArgs is called when production tupleTypeArgs is entered.
func (s *BasebapelListener) EnterTupleTypeArgs(ctx *TupleTypeArgsContext) {}

// ExitTupleTypeArgs is called when production tupleTypeArgs is exited.
func (s *BasebapelListener) ExitTupleTypeArgs(ctx *TupleTypeArgsContext) {}

// EnterVariantType is called when production variantType is entered.
func (s *BasebapelListener) EnterVariantType(ctx *VariantTypeContext) {}

// ExitVariantType is called when production variantType is exited.
func (s *BasebapelListener) ExitVariantType(ctx *VariantTypeContext) {}

// EnterTags is called when production tags is entered.
func (s *BasebapelListener) EnterTags(ctx *TagsContext) {}

// ExitTags is called when production tags is exited.
func (s *BasebapelListener) ExitTags(ctx *TagsContext) {}

// EnterTag is called when production tag is entered.
func (s *BasebapelListener) EnterTag(ctx *TagContext) {}

// ExitTag is called when production tag is exited.
func (s *BasebapelListener) ExitTag(ctx *TagContext) {}

// EnterExpression is called when production expression is entered.
func (s *BasebapelListener) EnterExpression(ctx *ExpressionContext) {}

// ExitExpression is called when production expression is exited.
func (s *BasebapelListener) ExitExpression(ctx *ExpressionContext) {}

// EnterExpressionWithoutBlock is called when production expressionWithoutBlock is entered.
func (s *BasebapelListener) EnterExpressionWithoutBlock(ctx *ExpressionWithoutBlockContext) {}

// ExitExpressionWithoutBlock is called when production expressionWithoutBlock is exited.
func (s *BasebapelListener) ExitExpressionWithoutBlock(ctx *ExpressionWithoutBlockContext) {}

// EnterExpressionWithBlock is called when production expressionWithBlock is entered.
func (s *BasebapelListener) EnterExpressionWithBlock(ctx *ExpressionWithBlockContext) {}

// ExitExpressionWithBlock is called when production expressionWithBlock is exited.
func (s *BasebapelListener) ExitExpressionWithBlock(ctx *ExpressionWithBlockContext) {}

// EnterAssignTerm is called when production assignTerm is entered.
func (s *BasebapelListener) EnterAssignTerm(ctx *AssignTermContext) {}

// ExitAssignTerm is called when production assignTerm is exited.
func (s *BasebapelListener) ExitAssignTerm(ctx *AssignTermContext) {}

// EnterReturnTerm is called when production returnTerm is entered.
func (s *BasebapelListener) EnterReturnTerm(ctx *ReturnTermContext) {}

// ExitReturnTerm is called when production returnTerm is exited.
func (s *BasebapelListener) ExitReturnTerm(ctx *ReturnTermContext) {}

// EnterIfTerm is called when production ifTerm is entered.
func (s *BasebapelListener) EnterIfTerm(ctx *IfTermContext) {}

// ExitIfTerm is called when production ifTerm is exited.
func (s *BasebapelListener) ExitIfTerm(ctx *IfTermContext) {}

// EnterForTerm is called when production forTerm is entered.
func (s *BasebapelListener) EnterForTerm(ctx *ForTermContext) {}

// ExitForTerm is called when production forTerm is exited.
func (s *BasebapelListener) ExitForTerm(ctx *ForTermContext) {}

// EnterLambdaTerm is called when production lambdaTerm is entered.
func (s *BasebapelListener) EnterLambdaTerm(ctx *LambdaTermContext) {}

// ExitLambdaTerm is called when production lambdaTerm is exited.
func (s *BasebapelListener) ExitLambdaTerm(ctx *LambdaTermContext) {}

// EnterMatchTerm is called when production matchTerm is entered.
func (s *BasebapelListener) EnterMatchTerm(ctx *MatchTermContext) {}

// ExitMatchTerm is called when production matchTerm is exited.
func (s *BasebapelListener) ExitMatchTerm(ctx *MatchTermContext) {}

// EnterMatchArms is called when production matchArms is entered.
func (s *BasebapelListener) EnterMatchArms(ctx *MatchArmsContext) {}

// ExitMatchArms is called when production matchArms is exited.
func (s *BasebapelListener) ExitMatchArms(ctx *MatchArmsContext) {}

// EnterMatchArm is called when production matchArm is entered.
func (s *BasebapelListener) EnterMatchArm(ctx *MatchArmContext) {}

// ExitMatchArm is called when production matchArm is exited.
func (s *BasebapelListener) ExitMatchArm(ctx *MatchArmContext) {}

// EnterSetTerm is called when production setTerm is entered.
func (s *BasebapelListener) EnterSetTerm(ctx *SetTermContext) {}

// ExitSetTerm is called when production setTerm is exited.
func (s *BasebapelListener) ExitSetTerm(ctx *SetTermContext) {}

// EnterBlockExpr is called when production blockExpr is entered.
func (s *BasebapelListener) EnterBlockExpr(ctx *BlockExprContext) {}

// ExitBlockExpr is called when production blockExpr is exited.
func (s *BasebapelListener) ExitBlockExpr(ctx *BlockExprContext) {}

// EnterBlockStatements is called when production blockStatements is entered.
func (s *BasebapelListener) EnterBlockStatements(ctx *BlockStatementsContext) {}

// ExitBlockStatements is called when production blockStatements is exited.
func (s *BasebapelListener) ExitBlockStatements(ctx *BlockStatementsContext) {}

// EnterStatements is called when production statements is entered.
func (s *BasebapelListener) EnterStatements(ctx *StatementsContext) {}

// ExitStatements is called when production statements is exited.
func (s *BasebapelListener) ExitStatements(ctx *StatementsContext) {}

// EnterStatement is called when production statement is entered.
func (s *BasebapelListener) EnterStatement(ctx *StatementContext) {}

// ExitStatement is called when production statement is exited.
func (s *BasebapelListener) ExitStatement(ctx *StatementContext) {}

// EnterLetStatement is called when production letStatement is entered.
func (s *BasebapelListener) EnterLetStatement(ctx *LetStatementContext) {}

// ExitLetStatement is called when production letStatement is exited.
func (s *BasebapelListener) ExitLetStatement(ctx *LetStatementContext) {}

// EnterExpressionStatement is called when production expressionStatement is entered.
func (s *BasebapelListener) EnterExpressionStatement(ctx *ExpressionStatementContext) {}

// ExitExpressionStatement is called when production expressionStatement is exited.
func (s *BasebapelListener) ExitExpressionStatement(ctx *ExpressionStatementContext) {}

// EnterOperatorExpr is called when production operatorExpr is entered.
func (s *BasebapelListener) EnterOperatorExpr(ctx *OperatorExprContext) {}

// ExitOperatorExpr is called when production operatorExpr is exited.
func (s *BasebapelListener) ExitOperatorExpr(ctx *OperatorExprContext) {}

// EnterLogicalOrExpr is called when production logicalOrExpr is entered.
func (s *BasebapelListener) EnterLogicalOrExpr(ctx *LogicalOrExprContext) {}

// ExitLogicalOrExpr is called when production logicalOrExpr is exited.
func (s *BasebapelListener) ExitLogicalOrExpr(ctx *LogicalOrExprContext) {}

// EnterLogicalAndExpr is called when production logicalAndExpr is entered.
func (s *BasebapelListener) EnterLogicalAndExpr(ctx *LogicalAndExprContext) {}

// ExitLogicalAndExpr is called when production logicalAndExpr is exited.
func (s *BasebapelListener) ExitLogicalAndExpr(ctx *LogicalAndExprContext) {}

// EnterEqualityExpr is called when production equalityExpr is entered.
func (s *BasebapelListener) EnterEqualityExpr(ctx *EqualityExprContext) {}

// ExitEqualityExpr is called when production equalityExpr is exited.
func (s *BasebapelListener) ExitEqualityExpr(ctx *EqualityExprContext) {}

// EnterComparisonExpr is called when production comparisonExpr is entered.
func (s *BasebapelListener) EnterComparisonExpr(ctx *ComparisonExprContext) {}

// ExitComparisonExpr is called when production comparisonExpr is exited.
func (s *BasebapelListener) ExitComparisonExpr(ctx *ComparisonExprContext) {}

// EnterAdditiveExpr is called when production additiveExpr is entered.
func (s *BasebapelListener) EnterAdditiveExpr(ctx *AdditiveExprContext) {}

// ExitAdditiveExpr is called when production additiveExpr is exited.
func (s *BasebapelListener) ExitAdditiveExpr(ctx *AdditiveExprContext) {}

// EnterMultiplicativeExpr is called when production multiplicativeExpr is entered.
func (s *BasebapelListener) EnterMultiplicativeExpr(ctx *MultiplicativeExprContext) {}

// ExitMultiplicativeExpr is called when production multiplicativeExpr is exited.
func (s *BasebapelListener) ExitMultiplicativeExpr(ctx *MultiplicativeExprContext) {}

// EnterUnaryExpr is called when production unaryExpr is entered.
func (s *BasebapelListener) EnterUnaryExpr(ctx *UnaryExprContext) {}

// ExitUnaryExpr is called when production unaryExpr is exited.
func (s *BasebapelListener) ExitUnaryExpr(ctx *UnaryExprContext) {}

// EnterApplicativeExpr is called when production applicativeExpr is entered.
func (s *BasebapelListener) EnterApplicativeExpr(ctx *ApplicativeExprContext) {}

// ExitApplicativeExpr is called when production applicativeExpr is exited.
func (s *BasebapelListener) ExitApplicativeExpr(ctx *ApplicativeExprContext) {}

// EnterTypeApplicativeExpr is called when production typeApplicativeExpr is entered.
func (s *BasebapelListener) EnterTypeApplicativeExpr(ctx *TypeApplicativeExprContext) {}

// ExitTypeApplicativeExpr is called when production typeApplicativeExpr is exited.
func (s *BasebapelListener) ExitTypeApplicativeExpr(ctx *TypeApplicativeExprContext) {}

// EnterTypeApplicativeArgs is called when production typeApplicativeArgs is entered.
func (s *BasebapelListener) EnterTypeApplicativeArgs(ctx *TypeApplicativeArgsContext) {}

// ExitTypeApplicativeArgs is called when production typeApplicativeArgs is exited.
func (s *BasebapelListener) ExitTypeApplicativeArgs(ctx *TypeApplicativeArgsContext) {}

// EnterBasePrimaryExpr is called when production basePrimaryExpr is entered.
func (s *BasebapelListener) EnterBasePrimaryExpr(ctx *BasePrimaryExprContext) {}

// ExitBasePrimaryExpr is called when production basePrimaryExpr is exited.
func (s *BasebapelListener) ExitBasePrimaryExpr(ctx *BasePrimaryExprContext) {}

// EnterPrimaryExpr is called when production primaryExpr is entered.
func (s *BasebapelListener) EnterPrimaryExpr(ctx *PrimaryExprContext) {}

// ExitPrimaryExpr is called when production primaryExpr is exited.
func (s *BasebapelListener) ExitPrimaryExpr(ctx *PrimaryExprContext) {}

// EnterProjectionExpr is called when production projectionExpr is entered.
func (s *BasebapelListener) EnterProjectionExpr(ctx *ProjectionExprContext) {}

// ExitProjectionExpr is called when production projectionExpr is exited.
func (s *BasebapelListener) ExitProjectionExpr(ctx *ProjectionExprContext) {}

// EnterDerefExpr is called when production derefExpr is entered.
func (s *BasebapelListener) EnterDerefExpr(ctx *DerefExprContext) {}

// ExitDerefExpr is called when production derefExpr is exited.
func (s *BasebapelListener) ExitDerefExpr(ctx *DerefExprContext) {}

// EnterInjectionExpr is called when production injectionExpr is entered.
func (s *BasebapelListener) EnterInjectionExpr(ctx *InjectionExprContext) {}

// ExitInjectionExpr is called when production injectionExpr is exited.
func (s *BasebapelListener) ExitInjectionExpr(ctx *InjectionExprContext) {}

// EnterStructExpr is called when production structExpr is entered.
func (s *BasebapelListener) EnterStructExpr(ctx *StructExprContext) {}

// ExitStructExpr is called when production structExpr is exited.
func (s *BasebapelListener) ExitStructExpr(ctx *StructExprContext) {}

// EnterLabelValues is called when production labelValues is entered.
func (s *BasebapelListener) EnterLabelValues(ctx *LabelValuesContext) {}

// ExitLabelValues is called when production labelValues is exited.
func (s *BasebapelListener) ExitLabelValues(ctx *LabelValuesContext) {}

// EnterLabelValue is called when production labelValue is entered.
func (s *BasebapelListener) EnterLabelValue(ctx *LabelValueContext) {}

// ExitLabelValue is called when production labelValue is exited.
func (s *BasebapelListener) ExitLabelValue(ctx *LabelValueContext) {}

// EnterTupleExpr is called when production tupleExpr is entered.
func (s *BasebapelListener) EnterTupleExpr(ctx *TupleExprContext) {}

// ExitTupleExpr is called when production tupleExpr is exited.
func (s *BasebapelListener) ExitTupleExpr(ctx *TupleExprContext) {}

// EnterTupleExprArgs is called when production tupleExprArgs is entered.
func (s *BasebapelListener) EnterTupleExprArgs(ctx *TupleExprArgsContext) {}

// ExitTupleExprArgs is called when production tupleExprArgs is exited.
func (s *BasebapelListener) ExitTupleExprArgs(ctx *TupleExprArgsContext) {}

// EnterVarExpr is called when production varExpr is entered.
func (s *BasebapelListener) EnterVarExpr(ctx *VarExprContext) {}

// ExitVarExpr is called when production varExpr is exited.
func (s *BasebapelListener) ExitVarExpr(ctx *VarExprContext) {}

// EnterId is called when production id is entered.
func (s *BasebapelListener) EnterId(ctx *IdContext) {}

// ExitId is called when production id is exited.
func (s *BasebapelListener) ExitId(ctx *IdContext) {}

// EnterIdTokens is called when production idTokens is entered.
func (s *BasebapelListener) EnterIdTokens(ctx *IdTokensContext) {}

// ExitIdTokens is called when production idTokens is exited.
func (s *BasebapelListener) ExitIdTokens(ctx *IdTokensContext) {}
