package ir

func Call(id IrTerm, types []IrType, args ...IrTerm) IrTerm {
	term := id
	for _, typ := range types {
		term = NewAppTypeTerm(term, typ)
	}
	if len(args) == 0 {
		return term
	}
	return NewAppTermTerm(term, NewTupleTerm(args))
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

func Lambda(tvars []VarKind, args []FunctionArg, body IrTerm) IrTerm {
	if len(tvars) > 0 {
		tvar := tvars[0]
		return NewTypeAbsTerm(tvar.Var, tvar.Kind, Lambda(tvars[1:], args, body))
	}

	if len(args) > 0 {
		return NewLambdaTerm(args[0], Lambda(nil, args[1:], body))
	}

	return body
}

func Tvar(tvar string) IrType {
	return NewVarType(tvar)
}
