package comp

import (
	"fmt"
	"strings"

	"github.com/jabolopes/bapel/ast"
	"github.com/jabolopes/bapel/ir"
)

// TODO: Finish String() methods of all types.

type OriginCase int

const (
	ImportOrigin OriginCase = iota
	ImplOrigin
	ImplicitOrigin
	ExplicitUndefinedOrigin
	ExplicitDefinedOrigin
)

func (c OriginCase) String() string {
	switch c {
	case ImportOrigin:
		return "import"
	case ImplOrigin:
		return "impl"
	case ImplicitOrigin:
		return "implicit"
	case ExplicitUndefinedOrigin:
		return "explicit undefined"
	case ExplicitDefinedOrigin:
		return "explicit defined"
	default:
		panic(fmt.Errorf("unhandled %T %d", c, c))
	}
}

type importOrigin struct {
	ModuleID ast.ID // e.g., 'core'
}

func (t *importOrigin) String() string {
	return fmt.Sprintf("%+s", t.ModuleID)
}

type implOrigin struct {
	ModuleFilename ast.ID // e.g., 'core_impl.bpl' or 'core_impl.cc'
}

func (t *implOrigin) String() string {
	return fmt.Sprintf("%+s", t.ModuleFilename)
}

type implicitOrigin struct {
	Definition ir.IrDecl // IrDecl of the definition that must match the Symbol's declaration.
}

func (t *implicitOrigin) String() string {
	return t.Definition.String()
}

type explicitUndefinedOrigin struct {
}

func (t *explicitUndefinedOrigin) String() string {
	return ""
}

type explicitDefinedOrigin struct {
	Definition ir.IrDecl // IrDecl of the definition that must match the Symbol's declaration.
}

func (t *explicitDefinedOrigin) String() string {
	return t.Definition.String()
}

type Origin struct {
	Case              OriginCase
	Import            *importOrigin
	Impl              *implOrigin
	Implicit          *implicitOrigin
	ExplicitUndefined *explicitUndefinedOrigin
	ExplicitDefined   *explicitDefinedOrigin
}

func (t Origin) String() string {
	if t.Case == 0 && t.Import == nil {
		return ""
	}

	switch t.Case {
	case ImportOrigin:
		return t.Import.String()
	case ImplOrigin:
		return t.Impl.String()
	case ImplicitOrigin:
		return t.Import.String()
	case ExplicitUndefinedOrigin:
		return t.ExplicitUndefined.String()
	case ExplicitDefinedOrigin:
		return t.ExplicitDefined.String()
	default:
		panic(fmt.Errorf("unhandled %T %d", t.Case, t.Case))
	}
}

type Symbol struct {
	Decl     ir.IrDecl
	IsExport bool
	Origin   Origin
}

func (s Symbol) String() string {
	var b strings.Builder
	if s.IsExport {
		b.WriteString("export ")
	}
	b.WriteString(s.Decl.String())
	b.WriteString(s.Origin.String())
	return b.String()
}

func NewImportSymbol(moduleID ast.ID, decl ir.IrDecl) Symbol {
	return Symbol{
		decl,
		false,
		Origin{
			Case:   ImportOrigin,
			Import: &importOrigin{moduleID},
		},
	}
}

func NewImplSymbol(moduleFilename ast.ID, decl ir.IrDecl) Symbol {
	return Symbol{
		decl,
		false,
		Origin{
			Case: ImplOrigin,
			Impl: &implOrigin{moduleFilename},
		},
	}
}

func NewImplicitSymbol(definitionDecl ir.IrDecl, decl ir.IrDecl) Symbol {
	return Symbol{
		decl,
		false,
		Origin{
			Case:     ImplicitOrigin,
			Implicit: &implicitOrigin{definitionDecl},
		},
	}
}

func NewExplicitUndefinedSymbol(decl ir.IrDecl) Symbol {
	return Symbol{
		decl,
		false,
		Origin{
			Case:              ExplicitUndefinedOrigin,
			ExplicitUndefined: &explicitUndefinedOrigin{},
		},
	}
}

func NewExplicitDefinedSymbol(definitionDecl ir.IrDecl, decl ir.IrDecl) Symbol {
	return Symbol{
		decl,
		false,
		Origin{
			Case:            ExplicitDefinedOrigin,
			ExplicitDefined: &explicitDefinedOrigin{definitionDecl},
		},
	}
}

type SymbolTable struct {
	// Symbols keyed by IrDecl.ID().
	Table map[string]Symbol
}

// TODO: Finish.
func (t *SymbolTable) Add(symbol Symbol) error {
	t.Table[symbol.Decl.ID()] = symbol
	return nil
}

// TODO: Finish.
func (t *SymbolTable) Export(decl ir.IrDecl) error {
	val := t.Table[decl.ID()]
	val.IsExport = true
	t.Table[decl.ID()] = val
	return nil
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{map[string]Symbol{}}
}
