package linearstats

import (
	"math"
)

const (
	WindowSize = 20
	// multiplierC is dynamic, see Predict()
	initialRange = 60.0
	minStdRange  = 1.8 // Lowered for more aggressive scoring in stable periods
	minWidth     = 1

	// Dynamic tuning parameters
	veryLowVarThreshold = 5.0  // Threshold for extremely stable data
	veryLowMultiplier   = 1.2  // High risk, massive score reward
	lowVarThreshold     = 25.0 // Threshold for stable data
	highVarThreshold    = 40.0 // Threshold for moderate volatility
	aggroMultiplier     = 1.4  // Tight bounds for high scoring
	balancedMultiplier  = 1.8  // Balanced risk/reward
	defensiveMultiplier = 2.1  // Widen bounds for safety
	extremeVarThreshold = 80.0 // Threshold for extremely volatile periods
	extremeMultiplier   = 2.4  // Maximum safety to prevent misses
)

// Predict calculates expected lower and upper bounds for the next number using a sliding window.
// It applies a dynamic multiplier based on data standard deviation to balance hit rate and score.
func Predict(window []float64, seen int, current float64, useMedian bool) (int64, int64) {
	// Initial phase: safe wide range for the first few inputs.
	// Scale by the value's magnitude so streams that open in the thousands
	// aren't capped by a ±60 net that misses on the very first samples.
	if seen < 5 {
		r := initialRange
		if scaled := math.Abs(current) * 0.5; scaled > r {
			r = scaled
		}
		return roundBounds(current-r, current+r)
	}

	m, b := LinearRegression(window)
	r := PearsonCorrelation(window)
	w := math.Abs(r)

	regC := m*float64(len(window)) + b
	robC := Average(window)
	if useMedian {
		robC = Median(window)
	}
	center := w*regC + (1-w)*robC
	std := StdDev(window)

	// Use a dynamic multiplier based on data volatility (standard deviation).
	// This allows us to be aggressive in stable periods and defensive in volatile ones,
	// improving the score without sacrificing too much hit rate.
	var dynamicMultiplier float64
	if std < veryLowVarThreshold {
		// Extremely stable data: take a big risk for a massive score.
		dynamicMultiplier = veryLowMultiplier
	} else if std < lowVarThreshold {
		// Low variance: tighten the range for a higher score.
		dynamicMultiplier = aggroMultiplier
	} else if std < highVarThreshold {
		// Medium variance: a balanced approach.
		dynamicMultiplier = balancedMultiplier
	} else if std < extremeVarThreshold {
		// High variance: a defensive approach to ensure a hit.
		dynamicMultiplier = defensiveMultiplier
	} else {
		// Extreme variance: play it very safe to avoid a miss on huge spikes.
		dynamicMultiplier = extremeMultiplier
	}

	margin := dynamicMultiplier * std * (1 - 0.25*w)
	if margin < minStdRange {
		margin = minStdRange
	}
	// Recent-jump guard: if the latest value is a genuine outlier (>1.5σ from center), the stream may be shifting.
	// Stretch the margin to cover the jump plus a small buffer so the next prediction doesn't miss.
	if delta := math.Abs(current - center); delta > 1.5*std {
		if jump := delta * 1.1; jump > margin {
			margin = jump
		}
	}
	return roundBounds(center-margin, center+margin)
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
