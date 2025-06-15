package stlc

import (
	"fmt"
	"log"

	"github.com/jabolopes/bapel/ir"
)

func probeType(term ir.IrTerm) (ir.IrType, bool) {
	if term.Type != nil {
		if term.Type.Is(ir.TupleType) {
			return term.Type.Tuple.Elems[0], true
		}

		return *term.Type, true
	}

	for _, elem := range term.Tuple.Elems {
		if elem.Type != nil {
			return *elem.Type, true
		}
	}

	return ir.IrType{}, false
}

type Inferencer struct {
	*log.Logger
	context Context
}

func (t *Inferencer) reduceType(typ ir.IrType) ir.IrType {
	reducer := typeReducer{t.Logger, t.context}
	return reducer.reduce(typ)
}

func (t *Inferencer) inferConstTerm(term, parentTerm *ir.IrTerm, expectType *ir.IrType) error {
	if !parentTerm.Is(ir.AppTypeTerm) && expectType != nil {
		// The parent term is not an AppTypeTerm, so inject an AppTypeTerm
		// that infers the type of the constant.
		*term = ir.NewAppTypeTerm(*term, *expectType)
		return t.infer(term, parentTerm, expectType)
	}

	typ := ir.Forall("a", ir.NewTypeKind(), ir.Tvar("a"))
	term.Type = &typ
	return nil
}

func (t *Inferencer) inferInjectionTerm(term *ir.IrTerm, expectType *ir.IrType) error {
	if !term.Is(ir.InjectionTerm) {
		panic(fmt.Errorf("expected %T %d", ir.InjectionTerm, ir.InjectionTerm))
	}

	c := term.Injection

	variantType := t.reduceType(c.VariantType)
	if !variantType.Is(ir.VariantType) {
		return fmt.Errorf("expected type %v to be a variant type", variantType)
	}

	_, tag, err := variantType.TagByTerm(c.Tag)
	if err != nil {
		return err
	}

	if err := t.infer(&c.Value, term, &tag.Type); err != nil {
		return err
	}

	term.Type = &variantType
	return nil
}

func (t *Inferencer) inferBlockTerm(term *ir.IrTerm, expectType *ir.IrType) error {
	if !term.Is(ir.BlockTerm) {
		panic(fmt.Errorf("expected %T %d", ir.BlockTerm, ir.BlockTerm))
	}

	c := term.Block

	origContext := t.context
	defer func() {
		t.context = origContext
	}()

	for i := range c.Terms {
		var actualExpectType *ir.IrType
		if i == len(c.Terms)-1 {
			actualExpectType = expectType
		}

		if err := t.infer(&c.Terms[i], term, actualExpectType); err != nil {
			return err
		}
	}

	// The grammar ensures that block terms are not empty.
	term.Type = c.Terms[len(c.Terms)-1].Type
	return nil
}

func (t *Inferencer) inferLambdaTerm(term *ir.IrTerm, expectType *ir.IrType) error {
	if !term.Is(ir.LambdaTerm) {
		panic(fmt.Errorf("expected %T %d", ir.LambdaTerm, ir.LambdaTerm))
	}

	c := term.Lambda

	origContext := t.context
	defer func() {
		t.context = origContext
	}()

	var err error
	if t.context, err = t.context.AddBind(NewTermBind(c.Arg, c.ArgType, DefSymbol)); err != nil {
		return err
	}

	if err := t.infer(&c.Body, term, nil /* expectType */); err != nil {
		return err
	}

	if c.Body.Type != nil {
		typ := ir.NewFunctionType(c.ArgType, *c.Body.Type)
		term.Type = &typ
	}

	return nil
}

func (t *Inferencer) inferLetTerm(term, parentTerm *ir.IrTerm, expectType *ir.IrType) error {
	c := term.Let

	var err error
	if t.context, err = t.context.AddBind(NewTermBind(c.Var, c.VarType, DefSymbol)); err != nil {
		return err
	}

	if err := t.infer(&c.Value, term, &c.VarType); err != nil {
		return err
	}

	term.Type = &c.VarType
	return nil
}

func (t *Inferencer) inferMatchTerm(term *ir.IrTerm, expectType *ir.IrType) error {
	if !term.Is(ir.MatchTerm) {
		panic(fmt.Errorf("expected %T %d", ir.MatchTerm, ir.MatchTerm))
	}

	c := term.Match

	if err := t.infer(&c.Term, term, nil /* expectType */); err != nil {
		return err
	}

	variantType := c.Term.Type
	if variantType == nil || !variantType.Is(ir.VariantType) {
		return nil
	}

	var matchType *ir.IrType
	for i := range c.Arms {
		arm := &c.Arms[i]

		_, tag, ok := variantType.TagByID(arm.Tag)
		if !ok {
			return fmt.Errorf("tag %q is not a valid tag of variant type %s", arm.Tag, variantType)
		}

		origContext := t.context

		var err error
		if t.context, err = t.context.AddBind(NewTermBind(arm.Arg, tag.Type, DefSymbol)); err != nil {
			return err
		}

		if err := t.infer(&arm.Body, term, matchType); err != nil {
			return err
		}

		if matchType == nil {
			matchType = arm.Body.Type
		}

		t.context = origContext
	}

	term.Type = matchType
	return nil
}

func (t *Inferencer) inferProjectionTerm(term, parentTerm *ir.IrTerm, expectType *ir.IrType) error {
	if !term.Is(ir.ProjectionTerm) {
		panic(fmt.Errorf("expected %T %d", ir.ProjectionTerm, ir.ProjectionTerm))
	}

	c := term.Projection

	if err := t.infer(&c.Term, term, nil /* expectType */); err != nil {
		return err
	}

	if err := t.infer(&c.Label, term, nil /* expectType */); err != nil {
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
		term.Type = &objType.Array.ElemType

	// Struct projected by number term.
	case objType.Is(ir.StructType):
		_, field, err := objType.FieldByTerm(c.Label)
		if err != nil {
			return err
		}

		term.Type = &field.Type

	// Tuple projected by number literal.
	case objType.Is(ir.TupleType) && labelIndex != nil:
		elem, ok := objType.ElemByIndex(*labelIndex)
		if !ok {
			return fmt.Errorf("index %d is not a valid element of tuple type %s", *labelIndex, objType)
		}

		term.Type = &elem

	// Variant projected by tag.
	case objType.Is(ir.VariantType):
		_, tag, err := objType.TagByTerm(c.Label)
		if err != nil {
			return err
		}

		term.Type = &tag.Type
	}

	return nil
}

func (t *Inferencer) inferStructTerm(term, parentTerm *ir.IrTerm, expectType *ir.IrType) error {
	c := term.Struct

	var structType *ir.IrType
	if expectType != nil {
		typ := t.reduceType(*expectType)
		if typ.Is(ir.StructType) {
			structType = &typ
		}
	}

	for i := range c.Values {
		value := &c.Values[i]

		var fieldType *ir.IrType
		if structType != nil {
			_, field, ok := structType.FieldByID(value.Label)
			if ok {
				fieldType = &field.Type
			}
		}

		if err := t.infer(&c.Values[i].Value, term, fieldType); err != nil {
			return err
		}
	}

	typ, ok := term.StructType()
	if ok {
		term.Type = &typ
	}

	return nil
}

func (t *Inferencer) inferImpl(term, parentTerm *ir.IrTerm, expectType *ir.IrType) error {
	switch {
	case term.Is(ir.AppTermTerm) && term.AppTerm.Fun.Is(ir.VarTerm) && ir.IsOperator(term.AppTerm.Fun.Var.ID) && expectType == nil:
		c := term.AppTerm
		if err := t.infer(&c.Fun, term, nil /* expectType */); err != nil {
			return err
		}
		if err := t.infer(&c.Arg, term, nil /* expectType */); err != nil {
			return err
		}

		typ, ok := probeType(c.Arg)
		if ok {
			*term = ir.NewAppTermTerm(ir.NewAppTypeTerm(c.Fun, typ), c.Arg)
			term.Type = &typ
		}

		return nil

	case term.Is(ir.AppTermTerm) && term.AppTerm.Fun.Is(ir.VarTerm) && ir.IsOperator(term.AppTerm.Fun.Var.ID) && expectType != nil:
		c := term.AppTerm
		if err := t.infer(&c.Fun, term, nil /* expectType */); err != nil {
			return err
		}

		argType := ir.NewTupleType([]ir.IrType{*expectType, *expectType})
		if err := t.infer(&c.Arg, term, &argType); err != nil {
			return err
		}

		typ, ok := probeType(c.Arg)
		if ok {
			*term = ir.NewAppTermTerm(ir.NewAppTypeTerm(c.Fun, typ), c.Arg)
			term.Type = &typ
		}
		return nil

	case term.Is(ir.AppTermTerm):
		c := term.AppTerm

		if err := t.infer(&c.Fun, term, nil /* expectType */); err != nil {
			return err
		}

		var argType *ir.IrType
		if c.Fun.Type != nil && c.Fun.Type.Is(ir.FunType) {
			argType = &c.Fun.Type.Fun.Arg
		}

		if err := t.infer(&c.Arg, term, argType); err != nil {
			return err
		}

		if c.Fun.Type != nil && c.Fun.Type.Is(ir.FunType) {
			typ := c.Fun.Type.Fun.Ret
			term.Type = &typ
		}

		return nil

	case term.Is(ir.AppTypeTerm):
		c := term.AppType

		if err := t.infer(&c.Fun, term, nil /* expectType */); err != nil {
			return err
		}

		if c.Fun.Type == nil || !c.Fun.Type.Is(ir.ForallType) {
			return nil
		}
		forallType := c.Fun.Type.Forall

		typ := ir.SubstituteType(forallType.Type, ir.NewVarType(forallType.Var), c.Arg)
		term.Type = &typ

		return nil

	case term.Is(ir.AssignTerm):
		c := term.Assign
		if err := t.infer(&c.Ret, term, nil /* expectType */); err != nil {
			return err
		}

		if err := t.infer(&c.Arg, term, c.Ret.Type); err != nil {
			return err
		}

		if c.Arg.Type != nil {
			term.Type = c.Arg.Type
		}
		return nil

	case term.Is(ir.BlockTerm):
		return t.inferBlockTerm(term, expectType)

	case term.Is(ir.ConstTerm):
		return t.inferConstTerm(term, parentTerm, expectType)

	case term.Is(ir.IfTerm):
		c := term.If

		if err := t.infer(&c.Condition, term, nil /* expectType */); err != nil {
			return err
		}
		if err := t.infer(&c.Then, term, nil /* expectType */); err != nil {
			return err
		}
		if c.Else != nil {
			if err := t.infer(c.Else, term, c.Then.Type); err != nil {
				return err
			}
		}
		return nil

	case term.Is(ir.InjectionTerm):
		return t.inferInjectionTerm(term, expectType)

	case term.Is(ir.IndexSetTerm):
		c := term.IndexSet
		if err := t.infer(&c.Obj, term, nil /* expectType */); err != nil {
			return err
		}
		if err := t.infer(&c.Index, term, nil /* expectType */); err != nil {
			return err
		}
		return t.infer(&c.Value, term, nil /* expectType */)

	case term.Is(ir.LambdaTerm):
		return t.inferLambdaTerm(term, expectType)

	case term.Is(ir.LetTerm):
		return t.inferLetTerm(term, parentTerm, expectType)

	case term.Is(ir.MatchTerm):
		return t.inferMatchTerm(term, expectType)

	case term.Is(ir.ProjectionTerm):
		return t.inferProjectionTerm(term, parentTerm, expectType)

	case term.Is(ir.ReturnTerm):
		c := term.Return

		// TODO: Pass function return type as expectType.
		return t.infer(&c.Expr, term, nil /* expectType */)

	case term.Is(ir.StructTerm):
		return t.inferStructTerm(term, parentTerm, expectType)

	case term.Is(ir.TupleTerm) &&
		expectType != nil && expectType.Is(ir.TupleType) &&
		len(expectType.Tuple.Elems) == len(term.Tuple.Elems):

		typ := func() *ir.IrType {
			t := ir.NewTupleType(nil)
			return &t
		}()

		for i := range term.Tuple.Elems {
			if err := t.infer(&term.Tuple.Elems[i], term, &expectType.Tuple.Elems[i]); err != nil {
				return err
			}

			if term.Tuple.Elems[i].Type == nil {
				typ = nil
			} else if typ != nil {
				typ.Tuple.Elems = append(typ.Tuple.Elems, *term.Tuple.Elems[i].Type)
			}
		}

		term.Type = typ
		return nil

	case term.Is(ir.TupleTerm):
		typ := func() *ir.IrType {
			t := ir.NewTupleType(nil)
			return &t
		}()

		for i := range term.Tuple.Elems {
			if err := t.infer(&term.Tuple.Elems[i], term, nil /* expectType */); err != nil {
				return err
			}

			if term.Tuple.Elems[i].Type == nil {
				typ = nil
			} else if typ != nil {
				typ.Tuple.Elems = append(typ.Tuple.Elems, *term.Tuple.Elems[i].Type)
			}
		}

		term.Type = typ
		return nil

	case term.Is(ir.TypeAbsTerm):
		c := term.TypeAbs

		var err error
		if t.context, err = t.context.AddBind(NewTypeVarBind(c.TypeVar, c.Kind)); err != nil {
			return err
		}

		if err := t.infer(&c.Body, term, nil /* expectType */); err != nil {
			return err
		}

		if c.Body.Type != nil {
			typ := ir.NewForallType(c.TypeVar, c.Kind, *c.Body.Type)
			term.Type = &typ
		}
		return nil

	case term.Is(ir.VarTerm):
		c := term.Var

		bind, err := t.context.getTermBind(c.ID)
		if err != nil {
			return nil
		}

		term.Type = &bind.Term.Type
		return nil

	default:
		panic(fmt.Errorf("unhandled %T %d", term.Case, term.Case))
	}
}

func (t *Inferencer) infer(term, parentTerm *ir.IrTerm, expectType *ir.IrType) error {
	if err := t.inferImpl(term, parentTerm, expectType); err != nil {
		return fmt.Errorf("%v\n  inferring %s", err, term)
	}

	if term.Type != nil {
		reduced := t.reduceType(*term.Type)
		term.Type = &reduced
	}

	if term.Type == nil {
		t.Printf("infer: %s |- %s : ?", t.context.StringNoImports(), term)
	} else {
		t.Printf("infer: %s |- %s", t.context.StringNoImports(), term)
	}

	return nil
}

func (t *Inferencer) inferFunction(function *ir.IrFunction) error {
	decl := function.Decl()

	var err error
	t.context, err = t.context.AddBind(NewTermBind(decl.Term.ID, decl.Term.Type, DefSymbol))
	if err != nil {
		return err
	}

	if t.context, err = t.context.enterFunction(function.TypeVars, function.Args); err != nil {
		return err
	}

	return t.infer(&function.Body, nil /* parentTerm */, &function.RetType)
}
