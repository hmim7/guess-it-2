package main

// Residual-structure analysis for guess-it-2.
// Question: are the linear-regression residuals white noise, or do they carry
// exploitable structure (autocorrelation / periodicity)? If white, 3/5 vs
// linear-regr is the true ceiling. If structured, an AR/periodic correction
// could give a real edge.
//
// Run:  go run _sim/structure.go .resources/guess-it-dockerized/data_sets

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const windowSize = 20

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
		return 0, sy / fn
	}
	m = (fn*sxy - sx*sy) / den
	b = (sy - m*sx) / fn
	return m, b
}

func regNext(w []float64) float64 {
	m, b := linReg(w)
	return m*float64(len(w)) + b
}

// acf returns the autocorrelation of x at the given lag.
func acf(x []float64, lag int) float64 {
	n := len(x)
	if lag <= 0 || lag >= n {
		return 0
	}
	var mean float64
	for _, v := range x {
		mean += v
	}
	mean /= float64(n)
	var num, den float64
	for i := 0; i < n; i++ {
		dv := x[i] - mean
		den += dv * dv
		if i+lag < n {
			num += dv * (x[i+lag] - mean)
		}
	}
	if den == 0 {
		return 0
	}
	return num / den
}

func stddev(x []float64) float64 {
	n := len(x)
	if n == 0 {
		return 0
	}
	var mean float64
	for _, v := range x {
		mean += v
	}
	mean /= float64(n)
	var s float64
	for _, v := range x {
		d := v - mean
		s += d * d
	}
	return math.Sqrt(s / float64(n))
}

// residuals: e_i = data[i+1] - regNext(window ending at i), window size 20.
func residuals(data []float64) (resid []float64, centers []float64) {
	n := len(data)
	window := make([]float64, 0, windowSize)
	for i := 0; i < n-1; i++ {
		if len(window) < windowSize {
			window = append(window, data[i])
		} else {
			window = append(window[1:], data[i])
		}
		c := regNext(window)
		centers = append(centers, c)
		resid = append(resid, data[i+1]-c)
	}
	return resid, centers
}

// score with a fixed ±20 range about the given centers.
func scoreCenters(data []float64, centers []float64) int {
	n := len(data)
	perHit := int(math.Round(1e7 / 41.0 / float64(n-1)))
	hits := 0
	for i := 0; i < len(centers); i++ {
		lo := math.Round(centers[i] - 20)
		hi := math.Round(centers[i] + 20)
		next := data[i+1]
		if next >= lo && next <= hi {
			hits++
		}
	}
	return hits * perHit
}

// AR(1)-corrected centers: c'_i = c_i + phi * e_{i-1}, e = data[i]-c_{i-1}.
func ar1Centers(data, centers []float64, phi float64) []float64 {
	adj := make([]float64, len(centers))
	prevResid := 0.0
	for i := 0; i < len(centers); i++ {
		adj[i] = centers[i] + phi*prevResid
		prevResid = data[i+1] - centers[i]
	}
	return adj
}

func main() {
	base := os.Args[1]
	fmt.Println("Residual-structure analysis — window-20 linear regression")
	fmt.Println(strings.Repeat("=", 92))
	fmt.Printf("%-10s %8s %8s | %7s %7s %7s | %7s %7s | %5s@%-4s | %8s %8s\n",
		"file", "resSD", "diffSD", "rACF1", "rACF2", "rACF3", "dACF1", "dACF2",
		"|pk|", "lag", "base", "AR(1)")
	fmt.Println(strings.Repeat("-", 92))

	var allRes []float64
	totBase, totAR := 0, 0
	for ds := 4; ds <= 5; ds++ {
		for f := 1; f <= 5; f++ {
			data := readData(fmt.Sprintf("%s/%d/%d.txt", base, ds, f))
			res, centers := residuals(data)
			allRes = append(allRes, res...)

			diffs := make([]float64, len(data)-1)
			for i := 0; i < len(data)-1; i++ {
				diffs[i] = data[i+1] - data[i]
			}

			// periodicity: strongest residual autocorrelation over lags 2..300
			pk, pkLag := 0.0, 0
			maxLag := 300
			if maxLag > len(res)/4 {
				maxLag = len(res) / 4
			}
			for lag := 2; lag <= maxLag; lag++ {
				a := math.Abs(acf(res, lag))
				if a > pk {
					pk, pkLag = a, lag
				}
			}

			phi := acf(res, 1)
			base := scoreCenters(data, centers)
			ar := scoreCenters(data, ar1Centers(data, centers, phi))
			totBase += base
			totAR += ar

			fmt.Printf("D%d/%d      %8.1f %8.1f | %7.3f %7.3f %7.3f | %7.3f %7.3f | %5.3f@%-4d | %8d %8d\n",
				ds, f, stddev(res), stddev(diffs),
				acf(res, 1), acf(res, 2), acf(res, 3),
				acf(diffs, 1), acf(diffs, 2), pk, pkLag, base, ar)
		}
	}
	fmt.Println(strings.Repeat("-", 92))
	fmt.Printf("%-10s %8.1f %50s %8d %8d\n", "ALL", stddev(allRes), "", totBase/10, totAR/10)
	fmt.Println(strings.Repeat("=", 92))
	fmt.Printf("Aggregate residual ACF: lag1=%.4f lag2=%.4f lag3=%.4f lag4=%.4f lag5=%.4f\n",
		acf(allRes, 1), acf(allRes, 2), acf(allRes, 3), acf(allRes, 4), acf(allRes, 5))
	white := 1.96 / math.Sqrt(float64(len(allRes)))
	fmt.Printf("White-noise 95%% band: ±%.4f  (|ACF| below this = indistinguishable from noise)\n", white)
	fmt.Printf("base mean=%d  AR(1) mean=%d  delta=%+d (%.2f%%)\n",
		totBase/10, totAR/10, (totAR-totBase)/10,
		100*float64(totAR-totBase)/float64(totBase))
}
