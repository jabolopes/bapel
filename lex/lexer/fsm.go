// lexer the finite state machinery for implementing lexs.
//
// Inspired by https://github.com/bbuck/go-lexer.
package lexer

import (
	"github.com/jabolopes/bapel/lex/scanner"
)

type StateFunc func() StateFunc

const (
	EOFRune rune = -1

	channelSize = 4096
)

type LexerFSM struct {
	*scanner.Scanner
	tokens chan Token
}

func (l *LexerFSM) run(startState StateFunc) {
	state := startState
	for state != nil {
		state = state()
	}
	close(l.tokens)
}

// Start begins executing the Lexer in an asynchronous manner (using a goroutine).
func (l *LexerFSM) Start(startState StateFunc) {
	l.tokens = make(chan Token, channelSize)
	go l.run(startState)
}

// Emit emits a token with the currently read input (if any).
func (l *LexerFSM) Emit(t TokenType) {
	l.tokens <- Token{l.LineNum(), t, l.Current()}
	l.Ignore()
}

// Peek peeks the next rune in the reader without removing from the next read.
func (l *LexerFSM) PeekRune() rune {
	r, ok := l.Scanner.PeekRune()
	if !ok {
		return EOFRune
	}
	return r
}

// ReadRune returns the next rune. Returns EOFRune if EOF was reached or an error
// was encountered.
func (l *LexerFSM) ReadRune() rune {
	r, ok := l.Scanner.ReadRune()
	if !ok {
		return EOFRune
	}

	return r
}

// NextToken returns the next token and 'true', or 'false' if the underlying
// reader reached EOF or an error was encountered. Use 'Err()' to obtain the
// error.
func (l *LexerFSM) NextToken() (Token, bool) {
	if tok, ok := <-l.tokens; ok {
		return tok, true
	} else {
		return Token{}, false
	}
}

func New(file string) *LexerFSM {
	return &LexerFSM{Scanner: scanner.New(file)}
}
