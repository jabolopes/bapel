package ir

type IrVarType int

const (
	ArgVar = IrVarType(iota)
	RetVar
	LocalVar
)

// irVar is a variable (e.g., argument, return or local).
type IrVar struct {
	VarType IrVarType // Type of this variable, e.g., arg, ret, local.
	Type    IrType    // Type of this variable, e.g., i8, i16, etc.
	offset  uint16    // Offset in bytes relative to frame pointer.
}

func (v IrVar) Size() int {
	size, err := SizeOfType(v.Type)
	if err != nil {
		panic(err)
	}
	return size
}
