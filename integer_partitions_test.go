package gocombinatorics

import (
	"fmt"
	"math/big"
	"slices"
	"testing"
)

func TestIntegerPartitionsErrors(t *testing.T) {
	for _, n := range []int{0, -1} {
		if _, err := NewIntegerPartitions(n); err == nil {
			t.Errorf("NewIntegerPartitions(%d) should have errored", n)
		}
	}
}

func TestIntegerPartitionsLength(t *testing.T) {
	// p(1)..p(10), OEIS A000041.
	want := []int64{1, 2, 3, 5, 7, 11, 15, 22, 30, 42}
	for i, w := range want {
		n := i + 1
		p, err := NewIntegerPartitions(n)
		if err != nil {
			t.Fatalf("NewIntegerPartitions(%d): %v", n, err)
		}
		if p.Length.Cmp(big.NewInt(w)) != 0 {
			t.Errorf("p(%d): got Length %v, want %d", n, p.Length, w)
		}
	}
}

func TestIntegerPartitionsAll(t *testing.T) {
	p, err := NewIntegerPartitions(5)
	if err != nil {
		t.Fatal(err)
	}
	want := [][]int{
		{5},
		{4, 1},
		{3, 2},
		{3, 1, 1},
		{2, 2, 1},
		{2, 1, 1, 1},
		{1, 1, 1, 1, 1},
	}
	got := slices.Collect(p.All())
	if len(got) != len(want) {
		t.Fatalf("got %d partitions, want %d", len(got), len(want))
	}
	for i := range want {
		if !slices.Equal(got[i], want[i]) {
			t.Errorf("partition %d: got %v, want %v", i, got[i], want[i])
		}
	}
}

func TestIntegerPartitionsProperties(t *testing.T) {
	for n := 1; n <= 12; n++ {
		p, err := NewIntegerPartitions(n)
		if err != nil {
			t.Fatal(err)
		}
		count := int64(0)
		seen := make(map[string]bool)
		for parts := range p.All() {
			count++
			// Every partition sums to n and is non-increasing.
			sum := 0
			for i, v := range parts {
				sum += v
				if v < 1 {
					t.Fatalf("n=%d: part %d is %d, want >= 1", n, i, v)
				}
				if i > 0 && parts[i-1] < v {
					t.Fatalf("n=%d: parts %v not non-increasing", n, parts)
				}
			}
			if sum != n {
				t.Fatalf("n=%d: parts %v sum to %d", n, parts, sum)
			}
			key := fmt.Sprint(parts)
			if seen[key] {
				t.Fatalf("n=%d: duplicate partition %v", n, parts)
			}
			seen[key] = true
		}
		if p.Length.Cmp(big.NewInt(count)) != 0 {
			t.Errorf("n=%d: iterated %d partitions, Length is %v", n, count, p.Length)
		}
	}
}

func TestIntegerPartitionsBorrowedMatchesAll(t *testing.T) {
	p, err := NewIntegerPartitions(8)
	if err != nil {
		t.Fatal(err)
	}
	all := slices.Collect(p.All())
	i := 0
	for parts := range p.AllBorrowed() {
		if !slices.Equal(parts, all[i]) {
			t.Errorf("iteration %d: AllBorrowed gave %v, All gave %v", i, parts, all[i])
		}
		i++
	}
	if i != len(all) {
		t.Errorf("AllBorrowed yielded %d partitions, All yielded %d", i, len(all))
	}
}
