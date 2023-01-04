package bin2txt

import (
	"github.com/jabolopes/bapel/ir"
)

func opHalt(base ir.OpCode) opFamilyMap {
	return opFamilyMap{
		base: func(disassembler *disassembler) error {
			disassembler.printf("halt\n")
			return nil
		},
	}
}

func opCall(base ir.OpCode) opFamilyMap {
	return opFamilyMap{
		base: func(disassembler *disassembler) error {
			disassembler.printf("call %d\n", disassembler.dec().GetI64())
			return nil
		},
	}
}

func opReturn(base ir.OpCode) opFamilyMap {
	return opFamilyMap{
		base: func(disassembler *disassembler) error {
			disassembler.printf("return %d\n", disassembler.dec().GetI16())
			return nil
		},
	}
}

func opIfThen(base ir.OpCode) opFamilyMap {
	return opFamilyMap{
		base: func(disassembler *disassembler) error {
			disassembler.printf("ifThen %d\n", disassembler.dec().GetI64())
			return nil
		},
	}
}

func opIfElse(base ir.OpCode) opFamilyMap {
	return opFamilyMap{
		base: func(disassembler *disassembler) error {
			disassembler.printf("ifElse %d\n", disassembler.dec().GetI64())
			return nil
		},
	}
}

func opElse(base ir.OpCode) opFamilyMap {
	return opFamilyMap{
		base: func(disassembler *disassembler) error {
			disassembler.printf("else %d\n", disassembler.dec().GetI64())
			return nil
		},
	}
}

func opSyscall(base ir.OpCode) opFamilyMap {
	return opFamilyMap{
		base: func(disassembler *disassembler) error {
			disassembler.printf("syscall %d\n", disassembler.dec().GetI32())
			return nil
		},
	}
}

func opWaitIO(base ir.OpCode) opFamilyMap {
	return opFamilyMap{
		base: func(disassembler *disassembler) error {
			disassembler.printf("io.wait\n")
			return nil
		},
	}
}

func opDoIO(base ir.OpCode) opFamilyMap {
	return opFamilyMap{
		base: func(disassembler *disassembler) error {
			disassembler.printf("io.do %d\n", disassembler.dec().GetI16())
			return nil
		},
	}
}
