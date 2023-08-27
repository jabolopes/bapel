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

func NewConstantDecl(id string, typ IrType) IrDecl {
	return IrDecl{ConstantDecl, id, typ}
}

func NewVarDecl(id string, typ IrType) IrDecl {
	return IrDecl{VarDecl, id, typ}
}
