package comp

import (
	"fmt"
	"io"
	"os"

	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/ir"
)

func compileSource(compiler *ir.Compiler, source bplparser.Source) error {
	switch source.Case {
	case bplparser.SectionSource:
		return compiler.Section(source.Section.ID, source.Section.Decls)
	case bplparser.EntitySource:
		return compiler.Entity(*source.Entity)
	case bplparser.FunctionSource:
		return compiler.Function(*source.Function)
	case bplparser.TermSource:
		return compiler.Term(*source.Term)
	case bplparser.TypeDefSource:
		return compiler.TypeDefinition(source.TypeDef.Type)
	default:
		panic(fmt.Errorf("unhandled %T %d", source.Case, source.Case))
	}
}

func compileFile(compiler *ir.Compiler, input *os.File) error {
	sources, err := bplparser.ParseFile(input)
	if err != nil {
		return err
	}

	if err := compiler.Module(); err != nil {
		return err
	}

	for _, source := range sources {
		if err := compileSource(compiler, source); err != nil {
			return err
		}
	}

	return compiler.EndModule()
}

func CompileFile(inputFile *os.File, output io.Writer) error {
	compiler := ir.NewCompiler(output)
	return compileFile(compiler, inputFile)
}
