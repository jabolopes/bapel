package ast

import (
	"cmp"
	"fmt"
	"path"
	"strings"

	"github.com/jabolopes/bapel/ir"
)

type Filename struct {
	Value string
	// File information (if any).
	Pos ir.Pos
}

func (i Filename) Format(f fmt.State, verb rune) {
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

func NewFilename(value string, pos ir.Pos) Filename {
	return Filename{path.Clean(value), pos}
}

func ValidateFilename(filename Filename) error {
	if len(filename.Value) == 0 {
		return fmt.Errorf("invalid filename %q. Filenames cannot be empty", filename)
	}

	if strings.Contains(filename.Value, `"`) {
		return fmt.Errorf(`invalid filename %q. Filenames cannot contain quotes ('"')`, filename)
	}

	return nil
}

func CompareFilename(id1, id2 Filename) int {
	return cmp.Compare(id1.Value, id2.Value)
}
