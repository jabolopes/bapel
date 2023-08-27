package ir

import (
	"fmt"
)

type IrIntType int

const (
	I8 = IrIntType(iota)
	I16
	I32
	I64
	maxIrIntType
)

func (t IrIntType) String() string {
	switch t {
	case I8:
		return "i8"
	case I16:
		return "i16"
	case I32:
		return "i32"
	case I64:
		return "i64"
	default:
		panic(fmt.Errorf("unhandled IrIntType %d", t))
	}
}

func ParseIntType(arg string) (IrIntType, error) {
	switch arg {
	case "i8":
		return I8, nil
	case "i16":
		return I16, nil
	case "i32":
		return I32, nil
	case "i64":
		return I64, nil
	default:
		return 0, fmt.Errorf("unhandled IrIntType %q", arg)
	}
}
