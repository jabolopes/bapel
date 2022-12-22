package vm

import "encoding/binary"

type Stack struct {
	data *[]byte
}

func (s Stack) PushN(value []byte) {
	*s.data = append(*s.data, value...)
}

func (s Stack) PopU8() byte {
	last := len(*s.data) - 1
	value := (*s.data)[last]
	*s.data = (*s.data)[:last]
	return value
}

func (s Stack) PopU16() uint16 {
	const size = 2
	last := len(*s.data) - size
	value := binary.LittleEndian.Uint16((*s.data)[last:])
	*s.data = (*s.data)[:last]
	return value
}

func (s Stack) PopU32() uint32 {
	const size = 4
	last := len(*s.data) - size
	value := binary.LittleEndian.Uint32((*s.data)[last:])
	*s.data = (*s.data)[:last]
	return value
}

func (s Stack) PopU64() uint64 {
	const size = 8
	last := len(*s.data) - size
	value := binary.LittleEndian.Uint64((*s.data)[last:])
	*s.data = (*s.data)[:last]
	return value
}
