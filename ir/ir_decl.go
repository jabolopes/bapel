package ir

import (
	"fmt"
	"strings"

	"github.com/jabolopes/bapel/parser"
)

type irDecl struct {
	id  string
	typ IrType
}

func ParseDecl(args []string) (irDecl, error) {
	id, args, err := parser.Shift(args, fmt.Errorf("expected identifier as first token in declaration; got %v", args))
	if err != nil {
		return irDecl{}, err
	}

	args, err = parser.ShiftIf(args, ":", fmt.Errorf("expected token ':' after the declaration's identifier; got %v", args))
	if err != nil {
		return irDecl{}, err
	}

	if len(args) == 0 {
		return irDecl{}, fmt.Errorf("expected type in declaration; got %v", args)
	}

	typStr := strings.Join(args, " ")
	for _, arg := range args {
		if arg == "->" {
			// Declare function.
			typ, err := ParseFunctionType(typStr)
			if err != nil {
				return irDecl{}, err
			}

			return NewFunctionDecl(id, typ), nil
		}
	}

	// Declare variable.
	typ, err := ParseType(typStr)
	if err != nil {
		return irDecl{}, err
	}

	return NewIntDecl(id, typ), nil
}

// matchesDecl determines if the types of the actual declaration are
// equal to the types of the formal declaration. The name of the
// callee is taken from the formal declaration and ignored in the
// actual declaration.
func matchesDeclImpl(formal, actual irDecl, widen bool) error {
	id := formal.id

	if err := MatchesType(formal.typ, actual.typ, widen); err != nil {
		return fmt.Errorf("symbol %q %v", id, err)
	}

	return nil
}

func matchesDecl(formal, actual irDecl) error {
	return matchesDeclImpl(formal, actual, false /* widen */)
}

func matchesDeclWiden(formal, actual irDecl) error {
	return matchesDeclImpl(formal, actual, true /* widen */)
}

func NewIntDecl(id string, typ IrIntType) irDecl {
	return irDecl{id, NewIntType(typ)}
}

func NewFunctionDecl(id string, typ IrFunctionType) irDecl {
	return irDecl{id, NewFunctionType(typ)}
}
