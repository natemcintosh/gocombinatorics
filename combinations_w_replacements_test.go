package gocombinatorics

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"testing"
)

func TestNewCombinationsWithReplacementErrors(t *testing.T) {
	testCases := []struct {
		desc        string
		n           int
		k           int
		want_struct *CombinationsWithReplacement[int]
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
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			data := stepped_range(0, tC.n, 1)
			_, got_err := NewCombinationsWithReplacement(data, tC.k)
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
			want: csv_to_2d_int_array("testdata/15_combo_w_replacement_5.csv"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			data := stepped_range(0, tC.n, 1)
			cwr, err := NewCombinationsWithReplacement(data, tC.k)
			if err != nil {
				t.Errorf("NewCombinationsWithReplacement(%d, %d) = %v, want nil", tC.n, tC.k, err)
			}
			got := make([][]int, 0)
			for inds := range cwr.All() {
				got = append(got, inds)
			}

			if !reflect.DeepEqual(got, tC.want) {
				t.Errorf("CombinationsWithReplacement(%d, %d) = %v, want %v", tC.n, tC.k, got, tC.want)
			}
		})
	}
}

func FuzzCWRAllBorrowedMatchesAll(f *testing.F) {
	f.Add(3, 2)
	f.Add(3, 3)
	f.Add(5, 1)
	f.Add(10, 3)
	f.Fuzz(func(t *testing.T, n int, k int) {
		if n <= 0 || k <= 0 || n > 50 || k > 50 {
			t.Skip()
		}
		length := num_combinations_w_replacement(n, k)
		if length.Cmp(big.NewInt(1_000_000)) > 0 {
			t.Skip()
		}

		data := stepped_range(0, n, 1)
		c, err := NewCombinationsWithReplacement(data, k)
		if err != nil {
			t.Fatal(err)
		}

		var allInds [][]int
		for inds := range c.All() {
			allInds = append(allInds, inds)
		}

		i := 0
		for inds := range c.AllBorrowed() {
			if i >= len(allInds) {
				t.Fatalf("AllBorrowed yielded more items than All (%d)", len(allInds))
			}
			if !reflect.DeepEqual(inds, allInds[i]) {
				t.Errorf("iteration %d: AllBorrowed=%v, All=%v", i, inds, allInds[i])
			}
			i++
		}
		if i != len(allInds) {
			t.Errorf("AllBorrowed yielded %d items, All yielded %d", i, len(allInds))
		}
	})
}

func TestCWRAllBorrowedMatchesAll(t *testing.T) {
	testCases := []struct {
		n, k int
	}{
		{3, 2}, {3, 3}, {5, 1}, {10, 3},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("n=%d,k=%d", tc.n, tc.k), func(t *testing.T) {
			data := stepped_range(0, tc.n, 1)
			c, err := NewCombinationsWithReplacement(data, tc.k)
			if err != nil {
				t.Fatal(err)
			}

			var allInds [][]int
			var allItems [][]int
			for inds, items := range c.All() {
				allInds = append(allInds, inds)
				allItems = append(allItems, items)
			}

			i := 0
			for inds, items := range c.AllBorrowed() {
				if !reflect.DeepEqual(inds, allInds[i]) {
					t.Errorf("iteration %d indices: AllBorrowed=%v, All=%v", i, inds, allInds[i])
				}
				if !reflect.DeepEqual(items, allItems[i]) {
					t.Errorf("iteration %d items: AllBorrowed=%v, All=%v", i, items, allItems[i])
				}
				i++
			}
			if i != len(allInds) {
				t.Errorf("AllBorrowed yielded %d items, All yielded %d", i, len(allInds))
			}
		})
	}
}

func BenchmarkCWRAllVsBorrowed(b *testing.B) {
	benchmarks := []struct {
		desc string
		n    int
		k    int
	}{
		{desc: "n=10,k=3", n: 10, k: 3},
		{desc: "n=15,k=5", n: 15, k: 5},
	}
	for _, bm := range benchmarks {
		data := stepped_range(0, bm.n, 1)
		c, err := NewCombinationsWithReplacement(data, bm.k)
		if err != nil {
			b.Fatal(err)
		}

		b.Run("All/"+bm.desc, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for range c.All() {
				}
			}
		})
		b.Run("AllBorrowed/"+bm.desc, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for range c.AllBorrowed() {
				}
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
			for b.Loop() {
				data := stepped_range(0, bm.n, 1)
				cwr, err := NewCombinationsWithReplacement(data, bm.k)
				if err != nil {
					b.Errorf("NewCombinationsWithReplacement(%d, %d) = %v, want nil", bm.n, bm.k, err)
				}
				for range cwr.All() {
				}
			}
		})
	}
}
