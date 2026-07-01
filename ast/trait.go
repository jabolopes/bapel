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
	Export  bool
	ID      string
	Methods []Signature
	Pos     ir.Pos
}

func (t Trait) Format(f fmt.State, verb rune) {
	if addMetadata := f.Flag('+'); addMetadata {
		t.Pos.Format(f, verb)
	}

	if t.Export {
		fmt.Fprint(f, "pub ")
	}
	fmt.Fprintf(f, "trait %s {\n", t.ID)
	for _, m := range t.Methods {
		fmt.Fprintf(f, "  %s\n", m)
	}
	fmt.Fprint(f, "}")
}

// Impl represents a trait implementation.
type Impl struct {
	TraitName string
	TypeName  ir.IrType
	Methods   []Function
	Pos       ir.Pos
}

func (t Impl) Format(f fmt.State, verb rune) {
	if addMetadata := f.Flag('+'); addMetadata {
		t.Pos.Format(f, verb)
	}

	fmt.Fprintf(f, "impl %s for %s {\n", t.TraitName, t.TypeName)
	for _, m := range t.Methods {
		fmt.Fprintf(f, "  %s\n", m)
	}
	fmt.Fprint(f, "}")
}

func NewSignature(pos ir.Pos, id string, args []ir.FunctionArg, retType ir.IrType) Signature {
	return Signature{id, args, retType, pos}
}

func NewTrait(pos ir.Pos, export bool, id string, methods []Signature) Trait {
	return Trait{export, id, methods, pos}
}

func NewImpl(pos ir.Pos, traitName string, typeName ir.IrType, methods []Function) Impl {
	return Impl{traitName, typeName, methods, pos}
}
