package ir

import (
	"encoding/binary"

	"golang.org/x/exp/slices"
)

// ByteGenerator is a low level code geneator that program bytecode
// into an array.
type ByteGenerator struct {
	data []byte
}

func (a *ByteGenerator) PutOpCode(opcode uint64) *ByteGenerator {
	a.data = binary.AppendUvarint(a.data, opcode)
	return a
}

func (a *ByteGenerator) PutI8(value byte) *ByteGenerator {
	a.data = append(a.data, value)
	return a
}

func (a *ByteGenerator) PutI16(value uint16) *ByteGenerator {
	a.data = binary.LittleEndian.AppendUint16(a.data, value)
	return a
}

func (a *ByteGenerator) PutI32(value uint32) *ByteGenerator {
	a.data = binary.LittleEndian.AppendUint32(a.data, value)
	return a
}

func (a *ByteGenerator) PutI64(value uint64) *ByteGenerator {
	a.data = binary.LittleEndian.AppendUint64(a.data, value)
	return a
}

func (a *ByteGenerator) PutN(data []byte) *ByteGenerator {
	a.data = append(a.data, data...)
	return a
}

func (a *ByteGenerator) RewriteI8(value uint8) *ByteGenerator {
	const size = 1
	offset := len(a.data) - size
	a.data[offset] = value
	return a
}

func (a *ByteGenerator) RewriteI16(value uint16) *ByteGenerator {
	const size = 2
	offset := len(a.data) - size
	binary.LittleEndian.PutUint16(a.data[offset:], value)
	return a
}

func (a *ByteGenerator) RewriteI32(value uint32) *ByteGenerator {
	const size = 4
	offset := len(a.data) - size
	binary.LittleEndian.PutUint32(a.data[offset:], value)
	return a
}

func (a *ByteGenerator) RewriteI64(value uint64) *ByteGenerator {
	const size = 8
	offset := len(a.data) - size
	binary.LittleEndian.PutUint64(a.data[offset:], value)
	return a
}

func (a *ByteGenerator) Data() []byte {
	a.data = slices.Clip(a.data)
	return a.data
}

func (a *ByteGenerator) Len() int {
	return len(a.data)
}

func NewByteGenerator() *ByteGenerator {
	return &ByteGenerator{nil}
}
