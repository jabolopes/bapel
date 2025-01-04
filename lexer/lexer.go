package lexer

import (
	"errors"
	"io"
	"strings"

	"github.com/jabolopes/bapel/lexerfsm"
)

type Lexer struct {
	// Lexer state machine.
	states *states
	// Filter that converts newlines into semicolons.
	lineFilter *lineFilter
}

func (l *Lexer) NextToken() (lexerfsm.Token, bool) { return l.lineFilter.NextToken() }

// ScanErr returns any errors that occurred while processing the
// data. It should be called when ShiftWord() returns 'false'.
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

	lexer := &Lexer{
		states,
		newLineFilter(states),
	}
	lexer.states.Start(states.initialState)
	return lexer
}
