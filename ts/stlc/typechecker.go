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
	bindPosition bool
}

func (t *Typechecker) withBindPosition(callback func() error) error {
	bind := t.bindPosition
	t.bindPosition = true
	defer func() { t.bindPosition = bind }()
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
		if len(types) != 0 {
			return ir.IrType{}, fmt.Errorf("expected no types when call non-parametric type %s; got %v", typ, types)
		}

		if err := t.typecheck(term); err != nil {
			return ir.IrType{}, err
		}

		if err := t.subtype(*term.Type, typ.Fun.Arg); err != nil {
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

func (t *Typechecker) typecheckImpl(term *ir.IrTerm) error {
	switch {
	case term.Case == ir.AssignTerm:
		c := term.Assign
		if err := t.withBindPosition(func() error {
			return t.typecheck(&c.Ret)
		}); err != nil {
			return err
		}

		if err := t.typecheck(&c.Arg); err != nil {
			return err
		}

		if err := t.subtype(*c.Ret.Type, *c.Arg.Type); err != nil {
			return err
		}

		term.Type = c.Ret.Type
		return nil

	case term.Case == ir.BlockTerm:
		c := term.Block
		for i := range c.Terms {
			if err := t.typecheck(&c.Terms[i]); err != nil {
				return err
			}
		}

		typ := ir.NewTupleType(nil)
		term.Type = &typ
		return nil

	case term.Case == ir.CallTerm:
		c := term.Call

		idTerm := ir.NewTokenTerm(parser.NewIDToken(c.ID))
		if err := t.typecheck(&idTerm); err != nil {
			return err
		}

		retType, err := t.synthesizeApply(*idTerm.Type, c.Types, &c.Arg)
		if err != nil {
			return err
		}

		term.Type = &retType
		return nil

	case term.Case == ir.IfTerm && len(term.If.Types) == 1:
		c := term.If

		if err := t.typecheck(&c.Condition); err != nil {
			return err
		}

		if err := t.subtype(c.Types[0], *c.Condition.Type); err != nil {
			return err
		}

		if err := t.typecheck(&c.Then); err != nil {
			return err
		}

		if c.Else != nil {
			if err := t.typecheck(c.Else); err != nil {
				return err
			}

			if err := t.subtype(*c.Then.Type, *c.Else.Type); err != nil {
				return err
			}
		}

		term.Type = c.Then.Type
		return nil

	case term.Case == ir.IndexGetTerm:
		c := term.IndexGet
		if err := t.typecheckFull(&c.Obj); err != nil {
			return err
		}

		var index *int64
		var fieldID *string
		if c.Index.Case == ir.TokenTerm {
			switch c.Index.Token.Case {
			case parser.NumberToken:
				index = &c.Index.Token.Value
			case parser.IDToken:
				fieldID = &c.Index.Token.Text
			}
		}

		objType := *c.Obj.Type
		switch {
		case objType.Is(ir.StructType) && index != nil:
			field, ok := objType.FieldByIndex(int(*index))
			if !ok {
				return fmt.Errorf("field %d is not a valid field of struct type %s", *index, objType)
			}

			c.Field = field.ID
			term.Type = &field.Type
			return nil

		case objType.Is(ir.StructType) && fieldID != nil:
			field, ok := objType.FieldByID(*fieldID)
			if !ok {
				return fmt.Errorf("field %q is not a valid field of struct type %s", *fieldID, objType)
			}

			c.Field = field.ID
			term.Type = &field.Type
			return nil

		case objType.Is(ir.StructType):
			return fmt.Errorf("expected field identifier or number literal to index struct %s", objType)

		case objType.Is(ir.ArrayType) && index != nil:
			if *index < 0 || *index >= int64(objType.Array.Size) {
				return fmt.Errorf("index %d is out of bounds", *index)
			}

			term.Type = &objType.Array.ElemType
			return nil

		case objType.Is(ir.ArrayType):
			if err := t.typecheck(&c.Index); err != nil {
				return err
			}

			if err := t.isNumber(*c.Index.Type); err != nil {
				return err
			}

			term.Type = &objType.Array.ElemType
			return nil

		default:
			return fmt.Errorf("expected indexable type (e.g., array, struct, etc); got %s", objType)
		}

	case term.Case == ir.IndexSetTerm:
		c := term.IndexSet

		var index *int64
		var fieldID *string
		if c.Index.Case == ir.TokenTerm {
			switch c.Index.Token.Case {
			// Set field by index.
			//
			// Example:
			//   Index.set x 0 value
			case parser.NumberToken:
				index = &c.Index.Token.Value
			// Set field by label.
			//
			// Example:
			//   Index.set x myfield value
			case parser.IDToken:
				fieldID = &c.Index.Token.Text
			}
		}

		if err := t.typecheckFull(&c.Obj); err != nil {
			return err
		}

		objType := *c.Obj.Type
		switch {
		case objType.Is(ir.StructType) && index != nil:
			field, ok := objType.FieldByIndex(int(*index))
			if !ok {
				return fmt.Errorf("field %d is not a valid field of struct type %s", *index, objType)
			}

			c.Field = field.ID
			term.Type = &field.Type
			return nil

		case objType.Is(ir.StructType) && fieldID != nil:
			field, ok := objType.FieldByID(*fieldID)
			if !ok {
				return fmt.Errorf("field %q is not a valid field of struct type %s", *fieldID, objType)
			}

			c.Field = field.ID
			term.Type = &field.Type
			return nil

		case objType.Is(ir.StructType):
			return fmt.Errorf("expected field identifier or number literal to index struct %s", objType)

		case objType.Is(ir.ArrayType) && index != nil:
			if *index < 0 || *index >= int64(objType.Array.Size) {
				return fmt.Errorf("index %d is out of bounds", *index)
			}

			term.Type = &objType.Array.ElemType
			return nil

		case objType.Is(ir.ArrayType):
			if err := t.typecheck(&c.Index); err != nil {
				return err
			}

			if err := t.isNumber(*c.Index.Type); err != nil {
				return err
			}

			if err := t.typecheck(&c.Value); err != nil {
				return err
			}

			if err := t.subtype(objType.Array.ElemType, *c.Value.Type); err != nil {
				return err
			}

			typ := ir.NewTupleType(nil)
			term.Type = &typ
			return nil

		default:
			return fmt.Errorf("expected indexable type (e.g., array); got %s", objType)
		}

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
		switch {
		case c.Case == parser.IDToken:
			bind, err := t.context.getBind(c.Text, FindAny)
			if err != nil {
				return err
			}

			if bind.Decl.Case != ir.TermDecl {
				return fmt.Errorf("expected term; got %s", bind.Decl)
			}

			typ := bind.Decl.Type()
			term.Type = &typ
			return nil

		case c.Case == parser.NumberToken && t.bindPosition:
			return fmt.Errorf("expected symbol declared as %s; got number literal", ir.TermDecl)

		case c.Case == parser.NumberToken && !t.bindPosition && term.Type != nil && t.isNumber(*term.Type) == nil:
			return nil

		case c.Case == parser.NumberToken && !t.bindPosition:
			return fmt.Errorf("cannot synthesize a type for a number")

		default:
			panic(fmt.Errorf("unhandled %T %d", c.Case, c.Case))
		}

	case term.Case == ir.TupleTerm:
		types := make([]ir.IrType, len(term.Tuple))
		for i := range term.Tuple {
			var err error
			if err = t.typecheck(&term.Tuple[i]); err != nil {
				return err
			}
			types[i] = *term.Tuple[i].Type
		}

		typ := ir.NewTupleType(types)
		term.Type = &typ
		return nil

	default:
		panic(fmt.Errorf("unhandled ir.IrTerm %d", term.Case))
	}
}

func (t *Typechecker) typecheck(term *ir.IrTerm) error {
	if err := t.typecheckImpl(term); err != nil {
		return fmt.Errorf("%v\n  typechecking %s", err, *term)
	}

	t.Printf("typecheck: %s |- %s", t.context.StringNoImports(), *term)
	return nil
}

func (t *Typechecker) typecheckFull(term *ir.IrTerm) error {
	if err := t.typecheck(term); err != nil {
		return err
	}

	typ, err := t.context.resolveTypeName(*term.Type)
	if err != nil {
		return err
	}

	term.Type = &typ
	return nil
}

func (t *Typechecker) TypecheckTerm(term *ir.IrTerm) error {
	return t.typecheck(term)
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
		false, /* bindPosition */
	}
}
