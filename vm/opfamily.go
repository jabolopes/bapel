package vm

type opFunction = func(*Machine) error

type opFamilyMap = map[OpCode]opFunction
