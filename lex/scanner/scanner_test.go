package scanner_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jabolopes/bapel/lex/scanner"
)

func TestScannerReadRune(t *testing.T) {
	s := scanner.New([]rune("hi\nyou\n"))

	tests := []struct {
		wantRune    rune
		wantOk      bool
		wantLineNum int
	}{
		{'h', true, 1},
		{'i', true, 1},
		{'\n', true, 1},
		{'y', true, 2},
		{'o', true, 2},
		{'u', true, 2},
		{'\n', true, 2},
		{0, false, 3},
	}

	for _, test := range tests {
		if got := s.LineNum(); got != test.wantLineNum {
			t.Fatalf("LineNum() = %d; want %d", got, test.wantLineNum)
		}

		if got, gotOk := s.ReadRune(); got != test.wantRune || gotOk != test.wantOk {
			t.Fatalf("ReadRune() = %c, %v; want %c, %v", got, gotOk, test.wantRune, test.wantOk)
		}
	}
}

func TestScannerPeekRune(t *testing.T) {
	s := scanner.New([]rune("hi\nyou\n"))

	if got, gotOk := s.PeekRune(); got != 'h' || !gotOk {
		t.Errorf("PeekRune() = %c, %v; want %c, %v", got, gotOk, 'h', true)
	}
}

func TestScannerPeekRunes(t *testing.T) {
	s := scanner.New([]rune("hi"))

	tests := []struct {
		n      int
		want   []rune
		wantOk bool
	}{
		{0, nil, true},
		{1, []rune{'h'}, true},
		{2, []rune{'h', 'i'}, true},
		{3, nil, false},
	}

	for _, test := range tests {
		if got, gotOk := s.PeekRunes(test.n); !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) || gotOk != test.wantOk {
			t.Errorf("PeekRunes(%d) = %v, %v; want %v, %v", test.n, got, gotOk, test.want, test.wantOk)
		}
	}
}
