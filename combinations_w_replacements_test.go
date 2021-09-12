package gocombinatorics

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewCombinationsWithReplacementErrors(t *testing.T) {
	testCases := []struct {
		desc        string
		n           int
		k           int
		want_struct *CombinationsWithReplacement
		want_err    error
	}{
		{
			desc:        "n <= 0",
			n:           0,
			k:           1,
			want_struct: nil,
			want_err:    errors.New("n must be greater than 0"),
		},
		{
			desc:        "k <= 0",
			n:           1,
			k:           0,
			want_struct: nil,
			want_err:    errors.New("k must be greater than 0"),
		},
		{
			desc:        "k > n",
			n:           1,
			k:           2,
			want_struct: nil,
			want_err:    errors.New("k must be less than or equal to n"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, got_err := NewCombinationsWithReplacement(tC.n, tC.k)
			if got_err == nil {
				t.Errorf("NewCombinationsWithReplacement() = %v, want %v", got_err, tC.want_err)
			}
		})
	}
}

func TestCombinationsWithReplacementNew(t *testing.T) {
	testCases := []struct {
		desc string
		n    int
		k    int
		want [][]int
	}{
		{
			desc: "n=3, k=2",
			n:    3,
			k:    2,
			want: [][]int{
				{0, 0},
				{0, 1},
				{0, 2},
				{1, 1},
				{1, 2},
				{2, 2},
			},
		},
		{
			desc: "n=3, k=3",
			n:    3,
			k:    3,
			want: [][]int{
				{0, 0, 0},
				{0, 0, 1},
				{0, 0, 2},
				{0, 1, 1},
				{0, 1, 2},
				{0, 2, 2},
				{1, 1, 1},
				{1, 1, 2},
				{1, 2, 2},
				{2, 2, 2},
			},
		},
		{
			desc: "n=5, k=1",
			n:    5,
			k:    1,
			want: [][]int{
				{0},
				{1},
				{2},
				{3},
				{4},
			},
		},
		{
			desc: "n = 15, k = 5",
			n:    15,
			k:    5,
			want: csv_to_2d_int_array("15_combo_w_replacement_5.csv"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			combinations_w_replacement, err := NewCombinationsWithReplacement(tC.n, tC.k)
			if err != nil {
				t.Errorf("NewCombinationsWithReplacement(%d, %d) = %v, want nil", tC.n, tC.k, err)
			}
			got := make([][]int, 0)
			for combinations_w_replacement.Next() {
				// We need to append a copy of combinations.Inds to got
				next_set_of_indices := make([]int, len(combinations_w_replacement.Inds))
				copy(next_set_of_indices, combinations_w_replacement.Inds)
				got = append(got, next_set_of_indices)
			}

			if !reflect.DeepEqual(got, tC.want) {
				t.Errorf("CombinationsWithReplacement(%d, %d) = %v, want %v", tC.n, tC.k, got, tC.want)
			}
		})
	}
}

func BenchmarkCombinationsWithReplacementNext(b *testing.B) {
	benchmarks := []struct {
		desc string
		n    int
		k    int
	}{
		{
			desc: "n = 3, k = 2",
			n:    3,
			k:    2,
		},
		{
			desc: "n = 4, k = 3",
			n:    4,
			k:    3,
		},
		{
			desc: "n = 5, k = 4",
			n:    5,
			k:    4,
		},
		{
			desc: "n = 10, k = 8",
			n:    10,
			k:    8,
		},
		{
			desc: "n = 10, k = 3",
			n:    10,
			k:    3,
		},
		{
			desc: "n = 15, k = 5",
			n:    15,
			k:    5,
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				combinations_w_replacement, err := NewCombinationsWithReplacement(bm.n, bm.k)
				if err != nil {
					b.Errorf("NewCombinations(%d, %d) = %v, want nil", bm.n, bm.k, err)
				}
				for combinations_w_replacement.Next() {

				}
			}

		})
	}
}
