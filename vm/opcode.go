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

	PushI8
	PushI16
	PushI32
	PushI64

	PushLocalI8
	PushLocalI16
	PushLocalI32
	PushLocalI64

	PopLocalI8
	PopLocalI16
	PopLocalI32
	PopLocalI64

	AddI8
	AddI16
	AddI32
	AddI64
)
