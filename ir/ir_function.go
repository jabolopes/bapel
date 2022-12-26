package ir

// irFunction is a function.
type irFunction struct {
	id     string // Name of function.
	offset uint64 // Offset of function relative to program data.
	vars   map[string]IrVar
}

func (f *irFunction) lookupVar(id string) (IrVar, error) {
	irvar, ok := f.vars[id]
	if !ok {
		return IrVar{}, fmt.Errorf("Undefined variable %q", id)
	}
	return irvar, nil
}

func (f *irFunction) addVar(id string, irvar IrVar) error {
	if _, ok := f.vars[id]; ok {
		return fmt.Errorf("Variable %q already defined in this context", id)
	}

	f.vars[id] = irvar
	return nil
}

// varsSize returns the size in bytes of the given variable type.
func (f *irFunction) varsSize(typ IrVarType) (uint16, error) {
	var size uint16

	for _, irvar := range f.vars {
		if irvar.VarType == typ {
			varsize, err := SizeOfType(irvar.Type)
			if err != nil {
				return 0, err
			}

			size += uint16(varsize)
		}
	}

	return size, nil
}
