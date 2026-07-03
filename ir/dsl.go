package ir

func Forall(tvar string, kind IrKind, typ IrType) IrType {
	return NewForallType(TypeParam{Var: tvar, Kind: kind}, typ)
}

// ForallVars creates a nested forall type for each type variable.
//
// Example:
//
//	ForallVars(['a, 'b, 'c], 'a -> 'b -> 'c) =
//	  forall 'a. (forall 'b. (forall 'c. 'a -> 'b -> 'c))
func ForallVars(tvars []TypeParam, typ IrType) IrType {
	for i := len(tvars) - 1; i >= 0; i-- {
		typ = NewForallType(tvars[i], typ)
	}
	return typ
}

func LambdaVars(tvars []TypeParam, typ IrType) IrType {
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

func Tvar(tvar string) IrType {
	return NewVarType(tvar)
}
