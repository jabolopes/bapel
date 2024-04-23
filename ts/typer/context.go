package typer

import (
	"fmt"
	"strings"

	"github.com/jabolopes/bapel/ts/list"
)

type Context struct {
	list  list.List[Bind]
	idgen int
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

func (c Context) lookupExistVar(evar existVar) (int, bool) {
	for it := c.list.Iterate(); ; {
		i, bind, ok := it.Next()
		if !ok {
			break
		}

		if bind.Is(TypeBind) && bind.Type.Type.Is(ExistVarType) && *bind.Type.Type.ExistVar == evar {
			return i, true
		}
	}

	return 0, false
}

func (c Context) lookupTypeVar(tvar typeVar) (int, bool) {
	for it := c.list.Iterate(); ; {
		i, bind, ok := it.Next()
		if !ok {
			break
		}

		if bind.Is(TypeBind) && bind.Type.Type.Is(VarType) && *bind.Type.Type.Var == tvar {
			return i, true
		}
	}

	return 0, false
}

func (c Context) lookupTermBind(id string) (Bind, bool) {
	for it := c.list.Iterate(); ; {
		_, bind, ok := it.Next()
		if !ok {
			break
		}

		if bind.Case == TermBind && bind.Term.ID == id {
			return bind, true
		}
	}

	return Bind{}, false
}

func (c Context) ContainsExistVarType(existVar existVar) bool {
	_, ok := c.lookupExistVar(existVar)
	return ok
}

func (c Context) ContainsExistVarTypesInOrder(existVar1, existVar2 existVar) bool {
	index1, ok1 := c.lookupExistVar(existVar1)
	index2, ok2 := c.lookupExistVar(existVar2)
	return ok1 && ok2 && index1 < index2
}

func (c Context) ContainsVarsInOrder(typeVar typeVar, existVar existVar) bool {
	index1, ok1 := c.lookupTypeVar(typeVar)
	index2, ok2 := c.lookupExistVar(existVar)
	return ok1 && ok2 && index1 < index2
}

func (c Context) ContainsNameType(name name) bool {
	for it := c.list.Iterate(); ; {
		_, bind, ok := it.Next()
		if !ok {
			break
		}

		if bind.Case != TypeBind {
			continue
		}

		typ := bind.Type.Type
		if typ.Case != NameType {
			continue
		}

		if *typ.Name == name {
			return true
		}
	}

	return false
}

func (c Context) ContainsVarType(tvar typeVar) bool {
	_, ok := c.lookupTypeVar(tvar)
	return ok
}

func (c Context) ContainsTermBind(id string) bool {
	_, ok := c.lookupTermBind(id)
	return ok
}

func (c Context) GetTermType(id string) Type {
	bind, ok := c.lookupTermBind(id)
	if !ok {
		panic(fmt.Errorf("%s is undefined", id))
	}
	return bind.Term.Type
}

func (c Context) GenFreshID() (Context, string) {
	id := fmt.Sprintf("%c", 97+c.idgen)
	c.idgen++
	return c, id
}

func (c Context) GenFreshExistVarType() (Context, Type) {
	c, id := c.GenFreshID()
	return c, NewExistVarType(id)
}

func (c Context) GenFreshVarType() (Context, Type) {
	c, id := c.GenFreshID()
	return c, NewVarType(id)
}

func (c Context) AddType(typ Type) Context {
	c.list = c.list.Add(NewTypeBind(typ))
	return c
}

func (c Context) AddFreshExistVarType() (Context, Type) {
	c, typ := c.GenFreshExistVarType()
	return c.AddType(typ), typ
}

func (c Context) AddFreshVarType() (Context, Type) {
	c, typ := c.GenFreshVarType()
	return c.AddType(typ), typ
}

func (c Context) AddJudge(judge Judge) Context {
	c.list = c.list.Add(NewJudgeBind(judge))
	return c
}

func (c Context) Pop() (Bind, Context) {
	bind, ok := c.list.Value()
	if !ok {
		panic("Context is empty")
	}

	c.list = c.list.Remove()
	return bind, c
}

// ReplaceExistVar replaces an existential variable contained in the context
// with the given types. It does not recurse into the binds, judgements, etc. If
// `types` is nil, the existential variable is simply removed from the context.
func (c Context) ReplaceExistVar(existVar existVar, replacements []Type) Context {
	index, ok := c.lookupExistVar(existVar)
	if !ok {
		return c
	}

	var replacementBinds []Bind
	if len(replacements) > 0 {
		replacementBinds = make([]Bind, len(replacements))
		for i := range replacements {
			replacementBinds[i] = NewTypeBind(replacements[i])
		}
	}

	{
		binds := c.list.Iterate().Collect()
		binds = append(binds[:index], append(replacementBinds, binds[index+1:]...)...)
		c.list = list.FromSlice(binds)
	}

	return c
}
