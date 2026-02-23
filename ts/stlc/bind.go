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
	// Type variable binding.
	TypeVarBind
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

type typeVarBind struct {
	Name string
	Kind ir.IrKind
}

func (b *typeVarBind) String() string {
	return fmt.Sprintf("type '%s", b.Name)
}

type Bind struct {
	Case     BindCase
	Alias    *aliasBind
	Const    *constBind
	Scope    *scopeBind
	TermDecl *termDeclBind
	TermDef  *termDefBind
	TypeVar  *typeVarBind
}

func (b Bind) String() string {
	if b.Case == 0 && b.Alias == nil {
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
	case TypeVarBind:
		return b.TypeVar.String()
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

func NewTypeVarBind(typeVar string, kind ir.IrKind) Bind {
	return Bind{
		Case:    TypeVarBind,
		TypeVar: &typeVarBind{typeVar, kind},
	}
}
