package ir

import (
	"fmt"

	"github.com/zyedidia/generic/stack"
)

type FindCase int

const (
	FindAny = FindCase(iota)
	FindDeclOnly
	FindDefOnly
)

type IrContext struct {
	symbols []IrSymbol
	scopes  *stack.Stack[*irFunction]
}

func (c *IrContext) enterFunction(id string, args, rets []IrDecl) {
	function := NewFunction(id, args, rets)
	c.scopes.Push(&function)
}

func (c *IrContext) leaveFunction() {
	c.scopes.Pop()
}

func (c *IrContext) currentFunction() *irFunction {
	if c.scopes.Size() <= 0 {
		return nil
	}
	return c.scopes.Peek()
}

func (c *IrContext) lookupSymbol(id string, findCase FindCase) (IrSymbol, bool) {
	if findCase == FindAny || findCase == FindDefOnly {
		if fun := c.currentFunction(); fun != nil {
			if decl, err := fun.lookupVar(id); err == nil {
				return NewSymbol(DefSymbol, decl), true
			}
		}

		for i := len(c.symbols) - 1; i >= 0; i-- {
			symbol := c.symbols[i]
			if symbol.Case != DeclSymbol && symbol.Case != ExportSymbol && symbol.Decl.ID == id {
				return symbol, true
			}
		}
	}

	if findCase == FindAny || findCase == FindDeclOnly {
		for i := len(c.symbols) - 1; i >= 0; i-- {
			symbol := c.symbols[i]
			if (symbol.Case == DeclSymbol || symbol.Case == ExportSymbol) && symbol.Decl.ID == id {
				return symbol, true
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

	c.symbols = append(c.symbols, symbol)
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

	c.symbols = append(c.symbols, NewSymbol(DefSymbol, decl))
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

	c.symbols = append(c.symbols, NewSymbol(DefSymbol, decl))
	return nil
}

func (c *IrContext) addLocal(decl IrDecl) error {
	fun := c.currentFunction()
	if fun == nil {
		return fmt.Errorf("cannot define local %q outside of function", decl.ID)
	}

	return fun.addLocal(decl)
}

func (c *IrContext) checkModule() error {
	// Check all exports and all declarations have a definition (i.e., there are
	// no undefined exports or declarations).
	exported := map[string]struct{}{}
	declared := map[string]struct{}{}
	for _, symbol := range c.symbols {
		switch symbol.Case {
		case ExportSymbol:
			exported[symbol.Decl.ID] = struct{}{}
		case DeclSymbol:
			declared[symbol.Decl.ID] = struct{}{}
		}
	}

	for _, symbol := range c.symbols {
		if symbol.Case == DefSymbol {
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
		[]IrSymbol{}, /* symbols */
		stack.New[*irFunction](),
	}
}
