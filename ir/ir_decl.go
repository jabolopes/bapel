package ir

import (
	"cmp"
	"fmt"
	"strings"
)

type IrDeclCase int

const (
	TermDecl IrDeclCase = iota
	AliasDecl
	NameDecl
	TraitDecl
)

func (c IrDeclCase) String() string {
	switch c {
	case TermDecl:
		return "term declaration"
	case AliasDecl:
		return "alias declaration"
	case NameDecl:
		return "name declaration"
	case TraitDecl:
		return "trait declaration"
	default:
		panic(fmt.Errorf("unhandled %T %d", c, c))
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

type IrSignature struct {
	ID      string
	Args    []FunctionArg
	RetType IrType
}

func (s IrSignature) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("fn %s(", s.ID))
	Interleave(s.Args, func() { b.WriteString(", ") }, func(_ int, arg FunctionArg) {
		b.WriteString(arg.String())
	})
	b.WriteString(fmt.Sprintf(") -> %s", s.RetType))
	return b.String()
}

type traitDecl struct {
	ID         string
	TypeParams []VarKind
	Methods    []IrSignature
}

func (d *traitDecl) String() string {
	var b strings.Builder
	b.WriteString("trait ")
	b.WriteString(d.ID)
	if len(d.TypeParams) > 0 {
		b.WriteString(" [")
		Interleave(d.TypeParams, func() { b.WriteString(", ") }, func(_ int, tv VarKind) {
			b.WriteString(fmt.Sprintf("'%s", tv.Var))
		})
		b.WriteString("]")
	}
	b.WriteString(" {\n")
	for _, m := range d.Methods {
		b.WriteString(fmt.Sprintf("  %s\n", m))
	}
	b.WriteString("}")
	return b.String()
}

type IrDecl struct {
	Case  IrDeclCase
	Term  *termDecl
	Alias *aliasDecl
	Name  *nameDecl
	Trait *traitDecl

	// Whether this is an export.
	Export bool
	// Position in source file.
	Pos Pos
}

func (d IrDecl) String() string {
	if d.Case == 0 && d.Term == nil && d.Trait == nil {
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
	case TraitDecl:
		b.WriteString(d.Trait.String())
	default:
		panic(fmt.Errorf("unhandled %T %d", d.Case, d.Case))
	}

	return b.String()
}

func (d IrDecl) Format(f fmt.State, verb rune) {
	if d.Case == 0 && d.Term == nil && d.Trait == nil {
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
	case TraitDecl:
		return d.Trait.ID
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

func NewTraitDecl(id string, typeParams []VarKind, methods []IrSignature, export bool) IrDecl {
	return IrDecl{
		Case:   TraitDecl,
		Trait:  &traitDecl{id, typeParams, methods},
		Export: export,
	}
}

func CompareDecl(d1, d2 IrDecl) int {
	if c := cmp.Compare(d1.Case, d2.Case); c != 0 {
		// Sort name decl before alias decl before term decl.
		return -c
	}

	switch d1.Case {
	case TermDecl:
		return cmp.Compare(d1.Term.ID, d2.Term.ID)
	case AliasDecl:
		return cmp.Compare(d1.Alias.ID, d2.Alias.ID)
	case NameDecl:
		return cmp.Compare(d1.Name.ID, d2.Name.ID)
	case TraitDecl:
		return cmp.Compare(d1.Trait.ID, d2.Trait.ID)
	default:
		return 0
	}
}

func (d IrDecl) Clone() IrDecl {
	dCopy := d
	if d.Term != nil {
		dCopy.Term = &termDecl{ID: d.Term.ID, Type: d.Term.Type}
	}
	if d.Alias != nil {
		dCopy.Alias = &aliasDecl{ID: d.Alias.ID, Kind: d.Alias.Kind, Type: d.Alias.Type}
	}
	if d.Name != nil {
		dCopy.Name = &nameDecl{ID: d.Name.ID, Kind: d.Name.Kind}
	}
	if d.Trait != nil {
		typeParams := make([]VarKind, len(d.Trait.TypeParams))
		copy(typeParams, d.Trait.TypeParams)
		methods := make([]IrSignature, len(d.Trait.Methods))
		for i, m := range d.Trait.Methods {
			args := make([]FunctionArg, len(m.Args))
			copy(args, m.Args)
			methods[i] = IrSignature{ID: m.ID, Args: args, RetType: m.RetType}
		}
		dCopy.Trait = &traitDecl{ID: d.Trait.ID, TypeParams: typeParams, Methods: methods}
	}
	return dCopy
}


