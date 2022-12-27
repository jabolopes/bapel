package vm

type OpCode = uint64

const (
	haltOpcode = OpCode(iota)
	callOpcode
	returnOpcode
	ifThenOpcode
	ifElseOpcode
	elseOpcode

	PopVarI8
	PopVarI16
	PopVarI32
	PopVarI64
)
