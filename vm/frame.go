package vm

import "encoding/binary"

type Frame struct {
	stack []byte
	fp    uint64
}

func (f Frame) LocalI8(offset uint64) byte {
	return f.stack[f.fp+offset]
}

func (f Frame) LocalI16(offset uint64) uint16 {
	return binary.LittleEndian.Uint16(f.stack[f.fp+offset:])
}

func (f Frame) LocalI32(offset uint64) uint32 {
	return binary.LittleEndian.Uint32(f.stack[f.fp+offset:])
}

func (f Frame) LocalI64(offset uint64) uint64 {
	return binary.LittleEndian.Uint64(f.stack[f.fp+offset:])
}

func (f Frame) SetLocalI8(offset uint64, value byte) Frame {
	f.stack[f.fp+offset] = value
	return f
}

func (f Frame) SetLocalI16(offset uint64, value uint16) Frame {
	binary.LittleEndian.PutUint16(f.stack[f.fp+offset:], value)
	return f
}

func (f Frame) SetLocalI32(offset uint64, value uint32) Frame {
	binary.LittleEndian.PutUint32(f.stack[f.fp+offset:], value)
	return f
}

func (f Frame) SetLocalI64(offset uint64, value uint64) Frame {
	binary.LittleEndian.PutUint64(f.stack[f.fp+offset:], value)
	return f
}
