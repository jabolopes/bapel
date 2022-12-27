package vm

type OpTable struct {
	ops   []Op
	add   OpCode
	print OpCode
	push  OpCode
	pop   OpCode
}

func (t OpTable) unaryOpCode(base OpCode, mode OpMode, typ OpType) OpCode {
	return base + uint64(mode)*uint64(maxOpType) + uint64(typ)
}

func (t OpTable) binaryOpCode(base OpCode, mode1, mode2 OpMode, typ OpType) OpCode {
	return base + uint64(mode1)*uint64(maxOpType)*uint64(maxOpMode) + uint64(mode2)*uint64(maxOpType) + uint64(typ)
}

func (t OpTable) Halt() OpCode   { return haltOpcode }
func (t OpTable) Call() OpCode   { return callOpcode }
func (t OpTable) Return() OpCode { return returnOpcode }
func (t OpTable) IfThen() OpCode { return ifThenOpcode }
func (t OpTable) IfElse() OpCode { return ifElseOpcode }
func (t OpTable) Else() OpCode   { return elseOpcode }

func (t OpTable) Add(mode1, mode2 OpMode, typ OpType) OpCode {
	return t.binaryOpCode(t.add, mode1, mode2, typ)
}

func (t OpTable) Print(mode OpMode, typ OpType) OpCode {
	return t.unaryOpCode(t.print, mode, typ)
}

func (t OpTable) Push(mode OpMode, typ OpType) OpCode {
	return t.unaryOpCode(t.push, mode, typ)
}

func (t OpTable) PopVar(typ OpType) OpCode {
	return t.unaryOpCode(t.pop, VarMode, typ)
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
		0, /* add */
		0, /* print */
		0, /* push */
		0, /* pop */
	}

	table.add = OpCode(len(table.ops))
	for _, f := range opAdd {
		table.ops = append(table.ops, Op{f})
	}

	table.print = OpCode(len(table.ops))
	for _, f := range opPrint {
		table.ops = append(table.ops, Op{f})
	}

	table.push = OpCode(len(table.ops))
	for _, f := range opPush {
		table.ops = append(table.ops, Op{f})
	}

	table.pop = OpCode(len(table.ops))
	for _, f := range opPop {
		table.ops = append(table.ops, Op{f})
	}

	return table
}
