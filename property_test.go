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
	"iter"
	"math/big"
	"math/rand"
	"testing"
)

// indexCounter iterates through a Seq2 iterator, counting how many times each index
// appears across all yielded index slices.
func indexCounter[T any](seq iter.Seq2[[]int, []T]) map[int]int {
	result := make(map[int]int)
	for inds := range seq {
		for _, idx := range inds {
			result[idx]++
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
			data := stepped_range(0, tC.n, 1)
			c, err := NewCombinations(data, tC.k)
			if err != nil {
				t.Errorf("Error creating combinations: %v", err)
			}

			counts := indexCounter(c.All())

			for num, count := range counts {
				if count != tC.num_want_to_see {
					t.Errorf("Expected %v to appear %v times, but it appeared %v times", num, tC.num_want_to_see, count)
				}
			}
		})
	}
}

func Test100RandomCombinations(t *testing.T) {
	for range 100 {

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
			data := stepped_range(0, int(n), 1)
			c, err := NewCombinations(data, int(k))
			if err != nil {
				t.Errorf("Error creating combinations: %v", err)
			}

			counts := indexCounter(c.All())

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
	for range 100 {

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
			data := stepped_range(0, int(n), 1)
			c, err := NewCombinationsWithReplacement(data, int(k))
			if err != nil {
				t.Errorf("Error creating CombinationsWithReplacement: %v", err)
			}

			counts := indexCounter(c.All())

			for num, count := range counts {
				count_big := big.NewInt(int64(count))
				if count_big.Cmp(times_we_see_each_index) != 0 {
					t.Errorf("Expected %v to appear %v times, but it appeared %v times", num, times_we_see_each_index, count)
				}
			}
		})
	}

}

func Test100RandomCombinationsBorrowed(t *testing.T) {
	for range 100 {
		n := rand.Int63n(50) + 1
		k := rand.Int63n(n) + 1

		total_times := nchoosek(uint64(n), uint64(k))
		one_less := nchoosek(uint64(n-1), uint64(k))
		times_we_see_each_index := big.NewInt(0).Sub(total_times, one_less)
		if times_we_see_each_index.Cmp(big.NewInt(10000000)) > 0 {
			t.Logf("Skipping test because we see each index more than 10_000_000 times")
			continue
		}

		run_name := fmt.Sprintf("n=%v, k=%v", n, k)
		t.Run(run_name, func(t *testing.T) {
			data := stepped_range(0, int(n), 1)
			c, err := NewCombinations(data, int(k))
			if err != nil {
				t.Errorf("Error creating combinations: %v", err)
			}

			counts := indexCounter(c.AllBorrowed())

			for num, count := range counts {
				count_big := big.NewInt(int64(count))
				if count_big.Cmp(times_we_see_each_index) != 0 {
					t.Errorf("Expected %v to appear %v times, but it appeared %v times", num, times_we_see_each_index, count)
				}
			}
		})
	}
}

func Test100RandomCombinationsWithReplacementBorrowed(t *testing.T) {
	for range 100 {
		n := rand.Int63n(15) + 1
		k := rand.Int63n(n) + 1

		times_we_see_each_index := elts_in_combo_w_replacement(int(n), int(k))
		if times_we_see_each_index.Cmp(big.NewInt(10000000)) > 0 {
			t.Logf("Skipping test because we see each index more than 10_000_000 times")
			continue
		}

		run_name := fmt.Sprintf("n=%v, k=%v", n, k)
		t.Run(run_name, func(t *testing.T) {
			data := stepped_range(0, int(n), 1)
			c, err := NewCombinationsWithReplacement(data, int(k))
			if err != nil {
				t.Errorf("Error creating CombinationsWithReplacement: %v", err)
			}

			counts := indexCounter(c.AllBorrowed())

			for num, count := range counts {
				count_big := big.NewInt(int64(count))
				if count_big.Cmp(times_we_see_each_index) != 0 {
					t.Errorf("Expected %v to appear %v times, but it appeared %v times", num, times_we_see_each_index, count)
				}
			}
		})
	}
}

func Test100RandomProduct(t *testing.T) {
	for range 100 {
		// 1 <= n <= 8, 1 <= k <= 6
		n := rand.Int63n(8) + 1
		k := rand.Int63n(6) + 1

		times_we_see_each_index := elts_in_product(int(n), int(k))
		if times_we_see_each_index.Cmp(big.NewInt(10000000)) > 0 {
			t.Logf("Skipping test because we see each index more than 10_000_000 times")
			continue
		}

		run_name := fmt.Sprintf("n=%v, k=%v", n, k)
		t.Run(run_name, func(t *testing.T) {
			data := stepped_range(0, int(n), 1)
			p, err := NewProduct(data, int(k))
			if err != nil {
				t.Errorf("Error creating Product: %v", err)
			}

			counts := indexCounter(p.All())

			for num, count := range counts {
				count_big := big.NewInt(int64(count))
				if count_big.Cmp(times_we_see_each_index) != 0 {
					t.Errorf("Expected %v to appear %v times, but it appeared %v times", num, times_we_see_each_index, count)
				}
			}
		})
	}
}

func Test100RandomProductBorrowed(t *testing.T) {
	for range 100 {
		n := rand.Int63n(8) + 1
		k := rand.Int63n(6) + 1

		times_we_see_each_index := elts_in_product(int(n), int(k))
		if times_we_see_each_index.Cmp(big.NewInt(10000000)) > 0 {
			t.Logf("Skipping test because we see each index more than 10_000_000 times")
			continue
		}

		run_name := fmt.Sprintf("n=%v, k=%v", n, k)
		t.Run(run_name, func(t *testing.T) {
			data := stepped_range(0, int(n), 1)
			p, err := NewProduct(data, int(k))
			if err != nil {
				t.Errorf("Error creating Product: %v", err)
			}

			counts := indexCounter(p.AllBorrowed())

			for num, count := range counts {
				count_big := big.NewInt(int64(count))
				if count_big.Cmp(times_we_see_each_index) != 0 {
					t.Errorf("Expected %v to appear %v times, but it appeared %v times", num, times_we_see_each_index, count)
				}
			}
		})
	}
}

// elts_in_product_of returns how many times we expect to see each index value
// across the whole product. Axis i contributes Length/len(axes[i]) occurrences
// of each value 0..len(axes[i])-1, so the expected count for value v is the
// sum of that over every axis long enough to contain v.
func elts_in_product_of(axis_lens []int) map[int]*big.Int {
	length := big.NewInt(1)
	for _, n := range axis_lens {
		length.Mul(length, big.NewInt(int64(n)))
	}
	expected := make(map[int]*big.Int)
	for _, n := range axis_lens {
		per_value := new(big.Int).Div(length, big.NewInt(int64(n)))
		for v := range n {
			if _, ok := expected[v]; !ok {
				expected[v] = big.NewInt(0)
			}
			expected[v].Add(expected[v], per_value)
		}
	}
	return expected
}

func random_product_of_axes() ([][]int, []int) {
	// 1 <= number of axes <= 5, each with 1 <= length <= 8
	n_axes := rand.Int63n(5) + 1
	axes := make([][]int, n_axes)
	axis_lens := make([]int, n_axes)
	for i := range axes {
		n := int(rand.Int63n(8) + 1)
		axes[i] = stepped_range(0, n, 1)
		axis_lens[i] = n
	}
	return axes, axis_lens
}

func Test100RandomProductOf(t *testing.T) {
	for range 100 {
		axes, axis_lens := random_product_of_axes()
		expected := elts_in_product_of(axis_lens)

		run_name := fmt.Sprintf("axis_lens=%v", axis_lens)
		t.Run(run_name, func(t *testing.T) {
			p, err := NewProductOf(axes...)
			if err != nil {
				t.Errorf("Error creating ProductOf: %v", err)
			}

			counts := indexCounter(p.All())

			for num, count := range counts {
				count_big := big.NewInt(int64(count))
				if count_big.Cmp(expected[num]) != 0 {
					t.Errorf("Expected %v to appear %v times, but it appeared %v times", num, expected[num], count)
				}
			}
		})
	}
}

func Test100RandomProductOfBorrowed(t *testing.T) {
	for range 100 {
		axes, axis_lens := random_product_of_axes()
		expected := elts_in_product_of(axis_lens)

		run_name := fmt.Sprintf("axis_lens=%v", axis_lens)
		t.Run(run_name, func(t *testing.T) {
			p, err := NewProductOf(axes...)
			if err != nil {
				t.Errorf("Error creating ProductOf: %v", err)
			}

			counts := indexCounter(p.AllBorrowed())

			for num, count := range counts {
				count_big := big.NewInt(int64(count))
				if count_big.Cmp(expected[num]) != 0 {
					t.Errorf("Expected %v to appear %v times, but it appeared %v times", num, expected[num], count)
				}
			}
		})
	}
}

func Test100RandomPowerset(t *testing.T) {
	for range 100 {
		// 1 <= n <= 15
		n := rand.Int63n(15) + 1

		times_we_see_each_index := elts_in_powerset(int(n))
		if times_we_see_each_index.Cmp(big.NewInt(10000000)) > 0 {
			t.Logf("Skipping test because we see each index more than 10_000_000 times")
			continue
		}

		run_name := fmt.Sprintf("n=%v", n)
		t.Run(run_name, func(t *testing.T) {
			data := stepped_range(0, int(n), 1)
			p, err := NewPowerset(data)
			if err != nil {
				t.Errorf("Error creating Powerset: %v", err)
			}

			counts := indexCounter(p.All())

			for num, count := range counts {
				count_big := big.NewInt(int64(count))
				if count_big.Cmp(times_we_see_each_index) != 0 {
					t.Errorf("Expected %v to appear %v times, but it appeared %v times", num, times_we_see_each_index, count)
				}
			}
		})
	}
}

func Test100RandomPowersetBorrowed(t *testing.T) {
	for range 100 {
		n := rand.Int63n(15) + 1

		times_we_see_each_index := elts_in_powerset(int(n))
		if times_we_see_each_index.Cmp(big.NewInt(10000000)) > 0 {
			t.Logf("Skipping test because we see each index more than 10_000_000 times")
			continue
		}

		run_name := fmt.Sprintf("n=%v", n)
		t.Run(run_name, func(t *testing.T) {
			data := stepped_range(0, int(n), 1)
			p, err := NewPowerset(data)
			if err != nil {
				t.Errorf("Error creating Powerset: %v", err)
			}

			counts := indexCounter(p.AllBorrowed())

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
	for range 100 {
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
			p, err := NewPermutations(data, int(k))
			if err != nil {
				t.Errorf("Error creating Permutations: %v", err)
			}

			counts := indexCounter(p.All())

			for num, count := range counts {
				count_big := big.NewInt(int64(count))
				if count_big.Cmp(times_we_see_each_index) != 0 {
					t.Errorf("Expected %v to appear %v times, but it appeared %v times", num, times_we_see_each_index, count)
				}
			}
		})
	}
}

func Test100RandomPermutationsBorrowed(t *testing.T) {
	for range 100 {
		n := rand.Int63n(15) + 1
		k := rand.Int63n(n) + 1
		data := stepped_range(0, int(n), 1)

		times_we_see_each_index := elts_in_permutations(int(n), int(k))
		if times_we_see_each_index.Cmp(big.NewInt(10000000)) > 0 {
			t.Logf("Skipping test because we see each index more than 10_000_000 times")
			continue
		}

		run_name := fmt.Sprintf("n=%v, k=%v", n, k)
		t.Run(run_name, func(t *testing.T) {
			p, err := NewPermutations(data, int(k))
			if err != nil {
				t.Errorf("Error creating Permutations: %v", err)
			}

			counts := indexCounter(p.AllBorrowed())

			for num, count := range counts {
				count_big := big.NewInt(int64(count))
				if count_big.Cmp(times_we_see_each_index) != 0 {
					t.Errorf("Expected %v to appear %v times, but it appeared %v times", num, times_we_see_each_index, count)
				}
			}
		})
	}
}
