package vm

import (
	"encoding/binary"
	"fmt"
)

type Op struct {
	callback func(*Machine) error
}

type Machine struct {
	ops     []Op
	program []byte
	pc      []byte
	stack   []byte
}

func (m *Machine) Run() error {
	for len(m.pc) > 0 {
		opcode, size := binary.Uvarint(m.pc)
		if size <= 0 {
			return fmt.Errorf("Failed to read opcode from %q", m.pc)
		}

		if opcode >= uint64(len(m.ops)) {
			return fmt.Errorf("Unknown opcode %d", opcode)
		}

		m.pc = m.pc[size:]
		if err := m.ops[opcode].callback(m); err != nil {
			return err
		}
	}

	return nil
}

func New(program []byte) *Machine {
	return &Machine{
		[]Op{
			PushU8:   {opPushGeneric[byte]()},
			PushU16:  {opPushGeneric[uint16]()},
			PushU32:  {opPushGeneric[uint32]()},
			PushU64:  {opPushGeneric[uint64]()},
			PrintU8:  {opPrintU8},
			PrintU16: {opPrintU16},
			PrintU32: {opPrintU32},
			PrintU64: {opPrintU64},
		},
		program,
		program,
		nil, /* stack */
	}
}
