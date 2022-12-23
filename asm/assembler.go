package asm

import (
	"encoding/binary"

	"golang.org/x/exp/slices"
)

type Assembler struct {
	data []byte
}

func (a *Assembler) append(data []byte) *Assembler {
	a.data = append(a.data, data...)
	return a
}

func (a *Assembler) PutOpCode(opcode uint64) *Assembler {
	a.data = binary.AppendUvarint(a.data, opcode)
	return a
}

func (a *Assembler) PutI8(value byte) *Assembler {
	a.data = append(a.data, value)
	return a
}

func (a *Assembler) PutI16(value uint16) *Assembler {
	a.data = binary.LittleEndian.AppendUint16(a.data, value)
	return a
}

func (a *Assembler) PutI32(value uint32) *Assembler {
	a.data = binary.LittleEndian.AppendUint32(a.data, value)
	return a
}

func (a *Assembler) PutI64(value uint64) *Assembler {
	a.data = binary.LittleEndian.AppendUint64(a.data, value)
	return a
}

func (a *Assembler) Data() []byte {
	a.data = slices.Clip(a.data)
	return a.data
}

func (a *Assembler) Len() int {
	return len(a.data)
}

func NewAssembler() *Assembler { return &Assembler{nil} }
