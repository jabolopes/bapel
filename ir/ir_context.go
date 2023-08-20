package ir

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

		for _, d := range c.imports {
			if d.id == id {
				return d, true
			}
		}
	}

	return irDecl{}, false
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
