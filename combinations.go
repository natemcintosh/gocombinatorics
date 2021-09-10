package gocombinatorics

import (
	"errors"
	"math/big"
)

var EndOfCombinations = errors.New("End of combinations")

// Combinations will give you the indices of all possible combinations of an input
// slice/array of length N, choosing K elements.
type Combinations struct {
	N, K          uint64
	Length        *big.Int
	Inds          []int
	current_combo *big.Int
}

// // NewCombinations creates a new combinations object.
func NewCombinations(n uint64, k uint64) (*Combinations, error) {
	// Check for cases where we can't do combinations
	if k > n {
		return nil, errors.New("k must be less than or equal to n")
	} else if n <= 0 {
		return nil, errors.New("n must be greater than 0")
	} else if k <= 0 {
		return nil, errors.New("k must be greater than 0")
	}

	len := nchoosek(n, k)
	inds := make([]int, k)
	current_combo := big.NewInt(0)
	return &Combinations{n, k, len, inds, current_combo}, nil
}

// nchoosek returns the number of combinations of n things taken k at a time.
// nchoosek(n, k) = n! / (k! * (n-k)!) if n > k
// nchoosek(n, k) = 0 if n < k
// nchoosek(n, k) = 1 if n == k
// nchoosek(n, k) = 0 if n <= 0 or k <= 0
func nchoosek(n, k uint64) *big.Int {
	if n <= 0 || k <= 0 {
		return big.NewInt(0)
	} else if k > n {
		// nchoosek(n, k) = 0 if n < k
		return big.NewInt(0)
	} else if k == n {
		// nchoosek(n, k) = 1 if n == k
		return big.NewInt(1)
	}
	// Calculate the numerator
	numerator := factorial(int64(n))

	// Calculate the denominator
	kfact := factorial(int64(k))
	nminuskfact := factorial(int64(n - k))
	denominator := big.NewInt(0)
	denominator = denominator.Mul(kfact, nminuskfact)

	// Calculate the result
	result := big.NewInt(0)
	return result.Div(numerator, denominator)
}

// factorial returns the factorial of a number, i.e. n! = n * (n-1) * (n-2) * ... * 1
func factorial(n int64) *big.Int {
	fact := big.NewInt(0)
	fact.MulRange(1, n)
	return fact
}
