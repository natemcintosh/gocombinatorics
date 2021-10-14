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

**If you are looking for a more production ready combinatorics library** I would suggest
using gonum's [combin](https://pkg.go.dev/gonum.org/v1/gonum@v0.9.3/stat/combin) library. Note however that gonum doesn't provide combinations with replacement functionality.

---
### On Offer:
- [X] Lazy Combinations
- [X] Lazy Combinations with replacement
- [X] Lazy Permutations

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
    // Now c.Indices() has the indices of the next combination
    fmt.Println(c.Indices())

    // Write yourself a helper function like `getStrElts` to get the elements from your slice
    this_combination := getStrElts(my_strings, c.Indices())

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

Here's another example getting combinations with replacement for a slice of People structs 
```go
type Person struct {
    Name string
    age  int
}

// The stooges
people := []Person{
    {"Larry", 20},
    {"Curly", 30},
    {"Moe", 40},
    {"Shemp", 50},
}

// We want to see all possible combinations of length 4, with replacement
combos, err := NewCombinationsWithReplacement(len(people), 4)
if err != nil {
    log.Fatal(err)
}

// Write yourself a helper function like `getPeopleElts` to get elements from the
// correct type of slice
getPeopleElts := func(people []Person, inds []int) []Person {
    result := make([]Person, len(inds))
    for i, ind := range inds {
        result[i] = people[ind]
    }
    return result
}

// Now iterate over the combinations with replacement
for combos.Next() {
    this_combo := getPeopleElts(people, combos.Indices())
    fmt.Println(this_combo)
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

---
### Recent Improvements:
Now no longer require the use of big Int math, which helps keep it fast, but still doesn't need to worry about integer overflow.
