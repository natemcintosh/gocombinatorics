package gocombinatorics

import (
	"reflect"
	"testing"
)

func TestNewPowersetErrors(t *testing.T) {
	data := stepped_range(0, 0, 1)
	if _, err := NewPowerset(data); err == nil {
		t.Errorf("NewPowerset(n=0) = nil error, want error")
	}
}

func TestPowersetNext(t *testing.T) {
	testCases := []struct {
		desc string
		n    int
		want [][]int
	}{
		{
			desc: "n = 1",
			n:    1,
			want: [][]int{
				{},
				{0},
			},
		},
		{
			desc: "n = 3",
			n:    3,
			want: [][]int{
				{},
				{0},
				{1},
				{2},
				{0, 1},
				{0, 2},
				{1, 2},
				{0, 1, 2},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			data := stepped_range(0, tC.n, 1)
			p, err := NewPowerset(data)
			if err != nil {
				t.Errorf("NewPowerset(n=%d) = %v, want nil", tC.n, err)
			}
			got := make([][]int, 0)
			for inds := range p.All() {
				got = append(got, inds)
			}
			if !reflect.DeepEqual(got, tC.want) {
				t.Errorf("Powerset(n=%d) = %v, want %v", tC.n, got, tC.want)
			}

			// Length should equal the number of subsets yielded.
			if int64(len(got)) != p.Length.Int64() {
				t.Errorf("Length = %v, but yielded %d subsets", p.Length, len(got))
			}
		})
	}
}

func BenchmarkPowersetAllVsBorrowed(b *testing.B) {
	benchmarks := []struct {
		desc string
		n    int
	}{
		{desc: "n=10", n: 10},
		{desc: "n=16", n: 16},
	}
	for _, bm := range benchmarks {
		data := stepped_range(0, bm.n, 1)
		p, err := NewPowerset(data)
		if err != nil {
			b.Fatal(err)
		}

		b.Run("All/"+bm.desc, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for range p.All() {
				}
			}
		})
		b.Run("AllBorrowed/"+bm.desc, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for range p.AllBorrowed() {
				}
			}
		})
	}
}
