package vm

type OpCode = uint64

const (
	haltOpcode = OpCode(iota)
	callOpcode
	returnOpcode
	ifThenOpcode
	ifElseOpcode
	elseOpcode
)

type opFamily = uint64

const (
	haltOpFamily = opFamily(iota)
	callOpFamily
	returnOpFamily
	ifThenOpFamily
	ifElseOpFamily
	elseOpFamily
	// Unary ops.
	printOpFamily
	pushOpFamily
	popOpFamily
	// Binary ops.
	addOpFamily
	maxOpFamily
)

type opFunction = func(*Machine) error

type opFamilyMap = map[OpCode]opFunction
