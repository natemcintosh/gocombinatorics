package gocombinatorics

import (
	"errors"
	"iter"
	"math/big"
	"slices"
)

// IntegerPartitions generates all partitions of n: the ways to write n as a
// sum of positive integers, order-independent. Each partition is yielded as a
// slice of parts in non-increasing order. Use NewIntegerPartitions to create
// one, then iterate with All().
type IntegerPartitions struct {
	n      int
	Length *big.Int
}

// NewIntegerPartitions creates a new IntegerPartitions iterator. Unlike the
// other iterators it takes no input data, only the integer n to partition.
// Length is the partition function p(n).
func NewIntegerPartitions(n int) (*IntegerPartitions, error) {
	if n <= 0 {
		return nil, errors.New("n must be greater than 0")
	}

	return &IntegerPartitions{
		n:      n,
		Length: num_integer_partitions(n),
	}, nil
}

// All returns an iterator over all partitions of n, in reverse lexicographic
// order (starting from [n] and ending at [1, 1, ..., 1]). Each iteration
// yields a freshly allocated slice of parts.
func (p *IntegerPartitions) All() iter.Seq[[]int] {
	return func(yield func([]int) bool) {
		for parts := range p.AllBorrowed() {
			if !yield(slices.Clone(parts)) {
				return
			}
		}
	}
}

// AllBorrowed returns an iterator like All, but reuses an internal buffer.
// Each iteration yields the same underlying slice, overwritten in place; its
// length varies with the number of parts. The caller must not retain or
// modify the yielded slice across iterations. Use All() if you need to store
// results.
func (p *IntegerPartitions) AllBorrowed() iter.Seq[[]int] {
	return func(yield func([]int) bool) {
		parts := make([]int, 1, p.n)
		parts[0] = p.n
		if !yield(parts) {
			return
		}

		for {
			// Find the rightmost part greater than 1, counting the trailing 1s.
			ones := 0
			k := -1
			for i := len(parts) - 1; i >= 0; i-- {
				if parts[i] == 1 {
					ones++
				} else {
					k = i
					break
				}
			}
			if k < 0 {
				// All parts are 1: this was the last partition.
				return
			}

			// Decrement that part, then redistribute it plus the trailing 1s
			// into chunks no larger than the new value.
			parts[k]--
			rem := ones + 1
			parts = parts[:k+1]
			for rem > parts[k] {
				parts = append(parts, parts[k])
				rem -= parts[k]
			}
			parts = append(parts, rem)

			if !yield(parts) {
				return
			}
		}
	}
}

// num_integer_partitions returns the partition function p(n): the number of
// ways to write n as a sum of positive integers, order-independent. Computed
// with the standard O(n^2) coin-style dynamic program.
func num_integer_partitions(n int) *big.Int {
	ways := make([]*big.Int, n+1)
	ways[0] = big.NewInt(1)
	for i := 1; i <= n; i++ {
		ways[i] = big.NewInt(0)
	}
	for part := 1; part <= n; part++ {
		for s := part; s <= n; s++ {
			ways[s].Add(ways[s], ways[s-part])
		}
	}
	return ways[n]
}
