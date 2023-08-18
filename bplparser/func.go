package bplparser

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func ParseFunc(args []string) (string, []ir.IrVar, []string, error) {
	id, args, err := parser.Shift(args, fmt.Errorf("expected identifier as first token in declaration; got %v", args))
	if err != nil {
		return "", nil, nil, err
	}

	vars, args, err := ParseTupleArrow(args, true /* named */)
	if err != nil {
		return "", nil, nil, err
	}

	return id, vars, args, nil
}
