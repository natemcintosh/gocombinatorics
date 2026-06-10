package gocombinatorics

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"
)

// check_random asserts that Random draws are always members of the iteration
// space, that all tuples of a small space get drawn eventually (coverage), and
// that the same seed yields the same sequence of draws (determinism).
func check_random(
	t *testing.T,
	all func(yield func([]int, []int) bool),
	random func(*rand.Rand) ([]int, []int),
) {
	t.Helper()

	space := map[string][]int{}
	for inds, items := range all {
		if !slices.Equal(inds, items) {
			t.Fatalf("test expects identity data, got indices %v with items %v", inds, items)
		}
		space[fmt.Sprint(inds)] = inds
	}

	r := rand.New(rand.NewSource(42))
	seen := map[string]bool{}
	n_draws := 200 * len(space)
	for range n_draws {
		inds, items := random(r)
		if !slices.Equal(inds, items) {
			t.Fatalf("Random returned indices %v but items %v", inds, items)
		}
		key := fmt.Sprint(inds)
		if _, ok := space[key]; !ok {
			t.Fatalf("Random returned %v, which All() never yields", inds)
		}
		seen[key] = true
	}
	if len(seen) != len(space) {
		t.Errorf("after %d draws only %d of %d tuples were seen", n_draws, len(seen), len(space))
	}

	r1 := rand.New(rand.NewSource(7))
	r2 := rand.New(rand.NewSource(7))
	for range 20 {
		inds1, _ := random(r1)
		inds2, _ := random(r2)
		if !slices.Equal(inds1, inds2) {
			t.Fatalf("same seed gave different draws: %v vs %v", inds1, inds2)
		}
	}
}

func TestCombinationsRandom(t *testing.T) {
	c, err := NewCombinations(iota_slice(6), 3)
	if err != nil {
		t.Fatal(err)
	}
	check_random(t, c.All(), c.Random)
}

func TestCombinationsWithReplacementRandom(t *testing.T) {
	c, err := NewCombinationsWithReplacement(iota_slice(4), 3)
	if err != nil {
		t.Fatal(err)
	}
	check_random(t, c.All(), c.Random)
}

func TestPermutationsRandom(t *testing.T) {
	p, err := NewPermutations(iota_slice(5), 3)
	if err != nil {
		t.Fatal(err)
	}
	check_random(t, p.All(), p.Random)
}

func TestProductRandom(t *testing.T) {
	p, err := NewProduct(iota_slice(3), 3)
	if err != nil {
		t.Fatal(err)
	}
	check_random(t, p.All(), p.Random)
}

func TestPowersetRandom(t *testing.T) {
	p, err := NewPowerset(iota_slice(5))
	if err != nil {
		t.Fatal(err)
	}
	check_random(t, p.All(), p.Random)
}

// Random must work on spaces far too large to enumerate.
func TestRandomLargeLength(t *testing.T) {
	p, err := NewPermutations(iota_slice(50), 50)
	if err != nil {
		t.Fatal(err)
	}
	r := rand.New(rand.NewSource(1))
	inds, items := p.Random(r)
	if len(inds) != 50 || len(items) != 50 {
		t.Fatalf("expected 50-length permutation, got %d indices", len(inds))
	}
	sorted := slices.Clone(inds)
	slices.Sort(sorted)
	if !slices.Equal(sorted, iota_slice(50)) {
		t.Fatalf("draw is not a permutation of 0..49: %v", inds)
	}
}
