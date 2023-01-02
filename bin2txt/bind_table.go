package bin2txt

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

type opDecoder = func(*disassembler) error

type opFamilyMap = map[ir.OpCode]opDecoder

func merge(ops *[]opDecoder, m opFamilyMap) {
	for opcode, f := range m {
		if opcode >= uint64(len(*ops)) {
			delta := opcode - uint64(len(*ops)) + 1
			*ops = append(*ops, make([]opDecoder, delta)...)
		}
		(*ops)[opcode] = f
	}
}

type bindTable struct {
	ops []opDecoder
}

func newBindTable() bindTable {
	factories := []func(ir.OpCode) opFamilyMap{
		opHalt,
		opCall,
		opReturn,
		opIfThen,
		opIfElse,
		opElse,
		opPrintU,
		opPrintS,
		opPush,
		opPop,
		opNeg,
		opAdd,
	}

	var ops []opDecoder
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
