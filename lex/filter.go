package lex

import (
	"github.com/emirpasic/gods/v2/stacks"
	"github.com/emirpasic/gods/v2/stacks/arraystack"
	"github.com/jabolopes/bapel/lex/lexer"
)

type tokenReader interface {
	NextToken() (lexer.Token, bool)
}

// lineFilter converts newlines into semicolons within blocks.
type lineFilter struct {
	tokenReader
	// Whether the filter is enabled or not.
	//
	// Given that it's mainly the terms and types that are ambiguous due
	// to term application and type application, the filter remains
	// disabled until it encounters the first `fn` keyword. That means
	// the first part of a module file (e.g., `module`, `imports`, etc)
	// can assume that there are no newlines.
	enabled bool
	// Block ID generator
	idgen int
	// Tracks open blocks. The toplevel block has ID 0 but it's never on the stack
	// since we don't insert semicolons in between the toplevel elements.
	blocks stacks.Stack[int]
	// Block ID of the previous token seen by the filter. Comparing against the
	// previous block ID we determine whether the block has changed. Semicolons
	// are inserted only for expressions within the same block.
	previousBlockID int
}

func (f *lineFilter) NextToken() (lexer.Token, bool) {
	token, ok := f.tokenReader.NextToken()
	if !ok {
		return token, ok
	}

	if token.Value == "fn" {
		f.enabled = true
	}

	if !f.enabled {
		switch token.Value {
		case "\n":
			return f.NextToken()
		default:
			return token, ok
		}
	}

	switch token.Value {
	case "\n":
		_, ok := f.blocks.Peek()
		if !ok {
			return f.NextToken()
		}

		token.Value = ";"

	case "{":
		f.idgen++
		f.blocks.Push(f.idgen)

	case "}":
		f.blocks.Pop()
	}

	return token, true
}

func newLineFilter(reader tokenReader) *lineFilter {
	return &lineFilter{
		reader,
		false, /* enabled */
		0,     /* idgen */
		arraystack.New[int](),
		0, /* previousBlockID */
	}
}
