package ir

type IrVarType int

const (
	ArgVar = IrVarType(iota)
	RetVar
	LocalVar
)

// IrVar is a variable (e.g., argument, return or local).
type IrVar struct {
	Id      string    // Name of this variable, e.g., 'var1'.
	VarType IrVarType // Type of this variable, e.g., arg, ret, local.
	Type    IrIntType // Type of this variable, e.g., i8, i16, etc.
	offset  int       // Offset in bytes relative to frame pointer.
}

func (v *IrVar) decl() irDecl {
	return NewDecl(v.Id, NewIntType(v.Type))
}
