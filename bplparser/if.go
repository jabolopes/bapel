package bplparser

import (
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func ParseIf(args []string) (ir.IrTerm, []string, error) {
	orig := args

	args, err := parser.ShiftToken(args, "if")
	if err != nil {
		return ir.IrTerm{}, orig, err
	}

	args, err = parser.ShiftTokenEnd(args, "{")
	if err != nil {
		return ir.IrTerm{}, orig, err
	}

	then := true
	if args, err = parser.ShiftTokenEnd(args, "else"); err == nil {
		then = false
	}

	condition, args, err := ParseCall2(args)
	if err != nil {
		return ir.IrTerm{}, orig, err
	}

	return ir.NewIfTerm(then, condition), args, nil
}

func ParseElse(args []string) ([]string, error) {
	orig := args

	args, err := parser.ShiftTokens(args, []string{"}", "else", "{"})
	if err != nil {
		return orig, err
	}

	if err := parser.EOL(args); err != nil {
		return orig, err
	}

	return nil, nil
}
