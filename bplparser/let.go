package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseLetImpl() (ir.IrTerm, error) {
	if err := p.shiftLiteral("let"); err != nil {
		return ir.IrTerm{}, err
	}

	id, err := p.shiftID()
	if err != nil {
		return ir.IrTerm{}, err
	}

	typ, err := p.parseQuantifiedType()
	if err != nil {
		return ir.IrTerm{}, err
	}

	var arg *ir.IrTerm
	if p.shiftLiteral("=") == nil {
		term, err := p.parseExpression()
		if err != nil {
			return ir.IrTerm{}, err
		}

		arg = &term
	}

	if err := p.eol(); err != nil {
		return ir.IrTerm{}, err
	}

	return ir.NewLetTerm(ir.NewTermDecl(id, typ), arg), nil
}

func (p *Parser) parseLet() (result ir.IrTerm, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseLetImpl()
		return err
	})
	return
}
