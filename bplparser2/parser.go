package bplparser2

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/lexer"
	"github.com/jabolopes/go-lalr1"
	"github.com/jabolopes/go-lalr1/grammar"
)

type Token struct {
	Pos  ir.Pos
	Text string
}

type ID struct {
	Pos   ir.Pos
	Value string
}

type Integer struct {
	Pos   ir.Pos
	Value int
}

func NewWithInitialSymbol(symbol string) (*lalr1.Parser, error) {
	production := fmt.Sprintf("program -> %s eof", symbol)
	return lalr1.NewParser(NewGrammar(grammar.ProductionLine{production, first()}))
}

type Parser struct {
	initialSymbol string
	filename      string
	reader        io.Reader
}

func NewParser() *Parser {
	return &Parser{"Anys", "", nil}
}

func (p *Parser) SetInitialSymbol(symbol string) {
	p.initialSymbol = symbol
}

func (p *Parser) Open(filename string, reader io.Reader) {
	p.filename = filename
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

	lexer := lexer.New()
	lexer.Open(np.reader)

	// TODO: Fix.
	channel := make(chan lalr1.Token, 10000)

	brackets := 0

	for lexer.Scan() {
		isSingleExpression := true
		isEmpty := true

		pos := ir.Pos{np.filename, lexer.LineNum(), lexer.LineNum(), lexer.Line()}

		for {
			text, ok := lexer.ShiftWord()
			if !ok {
				break
			}

			isEmpty = false

			switch text {
			case "{":
				isSingleExpression = false
				brackets++
			case "}":
				isSingleExpression = false
				brackets--
			}

			token := Token{pos, text}
			log.Printf("HERE %v", token)

			if tokenType, ok := parser.ParseTable().GetTokenType(text); ok {
				channel <- lalr1.Token{tokenType, token}
			} else {
				channel <- lalr1.Token{parser.ParseTable().TokenType("Token"), token}
			}
		}

		if brackets > 0 && !isEmpty && isSingleExpression {
			token := lalr1.Token{parser.ParseTable().TokenType(";"), Token{Pos: pos}}
			log.Printf("HERE %v", token)
			channel <- token
		}
	}

	{
		pos := ir.Pos{np.filename, lexer.LineNum(), lexer.LineNum(), lexer.Line()}
		token := lalr1.Token{parser.ParseTable().TokenType("eof"), Token{Pos: pos}}
		log.Printf("HERE %v", token)
		channel <- token
	}

	close(channel)

	if err := lexer.ScanErr(); err != nil {
		return t, err
	}

	ast, output, err := parser.Parse(channel)
	if err != nil {
		gotToken := output.Got.Data.(Token)

		fmt.Printf("%s:\n expected: %v\n got: %v\n",
			gotToken.Pos,
			output.Expected,
			output.GotSymbol.Name)

		return t, err
	}
	fmt.Printf("AST: %v\n", ast)

	return ast.(T), nil
}

func ParseFile(filename string, input io.Reader) ([]bplparser.Source, error) {
	parser := NewParser()
	parser.Open(filename, input)
	return Parse[[]bplparser.Source](parser)
}
