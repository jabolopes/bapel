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
		panic(fmt.Errorf("unhandled IrTypeCase %d", c))
	}
}

type IrType struct {
	Case  IrTypeCase
	Array *struct {
		ElemType IrType
		Size     int
	}
	Fun *struct {
		Args []IrType
		Rets []IrType
	}
	Int    IrIntType
	Struct IrStructType
	Tuple  []IrType
	ID     string
}

func (t IrType) String() string {
	switch t.Case {
	case ArrayType:
		return fmt.Sprintf("[%v]", t.Array.ElemType)

	case FunType:
		var builder strings.Builder
		builder.WriteString(NewTupleType(t.Fun.Args).String())
		builder.WriteString(" -> ")
		builder.WriteString(NewTupleType(t.Fun.Rets).String())
		return builder.String()

	case IntType:
		return t.Int.String()

	case StructType:
		return t.Struct.String()

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
		return t.ID

	default:
		panic(fmt.Errorf("unhandled IrType %d", t.Case))
	}
}

func (t IrType) Is(Case IrTypeCase) bool { return t.Case == Case }

func NewArrayType(elemType IrType, size int) IrType {
	t := IrType{}
	t.Case = ArrayType
	t.Array = &struct {
		ElemType IrType
		Size     int
	}{elemType, size}
	return t
}

func NewFunctionType(args, rets []IrType) IrType {
	typ := IrType{}
	typ.Case = FunType
	typ.Fun = &struct {
		Args []IrType
		Rets []IrType
	}{args, rets}
	return typ
}

func NewIntType(intType IrIntType) IrType {
	typ := IrType{}
	typ.Case = IntType
	typ.Int = intType
	return typ
}

func NewStructType(structType IrStructType) IrType {
	typ := IrType{}
	typ.Case = StructType
	typ.Struct = structType
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
	typ.ID = idType
	return typ
}
