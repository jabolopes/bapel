package ir

import (
	"fmt"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type IrTypeCase int

const (
	AliasType IrTypeCase = iota
	ArrayType
	ComponentType
	ForallType
	FunType
	NameType
	StructType
	TupleType
	VarType
)

func (c IrTypeCase) String() string {
	switch c {
	case AliasType:
		return "alias"
	case ArrayType:
		return "array"
	case ComponentType:
		return "component"
	case ForallType:
		return "forall"
	case FunType:
		return "function"
	case NameType:
		return "typename"
	case StructType:
		return "struct"
	case TupleType:
		return "tuple"
	case VarType:
		return "type variable"
	default:
		panic(fmt.Errorf("unhandled %T %d", c, c))
	}
}

type componentType struct {
	Name     string
	ElemType IrType
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
	Alias *struct {
		Name  IrType
		Value IrType
	}
	Array *struct {
		ElemType IrType
		Size     int
	}
	Component *componentType
	Forall    *struct {
		// Type variables. Cannot be empty.
		Vars []string
		Type IrType
	}
	Fun *struct {
		Arg IrType
		Ret IrType
	}
	Name   string // Typename, e.g., 'Hello'.
	Struct []StructField
	Var    string // Type variable.
	Tuple  []IrType
}

func (t IrType) String() string {
	if t.Case == 0 && t.Alias == nil {
		return ""
	}

	switch t.Case {
	case AliasType:
		return fmt.Sprintf("%s = %s", t.Alias.Name, t.Alias.Value)
	case ArrayType:
		return fmt.Sprintf("[%v]", t.Array.ElemType)
	case ComponentType:
		c := *t.Component
		return fmt.Sprintf("component %s<%s>", c.Name, c.ElemType)

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
	default:
		panic(fmt.Errorf("unhandled IrType %d", t.Case))
	}
}

func (t IrType) ID() (string, bool) {
	switch t.Case {
	case AliasType:
		return t.Alias.Name.ID()
	case ComponentType:
		return t.Component.Name, true
	case ArrayType, ForallType, FunType, StructType, TupleType:
		return "", false
	case NameType:
		return t.Name, true
	case VarType:
		return t.Var, true
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

func NewAliasType(name, value IrType) IrType {
	return IrType{
		Case: AliasType,
		Alias: &struct {
			Name  IrType
			Value IrType
		}{name, value},
	}
}

func NewArrayType(elemType IrType, size int) IrType {
	t := IrType{}
	t.Case = ArrayType
	t.Array = &struct {
		ElemType IrType
		Size     int
	}{elemType, size}
	return t
}

func NewComponentType(id string, elemType IrType) IrType {
	return IrType{
		Case:      ComponentType,
		Component: &componentType{id, elemType},
	}
}

func NewForallType(vars []string, typ IrType) IrType {
	if len(vars) == 0 {
		return typ
	}

	slices.Sort(vars)
	vars = slices.Compact(vars)

	if typ.Case == ForallType {
		vars = append(vars, typ.Forall.Vars...)
		return NewForallType(vars, typ.Forall.Type)
	}

	return IrType{
		Case: ForallType,
		Forall: &struct {
			Vars []string
			Type IrType
		}{vars, typ},
	}
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

func getFreeTypeVars(t IrType, bound map[string]struct{}, free *map[string]struct{}) {
	switch t.Case {
	case AliasType:
		getFreeTypeVars(t.Alias.Name, bound, free)
		getFreeTypeVars(t.Alias.Value, bound, free)

	case ArrayType:
		getFreeTypeVars(t.Array.ElemType, bound, free)

	case ComponentType:
		getFreeTypeVars(t.Component.ElemType, bound, free)

	case ForallType:
		for _, tvar := range t.Forall.Vars {
			bound[tvar] = struct{}{}
		}
		getFreeTypeVars(t.Forall.Type, bound, free)

	case FunType:
		getFreeTypeVars(t.Fun.Arg, bound, free)
		getFreeTypeVars(t.Fun.Ret, bound, free)

	case NameType:
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

	default:
		panic(fmt.Errorf("unhandled IrTypeCase %d", t.Case))
	}
}

func EqualsType(t1, t2 IrType) bool {
	if t1.Case != t2.Case {
		return false
	}

	switch t1.Case {
	case AliasType:
		return EqualsType(t1.Alias.Name, t2.Alias.Name) &&
			EqualsType(t1.Alias.Value, t2.Alias.Value)
	case ArrayType:
		return EqualsType(t1.Array.ElemType, t2.Array.ElemType) && t1.Array.Size == t2.Array.Size
	case ComponentType:
		return t1.Component.Name == t2.Component.Name && EqualsType(t1.Component.ElemType, t2.Component.ElemType)
	case ForallType:
		return slices.Equal(t1.Forall.Vars, t2.Forall.Vars) && EqualsType(t1.Forall.Type, t2.Forall.Type)
	case FunType:
		return EqualsType(t1.Fun.Arg, t2.Fun.Arg) && EqualsType(t1.Fun.Ret, t2.Fun.Ret)
	case NameType:
		return t1.Name == t2.Name

	case StructType:
		return slices.EqualFunc(t1.Struct, t2.Struct, func(f1, f2 StructField) bool {
			return f1.ID == f2.ID && EqualsType(f1.Type, f2.Type)
		})

	case TupleType:
		return slices.EqualFunc(t1.Tuple, t2.Tuple, EqualsType)
	case VarType:
		return t1.Var == t2.Var
	default:
		panic(fmt.Errorf("unhandled IrTypeCase %d", t1.Case))
	}
}

func SubstituteType(t, source, target IrType) IrType {
	if EqualsType(t, source) {
		return target
	}

	switch t.Case {
	case AliasType:
		return NewAliasType(
			SubstituteType(t.Alias.Name, source, target),
			SubstituteType(t.Alias.Value, source, target))
	case ArrayType:
		return NewArrayType(SubstituteType(t.Array.ElemType, source, target), t.Array.Size)
	case ComponentType:
		return NewComponentType(t.Component.Name, SubstituteType(t.Component.ElemType, source, target))
	case ForallType:
		return NewForallType(t.Forall.Vars, SubstituteType(t.Forall.Type, source, target))
	case FunType:
		return NewFunctionType(SubstituteType(t.Fun.Arg, source, target), SubstituteType(t.Fun.Ret, source, target))
	case NameType:
		return t

	case StructType:
		fields := make([]StructField, len(t.Struct))
		for i := range t.Struct {
			fields[i] = t.Struct[i]
			fields[i].Type = SubstituteType(fields[i].Type, source, target)
		}
		return NewStructType(fields)

	case TupleType:
		elems := make([]IrType, len(t.Tuple))
		for i := range t.Tuple {
			elems[i] = SubstituteType(t.Tuple[i], source, target)
		}
		return NewTupleType(elems)

	case VarType:
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
