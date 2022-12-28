package vm

func unaryOpCode(base OpCode, mode OpMode, typ OpType) OpCode {
	return base + uint64(mode)*uint64(maxOpType) + uint64(typ)
}

func binaryOpCode(base OpCode, mode1, mode2 OpMode, typ OpType) OpCode {
	return base + uint64(mode1)*uint64(maxOpType)*uint64(maxOpMode) + uint64(mode2)*uint64(maxOpType) + uint64(typ)
}

func merge(ops *[]Op, m map[OpCode]func(*Machine) error) {
	for opcode, f := range m {
		if opcode >= uint64(len(*ops)) {
			delta := opcode - uint64(len(*ops)) + 1
			*ops = append(*ops, make([]Op, delta)...)
		}
		(*ops)[opcode] = Op{f}
	}
}

type OpTable struct {
	ops []Op
	// Unary opcodes.
	print OpCode
	push  OpCode
	pop   OpCode
	// Binary opcodes.
	add OpCode
}

func (t OpTable) Halt() OpCode   { return haltOpcode }
func (t OpTable) Call() OpCode   { return callOpcode }
func (t OpTable) Return() OpCode { return returnOpcode }
func (t OpTable) IfThen() OpCode { return ifThenOpcode }
func (t OpTable) IfElse() OpCode { return ifElseOpcode }
func (t OpTable) Else() OpCode   { return elseOpcode }

func (t OpTable) Add(mode1, mode2 OpMode, typ OpType) OpCode {
	return binaryOpCode(t.add, mode1, mode2, typ)
}

func (t OpTable) Print(mode OpMode, typ OpType) OpCode {
	return unaryOpCode(t.print, mode, typ)
}

func (t OpTable) Push(mode OpMode, typ OpType) OpCode {
	return unaryOpCode(t.push, mode, typ)
}

func (t OpTable) PopVar(typ OpType) OpCode {
	return unaryOpCode(t.pop, VarMode, typ)
}

func (t OpTable) PopDiscard(typ OpType) OpCode {
	return unaryOpCode(t.pop, StackMode, typ)
}

func NewOpTable() OpTable {
	table := OpTable{
		[]Op{
			haltOpcode:   {opHalt},
			callOpcode:   {opCall},
			returnOpcode: {opReturn},
			ifThenOpcode: {opIfThen},
			ifElseOpcode: {opIfElse},
			elseOpcode:   {opElse},
		},
		0, /* print */
		0, /* push */
		0, /* pop */
		0, /* add */
	}

	table.print = OpCode(len(table.ops))
	merge(&table.ops, opPrint(table.print))

	table.push = OpCode(len(table.ops))
	merge(&table.ops, opPush(table.push))

	table.pop = OpCode(len(table.ops))
	merge(&table.ops, opPop(table.pop))

	table.add = OpCode(len(table.ops))
	merge(&table.ops, opAdd(table.add))

	return table
}
