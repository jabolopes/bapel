package ir

import "fmt"

type IrDeclCase int

const (
	TypeDecl = IrDeclCase(iota)
	ConstantDecl
	VarDecl
)

func (c IrDeclCase) String() string {
	switch c {
	case TypeDecl:
		return "type declaration"
	case ConstantDecl:
		return "constant declaration"
	case VarDecl:
		return "variable declaration"
	default:
		panic(fmt.Errorf("Unhandled IrDeclCase %d", c))
	}
}

type irDecl struct {
	Case IrDeclCase
	ID   string
	Type IrType
}

// TODO: Make struct public and delete type alias.
type IrDecl = irDecl

func NewTypeDecl(id string, typ IrType) irDecl {
	return irDecl{TypeDecl, id, typ}
}

func NewConstantDecl(id string, typ IrType) irDecl {
	return irDecl{ConstantDecl, id, typ}
}

func NewVarDecl(id string, typ IrType) irDecl {
	return irDecl{VarDecl, id, typ}
}
