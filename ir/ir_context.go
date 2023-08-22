package ir

import (
	"fmt"
)

type FindCase int

const (
	FindAny = FindCase(iota)
	FindDeclOnly
	FindDefOnly
	FindVarOnly
)

type IrContext struct {
	imports    []irDecl
	exports    []irDecl
	decls      []irDecl
	structDefs []irDecl
	functions  []irFunction
}

func (c *IrContext) fun() *irFunction {
	return &c.functions[len(c.functions)-1]
}

func (c *IrContext) lookupSymbol(id string, findCase FindCase) (IrSymbol, bool) {
	if findCase == FindAny || findCase == FindVarOnly {
		if len(c.functions) > 0 {
			if irvar, err := c.fun().lookupVar(id); err == nil {
				return NewSymbol(ReferenceSymbol, irvar.decl()), true
			}
		}
	}

	if findCase == FindAny || findCase == FindDefOnly {
		for _, d := range c.structDefs {
			if d.id == id {
				return NewSymbol(StructSymbol, d), true
			}
		}

		for _, f := range c.functions {
			if f.id == id {
				return NewSymbol(FunctionSymbol, f.decl()), true
			}
		}

		for _, d := range c.imports {
			if d.id == id {
				return NewSymbol(ImportSymbol, d), true
			}
		}
	}

	if findCase == FindAny || findCase == FindDeclOnly {
		for _, d := range c.decls {
			if d.id == id {
				return NewSymbol(DeclSymbol, d), true
			}
		}

		for _, d := range c.exports {
			if d.id == id {
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

	return decl.typ, nil
}

func (c *IrContext) isExport(id string) bool {
	for _, decl := range c.exports {
		if decl.id == id {
			return true
		}
	}
	return false
}

func (c *IrContext) addImport(decl irDecl) error {
	if _, ok := c.lookupSymbol(decl.id, FindAny); ok {
		return fmt.Errorf("symbol %q is already declared, imported, exported, or defined", decl.id)
	}

	c.imports = append(c.imports, decl)
	return nil
}

func (c *IrContext) addExport(decl irDecl) error {
	if _, ok := c.lookupSymbol(decl.id, FindAny); ok {
		return fmt.Errorf("symbol %q is already declared, imported, exported, or defined", decl.id)
	}

	c.exports = append(c.exports, decl)
	return nil
}

func (c *IrContext) addDecl(decl irDecl) error {
	if _, ok := c.lookupSymbol(decl.id, FindAny); ok {
		return fmt.Errorf("symbol %q is already declared, imported, exported, or defined", decl.id)
	}

	c.decls = append(c.decls, decl)
	return nil
}

func (c *IrContext) addFunction(function irFunction) error {
	if _, ok := c.lookupSymbol(function.id, FindDefOnly); ok {
		return fmt.Errorf("symbol %q already defined", function.id)
	}

	// Check function definition matches declaration (if any).
	if symbol, ok := c.lookupSymbol(function.id, FindDeclOnly); ok {
		if err := NewIrTypechecker(c).MatchesDecl(symbol.Decl, function.decl()); err != nil {
			return err
		}
	}

	c.functions = append(c.functions, function)
	return nil
}

func (c *IrContext) addStruct(decl irDecl) error {
	if _, ok := c.lookupSymbol(decl.id, FindDefOnly); ok {
		return fmt.Errorf("symbol %q already defined", decl.id)
	}

	// Check struct definition matches declaration (if any).
	if symbol, ok := c.lookupSymbol(decl.id, FindDeclOnly); ok {
		if err := NewIrTypechecker(c).MatchesDecl(symbol.Decl, decl); err != nil {
			return err
		}
	}

	c.structDefs = append(c.structDefs, decl)
	return nil
}

func (c *IrContext) checkModule() error {
	// Check all exports and all declarations have a definition (i.e., there are
	// no undefined exports or declarations).
	exported := map[string]struct{}{}
	for _, decl := range c.exports {
		exported[decl.id] = struct{}{}
	}

	declared := map[string]struct{}{}
	for _, decl := range c.decls {
		declared[decl.id] = struct{}{}
	}

	for _, decl := range c.structDefs {
		delete(exported, decl.id)
		delete(declared, decl.id)
	}

	for _, function := range c.functions {
		delete(exported, function.id)
		delete(declared, function.id)
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
		[]irDecl{},     /* imports */
		[]irDecl{},     /* exports */
		[]irDecl{},     /* decls */
		[]irDecl{},     /* structDefs */
		[]irFunction{}, /* functions */
	}
}
