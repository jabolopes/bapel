package ir

import "fmt"

type IrTypeCase int

const (
	IntType = IrTypeCase(iota)
	FunType
)

func (c IrTypeCase) String() string {
	switch c {
	case IntType:
		return "integer"
	case FunType:
		return "function"
	default:
		panic(fmt.Errorf("Unhandled IrTypeCase %d", c))
	}
}

type IrType struct {
	Case    IrTypeCase
	IntType IrIntType
	FunType IrFunctionType
}

func (t IrType) String() string {
	switch t.Case {
	case IntType:
		return t.IntType.String()
	case FunType:
		return t.FunType.String()
	default:
		panic(fmt.Errorf("Unhandled IR type %d", t))
	}
}

func (t IrType) Is(Case IrTypeCase) bool { return t.Case == Case }

func MatchesType(formal, actual IrType, widen bool) error {
	if formal.Case != actual.Case {
		return fmt.Errorf("expected type %s; got %s", formal.Case, actual.Case)
	}

	switch formal.Case {
	case IntType:
		return MatchesIntType(formal.IntType, actual.IntType, widen)
	case FunType:
		return MatchesFunctionType(formal.FunType, actual.FunType)
	default:
		panic(fmt.Errorf("Unhandled IrTypeCase %d", formal.Case))
	}
}

func NewIntType(typ IrIntType) IrType {
	return IrType{IntType, typ, IrFunctionType{}}
}

func NewFunctionType(typ IrFunctionType) IrType {
	return IrType{FunType, 0, typ}
}
