package bplparser2

import (
	"math"
	"strconv"

	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
	"github.com/jabolopes/go-lalr1/grammar"
)

func first(args []any) any { return args[0] }

func newUnaryOpTerm(id string, term ir.IrTerm) ir.IrTerm {
	if id == "-" {
		// 0 - $term
		args := []ir.IrTerm{ir.NewTokenTerm(parser.NewNumberToken(0)), term}
		return ir.CallPF(id, nil /* types */, ir.NewTupleTerm(args))
	}

	return ir.Call(id, term)
}

func newBinOpTerm(id string, t1, t2 ir.IrTerm) ir.IrTerm {
	return ir.Call(id, t1, t2)
}

func NewGrammar(initial grammar.ProductionLine) []grammar.ProductionLine {
	return []grammar.ProductionLine{
		initial,

		/* Source */

		{"Anys -> Anys Any", func(args []any) any {
			return append(args[0].([]bplparser.Source), args[1].(bplparser.Source))
		}},
		{"Anys -> Any", func(args []any) any {
			return []bplparser.Source{args[0].(bplparser.Source)}
		}},

		{"Any -> Import", first},
		{"Any -> DeclsSection", first},
		{"Any -> ExportsSection", first},
		{"Any -> Function", first},
		{"Any -> export Function", func(args []any) any {
			source := args[1].(bplparser.Source)
			source.Function.Export = true
			return source
		}},
		{"Any -> Struct", first},
		{"Any -> export Struct", func(args []any) any {
			source := args[1].(bplparser.Source)
			source.TypeDef.Export = true
			return source
		}},
		{"Any -> Component", first},
		{"Any -> Term", func(args []any) any {
			return bplparser.NewTermSource(args[0].(ir.IrTerm))
		}},

		/* Import */

		{"Import -> import ID", func(args []any) any {
			return bplparser.NewImportSource(args[1].(string))
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

		{"Decls -> Decls Decl", func(args []any) any {
			return append(args[0].([]ir.IrDecl), args[1].(ir.IrDecl))
		}},
		{"Decls -> Decl", func(args []any) any {
			return []ir.IrDecl{args[0].(ir.IrDecl)}
		}},

		{"Decl -> Struct", func(args []any) any {
			return args[0].(bplparser.Source).TypeDef.Decl
		}},
		{"Decl -> TermDecl", first},
		{"Decl -> TypeDecl", first},

		{"TermDecl -> ID : SingleQuantifiedType", func(args []any) any {
			id := args[0].(string)
			typ := args[2].(ir.IrType)
			return ir.NewTermDecl(id, typ)
		}},

		{"TypeDecl -> type ID ;", func(args []any) any {
			id := args[1].(string)
			return ir.NewNameDecl(id)
		}},

		/* Function */

		{"Function -> func ID TypeAbstraction FunctionBindList -> FunctionBindList Block", func(args []any) any {
			id := args[1].(string)
			tvars := args[2].([]ir.VarKind)
			funArgs := args[3].([]ir.IrDecl)
			funRets := args[5].([]ir.IrDecl)
			body := args[6].(ir.IrTerm)
			return bplparser.NewFunctionSource(
				ir.NewFunction(false /* export */, id, tvars, funArgs, funRets, body))
		}},
		{"Function -> func ID FunctionBindList -> FunctionBindList Block", func(args []any) any {
			id := args[1].(string)
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

		{"Args -> Args , Arg", func(args []any) any {
			sources := args[0].([]ir.IrDecl)
			return append(sources, args[2].(ir.IrDecl))
		}},
		{"Args -> Arg", func(args []any) any {
			return []ir.IrDecl{args[0].(ir.IrDecl)}
		}},

		{"Arg -> ID UnquantifiedType", func(args []any) any {
			return ir.NewTermDecl(args[0].(string), args[1].(ir.IrType))
		}},

		/* Struct */

		{"Struct -> struct ID TypeAbstraction StructType", func(args []any) any {
			id := args[1].(string)
			tvars := args[2].([]ir.VarKind)
			structType := args[3].(ir.IrType)

			lambdaType := ir.LambdaVars(tvars, structType)
			return bplparser.NewTypeDefSource(false /* export */, ir.NewAliasDecl(id, lambdaType))
		}},
		{"Struct -> struct ID StructType", func(args []any) any {
			id := args[1].(string)
			var tvars []ir.VarKind
			structType := args[2].(ir.IrType)

			lambdaType := ir.LambdaVars(tvars, structType)
			return bplparser.NewTypeDefSource(false /* export */, ir.NewAliasDecl(id, lambdaType))
		}},

		/* Component */

		{"Component -> component [ UnquantifiedType , Num ]", func(args []any) any {
			elemType := args[2].(ir.IrType)
			length := args[4].(int)
			return bplparser.NewComponentSource(ir.NewComponent(elemType, length))
		}},

		/* Type variables */

		{"TypeAbstraction -> [ Tvars ]", func(args []any) any {
			return args[1].([]ir.VarKind)
		}},

		{"Tvars -> Tvars , Tvar", func(args []any) any {
			values := args[0].([]ir.VarKind)
			return append(values, args[2].(ir.VarKind))
		}},
		{"Tvars -> Tvar", func(args []any) any {
			return []ir.VarKind{args[0].(ir.VarKind)}
		}},

		{"Tvar -> ' ID", func(args []any) any {
			return ir.VarKind{args[1].(string), ir.NewTypeKind()}
		}},

		/* Single quantified type */

		{"SingleQuantifiedType -> QuantifiedType ;", first},

		/* Type */

		{"QuantifiedType -> UnquantifiedType", func(args []any) any {
			typ := args[0].(ir.IrType)
			return ir.QuantifyType(typ)
		}},

		/* Type */

		{"UnquantifiedType -> ForallType", first},

		/* Forall type */

		{"ForallType -> forall TypeAbstraction FunctionType", func(args []any) any {
			tvars := args[1].([]ir.VarKind)
			subType := args[2].(ir.IrType)
			return ir.ForallVars(tvars, subType)
		}},
		{"ForallType -> FunctionType", first},

		/* Function type */

		{"FunctionType -> AppType -> FunctionType", func(args []any) any {
			arg := args[0].(ir.IrType)
			ret := args[2].(ir.IrType)
			return ir.NewFunctionType(arg, ret)
		}},
		{"FunctionType -> AppType", first},

		/* App type */

		{"AppType -> AppType PrimaryType", func(args []any) any {
			argType := args[0].(ir.IrType)
			retType := args[1].(ir.IrType)
			return ir.NewAppType(argType, retType)
		}},
		{"AppType -> PrimaryType", func(args []any) any {
			return args[0].(ir.IrType)
		}},

		/* Simple Type */

		{"PrimaryType -> ArrayType", first},
		{"PrimaryType -> StructType", first},
		{"PrimaryType -> TupleType", first},
		{"PrimaryType -> ' ID", func(args []any) any {
			return ir.NewVarType(args[1].(string))
		}},
		{"PrimaryType -> ID", func(args []any) any {
			return ir.NewNameType(args[0].(string))
		}},
		{"PrimaryType -> ( UnquantifiedType )", func(args []any) any {
			return args[1].(ir.IrType)
		}},

		/* Array type */

		{"ArrayType -> [ UnquantifiedType , Num ]", func(args []any) any {
			typ := args[1].(ir.IrType)
			length := args[3].(int)
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

		{"TupleTypeArgs -> TupleTypeArgs , UnquantifiedType", func(args []any) any {
			return append(args[0].([]ir.IrType), args[2].(ir.IrType))
		}},
		{"TupleTypeArgs -> UnquantifiedType , UnquantifiedType", func(args []any) any {
			return []ir.IrType{args[0].(ir.IrType), args[2].(ir.IrType)}
		}},

		/* Struct type */

		{"StructType -> { }", func(args []any) any {
			return ir.NewStructType(nil)
		}},
		{"StructType -> { Fields }", func(args []any) any {
			values := args[1].([]ir.StructField)
			return ir.NewStructType(values)
		}},

		{"Fields -> Fields , Field", func(args []any) any {
			values := args[0].([]ir.StructField)
			return append(values, args[2].(ir.StructField))
		}},
		{"Fields -> Field", func(args []any) any {
			return []ir.StructField{args[0].(ir.StructField)}
		}},

		{"Field -> ID UnquantifiedType", func(args []any) any {
			return ir.StructField{args[0].(string), args[1].(ir.IrType)}
		}},

		/* ID */

		{"ID -> Token", func(args []any) any {
			return args[0].(Token).Token.Text
		}},

		{"Num -> Token", func(args []any) any {
			value, err := strconv.Atoi(args[0].(Token).Token.Text)
			if err != nil {
				panic(err)
			}

			return value
		}},

		/* Term */

		{"Block -> { Terms }", func(args []any) any {
			return ir.NewBlockTerm(args[1].([]ir.IrTerm))
		}},
		{"Block -> { }", func(args []any) any {
			return ir.NewBlockTerm(nil)
		}},

		{"Terms -> Terms Term", func(args []any) any {
			return append(args[0].([]ir.IrTerm), args[1].(ir.IrTerm))
		}},
		{"Terms -> Term", func(args []any) any {
			return []ir.IrTerm{args[0].(ir.IrTerm)}
		}},

		{"Term -> IfTerm", first},
		{"Term -> StatementTerm", first},

		/* Statement */

		{"StatementTerm -> AssignTerm", first},
		{"StatementTerm -> LetTerm", first},
		{"StatementTerm -> SingleExpression", first},

		/* Assign term */

		{"AssignTerm -> Token <- SingleExpression", func(args []any) any {
			ret := args[0].(Token)
			term := args[2].(ir.IrTerm)
			return ir.NewAssignTerm(term, ir.NewTokenTerm(ret.Token))
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
			id := args[1].(string)
			typ := args[2].(ir.IrType)
			var arg *ir.IrTerm
			return ir.NewLetTerm(ir.NewTermDecl(id, typ), arg)
		}},
		{"LetTerm -> let ID QuantifiedType = SingleExpression", func(args []any) any {
			id := args[1].(string)
			typ := args[2].(ir.IrType)
			arg := args[4].(ir.IrTerm)
			return ir.NewLetTerm(ir.NewTermDecl(id, typ), &arg)
		}},

		/* Single expression */

		{"SingleExpression -> Expression ;", func(args []any) any {
			return args[0].(ir.IrTerm)
		}},

		/* Expression */

		{"Expression -> equality", first},

		{"equality -> equality != comparison", func(args []any) any {
			return newBinOpTerm("!=", args[0].(ir.IrTerm), args[2].(ir.IrTerm))
		}},
		{"equality -> equality == comparison", func(args []any) any {
			return newBinOpTerm("==", args[0].(ir.IrTerm), args[2].(ir.IrTerm))
		}},
		{"equality -> comparison", first},

		{"comparison -> comparison > additive", func(args []any) any {
			return newBinOpTerm(">", args[0].(ir.IrTerm), args[2].(ir.IrTerm))
		}},
		{"comparison -> comparison >= additive", func(args []any) any {
			return newBinOpTerm(">=", args[0].(ir.IrTerm), args[2].(ir.IrTerm))
		}},
		{"comparison -> comparison < additive", func(args []any) any {
			return newBinOpTerm("<", args[0].(ir.IrTerm), args[2].(ir.IrTerm))
		}},
		{"comparison -> comparison <= additive", func(args []any) any {
			return newBinOpTerm("<=", args[0].(ir.IrTerm), args[2].(ir.IrTerm))
		}},
		{"comparison -> additive", first},

		{"additive -> additive + multiplicative", func(args []any) any {
			return newBinOpTerm("+", args[0].(ir.IrTerm), args[2].(ir.IrTerm))
		}},
		{"additive -> additive - multiplicative", func(args []any) any {
			return newBinOpTerm("-", args[0].(ir.IrTerm), args[2].(ir.IrTerm))
		}},
		{"additive -> multiplicative", first},

		{"multiplicative -> multiplicative * unary", func(args []any) any {
			return newBinOpTerm("*", args[0].(ir.IrTerm), args[2].(ir.IrTerm))
		}},
		{"multiplicative -> multiplicative / unary", func(args []any) any {
			return newBinOpTerm("/", args[0].(ir.IrTerm), args[2].(ir.IrTerm))
		}},
		{"multiplicative -> unary", first},

		{"unary -> ! unary", func(args []any) any { return newUnaryOpTerm("!", args[1].(ir.IrTerm)) }},
		{"unary -> - unary", func(args []any) any { return newUnaryOpTerm("-", args[1].(ir.IrTerm)) }},
		{"unary -> Applicative", first},

		{"Applicative -> Index.get TypeApplicative TypeApplicative", func(args []any) any {
			return ir.NewIndexGetTerm(args[1].(ir.IrTerm), args[2].(ir.IrTerm))
		}},
		{"Applicative -> Index.set TypeApplicative TypeApplicative TypeApplicative", func(args []any) any {
			return ir.NewIndexSetTerm(args[1].(ir.IrTerm), args[2].(ir.IrTerm), args[3].(ir.IrTerm))
		}},
		{"Applicative -> Applicative TypeApplicative", func(args []any) any {
			return ir.NewAppTermTerm(args[0].(ir.IrTerm), args[1].(ir.IrTerm))
		}},
		{"Applicative -> TypeApplicative", first},

		/* Type applicative */

		{"TypeApplicative -> Primary TypeApplicativeArgs", func(args []any) any {
			term := args[0].(ir.IrTerm)
			for _, typ := range args[1].([]ir.IrType) {
				term = ir.NewAppTypeTerm(term, typ)
			}
			return term
		}},
		{"TypeApplicative -> Primary", first},

		{"TypeApplicativeArgs -> [ TupleTypeArgs ]", func(args []any) any {
			return args[1].([]ir.IrType)
		}},
		{"TypeApplicativeArgs -> [ UnquantifiedType ]", func(args []any) any {
			return []ir.IrType{args[1].(ir.IrType)}
		}},

		/* Primary */

		{"Primary -> TupleTerm", first},
		{"Primary -> Token", func(args []any) any {
			return ir.NewTokenTerm(args[0].(Token).Token)
		}},
		{"Primary -> ( Expression )", func(args []any) any { return args[1] }},

		/* Tuple term */

		{"TupleTerm -> ( )", func(args []any) any {
			return ir.NewTupleTerm(nil)
		}},
		{"TupleTerm -> ( TupleTermArgs )", func(args []any) any {
			return ir.NewTupleTerm(args[1].([]ir.IrTerm))
		}},

		{"TupleTermArgs -> TupleTermArgs , Expression", func(args []any) any {
			return append(args[0].([]ir.IrTerm), args[2].(ir.IrTerm))
		}},
		{"TupleTermArgs -> Expression , Expression", func(args []any) any {
			return []ir.IrTerm{args[0].(ir.IrTerm), args[2].(ir.IrTerm)}
		}},
	}
}
