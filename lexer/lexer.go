package lexer

import (
	"bufio"
	"fmt"
	"io"

	"github.com/jabolopes/bapel/bpllexer"
)

type Lexer struct {
	scanner *bufio.Scanner
	line    string
	words   []string
	lineNum int
	err     error
}

func (p *Lexer) Open(reader io.Reader) {
	p.scanner = bufio.NewScanner(reader)
	p.line = ""
	p.words = nil
	p.lineNum = 0
	p.err = nil
}

func (p *Lexer) Scan() bool {
	if p.scanner == nil {
		return false
	}

	for p.scanner.Scan() {
		line := p.scanner.Text()

		p.words = p.words[:0]

		lexer := bpllexer.New(line)
		for {
			token, ok := lexer.ShiftWord()
			if !ok {
				break
			}

			p.words = append(p.words, token.Value)
		}

		if err := lexer.ScanErr(); err != nil {
			p.err = err
			return false
		}

		p.lineNum++

		if line == "" {
			continue
		}

		p.line = line
		return true
	}

	return false
}

func (p *Lexer) ScanErr() error {
	if p.scanner == nil {
		return fmt.Errorf("must be initialized by calling Open() first")
	}

	if p.err != nil {
		return p.err
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
		nil, /* error */
	}
}
