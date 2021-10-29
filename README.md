# gocombinatorics
## Author: Nathan McIntosh

---
### About
Basic lazy combinatorics. It gives you the indices for the next combination/permutation
when you call `Next()`.

Will this repo become redundant once generics are added in 1.18? Probably so. It'll be
much easier to have a function with a signature like:
```go
Combination[T any]([]T, int) [][]T
```

However, this is still good learning practice.

**If you are looking for a more production ready combinatorics library** I would suggest
using gonum's [combin](https://pkg.go.dev/gonum.org/v1/gonum@v0.9.3/stat/combin) library. Note however that gonum doesn't provide combinations with replacement functionality.

---
### On Offer:
- [X] Lazy Combinations: create a `Combinations` struct with `NewCombinations()` function
- [X] Lazy Combinations with replacement: create a `CombinationsWithReplacement` struct with `NewCombinationsWithReplacement()` function
- [X] Lazy Permutations: create a `Permutations` struct with `NewPermutations()` function

Each of the above structs meets the interface
```go
type CombinationLike interface {
	Next() bool
	LenInds() int
	Indices() []int
}
```
- `Next()` is what you use to iterate forward
- `LenInds()` tells you how long the indices slice is (you could also get this from `len(c.Indices()))`
- `Indices()` gives you the slice containing the indices of the items for this iteration

---
### How to use:
Say you have a slice of strings: `["apple, "banana", "cherry"]` and you want to get all the combinations of 2 strings:
1. `["apple", "banana"]`
1. `["apple", "cherry"]`
1. `["banana", "cherry"]`
```go
package main

import (
	"fmt"
	"log"

	combo "github.com/natemcintosh/gocombinatorics"
)

func main() {
	my_strings := []string{"apple", "banana", "cherry"}
	c, err := combo.NewCombinations(len(my_strings), 2)
	if err != nil {
		log.Fatal(err)
	}

	for c.Next() {
		// Now c.Indices() has the indices of the next combination
		fmt.Println(c.Indices())

		// Write yourself a helper function like `getStrElts` to get the elements from your slice
		this_combination := getStrElts(my_strings, c.Indices())

		// Do something with this combination
		fmt.Println(this_combination)
	}
}

// getStrElts will get the elements of an string slice at the specified indices
func getStrElts(s []string, elts []int) []string {
	result := make([]string, len(elts))
	for i, e := range elts {
		result[i] = s[e]
	}
	return result
}

```

Here's another example getting combinations with replacement for a slice of People structs. For variety's sake, this example will have a reusable buffer to reduce the number of allocations.
```go
package main

import (
	"fmt"
	"log"

	combo "github.com/natemcintosh/gocombinatorics"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	// The stooges
	people := []Person{
		{"Larry", 20},
		{"Curly", 30},
		{"Moe", 40},
		{"Shemp", 50},
	}

	// We want to see all possible combinations of length 4, with replacement
	combos, err := combo.NewCombinationsWithReplacement(len(people), 4)
	if err != nil {
		log.Fatal(err)
	}

	// fill_person_buffer will fill up a buffer with the items from `data` at `indices`
	fill_person_buffer := func(buffer []Person, data []Person, indices []int) {
		if len(indices) != len(buffer) {
			log.Fatal("Buffer needs to be same size as indices")
		}

		for buff_idx, data_idx := range indices {
			buffer[buff_idx] = data[data_idx]
		}

	}

	// Create the buffer. Note that using a buffer may be faster, but will always overwrite
	// the last iteration
	buff := make([]Person, combos.LenInds())

	// Now iterate over the combinations with replacement
	for combos.Next() {
		fill_person_buffer(buff, people, combos.Indices())
		fmt.Println(buff)
	}
}
```

---
### How is this library tested?
There are a few basic test, including one testing a combination of length 1,313,400, one
testing a combination with replacement of length 11,628, one testing a permutation of
length 970,200.

The file `property_test.go` also performs some basic property testing (do we see the
number of elements we expect to) on 100 random inputs to combinations/combinations with
replacement/permutations every time `go test` is run.

---
### Thoughts for Future Improvements:
Once generics are added in go1.18, this package could provide a generic "fill up the buffer"
function as seen in the `fill_person_buffer` used above. It would probably look like
```go
// FillBuffer will fill up a buffer with the items from `data` at `indices`
func FillBuffer[T any](buffer []T, data []T, indices []int) {
	if len(indices) != len(buffer) {
		log.Fatal("Buffer needs to be same size as indices")
	}

	for buff_idx, data_idx := range indices {
		buffer[buff_idx] = data[data_idx]
	}

}
```

Or perhaps we could add the generic buffer as a field to each `CombinationLike` struct.
Then the user calls `c.Items()` to get the items for this iteration. This has the benefit
of being simpler for the user: they don't have to create the buffer. Just have to make 
sure that the user is clear on whether we are re-using the underlying array. See
https://pkg.go.dev/encoding/csv#Reader on `ReuseRecord`. Should I mimic that?
