package typer

import "fmt"

// substituteVar substitutes type variables or existential variables with the
// given replacement type.
func substituteVar(t Type, id string, replacement Type) Type {
	switch t.Case {
	case ExistVarType:
		c := *t.ExistVar
		if c.Var == id {
			return replacement
		}
		return t

	case ForallType:
		c := t.Forall
		if c.Var != id {
			return NewForallType(c.Var, substituteVar(c.Type, id, replacement))
		}
		return t

	case FunType:
		c := t.Fun
		return NewFunType(
			substituteVar(c.Arg, id, replacement),
			substituteVar(c.Ret, id, replacement))

	case NameType:
		return t

	case VarType:
		c := t.Var
		if c.Var == id {
			return replacement
		}
		return t

	default:
		panic(fmt.Errorf("unhandled %T %d", t.Case, t.Case))
	}
}
