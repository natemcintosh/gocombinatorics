package gocombinatorics

import (
	"errors"
	"math/big"
	"math/rand"
)

// random_rank returns a uniform-random rank in [0, length).
func random_rank(r *rand.Rand, length *big.Int) *big.Int {
	return new(big.Int).Rand(r, length)
}

// check_nth_bounds validates that i is a valid 0-based rank, i.e. 0 <= i < length.
func check_nth_bounds(i, length *big.Int) error {
	if i.Sign() < 0 {
		return errors.New("i must be non-negative")
	}
	if i.Cmp(length) >= 0 {
		return errors.New("i must be less than Length")
	}
	return nil
}

// fillBuf writes items from data at the given indices into buf.
func fillBuf[T any](buf []T, data []T, inds []int) {
	for i, idx := range inds {
		buf[i] = data[idx]
	}
}
