package MathExt

import (
	"math/big"

	luc "github.com/PlayerR9/lib_units/common"

	gcers "github.com/PlayerR9/go-commons/errors"
)

// Add adds two numbers of the same base. Both numbers are Least Significant Digit
// (LSD) first.
//
// Parameters:
//   - n1: The first number to add.
//   - n2: The second number to add.
//   - base: The base of the numbers.
//
// Returns:
//   - []int: The sum of the two numbers. Nil if the base is less than or equal to 0.
func Add(n1, n2 []int, base int) []int {
	if base <= 0 {
		return nil
	}

	if base == 1 {
		return make([]int, len(n1)+len(n2))
	}

	maxLen := len(n1)
	if len(n2) > maxLen {
		maxLen = len(n2)
	}

	// Add the two binary numbers.
	result := make([]int, maxLen)
	var carry int

	// Add the digits for the common length of n1 and n2
	for i := 0; i < len(n1) && i < len(n2); i++ {
		result[i] = n1[i] + n2[i] + carry
		carry = result[i] / base
		result[i] %= base
	}

	// Add the remaining digits of the longer number
	for i := len(n2); i < len(n1); i++ {
		result[i] = n1[i] + carry
		carry = result[i] / base
		result[i] %= base
	}

	for i := len(n1); i < len(n2); i++ {
		result[i] = n2[i] + carry
		carry = result[i] / base
		result[i] %= base
	}

	if carry > 0 {
		result = append(result, carry)
	}

	return result
}

// Subtract subtracts two numbers of the same base. Both numbers are Least Significant
// Digit (LSD) first.
//
// Parameters:
//   - n1: The number to subtract from.
//   - n2: The number to subtract.
//   - base: The base of the numbers.
//
// Returns:
//   - []int: The result of the subtraction.
//   - error: An error if the subtraction failed.
//
// Errors:
//   - *ErrSubtractionUnderflow: The subtraction resulted in a negative number.
//   - *errors.ErrInvalidParameter: The base is less than or equal to 0.
func Subtract(n1, n2 []int, base int) ([]int, error) {
	if base <= 0 {
		return nil, gcers.NewErrInvalidParameter("base", luc.NewErrGT(0))
	}

	if base == 1 {
		return make([]int, len(n1)), nil
	}

	// Subtract the two binary numbers.
	result := make([]int, len(n1))
	var borrow int

	// Subtract the digits for the common length of n1 and n2
	for i := 0; i < len(n1) && i < len(n2); i++ {
		result[i] = n1[i] - n2[i] - borrow

		if result[i] < 0 {
			result[i] += base
			borrow = 1
		} else {
			borrow = 0
		}
	}

	// Subtract the remaining digits of the longer number
	for i := len(n2); i < len(n1); i++ {
		result[i] = n1[i] - borrow

		if result[i] < 0 {
			result[i] += base
			borrow = 1
		} else {
			borrow = 0
		}
	}

	if borrow > 0 {
		return nil, NewErrSubtractionUnderflow()
	}

	// Remove leading zeros
	i := len(result) - 1
	for ; i >= 0 && result[i] == 0; i-- {
	}
	result = result[:i+1]

	if len(result) == 0 {
		result = []int{0}
	}

	return result, nil
}

// IntToBigInt converts an integer to a big.Int.
//
// Parameters:
//   - n: The integer to convert.
//
// Returns:
//   - *big.Int: The big.Int representation of the integer.
func IntToBigInt(n int) *big.Int {
	bi := new(big.Int)

	res := bi.SetInt64(int64(n))

	return res
}
