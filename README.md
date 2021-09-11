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
c, err := NewCombinations(len(my_strings), 2)
if err != nil {
    log.Fatal(err)
}

for c.Next() {
    // Now c.Inds has the indices of the next combination
    fmt.Println(c.Inds)

    // Write yourself a helper function like `getStrElts` to get the elements from your slice
    this_combination := getStrElts(my_strings, c.Inds)

    // Do something with this combination
    fmt.Println(this_combination)
}

// getStrElts will get the elements of an string slice at the specified indices
func getStrElts(s []string], elts []int]) []string] {
	result := make([]string], len(elts))
	for i, e := range elts {
		result[i] = s[e]
	}
	return result
}
```