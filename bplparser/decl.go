package bplparser

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func ParseDecl(args []string, named bool) (ir.IrDecl, []string, error) {
	id, args, err := parser.Shift(args, fmt.Errorf("expected identifier as first token in declaration; got %v", args))
	if err != nil {
		return ir.IrDecl{}, nil, err
	}

	args, err = parser.ShiftToken(args, ":")
	if err != nil {
		return ir.IrDecl{}, nil, err
	}

	if len(args) == 0 {
		return ir.IrDecl{}, nil, fmt.Errorf("expected type in declaration; got %v", args)
	}

	typ, args, err := ParseType(args, named)
	if err != nil {
		return ir.IrDecl{}, nil, err
	}

	return ir.NewDecl(id, typ), args, nil
}
