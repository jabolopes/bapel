package stlc_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jabolopes/bapel/bplparser"
	"github.com/jabolopes/bapel/bplparser2"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/ts/stlc"
)

func p[A any](a A) *A {
	return &a
}

type expectation struct {
	got  ir.IrTerm
	want ir.IrTerm
}

func newCallWithPolymorphicID() expectation {
	i8 := ir.Const("i8")

	got := ir.CallPF(ir.ID("print"), ir.TypesA(i8), ir.Terms(ir.Number(1)))
	want := ir.CallPF(ir.ID("print"), ir.TypesA(i8), ir.Terms(ir.Number(1)))

	want.AppTerm.Fun.AppType.Fun.Type = p(ir.Forall("a", ir.NewTypeKind(), ir.NewFunctionType(ir.Types(ir.Tvar("a")), ir.Types())))
	want.AppTerm.Fun.Type = p(ir.NewFunctionType(ir.Types(i8), ir.Types()))
	want.AppTerm.Arg.Type = p(i8)
	want.Type = p(ir.NewTupleType(nil))

	return expectation{got, want}
}

func newCallWithIDs() expectation {
	i8 := ir.Const("i8")

	got := ir.Call(ir.ID("+"), ir.Terms(ir.ID("i"), ir.ID("j")))
	want := ir.CallPF(ir.ID("+"), ir.TypesA(i8), ir.Terms(ir.ID("i"), ir.ID("j")))

	want.AppTerm.Arg.Tuple.Elems[0].Type = p(i8)
	want.AppTerm.Arg.Tuple.Elems[1].Type = p(i8)
	want.AppTerm.Arg.Type = p(ir.Types(i8, i8))
	want.Type = p(i8)

	return expectation{got, want}
}

func newAssignWithIDs() expectation {
	i8 := ir.Const("i8")

	got := ir.NewAssignTerm(
		ir.Call(ir.ID("+"), ir.Terms(ir.ID("i"), ir.ID("j"))),
		ir.ID("x"))

	want := ir.NewAssignTerm(
		ir.CallPF(ir.ID("+"), ir.TypesA(i8), ir.Terms(ir.ID("i"), ir.ID("j"))),
		ir.ID("x"))

	want.Assign.Arg.AppTerm.Arg.Tuple.Elems[0].Type = p(i8)
	want.Assign.Arg.AppTerm.Arg.Tuple.Elems[1].Type = p(i8)
	want.Assign.Arg.AppTerm.Arg.Type = p(ir.NewTupleType([]ir.IrType{i8, i8}))
	want.Assign.Arg.Type = p(i8)
	want.Assign.Ret.Type = p(i8)

	return expectation{got, want}
}

func newAssignWithIDAndLiterals() expectation {
	i8 := ir.Const("i8")

	got := ir.NewAssignTerm(
		ir.Call(ir.ID("+"), ir.Terms(ir.ID("i"), ir.Number(1))),
		ir.ID("x"))

	want := ir.NewAssignTerm(
		ir.CallPF(ir.ID("+"), ir.TypesA(i8), ir.Terms(ir.ID("i"), ir.Number(1))),
		ir.ID("x"))

	want.Assign.Arg.AppTerm.Arg.Tuple.Elems[0].Type = p(i8)
	want.Assign.Arg.AppTerm.Arg.Tuple.Elems[1].Type = p(i8)
	want.Assign.Arg.AppTerm.Arg.Type = p(ir.NewTupleType([]ir.IrType{i8, i8}))
	want.Assign.Arg.Type = p(i8)
	want.Assign.Ret.Type = p(i8)

	return expectation{got, want}
}

func newAssignWithLiterals() expectation {
	i8 := ir.Const("i8")

	got := ir.NewAssignTerm(
		ir.Call(ir.ID("+"), ir.Terms(ir.Number(1), ir.Number(2))),
		ir.ID("x"))

	want := ir.NewAssignTerm(
		ir.CallPF(ir.ID("+"), ir.TypesA(i8), ir.Terms(ir.Number(1), ir.Number(2))),
		ir.ID("x"))

	want.Assign.Arg.AppTerm.Arg.Tuple.Elems[0].Type = p(i8)
	want.Assign.Arg.AppTerm.Arg.Tuple.Elems[1].Type = p(i8)
	want.Assign.Arg.AppTerm.Arg.Type = p(ir.NewTupleType([]ir.IrType{i8, i8}))
	want.Assign.Arg.Type = p(i8)
	want.Assign.Ret.Type = p(i8)

	return expectation{got, want}
}

func newTypeCast() expectation {
	i8 := ir.Const("i8")

	got := ir.CallPF(ir.Number(1), []ir.IrType{i8})
	want := typed(ir.Number(1), i8)

	return expectation{got, want}
}

func TestInferTerm2(t *testing.T) {
	i8 := ir.Const("i8")

	context := stlc.NewContext()
	binds := []stlc.Bind{
		stlc.NewTermBind("print", ir.NewForallType("a", ir.NewTypeKind(), ir.NewFunctionType(ir.Types(ir.NewVarType("a")), ir.Types())), stlc.DefSymbol),
		stlc.NewConstBind("i8", stlc.DefSymbol),
		stlc.NewTermBind("x", i8, stlc.DefSymbol),
		stlc.NewTermBind("i", i8, stlc.DefSymbol),
		stlc.NewTermBind("j", i8, stlc.DefSymbol),
	}

	for _, bind := range binds {
		var err error
		context, err = context.AddBind(bind)
		if err != nil {
			t.Fatal(err)
		}
	}

	in, err := os.Open("inferencer_test1.in")
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()

	want, err := os.ReadFile("inferencer_test1.out")
	if err != nil {
		t.Fatal(err)
	}

	sources, err := bplparser2.ParseFile(in.Name(), in)
	if err != nil {
		t.Fatal(err)
	}

	var inFunction *ir.IrFunction
	for _, source := range sources {
		if !source.Is(bplparser.FunctionSource) {
			continue
		}

		inFunction = source.Function
		break
	}

	if inFunction == nil {
		t.Fatal("Missing in function")
	}

	typechecker := stlc.NewTypechecker(context)
	if err := typechecker.InferFunction(inFunction); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(fmt.Sprintf("%v\n", inFunction), string(want)); len(diff) > 0 {
		t.Errorf("Infer() diff = %v", diff)
	}
}
