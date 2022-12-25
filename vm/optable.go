package vm

type OpTable struct {
	ops   []Op
	print OpCode
	push  OpCode
}

func (t OpTable) unaryOpCode(base OpCode, mode OpMode, typ OpType) OpCode {
	return base + uint64(mode)*uint64(maxOpType) + uint64(typ)
}

func (t OpTable) Print(mode OpMode, typ OpType) OpCode {
	return t.unaryOpCode(t.print, mode, typ)
}

func (t OpTable) Push(mode OpMode, typ OpType) OpCode {
	return t.unaryOpCode(t.push, mode, typ)
}

func NewOpTable() OpTable {
	table := OpTable{
		[]Op{
			Halt: {opHalt},

			Call:   {opCall},
			Return: {opReturn},

			IfThen: {opIfThen},
			IfElse: {opIfElse},
			Else:   {opElse},

			StackAlloc: {opStackAlloc},

			PopLocalI8:  {opPopLocalI8},
			PopLocalI16: {opPopLocalI16},
			PopLocalI32: {opPopLocalI32},
			PopLocalI64: {opPopLocalI64},

			AddI8:  {opAddI8},
			AddI16: {opAddI16},
			AddI32: {opAddI32},
			AddI64: {opAddI64},
		},
		0, /* print */
		0, /* push */
	}

	table.print = OpCode(len(table.ops))
	for _, f := range opPrint {
		table.ops = append(table.ops, Op{f})
	}

	table.push = OpCode(len(table.ops))
	for _, f := range opPush {
		table.ops = append(table.ops, Op{f})
	}

	return table
}
