package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

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
		if c.ContainsNameBind(t.Name) {
			return nil
		}
		// TODO: Validate that the aliased type is also a WellformedType in its own
		// context.
		if c.ContainsAliasBind(t.Name) {
			return nil
		}
		return fmt.Errorf("%q is undefined", t)

	case ir.StructType:
		for _, typ := range t.FieldTypes() {
			if err := IsWellformedType(c, typ); err != nil {
				return err
			}
		}
		return nil

	case ir.TupleType:
		for _, typ := range t.Tuple {
			if err := IsWellformedType(c, typ); err != nil {
				return err
			}
		}
		return nil

	case ir.VarType:
		if c.ContainsTypeVarBind(t.Var) {
			return nil
		}
		return fmt.Errorf("%q is undefined", t)

	default:
		panic(fmt.Errorf("unhandled %T %d", t.Case, t.Case))
	}
}
