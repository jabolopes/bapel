package vm

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

type Machine struct {
	bindTable bindTable
	program   ir.IrProgram
	stack     []byte
	pc        uint64 // Program counter. Offest in program.Data. Avoid slice since it needs to be incremented / decremented by n.
	fp        uint64 // Framepointer. Offset in stack. Avoid slice since stack can be reallocated.
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
		opcode, err := m.Tape().GetOpCode()
		if err != nil {
			return err
		}

		if opcode >= uint64(len(m.bindTable.ops)) {
			return fmt.Errorf("unknown opcode %d", opcode)
		}

		if err := m.bindTable.ops[opcode](m); err != nil {
			return err
		}
	}

	return nil
}

func New(program ir.IrProgram) *Machine {
	return &Machine{
		newBindTable(),
		program,
		nil, /* stack */
		0,   /* pc */
		0,   /* fp */
	}
}
