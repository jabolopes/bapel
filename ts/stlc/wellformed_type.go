package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
	"golang.org/x/exp/slices"
)

func containsDuplicates(ids []string) bool {
	slices.Sort(ids)
	for i := 0; i < len(ids)-1; i++ {
		if ids[i] == ids[i+1] {
			return true
		}
	}
	return false
}

func IsWellformedType(c Context, t ir.IrType) error {
	switch t.Case {
	case ir.AppType:
		if err := IsWellformedType(c, t.App.Fun); err != nil {
			return err
		}
		return IsWellformedType(c, t.App.Arg)

	case ir.ArrayType:
		return IsWellformedType(c, t.Array.ElemType)

	case ir.ForallType:
		var bodyType ir.IrType
		var err error
		c, bodyType, err = c.AddFreshType(t)
		if err != nil {
			return err
		}
		return IsWellformedType(c, bodyType)

	case ir.FunType:
		if err := IsWellformedType(c, t.Fun.Arg); err != nil {
			return err
		}
		return IsWellformedType(c, t.Fun.Ret)

	case ir.LambdaType:
		var bodyType ir.IrType
		var err error
		c, bodyType, err = c.AddFreshType(t)
		if err != nil {
			return err
		}
		return IsWellformedType(c, bodyType)

	case ir.NameType:
		if c.containsConstBind(t.Name) {
			return nil
		}
		// TODO: Validate that the aliased type is also a WellformedType in its own
		// context.
		if c.containsAliasBind(t.Name) {
			return nil
		}
		return fmt.Errorf("%q is undefined", t)

	case ir.StructType:
		if containsDuplicates(t.FieldIDs()) {
			return fmt.Errorf("struct type %v contains duplicate fields", t)
		}
		for _, typ := range t.FieldTypes() {
			if err := IsWellformedType(c, typ); err != nil {
				return err
			}
		}
		return nil

	case ir.TupleType:
		for _, typ := range t.Tuple.Elems {
			if err := IsWellformedType(c, typ); err != nil {
				return err
			}
		}
		return nil

	case ir.VariantType:
		if containsDuplicates(t.TagIDs()) {
			return fmt.Errorf("variant type %v contains duplicate tags", t)
		}
		for _, typ := range t.TagTypes() {
			if err := IsWellformedType(c, typ); err != nil {
				return err
			}
		}
		return nil

	case ir.VarType:
		if c.containsTypeVarBind(t.Var) {
			return nil
		}
		return fmt.Errorf("%q is undefined", t)

	default:
		panic(fmt.Errorf("unhandled %T %d", t.Case, t.Case))
	}
}
