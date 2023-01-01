package ir

import (
	"encoding/binary"

	"golang.org/x/exp/slices"
)

// ByteArrayEncoder is a low level code generator that writes bytecode
// into a byte array.
type ByteArrayEncoder struct {
	data []byte
}

func (a *ByteArrayEncoder) PutOpCode(opcode uint64) *ByteArrayEncoder {
	a.data = binary.AppendUvarint(a.data, opcode)
	return a
}

func (a *ByteArrayEncoder) PutI8(value byte) *ByteArrayEncoder {
	a.data = append(a.data, value)
	return a
}

func (a *ByteArrayEncoder) PutI16(value uint16) *ByteArrayEncoder {
	a.data = binary.LittleEndian.AppendUint16(a.data, value)
	return a
}

func (a *ByteArrayEncoder) PutI32(value uint32) *ByteArrayEncoder {
	a.data = binary.LittleEndian.AppendUint32(a.data, value)
	return a
}

func (a *ByteArrayEncoder) PutI64(value uint64) *ByteArrayEncoder {
	a.data = binary.LittleEndian.AppendUint64(a.data, value)
	return a
}

func (a *ByteArrayEncoder) PutN(data []byte) *ByteArrayEncoder {
	a.data = append(a.data, data...)
	return a
}

func (a *ByteArrayEncoder) RewriteI8(value uint8) *ByteArrayEncoder {
	const size = 1
	offset := len(a.data) - size
	a.data[offset] = value
	return a
}

func (a *ByteArrayEncoder) RewriteI16(value uint16) *ByteArrayEncoder {
	const size = 2
	offset := len(a.data) - size
	binary.LittleEndian.PutUint16(a.data[offset:], value)
	return a
}

func (a *ByteArrayEncoder) RewriteI32(value uint32) *ByteArrayEncoder {
	const size = 4
	offset := len(a.data) - size
	binary.LittleEndian.PutUint32(a.data[offset:], value)
	return a
}

func (a *ByteArrayEncoder) RewriteI64(value uint64) *ByteArrayEncoder {
	const size = 8
	offset := len(a.data) - size
	binary.LittleEndian.PutUint64(a.data[offset:], value)
	return a
}

func (a *ByteArrayEncoder) Data() []byte {
	a.data = slices.Clip(a.data)
	return a.data
}

func (a *ByteArrayEncoder) Len() int {
	return len(a.data)
}

func NewByteArrayEncoder() *ByteArrayEncoder {
	return &ByteArrayEncoder{nil}
}
