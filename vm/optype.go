package vm

type OpType int

const (
	I8 = OpType(iota)
	I16
	I32
	I64
	maxOpType
)
