package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseTupleTermImpl() (ir.IrTerm, error) {
	if err := p.shiftLiteral("("); err != nil {
		return ir.IrTerm{}, err
	}

	if p.shiftLiteral(")") == nil {
		return ir.NewTupleTerm(nil), nil
	}

	var terms []ir.IrTerm
	for {
		term, err := p.parseExpression()
		if err != nil {
			return ir.IrTerm{}, err
		}

		terms = append(terms, term)

		if p.shiftLiteral(",") == nil {
			continue
		}

		break
	}

	if err := p.shiftLiteral(")"); err != nil {
		return ir.IrTerm{}, err
	}

	return ir.NewTupleTerm(terms), nil
}

func (p *Parser) parseTupleTerm() (result ir.IrTerm, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseTupleTermImpl()
		return err
	})
	return
}
