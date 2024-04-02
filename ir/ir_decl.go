package ir

import (
	"fmt"
)

type IrDeclCase int

const (
	TermDecl = IrDeclCase(iota)
	TypeDecl
)

func (c IrDeclCase) String() string {
	switch c {
	case TypeDecl:
		return "type declaration"
	case TermDecl:
		return "term declaration"
	default:
		panic(fmt.Errorf("unhandled IrDeclCase %d", c))
	}
}

type IrDecl struct {
	Case IrDeclCase
	Term *struct {
		ID   string
		Type IrType
	}
	AsType *struct {
		Type IrType
	}
}

func (d IrDecl) String() string {
	switch d.Case {
	case TermDecl:
		return fmt.Sprintf("%s : %s", d.Term.ID, d.Type())
	case TypeDecl:
		return fmt.Sprintf("type %s", d.Type())
	default:
		panic(fmt.Errorf("unhandled %T %d", d.Case, d.Case))
	}
}

func (d IrDecl) Type() IrType {
	switch d.Case {
	case TermDecl:
		return d.Term.Type
	case TypeDecl:
		return d.AsType.Type
	default:
		panic(fmt.Errorf("unhandled %T %d", d.Case, d.Case))
	}
}

func NewTermDecl(id string, typ IrType) IrDecl {
	return IrDecl{
		Case: TermDecl,
		Term: &struct {
			ID   string
			Type IrType
		}{id, typ},
	}
}

func NewTypeDecl(typ IrType) IrDecl {
	return IrDecl{
		Case:   TypeDecl,
		AsType: &struct{ Type IrType }{typ},
	}
}
