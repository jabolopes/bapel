package bplparser

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/jabolopes/bapel/parser"
)

type Parser struct {
	scanner *bufio.Scanner
	line    string
	words   []string
	lineNum int
}

func (p *Parser) Open(reader io.Reader) {
	p.scanner = bufio.NewScanner(reader)
	p.line = ""
	p.words = nil
	p.lineNum = 0
}

func (p *Parser) Scan() bool {
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
		p.words = parser.Words(line)
		return true
	}

	return false
}

func (p *Parser) ScanErr() error {
	if p.scanner == nil {
		return fmt.Errorf("must be initialized by calling Open() first")
	}

	return p.scanner.Err()
}

func (p *Parser) Line() string {
	return p.line
}

func (p *Parser) Words() []string {
	return p.words
}

func (p *Parser) LineNum() int {
	return p.lineNum
}

func (p *Parser) ShiftToken() (parser.Token, error) {
	token, words, err := parser.ShiftToken(p.words)
	if err != nil {
		return parser.Token{}, err
	}

	p.words = words
	return token, nil
}

func NewParser() *Parser {
	return &Parser{
		nil, /* scanner */
		"",  /* line */
		nil, /* words */
		0,   /* lineNum */
	}
}
