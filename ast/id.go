package ast

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

type ID struct {
	Value string
	// File information (if any).
	Pos ir.Pos
}

func (i ID) Format(f fmt.State, verb rune) {
	if addMetadata := f.Flag('+'); addMetadata {
		i.Pos.Format(f, verb)
	}

	fmt.Fprintf(f, fmt.FormatString(f, verb), i.Value)
}

func NewID(value string, pos ir.Pos) ID {
	return ID{value, pos}
}
