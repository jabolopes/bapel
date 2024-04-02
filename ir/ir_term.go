package ir

import (
	"fmt"
	"strings"

	"github.com/jabolopes/bapel/parser"
)

type IrTermCase int

const (
	AssignTerm IrTermCase = iota
	CallTerm
	IfTerm
	IndexGetTerm
	IndexSetTerm
	StatementTerm
	TokenTerm
	TupleTerm
	WidenTerm
)

func (c IrTermCase) String() string {
	switch c {
	case AssignTerm:
		return "assign"
	case CallTerm:
		return "call"
	case IfTerm:
		return "if"
	case IndexGetTerm:
		return "index get"
	case IndexSetTerm:
		return "index set"
	case StatementTerm:
		return "statement"
	case TokenTerm:
		return "token"
	case TupleTerm:
		return "tuple"
	case WidenTerm:
		return "widen"
	default:
		panic(fmt.Errorf("unhandled IrTermCase %d", c))
	}
}

type IrTerm struct {
	Case   IrTermCase
	Assign *struct {
		Arg IrTerm
		Ret IrTerm
	}
	Call *struct {
		ID    string
		Types []IrType
		Arg   IrTerm
	}
	If *struct {
		Then      bool
		Condition IrTerm
	}
	IndexGet *struct {
		Term  IrTerm
		Index IrTerm
		// Determines whether to generate C++ code using array notation ([]) or
		// field notation (.). If Field is set, this uses field notation and this
		// contains the name of the field to index. Set by the typechecker.
		Field string
	}
	IndexSet *struct {
		Ret   IrTerm
		Index IrTerm
		Arg   IrTerm
		// Determines whether to generate C++ code using array notation ([]) or
		// field notation (.). If Field is set, this uses field notation and this
		// contains the name of the field to index. Set by the typechecker.
		Field string
	}
	OpUnary *struct {
		ID   string
		Term IrTerm
	}
	Statement *struct{ Term IrTerm }
	Token     *parser.Token
	Tuple     []IrTerm
	Widen     *struct{ Term IrTerm }

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
	case CallTerm:
		return fmt.Sprintf("%s %s", t.Call.ID, t.Call.Arg)

	case IfTerm:
		if t.If.Then {
			return fmt.Sprintf("if %s", t.If.Condition)
		} else {
			return fmt.Sprintf("if !%s", t.If.Condition)
		}

	case IndexGetTerm:
		return fmt.Sprintf("Index.get %s %s", t.IndexGet.Term, t.IndexGet.Index)
	case IndexSetTerm:
		return fmt.Sprintf("Index.set %s %s %s", t.IndexSet.Ret, t.IndexSet.Index, t.IndexSet.Arg)
	case StatementTerm:
		return fmt.Sprintf("%s;", t.Statement.Term.String())
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

	case WidenTerm:
		return fmt.Sprintf("widen %s", t.Widen.Term)
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

func NewIfTerm(then bool, condition IrTerm) IrTerm {
	term := IrTerm{}
	term.Case = IfTerm
	term.If = &struct {
		Then      bool
		Condition IrTerm
	}{then, condition}
	return term
}

func NewIndexGetTerm(term IrTerm, index IrTerm) IrTerm {
	t := IrTerm{}
	t.Case = IndexGetTerm
	t.IndexGet = &struct {
		Term  IrTerm
		Index IrTerm
		Field string
	}{term, index, ""}
	return t
}

func NewIndexSetTerm(term IrTerm, index IrTerm, value IrTerm) IrTerm {
	t := IrTerm{}
	t.Case = IndexSetTerm
	t.IndexSet = &struct {
		Ret   IrTerm
		Index IrTerm
		Arg   IrTerm
		Field string
	}{term, index, value, ""}
	return t
}

func NewStatementTerm(expr IrTerm) IrTerm {
	term := IrTerm{}
	term.Case = StatementTerm
	term.Statement = &struct{ Term IrTerm }{expr}
	return term
}

func NewTokenTerm(token parser.Token) IrTerm {
	term := IrTerm{}
	term.Case = TokenTerm
	term.Token = &token
	return term
}

func NewTupleTerm(tuple []IrTerm) IrTerm {
	if len(tuple) == 1 {
		return tuple[0]
	}

	term := IrTerm{}
	term.Case = TupleTerm
	term.Tuple = tuple
	return term
}

func NewWidenTerm(widen IrTerm) IrTerm {
	term := IrTerm{}
	term.Case = WidenTerm
	term.Widen = &struct{ Term IrTerm }{widen}
	return term
}
