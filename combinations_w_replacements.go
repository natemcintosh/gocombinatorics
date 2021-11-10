package gocombinatorics

import (
	"errors"
	"math/big"
)

// CombinationsWithReplacement will give you the indices of all possible combinations
// with replacement of an input slice/array of length N, choosing K elements.
type CombinationsWithReplacement struct {
	N, K    int
	Length  *big.Int
	inds    []int
	isfirst bool
}

func NewCombinationsWithReplacement(n, k int) (*CombinationsWithReplacement, error) {
	// Check for cases where we can't do combinations with replacement
	if n <= 0 {
		return nil, errors.New("n must be greater than 0")
	} else if k <= 0 {
		return nil, errors.New("k must be greater than 0")
	}

	len := num_combinations_w_replacement(n, k)
	inds := make([]int, k)
	isfirst := true
	return &CombinationsWithReplacement{n, k, len, inds, isfirst}, nil
}

// Next returns the next combination of indices until the end, and then returns false.
// The correct indices are acces in the Inds field of the combinations object.
// This code was copied as much as possible from the python documentation itertools.combinations_with_replacement
// (https://docs.python.org/3/library/itertools.html#itertools.combinations_with_replacement)
func (c *CombinationsWithReplacement) Next() bool {
	// If it's the first combo, the indices are all 0
	if c.isfirst {
		for i := 0; i < c.K; i++ {
			c.inds[i] = 0
		}
		c.isfirst = false
		return true
	}

	what_is_i := -1
	// Go over the indices from (k-1) to 0 in reverse order
	for i := c.K - 1; i >= 0; i-- {
		if c.inds[i] != c.N-1 {
			what_is_i = i
			break
		} else if i == 0 {
			return false
		}
	}
	// This for loop mimics the python list slice
	new_val := c.inds[what_is_i] + 1
	for i := what_is_i; i < c.K; i++ {
		c.inds[i] = new_val
	}
	return true
}

func (c *CombinationsWithReplacement) LenInds() int {
	return c.K
}

func (c *CombinationsWithReplacement) Indices() []int {
	return c.inds
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
