package bplparser

import "github.com/jabolopes/bapel/ir"

func (p *Parser) parseEntityImpl() (Source, error) {
	if err := p.shiftLiteral("entity"); err != nil {
		return Source{}, err
	}

	id, err := p.shiftID()
	if err != nil {
		return Source{}, err
	}

	if err := p.shiftLiteral("{"); err != nil {
		return Source{}, err
	}

	length, err := shiftInteger[int](p)
	if err != nil {
		return Source{}, err
	}

	if err := p.shiftLiteral("}"); err != nil {
		return Source{}, err
	}

	if err := p.eol(); err != nil {
		return Source{}, err
	}

	return NewEntitySource(ir.NewEntity(id, length)), nil
}

func (p *Parser) parseEntity() (result Source, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseEntityImpl()
		return err
	})
	return
}
