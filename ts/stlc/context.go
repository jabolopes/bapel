package stlc

import (
	"fmt"
	"strings"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/ts/list"
)

type FindCase int

const (
	FindAny = FindCase(iota)
	FindDeclOnly
	FindDefOnly
)

type Context struct {
	list list.List[Bind]
}

func (c Context) String() string {
	var b strings.Builder
	if !c.list.Empty() {
		binds := c.list.Iterate().Collect()
		b.WriteString(binds[0].String())
		for _, bind := range binds[1:] {
			b.WriteString(", ")
			b.WriteString(bind.String())
		}
	}
	return b.String()
}

func (c Context) StringNoImports() string {
	var b strings.Builder
	if !c.list.Empty() {
		binds := c.list.Iterate().Collect()
		for _, bind := range binds {
			if symbol, ok := bind.Symbol(); !ok || symbol != DefSymbol {
				continue
			}

			b.WriteString(", ")
			b.WriteString(bind.String())
		}
	}
	return b.String()
}

func (c Context) GetAliasBind(name string) (Bind, error) {
	for it := c.list.Iterate(); ; {
		_, bind, ok := it.Next()
		if !ok {
			break
		}

		if bind.Is(AliasBind) && bind.Alias.Name == name {
			return bind, nil
		}
	}

	return Bind{}, fmt.Errorf("%q is undefined", name)
}

func (c Context) GetTermBind(name string) (Bind, error) {
	for it := c.list.Iterate(); ; {
		_, bind, ok := it.Next()
		if !ok {
			break
		}

		if bind.Is(TermBind) && bind.Term.Name == name {
			return bind, nil
		}
	}

	return Bind{}, fmt.Errorf("%q is undefined", name)
}

func (c Context) ContainsAliasBind(name string) bool {
	_, err := c.GetAliasBind(name)
	return err == nil
}

func (c Context) ContainsComponentBind(elemType ir.IrType) bool {
	for it := c.list.Iterate(); ; {
		_, bind, ok := it.Next()
		if !ok {
			break
		}

		if bind.Is(ComponentBind) && ir.EqualsType(bind.Component.ElemType, elemType) {
			return true
		}
	}

	return false
}

func (c Context) ContainsNameBind(name string) bool {
	for it := c.list.Iterate(); ; {
		_, bind, ok := it.Next()
		if !ok {
			break
		}

		if bind.Is(NameBind) && bind.Name.Name == name {
			return true
		}
	}

	return false
}

func (c Context) ContainsVarType(tvar string) bool {
	for it := c.list.Iterate(); ; {
		_, bind, ok := it.Next()
		if !ok {
			break
		}

		if bind.Is(TypeVarBind) && bind.TypeVar.Name == tvar {
			return true
		}
	}

	return false
}

func (c Context) Empty() bool {
	return c.list.Empty()
}

func (c Context) Pop() (Bind, Context) {
	bind, ok := c.list.Value()
	if !ok {
		panic("Context is empty")
	}

	c.list = c.list.Remove()
	return bind, c
}

func (c Context) Copy() Context {
	return c
}

func (c Context) LookupBind(id string, findCase FindCase) (Bind, bool) {
	for it := c.list.Iterate(); ; {
		_, bind, ok := it.Next()
		if !ok {
			break
		}

		if bindID, ok := bind.ID(); !ok || bindID != id {
			continue
		}

		symbol, symbolOk := bind.Symbol()

		switch {
		case findCase == FindDeclOnly && symbolOk && symbol != DefSymbol:
			return bind, true
		case findCase == FindDefOnly && !symbolOk:
			return bind, true
		case findCase == FindDefOnly && symbolOk && symbol == DefSymbol:
			return bind, true
		}
	}

	return Bind{}, false
}

func (c Context) resolveTypeName(typ ir.IrType) (ir.IrType, error) {
	switch typ.Case {
	case ir.NameType:
		if bind, err := c.GetAliasBind(typ.Name); err == nil {
			return c.resolveTypeName(bind.Alias.Type)
		}
		return typ, nil
	default:
		return typ, nil
	}
}

func (c Context) AddBind(bind Bind) (Context, error) {
	// TODO: Finish.

	// bindID, ok := bind.ID()
	// if ok {
	// 	if _, ok := c.LookupBind(bindID, FindDefOnly); ok {
	// 		return c, fmt.Errorf("%q is already defined", bindID)
	// 	}

	// 	if ok && bind.Symbol == DefSymbol {
	// 		// Check that definition (e.g., function, struct, etc) matches declaration (if any).
	// 		declaration, ok := c.LookupBind(bindID, FindDeclOnly)
	// 		if ok {
	// 			if err := NewTypechecker(c).subtype(declaration.Decl.Type(), bind.Decl.Type()); err != nil {
	// 				return c, err
	// 			}
	// 		}
	// 	}
	// }

	c.list = c.list.Add(bind)

	if err := IsWellformedContext(c); err != nil {
		return c, err
	}

	return c, nil
}

func (c Context) enterFunction(id string, typeVars []string, args, rets []ir.IrDecl) (Context, error) {
	for _, tvar := range typeVars {
		var err error
		if c, err = c.AddBind(NewTypeVarBind(tvar)); err != nil {
			return c, err
		}
	}

	for _, arg := range args {
		var err error
		if c, err = c.AddBind(NewTermBind(arg.Term.ID, arg.Term.Type, DefSymbol)); err != nil {
			return c, err
		}
	}

	for _, ret := range rets {
		var err error
		if c, err = c.AddBind(NewTermBind(ret.Term.ID, ret.Term.Type, DefSymbol)); err != nil {
			return c, err
		}
	}

	return c, nil
}

func (c Context) IsExport(id string) bool {
	bind, ok := c.LookupBind(id, FindDeclOnly)
	if !ok {
		return false
	}

	symbol, ok := bind.Symbol()
	if !ok {
		return false
	}

	return symbol == ExportSymbol
}

func (c Context) CheckModule() error {
	// Check all exports and all declarations have a definition (i.e., there are
	// no undefined exports or declarations).
	exported := map[string]struct{}{}
	declared := map[string]struct{}{}
	for it := c.list.Iterate(); ; {
		_, bind, ok := it.Next()
		if !ok {
			break
		}

		bindID, ok := bind.ID()
		if !ok {
			continue
		}

		symbol, symbolOK := bind.Symbol()
		switch {
		case symbolOK && symbol == ExportSymbol:
			exported[bindID] = struct{}{}
		case symbolOK && symbol == DeclSymbol:
			declared[bindID] = struct{}{}
		}
	}

	for it := c.list.Iterate(); ; {
		_, bind, ok := it.Next()
		if !ok {
			break
		}

		bindID, ok := bind.ID()
		if !ok {
			continue
		}

		if symbol, symbolOK := bind.Symbol(); !symbolOK || symbol == DefSymbol {
			delete(exported, bindID)
			delete(declared, bindID)
		}
	}

	if len(exported) > 0 {
		return fmt.Errorf("symbols %v are exported but not defined", exported)
	}

	if len(declared) > 0 {
		return fmt.Errorf("symbols %v are declared but not defined", declared)
	}

	return nil
}

func NewContext() Context {
	return Context{
		list.New[Bind](),
	}
}
