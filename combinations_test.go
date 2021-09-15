package gocombinatorics

import (
	"encoding/csv"
	"errors"
	"log"
	"math/big"
	"os"
	"reflect"
	"strconv"
	"testing"
)

func TestNewCombinationsErrors(t *testing.T) {
	testCases := []struct {
		desc        string
		n           int
		k           int
		want_struct *Combinations
		want_err    error
	}{
		{
			desc:        "n <= 0",
			n:           0,
			k:           1,
			want_struct: nil,
			want_err:    errors.New("n must be greater than 0"),
		},
		{
			desc:        "k <= 0",
			n:           1,
			k:           0,
			want_struct: nil,
			want_err:    errors.New("k must be greater than 0"),
		},
		{
			desc:        "k > n",
			n:           1,
			k:           2,
			want_struct: nil,
			want_err:    errors.New("k must be less than or equal to n"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, got_err := NewCombinations(tC.n, tC.k)
			if got_err == nil {
				t.Errorf("NewCombinations(%d, %d) = %v, want %v", tC.n, tC.k, got_err, tC.want_err)
			}
		})
	}
}

func TestCombinationsNext(t *testing.T) {
	testCases := []struct {
		desc string
		n    int
		k    int
		want [][]int
	}{
		{
			desc: "n = 3, k = 2",
			n:    3,
			k:    2,
			want: [][]int{
				{0, 1},
				{0, 2},
				{1, 2},
			},
		},
		{
			desc: "n = 5, k = 3",
			n:    5,
			k:    3,
			want: [][]int{
				{0, 1, 2},
				{0, 1, 3},
				{0, 1, 4},
				{0, 2, 3},
				{0, 2, 4},
				{0, 3, 4},
				{1, 2, 3},
				{1, 2, 4},
				{1, 3, 4},
				{2, 3, 4},
			},
		},
		{
			desc: "n = 2, k = 2",
			n:    2,
			k:    2,
			want: [][]int{
				{0, 1},
			},
		},
		{
			desc: "n = 10, k = 8",
			n:    10,
			k:    8,
			want: [][]int{
				{0, 1, 2, 3, 4, 5, 6, 7},
				{0, 1, 2, 3, 4, 5, 6, 8},
				{0, 1, 2, 3, 4, 5, 6, 9},
				{0, 1, 2, 3, 4, 5, 7, 8},
				{0, 1, 2, 3, 4, 5, 7, 9},
				{0, 1, 2, 3, 4, 5, 8, 9},
				{0, 1, 2, 3, 4, 6, 7, 8},
				{0, 1, 2, 3, 4, 6, 7, 9},
				{0, 1, 2, 3, 4, 6, 8, 9},
				{0, 1, 2, 3, 4, 7, 8, 9},
				{0, 1, 2, 3, 5, 6, 7, 8},
				{0, 1, 2, 3, 5, 6, 7, 9},
				{0, 1, 2, 3, 5, 6, 8, 9},
				{0, 1, 2, 3, 5, 7, 8, 9},
				{0, 1, 2, 3, 6, 7, 8, 9},
				{0, 1, 2, 4, 5, 6, 7, 8},
				{0, 1, 2, 4, 5, 6, 7, 9},
				{0, 1, 2, 4, 5, 6, 8, 9},
				{0, 1, 2, 4, 5, 7, 8, 9},
				{0, 1, 2, 4, 6, 7, 8, 9},
				{0, 1, 2, 5, 6, 7, 8, 9},
				{0, 1, 3, 4, 5, 6, 7, 8},
				{0, 1, 3, 4, 5, 6, 7, 9},
				{0, 1, 3, 4, 5, 6, 8, 9},
				{0, 1, 3, 4, 5, 7, 8, 9},
				{0, 1, 3, 4, 6, 7, 8, 9},
				{0, 1, 3, 5, 6, 7, 8, 9},
				{0, 1, 4, 5, 6, 7, 8, 9},
				{0, 2, 3, 4, 5, 6, 7, 8},
				{0, 2, 3, 4, 5, 6, 7, 9},
				{0, 2, 3, 4, 5, 6, 8, 9},
				{0, 2, 3, 4, 5, 7, 8, 9},
				{0, 2, 3, 4, 6, 7, 8, 9},
				{0, 2, 3, 5, 6, 7, 8, 9},
				{0, 2, 4, 5, 6, 7, 8, 9},
				{0, 3, 4, 5, 6, 7, 8, 9},
				{1, 2, 3, 4, 5, 6, 7, 8},
				{1, 2, 3, 4, 5, 6, 7, 9},
				{1, 2, 3, 4, 5, 6, 8, 9},
				{1, 2, 3, 4, 5, 7, 8, 9},
				{1, 2, 3, 4, 6, 7, 8, 9},
				{1, 2, 3, 5, 6, 7, 8, 9},
				{1, 2, 4, 5, 6, 7, 8, 9},
				{1, 3, 4, 5, 6, 7, 8, 9},
				{2, 3, 4, 5, 6, 7, 8, 9},
			},
		},
		{
			desc: "n = 10, k = 3",
			n:    10,
			k:    3,
			want: [][]int{
				{0, 1, 2},
				{0, 1, 3},
				{0, 1, 4},
				{0, 1, 5},
				{0, 1, 6},
				{0, 1, 7},
				{0, 1, 8},
				{0, 1, 9},
				{0, 2, 3},
				{0, 2, 4},
				{0, 2, 5},
				{0, 2, 6},
				{0, 2, 7},
				{0, 2, 8},
				{0, 2, 9},
				{0, 3, 4},
				{0, 3, 5},
				{0, 3, 6},
				{0, 3, 7},
				{0, 3, 8},
				{0, 3, 9},
				{0, 4, 5},
				{0, 4, 6},
				{0, 4, 7},
				{0, 4, 8},
				{0, 4, 9},
				{0, 5, 6},
				{0, 5, 7},
				{0, 5, 8},
				{0, 5, 9},
				{0, 6, 7},
				{0, 6, 8},
				{0, 6, 9},
				{0, 7, 8},
				{0, 7, 9},
				{0, 8, 9},
				{1, 2, 3},
				{1, 2, 4},
				{1, 2, 5},
				{1, 2, 6},
				{1, 2, 7},
				{1, 2, 8},
				{1, 2, 9},
				{1, 3, 4},
				{1, 3, 5},
				{1, 3, 6},
				{1, 3, 7},
				{1, 3, 8},
				{1, 3, 9},
				{1, 4, 5},
				{1, 4, 6},
				{1, 4, 7},
				{1, 4, 8},
				{1, 4, 9},
				{1, 5, 6},
				{1, 5, 7},
				{1, 5, 8},
				{1, 5, 9},
				{1, 6, 7},
				{1, 6, 8},
				{1, 6, 9},
				{1, 7, 8},
				{1, 7, 9},
				{1, 8, 9},
				{2, 3, 4},
				{2, 3, 5},
				{2, 3, 6},
				{2, 3, 7},
				{2, 3, 8},
				{2, 3, 9},
				{2, 4, 5},
				{2, 4, 6},
				{2, 4, 7},
				{2, 4, 8},
				{2, 4, 9},
				{2, 5, 6},
				{2, 5, 7},
				{2, 5, 8},
				{2, 5, 9},
				{2, 6, 7},
				{2, 6, 8},
				{2, 6, 9},
				{2, 7, 8},
				{2, 7, 9},
				{2, 8, 9},
				{3, 4, 5},
				{3, 4, 6},
				{3, 4, 7},
				{3, 4, 8},
				{3, 4, 9},
				{3, 5, 6},
				{3, 5, 7},
				{3, 5, 8},
				{3, 5, 9},
				{3, 6, 7},
				{3, 6, 8},
				{3, 6, 9},
				{3, 7, 8},
				{3, 7, 9},
				{3, 8, 9},
				{4, 5, 6},
				{4, 5, 7},
				{4, 5, 8},
				{4, 5, 9},
				{4, 6, 7},
				{4, 6, 8},
				{4, 6, 9},
				{4, 7, 8},
				{4, 7, 9},
				{4, 8, 9},
				{5, 6, 7},
				{5, 6, 8},
				{5, 6, 9},
				{5, 7, 8},
				{5, 7, 9},
				{5, 8, 9},
				{6, 7, 8},
				{6, 7, 9},
				{6, 8, 9},
				{7, 8, 9},
			},
		},
		{
			desc: "n = 200, k = 3",
			n:    200,
			k:    3,
			want: csv_to_2d_int_array("200_combo_3.csv"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			combinations, err := NewCombinations(tC.n, tC.k)
			if err != nil {
				t.Errorf("NewCombinations(%d, %d) = %v, want nil", tC.n, tC.k, err)
			}
			got := make([][]int, 0)
			for combinations.Next() {
				// We need to append a copy of combinations.Inds to got
				next_set_of_indices := make([]int, len(combinations.Inds))
				copy(next_set_of_indices, combinations.Indices())
				got = append(got, next_set_of_indices)
			}

			if !reflect.DeepEqual(got, tC.want) {
				t.Errorf("Combinations(%d, %d) = %v, want %v", tC.n, tC.k, got, tC.want)
			}
		})
	}
}

func BenchmarkCombinationsNext(b *testing.B) {
	benchmarks := []struct {
		desc string
		n    int
		k    int
	}{
		{
			desc: "n = 3, k = 2",
			n:    3,
			k:    2,
		},
		{
			desc: "n = 4, k = 3",
			n:    4,
			k:    3,
		},
		{
			desc: "n = 5, k = 4",
			n:    5,
			k:    4,
		},
		{
			desc: "n = 10, k = 8",
			n:    10,
			k:    8,
		},
		{
			desc: "n = 10, k = 3",
			n:    10,
			k:    3,
		},
		{
			desc: "n = 200, k = 3",
			n:    200,
			k:    3,
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				combinations, err := NewCombinations(bm.n, bm.k)
				if err != nil {
					b.Errorf("NewCombinations(%d, %d) = %v, want nil", bm.n, bm.k, err)
				}
				for combinations.Next() {

				}
			}

		})
	}
}

func TestFactorial(t *testing.T) {
	testCases := []struct {
		desc string
		n    int64
		want *big.Int
	}{
		{
			desc: "3",
			n:    3,
			want: big.NewInt(6),
		},
		{
			desc: "4",
			n:    4,
			want: big.NewInt(24),
		},
		{
			desc: "10",
			n:    10,
			want: big.NewInt(3628800),
		},
		{
			desc: "20",
			n:    20,
			want: bigIntFromString("2432902008176640000"),
		},
		{
			desc: "100",
			n:    100,
			want: bigIntFromString("93326215443944152681699238856266700490715968264381621468592963895217599993229915608941463976156518286253697920827223758251185210916864000000000000000000000000"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := factorial(tC.n)
			if got.Cmp(tC.want) != 0 {
				t.Errorf("Factorial(%d) = %v, want %v", tC.n, got, tC.want)
			}
		})
	}
}

func BenchmarkFactorial(b *testing.B) {
	benchmarks := []struct {
		desc string
		in   int64
	}{
		{
			desc: "3",
			in:   3,
		},
		{
			desc: "4",
			in:   4,
		},
		{
			desc: "10",
			in:   10,
		},
		{
			desc: "20",
			in:   20,
		},
		{
			desc: "100",
			in:   100,
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				factorial(bm.in)
			}

		})
	}
}

func TestNChooseK(t *testing.T) {
	testCases := []struct {
		desc string
		n, k uint64
		want *big.Int
	}{
		{
			desc: "3 choose 2",
			n:    3,
			k:    2,
			want: big.NewInt(3),
		},
		{
			desc: "4 choose 2",
			n:    4,
			k:    2,
			want: big.NewInt(6),
		},
		{
			desc: "10 choose 2",
			n:    10,
			k:    2,
			want: big.NewInt(45),
		},
		{
			desc: "20 choose 2",
			n:    20,
			k:    2,
			want: big.NewInt(190),
		},
		{
			desc: "100 choose 2",
			n:    100,
			k:    2,
			want: big.NewInt(4950),
		},
		{
			desc: "100 choose 34",
			n:    100,
			k:    34,
			want: bigIntFromString("580717429720889409486981450"),
		},
		{
			desc: "1 choose 1",
			n:    1,
			k:    1,
			want: big.NewInt(1),
		},
		{
			desc: "2 choose 3",
			n:    2,
			k:    3,
			want: big.NewInt(0),
		},
		{
			desc: "1000 choose 832",
			n:    1000,
			k:    832,
			want: bigIntFromString("1359578307154377147929220245480317022063628494700109214161339185998491705982730220827187173838995764985909775118253237706078706638334815696866772700709347742606303804123006462298452381996277060875"),
		},
		{
			desc: "0 choose 0",
			n:    0,
			k:    0,
			want: big.NewInt(0),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := nchoosek(tC.n, tC.k)
			if got.Cmp(tC.want) != 0 {
				t.Errorf("NChooseK(%d, %d) = %v, want %v", tC.n, tC.k, got, tC.want)
			}
		})
	}
}

func BenchmarkNChooseK(b *testing.B) {
	benchmarks := []struct {
		desc string
		n    uint64
		k    uint64
	}{
		{
			desc: "3 choose 2",
			n:    3,
			k:    2,
		},
		{
			desc: "4 choose 2",
			n:    4,
			k:    2,
		},
		{
			desc: "10 choose 2",
			n:    10,
			k:    2,
		},
		{
			desc: "20 choose 2",
			n:    20,
			k:    2,
		},
		{
			desc: "100 choose 2",
			n:    100,
			k:    2,
		},
		{
			desc: "100 choose 34",
			n:    100,
			k:    34,
		},
		{
			desc: "1 choose 1",
			n:    1,
			k:    1,
		},
		{
			desc: "2 choose 3",
			n:    2,
			k:    3,
		},
		{
			desc: "1000 choose 832",
			n:    1000,
			k:    832,
		},
		{
			desc: "0 choose 0",
			n:    0,
			k:    0,
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				nchoosek(bm.n, bm.k)
			}

		})
	}
}

func bigIntFromString(s string) *big.Int {
	i, _ := new(big.Int).SetString(s, 10)
	return i
}

// csv_to_2d_int_array will read in a csv file and return a 2d array of ints
func csv_to_2d_int_array(path string) [][]int {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Convert the strings to ints
	// First make a 2d slice of ints
	int_lines := make([][]int, len(lines))

	// For each line of strings, convert it to a slice of ints
	for line_idx, line := range lines {
		int_lines[line_idx] = make([]int, len(line))
		for col_idx, s := range line {
			int_lines[line_idx][col_idx], err = strconv.Atoi(s)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return int_lines

}
