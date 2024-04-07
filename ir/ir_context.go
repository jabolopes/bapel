package ir

import (
	"fmt"
	"strings"
)

type FindCase int

const (
	FindAny = FindCase(iota)
	FindDeclOnly
	FindDefOnly
)

type IrContext struct {
	binds []IrBind
}

func (c *IrContext) String() string {
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

func (c *IrContext) StringNoImports() string {
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

func (c *IrContext) lookupBind(id string, findCase FindCase) (IrBind, bool) {
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

	return IrBind{}, false
}

func (c *IrContext) getBind(id string, findCase FindCase) (IrBind, error) {
	bind, ok := c.lookupBind(id, findCase)
	if !ok {
		return IrBind{}, fmt.Errorf("%q is undefined", id)
	}

	return bind, nil
}

func (c *IrContext) lookupType(typ IrType) (IrBind, bool) {
	for i := len(c.binds) - 1; i >= 0; i-- {
		bind := c.binds[i]
		if bind.Case != DeclBind || bind.Decl.Case != TypeDecl {
			continue
		}

		if equalsType(bind.Decl.Type(), typ) {
			return bind, true
		}
	}

	return IrBind{}, false
}

func (c *IrContext) getType(typ IrType) (IrBind, error) {
	bind, ok := c.lookupType(typ)
	if !ok {
		return IrBind{}, fmt.Errorf("type %s is undefined", typ)
	}

	return bind, nil
}

func (c *IrContext) resolveTypeName(typ IrType) (IrType, error) {
	switch typ.Case {
	case AliasType:
		return c.resolveTypeName(typ.Alias.Value)

	case NameType:
		bind, ok := c.lookupBind(typ.Name, FindAny)
		if !ok {
			return IrType{}, fmt.Errorf("%q is undefined", typ.Name)
		}

		if equalsType(typ, bind.Decl.Type()) {
			return typ, nil
		}

		return c.resolveTypeName(bind.Decl.Type())

	default:
		// TODO: Should probably check if the type is defined in the context.
		return typ, nil
	}
}

func (c *IrContext) addBind(bind IrBind) error {
	bindID, ok := bind.ID()
	if ok {
		if _, ok := c.lookupBind(bindID, FindDefOnly); ok {
			return fmt.Errorf("%q is already defined", bindID)
		}

		if ok && bind.Symbol == DefSymbol {
			// Check that definition (e.g., function, struct, etc) matches declaration (if any).
			declaration, ok := c.lookupBind(bindID, FindDeclOnly)
			if ok {
				if err := NewIrTypechecker(c).subtype(declaration.Decl.Type(), bind.Decl.Type()); err != nil {
					return err
				}
			}
		}
	}

	c.binds = append(c.binds, bind)
	return nil
}

// TODO: Delete addBind.
func (c *IrContext) AddBind(bind IrBind) error {
	return c.addBind(bind)
}

func (c *IrContext) addMarker(id string) {
	// TODO: Call addBind instead and return propagate error.
	c.binds = append(c.binds, NewMarkerBind(id))
}

func (c *IrContext) removeTillMarker(id string) {
	for {
		// TODO: Check bounds and return an error.
		bind := c.binds[len(c.binds)-1]
		c.binds = c.binds[:len(c.binds)-1]

		if bind.Case == MarkerBind && *bind.Marker == id {
			return
		}
	}
}

func (c *IrContext) enterFunction(id string, typeVars []string, args, rets []IrDecl) {
	c.addMarker(id)

	for _, tvar := range typeVars {
		c.binds = append(c.binds, NewDeclBind(DefSymbol, NewTypeDecl(NewVarType(tvar))))
	}

	for _, arg := range args {
		c.binds = append(c.binds, NewDeclBind(DefSymbol, arg))
	}

	for _, ret := range rets {
		c.binds = append(c.binds, NewDeclBind(DefSymbol, ret))
	}
}

func (c *IrContext) isExport(id string) bool {
	symbol, ok := c.lookupBind(id, FindDeclOnly)
	return ok && symbol.Symbol == ExportSymbol
}

func (c *IrContext) checkModule() error {
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

func NewIrContext() *IrContext {
	return &IrContext{
		[]IrBind{}, /* binds */
	}
}

func isTypeWellFormed(c IrContext, t IrType) error {
	switch t.Case {
	case AliasType:
		if err := isTypeWellFormed(c, t.Alias.Name); err != nil {
			return err
		}
		return isTypeWellFormed(c, t.Alias.Value)

	case ArrayType:
		return isTypeWellFormed(c, t.Array.ElemType)

	case ForallType:
		for _, tvar := range t.Forall.Vars {
			c.binds = append(c.binds, NewDeclBind(DefSymbol, NewTypeDecl(NewVarType(tvar))))
		}
		return isTypeWellFormed(c, t.Forall.Type)

	case FunType:
		if err := isTypeWellFormed(c, t.Fun.Arg); err != nil {
			return err
		}
		return isTypeWellFormed(c, t.Fun.Ret)

	case NameType:
		_, err := c.getType(t)
		return err

	case NumberType:
		return nil

	case StructType:
		for _, typ := range t.FieldTypes() {
			if err := isTypeWellFormed(c, typ); err != nil {
				return err
			}
		}
		return nil

	case TupleType:
		for _, typ := range t.Tuple {
			if err := isTypeWellFormed(c, typ); err != nil {
				return err
			}
		}
		return nil

	case VarType:
		_, err := c.getType(t)
		return err

	default:
		panic(fmt.Errorf("unhandled IrTypeCase %d", t.Case))
	}
}
