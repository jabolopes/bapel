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

	var err error
	switch bind.Case {
	case TermBind:
		err = isWellformedTermBind(newContext, bind.Term)
	case AliasBind:
		err = isWellformedAliasBind(newContext, bind.Alias)
	case ConstBind:
		err = isWellformedConstBind(newContext, bind.Const)
	case TypeVarBind:
		err = isWellformedTypeVarBind(newContext, bind.TypeVar)
	default:
		panic(fmt.Errorf("unhandled %T %d", bind.Case, bind.Case))
	}

	if err != nil {
		return fmt.Errorf("context is not wellformed: %v", err)
	}

	return nil
}
