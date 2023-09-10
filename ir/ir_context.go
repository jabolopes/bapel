package ir

import (
	"fmt"
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

func (c *IrContext) addMarker(id string) {
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

func (c *IrContext) enterFunction(id string, args, rets []IrDecl) {
	c.addMarker(id)

	for _, arg := range args {
		c.binds = append(c.binds, NewTermBind(NewSymbolFromDecl(DefSymbol, arg)))
	}

	for _, ret := range rets {
		c.binds = append(c.binds, NewTermBind(NewSymbolFromDecl(DefSymbol, ret)))
	}
}

func (c *IrContext) lookupBind(id string, findCase FindCase) (IrBind, bool) {
	if findCase == FindAny || findCase == FindDefOnly {
		for i := len(c.binds) - 1; i >= 0; i-- {
			bind := c.binds[i]
			if bind.ID() != id {
				continue
			}

			switch bind.Case {
			case TermBind:
				if symbol := bind.Term; symbol.Case != DeclSymbol && symbol.Case != ExportSymbol {
					return bind, true
				}
			case TypeBind:
				return bind, true
			}
		}
	}

	if findCase == FindAny || findCase == FindDeclOnly {
		for i := len(c.binds) - 1; i >= 0; i-- {
			bind := c.binds[i]
			if bind.ID() != id {
				continue
			}

			switch bind.Case {
			case TermBind:
				if symbol := bind.Term; symbol.Case == DeclSymbol || symbol.Case == ExportSymbol {
					return bind, true
				}
			case TypeBind:
				return bind, true
			}
		}
	}

	return IrBind{}, false
}

func (c *IrContext) lookupSymbol(id string, findCase FindCase) (IrSymbol, bool) {
	bind, ok := c.lookupBind(id, findCase)
	if !ok {
		return IrSymbol{}, false
	}

	if bind.Case == TermBind {
		return *bind.Term, true
	}

	return IrSymbol{}, false
}

func (c *IrContext) getSymbol(id string, findCase FindCase) (IrSymbol, error) {
	if symbol, ok := c.lookupSymbol(id, findCase); ok {
		return symbol, nil
	}

	return IrSymbol{}, fmt.Errorf("undefined symbol %q", id)
}

func (c *IrContext) getType(id string, findCase FindCase) (IrType, error) {
	decl, err := c.getSymbol(id, findCase)
	if err != nil {
		return IrType{}, err
	}

	if decl.Type == nil {
		return IrType{}, fmt.Errorf("symbol is not assigned a type")
	}

	return *decl.Type, nil
}

func (c *IrContext) getDecl(id string, findCase FindCase) (IrDecl, error) {
	symbol, err := c.getSymbol(id, findCase)
	if err != nil {
		return IrDecl{}, err
	}

	if symbol.Type == nil {
		return IrDecl{}, fmt.Errorf("symbol %q is not assigned a type", id)
	}

	return IrDecl{symbol.DeclCase, symbol.ID, *symbol.Type}, nil
}

func (c *IrContext) setType(id string, typ IrType) error {
	for i := len(c.binds) - 1; i >= 0; i-- {
		bind := c.binds[i]
		if bind.Case != TermBind || bind.Term.ID != id {
			continue
		}

		if bind.Term.Type != nil {
			return fmt.Errorf("symbol %q is already assigned type %s", id, typ)
		}

		c.binds[i].Term.Type = &typ
		return nil
	}

	return fmt.Errorf("symbol %q is not defined", id)
}

func (c *IrContext) isExport(id string) bool {
	symbol, ok := c.lookupSymbol(id, FindDeclOnly)
	return ok && symbol.Case == ExportSymbol
}

func (c *IrContext) addDeclaration(symbol IrSymbol) error {
	if _, ok := c.lookupSymbol(symbol.ID, FindAny); ok {
		return fmt.Errorf("symbol %q is already declared, imported, exported, or defined", symbol.ID)
	}

	c.binds = append(c.binds, NewTermBind(symbol))
	return nil
}

func (c *IrContext) addDefinition(decl IrDecl) error {
	// TODO: Exclude imports, otherwise someone exporting a new symbol
	// will break someone else's code.
	if _, ok := c.lookupSymbol(decl.ID, FindDefOnly); ok {
		return fmt.Errorf("symbol %q already defined", decl.ID)
	}

	// Check definition (e.g., function, struct, etc) matches declaration (if any).
	if symbolDecl, err := c.getDecl(decl.ID, FindDeclOnly); err == nil {
		if err := NewIrTypechecker(c).MatchesDecl(symbolDecl, decl); err != nil {
			return err
		}
	}

	c.binds = append(c.binds, NewTermBind(NewSymbolFromDecl(DefSymbol, decl)))
	return nil
}

func (c *IrContext) checkModule() error {
	// Check all exports and all declarations have a definition (i.e., there are
	// no undefined exports or declarations).
	exported := map[string]struct{}{}
	declared := map[string]struct{}{}
	for _, bind := range c.binds {
		if bind.Case != TermBind {
			continue
		}

		switch symbol := bind.Term; symbol.Case {
		case ExportSymbol:
			exported[symbol.ID] = struct{}{}
		case DeclSymbol:
			declared[symbol.ID] = struct{}{}
		}
	}

	for _, bind := range c.binds {
		if bind.Case != TermBind {
			continue
		}

		if symbol := bind.Term; symbol.Case == DefSymbol {
			delete(exported, symbol.ID)
			delete(declared, symbol.ID)
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
