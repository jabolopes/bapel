package ir

import (
	"fmt"
)

type irDecl struct {
	id  string
	typ IrType
}

// TODO: Make struct public and delete type alias.
type IrDecl = irDecl

// matchesDecl determines if the types of the actual declaration are
// equal to the types of the formal declaration. The name of the
// callee is taken from the formal declaration and ignored in the
// actual declaration.
func matchesDeclImpl(formal, actual irDecl, widen bool) error {
	if err := MatchesType(formal.typ, actual.typ, widen); err != nil {
		return fmt.Errorf("symbol %q definition %v does not match its declaration type %v typ: %v", formal.id, actual.typ, formal.typ, err)
	}

	return nil
}

func matchesDecl(formal, actual irDecl) error {
	return matchesDeclImpl(formal, actual, false /* widen */)
}

func matchesDeclWiden(formal, actual irDecl) error {
	return matchesDeclImpl(formal, actual, true /* widen */)
}

func NewDecl(id string, typ IrType) irDecl {
	return irDecl{id, typ}
}
