package ir

import (
	"fmt"

	"github.com/jabolopes/bapel/parser"
)

func isOperator(id string) bool {
	return id == "+" || id == "-" || id == "*" || id == "/"
}

func probeType(term IrTerm) (IrType, bool) {
	if term.Type != nil {
		if term.Type.Case == TupleType {
			return term.Type.Tuple[0], true
		}

		return *term.Type, true
	}

	for _, elem := range term.Tuple {
		if elem.Type != nil {
			return *elem.Type, true
		}
	}

	return IrType{}, false
}

type IrInferencer struct {
	context *IrContext
}

func (t *IrInferencer) inferImpl(term *IrTerm, checkType *IrType) error {
	switch term.Case {
	case AssignTerm:
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

	case CallTerm:
		c := term.Call
		if err := t.inferImpl(&c.Arg, nil /* checkType */); err != nil {
			return err
		}

		if len(c.Types) == 0 {
			if isOperator(c.ID) {
				typ, ok := probeType(c.Arg)
				if ok {
					c.Types = []IrType{typ}
				} else if checkType != nil {
					c.Types = []IrType{*checkType}
				}
			}
		}
		return nil

	case IfTerm:
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

	case IndexGetTerm:
		c := term.IndexGet
		if err := t.inferImpl(&c.Obj, nil /* checkType */); err != nil {
			return err
		}
		return t.inferImpl(&c.Index, nil /* checkType */)

	case IndexSetTerm:
		c := term.IndexSet
		if err := t.inferImpl(&c.Obj, nil /* checkType */); err != nil {
			return err
		}
		if err := t.inferImpl(&c.Index, nil /* checkType */); err != nil {
			return err
		}
		return t.inferImpl(&c.Value, nil /* checkType */)

	case LetTerm:
		return nil

	case StatementTerm:
		c := term.Statement
		for i := range c.Terms {
			if err := t.inferImpl(&c.Terms[i], nil /* checkType */); err != nil {
				return err
			}
		}
		return nil

	case TokenTerm:
		c := term.Token
		if c.Case != parser.IDToken {
			return nil
		}

		bind, ok := t.context.lookupBind(c.Text, FindAny)
		if !ok || bind.Case != DeclBind || bind.Decl.Case != TermDecl {
			return nil
		}

		term.Type = &bind.Decl.Term.Type
		return nil

	case TupleTerm:
		typ := func() *IrType {
			t := NewTupleType(nil)
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

	case WidenTerm:
		return t.inferImpl(&term.Widen.Term, nil /* checkType */)

	default:
		panic(fmt.Errorf("unhandled %T %d", term.Case, term.Case))
	}
}

func (t *IrInferencer) Infer(term *IrTerm) error {
	return t.inferImpl(term, nil /* checkType */)
}

func NewInferencer(context *IrContext) *IrInferencer {
	return &IrInferencer{context}
}
