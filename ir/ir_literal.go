package ir

import "fmt"

type IrLiteralCase int

const (
	IntLiteral IrLiteralCase = iota
	FloatLiteral
	RuneLiteral
	StrLiteral
)

func (c IrLiteralCase) String() string {
	switch c {
	case IntLiteral:
		return "integer"
	case FloatLiteral:
		return "float"
	case RuneLiteral:
		return "rune"
	case StrLiteral:
		return "string"
	default:
		panic(fmt.Errorf("unhandled %T %d", c, c))
	}
}

type floatLiteral struct {
	Integer int64
	Decimal int64
}

func (l *floatLiteral) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "%d.%d", l.Integer, l.Decimal)
}

type IrLiteral struct {
	Case  IrLiteralCase
	Int   *int64
	Float *floatLiteral
	Rune  *string
	Str   *string
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
	case FloatLiteral:
		return fmt.Sprintf("%s", l.Float)
	case RuneLiteral:
		return fmt.Sprintf(`'%s'`, *l.Rune)
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

func NewIntLiteral(pos Pos, value int64) IrLiteral {
	return IrLiteral{
		Case: IntLiteral,
		Int:  &value,
		Pos:  pos,
	}
}

func NewFloatLiteral(pos Pos, integer, decimal int64) IrLiteral {
	return IrLiteral{
		Case:  FloatLiteral,
		Float: &floatLiteral{integer, decimal},
		Pos:   pos,
	}
}

func NewRuneLiteral(pos Pos, value string) IrLiteral {
	return IrLiteral{
		Case: RuneLiteral,
		Rune: &value,
		Pos:  pos,
	}
}

func NewStrLiteral(pos Pos, value string) IrLiteral {
	return IrLiteral{
		Case: StrLiteral,
		Str:  &value,
		Pos:  pos,
	}
}
