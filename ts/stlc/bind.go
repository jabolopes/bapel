package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

type BindCase int

const (
	AliasBind BindCase = iota
	// Type constant binding. A constant is, e.g., 'i8', 'i16', etc.
	ConstBind
	// Static scope.
	//
	// For example, blocks, functions, etc, introduce new static scopes.
	ScopeBind
	// Term declaration, e.g., 'x: () -> ()'
	TermDeclBind
	// Term definition, e.g., 'fn x() -> () ...'
	TermDefBind
	// Type parameter binding.
	TypeParamBind
	// Trait binding.
	TraitBind
	// Trait implementation binding.
	TraitImplBind
)

type aliasBind struct {
	Name string
	Type ir.IrType
}

func (b *aliasBind) String() string {
	return fmt.Sprintf("type %s = %s", b.Name, b.Type)
}

type constBind struct {
	Name string
	Kind ir.IrKind
}

func (b *constBind) String() string {
	return fmt.Sprintf("type %s :: %s", b.Name, b.Kind)
}

type scopeBind struct {
	Level int
}

func (b *scopeBind) String() string {
	return fmt.Sprintf("scope %d", b.Level)
}

type termDeclBind struct {
	Name string
	Type ir.IrType
}

func (b *termDeclBind) String() string {
	return fmt.Sprintf("%s: %s", b.Name, b.Type)
}

type termDefBind struct {
	Name string
	Type ir.IrType
}

func (b *termDefBind) String() string {
	return fmt.Sprintf("let %s: %s", b.Name, b.Type)
}

type typeParamBind struct {
	Name   string
	Kind   ir.IrKind
	Bounds []ir.IrType
}

func (b *typeParamBind) String() string {
	return fmt.Sprintf("type '%s", b.Name)
}

type traitBind struct {
	Name       string
	TypeParams []ir.TypeParam
	Methods    []ir.IrSignature
}

func (b *traitBind) String() string {
	return fmt.Sprintf("trait %s", b.Name)
}

type traitImplBind struct {
	TypeParams []ir.TypeParam
	TraitType  ir.IrType // Changed from TraitName string
	TypeName   ir.IrType
}

func (b *traitImplBind) String() string {
	return fmt.Sprintf("impl %s for %s", b.TraitType, b.TypeName)
}

type Bind struct {
	Case      BindCase
	Alias     *aliasBind
	Const     *constBind
	Scope     *scopeBind
	TermDecl  *termDeclBind
	TermDef   *termDefBind
	TypeParam *typeParamBind
	Trait     *traitBind
	TraitImpl *traitImplBind
}

func (b Bind) String() string {
	if b.Case == 0 && b.Alias == nil && b.Trait == nil && b.TraitImpl == nil {
		return ""
	}

	switch b.Case {
	case AliasBind:
		return b.Alias.String()
	case ConstBind:
		return b.Const.String()
	case ScopeBind:
		return b.Scope.String()
	case TermDeclBind:
		return b.TermDecl.String()
	case TermDefBind:
		return b.TermDef.String()
	case TypeParamBind:
		return b.TypeParam.String()
	case TraitBind:
		return b.Trait.String()
	case TraitImplBind:
		return b.TraitImpl.String()
	default:
		panic(fmt.Errorf("unhandled %T %d", b.Case, b.Case))
	}
}

func (b Bind) Is(c BindCase) bool {
	return b.Case == c
}

func NewAliasBind(name string, typ ir.IrType) Bind {
	return Bind{
		Case:  AliasBind,
		Alias: &aliasBind{name, typ},
	}
}

func NewConstBind(name string, kind ir.IrKind) Bind {
	return Bind{
		Case:  ConstBind,
		Const: &constBind{name, kind},
	}
}

func NewScopeBind(level int) Bind {
	return Bind{
		Case:  ScopeBind,
		Scope: &scopeBind{level},
	}
}

func NewTermDeclBind(name string, typ ir.IrType) Bind {
	return Bind{
		Case:     TermDeclBind,
		TermDecl: &termDeclBind{name, typ},
	}
}

func NewTermDefBind(name string, typ ir.IrType) Bind {
	return Bind{
		Case:    TermDefBind,
		TermDef: &termDefBind{name, typ},
	}
}

func NewTypeParamBind(tp ir.TypeParam) Bind {
	return Bind{
		Case:      TypeParamBind,
		TypeParam: &typeParamBind{tp.Var, tp.Kind, tp.Bounds},
	}
}

func NewTraitBind(name string, typeParams []ir.TypeParam, methods []ir.IrSignature) Bind {
	return Bind{
		Case:  TraitBind,
		Trait: &traitBind{name, typeParams, methods},
	}
}

func NewTraitImplBind(typeParams []ir.TypeParam, traitType ir.IrType, typeName ir.IrType) Bind {
	return Bind{
		Case:      TraitImplBind,
		TraitImpl: &traitImplBind{typeParams, traitType, typeName},
	}
}


