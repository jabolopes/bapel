package bplparser

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/jabolopes/bapel/parser"
)

type IsFunction interface {
	IsFunction(string) bool
}

type Parser struct {
	compiler IsFunction
	scanner  *bufio.Scanner
	line     string
	words    []string
}

func (p *Parser) Open(reader io.Reader) {
	p.scanner = bufio.NewScanner(reader)
	p.line = ""
	p.words = nil
}

func (p *Parser) Scan() bool {
	if p.scanner == nil {
		return false
	}

	for p.scanner.Scan() {
		line := strings.TrimSpace(p.scanner.Text())
		if line == "" {
			continue
		}

		p.words = parser.Words(line)
		return true
	}

	return false
}

func (p *Parser) ScanErr() error {
	if p.scanner == nil {
		return fmt.Errorf("parser was not initialize; Open() must be called first")
	}

	return p.scanner.Err()
}

func (p *Parser) Line() string {
	return p.line
}

func (p *Parser) Words() []string {
	return p.words
}

func NewParser(compiler IsFunction) *Parser {
	return &Parser{
		compiler,
		nil, /* scanner */
		"",  /* line */
		nil, /* words */
	}
}
