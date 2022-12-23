package vm

import (
	"fmt"
)

type Op struct {
	callback func(*Machine) error
}

type Machine struct {
	ops     []Op
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
	return &Machine{
		[]Op{
			Halt: {opHalt},

			Call:   {opCall},
			Return: {opReturn},

			IfThen: {opIfThen},
			IfElse: {opIfElse},
			Else:   {opElse},

			StackAlloc: {opStackAlloc},

			PushI8:  {opPushImmediate[byte]()},
			PushI16: {opPushImmediate[uint16]()},
			PushI32: {opPushImmediate[uint32]()},
			PushI64: {opPushImmediate[uint64]()},

			PushLocalI8:  {opPushLocalI8},
			PushLocalI16: {opPushLocalI16},
			PushLocalI32: {opPushLocalI32},
			PushLocalI64: {opPushLocalI64},

			PopLocalI8:  {opPopLocalI8},
			PopLocalI16: {opPopLocalI16},
			PopLocalI32: {opPopLocalI32},
			PopLocalI64: {opPopLocalI64},

			PrintI8:  {opPrintI8},
			PrintI16: {opPrintI16},
			PrintI32: {opPrintI32},
			PrintI64: {opPrintI64},

			AddI8:  {opAddI8},
			AddI16: {opAddI16},
			AddI32: {opAddI32},
			AddI64: {opAddI64},
		},
		program,
		nil, /* stack */
		0,   /* pc */
		0,   /* fp */
	}
}
