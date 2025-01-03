package lexerfsm_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jabolopes/bapel/lexerfsm"
)

const (
	NumberToken lexerfsm.TokenType = iota
	OpToken
	IdentToken
)

type Lexer struct {
	*lexerfsm.LexerFSM
	errs []error
}

func (l *Lexer) Error(err error) {
	l.errs = append(l.errs, err)
}

func (l *Lexer) NumberState() lexerfsm.StateFunc {
	l.Take("0123456789")
	l.Emit(NumberToken)
	if l.Peek() == '.' {
		l.Next()
		l.Emit(OpToken)
		return l.IdentState
	}

	return nil
}

func (l *Lexer) IdentState() lexerfsm.StateFunc {
	r := l.Next()
	for (r >= 'a' && r <= 'z') || r == '_' {
		r = l.Next()
	}
	l.Rewind()
	l.Emit(IdentToken)

	return l.WhitespaceState
}

func (l *Lexer) WhitespaceState() lexerfsm.StateFunc {
	r := l.Next()
	if r == lexerfsm.EOFRune {
		return nil
	}

	if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
		l.Error(fmt.Errorf("unexpected token %q", r))
		return nil
	}

	l.Take(" \t\n\r")
	l.Ignore()

	return l.NumberState
}

func newLexer(source string) *Lexer {
	return &Lexer{lexerfsm.New(source), nil /* errs */}
}

func TestNext(t *testing.T) {
	l := lexerfsm.New("123")
	run := []struct {
		s string
		r rune
	}{
		{"1", '1'},
		{"12", '2'},
		{"123", '3'},
		{"123", lexerfsm.EOFRune},
	}

	for _, test := range run {
		r := l.Next()
		if r != test.r {
			t.Fatalf("Expected %q but got %q", test.r, r)
		}

		if l.Current() != test.s {
			t.Fatalf("Expected %q but got %q", test.s, l.Current())
		}
	}
}

func TestNumbers(t *testing.T) {
	l := newLexer("123")
	l.Start(l.NumberState)

	tests := []struct {
		want   lexerfsm.Token
		wantOk bool
	}{
		{lexerfsm.Token{NumberToken, "123"}, true},
		{lexerfsm.Token{}, false},
	}

	for _, test := range tests {
		got, gotOk := l.NextToken()
		if !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) || gotOk != test.wantOk {
			t.Fatalf("NextToken() = %v, %v; want %v, %v", got, gotOk, test.want, test.wantOk)
		}
	}
}

func TestRewind(t *testing.T) {
	l := newLexer("1")

	r := l.Next()
	if r != '1' {
		t.Fatalf("Expected %q but got %q", '1', r)
	}

	if l.Current() != "1" {
		t.Fatalf("Expected %q but got %q", "1", l.Current())
	}

	l.Rewind()
	if l.Current() != "" {
		t.Fatalf("Expected empty string, but got %q", l.Current())
	}
}

func TestTokens(t *testing.T) {
	cases := []struct {
		want   lexerfsm.Token
		wantOk bool
	}{
		{lexerfsm.Token{NumberToken, "123"}, true},
		{lexerfsm.Token{OpToken, "."}, true},
		{lexerfsm.Token{IdentToken, "hello"}, true},
		{lexerfsm.Token{NumberToken, "675"}, true},
		{lexerfsm.Token{OpToken, "."}, true},
		{lexerfsm.Token{IdentToken, "world"}, true},
		{lexerfsm.Token{}, false},
	}

	l := newLexer("123.hello  675.world")
	l.Start(l.NumberState)

	for _, test := range cases {
		got, gotOk := l.NextToken()
		if !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) || gotOk != test.wantOk {
			t.Fatalf("NextToken() = %v, %v; want %v, %v", got, gotOk, test.want, test.wantOk)
		}
	}
}

func TestErrors(t *testing.T) {
	l := newLexer("1")
	l.Start(l.WhitespaceState)

	got, gotOk := l.NextToken()
	if !cmp.Equal(got, lexerfsm.Token{}, cmpopts.EquateEmpty()) || gotOk {
		t.Errorf("NextToken() = %v, %v; want %v, %v", got, gotOk, lexerfsm.Token{}, false)
	}

	var gotErr string
	if len(l.errs) > 0 {
		gotErr = l.errs[0].Error()
	}
	wantErr := "unexpected token '1'"
	if !cmp.Equal(gotErr, wantErr, cmpopts.EquateErrors()) {
		t.Errorf("l.Errs = %v; want %v", gotErr, wantErr)
		t.Errorf("Diff = %v", cmp.Diff(l.errs, wantErr, cmpopts.EquateErrors()))
	}
}
