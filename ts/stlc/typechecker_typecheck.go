package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

func (t *Typechecker) typecheckAppTypeTerm(term *ir.IrTerm) error {
	if !term.Is(ir.AppTypeTerm) {
		panic(fmt.Errorf("expected %T %d", ir.AppTypeTerm, ir.AppTypeTerm))
	}

	c := term.AppType

	if err := t.typecheck(&c.Fun); err != nil {
		return err
	}

	if !c.Fun.Type.Is(ir.ForallType) {
		return fmt.Errorf("expected term %v to have forall type instead of %v", c.Fun, c.Fun.Type)
	}
	funType := c.Fun.Type.Forall

	if err := isWellformedType(t.context, c.Arg); err != nil {
		return err
	}

	argKind, err := inferKind(t.context, term.AppType.Arg)
	if err != nil {
		return err
	}

	if !ir.EqualsKind(funType.Kind, argKind) {
		return fmt.Errorf("expected argument in type application (%s) to match forall type's kind (%s)", argKind, funType.Kind)
	}

	typ := ir.SubstituteType(funType.Type, ir.NewVarType(funType.Var), c.Arg)
	term.Type = &typ
	return nil
}

func (t *Typechecker) typecheckAppTermTerm(term *ir.IrTerm) error {
	if !term.Is(ir.AppTermTerm) {
		panic(fmt.Errorf("expected %T %d", ir.AppTermTerm, ir.AppTermTerm))
	}

	c := term.AppTerm

	if err := t.typecheck(&c.Fun); err != nil {
		return err
	}

	if err := t.typecheck(&c.Arg); err != nil {
		return err
	}

	if !c.Fun.Type.Is(ir.FunType) {
		return fmt.Errorf("expected term %v to have function type instead of %v", c.Fun, c.Fun.Type)
	}
	funType := c.Fun.Type.Fun

	if err := t.subtype(*c.Arg.Type, funType.Arg); err != nil {
		return err
	}

	term.Type = &funType.Ret
	return nil
}

func (t *Typechecker) typecheckIndexGetTerm(term *ir.IrTerm) error {
	c := term.IndexGet
	if err := t.typecheck(&c.Obj); err != nil {
		return err
	}

	var index *int64
	var label *string
	switch c.Index.Case {
	// Get field by index.
	//
	// Example:
	//   Index.get x 0
	case ir.ConstTerm:
		index = &c.Index.Const.Number

	// Get field by label.
	//
	// Example:
	//   Index.get x myfield
	case ir.VarTerm:
		label = &c.Index.Var.ID
	}

	objType := *c.Obj.Type
	switch {
	case objType.Is(ir.ArrayType) && index != nil:
		if *index < 0 || *index >= int64(objType.Array.Size) {
			return fmt.Errorf("index %d is out of bounds", *index)
		}

		term.Type = &objType.Array.ElemType
		return nil

	case objType.Is(ir.ArrayType):
		if err := t.typecheck(&c.Index); err != nil {
			return err
		}

		if err := t.isNumber(*c.Index.Type); err != nil {
			return err
		}

		term.Type = &objType.Array.ElemType
		return nil

	case objType.Is(ir.StructType) && index != nil:
		field, ok := objType.FieldByIndex(int(*index))
		if !ok {
			return fmt.Errorf("field %d is not a valid field of struct type %s", *index, objType)
		}

		c.Field = field.ID
		term.Type = &field.Type
		return nil

	case objType.Is(ir.StructType) && label != nil:
		field, ok := objType.FieldByID(*label)
		if !ok {
			return fmt.Errorf("field %q is not a valid field of struct type %s", *label, objType)
		}

		c.Field = field.ID
		term.Type = &field.Type
		return nil

	case objType.Is(ir.StructType):
		return fmt.Errorf("expected field identifier or number literal to index struct type %s", objType)

	case objType.Is(ir.TupleType) && index != nil:
		elem, ok := objType.ElemByIndex(int(*index))
		if !ok {
			return fmt.Errorf("index %d is not a valid element of tuple type %s", *index, objType)
		}

		term.Type = &elem
		return nil

	case objType.Is(ir.TupleType):
		return fmt.Errorf("expected number literal to index tuple type %s", objType)

	case objType.Is(ir.VariantType) && index != nil:
		tag, ok := objType.TagByIndex(int(*index))
		if !ok {
			return fmt.Errorf("tag %d is not a valid tag of variant type %s", *index, objType)
		}

		term.Type = &tag.Type
		return nil

	case objType.Is(ir.VariantType) && label != nil:
		index, tag, ok := objType.TagByID(*label)
		if !ok {
			return fmt.Errorf("tag %q is not a valid tag of variant type %s", *label, objType)
		}

		term.IndexGet.Index = ir.Number(int64(index))
		term.Type = &tag.Type
		return nil

	case objType.Is(ir.VariantType):
		return fmt.Errorf("expected tag identifier or number literal to index variant type %s", objType)

	default:
		return fmt.Errorf("expected indexable type (e.g., array, struct, etc); got %s", objType)
	}

	return nil
}

func (t *Typechecker) typecheckIndexSetTerm(term *ir.IrTerm) error {
	c := term.IndexSet

	var index *int64
	var label *string
	switch c.Index.Case {
	// Set field by index.
	//
	// Example:
	//   Index.set x 0 value
	case ir.ConstTerm:
		index = &c.Index.Const.Number

	// Set field by label.
	//
	// Example:
	//   Index.set x myfield value
	case ir.VarTerm:
		label = &c.Index.Var.ID
	}

	if err := t.typecheck(&c.Obj); err != nil {
		return err
	}

	objType := *c.Obj.Type
	switch {
	case objType.Is(ir.ArrayType) && index != nil:
		if *index < 0 || *index >= int64(objType.Array.Size) {
			return fmt.Errorf("index %d is out of bounds", *index)
		}

		term.Type = &objType.Array.ElemType
		return nil

	case objType.Is(ir.ArrayType):
		if err := t.typecheck(&c.Index); err != nil {
			return err
		}

		if err := t.isNumber(*c.Index.Type); err != nil {
			return err
		}

		if err := t.typecheck(&c.Value); err != nil {
			return err
		}

		if err := t.subtype(objType.Array.ElemType, *c.Value.Type); err != nil {
			return err
		}

		typ := ir.NewTupleType(nil)
		term.Type = &typ
		return nil

	case objType.Is(ir.StructType) && index != nil:
		field, ok := objType.FieldByIndex(int(*index))
		if !ok {
			return fmt.Errorf("field %d is not a valid field of struct type %s", *index, objType)
		}

		c.Field = field.ID
		term.Type = &field.Type
		return nil

	case objType.Is(ir.StructType) && label != nil:
		field, ok := objType.FieldByID(*label)
		if !ok {
			return fmt.Errorf("field %q is not a valid field of struct type %s", *label, objType)
		}

		c.Field = field.ID
		term.Type = &field.Type
		return nil

	case objType.Is(ir.StructType):
		return fmt.Errorf("expected field identifier or number literal to index struct type %s", objType)

	case objType.Is(ir.TupleType) && index != nil:
		elem, ok := objType.ElemByIndex(int(*index))
		if !ok {
			return fmt.Errorf("index %d is not a valid element of tuple type %s", *index, objType)
		}

		term.Type = &elem
		return nil

	case objType.Is(ir.TupleType):
		return fmt.Errorf("expected number literal to index tuple type %s", objType)

	case objType.Is(ir.VariantType) && index != nil:
		tag, ok := objType.TagByIndex(int(*index))
		if !ok {
			return fmt.Errorf("tag %d is not a valid tag of variant type %s", *index, objType)
		}

		tagIndex := int(*index)
		c.TagIndex = &tagIndex
		term.Type = &tag.Type
		return nil

	case objType.Is(ir.VariantType) && label != nil:
		index, tag, ok := objType.TagByID(*label)
		if !ok {
			return fmt.Errorf("tag %q is not a valid tag of variant type %s", *label, objType)
		}

		c.TagIndex = &index
		term.Type = &tag.Type
		return nil

	case objType.Is(ir.VariantType):
		return fmt.Errorf("expected tag identifier or number literal to index variant type %s", objType)

	default:
		return fmt.Errorf("expected indexable type (e.g., array); got %s", objType)
	}

}

func (t *Typechecker) typecheckImpl(term *ir.IrTerm) error {
	switch {
	case term.Is(ir.AppTypeTerm):
		return t.typecheckAppTypeTerm(term)

	case term.Is(ir.AppTermTerm):
		return t.typecheckAppTermTerm(term)

	case term.Is(ir.AssignTerm):
		c := term.Assign
		if err := t.withBindPosition(func() error {
			return t.typecheck(&c.Ret)
		}); err != nil {
			return err
		}

		if err := t.typecheck(&c.Arg); err != nil {
			return err
		}

		if err := t.subtype(*c.Ret.Type, *c.Arg.Type); err != nil {
			return err
		}

		term.Type = c.Ret.Type
		return nil

	case term.Is(ir.BlockTerm):
		c := term.Block
		for i := range c.Terms {
			if err := t.typecheck(&c.Terms[i]); err != nil {
				return err
			}
		}

		typ := ir.NewTupleType(nil)
		term.Type = &typ
		return nil

	case term.Is(ir.ConstTerm) && t.bindPosition:
		return fmt.Errorf("expected symbol declared as %s; got number literal", ir.TermDecl)

	case term.Is(ir.ConstTerm) && term.Type != nil:
		kind, err := inferKind(t.context, *term.Type)
		if err != nil {
			return err
		}

		if !ir.EqualsKind(kind, ir.NewTypeKind()) {
			return fmt.Errorf("expected %v with type %v and kind %v to have kind %v", term, *term.Type, kind, ir.NewTypeKind())
		}

		return nil

	case term.Is(ir.ConstTerm):
		typ := ir.Forall("a", ir.NewTypeKind(), ir.Tvar("a"))
		term.Type = &typ
		return nil

	case term.Is(ir.IfTerm):
		c := term.If

		if err := t.typecheck(&c.Condition); err != nil {
			return err
		}

		if err := t.isNumber(*c.Condition.Type); err != nil {
			return err
		}

		if err := t.typecheck(&c.Then); err != nil {
			return err
		}

		if c.Else != nil {
			if err := t.typecheck(c.Else); err != nil {
				return err
			}

			if err := t.subtype(*c.Then.Type, *c.Else.Type); err != nil {
				return err
			}
		}

		term.Type = c.Then.Type
		return nil

	case term.Is(ir.IndexGetTerm):
		return t.typecheckIndexGetTerm(term)

	case term.Is(ir.IndexSetTerm):
		return t.typecheckIndexSetTerm(term)

	case term.Is(ir.LetTerm):
		c := term.Let
		var err error
		if t.context, err = t.context.AddBind(NewTermBind(c.Decl.Term.ID, c.Decl.Term.Type, DefSymbol)); err != nil {
			return err
		}

		if c.Arg != nil {
			if err := t.typecheck(c.Arg); err != nil {
				return err
			}

			if err := t.subtype(*c.Arg.Type, c.Decl.Term.Type); err != nil {
				return err
			}
		}

		term.Type = &c.Decl.Term.Type
		return nil

	case term.Is(ir.TupleTerm):
		types := make([]ir.IrType, len(term.Tuple.Elems))
		for i := range term.Tuple.Elems {
			var err error
			if err = t.typecheck(&term.Tuple.Elems[i]); err != nil {
				return err
			}
			types[i] = *term.Tuple.Elems[i].Type
		}

		typ := ir.NewTupleType(types)
		term.Type = &typ
		return nil

	case term.Is(ir.VarTerm):
		c := term.Var

		bind, err := t.context.getTermBind(c.ID)
		if err != nil {
			return err
		}

		typ, err := t.reduceType(bind.Term.Type)
		if err != nil {
			return err
		}

		term.Type = &typ
		return nil

	default:
		panic(fmt.Errorf("unhandled %T %d", term.Case, term.Case))
	}
}

func (t *Typechecker) typecheck(term *ir.IrTerm) error {
	if err := t.typecheckImpl(term); err != nil {
		return fmt.Errorf("%v\n  typechecking %s", err, *term)
	}

	reduced, err := t.reduceType(*term.Type)
	if err != nil {
		return err
	}

	term.Type = &reduced

	t.Printf("typecheck: %s |- %s", t.context.StringNoImports(), *term)
	return nil
}
