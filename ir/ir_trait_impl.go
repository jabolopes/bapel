package ir

import (
	"fmt"
)

type IrTraitImpl struct {
	TraitName string
	TypeName  IrType
	Methods   []IrFunction
	Pos       Pos
}

func (t IrTraitImpl) Format(f fmt.State, verb rune) {
	if addMetadata := f.Flag('+'); addMetadata {
		t.Pos.Format(f, verb)
	}

	fmt.Fprintf(f, "impl %s for %s {\n", t.TraitName, t.TypeName)
	for _, m := range t.Methods {
		// We need to indent the function formatting.
		// For now, we just print it.
		fmt.Fprintf(f, "  %s\n", m)
	}
	fmt.Fprint(f, "}")
}

func NewTraitImpl(pos Pos, traitName string, typeName IrType, methods []IrFunction) IrTraitImpl {
	return IrTraitImpl{traitName, typeName, methods, pos}
}
