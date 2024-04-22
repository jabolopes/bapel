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
	binds []Bind

	list list.List[Bind]
}

func (c *Context) String() string {
	var b strings.Builder
	if len(c.binds) > 0 {
		b.WriteString(c.binds[0].String())
		for _, bind := range c.binds[1:] {
			b.WriteString(", ")
			b.WriteString(bind.String())
		}
	}

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

func (c *Context) StringNoImports() string {
	var b strings.Builder
	if len(c.binds) > 0 {
		if c.binds[0].Symbol == DefSymbol {
			b.WriteString(c.binds[0].String())
		}

		for _, bind := range c.binds[1:] {
			if bind.Symbol == DefSymbol {
				b.WriteString(", ")
				b.WriteString(bind.String())
			}
		}
	}

	if !c.list.Empty() {
		binds := c.list.Iterate().Collect()
		for _, bind := range binds {
			if bind.Symbol == DefSymbol {
				continue
			}

			b.WriteString(", ")
			b.WriteString(bind.String())
		}
	}
	return b.String()
}

func (c Context) ContainsVarType(tvar string) bool {
	for i := len(c.binds) - 1; i >= 0; i-- {
		bind := c.binds[i]
		if bind.Is(DeclBind) && bind.Decl.Type().Is(ir.VarType) && bind.Decl.Type().Var == tvar {
			return true
		}
	}

	for it := c.list.Iterate(); ; {
		bind, ok := it.Next()
		if !ok {
			break
		}

		if bind.Is(DeclBind) && bind.Decl.Type().Is(ir.VarType) && bind.Decl.Type().Var == tvar {
			return true
		}
	}

	return false
}

func (c Context) Pop() (Bind, Context) {
	if len(c.binds) > 0 {
		bind := c.binds[len(c.binds)-1]
		c.binds = c.binds[:len(c.binds)-1]
		return bind, c
	}

	bind, ok := c.list.Value()
	if !ok {
		panic("Context is empty")
	}

	c.list = c.list.Remove()
	return bind, c
}

func (c Context) Copy() Context {
	c.binds = append([]Bind{}, c.binds...)
	return c
}

// TODO: Merge with LookupBind().
func (c *Context) lookupBind(id string, findCase FindCase) (Bind, bool) {
	for i := len(c.binds) - 1; i >= 0; i-- {
		bind := c.binds[i]
		if bindID, ok := bind.ID(); !ok || bindID != id {
			continue
		}

		switch {
		case findCase == FindDeclOnly && bind.Symbol == DefSymbol:
			continue
		case findCase == FindDefOnly && bind.Symbol != DefSymbol:
			continue
		}

		return bind, true
	}

	for it := c.list.Iterate(); ; {
		bind, ok := it.Next()
		if !ok {
			break
		}

		if bindID, ok := bind.ID(); !ok || bindID != id {
			continue
		}

		switch {
		case findCase == FindDeclOnly && bind.Symbol == DefSymbol:
			continue
		case findCase == FindDefOnly && bind.Symbol != DefSymbol:
			continue
		}

		return bind, true
	}

	return Bind{}, false
}

func (c *Context) LookupBind(id string, findCase FindCase) (Bind, bool) {
	return c.lookupBind(id, findCase)
}

func (c *Context) getBind(id string, findCase FindCase) (Bind, error) {
	bind, ok := c.lookupBind(id, findCase)
	if !ok {
		return Bind{}, fmt.Errorf("%q is undefined", id)
	}

	return bind, nil
}

func (c *Context) lookupType(typ ir.IrType) (Bind, bool) {
	for i := len(c.binds) - 1; i >= 0; i-- {
		bind := c.binds[i]
		if bind.Case != DeclBind || bind.Decl.Case != ir.TypeDecl {
			continue
		}

		if ir.EqualsType(bind.Decl.Type(), typ) {
			return bind, true
		}
	}

	for it := c.list.Iterate(); ; {
		bind, ok := it.Next()
		if !ok {
			break
		}

		if bind.Case != DeclBind || bind.Decl.Case != ir.TypeDecl {
			continue
		}

		if ir.EqualsType(bind.Decl.Type(), typ) {
			return bind, true
		}
	}

	return Bind{}, false
}

func (c *Context) getType(typ ir.IrType) (Bind, error) {
	bind, ok := c.lookupType(typ)
	if !ok {
		return Bind{}, fmt.Errorf("type %s is undefined", typ)
	}

	return bind, nil
}

func (c *Context) resolveTypeName(typ ir.IrType) (ir.IrType, error) {
	switch typ.Case {
	case ir.AliasType:
		return c.resolveTypeName(typ.Alias.Value)

	case ir.NameType:
		bind, ok := c.lookupBind(typ.Name, FindAny)
		if !ok {
			return ir.IrType{}, fmt.Errorf("%q is undefined", typ.Name)
		}

		if ir.EqualsType(typ, bind.Decl.Type()) {
			return typ, nil
		}

		return c.resolveTypeName(bind.Decl.Type())

	default:
		// TODO: Should probably check if the type is defined in the context.
		return typ, nil
	}
}

func (c *Context) AddBind(bind Bind) error {
	bindID, ok := bind.ID()
	if ok {
		if _, ok := c.lookupBind(bindID, FindDefOnly); ok {
			return fmt.Errorf("%q is already defined", bindID)
		}

		if ok && bind.Symbol == DefSymbol {
			// Check that definition (e.g., function, struct, etc) matches declaration (if any).
			declaration, ok := c.lookupBind(bindID, FindDeclOnly)
			if ok {
				if err := NewTypechecker(c).subtype(declaration.Decl.Type(), bind.Decl.Type()); err != nil {
					return err
				}
			}
		}
	}

	// c.binds = append(c.binds, bind)
	c.list = c.list.Add(bind)

	return IsWellformedContext(*c)
}

func (c Context) enterFunction(id string, typeVars []string, args, rets []ir.IrDecl) Context {
	for _, tvar := range typeVars {
		if err := c.AddBind(NewDeclBind(DefSymbol, ir.NewTypeDecl(ir.NewVarType(tvar)))); err != nil {
			// TODO: Remove panic.
			panic(err)
		}
	}

	for _, arg := range args {
		if err := c.AddBind(NewDeclBind(DefSymbol, arg)); err != nil {
			// TODO: Remove panic.
			panic(err)
		}
	}

	for _, ret := range rets {
		if err := c.AddBind(NewDeclBind(DefSymbol, ret)); err != nil {
			// TODO: Remove panic.
			panic(err)
		}
	}

	return c
}

func (c *Context) IsExport(id string) bool {
	symbol, ok := c.lookupBind(id, FindDeclOnly)
	return ok && symbol.Symbol == ExportSymbol
}

func (c *Context) CheckModule() error {
	{
		// Check all exports and all declarations have a definition (i.e., there are
		// no undefined exports or declarations).
		exported := map[string]struct{}{}
		declared := map[string]struct{}{}
		for _, bind := range c.binds {
			bindID, ok := bind.ID()
			if !ok {
				continue
			}

			switch bind.Symbol {
			case ExportSymbol:
				exported[bindID] = struct{}{}
			case DeclSymbol:
				declared[bindID] = struct{}{}
			}
		}

		for _, bind := range c.binds {
			bindID, ok := bind.ID()
			if !ok {
				continue
			}

			if bind.Symbol == DefSymbol {
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

	{
		// Check all exports and all declarations have a definition (i.e., there are
		// no undefined exports or declarations).
		exported := map[string]struct{}{}
		declared := map[string]struct{}{}
		for it := c.list.Iterate(); ; {
			bind, ok := it.Next()
			if !ok {
				break
			}

			bindID, ok := bind.ID()
			if !ok {
				continue
			}

			switch bind.Symbol {
			case ExportSymbol:
				exported[bindID] = struct{}{}
			case DeclSymbol:
				declared[bindID] = struct{}{}
			}
		}

		for _, bind := range c.binds {
			bindID, ok := bind.ID()
			if !ok {
				continue
			}

			if bind.Symbol == DefSymbol {
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
}

func NewContext() *Context {
	return &Context{
		[]Bind{}, /* binds */
		list.New[Bind](),
	}
}
