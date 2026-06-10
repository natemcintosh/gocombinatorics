package gocombinatorics

import (
	"iter"
	"math/big"
	"slices"
	"testing"
)

// check_nth_matches_all asserts that nth(i) returns exactly what the i-th
// iteration of all yields, for every rank in [0, length).
func check_nth_matches_all(
	t *testing.T,
	length *big.Int,
	all iter.Seq2[[]int, []int],
	nth func(*big.Int) ([]int, []int, error),
) {
	t.Helper()
	rank := int64(0)
	for want_inds, want_items := range all {
		got_inds, got_items, err := nth(big.NewInt(rank))
		if err != nil {
			t.Fatalf("Nth(%d) returned error: %v", rank, err)
		}
		if !slices.Equal(got_inds, want_inds) {
			t.Fatalf("Nth(%d) indices = %v, want %v", rank, got_inds, want_inds)
		}
		if !slices.Equal(got_items, want_items) {
			t.Fatalf("Nth(%d) items = %v, want %v", rank, got_items, want_items)
		}
		rank++
	}
	if big.NewInt(rank).Cmp(length) != 0 {
		t.Fatalf("iterated %d tuples, but Length = %v", rank, length)
	}
}

// check_nth_errors asserts that out-of-range ranks return errors with nil slices.
func check_nth_errors(
	t *testing.T,
	length *big.Int,
	nth func(*big.Int) ([]int, []int, error),
) {
	t.Helper()
	for _, bad := range []*big.Int{big.NewInt(-1), length} {
		inds, items, err := nth(bad)
		if err == nil {
			t.Fatalf("Nth(%v) should have returned an error", bad)
		}
		if inds != nil || items != nil {
			t.Fatalf("Nth(%v) should return nil slices on error, got %v, %v", bad, inds, items)
		}
	}
}

func iota_slice(n int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = i
	}
	return s
}

func TestCombinationsNth(t *testing.T) {
	cases := []struct{ n, k int }{{5, 3}, {4, 4}, {6, 1}}
	for _, tc := range cases {
		c, err := NewCombinations(iota_slice(tc.n), tc.k)
		if err != nil {
			t.Fatal(err)
		}
		check_nth_matches_all(t, c.Length, c.All(), c.Nth)
		check_nth_errors(t, c.Length, c.Nth)
	}
}

func TestCombinationsWithReplacementNth(t *testing.T) {
	cases := []struct{ n, k int }{{3, 4}, {5, 2}, {1, 3}}
	for _, tc := range cases {
		c, err := NewCombinationsWithReplacement(iota_slice(tc.n), tc.k)
		if err != nil {
			t.Fatal(err)
		}
		check_nth_matches_all(t, c.Length, c.All(), c.Nth)
		check_nth_errors(t, c.Length, c.Nth)
	}
}

func TestPermutationsNth(t *testing.T) {
	cases := []struct{ n, k int }{{5, 3}, {4, 4}, {6, 1}}
	for _, tc := range cases {
		p, err := NewPermutations(iota_slice(tc.n), tc.k)
		if err != nil {
			t.Fatal(err)
		}
		check_nth_matches_all(t, p.Length, p.All(), p.Nth)
		check_nth_errors(t, p.Length, p.Nth)
	}
}

func TestPermutationsNthZeroK(t *testing.T) {
	p, err := NewPermutations(iota_slice(3), 0)
	if err != nil {
		t.Fatal(err)
	}
	inds, items, err := p.Nth(big.NewInt(0))
	if err != nil {
		t.Fatalf("Nth(0) returned error: %v", err)
	}
	if inds == nil || items == nil || len(inds) != 0 || len(items) != 0 {
		t.Fatalf("Nth(0) with k=0 should return empty non-nil slices, got %v, %v", inds, items)
	}
	check_nth_errors(t, p.Length, p.Nth)
}

func TestProductNth(t *testing.T) {
	cases := []struct{ n, k int }{{3, 4}, {2, 3}, {5, 1}, {1, 4}}
	for _, tc := range cases {
		p, err := NewProduct(iota_slice(tc.n), tc.k)
		if err != nil {
			t.Fatal(err)
		}
		check_nth_matches_all(t, p.Length, p.All(), p.Nth)
		check_nth_errors(t, p.Length, p.Nth)
	}
}

func TestPowersetNth(t *testing.T) {
	for _, n := range []int{1, 3, 5} {
		p, err := NewPowerset(iota_slice(n))
		if err != nil {
			t.Fatal(err)
		}
		check_nth_matches_all(t, p.Length, p.All(), p.Nth)
		check_nth_errors(t, p.Length, p.Nth)
	}
}

func TestPowersetNthEmptySet(t *testing.T) {
	p, err := NewPowerset(iota_slice(4))
	if err != nil {
		t.Fatal(err)
	}
	inds, items, err := p.Nth(big.NewInt(0))
	if err != nil {
		t.Fatalf("Nth(0) returned error: %v", err)
	}
	if inds == nil || items == nil || len(inds) != 0 || len(items) != 0 {
		t.Fatalf("Nth(0) should return empty non-nil slices, got %v, %v", inds, items)
	}
}

// Nth must not return slices sharing a backing array between calls, and must
// not mutate its argument.
func TestNthRetainSafety(t *testing.T) {
	c, err := NewCombinations(iota_slice(5), 3)
	if err != nil {
		t.Fatal(err)
	}
	i := big.NewInt(2)
	first_inds, first_items, err := c.Nth(i)
	if err != nil {
		t.Fatal(err)
	}
	if i.Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("Nth mutated its argument: %v", i)
	}
	want_inds := slices.Clone(first_inds)
	want_items := slices.Clone(first_items)
	second_inds, second_items, err := c.Nth(big.NewInt(2))
	if err != nil {
		t.Fatal(err)
	}
	second_inds[0] = 99
	second_items[0] = 99
	if !slices.Equal(first_inds, want_inds) || !slices.Equal(first_items, want_items) {
		t.Fatal("results of separate Nth calls share a backing array")
	}
}

// Spot checks for ranks beyond what iteration could reach in reasonable time,
// using closed-form first and last tuples.
func TestNthLargeLength(t *testing.T) {
	one := big.NewInt(1)

	// Product 10^30: first tuple all zeros, last all 9s.
	p, err := NewProduct(iota_slice(10), 30)
	if err != nil {
		t.Fatal(err)
	}
	inds, _, err := p.Nth(big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(inds, make([]int, 30)) {
		t.Fatalf("Product Nth(0) = %v, want all zeros", inds)
	}
	last := new(big.Int).Sub(p.Length, one)
	inds, _, err = p.Nth(last)
	if err != nil {
		t.Fatal(err)
	}
	want := make([]int, 30)
	for j := range want {
		want[j] = 9
	}
	if !slices.Equal(inds, want) {
		t.Fatalf("Product Nth(Length-1) = %v, want all 9s", inds)
	}

	// Combinations 100 choose 10: last tuple is [90..99].
	c, err := NewCombinations(iota_slice(100), 10)
	if err != nil {
		t.Fatal(err)
	}
	last = new(big.Int).Sub(c.Length, one)
	inds, _, err = c.Nth(last)
	if err != nil {
		t.Fatal(err)
	}
	want = make([]int, 10)
	for j := range want {
		want[j] = 90 + j
	}
	if !slices.Equal(inds, want) {
		t.Fatalf("Combinations Nth(Length-1) = %v, want %v", inds, want)
	}

	// Permutations 25 permute 25: last tuple is [24, 23, ..., 0].
	perm, err := NewPermutations(iota_slice(25), 25)
	if err != nil {
		t.Fatal(err)
	}
	last = new(big.Int).Sub(perm.Length, one)
	inds, _, err = perm.Nth(last)
	if err != nil {
		t.Fatal(err)
	}
	want = make([]int, 25)
	for j := range want {
		want[j] = 24 - j
	}
	if !slices.Equal(inds, want) {
		t.Fatalf("Permutations Nth(Length-1) = %v, want %v", inds, want)
	}

	// CombinationsWithReplacement n=50, k=20: last tuple is all 49s.
	cwr, err := NewCombinationsWithReplacement(iota_slice(50), 20)
	if err != nil {
		t.Fatal(err)
	}
	last = new(big.Int).Sub(cwr.Length, one)
	inds, _, err = cwr.Nth(last)
	if err != nil {
		t.Fatal(err)
	}
	want = make([]int, 20)
	for j := range want {
		want[j] = 49
	}
	if !slices.Equal(inds, want) {
		t.Fatalf("CombinationsWithReplacement Nth(Length-1) = %v, want all 49s", inds)
	}

	// Powerset of 200 elements: last subset is the full set.
	ps, err := NewPowerset(iota_slice(200))
	if err != nil {
		t.Fatal(err)
	}
	last = new(big.Int).Sub(ps.Length, one)
	inds, _, err = ps.Nth(last)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(inds, iota_slice(200)) {
		t.Fatalf("Powerset Nth(Length-1) should be the full set, got %v", inds)
	}
}
