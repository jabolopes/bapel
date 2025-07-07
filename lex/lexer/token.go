package lexer

type TokenType int

type Token struct {
	LineNum int
	Type    TokenType
	Value   string
}
