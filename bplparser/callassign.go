package bplparser

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func (p *Parser) parseAssignImpl() (ir.IrTerm, error) {
	isAssign := false
	var rets []string
	for {
		token, err := p.shift()
		if err != nil {
			break
		}

		if token == "<-" {
			isAssign = true
			break
		}

		rets = append(rets, token)
	}

	if !isAssign {
		return ir.IrTerm{}, fmt.Errorf("expected token '<-' in assignment term")
	}

	if len(rets) == 0 {
		return ir.IrTerm{}, fmt.Errorf("expected at least 1 return value before token '<-'")
	}

	callTerm, err := p.parseCall()
	if err != nil {
		return ir.IrTerm{}, err
	}

	retTokens, err := parser.ParseTokens(rets)
	if err != nil {
		return ir.IrTerm{}, err
	}

	retTerms := make([]ir.IrTerm, len(retTokens))
	for i := range retTokens {
		retTerms[i] = ir.NewTokenTerm(retTokens[i])
	}

	return ir.NewAssignTerm(callTerm, ir.NewTupleTerm(retTerms)), nil
}

func (p *Parser) parseAssign() (result ir.IrTerm, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseAssignImpl()
		return err
	})
	return
}

// ParseCallAssign parses call and assignment.
//
// Note that a call is an assignment without the '<-' and without any return
// values.
func (p *Parser) parseCallAssign() (ir.IrTerm, error) {
	if typ, err := p.parseAssign(); err == nil {
		return typ, nil
	}

	return p.parseCall()
}
