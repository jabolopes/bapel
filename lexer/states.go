package lexer

import (
	"fmt"
	"io"
	"unicode"

	"github.com/jabolopes/bapel/lexerfsm"
)

const (
	WordToken lexerfsm.TokenType = iota
	NumberToken
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
	"!=",
}

type states struct {
	*lexerfsm.LexerFSM
	outErrors []string
}

func (l *states) initialState() lexerfsm.StateFunc {
	switch c := l.Peek(); c {
	case lexerfsm.EOFRune:
		return nil

	case '"':
		return l.newStringState("string", '"')

	case '`':
		return l.newStringState("raw string", '`')

	case '\n':
		return l.newlineState()

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

		if unicode.IsSymbol(c) {
			return l.symbolState
		}

		if unicode.IsPrint(c) {
			l.Next()
			l.Emit(lexerfsm.TokenType(c))
			return l.initialState
		}

		l.error(fmt.Sprintf("unexpected token %q (%d)", c, c))
		return nil
	}
}

func (l *states) newlineState() lexerfsm.StateFunc {
	// Compress a sequence of newlines into a single newline.
	l.Next()
	l.Emit(WordToken)

	for l.Peek() == '\n' {
		l.Next()
		l.Ignore()
	}

	return l.initialState
}

func (l *states) lineCommentState() lexerfsm.StateFunc {
	l.Next()
	l.Next()
	for {
		c := l.Next()
		l.Ignore()

		if c == lexerfsm.EOFRune || c == '\n' {
			return l.initialState
		}
	}
}

func (l *states) blockCommentState() lexerfsm.StateFunc {
	l.Next()
	l.Next()
	for {
		c := l.Next()
		l.Ignore()

		if c == lexerfsm.EOFRune {
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

func (l *states) whitespaceState() lexerfsm.StateFunc {
	for unicode.IsSpace(l.Peek()) && l.Peek() != '\n' {
		l.Next()
		l.Ignore()
	}

	return l.initialState
}

func (l *states) wordState() lexerfsm.StateFunc {
	for unicode.IsLetter(l.Peek()) ||
		l.Peek() == '_' ||
		l.Peek() == '.' ||
		unicode.IsDigit(l.Peek()) {
		l.Next()
	}

	l.Emit(WordToken)
	return l.initialState
}

func (l *states) numberState() lexerfsm.StateFunc {
	for unicode.IsDigit(l.Peek()) {
		l.Next()
	}

	if l.TakeAll(".") {
		for unicode.IsDigit(l.Peek()) {
			l.Next()
		}
	}

	l.Emit(NumberToken)
	return l.initialState
}

func (l *states) symbolState() lexerfsm.StateFunc {
	for unicode.IsSymbol(l.Peek()) {
		l.Next()
	}

	if len(l.Current()) == 1 {
		l.Emit(lexerfsm.TokenType(l.Current()[0]))
	} else {
		l.Emit(SymbolToken)
	}
	return l.initialState
}

func (l *states) newStringState(name string, delimiter rune) func() lexerfsm.StateFunc {
	return func() lexerfsm.StateFunc {
		l.Next()
		for {
			switch c := l.Next(); c {
			case lexerfsm.EOFRune:
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

func newStates(reader io.Reader) *states {
	return &states{
		lexerfsm.New(reader),
		nil, /* outErrors */
	}
}
