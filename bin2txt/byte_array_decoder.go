package bin2txt

import (
	"encoding/binary"
	"errors"
)

type ByteArrayDecoder struct {
	data []byte
}

func (d *ByteArrayDecoder) GetOpCode() (uint64, error) {
	value, size := binary.Uvarint(d.data)
	if size <= 0 {
		return 0, errors.New("failed to decode opcode")
	}
	d.data = d.data[size:]
	return value, nil
}

func (d *ByteArrayDecoder) GetI8() byte {
	const size = 1
	value := d.data[0]
	d.data = d.data[size:]
	return value
}

func (d *ByteArrayDecoder) GetI16() uint16 {
	const size = 2
	value := binary.LittleEndian.Uint16(d.data[0:])
	d.data = d.data[size:]
	return value
}

func (d *ByteArrayDecoder) GetI32() uint32 {
	const size = 4
	value := binary.LittleEndian.Uint32(d.data[0:])
	d.data = d.data[size:]
	return value
}

func (d *ByteArrayDecoder) GetI64() uint64 {
	const size = 8
	value := binary.LittleEndian.Uint64(d.data[0:])
	d.data = d.data[size:]
	return value
}

func (d *ByteArrayDecoder) GetN(size uint64) []byte {
	value := d.data[:size]
	d.data = d.data[size:]
	return value
}

func (d *ByteArrayDecoder) Data() []byte {
	return d.data
}

func (d *ByteArrayDecoder) Len() int {
	return len(d.data)
}

func NewByteArrayDecoder(data []byte) *ByteArrayDecoder {
	return &ByteArrayDecoder{data}
}
