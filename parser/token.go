package parser

import (
	"fmt"
	"strings"

	"golang.org/x/exp/constraints"
)

type TokenCase int

const (
	IDToken TokenCase = iota
	NumberToken
)

func (c TokenCase) String() string {
	switch c {
	case IDToken:
		return "identifier"
	case NumberToken:
		return "number"
	default:
		panic(fmt.Errorf("unhandled TokenCase %d", c))
	}
}

type Token struct {
	Case  TokenCase
	Text  string
	Value int64
}

func (t Token) String() string {
	return t.Text
}

func NewIDToken(text string) Token {
	return Token{IDToken, text, 0}
}

func NewNumberToken(value int64) Token {
	return Token{NumberToken, fmt.Sprintf("%d", value), value}
}

func ParseNumber[T constraints.Integer](arg string) (T, error) {
	var value T

	if strings.HasPrefix(arg, "0x") {
		// Hexadecimal
		_, err := fmt.Sscanf(arg, "0x%x", &value)

		return value, err
	}

	// Decimal.
	_, err := fmt.Sscanf(arg, "%d", &value)
	return value, err
}

func ParseToken(text string) (Token, error) {
	if value, err := ParseNumber[int64](text); err == nil {
		return Token{NumberToken, text, value}, nil
	}
	return Token{IDToken, text, 0}, nil
}

func ParseTokens(texts []string) ([]Token, error) {
	tokens := make([]Token, len(texts))
	for i, text := range texts {
		token, err := ParseToken(text)
		if err != nil {
			return nil, err
		}

		tokens[i] = token
	}

	return tokens, nil
}
