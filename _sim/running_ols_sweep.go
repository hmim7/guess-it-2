package main

// Replicates the current Predictor exactly (running OLS over the full prefix,
// kicking in at seen >= 30, with a window-20 regression fallback below that)
// parametrised by half-width. Used to test the brainstorm's claim that h=46
// scores better than h=20 because the audit data is uniform[-50,+50] noise.
//
//   go run _sim/running_ols_sweep.go .resources/guess-it-dockerized/data_sets

import (
	"bufio"
	"fmt"
	"math"
	"os"
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

// RunningOLS — same incremental Σx, Σy, Σxy, Σx² update as linearstats.RunningOLS.
type RunningOLS struct{ n, sx, sy, sxy, sx2 float64 }

func (r *RunningOLS) Add(x, y float64) {
	r.n++
	r.sx += x
	r.sy += y
	r.sxy += x * y
	r.sx2 += x * x
}

func (r *RunningOLS) Fit() (m, b float64, ok bool) {
	if r.n < 2 {
		return 0, 0, false
	}
	den := r.n*r.sx2 - r.sx*r.sx
	if den == 0 {
		return 0, 0, false
	}
	m = (r.n*r.sxy - r.sx*r.sy) / den
	b = (r.sy - m*r.sx) / r.n
	return m, b, true
}

func windowReg(window []float64) float64 {
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

// scoreRunningOLS replicates Predictor.Next exactly with a configurable half.
func scoreRunningOLS(numbers []float64, half float64) (score, hits int) {
	const (
		windowSize    = 20
		olsMinSamples = 30
	)
	var ols RunningOLS
	window := make([]float64, 0, windowSize)
	n := len(numbers)
	for i := 0; i < n-1; i++ {
		cur := numbers[i]
		if len(window) < windowSize {
			window = append(window, cur)
		} else {
			window = append(window[1:], cur)
		}
		ols.Add(float64(i), cur)
		seen := i + 1

		var center float64
		if seen >= olsMinSamples {
			m, b, ok := ols.Fit()
			if ok {
				center = m*float64(seen) + b
			} else {
				center = windowReg(window)
			}
		} else {
			center = windowReg(window)
		}
		low := int64(math.Round(center - half))
		high := int64(math.Round(center + half))
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
	halves := []float64{15, 20, 25, 27, 30, 41, 45, 46, 47, 50}

	opp := map[string][2][5]int{
		"current_recorded": {{101900, 101640, 104220, 101820, 102580}, {102140, 98840, 101960, 100840, 99820}},
		"linear-regr":      {{102300, 100750, 104532, 101680, 101618}, {100006, 99696, 99448, 97216, 101370}},
		"mse":              {{100392, 101727, 97989, 100659, 104130}, {92115, 107067, 96921, 94785, 100392}},
		"nic":              {{111200, 98400, 106400, 88800, 97600}, {93600, 88000, 108000, 102400, 98400}},
		"big-range":        {{49992, 49996, 49996, 49996, 49996}, {49140, 49168, 49280, 49152, 49172}},
		"correlation-coef": {{95076, 94392, 93442, 92226, 93670}, {89566, 86830, 89718, 89376, 88654}},
	}

	files := make([][][]float64, 2)
	for ds := 0; ds < 2; ds++ {
		files[ds] = make([][]float64, 5)
		for f := 0; f < 5; f++ {
			files[ds][f] = readData(fmt.Sprintf("%s/%d/%d.txt", base, ds+4, f+1))
		}
	}

	fmt.Printf("%-7s %-5s | %-37s %-37s | %-9s %-9s\n",
		"half", "ds", "per-file scores", "hit %", "mean", "Δ vs h=20")

	// Baseline = h=20 for the delta column.
	var base20 [2]int
	for ds := 0; ds < 2; ds++ {
		tot := 0
		for f := 0; f < 5; f++ {
			s, _ := scoreRunningOLS(files[ds][f], 20)
			tot += s
		}
		base20[ds] = tot / 5
	}

	results := map[float64][2][5]int{}
	for _, h := range halves {
		var scoresAll [2][5]int
		var hitPctAll [2][5]float64
		for ds := 0; ds < 2; ds++ {
			for f := 0; f < 5; f++ {
				d := files[ds][f]
				s, hh := scoreRunningOLS(d, h)
				scoresAll[ds][f] = s
				hitPctAll[ds][f] = 100 * float64(hh) / float64(len(d)-1)
			}
			tot := 0
			for f := 0; f < 5; f++ {
				tot += scoresAll[ds][f]
			}
			mean := tot / 5
			fmt.Printf("±%-5.0f D%d   | %6d %6d %6d %6d %6d  %5.1f %5.1f %5.1f %5.1f %5.1f | %9d %+9d\n",
				h, ds+4,
				scoresAll[ds][0], scoresAll[ds][1], scoresAll[ds][2], scoresAll[ds][3], scoresAll[ds][4],
				hitPctAll[ds][0], hitPctAll[ds][1], hitPctAll[ds][2], hitPctAll[ds][3], hitPctAll[ds][4],
				mean, mean-base20[ds])
		}
		results[h] = scoresAll
	}

	fmt.Println("\nFile-wins vs each opponent (head-to-head, per dataset):")
	fmt.Printf("%-7s | %-13s %-13s %-13s %-13s %-13s\n",
		"half", "current_rec", "linear-regr", "mse", "nic", "all-hard")
	order := []string{"current_recorded", "linear-regr", "mse", "nic"}
	for _, h := range halves {
		s := results[h]
		fmt.Printf("±%-5.0f | ", h)
		hard4, hard5 := 0, 0
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
			fmt.Printf("%d/5 %d/5     ", w4, w5)
			if o == "linear-regr" || o == "mse" || o == "nic" {
				hard4 += w4
				hard5 += w5
			}
		}
		fmt.Printf("%d/15 %d/15\n", hard4, hard5)
	}
}
