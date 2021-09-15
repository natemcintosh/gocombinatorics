package gocombinatorics

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewPermutationErrors(t *testing.T) {
	testCases := []struct {
		desc        string
		n           int
		k           int
		want_struct *Permutations
		want_err    error
	}{
		{
			desc:        "k>n",
			n:           3,
			k:           4,
			want_struct: nil,
			want_err:    errors.New("k must be less than or equal to n"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got_struct, got_err := NewPermutations(tC.n, tC.k)
			if got_err == nil {
				t.Errorf("NewPermutation() = %v, want %v", got_struct, tC.want_struct)
			}
			if got_err.Error() != tC.want_err.Error() {
				t.Errorf("NewPermutation() = %v, want %v", got_err, tC.want_err)
			}
		})
	}
}

func TestPermutationsNext(t *testing.T) {
	testCases := []struct {
		desc string
		n    int
		k    int
		want [][]int
	}{
		{
			desc: "n=2, k=1",
			n:    2,
			k:    1,
			want: [][]int{
				{0},
				{1},
			},
		},
		{
			desc: "n=2, k=2",
			n:    2,
			k:    2,
			want: [][]int{
				{0, 1},
				{1, 0},
			},
		},
		{
			desc: "n=5, k=3",
			n:    5,
			k:    3,
			want: [][]int{
				{0, 1, 2},
				{0, 1, 3},
				{0, 1, 4},
				{0, 2, 1},
				{0, 2, 3},
				{0, 2, 4},
				{0, 3, 1},
				{0, 3, 2},
				{0, 3, 4},
				{0, 4, 1},
				{0, 4, 2},
				{0, 4, 3},
				{1, 0, 2},
				{1, 0, 3},
				{1, 0, 4},
				{1, 2, 0},
				{1, 2, 3},
				{1, 2, 4},
				{1, 3, 0},
				{1, 3, 2},
				{1, 3, 4},
				{1, 4, 0},
				{1, 4, 2},
				{1, 4, 3},
				{2, 0, 1},
				{2, 0, 3},
				{2, 0, 4},
				{2, 1, 0},
				{2, 1, 3},
				{2, 1, 4},
				{2, 3, 0},
				{2, 3, 1},
				{2, 3, 4},
				{2, 4, 0},
				{2, 4, 1},
				{2, 4, 3},
				{3, 0, 1},
				{3, 0, 2},
				{3, 0, 4},
				{3, 1, 0},
				{3, 1, 2},
				{3, 1, 4},
				{3, 2, 0},
				{3, 2, 1},
				{3, 2, 4},
				{3, 4, 0},
				{3, 4, 1},
				{3, 4, 2},
				{4, 0, 1},
				{4, 0, 2},
				{4, 0, 3},
				{4, 1, 0},
				{4, 1, 2},
				{4, 1, 3},
				{4, 2, 0},
				{4, 2, 1},
				{4, 2, 3},
				{4, 3, 0},
				{4, 3, 1},
				{4, 3, 2},
			},
		},
		{
			desc: "n=100, k=3",
			n:    100,
			k:    3,
			want: csv_to_2d_int_array("100_perm_3.csv"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			permutations, err := NewPermutations(tC.n, tC.k)
			if err != nil {
				t.Errorf("NewPermutation() = %v, want %v", err, nil)
			}
			got := make([][]int, 0)

			for permutations.Next() {
				// We need to append a copy of permutations.Indices() to got
				next_set_of_indices := make([]int, len(permutations.Indices()))
				copy(next_set_of_indices, permutations.Indices())
				got = append(got, next_set_of_indices)
			}

			if !reflect.DeepEqual(got, tC.want) {
				t.Errorf("Permutations(%d, %d) = %v, want %v", tC.n, tC.k, got, tC.want)
			}

		})
	}
}
