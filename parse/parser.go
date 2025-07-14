package parse

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/lex"
	"github.com/jabolopes/go-lalr1"
	"github.com/jabolopes/go-lalr1/grammar"
)

var (
	// Cached grammar to optimize parser construction.
	moduleGrammar = func() *lalr1.Parser {
		parser, err := lalr1.NewParser(NewGrammar(grammar.ProductionLine{"Program -> SourceFile EOF", first()}))
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
	if initialSymbol == "SourceFile" {
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
	return NewWithSymbol("SourceFile")
}

func Parse[T any](parser *Parser) (T, error) {
	var t T

	lex := lex.New(parser.reader)

	// TODO: Fix.
	channel := make(chan lalr1.Token, 10000)

	pos := ir.NewLinePos(parser.filename, 1)

	for {
		lexToken, ok := lex.NextToken()
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

	if err := lex.ScanErr(); err != nil {
		return t, err
	}

	parserLogger := log.New(io.Discard, "PARSER: ", 0)
	ast, output, err := parser.impl.Parse(channel, parserLogger)
	if err != nil {
		gotToken := output.Got.Data.(Token)
		return t, fmt.Errorf("%v %s:\n expected: %v\n got: %v",
			err,
			gotToken.Pos,
			output.Expected,
			output.GotSymbol.Name)
	}

	return ast.(T), nil
}

func ParseSourceFile(inputFilename string) (ast.SourceFile, error) {
	parser, err := New()
	if err != nil {
		return ast.SourceFile{}, err
	}

	inputFile, err := os.Open(inputFilename)
	if err != nil {
		return ast.SourceFile{}, err
	}
	defer inputFile.Close()

	parser.Open(inputFile.Name(), inputFile)

	sourceFile, err := Parse[ast.SourceFile](parser)
	if err != nil {
		return ast.SourceFile{}, err
	}

	sourceFile.Header.Filename = ir.NewFilename(inputFile.Name(), ir.Pos{})

	if validation := ast.ValidateSourceFile(&sourceFile); !validation.OK() {
		return ast.SourceFile{}, validation.Err()
	}

	return sourceFile, nil
}

func ParseWorkspace(inputFilename string) (ast.Workspace, error) {
	inputFile, err := os.Open(inputFilename)
	if err != nil {
		return ast.Workspace{}, err
	}
	defer inputFile.Close()

	parser, err := NewWithSymbol("Workspace")
	if err != nil {
		return ast.Workspace{}, err
	}

	parser.Open(inputFile.Name(), inputFile)

	workspace, err := Parse[ast.Workspace](parser)
	if err != nil {
		return ast.Workspace{}, err
	}

	if validation := ast.ValidateWorkspace(workspace); !validation.OK() {
		return ast.Workspace{}, validation.Err()
	}

	return workspace, nil
}
