package bplparser

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseSimpleTypeImpl() (ir.IrType, error) {
	if p.peek("(") {
		return p.parseTupleType()
	}

	if p.peek("{") {
		return p.parseStructType()
	}

	if p.peek("[") {
		return p.parseArrayType()
	}

	if p.peek("forall") {
		return p.parseForallType()
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

func (p *Parser) parseSimpleType() (result ir.IrType, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseSimpleTypeImpl()
		return err
	})
	return
}

func (p *Parser) parseTypeImpl() (ir.IrType, error) {
	if typ, err := p.parseFunctionType(); err == nil {
		return typ, nil
	}

	return p.parseSimpleType()
}

func (p *Parser) parseType() (result ir.IrType, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseTypeImpl()
		return err
	})
	return
}

func (p *Parser) parseQuantifiedType() (result ir.IrType, err error) {
	typ, err := p.parseType()
	if err != nil {
		return ir.IrType{}, err
	}

	return ir.QuantifyType(typ), nil
}
