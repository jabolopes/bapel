package vm

func merge(ops *[]Op, m map[OpCode]func(*Machine) error) {
	for opcode, f := range m {
		if opcode >= uint64(len(*ops)) {
			delta := opcode - uint64(len(*ops)) + 1
			*ops = append(*ops, make([]Op, delta)...)
		}
		(*ops)[opcode] = Op{f}
	}
}

type bindTable struct {
	ops []Op
}

func newBindTable() bindTable {
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

	var ops []Op
	baseOpcodes := make([]OpCode, maxOpFamily)
	for opFamily, factory := range opFactories {
		base := OpCode(len(ops))
		baseOpcodes[opFamily] = base
		merge(&ops, factory(base))
	}

	return bindTable{ops}
}
