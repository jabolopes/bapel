package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

func isWellformedTermBind(newContext Context, bind *termBind) error {
	if bind.Symbol == DeclSymbol {
		if _, ok := newContext.lookupTermBind(bind.Name); ok {
			return fmt.Errorf("context is not wellformed: term %q is already defined", bind.Name)
		}
	} else if bind.Symbol == DefSymbol {
		if _, ok := newContext.lookupTermBindWithSymbol(bind.Name, DefSymbol); ok {
			return fmt.Errorf("context is not wellformed: term %q is already defined", bind.Name)
		}

		if declBind, ok := newContext.lookupTermBindWithSymbol(bind.Name, DeclSymbol); ok {
			if err := NewTypechecker(newContext).subtype(bind.Type, declBind.Term.Type); err != nil {
				return fmt.Errorf("context is not wellformed: type %s does not match declaration type %s", bind.Type, declBind.Term.Type)
			}
		}
	}

	if err := isWellformedType(newContext, bind.Type); err != nil {
		return fmt.Errorf("context is not wellformed: type %s is not wellformed: %v", bind.Type, err)
	}

	kind, err := InferKind(newContext, bind.Type)
	if err != nil {
		return err
	}
	if !ir.EqualsKind(kind, ir.NewTypeKind()) {
		return fmt.Errorf("context is not wellformed: term %s with type %s must have kind %s instead of kind %s", bind.Name, bind.Type, ir.NewTypeKind(), kind)
	}

	return nil
}

func isWellformedAliasBind(newContext Context, bind *aliasBind) error {
	if _, ok := newContext.LookupBind(bind.Name, FindDefOnly); ok {
		return fmt.Errorf("context is not wellformed: term %q is defined more than once", bind.Name)
	}
	if err := isWellformedType(newContext, bind.Type); err != nil {
		return fmt.Errorf("context is not wellformed: aliased type %s is not wellformed: %v", bind.Type, err)
	}
	return nil
}

func isWellformedComponentBind(newContext Context, bind *componentBind) error {
	if ok := newContext.containsComponentBind(bind.ElemType); ok {
		return fmt.Errorf("context is not wellformed: component %q is defined more than once", bind.ElemType)
	}
	if err := isWellformedType(newContext, bind.ElemType); err != nil {
		return fmt.Errorf("context is not wellformed: component type %s is not wellformed: %v", bind.ElemType, err)
	}
	return nil
}

func isWellformedConstBind(newContext Context, bind *constBind) error {
	if _, ok := newContext.lookupConstBind(bind.Name); ok {
		return fmt.Errorf("context is not wellformed: type %q is defined more than once", bind.Name)
	}
	return nil
}

func isWellformedTypeVarBind(newContext Context, bind *typeVarBind) error {
	if newContext.containsTypeVarBind(bind.Name) {
		return fmt.Errorf("context is not wellformed: type variable %q is defined more than once\ncontext: %s", bind.Name, newContext.String())
	}
	return nil
}
