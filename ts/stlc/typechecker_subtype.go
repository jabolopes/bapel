package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

func (t *Typechecker) subtypeImpl(left, right ir.IrType) error {
	switch {
	case left.Case == ir.ArrayType && right.Case == ir.ArrayType:
		if err := t.subtype(left.Array.ElemType, right.Array.ElemType); err != nil {
			return fmt.Errorf("mismatch in array element types: %v", err)
		}

		if left.Array.Size != right.Array.Size {
			return fmt.Errorf("expected array with %d elements; got %d elements", left.Array.Size, right.Array.Size)
		}

		return nil

	case left.Case == ir.ForallType && right.Case == ir.ForallType:
		if len(left.Forall.Vars) != len(right.Forall.Vars) {
			return fmt.Errorf("expected forall type with %d variables (%v); got %d variables (%v)",
				len(left.Forall.Vars), left.Forall.Vars,
				len(right.Forall.Vars), right.Forall.Vars)
		}

		leftType := left.Forall.Type
		for i := range right.Forall.Vars {
			leftType = ir.SubstituteType(leftType, ir.NewVarType(right.Forall.Vars[i]), ir.NewVarType(right.Forall.Vars[i]))
		}

		return t.subtype(leftType, right.Forall.Type)

	// <:->
	case left.Case == ir.FunType && right.Case == ir.FunType:
		// B1 <: A1
		if err := t.subtype(right.Fun.Arg, left.Fun.Arg); err != nil {
			return err
		}

		// A2 <: B2
		if err := t.subtype(left.Fun.Ret, right.Fun.Ret); err != nil {
			return err
		}

		return nil

	case left.Case == ir.StructType && right.Case == ir.StructType:
		if len(left.Fields()) != len(right.Fields()) {
			return fmt.Errorf("expected %d fields; got %d", len(left.Fields()), len(right.Fields()))
		}

		for i := range left.Fields() {
			if left.Fields()[i].ID != right.Fields()[i].ID {
				return fmt.Errorf("expected field names %v; got %v", left.FieldIDs(), right.FieldIDs())
			}

			if err := t.subtype(left.Fields()[i].Type, right.Fields()[i].Type); err != nil {
				return err
			}
		}

		return nil

	case left.Case == ir.TupleType && right.Case == ir.TupleType:
		if len(left.Tuple) != len(right.Tuple) {
			return fmt.Errorf("expected %d elements; got %d", len(left.Tuple), len(right.Tuple))
		}

		for i := range left.Tuple {
			if err := t.subtype(left.Tuple[i], right.Tuple[i]); err != nil {
				return err
			}
		}

		return nil

	// <:Var
	case left.Case == ir.VarType && right.Case == ir.VarType && left.Var == right.Var:
		return nil

	// Typenames.
	case left.Case == ir.NameType && right.Case == ir.NameType && left.Name == right.Name:
		return nil

	default:
		return fmt.Errorf("expected type %s (%s); got %s (%s)", left.Case, left, right.Case, right)
	}
}

func (t *Typechecker) subtype(left, right ir.IrType) error {
	if err := t.subtypeImpl(left, right); err != nil {
		return fmt.Errorf("%s\n  subtyping %s and %s", err, left, right)
	}

	t.Printf("subtype: %s |- %s < %s", t.context.StringNoImports(), left, right)
	return nil
}
