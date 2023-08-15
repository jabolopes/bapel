package bplparser

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func ParseDecl(args []string) (ir.IrDecl, error) {
	id, args, err := parser.Shift(args, fmt.Errorf("expected identifier as first token in declaration; got %v", args))
	if err != nil {
		return ir.IrDecl{}, err
	}

	args, err = parser.ShiftIf(args, ":", fmt.Errorf("expected token ':' after the declaration's identifier; got %v", args))
	if err != nil {
		return ir.IrDecl{}, err
	}

	if len(args) == 0 {
		return ir.IrDecl{}, fmt.Errorf("expected type in declaration; got %v", args)
	}

	typ, err := ParseType(args)
	if err != nil {
		return ir.IrDecl{}, err
	}

	return ir.NewDecl(id, typ), nil
}
