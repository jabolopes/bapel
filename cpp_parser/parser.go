package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/jabolopes/bapel/cpp_parser/parser"
)

func ParseSymbol(symbol string, filename string, reader io.Reader) (any, error) {
	inputBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	inputStream := antlr.NewInputStream(string(inputBytes))
	lexer := parser.NewbapelLexer(inputStream)
	tokenStream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewbapelParser(tokenStream)

	p.BuildParseTrees = true

	errListener := parser.NewCustomErrorListener(filename)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errListener)
	p.RemoveErrorListeners()
	p.AddErrorListener(errListener)

	var tree antlr.ParseTree
	switch symbol {
	case "SourceFile":
		tree = p.SourceFile()
	case "Workspace":
		tree = p.Workspace()
	case "Decl":
		tree = p.Decl()
	default:
		return nil, fmt.Errorf("unsupported symbol %q", symbol)
	}

	if len(errListener.Errors()) > 0 {
		return nil, errors.New(errListener.Errors()[0])
	}

	visitor := parser.NewASTBuilder(filename)
	ast := tree.Accept(visitor)
	return ast, nil
}
