package gocombinatorics

// fillBuf writes items from data at the given indices into buf.
func fillBuf[T any](buf []T, data []T, inds []int) {
	for i, idx := range inds {
		buf[i] = data[idx]
	}
}
