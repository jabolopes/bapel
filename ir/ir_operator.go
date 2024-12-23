package ir

func IsOperator(id string) bool {
	return id == "+" || id == "-" || id == "*" || id == "/" || id == "!"
}
