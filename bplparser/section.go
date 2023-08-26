package bplparser

import (
	"fmt"

	"github.com/jabolopes/bapel/parser"
	"golang.org/x/exp/slices"
)

func (p *Parser) ParseSection(args []string) (string, []string, error) {
	orig := args

	section, args, err := parser.ShiftID(args)
	if err != nil {
		return "", orig, err
	}

	args, err = parser.ShiftToken(args, "{")
	if err != nil {
		return "", orig, err
	}

	sections := []string{"imports", "decls", "exports"}
	if !slices.Contains(sections, section) {
		return "", orig, fmt.Errorf("expected one of %v; got %s", sections, section)
	}

	if err := parser.EOL(args); err != nil {
		return "", orig, err
	}

	return section, nil, nil
}
