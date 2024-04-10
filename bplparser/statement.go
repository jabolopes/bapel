package bplparser

import "github.com/jabolopes/bapel/ir"

func (p *Parser) parseStatementImpl() (ir.IrTerm, error) {
	if p.peek("let") {
		term, err := p.parseLet()
		if err != nil {
			return ir.IrTerm{}, err
		}
		return ir.NewStatementTerm([]ir.IrTerm{term}), nil
	}

	term, err := p.parseCallAssign()
	if err != nil {
		return ir.IrTerm{}, err
	}

	return ir.NewStatementTerm([]ir.IrTerm{term}), nil
}

func (p *Parser) parseStatement() (result ir.IrTerm, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseStatementImpl()
		return err
	})
	return
}
