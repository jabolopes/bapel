package vm

import "github.com/jabolopes/bapel/ir"

type OpMode = ir.OpMode

const (
	ImmediateMode = ir.ImmediateMode
	VarMode       = ir.VarMode
	StackMode     = ir.StackMode
	maxOpMode     = ir.StackMode + 1
)
