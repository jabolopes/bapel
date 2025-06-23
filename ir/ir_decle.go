package ir

import "fmt"

// An IrDecl extended with the information whether it's exported.
//
// TODO: Can this be merged with IrDecl?
type IrDeclE struct {
	Decl   IrDecl
	Export bool
}

func (d IrDeclE) String() string {
	if d.Export {
		return fmt.Sprintf("export %s", d.Decl)
	}
	return d.Decl.String()
}

func NewDeclE(decl IrDecl, export bool) IrDeclE {
	return IrDeclE{decl, export}
}
