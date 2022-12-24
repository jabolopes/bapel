package ir

import "fmt"

type IrType int

const (
	I8 = IrType(iota)
	I16
	I32
	I64
)

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
		return 0, fmt.Errorf("Unhandled op type %q", arg)
	}
}
