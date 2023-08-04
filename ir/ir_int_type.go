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
		panic(fmt.Errorf("Unhandled IR type %d", t))
	}
}

func ParseType(arg string) (IrIntType, error) {
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

func SizeOfType(typ IrIntType) (int, error) {
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

func MatchesIntType(formal, actual IrIntType, widen bool) error {
	if widen {
		if formal < actual {
			return fmt.Errorf("expected type %s or wider; got %s", formal, actual)
		}
	} else {
		if formal != actual {
			return fmt.Errorf("expected type %s; got %s", formal, actual)
		}
	}
	return nil
}
