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
		inds := make([]int, c.k)
		for i := range c.k {
			inds[i] = i
		}

		if !yield(slices.Clone(inds), c.items(inds)) {
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
			if !yield(slices.Clone(inds), c.items(inds)) {
				return
			}
		}
	}
}

// items builds a fresh slice of items at the given indices.
func (c *Combinations[T]) items(inds []int) []T {
	result := make([]T, len(inds))
	for i, idx := range inds {
		result[i] = c.data[idx]
	}
	return result
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
