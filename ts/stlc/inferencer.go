package stlc

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/ts/list"
)

type expectReturn struct {
	RetType ir.IrType
}

type Inferencer struct {
	context       Context
	existVars     map[int]existVar
	expectReturns list.List[expectReturn]
}

func (t *Inferencer) newEvar() ir.IrType {
	evar := t.context.GenFreshExistVar()
	t.existVars[evar.ExistVar] = existVar{}
	return evar
}

func (t *Inferencer) reduceType(typ ir.IrType) ir.IrType {
	reducer := typeReducer{}
	return reducer.reduce(t.context, typ)
}

// predicateType predicates the type without reducing type aliases.
func (t *Inferencer) predicateType(typ ir.IrType) (ir.IrType, error) {
	predicator := typePredicator{t.context, nil /* tvars */}

	newType, err := predicator.predicate(typ)
	if err != nil {
		return ir.IrType{}, err
	}

	return ir.ForallVars(predicator.tvars, newType), nil
}

func (t *Inferencer) reduceAndPredicateType(typ ir.IrType) (ir.IrType, error) {
	return t.predicateType(t.reduceType(typ))
}

func (t *Inferencer) inferAppTermTerm(evar ir.IrType, term, parentTerm *ir.IrTerm, expectType *ir.IrType) error {
	if !term.Is(ir.AppTermTerm) {
		panic(fmt.Errorf("expected %T %d", ir.AppTermTerm, ir.AppTermTerm))
	}

	c := term.AppTerm

	if err := t.infer(&c.Fun, term, nil /* expectType */); err != nil {
		return err
	}

	switch {
	case c.Fun.Type == nil:
		if err := t.infer(&c.Arg, term, nil /* expectType */); err != nil {
			return err
		}
		return nil

	case c.Fun.Type.Is(ir.FunType):
		argType := c.Fun.Type.Fun.Arg
		retType := c.Fun.Type.Fun.Ret

		if err := t.infer(&c.Arg, term, &argType); err != nil {
			return err
		}

		t.unify(evar, retType)
		t.unify(*c.Fun.Type, ir.NewFunctionType(argType, evar))
		term.Type = &retType

		return nil

	case c.Fun.Type.Is(ir.ForallType) && !c.Fun.Is(ir.AppTermTerm):
		if err := t.infer(&c.Arg, term, nil /* expectType */); err != nil {
			return err
		}

		term.AppTerm.Fun = ir.NewAppTypeTerm(c.Fun, t.newEvar())
		return t.infer(term, parentTerm, expectType)

	case c.Fun.Type.Is(ir.ForallType):
		if err := t.infer(&c.Arg, term, nil /* expectType */); err != nil {
			return err
		}

		return t.infer(term, parentTerm, expectType)

	default:
		panic(fmt.Errorf("unhandled %T %v", c.Fun.Type, c.Fun.Type))
	}

	return nil
}

func (t *Inferencer) inferAppTypeTerm(evar ir.IrType, term, parentTerm *ir.IrTerm, expectType *ir.IrType) error {
	c := term.AppType

	if err := t.infer(&c.Fun, term, nil /* expectType */); err != nil {
		return err
	}
	if c.Fun.Type == nil || !c.Fun.Type.Is(ir.ForallType) {
		return nil
	}

	forallType := c.Fun.Type.Forall

	typ := ir.SubstituteType(forallType.Type, ir.NewVarType(forallType.Var), c.Arg)
	t.unify(evar, typ)
	term.Type = &typ

	return nil
}

func (t *Inferencer) inferAssignTerm(evar ir.IrType, term, parentTerm *ir.IrTerm, expectType *ir.IrType) error {
	if !term.Is(ir.AssignTerm) {
		panic(fmt.Errorf("expected %T %d", ir.AssignTerm, ir.AssignTerm))
	}

	c := term.Assign

	if err := t.infer(&c.Ret, term, nil /* expectType */); err != nil {
		return err
	}

	if err := t.infer(&c.Arg, term, c.Ret.Type); err != nil {
		return err
	}

	if c.Arg.Type != nil {
		t.unify(evar, *c.Arg.Type)
		term.Type = c.Arg.Type
	}

	return nil
}

func (t *Inferencer) inferBlockTerm(evar ir.IrType, term, parentTerm *ir.IrTerm, expectType *ir.IrType) error {
	if !term.Is(ir.BlockTerm) {
		panic(fmt.Errorf("expected %T %d", ir.BlockTerm, ir.BlockTerm))
	}

	c := term.Block

	origContext := t.context
	defer func() { t.context = origContext }()

	var err error
	t.context, err = t.context.enterScope()
	if err != nil {
		return err
	}

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
	lastTermType := c.Terms[len(c.Terms)-1].Type
	if lastTermType != nil {
		t.unify(evar, *lastTermType)
		term.Type = lastTermType
	}

	// Before closing the scope.
	return t.solveTerm(term)
}

func (t *Inferencer) inferConstTerm(evar ir.IrType, term, parentTerm *ir.IrTerm, expectType *ir.IrType) error {
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

func (t *Inferencer) inferInjectionTerm(evar ir.IrType, term *ir.IrTerm, expectType *ir.IrType) error {
	if !term.Is(ir.InjectionTerm) {
		panic(fmt.Errorf("expected %T %d", ir.InjectionTerm, ir.InjectionTerm))
	}

	c := term.Injection

	variantType, err := t.reduceAndPredicateType(c.VariantType)
	if err != nil {
		return err
	}
	if !variantType.Is(ir.VariantType) {
		return fmt.Errorf("expected type %v to be a variant type", variantType)
	}

	_, tag, err := variantType.TagByLabel(c.Tag)
	if err != nil {
		return err
	}

	if err := t.infer(&c.Value, term, &tag.Type); err != nil {
		return err
	}

	t.unify(evar, c.VariantType)
	term.Type = &c.VariantType

	return nil
}

func (t *Inferencer) inferLambdaTerm(evar ir.IrType, term *ir.IrTerm, expectType *ir.IrType) error {
	if !term.Is(ir.LambdaTerm) {
		panic(fmt.Errorf("expected %T %d", ir.LambdaTerm, ir.LambdaTerm))
	}

	c := term.Lambda

	origContext := t.context
	defer func() { t.context = origContext }()

	var err error
	if t.context, err = t.context.AddBind(NewTermDefBind(c.Arg, c.ArgType)); err != nil {
		return err
	}

	retEvar := t.newEvar()
	{
		orig := t.expectReturns
		t.expectReturns = t.expectReturns.Add(expectReturn{retEvar})
		defer func() { t.expectReturns = orig }()
	}

	if err := t.infer(&c.Body, term, nil /* expectType */); err != nil {
		return err
	}

	if c.Body.Type != nil {
		t.unify(retEvar, *c.Body.Type)
		typ := ir.NewFunctionType(c.ArgType, retEvar)
		t.unify(evar, typ)
		term.Type = &typ
	}

	return nil
}

func (t *Inferencer) inferLetTerm(evar ir.IrType, term, parentTerm *ir.IrTerm, expectType *ir.IrType) error {
	if !term.Is(ir.LetTerm) {
		panic(fmt.Errorf("expected %T %d", ir.LetTerm, ir.LetTerm))
	}

	c := term.Let

	if c.VarType == nil {
		if err := t.infer(&c.Value, term, nil /* expectType */); err != nil {
			return err
		}

		if c.Value.Type == nil {
			return fmt.Errorf("failed to infer type for let variable %v; please add a type annotation 'let %v : ... = ...'", c.Var, c.Var)
		}

		varType, err := t.predicateType(*c.Value.Type)
		if err != nil {
			return err
		}

		if t.context, err = t.context.AddBind(NewTermDefBind(c.Var, varType)); err != nil {
			return err
		}

		t.unify(evar, varType)
		c.VarType = &varType
		term.Type = &varType

		return nil
	}

	varType := *c.VarType

	var err error
	if t.context, err = t.context.AddBind(NewTermDefBind(c.Var, varType)); err != nil {
		return err
	}

	if err := t.infer(&c.Value, term, &varType); err != nil {
		return err
	}

	t.unify(evar, varType)
	term.Type = &varType
	return nil
}

func (t *Inferencer) inferMatchTerm(evar ir.IrType, term *ir.IrTerm, expectType *ir.IrType) error {
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
		if t.context, err = t.context.AddBind(NewTermDefBind(arm.Arg, tag.Type)); err != nil {
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

	if matchType != nil {
		t.unify(evar, *matchType)
		term.Type = matchType
	}

	return nil
}

func (t *Inferencer) inferProjectionTerm(evar ir.IrType, term, parentTerm *ir.IrTerm, expectType *ir.IrType) error {
	if !term.Is(ir.ProjectionTerm) {
		panic(fmt.Errorf("expected %T %d", ir.ProjectionTerm, ir.ProjectionTerm))
	}

	c := term.Projection

	if err := t.infer(&c.Term, term, nil /* expectType */); err != nil {
		return err
	}

	var objType *ir.IrType
	if c.Term.Type != nil {
		typ, err := t.reduceAndPredicateType(*c.Term.Type)
		if err != nil {
			return err
		}
		objType = &typ
	}

	switch {
	case objType == nil:
		break

	case objType.Is(ir.StructType):
		_, field, err := objType.FieldByLabel(c.Label)
		if err != nil {
			return err
		}

		t.unify(evar, field.Type)
		term.Type = &field.Type

	case objType.Is(ir.TupleType):
		_, elemType, err := objType.ElemByLabel(c.Label)
		if err != nil {
			return err
		}

		t.unify(evar, elemType)
		term.Type = &elemType

	case objType.Is(ir.VariantType):
		_, tag, err := objType.TagByLabel(c.Label)
		if err != nil {
			return err
		}

		t.unify(evar, tag.Type)
		term.Type = &tag.Type
	}

	return nil
}

func (t *Inferencer) inferReturnTerm(evar ir.IrType, term, parentTerm *ir.IrTerm, expectType *ir.IrType) error {
	if !term.Is(ir.ReturnTerm) {
		panic(fmt.Errorf("expected %T %d", ir.ReturnTerm, ir.ReturnTerm))
	}

	c := term.Return

	expectReturn, ok := t.expectReturns.Value()
	if !ok {
		return fmt.Errorf("term %v must be inside a function or lambda term", *term)
	}

	t.unify(evar, expectReturn.RetType)
	if err := t.infer(&c.Expr, term, &expectReturn.RetType); err != nil {
		return err
	}

	if c.Expr.Type != nil {
		t.unify(evar, *c.Expr.Type)
		term.Type = c.Expr.Type
	}

	return nil
}

func (t *Inferencer) inferSetTerm(evar ir.IrType, term, parentTerm *ir.IrTerm, expectType *ir.IrType) error {
	if !term.Is(ir.SetTerm) {
		panic(fmt.Errorf("expected %T %d", ir.SetTerm, ir.SetTerm))
	}

	c := term.Set

	if err := t.infer(&c.Term, parentTerm, expectType); err != nil {
		return err
	}

	var objType *ir.IrType
	if c.Term.Type != nil {
		typ, err := t.reduceAndPredicateType(*c.Term.Type)
		if err != nil {
			return err
		}
		objType = &typ
	}

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

	if c.Term.Type != nil {
		t.unify(evar, *c.Term.Type)
		term.Type = c.Term.Type
	}

	return nil
}

func (t *Inferencer) inferStructTerm(evar ir.IrType, term, parentTerm *ir.IrTerm, expectType *ir.IrType) error {
	if !term.Is(ir.StructTerm) {
		panic(fmt.Errorf("expected %T %d", ir.StructTerm, ir.StructTerm))
	}

	c := term.Struct

	var structType *ir.IrType
	if expectType != nil {
		typ, err := t.reduceAndPredicateType(*expectType)
		if err != nil {
			return err
		}
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

	if structType != nil {
		t.unify(evar, *structType)
		term.Type = structType
	} else {
		typ, ok := term.StructType()
		if ok {
			t.unify(evar, typ)
			term.Type = &typ
		}
	}

	return nil
}

func (t *Inferencer) inferTupleTerm(evar ir.IrType, term, parentTerm *ir.IrTerm, expectType *ir.IrType) error {
	if !term.Is(ir.TupleTerm) {
		panic(fmt.Errorf("expected %T %d", ir.TupleTerm, ir.TupleTerm))
	}

	c := term.Tuple

	var tupleType *ir.IrType
	if expectType != nil {
		typ, err := t.reduceAndPredicateType(*expectType)
		if err != nil {
			return err
		}
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

	if tupleType != nil {
		t.unify(evar, *tupleType)
		term.Type = tupleType
	} else {
		typ, ok := term.TupleType()
		if ok {
			t.unify(evar, typ)
			term.Type = &typ
		}
	}

	return nil
}

func (t *Inferencer) inferTypeAbsTerm(evar ir.IrType, term, parentTerm *ir.IrTerm, expectType *ir.IrType) error {
	if !term.Is(ir.TypeAbsTerm) {
		panic(fmt.Errorf("expected %T %d", ir.TypeAbsTerm, ir.TypeAbsTerm))
	}

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
		t.unify(evar, typ)
		term.Type = &typ
	}

	return nil
}

func (t *Inferencer) inferVarTerm(evar ir.IrType, term, parentTerm *ir.IrTerm, expectType *ir.IrType) error {
	c := term.Var

	switch bind, ok := t.context.lookupTermDeclOrDefBind(c.ID); {
	case ok && bind.Is(TermDeclBind):
		t.unify(evar, bind.TermDecl.Type)
		term.Type = &bind.TermDecl.Type
	case ok && bind.Is(TermDefBind):
		t.unify(evar, bind.TermDef.Type)
		term.Type = &bind.TermDef.Type
	}

	return nil
}

func (t *Inferencer) inferImpl(evar ir.IrType, term, parentTerm *ir.IrTerm, expectType *ir.IrType) error {
	switch {
	case term.Is(ir.AppTermTerm):
		return t.inferAppTermTerm(evar, term, parentTerm, expectType)

	case term.Is(ir.AppTypeTerm):
		return t.inferAppTypeTerm(evar, term, parentTerm, expectType)

	case term.Is(ir.AssignTerm):
		return t.inferAssignTerm(evar, term, parentTerm, expectType)

	case term.Is(ir.BlockTerm):
		return t.inferBlockTerm(evar, term, parentTerm, expectType)

	case term.Is(ir.ConstTerm):
		return t.inferConstTerm(evar, term, parentTerm, expectType)

	case term.Is(ir.InjectionTerm):
		return t.inferInjectionTerm(evar, term, expectType)

	case term.Is(ir.LambdaTerm):
		return t.inferLambdaTerm(evar, term, expectType)

	case term.Is(ir.LetTerm):
		return t.inferLetTerm(evar, term, parentTerm, expectType)

	case term.Is(ir.MatchTerm):
		return t.inferMatchTerm(evar, term, expectType)

	case term.Is(ir.ProjectionTerm):
		return t.inferProjectionTerm(evar, term, parentTerm, expectType)

	case term.Is(ir.ReturnTerm):
		return t.inferReturnTerm(evar, term, parentTerm, expectType)

	case term.Is(ir.SetTerm):
		return t.inferSetTerm(evar, term, parentTerm, expectType)

	case term.Is(ir.StructTerm):
		return t.inferStructTerm(evar, term, parentTerm, expectType)

	case term.Is(ir.TupleTerm):
		return t.inferTupleTerm(evar, term, parentTerm, expectType)

	case term.Is(ir.TypeAbsTerm):
		return t.inferTypeAbsTerm(evar, term, parentTerm, expectType)

	case term.Is(ir.VarTerm):
		return t.inferVarTerm(evar, term, parentTerm, expectType)

	default:
		panic(fmt.Errorf("unhandled %T %d", term.Case, term.Case))
	}
}

func (t *Inferencer) infer(term, parentTerm *ir.IrTerm, expectType *ir.IrType) error {
	evar := t.newEvar()

	if expectType != nil {
		t.unify(evar, *expectType)
	}

	if err := t.inferImpl(evar, term, parentTerm, expectType); err != nil {
		return fmt.Errorf("%v\n  inferring %s", err, term)
	}

	if term.Type != nil {
		{
			typ := t.solveType(*term.Type)
			term.Type = &typ
		}

		{
			typ, err := t.predicateType(*term.Type)
			if err != nil {
				return err
			}
			term.Type = &typ
		}
	}

	if term.Type == nil {
		glog.V(1).Infof("infer: %s\n  |- %s : ?", t.context, term)
	} else {
		glog.V(1).Infof("infer: %s\n  |- %s", t.context, term)
	}

	return nil
}

func (t *Inferencer) inferFunction(function *ir.IrFunction) (Context, error) {
	origContext := t.context

	decl := function.Decl()

	var err error
	t.context, err = t.context.AddBind(NewTermDefBind(decl.Term.ID, decl.Term.Type))
	if err != nil {
		return origContext, err
	}

	retContext := t.context

	if t.context, err = t.context.enterFunction(function.TypeVars, function.Args); err != nil {
		return origContext, err
	}

	{
		orig := t.expectReturns
		t.expectReturns = t.expectReturns.Add(expectReturn{function.RetType})
		defer func() { t.expectReturns = orig }()
	}

	if err := t.infer(&function.Body, nil /* parentTerm */, &function.RetType); err != nil {
		return origContext, err
	}

	return retContext, nil
}

func NewInferencer(context Context) *Inferencer {
	return &Inferencer{context, map[int]existVar{}, list.New[expectReturn]()}
}
