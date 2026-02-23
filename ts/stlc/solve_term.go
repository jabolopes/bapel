package stlc

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/jabolopes/bapel/ir"
)

func (t *Inferencer) solveAppTermTerm(term *ir.IrTerm) error {
	if !term.Is(ir.AppTermTerm) {
		panic(fmt.Errorf("expected %T %d", ir.AppTermTerm, ir.AppTermTerm))
	}

	c := term.AppTerm

	if err := t.solveTerm(&c.Fun); err != nil {
		return err
	}
	if err := t.solveTerm(&c.Arg); err != nil {
		return err
	}

	return nil
}

func (t *Inferencer) solveAppTypeTerm(term *ir.IrTerm) error {
	if !term.Is(ir.AppTypeTerm) {
		panic(fmt.Errorf("expected %T %d", ir.AppTypeTerm, ir.AppTypeTerm))
	}

	c := term.AppType

	if err := t.solveTerm(&c.Fun); err != nil {
		return err
	}

	c.Arg = t.solveType(c.Arg)

	return nil
}

func (t *Inferencer) solveBlockTerm(term *ir.IrTerm) error {
	if !term.Is(ir.BlockTerm) {
		panic(fmt.Errorf("expected %T %d", ir.BlockTerm, ir.BlockTerm))
	}

	c := term.Block

	for i := range c.Terms {
		if err := t.solveTerm(&c.Terms[i]); err != nil {
			return err
		}
	}

	return nil
}

func (t *Inferencer) solveIfTerm(term *ir.IrTerm) error {
	if !term.Is(ir.IfTerm) {
		panic(fmt.Errorf("expected %T %d", ir.IfTerm, ir.IfTerm))
	}

	c := term.If

	if err := t.solveTerm(&c.Condition); err != nil {
		return err
	}

	if err := t.solveTerm(&c.Then); err != nil {
		return err
	}

	if c.Else != nil {
		if err := t.solveTerm(c.Else); err != nil {
			return err
		}
	}

	return nil
}

func (t *Inferencer) solveInjectionTerm(term *ir.IrTerm) error {
	c := term.Injection

	c.VariantType = t.solveType(c.VariantType)

	return t.solveTerm(&c.Value)
}

func (t *Inferencer) solveLambdaTerm(term *ir.IrTerm) error {
	if !term.Is(ir.LambdaTerm) {
		panic(fmt.Errorf("expected %T %d", ir.LambdaTerm, ir.LambdaTerm))
	}

	c := term.Lambda

	c.ArgType = t.solveType(c.ArgType)

	if err := t.solveTerm(&c.Body); err != nil {
		return err
	}

	return nil
}

func (t *Inferencer) solveLetTerm(term *ir.IrTerm) error {
	if !term.Is(ir.LetTerm) {
		panic(fmt.Errorf("expected %T %d", ir.LetTerm, ir.LetTerm))
	}

	c := term.Let

	c.VarType = t.solveType(c.VarType)

	return t.solveTerm(&c.Value)
}

func (t *Inferencer) solveMatchTerm(term *ir.IrTerm) error {
	if !term.Is(ir.MatchTerm) {
		panic(fmt.Errorf("expected %T %d", ir.MatchTerm, ir.MatchTerm))
	}

	c := term.Match

	if err := t.solveTerm(&c.Term); err != nil {
		return err
	}

	for i := range c.Arms {
		if err := t.solveTerm(&c.Arms[i].Body); err != nil {
			return err
		}
	}

	return nil
}

func (t *Inferencer) solveProjectionTerm(term *ir.IrTerm) error {
	if !term.Is(ir.ProjectionTerm) {
		panic(fmt.Errorf("expected %T %d", ir.ProjectionTerm, ir.ProjectionTerm))
	}

	c := term.Projection

	return t.solveTerm(&c.Term)
}

func (t *Inferencer) solveReturnTerm(term *ir.IrTerm) error {
	if !term.Is(ir.ReturnTerm) {
		panic(fmt.Errorf("expected %T %d", ir.ReturnTerm, ir.ReturnTerm))
	}

	c := term.Return

	return t.solveTerm(&c.Expr)
}

func (t *Inferencer) solveSetTerm(term *ir.IrTerm) error {
	if !term.Is(ir.SetTerm) {
		panic(fmt.Errorf("expected %T %d", ir.SetTerm, ir.SetTerm))
	}

	c := term.Set

	if err := t.solveTerm(&c.Term); err != nil {
		return err
	}

	for i := range c.Values {
		if err := t.solveTerm(&c.Values[i].Value); err != nil {
			return err
		}
	}

	return nil
}

func (t *Inferencer) solveStructTerm(term *ir.IrTerm) error {
	c := term.Struct

	for i := range c.Values {
		if err := t.solveTerm(&c.Values[i].Value); err != nil {
			return err
		}
	}

	return nil
}

func (t *Inferencer) solveTupleTerm(term *ir.IrTerm) error {
	c := term.Tuple

	for i := range c.Elems {
		if err := t.solveTerm(&c.Elems[i]); err != nil {
			return err
		}
	}

	return nil
}

func (t *Inferencer) solveTypeAbsTerm(term *ir.IrTerm) error {
	c := term.TypeAbs

	return t.solveTerm(&c.Body)
}

func (t *Inferencer) solveTermImpl(term *ir.IrTerm) error {
	switch {
	case term.Is(ir.AppTermTerm):
		return t.solveAppTermTerm(term)

	case term.Is(ir.AppTypeTerm):
		return t.solveAppTypeTerm(term)

	case term.Is(ir.AssignTerm):
		c := term.Assign

		if err := t.solveTerm(&c.Ret); err != nil {
			return err
		}

		if err := t.solveTerm(&c.Arg); err != nil {
			return err
		}

		return nil

	case term.Is(ir.BlockTerm):
		return t.solveBlockTerm(term)

	case term.Is(ir.ConstTerm):
		return nil

	case term.Is(ir.IfTerm):
		return t.solveIfTerm(term)

	case term.Is(ir.InjectionTerm):
		return t.solveInjectionTerm(term)

	case term.Is(ir.LambdaTerm):
		return t.solveLambdaTerm(term)

	case term.Is(ir.LetTerm):
		return t.solveLetTerm(term)

	case term.Is(ir.MatchTerm):
		return t.solveMatchTerm(term)

	case term.Is(ir.ProjectionTerm):
		return t.solveProjectionTerm(term)

	case term.Is(ir.ReturnTerm):
		return t.solveReturnTerm(term)

	case term.Is(ir.SetTerm):
		return t.solveSetTerm(term)

	case term.Is(ir.StructTerm):
		return t.solveStructTerm(term)

	case term.Is(ir.TupleTerm):
		return t.solveTupleTerm(term)

	case term.Is(ir.TypeAbsTerm):
		return t.solveTypeAbsTerm(term)

	case term.Is(ir.VarTerm):
		return nil

	default:
		panic(fmt.Errorf("unhandled %T %d", term.Case, term.Case))
	}
}

func (t *Inferencer) solveTerm(term *ir.IrTerm) error {
	if err := t.solveTermImpl(term); err != nil {
		return fmt.Errorf("%v\n  solving %s", err, *term)
	}

	if term.Type != nil {
		typ := t.solveType(*term.Type)
		term.Type = &typ
	}

	glog.V(1).Infof("solveTerm: %s |- %s", t.context, *term)
	return nil
}
