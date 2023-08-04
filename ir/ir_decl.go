package ir

import (
	"fmt"
	"strings"

	"github.com/jabolopes/bapel/shift"
)

type irDeclType int

const (
	FunctionDecl = irDeclType(iota)
	VariableDecl
)

func (t irDeclType) String() string {
	switch t {
	case FunctionDecl:
		return "function"
	case VariableDecl:
		return "variable"
	default:
		panic(fmt.Errorf("Unhandled irDeclType %d", t))
	}
}

type irDecl struct {
	id       string
	declType irDeclType
	varType  IrIntType
	funType  IrFunctionType
}

func ParseDecl(args []string) (irDecl, error) {
	id, args, err := shift.Shift(args, fmt.Errorf("expected identifier as first token in declaration; got %v", args))
	if err != nil {
		return irDecl{}, err
	}

	args, err = shift.ShiftIf(args, ":", fmt.Errorf("expected token ':' after the declaration's identifier; got %v", args))
	if err != nil {
		return irDecl{}, err
	}

	if len(args) == 0 {
		return irDecl{}, fmt.Errorf("expected type in declaration; got %v", args)
	}

	typStr := strings.Join(args, " ")
	for _, arg := range args {
		if arg == "->" {
			// Declare function.
			typ, err := ParseFunctionType(typStr)
			if err != nil {
				return irDecl{}, err
			}

			return irDecl{id, FunctionDecl, 0, typ}, nil
		}
	}

	// Declare variable.
	typ, err := ParseType(typStr)
	if err != nil {
		return irDecl{}, err
	}

	return irDecl{id, VariableDecl, typ, IrFunctionType{}}, nil
}

// matchesDecl determines if the types of the actual declaration are
// equal to the types of the formal declaration. The name of the
// callee is taken from the formal declaration and ignored in the
// actual declaration.
func matchesDeclImpl(formal, actual irDecl, widen bool) error {
	id := formal.id

	if formal.declType != actual.declType {
		return fmt.Errorf("symbol %q expects declaration type %s; got %s", id, formal.declType, actual.declType)
	}

	if formal.declType == VariableDecl {
		if widen {
			if formal.varType < actual.varType {
				return fmt.Errorf("variable %q expects type %s or wider; got %s", id, formal.varType, actual.varType)
			}
		} else {
			if formal.varType != actual.varType {
				return fmt.Errorf("variable %q expects type %s; got %s", id, formal.varType, actual.varType)
			}
		}
		return nil
	}

	formalType := formal.funType
	actualType := actual.funType
	if len(formalType.Args) != len(actualType.Args) {
		return fmt.Errorf("function %q expects %d argument(s); got %q", id, len(formalType.Args), actualType.Args)
	}

	if len(formalType.Rets) != len(actualType.Rets) {
		return fmt.Errorf("function %q expects %d return value(s); got %q", id, len(formalType.Rets), actualType.Rets)
	}

	for i := range formalType.Args {
		if formalType.Args[i] != actualType.Args[i] {
			return fmt.Errorf("function %q expects argument %d with type %d; got %d", id, i, formalType.Args[i], actualType.Args[i])
		}
	}

	for i := range formalType.Rets {
		if formalType.Rets[i] != actualType.Rets[i] {
			return fmt.Errorf("function %q expects return value %d with type %d; got %d", id, i, formalType.Rets[i], actualType.Rets[i])
		}
	}

	return nil
}

func matchesDecl(formal, actual irDecl) error {
	return matchesDeclImpl(formal, actual, false /* widen */)
}

func matchesDeclWiden(formal, actual irDecl) error {
	return matchesDeclImpl(formal, actual, true /* widen */)
}
