package ir

import "fmt"

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

func (c *IrContext) lookupDecl(id string, findCase FindCase) (irDecl, bool) {
	if findCase == FindAny || findCase == FindVarOnly {
		if len(c.functions) > 0 {
			if irvar, err := c.fun().lookupVar(id); err == nil {
				return irvar.decl(), true
			}
		}
	}

	if findCase == FindAny || findCase == FindDefOnly {
		for _, d := range c.structDefs {
			if d.id == id {
				return d, true
			}
		}

		for _, f := range c.functions {
			if f.id == id {
				return f.decl(), true
			}
		}

		for _, d := range c.imports {
			if d.id == id {
				return d, true
			}
		}
	}

	if findCase == FindAny || findCase == FindDeclOnly {
		for _, d := range c.decls {
			if d.id == id {
				return d, true
			}
		}

		for _, d := range c.exports {
			if d.id == id {
				return d, true
			}
		}
	}

	return irDecl{}, false
}

func (c *IrContext) isExport(id string) bool {
	for _, decl := range c.exports {
		if decl.id == id {
			return true
		}
	}
	return false
}

func (c *IrContext) addFunction(function irFunction) error {
	if _, ok := c.lookupDecl(function.id, FindDefOnly); ok {
		return fmt.Errorf("symbol %q already defined", function.id)
	}

	// Check function definition matches declaration (if any).
	if decl, ok := c.lookupDecl(function.id, FindDeclOnly); ok {
		if err := matchesDecl(decl, function.decl()); err != nil {
			return err
		}
	}

	c.functions = append(c.functions, function)
	return nil
}

func (c *IrContext) addStruct(decl irDecl) error {
	if _, ok := c.lookupDecl(decl.id, FindDefOnly); ok {
		return fmt.Errorf("symbol %q already defined", decl.id)
	}

	// Check struct definition matches declaration (if any).
	if formalDecl, ok := c.lookupDecl(decl.id, FindDeclOnly); ok {
		if err := matchesDecl(formalDecl, decl); err != nil {
			return err
		}
	}

	c.structDefs = append(c.structDefs, decl)
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
