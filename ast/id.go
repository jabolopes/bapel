package ast

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

type ID struct {
	Pos   ir.Pos
	Value string
}

func (i ID) Format(f fmt.State, verb rune) {
	if addMetadata := f.Flag('+'); addMetadata {
		i.Pos.Format(f, verb)
	}
	fmt.Fprint(f, i.Value)
}
