package ir

import (
	"fmt"
	"strings"

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
	IfTerm
	InjectionTerm
	IndexGetTerm
	IndexSetTerm
	LambdaTerm
	LetTerm
	ProjectionTerm
	ReturnTerm
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
	case IfTerm:
		return "if"
	case InjectionTerm:
		return "injection"
	case IndexGetTerm:
		return "index get"
	case IndexSetTerm:
		return "index set"
	case LambdaTerm:
		return "lambda"
	case LetTerm:
		return "let"
	case ProjectionTerm:
		return "projection"
	case ReturnTerm:
		return "return"
	case StructTerm:
		return "struct"
	case TupleTerm:
		return "tuple"
	case TypeAbsTerm:
		return "type abstraction"
	case VarTerm:
		return "variable"
	default:
		panic(fmt.Errorf("unhandled IrTermCase %d", c))
	}
}

// Apply a term to a term.
//
// foo x
type appTermTerm struct {
	Fun IrTerm
	Arg IrTerm
}

func (t *appTermTerm) String() string {
	return fmt.Sprintf("%s %s", t.Fun, t.Arg)
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

func (t *appTypeTerm) String() string {
	return fmt.Sprintf("%s [%s]", t.Fun, t.Arg)
}

type assignTerm struct {
	Arg IrTerm
	Ret IrTerm
}

func (t *assignTerm) String() string {
	return fmt.Sprintf("%s <- %s", t.Ret, t.Arg)
}

type blockTerm struct {
	Terms []IrTerm
}

func (t *blockTerm) String() string {
	switch len(t.Terms) {
	case 0:
		return "{}"
	case 1:
		return fmt.Sprintf("{ %s }", t.Terms[0])
	default:
		var b strings.Builder
		b.WriteString("{\n")
		for _, term := range t.Terms {
			b.WriteString("  ")
			b.WriteString(term.String())
			b.WriteString("\n")
		}
		b.WriteString("}")
		return b.String()
	}
}

type constTerm struct {
	Value  string
	Number int64
}

func (t *constTerm) String() string {
	return t.Value
}

type ifTerm struct {
	Condition IrTerm
	Then      IrTerm
	Else      *IrTerm
}

func (t *ifTerm) String() string {
	var b strings.Builder
	b.WriteString("if ")
	b.WriteString(t.Condition.String())
	b.WriteString(" then ")
	b.WriteString(t.Then.String())
	if t.Else != nil {
		b.WriteString(" else ")
		b.WriteString(t.Else.String())
	}
	return b.String()
}

type injectionTerm struct {
	VariantType IrType
	Tag         IrTerm
	Value       IrTerm
	// Determines the index of the variant tag to generate C++ code
	// using std::in_place_index.
	TagIndex *int
}

func (t *injectionTerm) String() string {
	var b strings.Builder
	b.WriteString("{|")
	b.WriteString(t.VariantType.String())
	b.WriteString(" ")
	b.WriteString(t.Tag.String())
	b.WriteString(" = ")
	b.WriteString(t.Value.String())
	b.WriteString("|}")
	return b.String()
}

type indexGetTerm struct {
	Obj   IrTerm
	Index IrTerm
	// Determines whether to generate C++ code using array notation ([]) or
	// field notation (.). If Field is set, this uses field notation and this
	// contains the name of the field to index. Set by the typechecker.
	Field string
	// Determines the index of the variant tag to generate C++ code
	// using std::in_place_index.
	TagIndex *int
}

type indexSetTerm struct {
	Obj   IrTerm
	Index IrTerm
	Value IrTerm
	// Determines whether to generate C++ code using array notation ([]) or
	// field notation (.). If Field is set, this uses field notation and this
	// contains the name of the field to index. Set by the typechecker.
	Field string
	// Determines the index of the variant tag to generate C++ code
	// using std::in_place_index.
	TagIndex *int
}

// \ $arg $type = $body
type lambdaTerm struct {
	Arg     string
	ArgType IrType
	Body    IrTerm
}

func (t *lambdaTerm) String() string {
	return fmt.Sprintf(`\(%s : %s) -> %s`, t.Arg, t.ArgType, t.Body)
}

// let $var : $type = $value
type letTerm struct {
	Var     string
	VarType IrType
	Value   IrTerm
}

func (t *letTerm) String() string {
	return fmt.Sprintf("let %s : %s = %s", t.Var, t.VarType, t.Value)
}

type projectionTerm struct {
	Term IrTerm
	// Either a ConstTerm (index-based projection) or a VarTerm (label-based projection).
	Label IrTerm
	// The index of the label (if any). Set by the typechecker.
	//
	// When the term has tuple type, the index is always defined and it
	// corresponds to the element index.
	//
	// When the term has struct or variant type, the index is always defined and
	// it corresponds to the struct field index or the variant tag index.
	//
	// When the term has array type, the index is only defined if the Label is a
	// number literal. Otherwise, this is nil.
	Index *int
	// The name of the label (if any). Set by the typechecker.
	//
	// When the term has struct type or variant type, this is always defined and
	// it corresponds to the struct field name or the variant tag name.
	//
	// When the term has array or tuple, this is nil.
	LabelName *string
}

func (t *projectionTerm) String() string {
	return fmt.Sprintf("%s->%s", t.Term, t.Label)
}

type returnTerm struct {
	Expr IrTerm
}

func (t returnTerm) String() string {
	return fmt.Sprintf("return %s", t.Expr)
}

/* Struct term */

type LabelValue struct {
	Label string
	Value IrTerm
}

func (t LabelValue) String() string {
	return fmt.Sprintf("%s = %s", t.Label, t.Value)
}

type structTerm struct {
	Values []LabelValue
}

func (t structTerm) String() string {
	var b strings.Builder
	b.WriteString("{")
	if len(t.Values) > 0 {
		b.WriteString(t.Values[0].String())
		for _, term := range t.Values[1:] {
			b.WriteString(", ")
			b.WriteString(term.String())
		}
	}
	b.WriteString("}")
	return b.String()
}

/* Tuple term */

type tupleTerm struct {
	Elems []IrTerm
}

func (t *tupleTerm) String() string {
	var b strings.Builder
	b.WriteString("(")
	if len(t.Elems) > 0 {
		b.WriteString(t.Elems[0].String())
		for _, term := range t.Elems[1:] {
			b.WriteString(", ")
			b.WriteString(term.String())
		}
	}
	b.WriteString(")")
	return b.String()
}

/* Type abstraction term */

type typeAbsTerm struct {
	TypeVar string
	Kind    IrKind
	Body    IrTerm
}

func (t *typeAbsTerm) String() string {
	return fmt.Sprintf("Λ%s :: %s. %s", t.TypeVar, t.Kind, t.Body)
}

/* Variable term */

type varTerm struct {
	ID string
}

func (t *varTerm) String() string {
	return t.ID
}

type IrTerm struct {
	Case       IrTermCase
	AppTerm    *appTermTerm
	AppType    *appTypeTerm
	Assign     *assignTerm
	Block      *blockTerm
	Const      *constTerm
	If         *ifTerm
	Injection  *injectionTerm
	IndexGet   *indexGetTerm
	IndexSet   *indexSetTerm
	Lambda     *lambdaTerm
	Let        *letTerm
	Projection *projectionTerm
	Return     *returnTerm
	Struct     *structTerm
	Tuple      *tupleTerm
	TypeAbs    *typeAbsTerm
	Var        *varTerm

	// Position in source file.
	Pos Pos
	// Type of this term. Set by the typechecker.
	Type *IrType
	// Whether this is the last term of a function, which returns the
	// expression to the caller. Set by the typechecker.
	LastTerm bool
}

func (t IrTerm) stringImpl() string {
	if t.Case == 0 && t.AppTerm == nil {
		return ""
	}

	switch t.Case {
	case AppTermTerm:
		return t.AppTerm.String()
	case AppTypeTerm:
		return t.AppType.String()
	case AssignTerm:
		return t.Assign.String()
	case BlockTerm:
		return t.Block.String()
	case ConstTerm:
		return t.Const.String()
	case IfTerm:
		return t.If.String()
	case InjectionTerm:
		return t.Injection.String()
	case IndexGetTerm:
		return fmt.Sprintf("Index.get %s %s", t.IndexGet.Obj, t.IndexGet.Index)
	case IndexSetTerm:
		return fmt.Sprintf("Index.set %s %s %s", t.IndexSet.Obj, t.IndexSet.Index, t.IndexSet.Value)
	case LambdaTerm:
		return t.Lambda.String()
	case LetTerm:
		return t.Let.String()
	case ProjectionTerm:
		return t.Projection.String()
	case ReturnTerm:
		return t.Return.String()
	case StructTerm:
		return t.Struct.String()
	case TupleTerm:
		return t.Tuple.String()
	case TypeAbsTerm:
		return t.TypeAbs.String()
	case VarTerm:
		return t.Var.String()
	default:
		panic(fmt.Errorf("unhandled IrTermCase %d", t.Case))
	}
}

func (t IrTerm) String() string {
	if t.Type == nil {
		return t.stringImpl()
	}

	termNeedsParens := false
	switch t.Case {
	case AppTermTerm, AppTypeTerm, AssignTerm, IfTerm, InjectionTerm,
		IndexGetTerm, IndexSetTerm, LambdaTerm, LetTerm, ProjectionTerm,
		ReturnTerm, TypeAbsTerm:
		termNeedsParens = true
	}

	typeNeedsParens := false
	switch t.Type.Case {
	case AppType, ForallType, FunType, LambdaType:
		typeNeedsParens = true
	}

	var b strings.Builder
	if termNeedsParens {
		b.WriteString("(")
	}
	b.WriteString(t.stringImpl())
	if termNeedsParens {
		b.WriteString(")")
	}

	b.WriteString(":")

	if typeNeedsParens {
		b.WriteString("(")
	}
	b.WriteString(t.Type.String())
	if typeNeedsParens {
		b.WriteString(")")
	}

	return b.String()
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

func NewConstTerm(value string, number int64) IrTerm {
	return IrTerm{
		Case:  ConstTerm,
		Const: &constTerm{value, number},
	}
}

func NewIfTerm(condition IrTerm, then IrTerm, elseTerm *IrTerm) IrTerm {
	return IrTerm{
		Case: IfTerm,
		If:   &ifTerm{condition, then, elseTerm},
	}
}

func NewInjectionTerm(variantType IrType, tag, value IrTerm) IrTerm {
	return IrTerm{
		Case:      InjectionTerm,
		Injection: &injectionTerm{variantType, tag, value, nil /* TagIndex */},
	}
}

func NewIndexGetTerm(obj IrTerm, index IrTerm) IrTerm {
	return IrTerm{
		Case:     IndexGetTerm,
		IndexGet: &indexGetTerm{obj, index, "", nil},
	}
}

func NewIndexSetTerm(obj IrTerm, index IrTerm, value IrTerm) IrTerm {
	return IrTerm{
		Case:     IndexSetTerm,
		IndexSet: &indexSetTerm{obj, index, value, "", nil},
	}
}

func NewLambdaTerm(arg string, argType IrType, body IrTerm) IrTerm {
	return IrTerm{
		Case:   LambdaTerm,
		Lambda: &lambdaTerm{arg, argType, body},
	}
}

func NewLetTerm(varName string, varType IrType, value IrTerm) IrTerm {
	return IrTerm{
		Case: LetTerm,
		Let:  &letTerm{varName, varType, value},
	}
}

func NewProjectionTerm(term IrTerm, label IrTerm) IrTerm {
	return IrTerm{
		Case:       ProjectionTerm,
		Projection: &projectionTerm{term, label, nil, nil},
	}
}

func NewReturnTerm(expr IrTerm) IrTerm {
	return IrTerm{
		Case:   ReturnTerm,
		Return: &returnTerm{expr},
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
