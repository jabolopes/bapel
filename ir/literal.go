package ir

import (
	"fmt"
)

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

type Literal struct {
	Case   LiteralCase
	Text   string
	Number int64
}

func (t Literal) String() string {
	return t.Text
}

func (t Literal) Is(c LiteralCase) bool {
	return t.Case == c
}

func NewIDLiteral(text string) Literal {
	return Literal{
		Case: IDLiteral,
		Text: text,
	}
}

func NewNumberLiteral(text string, value int64) Literal {
	return Literal{
		Case:   NumberLiteral,
		Text:   text,
		Number: value,
	}
}

func NewDecimalLiteral(value int64) Literal {
	return Literal{
		Case:   NumberLiteral,
		Text:   fmt.Sprintf("%d", value),
		Number: value,
	}
}
