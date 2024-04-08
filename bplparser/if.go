package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseIfImpl() (Source, error) {
	if err := p.shiftLiteral("if"); err != nil {
		return Source{}, err
	}

	if err := p.shiftLiteralEnd("{"); err != nil {
		return Source{}, err
	}

	then := true
	if p.shiftLiteralEnd("else") == nil {
		then = false
	}

	condition, err := p.parseCall()
	if err != nil {
		return Source{}, err
	}

	return NewTermSource(ir.NewIfTerm(then, condition)), nil
}

func (p *Parser) parseIf() (result Source, err error) {
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
