package vm

type OpTable struct {
	ops   []Op
	print OpCode
}

func (t OpTable) Print(mode OpMode, typ OpType) OpCode {
	return t.print + uint64(mode)*uint64(maxOpType) + uint64(typ)
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

			AddI8:  {opAddI8},
			AddI16: {opAddI16},
			AddI32: {opAddI32},
			AddI64: {opAddI64},
		},
		0, /* print */
	}

	table.print = OpCode(len(table.ops))
	for _, f := range opPrint {
		table.ops = append(table.ops, Op{f})
	}

	return table
}
