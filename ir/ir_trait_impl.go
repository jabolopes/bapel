package ir

import (
	"fmt"
)

type ImplCase int

const (
	TraitImpl = ImplCase(iota)
	InherentImpl
)

type IrTraitImpl struct {
	Case      ImplCase
	TraitName string // Only valid if Case == TraitImpl
	TypeName  IrType
	Methods   []IrFunction
	Pos       Pos
}

func (t IrTraitImpl) Format(f fmt.State, verb rune) {
	if addMetadata := f.Flag('+'); addMetadata {
		t.Pos.Format(f, verb)
	}

	if t.Case == InherentImpl {
		fmt.Fprintf(f, "impl %s {\n", t.TypeName)
	} else {
		fmt.Fprintf(f, "impl %s for %s {\n", t.TraitName, t.TypeName)
	}
	for _, m := range t.Methods {
		// We need to indent the function formatting.
		// For now, we just print it.
		fmt.Fprintf(f, "  %s\n", m)
	}
	fmt.Fprint(f, "}")
}

func NewTraitImpl(pos Pos, traitName string, typeName IrType, methods []IrFunction) IrTraitImpl {
	return IrTraitImpl{Case: TraitImpl, TraitName: traitName, TypeName: typeName, Methods: methods, Pos: pos}
}

func NewInherentImpl(pos Pos, typeName IrType, methods []IrFunction) IrTraitImpl {
	return IrTraitImpl{Case: InherentImpl, TypeName: typeName, Methods: methods, Pos: pos}
}
