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

func (c Context) lookupTypeVarBindInScope(tvar string) (Bind, bool) {
	bind, ok := c.lookupBind(func(bind Bind) bool {
		return bind.Is(ScopeBind) || bind.Is(TypeVarBind) && bind.TypeVar.Name == tvar
	})
	if !ok || bind.Is(ScopeBind) {
		return Bind{}, false
	}
	return bind, true
}

func (c Context) lookupTraitBind(name string) (Bind, bool) {
	return c.lookupBind(func(bind Bind) bool {
		return bind.Is(TraitBind) && bind.Trait.Name == name
	})
}

func (c Context) containsTraitBind(name string) bool {
	_, ok := c.lookupTraitBind(name)
	return ok
}

func (c Context) lookupTypeOrTraitBind(name string) (Bind, bool) {
	return c.lookupBind(func(b Bind) bool {
		return (b.Is(AliasBind) && b.Alias.Name == name) ||
			(b.Is(ConstBind) && b.Const.Name == name) ||
			(b.Is(TraitBind) && b.Trait.Name == name)
	})
}

func (c Context) GetTraitBind(name string) (Bind, error) {
	bind, ok := c.lookupTraitBind(name)
	if !ok {
		return Bind{}, fmt.Errorf("trait %q is undefined", name)
	}
	return bind, nil
}

func (c Context) lookupTraitImplBind(traitType ir.IrType, typeName ir.IrType) (Bind, bool) {
	return c.lookupBind(func(bind Bind) bool {
		if !bind.Is(TraitImplBind) {
			return false
		}
		impl := bind.TraitImpl
		vars := make(map[string]bool)
		for _, tp := range impl.TypeParams {
			vars[tp.Var] = true
		}

		subs := make(map[string]ir.IrType)
		subs, ok := matchType(impl.TraitType, traitType, vars, subs)
		if !ok {
			return false
		}
		_, ok = matchType(impl.TypeName, typeName, vars, subs)
		return ok
	})
}

func matchType(pattern ir.IrType, target ir.IrType, vars map[string]bool, subs map[string]ir.IrType) (map[string]ir.IrType, bool) {
	if pattern.Is(ir.VarType) && vars[pattern.Var] {
		if existing, ok := subs[pattern.Var]; ok {
			if ir.EqualsType(existing, target) {
				return subs, true
			}
			return nil, false
		}
		newSubs := make(map[string]ir.IrType)
		for k, v := range subs {
			newSubs[k] = v
		}
		newSubs[pattern.Var] = target
		return newSubs, true
	}

	if pattern.Case != target.Case {
		return nil, false
	}

	switch pattern.Case {
	case ir.NameType:
		if pattern.Name == target.Name {
			return subs, true
		}
		return nil, false

	case ir.AppType:
		subs, ok := matchType(pattern.App.Fun, target.App.Fun, vars, subs)
		if !ok {
			return nil, false
		}
		return matchType(pattern.App.Arg, target.App.Arg, vars, subs)

	case ir.TupleType:
		if len(pattern.Tuple.Elems) != len(target.Tuple.Elems) {
			return nil, false
		}
		currSubs := subs
		var ok bool
		for i := range pattern.Tuple.Elems {
			currSubs, ok = matchType(pattern.Tuple.Elems[i], target.Tuple.Elems[i], vars, currSubs)
			if !ok {
				return nil, false
			}
		}
		return currSubs, true

	case ir.FunType:
		subs, ok := matchType(pattern.Fun.Arg, target.Fun.Arg, vars, subs)
		if !ok {
			return nil, false
		}
		return matchType(pattern.Fun.Ret, target.Fun.Ret, vars, subs)

	default:
		if ir.EqualsType(pattern, target) {
			return subs, true
		}
		return nil, false
	}
}

func (c Context) containsTraitImpl(traitType ir.IrType, typeName ir.IrType) bool {
	_, ok := c.lookupTraitImplBind(traitType, typeName)
	return ok
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

func (c Context) containsTypeVarBindInScope(tvar string) bool {
	_, ok := c.lookupTypeVarBindInScope(tvar)
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
		if c, err = c.AddBind(NewTypeVarBind(tvar)); err != nil {
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

func (c Context) AddFreshType(typ ir.IrType) (Context, ir.VarKind, ir.IrType, error) {
	if !typ.Is(ir.ForallType) && !typ.Is(ir.LambdaType) {
		return c, ir.VarKind{}, typ, nil
	}

	var renamed ir.IrType
	var err error
	c, renamed, err = renameTypeVars(c, typ)
	if err != nil {
		return c, ir.VarKind{}, ir.IrType{}, err
	}

	switch renamed.Case {
	case ir.ForallType:
		vk := ir.VarKind{Var: renamed.Forall.Var, Kind: renamed.Forall.Kind, Bounds: renamed.Forall.Bounds}
		newContext, err := c.AddBind(NewTypeVarBind(vk))
		if err != nil {
			return c, ir.VarKind{}, ir.IrType{}, err
		}
		return newContext, vk, renamed.Forall.Type, nil

	case ir.LambdaType:
		vk := ir.VarKind{Var: renamed.Lambda.Var, Kind: renamed.Lambda.Kind}
		newContext, err := c.AddBind(NewTypeVarBind(vk))
		if err != nil {
			return c, ir.VarKind{}, ir.IrType{}, err
		}
		return newContext, vk, renamed.Lambda.Type, nil

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
	case ir.TraitDecl:
		c, err = c.AddBind(NewTraitBind(decl.Trait.ID, decl.Trait.TypeParams, decl.Trait.Methods))
		if err != nil {
			return c, err
		}
		for _, m := range decl.Trait.Methods {
			var args []ir.IrType
			for _, arg := range m.Args {
				t := ir.SubstituteType(arg.Type, ir.NewNameType("Self"), ir.NewVarType("Self"))
				args = append(args, t)
			}
			ret := ir.SubstituteType(m.RetType, ir.NewNameType("Self"), ir.NewVarType("Self"))
			methodType := ir.NewFunctionType(ir.NewTupleType(args), ret)
			for i := len(decl.Trait.TypeParams) - 1; i >= 0; i-- {
				tp := decl.Trait.TypeParams[i]
				methodType = ir.Forall(tp.Var, tp.Kind, methodType)
			}
			methodType = ir.Forall("Self", ir.NewTypeKind(), methodType)
			c, err = c.AddBind(NewTermDeclBind(decl.Trait.ID+ "::" + m.ID, methodType))
			if err != nil {
				return c, err
			}
		}
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

