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

func (c *IrContext) enterFunction(id string, args, rets []IrDecl) {
	c.binds = append(c.binds, NewScopeBind(id))

	for _, arg := range args {
		c.binds = append(c.binds, NewSymbolBind(NewDefSymbol(arg)))
	}

	for _, ret := range rets {
		c.binds = append(c.binds, NewSymbolBind(NewDefSymbol(ret)))
	}
}

func (c *IrContext) leaveFunction(id string) {
	for {
		// TODO: Check bounds and return an error.
		bind := c.binds[len(c.binds)-1]
		c.binds = c.binds[:len(c.binds)-1]

		if bind.Case == ScopeBind && *bind.Scope == id {
			return
		}
	}
}

func (c *IrContext) lookupSymbol(id string, findCase FindCase) (IrSymbol, bool) {
	if findCase == FindAny || findCase == FindDefOnly {
		for i := len(c.binds) - 1; i >= 0; i-- {
			bind := c.binds[i]
			if bind.Case != SymbolBind {
				continue
			}

			symbol := bind.Symbol
			if symbol.Case != DeclSymbol && symbol.Case != ExportSymbol && symbol.Decl.ID == id {
				return *symbol, true
			}
		}
	}

	if findCase == FindAny || findCase == FindDeclOnly {
		for i := len(c.binds) - 1; i >= 0; i-- {
			bind := c.binds[i]
			if bind.Case != SymbolBind {
				continue
			}

			symbol := bind.Symbol
			if (symbol.Case == DeclSymbol || symbol.Case == ExportSymbol) && symbol.Decl.ID == id {
				return *symbol, true
			}
		}
	}

	return IrSymbol{}, false
}

func (c *IrContext) getSymbol(id string, findCase FindCase) (IrSymbol, error) {
	if symbol, ok := c.lookupSymbol(id, findCase); ok {
		return symbol, nil
	}

	return IrSymbol{}, fmt.Errorf("undefined symbol %q", id)
}

func (c *IrContext) getDecl(id string, findCase FindCase) (IrDecl, error) {
	symbol, err := c.getSymbol(id, findCase)
	if err != nil {
		return IrDecl{}, err
	}

	return symbol.Decl, nil
}

func (c *IrContext) getType(id string, findCase FindCase) (IrType, error) {
	decl, err := c.getDecl(id, findCase)
	if err != nil {
		return IrType{}, err
	}

	return decl.Type, nil
}

func (c *IrContext) isExport(id string) bool {
	symbol, ok := c.lookupSymbol(id, FindDeclOnly)
	return ok && symbol.Case == ExportSymbol
}

func (c *IrContext) addDeclaration(symbol IrSymbol) error {
	if _, ok := c.lookupSymbol(symbol.Decl.ID, FindAny); ok {
		return fmt.Errorf("symbol %q is already declared, imported, exported, or defined", symbol.Decl.ID)
	}

	c.binds = append(c.binds, NewSymbolBind(symbol))
	return nil
}

func (c *IrContext) addFunction(decl IrDecl) error {
	if _, ok := c.lookupSymbol(decl.ID, FindDefOnly); ok {
		return fmt.Errorf("symbol %q already defined", decl.ID)
	}

	// Check function definition matches declaration (if any).
	if symbol, ok := c.lookupSymbol(decl.ID, FindDeclOnly); ok {
		if err := NewIrTypechecker(c).MatchesDecl(symbol.Decl, decl); err != nil {
			return err
		}
	}

	c.binds = append(c.binds, NewSymbolBind(NewSymbol(DefSymbol, decl)))
	return nil
}

func (c *IrContext) addStruct(decl IrDecl) error {
	if _, ok := c.lookupSymbol(decl.ID, FindDefOnly); ok {
		return fmt.Errorf("symbol %q already defined", decl.ID)
	}

	// Check struct definition matches declaration (if any).
	if symbol, ok := c.lookupSymbol(decl.ID, FindDeclOnly); ok {
		if err := NewIrTypechecker(c).MatchesDecl(symbol.Decl, decl); err != nil {
			return err
		}
	}

	c.binds = append(c.binds, NewSymbolBind(NewSymbol(DefSymbol, decl)))
	return nil
}

func (c *IrContext) addLocal(decl IrDecl) error {
	// TODO: Exclude imports, otherwise someone exporting a new symbol
	// will break someone else's code.
	if _, ok := c.lookupSymbol(decl.ID, FindAny); ok {
		return fmt.Errorf("symbol %q already defined", decl.ID)
	}

	c.binds = append(c.binds, NewSymbolBind(NewSymbol(DefSymbol, decl)))
	return nil
}

func (c *IrContext) checkModule() error {
	// Check all exports and all declarations have a definition (i.e., there are
	// no undefined exports or declarations).
	exported := map[string]struct{}{}
	declared := map[string]struct{}{}
	for _, bind := range c.binds {
		if bind.Case != SymbolBind {
			continue
		}

		switch symbol := bind.Symbol; symbol.Case {
		case ExportSymbol:
			exported[symbol.Decl.ID] = struct{}{}
		case DeclSymbol:
			declared[symbol.Decl.ID] = struct{}{}
		}
	}

	for _, bind := range c.binds {
		if bind.Case != SymbolBind {
			continue
		}

		if symbol := bind.Symbol; symbol.Case == DefSymbol {
			delete(exported, symbol.Decl.ID)
			delete(declared, symbol.Decl.ID)
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
