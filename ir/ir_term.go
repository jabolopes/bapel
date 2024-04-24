package ir

import (
	"fmt"
	"strings"

	"github.com/jabolopes/bapel/parser"
)

type IrTermCase int

const (
	AssignTerm IrTermCase = iota
	BlockTerm
	CallTerm
	IfTerm
	IndexGetTerm
	IndexSetTerm
	LetTerm
	TokenTerm
	TupleTerm
)

func (c IrTermCase) String() string {
	switch c {
	case AssignTerm:
		return "assign"
	case BlockTerm:
		return "block"
	case CallTerm:
		return "call"
	case IfTerm:
		return "if"
	case IndexGetTerm:
		return "index get"
	case IndexSetTerm:
		return "index set"
	case LetTerm:
		return "let"
	case TokenTerm:
		return "token"
	case TupleTerm:
		return "tuple"
	default:
		panic(fmt.Errorf("unhandled IrTermCase %d", c))
	}
}

type ifTerm struct {
	Negate    bool
	Types     []IrType // Parametric polymorphism type arguments.
	Condition IrTerm
	Then      IrTerm
	Else      *IrTerm
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

type blockTerm struct {
	Terms []IrTerm
}

func (t blockTerm) String() string {
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

type letTerm struct {
	Decl IrDecl
}

type IrTerm struct {
	Case   IrTermCase
	Assign *struct {
		Arg IrTerm
		Ret IrTerm
	}
	Block *blockTerm
	Call  *struct {
		ID    string
		Types []IrType
		Arg   IrTerm
	}
	If       *ifTerm
	IndexGet *indexGetTerm
	IndexSet *indexSetTerm
	Let      *letTerm
	Token    *parser.Token
	Tuple    []IrTerm

	// Type of this term. Set by the typechecker.
	Type *IrType
}

func (t IrTerm) stringImpl() string {
	if t.Case == 0 && t.Assign == nil {
		return ""
	}

	switch t.Case {
	case AssignTerm:
		return fmt.Sprintf("%s <- %s", t.Assign.Ret, t.Assign.Arg)

	case BlockTerm:
		return t.Block.String()

	case CallTerm:
		return fmt.Sprintf("%s %v %s", t.Call.ID, t.Call.Types, t.Call.Arg)

	case IfTerm:
		c := t.If
		var b strings.Builder
		b.WriteString("if ")
		if c.Negate {
			b.WriteString("not ")
		}
		b.WriteString(c.Condition.String())
		b.WriteString(" then ")
		b.WriteString(c.Then.String())
		if c.Else != nil {
			b.WriteString(" else ")
			b.WriteString(c.Else.String())
		}
		return b.String()

	case IndexGetTerm:
		return fmt.Sprintf("Index.get %s %s", t.IndexGet.Obj, t.IndexGet.Index)
	case IndexSetTerm:
		return fmt.Sprintf("Index.set %s %s %s", t.IndexSet.Obj, t.IndexSet.Index, t.IndexSet.Value)
	case LetTerm:
		return fmt.Sprintf("let %s", t.Let.Decl)
	case TokenTerm:
		return t.Token.String()

	case TupleTerm:
		var b strings.Builder
		b.WriteString("(")
		if len(t.Tuple) > 0 {
			b.WriteString(t.Tuple[0].String())
			for _, term := range t.Tuple[1:] {
				b.WriteString(", ")
				b.WriteString(term.String())
			}
		}
		b.WriteString(")")
		return b.String()

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

func NewAssignTerm(arg, ret IrTerm) IrTerm {
	if ret.Case == TupleTerm && len(ret.Tuple) == 0 {
		return arg
	}

	term := IrTerm{}
	term.Case = AssignTerm
	term.Assign = &struct {
		Arg IrTerm
		Ret IrTerm
	}{arg, ret}
	return term
}

func NewBlockTerm(terms []IrTerm) IrTerm {
	return IrTerm{
		Case:  BlockTerm,
		Block: &blockTerm{terms},
	}
}

func NewCallTerm(id string, types []IrType, arg IrTerm) IrTerm {
	return IrTerm{
		Case: CallTerm,
		Call: &struct {
			ID    string
			Types []IrType
			Arg   IrTerm
		}{id, types, arg},
	}
}

func NewIfTerm(negate bool, types []IrType, condition IrTerm, then IrTerm, elseTerm *IrTerm) IrTerm {
	return IrTerm{
		Case: IfTerm,
		If:   &ifTerm{negate, types, condition, then, elseTerm},
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

func NewLetTerm(decl IrDecl) IrTerm {
	return IrTerm{
		Case: LetTerm,
		Let:  &letTerm{decl},
	}
}

func NewTokenTerm(token parser.Token) IrTerm {
	return IrTerm{
		Case:  TokenTerm,
		Token: &token,
	}
}

func NewTupleTerm(tuple []IrTerm) IrTerm {
	if len(tuple) == 1 {
		return tuple[0]
	}

	return IrTerm{
		Case:  TupleTerm,
		Tuple: tuple,
	}
}
