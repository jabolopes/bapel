package lexer_test

import (
	"fmt"
	"strings"
	"testing"
	"unicode"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jabolopes/bapel/lex/lexer"
)

const (
	NumberToken lexer.TokenType = iota
	OpToken
	IdentToken
)

type Lexer struct {
	*lexer.LexerFSM
	errs []error
}

func (l *Lexer) Error(err error) {
	l.errs = append(l.errs, err)
}

func (l *Lexer) NumberState() lexer.StateFunc {
	for unicode.IsDigit(l.Peek()) {
		l.Next()
	}
	l.Emit(NumberToken)
	if l.Peek() == '.' {
		l.Next()
		l.Emit(OpToken)
		return l.IdentState
	}

	return nil
}

func (l *Lexer) IdentState() lexer.StateFunc {
	for unicode.IsLetter(l.Peek()) || l.Peek() == '_' {
		l.Next()
	}
	l.Emit(IdentToken)

	return l.WhitespaceState
}

func (l *Lexer) WhitespaceState() lexer.StateFunc {
	r := l.Next()
	if r == lexer.EOFRune {
		return nil
	}

	if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
		l.Error(fmt.Errorf("unexpected token %q", r))
		return nil
	}

	if unicode.IsSpace(l.Peek()) {
		l.Next()
		l.Ignore()
	}

	return l.NumberState
}

func newLexer(source string) *Lexer {
	return &Lexer{lexer.New(strings.NewReader(source)), nil /* errs */}
}

func TestNext(t *testing.T) {
	l := lexer.New(strings.NewReader("123"))
	run := []struct {
		s string
		r rune
	}{
		{"1", '1'},
		{"12", '2'},
		{"123", '3'},
		{"123", lexer.EOFRune},
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
		want   lexer.Token
		wantOk bool
	}{
		{lexer.Token{1, NumberToken, "123"}, true},
		{lexer.Token{}, false},
	}

	for _, test := range tests {
		got, gotOk := l.NextToken()
		if !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) || gotOk != test.wantOk {
			t.Fatalf("NextToken() = %v, %v; want %v, %v", got, gotOk, test.want, test.wantOk)
		}
	}
}

func TestTokens(t *testing.T) {
	cases := []struct {
		want   lexer.Token
		wantOk bool
	}{
		{lexer.Token{1, NumberToken, "123"}, true},
		{lexer.Token{1, OpToken, "."}, true},
		{lexer.Token{1, IdentToken, "hello"}, true},
		{lexer.Token{1, NumberToken, "675"}, true},
		{lexer.Token{1, OpToken, "."}, true},
		{lexer.Token{1, IdentToken, "world"}, true},
		{lexer.Token{}, false},
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
	if !cmp.Equal(got, lexer.Token{}, cmpopts.EquateEmpty()) || gotOk {
		t.Errorf("NextToken() = %v, %v; want %v, %v", got, gotOk, lexer.Token{}, false)
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
