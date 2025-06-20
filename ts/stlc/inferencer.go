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
	reducer := typeReducer{t.Logger}
	return reducer.reduce(t.context, typ)
}

func (t *Inferencer) inferConstTerm(term, parentTerm *ir.IrTerm, expectType *ir.IrType) error {
	if !term.Is(ir.ConstTerm) {
		panic(fmt.Errorf("expected %T %d", ir.ConstTerm, ir.ConstTerm))
	}

	// TODO: Check parentTerm is not nil.
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

func (t *Inferencer) inferBlockTerm(term *ir.IrTerm, expectType *ir.IrType) error {
	if !term.Is(ir.BlockTerm) {
		panic(fmt.Errorf("expected %T %d", ir.BlockTerm, ir.BlockTerm))
	}

	c := term.Block

	origContext := t.context
	defer func() { t.context = origContext }()

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
	if !term.Is(ir.LetTerm) {
		panic(fmt.Errorf("expected %T %d", ir.LetTerm, ir.LetTerm))
	}

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

	objType := c.Term.Type
	switch {
	case objType == nil:
		break

	case objType.Is(ir.StructType):
		_, field, err := objType.FieldByLabel(c.Label)
		if err != nil {
			return err
		}

		term.Type = &field.Type

	case objType.Is(ir.TupleType):
		_, elemType, err := objType.ElemByLabel(c.Label)
		if err != nil {
			return err
		}

		term.Type = &elemType

	case objType.Is(ir.VariantType):
		_, tag, err := objType.TagByLabel(c.Label)
		if err != nil {
			return err
		}

		term.Type = &tag.Type
	}

	return nil
}

func (t *Inferencer) inferSetTerm(term, parentTerm *ir.IrTerm, expectType *ir.IrType) error {
	if !term.Is(ir.SetTerm) {
		panic(fmt.Errorf("expected %T %d", ir.SetTerm, ir.SetTerm))
	}

	c := term.Set

	if err := t.infer(&c.Term, parentTerm, expectType); err != nil {
		return err
	}

	objType := c.Term.Type
	switch {
	case objType == nil:
		for i := range c.Values {
			if err := t.infer(&c.Values[i].Value, parentTerm, nil /* expectType */); err != nil {
				return err
			}
		}

	case objType.Is(ir.StructType):
		for i := range c.Values {
			lv := &c.Values[i]

			_, field, err := objType.FieldByLabel(lv.Label)
			if err != nil {
				return err
			}

			if err := t.infer(&lv.Value, parentTerm, &field.Type); err != nil {
				return err
			}
		}

	case objType.Is(ir.TupleType):
		for i := range c.Values {
			lv := &c.Values[i]

			_, elemType, err := objType.ElemByLabel(lv.Label)
			if err != nil {
				return err
			}

			if err := t.infer(&lv.Value, parentTerm, &elemType); err != nil {
				return err
			}
		}
	}

	term.Type = c.Term.Type
	return nil
}

func (t *Inferencer) inferStructTerm(term, parentTerm *ir.IrTerm, expectType *ir.IrType) error {
	if !term.Is(ir.StructTerm) {
		panic(fmt.Errorf("expected %T %d", ir.StructTerm, ir.StructTerm))
	}

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

		if err := t.infer(&value.Value, term, fieldType); err != nil {
			return err
		}
	}

	typ, ok := term.StructType()
	if ok {
		term.Type = &typ
	}

	return nil
}

func (t *Inferencer) inferTupleTerm(term, parentTerm *ir.IrTerm, expectType *ir.IrType) error {
	if !term.Is(ir.TupleTerm) {
		panic(fmt.Errorf("expected %T %d", ir.TupleTerm, ir.TupleTerm))
	}

	c := term.Tuple

	var tupleType *ir.IrType
	if expectType != nil {
		typ := t.reduceType(*expectType)
		if typ.Is(ir.TupleType) {
			tupleType = &typ
		}
	}

	for i := range c.Elems {
		var elemType *ir.IrType
		if tupleType != nil {
			typ, ok := tupleType.ElemByIndex(i)
			if ok {
				elemType = &typ
			}
		}

		if err := t.infer(&c.Elems[i], term, elemType); err != nil {
			return err
		}
	}

	typ, ok := term.TupleType()
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

	case term.Is(ir.SetTerm):
		return t.inferSetTerm(term, parentTerm, expectType)

	case term.Is(ir.StructTerm):
		return t.inferStructTerm(term, parentTerm, expectType)

	case term.Is(ir.TupleTerm):
		return t.inferTupleTerm(term, parentTerm, expectType)

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

func (t *Inferencer) inferFunction(function *ir.IrFunction) (Context, error) {
	origContext := t.context

	decl := function.Decl()

	var err error
	t.context, err = t.context.AddBind(NewTermBind(decl.Term.ID, decl.Term.Type, DefSymbol))
	if err != nil {
		return origContext, err
	}

	retContext := t.context

	if t.context, err = t.context.enterFunction(function.TypeVars, function.Args); err != nil {
		return origContext, err
	}

	if err := t.infer(&function.Body, nil /* parentTerm */, &function.RetType); err != nil {
		return origContext, err
	}

	return retContext, nil
}
