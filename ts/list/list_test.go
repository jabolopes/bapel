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
		wantValue int
		wantEmpty bool
		wantSize  int
	}{
		{l, []int{}, 0, true, 0},
		{l1, []int{1}, 1, false, 1},
		{l2, []int{1, 2}, 2, false, 2},
		{l3, []int{1, 2, 3}, 3, false, 3},
		{l4, []int{1, 2}, 2, false, 2},
		{l5, []int{1}, 1, false, 1},
		{l6, []int{}, 0, true, 0},
		{l7, []int{}, 0, true, 0},
	}

	for _, test := range tests {
		if got := test.input.Iterate().Collect(); !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) {
			t.Errorf("Iterate(%v).Collect() = %v; want %v", test.input, got, test.want)
		}

		if got, gotOk := test.input.Value(); got != test.wantValue || gotOk != !test.wantEmpty {
			t.Errorf("Value(%v) = %v; want %v", test.input, got, test.wantValue)
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

func TestFromSlice(t *testing.T) {
	tests := []struct {
		input []int
	}{
		{nil},
		{[]int{1}},
		{[]int{1, 2}},
		{[]int{1, 2, 3}},
	}

	for _, test := range tests {
		if got := list.FromSlice(test.input).Iterate().Collect(); !cmp.Equal(got, test.input, cmpopts.EquateEmpty()) {
			t.Errorf("Collect(%v) = %v; want %v", test.input, got, test.input)
		}
	}
}

func TestIterate(t *testing.T) {
	l := list.New[int]().Add(1).Add(2).Add(3)

	it := l.Iterate()

	tests := []struct {
		wantIndex int
		want      int
	}{
		{2, 3},
		{1, 2},
		{0, 1},
	}

	for _, test := range tests {
		if gotIndex, got, gotOk := it.Next(); gotIndex != test.wantIndex || got != test.want || !gotOk {
			t.Errorf("Next() = %v, %v, %v; want %v, %v, %v", gotIndex, got, gotOk, test.wantIndex, test.want, true)
		}
	}

	for i := 0; i < 10; i++ {
		if _, _, gotOk := it.Next(); gotOk {
			t.Errorf("Next() = _, %v; want _, %v", gotOk, false)
		}
	}
}

func TestIterateCollect(t *testing.T) {
	tests := []struct {
		input list.List[int]
		want  []int
	}{
		{list.New[int](), nil},
		{list.New[int]().Add(1).Add(2).Add(3), []int{1, 2, 3}},
	}

	for _, test := range tests {
		if got := test.input.Iterate().Collect(); !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) {
			t.Errorf("Collect(%v) = %v; want %v", test.input, got, test.want)
		}
	}
}

func TestIterateCollectReverse(t *testing.T) {
	tests := []struct {
		input list.List[int]
		want  []int
	}{
		{list.New[int](), nil},
		{list.New[int]().Add(1).Add(2).Add(3), []int{3, 2, 1}},
	}

	for _, test := range tests {
		if got := test.input.Iterate().CollectReverse(); !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) {
			t.Errorf("CollectReverse(%v) = %v; want %v", test.input, got, test.want)
		}
	}
}
