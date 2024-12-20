package bplparser2

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/parser"
	"github.com/jabolopes/go-lalr1"
	"github.com/jabolopes/go-lalr1/grammar"
)

type Pos struct {
	lineNum int
	line    string
}

type Token struct {
	Pos   Pos
	Token parser.Token
}

func NewWithInitialSymbol(symbol string) (*lalr1.Parser, error) {
	production := fmt.Sprintf("program -> %s eof", symbol)
	return lalr1.NewParser(NewGrammar(grammar.ProductionLine{production, first}))
}

type Parser struct {
	initialSymbol string
	reader        io.Reader
}

func NewParser() *Parser {
	return &Parser{"Anys", nil}
}

func (p *Parser) SetInitialSymbol(symbol string) {
	p.initialSymbol = symbol
}

func (p *Parser) Open(reader io.Reader) {
	p.reader = reader
}

// TODO: Replace np with p.
func Parse[T any](np *Parser) (T, error) {
	var t T

	if np.reader == nil {
		return t, errors.New("Parser.Open() must be called with a valid reader prior to calling Parse()")
	}

	var parser *lalr1.Parser
	{
		p, err := NewWithInitialSymbol(np.initialSymbol)
		if err != nil {
			return t, err
		}
		parser = p
	}

	if conflicts := parser.Machine().Conflicts(); len(conflicts) > 0 {
		// Return an error if there are conflicts.
		var str strings.Builder
		str.WriteString("Grammar has conflicts:\n")
		for _, conflict := range conflicts {
			str.WriteString(fmt.Sprintf("  * %s\n", conflict))
		}
		return t, errors.New(str.String())
	}

	fmt.Println(parser.Machine())
	fmt.Println(parser.ParseTable())

	p := bplparser.NewParser()
	p.Open(np.reader)

	// TODO: Fix.
	channel := make(chan lalr1.Token, 10000)

	brackets := 0

	for p.Scan() {
		isSingleExpression := true
		isEmpty := true

		pos := Pos{p.LineNum(), p.Line()}

		for {
			parserToken, err := p.ShiftToken()
			if err == io.EOF {
				break
			}
			if err != nil {
				return t, fmt.Errorf("in line %d:\n  %s\n%v", p.LineNum(), p.Line(), err)
			}

			isEmpty = false

			token := Token{pos, parserToken}
			switch token.Token.Text {
			case "{":
				isSingleExpression = false
				brackets++
			case "}":
				isSingleExpression = false
				brackets--
			}

			log.Printf("HERE %v", token)

			if tokenType, ok := parser.ParseTable().GetTokenType(token.Token.Text); ok {
				channel <- lalr1.Token{tokenType, token}
			} else {
				channel <- lalr1.Token{parser.ParseTable().TokenType("Token"), token}
			}
		}

		log.Printf("HERE2 %v && !%v && %v", brackets, isEmpty, isSingleExpression)
		if brackets > 0 && !isEmpty && isSingleExpression {
			token := lalr1.Token{parser.ParseTable().TokenType(";"), Token{Pos: pos}}
			log.Printf("HERE %v", token)
			channel <- token
		}
	}

	{
		pos := Pos{p.LineNum(), p.Line()}
		token := lalr1.Token{parser.ParseTable().TokenType("eof"), Token{Pos: pos}}
		log.Printf("HERE %v", token)
		channel <- token
	}

	close(channel)

	ast, output, err := parser.Parse(channel)
	if err != nil {
		gotToken := output.Got.Data.(Token)

		fmt.Printf("In line %d:\n  %s\n expected: %v\n got: %v\n",
			gotToken.Pos.lineNum,
			gotToken.Pos.line,
			output.Expected,
			output.GotSymbol.Name)

		return t, err
	}
	fmt.Printf("AST: %v\n", ast)

	return ast.(T), nil
}

func ParseFile(input io.Reader) ([]bplparser.Source, error) {
	parser := NewParser()
	parser.Open(input)
	return Parse[[]bplparser.Source](parser)
}
