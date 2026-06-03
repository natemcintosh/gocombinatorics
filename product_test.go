package gocombinatorics

import (
	"reflect"
	"testing"
)

func TestNewProductErrors(t *testing.T) {
	testCases := []struct {
		desc string
		n    int
		k    int
	}{
		{desc: "n <= 0", n: 0, k: 1},
		{desc: "k <= 0", n: 1, k: 0},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			data := stepped_range(0, tC.n, 1)
			_, got_err := NewProduct(data, tC.k)
			if got_err == nil {
				t.Errorf("NewProduct(n=%d, k=%d) = nil error, want error", tC.n, tC.k)
			}
		})
	}
}

func TestProductNext(t *testing.T) {
	testCases := []struct {
		desc string
		n    int
		k    int
		want [][]int
	}{
		{
			desc: "n = 2, k = 2",
			n:    2,
			k:    2,
			want: [][]int{
				{0, 0},
				{0, 1},
				{1, 0},
				{1, 1},
			},
		},
		{
			desc: "n = 3, k = 1",
			n:    3,
			k:    1,
			want: [][]int{
				{0},
				{1},
				{2},
			},
		},
		{
			// k > n is allowed for products
			desc: "n = 2, k = 3",
			n:    2,
			k:    3,
			want: [][]int{
				{0, 0, 0},
				{0, 0, 1},
				{0, 1, 0},
				{0, 1, 1},
				{1, 0, 0},
				{1, 0, 1},
				{1, 1, 0},
				{1, 1, 1},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			data := stepped_range(0, tC.n, 1)
			p, err := NewProduct(data, tC.k)
			if err != nil {
				t.Errorf("NewProduct(n=%d, k=%d) = %v, want nil", tC.n, tC.k, err)
			}
			got := make([][]int, 0)
			for inds := range p.All() {
				got = append(got, inds)
			}
			if !reflect.DeepEqual(got, tC.want) {
				t.Errorf("Product(n=%d, k=%d) = %v, want %v", tC.n, tC.k, got, tC.want)
			}

			// Length should equal the number of tuples yielded.
			if int64(len(got)) != p.Length.Int64() {
				t.Errorf("Length = %v, but yielded %d tuples", p.Length, len(got))
			}
		})
	}
}

func BenchmarkProductAllVsBorrowed(b *testing.B) {
	benchmarks := []struct {
		desc string
		n    int
		k    int
	}{
		{desc: "n=10,k=3", n: 10, k: 3},
		{desc: "n=6,k=6", n: 6, k: 6},
	}
	for _, bm := range benchmarks {
		data := stepped_range(0, bm.n, 1)
		p, err := NewProduct(data, bm.k)
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
