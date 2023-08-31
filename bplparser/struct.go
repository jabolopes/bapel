package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseStruct() (Source, error) {
	if err := p.shiftToken("struct"); err != nil {
		return Source{}, err
	}

	id, err := p.shiftID()
	if err != nil {
		return Source{}, err
	}

	typ, err := p.ParseStructType(true /* named */)
	if err != nil {
		return Source{}, err
	}

	if err := p.eol(); err != nil {
		return Source{}, err
	}

	return NewDeclSource(ir.NewTypeDecl(id, typ)), nil
}

func (p *Parser) ParseStruct() (result Source, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseStruct()
		return err
	})
	return
}
