package ir

import (
	"fmt"
	"strings"

	"github.com/jabolopes/bapel/parser"
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

func ParseFunctionType(token string) (IrFunctionType, error) {
	splits := strings.SplitN(token, " -> ", 2)
	if len(splits) != 2 {
		return IrFunctionType{}, fmt.Errorf("invalid type; expected '(arg1 type1, ...) -> (ret1 type1, ...)'; got %q", token)
	}

	arg := splits[0]
	ret := splits[1]

	if err := parser.TrimPrefix(&arg, "(", fmt.Errorf("expected argument list in type; got %v", token)); err != nil {
		return IrFunctionType{}, err
	}

	if err := parser.TrimSuffix(&arg, ")", fmt.Errorf("expected argument list in type; got %v", token)); err != nil {
		return IrFunctionType{}, err
	}

	if err := parser.TrimPrefix(&ret, "(", fmt.Errorf("expected return value list in type; got %v", token)); err != nil {
		return IrFunctionType{}, err
	}

	if err := parser.TrimSuffix(&ret, ")", fmt.Errorf("expected return value list in type; got %v", token)); err != nil {
		return IrFunctionType{}, err
	}

	var args []string
	if len(arg) > 0 {
		args = strings.Split(arg, ", ")
	}

	var rets []string
	if len(ret) > 0 {
		rets = strings.Split(ret, ", ")
	}

	var argTypes []IrType
	for _, arg := range args {
		typ, err := ParseType(arg)
		if err != nil {
			return IrFunctionType{}, err
		}

		argTypes = append(argTypes, typ)
	}

	var retTypes []IrType
	for _, ret := range rets {
		typ, err := ParseType(ret)
		if err != nil {
			return IrFunctionType{}, err
		}

		retTypes = append(retTypes, typ)
	}

	return IrFunctionType{argTypes, retTypes}, nil
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
