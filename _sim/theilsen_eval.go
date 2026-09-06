package main

// Evaluates the proposed "Theil-Sen + empirical-error interval" hybrid
// predictor (.resources/Theil_Sen_estimator.md). Window-30 Theil-Sen point
// prediction, then an O(M^2) search over all pairs of the last 100 signed
// prediction errors, picking the [lowErr, highErr] band that maximises
// count/(1+width). Scored with the exact server.js formula and compared
// head-to-head against the recorded current-predictor and audit-opponent
// scores.
//
//   go run _sim/theilsen_eval.go .resources/guess-it-dockerized/data_sets

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func readData(path string) []float64 {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var d []float64
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		if v, e := strconv.ParseFloat(line, 64); e == nil {
			d = append(d, v)
		}
	}
	return d
}

func medianSorted(sorted []float64) float64 {
	n := len(sorted)
	if n == 0 {
		return 0
	}
	if n%2 == 1 {
		return sorted[n/2]
	}
	return (sorted[n/2-1] + sorted[n/2]) / 2.0
}

func mad(sorted []float64, med float64) float64 {
	dev := make([]float64, len(sorted))
	for i, v := range sorted {
		dev[i] = math.Abs(v - med)
	}
	sort.Float64s(dev)
	return medianSorted(dev)
}

// theilSen — verbatim from .resources/Theil_Sen_estimator.md: collect all
// pairwise slopes (y_j - y_i)/(j - i), take their median; then intercept =
// median(y_i - slope*i); extrapolate to x = len(window).
func theilSen(window []float64) float64 {
	n := len(window)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return window[0]
	}
	slopes := make([]float64, 0, n*(n-1)/2)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			slopes = append(slopes, (window[j]-window[i])/float64(j-i))
		}
	}
	sort.Float64s(slopes)
	slope := medianSorted(slopes)
	intercepts := make([]float64, n)
	for i, y := range window {
		intercepts[i] = y - slope*float64(i)
	}
	sort.Float64s(intercepts)
	intercept := medianSorted(intercepts)
	return intercept + slope*float64(n)
}

// scoreHybrid runs the .resources hybrid (Theil-Sen centre + empirical
// interval search) and returns server.js score, hits, mean width.
func scoreHybrid(numbers []float64) (score, hits int, meanWidth float64) {
	n := len(numbers)
	if n < 2 {
		return 0, 0, 0
	}
	const windowSize = 30
	const errorBuffer = 100
	recentErrors := []float64{}
	widthSum := 0
	for i := 0; i < n-1; i++ {
		start := i + 1 - windowSize
		if start < 0 {
			start = 0
		}
		window := numbers[start : i+1]
		pred := theilSen(window)

		var low, high int
		if len(recentErrors) < 5 {
			half := 40.0
			if len(recentErrors) >= 2 {
				sorted := make([]float64, len(recentErrors))
				copy(sorted, recentErrors)
				sort.Float64s(sorted)
				med := medianSorted(sorted)
				half = 1.5 * mad(sorted, med)
				if half < 20 {
					half = 20
				}
			}
			low = int(math.Floor(pred - half))
			high = int(math.Ceil(pred + half))
		} else {
			sorted := make([]float64, len(recentErrors))
			copy(sorted, recentErrors)
			sort.Float64s(sorted)
			bestScore := -1.0
			bestLowErr, bestHighErr := 0.0, 0.0
			for a := 0; a < len(sorted); a++ {
				for b := a; b < len(sorted); b++ {
					lowErr := sorted[a]
					highErr := sorted[b]
					count := b - a + 1
					l := int(math.Floor(pred + lowErr))
					h := int(math.Ceil(pred + highErr))
					width := h - l
					if width < 0 {
						width = 0
					}
					sc := float64(count) / float64(1+width)
					if sc > bestScore ||
						(sc == bestScore && (highErr-lowErr) < (bestHighErr-bestLowErr)) {
						bestScore = sc
						bestLowErr = lowErr
						bestHighErr = highErr
					}
				}
			}
			low = int(math.Floor(pred + bestLowErr))
			high = int(math.Ceil(pred + bestHighErr))
		}
		if low > high {
			low, high = high, low
		}

		actual := numbers[i+1]
		w := high - low
		widthSum += w
		if actual >= float64(low) && actual <= float64(high) {
			score += int(math.Round(1e7 / float64(1+w) / float64(n-1)))
			hits++
		}
		recentErrors = append(recentErrors, actual-pred)
		if len(recentErrors) > errorBuffer {
			recentErrors = recentErrors[len(recentErrors)-errorBuffer:]
		}
	}
	return score, hits, float64(widthSum) / float64(n-1)
}

// Also evaluate the Theil-Sen-only ablation (robust centre + symmetric ±20
// fixed width) — isolates whether the centre change alone helps.
func scoreTheilSenFixed(numbers []float64) int {
	n := len(numbers)
	if n < 2 {
		return 0
	}
	const windowSize = 30
	const half = 20.0
	score := 0
	for i := 0; i < n-1; i++ {
		start := i + 1 - windowSize
		if start < 0 {
			start = 0
		}
		window := numbers[start : i+1]
		pred := theilSen(window)
		low := int(math.Round(pred - half))
		high := int(math.Round(pred + half))
		if high-low < 1 {
			high = low + 1
		}
		actual := numbers[i+1]
		w := high - low
		if actual >= float64(low) && actual <= float64(high) {
			score += int(math.Round(1e7 / float64(1+w) / float64(n-1)))
		}
	}
	return score
}

func main() {
	base := os.Args[1]

	// Recorded per-file scores (WSL, exact server.js formula) — same source as
	// _sim/adaptive_eval.go. "current" = linear-v2 running-OLS predictor.
	opp := map[string][2][5]int{
		"current":          {{101900, 101640, 104220, 101820, 102580}, {102140, 98840, 101960, 100840, 99820}},
		"linear-regr":      {{102300, 100750, 104532, 101680, 101618}, {100006, 99696, 99448, 97216, 101370}},
		"mse":              {{100392, 101727, 97989, 100659, 104130}, {92115, 107067, 96921, 94785, 100392}},
		"nic":              {{111200, 98400, 106400, 88800, 97600}, {93600, 88000, 108000, 102400, 98400}},
		"big-range":        {{49992, 49996, 49996, 49996, 49996}, {49140, 49168, 49280, 49152, 49172}},
		"correlation-coef": {{95076, 94392, 93442, 92226, 93670}, {89566, 86830, 89718, 89376, 88654}},
	}

	var hyb, tsf [2][5]int
	var hitPct, widths [2][5]float64
	for ds := 0; ds < 2; ds++ {
		for f := 0; f < 5; f++ {
			d := readData(fmt.Sprintf("%s/%d/%d.txt", base, ds+4, f+1))
			s, h, w := scoreHybrid(d)
			hyb[ds][f] = s
			hitPct[ds][f] = 100 * float64(h) / float64(len(d)-1)
			widths[ds][f] = w
			tsf[ds][f] = scoreTheilSenFixed(d)
		}
	}

	for ds := 0; ds < 2; ds++ {
		hTot, tTot, cTot := 0, 0, 0
		fmt.Printf("DATA %d  hybrid:    ", ds+4)
		for f := 0; f < 5; f++ {
			hTot += hyb[ds][f]
			fmt.Printf("%6d ", hyb[ds][f])
		}
		fmt.Printf(" mean=%d\n", hTot/5)
		fmt.Printf("         TS+fixed:  ")
		for f := 0; f < 5; f++ {
			tTot += tsf[ds][f]
			fmt.Printf("%6d ", tsf[ds][f])
		}
		fmt.Printf(" mean=%d\n", tTot/5)
		fmt.Printf("         current:   ")
		for f := 0; f < 5; f++ {
			cTot += opp["current"][ds][f]
			fmt.Printf("%6d ", opp["current"][ds][f])
		}
		fmt.Printf(" mean=%d\n", cTot/5)
		fmt.Printf("         hybrid hit%%:")
		for f := 0; f < 5; f++ {
			fmt.Printf("%5.1f%% ", hitPct[ds][f])
		}
		fmt.Printf("\n         hybrid w:   ")
		for f := 0; f < 5; f++ {
			fmt.Printf("%6.2f ", widths[ds][f])
		}
		fmt.Println()
	}
	fmt.Println("-----------------------------------------------------------")
	order := []string{"current", "linear-regr", "mse", "nic", "big-range", "correlation-coef"}
	for _, o := range order {
		w4, w5 := 0, 0
		for f := 0; f < 5; f++ {
			if hyb[0][f] > opp[o][0][f] {
				w4++
			}
			if hyb[1][f] > opp[o][1][f] {
				w5++
			}
		}
		fmt.Printf("hybrid    vs %-17s  D4 %d/5   D5 %d/5\n", o, w4, w5)
	}
	fmt.Println()
	for _, o := range order {
		w4, w5 := 0, 0
		for f := 0; f < 5; f++ {
			if tsf[0][f] > opp[o][0][f] {
				w4++
			}
			if tsf[1][f] > opp[o][1][f] {
				w5++
			}
		}
		fmt.Printf("TS+fixed  vs %-17s  D4 %d/5   D5 %d/5\n", o, w4, w5)
	}
}
