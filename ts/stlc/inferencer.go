package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

func probeType(term ir.IrTerm) (ir.IrType, bool) {
	if term.Type != nil {
		if term.Type.Case == ir.TupleType {
			return term.Type.Tuple[0], true
		}

		return *term.Type, true
	}

	for _, elem := range term.Tuple {
		if elem.Type != nil {
			return *elem.Type, true
		}
	}

	return ir.IrType{}, false
}

type Inferencer struct {
	context Context
}

func (t *Inferencer) inferImpl(term *ir.IrTerm, checkType *ir.IrType) error {
	switch term.Case {
	case ir.AssignTerm:
		c := term.Assign
		if err := t.inferImpl(&c.Ret, nil /* checkType */); err != nil {
			return err
		}

		if err := t.inferImpl(&c.Arg, nil /* checkType */); err != nil {
			return err
		}

		if c.Arg.Type == nil && c.Ret.Type != nil {
			if err := t.inferImpl(&c.Arg, c.Ret.Type); err != nil {
				return err
			}
		}

		return nil

	case ir.BlockTerm:
		c := term.Block
		for i := range c.Terms {
			if err := t.inferImpl(&c.Terms[i], nil /* checkType */); err != nil {
				return err
			}
		}
		return nil

	case ir.CallTerm:
		c := term.Call
		if err := t.inferImpl(&c.Arg, nil /* checkType */); err != nil {
			return err
		}

		if len(c.Types) == 0 {
			if ir.IsOperator(c.ID) {
				typ, ok := probeType(c.Arg)
				if ok {
					c.Types = []ir.IrType{typ}
				} else if checkType != nil {
					c.Types = []ir.IrType{*checkType}
				}
			}
		}
		return nil

	case ir.IfTerm:
		c := term.If

		if err := t.inferImpl(&c.Condition, nil /* checkType */); err != nil {
			return err
		}
		if err := t.inferImpl(&c.Then, nil /* checkType */); err != nil {
			return err
		}
		if c.Else != nil {
			if err := t.inferImpl(c.Else, nil /* checkType */); err != nil {
				return err
			}
		}
		return nil

	case ir.IndexGetTerm:
		c := term.IndexGet
		if err := t.inferImpl(&c.Obj, nil /* checkType */); err != nil {
			return err
		}
		return t.inferImpl(&c.Index, nil /* checkType */)

	case ir.IndexSetTerm:
		c := term.IndexSet
		if err := t.inferImpl(&c.Obj, nil /* checkType */); err != nil {
			return err
		}
		if err := t.inferImpl(&c.Index, nil /* checkType */); err != nil {
			return err
		}
		return t.inferImpl(&c.Value, nil /* checkType */)

	case ir.LetTerm:
		return nil

	case ir.StatementTerm:
		c := term.Statement
		return t.inferImpl(&c.Term, nil /* checkType */)

	case ir.TokenTerm:
		c := term.Token
		if c.Case != parser.IDToken {
			return nil
		}

		bind, ok := t.context.LookupBind(c.Text, FindAny)
		if !ok || bind.Case != DeclBind || bind.Decl.Case != ir.TermDecl {
			return nil
		}

		term.Type = &bind.Decl.Term.Type
		return nil

	case ir.TupleTerm:
		typ := func() *ir.IrType {
			t := ir.NewTupleType(nil)
			return &t
		}()

		for i := range term.Tuple {
			if err := t.inferImpl(&term.Tuple[i], nil /* checkType */); err != nil {
				return err
			}

			if term.Tuple[i].Type == nil {
				typ = nil
			} else if typ != nil {
				typ.Tuple = append(typ.Tuple, *term.Tuple[i].Type)
			}
		}

		term.Type = typ
		return nil

	case ir.WidenTerm:
		return t.inferImpl(&term.Widen.Term, nil /* checkType */)

	default:
		panic(fmt.Errorf("unhandled %T %d", term.Case, term.Case))
	}
}

func (t *Inferencer) Infer(term *ir.IrTerm) error {
	return t.inferImpl(term, nil /* checkType */)
}

func NewInferencer(context Context) *Inferencer {
	return &Inferencer{context}
}
