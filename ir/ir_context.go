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
		if c.binds[0].Case != TermBind || c.binds[0].Term.Case == DefSymbol {
			b.WriteString(c.binds[0].String())
		}

		for _, bind := range c.binds[1:] {
			if bind.Case != TermBind || bind.Term.Case == DefSymbol {
				b.WriteString(", ")
				b.WriteString(bind.String())
			}
		}
	}
	return b.String()
}

func (c *IrContext) addBind(bind IrBind) error {
	// TODO: Check if the bind is already defined.
	c.binds = append(c.binds, bind)
	return nil
}

func (c *IrContext) addMarker(id string) {
	// TODO: Call addBind instead and return propagate error.
	c.binds = append(c.binds, NewMarkerBind(id))
}

func (c *IrContext) addDeclaration(symbol IrSymbol) error {
	if _, ok := c.lookupSymbol(symbol.Decl.ID, FindAny); ok {
		return fmt.Errorf("symbol %q is already declared, imported, exported, or defined", symbol.Decl.ID)
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

func (c *IrContext) lookupType(typ IrType) (IrBind, bool) {
	for i := len(c.binds) - 1; i >= 0; i-- {
		bind := c.binds[i]
		if bind.Case == TypeBind && equalsType(bind.Type.Type, typ) {
			return bind, true
		}
	}

	return IrBind{}, false
}

func (c *IrContext) lookupBind(id string, findCase FindCase) (IrBind, bool) {
	if findCase == FindAny || findCase == FindDefOnly {
		for i := len(c.binds) - 1; i >= 0; i-- {
			bind := c.binds[i]
			if bindID, ok := bind.ID(); !ok || bindID != id {
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
			if bindID, ok := bind.ID(); !ok || bindID != id {
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

func (c *IrContext) getType(id string, findCase FindCase) (IrType, error) {
	bind, ok := c.lookupBind(id, findCase)
	if !ok {
		return IrType{}, fmt.Errorf("id %q is undefined", id)
	}

	switch bind.Case {
	case TermBind:
		return bind.Term.Decl.Type, nil

	case TypeBind:
		if bind.Type.Solution != nil {
			if bind.Type.Solution.Case == VarExistType {
				return c.getType(bind.Type.Solution.VarExist.Var, findCase)
			}
			return *bind.Type.Solution, nil
		}

		return bind.Type.Type, nil

	default:
		return IrType{}, fmt.Errorf("id %q is not associated with a type", id)
	}
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

func (c *IrContext) isSolvedVar(id string) bool {
	bind, ok := c.lookupBind(id, FindAny)
	if !ok {
		return false
	}

	return bind.Case == TypeBind &&
		bind.Type.Type.Case == VarExistType &&
		bind.Type.Solution != nil
}

func (c *IrContext) isDefinedInOrder(id1, id2 string) bool {
	var i1 *int
	var i2 *int

	for i := len(c.binds) - 1; i >= 0; i-- {
		if i1 != nil && i2 != nil {
			break
		}

		bind := c.binds[i]
		if bindID, ok := bind.ID(); ok && bindID == id1 {
			j := i
			i1 = &j
		}

		if bindID, ok := bind.ID(); ok && bindID == id2 {
			j := i
			i2 = &j
		}
	}

	if i1 == nil || i2 == nil {
		return false
	}

	return *i1 < *i2
}

func (c *IrContext) setType(id string, typ IrType) error {
	if !IsMonotype(*c, typ) {
		return fmt.Errorf("cannot assign non-monotype %s to type variable %q", typ, id)
	}

	bind, ok := c.lookupBind(id, FindAny)
	if !ok {
		return fmt.Errorf("symbol %q is undefined", id)
	}

	// TODO: Check that 'typ' is defined in the context.

	switch bind.Case {
	case TermBind:
		return fmt.Errorf("cannot assign type to term binding %q", id)

	case TypeBind:
		if bind.Type.Type.Case != VarExistType {
			return fmt.Errorf("cannot assign a type to %s", bind.Type.Type)
		}

		bind.Type.Solution = &typ
		return nil

	default:
		return fmt.Errorf("cannot assign type to %q", id)
	}
}

func (c *IrContext) isExport(id string) bool {
	symbol, ok := c.lookupSymbol(id, FindDeclOnly)
	return ok && symbol.Case == ExportSymbol
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
			exported[symbol.Decl.ID] = struct{}{}
		case DeclSymbol:
			declared[symbol.Decl.ID] = struct{}{}
		}
	}

	for _, bind := range c.binds {
		if bind.Case != TermBind {
			continue
		}

		if symbol := bind.Term; symbol.Case == DefSymbol {
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

func IsMonotype(c IrContext, t IrType) bool {
	switch t.Case {
	case ArrayType:
		return IsMonotype(c, t.Array.ElemType)
	case ForallType:
		return false
	case FunType:
		return IsMonotype(c, t.Fun.Arg) && IsMonotype(c, t.Fun.Ret)

	case NameType:
		typ, err := c.getType(t.Name, FindAny)
		if err != nil {
			panic(fmt.Sprintf("typename %q is undefined", t.Name))
		}

		if equalsType(t, typ) {
			return true
		}

		return IsMonotype(c, typ)

	case NumberType:
		return true
	case StructType:
		return IsMonotype(c, NewTupleType(t.FieldTypes()))

	case TupleType:
		for _, typ := range t.Tuple {
			if !IsMonotype(c, typ) {
				return false
			}
		}
		return true

	case VarType:
		return true
	case VarExistType:
		return true

	default:
		panic(fmt.Errorf("unhandled IrTypeCase %d", t.Case))
	}
}

func isTypeWellformed(c IrContext, t IrType) bool {
	switch t.Case {
	case ArrayType:
		return isTypeWellformed(c, t.Array.ElemType)

	case ForallType:
		for _, tvar := range t.Forall.Vars {
			c.binds = append(c.binds, NewTypeBind(NewVarType(tvar), nil))
		}
		return isTypeWellformed(c, t.Forall.Type)

	case FunType:
		return isTypeWellformed(c, t.Fun.Arg) && isTypeWellformed(c, t.Fun.Ret)
	case NameType:
		_, ok := c.lookupBind(t.Name, FindAny)
		return ok
	case NumberType:
		return true

	case StructType:
		for _, typ := range t.FieldTypes() {
			if !isTypeWellformed(c, typ) {
				return false
			}
		}
		return true

	case TupleType:
		for _, typ := range t.Tuple {
			if !isTypeWellformed(c, typ) {
				return false
			}
		}
		return true

	case VarType:
		_, ok := c.lookupType(t)
		return ok

	case VarExistType:
		_, ok := c.lookupType(t)
		return ok

	default:
		panic(fmt.Errorf("unhandled IrTypeCase %d", t.Case))
	}
}

func sliceAtType(c IrContext, typ IrType) IrContext {
	for len(c.binds) > 0 {
		bind := c.binds[len(c.binds)-1]
		if bind.Case == TypeBind && equalsType(bind.Type.Type, typ) {
			break
		}

		c.binds = c.binds[:len(c.binds)-1]
	}

	return c
}
