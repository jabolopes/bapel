package ir

import (
	"fmt"
	"strings"
)

type IrDeclCase int

const (
	TermDecl IrDeclCase = iota
	AliasDecl
	NameDecl
)

func (c IrDeclCase) String() string {
	switch c {
	case TermDecl:
		return "term declaration"
	case AliasDecl:
		return "alias declaration"
	case NameDecl:
		return "name declaration"
	default:
		panic(fmt.Errorf("unhandled IrDeclCase %d", c))
	}
}

type termDecl struct {
	ID   string
	Type IrType
}

func (d *termDecl) String() string {
	return fmt.Sprintf("%s: %s", d.ID, d.Type)
}

type aliasDecl struct {
	ID   string
	Kind IrKind
	Type IrType
}

func (d *aliasDecl) String() string {
	if d.Kind.Is(TypeKind) {
		return fmt.Sprintf("type %s = %s", d.ID, d.Type)
	}
	return fmt.Sprintf("type %s :: %s = %s", d.ID, d.Kind, d.Type)
}

type nameDecl struct {
	ID   string
	Kind IrKind
}

func (d *nameDecl) String() string {
	if d.Kind.Is(TypeKind) {
		return fmt.Sprintf("type %s", d.ID)
	}
	return fmt.Sprintf("type %s :: %s", d.ID, d.Kind)
}

type IrDecl struct {
	Case  IrDeclCase
	Term  *termDecl
	Alias *aliasDecl
	Name  *nameDecl

	// Whether this is an export.
	Export bool
	// Position in source file.
	Pos Pos
}

func (d IrDecl) String() string {
	if d.Case == 0 && d.Term == nil {
		return ""
	}

	var b strings.Builder
	if d.Export {
		b.WriteString("export ")
	}

	switch d.Case {
	case TermDecl:
		b.WriteString(d.Term.String())
	case AliasDecl:
		b.WriteString(d.Alias.String())
	case NameDecl:
		b.WriteString(d.Name.String())
	default:
		panic(fmt.Errorf("unhandled %T %d", d.Case, d.Case))
	}

	return b.String()
}

func (d IrDecl) Format(f fmt.State, verb rune) {
	if d.Case == 0 && d.Term == nil {
		return
	}

	if addMetadata := f.Flag('+'); addMetadata {
		d.Pos.Format(f, verb)
	}

	fmt.Fprint(f, d.String())
}

func (d IrDecl) Is(c IrDeclCase) bool {
	return d.Case == c
}

func (d IrDecl) ID() string {
	switch d.Case {
	case TermDecl:
		return d.Term.ID
	case AliasDecl:
		return d.Alias.ID
	case NameDecl:
		return d.Name.ID
	default:
		panic(fmt.Errorf("unhandled %T %d", d.Case, d.Case))
	}
}

func NewTermDecl(id string, typ IrType, export bool) IrDecl {
	return IrDecl{
		Case:   TermDecl,
		Term:   &termDecl{id, typ},
		Export: export,
	}
}

func NewAliasDecl(id string, kind IrKind, typ IrType, export bool) IrDecl {
	return IrDecl{
		Case:   AliasDecl,
		Alias:  &aliasDecl{id, kind, typ},
		Export: export,
	}
}

func NewNameDecl(id string, kind IrKind, export bool) IrDecl {
	return IrDecl{
		Case:   NameDecl,
		Name:   &nameDecl{id, kind},
		Export: export,
	}
}
