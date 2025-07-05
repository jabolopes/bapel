package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

func isWellformedTermBind(context Context, bind *termBind) error {
	if bind.Symbol == DeclSymbol {
		if ok := context.containsTermBindInScope(bind.Name); ok {
			return fmt.Errorf("term %q is already defined", bind.Name)
		}
	} else if bind.Symbol == DefSymbol {
		if ok := context.containsTermBindInScopeWithSymbol(bind.Name, DefSymbol); ok {
			return fmt.Errorf("term %q is already defined", bind.Name)
		}

		if declBind, ok := context.lookupTermBindInScopeWithSymbol(bind.Name, DeclSymbol); ok {
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

func isWellformedConstBind(context Context, bind *constBind) error {
	if _, ok := context.lookupConstBind(bind.Name); ok {
		return fmt.Errorf("type %q is already defined", bind.Name)
	}
	return nil
}

func isWellformedScopeBind(context Context, bind *scopeBind) error {
	scopeBind, ok := context.lookupScopeBind()

	wantLevel := 1
	if ok {
		wantLevel = scopeBind.Scope.Level + 1
	}
	if bind.Level != wantLevel {
		return fmt.Errorf("expected scope %d; got %s", wantLevel, bind)
	}

	return nil
}

func isWellformedTypeVarBind(context Context, bind *typeVarBind) error {
	if context.containsTypeVarBind(bind.Name) {
		return fmt.Errorf("type variable %q is already defined", bind.Name)
	}
	return nil
}
