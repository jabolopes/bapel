package ir

import (
	"fmt"
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
	return fmt.Sprintf("%s : %s", d.ID, d.Type)
}

type aliasDecl struct {
	ID   string
	Type IrType
}

func (d *aliasDecl) String() string {
	return fmt.Sprintf("type %s = %s", d.ID, d.Type)
}

type nameDecl struct {
	ID string
}

func (d *nameDecl) String() string {
	return fmt.Sprintf("type %s", d.ID)
}

type IrDecl struct {
	Case  IrDeclCase
	Term  *termDecl
	Alias *aliasDecl
	Name  *nameDecl

	// Position in source file.
	Pos Pos
}

func (d IrDecl) String() string {
	if d.Case == 0 && d.Term == nil {
		return ""
	}

	switch d.Case {
	case TermDecl:
		return d.Term.String()
	case AliasDecl:
		return d.Alias.String()
	case NameDecl:
		return d.Name.String()
	default:
		panic(fmt.Errorf("unhandled %T %d", d.Case, d.Case))
	}
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

func NewTermDecl(id string, typ IrType) IrDecl {
	return IrDecl{
		Case: TermDecl,
		Term: &termDecl{id, typ},
	}
}

func NewAliasDecl(id string, typ IrType) IrDecl {
	return IrDecl{
		Case:  AliasDecl,
		Alias: &aliasDecl{id, typ},
	}
}

func NewNameDecl(id string) IrDecl {
	return IrDecl{
		Case: NameDecl,
		Name: &nameDecl{id},
	}
}
