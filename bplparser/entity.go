package bplparser

func (p *Parser) ParseEntity() (string, error) {
	if err := p.shiftToken("entity"); err != nil {
		return "", err
	}

	id, err := p.shiftID()
	if err != nil {
		return "", err
	}

	if err := p.shiftToken("{"); err != nil {
		return "", err
	}

	if err := p.shiftToken("}"); err != nil {
		return "", err
	}

	if err := p.eol(); err != nil {
		return "", err
	}

	return id, nil
}
