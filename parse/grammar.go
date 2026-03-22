package parse

import (
	"fmt"
	"math"
	"strings"
	"unicode"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/go-lalr1/grammar"
	"golang.org/x/exp/constraints"
)

func parseNumber[T constraints.Integer](arg string) (T, error) {
	var value T

	if strings.HasPrefix(arg, "0x") {
		// Hexadecimal
		_, err := fmt.Sscanf(arg, "0x%x", &value)

		return value, err
	}

	// Decimal.
	_, err := fmt.Sscanf(arg, "%d", &value)
	return value, err
}

func parseFloat[T constraints.Integer](arg string) (T, T, error) {
	splits := strings.Split(arg, ".")
	if len(splits) != 2 {
		return 0, 0, fmt.Errorf("invalid floating point number %q", arg)
	}

	integer, err := parseNumber[T](splits[0])
	if err != nil {
		return 0, 0, err
	}

	decimal, err := parseNumber[T](splits[1])
	if err != nil {
		return 0, 0, err
	}

	return integer, decimal, nil
}

func makePos(pos1, pos2 ir.Pos) ir.Pos {
	return ir.NewRangePos(pos1.Filename, pos1.BeginLineNum, pos2.EndLineNum)
}

func makePos2(args []any) ir.Pos {
	return makePos(args[0].(Token).Pos, args[len(args)-1].(Token).Pos)
}

func newUnaryOpExpr(id ast.Expr, typeArgs []ir.IrType, expr ast.Expr) (r ast.Expr) {
	pos := makePos(id.Pos, expr.Pos)

	if id.Is(ast.VarExpr) && id.Var.ID == "-" {
		// 0 - $expr
		return ast.Call(pos, id, typeArgs, ast.NewConstExpr(newIntLiteral(pos, 0)), expr)
	}

	return ast.Call(pos, id, typeArgs, expr)
}

func newBinOpExpr(id ast.Expr, typeArgs []ir.IrType, t1, t2 ast.Expr) ast.Expr {
	return ast.Call(makePos(t1.Pos, t2.Pos), id, typeArgs, t1, t2)
}

func newAliasDecl(id ast.ID, kind ir.IrKind, typ ir.IrType, export bool) ir.IrDecl {
	decl := ir.NewAliasDecl(id.Value, kind, typ, export)
	decl.Pos = makePos(id.Pos, typ.Pos)
	return decl
}

func newTermDecl(id ast.ID, typ ir.IrType, export bool) ir.IrDecl {
	decl := ir.NewTermDecl(id.Value, typ, export)
	decl.Pos = makePos(id.Pos, typ.Pos)
	return decl
}

func newNameDecl(id ast.ID, kind ir.IrKind, export bool) ir.IrDecl {
	decl := ir.NewNameDecl(id.Value, kind, export)
	decl.Pos = id.Pos
	return decl
}

func newQuantifiedType(typ ir.IrType) ir.IrType {
	quantified := ir.QuantifyType(typ)
	quantified.Pos = typ.Pos
	return quantified
}

func newForallType(pos ir.Pos, tvars []ir.VarKind, subType ir.IrType) ir.IrType {
	forall := ir.ForallVars(tvars, subType)
	forall.Pos = pos
	return forall
}

func newFunctionType(arg, ret ir.IrType) ir.IrType {
	typ := ir.NewFunctionType(arg, ret)
	typ.Pos = makePos(arg.Pos, ret.Pos)
	return typ
}

func newAppType(arg, ret ir.IrType) ir.IrType {
	typ := ir.NewAppType(arg, ret)
	typ.Pos = makePos(arg.Pos, ret.Pos)
	return typ
}

func newVarType(id ast.ID) ir.IrType {
	typ := ir.NewVarType(id.Value)
	typ.Pos = id.Pos
	return typ
}

func newNameType(id ast.ID) ir.IrType {
	typ := ir.NewNameType(id.Value)
	typ.Pos = id.Pos
	return typ
}

func newIntLiteral(pos ir.Pos, value int64) ir.IrLiteral {
	lit := ir.NewIntLiteral(value)
	lit.Pos = pos
	return lit
}

func newFloatLiteral(pos ir.Pos, integer, decimal int64) ir.IrLiteral {
	lit := ir.NewFloatLiteral(integer, decimal)
	lit.Pos = pos
	return lit
}

func newStrLiteral(pos ir.Pos, value string) ir.IrLiteral {
	lit := ir.NewStrLiteral(value)
	lit.Pos = pos
	return lit
}

func newArrayType(pos ir.Pos, elemType ir.IrType, length int) ir.IrType {
	typ := ir.NewArrayType(elemType, length)
	typ.Pos = pos
	return typ
}

func newStructType(pos ir.Pos, fields []ir.StructField) ir.IrType {
	typ := ir.NewStructType(fields)
	typ.Pos = pos
	return typ
}

func newTupleType(pos ir.Pos, values []ir.IrType) ir.IrType {
	typ := ir.NewTupleType(values)
	typ.Pos = pos
	return typ
}

func newVariantType(pos ir.Pos, fields []ir.VariantTag) ir.IrType {
	typ := ir.NewVariantType(fields)
	typ.Pos = pos
	return typ
}

func newAssignExpr(arg, ret ast.Expr) ast.Expr {
	return ast.NewAssignExpr(makePos(arg.Pos, ret.Pos), arg, ret)
}

func newIfExpr(pos ir.Pos, condition ast.Expr, then ast.Expr, elseExpr *ast.Expr) ast.Expr {
	if elseExpr == nil {
		return ast.NewAppTermExpr(
			pos,
			ast.NewVarExpr(ast.NewID("ifthen", pos)),
			ast.NewTupleExpr(pos, []ast.Expr{condition, then}))
	}

	return ast.NewAppTermExpr(
		pos,
		ast.NewVarExpr(ast.NewID("ifelse", pos)),
		ast.NewTupleExpr(pos, []ast.Expr{condition, then, *elseExpr}))
}

func newMatchArm(tag, arg ast.ID, body ast.Expr) ast.MatchArm {
	return ast.NewMatchArm(tag.Value, arg.Value, body)
}

func newLetExpr(varName ast.ID, varType *ir.IrType, value ast.Expr) ast.Expr {
	return ast.NewLetExpr(makePos(varName.Pos, value.Pos), varName.Value, varType, value)
}

func newAppTypeExpr(expr ast.Expr, types []ir.IrType) ast.Expr {
	pos := expr.Pos
	for _, typ := range types {
		expr = ast.NewAppTypeExpr(pos, expr, typ)
	}
	return expr
}

func newModuleID(tokens []Token) ir.ModuleID {
	var b strings.Builder
	ir.Interleave(tokens, func() { b.WriteString(ir.ModuleIDSeparator) }, func(_ int, token Token) {
		b.WriteString(token.Text)
	})
	return ir.NewModuleID(b.String(), makePos(tokens[0].Pos, tokens[len(tokens)-1].Pos))
}

func newNumberLiteral(token Token) ir.IrLiteral {
	if strings.Contains(token.Text, ".") {
		integer, decimal, err := parseFloat[int64](token.Text)
		if err != nil {
			panic(fmt.Errorf("expected floating point literal; got %q", token.Text))
		}

		return newFloatLiteral(token.Pos, integer, decimal)
	}

	value, err := parseNumber[int64](token.Text)
	if err != nil {
		// TODO: Avoid panic.
		panic(fmt.Errorf("expected integer; got %q", token.Text))
	}

	return newIntLiteral(token.Pos, value)
}

func newStringLiteral(token Token) ir.IrLiteral {
	text := token.Text

	if !strings.HasPrefix(text, `"`) {
		// TODO: Avoid panic.
		panic(fmt.Errorf(`expected string terminated with '"'; got %q`, token.Text))
	}

	if !strings.HasSuffix(text, `"`) {
		// TODO: Avoid panic.
		panic(fmt.Errorf(`expected string terminated with '"'; got %q`, token.Text))
	}

	text = strings.TrimPrefix(text, `"`)
	text = strings.TrimSuffix(text, `"`)
	return newStrLiteral(token.Pos, text)
}

type action = func(args []any) any

// TODO: Enhance this with a type check, e.g., first[T any] { ... .(T) }.
func first() action {
	return func(args []any) any {
		return args[0]
	}
}

func second() action {
	return func(args []any) any {
		return args[1]
	}
}

func listAppend[T any](arg1, arg2 int) action {
	return func(args []any) any {
		return append(args[arg1].([]T), args[arg2].(T))
	}
}

func list[T any](is ...int) action {
	return func(args []any) any {
		values := make([]T, 0, len(args))
		for _, i := range is {
			values = append(values, args[i].(T))
		}
		return values
	}
}

func listNil[T any]() action {
	return func(args []any) any {
		return []T{}
	}
}

func binOp() action {
	return func(args []any) any {
		operator := args[1].(Token)
		return newBinOpExpr(
			ast.NewVarExpr(ast.NewID(operator.Text, operator.Pos)),
			nil, /* typeArgs */
			args[0].(ast.Expr),
			args[2].(ast.Expr))
	}
}

func binOpTypeApplicative() action {
	return func(args []any) any {
		operand1 := args[0].(ast.Expr)
		operator := args[1].(Token)
		typeArgs := args[2].([]ir.IrType)
		operand2 := args[3].(ast.Expr)
		return newBinOpExpr(
			ast.NewVarExpr(ast.NewID(operator.Text, operator.Pos)),
			typeArgs,
			operand1,
			operand2)
	}
}

func unaryOp() action {
	return func(args []any) any {
		operator := args[0].(Token)
		return newUnaryOpExpr(
			ast.NewVarExpr(ast.NewID(operator.Text, operator.Pos)),
			nil, /* typeArgs */
			args[1].(ast.Expr))
	}
}

func unaryOpTypeApplicative() action {
	return func(args []any) any {
		operator := args[0].(Token)
		typeArgs := args[1].([]ir.IrType)
		operand := args[2].(ast.Expr)
		return newUnaryOpExpr(
			ast.NewVarExpr(ast.NewID(operator.Text, operator.Pos)),
			typeArgs,
			operand)
	}
}

func NewGrammar(initial grammar.ProductionLine) []grammar.ProductionLine {
	return []grammar.ProductionLine{
		initial,

		/* Workspace */

		{"Workspace -> workspace { WorkspacePackages }", func(args []any) any {
			return ast.NewWorkspace(args[2].(ast.Packages))
		}},

		{"WorkspacePackages -> PackagesSection", first()},

		{"PackagesSection -> packages { Packages }", func(args []any) any {
			packages := args[2].([]ast.Package)
			return ast.NewPackages(packages, makePos2(args))
		}},

		{"Packages -> Packages Package", listAppend[ast.Package](0, 1)},
		{"Packages -> Package", list[ast.Package](0)},

		{"Package -> prefix ModuleID in Filename", func(args []any) any {
			moduleID := args[1].(ir.ModuleID)
			filename := args[3].(ir.Filename)
			pos := makePos(args[0].(Token).Pos, filename.Pos)
			return ast.NewPrefixPackage(moduleID, filename, pos)
		}},
		{"Package -> module ModuleID in Filename", func(args []any) any {
			moduleID := args[1].(ir.ModuleID)
			filename := args[3].(ir.Filename)
			pos := makePos(args[0].(Token).Pos, filename.Pos)
			return ast.NewModulePackage(moduleID, filename, pos)
		}},

		{"ModuleID -> ModuleIDTokens", func(args []any) any {
			return newModuleID(args[0].([]Token))
		}},

		{"ModuleIDTokens -> ModuleIDTokens . Token", listAppend[Token](0, 2)},
		{"ModuleIDTokens -> Token", list[Token](0)},

		{"Filename -> StringLiteral", func(args []any) any {
			literal := args[0].(ir.IrLiteral)
			return ir.NewFilename(*literal.Str, literal.Pos)
		}},

		/* Base source file */

		{"SourceFile -> module ModuleID", func(args []any) any {
			id := args[1].(ir.ModuleID)
			return ast.SourceFile{Header: ast.NewBaseSourceFileHeader(id)}
		}},
		{"SourceFile -> module ModuleID SourceFileImports", func(args []any) any {
			id := args[1].(ir.ModuleID)
			sourceFile := args[2].(ast.SourceFile)
			sourceFile.Header = ast.NewBaseSourceFileHeader(id)
			return sourceFile
		}},

		/* Implementation source file */

		{"SourceFile -> implements ModuleID", func(args []any) any {
			id := args[1].(ir.ModuleID)
			return ast.SourceFile{Header: ast.NewImplSourceFileHeader(id)}
		}},
		{"SourceFile -> implements ModuleID SourceFileImports", func(args []any) any {
			id := args[1].(ir.ModuleID)
			sourceFile := args[2].(ast.SourceFile)
			sourceFile.Header = ast.NewImplSourceFileHeader(id)
			return sourceFile
		}},

		/* Source files */

		{"SourceFileImports -> ImportsSection SourceFileImpls", func(args []any) any {
			sourceFile := args[1].(ast.SourceFile)
			sourceFile.Imports = args[0].(ast.Imports)
			return sourceFile
		}},
		{"SourceFileImports -> ImportsSection", func(args []any) any {
			return ast.SourceFile{Imports: args[0].(ast.Imports)}
		}},
		{"SourceFileImports -> SourceFileImpls", first()},

		{"SourceFileImpls -> ImplsSection SourceFileFlags", func(args []any) any {
			sourceFile := args[1].(ast.SourceFile)
			sourceFile.Impls = args[0].(ast.Impls)
			return sourceFile
		}},
		{"SourceFileImpls -> ImplsSection", func(args []any) any {
			return ast.SourceFile{Impls: args[0].(ast.Impls)}
		}},
		{"SourceFileImpls -> SourceFileFlags", first()},

		{"SourceFileFlags -> FlagsSection SourceFileBody", func(args []any) any {
			sourceFile := args[1].(ast.SourceFile)
			sourceFile.Flags = args[0].(ast.Flags)
			return sourceFile
		}},
		{"SourceFileFlags -> FlagsSection", func(args []any) any {
			return ast.SourceFile{Flags: args[0].(ast.Flags)}
		}},
		{"SourceFileFlags -> SourceFileBody", first()},

		{"SourceFileBody -> Sources", func(args []any) any {
			return ast.SourceFile{Body: args[0].([]ast.Source)}
		}},

		/* Imports section */

		{"ImportsSection -> imports { ModuleIDs }", func(args []any) any {
			return ast.NewImports(args[2].([]ir.ModuleID), makePos2(args))
		}},

		{"ModuleIDs -> ModuleIDs ModuleID", listAppend[ir.ModuleID](0, 1)},
		{"ModuleIDs -> ModuleID", list[ir.ModuleID](0)},

		/* Impls section */

		{"ImplsSection -> impls { Filenames }", func(args []any) any {
			return ast.NewImpls(args[2].([]ir.Filename), makePos2(args))
		}},

		{"Filenames -> Filenames Filename", listAppend[ir.Filename](0, 1)},
		{"Filenames -> Filename", list[ir.Filename](0)},

		/* Flags section */

		{"FlagsSection -> flags { Filenames }", func(args []any) any {
			ids := args[2].([]ir.Filename)
			for i, id := range ids {
				id.Value = strings.TrimPrefix(id.Value, `"`)
				id.Value = strings.TrimSuffix(id.Value, `"`)
				ids[i] = id
			}
			return ast.NewFlags(ids, makePos2(args))
		}},

		/* Source */

		{"Sources -> Sources Source", listAppend[ast.Source](0, 1)},
		{"Sources -> Source", list[ast.Source](0)},

		{"Source -> DeclSource", first()},
		{"Source -> Function", first()},
		{"Source -> pub Function", func(args []any) any {
			source := args[1].(ast.Source)
			source.Function.Export = true
			return source
		}},

		/* Decl source */

		{"DeclSource -> DeclNoTerm", func(args []any) any {
			return ast.NewDeclSource(args[0].(ir.IrDecl))
		}},
		{"DeclSource -> decl TermDecl", func(args []any) any {
			return ast.NewDeclSource(args[1].(ir.IrDecl))
		}},
		{"DeclSource -> pub TermDecl", func(args []any) any {
			decl := args[1].(ir.IrDecl)
			decl.Export = true
			return ast.NewDeclSource(decl)
		}},

		/* Function */

		{"Function -> fn ID TypeAbstraction FunctionArgs -> AppType Block", func(args []any) any {
			id := args[1].(ast.ID)
			tvars := args[2].([]ir.VarKind)
			funArgs := args[3].([]ir.FunctionArg)
			retType := args[5].(ir.IrType)
			body := args[6].(ast.Expr)
			return ast.NewFunctionSource(
				ast.NewFunction(
					makePos(id.Pos, body.Pos),
					false /* export */, id.Value, tvars, funArgs, retType, body))
		}},
		{"Function -> fn ID FunctionArgs -> PrimaryType Block", func(args []any) any {
			id := args[1].(ast.ID)
			funArgs := args[2].([]ir.FunctionArg)
			retType := args[4].(ir.IrType)
			body := args[5].(ast.Expr)
			return ast.NewFunctionSource(
				ast.NewFunction(
					makePos(id.Pos, body.Pos),
					false /* export */, id.Value, nil /* tvars */, funArgs, retType, body))
		}},

		{"FunctionArgs -> ( Args )", second()},
		{"FunctionArgs -> ( )", listNil[ir.FunctionArg]()},

		{"Args -> Args , Arg", listAppend[ir.FunctionArg](0, 2)},
		{"Args -> Arg", list[ir.FunctionArg](0)},

		{"Arg -> ID : UnquantifiedType", func(args []any) any {
			return ir.FunctionArg{args[0].(ast.ID).Value, args[2].(ir.IrType)}
		}},

		/* Decl */

		// Decl is used by annotations and DeclNoTerm is used by DeclSource. It
		// would be great if both annotations and DeclSource could reuse Decl, but
		// the grammar becomes ambiguous if TermDecl is not preceded by `decl` (or
		// another equivalent solution).
		{"Decl -> pub UnexportedDecl", func(args []any) any {
			decl := args[1].(ir.IrDecl)
			decl.Export = true
			return decl
		}},
		{"Decl -> UnexportedDecl", first()},

		{"UnexportedDecl -> StructDecl", first()},
		{"UnexportedDecl -> TermDecl", first()},
		{"UnexportedDecl -> TupleDecl", first()},
		{"UnexportedDecl -> TypeDecl", first()},
		{"UnexportedDecl -> VariantDecl", first()},

		{"DeclNoTerm -> pub UnexportedDeclNoTerm", func(args []any) any {
			decl := args[1].(ir.IrDecl)
			decl.Export = true
			return decl
		}},
		{"DeclNoTerm -> UnexportedDeclNoTerm", first()},

		{"UnexportedDeclNoTerm -> StructDecl", first()},
		{"UnexportedDeclNoTerm -> TupleDecl", first()},
		{"UnexportedDeclNoTerm -> TypeDecl", first()},
		{"UnexportedDeclNoTerm -> VariantDecl", first()},

		/* Struct decl */

		{"StructDecl -> type ID TypeAbstraction = StructType", func(args []any) any {
			id := args[1].(ast.ID)
			tvars := args[2].([]ir.VarKind)
			structType := args[4].(ir.IrType)

			kind := ir.NewTypeKind()
			for i := 0; i < len(tvars); i++ {
				kind = ir.NewArrowKind(ir.NewTypeKind(), kind)
			}

			return newAliasDecl(id, kind, ir.LambdaVars(tvars, structType), false /* export */)
		}},
		{"StructDecl -> type ID = StructType", func(args []any) any {
			id := args[1].(ast.ID)
			var tvars []ir.VarKind
			structType := args[3].(ir.IrType)
			return newAliasDecl(id, ir.NewTypeKind(), ir.LambdaVars(tvars, structType), false /* export */)
		}},

		/* Term decl */

		{"TermDecl -> ID : QuantifiedType", func(args []any) any {
			return newTermDecl(args[0].(ast.ID), args[2].(ir.IrType), false /* export */)
		}},

		{"QuantifiedType -> UnquantifiedType", func(args []any) any {
			return newQuantifiedType(args[0].(ir.IrType))
		}},

		/* Tuple decl */

		{"TupleDecl -> type ID TypeAbstraction = TupleType", func(args []any) any {
			id := args[1].(ast.ID)
			tvars := args[2].([]ir.VarKind)
			tupleType := args[4].(ir.IrType)

			kind := ir.NewTypeKind()
			for i := 0; i < len(tvars); i++ {
				kind = ir.NewArrowKind(ir.NewTypeKind(), kind)
			}

			return newAliasDecl(id, kind, ir.LambdaVars(tvars, tupleType), false /* export */)
		}},
		{"TupleDecl -> type ID = TupleType", func(args []any) any {
			id := args[1].(ast.ID)
			var tvars []ir.VarKind
			tupleType := args[3].(ir.IrType)
			return newAliasDecl(id, ir.NewTypeKind(), ir.LambdaVars(tvars, tupleType), false /* export */)
		}},

		/* Type decl */

		{"TypeDecl -> type ID TypeAbstraction", func(args []any) any {
			tvars := args[2].([]ir.VarKind)

			kind := ir.NewTypeKind()
			for i := 0; i < len(tvars); i++ {
				kind = ir.NewArrowKind(ir.NewTypeKind(), kind)
			}

			return newNameDecl(args[1].(ast.ID), kind, false /* export */)
		}},
		{"TypeDecl -> type ID", func(args []any) any {
			return newNameDecl(args[1].(ast.ID), ir.NewTypeKind(), false /* export */)
		}},

		/* Variant decl */

		{"VariantDecl -> type ID TypeAbstraction = VariantType", func(args []any) any {
			id := args[1].(ast.ID)
			tvars := args[2].([]ir.VarKind)
			variantType := args[4].(ir.IrType)

			kind := ir.NewTypeKind()
			for i := 0; i < len(tvars); i++ {
				kind = ir.NewArrowKind(ir.NewTypeKind(), kind)
			}

			return newAliasDecl(id, kind, ir.LambdaVars(tvars, variantType), false /* export */)
		}},
		{"VariantDecl -> type ID = VariantType", func(args []any) any {
			id := args[1].(ast.ID)
			var tvars []ir.VarKind
			variantType := args[3].(ir.IrType)
			return newAliasDecl(id, ir.NewTypeKind(), ir.LambdaVars(tvars, variantType), false /* export */)
		}},

		/* Type variables */

		{"TypeAbstraction -> [ Tvars ]", second()},

		{"Tvars -> Tvars , Tvar", listAppend[ir.VarKind](0, 2)},
		{"Tvars -> Tvar", list[ir.VarKind](0)},

		{"Tvar -> ' ID", func(args []any) any {
			return ir.VarKind{args[1].(ast.ID).Value, ir.NewTypeKind()}
		}},

		/* Unquantified type */

		{"UnquantifiedType -> ForallType", first()},

		/* Forall type */

		{"ForallType -> forall TypeAbstraction FunctionType", func(args []any) any {
			tvars := args[1].([]ir.VarKind)
			subType := args[2].(ir.IrType)
			return newForallType(makePos(args[0].(Token).Pos, subType.Pos), tvars, subType)
		}},
		{"ForallType -> FunctionType", first()},

		/* Function type */

		{"FunctionType -> AppType -> FunctionType", func(args []any) any {
			return newFunctionType(args[0].(ir.IrType), args[2].(ir.IrType))
		}},
		{"FunctionType -> AppType", first()},

		/* App type */

		{"AppType -> AppType PrimaryType", func(args []any) any {
			return newAppType(args[0].(ir.IrType), args[1].(ir.IrType))
		}},
		{"AppType -> PrimaryType", first()},

		/* Simple Type */

		{"PrimaryType -> ArrayType", first()},
		{"PrimaryType -> StructType", first()},
		{"PrimaryType -> TupleType", first()},
		{"PrimaryType -> VariantType", first()},
		{"PrimaryType -> ' ID", func(args []any) any {
			return newVarType(args[1].(ast.ID))
		}},
		{"PrimaryType -> ID", func(args []any) any {
			return newNameType(args[0].(ast.ID))
		}},
		{"PrimaryType -> ( UnquantifiedType )", second()},

		/* Array type */

		{"ArrayType -> [ UnquantifiedType , IntLiteral ]", func(args []any) any {
			elemType := args[1].(ir.IrType)
			length := *args[3].(ir.IrLiteral).Int
			return newArrayType(makePos2(args), elemType, int(length))
		}},
		{"ArrayType -> [ UnquantifiedType ]", func(args []any) any {
			elemType := args[1].(ir.IrType)
			length := math.MaxInt
			return newArrayType(makePos2(args), elemType, length)
		}},

		/* Struct type */

		{"StructType -> struct { }", func(args []any) any {
			return newStructType(makePos2(args), nil)
		}},
		{"StructType -> struct { Fields }", func(args []any) any {
			return newStructType(makePos2(args), args[2].([]ir.StructField))
		}},
		{"StructType -> struct { Fields , }", func(args []any) any {
			return newStructType(makePos2(args), args[2].([]ir.StructField))
		}},
		{"StructType -> struct { Fields , ; }", func(args []any) any {
			return newStructType(makePos2(args), args[2].([]ir.StructField))
		}},

		{"Fields -> Fields , Field", listAppend[ir.StructField](0, 2)},
		{"Fields -> Field", list[ir.StructField](0)},

		{"Field -> ID UnquantifiedType", func(args []any) any {
			return ir.StructField{args[0].(ast.ID).Value, args[1].(ir.IrType)}
		}},
		{"Field -> ; ID UnquantifiedType", func(args []any) any {
			return ir.StructField{args[1].(ast.ID).Value, args[2].(ir.IrType)}
		}},

		/* Tuple type */

		{"TupleType -> ( )", func(args []any) any {
			return newTupleType(makePos2(args), nil)
		}},
		{"TupleType -> ( TupleTypeArgs )", func(args []any) any {
			return newTupleType(makePos2(args), args[1].([]ir.IrType))
		}},

		{"TupleTypeArgs -> TupleTypeArgs , UnquantifiedType", listAppend[ir.IrType](0, 2)},
		{"TupleTypeArgs -> UnquantifiedType , UnquantifiedType", list[ir.IrType](0, 2)},

		/* Variant type */

		{"VariantType -> variant { }", func(args []any) any {
			return newVariantType(makePos2(args), nil)
		}},
		{"VariantType -> variant { Tags }", func(args []any) any {
			return newVariantType(makePos2(args), args[2].([]ir.VariantTag))
		}},
		{"VariantType -> variant { Tags , }", func(args []any) any {
			return newVariantType(makePos2(args), args[2].([]ir.VariantTag))
		}},
		{"VariantType -> variant { Tags , ; }", func(args []any) any {
			return newVariantType(makePos2(args), args[2].([]ir.VariantTag))
		}},

		{"Tags -> Tags , Tag", listAppend[ir.VariantTag](0, 2)},
		{"Tags -> Tag", list[ir.VariantTag](0)},

		{"Tag -> ID UnquantifiedType", func(args []any) any {
			return ir.VariantTag{args[0].(ast.ID).Value, args[1].(ir.IrType)}
		}},
		{"Tag -> ; ID UnquantifiedType", func(args []any) any {
			return ir.VariantTag{args[1].(ast.ID).Value, args[2].(ir.IrType)}
		}},

		/* ID */

		{"ID -> IDTokens", func(args []any) any {
			tokens := args[0].([]Token)

			for _, token := range tokens {
				if len(token.Text) == 0 || unicode.IsDigit(rune(token.Text[0])) {
					// TODO: Avoid panic.
					panic(fmt.Errorf("expected identifier; got %q; identifiers must begin with a non-digit character", token.Text))
				}
			}

			var b strings.Builder
			ir.Interleave(tokens, func() { b.WriteString(ir.NamespaceSeparator) }, func(_ int, token Token) {
				b.WriteString(token.Text)
			})

			return ast.NewID(b.String(), makePos(tokens[0].Pos, tokens[len(tokens)-1].Pos))
		}},

		// TODO: To support array::set. Fix this.
		{"IDTokens -> IDTokens :: set", listAppend[Token](0, 2)},
		{"IDTokens -> IDTokens :: Token", listAppend[Token](0, 2)},
		{"IDTokens -> Token", list[Token](0)},

		// The parenthesis around operator IDs are not needed to make the
		// grammar unambiguous. They help with readability.
		//
		// For example, it would be hard to read:
		//
		//   + == +
		//
		// Better is:
		//
		//   (+) == (+)
		{"ID -> ( || )", func(args []any) any {
			token := args[1].(Token)
			return ast.NewID(token.Text, token.Pos)
		}},
		{"ID -> ( && )", func(args []any) any {
			token := args[1].(Token)
			return ast.NewID(token.Text, token.Pos)
		}},
		{"ID -> ( != )", func(args []any) any {
			token := args[1].(Token)
			return ast.NewID(token.Text, token.Pos)
		}},
		{"ID -> ( == )", func(args []any) any {
			token := args[1].(Token)
			return ast.NewID(token.Text, token.Pos)
		}},
		{"ID -> ( > )", func(args []any) any {
			token := args[1].(Token)
			return ast.NewID(token.Text, token.Pos)
		}},
		{"ID -> ( >= )", func(args []any) any {
			token := args[1].(Token)
			return ast.NewID(token.Text, token.Pos)
		}},
		{"ID -> ( < )", func(args []any) any {
			token := args[1].(Token)
			return ast.NewID(token.Text, token.Pos)
		}},
		{"ID -> ( <= )", func(args []any) any {
			token := args[1].(Token)
			return ast.NewID(token.Text, token.Pos)
		}},
		{"ID -> ( + )", func(args []any) any {
			token := args[1].(Token)
			return ast.NewID(token.Text, token.Pos)
		}},
		{"ID -> ( - )", func(args []any) any {
			token := args[1].(Token)
			return ast.NewID(token.Text, token.Pos)
		}},
		{"ID -> ( * )", func(args []any) any {
			token := args[1].(Token)
			return ast.NewID(token.Text, token.Pos)
		}},
		{"ID -> ( / )", func(args []any) any {
			token := args[1].(Token)
			return ast.NewID(token.Text, token.Pos)
		}},
		{"ID -> ( ! )", func(args []any) any {
			token := args[1].(Token)
			return ast.NewID(token.Text, token.Pos)
		}},

		/* Block */

		{"Block -> { Terms }", func(args []any) any {
			return ast.NewBlockExpr(makePos2(args), args[1].([]ast.Expr))
		}},
		{"Block -> { Terms ; }", func(args []any) any {
			return ast.NewBlockExpr(makePos2(args), args[1].([]ast.Expr))
		}},
		{"Block -> { Term }", func(args []any) any {
			expr := args[1].(ast.Expr)
			return ast.NewBlockExpr(makePos2(args), []ast.Expr{expr})
		}},
		{"Block -> { Term ; }", func(args []any) any {
			expr := args[1].(ast.Expr)
			return ast.NewBlockExpr(makePos2(args), []ast.Expr{expr})
		}},

		{"Terms -> Terms ; Term", listAppend[ast.Expr](0, 2)},
		{"Terms -> ; Term", list[ast.Expr](1)},

		/* Term */

		{"Term -> AssignTerm", first()},
		{"Term -> LetTerm", first()},
		{"Term -> ReturnTerm", first()},
		{"Term -> Expression", first()},

		/* Assign term */

		{"AssignTerm -> ID <- Expression", func(args []any) any {
			ret := args[0].(ast.ID)
			arg := args[2].(ast.Expr)
			return newAssignExpr(arg, ast.NewVarExpr(ret))
		}},
		{"AssignTerm -> TupleTerm <- Expression", func(args []any) any {
			ret := args[0].(ast.Expr)
			arg := args[2].(ast.Expr)
			return newAssignExpr(arg, ret)
		}},

		/* Let term */

		{"LetTerm -> let ID : UnquantifiedType = Expression", func(args []any) any {
			varName := args[1].(ast.ID)
			varType := args[3].(ir.IrType)
			value := args[5].(ast.Expr)
			return newLetExpr(varName, &varType, value)
		}},
		{"LetTerm -> let ID = Expression", func(args []any) any {
			varName := args[1].(ast.ID)
			value := args[3].(ast.Expr)
			return newLetExpr(varName, nil /* varType */, value)
		}},

		/* Return term */

		{"ReturnTerm -> return Expression", func(args []any) any {
			expr := args[1].(ast.Expr)
			return ast.NewReturnExpr(makePos(args[0].(Token).Pos, expr.Pos), expr)
		}},

		/* Expression */

		{"Expression -> Block", first()},
		{"Expression -> IfTerm", first()},
		{"Expression -> LambdaTerm", first()},
		{"Expression -> MatchTerm", first()},
		{"Expression -> SetTerm", first()},

		/* If term */

		{"IfTerm -> if Expression Block", func(args []any) any {
			condition := args[1].(ast.Expr)
			then := args[2].(ast.Expr)
			var elseExpr *ast.Expr
			return newIfExpr(
				makePos(args[0].(Token).Pos, then.Pos), condition, then, elseExpr)
		}},
		{"IfTerm -> if Expression Block else Block", func(args []any) any {
			condition := args[1].(ast.Expr)
			then := args[2].(ast.Expr)
			elseExpr := args[4].(ast.Expr)
			return newIfExpr(
				makePos(args[0].(Token).Pos, elseExpr.Pos), condition, then, &elseExpr)
		}},
		{"IfTerm -> if Expression Block else IfTerm", func(args []any) any {
			condition := args[1].(ast.Expr)
			then := args[2].(ast.Expr)
			elseExpr := args[4].(ast.Expr)
			return newIfExpr(
				makePos(args[0].(Token).Pos, elseExpr.Pos), condition, then, &elseExpr)
		}},

		/* Lambda term */

		{"LambdaTerm -> fn TypeAbstraction FunctionArgs Block", func(args []any) any {
			tvars := args[1].([]ir.VarKind)
			funArgs := args[2].([]ir.FunctionArg)
			body := args[3].(ast.Expr)
			return ast.Lambda(
				makePos(args[0].(Token).Pos, body.Pos),
				tvars, funArgs, body)
		}},
		{"LambdaTerm -> fn FunctionArgs Block", func(args []any) any {
			funArgs := args[1].([]ir.FunctionArg)
			body := args[2].(ast.Expr)
			return ast.Lambda(
				makePos(args[0].(Token).Pos, body.Pos),
				nil /* tvars */, funArgs, body)
		}},
		{"LambdaTerm -> Operator", first()},

		/* Match term */

		{"MatchTerm -> match Expression { MatchArms }", func(args []any) any {
			expr := args[1].(ast.Expr)
			arms := args[3].([]ast.MatchArm)
			return ast.NewMatchExpr(makePos2(args), expr, arms)
		}},
		{"MatchTerm -> match Expression { MatchArms , }", func(args []any) any {
			expr := args[1].(ast.Expr)
			arms := args[3].([]ast.MatchArm)
			return ast.NewMatchExpr(makePos2(args), expr, arms)
		}},
		{"MatchTerm -> match Expression { MatchArms , ; }", func(args []any) any {
			expr := args[1].(ast.Expr)
			arms := args[3].([]ast.MatchArm)
			return ast.NewMatchExpr(makePos2(args), expr, arms)
		}},

		{"MatchArms -> MatchArms , MatchArm", listAppend[ast.MatchArm](0, 2)},
		{"MatchArms -> MatchArm", list[ast.MatchArm](0)},

		{"MatchArm -> ID ID => Term", func(args []any) any {
			tag := args[0].(ast.ID)
			arg := args[1].(ast.ID)
			body := args[3].(ast.Expr)
			return newMatchArm(tag, arg, body)
		}},
		{"MatchArm -> ; ID ID => Term", func(args []any) any {
			tag := args[1].(ast.ID)
			arg := args[2].(ast.ID)
			body := args[4].(ast.Expr)
			return newMatchArm(tag, arg, body)
		}},

		/* Set term */

		// TODO: Get rid of 'set' keyword. This is only here to avoid grammar conflicts.
		{"SetTerm -> set Expression { LabelValues }", func(args []any) any {
			expr := args[1].(ast.Expr)
			values := args[3].([]ast.LabelValue)
			return ast.NewSetExpr(makePos2(args), expr, values)
		}},
		{"SetTerm -> set Expression { LabelValues , }", func(args []any) any {
			expr := args[1].(ast.Expr)
			values := args[3].([]ast.LabelValue)
			return ast.NewSetExpr(makePos2(args), expr, values)
		}},
		{"SetTerm -> set Expression { LabelValues , ; }", func(args []any) any {
			expr := args[1].(ast.Expr)
			values := args[3].([]ast.LabelValue)
			return ast.NewSetExpr(makePos2(args), expr, values)
		}},

		/* Operators */

		{"Operator -> LogicalOr", first()},

		{"LogicalOr -> LogicalOr || Equality", binOp()},
		{"LogicalOr -> LogicalAnd", first()},

		{"LogicalAnd -> LogicalAnd && Equality", binOp()},
		{"LogicalAnd -> Equality", first()},

		{"Equality -> Equality != TypeApplicativeArgs Comparison", binOpTypeApplicative()},
		{"Equality -> Equality != Comparison", binOp()},
		{"Equality -> Equality == TypeApplicativeArgs Comparison", binOpTypeApplicative()},
		{"Equality -> Equality == Comparison", binOp()},
		{"Equality -> Comparison", first()},

		{"Comparison -> Comparison > TypeApplicativeArgs Additive", binOpTypeApplicative()},
		{"Comparison -> Comparison > Additive", binOp()},
		{"Comparison -> Comparison >= TypeApplicativeArgs Additive", binOpTypeApplicative()},
		{"Comparison -> Comparison >= Additive", binOp()},
		{"Comparison -> Comparison < TypeApplicativeArgs Additive", binOpTypeApplicative()},
		{"Comparison -> Comparison < Additive", binOp()},
		{"Comparison -> Comparison <= TypeApplicativeArgs Additive", binOpTypeApplicative()},
		{"Comparison -> Comparison <= Additive", binOp()},
		{"Comparison -> Additive", first()},

		{"Additive -> Additive + TypeApplicativeArgs Multiplicative", binOpTypeApplicative()},
		{"Additive -> Additive + Multiplicative", binOp()},
		{"Additive -> Additive - TypeApplicativeArgs Multiplicative", binOpTypeApplicative()},
		{"Additive -> Additive - Multiplicative", binOp()},
		{"Additive -> Multiplicative", first()},

		{"Multiplicative -> Multiplicative * TypeApplicativeArgs Unary", binOpTypeApplicative()},
		{"Multiplicative -> Multiplicative * Unary", binOp()},
		{"Multiplicative -> Multiplicative / TypeApplicativeArgs Unary", binOpTypeApplicative()},
		{"Multiplicative -> Multiplicative / Unary", binOp()},
		{"Multiplicative -> Unary", first()},

		{"Unary -> ! TypeApplicativeArgs Unary", unaryOpTypeApplicative()},
		{"Unary -> ! Unary", unaryOp()},
		{"Unary -> - TypeApplicativeArgs Unary", unaryOpTypeApplicative()},
		{"Unary -> - Unary", unaryOp()},
		{"Unary -> Applicative", first()},

		/* Applicative */

		{"Applicative -> Applicative Primary", func(args []any) any {
			fun := args[0].(ast.Expr)
			arg := args[1].(ast.Expr)
			return ast.NewAppTermExpr(makePos(fun.Pos, arg.Pos), fun, arg)
		}},
		{"Applicative -> TypeApplicative", first()},

		/* Type applicative */

		{"TypeApplicative -> Primary TypeApplicativeArgs", func(args []any) any {
			return newAppTypeExpr(args[0].(ast.Expr), args[1].([]ir.IrType))
		}},
		{"TypeApplicative -> Primary", first()},

		{"TypeApplicativeArgs -> [ TupleTypeArgs ]", second()},
		{"TypeApplicativeArgs -> [ UnquantifiedType ]", list[ir.IrType](1)},

		/* Primary */

		{"Primary -> ProjectionTerm", first()},
		{"Primary -> IntLiteral", func(args []any) any {
			return ast.NewConstExpr(args[0].(ir.IrLiteral))
		}},
		{"Primary -> FloatLiteral", func(args []any) any {
			return ast.NewConstExpr(args[0].(ir.IrLiteral))
		}},

		/* Projection term */

		{"ProjectionTerm -> ProjectionTerm . IntLiteral", func(args []any) any {
			expr := args[0].(ast.Expr)
			label := args[2].(ir.IrLiteral)
			return ast.NewProjectionExpr(makePos(expr.Pos, label.Pos), expr, label.String())
		}},
		{"ProjectionTerm -> ProjectionTerm . Token", func(args []any) any {
			expr := args[0].(ast.Expr)
			label := args[2].(Token)
			return ast.NewProjectionExpr(makePos(expr.Pos, label.Pos), expr, label.Text)
		}},
		{"ProjectionTerm -> Deref", first()},

		/* Deref */

		{"Deref -> InjectionTerm", first()},
		{"Deref -> StringLiteral", func(args []any) any {
			return ast.NewConstExpr(args[0].(ir.IrLiteral))
		}},
		{"Deref -> StructTerm", first()},
		{"Deref -> TupleTerm", first()},
		{"Deref -> VarTerm", first()},
		{"Deref -> ( Expression )", second()},

		/* Injection term */

		{"InjectionTerm -> variant { PrimaryType LabelValue }", func(args []any) any {
			variantType := args[2].(ir.IrType)
			labelValue := args[3].(ast.LabelValue)
			return ast.NewInjectionExpr(makePos2(args), variantType, labelValue.Label, labelValue.Value)
		}},
		{"InjectionTerm -> variant { ; PrimaryType LabelValue ; }", func(args []any) any {
			variantType := args[3].(ir.IrType)
			labelValue := args[4].(ast.LabelValue)
			return ast.NewInjectionExpr(makePos2(args), variantType, labelValue.Label, labelValue.Value)
		}},

		/* Struct term */

		{"StructTerm -> struct { }", func(args []any) any { return ast.NewStructExpr(makePos2(args), nil) }},

		{"StructTerm -> struct { LabelValues }", func(args []any) any {
			return ast.NewStructExpr(makePos2(args), args[2].([]ast.LabelValue))
		}},
		{"StructTerm -> struct { LabelValues , }", func(args []any) any {
			return ast.NewStructExpr(makePos2(args), args[2].([]ast.LabelValue))
		}},
		{"StructTerm -> struct { LabelValues , ; }", func(args []any) any {
			return ast.NewStructExpr(makePos2(args), args[2].([]ast.LabelValue))
		}},

		{"LabelValues -> LabelValues , LabelValue", listAppend[ast.LabelValue](0, 2)},
		{"LabelValues -> LabelValue", list[ast.LabelValue](0)},

		{"LabelValue -> ID = Expression", func(args []any) any {
			label := args[0].(ast.ID)
			value := args[2].(ast.Expr)
			return ast.LabelValue{label.Value, value}
		}},
		{"LabelValue -> ; ID = Expression", func(args []any) any {
			label := args[1].(ast.ID)
			value := args[3].(ast.Expr)
			return ast.LabelValue{label.Value, value}
		}},
		{"LabelValue -> IntLiteral = Expression", func(args []any) any {
			label := *args[0].(ir.IrLiteral).Int
			value := args[2].(ast.Expr)
			return ast.LabelValue{fmt.Sprintf("%d", label), value}
		}},
		{"LabelValue -> ; IntLiteral = Expression", func(args []any) any {
			label := *args[1].(ir.IrLiteral).Int
			value := args[3].(ast.Expr)
			return ast.LabelValue{fmt.Sprintf("%d", label), value}
		}},

		/* Tuple term */

		{"TupleTerm -> ( )", func(args []any) any {
			return ast.NewTupleExpr(makePos2(args), nil)
		}},
		{"TupleTerm -> ( TupleTermArgs )", func(args []any) any {
			return ast.NewTupleExpr(makePos2(args), args[1].([]ast.Expr))
		}},

		{"TupleTermArgs -> TupleTermArgs , Expression", listAppend[ast.Expr](0, 2)},
		{"TupleTermArgs -> Expression , Expression", list[ast.Expr](0, 2)},

		/* Var term */

		{"VarTerm -> ID", func(args []any) any {
			return ast.NewVarExpr(args[0].(ast.ID))
		}},

		/* Literals */

		{"IntLiteral -> NumberToken", func(args []any) any {
			return newNumberLiteral(args[0].(Token))
		}},

		{"FloatLiteral -> IntLiteral . IntLiteral", func(args []any) any {
			integer := args[0].(ir.IrLiteral)
			decimal := args[2].(ir.IrLiteral)

			if !integer.Is(ir.IntLiteral) {
				// TODO: Avoid panic.
				panic(fmt.Errorf("expected integer for the integer part of the floating point number; got %v", integer))
			}

			if !decimal.Is(ir.IntLiteral) {
				// TODO: Avoid panic.
				panic(fmt.Errorf("expected integer for the decimal part of the floating point number; got %v", decimal))
			}

			return newFloatLiteral(makePos(integer.Pos, decimal.Pos), *integer.Int, *decimal.Int)
		}},

		{"StringLiteral -> StringToken", func(args []any) any {
			return newStringLiteral(args[0].(Token))
		}},
	}
}
