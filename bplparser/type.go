package bplparser

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseSimpleTypeImpl(named bool) (ir.IrType, error) {
	if p.peek("(") {
		return p.parseTupleType(false /* named */)
	}

	if p.peek("{") {
		return p.parseStructType(true /* named */)
	}

	if p.peek("[") {
		return p.parseArrayType(named)
	}

	if p.peek("forall") {
		return p.parseForallType(named)
	}

	if p.peekRune(func(r rune) bool { return r == '\'' }) {
		token, err := p.shiftID()
		if err != nil {
			return ir.IrType{}, err
		}

		return ir.NewVarType(strings.TrimPrefix(token, "'")), nil
	}

	if p.peekRune(unicode.IsLetter) {
		token, err := p.shiftID()
		if err != nil {
			return ir.IrType{}, err
		}

		return ir.NewNameType(token), nil
	}

	return ir.IrType{}, fmt.Errorf("expected type")
}

func (p *Parser) parseSimpleType(named bool) (result ir.IrType, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseSimpleTypeImpl(named)
		return err
	})
	return
}

func (p *Parser) parseTypeImpl(named bool) (ir.IrType, error) {
	if typ, err := p.parseFunctionType(named); err == nil {
		return typ, nil
	}

	return p.parseSimpleType(named)
}

func (p *Parser) parseType(named bool) (result ir.IrType, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseTypeImpl(named)
		return err
	})
	return
}

func (p *Parser) parseQuantifiedType(named bool) (result ir.IrType, err error) {
	typ, err := p.parseType(named)
	if err != nil {
		return ir.IrType{}, err
	}

	return ir.QuantifyType(typ), nil
}
