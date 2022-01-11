// property_test.go offers a very simple property test on the available combinatorics
// functionality.

// If we want all possible length 2 combinations of 10 items, then we can expect to get
// n choose k = n! / (k! * (n-k)!) = 10! / (2! * (10-2)!) = 45. In those 45 cases, we
// should expect to see each number exactly 45 - (9 choose 2) = 45 - 36 = 9 times.

// The above is an example with combinations. We can possibly also do them for other
// functionality as well
package gocombinatorics

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"
)

// Multiple types all adhere to this interface
type combinationLike interface {
	Next() bool
	LenInds() int
	Indices() []int
}

// combinationLikeValueCounter will iterate through something combination like and count
// each unique value, returning it in a map.
// It assumes that the combinationLike has already been created
func combinationLikeValueCounter(c combinationLike) map[int]int {
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
			data := stepped_range(0, tC.n, 1)
			c, err := NewCombinations(data, tC.k)
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

func Test100RandomCombinations(t *testing.T) {
	for i := 0; i < 100; i++ {

		// Generate two numbers, 1 <= n <= 50, 1 <= k <= n
		n := rand.Int63n(50) + 1
		k := rand.Int63n(n) + 1

		// If any index should appear more than 10_000_000 times, skip this iteration
		total_times := nchoosek(uint64(n), uint64(k))
		one_less := nchoosek(uint64(n-1), uint64(k))
		times_we_see_each_index := big.NewInt(0).Sub(total_times, one_less)
		if times_we_see_each_index.Cmp(big.NewInt(10000000)) > 0 {
			t.Logf("Skipping test because we see each index more than 10_000_000 times")
			continue
		}

		run_name := fmt.Sprintf("n=%v, k=%v", n, k)
		t.Run(run_name, func(t *testing.T) {
			// Create the combination
			data := stepped_range(0, int(n), 1)
			c, err := NewCombinations(data, int(k))
			if err != nil {
				t.Errorf("Error creating combinations: %v", err)
			}

			// Count the number of times each item appears
			counts := combinationLikeValueCounter(c)

			// Check that each value in counts appears t.num_want_to_see times
			for num, count := range counts {
				count_big := big.NewInt(int64(count))
				if count_big.Cmp(times_we_see_each_index) != 0 {
					t.Errorf("Expected %v to appear %v times, but it appeared %v times", num, times_we_see_each_index, count)
				}
			}
		})
	}

}

func Test100RandomCombinationsWithReplacement(t *testing.T) {
	for i := 0; i < 100; i++ {

		// Generate two numbers, 1 <= n <= 50, 1 <= k <= n
		n := rand.Int63n(15) + 1
		k := rand.Int63n(n) + 1

		// If any index should appear more than 10_000_000 times, skip this iteration
		times_we_see_each_index := elts_in_combo_w_replacement(int(n), int(k))
		if times_we_see_each_index.Cmp(big.NewInt(10000000)) > 0 {
			t.Logf("Skipping test because we see each index more than 10_000_000 times")
			continue
		}

		run_name := fmt.Sprintf("n=%v, k=%v", n, k)
		t.Run(run_name, func(t *testing.T) {
			// Create the combination
			c, err := NewCombinationsWithReplacement(int(n), int(k))
			if err != nil {
				t.Errorf("Error creating CombinationsWithReplacement: %v", err)
			}

			// Count the number of times each item appears
			counts := combinationLikeValueCounter(c)

			// Check that each value in counts appears t.num_want_to_see times
			for num, count := range counts {
				count_big := big.NewInt(int64(count))
				if count_big.Cmp(times_we_see_each_index) != 0 {
					t.Errorf("Expected %v to appear %v times, but it appeared %v times", num, times_we_see_each_index, count)
				}
			}
		})
	}

}

func Test100RandomPermutations(t *testing.T) {
	// Do 100 iterations
	for i := 0; i < 100; i++ {
		// Generate two numbers, 1 <= n <= 50, 1 <= k <= n
		n := rand.Int63n(15) + 1
		k := rand.Int63n(n) + 1
		data := stepped_range(0, int(n), 1)

		// If any index should appear more than 10_000_000 times, skip this iteration
		times_we_see_each_index := elts_in_permutations(int(n), int(k))
		if times_we_see_each_index.Cmp(big.NewInt(10000000)) > 0 {
			t.Logf("Skipping test because we see each index more than 10_000_000 times")
			continue
		}

		run_name := fmt.Sprintf("n=%v, k=%v", n, k)
		t.Run(run_name, func(t *testing.T) {
			// Create the permutation
			p, err := NewPermutations(data, int(k))
			if err != nil {
				t.Errorf("Error creating Permutations: %v", err)
			}

			// Count the number of times each item appears
			var counts map[int]int = combinationLikeValueCounter(p)

			// Check that each value in counts appears t.num_want_to_see times
			for num, count := range counts {
				count_big := big.NewInt(int64(count))
				if count_big.Cmp(times_we_see_each_index) != 0 {
					t.Errorf("Expected %v to appear %v times, but it appeared %v times", num, times_we_see_each_index, count)
				}
			}
		})
	}
}
