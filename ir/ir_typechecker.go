package ir

import (
	"fmt"
	"sort"

	"github.com/jabolopes/bapel/parser"
)

type IrTypechecker struct {
	context *IrContext
	widen   bool
}

func (t *IrTypechecker) MatchesArrayType(formal, actual IrArrayType) error {
	if err := t.MatchesType(formal.ElemType, actual.ElemType); err != nil {
		return fmt.Errorf("mismatch in array element types: %v", err)
	}

	if formal.Size != actual.Size {
		return fmt.Errorf("expected array with %d elements; got %d elements", formal.Size, actual.Size)
	}

	return nil
}

func (t *IrTypechecker) MatchesFunctionType(formal, actual IrFunctionType) error {
	if len(formal.Args) != len(actual.Args) {
		return fmt.Errorf("expected function with %d argument(s); got %q", len(formal.Args), actual.Args)
	}

	if len(formal.Rets) != len(actual.Rets) {
		return fmt.Errorf("expected function with %d return value(s); got %q", len(formal.Rets), actual.Rets)
	}

	for i := range formal.Args {
		if err := t.MatchesType(formal.Args[i], actual.Args[i]); err != nil {
			return fmt.Errorf("in function argument %d: %v", i+1, err)
		}
	}

	for i := range formal.Rets {
		if err := t.MatchesType(formal.Rets[i], actual.Rets[i]); err != nil {
			return fmt.Errorf("in return value %d: %v", i, err)
		}
	}

	return nil
}

func (t *IrTypechecker) MatchesIntType(formal, actual IrIntType) error {
	if t.widen {
		if formal < actual {
			return fmt.Errorf("expected type %s or wider; got %s", formal, actual)
		}
	} else {
		if formal != actual {
			return fmt.Errorf("expected type %s; got %s", formal, actual)
		}
	}
	return nil
}

func (t *IrTypechecker) MatchesStructType(formal, actual IrStructType) error {
	if len(formal.Fields) != len(actual.Fields) {
		return fmt.Errorf("expected %d fields; got %d", len(formal.Fields), len(actual.Fields))
	}

	formalFields := append([]StructField{}, formal.Fields...)
	actualFields := append([]StructField{}, actual.Fields...)

	sort.Slice(formalFields, func(i, j int) bool {
		return formalFields[i].Name < formalFields[j].Name
	})
	sort.Slice(actualFields, func(i, j int) bool {
		return actualFields[i].Name < actualFields[j].Name
	})

	for i := range formalFields {
		if formalFields[i].Name != actualFields[i].Name {
			return fmt.Errorf("expected field names %v; got %v", formal.Names(), actual.Names())
		}

		if err := t.MatchesType(formalFields[i].Type, actualFields[i].Type); err != nil {
			return err
		}
	}

	return nil
}

func (t *IrTypechecker) MatchesIDType(formal, actual string) error {
	formalDecl, err := t.context.getDecl(formal, FindAny)
	if err != nil {
		return err
	}

	actualDecl, err := t.context.getDecl(actual, FindAny)
	if err != nil {
		return err
	}

	return t.MatchesDecl(formalDecl, actualDecl)
}

func (t *IrTypechecker) MatchesType(formal, actual IrType) error {
	if formal.Case != actual.Case {
		return fmt.Errorf("expected type %s; got %s", formal.Case, actual.Case)
	}

	switch formal.Case {
	case ArrayType:
		return t.MatchesArrayType(*formal.ArrayType, *actual.ArrayType)
	case FunType:
		return t.MatchesFunctionType(formal.FunType, actual.FunType)
	case IntType:
		return t.MatchesIntType(formal.IntType, actual.IntType)
	case StructType:
		return t.MatchesStructType(formal.StructType, actual.StructType)
	case IDType:
		return t.MatchesIDType(formal.IDType, actual.IDType)
	default:
		panic(fmt.Errorf("Unhandled IrTypeCase %d", formal.Case))
	}
}

// MatchesDecl determines if the types of the actual declaration are equal to
// the types of the formal declaration. The name of the callee is taken from the
// formal declaration and ignored in the actual declaration.
func (t *IrTypechecker) MatchesDecl(formal, actual IrDecl) error {
	if formal.Case != actual.Case {
		return fmt.Errorf("in declaration %q: expected %s; got %s", formal.ID, formal.Case, actual.Case)
	}

	if err := t.MatchesType(formal.Type, actual.Type); err != nil {
		return fmt.Errorf("in declaration %q: %v", formal.ID, err)
	}

	return nil
}

func (t *IrTypechecker) MatchesDeclWiden(formal, actual IrDecl) error {
	t.widen = true
	defer func() { t.widen = false }()
	return t.MatchesDecl(formal, actual)
}

func (t *IrTypechecker) CheckCallArg(formal IrType, arg parser.Token) error {
	switch arg.Case {
	case parser.IDToken:
		actualType, err := t.context.getType(arg.Text, FindAny)
		if err != nil {
			return err
		}
		return t.MatchesType(formal, actualType)

	case parser.NumberToken:
		if !formal.Is(IntType) {
			return fmt.Errorf("expected type %s; got %d", formal, arg.Value)
		}

	default:
		panic(fmt.Errorf("Unhandled token %d", arg.Case))
	}

	return nil
}

func (t *IrTypechecker) CheckCallRet(formal IrType, arg string) error {
	actualDecl, err := t.context.getDecl(arg, FindAny)
	if err != nil {
		return err
	}

	if actualDecl.Case != VarDecl {
		return fmt.Errorf("expected return value declared as %s; got %q", VarDecl, actualDecl.Case)
	}

	return t.MatchesType(formal, actualDecl.Type)
}

func (t *IrTypechecker) CheckCall(id string, args []parser.Token, rets []string) error {
	formalType, err := t.context.getType(id, FindAny)
	if err != nil {
		return err
	}

	if formalType.Case != FunType {
		return fmt.Errorf("expected function type; got %s", formalType)
	}

	if len(formalType.FunType.Args) != len(args) {
		return fmt.Errorf("expected %d arguments; got %d", len(formalType.FunType.Args), len(args))
	}

	for i := range formalType.FunType.Args {
		formalArg := formalType.FunType.Args[i]
		actualArg := args[i]
		if err := t.CheckCallArg(formalArg, actualArg); err != nil {
			return fmt.Errorf("in argument %d of function %s: %v", i+1, id, err)
		}
	}

	if len(formalType.FunType.Rets) != len(rets) {
		return fmt.Errorf("expected %d return values; got %d", len(formalType.FunType.Rets), len(rets))
	}

	for i := range formalType.FunType.Rets {
		formalRet := formalType.FunType.Rets[i]
		actualRet := rets[i]
		if err := t.CheckCallRet(formalRet, actualRet); err != nil {
			return fmt.Errorf("in return value %d of function %s: %v", i+1, id, err)
		}
	}

	return nil
}

func (t *IrTypechecker) CheckIfVar(arg string) error {
	typ, err := t.context.getType(arg, FindAny)
	if err != nil {
		return err
	}

	if !typ.Is(IntType) {
		return fmt.Errorf("expected integer type; got %v", typ)
	}

	return nil
}

func (t *IrTypechecker) CheckSingleAssign(arg parser.Token, ret string) error {
	retDecl, err := t.context.getDecl(ret, FindAny)
	if err != nil {
		return err
	}
	if retDecl.Case != VarDecl {
		return fmt.Errorf("expected return value declared as %s; got %q", VarDecl, retDecl.Case)
	}

	return t.CheckCallArg(retDecl.Type, arg)
}

func (t *IrTypechecker) CheckWiden(arg parser.Token, ret string) error {
	retDecl, err := t.context.getDecl(ret, FindAny)
	if err != nil {
		return err
	}
	if retDecl.Case != VarDecl {
		return fmt.Errorf("expected return value declared as %s; got %q", VarDecl, retDecl.Case)
	}

	argDecl, err := t.context.getDecl(arg.Text, FindAny)
	if err != nil {
		return err
	}

	return t.MatchesDeclWiden(retDecl, argDecl)
}

func NewIrTypechecker(context *IrContext) *IrTypechecker {
	return &IrTypechecker{context, false /* widen */}
}
