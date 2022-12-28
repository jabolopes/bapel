package vm

import "github.com/jabolopes/bapel/ir"

func unaryOpCode(base OpCode, mode OpMode, typ OpType) OpCode {
	return ir.UnaryOpCode(base, mode, typ)
}

func binaryOpCode(base OpCode, mode1, mode2 OpMode, typ OpType) OpCode {
	return ir.BinaryOpCode(base, mode1, mode2, typ)
}

type OpTable = ir.OpTable

func NewOpTable() OpTable { return ir.NewOpTable() }
