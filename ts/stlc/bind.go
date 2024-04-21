package stlc

import (
	"fmt"

	"github.com/jabolopes/bapel/ir"
)

type BindCase int

const (
	DeclBind BindCase = iota
)

type Bind struct {
	Case   BindCase
	Symbol Symbol
	Decl   *ir.IrDecl
}

func (b Bind) String() string {
	switch b.Case {
	case DeclBind:
		return b.Decl.String()
	default:
		panic(fmt.Errorf("unhandled BindCase %d", b.Case))
	}
}

func (b Bind) ID() (string, bool) {
	switch b.Case {
	case DeclBind:
		switch b.Decl.Case {
		case ir.TypeDecl:
			return b.Decl.Type().ID()
		case ir.TermDecl:
			return b.Decl.Term.ID, true
		default:
			panic(fmt.Errorf("unhandled %T %d", b.Decl.Case, b.Decl.Case))
		}
	default:
		panic(fmt.Errorf("unhandled %T %d", b.Case, b.Case))
	}
}

func NewDeclBind(symbol Symbol, decl ir.IrDecl) Bind {
	return Bind{
		Case:   DeclBind,
		Symbol: symbol,
		Decl:   &decl,
	}
}
