package MathExt

import (
	"math/big"
	"strings"

	gcers "github.com/PlayerR9/go-commons/errors"
	gcint "github.com/PlayerR9/go-commons/ints"
)

// Serieser is an interface for series.
type Serieser interface {
	// Term returns the nth term of the series.
	//
	// Parameters:
	//   - n: The term number.
	//
	// Returns:
	//   - *big.Int: The nth term of the series.
	//   - error: An error if the term cannot be calculated.
	Term(n int) (*big.Int, error)
}

// ApproximateConvergence approximates the convergence of a series.
// It calculates the average of the last n values.
//
// Parameters:
//   - n: The number of values to calculate the average.
//
// Returns:
//   - *big.Float: The average of the last n values.
//   - error: An error if the calculation fails.
//
// Errors:
//   - *luc.ErrInvalidParameter: If n is less than or equal to 0 or if
//     there are not enough values to calculate the average.
func ApproximateConvergence(values []*big.Float, n int) (*big.Float, error) {
	if n <= 0 {
		return nil, gcers.NewErrInvalidParameter("n", gcint.NewErrGT(0))
	} else if len(values) < n {
		return nil, gcers.NewErrInvalidParameter(
			"n",
			gcint.NewErrOutOfBounds(n, 0, len(values)),
		)
	}

	sum := new(big.Float)
	for i := len(values) - n; i < len(values); i++ {
		sum.Add(sum, values[i])
	}

	average := new(big.Float).Quo(sum, new(big.Float).SetFloat64(float64(n)))
	return average, nil
}

// CalculateConvergence calculates the convergence of a series.
// It calculates the quotient of the ith term and the (i+delta)th term.
//
// Parameters:
//   - series: The series to calculate the convergence.
//   - upperLimit: The upper limit of the series to calculate the convergence.
//   - delta: The difference between the terms to calculate the convergence.
//
// Returns:
//   - *ConvergenceResult: The convergence result.
//   - error: An error if the calculation fails.
func CalculateConvergence(series Serieser, upperLimit int, delta int) (values []*big.Float, err error) {
	if series == nil {
		return nil, gcers.NewErrNilParameter("series")
	}

	for i := 0; i < upperLimit-delta; i++ {
		ithTerm, reason := series.Term(i)
		if reason != nil {
			err = gcint.NewErrAt(i+1, "term", reason)
			return
		}

		ithPlusDeltaTerm, reason := series.Term(i + delta)
		if reason != nil {
			err = gcint.NewErrAt(i+delta+1, "term", reason)
			return
		}

		quotient := new(big.Float).Quo(
			new(big.Float).SetInt(ithPlusDeltaTerm),
			new(big.Float).SetInt(ithTerm),
		)

		values = append(values, quotient)
	}

	return
}

// LinearRegression is a struct that holds the equation of a linear regression.
type LinearRegression struct {
	// A is the coefficient of the linear regression.
	A *big.Float

	// B is the exponent of the linear regression.
	B *big.Float
}

// String implements the fmt.Stringer interface.
//
// Format: y = a * x^b
func (lr *LinearRegression) String() string {
	values := []string{
		"y =",
		lr.A.String(),
		"* x ^",
		lr.B.String(),
	}

	str := strings.Join(values, " ")

	return str
}

// NewLinearRegression creates a new LinearRegression.
//
// Returns:
//   - LinearRegression: The new LinearRegression.
func NewLinearRegression() *LinearRegression {
	lr := &LinearRegression{
		A: new(big.Float).SetPrec(1000),
		B: new(big.Float).SetPrec(1000),
	}

	return lr
}

// FindEquation is a method of ConvergenceResult that finds the equation of the series
// that best fits the convergence values.
//
// The equation is of the form y = a * x^b.
//
// Returns:
//   - bool: False if there are less than 2 values to calculate the equation.
//     True otherwise.
func (l *LinearRegression) FindEquation(values []*big.Float) bool {
	if len(values) < 2 {
		return false
	}

	sumX := new(big.Float)
	sumY := new(big.Float)
	sumXY := new(big.Float)
	sumX2 := new(big.Float)
	n := big.NewFloat(float64(len(values)))

	for i, v := range values {
		x := big.NewFloat(float64(i))
		y := new(big.Float).Set(v)

		sumX.Add(sumX, x)
		sumY.Add(sumY, y)
		sumXY.Add(sumXY, new(big.Float).Mul(x, y))
		sumX2.Add(sumX2, new(big.Float).Mul(x, x))
	}

	// a = (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	l.A = new(big.Float).Sub(
		new(big.Float).Mul(n, sumXY),
		new(big.Float).Mul(sumX, sumY),
	)
	l.A = l.A.Quo(l.A, new(big.Float).Sub(
		new(big.Float).Mul(n, sumX2),
		new(big.Float).Mul(sumX, sumX),
	))
	l.A = new(big.Float).SetPrec(1000).Set(l.A)

	// b = (sumY - a*sumX) / n
	l.B = new(big.Float).Sub(
		sumY,
		new(big.Float).Mul(l.A, sumX),
	)
	l.B = l.B.Quo(l.B, n)
	l.B = new(big.Float).SetPrec(1000).Set(l.B)

	return true
}
