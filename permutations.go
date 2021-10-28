package gocombinatorics

import (
	"errors"
	"math/big"
)

type Permutations struct {
	N, K    int
	Length  *big.Int
	Inds    []int
	cycles  []int
	isfirst bool
}

func NewPermutations(n, k int) (*Permutations, error) {
	if k > n {
		return nil, errors.New("k must be less than or equal to n")
	}
	Len := n_permutations(n, k)
	Inds := make([]int, n)
	for i := 0; i < n; i++ {
		Inds[i] = i
	}

	// A cycles slice
	cycles := stepped_range(n, n-k, -1)
	isfirst := true
	return &Permutations{n, k, Len, Inds, cycles, isfirst}, nil
}

func (p *Permutations) Next() bool {
	// Check if we're at the first permutation
	if p.isfirst {
		// Update inds with 1,...,k
		for i := 0; i < p.K; i++ {
			p.Inds[i] = i
		}
		p.isfirst = false
		return true
	}

	for i := p.K - 1; i >= 0; i-- {
		p.cycles[i] -= 1
		if p.cycles[i] == 0 {
			// Move item at i to the end of the slice
			// Grab the element at i
			ith_elt := p.Inds[i]

			// Delete the element at i
			p.Inds = append(p.Inds[:i], p.Inds[i+1:]...)

			// Append ith_elt to the end of the slice
			p.Inds = append(p.Inds, ith_elt)

			p.cycles[i] = p.N - i
		} else {
			j := p.cycles[i]
			// BUG: the python code is indices[i], indices[-j] = indices[-j], indices[i]
			// And aparently the code below does not mimic it how I thought it did
			new_at_i := p.Inds[len(p.Inds)-j]
			new_at_minus_j := p.Inds[i]
			p.Inds[i] = new_at_i
			p.Inds[len(p.Inds)-j] = new_at_minus_j
			return true
		}
	}
	return false
}

func (p *Permutations) Indices() []int {
	return p.Inds[:p.K]
}

func (p *Permutations) LenInds() int {
	return p.K
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
