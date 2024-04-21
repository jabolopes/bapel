package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

func IsWellformedType(c Context, t ir.IrType) error {
	switch t.Case {
	case ir.AliasType:
		if err := IsWellformedType(c, t.Alias.Name); err != nil {
			return err
		}
		return IsWellformedType(c, t.Alias.Value)

	case ir.ArrayType:
		return IsWellformedType(c, t.Array.ElemType)

	case ir.ComponentType:
		return IsWellformedType(c, t.Component.ElemType)

	case ir.ForallType:
		c = c.Copy()
		for _, tvar := range t.Forall.Vars {
			c.binds = append(c.binds, NewDeclBind(DefSymbol, ir.NewTypeDecl(ir.NewVarType(tvar))))
		}
		return IsWellformedType(c, t.Forall.Type)

	case ir.FunType:
		if err := IsWellformedType(c, t.Fun.Arg); err != nil {
			return err
		}
		return IsWellformedType(c, t.Fun.Ret)

	case ir.NameType:
		if _, err := c.resolveTypeName(t); err == nil {
			return nil
		}

		_, err := c.getType(t)
		return err

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
		_, err := c.getType(t)
		return err

	default:
		panic(fmt.Errorf("unhandled %T %d", t.Case, t.Case))
	}
}
