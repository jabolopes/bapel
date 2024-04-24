package stlc

import (
	"fmt"
	"strings"

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

func (t *Inferencer) inferApply(term *ir.IrTerm, typ ir.IrType, types []ir.IrType) error {
	switch {
	case typ.Case == ir.ForallType && len(types) == len(typ.Forall.Vars):
		for i, tvar := range typ.Forall.Vars {
			tvar = strings.TrimPrefix(tvar, "'")
			typeVar := ir.NewVarType(tvar)
			typ = ir.SubstituteType(typ, typeVar, types[i])
		}

		return t.inferApply(term, typ.Forall.Type, nil /* types */)

	case typ.Case == ir.FunType && len(types) == 0:
		if err := t.inferImpl(&term.Call.Arg, &typ.Fun.Arg); err != nil {
			return err
		}

		typ := typ.Fun.Ret
		term.Type = &typ
		return nil

	default:
		return nil
	}
}

func (t *Inferencer) inferImpl(term *ir.IrTerm, expectType *ir.IrType) error {
	switch {
	case term.Case == ir.AssignTerm:
		c := term.Assign
		if err := t.inferImpl(&c.Ret, nil /* expectType */); err != nil {
			return err
		}

		if err := t.inferImpl(&c.Arg, c.Ret.Type); err != nil {
			return err
		}

		return nil

	case term.Case == ir.BlockTerm:
		c := term.Block
		for i := range c.Terms {
			if err := t.inferImpl(&c.Terms[i], nil /* expectType */); err != nil {
				return err
			}
		}
		return nil

	case term.Case == ir.CallTerm && ir.IsOperator(term.Call.ID) && len(term.Call.Types) == 1:
		c := term.Call

		argType := ir.NewTupleType([]ir.IrType{c.Types[0], c.Types[0]})
		if err := t.inferImpl(&c.Arg, &argType); err != nil {
			return err
		}

		typ, ok := probeType(c.Arg)
		if ok {
			c.Types = []ir.IrType{typ}
		}
		return nil

	case term.Case == ir.CallTerm && ir.IsOperator(term.Call.ID) && len(term.Call.Types) > 1:
		return nil

	case term.Case == ir.CallTerm && ir.IsOperator(term.Call.ID) && expectType != nil:
		c := term.Call

		argType := ir.NewTupleType([]ir.IrType{*expectType, *expectType})
		if err := t.inferImpl(&c.Arg, &argType); err != nil {
			return err
		}

		typ, ok := probeType(c.Arg)
		if ok {
			c.Types = []ir.IrType{typ}
		}
		return nil

	case term.Case == ir.CallTerm && ir.IsOperator(term.Call.ID) && expectType == nil:
		c := term.Call

		tokenTerm := ir.NewTokenTerm(parser.NewIDToken(c.ID))
		if err := t.inferImpl(&tokenTerm, nil /* expectType */); err != nil {
			return err
		}

		if err := t.inferImpl(&c.Arg, nil /* expectType */); err != nil {
			return err
		}

		typ, ok := probeType(c.Arg)
		if ok {
			c.Types = []ir.IrType{typ}
		}
		return nil

	case term.Case == ir.CallTerm:
		c := term.Call

		idTerm := ir.NewTokenTerm(parser.NewIDToken(c.ID))
		if err := t.inferImpl(&idTerm, nil /* expectType */); err != nil {
			return err
		}

		if err := t.inferImpl(&c.Arg, nil /* expectType */); err != nil {
			return err
		}

		if idTerm.Type != nil {
			return t.inferApply(term, *idTerm.Type, c.Types)
		}

		return nil

	case term.Case == ir.IfTerm:
		c := term.If

		if err := t.inferImpl(&c.Condition, nil /* expectType */); err != nil {
			return err
		}

		if len(c.Types) == 0 && c.Condition.Type != nil {
			c.Types = []ir.IrType{*c.Condition.Type}
		}

		if err := t.inferImpl(&c.Then, nil /* expectType */); err != nil {
			return err
		}
		if c.Else != nil {
			if err := t.inferImpl(c.Else, nil /* expectType */); err != nil {
				return err
			}
		}
		return nil

	case term.Case == ir.IndexGetTerm:
		c := term.IndexGet
		if err := t.inferImpl(&c.Obj, nil /* expectType */); err != nil {
			return err
		}
		return t.inferImpl(&c.Index, nil /* expectType */)

	case term.Case == ir.IndexSetTerm:
		c := term.IndexSet
		if err := t.inferImpl(&c.Obj, nil /* expectType */); err != nil {
			return err
		}
		if err := t.inferImpl(&c.Index, nil /* expectType */); err != nil {
			return err
		}
		return t.inferImpl(&c.Value, nil /* expectType */)

	case term.Case == ir.LetTerm:
		c := term.Let
		var err error
		if t.context, err = t.context.AddBind(NewDeclBind(DefSymbol, c.Decl)); err != nil {
			return err
		}
		typ := c.Decl.Type()
		term.Type = &typ
		return nil

	case term.Case == ir.TokenTerm:
		c := term.Token
		if c.Case != parser.IDToken {
			term.Type = expectType
			return nil
		}

		bind, ok := t.context.LookupBind(c.Text, FindAny)
		if !ok || bind.Case != DeclBind || bind.Decl.Case != ir.TermDecl {
			return nil
		}

		term.Type = &bind.Decl.Term.Type
		return nil

	case term.Case == ir.TupleTerm &&
		expectType != nil && expectType.Is(ir.TupleType) &&
		len(expectType.Tuple) == len(term.Tuple):

		typ := func() *ir.IrType {
			t := ir.NewTupleType(nil)
			return &t
		}()

		for i := range term.Tuple {
			if err := t.inferImpl(&term.Tuple[i], &expectType.Tuple[i]); err != nil {
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

	case term.Case == ir.TupleTerm:
		typ := func() *ir.IrType {
			t := ir.NewTupleType(nil)
			return &t
		}()

		for i := range term.Tuple {
			if err := t.inferImpl(&term.Tuple[i], nil /* expectType */); err != nil {
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

	default:
		panic(fmt.Errorf("unhandled %T %d", term.Case, term.Case))
	}
}

func (t *Inferencer) InferTerm(term *ir.IrTerm) error {
	return t.inferImpl(term, nil /* expectType */)
}

func (t *Inferencer) InferFunction(function *ir.IrFunction) error {
	var err error
	t.context, err = t.context.AddBind(NewDeclBind(DefSymbol, function.Decl()))
	if err != nil {
		return err
	}

	if t.context, err = t.context.enterFunction(function.ID, function.TypeVars, function.Args, function.Rets); err != nil {
		return err
	}

	return t.InferTerm(&function.Body)
}

func NewInferencer(context Context) *Inferencer {
	return &Inferencer{context}
}
