package bplparser

import (
	"github.com/jabolopes/bapel/parser"
)

func ParseIf(args []string) (bool, []parser.Token, []string, error) {
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

	return then, argTokens, nil, nil
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
