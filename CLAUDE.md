# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go test ./...          # run all tests (includes randomized property tests)
go test -run TestName  # run a single test by name
go build ./...         # build
```

There is no linter configured. No external dependencies beyond the standard library.

## Conventions

Use [Conventional Commits](https://www.conventionalcommits.org/) for all commit messages (e.g. `feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`).

## Architecture

Single-package Go library (`package gocombinatorics`) providing lazy iterators for combinatorics, modeled after Python's `itertools`. Uses Go 1.18 generics.

**Three iterator types**, all following the same pattern:
- `Combinations[T]` — n choose k
- `CombinationsWithReplacement[T]` — combinations allowing repeated elements
- `Permutations[T]` — k-length permutations of n elements (uses a `cycles` slice for state)

**Shared interface** (`CombinationLike[T]` in `general.go`): `Next() bool`, `LenInds() int`, `Indices() []int`, `Items() []T`. Declared but not used as a constraint anywhere — the types satisfy it structurally.

**Iteration pattern**: call `Next()` to advance, then `Items()` or `Indices()` to read the current state. Both return shared internal slices (documented as such) — callers must copy if they need to retain values across iterations.

**Helper functions** live alongside their primary consumer: `nchoosek`/`factorial` in `combinations.go`, `num_combinations_w_replacement`/`elts_in_combo_w_replacement` in `combinations_w_replacements.go`, `n_permutations`/`elts_in_permutations`/`stepped_range` in `permutations.go`. `fill_buffer` (the shared index-to-item mapper) is in `general.go`.

**Testing**: table-driven unit tests per type plus `property_test.go` which runs 100 randomized inputs per type, verifying that each index appears the mathematically expected number of times. CSV fixtures in `testdata/`.
