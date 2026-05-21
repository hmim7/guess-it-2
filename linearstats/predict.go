package linearstats

import (
	"math"
)

const (
	// WindowSize is the number of most-recent values kept for each prediction.
	WindowSize = 20
	// fixedRange is the half-width of every predicted range. A constant narrow
	// width maximises score-per-hit: under the audit's "smaller range scores
	// higher" rule, a tight band beats a wide one even at a lower hit rate.
	fixedRange = 20.0
	// regressionPhaseLimit is the input count below which the centre is the
	// linear-regression extrapolation; at or above it the centre is the window
	// median, which the dataset analysis showed scores higher in steady state.
	regressionPhaseLimit = 1000
	// minWidth is the smallest allowed gap between the lower and upper bounds.
	minWidth = 1
)

// Predict returns the lower and upper bounds for the next value in the stream.
// While fewer than regressionPhaseLimit values have been seen, the range is
// centred on the linear-regression extrapolation of the window; afterwards it
// is centred on the window median. The range half-width is always fixedRange.
func Predict(window []float64, seen int, current float64) (int64, int64) {
	var center float64
	if seen < regressionPhaseLimit {
		m, b := LinearRegression(window)
		center = m*float64(len(window)) + b
	} else {
		center = Median(window)
	}
	return roundBounds(center-fixedRange, center+fixedRange)
}

// roundBounds converts calculated floating-point boundaries into nearest integers.
func roundBounds(lower, upper float64) (int64, int64) {
	l := int64(math.Round(lower))
	u := int64(math.Round(upper))
	if l > u {
		l, u = u, l
	}
	if u-l < int64(minWidth) {
		u = l + int64(minWidth)
	}
	return l, u
}
