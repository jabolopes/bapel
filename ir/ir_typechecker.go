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

func (t *IrTypechecker) isNumber(typ IrType) error {
	if typ.Case == NameType &&
		(typ.Name == "i8" || typ.Name == "i16" || typ.Name == "i32" || typ.Name == "i64") {
		return nil
	}

	return fmt.Errorf("expected number type, e.g., i8, i16, i32, i64")
}

func (t *IrTypechecker) subtypeImpl(left, right IrType) error {
	switch {
	case left.Case == ArrayType && right.Case == ArrayType:
		if err := t.subtype(left.Array.ElemType, right.Array.ElemType); err != nil {
			return fmt.Errorf("mismatch in array element types: %v", err)
		}

		if left.Array.Size != right.Array.Size {
			return fmt.Errorf("expected array with %d elements; got %d elements", left.Array.Size, right.Array.Size)
		}

		return nil

	case left.Case == ForallType && right.Case == ForallType:
		if len(left.Forall.Vars) != len(right.Forall.Vars) {
			return fmt.Errorf("expected forall type with %d variables (%v); got %d variables (%v)",
				len(left.Forall.Vars), left.Forall.Vars,
				len(right.Forall.Vars), right.Forall.Vars)
		}

		leftType := left.Forall.Type
		for i := range right.Forall.Vars {
			leftType = substituteType(leftType, NewVarType(right.Forall.Vars[i]), NewVarType(right.Forall.Vars[i]))
		}

		return t.subtype(leftType, right.Forall.Type)

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

	// Typenames.
	case left.Case == NameType && right.Case == NameType && left.Name == right.Name:
		return nil

	default:
		return fmt.Errorf("expected type %s (%s); got %s (%s)", left.Case, left, right.Case, right)
	}
}

func (t *IrTypechecker) subtype(left, right IrType) error {
	if err := t.subtypeImpl(left, right); err != nil {
		return fmt.Errorf("%s\n  subtyping %s and %s", err, left, right)
	}

	t.Printf("subtype: %s |- %s < %s", t.context.StringNoImports(), left, right)
	return nil
}

func (t *IrTypechecker) synthesizeApplyImpl(typ IrType, types []IrType, term *IrTerm) (IrType, error) {
	switch typ.Case {
	case ForallType:
		if len(types) != len(typ.Forall.Vars) {
			return IrType{}, fmt.Errorf("expected %d types to call parametric type %s; got %v", len(typ.Forall.Vars), typ, types)
		}

		for _, typ := range types {
			if err := isTypeWellFormed(*t.context, typ); err != nil {
				return IrType{}, err
			}
		}

		for i, tvar := range typ.Forall.Vars {
			tvar = strings.TrimPrefix(tvar, "'")
			typeVar := NewVarType(tvar)
			typ = substituteType(typ, typeVar, types[i])
		}

		return t.synthesizeApply(typ.Forall.Type, nil /* types */, term)

	case FunType:
		if err := t.check(term, typ.Fun.Arg); err != nil {
			return IrType{}, err
		}

		return typ.Fun.Ret, nil

	default:
		panic(fmt.Errorf("unhandled IrType case %d", typ.Case))
	}
}

func (t *IrTypechecker) synthesizeApply(typ IrType, types []IrType, term *IrTerm) (IrType, error) {
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
		c := term.If

		condType, err := t.synthesizeFull(&c.Condition)
		if err != nil {
			return IrType{}, err
		}

		if err := t.isNumber(condType); err != nil {
			return IrType{}, err
		}

		if c.Else == nil {
			return t.synthesizeFull(&c.Then)
		}

		typ, err := t.synthesizeFull(&c.Then)
		if err != nil {
			return IrType{}, err
		}

		if err := t.check(c.Else, typ); err != nil {
			return IrType{}, err
		}

		return typ, nil

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
			indexType, err := t.synthesizeFull(&term.IndexGet.Index)
			if err != nil {
				return IrType{}, err
			}

			if err := t.isNumber(indexType); err != nil {
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
			// Set field by index.
			//
			// Example:
			//   Index.set x 0 value
			case parser.NumberToken:
				index = &term.IndexSet.Index.Token.Value
			// Set field by label.
			//
			// Example:
			//   Index.set x myfield value
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
			indexType, err := t.synthesizeFull(&term.IndexSet.Index)
			if err != nil {
				return IrType{}, err
			}

			if err := t.isNumber(indexType); err != nil {
				return IrType{}, err
			}

			if err := t.check(&term.IndexSet.Arg, indexableType.Array.ElemType); err != nil {
				return IrType{}, err
			}

			return NewTupleType(nil), nil

		default:
			return IrType{}, fmt.Errorf("expected indexable type (e.g., array); got %s", indexableType)
		}

	case LetTerm:
		c := term.Let
		if err := t.context.AddBind(NewDeclBind(DefSymbol, c.Decl)); err != nil {
			return IrType{}, err
		}
		return c.Decl.Type(), nil

	case StatementTerm:
		c := term.Statement
		for i := range c.Terms {
			if _, err := t.synthesize(&c.Terms[i]); err != nil {
				return IrType{}, err
			}
		}
		return NewTupleType(nil), nil

	case TokenTerm:
		token := term.Token
		switch token.Case {
		case parser.IDToken:
			bind, err := t.context.getBind(token.Text, FindAny)
			if err != nil {
				return IrType{}, err
			}

			if bind.Decl.Case != TermDecl {
				return IrType{}, fmt.Errorf("expected term; got %s", bind.Decl)
			}

			return bind.Decl.Type(), err

		case parser.NumberToken:
			return IrType{}, fmt.Errorf("cannot synthesize number token types")

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
		return IrType{}, fmt.Errorf("%v\n  synthesizing %s", err, *term)
	}

	term.Type = &typ
	t.Printf("synthesize: %s |- %s", t.context.StringNoImports(), *term)
	return typ, nil
}

// synthesizeFull synthesizes the type for a term and also resolves
// any type alias / type names to the final type.
func (t *IrTypechecker) synthesizeFull(term *IrTerm) (IrType, error) {
	typ, err := t.synthesize(term)
	if err != nil {
		return IrType{}, err
	}

	return t.context.resolveTypeName(typ)
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

	case term.Case == StatementTerm:
		c := term.Statement
		for i := range c.Terms {
			if _, err := t.synthesize(&c.Terms[i]); err != nil {
				return err
			}
			if err := t.subtype(NewTupleType(nil), typ); err != nil {
				return err
			}
		}
		return nil

	case term.Case == TokenTerm && t.bindPosition:
		switch token := term.Token; token.Case {
		case parser.IDToken:
			bind, err := t.context.getBind(token.Text, FindAny)
			if err != nil {
				return err
			}

			if bind.Decl.Case != TermDecl {
				return fmt.Errorf("expected term; got %s", bind.Decl)
			}

			return t.subtype(bind.Decl.Type(), typ)

		case parser.NumberToken:
			return fmt.Errorf("expected symbol declared as %s; got number literal", TermDecl)

		default:
			panic(fmt.Errorf("unhandled token %d", token.Case))
		}

	case term.Case == TokenTerm && !t.bindPosition && term.Token.Case == parser.NumberToken:
		return t.isNumber(typ)

	case term.Case == TupleTerm && typ.Case == TupleType:
		for i := range term.Tuple {
			if err := t.check(&term.Tuple[i], typ.Tuple[i]); err != nil {
				return err
			}
		}
		return nil

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
	t.Printf("check: %s |- %s <= %s", t.context.StringNoImports(), *term, typ)
	return nil
}

func (t *IrTypechecker) TypecheckTerm(term *IrTerm) error {
	return t.check(term, NewTupleType(nil))
}

func NewIrTypechecker(context *IrContext) *IrTypechecker {
	return &IrTypechecker{
		log.New(os.Stderr, "DEBUG ", 0),
		context,
		false, /* widen */
		false, /* bindPosition */
	}
}
