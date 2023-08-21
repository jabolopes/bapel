package bplparser

import (
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func ParseDecl(args []string, named bool) (ir.IrDecl, []string, error) {
	orig := args

	id, args, err := parser.ShiftID(args)
	if err != nil {
		return ir.IrDecl{}, orig, err
	}

	args, err = parser.ShiftToken(args, ":")
	if err != nil {
		return ir.IrDecl{}, orig, err
	}

	typ, args, err := ParseType(args, named)
	if err != nil {
		return ir.IrDecl{}, orig, err
	}

	if err := parser.EOL(args); err != nil {
		return ir.IrDecl{}, orig, err
	}

	return ir.NewDecl(id, typ), args, nil
}
