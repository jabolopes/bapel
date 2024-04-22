package list_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jabolopes/bapel/ts/list"
	"golang.org/x/exp/slices"
)

func TestList(t *testing.T) {
	l := list.New[int]()
	l1 := l.Add(1)
	l2 := l1.Add(2)
	l3 := l2.Add(3)
	l4 := l3.Remove()
	l5 := l4.Remove()
	l6 := l5.Remove()
	l7 := l6.Remove()

	tests := []struct {
		input     list.List[int]
		want      []int
		wantEmpty bool
		wantSize  int
	}{
		{l, []int{}, true, 0},
		{l1, []int{1}, false, 1},
		{l2, []int{1, 2}, false, 2},
		{l3, []int{1, 2, 3}, false, 3},
		{l4, []int{1, 2}, false, 2},
		{l5, []int{1}, false, 1},
		{l6, []int{}, true, 0},
		{l7, []int{}, true, 0},
	}

	for _, test := range tests {
		if got := test.input.Iterate().Collect(); !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) {
			t.Errorf("Iterate(%v).Collect() = %v; want %v", test.input, got, test.want)
		}

		if got := test.input.Empty(); got != test.wantEmpty {
			t.Errorf("Empty(%v) = %v; want %v", test.input, got, test.wantEmpty)
		}

		if got := test.input.Size(); got != test.wantSize {
			t.Errorf("Size(%v) = %v; want %v", test.input, got, test.wantSize)
		}

		slices.Reverse(test.want)
		if got := test.input.Iterate().CollectReverse(); !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) {
			t.Errorf("Iterate.CollectReverse(%v) = %v; want %v", test.input, got, test.want)
		}
	}
}
