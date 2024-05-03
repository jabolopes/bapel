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
	NameBind
	TypeVarBind
)

type termBind struct {
	Name   string
	Type   ir.IrType
	Symbol Symbol
}

func (b *termBind) String() string {
	return fmt.Sprintf("%s : %s", b.Name, b.Type)
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

type nameBind struct {
	Name   string
	Symbol Symbol
}

func (b *nameBind) String() string {
	return fmt.Sprintf("type %s", b.Name)
}

type typeVarBind struct {
	Name string
}

func (b *typeVarBind) String() string {
	return fmt.Sprintf("type '%s", b.Name)
}

type Bind struct {
	Case      BindCase
	Term      *termBind
	Alias     *aliasBind
	Component *componentBind
	Name      *nameBind
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
	case NameBind:
		return b.Name.String()
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
	case NameBind:
		return b.Name.Name, true
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
	case NameBind:
		return b.Name.Symbol, true
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

func NewNameBind(name string, symbol Symbol) Bind {
	return Bind{
		Case: NameBind,
		Name: &nameBind{name, symbol},
	}
}

func NewTypeVarBind(typeVar string) Bind {
	return Bind{
		Case:    TypeVarBind,
		TypeVar: &typeVarBind{typeVar},
	}
}
