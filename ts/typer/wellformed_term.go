package typer

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

func WellformedTerm(context Context, term ir.IrTerm) error {
	switch term.Case {
	case ir.AppTermTerm:
		c := term.AppTerm
		if err := WellformedTerm(context, c.Fun); err != nil {
			return err
		}
		return WellformedTerm(context, c.Arg)

	case ir.AppTypeTerm:
		c := term.AppType
		if err := WellformedTerm(context, c.Fun); err != nil {
			return err
		}

		// TODO: Finish when c.Arg and WellformedType use the same type
		// (e.g., IrType vs typer.Type).
		//
		// return WellformedType(context, c.Arg)

		return nil

	case ir.AssignTerm:
		c := term.Assign
		if err := WellformedTerm(context, c.Arg); err != nil {
			return err
		}
		return WellformedTerm(context, c.Ret)

	case ir.BlockTerm:
		c := term.Block
		for _, term := range c.Terms {
			if err := WellformedTerm(context, term); err != nil {
				return err
			}
		}
		return nil

	case ir.IfTerm:
		return WellformedTerm(context, term.If.Condition)

	case ir.IndexGetTerm:
		c := term.IndexGet
		if err := WellformedTerm(context, c.Obj); err != nil {
			return err
		}
		return WellformedTerm(context, c.Index)

	case ir.IndexSetTerm:
		c := term.IndexSet
		if err := WellformedTerm(context, c.Obj); err != nil {
			return err
		}
		if err := WellformedTerm(context, c.Index); err != nil {
			return err
		}
		return WellformedTerm(context, c.Value)

	case ir.LiteralTerm:
		if term.Literal.Is(ir.IDLiteral) && !context.ContainsTermBind(term.Literal.Text) {
			return fmt.Errorf("term %s is not wellformed", term)
		}
		return nil

	case ir.TupleTerm:
		for _, t := range term.Tuple {
			if err := WellformedTerm(context, t); err != nil {
				return err
			}
		}
		return nil

	default:
		panic(fmt.Errorf("unhandled %T %d", term.Case, term.Case))
	}
}
