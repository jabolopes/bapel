package ir

import (
	"fmt"
	"sort"
)

type IrTypechecker struct {
	context *IrContext
}

// matchesDecl determines if the types of the actual declaration are
// equal to the types of the formal declaration. The name of the
// callee is taken from the formal declaration and ignored in the
// actual declaration.
func (t *IrTypechecker) matchesDeclImpl(formal, actual irDecl, widen bool) error {
	if err := t.MatchesType(formal.typ, actual.typ, widen); err != nil {
		return fmt.Errorf("symbol %q definition %v does not match its declaration type %v typ: %v", formal.id, actual.typ, formal.typ, err)
	}

	return nil
}

func (t *IrTypechecker) MatchesArrayType(formal, actual IrArrayType, widen bool) error {
	if err := t.MatchesType(formal.ElemType, actual.ElemType, widen); err != nil {
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
		if err := t.MatchesType(formal.Args[i], actual.Args[i], false /* widen */); err != nil {
			return fmt.Errorf("in function argument %d: %v", i+1, err)
		}
	}

	for i := range formal.Rets {
		if err := t.MatchesType(formal.Rets[i], actual.Rets[i], false /* widen */); err != nil {
			return fmt.Errorf("in return value %d: %v", i, err)
		}
	}

	return nil
}

func (t *IrTypechecker) MatchesIntType(formal, actual IrIntType, widen bool) error {
	if widen {
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

func (t *IrTypechecker) MatchesStructType(formal, actual IrStructType, widen bool) error {
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

		if err := t.MatchesType(formalFields[i].Type, actualFields[i].Type, widen); err != nil {
			return err
		}
	}

	return nil
}

func (t *IrTypechecker) MatchesType(formal, actual IrType, widen bool) error {
	if formal.Case != actual.Case {
		return fmt.Errorf("expected type %s; got %s", formal.Case, actual.Case)
	}

	switch formal.Case {
	case ArrayType:
		return t.MatchesArrayType(*formal.ArrayType, *actual.ArrayType, widen)
	case FunType:
		return t.MatchesFunctionType(formal.FunType, actual.FunType)
	case IntType:
		return t.MatchesIntType(formal.IntType, actual.IntType, widen)
	case StructType:
		return t.MatchesStructType(formal.StructType, actual.StructType, widen)
	default:
		panic(fmt.Errorf("Unhandled IrTypeCase %d", formal.Case))
	}
}

func (t *IrTypechecker) MatchesDecl(formal, actual irDecl) error {
	return t.matchesDeclImpl(formal, actual, false /* widen */)
}

func (t *IrTypechecker) MatchesDeclWiden(formal, actual irDecl) error {
	return t.matchesDeclImpl(formal, actual, true /* widen */)
}

func NewIrTypechecker(context *IrContext) *IrTypechecker {
	return &IrTypechecker{context}
}
