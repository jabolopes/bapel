package typer

import "fmt"

type BindCase int

const (
	JudgeBind BindCase = iota
	TermBind
	TypeBind
)

type judgeBind struct {
	Judge Judge
}

type termBind struct {
	ID   string
	Type Type
}

type typeBind struct {
	Type Type
}

type Bind struct {
	Case  BindCase
	Judge *judgeBind
	Term  *termBind
	Type  *typeBind
}

func (b Bind) String() string {
	{
		var d Bind
		if b == d {
			return ""
		}
	}

	switch b.Case {
	case JudgeBind:
		return b.Judge.Judge.String()
	case TermBind:
		return fmt.Sprintf("%s : %s", b.Term.ID, b.Term.Type)
	case TypeBind:
		return b.Type.Type.String()
	default:
		panic(fmt.Errorf("unhandled %T %d", b.Case, b.Case))
	}
}

func (b Bind) Is(c BindCase) bool {
	return b.Case == c
}

func NewJudgeBind(judge Judge) Bind {
	return Bind{
		Case:  JudgeBind,
		Judge: &judgeBind{judge},
	}
}

func NewTermBind(id string, typ Type) Bind {
	return Bind{
		Case: TermBind,
		Term: &termBind{id, typ},
	}
}

func NewTypeBind(typ Type) Bind {
	return Bind{
		Case: TypeBind,
		Type: &typeBind{typ},
	}
}
