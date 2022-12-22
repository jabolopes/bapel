package vm

import "encoding/binary"

type Frame struct {
	stack []byte
	fp    uint64
}

func (t Frame) LocalI8(offset uint64) byte {
	return t.stack[t.fp+offset]
}

func (t Frame) LocalI16(offset uint64) uint16 {
	return binary.LittleEndian.Uint16(t.stack[t.fp+offset:])
}

func (t Frame) LocalI32(offset uint64) uint32 {
	return binary.LittleEndian.Uint32(t.stack[t.fp+offset:])
}

func (t Frame) LocalI64(offset uint64) uint64 {
	return binary.LittleEndian.Uint64(t.stack[t.fp+offset:])
}
