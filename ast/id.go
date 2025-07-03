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

	isQuoted := verb == 'q'
	if isQuoted {
		fmt.Fprint(f, `"`)
	}
	fmt.Fprint(f, i.Value)
	if isQuoted {
		fmt.Fprint(f, `"`)
	}
}

func NewID(value string, pos ir.Pos) ID {
	return ID{value, pos}
}

func ValidateID(id ID) error {
	if len(id.Value) == 0 {
		return fmt.Errorf("invalid ID %q. IDs cannot be empty", id)
	}

	return nil
}

func CompareID(id1, id2 ID) int {
	return cmp.Compare(id1.Value, id2.Value)
}
