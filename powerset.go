package gocombinatorics

import (
	"errors"
	"iter"
	"math/big"
	"slices"
)

// Powerset generates all 2^n subsets of the input data, including the empty set.
// Use NewPowerset to create one, then iterate with All().
type Powerset[T any] struct {
	data   []T
	n      int
	Length *big.Int
}

// NewPowerset creates a new Powerset iterator. Unlike the other iterators it
// takes no k, since a powerset ranges over subsets of every length 0..n.
func NewPowerset[T any](input_data []T) (*Powerset[T], error) {
	data := make([]T, len(input_data))
	copy(data, input_data)
	n := len(input_data)

	if n <= 0 {
		return nil, errors.New("len(input_data) must be greater than 0")
	}

	Length := num_powersets(n)

	return &Powerset[T]{
		data:   data,
		n:      n,
		Length: Length,
	}, nil
}

// All returns an iterator over all subsets. Each iteration yields a freshly
// allocated indices slice and items slice.
// This follows Python's itertools powerset recipe
// (https://docs.python.org/3/library/itertools.html#itertools-recipes)
func (p *Powerset[T]) All() iter.Seq2[[]int, []T] {
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
// in place; their length grows with the subset size. The caller must not retain
// or modify the yielded slices across iterations. Use All() if you need to store
// results.
func (p *Powerset[T]) AllBorrowed() iter.Seq2[[]int, []T] {
	return func(yield func([]int, []T) bool) {
		buf := make([]T, 0, p.n)
		for inds := range p.IndicesBorrowed() {
			buf = buf[:len(inds)]
			fillBuf(buf, p.data, inds)
			if !yield(inds, buf) {
				return
			}
		}
	}
}

// IndicesBorrowed yields the same index sequence as AllBorrowed, but never
// gathers items into a []T. The yielded indices slice is reused (borrowed)
// across iterations and its length grows with the subset size; clone it if you
// need to retain it. This is the fast path for callers that only need the
// combinatorial structure, not the items.
func (p *Powerset[T]) IndicesBorrowed() iter.Seq[[]int] {
	return func(yield func([]int) bool) {
		// The empty set, yielded first. NewCombinations rejects r <= 0, so it is
		// handled separately here rather than in the loop below.
		if !yield([]int{}) {
			return
		}

		for r := 1; r <= p.n; r++ {
			c, err := NewCombinations(p.data, r)
			if err != nil {
				return
			}
			for inds := range c.IndicesBorrowed() {
				if !yield(inds) {
					return
				}
			}
		}
	}
}

// num_powersets returns the number of subsets of an n-element set, i.e. 2^n.
func num_powersets(n int) *big.Int {
	return new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(n)), nil)
}

// elts_in_powerset is how many times we expect to see a given element across all
// subsets. Each element is present in exactly half of the 2^n subsets, i.e.
// 2^(n-1) times.
func elts_in_powerset(n int) *big.Int {
	return new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(n-1)), nil)
}
