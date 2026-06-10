package gocombinatorics_test

import (
	"fmt"
	"math/big"
	"math/rand"
	"slices"

	combo "github.com/natemcintosh/gocombinatorics"
)

// Example demonstrates the basic usage pattern shared by all iterators:
// construct with a New* function, then range over All().
func Example() {
	c, err := combo.NewCombinations([]string{"apple", "banana", "cherry"}, 2)
	if err != nil {
		panic(err)
	}
	for _, items := range c.All() {
		fmt.Println(items)
	}
	// Output:
	// [apple banana]
	// [apple cherry]
	// [banana cherry]
}

func ExampleNewCombinations() {
	c, err := combo.NewCombinations([]string{"a", "b", "c", "d"}, 2)
	if err != nil {
		panic(err)
	}
	fmt.Println("4 choose 2 =", c.Length)

	// k may not exceed the number of elements.
	_, err = combo.NewCombinations([]string{"a", "b"}, 3)
	fmt.Println(err)
	// Output:
	// 4 choose 2 = 6
	// k must be less than or equal to len(input_data)
}

func ExampleNewCombinationsWithReplacement() {
	c, err := combo.NewCombinationsWithReplacement([]string{"a", "b", "c"}, 2)
	if err != nil {
		panic(err)
	}
	// Unlike plain combinations, an element may be picked more than once.
	fmt.Println("Length:", c.Length)
	// Output:
	// Length: 6
}

func ExampleNewPermutations() {
	p, err := combo.NewPermutations([]string{"a", "b", "c"}, 3)
	if err != nil {
		panic(err)
	}
	fmt.Println("3! =", p.Length)
	// Output:
	// 3! = 6
}

func ExampleNewProduct() {
	// Unlike Combinations and Permutations, k may exceed the number of
	// elements: each position independently ranges over all of them.
	p, err := combo.NewProduct([]int{0, 1}, 3)
	if err != nil {
		panic(err)
	}
	fmt.Println("2^3 =", p.Length)
	// Output:
	// 2^3 = 8
}

func ExampleNewProductOf() {
	// ProductOf is variadic over distinct axes, which may have different
	// lengths. Each tuple takes one element from each axis.
	p, err := combo.NewProductOf([]string{"a", "b"}, []string{"x", "y", "z"})
	if err != nil {
		panic(err)
	}
	fmt.Println("2 * 3 =", p.Length)
	// Output:
	// 2 * 3 = 6
}

func ExampleNewPowerset() {
	// A powerset has 2^n subsets, including the empty set.
	p, err := combo.NewPowerset([]string{"a", "b", "c"})
	if err != nil {
		panic(err)
	}
	fmt.Println("2^3 =", p.Length)
	// Output:
	// 2^3 = 8
}

func ExampleNewIntegerPartitions() {
	// Length is the partition function p(n).
	p, err := combo.NewIntegerPartitions(5)
	if err != nil {
		panic(err)
	}
	fmt.Println("p(5) =", p.Length)
	// Output:
	// p(5) = 7
}

func ExampleNewSetPartitions() {
	// Length is the Bell number B(n).
	s, err := combo.NewSetPartitions([]string{"a", "b", "c", "d"})
	if err != nil {
		panic(err)
	}
	fmt.Println("B(4) =", s.Length)
	// Output:
	// B(4) = 15
}

func ExampleNewSetPartitionsK() {
	// Length is the Stirling number of the second kind S(n, k).
	s, err := combo.NewSetPartitionsK([]string{"a", "b", "c", "d"}, 2)
	if err != nil {
		panic(err)
	}
	fmt.Println("S(4, 2) =", s.Length)
	// Output:
	// S(4, 2) = 7
}

func ExampleCombinations_All() {
	c, err := combo.NewCombinations([]string{"a", "b", "c", "d"}, 2)
	if err != nil {
		panic(err)
	}
	for indices, items := range c.All() {
		fmt.Println(indices, items)
	}
	// Output:
	// [0 1] [a b]
	// [0 2] [a c]
	// [0 3] [a d]
	// [1 2] [b c]
	// [1 3] [b d]
	// [2 3] [c d]
}

func ExampleCombinations_AllBorrowed() {
	c, err := combo.NewCombinations([]string{"a", "b", "c", "d"}, 3)
	if err != nil {
		panic(err)
	}
	// AllBorrowed reuses its buffers, so the yielded slices must not be
	// retained across iterations: clone anything you want to keep.
	var kept [][]string
	for _, items := range c.AllBorrowed() {
		if items[0] == "b" {
			kept = append(kept, slices.Clone(items))
		}
	}
	fmt.Println(kept)
	// Output:
	// [[b c d]]
}

func ExampleCombinations_IndicesBorrowed() {
	c, err := combo.NewCombinations([]string{"a", "b", "c", "d"}, 2)
	if err != nil {
		panic(err)
	}
	// IndicesBorrowed skips gathering items: the fast path when only the
	// combinatorial structure is needed. The slice is reused each iteration.
	count := 0
	for indices := range c.IndicesBorrowed() {
		count += indices[1] - indices[0]
	}
	fmt.Println(count)
	// Output:
	// 10
}

func ExampleCombinations_Nth() {
	c, err := combo.NewCombinations([]string{"a", "b", "c", "d"}, 2)
	if err != nil {
		panic(err)
	}
	// Nth unranks: it returns what All() would yield on iteration i,
	// without iterating from the start.
	indices, items, err := c.Nth(big.NewInt(4))
	if err != nil {
		panic(err)
	}
	fmt.Println(indices, items)

	// Ranks at or beyond Length are an error.
	_, _, err = c.Nth(big.NewInt(6))
	fmt.Println(err)
	// Output:
	// [1 3] [b d]
	// i must be less than Length
}

func ExampleCombinations_Random() {
	c, err := combo.NewCombinations([]string{"a", "b", "c", "d"}, 2)
	if err != nil {
		panic(err)
	}
	// A fixed seed makes the draws deterministic.
	r := rand.New(rand.NewSource(42))
	_, items := c.Random(r)
	fmt.Println(items)
	// Output:
	// [b c]
}

func ExampleCombinationsWithReplacement_All() {
	c, err := combo.NewCombinationsWithReplacement([]string{"a", "b"}, 2)
	if err != nil {
		panic(err)
	}
	for indices, items := range c.All() {
		fmt.Println(indices, items)
	}
	// Output:
	// [0 0] [a a]
	// [0 1] [a b]
	// [1 1] [b b]
}

func ExamplePermutations_All() {
	p, err := combo.NewPermutations([]string{"a", "b", "c"}, 2)
	if err != nil {
		panic(err)
	}
	for indices, items := range p.All() {
		fmt.Println(indices, items)
	}
	// Output:
	// [0 1] [a b]
	// [0 2] [a c]
	// [1 0] [b a]
	// [1 2] [b c]
	// [2 0] [c a]
	// [2 1] [c b]
}

func ExampleProduct_All() {
	p, err := combo.NewProduct([]string{"a", "b"}, 2)
	if err != nil {
		panic(err)
	}
	for indices, items := range p.All() {
		fmt.Println(indices, items)
	}
	// Output:
	// [0 0] [a a]
	// [0 1] [a b]
	// [1 0] [b a]
	// [1 1] [b b]
}

func ExampleProductOf_All() {
	p, err := combo.NewProductOf([]string{"a", "b"}, []string{"x", "y", "z"})
	if err != nil {
		panic(err)
	}
	// Index j in the indices slice refers to an element of axis j.
	for indices, items := range p.All() {
		fmt.Println(indices, items)
	}
	// Output:
	// [0 0] [a x]
	// [0 1] [a y]
	// [0 2] [a z]
	// [1 0] [b x]
	// [1 1] [b y]
	// [1 2] [b z]
}

func ExamplePowerset_All() {
	p, err := combo.NewPowerset([]string{"a", "b"})
	if err != nil {
		panic(err)
	}
	// Subsets vary in length, starting with the empty set.
	for indices, items := range p.All() {
		fmt.Println(indices, items)
	}
	// Output:
	// [] []
	// [0] [a]
	// [1] [b]
	// [0 1] [a b]
}

func ExampleIntegerPartitions_All() {
	p, err := combo.NewIntegerPartitions(4)
	if err != nil {
		panic(err)
	}
	// Partitions are yielded in reverse lexicographic order, each as a
	// slice of non-increasing parts.
	for parts := range p.All() {
		fmt.Println(parts)
	}
	// Output:
	// [4]
	// [3 1]
	// [2 2]
	// [2 1 1]
	// [1 1 1 1]
}

func ExampleSetPartitions_All() {
	s, err := combo.NewSetPartitions([]string{"a", "b", "c"})
	if err != nil {
		panic(err)
	}
	// assignment[i] is the block index of element i; blocks are ordered by
	// their smallest element.
	for assignment, blocks := range s.All() {
		fmt.Println(assignment, blocks)
	}
	// Output:
	// [0 0 0] [[a b c]]
	// [0 0 1] [[a b] [c]]
	// [0 1 0] [[a c] [b]]
	// [0 1 1] [[a] [b c]]
	// [0 1 2] [[a] [b] [c]]
}

func ExampleSetPartitions_AssignmentsBorrowed() {
	s, err := combo.NewSetPartitionsK([]string{"a", "b", "c", "d"}, 2)
	if err != nil {
		panic(err)
	}
	// AssignmentsBorrowed skips gathering blocks: the fast path when only
	// the restricted growth strings are needed. The slice is reused each
	// iteration, so clone it to retain it.
	for assignment := range s.AssignmentsBorrowed() {
		fmt.Println(assignment)
	}
	// Output:
	// [0 0 0 1]
	// [0 0 1 0]
	// [0 0 1 1]
	// [0 1 0 0]
	// [0 1 0 1]
	// [0 1 1 0]
	// [0 1 1 1]
}
