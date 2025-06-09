package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

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

func (t *Typechecker) typecheckIndexSetTerm(term *ir.IrTerm) error {
	if !term.Is(ir.IndexSetTerm) {
		panic(fmt.Errorf("expected %T %d", ir.IndexSetTerm, ir.IndexSetTerm))
	}

	c := term.IndexSet

	var index *int64
	var label *string
	switch {
	// Set field by index.
	//
	// Example:
	//   Index.set x 0 value
	case c.Index.Is(ir.ConstTerm) && c.Index.Const.Is(ir.IntLiteral):
		index = c.Index.Const.Int

	// Set field by label.
	//
	// Example:
	//   Index.set x myfield value
	case c.Index.Is(ir.VarTerm):
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
		_, field, ok := objType.FieldByID(*label)
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

func (t *Typechecker) typecheckLambdaTerm(term *ir.IrTerm) error {
	if !term.Is(ir.LambdaTerm) {
		panic(fmt.Errorf("expected %T %d", ir.LambdaTerm, ir.LambdaTerm))
	}

	c := term.Lambda

	argKind, err := inferKind(t.context, c.ArgType)
	if err != nil {
		return err
	}
	if !ir.EqualsKind(argKind, ir.NewTypeKind()) {
		return fmt.Errorf("expected lambda argument (%v) to have kind %v instead of kind %v", c.Arg, ir.NewTypeKind(), argKind)
	}

	origContext := t.context
	defer func() {
		t.context = origContext
	}()

	if t.context, err = t.context.AddBind(NewTermBind(c.Arg, c.ArgType, DefSymbol)); err != nil {
		return err
	}

	if err := t.typecheck(&c.Body); err != nil {
		return err
	}

	if c.Body.Is(ir.BlockTerm) {
		// Check all return terms have the correct function return type.
		returns := allReturns(c.Body)
		for _, ret := range returns {
			returnType := *ret.Return.Expr.Type
			if err := t.subtype(*c.Body.Type, returnType); err != nil {
				return fmt.Errorf("%v:\n%v", ret.Pos, err)
			}
		}

		// Check all function exits have the correct type.
		last := lastTerms(&c.Body)
		for _, term := range last {
			if term.Is(ir.ReturnTerm) {
				return fmt.Errorf("%v:\n redundant 'return' statement as the last term of a function", term.Pos)
			}

			if err := t.subtype(*c.Body.Type, *term.Type); err != nil {
				return fmt.Errorf("%v:\n%v", term.Pos, err)
			}

			term.LastTerm = true
		}

		if len(last) == 0 {
			return fmt.Errorf("%v:\nexpected non-empty function block", c.Body.Pos)
		}
	}

	typ := ir.NewFunctionType(c.ArgType, *c.Body.Type)
	term.Type = &typ
	return nil
}

func (t *Typechecker) typecheckMatchTerm(term *ir.IrTerm) error {
	if !term.Is(ir.MatchTerm) {
		panic(fmt.Errorf("expected %T %d", ir.MatchTerm, ir.MatchTerm))
	}

	c := term.Match

	if err := t.typecheck(&c.Term); err != nil {
		return err
	}

	variantType := *c.Term.Type
	if !variantType.Is(ir.VariantType) {
		return fmt.Errorf("expected type %v to be a variant type", variantType)
	}

	var matchType *ir.IrType
	for i := range c.Arms {
		arm := &c.Arms[i]

		index, tag, ok := variantType.TagByID(arm.Tag)
		if !ok {
			return fmt.Errorf("tag %q is not a valid tag of variant type %s", arm.Tag, variantType)
		}

		arm.Index = &index

		origContext := t.context

		var err error
		if t.context, err = t.context.AddBind(NewTermBind(arm.Arg, tag.Type, DefSymbol)); err != nil {
			return err
		}

		if err := t.typecheck(&arm.Body); err != nil {
			return err
		}

		if matchType == nil {
			matchType = arm.Body.Type
		} else {
			if err := t.subtype(*arm.Body.Type, *matchType); err != nil {
				return err
			}
		}

		t.context = origContext
	}

	term.Type = matchType
	return nil
}

func (t *Typechecker) typecheckProjectionTerm(term *ir.IrTerm) error {
	if !term.Is(ir.ProjectionTerm) {
		panic(fmt.Errorf("expected %T %d", ir.ProjectionTerm, ir.ProjectionTerm))
	}

	c := term.Projection

	if err := t.typecheck(&c.Term); err != nil {
		return err
	}

	var labelIndex *int
	switch {
	case c.Label.Is(ir.ConstTerm) && c.Label.Const.Is(ir.IntLiteral):
		number := int(*c.Label.Const.Int)
		labelIndex = &number
	}

	objType := *c.Term.Type
	switch {
	// Array projected by number literal.
	case objType.Is(ir.ArrayType) && labelIndex != nil:
		if *labelIndex < 0 || *labelIndex >= objType.Array.Size {
			return fmt.Errorf("index %d is out of bounds of array %v", *labelIndex, objType)
		}

		term.Type = &objType.Array.ElemType
		c.Index = labelIndex
		return nil

	// Array projected by variable.
	case objType.Is(ir.ArrayType):
		return fmt.Errorf("expected number literal to index array type %s", objType)

	// Struct projected by number term.
	case objType.Is(ir.StructType):
		index, field, err := objType.FieldByTerm(c.Label)
		if err != nil {
			return err
		}

		term.Type = &field.Type
		c.Index = &index
		c.LabelName = &field.ID
		return nil

	// Tuple projected by number literal.
	case objType.Is(ir.TupleType) && labelIndex != nil:
		elem, ok := objType.ElemByIndex(*labelIndex)
		if !ok {
			return fmt.Errorf("index %d is not a valid element of tuple type %s", *labelIndex, objType)
		}

		term.Type = &elem
		c.Index = labelIndex
		return nil

	// Bad label for tuple.
	case objType.Is(ir.TupleType):
		return fmt.Errorf("expected number literal to index tuple type %s", objType)

	// Variant projected by tag.
	case objType.Is(ir.VariantType):
		index, tag, err := objType.TagByTerm(c.Label)
		if err != nil {
			return err
		}

		term.Type = &tag.Type
		c.Index = &index
		c.LabelName = &tag.ID
		return nil

	default:
		return fmt.Errorf("expected projectable type (e.g., array, struct, variant, tuple, etc); got %s", objType)
	}
}

func (t *Typechecker) typecheckImpl(term *ir.IrTerm) error {
	switch {
	case term.Is(ir.AppTermTerm):
		return t.typecheckAppTermTerm(term)

	case term.Is(ir.AppTypeTerm):
		return t.typecheckAppTypeTerm(term)

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

	case term.Is(ir.InjectionTerm):
		c := term.Injection

		variantType, err := t.reduceType(c.VariantType)
		if err != nil {
			return err
		}

		if !variantType.Is(ir.VariantType) {
			return fmt.Errorf("expected type %v to be a variant type", variantType)
		}

		variantKind, err := inferKind(t.context, variantType)
		if err != nil {
			return err
		}
		if !ir.EqualsKind(variantKind, ir.NewTypeKind()) {
			return fmt.Errorf("expected type %v to have kind %v instead of kind %v", variantType, ir.NewTypeKind(), variantKind)
		}

		index, tag, err := variantType.TagByTerm(c.Tag)
		if err != nil {
			return err
		}

		tagType := ir.NewFunctionType(tag.Type, c.VariantType)
		c.Tag.Type = &tagType

		if err := t.typecheck(&c.Value); err != nil {
			return err
		}

		if err := t.subtype(*c.Value.Type, tag.Type); err != nil {
			return err
		}

		c.TagIndex = &index
		term.Type = &variantType
		return nil

	case term.Is(ir.IndexSetTerm):
		return t.typecheckIndexSetTerm(term)

	case term.Is(ir.LambdaTerm):
		return t.typecheckLambdaTerm(term)

	case term.Is(ir.LetTerm):
		c := term.Let

		var err error
		if t.context, err = t.context.AddBind(NewTermBind(c.Var, c.VarType, DefSymbol)); err != nil {
			return err
		}

		if err := t.typecheck(&c.Value); err != nil {
			return err
		}

		if err := t.subtype(*c.Value.Type, c.VarType); err != nil {
			return err
		}

		term.Type = &c.VarType
		return nil

	case term.Is(ir.MatchTerm):
		return t.typecheckMatchTerm(term)

	case term.Is(ir.ProjectionTerm):
		return t.typecheckProjectionTerm(term)

	case term.Is(ir.ReturnTerm):
		c := term.Return

		if err := t.typecheck(&c.Expr); err != nil {
			return err
		}

		typ := ir.NewTupleType(nil)
		term.Type = &typ
		return nil

	case term.Is(ir.StructTerm):
		c := term.Struct

		for i := range c.Values {
			if err := t.typecheck(&c.Values[i].Value); err != nil {
				return err
			}
		}

		typ, ok := term.StructType()
		if !ok {
			panic(fmt.Errorf("failed to determine struct type of %v", term))
		}

		term.Type = &typ
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

	case term.Is(ir.TypeAbsTerm):
		c := term.TypeAbs

		var err error
		if t.context, err = t.context.AddBind(NewTypeVarBind(c.TypeVar, c.Kind)); err != nil {
			return err
		}

		if err := t.typecheck(&c.Body); err != nil {
			return err
		}

		typ := ir.NewForallType(c.TypeVar, c.Kind, *c.Body.Type)
		term.Type = &typ
		return nil

	case term.Is(ir.VarTerm):
		c := term.Var

		bind, err := t.context.getTermBind(c.ID)
		if err != nil {
			return err
		}

		term.Type = &bind.Term.Type
		return nil

	default:
		panic(fmt.Errorf("unhandled %T %d", term.Case, term.Case))
	}
}

func (t *Typechecker) typecheck(term *ir.IrTerm) error {
	origType := term.Type

	if err := t.typecheckImpl(term); err != nil {
		return fmt.Errorf("%v\n  typechecking %s", err, *term)
	}

	reduced, err := t.reduceType(*term.Type)
	if err != nil {
		return err
	}

	if origType != nil && !ir.EqualsType(*origType, reduced) {
		return fmt.Errorf("mismatched inferred type %s and typechecked type %s", *origType, reduced)
	}

	term.Type = &reduced

	t.Printf("typecheck: %s |- %s", t.context.StringNoImports(), *term)
	return nil
}
