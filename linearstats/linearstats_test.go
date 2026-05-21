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

// ---------- Predict ----------

func TestPredict(t *testing.T) {
	cases := []struct {
		name    string
		window  []float64
		seen    int
		current float64
		wantLo  int64
		wantHi  int64
	}{
		{
			name:   "regression phase: single value centres on itself",
			window: []float64{250}, seen: 1, current: 250,
			wantLo: 230, wantHi: 270,
		},
		{
			name:   "regression phase: perfect upward trend extrapolates",
			window: []float64{100, 101, 102, 103, 104}, seen: 5, current: 104,
			wantLo: 85, wantHi: 125, // m=1,b=100 -> centre = 1*5+100 = 105
		},
		{
			name:   "regression phase: perfect downward trend extrapolates",
			window: []float64{104, 103, 102, 101, 100}, seen: 5, current: 100,
			wantLo: 79, wantHi: 119, // m=-1,b=104 -> centre = -1*5+104 = 99
		},
		{
			name:   "regression phase: constant window centres on the value",
			window: []float64{50, 50, 50, 50, 50}, seen: 10, current: 50,
			wantLo: 30, wantHi: 70,
		},
		{
			name:   "regression phase: just below the threshold",
			window: []float64{200, 200, 200, 200}, seen: regressionPhaseLimit - 1, current: 200,
			wantLo: 180, wantHi: 220,
		},
		{
			name:   "median phase: odd-length window centres on the median",
			window: []float64{10, 20, 30, 40, 50}, seen: regressionPhaseLimit, current: 50,
			wantLo: 10, wantHi: 50, // median = 30
		},
		{
			name:   "median phase: even-length window averages the two middles",
			window: []float64{10, 20, 30, 40}, seen: regressionPhaseLimit + 500, current: 40,
			wantLo: 5, wantHi: 45, // median = 25
		},
		{
			name:   "median phase: ignores an extreme outlier",
			window: []float64{100, 100, 100, 100, 9000}, seen: regressionPhaseLimit, current: 9000,
			wantLo: 80, wantHi: 120, // median = 100
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			lo, hi := Predict(tc.window, tc.seen, tc.current)
			if lo != tc.wantLo || hi != tc.wantHi {
				t.Errorf("Predict(%v, seen=%d) = %d %d; want %d %d",
					tc.window, tc.seen, lo, hi, tc.wantLo, tc.wantHi)
			}
		})
	}
}

func TestPredictFixedWidth(t *testing.T) {
	// Every prediction must span exactly 2*fixedRange, regardless of phase.
	wantWidth := int64(2 * fixedRange)
	cases := []struct {
		name   string
		window []float64
		seen   int
	}{
		{"regression phase", []float64{1, 5, 9, 13, 17}, 5},
		{"median phase", []float64{1, 5, 9, 13, 17}, regressionPhaseLimit},
		{"regression phase, large magnitude", []float64{9e12, 9e12, 9e12}, 3},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			lo, hi := Predict(tc.window, tc.seen, tc.window[len(tc.window)-1])
			if hi-lo != wantWidth {
				t.Errorf("width = %d; want %d", hi-lo, wantWidth)
			}
		})
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
