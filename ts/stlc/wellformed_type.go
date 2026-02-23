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

func isWellformedType(c Context, t ir.IrType) error {
	switch {
	case t.Is(ir.AppType):
		if err := isWellformedType(c, t.App.Fun); err != nil {
			return err
		}
		return isWellformedType(c, t.App.Arg)

	case t.Is(ir.ArrayType):
		return isWellformedType(c, t.Array.ElemType)

	case t.Is(ir.ExistVarType):
		return nil

	case t.Is(ir.ForallType):
		var bodyType ir.IrType
		var err error
		c, _, bodyType, err = c.AddFreshType(t)
		if err != nil {
			return err
		}
		return isWellformedType(c, bodyType)

	case t.Is(ir.FunType):
		if err := isWellformedType(c, t.Fun.Arg); err != nil {
			return err
		}
		return isWellformedType(c, t.Fun.Ret)

	case t.Is(ir.LambdaType):
		var bodyType ir.IrType
		var err error
		c, _, bodyType, err = c.AddFreshType(t)
		if err != nil {
			return err
		}
		return isWellformedType(c, bodyType)

	case t.Is(ir.NameType) && c.containsAliasBind(t.Name):
		return nil

	case t.Is(ir.NameType) && c.containsConstBind(t.Name):
		return nil

	case t.Is(ir.NameType):
		return fmt.Errorf("type %q is undefined", t.Name)

	case t.Is(ir.StructType):
		if containsDuplicates(t.FieldIDs()) {
			return fmt.Errorf("struct type %v contains duplicate fields", t)
		}
		for _, typ := range t.FieldTypes() {
			if err := isWellformedType(c, typ); err != nil {
				return err
			}
		}
		return nil

	case t.Is(ir.TupleType):
		for _, typ := range t.Tuple.Elems {
			if err := isWellformedType(c, typ); err != nil {
				return err
			}
		}
		return nil

	case t.Is(ir.VariantType):
		if containsDuplicates(t.TagIDs()) {
			return fmt.Errorf("variant type %v contains duplicate tags", t)
		}
		for _, typ := range t.TagTypes() {
			if err := isWellformedType(c, typ); err != nil {
				return err
			}
		}
		return nil

	case t.Is(ir.VarType) && c.containsTypeVarBind(t.Var):
		return nil

	case t.Is(ir.VarType):
		return fmt.Errorf("%q is undefined", t)

	default:
		panic(fmt.Errorf("unhandled %T %d", t.Case, t.Case))
	}
}
