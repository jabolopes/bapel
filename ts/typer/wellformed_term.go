package typer

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func WellformedTerm(context Context, term ir.IrTerm) error {
	switch term.Case {
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

	case ir.CallTerm:
		c := term.Call
		if !context.ContainsTermBind(c.ID) {
			return fmt.Errorf("term %s is not wellformed: ID %s is not wellformed", term, c.ID)
		}

		// TODO: Finish when IrType is replaced with typer.Type.
		//
		// for _, typ := range c.Types {
		// 	if err := WellformedType(context, typ); err != nil {
		// 		return fmt.Sprintf("term %s is not wellformed: %v", err)
		// 	}
		// }

		return WellformedTerm(context, c.Arg)

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

	case ir.TokenTerm:
		if term.Token.Case == parser.IDToken &&
			!context.ContainsTermBind(term.Token.Text) {
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

	case ir.WidenTerm:
		return WellformedTerm(context, term.Widen.Term)

	default:
		panic(fmt.Errorf("unhandled %T %d", term.Case, term.Case))
	}
}
