package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

func (t *Typechecker) applyImpl(term *ir.IrTerm) (ir.IrType, error) {
	switch {
	case term.Is(ir.AppTypeTerm) &&
		term.AppType.Fun.Type != nil && term.AppType.Fun.Type.Is(ir.ForallType) &&
		(isWellformedType(t.context, term.AppType.Arg) == nil):
		forall := term.AppType.Fun.Type.Forall
		return ir.SubstituteType(forall.Type, ir.NewVarType(forall.Var), term.AppType.Arg), nil

	case term.Is(ir.AppTermTerm) &&
		term.AppTerm.Fun.Type != nil && term.AppTerm.Fun.Type.Is(ir.FunType):
		funType := term.AppTerm.Fun.Type.Fun
		if err := t.typecheck(&term.AppTerm.Arg); err != nil {
			return ir.IrType{}, err
		}

		if err := t.subtype(*term.AppTerm.Arg.Type, funType.Arg); err != nil {
			return ir.IrType{}, err
		}

		return funType.Ret, nil

	default:
		return ir.IrType{}, fmt.Errorf("failed to apply")
	}
}

func (t *Typechecker) apply(term *ir.IrTerm) (ir.IrType, error) {
	termType, err := t.applyImpl(term)
	if err != nil {
		return ir.IrType{}, err
	}

	if term != nil {
		term.Type = &termType
	}
	return termType, nil
}
