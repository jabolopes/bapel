package parser

import (
	"fmt"
)

type TokenCase int

const (
	IDToken TokenCase = iota
	NumberToken
)

func (c TokenCase) String() string {
	switch c {
	case IDToken:
		return "identifier"
	case NumberToken:
		return "number"
	default:
		panic(fmt.Errorf("unhandled TokenCase %d", c))
	}
}

type Token struct {
	Case  TokenCase
	Text  string
	Value int64
}

func (t Token) String() string {
	return t.Text
}

func (t Token) Is(c TokenCase) bool {
	return t.Case == c
}

func NewIDToken(text string) Token {
	return Token{IDToken, text, 0}
}

func NewNumberToken(value int64) Token {
	return Token{NumberToken, fmt.Sprintf("%d", value), value}
}
