package bplparser

type IsFunction interface {
	IsFunction(string) bool
}

type Parser struct {
	compiler IsFunction
}

func NewParser(compiler IsFunction) *Parser {
	return &Parser{compiler}
}
