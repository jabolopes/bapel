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

	condition, err := p.parseCall()
	if err != nil {
		return ir.IrTerm{}, err
	}

	then, err := p.parseBlock()
	if err != nil {
		return ir.IrTerm{}, err
	}

	var elseTerm *ir.IrTerm
	if p.shiftLiteral("else") == nil {
		term, err := p.parseBlock()
		if err != nil {
			return ir.IrTerm{}, err
		}

		if err := p.eol(); err != nil {
			return ir.IrTerm{}, err
		}

		elseTerm = &term
	}

	return ir.NewIfTerm(negate, condition, then, elseTerm), nil
}

func (p *Parser) parseIf() (result ir.IrTerm, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseIfImpl()
		return err
	})
	return
}
