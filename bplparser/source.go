package bplparser

import (
	"fmt"
	"strings"

	"github.com/jabolopes/bapel/ir"
)

type SourceCase int

const (
	SectionSource = SourceCase(iota)
	DeclSource
	EntitySource
	FunctionSource
	TermSource
	ElseSource
	EndSource
	PrintSource
)

type Source struct {
	Case     SourceCase
	Section  string
	Decl     *ir.IrDecl
	Entity   string
	Function *struct {
		ID   string
		Args []ir.IrDecl
		Rets []ir.IrDecl
	}
	Term  *ir.IrTerm
	Print *struct {
		Args []string
	}
}

func (s Source) String() string {
	switch s.Case {
	case SectionSource:
		return s.Section

	case DeclSource:
		return s.Decl.String()

	case EntitySource:
		return s.Entity

	case FunctionSource:
		var b strings.Builder
		b.WriteString(fmt.Sprintf("func %s(", s.Function.ID))
		if len(s.Function.Args) > 0 {
			b.WriteString(s.Function.Args[0].String())
			for _, arg := range s.Function.Args[1:] {
				b.WriteString(", ")
				b.WriteString(arg.String())
			}
		}
		b.WriteString(") -> (")
		if len(s.Function.Rets) > 0 {
			b.WriteString(s.Function.Rets[0].String())
			for _, ret := range s.Function.Rets[1:] {
				b.WriteString(", ")
				b.WriteString(ret.String())
			}
		}
		b.WriteString(")")
		return b.String()

	case TermSource:
		return s.Term.String()

	case ElseSource:
		return "else"

	case EndSource:
		return "end"

	case PrintSource:
		var b strings.Builder
		for _, arg := range s.Print.Args[1:] {
			b.WriteString(" ")
			b.WriteString(arg)
		}
		return b.String()

	default:
		panic(fmt.Errorf("unhandled Source case %d", s.Case))
	}
}

func NewSectionSource(section string) Source {
	s := Source{}
	s.Case = SectionSource
	s.Section = section
	return s
}

func NewDeclSource(decl ir.IrDecl) Source {
	s := Source{}
	s.Case = DeclSource
	s.Decl = &decl
	return s
}

func NewEntitySource(id string) Source {
	s := Source{}
	s.Case = EntitySource
	s.Entity = id
	return s
}

func NewFunctionSource(id string, args, rets []ir.IrDecl) Source {
	s := Source{}
	s.Case = FunctionSource
	s.Function = &struct {
		ID   string
		Args []ir.IrDecl
		Rets []ir.IrDecl
	}{id, args, rets}
	return s
}

func NewTermSource(term ir.IrTerm) Source {
	s := Source{}
	s.Case = TermSource
	s.Term = &term
	return s
}

func NewElseSource() Source {
	s := Source{}
	s.Case = ElseSource
	return s
}

func NewEndSource() Source {
	s := Source{}
	s.Case = EndSource
	return s
}

func NewPrintSource(args []string) Source {
	s := Source{}
	s.Case = PrintSource
	s.Print = &struct {
		Args []string
	}{args}
	return s
}

func (p *Parser) parseAny() (Source, error) {
	if source, err := p.parseSection(); err == nil {
		return source, nil
	}

	if p.peek("func") {
		return p.parseFunc()
	}

	if p.peek("struct") {
		return p.parseStruct()
	}

	if p.peek("let") {
		return p.parseLet()
	}

	if p.peek("if") {
		return p.parseIf()
	}

	if p.peek("}") {
		if err := p.parseElse(); err == nil {
			return NewElseSource(), nil
		}

		if err := p.parseEnd(); err != nil {
			return Source{}, err
		}

		return NewEndSource(), nil
	}

	if p.peek("entity") {
		return p.parseEntity()
	}

	if p.peek("printU") {
		return p.parsePrintU()
	}

	if p.peek("printS") {
		return p.parsePrintS()
	}

	if source, err := p.parseDecl(false /* named */); err == nil {
		return source, nil
	}

	return p.parseStatement()
}

func (p *Parser) ParseAny() (result Source, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseAny()
		return err
	})
	return
}
