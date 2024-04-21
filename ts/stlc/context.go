package stlc

import (
	"fmt"
	"strings"

	"github.com/jabolopes/bapel/ir"
)

type FindCase int

const (
	FindAny = FindCase(iota)
	FindDeclOnly
	FindDefOnly
)

type Context struct {
	binds []Bind
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
	return b.String()
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

	c.binds = append(c.binds, bind)
	return nil
}

func (c *Context) addMarker(id string) {
	// TODO: Call AddBind instead and return propagate error.
	c.binds = append(c.binds, NewMarkerBind(id))
}

func (c *Context) removeTillMarker(id string) {
	for {
		// TODO: Check bounds and return an error.
		bind := c.binds[len(c.binds)-1]
		c.binds = c.binds[:len(c.binds)-1]

		if bind.Case == MarkerBind && *bind.Marker == id {
			return
		}
	}
}

func (c *Context) enterFunction(id string, typeVars []string, args, rets []ir.IrDecl) {
	c.addMarker(id)

	for _, tvar := range typeVars {
		c.binds = append(c.binds, NewDeclBind(DefSymbol, ir.NewTypeDecl(ir.NewVarType(tvar))))
	}

	for _, arg := range args {
		c.binds = append(c.binds, NewDeclBind(DefSymbol, arg))
	}

	for _, ret := range rets {
		c.binds = append(c.binds, NewDeclBind(DefSymbol, ret))
	}
}

func (c *Context) IsExport(id string) bool {
	symbol, ok := c.lookupBind(id, FindDeclOnly)
	return ok && symbol.Symbol == ExportSymbol
}

func (c *Context) CheckModule() error {
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

func NewContext() *Context {
	return &Context{
		[]Bind{}, /* binds */
	}
}

func IsTypeWellFormed(c Context, t ir.IrType) error {
	switch t.Case {
	case ir.AliasType:
		if err := IsTypeWellFormed(c, t.Alias.Name); err != nil {
			return err
		}
		return IsTypeWellFormed(c, t.Alias.Value)

	case ir.ArrayType:
		return IsTypeWellFormed(c, t.Array.ElemType)

	case ir.ComponentType:
		return IsTypeWellFormed(c, t.Component.ElemType)

	case ir.ForallType:
		for _, tvar := range t.Forall.Vars {
			c.binds = append(c.binds, NewDeclBind(DefSymbol, ir.NewTypeDecl(ir.NewVarType(tvar))))
		}
		return IsTypeWellFormed(c, t.Forall.Type)

	case ir.FunType:
		if err := IsTypeWellFormed(c, t.Fun.Arg); err != nil {
			return err
		}
		return IsTypeWellFormed(c, t.Fun.Ret)

	case ir.NameType:
		if _, err := c.resolveTypeName(t); err == nil {
			return nil
		}

		_, err := c.getType(t)
		return err

	case ir.StructType:
		for _, typ := range t.FieldTypes() {
			if err := IsTypeWellFormed(c, typ); err != nil {
				return err
			}
		}
		return nil

	case ir.TupleType:
		for _, typ := range t.Tuple {
			if err := IsTypeWellFormed(c, typ); err != nil {
				return err
			}
		}
		return nil

	case ir.VarType:
		_, err := c.getType(t)
		return err

	default:
		panic(fmt.Errorf("unhandled %T %d", t.Case, t.Case))
	}
}
