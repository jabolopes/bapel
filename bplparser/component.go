package bplparser

import "github.com/jabolopes/bapel/ir"

func (p *Parser) parseComponentImpl() (Source, error) {
	if err := p.shiftLiteral("component"); err != nil {
		return Source{}, err
	}

	id, err := p.shiftID()
	if err != nil {
		return Source{}, err
	}

	if err := p.shiftLiteral("{"); err != nil {
		return Source{}, err
	}

	typeID, err := p.shiftID()
	if err != nil {
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

	return NewComponentSource(ir.NewComponent(id, typeID, length)), nil
}

func (p *Parser) parseComponent() (result Source, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseComponentImpl()
		return err
	})
	return
}
