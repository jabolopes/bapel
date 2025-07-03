package bplparser2

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

func makePos(pos1, pos2 ir.Pos) ir.Pos {
	return ir.NewRangePos(pos1.Filename, pos1.BeginLineNum, pos2.EndLineNum)
}

func makePos2(args []any) ir.Pos {
	return makePos(args[0].(Token).Pos, args[len(args)-1].(Token).Pos)
}

func newUnaryOpTerm(id ir.IrTerm, term ir.IrTerm) (r ir.IrTerm) {
	defer func() {
		r.Pos = makePos(id.Pos, term.Pos)
	}()

	if id.Is(ir.VarTerm) && id.Var.ID == "-" {
		// 0 - $term
		return ir.CallPF(id, nil /* types */, ir.Number(0), term)
	}

	return ir.Call(id, term)
}

func newBinOpTerm(id ir.IrTerm, t1, t2 ir.IrTerm) ir.IrTerm {
	term := ir.Call(id, t1, t2)
	term.Pos = makePos(t1.Pos, t2.Pos)
	return term
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

func newDeclSource(decl ir.IrDecl) ast.Source {
	source := ast.NewDeclSource(decl)
	source.Pos = decl.Pos
	return source
}

func newFunctionSource(pos ir.Pos, fun ir.IrFunction) ast.Source {
	fun.Pos = pos
	source := ast.NewFunctionSource(fun)
	source.Pos = pos
	return source
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

func newInjectionTerm(pos ir.Pos, variantType ir.IrType, tag, value ir.IrTerm) ir.IrTerm {
	typ := ir.NewInjectionTerm(variantType, tag, value)
	typ.Pos = pos
	return typ
}

func newSetTerm(pos ir.Pos, term ir.IrTerm, values []ir.LabelValue) ir.IrTerm {
	typ := ir.NewSetTerm(term, values)
	typ.Pos = pos
	return typ
}

func newStructTerm(pos ir.Pos, values []ir.LabelValue) ir.IrTerm {
	typ := ir.NewStructTerm(values)
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

func newBlockTerm(pos ir.Pos, terms []ir.IrTerm) ir.IrTerm {
	typ := ir.NewBlockTerm(terms)
	typ.Pos = pos
	return typ
}

func newIDTerm(id ast.ID) ir.IrTerm {
	term := ir.ID(id.Value)
	term.Pos = id.Pos
	return term
}

func newAssignTerm(arg, ret ir.IrTerm) ir.IrTerm {
	term := ir.NewAssignTerm(arg, ret)
	term.Pos = makePos(arg.Pos, ret.Pos)
	return term
}

func newIfTerm(pos ir.Pos, condition ir.IrTerm, then ir.IrTerm, elseTerm *ir.IrTerm) ir.IrTerm {
	term := ir.NewIfTerm(condition, then, elseTerm)
	term.Pos = pos
	return term
}

func newLambdaTerm(pos ir.Pos, tvars []ir.VarKind, args []ir.ArgType, body ir.IrTerm) ir.IrTerm {
	term := ir.Lambda(tvars, args, body)
	term.Pos = pos
	return term
}

func newMatchArm(tag, arg ast.ID, body ir.IrTerm) ir.MatchArm {
	return ir.NewMatchArm(tag.Value, arg.Value, body)
}

func newMatchTerm(pos ir.Pos, term ir.IrTerm, arms []ir.MatchArm) ir.IrTerm {
	match := ir.NewMatchTerm(term, arms)
	match.Pos = pos
	return match
}

func newLetTerm(varName ast.ID, varType ir.IrType, value ir.IrTerm) ir.IrTerm {
	term := ir.NewLetTerm(varName.Value, varType, value)
	term.Pos = makePos(varName.Pos, value.Pos)
	return term
}

func newProjectionTerm(pos ir.Pos, term ir.IrTerm, label string) ir.IrTerm {
	proj := ir.NewProjectionTerm(term, label)
	proj.Pos = pos
	return proj
}

func newReturnTerm(expr ir.IrTerm) ir.IrTerm {
	term := ir.NewReturnTerm(expr)
	term.Pos = expr.Pos
	return term
}

func newAppTermTerm(fun, arg ir.IrTerm) ir.IrTerm {
	term := ir.NewAppTermTerm(fun, arg)
	term.Pos = makePos(fun.Pos, arg.Pos)
	return term
}

func newAppTypeTerm(term ir.IrTerm, types []ir.IrType) ir.IrTerm {
	pos := term.Pos
	for _, typ := range types {
		term = ir.NewAppTypeTerm(term, typ)
	}
	term.Pos = pos
	return term
}

func newTupleTerm(pos ir.Pos, elems []ir.IrTerm) ir.IrTerm {
	term := ir.NewTupleTerm(elems)
	term.Pos = pos
	return term
}

func newVarTerm(pos ir.Pos, id string) ir.IrTerm {
	term := ir.NewVarTerm(id)
	term.Pos = pos
	return term
}

func newImportID(token Token) ast.ModuleID {
	text := token.Text

	if strings.HasPrefix(text, `"`) {
		if !strings.HasSuffix(text, `"`) {
			// TODO: Avoid panic.
			panic(fmt.Errorf(`expected string terminated with '"'; got %q`, token.Text))
		}

		text = strings.TrimPrefix(text, `"`)
		text = strings.TrimSuffix(text, `"`)
	}

	return ast.NewModuleID(text, token.Pos)
}

func newImplID(token Token) ast.ID {
	if text := token.Text; strings.HasPrefix(text, `"`) {
		if !strings.HasSuffix(text, `"`) {
			// TODO: Avoid panic.
			panic(fmt.Errorf(`expected string terminated with '"'; got %q`, token.Text))
		}

		text = strings.TrimPrefix(text, `"`)
		text = strings.TrimSuffix(text, `"`)
		return ast.ID{token.Pos, text}
	}

	return ast.ID{token.Pos, token.Text}
}

func newLiteralTerm(token Token) ir.IrTerm {
	if unicode.IsDigit(rune(token.Text[0])) {
		value, err := parseNumber[int64](token.Text)
		if err != nil {
			// TODO: Avoid panic.
			panic(fmt.Errorf("expected integer; got %q", token.Text))
		}

		return ir.NewConstTerm(newIntLiteral(token.Pos, value))
	}

	if text := token.Text; strings.HasPrefix(text, `"`) {
		if !strings.HasSuffix(text, `"`) {
			// TODO: Avoid panic.
			panic(fmt.Errorf(`expected string terminated with '"'; got %q`, token.Text))
		}

		text = strings.TrimPrefix(text, `"`)
		text = strings.TrimSuffix(text, `"`)
		return ir.NewConstTerm(newStrLiteral(token.Pos, text))
	}

	return newVarTerm(token.Pos, token.Text)
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
		return newBinOpTerm(
			newIDTerm(ast.ID{operator.Pos, operator.Text}),
			args[0].(ir.IrTerm),
			args[2].(ir.IrTerm))
	}
}

func binOpTypeApplicative() action {
	return func(args []any) any {
		operand1 := args[0].(ir.IrTerm)
		operator := args[1].(Token)
		typeApplicative := args[2].([]ir.IrType)
		operand2 := args[3].(ir.IrTerm)
		return newBinOpTerm(
			newAppTypeTerm(newIDTerm(ast.ID{operator.Pos, operator.Text}), typeApplicative),
			operand1,
			operand2)
	}
}

func unaryOp() action {
	return func(args []any) any {
		operator := args[0].(Token)
		return newUnaryOpTerm(
			newIDTerm(ast.ID{operator.Pos, operator.Text}),
			args[1].(ir.IrTerm))
	}
}

func unaryOpTypeApplicative() action {
	return func(args []any) any {
		operator := args[0].(Token)
		typeApplicative := args[1].([]ir.IrType)
		operand := args[2].(ir.IrTerm)
		return newUnaryOpTerm(
			newAppTypeTerm(newIDTerm(ast.ID{operator.Pos, operator.Text}), typeApplicative),
			operand)
	}
}

type Positional interface {
	GetPos() ir.Pos
	SetPos(ir.Pos)
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

		{"Packages -> Packages Package ;", listAppend[ast.Package](0, 1)},
		{"Packages -> Package ;", list[ast.Package](0)},

		{"Package -> module ImportID in ImportID", func(args []any) any {
			moduleID := args[1].(ast.ModuleID)
			filename := args[3].(ast.ModuleID)
			pos := makePos(args[0].(Token).Pos, filename.Pos)
			// TODO: Avoid conversion from ModuleID to ID.
			return ast.NewPackage(moduleID, ast.ID{filename.Pos, filename.Name}, pos)
		}},

		/* Module implementation file */

		{"Module -> module ImportID", func(args []any) any {
			id := args[1].(ast.ModuleID)
			return ast.Module{Header: ast.NewBaseFileHeader(id)}
		}},
		{"Module -> module ImportID ModuleImports", func(args []any) any {
			id := args[1].(ast.ModuleID)
			module := args[2].(ast.Module)
			module.Header = ast.NewBaseFileHeader(id)
			return module
		}},

		/* Module implementation file */

		{"Module -> implements ImportID", func(args []any) any {
			id := args[1].(ast.ModuleID)
			return ast.Module{Header: ast.NewImplementationFileHeader(id)}
		}},
		{"Module -> implements ImportID ModuleImports", func(args []any) any {
			id := args[1].(ast.ModuleID)
			module := args[2].(ast.Module)
			module.Header = ast.NewImplementationFileHeader(id)
			return module
		}},

		/* Modules */

		{"ModuleImports -> ImportsSection ModuleImpls", func(args []any) any {
			module := args[1].(ast.Module)
			module.Imports = args[0].(ast.Imports)
			return module
		}},
		{"ModuleImports -> ImportsSection", func(args []any) any {
			return ast.Module{Imports: args[0].(ast.Imports)}
		}},
		{"ModuleImports -> ModuleImpls", first()},

		{"ModuleImpls -> ImplsSection ModuleFlags", func(args []any) any {
			module := args[1].(ast.Module)
			module.Impls = args[0].(ast.Impls)
			return module
		}},
		{"ModuleImpls -> ImplsSection", func(args []any) any {
			return ast.Module{Impls: args[0].(ast.Impls)}
		}},
		{"ModuleImpls -> ModuleFlags", first()},

		{"ModuleFlags -> FlagsSection ModuleBody", func(args []any) any {
			module := args[1].(ast.Module)
			module.Flags = args[0].(ast.Flags)
			return module
		}},
		{"ModuleFlags -> FlagsSection", func(args []any) any {
			return ast.Module{Flags: args[0].(ast.Flags)}
		}},
		{"ModuleFlags -> ModuleBody", first()},

		/* Module body */

		{"ModuleBody -> Sources", func(args []any) any {
			return ast.Module{Body: args[0].([]ast.Source)}
		}},

		/* Imports section */

		{"ImportsSection -> imports { ImportIDs }", func(args []any) any {
			return ast.NewImports(args[2].([]ast.ModuleID), makePos2(args))
		}},

		{"ImportIDs -> ImportIDs ImportID ;", listAppend[ast.ModuleID](0, 1)},
		{"ImportIDs -> ImportID ;", list[ast.ModuleID](0)},

		{"ImportID -> Token", func(args []any) any {
			return newImportID(args[0].(Token))
		}},

		/* Impls section */

		{"ImplsSection -> impls { ImplIDs }", func(args []any) any {
			return ast.NewImpls(args[2].([]ast.ID), makePos2(args))
		}},

		{"ImplIDs -> ImplIDs ImplID ;", listAppend[ast.ID](0, 1)},
		{"ImplIDs -> ImplID ;", list[ast.ID](0)},

		{"ImplID -> Token", func(args []any) any {
			return newImplID(args[0].(Token))
		}},

		/* Flags section */

		{"FlagsSection -> flags { ImplIDs }", func(args []any) any {
			ids := args[2].([]ast.ID)
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
		{"Source -> export Function", func(args []any) any {
			source := args[1].(ast.Source)
			source.Function.Export = true
			return source
		}},

		/* Decl source */

		{"DeclSource -> DeclNoTerm", func(args []any) any {
			return newDeclSource(args[0].(ir.IrDecl))
		}},
		{"DeclSource -> decl TermDecl", func(args []any) any {
			return newDeclSource(args[1].(ir.IrDecl))
		}},
		{"DeclSource -> export TermDecl", func(args []any) any {
			decl := args[1].(ir.IrDecl)
			decl.Export = true
			return newDeclSource(decl)
		}},

		/* Function */

		{"Function -> fn ID TypeAbstraction FunctionArgs -> AppType Block", func(args []any) any {
			id := args[1].(ast.ID)
			tvars := args[2].([]ir.VarKind)
			funArgs := args[3].([]ir.IrDecl)
			retType := args[5].(ir.IrType)
			body := args[6].(ir.IrTerm)
			return newFunctionSource(
				makePos(id.Pos, body.Pos),
				ir.NewFunction(false /* export */, id.Value, tvars, funArgs, retType, body))
		}},
		{"Function -> fn ID FunctionArgs -> PrimaryType Block", func(args []any) any {
			id := args[1].(ast.ID)
			funArgs := args[2].([]ir.IrDecl)
			retType := args[4].(ir.IrType)
			body := args[5].(ir.IrTerm)
			return newFunctionSource(
				makePos(id.Pos, body.Pos),
				ir.NewFunction(false /* export */, id.Value, nil /* tvars */, funArgs, retType, body))
		}},

		{"FunctionArgs -> ( Args )", second()},
		{"FunctionArgs -> ( )", listNil[ir.IrDecl]()},

		{"Args -> Args , Arg", listAppend[ir.IrDecl](0, 2)},
		{"Args -> Arg", list[ir.IrDecl](0)},

		{"Arg -> ID : UnquantifiedType", func(args []any) any {
			return newTermDecl(args[0].(ast.ID), args[2].(ir.IrType), false /* export */)
		}},

		/* Decl */

		// Decl is used by annotations and DeclNoTerm is used by DeclSource. It
		// would be great if both annotations and DeclSource could reuse Decl, but
		// the grammar becomes ambiguous if TermDecl is not preceded by `decl` (or
		// another equivalent solution).
		{"Decl -> export UnexportedDecl", func(args []any) any {
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

		{"DeclNoTerm -> export UnexportedDeclNoTerm", func(args []any) any {
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

		{"ArrayType -> [ UnquantifiedType , Integer ]", func(args []any) any {
			elemType := args[1].(ir.IrType)
			length := args[3].(Integer).Value
			return newArrayType(makePos2(args), elemType, length)
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

		{"Fields -> Fields , Field", listAppend[ir.StructField](0, 2)},
		{"Fields -> Field", list[ir.StructField](0)},

		{"Field -> ID UnquantifiedType", func(args []any) any {
			return ir.StructField{args[0].(ast.ID).Value, args[1].(ir.IrType)}
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

		{"Tags -> Tags , Tag", listAppend[ir.VariantTag](0, 2)},
		{"Tags -> Tag", list[ir.VariantTag](0)},

		{"Tag -> ID UnquantifiedType", func(args []any) any {
			return ir.VariantTag{args[0].(ast.ID).Value, args[1].(ir.IrType)}
		}},

		/* ID */

		{"ID -> Token", func(args []any) any {
			token := args[0].(Token)

			if len(token.Text) == 0 || unicode.IsDigit(rune(token.Text[0])) {
				// TODO: Avoid panic.
				panic(fmt.Errorf("expected identifier; got %q; identifiers must begin with a non-digit character", token.Text))
			}

			return ast.ID{token.Pos, token.Text}
		}},

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
		{"ID -> ( != )", func(args []any) any {
			token := args[1].(Token)
			return ast.ID{token.Pos, token.Text}
		}},
		{"ID -> ( == )", func(args []any) any {
			token := args[1].(Token)
			return ast.ID{token.Pos, token.Text}
		}},
		{"ID -> ( > )", func(args []any) any {
			token := args[1].(Token)
			return ast.ID{token.Pos, token.Text}
		}},
		{"ID -> ( >= )", func(args []any) any {
			token := args[1].(Token)
			return ast.ID{token.Pos, token.Text}
		}},
		{"ID -> ( < )", func(args []any) any {
			token := args[1].(Token)
			return ast.ID{token.Pos, token.Text}
		}},
		{"ID -> ( <= )", func(args []any) any {
			token := args[1].(Token)
			return ast.ID{token.Pos, token.Text}
		}},
		{"ID -> ( + )", func(args []any) any {
			token := args[1].(Token)
			return ast.ID{token.Pos, token.Text}
		}},
		{"ID -> ( - )", func(args []any) any {
			token := args[1].(Token)
			return ast.ID{token.Pos, token.Text}
		}},
		{"ID -> ( * )", func(args []any) any {
			token := args[1].(Token)
			return ast.ID{token.Pos, token.Text}
		}},
		{"ID -> ( / )", func(args []any) any {
			token := args[1].(Token)
			return ast.ID{token.Pos, token.Text}
		}},
		{"ID -> ( ! )", func(args []any) any {
			token := args[1].(Token)
			return ast.ID{token.Pos, token.Text}
		}},

		/* Integer */

		{"Integer -> Token", func(args []any) any {
			token := args[0].(Token)

			value, err := parseNumber[int](token.Text)
			if err != nil {
				// TODO: Avoid panic.
				panic(fmt.Errorf("expected integer; got %q", token.Text))
			}

			return Integer{token.Pos, value}
		}},

		/* Term */

		{"Block -> { Terms }", func(args []any) any {
			return newBlockTerm(makePos2(args), args[1].([]ir.IrTerm))
		}},

		{"Terms -> Terms Term", listAppend[ir.IrTerm](0, 1)},
		{"Terms -> Term", list[ir.IrTerm](0)},

		{"Term -> IfTerm", first()},
		{"Term -> StatementTerm", first()},

		/* If term */

		{"IfTerm -> if Expression Block", func(args []any) any {
			condition := args[1].(ir.IrTerm)
			then := args[2].(ir.IrTerm)
			var elseTerm *ir.IrTerm
			return newIfTerm(
				makePos(args[0].(Token).Pos, then.Pos), condition, then, elseTerm)
		}},
		{"IfTerm -> if Expression Block else Block", func(args []any) any {
			condition := args[1].(ir.IrTerm)
			then := args[2].(ir.IrTerm)
			elseTerm := args[4].(ir.IrTerm)
			return newIfTerm(
				makePos(args[0].(Token).Pos, elseTerm.Pos), condition, then, &elseTerm)
		}},

		/* Statement */

		{"StatementTerm -> AssignTerm", first()},
		{"StatementTerm -> LetTerm", first()},
		{"StatementTerm -> ReturnTerm", first()},
		{"StatementTerm -> ExpressionNL", first()},

		/* Assign term */

		{"AssignTerm -> ID <- ExpressionNL", func(args []any) any {
			ret := args[0].(ast.ID)
			arg := args[2].(ir.IrTerm)
			return newAssignTerm(arg, newIDTerm(ret))
		}},
		{"AssignTerm -> TupleTerm <- ExpressionNL", func(args []any) any {
			ret := args[0].(ir.IrTerm)
			arg := args[2].(ir.IrTerm)
			return newAssignTerm(arg, ret)
		}},

		/* Let term */

		{"LetTerm -> let ID : UnquantifiedType = ExpressionNL", func(args []any) any {
			varName := args[1].(ast.ID)
			varType := args[3].(ir.IrType)
			value := args[5].(ir.IrTerm)
			return newLetTerm(varName, varType, value)
		}},

		/* Return term */

		{"ReturnTerm -> return ExpressionNL", func(args []any) any {
			return newReturnTerm(args[1].(ir.IrTerm))
		}},

		/* Expression NL */

		{"ExpressionNL -> MatchTerm", first()},
		{"ExpressionNL -> Expression ;", first()},

		/* Match term */

		{"MatchTerm -> case Expression { MatchArms }", func(args []any) any {
			term := args[1].(ir.IrTerm)
			arms := args[3].([]ir.MatchArm)
			return newMatchTerm(makePos(args[0].(Token).Pos, arms[len(arms)-1].Body.Pos), term, arms)
		}},

		{"MatchArms -> MatchArms | MatchArm", listAppend[ir.MatchArm](0, 2)},
		{"MatchArms -> | MatchArm", list[ir.MatchArm](1)},

		{"MatchArm -> ID ID = Block", func(args []any) any {
			tag := args[0].(ast.ID)
			arg := args[1].(ast.ID)
			body := args[3].(ir.IrTerm)
			return newMatchArm(tag, arg, body)
		}},
		{"MatchArm -> ID ID = Term", func(args []any) any {
			tag := args[0].(ast.ID)
			arg := args[1].(ast.ID)
			body := args[3].(ir.IrTerm)
			return newMatchArm(tag, arg, body)
		}},

		/* Expression */

		{"Expression -> LambdaTerm", first()},

		/* Lambda term */

		{"LambdaTerm -> \\ TypeAbstraction ID : UnquantifiedType = LambdaTerm", func(args []any) any {
			tvars := args[1].([]ir.VarKind)
			arg := args[2].(ast.ID)
			argType := args[4].(ir.IrType)
			body := args[6].(ir.IrTerm)
			return newLambdaTerm(
				makePos(args[0].(Token).Pos, body.Pos),
				tvars, []ir.ArgType{ir.ArgType{arg.Value, argType}}, body)
		}},
		{"LambdaTerm -> \\ ID : UnquantifiedType = LambdaTerm", func(args []any) any {
			var tvars []ir.VarKind
			arg := args[1].(ast.ID)
			argType := args[3].(ir.IrType)
			body := args[5].(ir.IrTerm)
			return newLambdaTerm(
				makePos(args[0].(Token).Pos, body.Pos),
				tvars, []ir.ArgType{ir.ArgType{arg.Value, argType}}, body)
		}},
		{"LambdaTerm -> Operator", first()},

		/* Operators */

		{"Operator -> Equality", first()},

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
			return newAppTermTerm(args[0].(ir.IrTerm), args[1].(ir.IrTerm))
		}},
		{"Applicative -> TypeApplicative", first()},

		/* Type applicative */

		{"TypeApplicative -> Primary TypeApplicativeArgs", func(args []any) any {
			return newAppTypeTerm(args[0].(ir.IrTerm), args[1].([]ir.IrType))
		}},
		{"TypeApplicative -> SetTerm", first()},

		{"TypeApplicativeArgs -> [ TupleTypeArgs ]", second()},
		{"TypeApplicativeArgs -> [ UnquantifiedType ]", list[ir.IrType](1)},

		/* Set term */

		// TODO: Get rid of 'set' keyword. This is only here to avoid grammar conflicts.
		{"SetTerm -> set Primary { SetValues }", func(args []any) any {
			term := args[1].(ir.IrTerm)
			values := args[3].([]ir.LabelValue)
			return newSetTerm(makePos2(args), term, values)
		}},
		{"SetTerm -> Primary", first()},

		{"SetValues -> SetValues , SetValue", listAppend[ir.LabelValue](0, 2)},
		{"SetValues -> SetValue", list[ir.LabelValue](0)},

		{"SetValue -> Token = Expression", func(args []any) any {
			token := args[0].(Token)
			value := args[2].(ir.IrTerm)
			return ir.LabelValue{token.Text, value}
		}},

		/* Primary */

		{"Primary -> InjectionTerm", first()},
		{"Primary -> LiteralTerm", first()},
		{"Primary -> ProjectionTerm", first()},
		{"Primary -> StructTerm", first()},
		{"Primary -> TupleTerm", first()},
		{"Primary -> ( Expression )", second()},

		/* Injection term */

		{"InjectionTerm -> variant { PrimaryType LiteralTerm = Expression }", func(args []any) any {
			variantType := args[2].(ir.IrType)
			tag := args[3].(ir.IrTerm)
			value := args[5].(ir.IrTerm)
			return newInjectionTerm(makePos2(args), variantType, tag, value)
		}},

		/* Struct term */

		{"StructTerm -> struct { }", func(args []any) any { return newStructTerm(makePos2(args), nil) }},
		{"StructTerm -> struct { LabelValues }", func(args []any) any {
			return newStructTerm(makePos2(args), args[2].([]ir.LabelValue))
		}},

		{"LabelValues -> LabelValues , LabelValue", listAppend[ir.LabelValue](0, 2)},
		{"LabelValues -> LabelValue", list[ir.LabelValue](0)},

		{"LabelValue -> ID = Expression", func(args []any) any {
			label := args[0].(ast.ID)
			value := args[2].(ir.IrTerm)
			return ir.LabelValue{label.Value, value}
		}},

		/* Literal term */

		{"LiteralTerm -> Token", func(args []any) any {
			return newLiteralTerm(args[0].(Token))
		}},

		/* Projection term */

		{"ProjectionTerm -> Primary -> Token", func(args []any) any {
			term := args[0].(ir.IrTerm)
			label := args[2].(Token)
			return newProjectionTerm(makePos(term.Pos, label.Pos), term, label.Text)
		}},

		/* Tuple term */

		{"TupleTerm -> ( )", func(args []any) any {
			return newTupleTerm(makePos2(args), nil)
		}},
		{"TupleTerm -> ( TupleTermArgs )", func(args []any) any {
			return newTupleTerm(makePos2(args), args[1].([]ir.IrTerm))
		}},

		{"TupleTermArgs -> TupleTermArgs , Expression", listAppend[ir.IrTerm](0, 2)},
		{"TupleTermArgs -> Expression , Expression", list[ir.IrTerm](0, 2)},
	}
}
