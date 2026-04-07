package lex

import (
	"fmt"
	"unicode"

	"github.com/jabolopes/bapel/lex/lexer"
)

const (
	WordToken lexer.TokenType = iota
	NumberToken
	RuneToken
	StringToken
	SymbolToken   = WordToken
	OperatorToken = WordToken
)

// The general rules below for symbols capture operators that are in
// the same Unicode group (e.g., symbol, punctuation, etc). Operators
// that are in different Unicode groups are listed here.
var operators = []string{
	"<-",
	"->",
	"||",
	"&&",
	"!=",
	"::",
}

type states struct {
	*lexer.LexerFSM
	outErrors []string
}

func (l *states) initialState() lexer.StateFunc {
	switch c := l.Peek(); c {
	case lexer.EOFRune:
		return nil

	case '"':
		return l.newStringState("string", '"')

	case '`':
		return l.newStringState("raw string", '`')

	default:
		if l.PeekAll("//") {
			return l.lineCommentState
		}

		if l.PeekAll("/*") {
			return l.blockCommentState
		}

		for _, operator := range operators {
			if l.TakeAll(operator) {
				l.Emit(OperatorToken)
				return l.initialState
			}
		}

		if unicode.IsSpace(c) {
			return l.whitespaceState
		}

		if unicode.IsLetter(c) || c == '_' {
			return l.wordState
		}

		if unicode.IsDigit(c) {
			return l.numberState
		}

		if l.PeekAll(`'\`) {
			return l.runeState
		}

		if unicode.IsSymbol(c) {
			return l.symbolState
		}

		if unicode.IsPrint(c) {
			l.Next()
			l.Emit(lexer.TokenType(c))
			return l.initialState
		}

		l.error(fmt.Sprintf("unexpected token %q (%d)", c, c))
		return nil
	}
}

func (l *states) lineCommentState() lexer.StateFunc {
	l.Next()
	l.Next()
	for {
		c := l.Next()
		l.Ignore()

		if c == lexer.EOFRune || c == '\n' {
			return l.initialState
		}
	}
}

func (l *states) blockCommentState() lexer.StateFunc {
	l.Next()
	l.Next()
	for {
		c := l.Next()
		l.Ignore()

		if c == lexer.EOFRune {
			l.error(fmt.Sprintf("unterminated block comment %s", l.Current()))
			return nil
		}

		if c == '*' && l.Peek() == '/' {
			l.Next()
			l.Ignore()
			return l.initialState
		}
	}
}

func (l *states) whitespaceState() lexer.StateFunc {
	for unicode.IsSpace(l.Peek()) {
		l.Next()
		l.Ignore()
	}

	return l.initialState
}

func (l *states) wordState() lexer.StateFunc {
	for unicode.IsLetter(l.Peek()) ||
		l.Peek() == '_' ||
		unicode.IsDigit(l.Peek()) {
		l.Next()
	}

	l.Emit(WordToken)
	return l.initialState
}

func (l *states) numberState() lexer.StateFunc {
	for unicode.IsDigit(l.Peek()) {
		l.Next()
	}

	l.Emit(NumberToken)
	return l.initialState
}

func (l *states) symbolState() lexer.StateFunc {
	for unicode.IsSymbol(l.Peek()) {
		l.Next()
	}

	if len(l.Current()) == 1 {
		l.Emit(lexer.TokenType(l.Current()[0]))
	} else {
		l.Emit(SymbolToken)
	}
	return l.initialState
}

func (l *states) runeState() lexer.StateFunc {
	// Consume '\''.
	l.Next()

	for l.Peek() != '\'' {
		l.Next()
	}

	// Consume '\''.
	l.Next()

	l.Emit(RuneToken)
	return l.initialState
}

func (l *states) newStringState(name string, delimiter rune) func() lexer.StateFunc {
	return func() lexer.StateFunc {
		l.Next()
		for {
			switch c := l.Next(); c {
			case lexer.EOFRune:
				l.error(fmt.Sprintf("unterminated %s %s", name, l.Current()))
				return nil

			case delimiter:
				l.Emit(StringToken)
				return l.initialState
			}
		}
	}
}

func (l *states) error(err string) {
	l.outErrors = append(l.outErrors, fmt.Sprint("Parse error: ", err))
}

func newStates(file []rune) *states {
	return &states{lexer.New(file), nil /* outErrors */}
}
