package bplparser

import (
	"github.com/jabolopes/bapel/ir"
	"golang.org/x/exp/slices"
)

func (p *Parser) parseStatementImpl() (ir.IrTerm, error) {
	if len(p.Words()) > 0 && slices.Contains(p.Words(), "<-") {
		term, err := p.parseAssign()
		if err != nil {
			return ir.IrTerm{}, err
		}
		return ir.NewStatementTerm(term), nil
	}

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
