package vm

import (
	"fmt"
)

type Op struct {
	callback func(*Machine) error
}

type Machine struct {
	optable OpTable
	program OpProgram
	stack   []byte

	pc uint64
	fp uint64 // Framepointer. Offset in stack. Avoid slice since stack can be reallocated.
}

type OpFunction struct {
	Locals uint64
	Offset uint64
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

func (m *Machine) Frame() Frame {
	return Frame{m.stack, m.fp}
}

func (m *Machine) Run() error {
	for m.pc < uint64(len(m.program.Data)) {
		opcode, err := m.Tape().GetUvarint()
		if err != nil {
			return err
		}

		if opcode >= uint64(len(m.optable.ops)) {
			return fmt.Errorf("Unknown opcode %d", opcode)
		}

		callback := m.optable.ops[opcode].callback
		if callback == nil {
			return fmt.Errorf("Unimplemented opcode %d", opcode)
		}

		if err := m.optable.ops[opcode].callback(m); err != nil {
			return err
		}
	}

	return nil
}

func New(program OpProgram) *Machine {
	return &Machine{
		NewOpTable(),
		program,
		nil, /* stack */
		0,   /* pc */
		0,   /* fp */
	}
}
