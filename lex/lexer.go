package lex

import (
	"errors"
	"strings"

	"github.com/jabolopes/bapel/lex/lexer"
)

type Lexer struct {
	// Lexer state machine.
	states *states
}

func (l *Lexer) NextToken() (lexer.Token, bool) { return l.states.NextToken() }

// ScanErr returns any errors that occurred while processing the
// data. It should be called when `NextToken` returns 'false'.
func (l *Lexer) ScanErr() error {
	var errs []string
	if len(l.states.outErrors) > 0 {
		errs = append(errs, l.states.outErrors...)
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}

func New(file string) *Lexer {
	states := newStates(file)

	lex := &Lexer{states}
	lex.states.Start(states.initialState)

	return lex
}
