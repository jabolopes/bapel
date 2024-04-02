package bplparser

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
	"golang.org/x/exp/slices"
)

func (p *Parser) parseSectionImpl() (Source, error) {
	id, err := p.shiftID()
	if err != nil {
		return Source{}, err
	}

	if err = p.shiftToken("{"); err != nil {
		return Source{}, err
	}

	sections := []string{"imports", "decls", "exports"}
	if !slices.Contains(sections, id) {
		return Source{}, fmt.Errorf("expected one of %v; got %s", sections, id)
	}

	if err := p.eol(); err != nil {
		return Source{}, err
	}

	var decls []ir.IrDecl
	for p.Scan() {
		if p.peek("}") {
			_ = p.shiftToken("}")

			if err := p.eol(); err != nil {
				return Source{}, err
			}

			break
		}

		decl, err := p.parseDecl(false /* named */)
		if err != nil {
			return Source{}, err
		}

		decls = append(decls, decl)
	}

	return NewSectionSource(id, decls), nil
}

func (p *Parser) parseSection() (result Source, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseSectionImpl()
		return err
	})
	return
}
