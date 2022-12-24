package ir

import (
	"encoding/binary"

	"golang.org/x/exp/slices"
)

// ByteAssembler is a low level assembler that assembles program bytes (words)
// into an array.
type ByteAssembler struct {
	data []byte
}

func (a *ByteAssembler) append(data []byte) *ByteAssembler {
	a.data = append(a.data, data...)
	return a
}

func (a *ByteAssembler) PutOpCode(opcode uint64) *ByteAssembler {
	a.data = binary.AppendUvarint(a.data, opcode)
	return a
}

func (a *ByteAssembler) PutI8(value byte) *ByteAssembler {
	a.data = append(a.data, value)
	return a
}

func (a *ByteAssembler) PutI16(value uint16) *ByteAssembler {
	a.data = binary.LittleEndian.AppendUint16(a.data, value)
	return a
}

func (a *ByteAssembler) PutI32(value uint32) *ByteAssembler {
	a.data = binary.LittleEndian.AppendUint32(a.data, value)
	return a
}

func (a *ByteAssembler) PutI64(value uint64) *ByteAssembler {
	a.data = binary.LittleEndian.AppendUint64(a.data, value)
	return a
}

func (a *ByteAssembler) Data() []byte {
	a.data = slices.Clip(a.data)
	return a.data
}

func (a *ByteAssembler) Len() int {
	return len(a.data)
}

func NewByteAssembler() *ByteAssembler {
	return &ByteAssembler{nil}
}
