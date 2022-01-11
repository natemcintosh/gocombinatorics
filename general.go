package gocombinatorics

import "errors"

// Multiple types all adhere to this interface
type CombinationLike[T any] interface {
	Next() bool
	LenInds() int
	Indices() []int
	Items() []T
}

var errBufferIndicesMismatch = errors.New("Length of buffer and indices did not match")

// fill_buffer will fill up a buffer with the items from `data` at `indices`
func fill_buffer[T any](buffer []T, data []T, indices []int) error {
	if len(indices) != len(buffer) {
		return errBufferIndicesMismatch
	}

	for buff_idx, data_idx := range indices {
		buffer[buff_idx] = data[data_idx]
	}
	return nil
}
