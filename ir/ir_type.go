package ir

import (
	"fmt"
	"math"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type IrTypeCase int

const (
	AppType IrTypeCase = iota // Type application.
	ArrayType
	ForallType
	FunType
	LambdaType // Type abstraction, i.e., a function over types.
	NameType
	StructType
	TupleType
	VariantType
	VarType
)

func (c IrTypeCase) String() string {
	switch c {
	case AppType:
		return "type application"
	case ArrayType:
		return "array type"
	case ForallType:
		return "forall type"
	case FunType:
		return "function type"
	case LambdaType:
		return "lambda type"
	case NameType:
		return "typename"
	case StructType:
		return "struct type"
	case TupleType:
		return "tuple type"
	case VariantType:
		return "variant type"
	case VarType:
		return "type variable"
	default:
		panic(fmt.Errorf("unhandled %T %d", c, c))
	}
}

type appType struct {
	Fun IrType
	Arg IrType
}

func (t *appType) String() string {
	return fmt.Sprintf("%s %s", t.Fun, t.Arg)
}

type arrayType struct {
	ElemType IrType
	Size     int
}

func (t *arrayType) String() string {
	if t.Size == math.MaxInt {
		return fmt.Sprintf("[%v]", t.ElemType)
	}
	return fmt.Sprintf("[%v, %d]", t.ElemType, t.Size)
}

// Forall type.
//
// Example:
//
//	forall ['a] 'a -> 'a
//	forall ['a :: *] 'a -> 'a
type forallType struct {
	// Type variable. It is not prefixed with "'" when stored in this
	// field. When printed, the character "'" is prepended.
	Var  string
	Kind IrKind
	Type IrType
}

func (t *forallType) String() string {
	return fmt.Sprintf("forall ['%s] %s", t.Var, t.Type)
}

type functionType struct {
	Arg IrType
	Ret IrType
}

func (t *functionType) String() string {
	return fmt.Sprintf("%s -> %s", t.Arg, t.Ret)
}

type lambdaType struct {
	Var  string
	Kind IrKind
	Type IrType
}

func (t *lambdaType) String() string {
	return fmt.Sprintf("fun (%s %s) (%s)", t.Var, t.Kind, t.Type)
}

/* Struct type */

type StructField struct {
	ID   string
	Type IrType
}

func (f StructField) String() string {
	return fmt.Sprintf("%s %s", f.ID, f.Type)
}

type structType struct {
	Fields []StructField
}

func (t *structType) String() string {
	var b strings.Builder
	b.WriteString("{")
	if len(t.Fields) > 0 {
		b.WriteString(t.Fields[0].String())
		for _, field := range t.Fields[1:] {
			b.WriteString(fmt.Sprintf(", %s", field))
		}
	}
	b.WriteString("}")
	return b.String()
}

/* Tuple type */

type tupleType struct {
	Elems []IrType
}

func (t *tupleType) String() string {
	var b strings.Builder
	b.WriteString("(")
	if len(t.Elems) > 0 {
		b.WriteString(t.Elems[0].String())
		for _, typ := range t.Elems[1:] {
			b.WriteString(fmt.Sprintf(", %s", typ.String()))
		}
	}
	b.WriteString(")")
	return b.String()
}

/* Variant type */

type VariantTag struct {
	ID   string
	Type IrType
}

func (t VariantTag) String() string {
	return fmt.Sprintf("%s %s", t.ID, t.Type)
}

type variantType struct {
	Tags []VariantTag
}

func (t *variantType) String() string {
	var b strings.Builder
	b.WriteString("{|")
	if len(t.Tags) > 0 {
		b.WriteString(t.Tags[0].String())
		for _, typ := range t.Tags[1:] {
			b.WriteString(fmt.Sprintf(", %s", typ.String()))
		}
	}
	b.WriteString("|}")
	return b.String()
}

/* Type */

type IrType struct {
	Case    IrTypeCase
	App     *appType
	Array   *arrayType
	Forall  *forallType
	Fun     *functionType
	Lambda  *lambdaType
	Name    string // Typename, e.g., 'Hello'.
	Struct  *structType
	Tuple   *tupleType
	Variant *variantType
	Var     string // Type variable.

	// Position in source file.
	Pos Pos
}

func (t IrType) String() string {
	if t.Case == 0 && t.App == nil {
		return ""
	}

	switch t.Case {
	case AppType:
		return t.App.String()
	case ArrayType:
		return t.Array.String()
	case ForallType:
		return t.Forall.String()
	case FunType:
		return t.Fun.String()
	case LambdaType:
		return t.Lambda.String()
	case NameType:
		return t.Name
	case StructType:
		return t.Struct.String()
	case TupleType:
		return t.Tuple.String()
	case VariantType:
		return t.Variant.String()
	case VarType:
		return fmt.Sprintf("'%s", t.Var)
	default:
		panic(fmt.Errorf("unhandled IrType %d", t.Case))
	}
}

func (t IrType) Is(Case IrTypeCase) bool { return t.Case == Case }

func (t IrType) AppArgs() []IrType {
	if !t.Is(AppType) {
		return nil
	}

	return append(t.App.Fun.AppArgs(), t.App.Arg)
}

// Returns the type variables of a forall type (including immediate forall
// types).
//
// For example, for the following type:
//
//	forall 'a. (forall 'b. (forall 'c. 'a -> 'b -> 'c))
//
// This returns:
//
//	['a, 'b, 'c]
func (t IrType) ForallVars() []string {
	if !t.Is(ForallType) {
		return nil
	}

	return append([]string{t.Forall.Var}, t.Forall.Type.ForallVars()...)
}

// Returns the subtype of a forall type (including immediate forall types).
//
// For example, for the following type:
//
//	forall 'a. (forall 'b. (forall 'c. 'a -> 'b -> 'c))
//
// This returns:
//
//	'a -> 'b -> 'c
func (t IrType) ForallBody() IrType {
	if t.Is(ForallType) {
		return t.Forall.Type.ForallBody()
	}
	return t
}

// Same as ForallVars but for LambdaType instead of ForallType.
func (t IrType) LambdaVars() []string {
	if !t.Is(LambdaType) {
		return nil
	}

	return append([]string{t.Lambda.Var}, t.Lambda.Type.LambdaVars()...)
}

// Same as ForallBody but for LambdaType instead of ForallType.
func (t IrType) LambdaBody() IrType {
	if t.Is(LambdaType) {
		return t.Lambda.Type.LambdaBody()
	}
	return t
}

func (t IrType) Fields() []StructField {
	if t.Case != StructType {
		return nil
	}

	return t.Struct.Fields
}

func (t IrType) FieldByIndex(index int) (StructField, bool) {
	if index >= 0 && index < len(t.Fields()) {
		return t.Fields()[index], true
	}
	return StructField{}, false
}

func (t IrType) FieldByID(id string) (int, StructField, bool) {
	for index, field := range t.Fields() {
		if field.ID == id {
			return index, field, true
		}
	}
	return 0, StructField{}, false
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

func (t IrType) FieldByTerm(term IrTerm) (int, StructField, error) {
	switch {
	case term.Is(ConstTerm) && term.Const.Is(IntLiteral):
		index := int(*term.Const.Int)

		field, ok := t.FieldByIndex(index)
		if !ok {
			return 0, StructField{}, fmt.Errorf("field %d is not a valid field index of struct type %s", index, t)
		}
		return index, field, nil

	case term.Is(VarTerm):
		label := term.Var.ID

		index, field, ok := t.FieldByID(label)
		if !ok {
			return 0, StructField{}, fmt.Errorf("field %q is not a valid field label of struct type %s", label, t)
		}
		return index, field, nil

	default:
		return 0, StructField{}, fmt.Errorf("expected literal term (e.g., label, number) instead of %v", t)
	}
}

func (t IrType) Tags() []VariantTag {
	if t.Case != VariantType {
		return nil
	}

	return t.Variant.Tags
}

func (t IrType) TagByIndex(index int) (VariantTag, bool) {
	if index >= 0 && index < len(t.Tags()) {
		return t.Tags()[index], true
	}
	return VariantTag{}, false
}

func (t IrType) TagByID(id string) (int, VariantTag, bool) {
	for index, tag := range t.Tags() {
		if tag.ID == id {
			return index, tag, true
		}
	}
	return 0, VariantTag{}, false
}

func (t IrType) TagByTerm(term IrTerm) (int, VariantTag, error) {
	switch {
	case term.Is(ConstTerm) && term.Const.Is(IntLiteral):
		index := int(*term.Const.Int)

		tag, ok := t.TagByIndex(index)
		if !ok {
			return 0, VariantTag{}, fmt.Errorf("tag %d is not a valid tag index of variant type %s", index, t)
		}
		return index, tag, nil

	case term.Is(VarTerm):
		label := term.Var.ID

		index, tag, ok := t.TagByID(label)
		if !ok {
			return 0, VariantTag{}, fmt.Errorf("tag %q is not a valid tag label of variant type %s", label, t)
		}
		return index, tag, nil

	default:
		return 0, VariantTag{}, fmt.Errorf("expected literal term (e.g., label, number) instead of %v", t)
	}
}

func (t IrType) TagIDs() []string {
	ids := make([]string, len(t.Tags()))
	for i, tag := range t.Tags() {
		ids[i] = tag.ID
	}
	return ids
}

func (t IrType) TagTypes() []IrType {
	ids := make([]IrType, len(t.Tags()))
	for i, tag := range t.Tags() {
		ids[i] = tag.Type
	}
	return ids
}

func (t IrType) ElemByIndex(index int) (IrType, bool) {
	if !t.Is(TupleType) {
		return IrType{}, false
	}
	if index >= 0 && index < len(t.Tuple.Elems) {
		return t.Tuple.Elems[index], true
	}
	return IrType{}, false
}

func NewAppType(fun, arg IrType) IrType {
	return IrType{
		Case: AppType,
		App:  &appType{fun, arg},
	}
}

func NewArrayType(elemType IrType, size int) IrType {
	return IrType{
		Case:  ArrayType,
		Array: &arrayType{elemType, size},
	}
}

func NewForallType(tvar string, kind IrKind, typ IrType) IrType {
	return IrType{
		Case:   ForallType,
		Forall: &forallType{tvar, kind, typ},
	}
}

func NewFunctionType(arg, ret IrType) IrType {
	return IrType{
		Case: FunType,
		Fun:  &functionType{arg, ret},
	}
}

func NewLambdaType(tvar string, kind IrKind, body IrType) IrType {
	return IrType{
		Case:   LambdaType,
		Lambda: &lambdaType{tvar, kind, body},
	}
}

func NewNameType(name string) IrType {
	t := IrType{}
	t.Case = NameType
	t.Name = name
	return t
}

func NewStructType(fields []StructField) IrType {
	return IrType{
		Case:   StructType,
		Struct: &structType{fields},
	}
}

func NewTupleType(elems []IrType) IrType {
	if len(elems) == 1 {
		return elems[0]
	}

	return IrType{
		Case:  TupleType,
		Tuple: &tupleType{elems},
	}
}

func NewVariantType(tags []VariantTag) IrType {
	return IrType{
		Case:    VariantType,
		Variant: &variantType{tags},
	}
}

func NewVarType(tvar string) IrType {
	t := IrType{}
	t.Case = VarType
	t.Var = tvar
	return t
}

func getFreeTypeVars(t IrType, bound map[string]struct{}, free *map[VarKind]struct{}) {
	switch t.Case {
	case AppType:
		getFreeTypeVars(t.App.Fun, bound, free)
		getFreeTypeVars(t.App.Arg, bound, free)

	case ArrayType:
		getFreeTypeVars(t.Array.ElemType, bound, free)

	case ForallType:
		bound[t.Forall.Var] = struct{}{}
		getFreeTypeVars(t.Forall.Type, bound, free)

	case FunType:
		getFreeTypeVars(t.Fun.Arg, bound, free)
		getFreeTypeVars(t.Fun.Ret, bound, free)

	case LambdaType:
		bound[t.Lambda.Var] = struct{}{}
		getFreeTypeVars(t.Lambda.Type, bound, free)

	case NameType:
		return

	case StructType:
		for _, typ := range t.FieldTypes() {
			getFreeTypeVars(typ, bound, free)
		}

	case TupleType:
		for _, typ := range t.Tuple.Elems {
			getFreeTypeVars(typ, bound, free)
		}

	case VariantType:
		for _, typ := range t.FieldTypes() {
			getFreeTypeVars(typ, bound, free)
		}

	case VarType:
		if _, ok := bound[t.Var]; !ok {
			(*free)[VarKind{t.Var, NewTypeKind()}] = struct{}{}
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
	case AppType:
		return EqualsType(t1.App.Fun, t2.App.Fun) && EqualsType(t1.App.Arg, t2.App.Arg)
	case ArrayType:
		return EqualsType(t1.Array.ElemType, t2.Array.ElemType) && t1.Array.Size == t2.Array.Size
	case ForallType:
		return t1.Forall.Var == t2.Forall.Var && EqualsType(t1.Forall.Type, t2.Forall.Type)
	case FunType:
		return EqualsType(t1.Fun.Arg, t2.Fun.Arg) && EqualsType(t1.Fun.Ret, t2.Fun.Ret)
	case LambdaType:
		return t1.Lambda.Var == t2.Lambda.Var &&
			t1.Lambda.Kind == t2.Lambda.Kind &&
			EqualsType(t1.Lambda.Type, t2.Lambda.Type)
	case NameType:
		return t1.Name == t2.Name
	case StructType:
		return slices.EqualFunc(t1.Struct.Fields, t2.Struct.Fields, func(f1, f2 StructField) bool {
			return f1.ID == f2.ID && EqualsType(f1.Type, f2.Type)
		})
	case TupleType:
		return slices.EqualFunc(t1.Tuple.Elems, t2.Tuple.Elems, EqualsType)
	case VariantType:
		return slices.EqualFunc(t1.Variant.Tags, t2.Variant.Tags, func(f1, f2 VariantTag) bool {
			return f1.ID == f2.ID && EqualsType(f1.Type, f2.Type)
		})
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
	case AppType:
		return NewAppType(SubstituteType(t.App.Fun, source, target), SubstituteType(t.App.Arg, source, target))
	case ArrayType:
		return NewArrayType(SubstituteType(t.Array.ElemType, source, target), t.Array.Size)
	case ForallType:
		return NewForallType(t.Forall.Var, t.Forall.Kind, SubstituteType(t.Forall.Type, source, target))
	case FunType:
		return NewFunctionType(SubstituteType(t.Fun.Arg, source, target), SubstituteType(t.Fun.Ret, source, target))
	case LambdaType:
		return NewLambdaType(t.Lambda.Var, t.Lambda.Kind, SubstituteType(t.Lambda.Type, source, target))
	case NameType:
		return t

	case StructType:
		fields := make([]StructField, len(t.Struct.Fields))
		for i := range t.Struct.Fields {
			fields[i] = t.Struct.Fields[i]
			fields[i].Type = SubstituteType(fields[i].Type, source, target)
		}
		return NewStructType(fields)

	case TupleType:
		elems := make([]IrType, len(t.Tuple.Elems))
		for i := range t.Tuple.Elems {
			elems[i] = SubstituteType(t.Tuple.Elems[i], source, target)
		}
		return NewTupleType(elems)

	case VariantType:
		tags := make([]VariantTag, len(t.Variant.Tags))
		for i := range t.Variant.Tags {
			tags[i] = t.Variant.Tags[i]
			tags[i].Type = SubstituteType(tags[i].Type, source, target)
		}
		return NewVariantType(tags)

	case VarType:
		return t
	default:
		panic(fmt.Errorf("unhandled IrTypeCase %d", t.Case))
	}
}

func QuantifyType(typ IrType) IrType {
	free := map[VarKind]struct{}{}
	getFreeTypeVars(typ, map[string]struct{}{}, &free)
	return ForallVars(maps.Keys(free), typ)
}
