package common

import (
	"cmp"
)

// Min is a function that takes two parameters, a and b, of any type T
// according to the cmp.Ordered interface and returns the smaller of the two values.
//
// Parameters:
//   - a, b: The two values to compare.
//
// Return:
//   - T: The smaller of the two values.
func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	} else {
		return b
	}
}

// Max is a function that takes two parameters, a and b, of any type T
// according to the cmp.Ordered function and returns the larger of the two values.
//
// Parameters:
//   - a, b: The two values to compare.
//
// Return:
//   - T: The larger of the two values.
func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	} else {
		return b
	}
}
