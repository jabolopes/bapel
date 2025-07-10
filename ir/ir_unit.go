package ir

type IrUnitCase int

const (
	BaseUnit IrUnitCase = iota
	ImplUnit
)

type IrImport struct {
	ModuleID string
}

func NewImport(moduleID string) IrImport {
	return IrImport{moduleID}
}

type IrImpl struct {
	RelativeFilename string
}

func NewImpl(relativeFilename string) IrImpl {
	return IrImpl{relativeFilename}
}

type IrUnit struct {
	Case        IrUnitCase
	ModuleID    string
	Filename    string
	Imports     []IrImport
	Impls       []IrImpl
	ImportDecls []IrDecl
	ImplDecls   []IrDecl
	Decls       []IrDecl
	Functions   []IrFunction
}
