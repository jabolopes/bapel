package ir

import (
	"fmt"
	"strings"
)

type IrFunctionType struct {
	Args []IrType
	Rets []IrType
}

func (t IrFunctionType) String() string {
	var builder strings.Builder

	switch len(t.Args) {
	case 0:
		builder.WriteString("()")
	case 1:
		builder.WriteString(fmt.Sprintf("(%s)", t.Args[0]))
	default:
		builder.WriteString(fmt.Sprintf("(%s", t.Args[0]))
		for _, typ := range t.Args[1:] {
			builder.WriteString(fmt.Sprintf(", %s", typ))
		}
		builder.WriteString(")")
	}

	builder.WriteString(" -> ")

	switch len(t.Rets) {
	case 0:
		builder.WriteString("()")
	case 1:
		builder.WriteString(fmt.Sprintf("(%s)", t.Rets[0]))
	default:
		builder.WriteString(fmt.Sprintf("(%s", t.Rets[0]))
		for _, typ := range t.Rets[1:] {
			builder.WriteString(fmt.Sprintf(", %s", typ))
		}
		builder.WriteString(")")
	}

	return builder.String()
}

func MatchesFunctionType(formal, actual IrFunctionType) error {
	if len(formal.Args) != len(actual.Args) {
		return fmt.Errorf("expected function with %d argument(s); got %q", len(formal.Args), actual.Args)
	}

	if len(formal.Rets) != len(actual.Rets) {
		return fmt.Errorf("expected function with %d return value(s); got %q", len(formal.Rets), actual.Rets)
	}

	for i := range formal.Args {
		if err := MatchesType(formal.Args[i], actual.Args[i], false /* widen */); err != nil {
			return fmt.Errorf("in function argument %d: %v", i+1, err)
		}
	}

	for i := range formal.Rets {
		if err := MatchesType(formal.Rets[i], actual.Rets[i], false /* widen */); err != nil {
			return fmt.Errorf("in return value %d: %v", i, err)
		}
	}

	return nil
}
