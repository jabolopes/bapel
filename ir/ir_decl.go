package ir

type irDecl struct {
	id  string
	typ IrType
}

// TODO: Make struct public and delete type alias.
type IrDecl = irDecl

func NewDecl(id string, typ IrType) irDecl {
	return irDecl{id, typ}
}
