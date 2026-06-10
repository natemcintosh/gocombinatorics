package gocombinatorics

import (
	"reflect"
	"testing"
)

func TestNewProductOfErrors(t *testing.T) {
	testCases := []struct {
		desc string
		axes [][]int
	}{
		{desc: "no axes", axes: [][]int{}},
		{desc: "empty axis", axes: [][]int{{0, 1}, {}}},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, got_err := NewProductOf(tC.axes...)
			if got_err == nil {
				t.Errorf("NewProductOf(%v) = nil error, want error", tC.axes)
			}
		})
	}
}

func TestProductOfNext(t *testing.T) {
	testCases := []struct {
		desc string
		axes [][]int
		want [][]int
	}{
		{
			desc: "2 x 3",
			axes: [][]int{{0, 1}, {0, 1, 2}},
			want: [][]int{
				{0, 0},
				{0, 1},
				{0, 2},
				{1, 0},
				{1, 1},
				{1, 2},
			},
		},
		{
			desc: "single axis",
			axes: [][]int{{0, 1, 2}},
			want: [][]int{
				{0},
				{1},
				{2},
			},
		},
		{
			desc: "2 x 1 x 2",
			axes: [][]int{{0, 1}, {0}, {0, 1}},
			want: [][]int{
				{0, 0, 0},
				{0, 0, 1},
				{1, 0, 0},
				{1, 0, 1},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p, err := NewProductOf(tC.axes...)
			if err != nil {
				t.Errorf("NewProductOf(%v) = %v, want nil", tC.axes, err)
			}
			got := make([][]int, 0)
			for inds := range p.All() {
				got = append(got, inds)
			}
			if !reflect.DeepEqual(got, tC.want) {
				t.Errorf("ProductOf(%v) = %v, want %v", tC.axes, got, tC.want)
			}

			// Length should equal the number of tuples yielded.
			if int64(len(got)) != p.Length.Int64() {
				t.Errorf("Length = %v, but yielded %d tuples", p.Length, len(got))
			}
		})
	}
}

// ProductOf items must come from the corresponding axis, not a shared slice.
func TestProductOfItems(t *testing.T) {
	p, err := NewProductOf([]string{"a", "b"}, []string{"x", "y", "z"})
	if err != nil {
		t.Fatal(err)
	}
	want := [][]string{
		{"a", "x"}, {"a", "y"}, {"a", "z"},
		{"b", "x"}, {"b", "y"}, {"b", "z"},
	}
	got := make([][]string, 0)
	for _, items := range p.All() {
		got = append(got, items)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ProductOf items = %v, want %v", got, want)
	}
}

func BenchmarkProductOfAllVsBorrowed(b *testing.B) {
	benchmarks := []struct {
		desc string
		axes [][]int
	}{
		{desc: "10x10x10", axes: [][]int{
			stepped_range(0, 10, 1), stepped_range(0, 10, 1), stepped_range(0, 10, 1),
		}},
		{desc: "2x12x3x9x5x4", axes: [][]int{
			stepped_range(0, 2, 1), stepped_range(0, 12, 1), stepped_range(0, 3, 1),
			stepped_range(0, 9, 1), stepped_range(0, 5, 1), stepped_range(0, 4, 1),
		}},
	}
	for _, bm := range benchmarks {
		p, err := NewProductOf(bm.axes...)
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
