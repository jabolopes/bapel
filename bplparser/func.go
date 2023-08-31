package bplparser

func (p *Parser) parseFunc() (Source, error) {
	if err := p.shiftToken("func"); err != nil {
		return Source{}, err
	}

	id, err := p.shiftID()
	if err != nil {
		return Source{}, err
	}

	argTuple, retTuple, err := p.ParseTupleArrow(true /* named */)
	if err != nil {
		return Source{}, err
	}

	if err = p.shiftToken("{"); err != nil {
		return Source{}, err
	}

	if err := p.eol(); err != nil {
		return Source{}, err
	}

	return NewFunctionSource(id, argTuple, retTuple), err
}

func (p *Parser) ParseFunc() (result Source, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseFunc()
		return err
	})
	return
}
