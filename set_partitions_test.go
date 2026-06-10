package gocombinatorics

import (
	"fmt"
	"math/big"
	"slices"
	"testing"
)

func TestSetPartitionsErrors(t *testing.T) {
	if _, err := NewSetPartitions([]int{}); err == nil {
		t.Error("NewSetPartitions with empty input should have errored")
	}
	if _, err := NewSetPartitionsK([]int{}, 1); err == nil {
		t.Error("NewSetPartitionsK with empty input should have errored")
	}
	if _, err := NewSetPartitionsK([]int{1, 2}, 0); err == nil {
		t.Error("NewSetPartitionsK with k=0 should have errored")
	}
	if _, err := NewSetPartitionsK([]int{1, 2}, 3); err == nil {
		t.Error("NewSetPartitionsK with k > n should have errored")
	}
}

func TestSetPartitionsLength(t *testing.T) {
	// Bell numbers B(1)..B(8), OEIS A000110.
	bells := []int64{1, 2, 5, 15, 52, 203, 877, 4140}
	for i, w := range bells {
		n := i + 1
		data := make([]int, n)
		s, err := NewSetPartitions(data)
		if err != nil {
			t.Fatalf("NewSetPartitions(n=%d): %v", n, err)
		}
		if s.Length.Cmp(big.NewInt(w)) != 0 {
			t.Errorf("B(%d): got Length %v, want %d", n, s.Length, w)
		}
	}

	// Stirling numbers of the second kind S(n, k).
	stirlings := []struct {
		n, k int
		want int64
	}{
		{1, 1, 1},
		{3, 2, 3},
		{4, 2, 7},
		{5, 3, 25},
		{6, 3, 90},
		{7, 4, 350},
		{10, 5, 42525},
	}
	for _, tc := range stirlings {
		data := make([]int, tc.n)
		s, err := NewSetPartitionsK(data, tc.k)
		if err != nil {
			t.Fatalf("NewSetPartitionsK(n=%d, k=%d): %v", tc.n, tc.k, err)
		}
		if s.Length.Cmp(big.NewInt(tc.want)) != 0 {
			t.Errorf("S(%d, %d): got Length %v, want %d", tc.n, tc.k, s.Length, tc.want)
		}
	}
}

func TestSetPartitionsAll(t *testing.T) {
	s, err := NewSetPartitions([]string{"a", "b", "c"})
	if err != nil {
		t.Fatal(err)
	}
	want := [][][]string{
		{{"a", "b", "c"}},
		{{"a", "b"}, {"c"}},
		{{"a", "c"}, {"b"}},
		{{"a"}, {"b", "c"}},
		{{"a"}, {"b"}, {"c"}},
	}
	i := 0
	for assignment, blocks := range s.All() {
		if i >= len(want) {
			t.Fatalf("yielded more than %d partitions", len(want))
		}
		if len(assignment) != 3 {
			t.Errorf("partition %d: assignment %v has length %d, want 3", i, assignment, len(assignment))
		}
		if fmt.Sprint(blocks) != fmt.Sprint(want[i]) {
			t.Errorf("partition %d: got %v, want %v", i, blocks, want[i])
		}
		i++
	}
	if i != len(want) {
		t.Errorf("yielded %d partitions, want %d", i, len(want))
	}
}

func TestSetPartitionsKAll(t *testing.T) {
	s, err := NewSetPartitionsK([]int{1, 2, 3, 4}, 2)
	if err != nil {
		t.Fatal(err)
	}
	want := [][][]int{
		{{1, 2, 3}, {4}},
		{{1, 2, 4}, {3}},
		{{1, 2}, {3, 4}},
		{{1, 3, 4}, {2}},
		{{1, 3}, {2, 4}},
		{{1, 4}, {2, 3}},
		{{1}, {2, 3, 4}},
	}
	got := [][][]int{}
	for _, blocks := range s.All() {
		got = append(got, blocks)
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSetPartitionsProperties(t *testing.T) {
	check := func(t *testing.T, s *SetPartitions[int], n, k int) {
		t.Helper()
		count := int64(0)
		seen := make(map[string]bool)
		for assignment, blocks := range s.All() {
			count++
			// Assignment is a restricted growth string.
			maxSeen := -1
			for i, b := range assignment {
				if b < 0 || b > maxSeen+1 {
					t.Fatalf("n=%d k=%d: assignment %v invalid at %d", n, k, assignment, i)
				}
				maxSeen = max(maxSeen, b)
			}
			if k > 0 && maxSeen+1 != k {
				t.Fatalf("n=%d k=%d: assignment %v has %d blocks", n, k, assignment, maxSeen+1)
			}
			if len(blocks) != maxSeen+1 {
				t.Fatalf("n=%d k=%d: %d blocks but assignment max %d", n, k, len(blocks), maxSeen)
			}
			// Blocks are non-empty and contain each element exactly once.
			total := 0
			for _, b := range blocks {
				if len(b) == 0 {
					t.Fatalf("n=%d k=%d: empty block in %v", n, k, blocks)
				}
				total += len(b)
			}
			if total != n {
				t.Fatalf("n=%d k=%d: blocks %v cover %d elements, want %d", n, k, blocks, total, n)
			}
			key := fmt.Sprint(assignment)
			if seen[key] {
				t.Fatalf("n=%d k=%d: duplicate partition %v", n, k, assignment)
			}
			seen[key] = true
		}
		if s.Length.Cmp(big.NewInt(count)) != 0 {
			t.Errorf("n=%d k=%d: iterated %d partitions, Length is %v", n, k, count, s.Length)
		}
	}

	for n := 1; n <= 8; n++ {
		data := make([]int, n)
		for i := range data {
			data[i] = i
		}
		s, err := NewSetPartitions(data)
		if err != nil {
			t.Fatal(err)
		}
		check(t, s, n, 0)
		for k := 1; k <= n; k++ {
			sk, err := NewSetPartitionsK(data, k)
			if err != nil {
				t.Fatal(err)
			}
			check(t, sk, n, k)
		}
	}
}

func TestSetPartitionsBorrowedMatchesAll(t *testing.T) {
	data := []int{10, 20, 30, 40, 50}
	s, err := NewSetPartitions(data)
	if err != nil {
		t.Fatal(err)
	}
	allAssignments := [][]int{}
	allBlocks := [][][]int{}
	for assignment, blocks := range s.All() {
		allAssignments = append(allAssignments, assignment)
		allBlocks = append(allBlocks, blocks)
	}
	i := 0
	for assignment, blocks := range s.AllBorrowed() {
		if !slices.Equal(assignment, allAssignments[i]) {
			t.Errorf("iteration %d: AllBorrowed assignment %v, All gave %v", i, assignment, allAssignments[i])
		}
		if fmt.Sprint(blocks) != fmt.Sprint(allBlocks[i]) {
			t.Errorf("iteration %d: AllBorrowed blocks %v, All gave %v", i, blocks, allBlocks[i])
		}
		i++
	}
	if i != len(allAssignments) {
		t.Errorf("AllBorrowed yielded %d partitions, All yielded %d", i, len(allAssignments))
	}
}
