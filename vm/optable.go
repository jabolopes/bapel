package vm

func unaryOpCode(base OpCode, mode OpMode, typ OpType) OpCode {
	return base + uint64(mode)*uint64(maxOpType) + uint64(typ)
}

func binaryOpCode(base OpCode, mode1, mode2 OpMode, typ OpType) OpCode {
	return base + uint64(mode1)*uint64(maxOpType)*uint64(maxOpMode) + uint64(mode2)*uint64(maxOpType) + uint64(typ)
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
	for opcode, f := range opPrint(table.print) {
		if opcode >= len(table.ops) {
			table.ops = append(table.ops, make([]Op, opcode-len(table.ops)+1)...)
		}
		table.ops[opcode] = Op{f}
	}

	table.push = OpCode(len(table.ops))
	for _, f := range opPush {
		table.ops = append(table.ops, Op{f})
	}

	table.pop = OpCode(len(table.ops))
	for _, f := range opPop {
		table.ops = append(table.ops, Op{f})
	}

	table.add = OpCode(len(table.ops))
	for opcode, f := range opAdd(table.add) {
		if opcode >= len(table.ops) {
			table.ops = append(table.ops, make([]Op, opcode-len(table.ops)+1)...)
		}
		table.ops[opcode] = Op{f}
	}

	return table
}
