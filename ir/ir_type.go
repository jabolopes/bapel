package ir

import (
	"fmt"
	"strings"

	"golang.org/x/exp/maps"
)

type IrTypeCase int

const (
	ArrayType = IrTypeCase(iota)
	ForallType
	FunType
	IntType
	StructType
	TupleType
	VarType
	IDType
)

func (c IrTypeCase) String() string {
	switch c {
	case ArrayType:
		return "array"
	case ForallType:
		return "forall"
	case FunType:
		return "function"
	case IntType:
		return "integer"
	case StructType:
		return "struct"
	case TupleType:
		return "tuple"
	case VarType:
		return "type variable"
	case IDType:
		return "id"
	default:
		panic(fmt.Errorf("unhandled IrTypeCase %d", c))
	}
}

type StructField struct {
	ID   string
	Type IrType
}

func (f StructField) String() string {
	return fmt.Sprintf("%s %s", f.ID, f.Type)
}

type IrType struct {
	Case  IrTypeCase
	Array *struct {
		ElemType IrType
		Size     int
	}
	Forall *struct {
		// Type variables. Cannot be empty.
		Vars []string
		Type IrType
	}
	Fun *struct {
		Args []IrType
		Rets []IrType
	}
	Int    IrIntType
	Struct []StructField
	Var    string // Type variable.
	Tuple  []IrType
	ID     string
}

func (t IrType) String() string {
	switch t.Case {
	case ArrayType:
		return fmt.Sprintf("[%v]", t.Array.ElemType)

	case ForallType:
		var b strings.Builder
		b.WriteString("(")
		b.WriteString(fmt.Sprintf("'%s", t.Forall.Vars[0]))
		for _, tvar := range t.Forall.Vars[1:] {
			b.WriteString(fmt.Sprintf("'%s", tvar))
		}
		b.WriteString(") => ")
		b.WriteString(t.Forall.Type.String())
		return b.String()

	case FunType:
		var builder strings.Builder
		builder.WriteString(NewTupleType(t.Fun.Args).String())
		builder.WriteString(" -> ")
		builder.WriteString(NewTupleType(t.Fun.Rets).String())
		return builder.String()

	case IntType:
		return t.Int.String()

	case StructType:
		var b strings.Builder
		b.WriteString("{")
		if len(t.Struct) > 0 {
			b.WriteString(t.Struct[0].String())
			for _, field := range t.Struct[1:] {
				b.WriteString(fmt.Sprintf(", %s", field))
			}
		}
		b.WriteString("}")
		return b.String()

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

	case VarType:
		return fmt.Sprintf("'%s", t.Var)

	case IDType:
		return t.ID

	default:
		panic(fmt.Errorf("unhandled IrType %d", t.Case))
	}
}

func (t IrType) Fields() []StructField {
	if t.Case != StructType {
		return nil
	}

	return t.Struct
}

func (t IrType) FieldByIndex(index int) (StructField, bool) {
	if index >= 0 && index < len(t.Fields()) {
		return t.Fields()[index], true
	}
	return StructField{}, false
}

func (t IrType) FieldByID(id string) (StructField, bool) {
	for _, field := range t.Fields() {
		if field.ID == id {
			return field, true
		}
	}
	return StructField{}, false
}

func (t IrType) FieldIDs() []string {
	ids := make([]string, len(t.Fields()))
	for i, field := range t.Fields() {
		ids[i] = field.ID
	}
	return ids
}

func (t IrType) FieldTypes() []IrType {
	ids := make([]IrType, len(t.Fields()))
	for i, field := range t.Fields() {
		ids[i] = field.Type
	}
	return ids
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

func NewForallType(vars []string, typ IrType) IrType {
	if len(vars) == 0 {
		return typ
	}

	t := IrType{}
	t.Case = ForallType
	t.Forall = &struct {
		Vars []string
		Type IrType
	}{vars, typ}
	return t
}

func NewFunctionType(args, rets []IrType) IrType {
	t := IrType{}
	t.Case = FunType
	t.Fun = &struct {
		Args []IrType
		Rets []IrType
	}{args, rets}
	return t
}

func NewIntType(intType IrIntType) IrType {
	t := IrType{}
	t.Case = IntType
	t.Int = intType
	return t
}

func NewStructType(fields []StructField) IrType {
	t := IrType{}
	t.Case = StructType
	t.Struct = fields
	return t
}

func NewTupleType(tuple []IrType) IrType {
	if len(tuple) == 1 {
		return tuple[0]
	}

	t := IrType{}
	t.Case = TupleType
	t.Tuple = tuple
	return t
}

func NewVarType(tvar string) IrType {
	t := IrType{}
	t.Case = VarType
	t.Var = tvar
	return t
}

func NewIDType(idType string) IrType {
	t := IrType{}
	t.Case = IDType
	t.ID = idType
	return t
}

func IsMonotype(t IrType) bool {
	switch t.Case {
	case ArrayType:
		return IsMonotype(t.Array.ElemType)

	case ForallType:
		return false

	case FunType:
		return IsMonotype(NewTupleType(t.Fun.Args)) && IsMonotype(NewTupleType(t.Fun.Rets))

	case IntType:
		return true

	case StructType:
		return IsMonotype(NewTupleType(t.FieldTypes()))

	case TupleType:
		for _, typ := range t.Tuple {
			if !IsMonotype(typ) {
				return false
			}
		}
		return true

	case VarType:
		return true

	case IDType:
		// TODO: This doesn't look correct since a type ID can
		// theoretically refer to a polymorphic type.
		return true

	default:
		panic(fmt.Errorf("unhandled IrTypeCase %d", t.Case))
	}
}

func getFreeTypeVars(t IrType, bound map[string]struct{}, free *map[string]struct{}) {
	switch t.Case {
	case ArrayType:
		getFreeTypeVars(t.Array.ElemType, bound, free)

	case ForallType:
		for _, tvar := range t.Forall.Vars {
			bound[tvar] = struct{}{}
		}
		getFreeTypeVars(t.Forall.Type, bound, free)

	case FunType:
		for _, arg := range t.Fun.Args {
			getFreeTypeVars(arg, bound, free)
		}
		for _, ret := range t.Fun.Rets {
			getFreeTypeVars(ret, bound, free)
		}

	case IntType:
		return

	case StructType:
		for _, typ := range t.FieldTypes() {
			getFreeTypeVars(typ, bound, free)
		}

	case TupleType:
		for _, typ := range t.Tuple {
			getFreeTypeVars(typ, bound, free)
		}

	case VarType:
		if _, ok := bound[t.Var]; !ok {
			(*free)[t.Var] = struct{}{}
		}

	case IDType:
		// TODO: This doesn't look correct since a type ID can
		// theoretically refer to a polymorphic type.
		return

	default:
		panic(fmt.Errorf("unhandled IrTypeCase %d", t.Case))
	}
}

func QuantifyType(typ IrType) IrType {
	free := map[string]struct{}{}
	getFreeTypeVars(typ, map[string]struct{}{}, &free)
	return NewForallType(maps.Keys(free), typ)
}
