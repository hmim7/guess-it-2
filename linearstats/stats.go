package linearstats

import (
	"math"
	"sort"
)

// Average calculates the mean of a slice of float64 values.
func Average(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	var sum float64
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

// Median returns the middle value of a sorted copy of data.
// For even-length slices it averages the two middle elements.
func Median(data []float64) float64 {
	n := len(data)
	if n == 0 {
		return 0
	}
	sorted := make([]float64, n)
	copy(sorted, data)
	sort.Float64s(sorted)
	if n%2 == 1 {
		return sorted[n/2]
	}
	return (sorted[n/2-1] + sorted[n/2]) / 2
}

// Variance computes the population variance using Welford's algorithm.
func Variance(data []float64) float64 {
	n := len(data)
	if n == 0 {
		return 0
	}
	var mean, m2 float64
	for i, v := range data {
		delta := v - mean
		mean += delta / float64(i+1)
		m2 += delta * (v - mean)
	}
	v := m2 / float64(n)
	if v < 0 {
		v = 0
	}
	return v
}

// StdDev returns the population standard deviation of the data.
func StdDev(data []float64) float64 {
	return math.Sqrt(Variance(data))
}

// LinearRegression returns slope m and intercept b for the line y = mx + b,
// where x = [0, 1, ..., n-1] and y = data values.
// Returns (0, 0) for empty slices and (0, data[0]) for single-element slices.
func LinearRegression(data []float64) (m, b float64) {
	n := len(data)
	if n == 0 {
		return 0, 0
	}
	if n == 1 {
		return 0, data[0]
	}
	fn := float64(n)
	var sumX, sumY, sumXY, sumX2 float64
	for i, y := range data {
		x := float64(i)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}
	denom := fn*sumX2 - sumX*sumX
	if denom == 0 {
		return 0, Average(data)
	}
	m = (fn*sumXY - sumX*sumY) / denom
	b = (sumY - m*sumX) / fn
	return m, b
}

// RunningOLS accumulates (x, y) pairs and yields the least-squares line in
// O(1) per call. Used by the stateful predictor so the trend is fit on the
// full prefix rather than a fixed 20-point window.
type RunningOLS struct {
	n                        int
	sumX, sumY, sumXY, sumXX float64
}

// Add includes one observation in the running sums.
func (o *RunningOLS) Add(x, y float64) {
	o.n++
	o.sumX += x
	o.sumY += y
	o.sumXY += x * y
	o.sumXX += x * x
}

// N returns the number of observations seen so far.
func (o *RunningOLS) N() int { return o.n }

// Fit returns slope and intercept of the least-squares line, plus ok=false if
// the fit is degenerate (fewer than 2 points or a constant x).
func (o *RunningOLS) Fit() (m, b float64, ok bool) {
	if o.n < 2 {
		return 0, 0, false
	}
	fn := float64(o.n)
	denom := fn*o.sumXX - o.sumX*o.sumX
	if denom == 0 {
		return 0, o.sumY / fn, false
	}
	m = (fn*o.sumXY - o.sumX*o.sumY) / denom
	b = (o.sumY - m*o.sumX) / fn
	return m, b, true
}

// PearsonCorrelation returns r, the Pearson correlation coefficient,
// where x = [0, 1, ..., n-1] and y = data values.
// Returns 0 for edge cases: n<=1, constant y, or overflow in the denominator.
func PearsonCorrelation(data []float64) float64 {
	n := len(data)
	if n <= 1 {
		return 0
	}
	fn := float64(n)
	var sumX, sumY, sumXY, sumX2, sumY2 float64
	for i, y := range data {
		x := float64(i)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
		sumY2 += y * y
	}
	numerator := fn*sumXY - sumX*sumY
	denomX := fn*sumX2 - sumX*sumX
	denomY := fn*sumY2 - sumY*sumY
	denom := math.Sqrt(denomX * denomY)
	if denom == 0 || math.IsNaN(denom) || math.IsInf(denom, 0) {
		return 0
	}
	r := numerator / denom
	if math.IsNaN(r) || math.IsInf(r, 0) {
		return 0
	}
	return r
}
