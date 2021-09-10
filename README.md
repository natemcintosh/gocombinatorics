# gocombinatorics
Basic combinatorics. It gives you the indices for the next combination/permutation/etc.

Will this repo become redundant once generics are added in 1.18? Probably so. It'll be
much easier to have a function with a signature like:
```go
Combination[T any]([]T, int) [][]T
```

However, this is still good practice, and github copilot helps me write a lot of the
boring tests and benchmark code.