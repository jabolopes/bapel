package vm

import "github.com/jabolopes/bapel/ir"

type opFunction = func(*Machine) error

type opFamilyMap = map[ir.OpCode]opFunction
