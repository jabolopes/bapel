package vm

import (
	"fmt"
)

type opFunction = func(*Machine) error

type opFamilyMap = map[OpCode]opFunction

func merge(ops *[]opFunction, m opFamilyMap) {
	for opcode, f := range m {
		if opcode >= uint64(len(*ops)) {
			delta := opcode - uint64(len(*ops)) + 1
			*ops = append(*ops, make([]opFunction, delta)...)
		}
		(*ops)[opcode] = f
	}
}

type bindTable struct {
	ops []opFunction
}

func newBindTable() bindTable {
	factories := []func(OpCode) opFamilyMap{
		opHalt,
		opCall,
		opReturn,
		opIfThen,
		opIfElse,
		opElse,
		opSyscall,
		opWaitIO,
		opDoIO,
		opPrintU,
		opPrintS,
		opPush,
		opPop,
		opNeg,
		opAdd,
	}

	var ops []opFunction
	baseOpcodes := make([]OpCode, len(factories))
	for family, factory := range factories {
		base := OpCode(len(ops))
		baseOpcodes[family] = base
		merge(&ops, factory(base))
	}

	if got := NewOpTable().Len(); len(ops) != got {
		panic(fmt.Errorf("invalid bind table; expected table size %d; got %d", len(ops), got))
	}

	return bindTable{ops}
}
