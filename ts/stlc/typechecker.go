package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

type Typechecker struct {
	context      Context
	bindPosition bool
}

func (t *Typechecker) withBindPosition(callback func() error) error {
	bind := t.bindPosition
	t.bindPosition = true
	defer func() { t.bindPosition = bind }()
	return callback()
}

func (t *Typechecker) isBool(typ ir.IrType) error {
	if !typ.Is(ir.NameType) || typ.Name != "bool" {
		return fmt.Errorf("expected bool; got %s", typ)
	}
	return nil
}

func (t *Typechecker) reduceAndPredicateType(typ ir.IrType) (ir.IrType, error) {
	reducer := typeReducer{}
	typ = reducer.reduce(t.context, typ)

	predicator := typePredicator{t.context, nil /* tvars */}

	newType, err := predicator.predicate(typ)
	if err != nil {
		return ir.IrType{}, err
	}

	return ir.ForallVars(predicator.tvars, newType), nil
}

func (t *Typechecker) InferFunction(function *ir.IrFunction) (Context, error) {
	inferencer := NewInferencer(t.context)

	context, err := inferencer.inferFunction(function)
	if err != nil {
		return context, fmt.Errorf("%v:\n%v", function.Pos, err)
	}

	return context, nil
}

func (t *Typechecker) TypecheckTerm(term *ir.IrTerm) error {
	if err := t.typecheck(term); err != nil {
		return fmt.Errorf("%v:\n%v", term.Pos, err)
	}
	return nil
}

func (t *Typechecker) TypecheckFunction(function *ir.IrFunction) (Context, error) {
	origContext := t.context

	decl := function.Decl()

	var err error
	retContext, err := t.context.AddBind(NewTermDefBind(decl.Term.ID, decl.Term.Type))
	if err != nil {
		return origContext, err
	}

	t.context = retContext
	if t.context, err = t.context.enterFunction(function.TypeVars, function.Args); err != nil {
		return origContext, err
	}

	if err := t.TypecheckTerm(&function.Body); err != nil {
		return origContext, err
	}

	{
		switch function.Body.Case {
		case ir.BlockTerm:
			// Check all return terms have the correct function return type.
			returns := allReturns(function.Body)
			for _, ret := range returns {
				returnType := *ret.Type
				if err := t.subtype(function.RetType, returnType); err != nil {
					return origContext, fmt.Errorf("%v:\n%v", ret.Pos, err)
				}
			}

			// Check all function exits have the correct type.
			last := lastTerms(&function.Body)
			for _, term := range last {
				if term.Is(ir.ReturnTerm) {
					return origContext, fmt.Errorf("%v:\n redundant 'return' statement as the last term of a function", term.Pos)
				}

				if err := t.subtype(function.RetType, *term.Type); err != nil {
					return origContext, fmt.Errorf("%v:\n%v", term.Pos, err)
				}
			}

			if len(last) == 0 {
				return origContext, fmt.Errorf("%v:\nexpected non-empty function block", function.Body.Pos)
			}

		default:
			if err := t.subtype(function.RetType, *function.Body.Type); err != nil {
				return origContext, fmt.Errorf("%v:\n%v", function.Body.Pos, err)
			}
		}
	}

	return retContext, nil
}

func NewTypechecker(context Context) *Typechecker {
	return &Typechecker{
		context,
		false, /* bindPosition */
	}
}
