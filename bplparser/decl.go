package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseDeclImpl(named bool) (Source, error) {
	isType := false
	if err := p.shiftToken("type"); err == nil {
		isType = true
	}

	id, err := p.shiftID()
	if err != nil {
		return Source{}, err
	}

	if err := p.shiftToken(":"); err != nil {
		return Source{}, err
	}

	typ, err := p.parseQuantifiedType(named)
	if err != nil {
		return Source{}, err
	}

	if err := p.eol(); err != nil {
		return Source{}, err
	}

	if isType {
		return NewDeclSource(ir.NewTypeDecl(id, typ)), nil
	}

	return NewDeclSource(ir.NewTermDecl(id, typ)), nil
}

func (p *Parser) parseDecl(named bool) (result Source, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseDeclImpl(named)
		return err
	})
	return
}
