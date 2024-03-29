package ir

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jabolopes/bapel/parser"
)

type IrTypechecker struct {
	*log.Logger
	context      *IrContext
	idgen        int
	widen        bool
	bindPosition bool
}

func (t *IrTypechecker) genID() string {
	id := t.idgen
	t.idgen++
	return fmt.Sprintf("%d", id)
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

func (t *IrTypechecker) instantiate(left, right IrType) error {
	switch {
	// InstLSolve
	case left.Case == VarExistType &&
		!t.context.isSolvedVar(left.VarExist.Var) &&
		isTypeWellformed(sliceAtType(*t.context, left), right) &&
		IsMonotype(*t.context, right):
		return t.context.setType(left.VarExist.Var, right)

	// InstLReach
	case left.Case == VarExistType &&
		right.Case == VarExistType &&
		!t.context.isSolvedVar(left.VarExist.Var) &&
		!t.context.isSolvedVar(right.VarExist.Var) &&
		t.context.isDefinedInOrder(left.VarExist.Var, right.VarExist.Var):
		return t.context.setType(right.VarExist.Var, left)

	// InstRSolve
	case right.Case == VarExistType &&
		!t.context.isSolvedVar(right.VarExist.Var) &&
		isTypeWellformed(sliceAtType(*t.context, right), left) &&
		IsMonotype(*t.context, left):
		return t.context.setType(right.VarExist.Var, left)

	// InstRReach
	case left.Case == VarExistType &&
		right.Case == VarExistType &&
		!t.context.isSolvedVar(left.VarExist.Var) &&
		!t.context.isSolvedVar(right.VarExist.Var) &&
		t.context.isDefinedInOrder(right.VarExist.Var, left.VarExist.Var):
		return t.context.setType(left.VarExist.Var, right)

	default:
		panic(fmt.Errorf("unhandled cases %s and %s in instantiate", left, right))
	}
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

	// <:->
	case left.Case == FunType && right.Case == FunType:
		// B1 <: A1
		if err := t.subtype(right.Fun.Arg, left.Fun.Arg); err != nil {
			return err
		}

		// A2 <: B2
		if err := t.subtype(left.Fun.Ret, right.Fun.Ret); err != nil {
			return err
		}

		return nil

	// <:Unit
	case left.Case == NumberType && right.Case == NumberType:
		return nil

	// TODO: Improve.
	case left.Case == NumberType &&
		right.Case == NameType &&
		(right.Name == "i8" || right.Name == "i16" || right.Name == "i32" || right.Name == "i64"):
		return nil

	// TODO: Improve.
	case left.Case == NameType && right.Case == NumberType &&
		(left.Name == "i8" || left.Name == "i16" || left.Name == "i32" || left.Name == "i64"):
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

	// <:Var
	case left.Case == VarType && right.Case == VarType && left.Var == right.Var:
		return nil

	// <:Exvar
	case left.Case == VarExistType &&
		right.Case == VarExistType &&
		*left.VarExist == *right.VarExist:
		return nil

	// <:InstantiateL
	case left.Case == VarExistType && !t.context.isSolvedVar(left.VarExist.Var):
		return t.instantiate(left, right)

	case left.Case == VarExistType:
		leftType, err := t.context.getType(left.VarExist.Var, FindAny)
		if err != nil {
			return err
		}
		return t.subtype(leftType, right)

	// <:InstantiateR
	case right.Case == VarExistType && !t.context.isSolvedVar(right.VarExist.Var):
		return t.instantiate(left, right)

	case right.Case == VarExistType:
		rightType, err := t.context.getType(right.VarExist.Var, FindAny)
		if err != nil {
			return err
		}

		return t.subtype(left, rightType)

	// Typenames.
	case left.Case == NameType && right.Case == NameType && left.Name == right.Name:
		return nil

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

func (t *IrTypechecker) synthesizeApplyImpl(typ IrType, types []string, term *IrTerm) (IrType, error) {
	switch typ.Case {
	case ForallType:
		if true {
			if len(types) != len(typ.Forall.Vars) {
				return IrType{}, fmt.Errorf("expected %d types to call parametric type %s; got %v", len(typ.Forall.Vars), typ, types)
			}

			for i, tvar := range typ.Forall.Vars {
				tvar = strings.TrimPrefix(tvar, "'")
				typeVar := NewVarType(tvar)
				typeInst := NewNameType(types[i])
				typ = substituteType(typ, typeVar, typeInst)
			}

			return t.synthesizeApply(typ.Forall.Type, nil /* types */, term)
		}

		for _, tvar := range typ.Forall.Vars {
			typeVar := NewVarType(tvar)
			existVar := NewVarExistType(t.genID())
			typ = substituteType(typ, typeVar, existVar)

			if err := t.context.addBind(NewTypeBind(DefSymbol, existVar, nil)); err != nil {
				return IrType{}, err
			}
		}

		typ, err := t.synthesizeApply(typ.Forall.Type, nil /* types */, term)
		if err != nil {
			return IrType{}, err
		}

		if typ.Is(VarExistType) {
			if resolvedType, err := t.context.getType(typ.VarExist.Var, FindAny); err == nil {
				typ = resolvedType
			}
		}

		return typ, nil

	case FunType:
		if err := t.check(term, typ.Fun.Arg); err != nil {
			return IrType{}, err
		}

		return typ.Fun.Ret, nil

	default:
		panic(fmt.Errorf("unhandled IrType case %d", typ.Case))
	}
}

func (t *IrTypechecker) synthesizeApply(typ IrType, types []string, term *IrTerm) (IrType, error) {
	termType, err := t.synthesizeApplyImpl(typ, types, term)
	if err != nil {
		return IrType{}, err
	}

	term.Type = &termType
	return termType, nil
}

func (t *IrTypechecker) synthesizeImpl(term *IrTerm) (IrType, error) {
	switch term.Case {
	case AssignTerm:
		retType, err := t.withBindPosition(func() (IrType, error) {
			return t.synthesize(&term.Assign.Ret)
		})
		if err != nil {
			return IrType{}, err
		}

		if err := t.check(&term.Assign.Arg, retType); err != nil {
			return IrType{}, err
		}

		return retType, nil

	case CallTerm:
		idTerm := NewTokenTerm(parser.NewIDToken(term.Call.ID))
		formal, err := t.synthesize(&idTerm)
		if err != nil {
			return IrType{}, err
		}

		return t.synthesizeApply(formal, term.Call.Types, &term.Call.Arg)

	case IfTerm:
		condType, err := t.synthesizeFull(&term.If.Condition)
		if err != nil {
			return IrType{}, err
		}

		if err := t.subtype(NewNumberType(), condType); err != nil {
			return IrType{}, err
		}

		return NewTupleType(nil), nil

	case IndexGetTerm:
		indexableType, err := t.synthesizeFull(&term.IndexGet.Term)
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
			if err := t.check(&term.IndexGet.Index, NewNumberType()); err != nil {
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

		indexableType, err := t.synthesizeFull(&term.IndexSet.Ret)
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
			if err := t.check(&term.IndexSet.Index, NewNumberType()); err != nil {
				return IrType{}, err
			}
			if err := t.check(&term.IndexSet.Arg, indexableType.Array.ElemType); err != nil {
				return IrType{}, err
			}
			return NewTupleType(nil), nil

		default:
			return IrType{}, fmt.Errorf("expected indexable type (e.g., array); got %s", indexableType)
		}

	case StatementTerm:
		if _, err := t.synthesize(&term.Statement.Term); err != nil {
			return IrType{}, err
		}
		return NewTupleType(nil), nil

	case TokenTerm:
		token := term.Token
		switch token.Case {
		case parser.IDToken:
			return t.context.getType(token.Text, FindAny)

		case parser.NumberToken:
			return NewNumberType(), nil

		default:
			panic(fmt.Errorf("unhandled token %d", token.Case))
		}

	case TupleTerm:
		types := make([]IrType, len(term.Tuple))
		for i := range term.Tuple {
			var err error
			types[i], err = t.synthesize(&term.Tuple[i])
			if err != nil {
				return IrType{}, err
			}
		}
		return NewTupleType(types), nil

	default:
		panic(fmt.Errorf("unhandled IrTerm %d", term.Case))
	}
}

func (t *IrTypechecker) synthesize(term *IrTerm) (IrType, error) {
	typ, err := t.synthesizeImpl(term)
	if err != nil {
		return IrType{}, err
	}

	term.Type = &typ
	t.Printf("synthesize: %s |- %s", t.context.StringNoImports(), *term)
	return typ, nil
}

func (t *IrTypechecker) synthesizeFull(term *IrTerm) (IrType, error) {
	typ, err := t.synthesize(term)
	if err != nil {
		return IrType{}, fmt.Errorf("%v\n  synthesizing %s", err, *term)
	}

	switch typ.Case {
	case NameType:
		return t.context.getType(typ.Name, FindAny)
	default:
		return typ, nil
	}
}

func (t *IrTypechecker) checkImpl(term *IrTerm, typ IrType) error {
	switch {
	case term.Case == AssignTerm:
		retType, err := t.withBindPosition(func() (IrType, error) {
			return t.synthesize(&term.Assign.Ret)
		})
		if err != nil {
			return err
		}

		return t.check(&term.Assign.Arg, retType)

	case term.Case == IfTerm:
		if err := t.check(&term.If.Condition, NewNumberType()); err != nil {
			return err
		}
		return t.subtype(NewTupleType(nil), typ)

	case term.Case == StatementTerm:
		if _, err := t.synthesize(&term.Statement.Term); err != nil {
			return err
		}
		return t.subtype(NewTupleType(nil), typ)

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

	case term.Case == WidenTerm:
		return t.withWiden(func() error {
			return t.check(&term.Widen.Term, typ)
		})

	default:
		// Sub:
		//   e <= B

		// e => A
		got, err := t.synthesize(term)
		if err != nil {
			return err
		}

		// A <: B
		return t.subtype(got, typ)
	}
}

func (t *IrTypechecker) check(term *IrTerm, typ IrType) error {
	if err := t.checkImpl(term, typ); err != nil {
		return fmt.Errorf("%s\n  checking %s with %s", err, *term, typ)
	}

	term.Type = &typ
	return nil
}

func (t *IrTypechecker) TypecheckTerm(term *IrTerm) error {
	return t.check(term, NewTupleType(nil))
}

func NewIrTypechecker(context *IrContext) *IrTypechecker {
	return &IrTypechecker{
		log.New(os.Stderr, "DEBUG ", 0),
		context,
		0,     /* idgen */
		false, /* widen */
		false, /* bindPosition */
	}
}
