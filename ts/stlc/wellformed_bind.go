package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

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

func isWellformedTermDeclBind(context Context, bind *termDeclBind) error {
	if ok := context.containsTermDeclBindInScope(bind.Name); ok {
		return fmt.Errorf("term %q is already declared", bind.Name)
	}

	if ok := context.containsTermDefBindInScope(bind.Name); ok {
		return fmt.Errorf("term %q is already defined", bind.Name)
	}

	if err := isWellformedType(context, bind.Type); err != nil {
		return fmt.Errorf("term %s has type %s that is not wellformed: %v", bind.Name, bind.Type, err)
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

func isWellformedTermDefBind(context Context, bind *termDefBind) error {
	if ok := context.containsTermDefBindInScope(bind.Name); ok {
		return fmt.Errorf("term %q is already defined", bind.Name)
	}

	if declBind, ok := context.lookupTermDeclBindInScope(bind.Name); ok {
		if err := NewTypechecker(context).subtype(bind.Type, declBind.TermDecl.Type); err != nil {
			return fmt.Errorf("term %s has type %s that does not match the declaration type %s\n  because %v", bind.Name, bind.Type, declBind.TermDecl.Type, err)
		}
	}

	if err := isWellformedType(context, bind.Type); err != nil {
		return fmt.Errorf("term %s has type %s that is not wellformed: %v", bind.Name, bind.Type, err)
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
	if context.containsTypeVarBindInScope(bind.Name) {
		return fmt.Errorf("type variable %q is already defined", bind.Name)
	}
	return nil
}
