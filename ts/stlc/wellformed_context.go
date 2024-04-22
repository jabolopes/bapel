package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

func IsWellformedContext(context Context) error {
	if context.Empty() {
		// Rule: EmptyCtx.
		return nil
	}

	bind, newContext := context.Pop()

	if err := IsWellformedContext(newContext); err != nil {
		return err
	}

	switch {
	case bind.Case == DeclBind && bind.Decl.Case == ir.TermDecl:
		if _, ok := newContext.LookupBind(bind.Decl.Term.ID, FindDefOnly); ok {
			return fmt.Errorf("context is not wellformed: term %q is defined more than once", bind.Decl.Term.ID)
		}
		if err := IsWellformedType(newContext, bind.Decl.Type()); err != nil {
			return fmt.Errorf("context is not wellformed: type %s is not wellformed: %v", bind.Decl.Type(), err)
		}
		return nil

	case bind.Case == DeclBind && bind.Decl.Case == ir.TypeDecl:
		c := bind.Decl.AsType
		switch typ := c.Type; typ.Case {
		case ir.AliasType:
			{
				newContext2, err := newContext.AddBind(NewDeclBind(bind.Symbol, ir.NewTypeDecl(typ.Alias.Name)))
				if err != nil {
					return err
				}

				if err := IsWellformedContext(newContext2); err != nil {
					return err
				}
			}
			if err := IsWellformedType(newContext, typ.Alias.Value); err != nil {
				return fmt.Errorf("context is not wellformed: aliased type is not wellformed: %v", err)
			}
			return nil

		case ir.ComponentType:
			if _, ok := newContext.LookupBind(typ.Component.Name, FindDefOnly); ok {
				return fmt.Errorf("context is not wellformed: component %q is defined more than once", typ.Component.Name)
			}
			if err := IsWellformedType(newContext, typ.Component.ElemType); err != nil {
				return fmt.Errorf("context is not wellformed: component type %s is not wellformed: %v", typ.Component.ElemType, err)
			}
			return nil

		case ir.NameType:
			if _, ok := newContext.LookupBind(typ.Name, FindDefOnly); ok {
				return fmt.Errorf("context is not wellformed: name %q is defined more than once", typ.Name)
			}

			// TODO: Check that the definition matches the declaration.
			return nil

		case ir.VarType:
			if newContext.ContainsVarType(typ.Var) {
				return fmt.Errorf("context is not wellformed: type variable %q is defined more than once\ncontext: %s", typ.Var, context.String())
			}
			return nil

		case ir.ArrayType:
		case ir.ForallType:
		case ir.FunType:
		case ir.StructType:
		case ir.TupleType:
			return fmt.Errorf("context should not contain %s type", typ.Case)

		default:
			panic(fmt.Errorf("unhandled %T %d", typ.Case, typ.Case))
		}

	default:
		panic(fmt.Errorf("unhandled %T %d", bind.Case, bind.Case))
	}

	return nil
}
