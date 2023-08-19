package bplparser

import (
	"github.com/jabolopes/bapel/parser"
)

func ParseEntity(args []string) (string, []string, error) {
	args, err := parser.ShiftToken(args, "entity")
	if err != nil {
		return "", args, err
	}

	id, args, err := parser.ShiftID(args)
	if err != nil {
		return "", args, err
	}

	args, err = parser.ShiftToken(args, "{")
	if err != nil {
		return "", args, err
	}

	args, err = parser.ShiftToken(args, "}")
	if err != nil {
		return "", args, err
	}

	return id, args, nil
}
