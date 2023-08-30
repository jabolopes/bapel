package ir

import "fmt"

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

func NewTypeDecl(id string, typ IrType) IrDecl {
	return IrDecl{TypeDecl, id, typ}
}

func NewTermDecl(id string, typ IrType) IrDecl {
	return IrDecl{TermDecl, id, typ}
}
