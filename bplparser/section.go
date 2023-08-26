package bplparser

import (
	"fmt"

	"golang.org/x/exp/slices"
)

func (p *Parser) parseSection() (string, error) {
	section, err := p.shiftID()
	if err != nil {
		return "", err
	}

	if err = p.shiftToken("{"); err != nil {
		return "", err
	}

	sections := []string{"imports", "decls", "exports"}
	if !slices.Contains(sections, section) {
		return "", fmt.Errorf("expected one of %v; got %s", sections, section)
	}

	if err := p.eol(); err != nil {
		return "", err
	}

	return section, nil
}

func (p *Parser) ParseSection() (result string, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseSection()
		return err
	})
	return
}
