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

Each type provides two iteration methods, both returning `iter.Seq2[[]int, []T]`:
- **`All()`** — yields freshly allocated index and item slices each iteration. Safe to retain across iterations.
- **`AllBorrowed()`** — yields shared internal buffers, overwritten each iteration. Much faster (see benchmarks below), but callers must not retain or modify the yielded slices.

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
## Low-Allocation Iteration with `AllBorrowed()`

`AllBorrowed()` reuses internal buffers instead of allocating fresh slices each iteration. Use it when you process each combination/permutation inline without storing it:

```go
c, _ := combo.NewCombinations(myData, 3)

// Fast path: process each combination without retaining it
for indices, items := range c.AllBorrowed() {
    // Use indices and items here, but do NOT store them —
    // they will be overwritten on the next iteration.
    fmt.Println(indices, items)
}
```

If you need to collect results, use `All()` instead (or copy the slices yourself).

### Benchmarks

Measured on an Intel i7-14700F. `All()` allocates 2 fresh slices per iteration; `AllBorrowed()` allocates a constant 2-3 slices total.

| Type | (n, k) | Iterations | `All()` allocs | `AllBorrowed()` allocs | Speedup |
|------|--------|-----------|---------------|----------------------|---------|
| Combinations | (10, 3) | 120 | 244 | 5 | ~6.5x |
| Combinations | (200, 3) | 1,313,400 | 2,626,804 | 5 | ~7.7x |
| Combinations | (26, 12) | 9,657,700 | 19,315,475 | 5 | ~5.5x |
| CombinationsWR | (10, 3) | 220 | 444 | 5 | ~6.4x |
| CombinationsWR | (15, 5) | 11,628 | 23,260 | 5 | ~7.4x |
| Permutations | (10, 3) | 720 | 1,445 | 6 | ~8.2x |
| Permutations | (10, 8) | 1,814,400 | 3,628,814 | 6 | ~6.7x |

---
## How is this library tested?
There are a few basic tests, including one testing a combination of length 1,313,400, one
testing a combination with replacement of length 11,628, one testing a permutation of
length 970,200.

The file `property_test.go` also performs some basic property testing (do we see the
number of elements we expect to) on 100 random inputs to combinations/combinations with
replacement/permutations every time `go test` is run.
