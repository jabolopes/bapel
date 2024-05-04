package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

func IsWellformedType(c Context, t ir.IrType) error {
	switch t.Case {
	case ir.ArrayType:
		return IsWellformedType(c, t.Array.ElemType)

	case ir.ForallType:
		var err error
		c, err = c.AddBind(NewTypeVarBind(t.Forall.Var))
		if err != nil {
			return err
		}
		return IsWellformedType(c, t.Forall.Type)

	case ir.FunType:
		if err := IsWellformedType(c, t.Fun.Arg); err != nil {
			return err
		}
		return IsWellformedType(c, t.Fun.Ret)

	case ir.NameType:
		if c.ContainsNameBind(t.Name) {
			return nil
		}
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
		if c.ContainsVarType(t.Var) {
			return nil
		}
		return fmt.Errorf("%q is undefined", t)

	default:
		panic(fmt.Errorf("unhandled %T %d", t.Case, t.Case))
	}
}
