package stlc

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
)

type Typechecker struct {
	*log.Logger
	context      Context
	widen        bool
	bindPosition bool
}

func (t *Typechecker) withBindPosition(callback func() (ir.IrType, error)) (ir.IrType, error) {
	bind := t.bindPosition
	t.bindPosition = true
	defer func() { t.bindPosition = bind }()
	return callback()
}

func (t *Typechecker) withWiden(callback func() error) error {
	widen := t.widen
	t.widen = true
	defer func() { t.widen = widen }()
	return callback()
}

func (t *Typechecker) isNumber(typ ir.IrType) error {
	if typ.Case == ir.NameType &&
		(typ.Name == "i8" || typ.Name == "i16" || typ.Name == "i32" || typ.Name == "i64") {
		return nil
	}

	return fmt.Errorf("expected number type, e.g., i8, i16, i32, i64; got %v", typ)
}

func (t *Typechecker) synthesizeApplyImpl(typ ir.IrType, types []ir.IrType, term *ir.IrTerm) (ir.IrType, error) {
	switch typ.Case {
	case ir.ForallType:
		if len(types) != len(typ.Forall.Vars) {
			return ir.IrType{}, fmt.Errorf("expected %d types to call parametric type %s; got %v", len(typ.Forall.Vars), typ, types)
		}

		for _, typ := range types {
			if err := IsWellformedType(t.context, typ); err != nil {
				return ir.IrType{}, err
			}
		}

		for i, tvar := range typ.Forall.Vars {
			tvar = strings.TrimPrefix(tvar, "'")
			typeVar := ir.NewVarType(tvar)
			typ = ir.SubstituteType(typ, typeVar, types[i])
		}

		return t.synthesizeApply(typ.Forall.Type, nil /* types */, term)

	case ir.FunType:
		if err := t.check(term, typ.Fun.Arg); err != nil {
			return ir.IrType{}, err
		}

		return typ.Fun.Ret, nil

	default:
		panic(fmt.Errorf("unhandled %T %d", typ.Case, typ.Case))
	}
}

func (t *Typechecker) synthesizeApply(typ ir.IrType, types []ir.IrType, term *ir.IrTerm) (ir.IrType, error) {
	termType, err := t.synthesizeApplyImpl(typ, types, term)
	if err != nil {
		return ir.IrType{}, err
	}

	term.Type = &termType
	return termType, nil
}

func (t *Typechecker) synthesizeImpl(term *ir.IrTerm) (ir.IrType, error) {
	switch term.Case {
	case ir.AssignTerm:
		retType, err := t.withBindPosition(func() (ir.IrType, error) {
			return t.synthesize(&term.Assign.Ret)
		})
		if err != nil {
			return ir.IrType{}, err
		}

		if err := t.check(&term.Assign.Arg, retType); err != nil {
			return ir.IrType{}, err
		}

		return retType, nil

	case ir.BlockTerm:
		c := term.Block
		for i := range c.Terms {
			if _, err := t.synthesizeFull(&c.Terms[i]); err != nil {
				return ir.IrType{}, err
			}
		}
		return ir.NewTupleType(nil), nil

	case ir.CallTerm:
		c := term.Call

		idTerm := ir.NewTokenTerm(parser.NewIDToken(c.ID))
		formal, err := t.synthesize(&idTerm)
		if err != nil {
			return ir.IrType{}, err
		}

		return t.synthesizeApply(formal, c.Types, &c.Arg)

	case ir.IfTerm:
		c := term.If

		condType, err := t.synthesizeFull(&c.Condition)
		if err != nil {
			return ir.IrType{}, err
		}

		if err := t.isNumber(condType); err != nil {
			return ir.IrType{}, err
		}

		if c.Else == nil {
			return t.synthesizeFull(&c.Then)
		}

		typ, err := t.synthesizeFull(&c.Then)
		if err != nil {
			return ir.IrType{}, err
		}

		if err := t.check(c.Else, typ); err != nil {
			return ir.IrType{}, err
		}

		return typ, nil

	case ir.IndexGetTerm:
		objType, err := t.synthesizeFull(&term.IndexGet.Obj)
		if err != nil {
			return ir.IrType{}, err
		}

		var index *int64
		var fieldID *string
		if term.IndexGet.Index.Case == ir.TokenTerm {
			switch term.IndexGet.Index.Token.Case {
			case parser.NumberToken:
				index = &term.IndexGet.Index.Token.Value
			case parser.IDToken:
				fieldID = &term.IndexGet.Index.Token.Text
			}
		}

		switch {
		case objType.Is(ir.StructType) && index != nil:
			field, ok := objType.FieldByIndex(int(*index))
			if !ok {
				return ir.IrType{}, fmt.Errorf("field %d is not a valid field of struct type %s", *index, objType)
			}

			term.IndexGet.Field = field.ID
			return field.Type, nil

		case objType.Is(ir.StructType) && fieldID != nil:
			field, ok := objType.FieldByID(*fieldID)
			if !ok {
				return ir.IrType{}, fmt.Errorf("field %q is not a valid field of struct type %s", *fieldID, objType)
			}

			term.IndexGet.Field = field.ID
			return field.Type, nil

		case objType.Is(ir.StructType):
			return ir.IrType{}, fmt.Errorf("expected field identifier or number literal to index struct %s", objType)

		case objType.Is(ir.ArrayType) && index != nil:
			if *index < 0 || *index >= int64(objType.Array.Size) {
				return ir.IrType{}, fmt.Errorf("index %d is out of bounds", *index)
			}
			return objType.Array.ElemType, nil

		case objType.Is(ir.ArrayType):
			indexType, err := t.synthesizeFull(&term.IndexGet.Index)
			if err != nil {
				return ir.IrType{}, err
			}

			if err := t.isNumber(indexType); err != nil {
				return ir.IrType{}, err
			}

			return objType.Array.ElemType, nil

		default:
			return ir.IrType{}, fmt.Errorf("expected indexable type (e.g., array, struct, etc); got %s", objType)
		}

	case ir.IndexSetTerm:
		var index *int64
		var fieldID *string
		if term.IndexSet.Index.Case == ir.TokenTerm {
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

		objType, err := t.synthesizeFull(&term.IndexSet.Obj)
		if err != nil {
			return ir.IrType{}, err
		}

		switch {
		case objType.Is(ir.StructType) && index != nil:
			field, ok := objType.FieldByIndex(int(*index))
			if !ok {
				return ir.IrType{}, fmt.Errorf("field %d is not a valid field of struct type %s", *index, objType)
			}

			term.IndexSet.Field = field.ID
			return field.Type, nil

		case objType.Is(ir.StructType) && fieldID != nil:
			field, ok := objType.FieldByID(*fieldID)
			if !ok {
				return ir.IrType{}, fmt.Errorf("field %q is not a valid field of struct type %s", *fieldID, objType)
			}

			term.IndexSet.Field = field.ID
			return field.Type, nil

		case objType.Is(ir.StructType):
			return ir.IrType{}, fmt.Errorf("expected field identifier or number literal to index struct %s", objType)

		case objType.Is(ir.ArrayType) && index != nil:
			if *index < 0 || *index >= int64(objType.Array.Size) {
				return ir.IrType{}, fmt.Errorf("index %d is out of bounds", *index)
			}
			return objType.Array.ElemType, nil

		case objType.Is(ir.ArrayType):
			indexType, err := t.synthesizeFull(&term.IndexSet.Index)
			if err != nil {
				return ir.IrType{}, err
			}

			if err := t.isNumber(indexType); err != nil {
				return ir.IrType{}, err
			}

			if err := t.check(&term.IndexSet.Value, objType.Array.ElemType); err != nil {
				return ir.IrType{}, err
			}

			return ir.NewTupleType(nil), nil

		default:
			return ir.IrType{}, fmt.Errorf("expected indexable type (e.g., array); got %s", objType)
		}

	case ir.LetTerm:
		c := term.Let
		var err error
		if t.context, err = t.context.AddBind(NewDeclBind(DefSymbol, c.Decl)); err != nil {
			return ir.IrType{}, err
		}
		return c.Decl.Type(), nil

	case ir.StatementTerm:
		c := term.Statement
		if _, err := t.synthesizeFull(&c.Term); err != nil {
			return ir.IrType{}, err
		}
		return ir.NewTupleType(nil), nil

	case ir.TokenTerm:
		token := term.Token
		switch token.Case {
		case parser.IDToken:
			bind, err := t.context.getBind(token.Text, FindAny)
			if err != nil {
				return ir.IrType{}, err
			}

			if bind.Decl.Case != ir.TermDecl {
				return ir.IrType{}, fmt.Errorf("expected term; got %s", bind.Decl)
			}

			return bind.Decl.Type(), err

		case parser.NumberToken:
			return ir.IrType{}, fmt.Errorf("cannot synthesize number token types")

		default:
			panic(fmt.Errorf("unhandled token %d", token.Case))
		}

	case ir.TupleTerm:
		types := make([]ir.IrType, len(term.Tuple))
		for i := range term.Tuple {
			var err error
			types[i], err = t.synthesize(&term.Tuple[i])
			if err != nil {
				return ir.IrType{}, err
			}
		}
		return ir.NewTupleType(types), nil

	default:
		panic(fmt.Errorf("unhandled ir.IrTerm %d", term.Case))
	}
}

func (t *Typechecker) synthesize(term *ir.IrTerm) (ir.IrType, error) {
	typ, err := t.synthesizeImpl(term)
	if err != nil {
		return ir.IrType{}, fmt.Errorf("%v\n  synthesizing %s", err, *term)
	}

	term.Type = &typ
	t.Printf("synthesize: %s |- %s", t.context.StringNoImports(), *term)
	return typ, nil
}

// synthesizeFull synthesizes the type for a term and also resolves
// any type alias / type names to the final type.
func (t *Typechecker) synthesizeFull(term *ir.IrTerm) (ir.IrType, error) {
	typ, err := t.synthesize(term)
	if err != nil {
		return ir.IrType{}, err
	}

	return t.context.resolveTypeName(typ)
}

func (t *Typechecker) checkImpl(term *ir.IrTerm, typ ir.IrType) error {
	switch {
	case term.Case == ir.AssignTerm:
		retType, err := t.withBindPosition(func() (ir.IrType, error) {
			return t.synthesize(&term.Assign.Ret)
		})
		if err != nil {
			return err
		}

		return t.check(&term.Assign.Arg, retType)

	case term.Case == ir.StatementTerm:
		c := term.Statement
		if _, err := t.synthesize(&c.Term); err != nil {
			return err
		}
		return t.subtype(ir.NewTupleType(nil), typ)

	case term.Case == ir.TokenTerm && t.bindPosition:
		switch token := term.Token; token.Case {
		case parser.IDToken:
			bind, err := t.context.getBind(token.Text, FindAny)
			if err != nil {
				return err
			}

			if bind.Decl.Case != ir.TermDecl {
				return fmt.Errorf("expected term; got %s", bind.Decl)
			}

			return t.subtype(bind.Decl.Type(), typ)

		case parser.NumberToken:
			return fmt.Errorf("expected symbol declared as %s; got number literal", ir.TermDecl)

		default:
			panic(fmt.Errorf("unhandled token %d", token.Case))
		}

	case term.Case == ir.TokenTerm && !t.bindPosition && term.Token.Case == parser.NumberToken:
		return t.isNumber(typ)

	case term.Case == ir.TupleTerm && typ.Case == ir.TupleType:
		for i := range term.Tuple {
			if err := t.check(&term.Tuple[i], typ.Tuple[i]); err != nil {
				return err
			}
		}
		return nil

	case term.Case == ir.WidenTerm:
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

func (t *Typechecker) check(term *ir.IrTerm, typ ir.IrType) error {
	if err := t.checkImpl(term, typ); err != nil {
		return fmt.Errorf("%s\n  checking %s with %s", err, *term, typ)
	}

	term.Type = &typ
	t.Printf("check: %s |- %s <= %s", t.context.StringNoImports(), *term, typ)
	return nil
}

func (t *Typechecker) TypecheckTerm(term *ir.IrTerm) error {
	return t.check(term, ir.NewTupleType(nil))
}

func (t *Typechecker) TypecheckFunction(function *ir.IrFunction) (Context, error) {
	origContext := t.context

	var err error
	retContext, err := t.context.AddBind(NewDeclBind(DefSymbol, function.Decl()))
	if err != nil {
		return origContext, err
	}

	t.context = retContext
	if t.context, err = t.context.enterFunction(function.ID, function.TypeVars, function.Args, function.Rets); err != nil {
		return origContext, err
	}

	if err := t.TypecheckTerm(&function.Body); err != nil {
		return origContext, err
	}

	return retContext, nil
}

func NewTypechecker(context Context) *Typechecker {
	return &Typechecker{
		log.New(os.Stderr, "DEBUG ", 0),
		context,
		false, /* widen */
		false, /* bindPosition */
	}
}
