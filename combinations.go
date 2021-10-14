package gocombinatorics

import (
	"errors"
	"fmt"
	"log"
	"math/big"
)

func ExampleCominations_Next() {
	my_strings := []string{"apple", "banana", "cherry"}
	c, err := NewCombinations(len(my_strings), 2)
	if err != nil {
		log.Fatal(err)
	}

	for c.Next() {
		// Now c.Inds has the indices of the next combination
		fmt.Println(c.Indices())

		// It's up to you to get the elements at those indices from the slice
	}

}

// Combinations will give you the indices of all possible combinations of an input
// slice/array of length N, choosing K elements.
type Combinations struct {
	N, K    int
	isfirst bool
	Inds    []int
	Length  *big.Int
}

// // NewCombinations creates a new combinations object.
func NewCombinations(n, k int) (*Combinations, error) {
	// Check for cases where we can't do combinations
	if k > n {
		return nil, errors.New("k must be less than or equal to n")
	} else if n <= 0 {
		return nil, errors.New("n must be greater than 0")
	} else if k <= 0 {
		return nil, errors.New("k must be greater than 0")
	}
	isfirst := true
	inds := make([]int, k)
	Length := nchoosek(uint64(n), uint64(k))
	return &Combinations{n, k, isfirst, inds, Length}, nil
}

// Next will return the next combination of indices, until it reaches the end, at which
// point it will return false
// The correct indices are acces in the Inds field of the combinations object.
// This code was copied as much as possible from the python documentation itertools.combinations
// (https://docs.python.org/3/library/itertools.html#itertools.combinations)
func (c *Combinations) Next() bool {
	// If this is the first combo, just get the first k elements
	if c.isfirst {
		for i := 0; i < c.K; i++ {
			c.Inds[i] = i
		}
		c.isfirst = false
		return true
	}

	what_is_i := -1
	// Go over possible indices from k to 0 in reverse order
	for i := c.K - 1; i >= 0; i-- {
		if c.Inds[i] != i+c.N-c.K {
			what_is_i = i
			break
		} else if i == 0 {
			return false
		}
	}
	c.Inds[what_is_i] = c.Inds[what_is_i] + 1
	for j := what_is_i + 1; j < c.K; j++ {
		c.Inds[j] = c.Inds[j-1] + 1
	}
	return true

}

func (c *Combinations) LenInds() int {
	return c.K
}

func (c *Combinations) Indices() []int {
	return c.Inds
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
