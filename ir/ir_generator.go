package ir

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/jabolopes/bapel/vm"
	"github.com/zyedidia/generic/stack"
)

type blockType int

const (
	moduleBlock = blockType(iota)
	functionBlock
	argsBlock
	retsBlock
	localsBlock
	ifThenBlock
	ifElseBlock
	elseBlock
)

type IrGenerator struct {
	generators     *stack.Stack[*ByteGenerator]
	blocks         *stack.Stack[blockType]
	functions      []IrFunction
	optable        vm.OpTable
	mainCallOffset uint64 // Callsite offset of the main function to be fixed when the module is closed.
}

func (a *IrGenerator) gen() *ByteGenerator {
	return a.generators.Peek()
}

func (a *IrGenerator) fun() *IrFunction {
	return &a.functions[len(a.functions)-1]
}

func (a *IrGenerator) endModule() error {
	if a.blocks.Pop() != moduleBlock {
		return errors.New("expected module block")
	}

	mainFunction, err := a.LookupFunction("main")
	if err != nil {
		return err
	}

	// Overwrite 'main' call site.
	binary.LittleEndian.PutUint64(a.gen().Data()[a.mainCallOffset:], mainFunction.offset)

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
	// TODO: Validate there's a current ongoing function.

	if a.blocks.Pop() != localsBlock {
		return errors.New("expected locals block")
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
		a.gen().PutOpCode(vm.IfThen)
	} else {
		a.gen().PutOpCode(vm.IfElse)
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
		PutOpCode(vm.Else).
		PutI64(uint64(len(elseGen.Data())))

	ifGen := a.generators.Pop()

	a.gen().
		PutOpCode(vm.IfThen).
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

func (a *IrGenerator) StackAlloc(size uint16) error {
	a.gen().
		PutOpCode(vm.StackAlloc).
		PutI16(size)
	return nil
}

func (a *IrGenerator) Module() error {
	if a.blocks.Size() != 0 {
		return fmt.Errorf("Modules can only be defined at the toplevel")
	}

	a.gen().PutOpCode(a.optable.Call())

	// Write placeholder operand for the address of the main function
	// which is not yet defined. Later, when the module is fully
	// defined, we will come back and overwrite this operand with the
	// correct address.
	{
		a.mainCallOffset = uint64(a.gen().Len())
		a.gen().PutI64(0)
	}

	a.gen().PutOpCode(vm.Halt)
	return nil
}

func (a *IrGenerator) Function(id string) error {
	if a.blocks.Peek() != moduleBlock {
		return fmt.Errorf("Can only be used within a module block")
	}

	a.blocks.Push(functionBlock)

	a.functions = append(a.functions, IrFunction{
		id,
		[]IrVar{}, /* vars */
		irFrame{}, /* frame */
		uint64(a.gen().Len()),
	})

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
	// TODO: Validate there's a current ongoing function.

	a.blocks.Push(argsBlock)
	return nil
}

func (a *IrGenerator) Rets() error {
	// TODO: Validate there's a current ongoing function.

	a.blocks.Push(retsBlock)
	return nil
}

func (a *IrGenerator) Locals() error {
	// TODO: Validate there's a current ongoing function.

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

	function, err := a.LookupFunction(id)
	if err != nil {
		return err
	}

	a.gen().
		PutOpCode(a.optable.Call()).
		PutI64(function.offset)
	return nil
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

	switch irvar.Type {
	case I8:
		a.gen().PutOpCode(vm.PopVarI8)
	case I16:
		a.gen().PutOpCode(vm.PopVarI16)
	case I32:
		a.gen().PutOpCode(vm.PopVarI32)
	case I64:
		a.gen().PutOpCode(vm.PopVarI64)
	default:
		return fmt.Errorf("Unhandled IR type %d", irvar.Type)
	}

	a.gen().PutI16(irvar.offset)
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
		[]IrFunction{},
		vm.NewOpTable(),
		0, /* mainCallOffset */
	}
	generator.generators.Push(NewByteGenerator())
	return generator
}
