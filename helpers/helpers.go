// Package Helpers provides a set of helper functions and types that
// can be used for automatic error handling and result evaluation.
//
// However, this is still Work In Progress and is not yet fully
// implemented.
package helpers

// DoFunc is a generic type that represents a function that takes a value
// and does something with it.
//
// Parameters:
//   - T: The type of the value.
type DoFunc[T any] func(T)

// DualDoFunc is a generic type that represents a function that takes two
// values and does something with them.
//
// Parameters:
//   - T: The type of the first value.
//   - U: The type of the second value.
type DualDoFunc[T, U any] func(T, U)

// Helper is an interface that represents a helper.
type Helperer[O any] interface {
	// Data returns the data of the element.
	//
	// Returns:
	//   - O: The data of the element.
	//   - error: The reason for the failure.
	Data() (O, error)

	// Weight returns the weight of the element.
	//
	// Returns:
	//   - float64: The weight of the element.
	Weight() float64
}

// DoIfSuccess executes a function for each successful helper.
//
// Parameters:
//   - S: slice of helpers.
//   - f: the function to execute.
func DoIfSuccess[T Helperer[O], O any](S []T, f DoFunc[O]) {
	if len(S) == 0 || f == nil {
		return
	}

	for _, h := range S {
		data, err := h.Data()
		if err == nil {
			f(data)
		}
	}
}

// DoIfFailure executes a function for each failed helper.
//
// Parameters:
//   - S: slice of helpers.
//   - f: the function to execute.
func DoIfFailure[T Helperer[O], O any](S []T, f DualDoFunc[O, error]) {
	if len(S) == 0 || f == nil {
		return
	}

	for _, h := range S {
		data, err := h.Data()
		if err != nil {
			f(data, err)
		}
	}
}

// ExtractResults extracts the results from the helpers. Unlike with the Data
// method, this function returns only the results and not the pair of results and
// errors.
//
// Parameters:
//   - S: slice of helpers.
//
// Returns:
//   - []O: slice of results.
//
// Behaviors:
//   - The results are returned regardless of whether the helper is successful or not.
func ExtractResults[T Helperer[O], O any](S []T) []O {
	if len(S) == 0 {
		return nil
	}

	results := make([]O, 0, len(S))

	for _, h := range S {
		data, _ := h.Data()

		results = append(results, data)
	}

	return results
}

// SimpleHelper is a type that represents the result of a function evaluation
// that can either be successful or a failure.
type SimpleHelper[O any] struct {
	// result is the result of the function evaluation.
	result O

	// reason is the error that occurred during the function evaluation.
	reason error
}

// Data implements the Helperer interface.
func (h *SimpleHelper[O]) Data() (O, error) {
	return h.result, h.reason
}

// Weight implements the Helperer interface.
//
// Always returns 0.0.
func (h *SimpleHelper[O]) Weight() float64 {
	return 0.0
}

// NewSimpleHelper creates a new SimpleHelper with the given result and reason.
//
// Parameters:
//   - result: The result of the function evaluation.
//   - reason: The error that occurred during the function evaluation.
//
// Returns:
//   - SimpleHelper: The new SimpleHelper.
func NewSimpleHelper[O any](result O, reason error) *SimpleHelper[O] {
	sh := &SimpleHelper[O]{
		result: result,
		reason: reason,
	}
	return sh
}

// WeightedHelper is a generic type that represents the result of a function
// evaluation.
type WeightedHelper[O any] struct {
	// result is the result of the function evaluation.
	result O

	// reason is the error that occurred during the function evaluation.
	reason error

	// weight is the weight of the result (i.e., the probability of the result being correct)
	// or the most likely error (if the result is an error).
	weight float64
}

// Data implements the Helperer interface.
func (h *WeightedHelper[O]) Data() (O, error) {
	return h.result, h.reason
}

// Weight implements the Helperer interface.
func (h *WeightedHelper[O]) Weight() float64 {
	return h.weight
}

// NewWeightedHelper creates a new WeightedHelper with the given result, reason, and weight.
//
// Parameters:
//   - result: The result of the function evaluation.
//   - reason: The error that occurred during the function evaluation.
//   - weight: The weight of the result. The higher the weight, the more likely the result
//     is correct.
//
// Returns:
//   - WeightedHelper: The new WeightedHelper.
func NewWeightedHelper[O any](result O, reason error, weight float64) *WeightedHelper[O] {
	we := &WeightedHelper[O]{
		result: result,
		reason: reason,
		weight: weight,
	}
	return we
}
