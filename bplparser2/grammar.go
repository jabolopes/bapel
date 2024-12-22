package bplparser2

import (
	"fmt"
	"math"
	"strconv"

	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
	"github.com/jabolopes/go-lalr1/grammar"
)

func newUnaryOpTerm(id string, term ir.IrTerm) ir.IrTerm {
	if id == "-" {
		// 0 - $term
		return ir.CallPF(id, nil /* types */, ir.Number(0), term)
	}

	return ir.Call(id, term)
}

func newBinOpTerm(id string, t1, t2 ir.IrTerm) ir.IrTerm {
	return ir.Call(id, t1, t2)
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

func binOp(operator string) action {
	return func(args []any) any {
		return newBinOpTerm(operator, args[0].(ir.IrTerm), args[2].(ir.IrTerm))
	}
}

func unaryOp(operator string) action {
	return func(args []any) any {
		return newUnaryOpTerm(operator, args[1].(ir.IrTerm))
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
		{"Any -> Struct", first()},
		{"Any -> export Struct", func(args []any) any {
			source := args[1].(bplparser.Source)
			source.TypeDef.Export = true
			return source
		}},
		{"Any -> Component", first()},
		{"Any -> Term", func(args []any) any {
			return bplparser.NewTermSource(args[0].(ir.IrTerm))
		}},

		/* Import */

		{"Import -> import ID", func(args []any) any {
			return bplparser.NewImportSource(args[1].(ID).Value)
		}},

		/* Section */

		{"DeclsSection -> decls { Decls }", func(args []any) any {
			decls := args[2].([]ir.IrDecl)
			return bplparser.NewSectionSource("decls", decls)
		}},
		{"ExportsSection -> exports { Decls }", func(args []any) any {
			decls := args[2].([]ir.IrDecl)
			return bplparser.NewSectionSource("exports", decls)
		}},

		{"Decls -> Decls Decl", listAppend[ir.IrDecl](0, 1)},
		{"Decls -> Decl", listCons[ir.IrDecl](0)},

		{"Decl -> Struct", func(args []any) any {
			return args[0].(bplparser.Source).TypeDef.Decl
		}},
		{"Decl -> TermDecl", first()},
		{"Decl -> TypeDecl", first()},

		{"TermDecl -> ID : SingleQuantifiedType", func(args []any) any {
			id := args[0].(ID).Value
			typ := args[2].(ir.IrType)
			return ir.NewTermDecl(id, typ)
		}},

		{"TypeDecl -> type ID ;", func(args []any) any {
			id := args[1].(ID).Value
			return ir.NewNameDecl(id)
		}},

		/* Function */

		{"Function -> func ID TypeAbstraction FunctionBindList -> FunctionBindList Block", func(args []any) any {
			id := args[1].(ID).Value
			tvars := args[2].([]ir.VarKind)
			funArgs := args[3].([]ir.IrDecl)
			funRets := args[5].([]ir.IrDecl)
			body := args[6].(ir.IrTerm)
			return bplparser.NewFunctionSource(
				ir.NewFunction(false /* export */, id, tvars, funArgs, funRets, body))
		}},
		{"Function -> func ID FunctionBindList -> FunctionBindList Block", func(args []any) any {
			id := args[1].(ID).Value
			funArgs := args[2].([]ir.IrDecl)
			funRets := args[4].([]ir.IrDecl)
			body := args[5].(ir.IrTerm)
			return bplparser.NewFunctionSource(
				ir.NewFunction(false /* export */, id, nil /* tvars */, funArgs, funRets, body))
		}},

		{"FunctionBindList -> ( Args )", func(args []any) any {
			return args[1].([]ir.IrDecl)
		}},
		{"FunctionBindList -> ( )", func(args []any) any {
			return []ir.IrDecl{}
		}},

		{"Args -> Args , Arg", listAppend[ir.IrDecl](0, 2)},
		{"Args -> Arg", listCons[ir.IrDecl](0)},

		{"Arg -> ID UnquantifiedType", func(args []any) any {
			return ir.NewTermDecl(args[0].(ID).Value, args[1].(ir.IrType))
		}},

		/* Struct */

		{"Struct -> struct ID TypeAbstraction StructType", func(args []any) any {
			id := args[1].(ID).Value
			tvars := args[2].([]ir.VarKind)
			structType := args[3].(ir.IrType)

			lambdaType := ir.LambdaVars(tvars, structType)
			return bplparser.NewTypeDefSource(false /* export */, ir.NewAliasDecl(id, lambdaType))
		}},
		{"Struct -> struct ID StructType", func(args []any) any {
			id := args[1].(ID).Value
			var tvars []ir.VarKind
			structType := args[2].(ir.IrType)

			lambdaType := ir.LambdaVars(tvars, structType)
			return bplparser.NewTypeDefSource(false /* export */, ir.NewAliasDecl(id, lambdaType))
		}},

		/* Component */

		{"Component -> component [ UnquantifiedType , Integer ]", func(args []any) any {
			elemType := args[2].(ir.IrType)
			length := args[4].(Integer).Value
			return bplparser.NewComponentSource(ir.NewComponent(elemType, length))
		}},

		/* Type variables */

		{"TypeAbstraction -> [ Tvars ]", func(args []any) any {
			return args[1].([]ir.VarKind)
		}},

		{"Tvars -> Tvars , Tvar", listAppend[ir.VarKind](0, 2)},
		{"Tvars -> Tvar", listCons[ir.VarKind](0)},

		{"Tvar -> ' ID", func(args []any) any {
			return ir.VarKind{args[1].(ID).Value, ir.NewTypeKind()}
		}},

		/* Single quantified type */

		{"SingleQuantifiedType -> QuantifiedType ;", first()},

		/* Type */

		{"QuantifiedType -> UnquantifiedType", func(args []any) any {
			return ir.QuantifyType(args[0].(ir.IrType))
		}},

		/* Type */

		{"UnquantifiedType -> ForallType", first()},

		/* Forall type */

		{"ForallType -> forall TypeAbstraction FunctionType", func(args []any) any {
			tvars := args[1].([]ir.VarKind)
			subType := args[2].(ir.IrType)
			return ir.ForallVars(tvars, subType)
		}},
		{"ForallType -> FunctionType", first()},

		/* Function type */

		{"FunctionType -> AppType -> FunctionType", func(args []any) any {
			arg := args[0].(ir.IrType)
			ret := args[2].(ir.IrType)
			return ir.NewFunctionType(arg, ret)
		}},
		{"FunctionType -> AppType", first()},

		/* App type */

		{"AppType -> AppType PrimaryType", func(args []any) any {
			argType := args[0].(ir.IrType)
			retType := args[1].(ir.IrType)
			return ir.NewAppType(argType, retType)
		}},
		{"AppType -> PrimaryType", first()},

		/* Simple Type */

		{"PrimaryType -> ArrayType", first()},
		{"PrimaryType -> StructType", first()},
		{"PrimaryType -> TupleType", first()},
		{"PrimaryType -> ' ID", func(args []any) any {
			return ir.NewVarType(args[1].(ID).Value)
		}},
		{"PrimaryType -> ID", func(args []any) any {
			return ir.NewNameType(args[0].(ID).Value)
		}},
		{"PrimaryType -> ( UnquantifiedType )", second()},

		/* Array type */

		{"ArrayType -> [ UnquantifiedType , Integer ]", func(args []any) any {
			typ := args[1].(ir.IrType)
			length := args[3].(Integer).Value
			return ir.NewArrayType(typ, length)
		}},
		{"ArrayType -> [ UnquantifiedType ]", func(args []any) any {
			typ := args[1].(ir.IrType)
			length := math.MaxInt
			return ir.NewArrayType(typ, length)
		}},

		/* Tuple type */

		{"TupleType -> ( )", func(args []any) any {
			return ir.NewTupleType(nil)
		}},
		{"TupleType -> ( TupleTypeArgs )", func(args []any) any {
			return ir.NewTupleType(args[1].([]ir.IrType))
		}},

		{"TupleTypeArgs -> TupleTypeArgs , UnquantifiedType", listAppend[ir.IrType](0, 2)},
		{"TupleTypeArgs -> UnquantifiedType , UnquantifiedType", listCons[ir.IrType](0, 2)},

		/* Struct type */

		{"StructType -> { }", func(args []any) any {
			return ir.NewStructType(nil)
		}},
		{"StructType -> { Fields }", func(args []any) any {
			return ir.NewStructType(args[1].([]ir.StructField))
		}},

		{"Fields -> Fields , Field", listAppend[ir.StructField](0, 2)},
		{"Fields -> Field", listCons[ir.StructField](0)},

		{"Field -> ID UnquantifiedType", func(args []any) any {
			return ir.StructField{args[0].(ID).Value, args[1].(ir.IrType)}
		}},

		/* ID */

		{"ID -> Token", func(args []any) any {
			token := args[0].(Token)
			return ID{token.Pos, token.Token.Text}
		}},

		{"Integer -> Token", func(args []any) any {
			token := args[0].(Token)

			value, err := strconv.Atoi(token.Token.Text)
			if err != nil {
				panic(err)
			}

			return Integer{token.Pos, value}
		}},

		/* Term */

		{"Block -> { Terms }", func(args []any) any {
			return ir.NewBlockTerm(args[1].([]ir.IrTerm))
		}},
		{"Block -> { }", func(args []any) any {
			return ir.NewBlockTerm(nil)
		}},

		{"Terms -> Terms Term", listAppend[ir.IrTerm](0, 1)},
		{"Terms -> Term", listCons[ir.IrTerm](0)},

		{"Term -> IfTerm", first()},
		{"Term -> StatementTerm", first()},

		/* Statement */

		{"StatementTerm -> AssignTerm", first()},
		{"StatementTerm -> LetTerm", first()},
		{"StatementTerm -> SingleExpression", first()},

		/* Assign term */

		{"AssignTerm -> ID <- SingleExpression", func(args []any) any {
			ret := args[0].(ID).Value
			term := args[2].(ir.IrTerm)
			return ir.NewAssignTerm(term, ir.ID(ret))
		}},
		{"AssignTerm -> TupleTerm <- SingleExpression", func(args []any) any {
			rets := args[0].(ir.IrTerm)
			term := args[2].(ir.IrTerm)
			return ir.NewAssignTerm(term, rets)
		}},

		/* If term */

		{"IfTerm -> if Expression Block", func(args []any) any {
			negate := false
			var types []ir.IrType
			condition := args[1].(ir.IrTerm)
			then := args[2].(ir.IrTerm)
			var elseTerm *ir.IrTerm
			return ir.NewIfTerm(negate, types, condition, then, elseTerm)
		}},
		{"IfTerm -> if Expression Block else Block", func(args []any) any {
			negate := false
			var types []ir.IrType
			condition := args[1].(ir.IrTerm)
			then := args[2].(ir.IrTerm)
			elseTerm := args[4].(ir.IrTerm)
			return ir.NewIfTerm(negate, types, condition, then, &elseTerm)
		}},
		{"IfTerm -> if TypeApplicativeArgs Expression Block", func(args []any) any {
			negate := false
			types := args[1].([]ir.IrType)
			condition := args[2].(ir.IrTerm)
			then := args[3].(ir.IrTerm)
			var elseTerm *ir.IrTerm
			return ir.NewIfTerm(negate, types, condition, then, elseTerm)
		}},
		{"IfTerm -> if TypeApplicativeArgs Expression Block else Block", func(args []any) any {
			negate := false
			types := args[1].([]ir.IrType)
			condition := args[2].(ir.IrTerm)
			then := args[3].(ir.IrTerm)
			elseTerm := args[5].(ir.IrTerm)
			return ir.NewIfTerm(negate, types, condition, then, &elseTerm)
		}},

		// TODO: Use rules above to get ! expression working.
		{"IfTerm -> if not Expression Block", func(args []any) any {
			negate := true
			var types []ir.IrType
			condition := args[2].(ir.IrTerm)
			then := args[3].(ir.IrTerm)
			var elseTerm *ir.IrTerm
			return ir.NewIfTerm(negate, types, condition, then, elseTerm)
		}},
		{"IfTerm -> if not Expression Block else Block", func(args []any) any {
			negate := true
			var types []ir.IrType
			condition := args[2].(ir.IrTerm)
			then := args[3].(ir.IrTerm)
			elseTerm := args[5].(ir.IrTerm)
			return ir.NewIfTerm(negate, types, condition, then, &elseTerm)
		}},

		/* Let term */

		{"LetTerm -> let ID SingleQuantifiedType", func(args []any) any {
			id := args[1].(ID).Value
			typ := args[2].(ir.IrType)
			var arg *ir.IrTerm
			return ir.NewLetTerm(ir.NewTermDecl(id, typ), arg)
		}},
		{"LetTerm -> let ID QuantifiedType = SingleExpression", func(args []any) any {
			id := args[1].(ID).Value
			typ := args[2].(ir.IrType)
			arg := args[4].(ir.IrTerm)
			return ir.NewLetTerm(ir.NewTermDecl(id, typ), &arg)
		}},

		/* Single expression */

		{"SingleExpression -> Expression ;", first()},

		/* Expression */

		{"Expression -> equality", first()},

		{"equality -> equality != comparison", binOp("!=")},
		{"equality -> equality == comparison", binOp("==")},
		{"equality -> comparison", first()},

		{"comparison -> comparison > additive", binOp(">")},
		{"comparison -> comparison >= additive", binOp(">=")},
		{"comparison -> comparison < additive", binOp("<")},
		{"comparison -> comparison <= additive", binOp("<=")},
		{"comparison -> additive", first()},

		{"additive -> additive + multiplicative", binOp("+")},
		{"additive -> additive - multiplicative", binOp("-")},
		{"additive -> multiplicative", first()},

		{"multiplicative -> multiplicative * unary", binOp("*")},
		{"multiplicative -> multiplicative / unary", binOp("/")},
		{"multiplicative -> unary", first()},

		{"unary -> ! unary", unaryOp("!")},
		{"unary -> - unary", unaryOp("-")},
		{"unary -> Applicative", first()},

		{"Applicative -> Index.get TypeApplicative TypeApplicative", func(args []any) any {
			return ir.NewIndexGetTerm(args[1].(ir.IrTerm), args[2].(ir.IrTerm))
		}},
		{"Applicative -> Index.set TypeApplicative TypeApplicative TypeApplicative", func(args []any) any {
			return ir.NewIndexSetTerm(args[1].(ir.IrTerm), args[2].(ir.IrTerm), args[3].(ir.IrTerm))
		}},
		{"Applicative -> Applicative TypeApplicative", func(args []any) any {
			return ir.NewAppTermTerm(args[0].(ir.IrTerm), args[1].(ir.IrTerm))
		}},
		{"Applicative -> TypeApplicative", first()},

		/* Type applicative */

		{"TypeApplicative -> Primary TypeApplicativeArgs", func(args []any) any {
			term := args[0].(ir.IrTerm)
			for _, typ := range args[1].([]ir.IrType) {
				term = ir.NewAppTypeTerm(term, typ)
			}
			return term
		}},
		{"TypeApplicative -> Primary", first()},

		{"TypeApplicativeArgs -> [ TupleTypeArgs ]", second()},
		{"TypeApplicativeArgs -> [ UnquantifiedType ]", listCons[ir.IrType](1)},

		/* Primary */

		{"Primary -> TupleTerm", first()},
		{"Primary -> Token", func(args []any) any {
			token := args[0].(Token).Token
			switch {
			case token.Case == parser.IDToken:
				return ir.ID(token.Text)
			case token.Case == parser.NumberToken:
				return ir.NewLiteralTerm(ir.NewNumberLiteral(token.Text, token.Value))
			default:
				panic(fmt.Errorf("unhandled %T %d", token.Case, token.Case))
			}
		}},
		{"Primary -> ( Expression )", second()},

		/* Tuple term */

		{"TupleTerm -> ( )", func(args []any) any {
			return ir.NewTupleTerm(nil)
		}},
		{"TupleTerm -> ( TupleTermArgs )", func(args []any) any {
			return ir.NewTupleTerm(args[1].([]ir.IrTerm))
		}},

		{"TupleTermArgs -> TupleTermArgs , Expression", listAppend[ir.IrTerm](0, 2)},
		{"TupleTermArgs -> Expression , Expression", listCons[ir.IrTerm](0, 2)},
	}
}
