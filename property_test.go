// property_test.go offers a very simple property test on the available combinatorics
// functionality.

// If we want all possible length 2 combinations of 10 items, then we can expect to get
// n choose k = n! / (k! * (n-k)!) = 10! / (2! * (10-2)!) = 45. In those 45 cases, we
// should expect to see each number exactly 45 - (9 choose 2) = 45 - 36 = 9 times.

// The above is an example with combinations. We can possibly also do them for other
// functionality as well
package gocombinatorics

import "testing"

// Multiple types all adhere to this interface
type CombinationLike interface {
	Next() bool
	LenInds() int
	Indices() []int
}

// combinationLikeValueCounter will iterate through something combination like and count
// each unique value, returning it in a map.
// It assumes that the combinationLike has already been created
func combinationLikeValueCounter(c CombinationLike) map[int]int {
	// Make the result map
	result := make(map[int]int, c.LenInds())

	// Iterate through the object
	for c.Next() {
		// Add each item in c.Inds to the result map
		for _, num := range c.Indices() {
			result[num]++
		}
	}

	return result

}

func TestCombinationsProperties(t *testing.T) {
	testCases := []struct {
		desc            string
		n               int
		k               int
		num_want_to_see int
	}{
		{
			desc:            "n=3, k=2",
			n:               3,
			k:               2,
			num_want_to_see: 2,
		},
		{
			desc:            "n=10, k=2",
			n:               10,
			k:               2,
			num_want_to_see: 9,
		},
		{
			desc:            "n=10, k=7",
			n:               10,
			k:               7,
			num_want_to_see: 120 - 36,
		},
		{
			desc:            "n=50, k=5",
			n:               50,
			k:               5,
			num_want_to_see: 211876,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// Create the combination
			c, err := NewCombinations(tC.n, tC.k)
			if err != nil {
				t.Errorf("Error creating combinations: %v", err)
			}

			// Count the number of times each item appears
			counts := combinationLikeValueCounter(c)

			// Check that each value in counts appears t.num_want_to_see times
			for num, count := range counts {
				if count != tC.num_want_to_see {
					t.Errorf("Expected %v to appear %v times, but it appeared %v times", num, tC.num_want_to_see, count)
				}
			}
		})
	}
}
