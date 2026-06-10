package gocombinatorics

import (
	"errors"
	"iter"
	"math/big"
	"math/rand"
	"slices"
)

// ProductOf generates the Cartesian product of multiple distinct slices
// (axes) of possibly different lengths, i.e. all tuples taking one element
// from each axis. Use NewProductOf to create one, then iterate with All().
type ProductOf[T any] struct {
	axes   [][]T
	Length *big.Int
}

// NewProductOf creates a new ProductOf iterator over the given axes. Each
// yielded tuple has one element per axis. At least one axis is required, and
// every axis must be non-empty.
func NewProductOf[T any](axes ...[]T) (*ProductOf[T], error) {
	if len(axes) == 0 {
		return nil, errors.New("at least one axis is required")
	}

	copied := make([][]T, len(axes))
	for i, axis := range axes {
		if len(axis) == 0 {
			return nil, errors.New("every axis must have at least one element")
		}
		copied[i] = slices.Clone(axis)
	}

	Length := num_products_of(copied)

	return &ProductOf[T]{
		axes:   copied,
		Length: Length,
	}, nil
}

// All returns an iterator over all products. Each iteration yields a freshly
// allocated indices slice and items slice. Index j in the indices slice
// refers to an element of axis j.
// This code follows the algorithm from Python's itertools.product
// (https://docs.python.org/3/library/itertools.html#itertools.product)
func (p *ProductOf[T]) All() iter.Seq2[[]int, []T] {
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
func (p *ProductOf[T]) AllBorrowed() iter.Seq2[[]int, []T] {
	return func(yield func([]int, []T) bool) {
		buf := make([]T, len(p.axes))
		for inds := range p.IndicesBorrowed() {
			fillBufAxes(buf, p.axes, inds)
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
func (p *ProductOf[T]) IndicesBorrowed() iter.Seq[[]int] {
	return func(yield func([]int) bool) {
		inds := make([]int, len(p.axes))

		if !yield(inds) {
			return
		}

		for {
			// Increment the rightmost index, carrying to the left like an
			// odometer; each position rolls over at its own axis's length.
			i := len(p.axes) - 1
			for i >= 0 && inds[i] == len(p.axes[i])-1 {
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

// Nth returns the indices and items that All() would yield on its i-th
// iteration (0-based), without iterating from the start. The returned slices
// are freshly allocated and safe to retain. Returns an error if i < 0 or
// i >= Length. The argument i is not modified.
func (p *ProductOf[T]) Nth(i *big.Int) ([]int, []T, error) {
	if err := check_nth_bounds(i, p.Length); err != nil {
		return nil, nil, err
	}
	// Mixed radix with radix len(axes[j]) at position j, rightmost digit
	// fastest, matching the odometer order.
	r := new(big.Int).Set(i)
	base := new(big.Int)
	rem := new(big.Int)
	inds := make([]int, len(p.axes))
	for j := len(p.axes) - 1; j >= 0; j-- {
		base.SetInt64(int64(len(p.axes[j])))
		r.DivMod(r, base, rem)
		inds[j] = int(rem.Int64())
	}
	items := make([]T, len(inds))
	fillBufAxes(items, p.axes, inds)
	return inds, items, nil
}

// Random returns a uniform-random element of the iteration space, drawn using
// r, without enumerating it. The returned slices are freshly allocated and
// safe to retain. r must be non-nil; pass a seeded *rand.Rand for
// deterministic results.
func (p *ProductOf[T]) Random(r *rand.Rand) ([]int, []T) {
	inds, items, _ := p.Nth(random_rank(r, p.Length))
	return inds, items
}

// num_products_of returns the number of tuples in the Cartesian product of
// the axes, i.e. the product of the axis lengths.
func num_products_of[T any](axes [][]T) *big.Int {
	total := big.NewInt(1)
	for _, axis := range axes {
		total.Mul(total, big.NewInt(int64(len(axis))))
	}
	return total
}
