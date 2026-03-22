package ast

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

type desugarer struct {
	err error
}

func (t *desugarer) desugarAppTermExpr(source *Expr) ir.IrTerm {
	if !source.Is(AppTermExpr) {
		panic(fmt.Errorf("expected %T %d", AppTermExpr, AppTermExpr))
	}

	c := source.AppTerm

	return ir.NewAppTermTerm(t.desugar(&c.Fun), t.desugar(&c.Arg))
}

func (t *desugarer) desugarAppTypeExpr(source *Expr) ir.IrTerm {
	if !source.Is(AppTypeExpr) {
		panic(fmt.Errorf("expected %T %d", AppTypeExpr, AppTypeExpr))
	}

	c := source.AppType

	return ir.NewAppTypeTerm(t.desugar(&c.Fun), c.Arg)
}

func (t *desugarer) desugarAssignExpr(source *Expr) ir.IrTerm {
	if !source.Is(AssignExpr) {
		panic(fmt.Errorf("expected %T %d", AssignExpr, AssignExpr))
	}

	c := source.Assign

	return ir.NewAssignTerm(t.desugar(&c.Arg), t.desugar(&c.Ret))
}

func (t *desugarer) desugarBlockExpr(source *Expr) ir.IrTerm {
	if !source.Is(BlockExpr) {
		panic(fmt.Errorf("expected %T %d", BlockExpr, BlockExpr))
	}

	c := source.Block

	terms := make([]ir.IrTerm, 0, len(c.Exprs))
	for _, source := range c.Exprs {
		terms = append(terms, t.desugar(&source))
	}

	return ir.NewBlockTerm(terms)
}

func (t *desugarer) desugarConstExpr(source *Expr) ir.IrTerm {
	if !source.Is(ConstExpr) {
		panic(fmt.Errorf("expected %T %d", ConstExpr, ConstExpr))
	}

	c := source.Const

	return ir.NewConstTerm(c.IrLiteral)
}

func (t *desugarer) desugarInjectionExpr(source *Expr) ir.IrTerm {
	if !source.Is(InjectionExpr) {
		panic(fmt.Errorf("expected %T %d", InjectionExpr, InjectionExpr))
	}

	c := source.Injection

	return ir.NewInjectionTerm(c.VariantType, c.Tag, t.desugar(&c.Expr))
}

func (t *desugarer) desugarLambdaExpr(source *Expr) ir.IrTerm {
	if !source.Is(LambdaExpr) {
		panic(fmt.Errorf("expected %T %d", LambdaExpr, LambdaExpr))
	}

	c := source.Lambda

	return ir.NewLambdaTerm(c.Arg, t.desugar(&c.Body))
}

func (t *desugarer) desugarLetExpr(source *Expr) ir.IrTerm {
	if !source.Is(LetExpr) {
		panic(fmt.Errorf("expected %T %d", LetExpr, LetExpr))
	}

	c := source.Let

	return ir.NewLetTerm(c.Var, c.VarType, t.desugar(&c.Expr))
}

func (t *desugarer) desugarMatchExpr(source *Expr) ir.IrTerm {
	if !source.Is(MatchExpr) {
		panic(fmt.Errorf("expected %T %d", MatchExpr, MatchExpr))
	}

	c := source.Match

	arms := make([]ir.MatchArm, 0, len(c.Arms))
	for _, source := range c.Arms {
		arms = append(arms, ir.NewMatchArm(source.Tag, source.Arg, t.desugar(&source.Body)))
	}

	return ir.NewMatchTerm(t.desugar(&c.Expr), arms)
}

func (t *desugarer) desugarProjectionExpr(source *Expr) ir.IrTerm {
	if !source.Is(ProjectionExpr) {
		panic(fmt.Errorf("expected %T %d", ProjectionExpr, ProjectionExpr))
	}

	c := source.Projection

	return ir.NewProjectionTerm(t.desugar(&c.Expr), c.Label)
}

func (t *desugarer) desugarReturnExpr(source *Expr) ir.IrTerm {
	if !source.Is(ReturnExpr) {
		panic(fmt.Errorf("expected %T %d", ReturnExpr, ReturnExpr))
	}

	c := source.Return

	return ir.NewReturnTerm(t.desugar(&c.Expr))
}

func (t *desugarer) desugarSetExpr(source *Expr) ir.IrTerm {
	if !source.Is(SetExpr) {
		panic(fmt.Errorf("expected %T %d", SetExpr, SetExpr))
	}

	c := source.Set

	values := make([]ir.LabelValue, 0, len(c.Values))
	for _, source := range c.Values {
		values = append(values, ir.LabelValue{source.Label, t.desugar(&source.Value)})
	}

	return ir.NewSetTerm(t.desugar(&c.Expr), values)
}

func (t *desugarer) desugarStructExpr(source *Expr) ir.IrTerm {
	if !source.Is(StructExpr) {
		panic(fmt.Errorf("expected %T %d", StructExpr, StructExpr))
	}

	c := source.Struct

	values := make([]ir.LabelValue, 0, len(c.Values))
	for _, source := range c.Values {
		values = append(values, ir.LabelValue{source.Label, t.desugar(&source.Value)})
	}

	return ir.NewStructTerm(values)
}

func (t *desugarer) desugarTupleExpr(source *Expr) ir.IrTerm {
	if !source.Is(TupleExpr) {
		panic(fmt.Errorf("expected %T %d", TupleExpr, TupleExpr))
	}

	c := source.Tuple

	elems := make([]ir.IrTerm, 0, len(c.Elems))
	for _, source := range c.Elems {
		elems = append(elems, t.desugar(&source))
	}

	return ir.NewTupleTerm(elems)
}

func (t *desugarer) desugarTypeAbsExpr(source *Expr) ir.IrTerm {
	if !source.Is(TypeAbsExpr) {
		panic(fmt.Errorf("expected %T %d", TypeAbsExpr, TypeAbsExpr))
	}

	c := source.TypeAbs

	return ir.NewTypeAbsTerm(c.Arg, t.desugar(&c.Body))
}

func (t *desugarer) desugarVarExpr(source *Expr) ir.IrTerm {
	if !source.Is(VarExpr) {
		panic(fmt.Errorf("expected %T %d", VarExpr, VarExpr))
	}

	c := source.Var

	return ir.NewVarTerm(c.ID)
}

func (t *desugarer) desugarImpl(source *Expr) ir.IrTerm {
	switch {
	case source.Is(AppTermExpr):
		return t.desugarAppTermExpr(source)

	case source.Is(AppTypeExpr):
		return t.desugarAppTypeExpr(source)

	case source.Is(AssignExpr):
		return t.desugarAssignExpr(source)

	case source.Is(BlockExpr):
		return t.desugarBlockExpr(source)

	case source.Is(ConstExpr):
		return t.desugarConstExpr(source)

	case source.Is(InjectionExpr):
		return t.desugarInjectionExpr(source)

	case source.Is(LambdaExpr):
		return t.desugarLambdaExpr(source)

	case source.Is(LetExpr):
		return t.desugarLetExpr(source)

	case source.Is(MatchExpr):
		return t.desugarMatchExpr(source)

	case source.Is(ProjectionExpr):
		return t.desugarProjectionExpr(source)

	case source.Is(ReturnExpr):
		return t.desugarReturnExpr(source)

	case source.Is(SetExpr):
		return t.desugarSetExpr(source)

	case source.Is(StructExpr):
		return t.desugarStructExpr(source)

	case source.Is(TupleExpr):
		return t.desugarTupleExpr(source)

	case source.Is(TypeAbsExpr):
		return t.desugarTypeAbsExpr(source)

	case source.Is(VarExpr):
		return t.desugarVarExpr(source)

	default:
		panic(fmt.Errorf("unhandled %T %d", source.Case, source.Case))
	}
}

func (t *desugarer) desugar(source *Expr) ir.IrTerm {
	if t.err != nil {
		return ir.IrTerm{}
	}

	return t.desugarImpl(source)
}

func desugarExpr(source *Expr) (ir.IrTerm, error) {
	t := &desugarer{}
	term := t.desugar(source)
	if t.err != nil {
		return ir.IrTerm{}, t.err
	}

	term.Pos = source.Pos
	return term, nil
}

func DesugarFunction(function Function) (ir.IrFunction, error) {
	body, err := desugarExpr(&function.Body)
	if err != nil {
		return ir.IrFunction{}, err
	}

	fun := ir.NewFunction(function.Export, function.ID, function.TypeVars, function.Args, function.RetType, body)
	fun.Pos = function.Pos
	return fun, nil
}
