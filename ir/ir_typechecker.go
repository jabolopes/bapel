package ir

import (
	"fmt"
	"sort"

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

func (t *IrTypechecker) MatchesArrayType(formal, actual IrArrayType) error {
	if err := t.MatchesType(formal.ElemType, actual.ElemType); err != nil {
		return fmt.Errorf("mismatch in array element types: %v", err)
	}

	if formal.Size != actual.Size {
		return fmt.Errorf("expected array with %d elements; got %d elements", formal.Size, actual.Size)
	}

	return nil
}

func (t *IrTypechecker) MatchesFunctionType(formal, actual IrFunctionType) error {
	if len(formal.Args) != len(actual.Args) {
		return fmt.Errorf("expected function with %d argument(s); got %q", len(formal.Args), actual.Args)
	}

	if len(formal.Rets) != len(actual.Rets) {
		return fmt.Errorf("expected function with %d return value(s); got %q", len(formal.Rets), actual.Rets)
	}

	for i := range formal.Args {
		if err := t.MatchesType(formal.Args[i], actual.Args[i]); err != nil {
			return fmt.Errorf("in function argument %d: %v", i+1, err)
		}
	}

	for i := range formal.Rets {
		if err := t.MatchesType(formal.Rets[i], actual.Rets[i]); err != nil {
			return fmt.Errorf("in return value %d: %v", i, err)
		}
	}

	return nil
}

func (t *IrTypechecker) MatchesIntType(formal, actual IrIntType) error {
	if t.widen {
		if formal < actual {
			return fmt.Errorf("expected type %s or wider; got %s", formal, actual)
		}
	} else {
		if formal != actual {
			return fmt.Errorf("expected type %s; got %s", formal, actual)
		}
	}
	return nil
}

func (t *IrTypechecker) MatchesStructType(formal, actual IrStructType) error {
	if len(formal.Fields) != len(actual.Fields) {
		return fmt.Errorf("expected %d fields; got %d", len(formal.Fields), len(actual.Fields))
	}

	formalFields := append([]StructField{}, formal.Fields...)
	actualFields := append([]StructField{}, actual.Fields...)

	sort.Slice(formalFields, func(i, j int) bool {
		return formalFields[i].Name < formalFields[j].Name
	})
	sort.Slice(actualFields, func(i, j int) bool {
		return actualFields[i].Name < actualFields[j].Name
	})

	for i := range formalFields {
		if formalFields[i].Name != actualFields[i].Name {
			return fmt.Errorf("expected field names %v; got %v", formal.Names(), actual.Names())
		}

		if err := t.MatchesType(formalFields[i].Type, actualFields[i].Type); err != nil {
			return err
		}
	}

	return nil
}

func (t *IrTypechecker) MatchesTupleType(formal, actual IrType) error {
	if formal.Case != TupleType || actual.Case != TupleType {
		panic(fmt.Errorf("expected tuple types"))
	}

	formalTuple := formal.Tuple
	actualTuple := actual.Tuple

	if len(formalTuple) != len(actualTuple) {
		return fmt.Errorf("expected %d elements; got %d", len(formalTuple), len(actualTuple))
	}

	for i := range formalTuple {
		f := formalTuple[i]
		a := actualTuple[i]
		if err := t.MatchesType(f, a); err != nil {
			return err
		}
	}

	return nil
}

func (t *IrTypechecker) MatchesIDType(formal, actual string) error {
	formalDecl, err := t.context.getDecl(formal, FindAny)
	if err != nil {
		return err
	}

	actualDecl, err := t.context.getDecl(actual, FindAny)
	if err != nil {
		return err
	}

	return t.MatchesDecl(formalDecl, actualDecl)
}

func (t *IrTypechecker) MatchesType(formal, actual IrType) error {
	if formal.Case != actual.Case {
		return fmt.Errorf("expected type %s; got %s", formal.Case, actual.Case)
	}

	switch formal.Case {
	case ArrayType:
		return t.MatchesArrayType(*formal.Array, *actual.Array)
	case FunType:
		return t.MatchesFunctionType(formal.Fun, actual.Fun)
	case IntType:
		return t.MatchesIntType(formal.Int, actual.Int)
	case StructType:
		return t.MatchesStructType(formal.Struct, actual.Struct)
	case TupleType:
		return t.MatchesTupleType(formal, actual)
	case IDType:
		return t.MatchesIDType(formal.ID, actual.ID)
	default:
		panic(fmt.Errorf("unhandled IrTypeCase %d", formal.Case))
	}
}

// MatchesDecl determines if the types of the actual declaration are equal to
// the types of the formal declaration. The name of the callee is taken from the
// formal declaration and ignored in the actual declaration.
func (t *IrTypechecker) MatchesDecl(formal, actual IrDecl) error {
	if formal.Case != actual.Case {
		return fmt.Errorf("in declaration %q: expected %s; got %s", formal.ID, formal.Case, actual.Case)
	}

	if err := t.MatchesType(formal.Type, actual.Type); err != nil {
		return fmt.Errorf("in declaration %q: %v", formal.ID, err)
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

		if err := t.CheckTerm(retType, assign.Arg); err != nil {
			return IrType{}, err
		}

		return retType, nil

	case CallTerm:
		call := term.Call

		formal, err := t.context.getType(call.ID, FindAny)
		if err != nil {
			return IrType{}, err
		}

		if formal.Case != FunType {
			return IrType{}, fmt.Errorf("expected function type; got %s", formal)
		}

		if len(formal.Fun.Args) != len(call.Args) {
			return IrType{}, fmt.Errorf("expected %d arguments; got %d", len(formal.Fun.Args), len(call.Args))
		}

		for i := range formal.Fun.Args {
			formalArg := formal.Fun.Args[i]
			actualArg := call.Args[i]
			if err := t.CheckTerm(formalArg, actualArg); err != nil {
				return IrType{}, fmt.Errorf("in argument %d of function %s: %v", i+1, call.ID, err)
			}
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
			field, ok := indexableType.Struct.FieldByIndex(int(*index))
			if !ok {
				return IrType{}, fmt.Errorf("field %d is not a valid field of struct type %s", *index, indexableType)
			}

			term.IndexGet.Field = field.Name
			return field.Type, nil

		case indexableType.Is(StructType) && fieldID != nil:
			field, ok := indexableType.Struct.FieldByID(*fieldID)
			if !ok {
				return IrType{}, fmt.Errorf("field %q is not a valid field of struct type %s", *fieldID, indexableType)
			}

			term.IndexGet.Field = field.Name
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
			if err := t.CheckTerm(NewIntType(I64), term.IndexGet.Index); err != nil {
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
			field, ok := indexableType.Struct.FieldByIndex(int(*index))
			if !ok {
				return IrType{}, fmt.Errorf("field %d is not a valid field of struct type %s", *index, indexableType)
			}

			term.IndexSet.Field = field.Name
			return field.Type, nil

		case indexableType.Is(StructType) && fieldID != nil:
			field, ok := indexableType.Struct.FieldByID(*fieldID)
			if !ok {
				return IrType{}, fmt.Errorf("field %q is not a valid field of struct type %s", *fieldID, indexableType)
			}

			term.IndexSet.Field = field.Name
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
			if err := t.CheckTerm(NewIntType(I64), term.IndexSet.Index); err != nil {
				return IrType{}, err
			}
			if err := t.CheckTerm(indexableType.Array.ElemType, term.IndexSet.Arg); err != nil {
				return IrType{}, err
			}
			return NewTupleType(nil), nil

		default:
			return IrType{}, fmt.Errorf("expected indexable type (e.g., array); got %s", indexableType)
		}

	case StatementTerm:
		if _, err := t.SynthesizeTerm(term.Statement.Expr); err != nil {
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
	}

	panic(fmt.Errorf("unhandled IrTerm %d", term.Case))
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

func (t *IrTypechecker) CheckTerm(formal IrType, term IrTerm) error {
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

		return t.MatchesType(formal, NewTupleType(nil))

	case term.Case == StatementTerm:
		if _, err := t.SynthesizeTerm(term.Statement.Expr); err != nil {
			return err
		}
		return t.MatchesType(formal, NewTupleType(nil))

	case term.Case == TokenTerm && !t.bindPosition:
		switch token := term.Token; token.Case {
		case parser.IDToken:
			actualType, err := t.context.getType(token.Text, FindAny)
			if err != nil {
				return err
			}
			return t.MatchesType(formal, actualType)

		case parser.NumberToken:
			if !formal.Is(IntType) {
				return fmt.Errorf("expected type %s; got %q", formal, token.Text)
			}
			return nil

		default:
			panic(fmt.Errorf("unhandled token %d", token.Case))
		}

	case term.Case == TokenTerm && t.bindPosition:
		switch token := term.Token; token.Case {
		case parser.IDToken:
			actualDecl, err := t.context.getDecl(token.Text, FindAny)
			if err != nil {
				return err
			}
			if actualDecl.Case != VarDecl {
				return fmt.Errorf("expected symbol declared as %s; got %q", VarDecl, actualDecl.Case)
			}
			return t.MatchesType(formal, actualDecl.Type)

		case parser.NumberToken:
			return fmt.Errorf("expected symbol declared as %s; got number literal", VarDecl)

		default:
			panic(fmt.Errorf("unhandled token %d", token.Case))
		}

	case term.Case == OpUnaryTerm:
		if !formal.Is(IntType) {
			return fmt.Errorf("expected integer type; got %s", formal)
		}

		return t.CheckTerm(formal, term.OpUnary.Term)

	case term.Case == OpBinaryTerm:
		if !formal.Is(IntType) {
			return fmt.Errorf("expected integer type; got %s", formal)
		}

		if err := t.CheckTerm(formal, term.OpBinary.Left); err != nil {
			return err
		}

		return t.CheckTerm(formal, term.OpBinary.Right)

	case term.Case == WidenTerm:
		return t.withWiden(func() error {
			return t.CheckTerm(formal, term.Widen.Term)
		})

	default:
		actual, err := t.SynthesizeTerm(term)
		if err != nil {
			return err
		}
		return t.MatchesType(formal, actual)
	}
}

func NewIrTypechecker(context *IrContext) *IrTypechecker {
	return &IrTypechecker{
		context,
		false, /* widen */
		false, /* bindPosition */
	}
}
