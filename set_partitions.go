package gocombinatorics

import (
	"errors"
	"iter"
	"math/big"
	"slices"
)

// SetPartitions generates all ways to divide the input data into non-empty,
// disjoint blocks. With NewSetPartitions the number of partitions is the Bell
// number B(n); with NewSetPartitionsK (exactly k blocks) it is the Stirling
// number of the second kind S(n, k). Use one of those constructors, then
// iterate with All().
type SetPartitions[T any] struct {
	data   []T
	n      int
	k      int // 0 means any number of blocks
	Length *big.Int
}

// NewSetPartitions creates a SetPartitions iterator over all partitions of
// input_data, into any number of blocks. Length is the Bell number B(n).
func NewSetPartitions[T any](input_data []T) (*SetPartitions[T], error) {
	data := make([]T, len(input_data))
	copy(data, input_data)
	n := len(input_data)

	if n <= 0 {
		return nil, errors.New("len(input_data) must be greater than 0")
	}

	return &SetPartitions[T]{
		data:   data,
		n:      n,
		k:      0,
		Length: bell_number(n),
	}, nil
}

// NewSetPartitionsK creates a SetPartitions iterator over the partitions of
// input_data into exactly k non-empty blocks. Length is the Stirling number
// of the second kind S(n, k).
func NewSetPartitionsK[T any](input_data []T, k int) (*SetPartitions[T], error) {
	data := make([]T, len(input_data))
	copy(data, input_data)
	n := len(input_data)

	if k > n {
		return nil, errors.New("k must be less than or equal to len(input_data)")
	} else if n <= 0 {
		return nil, errors.New("len(input_data) must be greater than 0")
	} else if k <= 0 {
		return nil, errors.New("k must be greater than 0")
	}

	return &SetPartitions[T]{
		data:   data,
		n:      n,
		k:      k,
		Length: stirling2(n, k),
	}, nil
}

// All returns an iterator over all set partitions. Each iteration yields a
// freshly allocated assignment slice (assignment[i] is the block index of
// element i, a restricted growth string) and a freshly allocated slice of
// blocks. Blocks are ordered by their smallest element, and elements within a
// block keep their input order.
func (s *SetPartitions[T]) All() iter.Seq2[[]int, [][]T] {
	return func(yield func([]int, [][]T) bool) {
		for assignment, blocks := range s.AllBorrowed() {
			cloned := make([][]T, len(blocks))
			for i, b := range blocks {
				cloned[i] = slices.Clone(b)
			}
			if !yield(slices.Clone(assignment), cloned) {
				return
			}
		}
	}
}

// AllBorrowed returns an iterator like All, but reuses internal buffers.
// Each iteration yields the same underlying assignment slice and block
// slices, overwritten in place; the number of blocks and their lengths vary.
// The caller must not retain or modify the yielded slices across iterations.
// Use All() if you need to store results.
func (s *SetPartitions[T]) AllBorrowed() iter.Seq2[[]int, [][]T] {
	return func(yield func([]int, [][]T) bool) {
		blockBufs := make([][]T, s.n)
		for i := range blockBufs {
			blockBufs[i] = make([]T, 0, s.n)
		}
		for assignment := range s.AssignmentsBorrowed() {
			numBlocks := 0
			for _, b := range assignment {
				if b+1 > numBlocks {
					numBlocks = b + 1
				}
			}
			blocks := blockBufs[:numBlocks]
			for i := range blocks {
				blocks[i] = blocks[i][:0]
			}
			for i, b := range assignment {
				blocks[b] = append(blocks[b], s.data[i])
			}
			if !yield(assignment, blocks) {
				return
			}
		}
	}
}

// AssignmentsBorrowed yields the same assignment sequence as AllBorrowed, but
// never gathers items into blocks. Each assignment is a restricted growth
// string: assignment[0] == 0 and each later entry is at most one greater than
// the maximum of the entries before it. The yielded slice is reused
// (borrowed) across iterations; clone it if you need to retain it. This is
// the fast path for callers that only need the combinatorial structure, not
// the items.
func (s *SetPartitions[T]) AssignmentsBorrowed() iter.Seq[[]int] {
	return func(yield func([]int) bool) {
		assignment := make([]int, s.n)
		s.next_assignment(assignment, 1, 0, yield)
	}
}

// next_assignment recursively fills assignment[i:] with every valid
// restricted growth string suffix, given that the maximum entry so far is m.
// When k is fixed, branches that cannot end with exactly k blocks are pruned.
// Returns false once yield does, to stop the whole recursion.
func (s *SetPartitions[T]) next_assignment(assignment []int, i, m int, yield func([]int) bool) bool {
	if i == s.n {
		return yield(assignment)
	}
	for v := 0; v <= m+1; v++ {
		if s.k > 0 && v > s.k-1 {
			break
		}
		newMax := max(m, v)
		// Each remaining position can open at most one new block, so prune
		// branches that can no longer reach k blocks.
		if s.k > 0 && s.k-1-newMax > s.n-1-i {
			continue
		}
		assignment[i] = v
		if !s.next_assignment(assignment, i+1, newMax, yield) {
			return false
		}
	}
	return true
}

// bell_number returns the Bell number B(n): the number of partitions of an
// n-element set. Computed with the Bell triangle.
func bell_number(n int) *big.Int {
	row := []*big.Int{big.NewInt(1)}
	for i := 1; i < n; i++ {
		next := make([]*big.Int, i+1)
		next[0] = row[len(row)-1]
		for j := 1; j <= i; j++ {
			next[j] = new(big.Int).Add(next[j-1], row[j-1])
		}
		row = next
	}
	return row[len(row)-1]
}

// stirling2 returns the Stirling number of the second kind S(n, k): the
// number of partitions of an n-element set into exactly k non-empty blocks.
// Uses the recurrence S(n, k) = k*S(n-1, k) + S(n-1, k-1).
func stirling2(n, k int) *big.Int {
	row := make([]*big.Int, k+1)
	row[0] = big.NewInt(1)
	for j := 1; j <= k; j++ {
		row[j] = big.NewInt(0)
	}
	for i := 1; i <= n; i++ {
		for j := min(i, k); j >= 1; j-- {
			row[j].Mul(row[j], big.NewInt(int64(j)))
			row[j].Add(row[j], row[j-1])
		}
		row[0] = big.NewInt(0)
	}
	return row[k]
}
