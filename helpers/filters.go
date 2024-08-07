package helpers

import (
	lus "github.com/PlayerR9/go-commons/slices"
)

// EvalOneFunc is a function that evaluates one element.
//
// Parameters:
//   - elem: The element to evaluate.
//
// Returns:
//   - R: The result of the evaluation.
//   - error: An error if the evaluation failed.
type EvalOneFunc[E, R any] func(elem E) (R, error)

// FilterIsSuccess filters any helper that is not successful.
//
// Parameters:
//   - h: The helper to filter.
//
// Returns:
//   - bool: True if the helper is successful, false otherwise.
//
// Behaviors:
//   - It assumes that the h is not nil.
func FilterIsSuccess[T Helperer[O], O any](h T) bool {
	_, err := h.Data()
	return err == nil
}

// FilterByPositiveWeight is a function that iterates over weight results and
// returns the elements with the maximum weight.
//
// Parameters:
//   - S: slice of weight results.
//
// Returns:
//   - []T: slice of elements with the maximum weight.
//
// Behaviors:
//   - If S is empty, the function returns a nil slice.
//   - If multiple elements have the same maximum weight, they are all returned.
//   - If S contains only one element, that element is returned.
func FilterByPositiveWeight[T Helperer[O], O any](S []T) []T {
	if len(S) == 0 {
		return nil
	}

	maxWeight := S[0].Weight()
	indices := []int{0}

	for i, e := range S[1:] {
		currentWeight := e.Weight()

		if currentWeight > maxWeight {
			maxWeight = currentWeight
			indices = []int{i + 1}
		} else if currentWeight == maxWeight {
			indices = append(indices, i+1)
		}
	}

	solution := make([]T, 0, len(indices))

	for _, index := range indices {
		solution = append(solution, S[index])
	}

	return solution
}

// FilterByNegativeWeight is a function that iterates over weight results and
// returns the elements with the minimum weight.
//
// Parameters:
//   - S: slice of weight results.
//
// Returns:
//   - []T: slice of elements with the minimum weight.
//
// Behaviors:
//   - If S is empty, the function returns a nil slice.
//   - If multiple elements have the same minimum weight, they are all returned.
//   - If S contains only one element, that element is returned.
func FilterByNegativeWeight[T Helperer[O], O any](S []T) []T {
	if len(S) == 0 {
		return nil
	}

	minWeight := S[0].Weight()
	indices := []int{0}

	for i, e := range S[1:] {
		currentWeight := e.Weight()

		if currentWeight < minWeight {
			minWeight = currentWeight
			indices = []int{i + 1}
		} else if currentWeight == minWeight {
			indices = append(indices, i+1)
		}
	}

	solution := make([]T, 0, len(indices))
	for _, index := range indices {
		solution = append(solution, S[index])
	}

	return solution
}

// SuccessOrFail returns the results with the maximum weight.
//
// Parameters:
//   - batch: The slice of results.
//   - useMax: True if the maximum weight should be used, false otherwise.
//
// Returns:
//   - []*luc.Pair[O, error]: The results with the maximum weight.
//   - bool: True if the slice was filtered, false otherwise.
//
// Behaviors:
//   - If the slice is empty, the function returns a nil slice and true.
//   - The result can either be the sucessful results or the original slice.
//     Nonetheless, the maximum weight is always applied.
func SuccessOrFail[T Helperer[O], O any](batch []T, useMax bool) ([]T, bool) {
	// 1. Remove nil elements.
	if len(batch) == 0 {
		return nil, true
	}

	success, fail := lus.SFSeparate(batch, FilterIsSuccess[T, O])

	var target, solution []T

	if len(success) == 0 {
		target = fail
	} else {
		target = success
	}

	if useMax {
		solution = FilterByPositiveWeight(target)
	} else {
		solution = FilterByNegativeWeight(target)
	}

	return solution, len(success) > 0
}

// EvaluateSimpleHelpers is a function that evaluates a batch of helpers and returns
// the results.
//
// Parameters:
//   - batch: The slice of helpers.
//   - f: The evaluation function.
//
// Returns:
//   - []*SimpleHelper[O]: The results of the evaluation.
//   - bool: True if the slice was filtered, false otherwise.
//
// Behaviors:
//   - This function returns either the successful results or the original slice.
func EvaluateSimpleHelpers[T, O any](batch []T, f EvalOneFunc[T, O]) ([]*SimpleHelper[O], bool) {
	if len(batch) == 0 || f == nil {
		return nil, true
	}

	solutions := make([]*SimpleHelper[O], 0, len(batch))

	for _, h := range batch {
		res, err := f(h)

		helper := NewSimpleHelper(res, err)
		solutions = append(solutions, helper)
	}

	success, fail := lus.SFSeparate(solutions, FilterIsSuccess)

	var result []*SimpleHelper[O]

	if len(success) == 0 {
		result = fail
	} else {
		result = success
	}

	return result, len(success) > 0
}

// EvaluateWeightHelpers is a function that evaluates a batch of helpers and returns
// the results.
//
// Parameters:
//   - batch: The slice of helpers.
//   - f: The evaluation function.
//   - wf: The weight function.
//   - useMax: True if the maximum weight should be used, false otherwise.
//
// Returns:
//   - []*WeightedHelper[O]: The results of the evaluation.
//   - bool: True if the slice was filtered, false otherwise.
//
// Behaviors:
//   - This function returns either the successful results or the original slice.
func EvaluateWeightHelpers[T, O any](batch []T, f EvalOneFunc[T, O], wf WeightFunc[T], useMax bool) ([]*WeightedHelper[O], bool) {
	if len(batch) == 0 || f == nil || wf == nil {
		return nil, true
	}

	solutions := make([]*WeightedHelper[O], 0, len(batch))

	for _, h := range batch {
		res, err := f(h)

		weight, ok := wf(h)
		if !ok {
			continue
		}

		h := NewWeightedHelper(res, err, weight)
		solutions = append(solutions, h)
	}

	success, fail := lus.SFSeparate(solutions, FilterIsSuccess)

	var target, result []*WeightedHelper[O]

	if len(success) == 0 {
		target = fail
	} else {
		target = success
	}

	if useMax {
		result = FilterByPositiveWeight(target)
	} else {
		result = FilterByNegativeWeight(target)
	}
	return result, len(success) > 0
}
