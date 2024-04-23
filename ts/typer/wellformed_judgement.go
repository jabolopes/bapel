package typer

import "fmt"

func WellformedJudge(context Context, judge Judge) error {
	switch judge.Case {
	case ApplicationInferenceJudge:
		c := judge.ApplicationInference
		if err := WellformedType(context, c.Type); err != nil {
			return err
		}
		if err := WellformedTerm(context, c.Term); err != nil {
			return err
		}
		context = context.AddType(NewVarType(c.Var))
		return WellformedJudge(context, c.Judge)

	case CheckJudge:
		c := judge.Check
		if err := WellformedTerm(context, c.Term); err != nil {
			return err
		}
		return WellformedType(context, c.Type)

	case InferenceJudge:
		c := judge.Inference
		if err := WellformedTerm(context, c.Term); err != nil {
			return err
		}

		context = context.AddType(NewVarType(c.Var))
		return WellformedJudge(context, c.Judge)

	case SubtypeJudge:
		c := judge.Subtype
		if err := WellformedType(context, c.Left); err != nil {
			return err
		}
		return WellformedType(context, c.Right)

	default:
		panic(fmt.Errorf("unhandled %T %d", judge.Case, judge.Case))
	}
}
