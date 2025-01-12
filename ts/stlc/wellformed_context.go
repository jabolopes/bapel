package stlc

import (
	"fmt"
)

func isWellformedContext(context Context) error {
	if context.empty() {
		// Rule: EmptyCtx.
		return nil
	}

	if context.wellformedSize == context.list.Size() {
		return nil
	}

	bind, newContext := context.pop()

	if err := isWellformedContext(newContext); err != nil {
		return err
	}

	switch bind.Case {
	case TermBind:
		return isWellformedTermBind(newContext, bind.Term)
	case AliasBind:
		return isWellformedAliasBind(newContext, bind.Alias)
	case ComponentBind:
		return isWellformedComponentBind(newContext, bind.Component)
	case ConstBind:
		return isWellformedConstBind(newContext, bind.Const)
	case TypeVarBind:
		return isWellformedTypeVarBind(newContext, bind.TypeVar)
	default:
		panic(fmt.Errorf("unhandled %T %d", bind.Case, bind.Case))
	}
}
