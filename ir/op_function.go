package ir

// opVar is a variable (e.g., argument, return or local).
type opVar struct {
	offset uint16 // Offset in bytes relative to frame pointer.
	size   uint16 // Size in bytes of this variable.
}

// opFunction is a function.
type opFunction struct {
	id     string // Name of function.
	offset uint64 // Offset of function relative to program data.
	args   map[string]opVar
	rets   map[string]opVar
	locals map[string]opVar
}

func (f opFunction) argsBytes() uint16 {
	var size uint16
	for _, arg := range f.args {
		size += arg.size
	}
	return size
}

func (f opFunction) retsBytes() uint16 {
	var size uint16
	for _, arg := range f.rets {
		size += arg.size
	}
	return size
}

func (f opFunction) localsBytes() uint16 {
	var size uint16
	for _, local := range f.locals {
		size += local.size
	}
	return size
}
