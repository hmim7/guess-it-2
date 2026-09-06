package main

// Characterise the residual distribution of the audit datasets:
// fit a global OLS line over the full prefix and tabulate the empirical
// distribution of residuals. Tests the brainstorm's claim that residuals
// are uniform on [-50, +50] with ~1% outliers (not Gaussian).
//
//   go run _sim/noise_dist.go .resources/guess-it-dockerized/data_sets

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

func globalOLS(d []float64) (slope, intercept float64) {
	n := float64(len(d))
	var sx, sy, sxy, sx2 float64
	for i, y := range d {
		x := float64(i)
		sx += x
		sy += y
		sxy += x * y
		sx2 += x * x
	}
	slope = (n*sxy - sx*sy) / (n*sx2 - sx*sx)
	intercept = (sy - slope*sx) / n
	return
}

func summarise(label string, residuals []float64) {
	n := len(residuals)
	var mean, sd float64
	for _, r := range residuals {
		mean += r
	}
	mean /= float64(n)
	for _, r := range residuals {
		sd += (r - mean) * (r - mean)
	}
	sd = math.Sqrt(sd / float64(n))
	sorted := append([]float64(nil), residuals...)
	sort.Float64s(sorted)
	p := func(q float64) float64 { return sorted[int(q*float64(n))] }
	fmt.Printf("%-12s n=%d mean=%6.2f sd=%5.2f  min=%5.0f p1=%5.0f p5=%4.0f p25=%4.0f p50=%4.0f p75=%4.0f p95=%4.0f p99=%4.0f max=%5.0f\n",
		label, n, mean, sd,
		sorted[0], p(0.01), p(0.05), p(0.25), p(0.50), p(0.75), p(0.95), p(0.99), sorted[n-1])

	// Empirical hit-rate at half-widths used in our sweep, and the
	// uniform[-50,50] reference rate h/50.
	hits := func(h float64) float64 {
		c := 0
		for _, r := range residuals {
			if r >= -h && r <= h {
				c++
			}
		}
		return float64(c) / float64(n)
	}
	fmt.Printf("            P(|r|≤h): h=15 %.4f  h=20 %.4f  h=25 %.4f  h=30 %.4f  h=46 %.4f  h=47 %.4f  h=50 %.4f\n",
		hits(15), hits(20), hits(25), hits(30), hits(46), hits(47), hits(50))
	fmt.Printf("            uniform reference (h/50): h=15 0.3000  h=20 0.4000  h=25 0.5000  h=30 0.6000  h=46 0.9200  h=50 1.0000\n")

	// Bulk SD: exclude |r| > 80 (outliers) and recompute.
	var bulk []float64
	for _, r := range residuals {
		if math.Abs(r) <= 80 {
			bulk = append(bulk, r)
		}
	}
	var bm, bsd float64
	for _, r := range bulk {
		bm += r
	}
	bm /= float64(len(bulk))
	for _, r := range bulk {
		bsd += (r - bm) * (r - bm)
	}
	bsd = math.Sqrt(bsd / float64(len(bulk)))
	outliers := n - len(bulk)
	fmt.Printf("            bulk (|r|≤80): n=%d sd=%.3f  (uniform[-50,+50] sd ≡ 50/√3 = 28.868)  outliers=%d (%.2f%%)\n\n",
		len(bulk), bsd, outliers, 100*float64(outliers)/float64(n))
}

func main() {
	base := os.Args[1]
	for ds := 0; ds < 2; ds++ {
		for f := 0; f < 5; f++ {
			d := readData(fmt.Sprintf("%s/%d/%d.txt", base, ds+4, f+1))
			m, b := globalOLS(d)
			res := make([]float64, len(d))
			for i, y := range d {
				res[i] = y - (m*float64(i) + b)
			}
			fmt.Printf("D%d/%d  slope=%.6f intercept=%.3f\n", ds+4, f+1, m, b)
			summarise(fmt.Sprintf("D%d/%d", ds+4, f+1), res)
		}
	}
}
