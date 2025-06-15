package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

func (t *Typechecker) subtypeImpl(left, right ir.IrType) error {
	left = t.reduceType(left)
	right = t.reduceType(right)

	switch {
	case left.Is(ir.AppType) && right.Is(ir.AppType):
		if err := t.subtype(left.App.Fun, right.App.Fun); err != nil {
			return fmt.Errorf("mismatch in function types: %v", err)
		}
		if err := t.subtype(left.App.Arg, right.App.Arg); err != nil {
			return fmt.Errorf("mismatch in argument types: %v", err)
		}
		return nil

	case left.Is(ir.ArrayType) && right.Is(ir.ArrayType):
		if err := t.subtype(left.Array.ElemType, right.Array.ElemType); err != nil {
			return fmt.Errorf("mismatch in array element types: %v", err)
		}
		if left.Array.Size != right.Array.Size {
			return fmt.Errorf("expected array with %d elements; got %d elements", left.Array.Size, right.Array.Size)
		}
		return nil

	case left.Is(ir.ForallType) && right.Is(ir.ForallType):
		leftType := ir.SubstituteType(left.Forall.Type, ir.NewVarType(left.Forall.Var), ir.NewVarType(right.Forall.Var))
		return t.subtype(leftType, right.Forall.Type)

	// <:->
	case left.Is(ir.FunType) && right.Is(ir.FunType):
		// B1 <: A1
		if err := t.subtype(right.Fun.Arg, left.Fun.Arg); err != nil {
			return err
		}

		// A2 <: B2
		if err := t.subtype(left.Fun.Ret, right.Fun.Ret); err != nil {
			return err
		}

		return nil

	case left.Is(ir.StructType) && right.Is(ir.StructType):
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

	case left.Is(ir.TupleType) && right.Is(ir.TupleType):
		if len(left.Tuple.Elems) != len(right.Tuple.Elems) {
			return fmt.Errorf("expected %d elements; got %d", len(left.Tuple.Elems), len(right.Tuple.Elems))
		}

		for i := range left.Tuple.Elems {
			if err := t.subtype(left.Tuple.Elems[i], right.Tuple.Elems[i]); err != nil {
				return err
			}
		}

		return nil

	case left.Is(ir.VariantType) && right.Is(ir.VariantType):
		if len(left.Tags()) != len(right.Tags()) {
			return fmt.Errorf("expected %d tags; got %d", len(left.Tags()), len(right.Tags()))
		}

		for i := range left.Tags() {
			if left.Tags()[i].ID != right.Tags()[i].ID {
				return fmt.Errorf("expected tag names %v; got %v", left.TagIDs(), right.TagIDs())
			}

			if err := t.subtype(left.Tags()[i].Type, right.Tags()[i].Type); err != nil {
				return err
			}
		}

		return nil

	// <:Var
	case left.Is(ir.VarType) && right.Is(ir.VarType) && left.Var == right.Var:
		return nil

	// Typenames.
	case left.Is(ir.NameType) && t.context.containsConstBind(left.Name) &&
		right.Is(ir.NameType) && t.context.containsConstBind(right.Name) &&
		left.Name == right.Name:
		return nil

	case left.Is(ir.NameType) && t.context.containsAliasBind(left.Name):
		bind, err := t.context.getAliasBind(left.Name)
		if err != nil {
			panic(err)
		}
		return t.subtype(bind.Alias.Type, right)

	case right.Is(ir.NameType) && t.context.containsAliasBind(right.Name):
		bind, err := t.context.getAliasBind(right.Name)
		if err != nil {
			panic(err)
		}
		return t.subtype(left, bind.Alias.Type)

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
