package gocombinatorics

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"testing"
)

func TestNewPermutationErrors(t *testing.T) {
	testCases := []struct {
		desc        string
		n           int
		k           int
		want_struct *Permutations[int]
		want_err    error
	}{
		{
			desc:        "k>n",
			n:           3,
			k:           4,
			want_struct: nil,
			want_err:    errors.New("k must be less than or equal to len(input_data)"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			data := stepped_range(0, tC.n, 1)
			got_struct, got_err := NewPermutations(data, tC.k)
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
			want: csv_to_2d_int_array("testdata/100_perm_3.csv"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			data := stepped_range(0, tC.n, 1)
			permutations, err := NewPermutations(data, tC.k)
			if err != nil {
				t.Errorf("NewPermutation() = %v, want %v", err, nil)
			}
			got := make([][]int, 0)
			for inds := range permutations.All() {
				got = append(got, inds)
			}

			if !reflect.DeepEqual(got, tC.want) {
				t.Errorf("Permutations(%d, %d) = %v, want %v", tC.n, tC.k, got, tC.want)
			}
		})
	}
}

func FuzzPermutationsAllBorrowedMatchesAll(f *testing.F) {
	f.Add(2, 1)
	f.Add(2, 2)
	f.Add(5, 3)
	f.Add(10, 3)
	f.Fuzz(func(t *testing.T, n int, k int) {
		if n <= 0 || k <= 0 || k > n || n > 50 {
			t.Skip()
		}
		length := n_permutations(n, k)
		if length.Cmp(big.NewInt(1_000_000)) > 0 {
			t.Skip()
		}

		data := stepped_range(0, n, 1)
		p, err := NewPermutations(data, k)
		if err != nil {
			t.Fatal(err)
		}

		var allInds [][]int
		for inds := range p.All() {
			allInds = append(allInds, inds)
		}

		i := 0
		for inds := range p.AllBorrowed() {
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

func TestPermutationsAllBorrowedMatchesAll(t *testing.T) {
	testCases := []struct {
		n, k int
	}{
		{2, 1}, {2, 2}, {5, 3}, {10, 3},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("n=%d,k=%d", tc.n, tc.k), func(t *testing.T) {
			data := stepped_range(0, tc.n, 1)
			p, err := NewPermutations(data, tc.k)
			if err != nil {
				t.Fatal(err)
			}

			var allInds [][]int
			var allItems [][]int
			for inds, items := range p.All() {
				allInds = append(allInds, inds)
				allItems = append(allItems, items)
			}

			i := 0
			for inds, items := range p.AllBorrowed() {
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

func BenchmarkPermutationsAllVsBorrowed(b *testing.B) {
	benchmarks := []struct {
		desc string
		n    int
		k    int
	}{
		{desc: "n=10,k=3", n: 10, k: 3},
		{desc: "n=10,k=8", n: 10, k: 8},
	}
	for _, bm := range benchmarks {
		data := stepped_range(0, bm.n, 1)
		p, err := NewPermutations(data, bm.k)
		if err != nil {
			b.Fatal(err)
		}

		b.Run("All/"+bm.desc, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for range p.All() {
				}
			}
		})
		b.Run("AllBorrowed/"+bm.desc, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for range p.AllBorrowed() {
				}
			}
		})
	}
}

func TestPermutationsItems(t *testing.T) {
	data := []string{"a", "b", "c", "d", "e"}
	p, err := NewPermutations(data, 2)
	if err != nil {
		t.Fatalf("NewPermutations() error: %v", err)
	}

	// Get the first permutation
	for _, items := range p.All() {
		if len(items) != 2 {
			t.Errorf("len(Items()) = %d, want 2", len(items))
		}

		want := []string{"a", "b"}
		if !reflect.DeepEqual(items, want) {
			t.Errorf("Items() = %v, want %v", items, want)
		}
		break // only check the first permutation
	}
}
