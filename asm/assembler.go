package asm

import (
	"encoding/binary"

	"golang.org/x/exp/slices"
)

type Assembler struct {
	data []byte
}

func (a *Assembler) PutOpCode(opcode uint64) {
	a.data = binary.AppendUvarint(a.data, opcode)
}

func (a *Assembler) PutU8(value byte) {
	a.data = append(a.data, value)
}

func (a *Assembler) PutU16(value uint16) {
	a.data = binary.LittleEndian.AppendUint16(a.data, value)
}

func (a *Assembler) PutU32(value uint32) {
	a.data = binary.LittleEndian.AppendUint32(a.data, value)
}

func (a *Assembler) PutU64(value uint64) {
	a.data = binary.LittleEndian.AppendUint64(a.data, value)
}

func (a *Assembler) Data() []byte {
	a.data = slices.Clip(a.data)
	return a.data
}

func New() *Assembler { return &Assembler{nil} }
