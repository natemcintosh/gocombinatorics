package gocombinatorics

import (
	"errors"
	"iter"
	"math/big"
	"math/rand"
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
func (p *Permutations[T]) IndicesBorrowed() iter.Seq[[]int] {
	return func(yield func([]int) bool) {
		inds := make([]int, p.n)
		for i := range p.n {
			inds[i] = i
		}
		cycles := stepped_range(p.n, p.n-p.k, -1)

		if !yield(inds[:p.k]) {
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
					if !yield(inds[:p.k]) {
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

// Nth returns the indices and items that All() would yield on its i-th
// iteration (0-based), without iterating from the start. The returned slices
// are freshly allocated and safe to retain. Returns an error if i < 0 or
// i >= Length. The argument i is not modified.
func (p *Permutations[T]) Nth(i *big.Int) ([]int, []T, error) {
	if err := check_nth_bounds(i, p.Length); err != nil {
		return nil, nil, err
	}
	// Falling-factorial (Lehmer) radix: digit j has place value
	// (n-1-j) permute (k-1-j). The iteration order is lexicographic in the
	// index tuples, so each digit selects from the remaining available indices.
	r := new(big.Int).Set(i)
	avail := make([]int, p.n)
	for j := range avail {
		avail[j] = j
	}
	inds := make([]int, max(p.k, 0))
	for j := range inds {
		base := n_permutations(p.n-1-j, p.k-1-j)
		d := new(big.Int)
		d.DivMod(r, base, r)
		di := int(d.Int64())
		inds[j] = avail[di]
		avail = append(avail[:di], avail[di+1:]...)
	}
	items := make([]T, len(inds))
	fillBuf(items, p.data, inds)
	return inds, items, nil
}

// Random returns a uniform-random element of the iteration space, drawn using
// r, without enumerating it. The returned slices are freshly allocated and
// safe to retain. r must be non-nil; pass a seeded *rand.Rand for
// deterministic results.
func (p *Permutations[T]) Random(r *rand.Rand) ([]int, []T) {
	inds, items, _ := p.Nth(random_rank(r, p.Length))
	return inds, items
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
