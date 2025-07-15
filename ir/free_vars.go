package ir

import (
	"fmt"

	"golang.org/x/exp/slices"
)

type freeVars struct {
	boundNames    map[string]struct{}
	boundTypeVars map[string]struct{}
	free          []IrType
}

func (v *freeVars) getFromType(typ IrType) {
	switch typ.Case {
	case AppType:
		v.getFromType(typ.App.Fun)
		v.getFromType(typ.App.Arg)

	case ArrayType:
		v.getFromType(typ.Array.ElemType)

	case ForallType:
		v.boundTypeVars[typ.Forall.Var] = struct{}{}
		v.getFromType(typ.Forall.Type)

	case FunType:
		v.getFromType(typ.Fun.Arg)
		v.getFromType(typ.Fun.Ret)

	case LambdaType:
		v.boundTypeVars[typ.Lambda.Var] = struct{}{}
		v.getFromType(typ.Lambda.Type)

	case NameType:
		if _, ok := v.boundNames[typ.Name]; !ok {
			v.boundNames[typ.Name] = struct{}{}
			v.free = append(v.free, typ)
		}

	case StructType:
		for _, typ := range typ.FieldTypes() {
			v.getFromType(typ)
		}

	case TupleType:
		for _, typ := range typ.Tuple.Elems {
			v.getFromType(typ)
		}

	case VariantType:
		for _, typ := range typ.FieldTypes() {
			v.getFromType(typ)
		}

	case VarType:
		if _, ok := v.boundTypeVars[typ.Var]; !ok {
			v.boundTypeVars[typ.Name] = struct{}{}
			v.free = append(v.free, typ)
		}

	default:
		panic(fmt.Errorf("unhandled %T %d", typ.Case, typ.Case))
	}
}

func newFreeVars() *freeVars {
	return &freeVars{
		map[string]struct{}{}, /* boundNames */
		map[string]struct{}{}, /* boundTypeVars */
		nil,                   /* free */
	}
}

func getFreeVarsFromType(typ IrType) []IrType {
	vars := newFreeVars()
	vars.getFromType(typ)
	return vars.free
}

// getFreeTypeVars returns the free type variables of `typ`. These
// includes only type variables and not types in general (e.g., name
// types).
func getFreeTypeVars(typ IrType) []VarKind {
	freeVars := getFreeVarsFromType(typ)

	var varKinds []VarKind
	for _, fvar := range freeVars {
		if !fvar.Is(VarType) {
			continue
		}

		varKinds = append(varKinds, VarKind{fvar.Var, NewTypeKind()})
	}
	slices.SortFunc(varKinds, CompareVarKind)

	return varKinds
}
