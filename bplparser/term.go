package bplparser

import "github.com/jabolopes/bapel/ir"

func (p *Parser) parseTermImpl() (ir.IrTerm, error) {
	if p.peek("if") {
		return p.parseIf()
	}

	if p.peek("let") {
		term, err := p.parseLet()
		if err != nil {
			return ir.IrTerm{}, err
		}
		return ir.NewStatementTerm(term), nil
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
