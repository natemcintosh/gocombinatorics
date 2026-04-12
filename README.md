[![Go Reference](https://pkg.go.dev/badge/github.com/natemcintosh/gocombinatorics.svg)](https://pkg.go.dev/github.com/natemcintosh/gocombinatorics)

# gocombinatorics
**Author: Nathan McIntosh**

## About
Lazy combinatorics using Go 1.23+ iterators. Each iterator's `All()` method returns an `iter.Seq2[[]int, []T]`, yielding freshly allocated index and item slices on each iteration.

Uses Go 1.18 generics. No external dependencies beyond the standard library.

## On Offer:
- [X] Lazy Combinations: create a `Combinations` struct with `NewCombinations()` function
- [X] Lazy Combinations with replacement: create a `CombinationsWithReplacement` struct with `NewCombinationsWithReplacement()` function
- [X] Lazy Permutations: create a `Permutations` struct with `NewPermutations()` function

Each type provides an `All()` method that returns an `iter.Seq2[[]int, []T]`. Use it with a `for range` loop to iterate over all combinations/permutations. Each iteration yields two values:
- The indices into the original data
- The corresponding items

---
## How to use:
Say you have a slice of strings: `["apple", "banana", "cherry"]` and you want to get all the combinations of 2 strings:
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
	c, err := combo.NewCombinations(my_strings, 2)
	if err != nil {
		log.Fatal(err)
	}

	for indices, items := range c.All() {
		fmt.Println(indices, items)
	}
}
```

Here's another example getting combinations with replacement for a slice of People structs.
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
	combos, err := combo.NewCombinationsWithReplacement(people, 4)
	if err != nil {
		log.Fatal(err)
	}

	for _, items := range combos.All() {
		fmt.Println(items)
	}
}
```

---
## How is this library tested?
There are a few basic tests, including one testing a combination of length 1,313,400, one
testing a combination with replacement of length 11,628, one testing a permutation of
length 970,200.

The file `property_test.go` also performs some basic property testing (do we see the
number of elements we expect to) on 100 random inputs to combinations/combinations with
replacement/permutations every time `go test` is run.
