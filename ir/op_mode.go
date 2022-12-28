package ir

import "github.com/jabolopes/bapel/vm"

type OpMode = vm.OpMode

const (
	ImmediateMode = vm.ImmediateMode
	VarMode       = vm.VarMode
	StackMode     = vm.StackMode
	maxOpMode     = vm.StackMode + 1
)
