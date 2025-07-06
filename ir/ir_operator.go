package ir

func IsOperator(id string) bool {
	switch id {
	case "&&":
	case "!=":
	case "==":
	case ">":
	case ">=":
	case "<":
	case "<=":
	case "+":
	case "-":
	case "*":
	case "/":
	case "!":
	default:
		return false
	}
	return true
}

func OperatorType(id string) IrType {
	// forall 'a. ('a, 'a) -> bool
	comparison := Forall("a", NewTypeKind(), NewFunctionType(NewTupleType([]IrType{NewVarType("a"), NewVarType("a")}), NewNameType("bool")))

	// forall 'a. ('a, 'a) -> 'a
	additive := Forall("a", NewTypeKind(), NewFunctionType(NewTupleType([]IrType{NewVarType("a"), NewVarType("a")}), NewVarType("a")))

	// bool -> bool
	logicalUnary := NewFunctionType(NewNameType("bool"), NewNameType("bool"))

	// (bool, bool) -> bool
	logicalBinary := NewFunctionType(NewTupleType([]IrType{NewNameType("bool"), NewNameType("bool")}), NewNameType("bool"))

	switch id {
	case "&&":
		return logicalBinary
	case "!=", "==", ">", ">=", "<", "<=":
		return comparison
	case "+", "-", "*", "/":
		return additive
	case "!":
		return logicalUnary
	default:
		return NewTupleType(nil)
	}
}
