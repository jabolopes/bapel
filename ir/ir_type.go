package ir

import (
	"fmt"
)

type IrTypeCase int

const (
	ArrayType = IrTypeCase(iota)
	FunType
	IntType
	StructType
)

func (c IrTypeCase) String() string {
	switch c {
	case ArrayType:
		return "array"
	case FunType:
		return "function"
	case IntType:
		return "integer"
	case StructType:
		return "struct"
	default:
		panic(fmt.Errorf("Unhandled IrTypeCase %d", c))
	}
}

type IrType struct {
	Case       IrTypeCase
	ArrayType  *IrArrayType
	FunType    IrFunctionType
	IntType    IrIntType
	StructType IrStructType
}

func (t IrType) String() string {
	switch t.Case {
	case ArrayType:
		return t.ArrayType.String()
	case FunType:
		return t.FunType.String()
	case IntType:
		return t.IntType.String()
	case StructType:
		return t.StructType.String()
	default:
		panic(fmt.Errorf("Unhandled IR type %d", t.Case))
	}
}

func (t IrType) Is(Case IrTypeCase) bool { return t.Case == Case }

func MatchesType(formal, actual IrType, widen bool) error {
	if formal.Case != actual.Case {
		return fmt.Errorf("expected type %s; got %s", formal.Case, actual.Case)
	}

	switch formal.Case {
	case ArrayType:
		return MatchesArrayType(*formal.ArrayType, *actual.ArrayType, widen)
	case FunType:
		return MatchesFunctionType(formal.FunType, actual.FunType)
	case IntType:
		return MatchesIntType(formal.IntType, actual.IntType, widen)
	case StructType:
		return MatchesStructType(formal.StructType, actual.StructType, widen)
	default:
		panic(fmt.Errorf("Unhandled IrTypeCase %d", formal.Case))
	}
}

func SizeOfType(typ IrType) int {
	switch typ.Case {
	case ArrayType:
		return SizeOfArrayType(*typ.ArrayType)
	case FunType:
		return SizeOfIntType(I64)
	case IntType:
		return SizeOfIntType(typ.IntType)
	case StructType:
		return SizeOfStructType(typ.StructType)
	default:
		panic(fmt.Errorf("Unhandled IrTypeCase %d", typ.Case))
	}
}

func NewArrayType(typ IrArrayType) IrType {
	return IrType{ArrayType, &typ, IrFunctionType{}, 0, IrStructType{}}
}

func NewFunctionType(typ IrFunctionType) IrType {
	return IrType{FunType, nil, typ, 0, IrStructType{}}
}

func NewIntType(typ IrIntType) IrType {
	return IrType{IntType, nil, IrFunctionType{}, typ, IrStructType{}}
}

func NewStructType(typ IrStructType) IrType {
	return IrType{StructType, nil, IrFunctionType{}, 0, typ}
}
