package stlc_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/parser"
	"github.com/jabolopes/bapel/ts/stlc"
)

type expectation struct {
	got  ir.IrTerm
	want ir.IrTerm
}

func newCallWithIDs() expectation {
	got := ir.NewCallTerm(
		"+", nil, /* types */
		ir.NewTupleTerm([]ir.IrTerm{
			ir.NewTokenTerm(parser.NewIDToken("i")),
			ir.NewTokenTerm(parser.NewIDToken("j")),
		}))

	want := ir.NewCallTerm(
		"+", []ir.IrType{ir.NewNameType("i8")},
		ir.NewTupleTerm([]ir.IrTerm{
			ir.NewTokenTerm(parser.NewIDToken("i")),
			ir.NewTokenTerm(parser.NewIDToken("j")),
		}))

	{
		typ1 := ir.NewNameType("i8")
		want.Call.Arg.Tuple[0].Type = &typ1
		want.Call.Arg.Tuple[1].Type = &typ1
		typ2 := ir.NewTupleType([]ir.IrType{typ1, typ1})
		want.Call.Arg.Type = &typ2
	}

	return expectation{got, want}
}

func newAssignWithIDs() expectation {
	got := ir.NewAssignTerm(
		ir.NewCallTerm(
			"+", nil, /* types */
			ir.NewTupleTerm([]ir.IrTerm{
				ir.NewTokenTerm(parser.NewIDToken("i")),
				ir.NewTokenTerm(parser.NewIDToken("j")),
			})),
		ir.NewTokenTerm(parser.NewIDToken("x")))

	want := ir.NewAssignTerm(
		ir.NewCallTerm(
			"+", []ir.IrType{ir.NewNameType("i8")},
			ir.NewTupleTerm([]ir.IrTerm{
				ir.NewTokenTerm(parser.NewIDToken("i")),
				ir.NewTokenTerm(parser.NewIDToken("j")),
			})),
		ir.NewTokenTerm(parser.NewIDToken("x")))

	{
		typ1 := ir.NewNameType("i8")
		want.Assign.Arg.Call.Arg.Tuple[0].Type = &typ1
		want.Assign.Arg.Call.Arg.Tuple[1].Type = &typ1
		typ2 := ir.NewTupleType([]ir.IrType{typ1, typ1})
		want.Assign.Arg.Call.Arg.Type = &typ2
	}

	{
		typ := ir.NewNameType("i8")
		want.Assign.Ret.Type = &typ
	}

	return expectation{got, want}
}

func newAssignWithIDAndLiterals() expectation {
	got := ir.NewAssignTerm(
		ir.NewCallTerm(
			"+", nil, /* types */
			ir.NewTupleTerm([]ir.IrTerm{
				ir.NewTokenTerm(parser.NewIDToken("i")),
				ir.NewTokenTerm(parser.NewNumberToken(1)),
			})),
		ir.NewTokenTerm(parser.NewIDToken("x")))

	want := ir.NewAssignTerm(
		ir.NewCallTerm(
			"+", []ir.IrType{ir.NewNameType("i8")},
			ir.NewTupleTerm([]ir.IrTerm{
				ir.NewTokenTerm(parser.NewIDToken("i")),
				ir.NewTokenTerm(parser.NewNumberToken(1)),
			})),
		ir.NewTokenTerm(parser.NewIDToken("x")))

	{
		typ := ir.NewNameType("i8")
		want.Assign.Arg.Call.Arg.Tuple[0].Type = &typ
	}

	{
		typ := ir.NewNameType("i8")
		want.Assign.Ret.Type = &typ
	}

	return expectation{got, want}
}

func newAssignWithLiterals() expectation {
	got := ir.NewAssignTerm(
		ir.NewCallTerm(
			"+", nil, /* types */
			ir.NewTupleTerm([]ir.IrTerm{
				ir.NewTokenTerm(parser.NewNumberToken(1)),
				ir.NewTokenTerm(parser.NewNumberToken(2)),
			})),
		ir.NewTokenTerm(parser.NewIDToken("x")))

	want := ir.NewAssignTerm(
		ir.NewCallTerm(
			"+", []ir.IrType{ir.NewNameType("i8")},
			ir.NewTupleTerm([]ir.IrTerm{
				ir.NewTokenTerm(parser.NewNumberToken(1)),
				ir.NewTokenTerm(parser.NewNumberToken(2)),
			})),
		ir.NewTokenTerm(parser.NewIDToken("x")))

	{
		typ := ir.NewNameType("i8")
		want.Assign.Ret.Type = &typ
	}

	return expectation{got, want}
}

func TestInference(t *testing.T) {
	tests := []expectation{
		newCallWithIDs(),
		newAssignWithIDs(),
		newAssignWithIDAndLiterals(),
		newAssignWithLiterals(),
	}

	context := stlc.NewContext()
	binds := []stlc.Bind{
		stlc.NewDeclBind(stlc.DefSymbol, ir.NewTypeDecl(ir.NewNameType("i8"))),
		stlc.NewDeclBind(stlc.DefSymbol, ir.NewTermDecl("x", ir.NewNameType("i8"))),
		stlc.NewDeclBind(stlc.DefSymbol, ir.NewTermDecl("i", ir.NewNameType("i8"))),
		stlc.NewDeclBind(stlc.DefSymbol, ir.NewTermDecl("j", ir.NewNameType("i8"))),
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
		if err := inferencer.Infer(&got); !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) || err != nil {
			t.Errorf("Infer() = %v, %v; want %v, %v", got, err, test.want, nil)
		}
	}
}
