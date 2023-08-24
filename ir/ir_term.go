package ir

import "github.com/jabolopes/bapel/parser"

type IrTermCase int

const (
	CallTerm = IrTermCase(iota)
	TokenTerm
)

type IrTerm struct {
	Case IrTermCase
	Call *struct {
		ID   string
		Args []parser.Token
	}
	Token *parser.Token
}

func NewCallTerm(id string, args []parser.Token) IrTerm {
	term := IrTerm{}
	term.Case = CallTerm
	term.Call = &struct {
		ID   string
		Args []parser.Token
	}{id, args}
	return term
}

func NewTokenTerm(token parser.Token) IrTerm {
	term := IrTerm{}
	term.Case = TokenTerm
	term.Token = &token
	return term
}
