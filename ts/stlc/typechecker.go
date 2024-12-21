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

func (t *Typechecker) TypecheckTerm(term *ir.IrTerm) error {
	return t.typecheck(term)
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

	return retContext, nil
}

func NewTypechecker(context Context) *Typechecker {
	return &Typechecker{
		log.New(os.Stderr, "DEBUG ", 0),
		context,
		false, /* bindPosition */
	}
}
