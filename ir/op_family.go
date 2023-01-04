package ir

type opFamily = uint64

const (
	haltOpFamily = opFamily(iota)
	callOpFamily
	returnOpFamily
	ifThenOpFamily
	ifElseOpFamily
	elseOpFamily
	// Syscall ops.
	syscallOpFamily
	// IO ops.
	ioWaitOpFamily
	ioDoOpFamily
	// Unary ops.
	printUOpFamily
	printSOpFamily
	pushOpFamily
	popOpFamily
	negOpFamily
	// Binary ops.
	addOpFamily
	maxOpFamily
)
