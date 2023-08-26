package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseIf() (ir.IrTerm, error) {
	if err := p.shiftToken("if"); err != nil {
		return ir.IrTerm{}, err
	}

	if err := p.shiftTokenEnd("{"); err != nil {
		return ir.IrTerm{}, err
	}

	then := true
	if err := p.shiftTokenEnd("else"); err == nil {
		then = false
	}

	condition, err := p.ParseCall()
	if err != nil {
		return ir.IrTerm{}, err
	}

	return ir.NewIfTerm(then, condition), nil
}

func (p *Parser) ParseIf() (result ir.IrTerm, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseIf()
		return err
	})
	return
}

func (p *Parser) parseElse() error {
	if err := p.shiftToken("}"); err != nil {
		return err
	}

	if err := p.shiftToken("else"); err != nil {
		return err
	}

	if err := p.shiftToken("{"); err != nil {
		return err
	}

	return p.eol()
}

func (p *Parser) ParseElse() error {
	return p.withCheckpoint(p.parseElse)
}
