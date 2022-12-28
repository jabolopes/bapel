package vm

import "github.com/jabolopes/bapel/ir"

type OpType = ir.IrType

const (
	I8        = ir.I8
	I16       = ir.I16
	I32       = ir.I32
	I64       = ir.I64
	maxOpType = ir.I64 + 1
)
