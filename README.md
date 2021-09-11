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

---
### On Offer:
- [X] Lazy Combinations
- [ ] Lazy Combinations with replacement
- [ ] Lazy Permutations
- [ ] Lazy Cartesian Product

---
### How to use:
Say you have a slice of strings: `["apple, "banana", "cherry"]` and you want to get all the combinations of 2 strings:
1. `["apple", "banana"]`
1. `["apple", "cherry"]`
1. `["banana", "cherry"]`
```go
my_strings := []string{"apple", "banana", "cherry"}
c := NewCombination(len(my_strings), 2)

var err error
for {
    err = c.Next()
    if err == ErrEndOfCombinations {
        break
    }
    // Now c.Inds has the indices of the next combination
    fmt.Println(c.Inds)

    // Now it's up to you to get the elements at those indices from the slice
}
```