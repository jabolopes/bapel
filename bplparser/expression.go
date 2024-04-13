package bplparser

import (
	"unicode"

	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseExpressionImpl() (ir.IrTerm, error) {
	if p.peekRune(unicode.IsDigit) {
		return p.parseNumberTerm()
	}

	if p.peek("(") {
		return p.parseTupleTerm()
	}

	return p.parseCall()
}

func (p *Parser) parseExpression() (result ir.IrTerm, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseExpressionImpl()
		return err
	})
	return
}
