package vm

import "encoding/binary"

// relative returns an address that is relative to the fp. The stack
// conceptually grows downwards, but the current implementation uses a
// byte array ([]byte) which actually grows upwards. As a result, args
// ad rets use negative offsets, and locals use positive offsets.
func relative(fp uint64, offset int64) uint64 {
	return fp + uint64(offset)
}

type Frame struct {
	stack []byte
	fp    uint64
}

func (f Frame) VarI8(offset int64) byte {
	return f.stack[relative(f.fp, offset)]
}

func (f Frame) VarI16(offset int64) uint16 {
	return binary.LittleEndian.Uint16(f.stack[relative(f.fp, offset):])
}

func (f Frame) VarI32(offset int64) uint32 {
	return binary.LittleEndian.Uint32(f.stack[relative(f.fp, offset):])
}

func (f Frame) VarI64(offset int64) uint64 {
	return binary.LittleEndian.Uint64(f.stack[relative(f.fp, offset):])
}

func (f Frame) SetVarI8(offset int64, value byte) Frame {
	f.stack[relative(f.fp, offset)] = value
	return f
}

func (f Frame) SetVarI16(offset int64, value uint16) Frame {
	binary.LittleEndian.PutUint16(f.stack[relative(f.fp, offset):], value)
	return f
}

func (f Frame) SetVarI32(offset int64, value uint32) Frame {
	binary.LittleEndian.PutUint32(f.stack[relative(f.fp, offset):], value)
	return f
}

func (f Frame) SetVarI64(offset int64, value uint64) Frame {
	binary.LittleEndian.PutUint64(f.stack[relative(f.fp, offset):], value)
	return f
}

// varPcI8 returns a variable's value, whose offset is pointed to by
// the pc.
func varPcI8(machine *Machine) uint8 {
	// Important: convert to int16 to ensure this is interpreted as a
	// negative number. This step needs to be done before extending to a
	// wider type (e.g., int64) otherwise the sign is lost.
	offset := int16(machine.Tape().GetI16())
	return machine.Frame().VarI8(int64(offset))
}

func varPcI16(machine *Machine) uint16 {
	// Important: convert to int16 to ensure this is interpreted as a
	// negative number. This step needs to be done before extending to a
	// wider type (e.g., int64) otherwise the sign is lost.
	offset := int16(machine.Tape().GetI16())
	return machine.Frame().VarI16(int64(offset))
}

func varPcI32(machine *Machine) uint32 {
	// Important: convert to int16 to ensure this is interpreted as a
	// negative number. This step needs to be done before extending to a
	// wider type (e.g., int64) otherwise the sign is lost.
	offset := int16(machine.Tape().GetI16())
	return machine.Frame().VarI32(int64(offset))
}

func varPcI64(machine *Machine) uint64 {
	// Important: convert to int16 to ensure this is interpreted as a
	// negative number. This step needs to be done before extending to a
	// wider type (e.g., int64) otherwise the sign is lost.
	offset := int16(machine.Tape().GetI16())
	return machine.Frame().VarI64(int64(offset))
}

// setVarPcI8 set a variable to a value. The offset of the variable is
// pointed to by the pc.
func setVarPcI8(machine *Machine, value uint8) {
	// Important: convert to int16 to ensure this is interpreted as a
	// negative number. This step needs to be done before extending to a
	// wider type (e.g., int64) otherwise the sign is lost.
	offset := int16(machine.Tape().GetI16())
	machine.Frame().SetVarI8(int64(offset), value)
}

func setVarPcI16(machine *Machine, value uint16) {
	// Important: convert to int16 to ensure this is interpreted as a
	// negative number. This step needs to be done before extending to a
	// wider type (e.g., int64) otherwise the sign is lost.
	offset := int16(machine.Tape().GetI16())
	machine.Frame().SetVarI16(int64(offset), value)
}

func setVarPcI32(machine *Machine, value uint32) {
	// Important: convert to int16 to ensure this is interpreted as a
	// negative number. This step needs to be done before extending to a
	// wider type (e.g., int64) otherwise the sign is lost.
	offset := int16(machine.Tape().GetI16())
	machine.Frame().SetVarI32(int64(offset), value)
}

func setVarPcI64(machine *Machine, value uint64) {
	// Important: convert to int16 to ensure this is interpreted as a
	// negative number. This step needs to be done before extending to a
	// wider type (e.g., int64) otherwise the sign is lost.
	offset := int16(machine.Tape().GetI16())
	machine.Frame().SetVarI64(int64(offset), value)
}
