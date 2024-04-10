package bplparser

import "github.com/jabolopes/bapel/ir"

func (p *Parser) parseTermImpl() (ir.IrTerm, error) {
	if p.peek("if") {
		return p.parseIf()
	}

	return p.parseStatement()
}

func (p *Parser) parseTerm() (result ir.IrTerm, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseTermImpl()
		return err
	})
	return
}
