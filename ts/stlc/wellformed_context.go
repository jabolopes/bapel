package stlc

import (
	"fmt"
)

func IsWellformedContext(context Context) error {
	if context.empty() {
		// Rule: EmptyCtx.
		return nil
	}

	if context.wellformedSize == context.list.Size() {
		return nil
	}

	bind, newContext := context.pop()

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
	case ConstBind:
		return IsWellformedConstBind(newContext, bind.Name)
	case TypeVarBind:
		return IsWellformedTypeVarBind(newContext, bind.TypeVar)
	default:
		panic(fmt.Errorf("unhandled %T %d", bind.Case, bind.Case))
	}
}
