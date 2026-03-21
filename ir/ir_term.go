package ir

import (
	"fmt"

	"golang.org/x/exp/slices"
)

type IrTermCase int

const (
	AppTermTerm IrTermCase = iota
	AppTypeTerm
	AssignTerm
	BlockTerm
	// Constant term, e.g., number, string, etc.
	ConstTerm
	InjectionTerm
	LambdaTerm
	LetTerm
	MatchTerm
	ProjectionTerm
	ReturnTerm
	SetTerm
	StructTerm
	TupleTerm
	TypeAbsTerm
	// Variable term, e.g., identifier.
	VarTerm
)

func (c IrTermCase) String() string {
	switch c {
	case AppTermTerm:
		return "apply term to term"
	case AppTypeTerm:
		return "apply type to term"
	case AssignTerm:
		return "assign"
	case BlockTerm:
		return "block"
	case ConstTerm:
		return "constant"
	case InjectionTerm:
		return "injection"
	case LambdaTerm:
		return "lambda"
	case LetTerm:
		return "let"
	case MatchTerm:
		return "match"
	case ProjectionTerm:
		return "projection"
	case ReturnTerm:
		return "return"
	case SetTerm:
		return "set"
	case StructTerm:
		return "struct"
	case TupleTerm:
		return "tuple"
	case TypeAbsTerm:
		return "type abstraction"
	case VarTerm:
		return "variable"
	default:
		panic(fmt.Errorf("unhandled %T %d", c, c))
	}
}

// Apply a term to a term.
//
// foo x
type appTermTerm struct {
	Fun IrTerm
	Arg IrTerm
}

func (t *appTermTerm) Format(f fmt.State, verb rune) {
	funParenL := ""
	funParenR := ""
	if t.Fun.Is(LambdaTerm) {
		funParenL = "("
		funParenR = ")"
	}

	argParenL := ""
	argParenR := ""
	switch t.Arg.Case {
	case AppTermTerm, LambdaTerm:
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
type appTypeTerm struct {
	Fun IrTerm
	Arg IrType
}

func (t *appTypeTerm) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "%s [%s]", t.Fun, t.Arg)
}

type assignTerm struct {
	Arg IrTerm
	Ret IrTerm
}

func (t *assignTerm) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "%s <- %s", t.Ret, t.Arg)
}

type blockTerm struct {
	Terms []IrTerm
}

func (t *blockTerm) Format(f fmt.State, verb rune) {
	i := NewIndent(f)

	if p, ok := f.Precision(); ok && p == 0 {
		fmt.Fprintln(f, "{")
	} else {
		i.Println("{")
	}

	i.Inc()
	for _, term := range t.Terms {
		i.Printf(fmt.FormatString(f, verb), term)
		i.Println()
	}
	i.Dec().Print("}")
}

type constTerm struct {
	IrLiteral
}

type injectionTerm struct {
	VariantType IrType
	Tag         string
	Value       IrTerm
	// Determines the index of the variant tag to generate C++ code
	// using std::in_place_index.
	TagIndex *int
}

func (t *injectionTerm) Format(f fmt.State, verb rune) {
	typeNeedsParens := false
	switch t.VariantType.Case {
	case AppType, ForallType, FunType, LambdaType:
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

	fmt.Fprintf(f, "variant{%s%s%s %s = %s}", lparen, t.VariantType, rparen, t.Tag, t.Value)
}

// \ $arg $type = $body
type lambdaTerm struct {
	Arg     string
	ArgType IrType
	Body    IrTerm
}

func (t *lambdaTerm) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, `\(%s: %s) -> %s`, t.Arg, t.ArgType, t.Body)
}

// let $var : $type = $value
type letTerm struct {
	Var     string
	VarType *IrType
	Value   IrTerm
}

func (t *letTerm) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "let %s: %s = %s", t.Var, t.VarType, t.Value)
}

type MatchArm struct {
	Tag  string
	Arg  string
	Body IrTerm
	// The index of the tag (if any). Set by the typechecker.
	Index *int
}

func (t MatchArm) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "%s %s -> %s", t.Tag, t.Arg, t.Body)
}

type matchTerm struct {
	Term IrTerm
	Arms []MatchArm
}

func (t *matchTerm) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "case %s {", t.Term)
	switch len(t.Arms) {
	case 1:
		fmt.Fprintf(f, " %s }", t.Arms[0])

	default:
		fmt.Fprint(f, "\n")
		// TODO: Avoid the double Inc() hack.
		i := NewIndent(f).Inc().Inc()
		Interleave(t.Arms, func() { i.Println() }, func(_ int, arm MatchArm) {
			i.Printf("%s", arm)
		})
		i.Dec().Print("}")
	}
}

type projectionTerm struct {
	Term IrTerm
	// Either an integer (index-based projection) or an identifier (label-based projection).
	Label string
}

func (t *projectionTerm) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "%s.%s", t.Term, t.Label)
}

type returnTerm struct {
	Expr IrTerm
}

func (t returnTerm) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "return %s", t.Expr)
}

/* Set term */

type setTerm struct {
	Term   IrTerm
	Values []LabelValue
}

func (t setTerm) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "set %s {", t.Term)
	Interleave(t.Values, func() { fmt.Fprint(f, ", ") }, func(_ int, lv LabelValue) {
		fmt.Fprintf(f, "%s", lv)
	})
	fmt.Fprint(f, "}")
}

/* Struct term */

type LabelValue struct {
	Label string
	Value IrTerm
}

func (t LabelValue) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "%s = %s", t.Label, t.Value)
}

type structTerm struct {
	Values []LabelValue
}

func (t structTerm) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "struct{")
	Interleave(t.Values, func() { fmt.Fprint(f, ", ") }, func(_ int, lv LabelValue) {
		fmt.Fprintf(f, "%s", lv)
	})
	fmt.Fprintf(f, "}")
}

/* Tuple term */

type tupleTerm struct {
	Elems []IrTerm
}

func (t *tupleTerm) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "(")
	Interleave(t.Elems, func() { fmt.Fprintf(f, ", ") }, func(_ int, term IrTerm) {
		fmt.Fprintf(f, "%s", term)
	})
	fmt.Fprintf(f, ")")
}

/* Type abstraction term */

type typeAbsTerm struct {
	TypeVar string
	Kind    IrKind
	Body    IrTerm
}

func (t *typeAbsTerm) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "Λ%s :: %s. %s", t.TypeVar, t.Kind, t.Body)
}

/* Variable term */

type varTerm struct {
	ID string
}

func (t *varTerm) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "%s", t.ID)
}

type IrTerm struct {
	Case       IrTermCase
	AppTerm    *appTermTerm
	AppType    *appTypeTerm
	Assign     *assignTerm
	Block      *blockTerm
	Const      *constTerm
	Injection  *injectionTerm
	Lambda     *lambdaTerm
	Let        *letTerm
	Match      *matchTerm
	Projection *projectionTerm
	Return     *returnTerm
	Set        *setTerm
	Struct     *structTerm
	Tuple      *tupleTerm
	TypeAbs    *typeAbsTerm
	Var        *varTerm

	// Position in source file.
	Pos Pos
	// Type of this term. Set by the typechecker.
	Type *IrType
}

func (t IrTerm) formatImpl(f fmt.State, verb rune) {
	if t.Case == 0 && t.AppTerm == nil {
		return
	}

	switch t.Case {
	case AppTermTerm:
		t.AppTerm.Format(f, verb)
	case AppTypeTerm:
		t.AppType.Format(f, verb)
	case AssignTerm:
		t.Assign.Format(f, verb)
	case BlockTerm:
		t.Block.Format(f, verb)
	case ConstTerm:
		t.Const.Format(f, verb)
	case InjectionTerm:
		t.Injection.Format(f, verb)
	case LambdaTerm:
		t.Lambda.Format(f, verb)
	case LetTerm:
		t.Let.Format(f, verb)
	case MatchTerm:
		t.Match.Format(f, verb)
	case ProjectionTerm:
		t.Projection.Format(f, verb)
	case ReturnTerm:
		t.Return.Format(f, verb)
	case SetTerm:
		t.Set.Format(f, verb)
	case StructTerm:
		t.Struct.Format(f, verb)
	case TupleTerm:
		t.Tuple.Format(f, verb)
	case TypeAbsTerm:
		t.TypeAbs.Format(f, verb)
	case VarTerm:
		t.Var.Format(f, verb)
	default:
		panic(fmt.Errorf("unhandled IrTermCase %d", t.Case))
	}
}

func (t IrTerm) Format(f fmt.State, verb rune) {
	if t.Type == nil {
		t.formatImpl(f, verb)
		return
	}

	termNeedsParens := false
	switch t.Case {
	case AppTermTerm, AppTypeTerm, AssignTerm, InjectionTerm, LambdaTerm,
		LetTerm, MatchTerm, ProjectionTerm, ReturnTerm, TypeAbsTerm:
		termNeedsParens = true
	}

	typeNeedsParens := false
	switch t.Type.Case {
	case AppType, ForallType, FunType, LambdaType:
		typeNeedsParens = true
	}

	if termNeedsParens {
		fmt.Fprintf(f, "(")
	}
	t.formatImpl(f, verb)
	if termNeedsParens {
		fmt.Fprintf(f, ")")
	}

	fmt.Fprintf(f, ":")

	if typeNeedsParens {
		fmt.Fprintf(f, "(")
	}
	fmt.Fprintf(f, "%s", t.Type)
	if typeNeedsParens {
		fmt.Fprintf(f, ")")
	}
}

func (t IrTerm) Is(c IrTermCase) bool {
	return t.Case == c
}

func (t IrTerm) AppTypes() (IrTerm, []IrType) {
	if !t.Is(AppTypeTerm) {
		return t, nil
	}

	var types []IrType
	for t.Is(AppTypeTerm) {
		types = append(types, t.AppType.Arg)
		t = t.AppType.Fun
	}
	slices.Reverse(types)

	return t, types
}

func (t IrTerm) AppArgs() (IrTerm, []IrType, IrTerm) {
	var arg IrTerm
	if !t.Is(AppTermTerm) {
		return IrTerm{}, nil, IrTerm{}
	}
	arg = t.AppTerm.Arg
	t = t.AppTerm.Fun

	var types []IrType
	for t.Is(AppTypeTerm) {
		types = append(types, t.AppType.Arg)
		t = t.AppType.Fun
	}
	slices.Reverse(types)

	return t, types, arg
}

func (t IrTerm) ToFunction() ([]string, []string, []IrType, IrTerm) {
	// Type variables from the type abstraction term, e.g., 'a' in '/\ a :: k. t'.
	var typeVars []string
	// Variables and their types from the abstraction term, e.g., 'x' and 'a' in '\ x : a. t'.
	var args []string
	var argTypes []IrType

	term := t

	for {
		if !term.Is(TypeAbsTerm) {
			break
		}

		typeVars = append(typeVars, term.TypeAbs.TypeVar)

		term = term.TypeAbs.Body
	}

	for {
		if !term.Is(LambdaTerm) {
			break
		}

		args = append(args, term.Lambda.Arg)
		argTypes = append(argTypes, term.Lambda.ArgType)

		term = term.Lambda.Body
	}

	return typeVars, args, argTypes, term
}

// StructType returns the type of a StructTerm (if any).
func (t IrTerm) StructType() (IrType, bool) {
	if !t.Is(StructTerm) {
		return IrType{}, false
	}

	fields := make([]StructField, 0, len(t.Struct.Values))
	for _, value := range t.Struct.Values {
		if value.Value.Type == nil {
			return IrType{}, false
		}
		fields = append(fields, StructField{value.Label, *value.Value.Type})
	}

	return NewStructType(fields), true
}

// TupleType returns the type of a TupleTerm (if any).
func (t IrTerm) TupleType() (IrType, bool) {
	if !t.Is(TupleTerm) {
		return IrType{}, false
	}

	elems := make([]IrType, 0, len(t.Tuple.Elems))
	for _, elem := range t.Tuple.Elems {
		if elem.Type == nil {
			return IrType{}, false
		}
		elems = append(elems, *elem.Type)
	}

	return NewTupleType(elems), true
}

func NewAppTermTerm(fun, arg IrTerm) IrTerm {
	return IrTerm{
		Case:    AppTermTerm,
		AppTerm: &appTermTerm{fun, arg},
	}
}

func NewAppTypeTerm(fun IrTerm, arg IrType) IrTerm {
	return IrTerm{
		Case:    AppTypeTerm,
		AppType: &appTypeTerm{fun, arg},
	}
}

func NewAssignTerm(arg, ret IrTerm) IrTerm {
	if ret.Is(TupleTerm) && len(ret.Tuple.Elems) == 0 {
		return arg
	}

	return IrTerm{
		Case:   AssignTerm,
		Assign: &assignTerm{arg, ret},
	}
}

func NewBlockTerm(terms []IrTerm) IrTerm {
	return IrTerm{
		Case:  BlockTerm,
		Block: &blockTerm{terms},
	}
}

func NewConstTerm(literal IrLiteral) IrTerm {
	return IrTerm{
		Case:  ConstTerm,
		Const: &constTerm{literal},
		Pos:   literal.Pos,
	}
}

func NewInjectionTerm(variantType IrType, tag string, value IrTerm) IrTerm {
	return IrTerm{
		Case:      InjectionTerm,
		Injection: &injectionTerm{variantType, tag, value, nil /* TagIndex */},
	}
}

func NewLambdaTerm(arg string, argType IrType, body IrTerm) IrTerm {
	return IrTerm{
		Case:   LambdaTerm,
		Lambda: &lambdaTerm{arg, argType, body},
	}
}

func NewLetTerm(varName string, varType *IrType, value IrTerm) IrTerm {
	return IrTerm{
		Case: LetTerm,
		Let:  &letTerm{varName, varType, value},
	}
}

func NewMatchArm(tag, arg string, body IrTerm) MatchArm {
	return MatchArm{tag, arg, body, nil /* Index */}
}

func NewMatchTerm(term IrTerm, arms []MatchArm) IrTerm {
	return IrTerm{
		Case:  MatchTerm,
		Match: &matchTerm{term, arms},
	}
}

func NewProjectionTerm(term IrTerm, label string) IrTerm {
	return IrTerm{
		Case:       ProjectionTerm,
		Projection: &projectionTerm{term, label},
	}
}

func NewReturnTerm(expr IrTerm) IrTerm {
	return IrTerm{
		Case:   ReturnTerm,
		Return: &returnTerm{expr},
	}
}

func NewSetTerm(term IrTerm, values []LabelValue) IrTerm {
	return IrTerm{
		Case: SetTerm,
		Set:  &setTerm{term, values},
	}
}

func NewStructTerm(values []LabelValue) IrTerm {
	return IrTerm{
		Case:   StructTerm,
		Struct: &structTerm{values},
	}
}

func NewTupleTerm(elems []IrTerm) IrTerm {
	if len(elems) == 1 {
		return elems[0]
	}

	return IrTerm{
		Case:  TupleTerm,
		Tuple: &tupleTerm{elems},
	}
}

func NewTypeAbsTerm(tvar string, kind IrKind, term IrTerm) IrTerm {
	return IrTerm{
		Case:    TypeAbsTerm,
		TypeAbs: &typeAbsTerm{tvar, kind, term},
	}
}

func NewVarTerm(id string) IrTerm {
	return IrTerm{
		Case: VarTerm,
		Var:  &varTerm{id},
	}
}
