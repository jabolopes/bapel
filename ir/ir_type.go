package ir

import (
	"fmt"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type IrTypeCase int

const (
	ArrayType = IrTypeCase(iota)
	ForallType
	FunType
	InstanceType
	NameType
	NumberType
	StructType
	TupleType
	VarType
	VarExistType
)

func (c IrTypeCase) String() string {
	switch c {
	case ArrayType:
		return "array"
	case ForallType:
		return "forall"
	case FunType:
		return "function"
	case InstanceType:
		return "instance"
	case NameType:
		return "typename"
	case NumberType:
		return "number"
	case StructType:
		return "struct"
	case TupleType:
		return "tuple"
	case VarType:
		return "type variable"
	case VarExistType:
		return "existential type variable"
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
		Arg IrType
		Ret IrType
	}
	Name     string // Typename, e.g., 'Hello'.
	Struct   []StructField
	Var      string // Type variable.
	VarExist *struct {
		Var string
	}
	Tuple []IrType
}

func (t IrType) String() string {
	switch t.Case {
	case ArrayType:
		return fmt.Sprintf("[%v]", t.Array.ElemType)

	case ForallType:
		var b strings.Builder
		b.WriteString(fmt.Sprintf("'%s", t.Forall.Vars[0]))
		for _, tvar := range t.Forall.Vars[1:] {
			b.WriteString(fmt.Sprintf(" '%s", tvar))
		}
		b.WriteString(" => ")
		b.WriteString(t.Forall.Type.String())
		return b.String()

	case FunType:
		return fmt.Sprintf("%s -> %s", t.Fun.Arg, t.Fun.Ret)
	case NameType:
		return t.Name
	case NumberType:
		return "Number"

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
	case VarExistType:
		return fmt.Sprintf("^%s", t.VarExist.Var)
	default:
		panic(fmt.Errorf("unhandled IrType %d", t.Case))
	}
}

// TODO: Should be called ID() but ID is already a field.
func (t IrType) TypeID() string {
	switch t.Case {
	case ArrayType:
		return ""
	case ForallType:
		return ""
	case FunType:
		return ""
	case NameType:
		return t.Name
	case NumberType:
		return ""
	case StructType:
		return ""
	case TupleType:
		return ""
	case VarType:
		return t.Var
	case VarExistType:
		return t.VarExist.Var
	default:
		panic(fmt.Errorf("unhandled IrTypeCase %d", t.Case))
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

func NewFunctionType(arg, ret IrType) IrType {
	t := IrType{}
	t.Case = FunType
	t.Fun = &struct {
		Arg IrType
		Ret IrType
	}{arg, ret}
	return t
}

func NewNameType(name string) IrType {
	t := IrType{}
	t.Case = NameType
	t.Name = name
	return t
}

func NewNumberType() IrType {
	t := IrType{}
	t.Case = NumberType
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

func NewVarExistType(tvar string) IrType {
	t := IrType{}
	t.Case = VarExistType
	t.VarExist = &struct {
		Var string
	}{tvar}
	return t
}

func IsMonotype(t IrType) bool {
	switch t.Case {
	case ArrayType:
		return IsMonotype(t.Array.ElemType)
	case ForallType:
		return false
	case FunType:
		return IsMonotype(t.Fun.Arg) && IsMonotype(t.Fun.Ret)
	case NameType:
		// TODO: This doesn't look correct since a type ID can
		// theoretically refer to a polymorphic type.
		return true
	case NumberType:
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
	case VarExistType:
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
		getFreeTypeVars(t.Fun.Arg, bound, free)
		getFreeTypeVars(t.Fun.Ret, bound, free)
	case NameType, NumberType:
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

	case VarExistType:
		if _, ok := bound[t.VarExist.Var]; !ok {
			(*free)[t.VarExist.Var] = struct{}{}
		}

	default:
		panic(fmt.Errorf("unhandled IrTypeCase %d", t.Case))
	}
}

func equalsType(t1, t2 IrType) bool {
	if t1.Case != t2.Case {
		return false
	}

	switch t1.Case {
	case ArrayType:
		return equalsType(t1.Array.ElemType, t2.Array.ElemType) && t1.Array.Size == t2.Array.Size
	case ForallType:
		return slices.Equal(t1.Forall.Vars, t2.Forall.Vars) && equalsType(t1.Forall.Type, t2.Forall.Type)
	case FunType:
		return equalsType(t1.Fun.Arg, t2.Fun.Arg) && equalsType(t1.Fun.Ret, t2.Fun.Ret)
	case NameType:
		return t1.Name == t2.Name
	case NumberType:
		return true

	case StructType:
		return slices.EqualFunc(t1.Struct, t2.Struct, func(f1, f2 StructField) bool {
			return f1.ID == f2.ID && equalsType(f1.Type, f2.Type)
		})

	case TupleType:
		return slices.EqualFunc(t1.Tuple, t2.Tuple, equalsType)
	case VarType:
		return t1.Var == t2.Var
	case VarExistType:
		return t1.VarExist.Var == t2.VarExist.Var
	default:
		panic(fmt.Errorf("unhandled IrTypeCase %d", t1.Case))
	}
}

func substituteType(t, source, target IrType) IrType {
	if equalsType(t, source) {
		return target
	}

	switch t.Case {
	case ArrayType:
		return NewArrayType(substituteType(t.Array.ElemType, source, target), t.Array.Size)
	case ForallType:
		return NewForallType(t.Forall.Vars, substituteType(t.Forall.Type, source, target))
	case FunType:
		return NewFunctionType(substituteType(t.Fun.Arg, source, target), substituteType(t.Fun.Ret, source, target))
	case NameType, NumberType:
		return t

	case StructType:
		fields := make([]StructField, len(t.Struct))
		for i := range t.Struct {
			fields[i] = t.Struct[i]
			fields[i].Type = substituteType(fields[i].Type, source, target)
		}
		return NewStructType(fields)

	case TupleType:
		elems := make([]IrType, len(t.Tuple))
		for i := range t.Tuple {
			elems[i] = substituteType(t.Tuple[i], source, target)
		}
		return NewTupleType(elems)

	case VarType:
		return t
	case VarExistType:
		return t
	default:
		panic(fmt.Errorf("unhandled IrTypeCase %d", t.Case))
	}
}

func QuantifyType(typ IrType) IrType {
	free := map[string]struct{}{}
	getFreeTypeVars(typ, map[string]struct{}{}, &free)
	return NewForallType(maps.Keys(free), typ)
}
