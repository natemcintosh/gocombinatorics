// Command go-driver is a hyperfine workload for the cross-language benchmark in
// issue #8. It iterates the full sequence of one combinatorics operation,
// sum-accumulating (wrapping uint64) a data-dependent value per item, and
// prints the accumulator once at the end. The data-dependent sum defeats
// dead-code elimination and gives a cheap, observable per-iteration cost that
// matches the Python and Rust drivers.
//
// Usage: go-driver [--borrowed] <op> <args...>
//
//	combinations | cwr | permutations | product  <n> <k>
//	powerset                                     <n>
//	productof                                    <size>...   (one axis per size)
//	intpartitions                                <n>
//	setpartitions                                <n> [k]
//
// For the index-yielding ops the accumulator is the sum of every index. For
// intpartitions it is the sum of every part. For setpartitions it is the sum
// of assignment[i] * i over the restricted-growth-string assignment, which is
// sensitive to block numbering (plain element sums are constant per
// partition).
//
// With --borrowed the library's buffer-reuse AllBorrowed() iterators are used
// (AssignmentsBorrowed() for setpartitions); otherwise the allocating All()
// iterators are used.
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
		log.Fatal("usage: go-driver [--borrowed] <op> <args...>")
	}

	op := args[0]
	nums := make([]int, len(args)-1)
	for i, a := range args[1:] {
		v, err := strconv.Atoi(a)
		if err != nil {
			log.Fatalf("invalid argument %q: %v", a, err)
		}
		nums[i] = v
	}

	acc, err := accumulate(op, nums, borrowed)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(acc)
}

// accumulate runs the full sequence for op and returns the wrapping-uint64
// accumulator. nums holds the op's numeric arguments (see the usage comment).
func accumulate(op string, nums []int, borrowed bool) (uint64, error) {
	switch op {
	case "combinations", "cwr", "permutations", "product":
		if len(nums) != 2 {
			return 0, fmt.Errorf("op %q requires <n> <k>", op)
		}
		seq, err := indexIterator(op, nums[0], nums[1], borrowed)
		if err != nil {
			return 0, err
		}
		return sumIndices(seq), nil

	case "powerset":
		if len(nums) != 1 {
			return 0, fmt.Errorf("op %q requires <n>", op)
		}
		p, err := combo.NewPowerset(intRange(nums[0]))
		if err != nil {
			return 0, err
		}
		return sumIndices(indices(p.All(), p.AllBorrowed(), borrowed)), nil

	case "productof":
		if len(nums) < 1 {
			return 0, fmt.Errorf("op %q requires <size>...", op)
		}
		axes := make([][]int, len(nums))
		for i, size := range nums {
			axes[i] = intRange(size)
		}
		p, err := combo.NewProductOf(axes...)
		if err != nil {
			return 0, err
		}
		return sumIndices(indices(p.All(), p.AllBorrowed(), borrowed)), nil

	case "intpartitions":
		if len(nums) != 1 {
			return 0, fmt.Errorf("op %q requires <n>", op)
		}
		p, err := combo.NewIntegerPartitions(nums[0])
		if err != nil {
			return 0, err
		}
		seq := p.All()
		if borrowed {
			seq = p.AllBorrowed()
		}
		return sumIndices(seq), nil

	case "setpartitions":
		if len(nums) != 1 && len(nums) != 2 {
			return 0, fmt.Errorf("op %q requires <n> [k]", op)
		}
		var s *combo.SetPartitions[int]
		var err error
		if len(nums) == 2 {
			s, err = combo.NewSetPartitionsK(intRange(nums[0]), nums[1])
		} else {
			s, err = combo.NewSetPartitions(intRange(nums[0]))
		}
		if err != nil {
			return 0, err
		}
		var acc uint64
		sum := func(assignment []int) {
			for i, b := range assignment {
				acc += uint64(b) * uint64(i)
			}
		}
		if borrowed {
			for assignment := range s.AssignmentsBorrowed() {
				sum(assignment)
			}
		} else {
			for assignment := range s.All() {
				sum(assignment)
			}
		}
		return acc, nil

	default:
		return 0, fmt.Errorf("unknown op %q", op)
	}
}

// indexIterator builds the index sequence for the four fixed-arity (n, k) ops.
func indexIterator(op string, n, k int, borrowed bool) (iter.Seq[[]int], error) {
	data := intRange(n)
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
	default:
		return nil, fmt.Errorf("unknown op %q", op)
	}
}

// sumIndices folds every value of every yielded slice into a wrapping uint64.
func sumIndices(seq iter.Seq[[]int]) uint64 {
	var acc uint64
	for vals := range seq {
		for _, v := range vals {
			acc += uint64(v)
		}
	}
	return acc
}

// intRange returns the ints 0..n-1.
func intRange(n int) []int {
	data := make([]int, n)
	for i := range data {
		data[i] = i
	}
	return data
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
