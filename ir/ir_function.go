package ir

// irVar is a variable (e.g., argument, return or local).
type irVar struct {
	offset uint16 // Offset in bytes relative to frame pointer.
	size   uint16 // Size in bytes of this variable.
}

// irFunction is a function.
type irFunction struct {
	id     string // Name of function.
	offset uint64 // Offset of function relative to program data.
	args   map[string]irVar
	rets   map[string]irVar
	locals map[string]irVar
}

func (f irFunction) argsBytes() uint16 {
	var size uint16
	for _, arg := range f.args {
		size += arg.size
	}
	return size
}

func (f irFunction) retsBytes() uint16 {
	var size uint16
	for _, arg := range f.rets {
		size += arg.size
	}
	return size
}

func (f irFunction) localsBytes() uint16 {
	var size uint16
	for _, local := range f.locals {
		size += local.size
	}
	return size
}
