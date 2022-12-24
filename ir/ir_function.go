package ir

// irFunction is a function.
type irFunction struct {
	id     string // Name of function.
	offset uint64 // Offset of function relative to program data.
	vars   map[string]IrVar
}

func (f irFunction) argsBytes() uint16 {
	var size uint16
	for _, irvar := range f.vars {
		if irvar.VarType == ArgVar {
			size += uint16(irvar.Size())
		}
	}
	return size
}

func (f irFunction) retsBytes() uint16 {
	var size uint16
	for _, irvar := range f.vars {
		if irvar.VarType == RetVar {
			size += uint16(irvar.Size())
		}
	}
	return size
}

func (f irFunction) localsBytes() uint16 {
	var size uint16
	for _, irvar := range f.vars {
		if irvar.VarType == LocalVar {
			size += uint16(irvar.Size())
		}
	}
	return size
}
