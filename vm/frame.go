package vm

import "encoding/binary"

// relative returns an address that is relative to the fp. The stack
// conceptually grows downwards, but the current implementation uses a
// byte array ([]byte) which actually grows upwards, therefore, the
// offsets are subtracted (- offset) instead of added (+ offset).
func relative(fp, offset uint64) uint64 {
	return fp - offset
}

type Frame struct {
	stack []byte
	fp    uint64
}

func (f Frame) VarI8(offset uint64) byte {
	return f.stack[relative(f.fp, offset)]
}

func (f Frame) VarI16(offset uint64) uint16 {
	return binary.LittleEndian.Uint16(f.stack[relative(f.fp, offset):])
}

func (f Frame) VarI32(offset uint64) uint32 {
	return binary.LittleEndian.Uint32(f.stack[relative(f.fp, offset):])
}

func (f Frame) VarI64(offset uint64) uint64 {
	return binary.LittleEndian.Uint64(f.stack[relative(f.fp, offset):])
}

func (f Frame) SetVarI8(offset uint64, value byte) Frame {
	f.stack[relative(f.fp, offset)] = value
	return f
}

func (f Frame) SetVarI16(offset uint64, value uint16) Frame {
	binary.LittleEndian.PutUint16(f.stack[relative(f.fp, offset):], value)
	return f
}

func (f Frame) SetVarI32(offset uint64, value uint32) Frame {
	binary.LittleEndian.PutUint32(f.stack[relative(f.fp, offset):], value)
	return f
}

func (f Frame) SetVarI64(offset uint64, value uint64) Frame {
	binary.LittleEndian.PutUint64(f.stack[relative(f.fp, offset):], value)
	return f
}
