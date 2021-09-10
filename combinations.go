package gocombinatorics

import (
	"math/big"
)

// Combinations will give you the indices of all possible combinations of an input
// slice/array of length N, choosing K elements.
// type Combinations struct {
// 	N, K   uint64
// 	Length *big.Int
// }

// // NewCombinations creates a new combinations object.
// func NewCombinations(N, K int) *Combinations {
// 	len := nchoosek(N, K)
// 	return &Combinations{N, K, len}
// }

// nchoosek returns the number of combinations of n things taken k at a time.
// nchoosek(n, k) = n! / (k! * (n-k)!) if n > k
// nchoosek(n, k) = 0 if n < k
// nchoosek(n, k) = 1 if n == k
// nchoosek(n, k) = 0 if n <= 0 or k= < 0
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
	var fact big.Int
	fact.MulRange(1, n)
	return &fact
}
