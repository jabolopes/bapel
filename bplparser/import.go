package bplparser

func (p *Parser) parseImportImpl() (Source, error) {
	if err := p.shiftLiteral("import"); err != nil {
		return Source{}, err
	}

	id, err := p.shiftID()
	if err != nil {
		return Source{}, err
	}

	if err := p.eol(); err != nil {
		return Source{}, err
	}

	return NewImportSource(id), nil
}

func (p *Parser) parseImport() (result Source, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseImportImpl()
		return err
	})
	return
}
