package scanner_test

import (
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jabolopes/bapel/scanner"
)

func TestScannerReadRune(t *testing.T) {
	s := scanner.New(strings.NewReader("hi\nyou\n"))

	tests := []struct {
		wantRune    rune
		wantErr     error
		wantLineNum int
	}{
		{'h', nil, 1},
		{'i', nil, 1},
		{'\n', nil, 1},
		{'y', nil, 2},
		{'o', nil, 2},
		{'u', nil, 2},
		{'\n', nil, 2},
		{0, io.EOF, 3},
	}

	for _, test := range tests {
		if got := s.LineNum(); got != test.wantLineNum {
			t.Fatalf("LineNum() = %d; want %d", got, test.wantLineNum)
		}

		if got, gotErr := s.ReadRune(); got != test.wantRune || gotErr != test.wantErr {
			t.Fatalf("ReadRune() = %c, %v; want %c, %v", got, gotErr, test.wantRune, test.wantErr)
		}
	}
}

func TestScannerPeekRune(t *testing.T) {
	s := scanner.New(strings.NewReader("hi\nyou\n"))

	if got, gotOk := s.PeekRune(); got != 'h' || !gotOk {
		t.Errorf("PeekRune() = %c, %v; want %c, %v", got, gotOk, 'h', true)
	}
}

func TestScannerPeekRunes(t *testing.T) {
	s := scanner.New(strings.NewReader("hi"))

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
