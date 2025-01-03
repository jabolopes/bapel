// lexerfsm the finite state machinery for implementing lexers.
//
// Inspired by https://github.com/bbuck/go-lexer.
package lexerfsm

import (
	"io"
	"strings"

	"github.com/jabolopes/bapel/scanner"
)

type StateFunc func() StateFunc

const (
	EOFRune rune = -1

	channelSize = 4096
)

type LexerFSM struct {
	scanner *scanner.Scanner
	tokens  chan Token
	lineNum int
	read    strings.Builder
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

// Returns the underlying error (if any). Returns io.EOF if EOF was reached.
func (l *LexerFSM) Err() error { return l.scanner.Err() }

// Current returns the value being being analyzed at this moment.
func (l *LexerFSM) Current() string {
	return l.read.String()
}

// Emit emits a token with the currently read input (if any).
func (l *LexerFSM) Emit(t TokenType) {
	l.tokens <- Token{l.lineNum, t, l.read.String()}
	l.read.Reset()
}

// Ignore discards any read input (via Next()) so that it doesn't become part of
// the emitted token.
func (l *LexerFSM) Ignore() {
	l.read.Reset()
}

// Peek peeks the next rune in the reader without removing from the next read.
func (l *LexerFSM) Peek() rune {
	r, ok := l.scanner.PeekRune()
	if !ok {
		return EOFRune
	}
	return r
}

// PeekAll peeks all the runes in the given string in the reader without
// removing it from the next read(s). Returns true if the whole string matches,
// false otherwise.
func (l *LexerFSM) PeekAll(str string) bool {
	rs, ok := l.scanner.PeekRunes(len(str))
	if !ok || len(str) != len(rs) {
		return false
	}

	i := 0
	for _, r := range str {
		if r != rs[i] {
			return false
		}
		i++
	}

	return true
}

// Next returns the next rune. Returns EOFRune if EOF was reached or an error
// was encountered.
func (l *LexerFSM) Next() rune {
	if l.read.Len() == 0 {
		l.lineNum = l.scanner.LineNum()
	}

	r, err := l.scanner.ReadRune()
	if err != nil {
		return EOFRune
	}

	l.read.WriteRune(r)
	return r
}

// TakeAll takes the given string if it matches the input, otherwise takes
// nothing. The string must match completely. This does not take a prefix of the
// string. Returns true if the input was taken, false otherwise.
func (l *LexerFSM) TakeAll(str string) bool {
	if l.PeekAll(str) {
		for range str {
			l.Next()
		}
		return true
	}
	return false
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

func New(reader io.Reader) *LexerFSM {
	return &LexerFSM{
		scanner: scanner.New(reader),
		lineNum: 1,
	}
}
