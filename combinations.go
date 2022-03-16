package gocombinatorics

import (
	"errors"
	"math/big"
)

// Combinations will give you the indices of all possible combinations of an input
// slice/array of length n, choosing k elements.
type Combinations[T any] struct {
	data    []T
	n, k    int
	isfirst bool
	inds    []int
	Length  *big.Int
	buffer  []T
}

// // NewCombinations creates a new combinations object.
func NewCombinations[T any](input_data []T, k int) (*Combinations[T], error) {
	data := make([]T, len(input_data))
	copy(data, input_data)
	n := len(input_data)

	// Check for cases where we can't do combinations
	if k > n {
		return nil, errors.New("k must be less than or equal to len(input_data)")
	} else if n <= 0 {
		return nil, errors.New("len(input_data) must be greater than 0")
	} else if k <= 0 {
		return nil, errors.New("k must be greater than 0")
	}
	isfirst := true
	inds := make([]int, k)
	Length := nchoosek(uint64(n), uint64(k))

	// Make the buffer slice
	buffer := make([]T, k)
	fill_buffer(buffer, data, inds)

	return &Combinations[T]{
		data:    data,
		n:       n,
		k:       k,
		isfirst: isfirst,
		inds:    inds,
		Length:  Length,
		buffer:  buffer,
	}, nil
}

// Next will return the next combination of indices, until it reaches the end, at which
// point it will return false
// The correct indices are acces in the Inds field of the combinations object.
// This code was copied as much as possible from the python documentation itertools.combinations
// (https://docs.python.org/3/library/itertools.html#itertools.combinations)
func (c *Combinations[T]) Next() bool {
	// If this is the first combo, just get the first k elements
	if c.isfirst {
		for i := 0; i < c.k; i++ {
			c.inds[i] = i
		}
		c.isfirst = false
		return true
	}

	what_is_i := -1
	// Go over possible indices from k to 0 in reverse order
	for i := c.k - 1; i >= 0; i-- {
		if c.inds[i] != i+c.n-c.k {
			what_is_i = i
			break
		} else if i == 0 {
			return false
		}
	}
	c.inds[what_is_i]++
	for j := what_is_i + 1; j < c.k; j++ {
		c.inds[j] = c.inds[j-1] + 1
	}
	return true

}

func (c *Combinations[T]) LenInds() int {
	return c.k
}

func (c *Combinations[T]) Indices() []int {
	return c.inds
}

// Items is how you get the items in this combination. You iterate with `c.Next()`, and
// then get the combination with `c.Items()`. The data in the slice returned will be
// overwritten every iteration. If you need to keep the data from each iteration, be
// sure to make a copy.
func (c *Combinations[T]) Items() []T {
	fill_buffer(c.buffer, c.data, c.inds)
	return c.buffer
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
