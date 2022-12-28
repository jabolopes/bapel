package vm

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
