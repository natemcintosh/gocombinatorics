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
- [X] Lazy Product (k-fold Cartesian product, `n^k` tuples): create a `Product` struct with `NewProduct()` function
- [X] Lazy Powerset (all `2^n` subsets, including the empty set): create a `Powerset` struct with `NewPowerset()` function

Each type provides three iteration methods:
- **`All() iter.Seq2[[]int, []T]`** — yields freshly allocated index and item slices each iteration. Safe to retain across iterations.
- **`AllBorrowed() iter.Seq2[[]int, []T]`** — yields shared internal buffers, overwritten each iteration. Much faster (see benchmarks below), but callers must not retain or modify the yielded slices.
- **`IndicesBorrowed() iter.Seq[[]int]`** — yields only the index slice and never gathers items. The fastest path when you need only the combinatorial structure — up to ~10× faster, and constant regardless of element size ([see below](#index-only-iteration-with-indicesborrowed)).

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

`Product` gives the k-fold Cartesian product of a slice with itself — equivalent to
`itertools.product(data, repeat=k)`. Unlike `Combinations`/`Permutations`, `k` may exceed
`len(data)`, since each of the `k` positions independently ranges over every element:
```go
p, err := combo.NewProduct([]int{0, 1, 2}, 2)
if err != nil {
	log.Fatal(err)
}
for indices, items := range p.All() {
	fmt.Println(indices, items) // [0 0] [0 0], [0 1] [0 1], ... [2 2] [2 2] — 9 in all
}
```

`Powerset` yields all `2^n` subsets, starting with the empty set, matching Python's
`powerset` recipe. It takes no `k`:
```go
ps, err := combo.NewPowerset([]string{"a", "b", "c"})
if err != nil {
	log.Fatal(err)
}
for indices, items := range ps.All() {
	fmt.Println(indices, items) // [] [], [0] [a], [1] [b], ... [0 1 2] [a b c] — 8 in all
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

Measured on an Intel i7-14700F. `All()` allocates 2 fresh slices per iteration. For the fixed-length iterators, `AllBorrowed()` allocates a constant handful of slices total. `Powerset` is the exception: it constructs a `Combinations` iterator per subset size, so its `AllBorrowed()` allocations grow linearly with `n` — still far below `All()`, but not constant.

| Type | (n, k) | Iterations | `All()` allocs | `AllBorrowed()` allocs | Speedup |
|------|--------|-----------|---------------|----------------------|---------|
| Combinations | (10, 3) | 120 | 244 | 5 | ~6.5x |
| Combinations | (200, 3) | 1,313,400 | 2,626,804 | 5 | ~7.7x |
| Combinations | (26, 12) | 9,657,700 | 19,315,475 | 5 | ~5.5x |
| CombinationsWR | (10, 3) | 220 | 444 | 5 | ~6.4x |
| CombinationsWR | (15, 5) | 11,628 | 23,260 | 5 | ~7.4x |
| Permutations | (10, 3) | 720 | 1,445 | 6 | ~8.2x |
| Permutations | (10, 8) | 1,814,400 | 3,628,814 | 6 | ~6.7x |
| Product | (10, 3) | 1,000 | 2,005 | 5 | ~9.5x |
| Product | (6, 6) | 46,656 | 93,317 | 5 | ~9.4x |
| Powerset | (10, —) | 1,024 | 2,478 | 423 | ~3.9x |
| Powerset | (16, —) | 65,536 | 132,144 | 1,059 | ~7.3x |

For a whole-process, cross-language wall-clock comparison against Python's stdlib `itertools` and the Rust `itertools` crate (run with [hyperfine](https://github.com/sharkdp/hyperfine) via `just bench-compare`), see [`bench/RESULTS.md`](bench/RESULTS.md).

---
## Index-Only Iteration with `IndicesBorrowed()`

Often you don't need the items at all — only the combinatorial *structure* (the `[]int` indices): for counting, scoring an index-space cost function, driving a struct-of-arrays (SoA) layout, or indexing your own columns. In those cases, materializing `[]T` items every iteration is pure wasted work.

`IndicesBorrowed()` yields the exact same index sequence as `AllBorrowed()`, but never gathers items:

```go
c, _ := combo.NewCombinations(myData, 3)

// Items path: gathers myData[i] into a []T every iteration.
for indices, items := range c.AllBorrowed() {
    _ = items // ... use items
}

// Index-only fast path: no item gather at all.
for indices := range c.IndicesBorrowed() {
    // indices is borrowed — do NOT retain it across iterations (clone if needed).
    score(indices) // e.g. index your own columns: colA[indices[0]], colB[indices[1]], ...
}
```

Like `AllBorrowed()`, the yielded slice is reused in place — clone it if you need to keep it.

### Why it's faster

The item gather is a random-access scatter into your data slice that the index-only loop skips entirely. Two effects compound:

- **Even for the smallest payload (`int`), index-only is faster** — the gather costs even when `T` is a single word.
- **Index-only is payload-agnostic**: its cost is constant regardless of `sizeof(T)`, while item materialization scales with the payload. The bigger your `T`, the more you save.

Microbenchmark, `C(22, 11) = 705,432` combinations, full iteration each op (Intel i7-14700F, Go 1.24, linux/amd64). `Items` = `AllBorrowed()`; `Indices` = `IndicesBorrowed()`:

| Payload `T` | `Items` ns/op | `Indices` ns/op | Speedup |
|-------------|--------------:|----------------:|--------:|
| `int`       |     5,209,111 |       2,156,814 | ~2.4×   |
| `[192]byte` |    20,947,620 |       1,997,375 | ~10.5×  |

Note how `Indices` stays ~constant (~2 ms) while `Items` grows ~4× from `int` to `[192]byte`. This pairs naturally with an **SoA** design: run combinatorics over a slice of `int` ids and index your own columns caller-side, rather than building an array-of-structs of fat `T`. Reproduce with `go test -bench BenchmarkCombinationsIndicesVsItems -benchmem`.

---
## How is this library tested?
There are a few basic tests, including one testing a combination of length 1,313,400, one
testing a combination with replacement of length 11,628, one testing a permutation of
length 970,200.

The file `property_test.go` also performs some basic property testing (do we see the
number of elements we expect to) on 100 random inputs to combinations/combinations with
replacement/permutations every time `go test` is run.
