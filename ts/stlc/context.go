package stlc

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/ts/list"
)

type FindCase int

const (
	FindAny = FindCase(iota)
	FindDeclOnly
	FindDefOnly
)

type Context struct {
	list           list.List[Bind]
	wellformedSize int
}

func (c Context) String() string {
	var b strings.Builder
	if !c.list.Empty() {
		binds := c.list.Iterate().Collect()
		b.WriteString(binds[0].String())
		for _, bind := range binds[1:] {
			b.WriteString(", ")
			b.WriteString(bind.String())
		}
	}
	return b.String()
}

func (c Context) empty() bool {
	return c.list.Empty()
}

func (c Context) lookupBind(is func(Bind) bool) (Bind, bool) {
	for it := c.list.Iterate(); ; {
		_, bind, ok := it.Next()
		if !ok {
			break
		}

		if is(bind) {
			return bind, true
		}
	}

	return Bind{}, false
}

func (c Context) lookupAliasBind(name string) (Bind, bool) {
	return c.lookupBind(func(bind Bind) bool {
		return bind.Is(AliasBind) && bind.Alias.Name == name
	})
}

func (c Context) lookupConstBind(name string) (Bind, bool) {
	return c.lookupBind(func(bind Bind) bool {
		return bind.Is(ConstBind) && bind.Const.Name == name
	})
}

func (c Context) lookupTermBind(name string) (Bind, bool) {
	return c.lookupBind(func(bind Bind) bool {
		return bind.Is(TermBind) && bind.Term.Name == name
	})
}

func (c Context) lookupTermBindInScope(name string) (Bind, bool) {
	bind, ok := c.lookupBind(func(bind Bind) bool {
		return bind.Is(ScopeBind) || (bind.Is(TermBind) && bind.Term.Name == name)
	})
	if !ok || bind.Is(ScopeBind) {
		return Bind{}, false
	}
	return bind, true
}

func (c Context) lookupTermBindInScopeWithSymbol(name string, symbol Symbol) (Bind, bool) {
	bind, ok := c.lookupBind(func(bind Bind) bool {
		return bind.Is(ScopeBind) || (bind.Is(TermBind) && bind.Term.Name == name && bind.Term.Symbol == symbol)
	})
	if !ok || bind.Is(ScopeBind) {
		return Bind{}, false
	}
	return bind, true
}

func (c Context) lookupScopeBind() (Bind, bool) {
	return c.lookupBind(func(bind Bind) bool {
		return bind.Is(ScopeBind)
	})
}

func (c Context) lookupTypeVarBind(tvar string) (Bind, bool) {
	return c.lookupBind(func(bind Bind) bool {
		return bind.Is(TypeVarBind) && bind.TypeVar.Name == tvar
	})
}

func (c Context) pop() (Bind, Context) {
	bind, ok := c.list.Value()
	if !ok {
		panic("Context is empty")
	}

	c.list = c.list.Remove()
	c.wellformedSize = min(c.wellformedSize, c.list.Size())
	return bind, c
}

func (c Context) containsAliasBind(name string) bool {
	_, ok := c.lookupAliasBind(name)
	return ok
}

func (c Context) containsConstBind(name string) bool {
	_, ok := c.lookupConstBind(name)
	return ok
}

func (c Context) containsTermBindInScope(name string) bool {
	_, ok := c.lookupTermBindInScope(name)
	return ok
}

func (c Context) containsTermBindInScopeWithSymbol(name string, symbol Symbol) bool {
	_, ok := c.lookupTermBindInScopeWithSymbol(name, symbol)
	return ok
}

func (c Context) containsTypeVarBind(tvar string) bool {
	_, ok := c.lookupTypeVarBind(tvar)
	return ok
}

func (c Context) getAliasBind(name string) (Bind, error) {
	bind, ok := c.lookupAliasBind(name)
	if !ok {
		return Bind{}, fmt.Errorf("type %q is undefined", name)
	}
	return bind, nil
}

func (c Context) getConstBind(name string) (Bind, error) {
	bind, ok := c.lookupConstBind(name)
	if !ok {
		return Bind{}, fmt.Errorf("type %q is undefined", name)
	}
	return bind, nil
}

func (c Context) getTermBind(name string) (Bind, error) {
	bind, ok := c.lookupTermBind(name)
	if !ok {
		return Bind{}, fmt.Errorf("%q is undefined", name)
	}
	return bind, nil
}

func (c Context) getTypeVarBind(tvar string) (Bind, error) {
	bind, ok := c.lookupTypeVarBind(tvar)
	if !ok {
		return Bind{}, fmt.Errorf("type variable %q is undefined", tvar)
	}
	return bind, nil
}

func (c Context) enterScope() (Context, error) {
	if bind, ok := c.lookupScopeBind(); ok {
		return c.AddBind(NewScopeBind(bind.Scope.Level + 1))
	}

	return c.AddBind(NewScopeBind(1))
}

func (c Context) enterFunction(typeVars []ir.VarKind, args []ir.IrDecl) (Context, error) {
	var err error
	c, err = c.enterScope()
	if err != nil {
		return c, err
	}

	for _, tvar := range typeVars {
		if c, err = c.AddBind(NewTypeVarBind(tvar.Var, tvar.Kind)); err != nil {
			return c, err
		}
	}

	for _, arg := range args {
		if c, err = c.AddBind(NewTermBind(arg.Term.ID, arg.Term.Type, DefSymbol)); err != nil {
			return c, err
		}
	}

	return c, nil
}

func (c Context) GenFreshVarType() ir.IrType {
	free := rune(97)

	for it := c.list.Iterate(); ; {
		_, bind, ok := it.Next()
		if !ok {
			break
		}

		if !bind.Is(TypeVarBind) {
			continue
		}

		if len(bind.TypeVar.Name) != 1 {
			continue
		}

		if r, _ := utf8.DecodeRuneInString(bind.TypeVar.Name); r >= free {
			free = r + 1
		}
	}

	return ir.NewVarType(string(free))
}

func (c Context) AddFreshType(typ ir.IrType) (Context, ir.IrType, ir.IrType, error) {
	switch typ.Case {
	case ir.ForallType:
		tvar := c.GenFreshVarType()
		newContext, err := c.AddBind(NewTypeVarBind(tvar.Var, typ.Forall.Kind))
		if err != nil {
			return c, ir.IrType{}, ir.IrType{}, err
		}
		return newContext, tvar, ir.SubstituteType(typ.Forall.Type, ir.NewVarType(typ.Forall.Var), tvar), nil

	case ir.LambdaType:
		tvar := c.GenFreshVarType()
		newContext, err := c.AddBind(NewTypeVarBind(tvar.Var, typ.Lambda.Kind))
		if err != nil {
			return c, ir.IrType{}, ir.IrType{}, err
		}
		return newContext, tvar, ir.SubstituteType(typ.Lambda.Type, ir.NewVarType(typ.Lambda.Var), tvar), nil

	default:
		return c, ir.IrType{}, typ, nil
	}
}

func (c Context) AddBind(bind Bind) (Context, error) {
	origC := c

	c.list = c.list.Add(bind)
	if err := isWellformedContext(c); err != nil {
		return origC, err
	}

	c.wellformedSize = c.list.Size()
	return c, nil
}

func (c Context) AddSymbol(decl ir.IrDecl, symbol Symbol) (Context, error) {
	var err error
	switch decl.Case {
	case ir.TermDecl:
		c, err = c.AddBind(NewTermBind(decl.Term.ID, decl.Term.Type, symbol))
	case ir.AliasDecl:
		c, err = c.AddBind(NewAliasBind(decl.Alias.ID, decl.Alias.Type, symbol))
	case ir.NameDecl:
		c, err = c.AddBind(NewConstBind(decl.Name.ID, decl.Name.Kind, symbol))
	default:
		panic(fmt.Errorf("unhandled %T %d", decl.Case, decl.Case))
	}

	return c, err
}

func (c Context) EnterScope() (Context, error) {
	return c.enterScope()
}

func (c Context) wellformedUnderTvar(tvar, typ ir.IrType) (bool, error) {
	if !tvar.Is(ir.VarType) {
		return false, fmt.Errorf("expected type variable; got %s", tvar)
	}

	for !c.empty() {
		var bind Bind
		bind, c = c.pop()

		if bind.Is(TypeVarBind) && bind.TypeVar.Name == tvar.Var {
			break
		}
	}

	if err := isWellformedType(c, typ); err != nil {
		return false, nil
	}

	return true, nil
}

func NewContext() Context {
	return Context{
		list.New[Bind](),
		0, /* wellformedSize */
	}
}
