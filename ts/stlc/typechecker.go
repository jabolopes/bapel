package stlc

import (
	"fmt"
	"log"
	"os"

	"github.com/jabolopes/bapel/ir"
)

type Typechecker struct {
	*log.Logger
	context      Context
	bindPosition bool
}

func (t *Typechecker) withBindPosition(callback func() error) error {
	bind := t.bindPosition
	t.bindPosition = true
	defer func() { t.bindPosition = bind }()
	return callback()
}

func (t *Typechecker) isNumber(typ ir.IrType) error {
	if typ.Is(ir.NameType) &&
		(typ.Name == "i8" || typ.Name == "i16" || typ.Name == "i32" || typ.Name == "i64") {
		return nil
	}

	return fmt.Errorf("expected number type, e.g., i8, i16, i32, i64; got %v", typ)
}

func (t *Typechecker) reduceType(typ ir.IrType) (ir.IrType, error) {
	reducer := typeReducer{t.Logger, t.context}
	return reducer.reduce(typ)
}

func (t *Typechecker) InferTerm(term *ir.IrTerm) error {
	inferencer := Inferencer{t.Logger, t.context}
	return inferencer.inferTerm(term)
}

func (t *Typechecker) InferFunction(function *ir.IrFunction) error {
	inferencer := Inferencer{t.Logger, t.context}
	return inferencer.inferFunction(function)
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
	retContext, err := t.context.AddBind(NewTermBind(decl.Term.ID, decl.Term.Type, DefSymbol))
	if err != nil {
		return origContext, err
	}

	t.context = retContext
	if t.context, err = t.context.enterFunction(function.ID, function.TypeVars, function.Args, function.Rets); err != nil {
		return origContext, err
	}

	if err := t.TypecheckTerm(&function.Body); err != nil {
		return origContext, err
	}

	{
		var retTypes []ir.IrType
		for _, ret := range function.Rets {
			retTypes = append(retTypes, ret.Term.Type)
		}
		functionReturnType := ir.NewTupleType(retTypes)

		switch function.Body.Case {
		case ir.BlockTerm:
			// Check all return terms have the correct function return type.
			var retTypes []ir.IrType
			for _, ret := range function.Rets {
				retTypes = append(retTypes, ret.Term.Type)
			}
			functionReturnType := ir.NewTupleType(retTypes)

			returns := allReturns(function.Body)
			for _, ret := range returns {
				returnType := *ret.Return.Expr.Type
				if err := t.subtype(functionReturnType, returnType); err != nil {
					return origContext, fmt.Errorf("%v:\n%v", ret.Pos, err)
				}
			}

			// Check all function exits have the correct type.
			last := lastTerms(function.Body)
			for _, term := range last {
				if term.Is(ir.ReturnTerm) {
					return origContext, fmt.Errorf("%v:\n redundant 'return' statement as the last term of a function", term.Pos)
				}

				if err := t.subtype(functionReturnType, *term.Type); err != nil {
					return origContext, fmt.Errorf("%v:\n%v", term.Pos, err)
				}
			}

			if len(last) == 0 {
				return origContext, fmt.Errorf("%v:\nexpected non-empty function block", function.Body.Pos)
			}

		default:
			if err := t.subtype(functionReturnType, *function.Body.Type); err != nil {
				return origContext, fmt.Errorf("%v:\n%v", function.Body.Pos, err)
			}
		}
	}

	return retContext, nil
}

func NewTypechecker(context Context) *Typechecker {
	return &Typechecker{
		log.New(os.Stderr, "DEBUG ", 0),
		context,
		false, /* bindPosition */
	}
}
