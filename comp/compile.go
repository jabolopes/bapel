package comp

import (
	"fmt"
	"io"
	"os"

	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/ir"
)

type Context struct {
	parser   *bplparser.Parser
	compiler *ir.Compiler
}

func compileSource(context *Context, source bplparser.Source) error {
	switch source.Case {
	case bplparser.SectionSource:
		return context.compiler.Section(source.Section.ID, source.Section.Decls)
	case bplparser.EntitySource:
		return context.compiler.Entity(*source.Entity)
	case bplparser.FunctionSource:
		return context.compiler.Function(*source.Function)
	case bplparser.TermSource:
		return context.compiler.Term(*source.Term)
	case bplparser.TypeDefSource:
		return context.compiler.TypeDefinition(source.TypeDef.Type)
	default:
		panic(fmt.Errorf("unhandled %T %d", source.Case, source.Case))
	}
}

func compileFile(context *Context, input *os.File) error {
	sources, err := bplparser.ParseFile(input)
	if err != nil {
		return err
	}

	if err := context.compiler.Module(); err != nil {
		return err
	}

	for _, source := range sources {
		if err := compileSource(context, source); err != nil {
			return err
		}
	}

	return context.compiler.End()
}

func CompileFile(inputFile *os.File, output io.Writer) error {
	compiler := ir.NewCompiler(output)
	context := &Context{bplparser.NewParser(), compiler}
	return compileFile(context, inputFile)
}
