package bplparser

import "github.com/jabolopes/bapel/parser"

func (p *Parser) ParseEnd(args []string) ([]string, error) {
	orig := args

	args, err := parser.ShiftToken(args, "}")
	if err != nil {
		return orig, err
	}

	if err := parser.EOL(args); err != nil {
		return orig, err
	}

	return nil, nil
}
