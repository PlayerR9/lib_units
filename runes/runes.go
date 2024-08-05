package runes

import (
	gcers "github.com/PlayerR9/go-commons/errors"
	gcint "github.com/PlayerR9/go-commons/ints"
	luc "github.com/PlayerR9/lib_units/common"
)

// LevenshteinTable is a table of words for the Levenshtein distance.
type LavenshteinTable struct {
	// words is the list of words.
	word_list [][]rune

	// word_length_list is the list of word lengths.
	word_length_list []int
}

// NewLevenshteinTable creates a new Levenshtein table
// with the given words.
//
// Parameters:
//   - words: The words to add to the table.
//
// Returns:
//   - *LevenshteinTable: The new Levenshtein table.
//   - error: An error if any of the words cannot be added to the table.
//
// Errors:
//   - *common.ErrAt: Whenever a word is not valid UTF-8.
//
// It is the same as creating an empty table and then adding the words to it.
func NewLevenshteinTable(words ...string) (*LavenshteinTable, error) {
	lt := &LavenshteinTable{}

	for i, word := range words {
		err := lt.AddWord(word)
		if err != nil {
			return nil, gcint.NewErrAt(i+1, word, err)
		}
	}

	return lt, nil
}

// AddWord adds a word to the table.
//
// Parameters:
//   - word: The word to add.
//
// Returns:
//   - error: An error of type *ErrInvalidUTF8Encoding if the word is not
//     valid UTF-8.
func (lt *LavenshteinTable) AddWord(word string) error {
	chars, err := StringToUtf8(word)
	if err != nil {
		return err
	}

	lt.word_list = append(lt.word_list, chars)
	lt.word_length_list = append(lt.word_length_list, len(chars))

	return nil
}

// Closest gets the closest word to a target.
//
// Parameters:
//   - target: The target.
//
// Returns:
//   - string: The closest word.
//   - error: The error if any occurs.
//
// Errors:
//   - *common.ErrInvalidParameter: If the target is empty.
//   - *ErrNoClosestWordFound: If no closest word is found.
func (lt *LavenshteinTable) Closest(target []rune) (string, error) {
	if len(target) == 0 {
		return "", gcers.NewErrInvalidParameter("target", gcers.NewErrEmpty("slice of runes"))
	}

	target_len := len(target)

	closest_idx := -1
	var min int

	for i, word := range lt.word_list {
		d := levenshtein_distance(target, target_len, word, lt.word_length_list[i])

		if closest_idx == -1 || d < min {
			min = d
			closest_idx = i
		}
	}

	if closest_idx == -1 {
		return "", NewErrNoClosestWordFound()
	}

	word := lt.word_list[closest_idx]

	return string(word), nil
}

// levenshteinDistance calculates the Levenshtein distance between two strings.
//
// Parameters:
//   - target: The target.
//   - target_len: The target length.
//   - other: The other.
//   - other_len: The other length.
//
// Returns:
//   - int: The Levenshtein distance.
func levenshtein_distance(target []rune, target_len int, other []rune, other_len int) int {
	matrix := make([][]int, 0, target_len+1)

	for i := 0; i <= target_len; i++ {
		row := make([]int, other_len+1)

		matrix = append(matrix, row)
	}

	// Initialize the matrix
	for i := 0; i <= target_len; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= other_len; j++ {
		matrix[0][j] = j
	}

	// Compute the distances
	for i := 1; i <= target_len; i++ {
		for j := 1; j <= other_len; j++ {
			if target[i-1] == other[j-1] {
				matrix[i][j] = matrix[i-1][j-1] // No operation needed
			} else {
				deletion := matrix[i-1][j] + 1
				insertion := matrix[i][j-1] + 1
				substitution := matrix[i-1][j-1] + 1

				min_first := luc.Min(deletion, insertion)
				min_second := luc.Min(min_first, substitution)
				matrix[i][j] = min_second
			}
		}
	}

	d := matrix[target_len][other_len]

	return d
}
