package vm

import (
	"encoding/binary"
	"errors"
)

type Tape struct {
	data []byte
	pc   *uint64
}

func (t Tape) GetOpCode() (uint64, error) {
	value, size := binary.Uvarint(t.data[*t.pc:])
	if size <= 0 {
		return 0, errors.New("failed to unsigned variable integer")
	}
	*t.pc += uint64(size)
	return value, nil
}

func (t Tape) GetI8() byte {
	const size = 1
	value := t.data[*t.pc:]
	*t.pc += size
	return value[0]
}

func (t Tape) GetI16() uint16 {
	const size = 2
	value := binary.LittleEndian.Uint16(t.data[*t.pc:])
	*t.pc += size
	return value
}

func (t Tape) GetI32() uint32 {
	const size = 4
	value := binary.LittleEndian.Uint32(t.data[*t.pc:])
	*t.pc += size
	return value
}

func (t Tape) GetI64() uint64 {
	const size = 8
	value := binary.LittleEndian.Uint64(t.data[*t.pc:])
	*t.pc += size
	return value
}
