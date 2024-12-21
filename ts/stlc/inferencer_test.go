package stlc_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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

	got := ir.CallPF("print", ir.TypesA(i8), ir.Terms(ir.Number(1)))
	want := ir.CallPF("print", ir.TypesA(i8), ir.Terms(ir.Number(1)))

	want.AppTerm.Fun.AppType.Fun.Type = p(ir.Forall("a", ir.NewTypeKind(), ir.NewFunctionType(ir.Types(ir.Tvar("a")), ir.Types())))
	want.AppTerm.Fun.Type = p(ir.NewFunctionType(ir.Types(i8), ir.Types()))
	want.AppTerm.Arg.Type = p(i8)
	want.Type = p(ir.NewTupleType(nil))

	return expectation{got, want}
}

func newCallWithIDs() expectation {
	i8 := ir.Const("i8")

	got := ir.Call("+", ir.Terms(ir.ID("i"), ir.ID("j")))
	want := ir.CallPF("+", ir.TypesA(i8), ir.Terms(ir.ID("i"), ir.ID("j")))

	want.AppTerm.Arg.Tuple.Elems[0].Type = p(i8)
	want.AppTerm.Arg.Tuple.Elems[1].Type = p(i8)
	want.AppTerm.Arg.Type = p(ir.Types(i8, i8))
	want.Type = p(i8)

	return expectation{got, want}
}

func newAssignWithIDs() expectation {
	i8 := ir.Const("i8")

	got := ir.NewAssignTerm(
		ir.Call("+", ir.Terms(ir.ID("i"), ir.ID("j"))),
		ir.ID("x"))

	want := ir.NewAssignTerm(
		ir.CallPF("+", ir.TypesA(i8), ir.Terms(ir.ID("i"), ir.ID("j"))),
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
		ir.Call("+", ir.Terms(ir.ID("i"), ir.Number(1))),
		ir.ID("x"))

	want := ir.NewAssignTerm(
		ir.CallPF("+", ir.TypesA(i8), ir.Terms(ir.ID("i"), ir.Number(1))),
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
		ir.Call("+", ir.Terms(ir.Number(1), ir.Number(2))),
		ir.ID("x"))

	want := ir.NewAssignTerm(
		ir.CallPF("+", ir.TypesA(i8), ir.Terms(ir.Number(1), ir.Number(2))),
		ir.ID("x"))

	want.Assign.Arg.AppTerm.Arg.Tuple.Elems[0].Type = p(i8)
	want.Assign.Arg.AppTerm.Arg.Tuple.Elems[1].Type = p(i8)
	want.Assign.Arg.AppTerm.Arg.Type = p(ir.NewTupleType([]ir.IrType{i8, i8}))
	want.Assign.Arg.Type = p(i8)
	want.Assign.Ret.Type = p(i8)

	return expectation{got, want}
}

func TestInferTerm(t *testing.T) {
	i8 := ir.Const("i8")

	tests := []expectation{
		newCallWithPolymorphicID(),
		newCallWithIDs(),
		newAssignWithIDs(),
		newAssignWithIDAndLiterals(),
		newAssignWithLiterals(),
	}

	context := stlc.NewContext()
	binds := []stlc.Bind{
		stlc.NewTermBind("print", ir.NewForallType("a", ir.NewTypeKind(), ir.NewFunctionType(ir.Types(ir.NewVarType("a")), ir.Types())), stlc.DefSymbol),
		stlc.NewNameBind("i8", stlc.DefSymbol),
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

	inferencer := stlc.NewInferencer(context)
	for _, test := range tests {
		got := test.got
		if err := inferencer.InferTerm(&got); !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) || err != nil {
			t.Errorf("Infer() = %v, %v; want %v, %v", got, err, test.want, nil)
		}
	}
}
