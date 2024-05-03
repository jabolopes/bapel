package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseStructImpl() (Source, error) {
	if err := p.shiftLiteral("struct"); err != nil {
		return Source{}, err
	}

	id, err := p.shiftID()
	if err != nil {
		return Source{}, err
	}

	structType, err := p.parseStructType()
	if err != nil {
		return Source{}, err
	}

	if err := p.eol(); err != nil {
		return Source{}, err
	}

	return NewTypeDefSource(ir.NewAliasDecl(id, structType)), nil
}

func (p *Parser) parseStruct() (result Source, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseStructImpl()
		return err
	})
	return
}
