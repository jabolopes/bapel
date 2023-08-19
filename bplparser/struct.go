package bplparser

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func ParseStruct(args []string) (string, ir.IrStructType, []string, error) {
	args, err := parser.ShiftIf(args, "struct", fmt.Errorf("expected 'struct'; got %v", args))
	if err != nil {
		return "", ir.IrStructType{}, args, err
	}

	id, args, err := parser.Shift(args, fmt.Errorf("expected identifier; got %v", args))
	if err != nil {
		return "", ir.IrStructType{}, args, err
	}

	typ, args, err := ParseStructType(args, true /* named */)
	if err != nil {
		return "", ir.IrStructType{}, args, err
	}

	return id, typ, args, err
}
