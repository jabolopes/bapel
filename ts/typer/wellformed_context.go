package typer

import "fmt"

func WellformedContext(context Context) error {
	if context.list.Empty() {
		return nil
	}

	bind, newContext := context.Pop()
	if err := WellformedContext(newContext); err != nil {
		return err
	}

	switch bind.Case {
	case JudgeBind:
		return WellformedJudge(newContext, bind.Judge.Judge)
	case TermBind:
		return WellformedType(newContext, bind.Term.Type)
	case TypeBind:
		return nil
	default:
		panic(fmt.Errorf("unhandled %T %d", bind.Case, bind.Case))
	}
}
