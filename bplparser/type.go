package bplparser

import (
	"fmt"
	"unicode"

	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseTypeVariable() (ir.IrType, error) {
	if err := p.shiftLiteral("'"); err != nil {
		return ir.IrType{}, err
	}

	id, err := p.shiftID()
	if err != nil {
		return ir.IrType{}, err
	}

	return ir.NewVarType(id), nil
}

func (p *Parser) parseAtomTypeImpl() (ir.IrType, error) {
	if p.peek("'") {
		return p.parseTypeVariable()
	}

	if p.peekRune(unicode.IsLetter) {
		token, err := p.shiftID()
		if err != nil {
			return ir.IrType{}, err
		}

		return ir.NewNameType(token), nil
	}

	return ir.IrType{}, fmt.Errorf("expected atom type")
}

func (p *Parser) parseAtomType() (result ir.IrType, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseAtomTypeImpl()
		return err
	})
	return
}

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

	{
		var types []ir.IrType
		for {
			atom, err := p.parseAtomType()
			if err != nil {
				break
			}

			types = append(types, atom)
		}

		switch len(types) {
		case 0:
			return ir.IrType{}, fmt.Errorf("expected type")
		case 1:
			return types[0], nil
		default:
			retType := types[0]
			for _, typ := range types[1:] {
				retType = ir.NewAppType(retType, typ)
			}
			return retType, nil
		}
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
