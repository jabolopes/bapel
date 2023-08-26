package bplparser

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/jabolopes/bapel/parser"
	"golang.org/x/exp/constraints"
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

func (p *Parser) withCheckpoint(callback func() error) {
	orig := p.words
	if err := callback(); err != nil {
		p.words = orig
	}
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

func (p *Parser) SetLine(line string) {
	p.line = line
	p.words = parser.Words(line)
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

func (p *Parser) shiftID() (string, error) {
	id, words, err := parser.ShiftID(p.words)
	if err != nil {
		return "", err
	}

	p.words = words
	return id, nil
}

func (p *Parser) shiftToken(token string) error {
	words, err := parser.ShiftToken(p.words, token)
	if err != nil {
		return err
	}

	p.words = words
	return nil
}

func shiftInteger[T constraints.Integer](p *Parser) (T, error) {
	integer, words, err := parser.ShiftNumber[T](p.words)
	if err != nil {
		var t T
		return t, err
	}

	p.words = words
	return integer, nil
}

func (p *Parser) eol() error {
	return parser.EOL(p.words)
}
