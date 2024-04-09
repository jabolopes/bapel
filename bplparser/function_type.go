package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

// func (p *Parser) parseFunctionTypeImpl(named bool) (ir.IrType, error) {
// 	argTuple, retTuple, err := p.parseTupleArrow(named)
// 	if err != nil {
// 		return ir.IrType{}, err
// 	}

// 	argTypes := make([]ir.IrType, len(argTuple))
// 	for i := range argTuple {
// 		argTypes[i] = argTuple[i].Type()
// 	}

// 	retTypes := make([]ir.IrType, len(retTuple))
// 	for i := range retTuple {
// 		retTypes[i] = retTuple[i].Type()
// 	}

// 	return ir.NewFunctionType(ir.NewTupleType(argTypes), ir.NewTupleType(retTypes)), nil
// }

func (p *Parser) parseFunctionTypeImpl(named bool) (ir.IrType, error) {
	arg, err := p.parseSimpleType(named)
	if err != nil {
		return ir.IrType{}, err
	}

	if err := p.shiftLiteral("->"); err != nil {
		return ir.IrType{}, err
	}

	ret, err := p.parseType(named)
	if err != nil {
		return ir.IrType{}, err
	}

	return ir.NewFunctionType(arg, ret), nil
}

func (p *Parser) parseFunctionType(named bool) (result ir.IrType, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseFunctionTypeImpl(named)
		return err
	})
	return
}
