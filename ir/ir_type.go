package ir

import (
	"fmt"
	"strings"
)

type IrTypeCase int

const (
	ArrayType = IrTypeCase(iota)
	FunType
	IntType
	StructType
	TupleType
	IDType
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
	case TupleType:
		return "tuple"
	case IDType:
		return "id"
	default:
		panic(fmt.Errorf("Unhandled IrTypeCase %d", c))
	}
}

type IrType struct {
	Case       IrTypeCase
	Array      *IrArrayType
	Fun        IrFunctionType
	IntType    IrIntType
	StructType IrStructType
	Tuple      []IrType
	IDType     string
}

func (t IrType) String() string {
	switch t.Case {
	case ArrayType:
		return t.Array.String()
	case FunType:
		return t.Fun.String()
	case IntType:
		return t.IntType.String()
	case StructType:
		return t.StructType.String()
	case TupleType:
		tuple := t.Tuple
		var b strings.Builder
		b.WriteString("(")
		if len(tuple) > 0 {
			b.WriteString(tuple[0].String())
			for _, typ := range tuple[1:] {
				b.WriteString(fmt.Sprintf(", %s", typ.String()))
			}
		}
		b.WriteString(")")
		return b.String()
	case IDType:
		return t.IDType
	default:
		panic(fmt.Errorf("Unhandled IR type %d", t.Case))
	}
}

func (t IrType) Is(Case IrTypeCase) bool { return t.Case == Case }

func NewArrayType(array IrArrayType) IrType {
	typ := IrType{}
	typ.Case = ArrayType
	typ.Array = &array
	return typ
}

func NewFunctionType(fun IrFunctionType) IrType {
	typ := IrType{}
	typ.Case = FunType
	typ.Fun = fun
	return typ
}

func NewIntType(intType IrIntType) IrType {
	typ := IrType{}
	typ.Case = IntType
	typ.IntType = intType
	return typ
}

func NewStructType(structType IrStructType) IrType {
	typ := IrType{}
	typ.Case = StructType
	typ.StructType = structType
	return typ
}

func NewTupleType(tuple []IrType) IrType {
	if len(tuple) == 1 {
		return tuple[0]
	}

	typ := IrType{}
	typ.Case = TupleType
	typ.Tuple = tuple
	return typ
}

func NewIDType(idType string) IrType {
	typ := IrType{}
	typ.Case = IDType
	typ.IDType = idType
	return typ
}
