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
	IfTerm
	IndexGetTerm
	IndexSetTerm
	LetTerm
	LiteralTerm
	TupleTerm
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
	case IfTerm:
		return "if"
	case IndexGetTerm:
		return "index get"
	case IndexSetTerm:
		return "index set"
	case LetTerm:
		return "let"
	case LiteralTerm:
		return "literal"
	case TupleTerm:
		return "tuple"
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
	case 2:
		return fmt.Sprintf("{ %s %s }", t.Terms[0], t.Terms[1])
	case 3:
		return fmt.Sprintf("{ %s %s %s }", t.Terms[0], t.Terms[1], t.Terms[2])
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

type ifTerm struct {
	Types     []IrType // Parametric polymorphism type arguments.
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

type indexGetTerm struct {
	Obj   IrTerm
	Index IrTerm
	// Determines whether to generate C++ code using array notation ([]) or
	// field notation (.). If Field is set, this uses field notation and this
	// contains the name of the field to index. Set by the typechecker.
	Field string
}

type indexSetTerm struct {
	Obj   IrTerm
	Index IrTerm
	Value IrTerm
	// Determines whether to generate C++ code using array notation ([]) or
	// field notation (.). If Field is set, this uses field notation and this
	// contains the name of the field to index. Set by the typechecker.
	Field string
}

// let $decl = $arg
type letTerm struct {
	Decl IrDecl
	Arg  *IrTerm
}

func (t *letTerm) String() string {
	if t.Arg == nil {
		return fmt.Sprintf("let %s", t.Decl)
	}
	return fmt.Sprintf("let %s = %s", t.Decl, *t.Arg)
}

/* Literal */

type LiteralCase int

const (
	IDLiteral LiteralCase = iota
	NumberLiteral
)

func (c LiteralCase) String() string {
	switch c {
	case IDLiteral:
		return "identifier"
	case NumberLiteral:
		return "number"
	default:
		panic(fmt.Errorf("unhandled LiteralCase %d", c))
	}
}

type literalTerm struct {
	Case   LiteralCase
	Text   string
	Number int64
}

func (t *literalTerm) String() string {
	return t.Text
}

func (t *literalTerm) Is(c LiteralCase) bool {
	return t.Case == c
}

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

type IrTerm struct {
	Case     IrTermCase
	AppTerm  *appTermTerm
	AppType  *appTypeTerm
	Assign   *assignTerm
	Block    *blockTerm
	If       *ifTerm
	IndexGet *indexGetTerm
	IndexSet *indexSetTerm
	Let      *letTerm
	Literal  *literalTerm
	Tuple    *tupleTerm

	// Position in source file.
	Pos Pos
	// Type of this term. Set by the typechecker.
	Type *IrType
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
	case IfTerm:
		return t.If.String()
	case IndexGetTerm:
		return fmt.Sprintf("Index.get %s %s", t.IndexGet.Obj, t.IndexGet.Index)
	case IndexSetTerm:
		return fmt.Sprintf("Index.set %s %s %s", t.IndexSet.Obj, t.IndexSet.Index, t.IndexSet.Value)
	case LetTerm:
		return t.Let.String()
	case LiteralTerm:
		return t.Literal.String()
	case TupleTerm:
		return t.Tuple.String()
	default:
		panic(fmt.Errorf("unhandled IrTermCase %d", t.Case))
	}
}

func (t IrTerm) String() string {
	if t.Type != nil {
		return fmt.Sprintf("(%s:%s)", t.stringImpl(), t.Type)
	}

	return t.stringImpl()
}

func (t IrTerm) Is(c IrTermCase) bool {
	return t.Case == c
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

func NewIfTerm(types []IrType, condition IrTerm, then IrTerm, elseTerm *IrTerm) IrTerm {
	return IrTerm{
		Case: IfTerm,
		If:   &ifTerm{types, condition, then, elseTerm},
	}
}

func NewIndexGetTerm(obj IrTerm, index IrTerm) IrTerm {
	return IrTerm{
		Case:     IndexGetTerm,
		IndexGet: &indexGetTerm{obj, index, ""},
	}
}

func NewIndexSetTerm(obj IrTerm, index IrTerm, value IrTerm) IrTerm {
	return IrTerm{
		Case:     IndexSetTerm,
		IndexSet: &indexSetTerm{obj, index, value, ""},
	}
}

func NewLetTerm(decl IrDecl, arg *IrTerm) IrTerm {
	return IrTerm{
		Case: LetTerm,
		Let:  &letTerm{decl, arg},
	}
}

func NewLiteralTerm(c LiteralCase, text string, value int64) IrTerm {
	return IrTerm{
		Case:    LiteralTerm,
		Literal: &literalTerm{c, text, value},
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
