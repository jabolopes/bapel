package ir

import "github.com/jabolopes/bapel/parser"

type IrTermCase int

const (
	AssignTerm = IrTermCase(iota)
	CallTerm
	IfTerm
	TokenTerm
	TupleTerm
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
	Token *parser.Token
	Tuple []IrTerm
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
