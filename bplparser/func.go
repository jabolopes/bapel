package bplparser

import (
	"fmt"
	"strings"

	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseTypeAbstraction() ([]string, error) {
	if err := p.shiftLiteral("["); err != nil {
		return nil, err
	}

	vars := []string{}
	for {
		if !p.peekRune(func(r rune) bool { return r == '\'' }) {
			return nil, fmt.Errorf(`expected token "'"`)
		}

		id, err := p.shiftID()
		if err != nil {
			return nil, err
		}

		vars = append(vars, strings.TrimPrefix(id, "'"))

		if p.shiftLiteral(",") == nil {
			continue
		}

		break
	}

	if err := p.shiftLiteral("]"); err != nil {
		return nil, err
	}

	return vars, nil
}

func (p *Parser) parseFuncBindList() ([]ir.IrDecl, error) {
	if err := p.shiftLiteral("("); err != nil {
		return nil, err
	}

	if p.shiftLiteral(")") == nil {
		return nil, nil
	}

	var decls []ir.IrDecl
	for {
		id, err := p.shiftID()
		if err != nil {
			return nil, err
		}

		typ, err := p.parseType()
		if err != nil {
			return nil, err
		}

		decls = append(decls, ir.NewTermDecl(id, typ))

		if p.shiftLiteral(",") == nil {
			continue
		}

		break
	}

	if err := p.shiftLiteral(")"); err != nil {
		return nil, err
	}

	return decls, nil
}

func (p *Parser) parseFuncArgsAndRets() ([]ir.IrDecl, []ir.IrDecl, error) {
	argTuple, err := p.parseFuncBindList()
	if err != nil {
		return nil, nil, fmt.Errorf("in argument list: %v", err)
	}

	if err := p.shiftLiteral("->"); err != nil {
		return nil, nil, err
	}

	retTuple, err := p.parseFuncBindList()
	if err != nil {
		return nil, nil, fmt.Errorf("in return list: %v", err)
	}

	return argTuple, retTuple, nil
}

func (p *Parser) parseFuncImpl() (Source, error) {
	if err := p.shiftLiteral("func"); err != nil {
		return Source{}, err
	}

	id, err := p.shiftID()
	if err != nil {
		return Source{}, err
	}

	var vars []string
	if p.peek("[") {
		var err error
		vars, err = p.parseTypeAbstraction()
		if err != nil {
			return Source{}, err
		}
	}

	argTuple, retTuple, err := p.parseFuncArgsAndRets()
	if err != nil {
		return Source{}, err
	}

	if err = p.shiftLiteral("{"); err != nil {
		return Source{}, err
	}

	if err := p.eol(); err != nil {
		return Source{}, err
	}

	return NewFunctionSource(ir.NewFunction(id, vars, argTuple, retTuple)), nil
}

func (p *Parser) parseFunc() (result Source, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseFuncImpl()
		return err
	})
	return
}
