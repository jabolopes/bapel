package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

// TODO: Check that the definition matches the declaration.
func IsWellformedTermBind(newContext Context, bind *termBind) error {
	if _, ok := newContext.LookupBind(bind.Name, FindDefOnly); ok {
		return fmt.Errorf("context is not wellformed: term %q is defined more than once", bind.Name)
	}

	if err := IsWellformedType(newContext, bind.Type); err != nil {
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

func IsWellformedAliasBind(newContext Context, bind *aliasBind) error {
	if _, ok := newContext.LookupBind(bind.Name, FindDefOnly); ok {
		return fmt.Errorf("context is not wellformed: term %q is defined more than once", bind.Name)
	}
	if err := IsWellformedType(newContext, bind.Type); err != nil {
		return fmt.Errorf("context is not wellformed: aliased type %s is not wellformed: %v", bind.Type, err)
	}
	return nil
}

func IsWellformedComponentBind(newContext Context, bind *componentBind) error {
	if ok := newContext.ContainsComponentBind(bind.ElemType); ok {
		return fmt.Errorf("context is not wellformed: component %q is defined more than once", bind.ElemType)
	}
	if err := IsWellformedType(newContext, bind.ElemType); err != nil {
		return fmt.Errorf("context is not wellformed: component type %s is not wellformed: %v", bind.ElemType, err)
	}
	return nil
}

// TODO: Check that the definition matches the declaration.
func IsWellformedNameBind(newContext Context, bind *nameBind) error {
	if _, ok := newContext.LookupBind(bind.Name, FindDefOnly); ok {
		return fmt.Errorf("context is not wellformed: type %q is defined more than once", bind.Name)
	}
	return nil
}

func IsWellformedTypeVarBind(newContext Context, bind *typeVarBind) error {
	if newContext.ContainsTypeVarBind(bind.Name) {
		return fmt.Errorf("context is not wellformed: type variable %q is defined more than once\ncontext: %s", bind.Name, newContext.String())
	}
	return nil
}

func IsWellformedContext(context Context) error {
	if context.Empty() {
		// Rule: EmptyCtx.
		return nil
	}

	if context.wellformedSize == context.list.Size() {
		return nil
	}

	bind, newContext := context.Pop()

	if err := IsWellformedContext(newContext); err != nil {
		return err
	}

	switch bind.Case {
	case TermBind:
		return IsWellformedTermBind(newContext, bind.Term)
	case AliasBind:
		return IsWellformedAliasBind(newContext, bind.Alias)
	case ComponentBind:
		return IsWellformedComponentBind(newContext, bind.Component)
	case NameBind:
		return IsWellformedNameBind(newContext, bind.Name)
	case TypeVarBind:
		return IsWellformedTypeVarBind(newContext, bind.TypeVar)
	default:
		panic(fmt.Errorf("unhandled %T %d", bind.Case, bind.Case))
	}
}
