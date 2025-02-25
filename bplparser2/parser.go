package bplparser2

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/lexer"
	"github.com/jabolopes/go-lalr1"
	"github.com/jabolopes/go-lalr1/grammar"
)

var (
	// Cached grammar to optimize parser construction.
	moduleGrammar = func() *lalr1.Parser {
		parser, err := lalr1.NewParser(NewGrammar(grammar.ProductionLine{"Program -> Module EOF", first()}))
		if err != nil {
			panic(err)
		}
		return parser
	}()
)

type Token struct {
	Pos  ir.Pos
	Text string
}

type Integer struct {
	Pos   ir.Pos
	Value int
}

type Parser struct {
	impl     *lalr1.Parser
	filename string
	reader   io.Reader
}

func newParserImpl(initialSymbol string) (*lalr1.Parser, error) {
	var impl *lalr1.Parser
	if initialSymbol == "Module" {
		impl = moduleGrammar
	} else {
		grammar := NewGrammar(grammar.ProductionLine{fmt.Sprintf("Program -> %s EOF", initialSymbol), first()})
		var err error
		impl, err = lalr1.NewParser(grammar)
		if err != nil {
			return nil, err
		}
	}

	if conflicts := impl.Machine().Conflicts(); len(conflicts) > 0 {
		// Return an error if there are conflicts.
		var str strings.Builder

		str.WriteString(impl.Machine().String())
		str.WriteString("\n")
		str.WriteString(impl.ParseTable().String())
		str.WriteString("\n")

		str.WriteString("Grammar has conflicts:\n")
		for _, conflict := range conflicts {
			str.WriteString(fmt.Sprintf("  * %s\n", conflict))
		}

		return nil, errors.New(str.String())
	}

	return impl, nil
}

func (p *Parser) Open(filename string, reader io.Reader) {
	p.filename = filename
	p.reader = reader
}

func NewWithSymbol(initialSymbol string) (*Parser, error) {
	impl, err := newParserImpl(initialSymbol)
	if err != nil {
		return nil, err
	}

	return &Parser{impl, "", bytes.NewReader(nil)}, nil
}

func New() (*Parser, error) {
	return NewWithSymbol("Module")
}

func Parse[T any](parser *Parser) (T, error) {
	var t T

	lexer := lexer.New(parser.reader)

	// TODO: Fix.
	channel := make(chan lalr1.Token, 10000)

	pos := ir.Pos{parser.filename, 1, 1, ""}

	for {
		lexToken, ok := lexer.NextToken()
		if !ok {
			break
		}

		pos.BeginLineNum = lexToken.LineNum
		pos.EndLineNum = lexToken.LineNum

		token := Token{pos, lexToken.Value}
		if tokenType, ok := parser.impl.ParseTable().GetTokenType(lexToken.Value); ok {
			channel <- lalr1.Token{tokenType, token}
		} else {
			channel <- lalr1.Token{parser.impl.ParseTable().TokenType("Token"), token}
		}
	}

	{
		token := lalr1.Token{parser.impl.ParseTable().TokenType("EOF"), Token{Pos: pos}}
		channel <- token
	}

	close(channel)

	if err := lexer.ScanErr(); err != nil {
		return t, err
	}

	parserLogger := log.New(io.Discard, "PARSER: ", 0)
	ast, output, err := parser.impl.Parse(channel, parserLogger)
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

func ParseWith(parser *Parser, filename string, input io.Reader) (ast.Module, error) {
	parser.Open(filename, input)

	module, err := Parse[ast.Module](parser)
	if err != nil {
		return ast.Module{}, err
	}

	module.Header.Name = TrimExtension(filename)
	ast.ValidateModule(&module)

	return module, nil
}

func ParseFile(filename string, input io.Reader) (ast.Module, error) {
	parser, err := New()
	if err != nil {
		return ast.Module{}, err
	}

	return ParseWith(parser, filename, input)
}
