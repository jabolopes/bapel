package vm

type OpCode = uint64

const (
	Halt = OpCode(iota)

	callOpcode
	returnOpcode
	enterOpcode
	leaveOpcode

	IfThen
	IfElse
	Else

	StackAlloc

	PopVarI8
	PopVarI16
	PopVarI32
	PopVarI64

	AddI8
	AddI16
	AddI32
	AddI64
)
