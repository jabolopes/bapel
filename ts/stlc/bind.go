package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

type BindCase int

const (
	TermBind BindCase = iota
	AliasBind
	ComponentBind
	// Type constant binding. A constant is, e.g., 'i8', 'i16', etc.
	ConstBind
	TypeVarBind
)

type termBind struct {
	Name   string
	Type   ir.IrType
	Symbol Symbol
}

func (b *termBind) String() string {
	return fmt.Sprintf("%s: %s", b.Name, b.Type)
}

type aliasBind struct {
	Name   string
	Type   ir.IrType
	Symbol Symbol
}

func (b *aliasBind) String() string {
	return fmt.Sprintf("type %s = %s", b.Name, b.Type)
}

type componentBind struct {
	ElemType ir.IrType
}

func (b *componentBind) String() string {
	return fmt.Sprintf("component %s", b.ElemType)
}

type constBind struct {
	Name   string
	Kind   ir.IrKind
	Symbol Symbol
}

func (b *constBind) String() string {
	return fmt.Sprintf("type %s :: %s", b.Name, b.Kind)
}

type typeVarBind struct {
	Name string
	Kind ir.IrKind
}

func (b *typeVarBind) String() string {
	return fmt.Sprintf("type '%s", b.Name)
}

type Bind struct {
	Case      BindCase
	Term      *termBind
	Alias     *aliasBind
	Component *componentBind
	Const     *constBind
	TypeVar   *typeVarBind
}

func (b Bind) String() string {
	if b.Case == 0 && b.Term == nil {
		return ""
	}

	switch b.Case {
	case TermBind:
		return b.Term.String()
	case AliasBind:
		return b.Alias.String()
	case ComponentBind:
		return b.Component.String()
	case ConstBind:
		return b.Const.String()
	case TypeVarBind:
		return b.TypeVar.String()
	default:
		panic(fmt.Errorf("unhandled %T %d", b.Case, b.Case))
	}
}

func (b Bind) Is(c BindCase) bool {
	return b.Case == c
}

func (b Bind) ID() (string, bool) {
	switch b.Case {
	case TermBind:
		return b.Term.Name, true
	case AliasBind:
		return b.Alias.Name, true
	case ComponentBind:
		return "", false
	case ConstBind:
		return b.Const.Name, true
	case TypeVarBind:
		return b.TypeVar.Name, true
	default:
		panic(fmt.Errorf("unhandled %T %d", b.Case, b.Case))
	}
}

func (b Bind) Symbol() (Symbol, bool) {
	switch b.Case {
	case TermBind:
		return b.Term.Symbol, true
	case AliasBind:
		return b.Alias.Symbol, true
	case ComponentBind:
		return Symbol(0), false
	case ConstBind:
		return b.Const.Symbol, true
	case TypeVarBind:
		return Symbol(0), false
	default:
		panic(fmt.Errorf("unhandled %T %d", b.Case, b.Case))
	}
}

func NewTermBind(name string, typ ir.IrType, symbol Symbol) Bind {
	return Bind{
		Case: TermBind,
		Term: &termBind{name, typ, symbol},
	}
}

func NewAliasBind(name string, typ ir.IrType, symbol Symbol) Bind {
	return Bind{
		Case:  AliasBind,
		Alias: &aliasBind{name, typ, symbol},
	}
}

func NewComponentBind(elemType ir.IrType) Bind {
	return Bind{
		Case:      ComponentBind,
		Component: &componentBind{elemType},
	}
}

func NewConstBind(name string, kind ir.IrKind, symbol Symbol) Bind {
	return Bind{
		Case:  ConstBind,
		Const: &constBind{name, kind, symbol},
	}
}

func NewTypeVarBind(typeVar string, kind ir.IrKind) Bind {
	return Bind{
		Case:    TypeVarBind,
		TypeVar: &typeVarBind{typeVar, kind},
	}
}
