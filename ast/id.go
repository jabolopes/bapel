package ast

import (
	"cmp"
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
	fmt.Fprint(f, i.Value)
}

func CompareID(id1, id2 ID) int {
	return cmp.Compare(id1.Value, id2.Value)
}

func NewID(name string, pos ir.Pos) ID {
	return ID{name, pos}
}
