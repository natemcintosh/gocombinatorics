[![Go Reference](https://pkg.go.dev/badge/github.com/natemcintosh/gocombinatorics.svg)](https://pkg.go.dev/github.com/natemcintosh/gocombinatorics)

# gocombinatorics
## Author: Nathan McIntosh

---
### About
Basic lazy combinatorics. It gives you the next combination/permutation when you call 
`Next()`. Access the items with `Items()`.

This library has been updated to use generics. If you require a version of go 
<1.18, please use version 0.2.0 of this library.

---
### On Offer:
- [X] Lazy Combinations: create a `Combinations` struct with `NewCombinations()` function
- [X] Lazy Combinations with replacement: create a `CombinationsWithReplacement` struct with `NewCombinationsWithReplacement()` function
- [X] Lazy Permutations: create a `Permutations` struct with `NewPermutations()` function

Each of the above structs meets the interface (but the interface is not actually used anywhere)
```go
type CombinationLike[T any] interface {
	Next() bool
	LenInds() int
	Indices() []int
	Items() []T
}
```
- `Next()` is what you use to iterate forward
- `Items()` will return a slice of the items in this combination/permutation. Note that this buffer is re-used every iteration. If you require the results of every iteration, make a copy of the slice returned by `Items()` every iteration.
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
	c, err := combo.NewCombinations(my_strings, 2)
	if err != nil {
		log.Fatal(err)
	}

	for c.Next() {
		// Do something with this combination
		fmt.Println(c.Items())
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
	
	// Now iterate over the combinations with replacement
	for combos.Next() {
		fmt.Println(combos.Items())
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

