package vm

type OpCode = uint64

const (
	haltOpcode = OpCode(iota)
	callOpcode
	returnOpcode

	IfThen
	IfElse
	Else

	PopVarI8
	PopVarI16
	PopVarI32
	PopVarI64
)
