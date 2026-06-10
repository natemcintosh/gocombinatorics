[![Go Reference](https://pkg.go.dev/badge/github.com/natemcintosh/gocombinatorics.svg)](https://pkg.go.dev/github.com/natemcintosh/gocombinatorics)

# gocombinatorics
**Author: Nathan McIntosh**

## About
Lazy combinatorics using Go 1.24+ iterators. Each iterator's `All()` method returns an `iter.Seq2[[]int, []T]`, yielding freshly allocated index and item slices on each iteration.

Uses Go 1.18 generics. No external dependencies beyond the standard library.

## On Offer:
- Lazy Combinations: create a `Combinations` struct with `NewCombinations()` function
- Lazy Combinations with replacement: create a `CombinationsWithReplacement` struct with `NewCombinationsWithReplacement()` function
- Lazy Permutations: create a `Permutations` struct with `NewPermutations()` function
- Lazy Product (k-fold Cartesian product, `n^k` tuples): create a `Product` struct with `NewProduct()` function
- Lazy ProductOf (Cartesian product of multiple distinct slices): create a `ProductOf` struct with `NewProductOf()` function
- Lazy Powerset (all `2^n` subsets, including the empty set): create a `Powerset` struct with `NewPowerset()` function
- Lazy Integer Partitions (the ways to write `n` as a sum of positive integers): create an `IntegerPartitions` struct with `NewIntegerPartitions()` function
- Lazy Set Partitions (the ways to divide a slice into non-empty, disjoint blocks): create a `SetPartitions` struct with `NewSetPartitions()` (all partitions, Bell numbers) or `NewSetPartitionsK()` (exactly `k` blocks, Stirling numbers of the second kind)

The two partition types yield differently shaped values (see [Partitions](#partitions) below); everything in the next three paragraphs applies to the other six types. Each of those provides three iteration methods:
- **`All() iter.Seq2[[]int, []T]`** — yields freshly allocated index and item slices each iteration. Safe to retain across iterations.
- **`AllBorrowed() iter.Seq2[[]int, []T]`** — yields shared internal buffers, overwritten each iteration. Much faster (see benchmarks below), but callers must not retain or modify the yielded slices.
- **`IndicesBorrowed() iter.Seq[[]int]`** — yields only the index slice and never gathers items. The fastest path when you need only the combinatorial structure — up to ~8× faster, and constant regardless of element size ([see below](#index-only-iteration-with-indicesborrowed)).

Each type also exposes a `Length *big.Int` field, precomputed by its constructor, holding the total number of items the iterator will yield (e.g. `n choose k` for `Combinations`, `2^n` for `Powerset`). It's a `*big.Int` because these counts can overflow `uint64` for even moderate inputs.

Each type also provides unranking and random-sampling methods:
- **`Nth(i *big.Int) ([]int, []T, error)`** — returns the indices and items that `All()` would yield on its *i*-th iteration (0-based), computed directly without iterating from the start. Returns freshly allocated slices, and an error if `i < 0` or `i >= Length`. Useful for sharded/parallel iteration, random sampling, or resuming from a checkpoint.
- **`Random(r *rand.Rand) ([]int, []T)`** — returns a uniform-random element of the iteration space without enumerating it, by drawing a random rank in `[0, Length)` and unranking it. Pass a seeded `*rand.Rand` (from `math/rand`) for deterministic results.

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

`ProductOf` gives the Cartesian product of multiple distinct slices of possibly
different lengths — equivalent to `itertools.product(A, B, C)`. All axes must share one
element type, which makes it a natural fit for parameter sweeps / grid searches. Index
`j` of each yielded indices slice refers to an element of axis `j`:
```go
p, err := combo.NewProductOf([]int{0, 1}, []int{10, 20, 30})
if err != nil {
	log.Fatal(err)
}
for indices, items := range p.All() {
	fmt.Println(indices, items) // [0 0] [0 10], [0 1] [0 20], ... [1 2] [1 30] — 6 in all
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
## Partitions

`IntegerPartitions` yields each partition of `n` as a single `[]int` of parts in
non-increasing order, in reverse lexicographic order. Its `Length` is the partition
function `p(n)`:
```go
ip, err := combo.NewIntegerPartitions(5)
if err != nil {
	log.Fatal(err)
}
for parts := range ip.All() {
	fmt.Println(parts) // [5], [4 1], [3 2], [3 1 1], [2 2 1], [2 1 1 1], [1 1 1 1 1]
}
```

`SetPartitions` yields each partition as an assignment slice (a restricted growth
string: `assignment[i]` is the block index of element `i`) plus the blocks themselves
as a `[][]T`. Blocks are ordered by their smallest element, and elements keep their
input order within a block. `Length` is the Bell number `B(n)`, or the Stirling number
of the second kind `S(n, k)` when constructed with `NewSetPartitionsK`:
```go
sp, err := combo.NewSetPartitions([]string{"a", "b", "c"})
if err != nil {
	log.Fatal(err)
}
for assignment, blocks := range sp.All() {
	fmt.Println(assignment, blocks) // [0 0 0] [[a b c]], [0 0 1] [[a b] [c]], ... [0 1 2] [[a] [b] [c]] — 5 in all
}

sp2, _ := combo.NewSetPartitionsK([]int{1, 2, 3, 4}, 2) // the 7 ways to split into exactly 2 blocks
```

Both types provide `All()` and `AllBorrowed()` with the usual fresh-vs-reused buffer
contract. `SetPartitions` also provides `AssignmentsBorrowed() iter.Seq[[]int]`, the
index-only fast path analogous to `IndicesBorrowed()`. Neither type currently
implements `Nth`/`Random`.

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
| Combinations | (10, 3) | 120 | 240 | 2 | ~7.9x |
| Combinations | (200, 3) | 1,313,400 | 2,626,806 | 2 | ~8.7x |
| Combinations | (26, 12) | 9,657,700 | 19,315,456 | 2 | ~7.2x |
| CombinationsWR | (10, 3) | 220 | 440 | 2 | ~7.3x |
| CombinationsWR | (15, 5) | 11,628 | 23,258 | 2 | ~7.2x |
| Permutations | (10, 3) | 720 | 1,441 | 2 | ~7.9x |
| Permutations | (10, 8) | 1,814,400 | 3,628,812 | 3 | ~8.2x |
| Product | (10, 3) | 1,000 | 2,000 | 2 | ~10.1x |
| Product | (6, 6) | 46,656 | 93,314 | 2 | ~8.8x |
| Powerset | (10, —) | 1,024 | 2,462 | 420 | ~3.8x |
| Powerset | (16, —) | 65,536 | 132,129 | 1,056 | ~6.0x |

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

Microbenchmark, `C(22, 11) = 705,432` combinations, full iteration each op (Intel i7-14700F, Go 1.26.3, linux/amd64). `Items` = `AllBorrowed()`; `Indices` = `IndicesBorrowed()`:

| Payload `T` | `Items` ns/op | `Indices` ns/op | Speedup |
|-------------|--------------:|----------------:|--------:|
| `int`       |     5,215,511 |       2,128,637 | ~2.4×   |
| `[192]byte` |    16,788,074 |       2,137,282 | ~7.9×   |

Note how `Indices` stays ~constant (~2 ms) while `Items` grows ~3× from `int` to `[192]byte`. This pairs naturally with an **SoA** design: run combinatorics over a slice of `int` ids and index your own columns caller-side, rather than building an array-of-structs of fat `T`. Reproduce with `go test -bench BenchmarkCombinationsIndicesVsItems -benchmem`.

---
## How is this library tested?
The suite has **62 test functions**, which expand to roughly **1,200 executed cases** on
each `go test` run — the property tests alone draw 100 random inputs per iterator type.
The exact total varies run to run, because those inputs are random and very large cases
are skipped by a 10,000,000-occurrence guard. The testing happens at a few layers:

- **Golden fixture tests.** Exact output is compared against CSV fixtures in `testdata/`,
  including a combination of length 1,313,400 (200 choose 3), a combination with
  replacement of length 11,628 (C(19, 5)), and a permutation of length 970,200
  (100 · 99 · 98).
- **Per-type unit tests.** Table-driven tests cover constructor validation, small
  hand-verified output sequences (including string payloads), and the `nchoosek` /
  `factorial` math helpers.
- **Iteration-path agreement.** All three iteration paths — `All()`, `AllBorrowed()`, and
  `IndicesBorrowed()` — are cross-checked against each other so the low-allocation and
  index-only paths can't silently diverge from `All()`. `indices_test.go` verifies that
  `IndicesBorrowed()` matches `AllBorrowed()` for all six iterator types.
- **Property tests.** `property_test.go` runs 100 random inputs through both `All()` and
  `AllBorrowed()` for each of the six iterator types — combinations, combinations with
  replacement, permutations, product, product-of, and powerset — checking that every index appears
  exactly the number of times the math predicts. The partition types get their own
  exhaustive sweeps: every integer partition for `n` up to 12 and every set partition
  (for every valid `k`) for `n` up to 8 is checked for validity, uniqueness, and
  agreement with `Length`.
