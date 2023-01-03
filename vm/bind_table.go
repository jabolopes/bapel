package vm

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

type opFunction = func(*Machine) error

type opFamilyMap = map[ir.OpCode]opFunction

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
	factories := []func(ir.OpCode) opFamilyMap{
		opHalt,
		opCall,
		opReturn,
		opIfThen,
		opIfElse,
		opElse,
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
	baseOpcodes := make([]ir.OpCode, len(factories))
	for family, factory := range factories {
		base := ir.OpCode(len(ops))
		baseOpcodes[family] = base
		merge(&ops, factory(base))
	}

	if got := ir.NewOpTable().Len(); len(ops) != got {
		panic(fmt.Errorf("Invalid bind table; expected table size %d; got %d", len(ops), got))
	}

	return bindTable{ops}
}
