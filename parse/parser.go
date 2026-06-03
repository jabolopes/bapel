package parse

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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

func FindLexerBin() (string, error) {
	// Try relative to CWD first
	path := "bootstrap/lexer"
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	// Walk up to find it (useful for tests running in subdirectories)
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		path = filepath.Join(cwd, "bootstrap/lexer")
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
		parent := filepath.Dir(cwd)
		if parent == cwd {
			break // reached root
		}
		cwd = parent
	}

	return "", fmt.Errorf("could not find bootstrap/lexer binary")
}

func (p *Parser) readAllTokens() ([]lalr1.Token, error) {
	file, err := io.ReadAll(p.reader)
	if err != nil {
		return nil, err
	}

	lexerBin, err := FindLexerBin()
	if err != nil {
		return nil, err
	}

	// Spawn C++ lexer with --raw and --filename flags
	cmd := exec.Command(lexerBin, "--raw", "--filename", p.filename)
	cmd.Stdin = bytes.NewReader(file)
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start C++ lexer: %v", err)
	}

	tokens := []lalr1.Token{}
	pos := ir.NewLinePos(p.filename, 1)

	reader := bufio.NewReader(stdout)
	for {
		// 1. Read header line (line type size)
		header, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to read token header: %v", err)
		}

		header = strings.TrimSuffix(header, "\n")
		if len(header) == 0 {
			continue
		}

		var line, tokType, size int
		_, err = fmt.Sscanf(header, "%d %d %d", &line, &tokType, &size)
		if err != nil {
			return nil, fmt.Errorf("failed to parse token header %q: %v", header, err)
		}

		// 2. Read exactly 'size' bytes of the raw token value
		valueBuf := make([]byte, size)
		_, err = io.ReadFull(reader, valueBuf)
		if err != nil {
			return nil, fmt.Errorf("failed to read token value of size %d: %v", size, err)
		}
		value := string(valueBuf)

		pos.BeginLineNum = line
		pos.EndLineNum = line

		token := Token{pos, value}
		if tokType == int(lex.NumberToken) {
			tokens = append(tokens, lalr1.Token{p.impl.ParseTable().TokenType("NumberToken"), token})
		} else if tokType == int(lex.RuneToken) {
			tokens = append(tokens, lalr1.Token{p.impl.ParseTable().TokenType("RuneToken"), token})
		} else if tokType == int(lex.StringToken) {
			tokens = append(tokens, lalr1.Token{p.impl.ParseTable().TokenType("StringToken"), token})
		} else if tokenType, ok := p.impl.ParseTable().GetTokenType(value); ok {
			tokens = append(tokens, lalr1.Token{tokenType, token})
		} else {
			tokens = append(tokens, lalr1.Token{p.impl.ParseTable().TokenType("Token"), token})
		}
	}

	if err := cmd.Wait(); err != nil {
		if stderr.Len() > 0 {
			return nil, errors.New(strings.TrimSpace(stderr.String()))
		}
		return nil, fmt.Errorf("c++ lexer failed: %v", err)
	}

	// Append EOF
	tokens = append(tokens, lalr1.Token{p.impl.ParseTable().TokenType("EOF"), Token{Pos: pos}})
	return tokens, nil
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

	tokens, err := parser.readAllTokens()
	if err != nil {
		return t, err
	}

	parserLogger := log.New(io.Discard, "PARSER: ", 0)
	ast, output, err := parser.impl.Parse(tokens, parserLogger)
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
