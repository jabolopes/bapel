package ast

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

// Signature represents a function signature (without a body).
type Signature struct {
	ID      string
	Args    []ir.FunctionArg
	RetType ir.IrType
	Pos     ir.Pos
}

func (s Signature) Format(f fmt.State, verb rune) {
	if addMetadata := f.Flag('+'); addMetadata {
		s.Pos.Format(f, verb)
	}

	fmt.Fprintf(f, "fn %s", s.ID)
	fmt.Fprint(f, "(")
	ir.Interleave(s.Args, func() { fmt.Fprint(f, ", ") }, func(_ int, arg ir.FunctionArg) {
		fmt.Fprint(f, arg.String())
	})
	fmt.Fprintf(f, ") -> %s", s.RetType)
}

// Trait represents a trait declaration.
type Trait struct {
	Export     bool
	ID         string
	TypeParams []ir.TypeParam
	Methods    []Signature
	Pos        ir.Pos
}

func (t Trait) Format(f fmt.State, verb rune) {
	if addMetadata := f.Flag('+'); addMetadata {
		t.Pos.Format(f, verb)
	}

	if t.Export {
		fmt.Fprint(f, "pub ")
	}
	fmt.Fprintf(f, "trait %s", t.ID)
	if len(t.TypeParams) > 0 {
		fmt.Fprint(f, " [")
		ir.Interleave(t.TypeParams, func() { fmt.Fprint(f, ", ") }, func(_ int, tv ir.TypeParam) {
			fmt.Fprintf(f, "'%s", tv.Var)
		})
		fmt.Fprint(f, "]")
	}
	fmt.Fprint(f, " {\n")
	for _, m := range t.Methods {
		fmt.Fprintf(f, "  %s\n", m)
	}
	fmt.Fprint(f, "}")
}

type ImplCase int

const (
	TraitImpl = ImplCase(iota)
	InherentImpl
)

// Impl represents a trait or inherent implementation.
type Impl struct {
	Case       ImplCase
	TypeParams []ir.TypeParam
	TraitType  ir.IrType // Changed from TraitName string
	TypeName   ir.IrType
	Methods    []Function
	Pos        ir.Pos
}

func (t Impl) Format(f fmt.State, verb rune) {
	if addMetadata := f.Flag('+'); addMetadata {
		t.Pos.Format(f, verb)
	}

	fmt.Fprint(f, "impl")
	if len(t.TypeParams) > 0 {
		fmt.Fprint(f, " [")
		ir.Interleave(t.TypeParams, func() { fmt.Fprint(f, ", ") }, func(_ int, tv ir.TypeParam) {
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
		fmt.Fprintf(f, "  %s\n", m)
	}
	fmt.Fprint(f, "}")
}

func NewSignature(pos ir.Pos, id string, args []ir.FunctionArg, retType ir.IrType) Signature {
	return Signature{id, args, retType, pos}
}

func NewTrait(pos ir.Pos, export bool, id string, typeParams []ir.TypeParam, methods []Signature) Trait {
	return Trait{export, id, typeParams, methods, pos}
}

func NewTraitImpl(pos ir.Pos, typeParams []ir.TypeParam, traitType ir.IrType, typeName ir.IrType, methods []Function) Impl {
	return Impl{Case: TraitImpl, TypeParams: typeParams, TraitType: traitType, TypeName: typeName, Methods: methods, Pos: pos}
}

func NewInherentImpl(pos ir.Pos, typeParams []ir.TypeParam, typeName ir.IrType, methods []Function) Impl {
	return Impl{Case: InherentImpl, TypeParams: typeParams, TypeName: typeName, Methods: methods, Pos: pos}
}

func (t Trait) Decl() ir.IrDecl {
	var irMethods []ir.IrSignature
	for _, m := range t.Methods {
		irArgs := make([]ir.FunctionArg, len(m.Args))
		for i, arg := range m.Args {
			irArgs[i] = ir.FunctionArg{ID: arg.ID, Type: arg.Type}
		}
		irMethods = append(irMethods, ir.IrSignature{
			ID:      m.ID,
			Args:    irArgs,
			RetType: m.RetType,
		})
	}

	decl := ir.NewTraitDecl(t.ID, t.TypeParams, irMethods, t.Export)
	decl.Pos = t.Pos
	return decl
}
func (impl Impl) ToIr() (ir.IrTraitImpl, error) {
	var irMethods []ir.IrFunction
	for _, m := range impl.Methods {
		function, err := DesugarFunction(m)
		if err != nil {
			return ir.IrTraitImpl{}, err
		}
		irMethods = append(irMethods, function)
	}

	if impl.Case == InherentImpl {
		return ir.NewInherentImpl(impl.Pos, impl.TypeParams, impl.TypeName, irMethods), nil
	}
	return ir.NewTraitImpl(impl.Pos, impl.TypeParams, impl.TraitType, impl.TypeName, irMethods), nil
}
