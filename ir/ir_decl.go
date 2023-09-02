package ir

import (
	"fmt"
	"strings"
)

type IrDeclCase int

const (
	TypeDecl = IrDeclCase(iota)
	TermDecl
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
	ID   string
	Type IrType
}

func (d IrDecl) String() string {
	switch d.Case {
	case TypeDecl:
		var b strings.Builder
		b.WriteString(fmt.Sprintf("type %s : ", d.ID))
		b.WriteString(d.Type.String())
		return b.String()

	case TermDecl:
		var b strings.Builder
		b.WriteString(fmt.Sprintf("%s : ", d.ID))
		b.WriteString(d.Type.String())
		return b.String()

	default:
		panic(fmt.Errorf("unhandled IrDeclCase %d", d.Case))
	}
}

func NewTypeDecl(id string, typ IrType) IrDecl {
	return IrDecl{TypeDecl, id, typ}
}

func NewTermDecl(id string, typ IrType) IrDecl {
	return IrDecl{TermDecl, id, typ}
}
