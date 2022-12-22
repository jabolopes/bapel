package vm

import (
	"fmt"

	"github.com/jabolopes/go-vm/asm"
)

type Op struct {
	callback func(*Machine) error
}

type Machine struct {
	ops     []Op
	boot    []byte
	program OpProgram
	stack   []byte

	pc     int
	locals uint64 // Offset in stack. Avoid slice since stack can be reallocated.
}

type OpFunction struct {
	Locals uint16
	Offset []int
}

type OpProgram struct {
	Data      []byte
	Functions []OpFunction
}

func (m *Machine) Tape() Tape {
	return Tape{m.program.Data, &m.pc}
}

func (m *Machine) Stack() Stack {
	return Stack{&m.stack}
}

func (m *Machine) Run() error {
	for m.pc < len(m.program.Data) {
		tape := Tape{m.program.Data, &m.pc}
		opcode, err := tape.GetUvarint()
		if err != nil {
			return err
		}

		if opcode >= uint64(len(m.ops)) {
			return fmt.Errorf("Unknown opcode %d", opcode)
		}

		callback := m.ops[opcode].callback
		if callback == nil {
			return fmt.Errorf("Unimplemented opcode %d", opcode)
		}

		if err := m.ops[opcode].callback(m); err != nil {
			return err
		}
	}

	return nil
}

func New(program OpProgram) *Machine {
	var boot []byte
	{
		assembler := asm.New()
		assembler.PutOpCode(Call)
		assembler.PutU16(0)
		boot = assembler.Data()
	}

	return &Machine{
		[]Op{

			PushU8:   {opPushImmediate[byte]()},
			PushU16:  {opPushImmediate[uint16]()},
			PushU32:  {opPushImmediate[uint32]()},
			PushU64:  {opPushImmediate[uint64]()},
			PushL8:   {opPushLocalU8},
			PrintU8:  {opPrintU8},
			PrintU16: {opPrintU16},
			PrintU32: {opPrintU32},
			PrintU64: {opPrintU64},
		},
		boot,
		program,
		nil, /* stack */
		0,   /* pc */
		0,   /* locals */
	}
}
