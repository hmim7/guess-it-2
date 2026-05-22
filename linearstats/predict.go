package linearstats

import (
	"math"
)

const (
	WindowSize           = 20   // most-recent values kept per prediction
	fixedRange           = 20.0 // half-width of predicted range (narrow = higher score)
	regressionPhaseLimit = 1000 // below: regression centre; at/above: median centre
	minWidth             = 1    // minimum gap between lower and upper bounds
)

// Predict returns [lower, upper] for the next value: regression-centred below
// regressionPhaseLimit inputs, median-centred above; half-width is fixedRange.
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
