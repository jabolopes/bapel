package lexerfsm

type TokenType int

type Token struct {
	Type  TokenType
	Value string
}
