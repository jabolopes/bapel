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
	Case       ImplCase
	TypeParams []VarKind
	TraitType  IrType // Changed from TraitName string
	TypeName   IrType
	Methods    []IrFunction
	Pos        Pos
}

func (t IrTraitImpl) Format(f fmt.State, verb rune) {
	if addMetadata := f.Flag('+'); addMetadata {
		t.Pos.Format(f, verb)
	}

	fmt.Fprint(f, "impl")
	if len(t.TypeParams) > 0 {
		fmt.Fprint(f, " [")
		Interleave(t.TypeParams, func() { fmt.Fprint(f, ", ") }, func(_ int, tv VarKind) {
			fmt.Fprintf(f, "'%s", tv.Var)
		})
		fmt.Fprint(f, "]")
	}

	if t.Case == InherentImpl {
		fmt.Fprintf(f, " %s {\n", t.TypeName)
	} else {
		fmt.Fprintf(f, " %s for %s {\n", t.TraitType, t.TypeName)
	}
	for _, m := range t.Methods {
		// We need to indent the function formatting.
		// For now, we just print it.
		fmt.Fprintf(f, "  %s\n", m)
	}
	fmt.Fprint(f, "}")
}

func NewTraitImpl(pos Pos, typeParams []VarKind, traitType IrType, typeName IrType, methods []IrFunction) IrTraitImpl {
	return IrTraitImpl{Case: TraitImpl, TypeParams: typeParams, TraitType: traitType, TypeName: typeName, Methods: methods, Pos: pos}
}

func NewInherentImpl(pos Pos, typeParams []VarKind, typeName IrType, methods []IrFunction) IrTraitImpl {
	return IrTraitImpl{Case: InherentImpl, TypeParams: typeParams, TypeName: typeName, Methods: methods, Pos: pos}
}
