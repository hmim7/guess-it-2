package main

// Focused benchmark: Theil-Sen centre (window=30) + fixed ±20 half-width,
// scored with the exact server.js formula, against the recorded per-file
// scores of every audit opponent and the current running-OLS predictor.
//
//   go run _sim/theilsen_fixed_eval.go .resources/guess-it-dockerized/data_sets

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

func medianSorted(s []float64) float64 {
	n := len(s)
	if n == 0 {
		return 0
	}
	if n%2 == 1 {
		return s[n/2]
	}
	return (s[n/2-1] + s[n/2]) / 2.0
}

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
	if den == 0 {
		return window[n-1]
	}
	slope := (float64(n)*sxy - sx*sy) / den
	icpt := (sy - slope*sx) / float64(n)
	return icpt + slope*float64(n)
}

func scoreOLSFixed(numbers []float64, half float64) (score, hits int) {
	n := len(numbers)
	if n < 2 {
		return 0, 0
	}
	const windowSize = 30
	for i := 0; i < n-1; i++ {
		start := i + 1 - windowSize
		if start < 0 {
			start = 0
		}
		window := numbers[start : i+1]
		pred := linreg(window)
		low := int(math.Round(pred - half))
		high := int(math.Round(pred + half))
		if high-low < 1 {
			high = low + 1
		}
		actual := numbers[i+1]
		w := high - low
		if actual >= float64(low) && actual <= float64(high) {
			score += int(math.Round(1e7 / float64(1+w) / float64(n-1)))
			hits++
		}
	}
	return
}

func scoreTSFixed(numbers []float64, half float64) (score, hits int) {
	n := len(numbers)
	if n < 2 {
		return 0, 0
	}
	const windowSize = 30
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
			hits++
		}
	}
	return
}

func main() {
	base := os.Args[1]
	half := 20.0
	if len(os.Args) > 2 {
		v, err := strconv.ParseFloat(os.Args[2], 64)
		if err != nil {
			panic(err)
		}
		half = v
	}

	opp := map[string][2][5]int{
		"current":          {{101900, 101640, 104220, 101820, 102580}, {102140, 98840, 101960, 100840, 99820}},
		"linear-regr":      {{102300, 100750, 104532, 101680, 101618}, {100006, 99696, 99448, 97216, 101370}},
		"mse":              {{100392, 101727, 97989, 100659, 104130}, {92115, 107067, 96921, 94785, 100392}},
		"nic":              {{111200, 98400, 106400, 88800, 97600}, {93600, 88000, 108000, 102400, 98400}},
		"big-range":        {{49992, 49996, 49996, 49996, 49996}, {49140, 49168, 49280, 49152, 49172}},
		"correlation-coef": {{95076, 94392, 93442, 92226, 93670}, {89566, 86830, 89718, 89376, 88654}},
	}

	var ts, ols [2][5]int
	var hits [2][5]int
	for ds := 0; ds < 2; ds++ {
		for f := 0; f < 5; f++ {
			d := readData(fmt.Sprintf("%s/%d/%d.txt", base, ds+4, f+1))
			ts[ds][f], hits[ds][f] = scoreTSFixed(d, half)
			ols[ds][f], _ = scoreOLSFixed(d, half)
		}
	}

	dsMean := func(a [5]int) int {
		s := 0
		for _, v := range a {
			s += v
		}
		return s / 5
	}

	fmt.Printf("=== Theil-Sen (window=30) + fixed ±%g — per-file scores ===\n", half)
	fmt.Printf("%-20s %7s %7s %7s %7s %7s   %7s\n", "predictor", "f1", "f2", "f3", "f4", "f5", "mean")
	rows := []string{"current", "Theil-Sen+±20", "linear-regr", "mse", "nic", "big-range", "correlation-coef"}
	for ds := 0; ds < 2; ds++ {
		fmt.Printf("--- DATA %d ---\n", ds+4)
		tsName := fmt.Sprintf("Theil-Sen+±%g", half)
		olsName := fmt.Sprintf("OLS(30)+±%g", half)
		rowsLocal := append([]string(nil), rows[:1]...)
		rowsLocal = append(rowsLocal, tsName, olsName)
		rowsLocal = append(rowsLocal, rows[2:]...)
		for _, name := range rowsLocal {
			var row [5]int
			if name == tsName {
				row = ts[ds]
			} else if name == olsName {
				row = ols[ds]
			} else {
				row = opp[name][ds]
			}
			fmt.Printf("%-20s %7d %7d %7d %7d %7d   %7d\n",
				name, row[0], row[1], row[2], row[3], row[4], dsMean(row))
		}
	}

	order := []string{"current", "linear-regr", "mse", "nic", "big-range", "correlation-coef"}
	printWins := func(label string, s [2][5]int) {
		fmt.Printf("\n=== %s file-wins (head-to-head, per dataset) ===\n", label)
		fmt.Printf("%-20s  %-7s %-7s\n", "opponent", "D4", "D5")
		totalHard := 0
		for _, o := range order {
			w4, w5 := 0, 0
			for f := 0; f < 5; f++ {
				if s[0][f] > opp[o][0][f] {
					w4++
				}
				if s[1][f] > opp[o][1][f] {
					w5++
				}
			}
			fmt.Printf("%-20s  %d/5     %d/5\n", o, w4, w5)
			if o == "linear-regr" || o == "mse" || o == "nic" {
				totalHard += w4 + w5
			}
		}
		fmt.Printf("Hard-opponent file-wins (linear-regr + mse + nic, 30 files): %d/30\n", totalHard)
	}
	printWins(fmt.Sprintf("Theil-Sen +±%g", half), ts)
	printWins(fmt.Sprintf("OLS(30) +±%g", half), ols)
}
