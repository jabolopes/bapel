package lexer

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func words(text string) []string {
	tokens := []string{}

	var s int
	var n int
	var ch rune
	for n, ch = range text {
		switch {
		case ch == ' ':
			if n > s {
				tokens = append(tokens, text[s:n])
			}
			s = n + 1

		case strings.ContainsAny(text[n:n+1], "()[]{},\n'!|"):
			if n > s {
				tokens = append(tokens, text[s:n])
			}
			tokens = append(tokens, string(ch))
			s = n + 1
		}
	}

	if len(text) > s {
		tokens = append(tokens, text[s:])
	}

	return tokens
}

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

func (p *Lexer) LineNum() int {
	return p.lineNum
}

func (p *Lexer) ShiftWord() (string, bool) {
	if len(p.words) == 0 {
		return "", false
	}

	word := p.words[0]
	p.words = p.words[1:]
	return word, true
}

func New() *Lexer {
	return &Lexer{
		nil, /* scanner */
		"",  /* line */
		nil, /* words */
		0,   /* lineNum */
	}
}
