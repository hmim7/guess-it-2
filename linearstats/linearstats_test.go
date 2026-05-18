package linearstats

import (
	"math"
	"testing"
)

// ---------- Stats ----------

func TestAverage(t *testing.T) {
	cases := []struct {
		name string
		data []float64
		want float64
	}{
		{"simple", []float64{10, 20, 30, 40, 50}, 30.0},
		{"empty", []float64{}, 0},
		{"single", []float64{5}, 5},
		{"negatives", []float64{-10, 0, 10, 20}, 5},
		{"all equal", []float64{7, 7, 7, 7}, 7},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Average(tc.data); got != tc.want {
				t.Errorf("Average(%v) = %f; want %f", tc.data, got, tc.want)
			}
		})
	}
}

func TestMedian(t *testing.T) {
	cases := []struct {
		name string
		data []float64
		want float64
	}{
		{"odd length", []float64{3, 1, 4, 1, 5}, 3},
		{"even length", []float64{8, 4, 1, 6}, 5},
		{"empty", []float64{}, 0},
		{"single", []float64{42}, 42},
		{"negatives", []float64{-10, 20, -30, 0}, -5},
		{"duplicates", []float64{5, 2, 5, 1, 5, 3}, 4},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Median(tc.data); got != tc.want {
				t.Errorf("Median(%v) = %f; want %f", tc.data, got, tc.want)
			}
		})
	}
}

func TestVariance(t *testing.T) {
	cases := []struct {
		name string
		data []float64
		want float64
	}{
		{"standard", []float64{1, 2, 3, 4, 5}, 2.0},
		{"empty", []float64{}, 0},
		{"single", []float64{100}, 0},
		{"all equal", []float64{8, 8, 8, 8}, 0},
		{"three values", []float64{1, 2, 3}, 2.0 / 3.0},
		{"negatives", []float64{-2, -1, 0, 1, 2}, 2.0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Variance(tc.data); math.Abs(got-tc.want) > 1e-9 {
				t.Errorf("Variance(%v) = %f; want %f", tc.data, got, tc.want)
			}
		})
	}
}

func TestStdDev(t *testing.T) {
	cases := []struct {
		name string
		data []float64
		want float64
	}{
		{"standard", []float64{1, 2, 3, 4, 5}, math.Sqrt(2.0)},
		{"empty", []float64{}, 0},
		{"single", []float64{100}, 0},
		{"all equal", []float64{8, 8, 8, 8}, 0},
		{"three values", []float64{1, 2, 3}, math.Sqrt(2.0 / 3.0)},
		{"negatives", []float64{-2, -1, 0, 1, 2}, math.Sqrt(2.0)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := StdDev(tc.data); math.Abs(got-tc.want) > 1e-9 {
				t.Errorf("StdDev(%v) = %f; want %f", tc.data, got, tc.want)
			}
		})
	}
}

func TestLinearRegression(t *testing.T) {
	cases := []struct {
		name      string
		data      []float64
		wantM     float64
		wantB     float64
		tolerance float64
	}{
		{name: "empty", data: []float64{}, wantM: 0, wantB: 0, tolerance: 1e-9},
		{name: "single value", data: []float64{42}, wantM: 0, wantB: 42, tolerance: 1e-9},
		{name: "perfect positive slope", data: []float64{1, 2, 3, 4, 5}, wantM: 1.0, wantB: 1.0, tolerance: 1e-9},
		{name: "perfect negative slope", data: []float64{5, 4, 3, 2, 1}, wantM: -1.0, wantB: 5.0, tolerance: 1e-9},
		{name: "constant y", data: []float64{10, 10, 10, 10}, wantM: 0.0, wantB: 10.0, tolerance: 1e-9},
		{name: "two values", data: []float64{10, 20}, wantM: 10.0, wantB: 10.0, tolerance: 1e-9},
		{
			name: "standard dataset", data: []float64{189, 113, 121, 114, 145, 110},
			wantM: -8.742857142857143, wantB: 153.85714285714286, tolerance: 1e-6,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m, b := LinearRegression(tc.data)
			if math.Abs(m-tc.wantM) > tc.tolerance {
				t.Errorf("slope: got %.10f; want %.10f", m, tc.wantM)
			}
			if math.Abs(b-tc.wantB) > tc.tolerance {
				t.Errorf("intercept: got %.10f; want %.10f", b, tc.wantB)
			}
		})
	}
}

func TestPearsonCorrelation(t *testing.T) {
	cases := []struct {
		name      string
		data      []float64
		wantR     float64
		tolerance float64
	}{
		{name: "empty", data: []float64{}, wantR: 0, tolerance: 1e-9},
		{name: "single value", data: []float64{42}, wantR: 0, tolerance: 1e-9},
		{name: "perfect positive correlation", data: []float64{1, 2, 3, 4, 5}, wantR: 1.0, tolerance: 1e-9},
		{name: "perfect negative correlation", data: []float64{5, 4, 3, 2, 1}, wantR: -1.0, tolerance: 1e-9},
		{name: "constant y (undefined → 0)", data: []float64{10, 10, 10, 10}, wantR: 0, tolerance: 1e-9},
		{name: "two values same", data: []float64{5, 5}, wantR: 0, tolerance: 1e-9},
		{name: "two values different", data: []float64{10, 20}, wantR: 1.0, tolerance: 1e-9},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := PearsonCorrelation(tc.data)
			if math.Abs(r-tc.wantR) > tc.tolerance {
				t.Errorf("got %.10f; want %.10f", r, tc.wantR)
			}
		})
	}
}

// ---------- Predict / roundBounds ----------

// applyInput simulates one iteration of the streaming loop.
func applyInput(window []float64, seen int, v float64, useMedian bool) ([]float64, int, int64, int64) {
	seen++
	if len(window) < WindowSize {
		window = append(window, v)
	} else {
		window = append(window[1:], v)
	}
	l, u := Predict(window, seen, v, useMedian)
	return window, seen, l, u
}

func TestPredictInitialPhase(t *testing.T) {
	window := []float64{250}
	current := 250.0

	l, u := Predict(window, 1, current, false)

	r := initialRange
	if scaled := current * 0.5; scaled > r {
		r = scaled
	}
	wantL, wantU := roundBounds(current-r, current+r)
	if l != wantL || u != wantU {
		t.Fatalf("initial phase: got %d %d; want %d %d", l, u, wantL, wantU)
	}
}

func TestPredictZeroStdDev(t *testing.T) {
	window := []float64{42, 42, 42, 42, 42}
	l, u := Predict(window, WindowSize, 42, false)
	if l != 40 || u != 44 {
		t.Fatalf("zero stddev: got %d %d; want 40 44", l, u)
	}
}

func TestPredictMedianCenter(t *testing.T) {
	// Window has Median (100) != Average (80.2). The hybrid predictor blends
	// the chosen robust center with regression extrapolation via |r|, so we
	// assert that toggling useMedian changes the output (the robust center
	// path is exercised) rather than hard-coding numerical bounds.
	window := []float64{1, 100, 100, 100, 100}
	current := 100.0

	if Median(window) == Average(window) {
		t.Fatal("test setup: median must differ from average")
	}

	lMed, uMed := Predict(window, len(window), current, true)
	lMean, uMean := Predict(window, len(window), current, false)
	if lMed == lMean && uMed == uMean {
		t.Errorf("median path not exercised: median=%d,%d mean=%d,%d", lMed, uMed, lMean, uMean)
	}
}

func TestGoldenInitialWindowFill(t *testing.T) {
	steps := []struct {
		name  string
		input float64
		lower int64
		upper int64
	}{
		{name: "Input 100", input: 100, lower: 40, upper: 160},
		{name: "Input 102", input: 102, lower: 42, upper: 162},
		{name: "Input 104", input: 104, lower: 44, upper: 164},
		{name: "Input 101", input: 101, lower: 41, upper: 161},
		{name: "Input 103", input: 103, lower: 101, upper: 105},
	}
	window := make([]float64, 0, WindowSize)
	seen := 0
	for _, s := range steps {
		t.Run(s.name, func(t *testing.T) {
			var l, u int64
			window, seen, l, u = applyInput(window, seen, s.input, false)
			if l != s.lower || u != s.upper {
				t.Errorf("got %d %d; want %d %d", l, u, s.lower, s.upper)
			}
		})
	}
}

func TestGoldenStableSlidingWindow(t *testing.T) {
	window := make([]float64, 0, WindowSize)
	seen := 0
	for _, v := range []float64{100, 102, 104, 101, 103} {
		window, seen, _, _ = applyInput(window, seen, v, false)
	}

	steps := []struct {
		name  string
		input float64
		lower int64
		upper int64
	}{
		{name: "Step 6: Input 105", input: 105, lower: 102, upper: 106},
		{name: "Step 7: Input 106", input: 106, lower: 104, upper: 108},
	}
	for _, s := range steps {
		t.Run(s.name, func(t *testing.T) {
			var l, u int64
			window, seen, l, u = applyInput(window, seen, s.input, false)
			if l != s.lower || u != s.upper {
				t.Errorf("got %d %d; want %d %d", l, u, s.lower, s.upper)
			}
		})
	}
}

func TestGoldenZeroStdDevScenario(t *testing.T) {
	window := []float64{50, 50, 50, 50, 50}
	l, u := Predict(window, WindowSize, 50, false)
	if l != 48 || u != 52 {
		t.Fatalf("GT03: got %d %d; want 48 52", l, u)
	}
}

func TestGoldenSuddenSpike(t *testing.T) {
	window := make([]float64, 0, WindowSize)
	seen := 0
	for _, v := range []float64{100, 101, 102, 103, 104} {
		window, seen, _, _ = applyInput(window, seen, v, false)
	}

	t.Run("Step 1: Predict on stable state", func(t *testing.T) {
		l, u := Predict(window, seen, 104, false)
		if l != 103 || u != 107 {
			t.Errorf("got %d %d; want 103 107", l, u)
		}
	})

	steps := []struct {
		name  string
		input float64
		lower int64
		upper int64
	}{
		{name: "Step 2: Input 500 (Sudden Spike)", input: 500, lower: 4, upper: 598},
		{name: "Step 3: Input 105 (Recovery)", input: 105, lower: -91, upper: 507},
	}
	for _, s := range steps {
		t.Run(s.name, func(t *testing.T) {
			var l, u int64
			window, seen, l, u = applyInput(window, seen, s.input, false)
			if l != s.lower || u != s.upper {
				t.Errorf("got %d %d; want %d %d", l, u, s.lower, s.upper)
			}
		})
	}
}

func TestEdgeOscillatingValues(t *testing.T) {
	window := make([]float64, 0, WindowSize)
	seen := 0
	for _, v := range []float64{10, 100, 10, 100, 10} {
		window, seen, _, _ = applyInput(window, seen, v, false)
	}

	steps := []struct {
		name  string
		input float64
		lower int64
		upper int64
	}{
		{name: "Predict after multiple oscillations", input: 100, lower: -25, upper: 150},
	}
	for _, s := range steps {
		t.Run(s.name, func(t *testing.T) {
			var l, u int64
			window, seen, l, u = applyInput(window, seen, s.input, false)
			if l > 10 || u < 100 {
				t.Errorf("bounds %d %d do not encompass oscillations 10 and 100", l, u)
			}
			if l != s.lower || u != s.upper {
				t.Errorf("got %d %d; want %d %d", l, u, s.lower, s.upper)
			}
		})
	}
}

func TestDynamicMultiplier(t *testing.T) {
	cases := []struct {
		name string
		std  float64
		want float64
	}{
		{"very low", 1.0, veryLowMultiplier},
		{"low (boundary)", veryLowVarThreshold, aggroMultiplier},
		{"low", 10.0, aggroMultiplier},
		{"balanced (boundary)", lowVarThreshold, balancedMultiplier},
		{"balanced", 30.0, balancedMultiplier},
		{"defensive (boundary)", highVarThreshold, defensiveMultiplier},
		{"defensive", 50.0, defensiveMultiplier},
		{"extreme (boundary)", extremeVarThreshold, extremeMultiplier},
		{"extreme", 200.0, extremeMultiplier},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := dynamicMultiplier(tc.std); got != tc.want {
				t.Errorf("dynamicMultiplier(%v)=%v; want %v", tc.std, got, tc.want)
			}
		})
	}
}

// fillWindow streams a slice through applyInput, returning the final state.
func fillWindow(t *testing.T, values []float64, useMedian bool) ([]float64, int) {
	t.Helper()
	window := make([]float64, 0, WindowSize)
	seen := 0
	for _, v := range values {
		window, seen, _, _ = applyInput(window, seen, v, useMedian)
	}
	return window, seen
}

func TestPredictLowRFallback(t *testing.T) {
	// Oscillating window → PearsonCorrelation ≈ 0, so the hybrid blend
	// collapses to the robust center and the multiplier scaling factor
	// (1 - 0.25*|r|) is ~1. The range must encompass both extremes.
	values := []float64{10, 100, 10, 100, 10, 100, 10, 100, 10, 100}
	window, _ := fillWindow(t, values, true)

	r := PearsonCorrelation(window)
	if math.Abs(r) > 0.2 {
		t.Fatalf("setup: oscillating window should have |r| ≈ 0, got %v", r)
	}

	l, u := Predict(window, len(window), 100, true)
	if int64(10) < l || int64(100) > u {
		t.Errorf("low-|r| fallback: bounds [%d,%d] must encompass 10 and 100", l, u)
	}
}

func TestPredictHighRTightening(t *testing.T) {
	// Clean linear trend → |r| ≈ 1. The hybrid centre tracks the regression
	// extrapolation and the multiplier is scaled by (1 - 0.25*|r|) ≈ 0.75,
	// producing a tighter band than the same window with |r| forced to 0.
	values := []float64{100, 101, 102, 103, 104, 105, 106, 107, 108, 109}
	window, _ := fillWindow(t, values, true)

	r := PearsonCorrelation(window)
	if r < 0.99 {
		t.Fatalf("setup: trended window should have |r| ≈ 1, got %v", r)
	}

	l, u := Predict(window, len(window), 110, true)
	width := u - l
	if width > 6 {
		t.Errorf("high-|r| tightening: band width %d unexpectedly wide (expected ≤ 6)", width)
	}
	if int64(110) < l || int64(110) > u {
		t.Errorf("high-|r| tightening: bounds [%d,%d] should encompass next value 110", l, u)
	}
}

func TestPredictMidRBlend(t *testing.T) {
	// Trended-but-noisy data → 0.4 < |r| < 0.95. The blended centre should
	// sit strictly between the pure median and the pure regression
	// extrapolation, exercising the blend path.
	values := []float64{100, 102, 99, 105, 103, 108, 106, 110, 108, 112}
	window, _ := fillWindow(t, values, true)

	r := PearsonCorrelation(window)
	if math.Abs(r) < 0.4 || math.Abs(r) > 0.95 {
		t.Fatalf("setup: mid-correlation window should have 0.4 < |r| < 0.95, got %v", r)
	}

	l, u := Predict(window, len(window), 113, true)
	if u-l < int64(minWidth) {
		t.Errorf("mid-|r| blend: bounds [%d,%d] violate minWidth", l, u)
	}
}

func TestPredictConstantWindowZeroR(t *testing.T) {
	// All-equal window: PearsonCorrelation returns 0 (constant-y guard),
	// StdDev is 0, margin floored to minStdRange. Centre = constant value.
	values := []float64{80, 80, 80, 80, 80}
	window, _ := fillWindow(t, values, true)

	if r := PearsonCorrelation(window); r != 0 {
		t.Fatalf("setup: constant window should have r = 0, got %v", r)
	}
	if std := StdDev(window); std != 0 {
		t.Fatalf("setup: constant window should have std = 0, got %v", std)
	}

	l, u := Predict(window, len(window), 80, true)
	if l != 78 || u != 82 {
		t.Errorf("constant-window-r=0: got %d %d; want 78 82", l, u)
	}
}

func TestRoundBounds(t *testing.T) {
	cases := []struct {
		name  string
		inL   float64
		inU   float64
		wantL int64
		wantU int64
	}{
		{"standard rounding", 10.2, 20.8, 10, 21},
		{"swapped bounds", 20.8, 10.2, 10, 21},
		{"min width enforcement (equal bounds)", 15.1, 15.4, 15, 16},
		{"min width enforcement (zero width)", 15.0, 15.0, 15, 16},
		{"negative numbers", -10.6, -5.2, -11, -5},
		{"mixed positive and negative", -2.4, 2.4, -2, 2},
		{"bounds already have min width", 10.0, 11.0, 10, 11},
		{"extremely large float values", 9999999999999.2, 9999999999999.8, 9999999999999, 10000000000000},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l, u := roundBounds(tc.inL, tc.inU)
			if l != tc.wantL || u != tc.wantU {
				t.Errorf("roundBounds(%.2f,%.2f)=(%d,%d); want (%d,%d)",
					tc.inL, tc.inU, l, u, tc.wantL, tc.wantU)
			}
		})
	}
}
