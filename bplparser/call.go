package bplparser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

// Parse types passing into the call.
//
// Example:
//   f [MyType]
func (p *Parser) parseCallTypes() ([]ir.IrType, error) {
	if err := p.shiftLiteral("["); err != nil {
		return nil, err
	}

	var types []ir.IrType
	for {
		typ, err := p.parseType()
		if err != nil {
			return nil, err
		}

		types = append(types, typ)

		if p.shiftLiteral(",") == nil {
			continue
		}

		break
	}

	if err := p.shiftLiteral("]"); err != nil {
		return nil, err
	}

	return types, nil
}

func (p *Parser) parseCallTypesOpt() ([]ir.IrType, error) {
	if p.peek("[") {
		return p.parseCallTypes()
	}
	return nil, nil
}

func (p *Parser) parseTerms() ([]ir.IrTerm, error) {
	var terms []ir.IrTerm
	for {
		if p.eol() == nil {
			break
		}

		token, err := p.shiftToken()
		if err != nil {
			return nil, err
		}

		terms = append(terms, ir.NewTokenTerm(token))
	}

	return terms, nil
}

func (p *Parser) parseIDAndArgs(id string) ([]ir.IrType, []ir.IrTerm, error) {
	if err := p.shiftLiteral(id); err != nil {
		return nil, nil, err
	}

	types, err := p.parseCallTypesOpt()
	if err != nil {
		return nil, nil, err
	}

	terms, err := p.parseTerms()
	if err != nil {
		return nil, nil, err
	}

	return types, terms, nil
}

func (p *Parser) parseOpUnary() (ir.IrTerm, error) {
	const id = "-"

	types, terms, err := p.parseIDAndArgs(id)
	if err != nil {
		return ir.IrTerm{}, err
	}

	if len(types) > 0 {
		return ir.IrTerm{}, fmt.Errorf("expected no call types; got %v", types)
	}

	// 0 - $terms
	terms = append([]ir.IrTerm{ir.NewTokenTerm(parser.NewNumberToken(0))}, terms...)

	return ir.NewCallTerm(id, nil /* types */, ir.NewTupleTerm(terms)), nil
}

func (p *Parser) parseOpBinary() (ir.IrTerm, error) {
	terms, err := p.parseTerms()
	if err != nil {
		return ir.IrTerm{}, err
	}

	if len(terms) != 3 {
		return ir.IrTerm{}, fmt.Errorf("binary operation expects 3 arguments; got %v", terms)
	}

	idTerm := terms[1]
	terms = append(terms[0:1], terms[2:]...)

	return ir.NewCallTerm(idTerm.Token.Text, nil /* types */, ir.NewTupleTerm(terms)), nil
}

func (p *Parser) parseIndexGet() (ir.IrTerm, error) {
	const id = "Index.get"

	types, terms, err := p.parseIDAndArgs(id)
	if err != nil {
		return ir.IrTerm{}, err
	}

	if len(types) > 0 {
		return ir.IrTerm{}, fmt.Errorf("expected no call types; got %v", types)
	}

	if len(terms) != 2 {
		return ir.IrTerm{}, fmt.Errorf("%q expects 2 arguments; got %v", id, terms)
	}

	return ir.NewIndexGetTerm(terms[0], terms[1]), nil
}

func (p *Parser) parseIndexSet() (ir.IrTerm, error) {
	const id = "Index.set"

	types, terms, err := p.parseIDAndArgs(id)
	if err != nil {
		return ir.IrTerm{}, err
	}

	if len(types) > 0 {
		return ir.IrTerm{}, fmt.Errorf("expected no call types; got %v", types)
	}

	if len(terms) != 3 {
		return ir.IrTerm{}, fmt.Errorf("%q expects 3 arguments; got %v", id, terms)
	}

	return ir.NewIndexSetTerm(terms[0], terms[1], terms[2]), nil
}

func (p *Parser) parseWiden() (ir.IrTerm, error) {
	const id = "widen"

	types, terms, err := p.parseIDAndArgs(id)
	if err != nil {
		return ir.IrTerm{}, err
	}

	if len(types) > 0 {
		return ir.IrTerm{}, fmt.Errorf("expected no call types; got %v", types)
	}

	return ir.NewWidenTerm(ir.NewTupleTerm(terms)), nil
}

func (p *Parser) parseFunctionCall() (ir.IrTerm, error) {
	id, ok := p.getPeek()
	if !ok {
		return ir.IrTerm{}, errors.New("failed to parse function call; expected identifier")
	}

	types, terms, err := p.parseIDAndArgs(id)
	if err != nil {
		return ir.IrTerm{}, err
	}

	return ir.NewCallTerm(id, types, ir.NewTupleTerm(terms)), nil
}

func (p *Parser) parseCallImpl() (ir.IrTerm, error) {
	if p.peek("-") {
		return p.parseOpUnary()
	}

	if len(p.Words()) > 1 && strings.ContainsAny(p.Words()[1], "+-*/") {
		return p.parseOpBinary()
	}

	if p.peek("Index.get") {
		return p.parseIndexGet()
	}

	if p.peek("Index.set") {
		return p.parseIndexSet()
	}

	if p.peek("widen") {
		return p.parseWiden()
	}

	if id, ok := p.getPeek(); ok && p.compiler.IsFunction(id) {
		return p.parseFunctionCall()
	}

	terms, err := p.parseTerms()
	if err != nil {
		return ir.IrTerm{}, err
	}

	return ir.NewTupleTerm(terms), nil
}

func (p *Parser) parseCall() (result ir.IrTerm, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseCallImpl()
		return err
	})
	return
}
