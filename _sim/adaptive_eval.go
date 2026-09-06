package main

// Evaluates the proposed "adaptive empirical-error interval" predictor:
// window-30 linear regression centre, then an O(M^2) search over all pairs
// of the last 100 signed prediction errors, picking the [lowErr, highErr]
// band that maximises count/(1+width). Scored with the exact server.js
// formula and compared head-to-head against the audit opponents.
//
//   go run _sim/adaptive_eval.go .resources/guess-it-dockerized/data_sets

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

// linreg fits x = 0..len-1 and extrapolates to x = len.
func linreg(window []float64) float64 {
	n := len(window)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return window[0]
	}
	var sx, sy, sxy, sx2 float64
	for i, y := range window {
		x := float64(i)
		sx += x
		sy += y
		sxy += x * y
		sx2 += x * x
	}
	den := float64(n)*sx2 - sx*sx
	slope := (float64(n)*sxy - sx*sy) / den
	icpt := (sy - slope*sx) / float64(n)
	return icpt + slope*float64(n)
}

// scoreAdaptive replicates the proposed program exactly and returns the
// server.js score, the hit count, and the mean emitted width.
func scoreAdaptive(numbers []float64) (score, hits int, meanWidth float64) {
	n := len(numbers)
	if n < 2 {
		return 0, 0, 0
	}
	recentErrors := []float64{}
	widthSum := 0
	for i := 0; i < n-1; i++ {
		const windowSize = 30
		start := i + 1 - windowSize
		if start < 0 {
			start = 0
		}
		window := numbers[start : i+1]
		pred := linreg(window)

		var low, high int
		if len(recentErrors) < 5 {
			const half = 40.0
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
		if len(recentErrors) > 100 {
			recentErrors = recentErrors[len(recentErrors)-100:]
		}
	}
	return score, hits, float64(widthSum) / float64(n-1)
}

func main() {
	base := os.Args[1]

	// Opponent + current-predictor per-file scores (WSL, exact server.js formula).
	opp := map[string][2][5]int{
		"current":          {{101900, 101640, 104220, 101820, 102580}, {102140, 98840, 101960, 100840, 99820}},
		"linear-regr":      {{102300, 100750, 104532, 101680, 101618}, {100006, 99696, 99448, 97216, 101370}},
		"mse":              {{100392, 101727, 97989, 100659, 104130}, {92115, 107067, 96921, 94785, 100392}},
		"nic":              {{111200, 98400, 106400, 88800, 97600}, {93600, 88000, 108000, 102400, 98400}},
		"big-range":        {{49992, 49996, 49996, 49996, 49996}, {49140, 49168, 49280, 49152, 49172}},
		"correlation-coef": {{95076, 94392, 93442, 92226, 93670}, {89566, 86830, 89718, 89376, 88654}},
	}

	var adp [2][5]int
	var hitPct, widths [2][5]float64
	for ds := 0; ds < 2; ds++ {
		for f := 0; f < 5; f++ {
			d := readData(fmt.Sprintf("%s/%d/%d.txt", base, ds+4, f+1))
			s, h, w := scoreAdaptive(d)
			adp[ds][f] = s
			hitPct[ds][f] = 100 * float64(h) / float64(len(d)-1)
			widths[ds][f] = w
		}
	}

	for ds := 0; ds < 2; ds++ {
		tot := 0
		fmt.Printf("DATA %d  adaptive: ", ds+4)
		for f := 0; f < 5; f++ {
			tot += adp[ds][f]
			fmt.Printf("%6d ", adp[ds][f])
		}
		fmt.Printf(" mean=%d\n", tot/5)
		fmt.Printf("         hit%%:     ")
		for f := 0; f < 5; f++ {
			fmt.Printf("%5.1f%% ", hitPct[ds][f])
		}
		fmt.Printf("\n         width:    ")
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
			if adp[0][f] > opp[o][0][f] {
				w4++
			}
			if adp[1][f] > opp[o][1][f] {
				w5++
			}
		}
		fmt.Printf("adaptive vs %-17s  D4 %d/5   D5 %d/5\n", o, w4, w5)
	}
}
