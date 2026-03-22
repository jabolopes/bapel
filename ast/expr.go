package ast

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

type ExprCase int

const (
	AppTermExpr ExprCase = iota
	AppTypeExpr
	AssignExpr
	BlockExpr
	// Constant term, e.g., number, string, etc.
	ConstExpr
	InjectionExpr
	LambdaExpr
	LetExpr
	MatchExpr
	ProjectionExpr
	ReturnExpr
	SetExpr
	StructExpr
	TupleExpr
	TypeAbsExpr
	// Variable term, e.g., identifier.
	VarExpr
)

// Apply a term to a term.
//
// foo x
type appTermExpr struct {
	Fun Expr
	Arg Expr
}

func (t *appTermExpr) Format(f fmt.State, verb rune) {
	funParenL := ""
	funParenR := ""
	if t.Fun.Is(LambdaExpr) {
		funParenL = "("
		funParenR = ")"
	}

	argParenL := ""
	argParenR := ""
	switch t.Arg.Case {
	case AppTermExpr, LambdaExpr:
		argParenL = "("
		argParenR = ")"
	}

	fmt.Fprintf(f, "%s%s%s %s%s%s", funParenL, t.Fun, funParenR, argParenL, t.Arg, argParenR)
}

// Apply a type to a term.
//
// foo [i8]
//
// C |- x : forall 'a. B
// -----------------
// C |- x A : B[A/'a]
type appTypeExpr struct {
	Fun Expr
	Arg ir.IrType
}

func (t *appTypeExpr) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "%s [%s]", t.Fun, t.Arg)
}

type assignExpr struct {
	Arg Expr
	Ret Expr
}

func (t *assignExpr) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "%s <- %s", t.Ret, t.Arg)
}

type blockExpr struct {
	Exprs []Expr
}

func (t *blockExpr) Format(f fmt.State, verb rune) {
	i := ir.NewIndent(f)

	if p, ok := f.Precision(); ok && p == 0 {
		fmt.Fprintln(f, "{")
	} else {
		i.Println("{")
	}

	i.Inc()
	for _, term := range t.Exprs {
		i.Printf(fmt.FormatString(f, verb), term)
		i.Println()
	}
	i.Dec().Print("}")
}

type constExpr struct {
	ir.IrLiteral
}

type injectionExpr struct {
	VariantType ir.IrType
	Tag         string
	Expr        Expr
	// Determines the index of the variant tag to generate C++ code
	// using std::in_place_index.
	TagIndex *int
}

func (t *injectionExpr) Format(f fmt.State, verb rune) {
	typeNeedsParens := false
	switch t.VariantType.Case {
	case ir.AppType, ir.ForallType, ir.FunType, ir.LambdaType:
		typeNeedsParens = true
	}

	lparen := ""
	if typeNeedsParens {
		lparen = "("
	}

	rparen := ""
	if typeNeedsParens {
		rparen = ")"
	}

	fmt.Fprintf(f, "variant{%s%s%s %s = %s}", lparen, t.VariantType, rparen, t.Tag, t.Expr)
}

// \ $arg $type = $body
type lambdaExpr struct {
	Arg  ir.FunctionArg
	Body Expr
}

func (t *lambdaExpr) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, `\(%s) -> %s`, t.Arg, t.Body)
}

// let $var : $type = $value
type letExpr struct {
	Var     string
	VarType *ir.IrType
	Expr    Expr
}

func (t *letExpr) Format(f fmt.State, verb rune) {
	if t.VarType == nil {
		fmt.Fprintf(f, "let %s = %s", t.Var, t.Expr)
	} else {
		fmt.Fprintf(f, "let %s: %s = %s", t.Var, t.VarType, t.Expr)
	}
}

type MatchArm struct {
	Tag  string
	Arg  string
	Body Expr
	// The index of the tag (if any). Set by the typechecker.
	Index *int
}

func (t MatchArm) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "%s %s -> %s", t.Tag, t.Arg, t.Body)
}

type matchExpr struct {
	Expr Expr
	Arms []MatchArm
}

func (t *matchExpr) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "case %s {", t.Expr)
	switch len(t.Arms) {
	case 1:
		fmt.Fprintf(f, " %s }", t.Arms[0])

	default:
		fmt.Fprint(f, "\n")
		// TODO: Avoid the double Inc() hack.
		i := ir.NewIndent(f).Inc().Inc()
		ir.Interleave(t.Arms, func() { i.Println() }, func(_ int, arm MatchArm) {
			i.Printf("%s", arm)
		})
		i.Dec().Print("}")
	}
}

type projectionExpr struct {
	Expr Expr
	// Either an integer (index-based projection) or an identifier (label-based projection).
	Label string
}

func (t *projectionExpr) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "%s.%s", t.Expr, t.Label)
}

type returnExpr struct {
	Expr Expr
}

func (t returnExpr) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "return %s", t.Expr)
}

/* Set term */

type setExpr struct {
	Expr   Expr
	Values []LabelValue
}

func (t setExpr) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "set %s {", t.Expr)
	ir.Interleave(t.Values, func() { fmt.Fprint(f, ", ") }, func(_ int, lv LabelValue) {
		fmt.Fprintf(f, "%s", lv)
	})
	fmt.Fprint(f, "}")
}

/* Struct term */

type LabelValue struct {
	Label string
	Value Expr
}

func (t LabelValue) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "%s = %s", t.Label, t.Value)
}

type structExpr struct {
	Values []LabelValue
}

func (t structExpr) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "struct{")
	ir.Interleave(t.Values, func() { fmt.Fprint(f, ", ") }, func(_ int, lv LabelValue) {
		fmt.Fprintf(f, "%s", lv)
	})
	fmt.Fprintf(f, "}")
}

/* Tuple term */

type tupleExpr struct {
	Elems []Expr
}

func (t *tupleExpr) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "(")
	ir.Interleave(t.Elems, func() { fmt.Fprintf(f, ", ") }, func(_ int, term Expr) {
		fmt.Fprintf(f, "%s", term)
	})
	fmt.Fprintf(f, ")")
}

/* Type abstraction term */

type typeAbsExpr struct {
	Arg  ir.VarKind
	Body Expr
}

func (t *typeAbsExpr) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "Λ%s. %s", t.Arg, t.Body)
}

/* Variable term */

type varExpr struct {
	ID string
}

func (t *varExpr) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "%s", t.ID)
}

type Expr struct {
	Case       ExprCase
	AppTerm    *appTermExpr
	AppType    *appTypeExpr
	Assign     *assignExpr
	Block      *blockExpr
	Const      *constExpr
	Injection  *injectionExpr
	Lambda     *lambdaExpr
	Let        *letExpr
	Match      *matchExpr
	Projection *projectionExpr
	Return     *returnExpr
	Set        *setExpr
	Struct     *structExpr
	Tuple      *tupleExpr
	TypeAbs    *typeAbsExpr
	Var        *varExpr

	Pos ir.Pos
}

func (t Expr) formatImpl(f fmt.State, verb rune) {
	if t.Case == 0 && t.AppTerm == nil {
		return
	}

	switch t.Case {
	case AppTermExpr:
		t.AppTerm.Format(f, verb)
	case AppTypeExpr:
		t.AppType.Format(f, verb)
	case AssignExpr:
		t.Assign.Format(f, verb)
	case BlockExpr:
		t.Block.Format(f, verb)
	case ConstExpr:
		t.Const.Format(f, verb)
	case InjectionExpr:
		t.Injection.Format(f, verb)
	case LambdaExpr:
		t.Lambda.Format(f, verb)
	case LetExpr:
		t.Let.Format(f, verb)
	case MatchExpr:
		t.Match.Format(f, verb)
	case ProjectionExpr:
		t.Projection.Format(f, verb)
	case ReturnExpr:
		t.Return.Format(f, verb)
	case SetExpr:
		t.Set.Format(f, verb)
	case StructExpr:
		t.Struct.Format(f, verb)
	case TupleExpr:
		t.Tuple.Format(f, verb)
	case TypeAbsExpr:
		t.TypeAbs.Format(f, verb)
	case VarExpr:
		t.Var.Format(f, verb)
	default:
		panic(fmt.Errorf("unhandled ExprCase %d", t.Case))
	}
}

func (t Expr) Format(f fmt.State, verb rune) {
	t.formatImpl(f, verb)
}

func (s Expr) Is(c ExprCase) bool {
	return s.Case == c
}

func NewAppTermExpr(pos ir.Pos, fun, arg Expr) Expr {
	return Expr{
		Case:    AppTermExpr,
		AppTerm: &appTermExpr{fun, arg},
		Pos:     pos,
	}
}

func NewAppTypeExpr(pos ir.Pos, fun Expr, arg ir.IrType) Expr {
	return Expr{
		Case:    AppTypeExpr,
		AppType: &appTypeExpr{fun, arg},
		Pos:     pos,
	}
}

func NewAssignExpr(pos ir.Pos, arg, ret Expr) Expr {
	if ret.Is(TupleExpr) && len(ret.Tuple.Elems) == 0 {
		return arg
	}

	return Expr{
		Case:   AssignExpr,
		Assign: &assignExpr{arg, ret},
		Pos:    pos,
	}
}

func NewBlockExpr(pos ir.Pos, sources []Expr) Expr {
	return Expr{
		Case:  BlockExpr,
		Block: &blockExpr{sources},
		Pos:   pos,
	}
}

func NewConstExpr(literal ir.IrLiteral) Expr {
	return Expr{
		Case:  ConstExpr,
		Const: &constExpr{literal},
		Pos:   literal.Pos,
	}
}

func NewInjectionExpr(pos ir.Pos, variantType ir.IrType, tag string, value Expr) Expr {
	return Expr{
		Case:      InjectionExpr,
		Injection: &injectionExpr{variantType, tag, value, nil /* TagIndex */},
		Pos:       pos,
	}
}

func NewLambdaExpr(pos ir.Pos, arg ir.FunctionArg, body Expr) Expr {
	return Expr{
		Case:   LambdaExpr,
		Lambda: &lambdaExpr{arg, body},
		Pos:    pos,
	}
}

func NewLetExpr(pos ir.Pos, varName string, varType *ir.IrType, value Expr) Expr {
	return Expr{
		Case: LetExpr,
		Let:  &letExpr{varName, varType, value},
		Pos:  pos,
	}
}

func NewMatchArm(tag, arg string, body Expr) MatchArm {
	return MatchArm{tag, arg, body, nil /* Index */}
}

func NewMatchExpr(pos ir.Pos, expr Expr, arms []MatchArm) Expr {
	return Expr{
		Case:  MatchExpr,
		Match: &matchExpr{expr, arms},
		Pos:   pos,
	}
}

func NewProjectionExpr(pos ir.Pos, expr Expr, label string) Expr {
	return Expr{
		Case:       ProjectionExpr,
		Projection: &projectionExpr{expr, label},
		Pos:        pos,
	}
}

func NewReturnExpr(pos ir.Pos, expr Expr) Expr {
	return Expr{
		Case:   ReturnExpr,
		Return: &returnExpr{expr},
		Pos:    pos,
	}
}

func NewSetExpr(pos ir.Pos, expr Expr, values []LabelValue) Expr {
	return Expr{
		Case: SetExpr,
		Set:  &setExpr{expr, values},
		Pos:  pos,
	}
}

func NewStructExpr(pos ir.Pos, values []LabelValue) Expr {
	return Expr{
		Case:   StructExpr,
		Struct: &structExpr{values},
		Pos:    pos,
	}
}

func NewTupleExpr(pos ir.Pos, elems []Expr) Expr {
	if len(elems) == 1 {
		return elems[0]
	}

	return Expr{
		Case:  TupleExpr,
		Tuple: &tupleExpr{elems},
		Pos:   pos,
	}
}

func NewTypeAbsExpr(pos ir.Pos, arg ir.VarKind, body Expr) Expr {
	return Expr{
		Case:    TypeAbsExpr,
		TypeAbs: &typeAbsExpr{arg, body},
		Pos:     pos,
	}
}

func NewVarExpr(id ID) Expr {
	return Expr{
		Case: VarExpr,
		Var:  &varExpr{id.Value},
		Pos:  id.Pos,
	}
}
