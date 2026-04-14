package scanner_test

import (
	"testing"
	"unicode/utf8"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jabolopes/bapel/lex/scanner"
)

func TestScannerReadRune(t *testing.T) {
	s := scanner.New("test", "hi\nyou\n")

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
		{utf8.RuneError, false, 3},
	}

	for _, test := range tests {
		if got := s.LineNum(); got != test.wantLineNum {
			t.Fatalf("LineNum() = %d; want %d", got, test.wantLineNum)
		}

		if got, gotOk := s.ReadRune(); got != test.wantRune || gotOk != test.wantOk {
			t.Fatalf("ReadRune() = %c, %v; want %c, %v", got, gotOk, test.wantRune, test.wantOk)
		}

		s.Ignore()
	}
}

func TestScannerPeekRune(t *testing.T) {
	s := scanner.New("test", "hi\nyou\n")

	if got, gotOk := s.PeekRune(); got != 'h' || !gotOk {
		t.Errorf("PeekRune() = %c, %v; want %c, %v", got, gotOk, 'h', true)
	}
}

func TestScannerPeekString(t *testing.T) {
	s := scanner.New("test", "hi")

	tests := []struct {
		str  string
		want bool
	}{
		{"", true},
		{"h", true},
		{"hi", true},
		{"hio", false},
	}

	for _, test := range tests {
		if got := s.PeekString(test.str); !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) {
			t.Errorf("PeekString(%q) = %v; want %v", test.str, got, test.want)
		}
	}
}
