package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func average(d []float64) float64 {
	if len(d) == 0 {
		return 0
	}
	var s float64
	for _, v := range d {
		s += v
	}
	return s / float64(len(d))
}

func median(d []float64) float64 {
	n := len(d)
	if n == 0 {
		return 0
	}
	c := make([]float64, n)
	copy(c, d)
	sort.Float64s(c)
	if n%2 == 1 {
		return c[n/2]
	}
	return (c[n/2-1] + c[n/2]) / 2
}

func linReg(d []float64) (m, b float64) {
	n := len(d)
	if n == 0 {
		return 0, 0
	}
	if n == 1 {
		return 0, d[0]
	}
	fn := float64(n)
	var sx, sy, sxy, sx2 float64
	for i, y := range d {
		x := float64(i)
		sx += x
		sy += y
		sxy += x * y
		sx2 += x * x
	}
	den := fn*sx2 - sx*sx
	if den == 0 {
		return 0, average(d)
	}
	m = (fn*sxy - sx*sy) / den
	b = (sy - m*sx) / fn
	return m, b
}

func roundBounds(lo, hi float64) (int64, int64) {
	l := int64(math.Round(lo))
	u := int64(math.Round(hi))
	if l > u {
		l, u = u, l
	}
	if u-l < 1 {
		u = l + 1
	}
	return l, u
}

func lastK(w []float64, k int) []float64 {
	if k >= len(w) {
		return w
	}
	return w[len(w)-k:]
}

type config struct {
	name string
	half float64
	ctr  func(w []float64, seen int) float64
}

func regNext(w []float64) float64 {
	m, b := linReg(w)
	return m*float64(len(w)) + b
}

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
		if v, err := strconv.ParseFloat(line, 64); err == nil {
			d = append(d, v)
		}
	}
	return d
}

const windowSize = 20

func scoreFile(data []float64, c config) int {
	n := len(data)
	window := make([]float64, 0, windowSize)
	var score int
	for i := 0; i < n-1; i++ {
		if len(window) < windowSize {
			window = append(window, data[i])
		} else {
			window = append(window[1:], data[i])
		}
		ctr := c.ctr(window, i+1)
		lo, hi := roundBounds(ctr-c.half, ctr+c.half)
		next := data[i+1]
		if next >= float64(lo) && next <= float64(hi) {
			w := hi - lo
			score += int(math.Round(10000000.0 / float64(1+w) / float64(n-1)))
		}
	}
	return score
}

// Opponent per-file scores measured in WSL (server.js formula).
var opp = map[string][2][5]int{
	"linear-regr":      {{102300, 100750, 104532, 101680, 101618}, {100006, 99696, 99448, 97216, 101370}},
	"mse":              {{100392, 101727, 97989, 100659, 104130}, {92115, 107067, 96921, 94785, 100392}},
	"nic":              {{111200, 98400, 106400, 88800, 97600}, {93600, 88000, 108000, 102400, 98400}},
	"big-range":        {{49992, 49996, 49996, 49996, 49996}, {49140, 49168, 49280, 49152, 49172}},
	"correlation-coef": {{95076, 94392, 93442, 92226, 93670}, {89566, 86830, 89718, 89376, 88654}},
}

func main() {
	base := os.Args[1]
	var files [2][5][]float64
	for ds := 0; ds < 2; ds++ {
		for f := 0; f < 5; f++ {
			files[ds][f] = readData(fmt.Sprintf("%s/%d/%d.txt", base, ds+4, f+1))
		}
	}

	mk := func(name string, half float64, ctr func([]float64, int) float64) config {
		return config{name, half, ctr}
	}
	medAll := func(k int) func([]float64, int) float64 {
		return func(w []float64, seen int) float64 { return median(lastK(w, k)) }
	}
	avgAll := func(k int) func([]float64, int) float64 {
		return func(w []float64, seen int) float64 { return average(lastK(w, k)) }
	}
	split := func(w []float64, seen int) float64 {
		if seen < 1000 {
			return regNext(w)
		}
		return median(w)
	}

	var configs []config
	configs = append(configs, mk("current split ±20", 20, split))
	configs = append(configs, mk("avg(10) ±0.5", 0.5, avgAll(10)))
	for _, k := range []int{10, 15, 20} {
		for _, hw := range []float64{13, 16, 20, 24} {
			configs = append(configs, mk(fmt.Sprintf("median(%d) ±%.0f", k, hw), hw, medAll(k)))
		}
	}
	for _, k := range []int{8, 10, 15} {
		for _, hw := range []float64{16, 20, 24} {
			configs = append(configs, mk(fmt.Sprintf("avg(%d) ±%.0f", k, hw), hw, avgAll(k)))
		}
	}
	for _, hw := range []float64{16, 20, 24} {
		configs = append(configs, mk(fmt.Sprintf("reg(20) ±%.0f", hw), hw, func(w []float64, seen int) float64 { return regNext(w) }))
	}

	hard := []string{"linear-regr", "mse", "nic"}
	type res struct {
		name    string
		stuMean int
		winsD4  map[string]int
		winsD5  map[string]int
		hardW   int // total file-wins vs the 3 hard opponents over 10 files
	}
	var results []res
	for _, c := range configs {
		var stu [2][5]int
		tot := 0
		for ds := 0; ds < 2; ds++ {
			for f := 0; f < 5; f++ {
				stu[ds][f] = scoreFile(files[ds][f], c)
				tot += stu[ds][f]
			}
		}
		r := res{name: c.name, stuMean: tot / 10, winsD4: map[string]int{}, winsD5: map[string]int{}}
		for _, o := range []string{"linear-regr", "mse", "nic", "big-range", "correlation-coef"} {
			for f := 0; f < 5; f++ {
				if stu[0][f] > opp[o][0][f] {
					r.winsD4[o]++
				}
				if stu[1][f] > opp[o][1][f] {
					r.winsD5[o]++
				}
			}
		}
		for _, o := range hard {
			r.hardW += r.winsD4[o] + r.winsD5[o]
		}
		results = append(results, r)
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].hardW != results[j].hardW {
			return results[i].hardW > results[j].hardW
		}
		return results[i].stuMean > results[j].stuMean
	})

	fmt.Printf("%-22s %8s | %-13s %-13s %-13s\n",
		"config", "stuMean", "linear-regr", "mse", "nic")
	fmt.Printf("%-22s %8s | %-13s %-13s %-13s\n",
		"", "", "D4 / D5", "D4 / D5", "D5 / D5")
	for _, r := range results {
		fmt.Printf("%-22s %8d | %d/5 %d/5      %d/5 %d/5      %d/5 %d/5     (hard wins %2d/30)\n",
			r.name, r.stuMean,
			r.winsD4["linear-regr"], r.winsD5["linear-regr"],
			r.winsD4["mse"], r.winsD5["mse"],
			r.winsD4["nic"], r.winsD5["nic"], r.hardW)
	}
}
