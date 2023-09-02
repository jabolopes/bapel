package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseLetImpl() (Source, error) {
	if err := p.shiftToken("let"); err != nil {
		return Source{}, err
	}

	id, err := p.shiftID()
	if err != nil {
		return Source{}, err
	}

	typ, err := p.parseType(false /* named */)
	if err != nil {
		return Source{}, err
	}

	if err := p.eol(); err != nil {
		return Source{}, err
	}

	return NewDeclSource(ir.NewTermDecl(id, typ)), nil
}

func (p *Parser) parseLet() (result Source, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseLetImpl()
		return err
	})
	return
}
