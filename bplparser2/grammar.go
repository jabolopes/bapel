package bplparser2

import (
	"fmt"
	"math"
	"strings"
	"unicode"

	"github.com/jabolopes/bapel/bplparser"
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
	pos1.EndLineNum = pos2.EndLineNum
	return pos1
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

func newImportSource(id ID) bplparser.Source {
	source := bplparser.NewImportSource(id.Value)
	source.Pos = id.Pos
	return source
}

func newSectionSource(id Token, decls []ir.IrDecl, endPos ir.Pos) bplparser.Source {
	source := bplparser.NewSectionSource(id.Text, decls)
	source.Pos = makePos(id.Pos, endPos)
	return source
}

func newAliasDecl(id ID, typ ir.IrType) ir.IrDecl {
	decl := ir.NewAliasDecl(id.Value, typ)
	decl.Pos = makePos(id.Pos, typ.Pos)
	return decl
}

func newTermDecl(id ID, typ ir.IrType) ir.IrDecl {
	decl := ir.NewTermDecl(id.Value, typ)
	decl.Pos = makePos(id.Pos, typ.Pos)
	return decl
}

func newNameDecl(id ID) ir.IrDecl {
	decl := ir.NewNameDecl(id.Value)
	decl.Pos = id.Pos
	return decl
}

func newFunctionSource(pos ir.Pos, fun ir.IrFunction) bplparser.Source {
	fun.Pos = pos
	source := bplparser.NewFunctionSource(fun)
	source.Pos = pos
	return source
}

func newTypeDefSource(export bool, decl ir.IrDecl) bplparser.Source {
	source := bplparser.NewTypeDefSource(export, decl)
	source.Pos = decl.Pos
	return source
}

func newComponentSource(pos ir.Pos, component ir.IrComponent) bplparser.Source {
	source := bplparser.NewComponentSource(component)
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

func newVarType(id ID) ir.IrType {
	typ := ir.NewVarType(id.Value)
	typ.Pos = id.Pos
	return typ
}

func newNameType(id ID) ir.IrType {
	typ := ir.NewNameType(id.Value)
	typ.Pos = id.Pos
	return typ
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

func newConstTerm(pos ir.Pos, text string, value int64) ir.IrTerm {
	term := ir.NewConstTerm(text, value)
	term.Pos = pos
	return term
}

func newIDTerm(id ID) ir.IrTerm {
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

func newLetTerm(varName ID, varType ir.IrType, value ir.IrTerm) ir.IrTerm {
	term := ir.NewLetTerm(varName.Value, varType, value)
	term.Pos = makePos(varName.Pos, value.Pos)
	return term
}

func newProjectionTerm(term, label ir.IrTerm) ir.IrTerm {
	proj := ir.NewProjectionTerm(term, label)
	proj.Pos = makePos(term.Pos, label.Pos)
	return proj
}

func newReturnTerm(expr ir.IrTerm) ir.IrTerm {
	term := ir.NewReturnTerm(expr)
	term.Pos = expr.Pos
	return term
}

func newIndexGetTerm(arg1, arg2 ir.IrTerm) ir.IrTerm {
	term := ir.NewIndexGetTerm(arg1, arg2)
	term.Pos = makePos(arg1.Pos, arg2.Pos)
	return term
}

func newIndexSetTerm(arg1, arg2, arg3 ir.IrTerm) ir.IrTerm {
	term := ir.NewIndexSetTerm(arg1, arg2, arg3)
	term.Pos = makePos(arg1.Pos, arg3.Pos)
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

func newLiteralTerm(token Token) ir.IrTerm {
	if unicode.IsDigit(rune(token.Text[0])) {
		value, err := parseNumber[int64](token.Text)
		if err != nil {
			// TODO: Avoid panic.
			panic(fmt.Errorf("expected integer; got %q", token.Text))
		}

		return newConstTerm(token.Pos, token.Text, value)
	}

	return newVarTerm(token.Pos, token.Text)
}

type action = func(args []any) any

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

func listCons[T any](is ...int) action {
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
		return newBinOpTerm(newIDTerm(ID{operator.Pos, operator.Text}), args[0].(ir.IrTerm), args[2].(ir.IrTerm))
	}
}

func unaryOp() action {
	return func(args []any) any {
		operator := args[0].(Token)
		return newUnaryOpTerm(newIDTerm(ID{operator.Pos, operator.Text}), args[1].(ir.IrTerm))
	}
}

type Positional interface {
	GetPos() ir.Pos
	SetPos(ir.Pos)
}

func NewGrammar(initial grammar.ProductionLine) []grammar.ProductionLine {
	return []grammar.ProductionLine{
		initial,

		/* Source */

		{"Anys -> Anys Any", listAppend[bplparser.Source](0, 1)},
		{"Anys -> Any", listCons[bplparser.Source](0)},

		{"Any -> Import", first()},
		{"Any -> DeclsSection", first()},
		{"Any -> ExportsSection", first()},
		{"Any -> Function", first()},
		{"Any -> export Function", func(args []any) any {
			source := args[1].(bplparser.Source)
			source.Function.Export = true
			return source
		}},
		{"Any -> StructSource", first()},
		{"Any -> VariantSource", first()},
		{"Any -> Component", first()},

		/* Import */

		{"Import -> import ID", func(args []any) any {
			return newImportSource(args[1].(ID))
		}},

		/* Section */

		{"DeclsSection -> decls { Decls }", func(args []any) any {
			return newSectionSource(args[0].(Token), args[2].([]ir.IrDecl), args[3].(Token).Pos)
		}},
		{"ExportsSection -> exports { Decls }", func(args []any) any {
			return newSectionSource(args[0].(Token), args[2].([]ir.IrDecl), args[3].(Token).Pos)
		}},

		/* Decls */

		{"Decls -> Decls Decl", listAppend[ir.IrDecl](0, 1)},
		{"Decls -> Decl", listCons[ir.IrDecl](0)},

		{"Decl -> StructDecl ;", first()},
		{"Decl -> TermDecl ;", first()},
		{"Decl -> TypeDecl ;", first()},
		{"Decl -> VariantDecl ;", first()},

		{"StructDecl -> type ID TypeAbstraction = StructType", func(args []any) any {
			id := args[1].(ID)
			tvars := args[2].([]ir.VarKind)
			structType := args[4].(ir.IrType)
			return newAliasDecl(id, ir.LambdaVars(tvars, structType))
		}},
		{"StructDecl -> type ID = StructType", func(args []any) any {
			id := args[1].(ID)
			var tvars []ir.VarKind
			structType := args[3].(ir.IrType)
			return newAliasDecl(id, ir.LambdaVars(tvars, structType))
		}},

		{"TermDecl -> ID : QuantifiedType", func(args []any) any {
			return newTermDecl(args[0].(ID), args[2].(ir.IrType))
		}},

		{"TypeDecl -> type ID", func(args []any) any {
			return newNameDecl(args[1].(ID))
		}},

		{"VariantDecl -> type ID = VariantType", func(args []any) any {
			id := args[1].(ID)
			var tvars []ir.VarKind
			variantType := args[3].(ir.IrType)
			return newAliasDecl(id, ir.LambdaVars(tvars, variantType))
		}},
		{"VariantDecl -> type ID TypeAbstraction = VariantType", func(args []any) any {
			id := args[1].(ID)
			tvars := args[2].([]ir.VarKind)
			variantType := args[4].(ir.IrType)
			return newAliasDecl(id, ir.LambdaVars(tvars, variantType))
		}},

		/* Function */

		{"Function -> fn ID TypeAbstraction FunctionArgs -> PrimaryType Block", func(args []any) any {
			id := args[1].(ID)
			tvars := args[2].([]ir.VarKind)
			funArgs := args[3].([]ir.IrDecl)
			retType := args[5].(ir.IrType)
			body := args[6].(ir.IrTerm)
			return newFunctionSource(
				makePos(id.Pos, body.Pos),
				ir.NewFunction(false /* export */, id.Value, tvars, funArgs, retType, body))
		}},
		{"Function -> fn ID FunctionArgs -> PrimaryType Block", func(args []any) any {
			id := args[1].(ID)
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
		{"Args -> Arg", listCons[ir.IrDecl](0)},

		{"Arg -> ID : UnquantifiedType", func(args []any) any {
			return newTermDecl(args[0].(ID), args[2].(ir.IrType))
		}},

		/* Struct source */

		{"StructSource -> StructDecl", func(args []any) any {
			return newTypeDefSource(false /* export */, args[0].(ir.IrDecl))
		}},
		{"StructSource -> export StructDecl", func(args []any) any {
			return newTypeDefSource(true /* export */, args[1].(ir.IrDecl))
		}},

		/* Variant source */

		{"VariantSource -> VariantDecl", func(args []any) any {
			return newTypeDefSource(false /* export */, args[0].(ir.IrDecl))
		}},
		{"VariantSource -> export VariantDecl", func(args []any) any {
			return newTypeDefSource(true /* export */, args[1].(ir.IrDecl))
		}},

		/* Component */

		{"Component -> component [ UnquantifiedType , Integer ]", func(args []any) any {
			elemType := args[2].(ir.IrType)
			length := args[4].(Integer).Value
			return newComponentSource(
				makePos(args[0].(Token).Pos, args[5].(Token).Pos),
				ir.NewComponent(elemType, length))
		}},

		/* Type variables */

		{"TypeAbstraction -> [ Tvars ]", second()},

		{"Tvars -> Tvars , Tvar", listAppend[ir.VarKind](0, 2)},
		{"Tvars -> Tvar", listCons[ir.VarKind](0)},

		{"Tvar -> ' ID", func(args []any) any {
			return ir.VarKind{args[1].(ID).Value, ir.NewTypeKind()}
		}},

		/* Quantified type */

		{"QuantifiedType -> UnquantifiedType", func(args []any) any {
			return newQuantifiedType(args[0].(ir.IrType))
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
			return newVarType(args[1].(ID))
		}},
		{"PrimaryType -> ID", func(args []any) any {
			return newNameType(args[0].(ID))
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
		{"Fields -> Field", listCons[ir.StructField](0)},

		{"Field -> ID UnquantifiedType", func(args []any) any {
			return ir.StructField{args[0].(ID).Value, args[1].(ir.IrType)}
		}},

		/* Tuple type */

		{"TupleType -> ( )", func(args []any) any {
			return newTupleType(makePos2(args), nil)
		}},
		{"TupleType -> ( TupleTypeArgs )", func(args []any) any {
			return newTupleType(makePos2(args), args[1].([]ir.IrType))
		}},

		{"TupleTypeArgs -> TupleTypeArgs , UnquantifiedType", listAppend[ir.IrType](0, 2)},
		{"TupleTypeArgs -> UnquantifiedType , UnquantifiedType", listCons[ir.IrType](0, 2)},

		/* Variant type */

		{"VariantType -> variant { }", func(args []any) any {
			return newVariantType(makePos2(args), nil)
		}},
		{"VariantType -> variant { Tags }", func(args []any) any {
			return newVariantType(makePos2(args), args[2].([]ir.VariantTag))
		}},

		{"Tags -> Tags , Tag", listAppend[ir.VariantTag](0, 2)},
		{"Tags -> Tag", listCons[ir.VariantTag](0)},

		{"Tag -> ID UnquantifiedType", func(args []any) any {
			return ir.VariantTag{args[0].(ID).Value, args[1].(ir.IrType)}
		}},

		/* ID */

		{"ID -> Token", func(args []any) any {
			token := args[0].(Token)

			if len(token.Text) == 0 || unicode.IsDigit(rune(token.Text[0])) {
				// TODO: Avoid panic.
				panic(fmt.Errorf("expected identifier; got %q; identifiers must begin with a non-digit character", token.Text))
			}

			return ID{token.Pos, token.Text}
		}},

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
		{"Terms -> Term", listCons[ir.IrTerm](0)},

		{"Term -> IfTerm", first()},
		{"Term -> StatementTerm", first()},

		/* Statement */

		{"StatementTerm -> AssignTerm ;", first()},
		{"StatementTerm -> LetTerm ;", first()},
		{"StatementTerm -> ReturnTerm ;", first()},
		{"StatementTerm -> Expression ;", first()},

		/* Assign term */

		{"AssignTerm -> ID <- Expression", func(args []any) any {
			ret := args[0].(ID)
			arg := args[2].(ir.IrTerm)
			return newAssignTerm(arg, newIDTerm(ret))
		}},
		{"AssignTerm -> TupleTerm <- Expression", func(args []any) any {
			ret := args[0].(ir.IrTerm)
			arg := args[2].(ir.IrTerm)
			return newAssignTerm(arg, ret)
		}},

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

		/* Let term */

		{"LetTerm -> let ID : UnquantifiedType = Expression", func(args []any) any {
			varName := args[1].(ID)
			varType := args[3].(ir.IrType)
			value := args[5].(ir.IrTerm)
			return newLetTerm(varName, varType, value)
		}},

		/* Return term */

		{"ReturnTerm -> return Expression", func(args []any) any {
			return newReturnTerm(args[1].(ir.IrTerm))
		}},

		/* Expression */

		{"Expression -> equality", first()},

		{"equality -> equality != comparison", binOp()},
		{"equality -> equality == comparison", binOp()},
		{"equality -> comparison", first()},

		{"comparison -> comparison > additive", binOp()},
		{"comparison -> comparison >= additive", binOp()},
		{"comparison -> comparison < additive", binOp()},
		{"comparison -> comparison <= additive", binOp()},
		{"comparison -> additive", first()},

		{"additive -> additive + multiplicative", binOp()},
		{"additive -> additive - multiplicative", binOp()},
		{"additive -> multiplicative", first()},

		{"multiplicative -> multiplicative * unary", binOp()},
		{"multiplicative -> multiplicative / unary", binOp()},
		{"multiplicative -> unary", first()},

		{"unary -> ! unary", unaryOp()},
		{"unary -> - unary", unaryOp()},
		{"unary -> Applicative", first()},

		{"Applicative -> Applicative Primary", func(args []any) any {
			return newAppTermTerm(args[0].(ir.IrTerm), args[1].(ir.IrTerm))
		}},
		{"Applicative -> TypeApplicative", first()},

		/* Type applicative */

		{"TypeApplicative -> Primary TypeApplicativeArgs", func(args []any) any {
			return newAppTypeTerm(args[0].(ir.IrTerm), args[1].([]ir.IrType))
		}},
		{"TypeApplicative -> ProjectionTerm", first()},

		{"TypeApplicativeArgs -> [ TupleTypeArgs ]", second()},
		{"TypeApplicativeArgs -> [ UnquantifiedType ]", listCons[ir.IrType](1)},

		/* Projection term */

		{"ProjectionTerm -> ProjectionTerm -> Primary", func(args []any) any {
			return newProjectionTerm(args[0].(ir.IrTerm), args[2].(ir.IrTerm))
		}},
		{"ProjectionTerm -> Primary", first()},

		/* Primary */

		{"Primary -> InjectionTerm", first()},
		{"Primary -> Index.get Primary LiteralTerm", func(args []any) any {
			return newIndexGetTerm(args[1].(ir.IrTerm), args[2].(ir.IrTerm))
		}},
		{"Primary -> Index.set Primary LiteralTerm Primary", func(args []any) any {
			return newIndexSetTerm(args[1].(ir.IrTerm), args[2].(ir.IrTerm), args[3].(ir.IrTerm))
		}},
		{"Primary -> LambdaTerm", first()},
		{"Primary -> LiteralTerm", first()},
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
		{"LabelValues -> LabelValue", listCons[ir.LabelValue](0)},

		{"LabelValue -> ID = Expression", func(args []any) any {
			label := args[0].(ID)
			value := args[2].(ir.IrTerm)
			return ir.LabelValue{label.Value, value}
		}},

		/* Lambda term */

		{"LambdaTerm -> \\ TypeAbstraction ID : UnquantifiedType = Primary", func(args []any) any {
			tvars := args[1].([]ir.VarKind)
			arg := args[2].(ID)
			argType := args[4].(ir.IrType)
			body := args[6].(ir.IrTerm)
			return newLambdaTerm(
				makePos(args[0].(Token).Pos, body.Pos),
				tvars, []ir.ArgType{ir.ArgType{arg.Value, argType}}, body)
		}},
		{"LambdaTerm -> \\ ID : UnquantifiedType = Primary", func(args []any) any {
			var tvars []ir.VarKind
			arg := args[1].(ID)
			argType := args[3].(ir.IrType)
			body := args[5].(ir.IrTerm)
			return newLambdaTerm(
				makePos(args[0].(Token).Pos, body.Pos),
				tvars, []ir.ArgType{ir.ArgType{arg.Value, argType}}, body)
		}},

		/* Literal term */

		{"LiteralTerm -> Token", func(args []any) any {
			return newLiteralTerm(args[0].(Token))
		}},

		/* Tuple term */

		{"TupleTerm -> ( )", func(args []any) any {
			return newTupleTerm(makePos2(args), nil)
		}},
		{"TupleTerm -> ( TupleTermArgs )", func(args []any) any {
			return newTupleTerm(makePos2(args), args[1].([]ir.IrTerm))
		}},

		{"TupleTermArgs -> TupleTermArgs , Expression", listAppend[ir.IrTerm](0, 2)},
		{"TupleTermArgs -> Expression , Expression", listCons[ir.IrTerm](0, 2)},
	}
}
