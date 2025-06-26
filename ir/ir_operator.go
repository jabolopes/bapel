package ir

func IsOperator(id string) bool {
	switch id {
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
