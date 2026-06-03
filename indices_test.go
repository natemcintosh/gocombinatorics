package gocombinatorics

import (
	"iter"
	"slices"
	"testing"
)

// collectIndices drains an index-only iterator, cloning each borrowed slice.
func collectIndices(seq iter.Seq[[]int]) [][]int {
	var out [][]int
	for inds := range seq {
		out = append(out, slices.Clone(inds))
	}
	return out
}

// indicesFromAll drains an AllBorrowed-style iterator, keeping only the indices.
func indicesFromAll[T any](seq iter.Seq2[[]int, []T]) [][]int {
	var out [][]int
	for inds := range seq {
		out = append(out, slices.Clone(inds))
	}
	return out
}

func assertSameIndices(t *testing.T, name string, got, want [][]int) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s: IndicesBorrowed yielded %d rows, AllBorrowed yielded %d", name, len(got), len(want))
	}
	for i := range want {
		if !slices.Equal(got[i], want[i]) {
			t.Errorf("%s: row %d: IndicesBorrowed=%v, AllBorrowed=%v", name, i, got[i], want[i])
		}
	}
}

// TestIndicesBorrowedMatchesAllBorrowed verifies that the index-only fast path
// yields exactly the same index sequence as AllBorrowed for every iterator type.
func TestIndicesBorrowedMatchesAllBorrowed(t *testing.T) {
	data := stepped_range(0, 6, 1) // [0 1 2 3 4 5]
	k := 3

	comb, err := NewCombinations(data, k)
	if err != nil {
		t.Fatal(err)
	}
	assertSameIndices(t, "Combinations",
		collectIndices(comb.IndicesBorrowed()), indicesFromAll(comb.AllBorrowed()))

	cwr, err := NewCombinationsWithReplacement(data, k)
	if err != nil {
		t.Fatal(err)
	}
	assertSameIndices(t, "CombinationsWithReplacement",
		collectIndices(cwr.IndicesBorrowed()), indicesFromAll(cwr.AllBorrowed()))

	perm, err := NewPermutations(data, k)
	if err != nil {
		t.Fatal(err)
	}
	assertSameIndices(t, "Permutations",
		collectIndices(perm.IndicesBorrowed()), indicesFromAll(perm.AllBorrowed()))

	prod, err := NewProduct(data, k)
	if err != nil {
		t.Fatal(err)
	}
	assertSameIndices(t, "Product",
		collectIndices(prod.IndicesBorrowed()), indicesFromAll(prod.AllBorrowed()))

	ps, err := NewPowerset(data)
	if err != nil {
		t.Fatal(err)
	}
	assertSameIndices(t, "Powerset",
		collectIndices(ps.IndicesBorrowed()), indicesFromAll(ps.AllBorrowed()))
}

// benchCombinationsIndicesVsItems compares the cost of materializing items
// (AllBorrowed) against the index-only fast path (IndicesBorrowed) for a given
// payload type T.
func benchCombinationsIndicesVsItems[T any](b *testing.B, data []T, k int) {
	c, err := NewCombinations(data, k)
	if err != nil {
		b.Fatal(err)
	}
	b.Run("Items", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			for range c.AllBorrowed() {
			}
		}
	})
	b.Run("Indices", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			for range c.IndicesBorrowed() {
			}
		}
	})
}

// BenchmarkCombinationsIndicesVsItems demonstrates that index-only iteration is
// faster than item materialization and, unlike materialization, is independent
// of the payload size: C(22, 11) = 705,432 combinations of int vs [192]byte.
func BenchmarkCombinationsIndicesVsItems(b *testing.B) {
	const n, k = 22, 11
	b.Run("T=int", func(b *testing.B) {
		benchCombinationsIndicesVsItems(b, stepped_range(0, n, 1), k)
	})
	b.Run("T=[192]byte", func(b *testing.B) {
		benchCombinationsIndicesVsItems(b, make([][192]byte, n), k)
	})
}
