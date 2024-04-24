package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseStatementImpl() (ir.IrTerm, error) {
	term, err := p.parseExpression()
	if err != nil {
		return ir.IrTerm{}, err
	}

	return ir.NewStatementTerm(term), nil
}

func (p *Parser) parseStatement() (result ir.IrTerm, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseStatementImpl()
		return err
	})
	return
}
