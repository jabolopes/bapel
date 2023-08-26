package bplparser

import (
	"github.com/jabolopes/bapel/parser"
)

func (p *Parser) ParseEntity(args []string) (string, []string, error) {
	orig := args

	args, err := parser.ShiftToken(args, "entity")
	if err != nil {
		return "", orig, err
	}

	id, args, err := parser.ShiftID(args)
	if err != nil {
		return "", orig, err
	}

	args, err = parser.ShiftToken(args, "{")
	if err != nil {
		return "", orig, err
	}

	args, err = parser.ShiftToken(args, "}")
	if err != nil {
		return "", orig, err
	}

	if err := parser.EOL(args); err != nil {
		return "", orig, err
	}

	return id, args, nil
}
