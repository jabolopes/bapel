package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

type BindCase int

const (
	TermBind BindCase = iota
	AliasBind
	// Type constant binding. A constant is, e.g., 'i8', 'i16', etc.
	ConstBind
	ScopeBind
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

type constBind struct {
	Name   string
	Kind   ir.IrKind
	Symbol Symbol
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

type typeVarBind struct {
	Name string
	Kind ir.IrKind
}

func (b *typeVarBind) String() string {
	return fmt.Sprintf("type '%s", b.Name)
}

type Bind struct {
	Case    BindCase
	Term    *termBind
	Alias   *aliasBind
	Const   *constBind
	Scope   *scopeBind
	TypeVar *typeVarBind
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
	case ConstBind:
		return b.Const.String()
	case ScopeBind:
		return b.Scope.String()
	case TypeVarBind:
		return b.TypeVar.String()
	default:
		panic(fmt.Errorf("unhandled %T %d", b.Case, b.Case))
	}
}

func (b Bind) Is(c BindCase) bool {
	return b.Case == c
}

func (b Bind) ID() string {
	switch b.Case {
	case TermBind:
		return b.Term.Name
	case AliasBind:
		return b.Alias.Name
	case ConstBind:
		return b.Const.Name
	case ScopeBind:
		return b.Scope.String() // Not entirely correct.
	case TypeVarBind:
		return b.TypeVar.Name
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

func NewConstBind(name string, kind ir.IrKind, symbol Symbol) Bind {
	return Bind{
		Case:  ConstBind,
		Const: &constBind{name, kind, symbol},
	}
}

func NewScopeBind(level int) Bind {
	return Bind{
		Case:  ScopeBind,
		Scope: &scopeBind{level},
	}
}

func NewTypeVarBind(typeVar string, kind ir.IrKind) Bind {
	return Bind{
		Case:    TypeVarBind,
		TypeVar: &typeVarBind{typeVar, kind},
	}
}
