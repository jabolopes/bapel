package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseIfImpl() (ir.IrTerm, error) {
	if err := p.shiftLiteral("if"); err != nil {
		return ir.IrTerm{}, err
	}

	negate := false
	if p.shiftLiteral("not") == nil {
		negate = true
	}

	if err := p.shiftLiteralEnd("{"); err != nil {
		return ir.IrTerm{}, err
	}

	condition, err := p.parseCall()
	if err != nil {
		return ir.IrTerm{}, err
	}

	return ir.NewIfTerm(negate, condition), nil
}

func (p *Parser) parseIf() (result ir.IrTerm, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseIfImpl()
		return err
	})
	return
}

func (p *Parser) parseElseImpl() error {
	if err := p.shiftLiteral("}"); err != nil {
		return err
	}

	if err := p.shiftLiteral("else"); err != nil {
		return err
	}

	if err := p.shiftLiteral("{"); err != nil {
		return err
	}

	return p.eol()
}

func (p *Parser) parseElse() error {
	return p.withCheckpoint(p.parseElseImpl)
}
