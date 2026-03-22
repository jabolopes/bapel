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

func (c Context) lookupTermDeclOrDefBind(name string) (Bind, bool) {
	return c.lookupBind(func(bind Bind) bool {
		return (bind.Is(TermDeclBind) && bind.TermDecl.Name == name) ||
			(bind.Is(TermDefBind) && bind.TermDef.Name == name)
	})
}

func (c Context) lookupTermDeclBindInScope(name string) (Bind, bool) {
	bind, ok := c.lookupBind(func(bind Bind) bool {
		return bind.Is(ScopeBind) || (bind.Is(TermDeclBind) && bind.TermDecl.Name == name)
	})
	if !ok || bind.Is(ScopeBind) {
		return Bind{}, false
	}
	return bind, true
}

func (c Context) lookupTermDefBindInScope(name string) (Bind, bool) {
	bind, ok := c.lookupBind(func(bind Bind) bool {
		return bind.Is(ScopeBind) || (bind.Is(TermDefBind) && bind.TermDef.Name == name)
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

func (c Context) containsTermDeclBindInScope(name string) bool {
	_, ok := c.lookupTermDeclBindInScope(name)
	return ok
}

func (c Context) containsTermDefBindInScope(name string) bool {
	_, ok := c.lookupTermDefBindInScope(name)
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

func (c Context) getTermDeclOrDefBind(name string) (Bind, error) {
	bind, ok := c.lookupTermDeclOrDefBind(name)
	if !ok {
		return Bind{}, fmt.Errorf("term %q is undefined", name)
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

func (c Context) enterFunction(typeVars []ir.VarKind, args []ir.FunctionArg) (Context, error) {
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
		if c, err = c.AddBind(NewTermDefBind(arg.ID, arg.Type)); err != nil {
			return c, err
		}
	}

	return c, nil
}

var (
	tvarBase = "ꞇ"
	tvarGen  = 0
)

func (c Context) GenFreshVarType() ir.IrType {
	a := rune(97)
	end := rune(122) + 1 // 'z' + 1

	shortNameUsed := make([]bool, end-a)

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

		r, _ := utf8.DecodeRuneInString(bind.TypeVar.Name)

		c := int(r - a)
		if c >= 0 && c < len(shortNameUsed) {
			shortNameUsed[c] = true
		}
	}

	for i, used := range shortNameUsed {
		if !used {
			return ir.NewVarType(string(rune(int(a) + i)))
		}
	}

	typ := ir.NewVarType(fmt.Sprintf("%s%d", tvarBase, tvarGen))
	tvarGen++
	return typ
}

func (c Context) GenFreshExistVar() ir.IrType {
	typ := ir.NewExistVarType(tvarGen)
	tvarGen++
	return typ
}

func (c Context) AddFreshType(typ ir.IrType) (Context, ir.IrType, ir.IrType, error) {
	if !typ.Is(ir.ForallType) && !typ.Is(ir.LambdaType) {
		return c, ir.IrType{}, typ, nil
	}

	var renamed ir.IrType
	var err error
	c, renamed, err = renameTypeVars(c, typ)
	if err != nil {
		return c, ir.IrType{}, ir.IrType{}, err
	}

	switch renamed.Case {
	case ir.ForallType:
		newContext, err := c.AddBind(NewTypeVarBind(renamed.Forall.Var, renamed.Forall.Kind))
		if err != nil {
			return c, ir.IrType{}, ir.IrType{}, err
		}
		return newContext, ir.NewVarType(renamed.Forall.Var), renamed.Forall.Type, nil

	case ir.LambdaType:
		newContext, err := c.AddBind(NewTypeVarBind(renamed.Lambda.Var, renamed.Lambda.Kind))
		if err != nil {
			return c, ir.IrType{}, ir.IrType{}, err
		}
		return newContext, ir.NewVarType(renamed.Lambda.Var), renamed.Lambda.Type, nil

	default:
		panic(fmt.Errorf("unhandled %T %d", renamed.Case, renamed.Case))
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

func (c Context) AddSymbol(decl ir.IrDecl) (Context, error) {
	var err error
	switch decl.Case {
	case ir.TermDecl:
		c, err = c.AddBind(NewTermDeclBind(decl.Term.ID, decl.Term.Type))
	case ir.AliasDecl:
		c, err = c.AddBind(NewAliasBind(decl.Alias.ID, decl.Alias.Type))
	case ir.NameDecl:
		c, err = c.AddBind(NewConstBind(decl.Name.ID, decl.Name.Kind))
	default:
		panic(fmt.Errorf("unhandled %T %d", decl.Case, decl.Case))
	}

	return c, err
}

func (c Context) EnterScope() (Context, error) {
	return c.enterScope()
}

func NewContext() Context {
	return Context{
		list.New[Bind](),
		0, /* wellformedSize */
	}
}
