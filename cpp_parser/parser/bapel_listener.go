// Code generated from cpp_parser/bapel.g4 by ANTLR 4.9.2. DO NOT EDIT.

package parser // bapel

import "github.com/antlr/antlr4/runtime/Go/antlr"

// bapelListener is a complete listener for a parse tree produced by bapelParser.
type bapelListener interface {
	antlr.ParseTreeListener

	// EnterBaseSourceFile is called when entering the baseSourceFile production.
	EnterBaseSourceFile(c *BaseSourceFileContext)

	// EnterImplSourceFile is called when entering the implSourceFile production.
	EnterImplSourceFile(c *ImplSourceFileContext)

	// EnterModuleHeader is called when entering the moduleHeader production.
	EnterModuleHeader(c *ModuleHeaderContext)

	// EnterImplementsHeader is called when entering the implementsHeader production.
	EnterImplementsHeader(c *ImplementsHeaderContext)

	// EnterWorkspace is called when entering the workspace production.
	EnterWorkspace(c *WorkspaceContext)

	// EnterPackagesSection is called when entering the packagesSection production.
	EnterPackagesSection(c *PackagesSectionContext)

	// EnterPackageRule is called when entering the packageRule production.
	EnterPackageRule(c *PackageRuleContext)

	// EnterImportsSection is called when entering the importsSection production.
	EnterImportsSection(c *ImportsSectionContext)

	// EnterImplsSection is called when entering the implsSection production.
	EnterImplsSection(c *ImplsSectionContext)

	// EnterFlagsSection is called when entering the flagsSection production.
	EnterFlagsSection(c *FlagsSectionContext)

	// EnterModuleID is called when entering the moduleID production.
	EnterModuleID(c *ModuleIDContext)

	// EnterFilename is called when entering the filename production.
	EnterFilename(c *FilenameContext)

	// EnterSources is called when entering the sources production.
	EnterSources(c *SourcesContext)

	// EnterSource is called when entering the source production.
	EnterSource(c *SourceContext)

	// EnterTraitDecl is called when entering the traitDecl production.
	EnterTraitDecl(c *TraitDeclContext)

	// EnterTraitMethod is called when entering the traitMethod production.
	EnterTraitMethod(c *TraitMethodContext)

	// EnterImplBlock is called when entering the implBlock production.
	EnterImplBlock(c *ImplBlockContext)

	// EnterDeclNoExport is called when entering the declNoExport production.
	EnterDeclNoExport(c *DeclNoExportContext)

	// EnterFunctionNoExport is called when entering the functionNoExport production.
	EnterFunctionNoExport(c *FunctionNoExportContext)

	// EnterFunctionArgs is called when entering the functionArgs production.
	EnterFunctionArgs(c *FunctionArgsContext)

	// EnterArg is called when entering the arg production.
	EnterArg(c *ArgContext)

	// EnterDecl is called when entering the decl production.
	EnterDecl(c *DeclContext)

	// EnterUnexportedDecl is called when entering the unexportedDecl production.
	EnterUnexportedDecl(c *UnexportedDeclContext)

	// EnterDeclNoTerm is called when entering the declNoTerm production.
	EnterDeclNoTerm(c *DeclNoTermContext)

	// EnterTermDecl is called when entering the termDecl production.
	EnterTermDecl(c *TermDeclContext)

	// EnterTypeDecl is called when entering the typeDecl production.
	EnterTypeDecl(c *TypeDeclContext)

	// EnterTypeAbstraction is called when entering the typeAbstraction production.
	EnterTypeAbstraction(c *TypeAbstractionContext)

	// EnterTvar is called when entering the tvar production.
	EnterTvar(c *TvarContext)

	// EnterType_ is called when entering the type_ production.
	EnterType_(c *Type_Context)

	// EnterForallType is called when entering the forallType production.
	EnterForallType(c *ForallTypeContext)

	// EnterFunctionType is called when entering the functionType production.
	EnterFunctionType(c *FunctionTypeContext)

	// EnterPtrType is called when entering the ptrType production.
	EnterPtrType(c *PtrTypeContext)

	// EnterAppType is called when entering the appType production.
	EnterAppType(c *AppTypeContext)

	// EnterPrimaryType is called when entering the primaryType production.
	EnterPrimaryType(c *PrimaryTypeContext)

	// EnterArrayType is called when entering the arrayType production.
	EnterArrayType(c *ArrayTypeContext)

	// EnterStructType is called when entering the structType production.
	EnterStructType(c *StructTypeContext)

	// EnterFields is called when entering the fields production.
	EnterFields(c *FieldsContext)

	// EnterField is called when entering the field production.
	EnterField(c *FieldContext)

	// EnterTupleType is called when entering the tupleType production.
	EnterTupleType(c *TupleTypeContext)

	// EnterTupleTypeArgs is called when entering the tupleTypeArgs production.
	EnterTupleTypeArgs(c *TupleTypeArgsContext)

	// EnterVariantType is called when entering the variantType production.
	EnterVariantType(c *VariantTypeContext)

	// EnterTags is called when entering the tags production.
	EnterTags(c *TagsContext)

	// EnterTag is called when entering the tag production.
	EnterTag(c *TagContext)

	// EnterExpression is called when entering the expression production.
	EnterExpression(c *ExpressionContext)

	// EnterExpressionWithoutBlock is called when entering the expressionWithoutBlock production.
	EnterExpressionWithoutBlock(c *ExpressionWithoutBlockContext)

	// EnterExpressionWithBlock is called when entering the expressionWithBlock production.
	EnterExpressionWithBlock(c *ExpressionWithBlockContext)

	// EnterAssignTerm is called when entering the assignTerm production.
	EnterAssignTerm(c *AssignTermContext)

	// EnterReturnTerm is called when entering the returnTerm production.
	EnterReturnTerm(c *ReturnTermContext)

	// EnterIfTerm is called when entering the ifTerm production.
	EnterIfTerm(c *IfTermContext)

	// EnterForTerm is called when entering the forTerm production.
	EnterForTerm(c *ForTermContext)

	// EnterLambdaTerm is called when entering the lambdaTerm production.
	EnterLambdaTerm(c *LambdaTermContext)

	// EnterMatchTerm is called when entering the matchTerm production.
	EnterMatchTerm(c *MatchTermContext)

	// EnterMatchArms is called when entering the matchArms production.
	EnterMatchArms(c *MatchArmsContext)

	// EnterMatchArm is called when entering the matchArm production.
	EnterMatchArm(c *MatchArmContext)

	// EnterSetTerm is called when entering the setTerm production.
	EnterSetTerm(c *SetTermContext)

	// EnterBlockExpr is called when entering the blockExpr production.
	EnterBlockExpr(c *BlockExprContext)

	// EnterBlockStatements is called when entering the blockStatements production.
	EnterBlockStatements(c *BlockStatementsContext)

	// EnterStatements is called when entering the statements production.
	EnterStatements(c *StatementsContext)

	// EnterStatement is called when entering the statement production.
	EnterStatement(c *StatementContext)

	// EnterLetStatement is called when entering the letStatement production.
	EnterLetStatement(c *LetStatementContext)

	// EnterExpressionStatement is called when entering the expressionStatement production.
	EnterExpressionStatement(c *ExpressionStatementContext)

	// EnterOperatorExpr is called when entering the operatorExpr production.
	EnterOperatorExpr(c *OperatorExprContext)

	// EnterLogicalOrExpr is called when entering the logicalOrExpr production.
	EnterLogicalOrExpr(c *LogicalOrExprContext)

	// EnterLogicalAndExpr is called when entering the logicalAndExpr production.
	EnterLogicalAndExpr(c *LogicalAndExprContext)

	// EnterEqualityExpr is called when entering the equalityExpr production.
	EnterEqualityExpr(c *EqualityExprContext)

	// EnterComparisonExpr is called when entering the comparisonExpr production.
	EnterComparisonExpr(c *ComparisonExprContext)

	// EnterAdditiveExpr is called when entering the additiveExpr production.
	EnterAdditiveExpr(c *AdditiveExprContext)

	// EnterMultiplicativeExpr is called when entering the multiplicativeExpr production.
	EnterMultiplicativeExpr(c *MultiplicativeExprContext)

	// EnterUnaryExpr is called when entering the unaryExpr production.
	EnterUnaryExpr(c *UnaryExprContext)

	// EnterApplicativeExpr is called when entering the applicativeExpr production.
	EnterApplicativeExpr(c *ApplicativeExprContext)

	// EnterTypeApplicativeExpr is called when entering the typeApplicativeExpr production.
	EnterTypeApplicativeExpr(c *TypeApplicativeExprContext)

	// EnterTypeApplicativeArgs is called when entering the typeApplicativeArgs production.
	EnterTypeApplicativeArgs(c *TypeApplicativeArgsContext)

	// EnterBasePrimaryExpr is called when entering the basePrimaryExpr production.
	EnterBasePrimaryExpr(c *BasePrimaryExprContext)

	// EnterPrimaryExpr is called when entering the primaryExpr production.
	EnterPrimaryExpr(c *PrimaryExprContext)

	// EnterProjectionExpr is called when entering the projectionExpr production.
	EnterProjectionExpr(c *ProjectionExprContext)

	// EnterDerefExpr is called when entering the derefExpr production.
	EnterDerefExpr(c *DerefExprContext)

	// EnterInjectionExpr is called when entering the injectionExpr production.
	EnterInjectionExpr(c *InjectionExprContext)

	// EnterStructExpr is called when entering the structExpr production.
	EnterStructExpr(c *StructExprContext)

	// EnterLabelValues is called when entering the labelValues production.
	EnterLabelValues(c *LabelValuesContext)

	// EnterLabelValue is called when entering the labelValue production.
	EnterLabelValue(c *LabelValueContext)

	// EnterTupleExpr is called when entering the tupleExpr production.
	EnterTupleExpr(c *TupleExprContext)

	// EnterTupleExprArgs is called when entering the tupleExprArgs production.
	EnterTupleExprArgs(c *TupleExprArgsContext)

	// EnterVarExpr is called when entering the varExpr production.
	EnterVarExpr(c *VarExprContext)

	// EnterId is called when entering the id production.
	EnterId(c *IdContext)

	// EnterIdTokens is called when entering the idTokens production.
	EnterIdTokens(c *IdTokensContext)

	// ExitBaseSourceFile is called when exiting the baseSourceFile production.
	ExitBaseSourceFile(c *BaseSourceFileContext)

	// ExitImplSourceFile is called when exiting the implSourceFile production.
	ExitImplSourceFile(c *ImplSourceFileContext)

	// ExitModuleHeader is called when exiting the moduleHeader production.
	ExitModuleHeader(c *ModuleHeaderContext)

	// ExitImplementsHeader is called when exiting the implementsHeader production.
	ExitImplementsHeader(c *ImplementsHeaderContext)

	// ExitWorkspace is called when exiting the workspace production.
	ExitWorkspace(c *WorkspaceContext)

	// ExitPackagesSection is called when exiting the packagesSection production.
	ExitPackagesSection(c *PackagesSectionContext)

	// ExitPackageRule is called when exiting the packageRule production.
	ExitPackageRule(c *PackageRuleContext)

	// ExitImportsSection is called when exiting the importsSection production.
	ExitImportsSection(c *ImportsSectionContext)

	// ExitImplsSection is called when exiting the implsSection production.
	ExitImplsSection(c *ImplsSectionContext)

	// ExitFlagsSection is called when exiting the flagsSection production.
	ExitFlagsSection(c *FlagsSectionContext)

	// ExitModuleID is called when exiting the moduleID production.
	ExitModuleID(c *ModuleIDContext)

	// ExitFilename is called when exiting the filename production.
	ExitFilename(c *FilenameContext)

	// ExitSources is called when exiting the sources production.
	ExitSources(c *SourcesContext)

	// ExitSource is called when exiting the source production.
	ExitSource(c *SourceContext)

	// ExitTraitDecl is called when exiting the traitDecl production.
	ExitTraitDecl(c *TraitDeclContext)

	// ExitTraitMethod is called when exiting the traitMethod production.
	ExitTraitMethod(c *TraitMethodContext)

	// ExitImplBlock is called when exiting the implBlock production.
	ExitImplBlock(c *ImplBlockContext)

	// ExitDeclNoExport is called when exiting the declNoExport production.
	ExitDeclNoExport(c *DeclNoExportContext)

	// ExitFunctionNoExport is called when exiting the functionNoExport production.
	ExitFunctionNoExport(c *FunctionNoExportContext)

	// ExitFunctionArgs is called when exiting the functionArgs production.
	ExitFunctionArgs(c *FunctionArgsContext)

	// ExitArg is called when exiting the arg production.
	ExitArg(c *ArgContext)

	// ExitDecl is called when exiting the decl production.
	ExitDecl(c *DeclContext)

	// ExitUnexportedDecl is called when exiting the unexportedDecl production.
	ExitUnexportedDecl(c *UnexportedDeclContext)

	// ExitDeclNoTerm is called when exiting the declNoTerm production.
	ExitDeclNoTerm(c *DeclNoTermContext)

	// ExitTermDecl is called when exiting the termDecl production.
	ExitTermDecl(c *TermDeclContext)

	// ExitTypeDecl is called when exiting the typeDecl production.
	ExitTypeDecl(c *TypeDeclContext)

	// ExitTypeAbstraction is called when exiting the typeAbstraction production.
	ExitTypeAbstraction(c *TypeAbstractionContext)

	// ExitTvar is called when exiting the tvar production.
	ExitTvar(c *TvarContext)

	// ExitType_ is called when exiting the type_ production.
	ExitType_(c *Type_Context)

	// ExitForallType is called when exiting the forallType production.
	ExitForallType(c *ForallTypeContext)

	// ExitFunctionType is called when exiting the functionType production.
	ExitFunctionType(c *FunctionTypeContext)

	// ExitPtrType is called when exiting the ptrType production.
	ExitPtrType(c *PtrTypeContext)

	// ExitAppType is called when exiting the appType production.
	ExitAppType(c *AppTypeContext)

	// ExitPrimaryType is called when exiting the primaryType production.
	ExitPrimaryType(c *PrimaryTypeContext)

	// ExitArrayType is called when exiting the arrayType production.
	ExitArrayType(c *ArrayTypeContext)

	// ExitStructType is called when exiting the structType production.
	ExitStructType(c *StructTypeContext)

	// ExitFields is called when exiting the fields production.
	ExitFields(c *FieldsContext)

	// ExitField is called when exiting the field production.
	ExitField(c *FieldContext)

	// ExitTupleType is called when exiting the tupleType production.
	ExitTupleType(c *TupleTypeContext)

	// ExitTupleTypeArgs is called when exiting the tupleTypeArgs production.
	ExitTupleTypeArgs(c *TupleTypeArgsContext)

	// ExitVariantType is called when exiting the variantType production.
	ExitVariantType(c *VariantTypeContext)

	// ExitTags is called when exiting the tags production.
	ExitTags(c *TagsContext)

	// ExitTag is called when exiting the tag production.
	ExitTag(c *TagContext)

	// ExitExpression is called when exiting the expression production.
	ExitExpression(c *ExpressionContext)

	// ExitExpressionWithoutBlock is called when exiting the expressionWithoutBlock production.
	ExitExpressionWithoutBlock(c *ExpressionWithoutBlockContext)

	// ExitExpressionWithBlock is called when exiting the expressionWithBlock production.
	ExitExpressionWithBlock(c *ExpressionWithBlockContext)

	// ExitAssignTerm is called when exiting the assignTerm production.
	ExitAssignTerm(c *AssignTermContext)

	// ExitReturnTerm is called when exiting the returnTerm production.
	ExitReturnTerm(c *ReturnTermContext)

	// ExitIfTerm is called when exiting the ifTerm production.
	ExitIfTerm(c *IfTermContext)

	// ExitForTerm is called when exiting the forTerm production.
	ExitForTerm(c *ForTermContext)

	// ExitLambdaTerm is called when exiting the lambdaTerm production.
	ExitLambdaTerm(c *LambdaTermContext)

	// ExitMatchTerm is called when exiting the matchTerm production.
	ExitMatchTerm(c *MatchTermContext)

	// ExitMatchArms is called when exiting the matchArms production.
	ExitMatchArms(c *MatchArmsContext)

	// ExitMatchArm is called when exiting the matchArm production.
	ExitMatchArm(c *MatchArmContext)

	// ExitSetTerm is called when exiting the setTerm production.
	ExitSetTerm(c *SetTermContext)

	// ExitBlockExpr is called when exiting the blockExpr production.
	ExitBlockExpr(c *BlockExprContext)

	// ExitBlockStatements is called when exiting the blockStatements production.
	ExitBlockStatements(c *BlockStatementsContext)

	// ExitStatements is called when exiting the statements production.
	ExitStatements(c *StatementsContext)

	// ExitStatement is called when exiting the statement production.
	ExitStatement(c *StatementContext)

	// ExitLetStatement is called when exiting the letStatement production.
	ExitLetStatement(c *LetStatementContext)

	// ExitExpressionStatement is called when exiting the expressionStatement production.
	ExitExpressionStatement(c *ExpressionStatementContext)

	// ExitOperatorExpr is called when exiting the operatorExpr production.
	ExitOperatorExpr(c *OperatorExprContext)

	// ExitLogicalOrExpr is called when exiting the logicalOrExpr production.
	ExitLogicalOrExpr(c *LogicalOrExprContext)

	// ExitLogicalAndExpr is called when exiting the logicalAndExpr production.
	ExitLogicalAndExpr(c *LogicalAndExprContext)

	// ExitEqualityExpr is called when exiting the equalityExpr production.
	ExitEqualityExpr(c *EqualityExprContext)

	// ExitComparisonExpr is called when exiting the comparisonExpr production.
	ExitComparisonExpr(c *ComparisonExprContext)

	// ExitAdditiveExpr is called when exiting the additiveExpr production.
	ExitAdditiveExpr(c *AdditiveExprContext)

	// ExitMultiplicativeExpr is called when exiting the multiplicativeExpr production.
	ExitMultiplicativeExpr(c *MultiplicativeExprContext)

	// ExitUnaryExpr is called when exiting the unaryExpr production.
	ExitUnaryExpr(c *UnaryExprContext)

	// ExitApplicativeExpr is called when exiting the applicativeExpr production.
	ExitApplicativeExpr(c *ApplicativeExprContext)

	// ExitTypeApplicativeExpr is called when exiting the typeApplicativeExpr production.
	ExitTypeApplicativeExpr(c *TypeApplicativeExprContext)

	// ExitTypeApplicativeArgs is called when exiting the typeApplicativeArgs production.
	ExitTypeApplicativeArgs(c *TypeApplicativeArgsContext)

	// ExitBasePrimaryExpr is called when exiting the basePrimaryExpr production.
	ExitBasePrimaryExpr(c *BasePrimaryExprContext)

	// ExitPrimaryExpr is called when exiting the primaryExpr production.
	ExitPrimaryExpr(c *PrimaryExprContext)

	// ExitProjectionExpr is called when exiting the projectionExpr production.
	ExitProjectionExpr(c *ProjectionExprContext)

	// ExitDerefExpr is called when exiting the derefExpr production.
	ExitDerefExpr(c *DerefExprContext)

	// ExitInjectionExpr is called when exiting the injectionExpr production.
	ExitInjectionExpr(c *InjectionExprContext)

	// ExitStructExpr is called when exiting the structExpr production.
	ExitStructExpr(c *StructExprContext)

	// ExitLabelValues is called when exiting the labelValues production.
	ExitLabelValues(c *LabelValuesContext)

	// ExitLabelValue is called when exiting the labelValue production.
	ExitLabelValue(c *LabelValueContext)

	// ExitTupleExpr is called when exiting the tupleExpr production.
	ExitTupleExpr(c *TupleExprContext)

	// ExitTupleExprArgs is called when exiting the tupleExprArgs production.
	ExitTupleExprArgs(c *TupleExprArgsContext)

	// ExitVarExpr is called when exiting the varExpr production.
	ExitVarExpr(c *VarExprContext)

	// ExitId is called when exiting the id production.
	ExitId(c *IdContext)

	// ExitIdTokens is called when exiting the idTokens production.
	ExitIdTokens(c *IdTokensContext)
}
