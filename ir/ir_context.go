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

func (c *IrContext) addBind(bind IrBind) error {
	bindID, ok := bind.ID()
	if ok {
		if _, ok := c.lookupBind(bindID, FindDefOnly); ok {
			return fmt.Errorf("%q is already defined", bindID)
		}
	}

	c.binds = append(c.binds, bind)
	return nil
}

func (c *IrContext) addMarker(id string) {
	// TODO: Call addBind instead and return propagate error.
	c.binds = append(c.binds, NewMarkerBind(id))
}

func (c *IrContext) addDefinition(symbol IrSymbol, decl IrDecl) error {
	// Check definition (e.g., function, struct, etc) matches declaration (if any).
	if symbolDecl, err := c.getDecl(decl.ID, FindDeclOnly); err == nil {
		if err := NewIrTypechecker(c).MatchesDecl(symbolDecl, decl); err != nil {
			return err
		}
	}

	return c.addBind(NewBindFromDecl(symbol, decl))
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
		c.binds = append(c.binds, NewTermBind(DefSymbol, arg))
	}

	for _, ret := range rets {
		c.binds = append(c.binds, NewTermBind(DefSymbol, ret))
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

// TODO: Deduplicate with 'getType()'.
func (c *IrContext) getDecl(id string, findCase FindCase) (IrDecl, error) {
	bind, ok := c.lookupBind(id, findCase)
	if !ok {
		return IrDecl{}, fmt.Errorf("id %q is undefined", id)
	}

	typ, err := c.getType(id, findCase)
	if err != nil {
		return IrDecl{}, err
	}

	switch bind.Case {
	case TermBind:
		return NewTermDecl(id, typ), nil

	case TypeBind:
		return NewTypeDecl(id, typ), nil

	default:
		return IrDecl{}, fmt.Errorf("id %q is not associated with a type", id)
	}
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

	if !isTypeWellformed(*c, typ) {
		return fmt.Errorf("type %s is not wellformed", typ)
	}

	switch bind.Case {
	case TermBind:
		return fmt.Errorf("cannot assign type to term binding %q", id)

	case TypeBind:
		if bind.Type.Type.Case != NameType && bind.Type.Type.Case != VarExistType {
			return fmt.Errorf("cannot assign a type to %s", bind.Type.Type)
		}

		bind.Type.Solution = &typ
		return nil

	default:
		return fmt.Errorf("cannot assign type to %q", id)
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
			c.binds = append(c.binds, NewTypeBind(DefSymbol, NewVarType(tvar), nil))
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
