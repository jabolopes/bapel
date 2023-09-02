package bplparser

func (p *Parser) parseEntityImpl() (Source, error) {
	if err := p.shiftToken("entity"); err != nil {
		return Source{}, err
	}

	id, err := p.shiftID()
	if err != nil {
		return Source{}, err
	}

	if err := p.shiftToken("{"); err != nil {
		return Source{}, err
	}

	if err := p.shiftToken("}"); err != nil {
		return Source{}, err
	}

	if err := p.eol(); err != nil {
		return Source{}, err
	}

	return NewEntitySource(id), nil
}

func (p *Parser) parseEntity() (result Source, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseEntityImpl()
		return err
	})
	return
}
