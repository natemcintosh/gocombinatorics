// property_test.go offers a very simple property test on the available combinatorics
// functionality.

// If we want all possible length 2 combinations of 10 items, then we can expect to get
// n choose k = n! / (k! * (n-k)!) = 10! / (2! * (10-2)!) = 45. In those 45 cases, we
// should expect to see each number exactly 45 - (9 choose 2) = 45 - 36 = 9 times.

// The above is an example with combinations. We can possibly also do them for other
// functionality as well
package gocombinatorics

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
