package ir

import (
	"fmt"
)

type IrType int

const (
	I8 = IrType(iota)
	I16
	I32
	I64
	maxIrType
)

func (t IrType) String() string {
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
		panic(fmt.Errorf("Unhandled IR type %d", t))
	}
}

func ParseType(arg string) (IrType, error) {
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
		return 0, fmt.Errorf("Unhandled IR type %q", arg)
	}
}

func SizeOfType(typ IrType) (int, error) {
	switch typ {
	case I8:
		return 1, nil
	case I16:
		return 2, nil
	case I32:
		return 4, nil
	case I64:
		return 8, nil
	default:
		return 0, fmt.Errorf("Unhandled IR type %q", typ)
	}
}
