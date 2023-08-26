package ir

import "github.com/jabolopes/bapel/parser"

type IrTermCase int

const (
	AssignTerm = IrTermCase(iota)
	CallTerm
	IfTerm
	IndexGetTerm
	IndexSetTerm
	OpUnaryTerm
	OpBinaryTerm
	StatementTerm
	TokenTerm
	TupleTerm
	WidenTerm
)

type IrTerm struct {
	Case   IrTermCase
	Assign *struct {
		Arg IrTerm
		Ret IrTerm
	}
	Call *struct {
		ID   string
		Args []IrTerm
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
	OpBinary *struct {
		ID    string
		Left  IrTerm
		Right IrTerm
	}
	Statement *struct{ Expr IrTerm }
	Token     *parser.Token
	Tuple     []IrTerm
	Widen     *struct{ Term IrTerm }
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

func NewCallTerm(id string, args []IrTerm) IrTerm {
	term := IrTerm{}
	term.Case = CallTerm
	term.Call = &struct {
		ID   string
		Args []IrTerm
	}{id, args}
	return term
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

func NewOpUnaryTerm(id string, term IrTerm) IrTerm {
	t := IrTerm{}
	t.Case = OpUnaryTerm
	t.OpUnary = &struct {
		ID   string
		Term IrTerm
	}{id, term}
	return t
}

func NewOpBinaryTerm(id string, left, right IrTerm) IrTerm {
	t := IrTerm{}
	t.Case = OpBinaryTerm
	t.OpBinary = &struct {
		ID    string
		Left  IrTerm
		Right IrTerm
	}{id, left, right}
	return t
}

func NewStatementTerm(expr IrTerm) IrTerm {
	term := IrTerm{}
	term.Case = StatementTerm
	term.Statement = &struct{ Expr IrTerm }{expr}
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
