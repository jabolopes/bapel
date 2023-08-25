package bplparser

import (
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func ParseIf(args []string) (bool, []ir.IrTerm, []string, error) {
	orig := args

	args, err := parser.ShiftToken(args, "if")
	if err != nil {
		return false, nil, orig, err
	}

	args, err = parser.ShiftTokenEnd(args, "{")
	if err != nil {
		return false, nil, orig, err
	}

	then := true
	if args, err = parser.ShiftTokenEnd(args, "else"); err == nil {
		then = false
	}

	argTokens, args, err := ParseCall(args)
	if err != nil {
		return false, nil, orig, err
	}

	argTerms := make([]ir.IrTerm, len(argTokens))
	for i := range argTokens {
		argTerms[i] = ir.NewTokenTerm(argTokens[i])
	}

	return then, argTerms, nil, nil
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
