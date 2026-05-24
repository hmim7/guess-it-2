package linearstats

import (
	"math"
)

const (
	WindowSize           = 20   // most-recent values kept per prediction
	fixedRange           = 46.0 // half-width of predicted range (narrow = higher score)
	regressionPhaseLimit = 1000 // below: regression centre; at/above: median centre
	minWidth             = 1    // minimum gap between lower and upper bounds

	// olsMinSamples is the prefix size at which the running OLS fit is trusted.
	// Below this the predictor falls back to the original window-based centre.
	olsMinSamples = 30
)

// Predictor uses a full-prefix running OLS for stable, unbiased trend extrapolation.
// It maintains the constant ±fixedRange half-width from the baseline model.
type Predictor struct {
	window []float64
	ols    RunningOLS
	seen   int
}

// NewPredictor builds a fresh predictor.
func NewPredictor() *Predictor {
	return &Predictor{window: make([]float64, 0, WindowSize)}
}

// Next consumes one sample y_n and returns the predicted [lower, upper] bounds for y_{n+1}.
func (p *Predictor) Next(current float64) (int64, int64) {
	if len(p.window) < WindowSize {
		p.window = append(p.window, current)
	} else {
		p.window = append(p.window[1:], current)
	}
	p.ols.Add(float64(p.seen), current)
	p.seen++

	if p.seen >= olsMinSamples {
		if m, b, ok := p.ols.Fit(); ok {
			center := m*float64(p.seen) + b
			return roundBounds(center-fixedRange, center+fixedRange)
		}
	}
	return Predict(p.window, p.seen, current)
}

// Predict implements the legacy split-prediction logic (window regression vs. median).
// It serves as a cold-start fallback for the stateful Predictor.
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
