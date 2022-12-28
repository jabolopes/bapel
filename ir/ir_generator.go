package ir

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/jabolopes/bapel/vm"
	"github.com/zyedidia/generic/stack"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type blockType int

const (
	moduleBlock = blockType(iota)
	declsBlock
	functionBlock
	argsBlock
	retsBlock
	localsBlock
	ifThenBlock
	ifElseBlock
	elseBlock
)

type IrGenerator struct {
	generators *stack.Stack[*ByteGenerator]
	blocks     *stack.Stack[blockType]
	decls      []irDecl
	functions  []IrFunction
	optable    vm.OpTable
	callsites  map[string]irCallsite // Callsites indexed by function name.
}

func (a *IrGenerator) gen() *ByteGenerator {
	return a.generators.Peek()
}

func (a *IrGenerator) fun() *IrFunction {
	return &a.functions[len(a.functions)-1]
}

func (a *IrGenerator) lookupDecl(id string) (irDecl, bool) {
	for _, d := range a.decls {
		if d.id == id {
			return d, true
		}
	}

	return irDecl{}, false
}

func (a *IrGenerator) callInternal(id string) error {
	function, err := a.LookupFunction(id)
	if err == nil {
		// Make regular call.
		a.gen().
			PutOpCode(a.optable.Call()).
			PutI64(function.offset)
		return nil
	}

	if _, ok := a.lookupDecl(id); ok {
		// Make call with placeholder address.
		a.gen().PutOpCode(a.optable.Call())
		callsiteOffset := uint64(a.gen().Len())
		a.gen().PutI64(0)

		// Record callsite to be fixed later.
		callsite := a.callsites[id]
		callsite.offsets = append(callsite.offsets, callsiteOffset)
		a.callsites[id] = callsite
		return nil
	}

	return err
}

func (a *IrGenerator) endModule() error {
	if a.blocks.Pop() != moduleBlock {
		return errors.New("expected module block")
	}

	{
		// Check there are no undefined declarations.
		for _, decl := range a.decls {
			if _, err := a.LookupFunction(decl.id); err != nil {
				return fmt.Errorf("Symbol %q is declared but it is not defined", decl.id)
			}
		}
	}

	{
		// Check there are no unresolved callsites.
		if len(a.callsites) > 0 {
			return fmt.Errorf("There are unresolved callsites for symbols %v", maps.Keys(a.callsites))
		}
	}

	return nil
}

func (a *IrGenerator) endDecls() error {
	if a.blocks.Pop() != declsBlock {
		return errors.New("expected declarations block")
	}
	return nil
}

func (a *IrGenerator) endFunction() error {
	if a.blocks.Peek() != functionBlock {
		return errors.New("expected function block")
	}
	defer a.blocks.Pop()

	return a.Return()
}

func (a *IrGenerator) endArgs() error {
	if a.blocks.Pop() != argsBlock {
		return errors.New("expected args block")
	}
	return nil
}

func (a *IrGenerator) endRets() error {
	if a.blocks.Pop() != retsBlock {
		return errors.New("expected rets block")
	}

	return nil
}

func (a *IrGenerator) endLocals() error {
	if a.blocks.Pop() != localsBlock {
		return errors.New("expected locals block")
	}

	{
		// Check function definition matches declaration (if any).
		decl, ok := a.lookupDecl(a.fun().id)
		if ok {
			var argTypes []IrType
			var retTypes []IrType
			for _, irvar := range a.fun().Vars() {
				if irvar.VarType == ArgVar {
					argTypes = append(argTypes, irvar.Type)
				} else if irvar.VarType == RetVar {
					retTypes = append(retTypes, irvar.Type)
				}
			}

			if !slices.Equal(decl.args, argTypes) || !slices.Equal(decl.rets, retTypes) {
				return fmt.Errorf("definition of function %q does not match its declaration type", a.fun().id)
			}
		}
	}

	if err := a.fun().computeFrame(); err != nil {
		return err
	}

	a.gen().PutI16(a.fun().frame.enterSize())

	fmt.Printf("DEBUG function %s %d %v\n", a.fun().id, a.fun().offset, a.fun().frame)
	return nil
}

func (a *IrGenerator) endIf() error {
	block := a.blocks.Pop()
	if block != ifThenBlock && block != ifElseBlock {
		return errors.New("expected if block")
	}

	nested := a.generators.Pop()

	if block == ifThenBlock {
		a.gen().PutOpCode(a.optable.IfThen())
	} else {
		a.gen().PutOpCode(a.optable.IfElse())
	}

	a.gen().
		PutI64(uint64(len(nested.Data()))).
		append(nested.Data())
	return nil
}

func (a *IrGenerator) endElse() error {
	if a.blocks.Pop() != elseBlock {
		return errors.New("expected else block")
	}

	elseGen := a.generators.Pop()

	// Finish the 'if' section by adding the 'else' (aka jump) instruction. This
	// is important so that the length of the 'if' section is correct when jumping
	// to the else branch.
	a.gen().
		PutOpCode(a.optable.Else()).
		PutI64(uint64(len(elseGen.Data())))

	ifGen := a.generators.Pop()

	a.gen().
		PutOpCode(a.optable.IfThen()).
		PutI64(uint64(len(ifGen.Data()))).
		append(ifGen.Data()).
		append(elseGen.Data())

	return nil
}

func (a *IrGenerator) putImmediate(typ IrType, value uint64) error {
	// TODO: Validate whether typecast truncates the value and return an
	// error in that case.

	switch typ {
	case I8:
		a.gen().PutI8(byte(value))
	case I16:
		a.gen().PutI16(uint16(value))
	case I32:
		a.gen().PutI32(uint32(value))
	case I64:
		a.gen().PutI64(value)
	default:
		return fmt.Errorf("Unhandled optype %d", typ)
	}

	return nil
}

func (a *IrGenerator) Module() error {
	if a.blocks.Size() != 0 {
		return fmt.Errorf("Modules can only be defined at the toplevel")
	}

	a.decls = append(a.decls, irDecl{"main", nil, nil})
	if err := a.callInternal("main"); err != nil {
		return err
	}

	a.gen().PutOpCode(a.optable.Halt())
	return nil
}

func (a *IrGenerator) Decls() error {
	if a.blocks.Peek() != moduleBlock {
		return fmt.Errorf("Can only start a 'decls' block within a module block")
	}
	a.blocks.Push(declsBlock)
	return nil
}

func (a *IrGenerator) Declare(id string, args []IrType, rets []IrType) error {
	fmt.Printf("HERE DECL %q %v %v\n", id, args, rets)

	if _, ok := a.lookupDecl(id); ok {
		return fmt.Errorf("Symbol %q is already declared in this module", id)
	}

	a.decls = append(a.decls, irDecl{id, args, rets})
	return nil
}

func (a *IrGenerator) Function(id string) error {
	if a.blocks.Peek() != moduleBlock {
		return fmt.Errorf("Can only be used within a module block")
	}

	a.blocks.Push(functionBlock)

	a.functions = append(a.functions, IrFunction{
		id,
		[]IrVar{},             /* vars */
		irFrame{},             /* frame */
		uint64(a.gen().Len()), /* offset */
	})

	{
		// Resolve callsites (if any).
		callsite, ok := a.callsites[id]
		if ok {
			fmt.Printf("DEBUG LINK %s %v = %d\n", id, callsite, a.fun().offset)

			for _, offset := range callsite.offsets {
				binary.LittleEndian.PutUint64(a.gen().Data()[offset:], a.fun().offset)
			}
			delete(a.callsites, id)
		}
	}

	return nil
}

func (a *IrGenerator) LookupFunction(id string) (IrFunction, error) {
	for _, f := range a.functions {
		if f.id == id {
			return f, nil
		}
	}

	return IrFunction{}, fmt.Errorf("Undefined function %q", id)
}

func (a *IrGenerator) Args() error {
	if a.blocks.Peek() != functionBlock {
		return fmt.Errorf("Can only start an 'args' block within a function block")
	}
	a.blocks.Push(argsBlock)
	return nil
}

func (a *IrGenerator) Rets() error {
	if a.blocks.Peek() != functionBlock {
		return fmt.Errorf("Can only start a 'rets' block within a function block")
	}
	a.blocks.Push(retsBlock)
	return nil
}

func (a *IrGenerator) Locals() error {
	if a.blocks.Peek() != functionBlock {
		return fmt.Errorf("Can only start a 'locals' block within a function block")
	}
	a.blocks.Push(localsBlock)
	return nil
}

func (a *IrGenerator) DefineVar(id string, typ IrType) error {
	var vartype IrVarType
	switch block := a.blocks.Peek(); block {
	case argsBlock:
		vartype = ArgVar
	case retsBlock:
		vartype = RetVar
	case localsBlock:
		vartype = LocalVar
	default:
		return fmt.Errorf("Cannot declare variable inside block type %d", block)
	}

	return a.fun().addVar(id, IrVar{id, vartype, typ, 0 /* offset */})
}

func (a *IrGenerator) LookupVar(id string) (IrVar, error) {
	return a.fun().lookupVar(id)
}

func (a *IrGenerator) Call(id string) error {
	if a.blocks.Peek() != functionBlock {
		return fmt.Errorf("Can only be used within a function block")
	}
	return a.callInternal(id)
}

func (a *IrGenerator) Return() error {
	if a.blocks.Peek() != functionBlock {
		return fmt.Errorf("Can only be used within a function block")
	}

	a.gen().
		PutOpCode(a.optable.Return()).
		PutI16(a.fun().frame.leaveSize())
	return nil
}

func (a *IrGenerator) IfThen() error {
	// TODO: Validate there's a current ongoing function.

	a.generators.Push(NewByteGenerator())
	a.blocks.Push(ifThenBlock)
	return nil
}

func (a *IrGenerator) IfElse() error {
	// TODO: Validate there's a current ongoing function.

	a.generators.Push(NewByteGenerator())
	a.blocks.Push(ifElseBlock)
	return nil
}

func (a *IrGenerator) Else() error {
	if a.blocks.Pop() != ifThenBlock {
		return errors.New("expected if block")
	}

	a.generators.Push(NewByteGenerator())
	a.blocks.Push(elseBlock)
	return nil
}

func (a *IrGenerator) End() error {
	switch block := a.blocks.Peek(); block {
	case moduleBlock:
		return a.endModule()
	case declsBlock:
		return a.endDecls()
	case functionBlock:
		return a.endFunction()
	case argsBlock:
		return a.endArgs()
	case retsBlock:
		return a.endRets()
	case localsBlock:
		return a.endLocals()
	case ifThenBlock, ifElseBlock:
		return a.endIf()
	case elseBlock:
		return a.endElse()
	default:
		return fmt.Errorf("Unexpected block type %d", block)
	}
}

func (a *IrGenerator) PushImmediate(typ IrType, value uint64) error {
	a.gen().PutOpCode(a.optable.Push(vm.ImmediateMode, typ))
	return a.putImmediate(typ, value)
}

func (a *IrGenerator) PushVar(id string) error {
	if a.blocks.Peek() != functionBlock {
		return fmt.Errorf("Can only be used within a function block")
	}

	// TODO: Validate there's a current ongoing function.

	irvar, err := a.fun().lookupVar(id)
	if err != nil {
		return err
	}

	a.gen().
		PutOpCode(a.optable.Push(vm.VarMode, irvar.Type)).
		PutI16(irvar.offset)
	return nil
}

func (a *IrGenerator) PopVar(id string) error {
	// TODO: Validate there's a current ongoing function.

	irvar, err := a.fun().lookupVar(id)
	if err != nil {
		return err
	}

	a.gen().
		PutOpCode(a.optable.PopVar(irvar.Type)).
		PutI16(irvar.offset)
	return nil
}

func (a *IrGenerator) Add(typ IrType) error {
	a.gen().PutOpCode(a.optable.Add(vm.StackMode, vm.StackMode, typ))
	return nil
}

func (a *IrGenerator) PrintImmediate(typ IrType, value uint64) error {
	a.gen().PutOpCode(a.optable.Print(vm.ImmediateMode, typ))
	return a.putImmediate(typ, value)
}

func (a *IrGenerator) PrintVar(id string) error {
	// TODO: Validate there's a current ongoing function.

	irvar, err := a.fun().lookupVar(id)
	if err != nil {
		return err
	}

	a.gen().
		PutOpCode(a.optable.Print(vm.VarMode, irvar.Type)).
		PutI16(irvar.offset)
	return nil
}

func (a *IrGenerator) PrintStack(typ IrType) error {
	a.gen().PutOpCode(a.optable.Print(vm.StackMode, typ))
	return nil
}

func (a *IrGenerator) Program() vm.OpProgram {
	return vm.OpProgram{
		a.gen().Data(),
		[]vm.OpFunction{},
	}
}

func New() *IrGenerator {
	generator := &IrGenerator{
		stack.New[*ByteGenerator](), /* generators */
		stack.New[blockType](),      /* blocks */
		[]irDecl{},                  /* decls */
		[]IrFunction{},              /* functions */
		vm.NewOpTable(),
		map[string]irCallsite{}, /* callsites */
	}
	generator.generators.Push(NewByteGenerator())
	return generator
}
