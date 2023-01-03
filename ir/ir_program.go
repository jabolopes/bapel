package ir

type Symbol struct {
	Id     string
	Offset uint64
}

type IrHeader struct {
	Symbols []Symbol
}

type IrProgram struct {
	Header IrHeader
	Data   []byte
}
