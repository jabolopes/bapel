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

func (s Filename) Format(f fmt.State, verb rune) {
	if addMetadata := f.Flag('+'); addMetadata {
		s.Pos.Format(f, verb)
	}

	fmt.Fprintf(f, fmt.FormatString(f, verb), s.Value)
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
