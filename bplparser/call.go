package bplparser

import (
	"fmt"
	"strings"
	"unicode"

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

func (p *Parser) parseExpressions() []ir.IrTerm {
	var terms []ir.IrTerm
	for {
		term, err := p.parseExpression()
		if err != nil {
			break
		}

		terms = append(terms, term)
	}

	return terms
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

func (p *Parser) parseIDAndArgs(id string) ([]ir.IrType, ir.IrTerm, error) {
	if err := p.shiftLiteral(id); err != nil {
		return nil, ir.IrTerm{}, err
	}

	types, err := p.parseCallTypesOpt()
	if err != nil {
		return nil, ir.IrTerm{}, err
	}

	terms, err := p.parseTerms()
	if err != nil {
		return nil, ir.IrTerm{}, err
	}

	return types, ir.NewTupleTerm(terms), nil
}

func (p *Parser) parseOpUnary() (ir.IrTerm, error) {
	id, err := p.shiftID()
	if err != nil {
		return ir.IrTerm{}, err
	}

	types, err := p.parseCallTypesOpt()
	if err != nil {
		return ir.IrTerm{}, err
	}

	term, err := p.parseExpression()
	if err != nil {
		return ir.IrTerm{}, err
	}

	// 0 - $term
	args := []ir.IrTerm{ir.NewTokenTerm(parser.NewNumberToken(0)), term}
	return ir.NewCallTerm(id, types, ir.NewTupleTerm(args)), nil
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

	types, term, err := p.parseIDAndArgs(id)
	if err != nil {
		return ir.IrTerm{}, err
	}

	if len(types) > 0 {
		return ir.IrTerm{}, fmt.Errorf("expected no call types; got %v", types)
	}

	if !term.Is(ir.TupleTerm) || len(term.Tuple) != 2 {
		return ir.IrTerm{}, fmt.Errorf("%q expects 2 arguments; got %v", id, term)
	}

	return ir.NewIndexGetTerm(term.Tuple[0], term.Tuple[1]), nil
}

func (p *Parser) parseIndexSet() (ir.IrTerm, error) {
	const id = "Index.set"

	types, term, err := p.parseIDAndArgs(id)
	if err != nil {
		return ir.IrTerm{}, err
	}

	if len(types) > 0 {
		return ir.IrTerm{}, fmt.Errorf("expected no call types; got %v", types)
	}

	if !term.Is(ir.TupleTerm) || len(term.Tuple) != 3 {
		return ir.IrTerm{}, fmt.Errorf("%q expects 3 arguments; got %v", id, term)
	}

	return ir.NewIndexSetTerm(term.Tuple[0], term.Tuple[1], term.Tuple[2]), nil
}

func (p *Parser) parseWiden() (ir.IrTerm, error) {
	const id = "widen"

	types, term, err := p.parseIDAndArgs(id)
	if err != nil {
		return ir.IrTerm{}, err
	}

	if len(types) > 0 {
		return ir.IrTerm{}, fmt.Errorf("expected no call types; got %v", types)
	}

	return ir.NewWidenTerm(term), nil
}

func (p *Parser) parseFunctionCall() (ir.IrTerm, error) {
	id, err := p.shiftID()
	if err != nil {
		return ir.IrTerm{}, err
	}

	types, err := p.parseCallTypesOpt()
	if err != nil {
		return ir.IrTerm{}, err
	}

	terms := p.parseExpressions()
	if len(types) == 0 && len(terms) == 0 {
		return ir.NewTokenTerm(parser.NewIDToken(id)), nil
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

	if p.peekRune(unicode.IsLetter) {
		term, err := p.parseFunctionCall()
		if err == nil {
			return term, nil
		}
	}

	return ir.IrTerm{}, fmt.Errorf("expected function call")
}

func (p *Parser) parseCall() (result ir.IrTerm, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseCallImpl()
		return err
	})
	return
}
