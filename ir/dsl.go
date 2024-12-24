package ir

import "fmt"

func CallID(id string, args ...IrTerm) IrTerm {
	return Call(ID(id), args...)
}

func Call(id IrTerm, args ...IrTerm) IrTerm {
	return NewAppTermTerm(id, NewTupleTerm(args))
}

func CallPF(id IrTerm, types []IrType, args ...IrTerm) IrTerm {
	term := id
	for _, typ := range types {
		term = NewAppTypeTerm(term, typ)
	}
	if len(args) == 0 {
		return term
	}
	return NewAppTermTerm(term, NewTupleTerm(args))
}

func Const(id string) IrType {
	return NewNameType(id)
}

func Forall(tvar string, kind IrKind, typ IrType) IrType {
	return NewForallType(tvar, kind, typ)
}

type VarKind struct {
	Var  string
	Kind IrKind
}

// ForallVars creates a nested forall type for each type variable.
//
// Example:
//   NewForallVarsType(['a, 'b, 'c], 'a -> 'b -> 'c) =
//     forall 'a. (forall 'b. (forall 'c. 'a -> 'b -> 'c))
func ForallVars(tvars []VarKind, typ IrType) IrType {
	for i := len(tvars) - 1; i >= 0; i-- {
		typ = NewForallType(tvars[i].Var, tvars[i].Kind, typ)
	}
	return typ
}

func LambdaVars(tvars []VarKind, typ IrType) IrType {
	if len(tvars) == 0 {
		return typ
	}

	for i := len(tvars) - 1; i >= 0; i-- {
		typ = NewLambdaType(tvars[i].Var, tvars[i].Kind, typ)
	}
	return typ
}

func Tvar(tvar string) IrType {
	return NewVarType(tvar)
}

func Terms(terms ...IrTerm) IrTerm {
	return NewTupleTerm(terms)
}

func Types(types ...IrType) IrType {
	return NewTupleType(types)
}

func TypesA(types ...IrType) []IrType {
	return append([]IrType{}, types...)
}

func ID(id string) IrTerm {
	return NewLiteralTerm(IDLiteral, id, 0)
}

func Number(value int64) IrTerm {
	return NewLiteralTerm(NumberLiteral, fmt.Sprintf("%d", value), value)
}
