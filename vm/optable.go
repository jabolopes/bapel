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
	opcodes []OpCode
	ops     []Op
}

func (t OpTable) Halt() OpCode   { return haltOpcode }
func (t OpTable) Call() OpCode   { return callOpcode }
func (t OpTable) Return() OpCode { return returnOpcode }
func (t OpTable) IfThen() OpCode { return ifThenOpcode }
func (t OpTable) IfElse() OpCode { return ifElseOpcode }
func (t OpTable) Else() OpCode   { return elseOpcode }

func (t OpTable) Add(mode1, mode2 OpMode, typ OpType) OpCode {
	return binaryOpCode(t.opcodes[addOpFamily], mode1, mode2, typ)
}

func (t OpTable) Print(mode OpMode, typ OpType) OpCode {
	return unaryOpCode(t.opcodes[printOpFamily], mode, typ)
}

func (t OpTable) Push(mode OpMode, typ OpType) OpCode {
	return unaryOpCode(t.opcodes[pushOpFamily], mode, typ)
}

func (t OpTable) Pop(mode OpMode, typ OpType) OpCode {
	return unaryOpCode(t.opcodes[popOpFamily], mode, typ)
}

func NewOpTable() OpTable {
	table := OpTable{
		make([]OpCode, maxOpFamily), /* opcodes */
		nil,                         /* ops */
	}

	opFactories := []func(OpCode) opFamilyMap{
		haltOpFamily:   opHalt,
		callOpFamily:   opCall,
		returnOpFamily: opReturn,
		ifThenOpFamily: opIfThen,
		ifElseOpFamily: opIfElse,
		elseOpFamily:   opElse,
		printOpFamily:  opPrint,
		pushOpFamily:   opPush,
		popOpFamily:    opPop,
		addOpFamily:    opAdd,
	}

	for opFamily, factory := range opFactories {
		base := OpCode(len(table.ops))
		table.opcodes[opFamily] = base
		merge(&table.ops, factory(base))
	}

	return table
}
