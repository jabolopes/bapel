package ir

import (
	"fmt"

	"github.com/jabolopes/bapel/parser"
)

type IrTypechecker struct {
	context      *IrContext
	widen        bool
	bindPosition bool
}

func (t *IrTypechecker) withBindPosition(callback func() (IrType, error)) (IrType, error) {
	bind := t.bindPosition
	t.bindPosition = true
	defer func() { t.bindPosition = bind }()
	return callback()
}

func (t *IrTypechecker) withWiden(callback func() error) error {
	widen := t.widen
	t.widen = true
	defer func() { t.widen = widen }()
	return callback()
}

func (t *IrTypechecker) subtype(left, right IrType) error {
	switch {
	case left.Case == ArrayType && right.Case == ArrayType:
		if err := t.subtype(left.Array.ElemType, right.Array.ElemType); err != nil {
			return fmt.Errorf("mismatch in array element types: %v", err)
		}

		if left.Array.Size != right.Array.Size {
			return fmt.Errorf("expected array with %d elements; got %d elements", left.Array.Size, right.Array.Size)
		}

		return nil

	case left.Case == FunType && right.Case == FunType:
		// B1 <: A1
		if err := t.subtype(NewTupleType(right.Fun.Args), NewTupleType(left.Fun.Args)); err != nil {
			return err
		}

		// A2 <: B2
		if err := t.subtype(NewTupleType(left.Fun.Rets), NewTupleType(right.Fun.Rets)); err != nil {
			return err
		}

		return nil

	case left.Case == IntType && right.Case == IntType:
		if t.widen {
			if left.Int < right.Int {
				return fmt.Errorf("expected type %s or wider; got %s", left.Int, right.Int)
			}
		} else {
			if left.Int != right.Int {
				return fmt.Errorf("expected type %s; got %s", left.Int, right.Int)
			}
		}

		return nil

	case left.Case == StructType && right.Case == StructType:
		if len(left.Fields()) != len(right.Fields()) {
			return fmt.Errorf("expected %d fields; got %d", len(left.Fields()), len(right.Fields()))
		}

		for i := range left.Fields() {
			if left.Fields()[i].ID != right.Fields()[i].ID {
				return fmt.Errorf("expected field names %v; got %v", left.FieldIDs(), right.FieldIDs())
			}

			if err := t.subtype(left.Fields()[i].Type, right.Fields()[i].Type); err != nil {
				return err
			}
		}

		return nil

	case left.Case == TupleType && right.Case == TupleType:
		if len(left.Tuple) != len(right.Tuple) {
			return fmt.Errorf("expected %d elements; got %d", len(left.Tuple), len(right.Tuple))
		}

		for i := range left.Tuple {
			if err := t.subtype(left.Tuple[i], right.Tuple[i]); err != nil {
				return err
			}
		}

		return nil

	case left.Case == VarType && right.Case == VarType && left.Var == right.Var:
		return nil

	case left.Case == IDType && right.Case == IDType:
		leftDecl, err := t.context.getDecl(left.ID, FindAny)
		if err != nil {
			return err
		}

		rightDecl, err := t.context.getDecl(right.ID, FindAny)
		if err != nil {
			return err
		}

		return t.MatchesDecl(leftDecl, rightDecl)

	default:
		return fmt.Errorf("expected type %s (%s); got %s (%s)", left.Case, left, right.Case, right)
	}
}

func (t *IrTypechecker) MatchesDecl(left, right IrDecl) error {
	if left.ID != right.ID {
		return fmt.Errorf("expected declaration of %s; got %s", left.ID, right.ID)
	}

	if left.Case != right.Case {
		return fmt.Errorf("in declaration of %q: expected %s; got %s", left.ID, left.Case, right.Case)
	}

	if err := t.subtype(left.Type, right.Type); err != nil {
		return fmt.Errorf("in declaration of %q: %v", left.ID, err)
	}

	return nil
}

func (t *IrTypechecker) SynthesizeTerm(term IrTerm) (IrType, error) {
	switch term.Case {
	case AssignTerm:
		assign := term.Assign

		retType, err := t.withBindPosition(func() (IrType, error) {
			return t.SynthesizeTerm(assign.Ret)
		})
		if err != nil {
			return IrType{}, err
		}

		if err := t.CheckTerm(assign.Arg, retType); err != nil {
			return IrType{}, err
		}

		return retType, nil

	case CallTerm:
		call := term.Call

		formal, err := t.SynthesizeTerm(NewTokenTerm(parser.NewIDToken(call.ID)))
		if err != nil {
			return IrType{}, err
		}

		if formal.Case != FunType {
			return IrType{}, fmt.Errorf("expected function type; got %s", formal)
		}

		if err := t.CheckTerm(NewTupleTerm(call.Args), NewTupleType(formal.Fun.Args)); err != nil {
			return IrType{}, err
		}

		return NewTupleType(formal.Fun.Rets), nil

	case IndexGetTerm:
		indexableType, err := t.SynthesizeTermFull(term.IndexGet.Term)
		if err != nil {
			return IrType{}, err
		}

		var index *int64
		var fieldID *string
		if term.IndexGet.Index.Case == TokenTerm {
			switch term.IndexGet.Index.Token.Case {
			case parser.NumberToken:
				index = &term.IndexGet.Index.Token.Value
			case parser.IDToken:
				fieldID = &term.IndexGet.Index.Token.Text
			}
		}

		switch {
		case indexableType.Is(StructType) && index != nil:
			field, ok := indexableType.FieldByIndex(int(*index))
			if !ok {
				return IrType{}, fmt.Errorf("field %d is not a valid field of struct type %s", *index, indexableType)
			}

			term.IndexGet.Field = field.ID
			return field.Type, nil

		case indexableType.Is(StructType) && fieldID != nil:
			field, ok := indexableType.FieldByID(*fieldID)
			if !ok {
				return IrType{}, fmt.Errorf("field %q is not a valid field of struct type %s", *fieldID, indexableType)
			}

			term.IndexGet.Field = field.ID
			return field.Type, nil

		case indexableType.Is(StructType):
			return IrType{}, fmt.Errorf("expected field identifier or number literal to index struct %s", indexableType)

		case indexableType.Is(ArrayType) && index != nil:
			if *index < 0 || *index >= int64(indexableType.Array.Size) {
				return IrType{}, fmt.Errorf("index %d is out of bounds", *index)
			}
			return indexableType.Array.ElemType, nil

		case indexableType.Is(ArrayType):
			// TODO: This should check any integer (or Number) instead of just i64.
			if err := t.CheckTerm(term.IndexGet.Index, NewIntType(I64)); err != nil {
				return IrType{}, err
			}

			return indexableType.Array.ElemType, nil

		default:
			return IrType{}, fmt.Errorf("expected indexable type (e.g., array, struct, etc); got %s", indexableType)
		}

	case IndexSetTerm:
		var index *int64
		var fieldID *string
		if term.IndexSet.Index.Case == TokenTerm {
			switch term.IndexSet.Index.Token.Case {
			case parser.NumberToken:
				index = &term.IndexSet.Index.Token.Value
			case parser.IDToken:
				fieldID = &term.IndexSet.Index.Token.Text
			}
		}

		indexableType, err := t.SynthesizeTermFull(term.IndexSet.Ret)
		if err != nil {
			return IrType{}, err
		}

		switch {
		case indexableType.Is(StructType) && index != nil:
			field, ok := indexableType.FieldByIndex(int(*index))
			if !ok {
				return IrType{}, fmt.Errorf("field %d is not a valid field of struct type %s", *index, indexableType)
			}

			term.IndexSet.Field = field.ID
			return field.Type, nil

		case indexableType.Is(StructType) && fieldID != nil:
			field, ok := indexableType.FieldByID(*fieldID)
			if !ok {
				return IrType{}, fmt.Errorf("field %q is not a valid field of struct type %s", *fieldID, indexableType)
			}

			term.IndexSet.Field = field.ID
			return field.Type, nil

		case indexableType.Is(StructType):
			return IrType{}, fmt.Errorf("expected field identifier or number literal to index struct %s", indexableType)

		case indexableType.Is(ArrayType) && index != nil:
			if *index < 0 || *index >= int64(indexableType.Array.Size) {
				return IrType{}, fmt.Errorf("index %d is out of bounds", *index)
			}
			return indexableType.Array.ElemType, nil

		case indexableType.Is(ArrayType):
			// TODO: This should check any integer (or Number) instead of just i64.
			if err := t.CheckTerm(term.IndexSet.Index, NewIntType(I64)); err != nil {
				return IrType{}, err
			}
			if err := t.CheckTerm(term.IndexSet.Arg, indexableType.Array.ElemType); err != nil {
				return IrType{}, err
			}
			return NewTupleType(nil), nil

		default:
			return IrType{}, fmt.Errorf("expected indexable type (e.g., array); got %s", indexableType)
		}

	case StatementTerm:
		if _, err := t.SynthesizeTerm(term.Statement.Term); err != nil {
			return IrType{}, err
		}
		return NewTupleType(nil), nil

	case TokenTerm:
		token := term.Token
		switch token.Case {
		case parser.IDToken:
			return t.context.getType(token.Text, FindAny)

		case parser.NumberToken:
			return IrType{}, fmt.Errorf("cannot synthesize type for number token")

		default:
			panic(fmt.Errorf("unhandled token %d", token.Case))
		}

	case TupleTerm:
		types := make([]IrType, len(term.Tuple))
		for i := range term.Tuple {
			var err error
			types[i], err = t.SynthesizeTerm(term.Tuple[i])
			if err != nil {
				return IrType{}, err
			}
		}
		return NewTupleType(types), nil

	default:
		panic(fmt.Errorf("unhandled IrTerm %d", term.Case))
	}
}

func (t *IrTypechecker) SynthesizeTermFull(term IrTerm) (IrType, error) {
	typ, err := t.SynthesizeTerm(term)
	if err != nil {
		return IrType{}, err
	}

	switch typ.Case {
	case IDType:
		return t.context.getType(typ.ID, FindAny)
	default:
		return typ, nil
	}
}

func (t *IrTypechecker) CheckTerm(term IrTerm, typ IrType) error {
	switch {
	// Case AssignTerm: handled by default case.

	case term.Case == IfTerm:
		condition := term.If.Condition

		conditionType, err := t.SynthesizeTerm(condition)
		if err != nil {
			return err
		}

		if !conditionType.Is(IntType) {
			return fmt.Errorf("expected integer type; got %s", conditionType)
		}

		return t.subtype(NewTupleType(nil), typ)

	case term.Case == StatementTerm:
		if _, err := t.SynthesizeTerm(term.Statement.Term); err != nil {
			return err
		}
		return t.subtype(NewTupleType(nil), typ)

	case term.Case == TokenTerm && term.Token.Case == parser.NumberToken && !t.bindPosition:
		if !typ.Is(IntType) {
			return fmt.Errorf("expected type %s; got %q", typ, term.Token.Text)
		}
		return nil

	case term.Case == TokenTerm && t.bindPosition:
		switch token := term.Token; token.Case {
		case parser.IDToken:
			actualDecl, err := t.context.getDecl(token.Text, FindAny)
			if err != nil {
				return err
			}
			if actualDecl.Case != TermDecl {
				return fmt.Errorf("expected symbol declared as %s; got %q", TermDecl, actualDecl.Case)
			}
			return t.subtype(actualDecl.Type, typ)

		case parser.NumberToken:
			return fmt.Errorf("expected symbol declared as %s; got number literal", TermDecl)

		default:
			panic(fmt.Errorf("unhandled token %d", token.Case))
		}

	case term.Case == OpUnaryTerm:
		if !typ.Is(IntType) {
			return fmt.Errorf("expected integer type; got %s", typ)
		}

		return t.CheckTerm(term.OpUnary.Term, typ)

	case term.Case == OpBinaryTerm:
		if !typ.Is(IntType) {
			return fmt.Errorf("expected integer type; got %s", typ)
		}

		if err := t.CheckTerm(term.OpBinary.Left, typ); err != nil {
			return err
		}

		return t.CheckTerm(term.OpBinary.Right, typ)

	case term.Case == WidenTerm:
		return t.withWiden(func() error {
			return t.CheckTerm(term.Widen.Term, typ)
		})

	default:
		// Sub:
		//   e <= B

		// e => A
		got, err := t.SynthesizeTerm(term)
		if err != nil {
			return err
		}

		// A <: B
		return t.subtype(got, typ)
	}
}

func (t *IrTypechecker) TypecheckTerm(term IrTerm) error {
	return t.CheckTerm(term, NewTupleType(nil))
}

func NewIrTypechecker(context *IrContext) *IrTypechecker {
	return &IrTypechecker{
		context,
		false, /* widen */
		false, /* bindPosition */
	}
}
