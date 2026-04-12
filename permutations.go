package gocombinatorics

import (
	"errors"
	"iter"
	"math/big"
	"slices"
)

// Permutations generates all k-length permutations of the input data.
// Use NewPermutations to create one, then iterate with All().
type Permutations[T any] struct {
	data   []T
	n, k   int
	Length *big.Int
}

// NewPermutations creates a new Permutations iterator.
func NewPermutations[T any](input_data []T, k int) (*Permutations[T], error) {
	data := make([]T, len(input_data))
	copy(data, input_data)
	n := len(input_data)
	if k > n {
		return nil, errors.New("k must be less than or equal to len(input_data)")
	}
	Length := n_permutations(n, k)

	return &Permutations[T]{
		data:   data,
		n:      n,
		k:      k,
		Length: Length,
	}, nil
}

// All returns an iterator over all permutations. Each iteration yields
// a freshly allocated indices slice and items slice.
// This code follows the algorithm from Python's itertools.permutations
// (https://docs.python.org/3/library/itertools.html#itertools.permutations)
func (p *Permutations[T]) All() iter.Seq2[[]int, []T] {
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
func (p *Permutations[T]) AllBorrowed() iter.Seq2[[]int, []T] {
	return func(yield func([]int, []T) bool) {
		inds := make([]int, p.n)
		for i := range p.n {
			inds[i] = i
		}
		cycles := stepped_range(p.n, p.n-p.k, -1)
		buf := make([]T, p.k)

		fillBuf(buf, p.data, inds[:p.k])
		if !yield(inds[:p.k], buf) {
			return
		}

		for {
			found := false
			for i := p.k - 1; i >= 0; i-- {
				cycles[i]--
				if cycles[i] == 0 {
					// Rotate element at i to the end
					ith := inds[i]
					copy(inds[i:], inds[i+1:])
					inds[len(inds)-1] = ith
					cycles[i] = p.n - i
				} else {
					j := cycles[i]
					inds[i], inds[len(inds)-j] = inds[len(inds)-j], inds[i]
					fillBuf(buf, p.data, inds[:p.k])
					if !yield(inds[:p.k], buf) {
						return
					}
					found = true
					break
				}
			}
			if !found {
				return
			}
		}
	}
}

func n_permutations(n, k int) *big.Int {
	numerator := factorial(int64(n))
	denominator := factorial(int64(n - k))
	result := new(big.Int).Div(numerator, denominator)
	return result
}

func elts_in_permutations(n, k int) *big.Int {
	if n == k {
		return n_permutations(n, k)
	}
	total_perms := n_permutations(n, k)
	n_minus_1_perms := n_permutations(n-1, k)
	return big.NewInt(0).Sub(total_perms, n_minus_1_perms)
}

// Mimics python's range() with a step argument
func stepped_range(start int, stop int, step int) []int {
	approx_size := (stop - start) / step
	if step == 0 {
		return make([]int, 0, approx_size)
	}
	result := make([]int, 0, approx_size)
	val := -1
	for {
		val++
		new_val := start + (val * step)
		if (step > 0) && (new_val >= stop) {
			break
		} else if (step < 0) && (new_val <= stop) {
			break
		} else {
			result = append(result, new_val)
		}
	}
	return result
}
