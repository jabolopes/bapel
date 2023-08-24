package ir

import "github.com/jabolopes/bapel/parser"

type IrTermCase int

const (
	CallTerm = IrTermCase(iota)
	TokenTerm
	TupleTerm
)

type IrTerm struct {
	Case IrTermCase
	Call *struct {
		ID   string
		Args []IrTerm
	}
	Token *parser.Token
	Tuple []IrTerm
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
