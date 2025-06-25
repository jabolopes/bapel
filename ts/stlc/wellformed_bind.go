package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

func isWellformedTermBind(context Context, bind *termBind) error {
	if bind.Symbol == DeclSymbol {
		if _, ok := context.lookupTermBind(bind.Name); ok {
			return fmt.Errorf("term %q is already defined", bind.Name)
		}
	} else if bind.Symbol == DefSymbol {
		if _, ok := context.lookupTermBindWithSymbol(bind.Name, DefSymbol); ok {
			return fmt.Errorf("term %q is already defined", bind.Name)
		}

		if declBind, ok := context.lookupTermBindWithSymbol(bind.Name, DeclSymbol); ok {
			if err := NewTypechecker(context).subtype(bind.Type, declBind.Term.Type); err != nil {
				return fmt.Errorf("type %s does not match declaration type %s\n  because %v", bind.Type, declBind.Term.Type, err)
			}
		}
	}

	if err := isWellformedType(context, bind.Type); err != nil {
		return fmt.Errorf("type %s is not wellformed: %v", bind.Type, err)
	}

	kind, err := inferKind(context, bind.Type)
	if err != nil {
		return err
	}
	if !ir.EqualsKind(kind, ir.NewTypeKind()) {
		return fmt.Errorf("term %s with type %s must have kind %s instead of kind %s", bind.Name, bind.Type, ir.NewTypeKind(), kind)
	}

	return nil
}

func isWellformedAliasBind(context Context, bind *aliasBind) error {
	if _, ok := context.lookupAliasBind(bind.Name); ok {
		return fmt.Errorf("type %q is already defined", bind.Name)
	}

	if constBind, ok := context.lookupConstBind(bind.Name); ok {
		kind, err := inferKind(context, bind.Type)
		if err != nil {
			return err
		}
		if !ir.EqualsKind(kind, constBind.Const.Kind) {
			return fmt.Errorf("type %s is defined with kind %s that does not match the declaration kind %s", bind.Type, kind, constBind.Const.Kind)
		}
	}

	if err := isWellformedType(context, bind.Type); err != nil {
		return fmt.Errorf("aliased type %s is not wellformed: %v", bind.Type, err)
	}
	return nil
}

func isWellformedComponentBind(context Context, bind *componentBind) error {
	if ok := context.containsComponentBind(bind.ElemType); ok {
		return fmt.Errorf("component %q is already defined", bind.ElemType)
	}
	if err := isWellformedType(context, bind.ElemType); err != nil {
		return fmt.Errorf("component type %s is not wellformed: %v", bind.ElemType, err)
	}
	return nil
}

func isWellformedConstBind(context Context, bind *constBind) error {
	if _, ok := context.lookupConstBind(bind.Name); ok {
		return fmt.Errorf("type %q is already defined", bind.Name)
	}
	return nil
}

func isWellformedTypeVarBind(context Context, bind *typeVarBind) error {
	if context.containsTypeVarBind(bind.Name) {
		return fmt.Errorf("type variable %q is already defined", bind.Name)
	}
	return nil
}
