package stlc_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jabolopes/bapel/comp"
	"github.com/jabolopes/bapel/ir"
	"github.com/jabolopes/bapel/query"
	"github.com/jabolopes/bapel/tests"
)

func checkUnit(filename string, typecheck bool) (ir.IrUnit, error) {
	querier, err := query.New()
	if err != nil {
		return ir.IrUnit{}, err
	}

	options := comp.TypecheckOptions{
		SkipDefaultContext:      true,
		SkipTermTypechecker:     !typecheck,
		SkipUndefinedTermChecks: true,
	}

	return comp.TypecheckSourceFile(querier, options, filename)
}

func TestInferTerm(t *testing.T) {
	matches, err := tests.Glob("*.in")
	if err != nil {
		t.Fatal(err)
	}

	for _, inFile := range matches {
		wantFile := fmt.Sprintf("%s.out", strings.TrimSuffix(inFile, ".in"))

		sourceFile, err := checkUnit(inFile, false /* typecheck */)
		if err != nil {
			t.Fatalf("in test %s: %v", inFile, err)
		}

		got := fmt.Sprintf("%+s\n", sourceFile)

		diff, err := tests.DiffOutRegen(got, wantFile)
		if err != nil {
			t.Fatalf("in test %s: %v", inFile, err)
		}
		if len(diff) > 0 {
			t.Errorf("in test %s: diff = %s", inFile, diff)
		}
	}
}

func TestTypecheckTerm(t *testing.T) {
	matches, err := tests.Glob("*.in")
	if err != nil {
		t.Fatal(err)
	}

	for _, inFile := range matches {
		if _, err := checkUnit(inFile, true /* typecheck */); err != nil {
			t.Fatalf("in test %s: %v", inFile, err)
		}
	}
}
