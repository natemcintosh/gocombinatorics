package gocombinatorics

import (
	"errors"
	"math/big"
)

// Permutations allows the user to iteratively generate all permutations of a certain length
// for any slice of data. Instantiate it with `NewPermutations`, iterate with the `.Next()`
// method, and access the data with the `.Items()` method.
// Permutations meets the `CombinationLike` interface
type Permutations[T any] struct {
	data    []T
	n, k    int
	Length  *big.Int
	inds    []int
	cycles  []int
	isfirst bool
	buffer  []T
}

// NewPermutations will return an instance of the `Permutations` struct
func NewPermutations[T any](input_data []T, k int) (*Permutations[T], error) {
	data := make([]T, len(input_data))
	copy(data, input_data)
	n := len(input_data)
	if k > n {
		return nil, errors.New("k must be less than or equal to len(input_data)")
	}
	Len := n_permutations(n, k)
	inds := make([]int, n)
	for i := 0; i < n; i++ {
		inds[i] = i
	}

	// A cycles slice
	cycles := stepped_range(n, n-k, -1)
	isfirst := true

	// The buffer slice
	buffer := make([]T, len(inds))
	fill_buffer(buffer, data, inds)

	// Return the Permutations struct
	return &Permutations[T]{data, n, k, Len, inds, cycles, isfirst, buffer}, nil
}

// Next will return true if there is another iteration to go, and false if not. It will
// update the state of the Permutations struct. The new permutation can be accessed with
// p.Items()
func (p *Permutations[T]) Next() bool {
	// Check if we're at the first permutation
	if p.isfirst {
		// Update inds with 1,...,k
		for i := 0; i < p.k; i++ {
			p.inds[i] = i
		}
		p.isfirst = false
		return true
	}

	for i := p.k - 1; i >= 0; i-- {
		p.cycles[i] -= 1
		if p.cycles[i] == 0 {
			// Move item at i to the end of the slice
			// Grab the element at i
			ith_elt := p.inds[i]

			// Delete the element at i
			p.inds = append(p.inds[:i], p.inds[i+1:]...)

			// Append ith_elt to the end of the slice
			p.inds = append(p.inds, ith_elt)

			p.cycles[i] = p.n - i
		} else {
			j := p.cycles[i]
			// BUG: the python code is indices[i], indices[-j] = indices[-j], indices[i]
			// And aparently the code below does not mimic it how I thought it did
			new_at_i := p.inds[len(p.inds)-j]
			new_at_minus_j := p.inds[i]
			p.inds[i] = new_at_i
			p.inds[len(p.inds)-j] = new_at_minus_j
			return true
		}
	}
	return false
}

// Indices tells you the current indices used to get this permutation
func (p *Permutations[T]) Indices() []int {
	return p.inds[:p.k]
}

// LenInds gives you how many items you want in each permutation
func (p *Permutations[T]) LenInds() int {
	return p.k
}

// Items is how you get the items in this permutation. You iterate with `p.Next()`, and
// then get the permutation with `p.Items()`. The data in the slice returned will be
// overwritten every iteration. If you need to keep the data from each iteration, be
// sure to make a copy.
func (p *Permutations[T]) Items() []T {
	fill_buffer(p.buffer, p.data, p.inds)
	return p.buffer
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

	total_perms := n_permutations(int(n), int(k))
	n_minus_1_perms := n_permutations(int(n-1), int(k))
	return big.NewInt(0).Sub(total_perms, n_minus_1_perms)
}

// Mimics python's range() with a step argument
func stepped_range(start int, stop int, step int) []int {
	if step == 0 {
		return make([]int, 0)
	}
	result := make([]int, 0)
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
