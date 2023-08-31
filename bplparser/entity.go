package bplparser

func (p *Parser) parseEntity() (Source, error) {
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

func (p *Parser) ParseEntity() (result Source, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseEntity()
		return err
	})
	return
}
