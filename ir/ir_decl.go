package ir

import "fmt"

type irDecl struct {
	id   string
	args []IrType
	rets []IrType
}

// matchesDecl determines if the types of the actual declaration are
// equal to the types of the formal declaration. The name of the
// callee is taken from the formal declaration and ignored in the
// actual declaration.
func matchesDecl(formal, actual irDecl) error {
	id := formal.id

	if len(formal.args) != len(actual.args) {
		return fmt.Errorf("Function %q expects %d argument(s); got %q", id, len(formal.args), len(actual.args))
	}

	if len(formal.rets) != len(actual.rets) {
		return fmt.Errorf("Function %q expects %d return value(s); got %q", id, len(formal.rets), len(actual.rets))
	}

	for i := range formal.args {
		if formal.args[i] != actual.args[i] {
			return fmt.Errorf("Function %q expects argument %d with type %d; got %d", id, i, formal.args[i], actual.args[i])
		}
	}

	for i := range formal.rets {
		if formal.rets[i] != actual.rets[i] {
			return fmt.Errorf("Function %q expects return value %d with type %d; got %d", id, i, formal.rets[i], actual.rets[i])
		}
	}

	return nil
}
