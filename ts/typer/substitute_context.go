package typer

import (
	"fmt"

	"github.com/jabolopes/bapel/ts/list"
)

func substituteVarInJudge(j Judge, id string, replacement Type) Judge {
	switch j.Case {
	case ApplicationInferenceJudge:
		c := *j.ApplicationInference
		return NewApplicationInferenceJudge(
			substituteVar(c.Type, id, replacement),
			c.Term,
			c.Var,
			substituteVarInJudge(c.Judge, id, replacement))

	case CheckJudge:
		c := *j.Check
		return NewCheckJudge(
			c.Term,
			substituteVar(c.Type, id, replacement))

	case InferenceJudge:
		c := *j.Inference
		return NewInferenceJudge(
			c.Term,
			c.Var,
			substituteVarInJudge(c.Judge, id, replacement))

	case SubtypeJudge:
		c := *j.Subtype
		return NewSubtypeJudge(
			substituteVar(c.Left, id, replacement),
			substituteVar(c.Right, id, replacement))

	default:
		panic(fmt.Errorf("unhandled %T %d", j.Case, j.Case))
	}
}

func substituteVarInBind(b Bind, id string, replacement Type) Bind {
	switch b.Case {
	case JudgeBind:
		c := *b.Judge
		return NewJudgeBind(substituteVarInJudge(c.Judge, id, replacement))

	case TermBind:
		c := *b.Term
		return NewTermBind(c.ID, substituteVar(c.Type, id, replacement))

	case TypeBind:
		c := *b.Type
		return NewTypeBind(substituteVar(c.Type, id, replacement))

	default:
		panic(fmt.Errorf("unhandled %T %d", b.Case, b.Case))
	}
}

// SubstituteContext substitutes all occurences of existVar in the context
// (including those inside bindings, judgements, etc) with the given type.
func substituteVarInContext(c Context, id string, replacement Type) Context {
	c.list = list.Map(func(bind Bind) Bind {
		return substituteVarInBind(bind, id, replacement)
	}, c.list)
	return c
}
