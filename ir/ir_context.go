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
	imports      []IrDecl
	exports      []IrDecl
	decls        []IrDecl
	structDefs   []IrDecl
	functionDefs []IrDecl
	scopes       *stack.Stack[*irFunction]
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

		for _, d := range c.structDefs {
			if d.ID == id {
				return NewSymbol(DefSymbol, d), true
			}
		}

		for _, d := range c.functionDefs {
			if d.ID == id {
				return NewSymbol(DefSymbol, d), true
			}
		}

		for _, d := range c.imports {
			if d.ID == id {
				return NewSymbol(ImportSymbol, d), true
			}
		}
	}

	if findCase == FindAny || findCase == FindDeclOnly {
		for _, d := range c.decls {
			if d.ID == id {
				return NewSymbol(DeclSymbol, d), true
			}
		}

		for _, d := range c.exports {
			if d.ID == id {
				return NewSymbol(ExportSymbol, d), true
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
	for _, decl := range c.exports {
		if decl.ID == id {
			return true
		}
	}
	return false
}

func (c *IrContext) addImport(decl IrDecl) error {
	if _, ok := c.lookupSymbol(decl.ID, FindAny); ok {
		return fmt.Errorf("symbol %q is already declared, imported, exported, or defined", decl.ID)
	}

	c.imports = append(c.imports, decl)
	return nil
}

func (c *IrContext) addExport(decl IrDecl) error {
	if _, ok := c.lookupSymbol(decl.ID, FindAny); ok {
		return fmt.Errorf("symbol %q is already declared, imported, exported, or defined", decl.ID)
	}

	c.exports = append(c.exports, decl)
	return nil
}

func (c *IrContext) addDecl(decl IrDecl) error {
	if _, ok := c.lookupSymbol(decl.ID, FindAny); ok {
		return fmt.Errorf("symbol %q is already declared, imported, exported, or defined", decl.ID)
	}

	c.decls = append(c.decls, decl)
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

	c.functionDefs = append(c.functionDefs, decl)
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

	c.structDefs = append(c.structDefs, decl)
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
	for _, decl := range c.exports {
		exported[decl.ID] = struct{}{}
	}

	declared := map[string]struct{}{}
	for _, decl := range c.decls {
		declared[decl.ID] = struct{}{}
	}

	for _, decl := range c.structDefs {
		delete(exported, decl.ID)
		delete(declared, decl.ID)
	}

	for _, decl := range c.functionDefs {
		delete(exported, decl.ID)
		delete(declared, decl.ID)
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
		[]IrDecl{}, /* imports */
		[]IrDecl{}, /* exports */
		[]IrDecl{}, /* decls */
		[]IrDecl{}, /* structDefs */
		[]IrDecl{}, /* functionDefs */
		stack.New[*irFunction](),
	}
}
