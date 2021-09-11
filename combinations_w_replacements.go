package gocombinatorics

import (
	"errors"
	"math/big"
)

// CombinationsWithReplacement will give you the indices of all possible combinations
// with replacement of an input slice/array of length N, choosing K elements.
type CombinationsWithReplacement struct {
	N, K          int
	Length        *big.Int
	Inds          []int
	current_combo *big.Int
}

func NewCombinationsWithReplacement(n, k int) (*CombinationsWithReplacement, error) {
	// Check for cases where we can't do combinations with replacement
	if k > n {
		return nil, errors.New("k must be less than or equal to n")
	} else if n <= 0 {
		return nil, errors.New("n must be greater than 0")
	} else if k <= 0 {
		return nil, errors.New("k must be greater than 0")
	}

	len := num_combinations_w_replacement(n, k)
	inds := make([]int, k)
	current_combo := big.NewInt(0)
	return &CombinationsWithReplacement{n, k, len, inds, current_combo}, nil
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
