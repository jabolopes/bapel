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
