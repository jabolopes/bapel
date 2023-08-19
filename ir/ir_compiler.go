package ir

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/jabolopes/bapel/parser"
	"github.com/zyedidia/generic/stack"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type FindCase int

const (
	FindAny = FindCase(iota)
	FindDeclOnly
	FindDefOnly
	FindVarOnly
)

func toID(id string) string {
	return strings.Replace(id, ".", "::", -1)
}

type Compiler struct {
	output     io.Writer
	blocks     *stack.Stack[blockType]
	imports    []irDecl
	exports    []irDecl
	decls      []irDecl
	structDefs []irDecl
	functions  []irFunction
	optable    OpTable
	callsites  map[string]irCallsite // Callsites indexed by function name.
}

func (a *Compiler) out() io.Writer {
	return a.output
}

func (a *Compiler) fun() *irFunction {
	return &a.functions[len(a.functions)-1]
}

func (a *Compiler) isFunctionBlock() bool {
	allowedBlocks := []blockType{functionBlock, ifThenBlock, ifElseBlock, elseBlock}

	for _, allowed := range allowedBlocks {
		if a.blocks.Peek() == allowed {
			return true
		}
	}

	return false
}

func (a *Compiler) lookupDecl(id string, findCase FindCase) (irDecl, bool) {
	if findCase == FindAny || findCase == FindVarOnly {
		if irvar, err := a.LookupVar(id); err == nil {
			return irvar.decl(), true
		}
	}

	if findCase == FindAny || findCase == FindDefOnly {
		for _, d := range a.structDefs {
			if d.id == id {
				return d, true
			}
		}

		for _, f := range a.functions {
			if f.id == id {
				return f.decl(), true
			}
		}
	}

	if findCase == FindAny || findCase == FindDeclOnly {
		for _, d := range a.decls {
			if d.id == id {
				return d, true
			}
		}

		for _, d := range a.exports {
			if d.id == id {
				return d, true
			}
		}

		for _, d := range a.imports {
			if d.id == id {
				return d, true
			}
		}
	}

	return irDecl{}, false
}

func (a *Compiler) lookupFunction(id string) (irFunction, error) {
	for _, f := range a.functions {
		if f.id == id {
			return f, nil
		}
	}

	return irFunction{}, fmt.Errorf("Undefined function %q", id)
}

func (a *Compiler) printType(typ IrType) {
	switch typ.Case {
	case ArrayType:
		fmt.Fprintf(a.out(), "std::array<")
		a.printType(typ.ArrayType.ElemType)
		fmt.Fprintf(a.out(), ", %d>", typ.ArrayType.Size)
	case FunType:
		panic(fmt.Errorf("printType: Uniplemented function type"))
	case IntType:
		switch typ.IntType {
		case I8:
			fmt.Fprintf(a.out(), "char")
		case I16:
			fmt.Fprintf(a.out(), "int16_t")
		case I32:
			fmt.Fprintf(a.out(), "int32_t")
		case I64:
			fmt.Fprintf(a.out(), "int64_t")
		}
	default:
		panic(fmt.Errorf("printType: Unhandled case %d", typ.Case))
	}
}

func (a *Compiler) printDecl(decl irDecl) {
	if decl.typ.Is(IntType) {
		a.printType(decl.typ)
		fmt.Fprintf(a.out(), " %s", decl.id)
		return
	}

	typ := decl.typ.FunType

	// Print rets.
	switch len(typ.Rets) {
	case 0:
		fmt.Fprintf(a.out(), "void")
	case 1:
		a.printType(typ.Rets[0])
	default:
		fmt.Fprintf(a.out(), "std::tuple<")
		a.printType(typ.Rets[0])
		for _, ret := range typ.Rets[1:] {
			fmt.Fprintf(a.out(), ", ")
			a.printType(ret)
		}
		fmt.Fprintf(a.out(), ">")
	}

	// Print id.
	fmt.Fprintf(a.out(), " %s(", decl.id)

	// Print args.
	switch len(typ.Args) {
	case 0:
		break
	case 1:
		a.printType(typ.Args[0])
	default:
		a.printType(typ.Args[0])
		for _, arg := range typ.Args[1:] {
			fmt.Fprintf(a.out(), ", ")
			a.printType(arg)
		}
	}

	fmt.Fprintf(a.out(), ")")
}

func (a *Compiler) printFunctionSignature(function irFunction) {
	id := function.id

	for _, d := range a.exports {
		if d.id == id {
			fmt.Fprintf(a.out(), "export ")
			break
		}
	}

	if strings.Contains(id, ".") {
		fmt.Fprintf(a.out(), "namespace ")

		tokens := strings.Split(id, ".")
		tokens, id = tokens[:len(tokens)-1], tokens[len(tokens)-1]

		fmt.Fprintf(a.out(), "%s", tokens[0])
		for _, token := range tokens[1:] {
			fmt.Fprintf(a.out(), "::%s", token)
		}

		fmt.Fprintf(a.out(), "{")
	}

	// Print rets.
	switch rets := function.rets(); len(rets) {
	case 0:
		fmt.Fprintf(a.out(), "void")
	case 1:
		a.printType(rets[0].Type)
	default:
		fmt.Fprintf(a.out(), "std::tuple<")
		a.printType(rets[0].Type)
		for _, ret := range rets[1:] {
			fmt.Fprintf(a.out(), ", ")
			a.printType(ret.Type)
		}
		fmt.Fprintf(a.out(), ">")
	}

	// Print id.
	fmt.Fprintf(a.out(), " %s(", id)

	// Print args.
	switch args := function.args(); len(args) {
	case 0:
		break
	case 1:
		a.printType(args[0].Type)
		fmt.Fprintf(a.out(), " %s", args[0].Id)
	default:
		a.printType(args[0].Type)
		fmt.Fprintf(a.out(), " %s", args[0].Id)
		for _, arg := range args[1:] {
			fmt.Fprintf(a.out(), ", ")
			a.printType(arg.Type)
			fmt.Fprintf(a.out(), " %s", arg.Id)
		}
	}

	fmt.Fprintf(a.out(), ")")
}

func (a *Compiler) callImpl(id string, args []parser.Token, rets []string) error {
	// Get function type.
	formalDecl, ok := a.lookupDecl(id, FindAny)
	if !ok {
		return fmt.Errorf("Undefined function %q", id)
	}

	// Compute type at callsite.
	//
	// TODO: Improve since code assumes function types and int types.
	actualType := IrFunctionType{}
	{
		for i, arg := range args {
			switch arg.Case {
			case parser.IDToken:
				decl, err := a.LookupDecl(arg.Text, FindVarOnly)
				if err != nil {
					return err
				}
				actualType.Args = append(actualType.Args, decl.typ)
			case parser.NumberToken:
				typ := NewIntType(I64)
				if i < len(formalDecl.typ.FunType.Args) {
					typ = formalDecl.typ.FunType.Args[i]
				}
				actualType.Args = append(actualType.Args, typ)
			}
		}

		for _, ret := range rets {
			decl, err := a.LookupDecl(ret, FindVarOnly)
			if err != nil {
				return err
			}
			actualType.Rets = append(actualType.Rets, decl.typ)
		}
	}

	// Check whether actual decl matches the formal decl.
	actualDecl := NewDecl(id, NewFunctionType(actualType))
	if err := matchesDecl(formalDecl, actualDecl); err != nil {
		return err
	}

	switch len(rets) {
	case 0:
		break
	case 1:
		fmt.Fprintf(a.out(), "%s = ", rets[0])
	default:
		fmt.Fprintf(a.out(), "std::tie(%s", rets[0])
		for _, ret := range rets[1:] {
			fmt.Fprintf(a.out(), ", %s", ret)
		}
		fmt.Fprintf(a.out(), ") = ")
	}

	fmt.Fprintf(a.out(), "%s(", toID(id))

	switch len(args) {
	case 0:
		break
	case 1:
		fmt.Fprintf(a.out(), "%s", args[0].Text)
	default:
		fmt.Fprintf(a.out(), "%s", args[0].Text)
		for _, arg := range args[1:] {
			fmt.Fprintf(a.out(), ", %s", arg.Text)
		}
	}

	fmt.Fprintf(a.out(), ")")
	return nil
}

func (a *Compiler) endModule() error {
	if a.blocks.Pop() != moduleBlock {
		return errors.New("expected module block")
	}

	{
		// Check there are no undefined declarations.
		for _, decl := range a.decls {
			if _, ok := a.lookupDecl(decl.id, FindDefOnly); !ok {
				return fmt.Errorf("Symbol %q is declared but it is not defined", decl.id)
			}
		}
	}

	{
		// Check there are no undefined exports.
		for _, decl := range a.exports {
			if _, ok := a.lookupDecl(decl.id, FindDefOnly); !ok {
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

func (a *Compiler) endImports() error {
	if a.blocks.Pop() != importsBlock {
		return errors.New("expected imports block")
	}
	fmt.Fprintln(a.out())
	return nil
}

func (a *Compiler) endExports() error {
	if a.blocks.Pop() != exportsBlock {
		return errors.New("expected exports block")
	}
	fmt.Fprintln(a.out())
	return nil
}

func (a *Compiler) endDecls() error {
	if a.blocks.Pop() != declsBlock {
		return errors.New("expected declarations block")
	}
	fmt.Fprintln(a.out())
	return nil
}

func (a *Compiler) endFunction() error {
	if a.blocks.Peek() != functionBlock {
		return errors.New("expected function block")
	}

	if err := a.Return(); err != nil {
		return err
	}

	fmt.Fprintf(a.out(), "}\n")

	if strings.Contains(a.fun().id, ".") {
		// This function was defined inside a namespace, so close the namespace.
		fmt.Fprintf(a.out(), "}\n")
	}

	a.blocks.Pop()
	return nil
}

func (a *Compiler) endIf() error {
	if block := a.blocks.Pop(); block != ifThenBlock && block != ifElseBlock {
		return errors.New("expected if block")
	}

	fmt.Fprintf(a.out(), "}\n")
	return nil
}

func (a *Compiler) endElse() error {
	if a.blocks.Pop() != elseBlock {
		return errors.New("expected else block")
	}

	fmt.Fprintf(a.out(), "}\n")
	return nil
}

func (a *Compiler) Module() error {
	if a.blocks.Size() != 0 {
		return fmt.Errorf("Modules can only be defined at the toplevel")
	}

	fmt.Fprintf(a.out(), "module;\n")
	fmt.Fprintf(a.out(), "\n")
	fmt.Fprintf(a.out(), "import <array>;\n")
	fmt.Fprintf(a.out(), "import <cstdlib>;\n")
	fmt.Fprintf(a.out(), "import <iostream>;\n")
	fmt.Fprintf(a.out(), "import <tuple>;\n")
	fmt.Fprintf(a.out(), "import <vector>;\n")
	fmt.Fprintf(a.out(), "\n")
	fmt.Fprintf(a.out(), "import c;\n")
	fmt.Fprintf(a.out(), "\n")
	fmt.Fprintf(a.out(), "export module bpl;\n")
	fmt.Fprintf(a.out(), "\n")

	return nil
}

func (a *Compiler) Imports() error {
	if a.blocks.Peek() != moduleBlock {
		return fmt.Errorf("can only start a 'imports' block within a module block")
	}
	a.blocks.Push(importsBlock)
	fmt.Fprintf(a.out(), "// IMPORTS\n")
	return nil
}

func (a *Compiler) Exports() error {
	if a.blocks.Peek() != moduleBlock {
		return fmt.Errorf("can only start a 'exports' block within a module block")
	}
	a.blocks.Push(exportsBlock)
	return nil
}

func (a *Compiler) Decls() error {
	if a.blocks.Peek() != moduleBlock {
		return fmt.Errorf("can only start a 'decls' block within a module block")
	}
	a.blocks.Push(declsBlock)

	fmt.Fprintf(a.out(), "// HEADER\n")
	return nil
}

func (a *Compiler) Declare(decl irDecl) error {
	if block := a.blocks.Peek(); block != importsBlock && block != exportsBlock && block != declsBlock {
		return fmt.Errorf("declarations can occur only within an 'imports', an 'exports', or a 'decls' block.")
	}

	if _, ok := a.lookupDecl(decl.id, FindAny); ok {
		return fmt.Errorf("Symbol %q is already declared in this module", decl.id)
	}

	switch a.blocks.Peek() {
	case importsBlock:
		a.imports = append(a.imports, decl)
	case exportsBlock:
		a.exports = append(a.exports, decl)
	case declsBlock:
		a.decls = append(a.decls, decl)
		a.printDecl(decl)
		fmt.Fprintf(a.out(), ";\n")
	}

	return nil
}

func (a *Compiler) Function(id string, vars []IrVar) error {
	if a.blocks.Peek() != moduleBlock {
		return fmt.Errorf("can only be used within a module block")
	}

	if _, ok := a.lookupDecl(id, FindDefOnly); ok {
		return fmt.Errorf("symbol %q already defined", id)
	}

	{
		var args []IrVar
		var rets []IrVar
		for _, irvar := range vars {
			if irvar.VarType == ArgVar {
				args = append(args, irvar)
			} else if irvar.VarType == RetVar {
				rets = append(rets, irvar)
			} else {
				return fmt.Errorf("locals should be defined in the function body")
			}
		}

		vars = append(args, rets...)
	}

	function := irFunction{
		id,
		vars,
		irFrame{}, /* frame */
		0,         /* offset */
	}
	a.functions = append(a.functions, function)

	// Check function definition matches declaration (if any).
	if decl, ok := a.lookupDecl(a.fun().id, FindDeclOnly); ok {
		if err := matchesDecl(decl, a.fun().decl()); err != nil {
			return fmt.Errorf("definition of function %q does not match its declaration type: %w", a.fun().id, err)
		}
	}

	// Compute frame with args and rets.
	if err := a.fun().computeFrame(); err != nil {
		return err
	}

	a.blocks.Push(functionBlock)

	a.printFunctionSignature(function)
	fmt.Fprintf(a.out(), " {\n")

	for _, ret := range function.rets() {
		a.printType(ret.Type)
		fmt.Fprintf(a.out(), " %s;\n", ret.Id)
	}

	return nil
}

func (a *Compiler) Struct(id string, typ IrStructType) error {
	if a.blocks.Peek() != moduleBlock {
		return fmt.Errorf("can only be used within a module block")
	}

	if _, ok := a.lookupDecl(id, FindDefOnly); ok {
		return fmt.Errorf("symbol %q already defined", id)
	}

	actualDecl := irDecl{id, NewStructType(typ)}
	if formalDecl, ok := a.lookupDecl(id, FindDeclOnly); ok {
		if err := matchesDecl(formalDecl, actualDecl); err != nil {
			return fmt.Errorf("struct %q type %v does not match its declaration type %v", id, typ, formalDecl.typ)
		}
	}

	a.structDefs = append(a.structDefs, actualDecl)

	fmt.Fprintf(a.out(), "struct %s {\n", id)
	for _, field := range typ.Fields {
		a.printType(field.Type)
		fmt.Fprintf(a.out(), " %s;\n", field.Name)
	}
	fmt.Fprintf(a.out(), "};\n")

	return nil
}

func (a *Compiler) Entity(id string) error {
	if a.blocks.Peek() != moduleBlock {
		return fmt.Errorf("can only be used within a module block")
	}

	if _, ok := a.lookupDecl(id, FindDefOnly); !ok {
		return fmt.Errorf("entity %q must have a previously defined type (e.g., struct)", id)
	}

	fmt.Fprintf(a.out(), "ecs::StaticPool<int, struct %s, 1024> %s_entity;\n", id, id)
	return nil
}

func (a *Compiler) DefineLocal(decl irDecl) error {
	if !a.isFunctionBlock() {
		return fmt.Errorf("can only define local variables inside a function")
	}

	if err := a.fun().addVar(decl.id, IrVar{decl.id, LocalVar, decl.typ, 0 /* offset */}); err != nil {
		return err
	}

	a.printType(decl.typ)
	fmt.Fprintf(a.out(), " %s;\n", decl.id)
	return nil
}

func (a *Compiler) LookupDecl(id string, findCase FindCase) (irDecl, error) {
	if decl, ok := a.lookupDecl(id, findCase); ok {
		return decl, nil
	}

	return irDecl{}, fmt.Errorf("Undefined symbol %q", id)
}

func (a *Compiler) LookupVar(id string) (IrVar, error) {
	if len(a.functions) <= 0 {
		return IrVar{}, fmt.Errorf("Undefined variable %q", id)
	}

	return a.fun().lookupVar(id)
}

func (a *Compiler) Assign(args []parser.Token, rets []string) error {
	if !a.isFunctionBlock() {
		return errors.New("op 'call' can only be used in a function block")
	}

	if len(args) == 0 {
		return fmt.Errorf("expected at least 1 argument; got %q", args)
	}

	switch args[0].Text {
	case "call":
		// ret1 [ret2 ...] <- call funID [arg1 ...]
		//
		// Examples:
		//   x <- call f
		//   x y <- call f a b c
		args = args[1:]

		id, args, err := parser.Shift(args, fmt.Errorf("expected identifier as first token; got %v", args))
		if err != nil {
			return err
		}

		if id.Case != parser.IDToken {
			return fmt.Errorf("expected identifier as first token; got %v", id)
		}

		if err := a.callImpl(id.Text, args, rets); err != nil {
			return err
		}
		fmt.Fprintf(a.out(), ";\n")
		return nil

	case "array::get":
		// ret <- array::get array index
		//
		// Examples:
		//   x <- array::get myarray 10
		args = args[1:]

		if len(rets) != 1 {
			return fmt.Errorf("expected exactly 1 return variable; got %q", rets)
		}

		id, args, err := parser.Shift(args, fmt.Errorf("expected identifier as first token; got %v", args))
		if err != nil {
			return err
		}

		if id.Case != parser.IDToken {
			return fmt.Errorf("expected identifier as first token; got %v", id)
		}

		index, args, err := parser.Shift(args, fmt.Errorf("expected number as second token; got %v", args))
		if err != nil {
			return err
		}

		// TODO: Check types.
		fmt.Fprintf(a.out(), "%s = %s[%s];\n", rets[0], id.Text, index.Text)
		return nil

	case "array::set":
		// array::set array index value
		//
		// Examples:
		//   array::set myarray 10 myvalue
		args = args[1:]

		if len(rets) != 0 {
			return fmt.Errorf("expected no return variables; got %q", rets)
		}

		id, args, err := parser.Shift(args, fmt.Errorf("expected identifier as first token; got %v", args))
		if err != nil {
			return err
		}

		if id.Case != parser.IDToken {
			return fmt.Errorf("expected identifier as first token; got %v", id)
		}

		index, args, err := parser.Shift(args, fmt.Errorf("expected number as second token; got %v", args))
		if err != nil {
			return err
		}

		value, args, err := parser.Shift(args, fmt.Errorf("expected value as third argument; got %v", args))
		if err != nil {
			return err
		}

		// TODO: Check types.
		fmt.Fprintf(a.out(), "%s[%s] = %s;\n", id.Text, index.Text, value.Text)
		return nil

	case "widen":
		// x <- widen y
		args = args[1:]

		if len(rets) != 1 {
			return fmt.Errorf("expected exactly 1 return variable; got %q", rets)
		}

		if len(args) != 1 {
			return fmt.Errorf("expected exactly 1 argument variable; got %q", args)
		}

		arg := args[0]
		ret := rets[0]

		retDecl, err := a.LookupDecl(ret, FindVarOnly)
		if err != nil {
			return err
		}

		argDecl, err := a.LookupDecl(arg.Text, FindVarOnly)
		if err != nil {
			return err
		}

		if err := matchesDeclWiden(retDecl, argDecl); err != nil {
			return err
		}

		fmt.Fprintf(a.out(), "%s = %s;\n", toID(ret), toID(arg.Text))
		return nil
	}

	if len(rets) != 1 {
		return fmt.Errorf("expected exactly 1 return variable; got %q", rets)
	}

	switch len(args) {
	case 1:
		// x <- y
		// x <- 123

		arg := args[0]
		ret := rets[0]

		retDecl, err := a.LookupDecl(ret, FindVarOnly)
		if err != nil {
			return err
		}

		switch arg.Case {
		case parser.IDToken:
			argDecl, err := a.LookupDecl(arg.Text, FindVarOnly)
			if err != nil {
				return err
			}

			if err := matchesDecl(retDecl, argDecl); err != nil {
				return err
			}
		}

		fmt.Fprintf(a.out(), "%s = %s;\n", rets[0], args[0].Text)
		return nil

	case 2:
		// x <- <unaryOp> y
		// x <- <unaryOp> 123

		// TODO: Check if argument is immediate or variable, and validate
		// variables are defined.
		//
		// TODO: Validate operation is defined.
		fmt.Fprintf(a.out(), "%s = %s %s;\n", rets[0], args[0].Text, args[1].Text)
		return nil

	case 3:
		// x <- y   <binaryOp> z
		// x <- 123 <binaryOp> 456
		// x <- y   <binaryOp> 123
		// x <- 123 <binaryOp> y

		// TODO: Check if argument is immediate or variable, and validate
		// variables are defined.
		//
		// TODO: Validate operation is defined.
		fmt.Fprintf(a.out(), "%s = %s %s %s;\n", rets[0], args[0].Text, args[1].Text, args[2].Text)
		return nil
	default:
		return fmt.Errorf("expected 1, 2 or 3 arguments; got %q", args)
	}
}

func (a *Compiler) Return() error {
	if a.blocks.Peek() != functionBlock {
		return fmt.Errorf("can only be used within a function block")
	}

	fmt.Fprintf(a.out(), "return ")

	switch rets := a.fun().rets(); len(rets) {
	case 0:
		break
	case 1:
		fmt.Fprintf(a.out(), "%s", rets[0].Id)
	default:
		fmt.Fprintf(a.out(), "{%s", rets[0].Id)
		for _, ret := range rets[1:] {
			fmt.Fprintf(a.out(), ", %s", ret.Id)
		}
		fmt.Fprintf(a.out(), "}")
	}

	fmt.Fprintf(a.out(), ";\n")
	return nil
}

func (a *Compiler) IfThen(arg string) error {
	if !a.isFunctionBlock() {
		return errors.New("op 'if then' can only be used in a function block")
	}

	if _, err := a.LookupVar(arg); err != nil {
		return err
	}

	fmt.Fprintf(a.out(), "if (%s) {\n", arg)
	a.blocks.Push(ifThenBlock)
	return nil
}

func (a *Compiler) IfElse(arg string) error {
	if !a.isFunctionBlock() {
		return errors.New("op 'if else' can only be used in a function block")
	}

	if _, err := a.LookupVar(arg); err != nil {
		return err
	}

	fmt.Fprintf(a.out(), "if (!%s) {\n", arg)
	a.blocks.Push(ifElseBlock)
	return nil
}

func (a *Compiler) Else() error {
	if a.blocks.Pop() != ifThenBlock {
		return errors.New("expected if block")
	}

	// After the opcode, put a placeholder offset to be rewritten by
	// 'endElse'.
	fmt.Fprintf(a.out(), "} else {\n")
	a.blocks.Push(elseBlock)
	return nil
}

func (a *Compiler) End() error {
	switch block := a.blocks.Peek(); block {
	case moduleBlock:
		return a.endModule()
	case importsBlock:
		return a.endImports()
	case exportsBlock:
		return a.endExports()
	case declsBlock:
		return a.endDecls()
	case functionBlock:
		return a.endFunction()
	case ifThenBlock, ifElseBlock:
		return a.endIf()
	case elseBlock:
		return a.endElse()
	default:
		return fmt.Errorf("Unexpected block type %d", block)
	}
}

func (a *Compiler) PrintImmediate(typ IrIntType, sign Sign, value uint64) error {
	if !a.isFunctionBlock() {
		return errors.New("op 'print immediate' can only be used in a function block")
	}

	// TODO: Handle signed and unsigned.
	fmt.Fprintf(a.out(), "std::cout << %d << std::endl;\n", value)
	return nil
}

func (a *Compiler) PrintVar(sign Sign, id string) error {
	if !a.isFunctionBlock() {
		return errors.New("op 'print var' can only be used in a function block")
	}

	if _, err := a.LookupVar(id); err != nil {
		return err
	}

	fmt.Fprintf(a.out(), "std::cout << %s << std::endl;\n", id)
	return nil
}

func (a *Compiler) PrintStack(typ IrIntType, sign Sign) error {
	if !a.isFunctionBlock() {
		return errors.New("op 'print stack' can only be used in a function block")
	}

	return errors.New("PrintStack is not implemented")
}

func (a *Compiler) Program() IrProgram {
	symbols := make([]Symbol, len(a.functions))
	for i, f := range a.functions {
		symbols[i].Id = f.id
		symbols[i].Offset = f.offset
	}

	slices.SortFunc(symbols, func(a, b Symbol) bool {
		return a.Offset < b.Offset
	})

	return IrProgram{IrHeader{symbols}, nil}
}

func NewCompiler(output io.Writer) *Compiler {
	compiler := &Compiler{
		output,
		stack.New[blockType](), /* blocks */
		[]irDecl{},             /* imports */
		[]irDecl{},             /* exports */
		[]irDecl{},             /* decls */
		[]irDecl{},             /* structDefs */
		[]irFunction{},         /* functions */
		NewOpTable(),
		map[string]irCallsite{}, /* callsites */
	}
	return compiler
}
