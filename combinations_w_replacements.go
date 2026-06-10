package gocombinatorics

import (
	"errors"
	"iter"
	"math/big"
	"slices"
)

// CombinationsWithReplacement generates all combinations with replacement
// of k elements from the input data. Use NewCombinationsWithReplacement
// to create one, then iterate with All().
type CombinationsWithReplacement[T any] struct {
	data   []T
	n, k   int
	Length *big.Int
}

// NewCombinationsWithReplacement creates a new CombinationsWithReplacement iterator.
func NewCombinationsWithReplacement[T any](input_data []T, k int) (*CombinationsWithReplacement[T], error) {
	data := make([]T, len(input_data))
	copy(data, input_data)
	n := len(input_data)

	if n <= 0 {
		return nil, errors.New("len(input_data) must be greater than 0")
	} else if k <= 0 {
		return nil, errors.New("k must be greater than 0")
	}

	Length := num_combinations_w_replacement(n, k)

	return &CombinationsWithReplacement[T]{
		data:   data,
		n:      n,
		k:      k,
		Length: Length,
	}, nil
}

// All returns an iterator over all combinations with replacement. Each
// iteration yields a freshly allocated indices slice and items slice.
// This code follows the algorithm from Python's itertools.combinations_with_replacement
// (https://docs.python.org/3/library/itertools.html#itertools.combinations_with_replacement)
func (c *CombinationsWithReplacement[T]) All() iter.Seq2[[]int, []T] {
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
func (c *CombinationsWithReplacement[T]) AllBorrowed() iter.Seq2[[]int, []T] {
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
func (c *CombinationsWithReplacement[T]) IndicesBorrowed() iter.Seq[[]int] {
	return func(yield func([]int) bool) {
		inds := make([]int, c.k)

		if !yield(inds) {
			return
		}

		for {
			what_is_i := -1
			for i := c.k - 1; i >= 0; i-- {
				if inds[i] != c.n-1 {
					what_is_i = i
					break
				} else if i == 0 {
					return
				}
			}
			new_val := inds[what_is_i] + 1
			for i := what_is_i; i < c.k; i++ {
				inds[i] = new_val
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
func (c *CombinationsWithReplacement[T]) Nth(i *big.Int) ([]int, []T, error) {
	if err := check_nth_bounds(i, c.Length); err != nil {
		return nil, nil, err
	}
	// Non-descending k-tuples over n values biject order-preservingly with
	// k-combinations of n+k-1 via comb[j] = inds[j] + j (stars and bars).
	inds := unrank_combination(new(big.Int).Set(i), c.n+c.k-1, c.k)
	for j := range inds {
		inds[j] -= j
	}
	items := make([]T, len(inds))
	fillBuf(items, c.data, inds)
	return inds, items, nil
}

// num_combinations_w_replacement returns (n+k-1)! / (k! * (n-1)!)
func num_combinations_w_replacement(n, k int) *big.Int {
	numerator := factorial(int64(n + k - 1))
	k_fact := factorial(int64(k))
	n_minus_1_fact := factorial(int64(n - 1))
	denominator := new(big.Int).Mul(k_fact, n_minus_1_fact)
	result := new(big.Int)
	return result.Div(numerator, denominator)
}

// elts_in_combo_w_replacement is how many times we expect to see a given element when
// carrying out combinations with replacement. In python style code, it is:
// sum((k-(i-1)*num_combinations_w_replacement(n-1, i-1)) for i in range(1,k+1))
func elts_in_combo_w_replacement(n, k int) *big.Int {
	sum := big.NewInt(0)
	for i := 1; i <= k; i++ {
		n_cols := big.NewInt(int64(k - (i - 1)))
		n_rows := num_combinations_w_replacement(n-1, i-1)
		this_num := new(big.Int).Mul(n_cols, n_rows)
		sum.Add(sum, this_num)
	}
	return sum
}
