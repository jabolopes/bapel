package lex

type TokenType int

const (
	WordToken TokenType = iota
	NumberToken
	RuneToken
	StringToken
	SymbolToken   = WordToken
	OperatorToken = WordToken
)
