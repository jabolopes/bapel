package bplparser

import (
	"fmt"

	"golang.org/x/exp/slices"
)

func (p *Parser) parseSection() (Source, error) {
	section, err := p.shiftID()
	if err != nil {
		return Source{}, err
	}

	if err = p.shiftToken("{"); err != nil {
		return Source{}, err
	}

	sections := []string{"imports", "decls", "exports"}
	if !slices.Contains(sections, section) {
		return Source{}, fmt.Errorf("expected one of %v; got %s", sections, section)
	}

	if err := p.eol(); err != nil {
		return Source{}, err
	}

	return NewSectionSource(section), nil
}

func (p *Parser) ParseSection() (result Source, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseSection()
		return err
	})
	return
}
