package bplparser

import (
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func ParseFunc(args []string) (string, []ir.IrDecl, []ir.IrDecl, []string, error) {
	orig := args

	args, err := parser.ShiftToken(args, "func")
	if err != nil {
		return "", nil, nil, orig, err
	}

	id, args, err := parser.ShiftID(args)
	if err != nil {
		return "", nil, nil, orig, err
	}

	argTuple, retTuple, args, err := ParseTupleArrow(args, true /* named */)
	if err != nil {
		return "", nil, nil, orig, err
	}

	args, err = parser.ShiftToken(args, "{")
	if err != nil {
		return "", nil, nil, orig, err
	}

	if err := parser.EOL(args); err != nil {
		return "", nil, nil, orig, err
	}

	return id, argTuple, retTuple, args, nil
}
