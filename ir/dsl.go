package ir

type VarKind struct {
	Var  string
	Kind IrKind
}

type ArgType struct {
	Var  string
	Type IrType
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

func Number(value int64) IrTerm {
	return NewConstTerm(NewIntLiteral(value))
}

func ID(id string) IrTerm {
	return NewVarTerm(id)
}

func Forall(tvar string, kind IrKind, typ IrType) IrType {
	return NewForallType(tvar, kind, typ)
}

// ForallVars creates a nested forall type for each type variable.
//
// Example:
//
//	NewForallVarsType(['a, 'b, 'c], 'a -> 'b -> 'c) =
//	  forall 'a. (forall 'b. (forall 'c. 'a -> 'b -> 'c))
func ForallVars(tvars []VarKind, typ IrType) IrType {
	for i := len(tvars) - 1; i >= 0; i-- {
		typ = NewForallType(tvars[i].Var, tvars[i].Kind, typ)
	}
	return typ
}

func LambdaVars(tvars []VarKind, typ IrType) IrType {
	pos := typ.Pos

	if len(tvars) == 0 {
		return typ
	}

	for i := len(tvars) - 1; i >= 0; i-- {
		typ = NewLambdaType(tvars[i].Var, tvars[i].Kind, typ)
		typ.Pos = pos
	}
	return typ
}

func Lambda(tvars []VarKind, args []ArgType, body IrTerm) IrTerm {
	if len(tvars) > 0 {
		tvar := tvars[0]
		return NewTypeAbsTerm(tvar.Var, tvar.Kind, Lambda(tvars[1:], args, body))
	}

	if len(args) > 0 {
		arg := args[0]
		return NewLambdaTerm(arg.Var, arg.Type, Lambda(nil, args[1:], body))
	}

	return body
}

func Tvar(tvar string) IrType {
	return NewVarType(tvar)
}
