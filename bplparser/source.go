package bplparser

import (
	"fmt"
	"strings"

	"github.com/jabolopes/bapel/ir"
)

type SourceCase int

const (
	SectionSource SourceCase = iota
	DeclSource
	EntitySource
	FunctionSource
	TermSource
	ElseSource
	EndSource
)

type Source struct {
	Case    SourceCase
	Section *struct {
		ID    string
		Decls []ir.IrDecl
	}
	Decl     *ir.IrDecl
	Entity   *ir.IrEntity
	Function *ir.IrFunction
	Term     *ir.IrTerm
}

func (s Source) String() string {
	if s.Case == 0 && s.Section == nil {
		return ""
	}

	switch s.Case {
	case SectionSource:
		var b strings.Builder
		b.WriteString(s.Section.ID)
		b.WriteString(" {\n")
		for _, decl := range s.Section.Decls {
			b.WriteString(decl.String())
			b.WriteString("\n")
		}
		b.WriteString("}\n")
		return b.String()

	case DeclSource:
		return s.Decl.String()

	case EntitySource:
		return s.Entity.ID

	case FunctionSource:
		var b strings.Builder
		b.WriteString(fmt.Sprintf("func %s", s.Function.ID))
		if len(s.Function.TypeVars) > 0 {
			b.WriteString("[")
			b.WriteString("'")
			b.WriteString(s.Function.TypeVars[0])
			for _, tvar := range s.Function.TypeVars[1:] {
				b.WriteString(", ")
				b.WriteString("'")
				b.WriteString(tvar)
			}
			b.WriteString("]")
		}
		b.WriteString("(")
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

	default:
		panic(fmt.Errorf("unhandled Source case %d", s.Case))
	}
}

func NewSectionSource(id string, decls []ir.IrDecl) Source {
	return Source{
		Case: SectionSource,
		Section: &struct {
			ID    string
			Decls []ir.IrDecl
		}{id, decls},
	}
}

func NewDeclSource(decl ir.IrDecl) Source {
	s := Source{}
	s.Case = DeclSource
	s.Decl = &decl
	return s
}

func NewEntitySource(entity ir.IrEntity) Source {
	s := Source{}
	s.Case = EntitySource
	s.Entity = &entity
	return s
}

func NewFunctionSource(function ir.IrFunction) Source {
	return Source{
		Case:     FunctionSource,
		Function: &function,
	}
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

func (p *Parser) parseAnyImpl() (Source, error) {
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

	if decl, err := p.parseDecl(false /* named */); err == nil {
		return NewDeclSource(decl), nil
	}

	return p.parseStatement()
}

func (p *Parser) parseAny() (result Source, err error) {
	p.withCheckpoint(func() error {
		result, err = p.parseAnyImpl()
		return err
	})
	return
}

func (p *Parser) ParseAny() (result Source, err error) {
	return p.parseAny()
}
