package bplparser

import (
	"github.com/jabolopes/bapel/ir"
)

func (p *Parser) parseFunctionType(named bool) (ir.IrFunctionType, error) {
	argTuple, retTuple, err := p.ParseTupleArrow(named)
	if err != nil {
		return ir.IrFunctionType{}, err
	}

	argTypes := make([]ir.IrType, len(argTuple))
	for i := range argTuple {
		argTypes[i] = argTuple[i].Type
	}

	retTypes := make([]ir.IrType, len(retTuple))
	for i := range retTuple {
		retTypes[i] = retTuple[i].Type
	}

	return ir.IrFunctionType{argTypes, retTypes}, nil
}

func (p *Parser) ParseFunctionType(named bool) (result ir.IrFunctionType, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseFunctionType(named)
		return err
	})
	return
}
