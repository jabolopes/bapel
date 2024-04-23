package typer

import "fmt"

func freeVarsImpl(t Type, bound map[string]struct{}, free *map[string]struct{}) {
	switch t.Case {
	case ExistVarType:
		c := t.ExistVar
		(*free)[c.Var] = struct{}{}

	case ForallType:
		c := t.Forall
		bound[c.Var] = struct{}{}
		freeVarsImpl(c.Type, bound, free)

	case FunType:
		c := t.Fun
		freeVarsImpl(c.Arg, bound, free)
		freeVarsImpl(c.Ret, bound, free)

	case NameType:
		break

	case VarType:
		c := t.Var
		if _, ok := bound[c.Var]; !ok {
			(*free)[c.Var] = struct{}{}
		}

	default:
		panic(fmt.Errorf("unhandled %T %d", t.Case, t.Case))
	}
}

func FreeVars(t Type) map[string]struct{} {
	free := map[string]struct{}{}
	freeVarsImpl(t, map[string]struct{}{}, &free)
	return free
}

func ContainsFreeVar(t Type, id string) bool {
	freeVars := FreeVars(t)
	_, ok := freeVars[id]
	return ok
}
