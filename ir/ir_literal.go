package ir

import "fmt"

type IrLiteralCase int

const (
	IntLiteral IrLiteralCase = iota
	StrLiteral
)

func (c IrLiteralCase) String() string {
	switch c {
	case IntLiteral:
		return "integer"
	case StrLiteral:
		return "string"
	default:
		panic(fmt.Errorf("unhandled %T %d", c, c))
	}
}

type IrLiteral struct {
	Case IrLiteralCase
	Int  *int64
	Str  *string
	// Position in source file.
	Pos Pos
}

func (l IrLiteral) Is(c IrLiteralCase) bool {
	return l.Case == c
}

func (l IrLiteral) String() string {
	if l.Case == 0 && l.Int == nil {
		return ""
	}

	switch l.Case {
	case IntLiteral:
		return fmt.Sprintf("%d", *l.Int)
	case StrLiteral:
		return fmt.Sprintf(`"%s"`, *l.Str)
	default:
		panic(fmt.Errorf("unhandled %T %d", l.Case, l.Case))
	}
}

func (l IrLiteral) Format(f fmt.State, verb rune) {
	if l.Case == 0 && l.Int == nil {
		return
	}

	if addMetadata := f.Flag('+'); addMetadata {
		l.Pos.Format(f, verb)
	}

	fmt.Fprint(f, l.String())
}

func NewIntLiteral(value int64) IrLiteral {
	return IrLiteral{
		Case: IntLiteral,
		Int:  &value,
	}
}

func NewStrLiteral(value string) IrLiteral {
	return IrLiteral{
		Case: StrLiteral,
		Str:  &value,
	}
}
