package ir

import "fmt"

type IrVarType int

const (
	ArgVar = IrVarType(iota)
	RetVar
	LocalVar
)

func (t IrVarType) String() string {
	switch t {
	case ArgVar:
		return "arg"
	case RetVar:
		return "ret"
	case LocalVar:
		return "local"
	default:
		panic(fmt.Errorf("Unhandled IrVarType %d", t))
	}
}

// IrVar is a variable (e.g., argument, return or local).
type IrVar struct {
	Id      string    // Name of this variable, e.g., 'var1'.
	VarType IrVarType // Type of this variable, e.g., arg, ret, local.
	Type    IrType    // Type of this variable, e.g., i8, i16, etc.
}

func (v IrVar) String() string {
	return fmt.Sprintf("%s %s %s", v.Id, v.VarType, v.Type)
}

func (v *IrVar) decl() irDecl {
	return NewVarDecl(v.Id, v.Type)
}

func NewVar(id string, varType IrVarType, typ IrType) IrVar {
	return IrVar{id, varType, typ}
}
