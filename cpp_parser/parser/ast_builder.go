package parser

import (
	"math"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/ir"
)

type ASTBuilder struct {
	BasebapelVisitor
	filename string
}

func NewASTBuilder(filename string) *ASTBuilder {
	return &ASTBuilder{
		BasebapelVisitor: BasebapelVisitor{
			BaseParseTreeVisitor: &antlr.BaseParseTreeVisitor{},
		},
		filename: filename,
	}
}

var _ bapelVisitor = (*ASTBuilder)(nil)

func (b *ASTBuilder) Visit(tree antlr.ParseTree) interface{} {
	if tree == nil {
		return nil
	}
	return tree.Accept(b)
}

func posFromToken(filename string, token antlr.Token) ir.Pos {
	line := token.GetLine()
	return ir.NewLinePos(filename, line)
}

func posFromContext(filename string, ctx antlr.ParserRuleContext) ir.Pos {
	start := ctx.GetStart()
	stop := ctx.GetStop()
	return ir.NewRangePos(filename, start.GetLine(), stop.GetLine())
}

func (b *ASTBuilder) VisitBaseSourceFile(ctx *BaseSourceFileContext) interface{} {
	header := b.Visit(ctx.ModuleHeader()).(ast.SourceFileHeader)
	var imports ast.Imports
	if ctx.ImportsSection() != nil {
		imports = b.Visit(ctx.ImportsSection()).(ast.Imports)
	}
	var impls ast.Impls
	if ctx.ImplsSection() != nil {
		impls = b.Visit(ctx.ImplsSection()).(ast.Impls)
	}
	var flags ast.Flags
	if ctx.FlagsSection() != nil {
		flags = b.Visit(ctx.FlagsSection()).(ast.Flags)
	}
	var body []ast.Source
	if ctx.Sources() != nil {
		body = b.Visit(ctx.Sources()).([]ast.Source)
	}

	return ast.SourceFile{
		Header:  header,
		Imports: imports,
		Impls:   impls,
		Flags:   flags,
		Body:    body,
	}
}

func (b *ASTBuilder) VisitImplSourceFile(ctx *ImplSourceFileContext) interface{} {
	header := b.Visit(ctx.ImplementsHeader()).(ast.SourceFileHeader)
	var imports ast.Imports
	if ctx.ImportsSection() != nil {
		imports = b.Visit(ctx.ImportsSection()).(ast.Imports)
	}
	var impls ast.Impls
	if ctx.ImplsSection() != nil {
		impls = b.Visit(ctx.ImplsSection()).(ast.Impls)
	}
	var flags ast.Flags
	if ctx.FlagsSection() != nil {
		flags = b.Visit(ctx.FlagsSection()).(ast.Flags)
	}
	var body []ast.Source
	if ctx.Sources() != nil {
		body = b.Visit(ctx.Sources()).([]ast.Source)
	}

	return ast.SourceFile{
		Header:  header,
		Imports: imports,
		Impls:   impls,
		Flags:   flags,
		Body:    body,
	}
}

func (b *ASTBuilder) VisitModuleHeader(ctx *ModuleHeaderContext) interface{} {
	moduleID := b.Visit(ctx.ModuleID()).(ir.ModuleID)
	return ast.NewBaseSourceFileHeader(moduleID)
}

func (b *ASTBuilder) VisitImplementsHeader(ctx *ImplementsHeaderContext) interface{} {
	moduleID := b.Visit(ctx.ModuleID()).(ir.ModuleID)
	return ast.NewImplSourceFileHeader(moduleID)
}

func (b *ASTBuilder) VisitWorkspace(ctx *WorkspaceContext) interface{} {
	packages := b.Visit(ctx.PackagesSection()).(ast.Packages)
	return ast.NewWorkspace(packages)
}

func (b *ASTBuilder) VisitPackagesSection(ctx *PackagesSectionContext) interface{} {
	var pkgs []ast.Package
	for _, pkgCtx := range ctx.AllPackageRule() {
		pkgs = append(pkgs, b.Visit(pkgCtx).(ast.Package))
	}
	return ast.NewPackages(pkgs, posFromContext(b.filename, ctx))
}

func (b *ASTBuilder) VisitPackageRule(ctx *PackageRuleContext) interface{} {
	moduleID := b.Visit(ctx.ModuleID()).(ir.ModuleID)
	filename := b.Visit(ctx.Filename()).(ir.Filename)
	pos := posFromContext(b.filename, ctx)
	if ctx.GetStart().GetText() == "prefix" {
		return ast.NewPrefixPackage(moduleID, filename, pos)
	}
	return ast.NewModulePackage(moduleID, filename, pos)
}

func (b *ASTBuilder) VisitImportsSection(ctx *ImportsSectionContext) interface{} {
	var ids []ir.ModuleID
	for _, idCtx := range ctx.AllModuleID() {
		ids = append(ids, b.Visit(idCtx).(ir.ModuleID))
	}
	return ast.NewImports(ids, posFromContext(b.filename, ctx))
}

func (b *ASTBuilder) VisitImplsSection(ctx *ImplsSectionContext) interface{} {
	var filenames []ir.Filename
	for _, fnCtx := range ctx.AllFilename() {
		filenames = append(filenames, b.Visit(fnCtx).(ir.Filename))
	}
	return ast.NewImpls(filenames, posFromContext(b.filename, ctx))
}

func (b *ASTBuilder) VisitFlagsSection(ctx *FlagsSectionContext) interface{} {
	var filenames []ir.Filename
	for _, fnCtx := range ctx.AllFilename() {
		filenames = append(filenames, b.Visit(fnCtx).(ir.Filename))
	}
	return ast.NewFlags(filenames, posFromContext(b.filename, ctx))
}

func (b *ASTBuilder) VisitModuleID(ctx *ModuleIDContext) interface{} {
	tokens := ctx.AllIDENTIFIER()
	var bstr strings.Builder
	bstr.WriteString(tokens[0].GetText())
	for _, tok := range tokens[1:] {
		bstr.WriteString(ir.ModuleIDSeparator)
		bstr.WriteString(tok.GetText())
	}
	return ir.NewModuleID(bstr.String(), posFromContext(b.filename, ctx))
}

func (b *ASTBuilder) VisitFilename(ctx *FilenameContext) interface{} {
	literalStr := ctx.STRING_LITERAL().GetText()
	literalStr = strings.TrimPrefix(literalStr, `"`)
	literalStr = strings.TrimSuffix(literalStr, `"`)
	return ir.NewFilename(literalStr, posFromContext(b.filename, ctx))
}

func (b *ASTBuilder) VisitSources(ctx *SourcesContext) interface{} {
	var sources []ast.Source
	for _, srcCtx := range ctx.AllSource() {
		sources = append(sources, b.Visit(srcCtx).(ast.Source))
	}
	return sources
}

func (b *ASTBuilder) VisitSource(ctx *SourceContext) interface{} {
	if ctx.DeclNoExport() != nil {
		return b.Visit(ctx.DeclNoExport()).(ast.Source)
	}
	if ctx.FunctionNoExport() != nil {
		fnSrc := b.Visit(ctx.FunctionNoExport()).(ast.Source)
		if ctx.GetStart().GetText() == "pub" {
			fnSrc.Function.Export = true
		}
		return fnSrc
	}
	if ctx.TraitDecl() != nil {
		trait := b.Visit(ctx.TraitDecl()).(ast.Trait)
		if ctx.GetStart().GetText() == "pub" {
			trait.Export = true
		}
		return ast.NewTraitSource(trait)
	}
	if ctx.ImplBlock() != nil {
		return b.Visit(ctx.ImplBlock()).(ast.Source)
	}
	panic("unreachable")
}

func (b *ASTBuilder) VisitDeclNoExport(ctx *DeclNoExportContext) interface{} {
	var decl ir.IrDecl
	if ctx.DeclNoTerm() != nil {
		decl = b.Visit(ctx.DeclNoTerm()).(ir.IrDecl)
	} else {
		decl = b.Visit(ctx.TermDecl()).(ir.IrDecl)
		if ctx.GetStart().GetText() == "pub" {
			decl.Export = true
		}
	}
	return ast.NewDeclSource(decl)
}

func (b *ASTBuilder) VisitFunctionNoExport(ctx *FunctionNoExportContext) interface{} {
	id := b.Visit(ctx.Id()).(ast.ID)

	var tvars []ir.VarKind
	if ctx.TypeAbstraction() != nil {
		tvars = b.Visit(ctx.TypeAbstraction()).([]ir.VarKind)
	}

	funArgs := b.Visit(ctx.FunctionArgs()).([]ir.FunctionArg)
	retType := b.Visit(ctx.Type_()).(ir.IrType)
	body := b.Visit(ctx.BlockExpr()).(ast.Expr)

	return ast.NewFunctionSource(
		ast.NewFunction(
			posFromContext(b.filename, ctx),
			false /* export */, id.Value, tvars, funArgs, retType, body))
}

func (b *ASTBuilder) VisitFunctionArgs(ctx *FunctionArgsContext) interface{} {
	var args []ir.FunctionArg
	for _, argCtx := range ctx.AllArg() {
		args = append(args, b.Visit(argCtx).(ir.FunctionArg))
	}
	return args
}

func (b *ASTBuilder) VisitArg(ctx *ArgContext) interface{} {
	idToken := ctx.IDENTIFIER().GetSymbol()
	typ := b.Visit(ctx.Type_()).(ir.IrType)
	return ir.FunctionArg{ID: idToken.GetText(), Type: typ}
}

func (b *ASTBuilder) VisitDecl(ctx *DeclContext) interface{} {
	decl := b.Visit(ctx.UnexportedDecl()).(ir.IrDecl)
	if ctx.GetStart().GetText() == "pub" {
		decl.Export = true
	}
	return decl
}

func (b *ASTBuilder) VisitUnexportedDecl(ctx *UnexportedDeclContext) interface{} {
	if ctx.TermDecl() != nil {
		return b.Visit(ctx.TermDecl()).(ir.IrDecl)
	}
	return b.Visit(ctx.TypeDecl()).(ir.IrDecl)
}

func (b *ASTBuilder) VisitDeclNoTerm(ctx *DeclNoTermContext) interface{} {
	decl := b.Visit(ctx.TypeDecl()).(ir.IrDecl)
	if ctx.GetStart().GetText() == "pub" {
		decl.Export = true
	}
	return decl
}

func (b *ASTBuilder) VisitTermDecl(ctx *TermDeclContext) interface{} {
	id := b.Visit(ctx.Id()).(ast.ID)
	typ := b.Visit(ctx.Type_()).(ir.IrType)
	qtyp := newQuantifiedType(typ)
	return newTermDecl(id, qtyp, false /* export */)
}

func (b *ASTBuilder) VisitTypeDecl(ctx *TypeDeclContext) interface{} {
	id := b.Visit(ctx.Id()).(ast.ID)

	var tvars []ir.VarKind
	if ctx.TypeAbstraction() != nil {
		tvars = b.Visit(ctx.TypeAbstraction()).([]ir.VarKind)
	}

	kind := ir.NewTypeKind()
	for i := 0; i < len(tvars); i++ {
		kind = ir.NewArrowKind(ir.NewTypeKind(), kind)
	}

	if ctx.Type_() != nil {
		typ := b.Visit(ctx.Type_()).(ir.IrType)
		return newAliasDecl(id, kind, ir.LambdaVars(tvars, typ), false /* export */)
	} else {
		return newNameDecl(id, kind, false /* export */)
	}
}

func (b *ASTBuilder) VisitTypeAbstraction(ctx *TypeAbstractionContext) interface{} {
	var tvars []ir.VarKind
	for _, tvarCtx := range ctx.AllTvar() {
		tvars = append(tvars, b.Visit(tvarCtx).(ir.VarKind))
	}
	return tvars
}

func (b *ASTBuilder) VisitTvar(ctx *TvarContext) interface{} {
	idToken := ctx.IDENTIFIER().GetSymbol()
	return ir.VarKind{Var: idToken.GetText(), Kind: ir.NewTypeKind()}
}

func (b *ASTBuilder) VisitType_(ctx *Type_Context) interface{} {
	return b.Visit(ctx.ForallType())
}

func (b *ASTBuilder) VisitForallType(ctx *ForallTypeContext) interface{} {
	if ctx.TypeAbstraction() != nil {
		tvars := b.Visit(ctx.TypeAbstraction()).([]ir.VarKind)
		subType := b.Visit(ctx.FunctionType()).(ir.IrType)
		return newForallType(posFromContext(b.filename, ctx), tvars, subType)
	}
	return b.Visit(ctx.FunctionType())
}

func (b *ASTBuilder) VisitFunctionType(ctx *FunctionTypeContext) interface{} {
	if ctx.FunctionType() != nil {
		arg := b.Visit(ctx.PtrType()).(ir.IrType)
		ret := b.Visit(ctx.FunctionType()).(ir.IrType)
		return newFunctionType(arg, ret)
	}
	return b.Visit(ctx.PtrType())
}

func (b *ASTBuilder) VisitAppType(ctx *AppTypeContext) interface{} {
	if ctx.AppType() != nil {
		fun := b.Visit(ctx.AppType()).(ir.IrType)
		arg := b.Visit(ctx.PrimaryType()).(ir.IrType)
		return newAppType(fun, arg)
	}
	return b.Visit(ctx.PrimaryType())
}

func (b *ASTBuilder) VisitPtrType(ctx *PtrTypeContext) interface{} {
	if ctx.AMP() != nil {
		id := ast.NewID("Ptr", posFromToken(b.filename, ctx.AMP().GetSymbol()))
		typ := b.Visit(ctx.PtrType()).(ir.IrType)
		return newAppType(newNameType(id), typ)
	}
	return b.Visit(ctx.AppType())
}

func (b *ASTBuilder) VisitPrimaryType(ctx *PrimaryTypeContext) interface{} {
	if ctx.ArrayType() != nil {
		return b.Visit(ctx.ArrayType())
	}
	if ctx.StructType() != nil {
		return b.Visit(ctx.StructType())
	}
	if ctx.TupleType() != nil {
		return b.Visit(ctx.TupleType())
	}
	if ctx.VariantType() != nil {
		return b.Visit(ctx.VariantType())
	}
	if ctx.SINGLE_QUOTE() != nil {
		idToken := ctx.IDENTIFIER().GetSymbol()
		id := ast.NewID(idToken.GetText(), posFromToken(b.filename, idToken))
		return newVarType(id)
	}
	if ctx.Id() != nil {
		id := b.Visit(ctx.Id()).(ast.ID)
		return newNameType(id)
	}
	return b.Visit(ctx.Type_())
}

func (b *ASTBuilder) VisitArrayType(ctx *ArrayTypeContext) interface{} {
	elemType := b.Visit(ctx.Type_()).(ir.IrType)
	length := math.MaxInt
	if ctx.INT_LITERAL() != nil {
		val, err := parseNumber[int64](ctx.INT_LITERAL().GetText())
		if err != nil {
			panic(err)
		}
		length = int(val)
		if ctx.MINUS() != nil {
			length = -length
		}
	}
	return newArrayType(posFromContext(b.filename, ctx), elemType, length)
}

func (b *ASTBuilder) VisitStructType(ctx *StructTypeContext) interface{} {
	var fields []ir.StructField
	if ctx.Fields() != nil {
		fields = b.Visit(ctx.Fields()).([]ir.StructField)
	}
	return newStructType(posFromContext(b.filename, ctx), fields)
}

func (b *ASTBuilder) VisitFields(ctx *FieldsContext) interface{} {
	var fields []ir.StructField
	for _, fieldCtx := range ctx.AllField() {
		fields = append(fields, b.Visit(fieldCtx).(ir.StructField))
	}
	return fields
}

func (b *ASTBuilder) VisitField(ctx *FieldContext) interface{} {
	id := b.Visit(ctx.Id()).(ast.ID)
	typ := b.Visit(ctx.Type_()).(ir.IrType)
	return ir.StructField{ID: id.Value, Type: typ}
}

func (b *ASTBuilder) VisitTupleType(ctx *TupleTypeContext) interface{} {
	var elems []ir.IrType
	if ctx.TupleTypeArgs() != nil {
		elems = b.Visit(ctx.TupleTypeArgs()).([]ir.IrType)
	}
	return newTupleType(posFromContext(b.filename, ctx), elems)
}

func (b *ASTBuilder) VisitTupleTypeArgs(ctx *TupleTypeArgsContext) interface{} {
	var elems []ir.IrType
	for _, typeCtx := range ctx.AllType_() {
		elems = append(elems, b.Visit(typeCtx).(ir.IrType))
	}
	return elems
}

func (b *ASTBuilder) VisitVariantType(ctx *VariantTypeContext) interface{} {
	var tags []ir.VariantTag
	if ctx.Tags() != nil {
		tags = b.Visit(ctx.Tags()).([]ir.VariantTag)
	}
	return newVariantType(posFromContext(b.filename, ctx), tags)
}

func (b *ASTBuilder) VisitTags(ctx *TagsContext) interface{} {
	var tags []ir.VariantTag
	for _, tagCtx := range ctx.AllTag() {
		tags = append(tags, b.Visit(tagCtx).(ir.VariantTag))
	}
	return tags
}

func (b *ASTBuilder) VisitTag(ctx *TagContext) interface{} {
	id := b.Visit(ctx.Id()).(ast.ID)
	typ := b.Visit(ctx.Type_()).(ir.IrType)
	return ir.VariantTag{ID: id.Value, Type: typ}
}

func (b *ASTBuilder) VisitExpression(ctx *ExpressionContext) interface{} {
	if ctx.ExpressionWithoutBlock() != nil {
		return b.Visit(ctx.ExpressionWithoutBlock())
	}
	return b.Visit(ctx.ExpressionWithBlock())
}

func (b *ASTBuilder) VisitExpressionWithoutBlock(ctx *ExpressionWithoutBlockContext) interface{} {
	if ctx.AssignTerm() != nil {
		return b.Visit(ctx.AssignTerm())
	}
	if ctx.OperatorExpr() != nil {
		return b.Visit(ctx.OperatorExpr())
	}
	return b.Visit(ctx.ReturnTerm())
}

func (b *ASTBuilder) VisitExpressionWithBlock(ctx *ExpressionWithBlockContext) interface{} {
	if ctx.BlockExpr() != nil {
		return b.Visit(ctx.BlockExpr())
	}
	if ctx.IfTerm() != nil {
		return b.Visit(ctx.IfTerm())
	}
	if ctx.ForTerm() != nil {
		return b.Visit(ctx.ForTerm())
	}
	if ctx.LambdaTerm() != nil {
		return b.Visit(ctx.LambdaTerm())
	}
	if ctx.MatchTerm() != nil {
		return b.Visit(ctx.MatchTerm())
	}
	return b.Visit(ctx.SetTerm())
}

func (b *ASTBuilder) VisitAssignTerm(ctx *AssignTermContext) interface{} {
	var ret ast.Expr
	if ctx.Id() != nil {
		id := b.Visit(ctx.Id()).(ast.ID)
		ret = ast.NewVarExpr(id)
	} else {
		ret = b.Visit(ctx.TupleExpr()).(ast.Expr)
	}
	arg := b.Visit(ctx.Expression()).(ast.Expr)
	return newAssignExpr(arg, ret)
}

func (b *ASTBuilder) VisitReturnTerm(ctx *ReturnTermContext) interface{} {
	expr := b.Visit(ctx.ExpressionWithoutBlock()).(ast.Expr)
	return ast.NewReturnExpr(posFromContext(b.filename, ctx), expr)
}

func (b *ASTBuilder) VisitIfTerm(ctx *IfTermContext) interface{} {
	condition := b.Visit(ctx.ExpressionWithoutBlock()).(ast.Expr)
	then := b.Visit(ctx.BlockExpr(0)).(ast.Expr)
	var elseExpr *ast.Expr
	if ctx.BlockExpr(1) != nil {
		ee := b.Visit(ctx.BlockExpr(1)).(ast.Expr)
		elseExpr = &ee
	} else if ctx.IfTerm() != nil {
		ee := b.Visit(ctx.IfTerm()).(ast.Expr)
		elseExpr = &ee
	}
	return newIfExpr(posFromContext(b.filename, ctx), condition, then, elseExpr)
}

func (b *ASTBuilder) VisitForTerm(ctx *ForTermContext) interface{} {
	condition := b.Visit(ctx.ExpressionWithoutBlock()).(ast.Expr)
	body := b.Visit(ctx.BlockExpr()).(ast.Expr)
	return ast.NewForExpr(posFromContext(b.filename, ctx), condition, body)
}

func (b *ASTBuilder) VisitLambdaTerm(ctx *LambdaTermContext) interface{} {
	var tvars []ir.VarKind
	if ctx.TypeAbstraction() != nil {
		tvars = b.Visit(ctx.TypeAbstraction()).([]ir.VarKind)
	}
	funArgs := b.Visit(ctx.FunctionArgs()).([]ir.FunctionArg)
	body := b.Visit(ctx.BlockExpr()).(ast.Expr)
	return ast.Lambda(posFromContext(b.filename, ctx), tvars, funArgs, body)
}

func (b *ASTBuilder) VisitMatchTerm(ctx *MatchTermContext) interface{} {
	expr := b.Visit(ctx.Expression()).(ast.Expr)
	arms := b.Visit(ctx.MatchArms()).([]ast.MatchArm)
	return ast.NewMatchExpr(posFromContext(b.filename, ctx), expr, arms)
}

func (b *ASTBuilder) VisitMatchArms(ctx *MatchArmsContext) interface{} {
	var arms []ast.MatchArm
	for _, armCtx := range ctx.AllMatchArm() {
		arms = append(arms, b.Visit(armCtx).(ast.MatchArm))
	}
	return arms
}

func (b *ASTBuilder) VisitMatchArm(ctx *MatchArmContext) interface{} {
	tag := b.Visit(ctx.Id()).(ast.ID)
	argToken := ctx.IDENTIFIER().GetSymbol()
	arg := ast.NewID(argToken.GetText(), posFromToken(b.filename, argToken))
	body := b.Visit(ctx.Expression()).(ast.Expr)
	return newMatchArm(tag, arg, body)
}

func (b *ASTBuilder) VisitSetTerm(ctx *SetTermContext) interface{} {
	expr := b.Visit(ctx.Expression()).(ast.Expr)
	values := b.Visit(ctx.LabelValues()).([]ast.LabelValue)
	return ast.NewSetExpr(posFromContext(b.filename, ctx), expr, values)
}

func (b *ASTBuilder) VisitBlockExpr(ctx *BlockExprContext) interface{} {
	exprs := b.Visit(ctx.BlockStatements()).([]ast.Expr)
	return ast.NewBlockExpr(posFromContext(b.filename, ctx), exprs)
}

func (b *ASTBuilder) VisitBlockStatements(ctx *BlockStatementsContext) interface{} {
	if ctx.Statements() != nil {
		exprs := b.Visit(ctx.Statements()).([]ast.Expr)
		if ctx.ExpressionWithoutBlock() != nil {
			exprs = append(exprs, b.Visit(ctx.ExpressionWithoutBlock()).(ast.Expr))
		}
		return exprs
	}
	return []ast.Expr{b.Visit(ctx.ExpressionWithoutBlock()).(ast.Expr)}
}

func (b *ASTBuilder) VisitStatements(ctx *StatementsContext) interface{} {
	var exprs []ast.Expr
	for _, stmtCtx := range ctx.AllStatement() {
		exprs = append(exprs, b.Visit(stmtCtx).(ast.Expr))
	}
	return exprs
}

func (b *ASTBuilder) VisitStatement(ctx *StatementContext) interface{} {
	if ctx.LetStatement() != nil {
		return b.Visit(ctx.LetStatement())
	}
	return b.Visit(ctx.ExpressionStatement())
}

func (b *ASTBuilder) VisitLetStatement(ctx *LetStatementContext) interface{} {
	id := b.Visit(ctx.Id()).(ast.ID)
	value := b.Visit(ctx.Expression()).(ast.Expr)
	var varType *ir.IrType
	if ctx.Type_() != nil {
		t := b.Visit(ctx.Type_()).(ir.IrType)
		varType = &t
	}
	return newLetExpr(id, varType, value)
}

func (b *ASTBuilder) VisitExpressionStatement(ctx *ExpressionStatementContext) interface{} {
	if ctx.ExpressionWithoutBlock() != nil {
		return b.Visit(ctx.ExpressionWithoutBlock())
	}
	return b.Visit(ctx.ExpressionWithBlock())
}

func (b *ASTBuilder) VisitOperatorExpr(ctx *OperatorExprContext) interface{} {
	return b.Visit(ctx.LogicalOrExpr())
}

func (b *ASTBuilder) VisitLogicalOrExpr(ctx *LogicalOrExprContext) interface{} {
	if ctx.LogicalOrExpr() != nil {
		left := b.Visit(ctx.LogicalOrExpr()).(ast.Expr)
		right := b.Visit(ctx.LogicalAndExpr()).(ast.Expr)
		opToken := ctx.OR().GetSymbol()
		return newBinOpExpr(
			ast.NewVarExpr(ast.NewID(opToken.GetText(), posFromToken(b.filename, opToken))),
			nil, left, right)
	}
	return b.Visit(ctx.LogicalAndExpr())
}

func (b *ASTBuilder) VisitLogicalAndExpr(ctx *LogicalAndExprContext) interface{} {
	if ctx.LogicalAndExpr() != nil {
		left := b.Visit(ctx.LogicalAndExpr()).(ast.Expr)
		right := b.Visit(ctx.EqualityExpr()).(ast.Expr)
		opToken := ctx.AND().GetSymbol()
		return newBinOpExpr(
			ast.NewVarExpr(ast.NewID(opToken.GetText(), posFromToken(b.filename, opToken))),
			nil, left, right)
	}
	return b.Visit(ctx.EqualityExpr())
}

func (b *ASTBuilder) VisitEqualityExpr(ctx *EqualityExprContext) interface{} {
	if ctx.EqualityExpr() != nil {
		left := b.Visit(ctx.EqualityExpr()).(ast.Expr)
		right := b.Visit(ctx.ComparisonExpr()).(ast.Expr)
		var opToken antlr.Token
		if ctx.NE() != nil {
			opToken = ctx.NE().GetSymbol()
		} else {
			opToken = ctx.EQ().GetSymbol()
		}
		var typeArgs []ir.IrType
		if ctx.TypeApplicativeArgs() != nil {
			typeArgs = b.Visit(ctx.TypeApplicativeArgs()).([]ir.IrType)
		}
		return newBinOpExpr(
			ast.NewVarExpr(ast.NewID(opToken.GetText(), posFromToken(b.filename, opToken))),
			typeArgs, left, right)
	}
	return b.Visit(ctx.ComparisonExpr())
}

func (b *ASTBuilder) VisitComparisonExpr(ctx *ComparisonExprContext) interface{} {
	if ctx.ComparisonExpr() != nil {
		left := b.Visit(ctx.ComparisonExpr()).(ast.Expr)
		right := b.Visit(ctx.AdditiveExpr()).(ast.Expr)
		var opToken antlr.Token
		if ctx.GT() != nil {
			opToken = ctx.GT().GetSymbol()
		} else if ctx.GE() != nil {
			opToken = ctx.GE().GetSymbol()
		} else if ctx.LT() != nil {
			opToken = ctx.LT().GetSymbol()
		} else {
			opToken = ctx.LE().GetSymbol()
		}
		var typeArgs []ir.IrType
		if ctx.TypeApplicativeArgs() != nil {
			typeArgs = b.Visit(ctx.TypeApplicativeArgs()).([]ir.IrType)
		}
		return newBinOpExpr(
			ast.NewVarExpr(ast.NewID(opToken.GetText(), posFromToken(b.filename, opToken))),
			typeArgs, left, right)
	}
	return b.Visit(ctx.AdditiveExpr())
}

func (b *ASTBuilder) VisitAdditiveExpr(ctx *AdditiveExprContext) interface{} {
	if ctx.AdditiveExpr() != nil {
		left := b.Visit(ctx.AdditiveExpr()).(ast.Expr)
		right := b.Visit(ctx.MultiplicativeExpr()).(ast.Expr)
		var opToken antlr.Token
		if ctx.PLUS() != nil {
			opToken = ctx.PLUS().GetSymbol()
		} else {
			opToken = ctx.MINUS().GetSymbol()
		}
		var typeArgs []ir.IrType
		if ctx.TypeApplicativeArgs() != nil {
			typeArgs = b.Visit(ctx.TypeApplicativeArgs()).([]ir.IrType)
		}
		return newBinOpExpr(
			ast.NewVarExpr(ast.NewID(opToken.GetText(), posFromToken(b.filename, opToken))),
			typeArgs, left, right)
	}
	return b.Visit(ctx.MultiplicativeExpr())
}

func (b *ASTBuilder) VisitMultiplicativeExpr(ctx *MultiplicativeExprContext) interface{} {
	if ctx.MultiplicativeExpr() != nil {
		left := b.Visit(ctx.MultiplicativeExpr()).(ast.Expr)
		right := b.Visit(ctx.UnaryExpr()).(ast.Expr)
		var opToken antlr.Token
		if ctx.MUL() != nil {
			opToken = ctx.MUL().GetSymbol()
		} else {
			opToken = ctx.DIV().GetSymbol()
		}
		var typeArgs []ir.IrType
		if ctx.TypeApplicativeArgs() != nil {
			typeArgs = b.Visit(ctx.TypeApplicativeArgs()).([]ir.IrType)
		}
		return newBinOpExpr(
			ast.NewVarExpr(ast.NewID(opToken.GetText(), posFromToken(b.filename, opToken))),
			typeArgs, left, right)
	}
	return b.Visit(ctx.UnaryExpr())
}

func (b *ASTBuilder) VisitUnaryExpr(ctx *UnaryExprContext) interface{} {
	if ctx.UnaryExpr() != nil {
		operand := b.Visit(ctx.UnaryExpr()).(ast.Expr)
		var opToken antlr.Token
		if ctx.NOT() != nil {
			opToken = ctx.NOT().GetSymbol()
		} else {
			opToken = ctx.MINUS().GetSymbol()
		}
		var typeArgs []ir.IrType
		if ctx.TypeApplicativeArgs() != nil {
			typeArgs = b.Visit(ctx.TypeApplicativeArgs()).([]ir.IrType)
		}
		return newUnaryOpExpr(
			ast.NewVarExpr(ast.NewID(opToken.GetText(), posFromToken(b.filename, opToken))),
			typeArgs, operand)
	}
	return b.Visit(ctx.ApplicativeExpr())
}

func (b *ASTBuilder) VisitApplicativeExpr(ctx *ApplicativeExprContext) interface{} {
	if ctx.ApplicativeExpr() != nil {
		fun := b.Visit(ctx.ApplicativeExpr()).(ast.Expr)
		arg := b.Visit(ctx.BasePrimaryExpr()).(ast.Expr)
		return ast.NewAppTermExpr(makePos(fun.Pos, arg.Pos), fun, arg)
	}
	return b.Visit(ctx.TypeApplicativeExpr())
}

func (b *ASTBuilder) VisitTypeApplicativeExpr(ctx *TypeApplicativeExprContext) interface{} {
	primary := b.Visit(ctx.PrimaryExpr()).(ast.Expr)
	if ctx.TypeApplicativeArgs() != nil {
		typeArgs := b.Visit(ctx.TypeApplicativeArgs()).([]ir.IrType)
		return newAppTypeExpr(primary, typeArgs)
	}
	return primary
}

func (b *ASTBuilder) VisitTypeApplicativeArgs(ctx *TypeApplicativeArgsContext) interface{} {
	if ctx.TupleTypeArgs() != nil {
		return b.Visit(ctx.TupleTypeArgs())
	}
	return []ir.IrType{b.Visit(ctx.Type_()).(ir.IrType)}
}

func (b *ASTBuilder) VisitPrimaryExpr(ctx *PrimaryExprContext) interface{} {
	if ctx.MUL() != nil {
		mulPos := posFromToken(b.filename, ctx.MUL().GetSymbol())
		id := ast.NewVarExpr(ast.NewID("Ptr::get", mulPos))
		expr := b.Visit(ctx.PrimaryExpr()).(ast.Expr)
		return ast.Call(makePos(id.Pos, expr.Pos), id, nil /* typeArgs */, expr)
	}
	if ctx.BasePrimaryExpr() != nil {
		return b.Visit(ctx.BasePrimaryExpr())
	}
	return nil
}

func (b *ASTBuilder) VisitBasePrimaryExpr(ctx *BasePrimaryExprContext) interface{} {
	if ctx.AMP() != nil {
		ampPos := posFromToken(b.filename, ctx.AMP().GetSymbol())
		id := ast.NewVarExpr(ast.NewID("Ptr::mk", ampPos))
		targetID := b.Visit(ctx.Id()).(ast.ID)
		varExpr := ast.NewVarExpr(targetID)
		return ast.Call(makePos(id.Pos, varExpr.Pos), id, nil /* typeArgs */, varExpr)
	}
	if ctx.ProjectionExpr() != nil {
		return b.Visit(ctx.ProjectionExpr())
	}
	if ctx.INT_LITERAL() != nil {
		valToken := ctx.INT_LITERAL().GetSymbol()
		tok := Token{Pos: posFromToken(b.filename, valToken), Text: valToken.GetText()}
		return ast.NewConstExpr(newNumberLiteral(tok))
	}
	if ctx.FLOAT_LITERAL() != nil {
		valToken := ctx.FLOAT_LITERAL().GetSymbol()
		tok := Token{Pos: posFromToken(b.filename, valToken), Text: valToken.GetText()}
		return ast.NewConstExpr(newNumberLiteral(tok))
	}
	return nil
}

func (b *ASTBuilder) VisitProjectionExpr(ctx *ProjectionExprContext) interface{} {
	if ctx.ProjectionExpr() != nil {
		expr := b.Visit(ctx.ProjectionExpr()).(ast.Expr)
		var label string
		var labelPos ir.Pos
		if ctx.INT_LITERAL() != nil {
			tok := ctx.INT_LITERAL().GetSymbol()
			label = tok.GetText()
			labelPos = posFromToken(b.filename, tok)
		} else {
			tok := ctx.IDENTIFIER().GetSymbol()
			label = tok.GetText()
			labelPos = posFromToken(b.filename, tok)
		}
		return ast.NewProjectionExpr(makePos(expr.Pos, labelPos), expr, label)
	}
	return b.Visit(ctx.DerefExpr())
}

func (b *ASTBuilder) VisitDerefExpr(ctx *DerefExprContext) interface{} {
	if ctx.InjectionExpr() != nil {
		return b.Visit(ctx.InjectionExpr())
	}
	if ctx.RUNE_LITERAL() != nil {
		tok := ctx.RUNE_LITERAL().GetSymbol()
		t := Token{Pos: posFromToken(b.filename, tok), Text: tok.GetText()}
		return ast.NewConstExpr(newRuneLiteral(t))
	}
	if ctx.STRING_LITERAL() != nil {
		tok := ctx.STRING_LITERAL().GetSymbol()
		t := Token{Pos: posFromToken(b.filename, tok), Text: tok.GetText()}
		return ast.NewConstExpr(newStringLiteral(t))
	}
	if ctx.StructExpr() != nil {
		return b.Visit(ctx.StructExpr())
	}
	if ctx.TupleExpr() != nil {
		return b.Visit(ctx.TupleExpr())
	}
	if ctx.VarExpr() != nil {
		return b.Visit(ctx.VarExpr())
	}
	return b.Visit(ctx.Expression())
}

func (b *ASTBuilder) VisitInjectionExpr(ctx *InjectionExprContext) interface{} {
	variantType := b.Visit(ctx.Type_()).(ir.IrType)
	labelValue := b.Visit(ctx.LabelValue()).(ast.LabelValue)
	return ast.NewInjectionExpr(posFromContext(b.filename, ctx), variantType, labelValue.Label, labelValue.Value)
}

func (b *ASTBuilder) VisitStructExpr(ctx *StructExprContext) interface{} {
	var values []ast.LabelValue
	if ctx.LabelValues() != nil {
		values = b.Visit(ctx.LabelValues()).([]ast.LabelValue)
	}
	return ast.NewStructExpr(posFromContext(b.filename, ctx), values)
}

func (b *ASTBuilder) VisitLabelValues(ctx *LabelValuesContext) interface{} {
	var values []ast.LabelValue
	for _, lvCtx := range ctx.AllLabelValue() {
		values = append(values, b.Visit(lvCtx).(ast.LabelValue))
	}
	return values
}

func (b *ASTBuilder) VisitLabelValue(ctx *LabelValueContext) interface{} {
	value := b.Visit(ctx.Expression()).(ast.Expr)
	if ctx.Id() != nil {
		id := b.Visit(ctx.Id()).(ast.ID)
		return ast.LabelValue{Label: id.Value, Value: value}
	}
	labelVal := ctx.INT_LITERAL().GetText()
	return ast.LabelValue{Label: labelVal, Value: value}
}

func (b *ASTBuilder) VisitTupleExpr(ctx *TupleExprContext) interface{} {
	var elems []ast.Expr
	if ctx.TupleExprArgs() != nil {
		elems = b.Visit(ctx.TupleExprArgs()).([]ast.Expr)
	}
	return ast.NewTupleExpr(posFromContext(b.filename, ctx), elems)
}

func (b *ASTBuilder) VisitTupleExprArgs(ctx *TupleExprArgsContext) interface{} {
	var elems []ast.Expr
	for _, exprCtx := range ctx.AllExpression() {
		elems = append(elems, b.Visit(exprCtx).(ast.Expr))
	}
	return elems
}

func (b *ASTBuilder) VisitVarExpr(ctx *VarExprContext) interface{} {
	id := b.Visit(ctx.Id()).(ast.ID)
	return ast.NewVarExpr(id)
}

func (b *ASTBuilder) VisitId(ctx *IdContext) interface{} {
	if ctx.IdTokens() != nil {
		return b.Visit(ctx.IdTokens())
	}
	opTok := ctx.GetChild(1).(antlr.TerminalNode).GetSymbol()
	return ast.NewID(opTok.GetText(), posFromToken(b.filename, opTok))
}

func (b *ASTBuilder) VisitIdTokens(ctx *IdTokensContext) interface{} {
	return ast.NewID(ctx.GetText(), posFromContext(b.filename, ctx))
}

func (b *ASTBuilder) VisitTraitDecl(ctx *TraitDeclContext) interface{} {
	id := b.Visit(ctx.Id()).(ast.ID)
	var methods []ast.Signature
	for _, methodCtx := range ctx.AllTraitMethod() {
		methods = append(methods, b.Visit(methodCtx).(ast.Signature))
	}
	return ast.NewTrait(posFromContext(b.filename, ctx), false /* export */, id.Value, methods)
}

func (b *ASTBuilder) VisitTraitMethod(ctx *TraitMethodContext) interface{} {
	id := b.Visit(ctx.Id()).(ast.ID)
	funArgs := b.Visit(ctx.FunctionArgs()).([]ir.FunctionArg)
	retType := b.Visit(ctx.Type_()).(ir.IrType)
	return ast.NewSignature(posFromContext(b.filename, ctx), id.Value, funArgs, retType)
}

func (b *ASTBuilder) VisitImplBlock(ctx *ImplBlockContext) interface{} {
	traitId := b.Visit(ctx.Id()).(ast.ID)
	targetType := b.Visit(ctx.Type_()).(ir.IrType)
	var methods []ast.Function
	for _, fnCtx := range ctx.AllFunctionNoExport() {
		fnSrc := b.Visit(fnCtx).(ast.Source)
		methods = append(methods, fnSrc.Function.Function)
	}
	return ast.NewImplSource(ast.NewImpl(posFromContext(b.filename, ctx), traitId.Value, targetType, methods))
}

