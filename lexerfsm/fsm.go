// lexerfsm the finite state machinery for implementing lexers.
//
// Inspired on https://github.com/bbuck/go-lexer.
package lexerfsm

import (
	"strings"
	"unicode/utf8"

	"github.com/emirpasic/gods/v2/stacks"
	"github.com/emirpasic/gods/v2/stacks/arraystack"
)

type StateFunc func() StateFunc

const (
	EOFRune rune = -1

	channelSize = 4096
)

type LexerFSM struct {
	source          string
	start, position int
	tokens          chan Token
	rewind          stacks.Stack[rune]
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

// Current returns the value being being analyzed at this moment.
func (l *LexerFSM) Current() string {
	return l.source[l.start:l.position]
}

// Emit will receive a token type and push a new token with the current analyzed
// value into the tokens channel.
func (l *LexerFSM) Emit(t TokenType) {
	tok := Token{
		Type:  t,
		Value: l.Current(),
	}
	l.tokens <- tok
	l.start = l.position
	l.rewind.Clear()
}

// Ignore clears the rewind stack and then sets the current beginning position
// to the current position in the source which effectively ignores the section
// of the source being analyzed.
func (l *LexerFSM) Ignore() {
	l.rewind.Clear()
	l.start = l.position
}

// Peek performs a Next operation immediately followed by a Rewind returning the
// peeked rune.
func (l *LexerFSM) Peek() rune {
	r := l.Next()
	l.Rewind()
	return r
}

// PeekAll peeks all the characters in the given string. Returns true if they
// match, false otherwise.
func (l *LexerFSM) PeekAll(str string) bool {
	match := true
	nexts := 0
	for _, c := range str {
		d := l.Next()
		nexts++
		if c != d {
			match = false
			break
		}
	}

	for i := 0; i < nexts; i++ {
		l.Rewind()
	}

	return match
}

// Rewind will take the last rune read (if any) and rewind back. Rewinds can
// occur more than once per call to Next but you can never rewind past the
// last point a token was emitted.
func (l *LexerFSM) Rewind() {
	r, ok := l.rewind.Pop()
	if ok && r > EOFRune {
		size := utf8.RuneLen(r)
		l.position -= size
		if l.position < l.start {
			l.position = l.start
		}
	}
}

// Next pulls the next rune from the Lexer and returns it, moving the position
// forward in the source.
func (l *LexerFSM) Next() rune {
	var (
		r rune
		s int
	)
	str := l.source[l.position:]
	if len(str) == 0 {
		r, s = EOFRune, 0
	} else {
		r, s = utf8.DecodeRuneInString(str)
	}
	l.position += s
	l.rewind.Push(r)

	return r
}

// Take receives a string containing all acceptable strings and will continue
// over each consecutive character in the source until a token not in the given
// string is encountered. This should be used to quickly pull token parts.
func (l *LexerFSM) Take(chars string) {
	r := l.Next()
	for strings.ContainsRune(chars, r) {
		r = l.Next()
	}
	l.Rewind() // last next wasn't a match
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

// NextToken returns the next token from the lexer and a value to denote whether
// or not the token is finished.
func (l *LexerFSM) NextToken() (Token, bool) {
	if tok, ok := <-l.tokens; ok {
		return tok, true
	} else {
		return Token{}, false
	}
}

// New creates a returns a lexer ready to parse the given source code.
func New(src string) *LexerFSM {
	return &LexerFSM{
		source:   src,
		start:    0,
		position: 0,
		rewind:   arraystack.New[rune](),
	}
}
