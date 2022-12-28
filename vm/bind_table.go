package vm

import "fmt"

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
	factories := []func(OpCode) opFamilyMap{
		opHalt,
		opCall,
		opReturn,
		opIfThen,
		opIfElse,
		opElse,
		opPrint,
		opPush,
		opPop,
		opAdd,
	}

	var ops []Op
	baseOpcodes := make([]OpCode, len(factories))
	for family, factory := range factories {
		base := OpCode(len(ops))
		baseOpcodes[family] = base
		merge(&ops, factory(base))
	}

	if got := NewOpTable().Len(); len(ops) != got {
		panic(fmt.Errorf("Invalid bind table; expected table size %d; got %d", len(ops), got))
	}

	return bindTable{ops}
}
