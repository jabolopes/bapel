package bplparser

func (p *Parser) parsePrintU() (Source, error) {
	if err := p.shiftToken("printU"); err != nil {
		return Source{}, err
	}

	return NewPrintSource(p.Words()), nil
}

func (p *Parser) parsePrintS() (Source, error) {
	if err := p.shiftToken("printS"); err != nil {
		return Source{}, err
	}

	return NewPrintSource(p.Words()), nil
}
