package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

func isWellformedTermBind(context Context, bind *termBind) error {
	if bind.Symbol == DeclSymbol {
		if _, ok := context.lookupTermBind(bind.Name); ok {
			return fmt.Errorf("context is not wellformed: term %q is already defined", bind.Name)
		}
	} else if bind.Symbol == DefSymbol {
		if _, ok := context.lookupTermBindWithSymbol(bind.Name, DefSymbol); ok {
			return fmt.Errorf("context is not wellformed: term %q is already defined", bind.Name)
		}

		if declBind, ok := context.lookupTermBindWithSymbol(bind.Name, DeclSymbol); ok {
			if err := NewTypechecker(context).subtype(bind.Type, declBind.Term.Type); err != nil {
				return fmt.Errorf("context is not wellformed: type %s does not match declaration type %s", bind.Type, declBind.Term.Type)
			}
		}
	}

	if err := isWellformedType(context, bind.Type); err != nil {
		return fmt.Errorf("context is not wellformed: type %s is not wellformed: %v", bind.Type, err)
	}

	kind, err := inferKind(context, bind.Type)
	if err != nil {
		return err
	}
	if !ir.EqualsKind(kind, ir.NewTypeKind()) {
		return fmt.Errorf("context is not wellformed: term %s with type %s must have kind %s instead of kind %s", bind.Name, bind.Type, ir.NewTypeKind(), kind)
	}

	return nil
}

func isWellformedAliasBind(context Context, bind *aliasBind) error {
	if _, ok := context.LookupBind(bind.Name, FindDefOnly); ok {
		return fmt.Errorf("context is not wellformed: term %q is defined more than once", bind.Name)
	}
	if err := isWellformedType(context, bind.Type); err != nil {
		return fmt.Errorf("context is not wellformed: aliased type %s is not wellformed: %v", bind.Type, err)
	}
	return nil
}

func isWellformedComponentBind(context Context, bind *componentBind) error {
	if ok := context.containsComponentBind(bind.ElemType); ok {
		return fmt.Errorf("context is not wellformed: component %q is defined more than once", bind.ElemType)
	}
	if err := isWellformedType(context, bind.ElemType); err != nil {
		return fmt.Errorf("context is not wellformed: component type %s is not wellformed: %v", bind.ElemType, err)
	}
	return nil
}

func isWellformedConstBind(context Context, bind *constBind) error {
	if _, ok := context.lookupConstBind(bind.Name); ok {
		return fmt.Errorf("context is not wellformed: type %q is defined more than once", bind.Name)
	}
	return nil
}

func isWellformedTypeVarBind(context Context, bind *typeVarBind) error {
	if context.containsTypeVarBind(bind.Name) {
		return fmt.Errorf("context is not wellformed: type variable %q is defined more than once\ncontext: %s", bind.Name, context.String())
	}
	return nil
}
