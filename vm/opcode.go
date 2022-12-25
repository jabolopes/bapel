package vm

type OpCode = uint64

const (
	Halt = OpCode(iota)

	Call
	Return

	IfThen
	IfElse
	Else

	StackAlloc

	PopLocalI8
	PopLocalI16
	PopLocalI32
	PopLocalI64

	AddI8
	AddI16
	AddI32
	AddI64
)
