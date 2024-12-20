package bplparser

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/jabolopes/bapel/parser"
	"golang.org/x/exp/constraints"
)

type Lexer struct {
	scanner *bufio.Scanner
	line    string
	words   []string
	lineNum int
}

func (p *Lexer) Open(reader io.Reader) {
	p.scanner = bufio.NewScanner(reader)
	p.line = ""
	p.words = nil
	p.lineNum = 0
}

func (p *Lexer) Scan() bool {
	if p.scanner == nil {
		return false
	}

	for p.scanner.Scan() {
		line := strings.TrimSpace(p.scanner.Text())

		p.lineNum++

		if line == "" {
			continue
		}

		p.line = line
		p.words = words(line)
		return true
	}

	return false
}

func (p *Lexer) ScanErr() error {
	if p.scanner == nil {
		return fmt.Errorf("must be initialized by calling Open() first")
	}

	return p.scanner.Err()
}

func (p *Lexer) Line() string {
	return p.line
}

func (p *Lexer) Words() []string {
	return p.words
}

func (p *Lexer) LineNum() int {
	return p.lineNum
}

func (p *Lexer) ShiftToken() (parser.Token, error) {
	token, words, err := shiftToken(p.words)
	if err != nil {
		return parser.Token{}, err
	}

	p.words = words
	return token, nil
}

func NewLexer() *Lexer {
	return &Lexer{
		nil, /* scanner */
		"",  /* line */
		nil, /* words */
		0,   /* lineNum */
	}
}

func parseNumber[T constraints.Integer](arg string) (T, error) {
	var value T

	if strings.HasPrefix(arg, "0x") {
		// Hexadecimal
		_, err := fmt.Sscanf(arg, "0x%x", &value)

		return value, err
	}

	// Decimal.
	_, err := fmt.Sscanf(arg, "%d", &value)
	return value, err
}

func parseToken(text string) (parser.Token, error) {
	if value, err := parseNumber[int64](text); err == nil {
		return parser.Token{parser.NumberToken, text, value}, nil
	}
	return parser.NewIDToken(text), nil
}

func shiftToken(args []string) (parser.Token, []string, error) {
	if len(args) == 0 {
		return parser.Token{}, args, io.EOF
	}

	token, err := parseToken(args[0])
	if err != nil {
		return parser.Token{}, args, err
	}

	return token, args[1:], nil
}

func words(text string) []string {
	tokens := []string{}

	var s int
	var n int
	var ch rune
	for n, ch = range text {
		switch ch {
		case '(', ')', '[', ']', '{', '}', ',', '\n', '\'', '!':
			if n > s {
				tokens = append(tokens, text[s:n])
			}
			tokens = append(tokens, string(ch))
			s = n + 1
		case ' ':
			if n > s {
				tokens = append(tokens, text[s:n])
			}
			s = n + 1
		}
	}

	if len(text) > s {
		tokens = append(tokens, text[s:])
	}

	return tokens
}
