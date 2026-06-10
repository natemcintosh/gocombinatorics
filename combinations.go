package gocombinatorics

import (
	"errors"
	"iter"
	"math/big"
	"slices"
)

// Combinations generates all combinations of k elements from the input data.
// Use NewCombinations to create one, then iterate with All().
type Combinations[T any] struct {
	data   []T
	n, k   int
	Length *big.Int
}

// NewCombinations creates a new Combinations iterator.
func NewCombinations[T any](input_data []T, k int) (*Combinations[T], error) {
	data := make([]T, len(input_data))
	copy(data, input_data)
	n := len(input_data)

	if k > n {
		return nil, errors.New("k must be less than or equal to len(input_data)")
	} else if n <= 0 {
		return nil, errors.New("len(input_data) must be greater than 0")
	} else if k <= 0 {
		return nil, errors.New("k must be greater than 0")
	}

	Length := nchoosek(uint64(n), uint64(k))

	return &Combinations[T]{
		data:   data,
		n:      n,
		k:      k,
		Length: Length,
	}, nil
}

// All returns an iterator over all combinations. Each iteration yields
// a freshly allocated indices slice and items slice.
// This code follows the algorithm from Python's itertools.combinations
// (https://docs.python.org/3/library/itertools.html#itertools.combinations)
func (c *Combinations[T]) All() iter.Seq2[[]int, []T] {
	return func(yield func([]int, []T) bool) {
		for inds, items := range c.AllBorrowed() {
			if !yield(slices.Clone(inds), slices.Clone(items)) {
				return
			}
		}
	}
}

// AllBorrowed returns an iterator like All, but reuses internal buffers.
// Each iteration yields the same underlying index and item slices, overwritten
// in place. The caller must not retain or modify the yielded slices across
// iterations. Use All() if you need to store results.
func (c *Combinations[T]) AllBorrowed() iter.Seq2[[]int, []T] {
	return func(yield func([]int, []T) bool) {
		buf := make([]T, c.k)
		for inds := range c.IndicesBorrowed() {
			fillBuf(buf, c.data, inds)
			if !yield(inds, buf) {
				return
			}
		}
	}
}

// IndicesBorrowed yields the same index sequence as AllBorrowed, but never
// gathers items into a []T. The yielded indices slice is reused (borrowed)
// across iterations; clone it if you need to retain it. This is the fast path
// for callers that only need the combinatorial structure, not the items.
func (c *Combinations[T]) IndicesBorrowed() iter.Seq[[]int] {
	return func(yield func([]int) bool) {
		inds := make([]int, c.k)
		for i := range c.k {
			inds[i] = i
		}

		if !yield(inds) {
			return
		}

		for {
			what_is_i := -1
			for i := c.k - 1; i >= 0; i-- {
				if inds[i] != i+c.n-c.k {
					what_is_i = i
					break
				} else if i == 0 {
					return
				}
			}
			inds[what_is_i]++
			for j := what_is_i + 1; j < c.k; j++ {
				inds[j] = inds[j-1] + 1
			}
			if !yield(inds) {
				return
			}
		}
	}
}

// Nth returns the indices and items that All() would yield on its i-th
// iteration (0-based), without iterating from the start. The returned slices
// are freshly allocated and safe to retain. Returns an error if i < 0 or
// i >= Length. The argument i is not modified.
func (c *Combinations[T]) Nth(i *big.Int) ([]int, []T, error) {
	if err := check_nth_bounds(i, c.Length); err != nil {
		return nil, nil, err
	}
	inds := unrank_combination(new(big.Int).Set(i), c.n, c.k)
	items := make([]T, len(inds))
	fillBuf(items, c.data, inds)
	return inds, items, nil
}

// unrank_combination returns the i-th k-combination of {0..n-1} in
// lexicographic order, using the combinatorial number system. It assumes
// 0 <= i < nchoosek(n, k) and consumes (mutates) i.
func unrank_combination(i *big.Int, n, k int) []int {
	inds := make([]int, k)
	v := 0
	for p := range k {
		for {
			// Number of combinations that have v at position p.
			c := nchoosek(uint64(n-1-v), uint64(k-1-p))
			if i.Cmp(c) < 0 {
				inds[p] = v
				v++
				break
			}
			i.Sub(i, c)
			v++
		}
	}
	return inds
}

// nchoosek returns the number of combinations of n things taken k at a time.
// nchoosek(n, k) = n! / (k! * (n-k)!) if n > k
// nchoosek(n, k) = 0 if k > n
// nchoosek(n, k) = 1 if k == 0 or k == n
func nchoosek(n, k uint64) *big.Int {
	if k > n {
		return big.NewInt(0)
	}
	if k == 0 || k == n {
		return big.NewInt(1)
	}
	numerator := factorial(int64(n))
	kfact := factorial(int64(k))
	nminuskfact := factorial(int64(n - k))
	denominator := big.NewInt(0)
	denominator = denominator.Mul(kfact, nminuskfact)
	result := big.NewInt(0)
	return result.Div(numerator, denominator)
}

// factorial returns the factorial of a number, i.e. n! = n * (n-1) * (n-2) * ... * 1
func factorial(n int64) *big.Int {
	fact := big.NewInt(0)
	fact.MulRange(1, n)
	return fact
}
