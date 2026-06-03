// Command go-driver is a hyperfine workload for the cross-language benchmark in
// issue #8. It iterates the full sequence of one combinatorics operation,
// sum-accumulating (wrapping uint64) every index it sees, and prints the
// accumulator once at the end. The data-dependent sum defeats dead-code
// elimination and gives a cheap, observable per-iteration cost that matches the
// Python and Rust drivers.
//
// Usage: go-driver [--borrowed] <op> <n> [k]
//
//	op  ∈ combinations | cwr | permutations | product | powerset
//	n   size of the input slice (elements are the ints 0..n-1)
//	k   required for every op except powerset
//
// With --borrowed the library's buffer-reuse AllBorrowed() iterator is used;
// otherwise the allocating All() iterator is used.
package main

import (
	"fmt"
	"iter"
	"log"
	"os"
	"strconv"

	combo "github.com/natemcintosh/gocombinatorics"
)

func main() {
	args := os.Args[1:]
	borrowed := false
	if len(args) > 0 && args[0] == "--borrowed" {
		borrowed = true
		args = args[1:]
	}

	if len(args) < 2 {
		log.Fatal("usage: go-driver [--borrowed] <op> <n> [k]")
	}

	op := args[0]
	n, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatalf("invalid n: %v", err)
	}

	k := 0
	if op != "powerset" {
		if len(args) < 3 {
			log.Fatalf("op %q requires k", op)
		}
		k, err = strconv.Atoi(args[2])
		if err != nil {
			log.Fatalf("invalid k: %v", err)
		}
	}

	data := make([]int, n)
	for i := range data {
		data[i] = i
	}

	seq, err := iterator(op, data, k, borrowed)
	if err != nil {
		log.Fatal(err)
	}

	var acc uint64
	for inds := range seq {
		for _, idx := range inds {
			acc += uint64(idx)
		}
	}
	fmt.Println(acc)
}

// iterator returns an index-only sequence for the requested op, selecting
// All() or AllBorrowed() based on borrowed.
func iterator(op string, data []int, k int, borrowed bool) (iter.Seq[[]int], error) {
	switch op {
	case "combinations":
		c, err := combo.NewCombinations(data, k)
		if err != nil {
			return nil, err
		}
		return indices(c.All(), c.AllBorrowed(), borrowed), nil
	case "cwr":
		c, err := combo.NewCombinationsWithReplacement(data, k)
		if err != nil {
			return nil, err
		}
		return indices(c.All(), c.AllBorrowed(), borrowed), nil
	case "permutations":
		p, err := combo.NewPermutations(data, k)
		if err != nil {
			return nil, err
		}
		return indices(p.All(), p.AllBorrowed(), borrowed), nil
	case "product":
		p, err := combo.NewProduct(data, k)
		if err != nil {
			return nil, err
		}
		return indices(p.All(), p.AllBorrowed(), borrowed), nil
	case "powerset":
		p, err := combo.NewPowerset(data)
		if err != nil {
			return nil, err
		}
		return indices(p.All(), p.AllBorrowed(), borrowed), nil
	default:
		return nil, fmt.Errorf("unknown op %q", op)
	}
}

// indices drops the items half of an iter.Seq2 and picks the allocating or
// borrowed variant, yielding just the index slices.
func indices(all, borrowedSeq iter.Seq2[[]int, []int], borrowed bool) iter.Seq[[]int] {
	src := all
	if borrowed {
		src = borrowedSeq
	}
	return func(yield func([]int) bool) {
		for inds := range src {
			if !yield(inds) {
				return
			}
		}
	}
}
