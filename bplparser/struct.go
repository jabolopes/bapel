package bplparser

import (
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func (p *Parser) ParseStruct(args []string) (string, ir.IrStructType, []string, error) {
	orig := args

	args, err := parser.ShiftToken(args, "struct")
	if err != nil {
		return "", ir.IrStructType{}, orig, err
	}

	id, args, err := parser.ShiftID(args)
	if err != nil {
		return "", ir.IrStructType{}, orig, err
	}

	typ, args, err := p.ParseStructType(args, true /* named */)
	if err != nil {
		return "", ir.IrStructType{}, orig, err
	}

	if err := parser.EOL(args); err != nil {
		return "", ir.IrStructType{}, orig, err
	}

	return id, typ, args, err
}
