package typer

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

// JudgeCase is a judgement case. Denoted by w.
type JudgeCase int

const (
	// Application inference judgement.
	//
	// A • term =>=>a w
	ApplicationInferenceJudge JudgeCase = iota
	// Checking judgement.
	//
	// term <= A
	CheckJudge
	// Inference judgement.
	//
	// e =>a w
	InferenceJudge
	// Subtype judgement.
	//
	// A ≤ B
	SubtypeJudge
)

type applicationInference struct {
	Type  Type
	Term  ir.IrTerm
	Var   string
	Judge Judge
}

type check struct {
	Term ir.IrTerm
	Type Type
}

type inference struct {
	Term  ir.IrTerm
	Var   string
	Judge Judge
}

type subtype struct {
	Left  Type
	Right Type
}

// Judge is a judgement.
type Judge struct {
	Case                 JudgeCase
	ApplicationInference *applicationInference
	Check                *check
	Inference            *inference
	Subtype              *subtype
}

func (j Judge) String() string {
	{
		var d Judge
		if j == d {
			return ""
		}
	}

	switch j.Case {
	case ApplicationInferenceJudge:
		c := j.ApplicationInference
		return fmt.Sprintf("%s • %s =>=>%s (%s)", c.Type, c.Term, c.Var, c.Judge)
	case CheckJudge:
		c := j.Check
		return fmt.Sprintf("%s <= %s", c.Term, c.Type)
	case InferenceJudge:
		c := j.Inference
		return fmt.Sprintf("%s =>%s %s", c.Term, c.Var, c.Judge)
	case SubtypeJudge:
		c := j.Subtype
		return fmt.Sprintf("%s ≤ %s", c.Left, c.Right)
	default:
		panic(fmt.Errorf("unhandled %T %d", j.Case, j.Case))
	}
}

func (j Judge) Is(c JudgeCase) bool {
	return j.Case == c
}

func NewApplicationInferenceJudge(typ Type, term ir.IrTerm, tvar string, judge Judge) Judge {
	return Judge{
		Case:                 ApplicationInferenceJudge,
		ApplicationInference: &applicationInference{typ, term, tvar, judge},
	}
}

func NewCheckJudge(term ir.IrTerm, typ Type) Judge {
	return Judge{
		Case:  CheckJudge,
		Check: &check{term, typ},
	}
}

func NewInferenceJudge(term ir.IrTerm, tvar string, judge Judge) Judge {
	return Judge{
		Case:      InferenceJudge,
		Inference: &inference{term, tvar, judge},
	}
}

func NewSubtypeJudge(a, b Type) Judge {
	return Judge{
		Case:    SubtypeJudge,
		Subtype: &subtype{a, b},
	}
}
