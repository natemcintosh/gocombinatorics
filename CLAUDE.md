# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go test ./...          # run all tests (includes randomized property tests)
go test -run TestName  # run a single test by name
go build ./...         # build
gh issue               # for interacting with GitHub issues
```

There is no linter configured. No external dependencies beyond the standard library.

## Conventions

Use [Conventional Commits](https://www.conventionalcommits.org/) for all commit messages (e.g. `feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`).

## Architecture

Single-package Go library (`package gocombinatorics`) providing lazy iterators for combinatorics, modeled after Python's `itertools`. Uses generics and `iter.Seq2` range-over-func iterators. The library itself only needs Go 1.23+, but `go.mod` declares 1.24 because the benchmarks use `testing.B.Loop()` (added in 1.24).

**Five iterator types**, all following the same pattern:
- `Combinations[T]` — n choose k
- `CombinationsWithReplacement[T]` — combinations allowing repeated elements
- `Permutations[T]` — k-length permutations of n elements
- `Product[T]` — k-fold Cartesian product of the slice with itself (`n^k` tuples)
- `Powerset[T]` — all `2^n` subsets, including the empty set

Each is constructed with a `NewX(input_data []T, k int) (*X[T], error)` constructor that validates `k`/`n`, copies the input, and precomputes the total count into the exported `Length *big.Int` field. Two iterators deviate slightly: `Product` does **not** reject `k > n` (each position independently ranges over all n elements), and `Powerset` takes no `k` (`NewPowerset(input_data []T)`) and yields **variable-length** slices — its `AllBorrowed()` re-slices a buffer that grows with the subset size, and it reuses the `Combinations` iterator internally for each subset size `r = 1..n`.

**Iteration**: each type exposes two `iter.Seq2[[]int, []T]` methods (range over them with Go's range-over-func):
- `All()` — yields a freshly allocated indices slice and items slice each iteration; safe to retain.
- `AllBorrowed()` — yields the same internal index/item buffers, overwritten in place each iteration; the caller must not retain or mutate them across iterations. Use for low-allocation hot loops.

`All()` is implemented by wrapping `AllBorrowed()` and `slices.Clone`-ing each pair, so the iteration algorithm lives in `AllBorrowed()`. `Permutations.AllBorrowed()` keeps a local `cycles` slice for state.

**Helper functions** live alongside their primary consumer: `nchoosek`/`factorial` in `combinations.go`, `num_combinations_w_replacement`/`elts_in_combo_w_replacement` in `combinations_w_replacements.go`, `n_permutations`/`elts_in_permutations`/`stepped_range` in `permutations.go`, `num_products`/`elts_in_product` in `product.go`, `num_powersets`/`elts_in_powerset` in `powerset.go`. `fillBuf` (the shared index-to-item mapper) is in `general.go`.

**Testing**: table-driven unit tests per type plus `property_test.go` which runs 100 randomized inputs per type, verifying that each index appears the mathematically expected number of times. CSV fixtures in `testdata/`.
