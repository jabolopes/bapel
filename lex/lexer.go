package lex

import (
	"errors"
	"io"
	"strings"

	"github.com/jabolopes/bapel/lex/lexer"
)

type Lexer struct {
	// Lexer state machine.
	states *states
	// Filter that converts newlines into semicolons.
	lineFilter *lineFilter
}

func (l *Lexer) NextToken() (lexer.Token, bool) { return l.lineFilter.NextToken() }

// ScanErr returns any errors that occurred while processing the
// data. It should be called when `NextToken` returns 'false'.
func (l *Lexer) ScanErr() error {
	var errs []string
	if l.states.Err() != io.EOF && l.states.Err() != nil {
		errs = append(errs, l.states.Err().Error())
	}

	if len(l.states.outErrors) > 0 {
		errs = append(errs, l.states.outErrors...)
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}

func New(reader io.Reader) *Lexer {
	states := newStates(reader)

	lex := &Lexer{
		states,
		newLineFilter(states),
	}
	lex.states.Start(states.initialState)
	return lex
}
