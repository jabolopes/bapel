package bplparser

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func ParseFunc(args []string) (string, []ir.IrVar, []string, error) {
	orig := args

	args, err := parser.ShiftToken(args, "func")
	if err != nil {
		return "", nil, orig, err
	}

	id, args, err := parser.Shift(args, fmt.Errorf("expected identifier as first token in declaration; got %v", args))
	if err != nil {
		return "", nil, orig, err
	}

	vars, args, err := ParseTupleArrow(args, true /* named */)
	if err != nil {
		return "", nil, orig, err
	}

	args, err = parser.ShiftToken(args, "{")
	if err != nil {
		return "", nil, orig, err
	}

	if err := parser.EOL(args); err != nil {
		return "", nil, orig, err
	}

	return id, vars, args, nil
}
