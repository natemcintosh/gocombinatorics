package gocombinatorics

import (
	"errors"
	"iter"
	"math/big"
	"slices"
)

// Product generates the k-fold Cartesian product of the input data with itself,
// i.e. all n^k tuples of length k. Use NewProduct to create one, then iterate
// with All().
type Product[T any] struct {
	data   []T
	n, k   int
	Length *big.Int
}

// NewProduct creates a new Product iterator. Unlike Combinations and
// Permutations, k may exceed n, since each position independently ranges over
// all n elements.
func NewProduct[T any](input_data []T, k int) (*Product[T], error) {
	data := make([]T, len(input_data))
	copy(data, input_data)
	n := len(input_data)

	if n <= 0 {
		return nil, errors.New("len(input_data) must be greater than 0")
	} else if k <= 0 {
		return nil, errors.New("k must be greater than 0")
	}

	Length := num_products(n, k)

	return &Product[T]{
		data:   data,
		n:      n,
		k:      k,
		Length: Length,
	}, nil
}

// All returns an iterator over all products. Each iteration yields a freshly
// allocated indices slice and items slice.
// This code follows the algorithm from Python's itertools.product
// (https://docs.python.org/3/library/itertools.html#itertools.product)
func (p *Product[T]) All() iter.Seq2[[]int, []T] {
	return func(yield func([]int, []T) bool) {
		for inds, items := range p.AllBorrowed() {
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
func (p *Product[T]) AllBorrowed() iter.Seq2[[]int, []T] {
	return func(yield func([]int, []T) bool) {
		buf := make([]T, p.k)
		for inds := range p.IndicesBorrowed() {
			fillBuf(buf, p.data, inds)
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
func (p *Product[T]) IndicesBorrowed() iter.Seq[[]int] {
	return func(yield func([]int) bool) {
		inds := make([]int, p.k)

		if !yield(inds) {
			return
		}

		for {
			// Increment the rightmost index, carrying to the left like an odometer.
			i := p.k - 1
			for i >= 0 && inds[i] == p.n-1 {
				inds[i] = 0
				i--
			}
			if i < 0 {
				return
			}
			inds[i]++
			if !yield(inds) {
				return
			}
		}
	}
}

// num_products returns the number of k-fold products of n elements, i.e. n^k.
func num_products(n, k int) *big.Int {
	return new(big.Int).Exp(big.NewInt(int64(n)), big.NewInt(int64(k)), nil)
}

// elts_in_product is how many times we expect to see a given element across all
// k-fold products of n elements. Each of the k positions independently takes
// each element n^(k-1) times, so the count is k * n^(k-1).
func elts_in_product(n, k int) *big.Int {
	per_position := new(big.Int).Exp(big.NewInt(int64(n)), big.NewInt(int64(k-1)), nil)
	return per_position.Mul(per_position, big.NewInt(int64(k)))
}
