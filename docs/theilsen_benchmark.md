# Theil-Sen Predictor Benchmark

**Date:** 2026-05-23
**Question:** Does swapping the running-OLS centre for a Theil-Sen estimator
([.resources/Theil_Sen_estimator.md](../.resources/Theil_Sen_estimator.md))
improve the audit score?

**Short answer (Theil-Sen):** No — the Theil-Sen swap itself is a wash at every
width tested.

**Short answer (half-width):** ±20 is **NOT** the optimum. ±46 is, by ~+2% mean
and +3 hard-opponent file-wins. See §6 below — the original ±20 sweep used a
window-N centre that added enough extrapolation variance to mask a rounding
cliff at ±47. With the actual running-OLS-prefix centre (σ_ŷ → 0 with growing
prefix), the residuals are pure uniform[-50,+50] and h=46 is the last width
before `round(10⁷/(2h+1)/(N-1))` drops from 9 to 8 per hit. ±15 and ±25 results
are mirrored on a fair-comparison OLS(30) variant — both centres take the
same hit, confirming the width is the lever, not the centre. The hybrid (TS
+ adaptive empirical interval) collapses to ~90k
([prior result](#footnote-hybrid)).

Scored with the exact `server.js` formula via
[_sim/theilsen_fixed_eval.go](../_sim/theilsen_fixed_eval.go). Opponent and
current-predictor scores are the recorded WSL values from
[predictor_benchmark_linear_v2.md](predictor_benchmark_linear_v2.md).

---

## 1. Theil-Sen (window = 30) + fixed ±20

Run: `go run _sim/theilsen_fixed_eval.go .resources/guess-it-dockerized/data_sets 20`

### Per-file scores

**Data 4**

| predictor          |    f1   |    f2   |    f3   |    f4   |    f5   |  mean  |
|--------------------|--------:|--------:|--------:|--------:|--------:|-------:|
| current (running-OLS) | 101,900 | 101,640 | 104,220 | 101,820 | 102,580 | **102,432** |
| **Theil-Sen ±20**     | 102,620 | 102,680 | 104,280 | 103,360 | 100,880 | **102,764** |
| linear-regr           | 102,300 | 100,750 | 104,532 | 101,680 | 101,618 | 102,176 |
| mse                   | 100,392 | 101,727 |  97,989 | 100,659 | 104,130 | 100,979 |
| nic                   | 111,200 |  98,400 | 106,400 |  88,800 |  97,600 | 100,480 |
| big-range             |  49,992 |  49,996 |  49,996 |  49,996 |  49,996 |  49,995 |
| correlation-coef      |  95,076 |  94,392 |  93,442 |  92,226 |  93,670 |  93,761 |

**Data 5**

| predictor          |    f1   |    f2   |    f3   |    f4   |    f5   |  mean  |
|--------------------|--------:|--------:|--------:|--------:|--------:|-------:|
| current (running-OLS) | 102,140 |  98,840 | 101,960 | 100,840 |  99,820 | **100,720** |
| **Theil-Sen ±20**     | 102,140 | 100,660 | 100,740 |  99,740 | 101,598 | **100,975** |
| linear-regr           | 100,006 |  99,696 |  99,448 |  97,216 | 101,370 |  99,547 |
| mse                   |  92,115 | 107,067 |  96,921 |  94,785 | 100,392 |  98,256 |
| nic                   |  93,600 |  88,000 | 108,000 | 102,400 |  98,400 |  98,080 |
| big-range             |  49,140 |  49,168 |  49,280 |  49,152 |  49,172 |  49,182 |
| correlation-coef      |  89,566 |  86,830 |  89,718 |  89,376 |  88,654 |  88,828 |

### File-wins vs each opponent

| opponent            | D4  | D5  |
|---------------------|----:|----:|
| current             | 4/5 | 2/5 |
| linear-regr         | 3/5 | 5/5 |
| mse                 | 4/5 | 4/5 |
| nic                 | 3/5 | 3/5 |
| big-range           | 5/5 | 5/5 |
| correlation-coef    | 5/5 | 5/5 |

**Hard-opponent file-wins (linear-regr + mse + nic):** 22/30
— identical to the current predictor's 22/30.

### Read

- **Mean delta vs current: +332 (D4), +255 (D5)** — ≈0.3%, inside the noise band
  ([predictor_analysis.md](predictor_analysis.md) documents the score-landscape
  as flat in a 95k–103k band).
- **Audit-relevant outcome is identical** — 5/5 vs `big-range` and
  `correlation-coef` on both datasets, 3/5 vs `nic`, 3-5/5 vs `linear-regr`/`mse`.
  Same audit-pass probability.
- The 6/10 "wins" vs current are coin-flips on the noise: per-file deltas
  (e.g. D4f4 +1540, D4f5 −1700) are residual-SD-scale (~32 on D4, ~48 on D5),
  exactly what the white-residual proof predicts.
- **Cost:** ~30× per-step CPU (435 pairwise slopes per window vs O(1) running-OLS
  update). No statistical gain to justify the swap.

---

## 2. Theil-Sen (window = 30) + fixed ±30

Run: `go run _sim/theilsen_fixed_eval.go .resources/guess-it-dockerized/data_sets 30`

### Per-file scores

**Data 4**

| predictor          |    f1   |    f2   |    f3   |    f4   |    f5   |  mean  |
|--------------------|--------:|--------:|--------:|--------:|--------:|-------:|
| current (running-OLS) | 101,900 | 101,640 | 104,220 | 101,820 | 102,580 | **102,432** |
| **Theil-Sen ±30**     |  97,422 |  98,098 |  98,865 |  97,799 |  97,058 | **97,848** |
| linear-regr           | 102,300 | 100,750 | 104,532 | 101,680 | 101,618 | 102,176 |
| mse                   | 100,392 | 101,727 |  97,989 | 100,659 | 104,130 | 100,979 |
| nic                   | 111,200 |  98,400 | 106,400 |  88,800 |  97,600 | 100,480 |
| big-range             |  49,992 |  49,996 |  49,996 |  49,996 |  49,996 |  49,995 |
| correlation-coef      |  95,076 |  94,392 |  93,442 |  92,226 |  93,670 |  93,761 |

**Data 5**

| predictor          |    f1   |    f2   |    f3   |    f4   |    f5   |  mean  |
|--------------------|--------:|--------:|--------:|--------:|--------:|-------:|
| current (running-OLS) | 102,140 |  98,840 | 101,960 | 100,840 |  99,820 | **100,720** |
| **Theil-Sen ±30**     |  97,097 |  96,161 |  96,902 |  96,148 |  97,201 | **96,701** |
| linear-regr           | 100,006 |  99,696 |  99,448 |  97,216 | 101,370 |  99,547 |
| mse                   |  92,115 | 107,067 |  96,921 |  94,785 | 100,392 |  98,256 |
| nic                   |  93,600 |  88,000 | 108,000 | 102,400 |  98,400 |  98,080 |
| big-range             |  49,140 |  49,168 |  49,280 |  49,152 |  49,172 |  49,182 |
| correlation-coef      |  89,566 |  86,830 |  89,718 |  89,376 |  88,654 |  88,828 |

### File-wins vs each opponent

| opponent            | D4  | D5  |
|---------------------|----:|----:|
| current             | 0/5 | 0/5 |
| linear-regr         | 0/5 | 0/5 |
| mse                 | 1/5 | 2/5 |
| nic                 | 1/5 | 2/5 |
| big-range           | 5/5 | 5/5 |
| correlation-coef    | 5/5 | 5/5 |

**Hard-opponent file-wins (linear-regr + mse + nic):** 6/30
— a 16-file regression versus the current 22/30. Loses 0/10 versus the current
predictor.

### Read

- **Mean delta vs current: −4,584 (D4), −4,019 (D5)** — ≈ −4.5%.
- ±30 yields width 60. Hit rate gains ~50% (band is 50% wider), but score per
  hit drops by `(1+40)/(1+60) ≈ 0.67`, a 33% per-hit penalty. The two cancel
  imperfectly: net ~-4% because higher hit rate doesn't quite compensate, and
  any hits above 100% saturate.
- Audit pass probability against `linear-regr`/`mse` collapses from ~65-100% to
  near-zero on D4; survives on big-range/correlation-coef only.
- **Verdict:** ±30 is strictly worse on every audit-relevant axis. The fixed-±20
  hyperparameter chosen by the current predictor's sweep is at the score peak.

---

## 3. Theil-Sen and OLS(30) at fixed ±15 (apples-to-apples width sweep)

Run: `go run _sim/theilsen_fixed_eval.go .resources/guess-it-dockerized/data_sets 15`

Tightening the band from ±20 → ±15 lets us isolate the **centre** effect: both
predictors share the same window (30) and the same half-width, so any score
gap between them is purely the OLS-vs-Theil-Sen choice. The current
running-OLS predictor (full-prefix OLS, ±20) is included for absolute
reference.

### Per-file scores

**Data 4**

| predictor          |    f1   |    f2   |    f3   |    f4   |    f5   |  mean  |
|--------------------|--------:|--------:|--------:|--------:|--------:|-------:|
| current (running-OLS ±20) | 101,900 | 101,640 | 104,220 | 101,820 | 102,580 | **102,432** |
| **Theil-Sen ±15**         | 101,738 | 100,750 | 103,506 | 102,310 | 100,646 | **101,790** |
| **OLS(30) ±15**           | 101,218 | 100,880 | 103,350 | 102,050 | 101,192 | **101,738** |
| linear-regr               | 102,300 | 100,750 | 104,532 | 101,680 | 101,618 | 102,176 |
| mse                       | 100,392 | 101,727 |  97,989 | 100,659 | 104,130 | 100,979 |
| nic                       | 111,200 |  98,400 | 106,400 |  88,800 |  97,600 | 100,480 |
| big-range                 |  49,992 |  49,996 |  49,996 |  49,996 |  49,996 |  49,995 |
| correlation-coef          |  95,076 |  94,392 |  93,442 |  92,226 |  93,670 |  93,761 |

**Data 5**

| predictor          |    f1   |    f2   |    f3   |    f4   |    f5   |  mean  |
|--------------------|--------:|--------:|--------:|--------:|--------:|-------:|
| current (running-OLS ±20) | 102,140 |  98,840 | 101,960 | 100,840 |  99,820 | **100,720** |
| **Theil-Sen ±15**         | 101,894 |  98,462 |  99,684 |  98,748 |  99,605 | **99,678** |
| **OLS(30) ±15**           | 100,308 |  97,214 |  97,344 |  97,604 |  96,382 | **97,770** |
| linear-regr               | 100,006 |  99,696 |  99,448 |  97,216 | 101,370 |  99,547 |
| mse                       |  92,115 | 107,067 |  96,921 |  94,785 | 100,392 |  98,256 |
| nic                       |  93,600 |  88,000 | 108,000 | 102,400 |  98,400 |  98,080 |
| big-range                 |  49,140 |  49,168 |  49,280 |  49,152 |  49,172 |  49,182 |
| correlation-coef          |  89,566 |  86,830 |  89,718 |  89,376 |  88,654 |  88,828 |

### File-wins vs each opponent

| opponent            | Theil-Sen ±15 D4 | D5  | OLS(30) ±15 D4 | D5  |
|---------------------|---:|---:|---:|---:|
| current             | 1/5 | 0/5 | 1/5 | 0/5 |
| linear-regr         | 1/5 | 3/5 | 2/5 | 2/5 |
| mse                 | 3/5 | 3/5 | 3/5 | 3/5 |
| nic                 | 3/5 | 3/5 | 3/5 | 2/5 |
| big-range           | 5/5 | 5/5 | 5/5 | 5/5 |
| correlation-coef    | 5/5 | 5/5 | 5/5 | 5/5 |

**Hard-opponent file-wins:** Theil-Sen ±15 = **16/30** · OLS(30) ±15 = **15/30**
— both regress from the 22/30 baseline by roughly the same margin.

### Read

- **Mean delta vs current:** Theil-Sen ±15 = −642 (D4), −1,042 (D5)
  (≈ −0.8%); OLS(30) ±15 = −694 (D4), −2,950 (D5) (≈ −1.6%).
- **Centre choice is irrelevant.** With identical window=30 and half=15, the
  two centres produce mean scores within 50 points (D4) and 1,900 points (D5);
  the file-by-file wins are also nearly identical (16/30 vs 15/30). Theil-Sen
  has no detectable edge over plain OLS — confirming the white-residual
  prediction.
- **±15 is past the score peak.** Width 30 trims +25% off the per-hit score
  (`(1+40)/(1+30) ≈ 1.32`) but cuts hit rate by ~25% on iid-Gaussian residuals
  with SD ≈ 32 (D4) / 48 (D5). Net is a small loss on D4 and a larger loss on
  D5 (noisier data → more hits drop). Confirms ±20 sits at the optimum.
- **Audit-pass probability:** 5/5 vs `big-range` and `correlation-coef` holds.
  Loses 3/5 vs `linear-regr` on D4 and 2/5 vs `nic` on D5 versus current — not
  audit-safe in the same configuration as the baseline.

---

## 4. Theil-Sen and OLS(30) at fixed ±25 (apples-to-apples width sweep)

Run: `go run _sim/theilsen_fixed_eval.go .resources/guess-it-dockerized/data_sets 25`

Widening from ±20 → ±25 — the other side of the score peak from §3 — to
confirm the optimum is at ±20 from both directions.

### Per-file scores

**Data 4**

| predictor          |    f1   |    f2   |    f3   |    f4   |    f5   |  mean  |
|--------------------|--------:|--------:|--------:|--------:|--------:|-------:|
| current (running-OLS ±20) | 101,900 | 101,640 | 104,220 | 101,820 | 102,580 | **102,432** |
| **Theil-Sen ±25**         | 101,744 | 101,807 | 102,128 | 101,408 | 100,688 | **101,555** |
| **OLS(30) ±25**           | 101,616 | 102,368 | 103,328 | 101,904 | 101,600 | **102,163** |
| linear-regr               | 102,300 | 100,750 | 104,532 | 101,680 | 101,618 | 102,176 |
| mse                       | 100,392 | 101,727 |  97,989 | 100,659 | 104,130 | 100,979 |
| nic                       | 111,200 |  98,400 | 106,400 |  88,800 |  97,600 | 100,480 |
| big-range                 |  49,992 |  49,996 |  49,996 |  49,996 |  49,996 |  49,995 |
| correlation-coef          |  95,076 |  94,392 |  93,442 |  92,226 |  93,670 |  93,761 |

**Data 5**

| predictor          |    f1   |    f2   |    f3   |    f4   |    f5   |  mean  |
|--------------------|--------:|--------:|--------:|--------:|--------:|-------:|
| current (running-OLS ±20) | 102,140 |  98,840 | 101,960 | 100,840 |  99,820 | **100,720** |
| **Theil-Sen ±25**         | 100,528 |  99,598 | 100,736 |  99,520 | 101,086 | **100,293** |
| **OLS(30) ±25**           |  99,840 |  98,752 |  99,232 |  97,936 |  97,584 | **98,668** |
| linear-regr               | 100,006 |  99,696 |  99,448 |  97,216 | 101,370 |  99,547 |
| mse                       |  92,115 | 107,067 |  96,921 |  94,785 | 100,392 |  98,256 |
| nic                       |  93,600 |  88,000 | 108,000 | 102,400 |  98,400 |  98,080 |
| big-range                 |  49,140 |  49,168 |  49,280 |  49,152 |  49,172 |  49,182 |
| correlation-coef          |  89,566 |  86,830 |  89,718 |  89,376 |  88,654 |  88,828 |

### File-wins vs each opponent

| opponent            | Theil-Sen ±25 D4 | D5  | OLS(30) ±25 D4 | D5  |
|---------------------|---:|---:|---:|---:|
| current             | 1/5 | 2/5 | 2/5 | 0/5 |
| linear-regr         | 1/5 | 3/5 | 2/5 | 1/5 |
| mse                 | 4/5 | 4/5 | 4/5 | 3/5 |
| nic                 | 3/5 | 3/5 | 3/5 | 2/5 |
| big-range           | 5/5 | 5/5 | 5/5 | 5/5 |
| correlation-coef    | 5/5 | 5/5 | 5/5 | 5/5 |

**Hard-opponent file-wins:** Theil-Sen ±25 = **18/30** · OLS(30) ±25 = **15/30**
— both regress from the 22/30 baseline; ±25 is a smaller hit than ±15 on D4 but
about the same on D5.

### Read

- **Mean delta vs current:** Theil-Sen ±25 = −877 (D4), −427 (D5) (≈ −0.7%);
  OLS(30) ±25 = −269 (D4), −2,052 (D5) (≈ −1.2%). The score peak sits between
  ±20 and ±25 — closer to ±20 on D4, slightly less penalty on D5 because the
  larger residual SD (~48) makes a wider band relatively cheaper.
- **Centre choice still irrelevant.** Theil-Sen and OLS(30) at the same width
  produce mean scores within 600 pts (D4) and 1,600 pts (D5), with file-wins
  18/30 vs 15/30 — same noise band as §3. Robust centre buys nothing.
- **Confirms ±20 is the peak.** Both ±15 (§3) and ±25 regress vs the baseline
  ±20; the score landscape is concave and the optimum is at ±20. Both
  directions cost roughly 0.5-1.5% on the mean and drop hard-opponent wins
  from 22/30 into the 15-18/30 band.
- **Audit-pass probability:** holds against `big-range` and `correlation-coef`
  (5/5 both datasets) but loses against `linear-regr`/`nic` on D4 — not
  audit-safe in the same configuration as the baseline.

---

## 5. Theil-Sen and OLS(30) at fixed ±46 (extreme-width sanity check)

Run: `go run _sim/theilsen_fixed_eval.go .resources/guess-it-dockerized/data_sets 46`

±46 is far past the peak. At width 92 the per-hit score collapses to
`10⁷/(1+92)/(N−1)`, ≈ half the ±20 per-hit value, while hit rate saturates
near 100% on D4 (residual SD ≈ 32, so ±46 ≈ 1.4σ ≈ ~84% coverage; rounding +
edge buffer pushes it up). This sets a lower bound on the wide-band regime.

### Per-file scores

**Data 4**

| predictor          |    f1   |    f2   |    f3   |    f4   |    f5   |  mean  |
|--------------------|--------:|--------:|--------:|--------:|--------:|-------:|
| current (running-OLS ±20) | 101,900 | 101,640 | 104,220 | 101,820 | 102,580 | **102,432** |
| **Theil-Sen ±46**         |  95,805 |  95,940 |  96,462 |  95,994 |  96,219 | **96,084** |
| **OLS(30) ±46**           |  97,911 |  97,929 |  98,262 |  98,298 |  98,307 | **98,141** |
| linear-regr               | 102,300 | 100,750 | 104,532 | 101,680 | 101,618 | 102,176 |
| mse                       | 100,392 | 101,727 |  97,989 | 100,659 | 104,130 | 100,979 |
| nic                       | 111,200 |  98,400 | 106,400 |  88,800 |  97,600 | 100,480 |
| big-range                 |  49,992 |  49,996 |  49,996 |  49,996 |  49,996 |  49,995 |
| correlation-coef          |  95,076 |  94,392 |  93,442 |  92,226 |  93,670 |  93,761 |

**Data 5**

| predictor          |    f1   |    f2   |    f3   |    f4   |    f5   |  mean  |
|--------------------|--------:|--------:|--------:|--------:|--------:|-------:|
| current (running-OLS ±20) | 102,140 |  98,840 | 101,960 | 100,840 |  99,820 | **100,720** |
| **Theil-Sen ±46**         |  95,418 |  95,031 |  95,346 |  94,959 |  95,112 | **95,173** |
| **OLS(30) ±46**           |  94,761 |  93,852 |  94,527 |  94,194 |  93,924 | **94,251** |
| linear-regr               | 100,006 |  99,696 |  99,448 |  97,216 | 101,370 |  99,547 |
| mse                       |  92,115 | 107,067 |  96,921 |  94,785 | 100,392 |  98,256 |
| nic                       |  93,600 |  88,000 | 108,000 | 102,400 |  98,400 |  98,080 |
| big-range                 |  49,140 |  49,168 |  49,280 |  49,152 |  49,172 |  49,182 |
| correlation-coef          |  89,566 |  86,830 |  89,718 |  89,376 |  88,654 |  88,828 |

### File-wins vs each opponent

| opponent            | Theil-Sen ±46 D4 | D5  | OLS(30) ±46 D4 | D5  |
|---------------------|---:|---:|---:|---:|
| current             | 0/5 | 0/5 | 0/5 | 0/5 |
| linear-regr         | 0/5 | 0/5 | 0/5 | 0/5 |
| mse                 | 0/5 | 2/5 | 1/5 | 1/5 |
| nic                 | 1/5 | 2/5 | 2/5 | 2/5 |
| big-range           | 5/5 | 5/5 | 5/5 | 5/5 |
| correlation-coef    | 5/5 | 5/5 | 5/5 | 5/5 |

**Hard-opponent file-wins:** Theil-Sen ±46 = **5/30** · OLS(30) ±46 = **6/30**
— catastrophic, 17 files worse than the 22/30 baseline. Loses 0/10 vs both
`current` and `linear-regr`.

### Read

- **Mean delta vs current:** Theil-Sen ±46 = −6,348 (D4), −5,547 (D5)
  (≈ −6%); OLS(30) ±46 = −4,291 (D4), −6,469 (D5) (≈ −5%). The two centres
  are within ~2k of each other again — centre choice still irrelevant; only
  the width matters.
- **Beats only the trivial opponents** (`big-range`, `correlation-coef`).
  Loses every file vs `current` and `linear-regr` on both datasets — audit
  would fail outright.
- **Confirms the concave landscape.** The width sweep now reads:

  | half-width | TS mean (D4) | TS mean (D5) | TS hard wins |
  |---:|---:|---:|---:|
  | ±15 | 101,790 |  99,678 | 16/30 |
  | ±20 | 102,764 | 100,975 | **22/30** |
  | ±25 | 101,555 | 100,293 | 18/30 |
  | ±30 |  97,848 |  96,701 |  6/30 |
  | ±46 |  96,084 |  95,173 |  5/30 |

  Monotone descent away from ±20 in both directions, with the wide side
  asymptoting toward the `big-range` floor (~50k at ±5000). ±20 is the global
  optimum on this scoring formula and these datasets.
- **Why ±46 isn't a saving floor.** Width 92 → per-hit = `10⁷/93/(N−1)` ≈
  43% of ±20's per-hit. Even at 100% hit rate (which it doesn't reach: D4
  ~76%, D5 ~52%) the ceiling is ~95-97k on D4 and ~93-95k on D5, exactly
  what the measurements show.

---

## 6. Corrigendum — the ±46 result with the correct centre

§§ 1-5 above use a **window-30 OLS / Theil-Sen** centre. The current
production predictor uses **running OLS over the full prefix** (see
[linearstats/predict.go](../linearstats/predict.go) — `Predictor.Next` at
`seen >= 30` evaluates `m*seen + b` from the full-history fit). At wider
half-widths these two centres are NOT equivalent, and the difference is what
made §5 conclude (wrongly) that ±46 was catastrophic.

### Why §5 was wrong: window-N centre variance

For a window of size N predicting one step past the last point, OLS prediction
variance at `x = N` given a fit on `x = 0..N-1` is

  σ²_ŷ = σ²_ε · (1/N + (N − x̄)² / Σ(x_i − x̄)²)

With N=30, x̄=14.5, Σ(x_i − x̄)² = 2247.5, the factor is **0.140**. So
σ²_ŷ ≈ 0.14·σ²_ε. The total prediction-error variance is σ²_ε + σ²_ŷ ≈
1.14·σ²_ε, and the error distribution becomes the **sum** of uniform[-50,50]
and a near-Gaussian σ_ŷ ≈ 10.8 — it is no longer uniform.

P(|error| ≤ 46) for that mixture is ~0.76 (not 0.92), and the score at ±46
collapses because per-hit value × hit-rate stays at uniform-bulk values
(~0.92 × 9) only when the centre is essentially noise-free.

The running-OLS prefix centre has variance ~σ²_ε/n → 0 as the prefix grows.
After ~30 samples σ_ŷ is already <2 and falls further; the error stays
uniform[-50,+50] for the entire stream, and the rounding-cliff geometry the
brainstorm identified applies.

### Verifying the noise distribution

Run: `go run _sim/noise_dist.go .resources/guess-it-dockerized/data_sets`

Global OLS fit on each file (slope 1.0000 ± 0.0002, intercept ≈ 148-150 on
D5, ≈ 150 on D4) gives residuals with these properties:

| file | bulk SD | uniform[-50,+50] SD (50/√3) | outliers (|r|>80) | P(|r|≤46) measured | P(|r|≤46) uniform |
|------|--------:|----------------------------:|--------:|--------:|--------:|
| D4/1 | 28.896 | 28.868 | 0 (0.00%) | 0.9198 | 0.9200 |
| D4/2 | 28.946 | 28.868 | 0 (0.00%) | 0.9163 | 0.9200 |
| D4/3 | 28.834 | 28.868 | 0 (0.00%) | 0.9208 | 0.9200 |
| D4/4 | 28.867 | 28.868 | 0 (0.00%) | 0.9194 | 0.9200 |
| D4/5 | 28.950 | 28.868 | 0 (0.00%) | 0.9199 | 0.9200 |
| D5/1 | 28.754 | 28.868 | 114 (0.91%) | 0.9093 | 0.9200 |
| D5/2 | 29.103 | 28.868 | 110 (0.88%) | 0.9078 | 0.9200 |
| D5/3 | 28.921 | 28.868 | 92 (0.74%) | 0.9116 | 0.9200 |
| D5/4 | 29.009 | 28.868 | 109 (0.87%) | 0.9102 | 0.9200 |
| D5/5 | 28.844 | 28.868 | 108 (0.86%) | 0.9156 | 0.9200 |

D4 is **pure uniform[-50,+50]**, zero outliers. D5 is uniform[-50,+50] bulk
with ~0.85% outliers (|r| up to ~650). Bulk SD on every file matches
50/√3 = 28.868 to two decimals. The "residual SD ≈ 32 / ≈ 48" reported in
[predictor_analysis.md](predictor_analysis.md) is **inflated by window-20
centre variance** — under the same global-OLS centre as here, both datasets
share the same bulk noise.

### Width sweep with the production running-OLS centre

Run: `go run _sim/running_ols_sweep.go .resources/guess-it-dockerized/data_sets`

| half | D4 mean | D5 mean | hit % D4 | hit % D5 | per-hit pts | D4-Δ vs ±20 | D5-Δ vs ±20 |
|---:|--------:|--------:|--------:|--------:|--------:|--------:|--------:|
| ±15 | 101,446 | 100,068 | 31.2 | 30.8 | 25.81 → 26 | −990 | −1,552 |
| **±20** | **102,436** | **101,620** | **41.0** | **40.7** | **19.51 → 20** | **0** | **0** |
| ±25 | 101,798 | 101,072 | 50.9 | 50.6 | 15.69 → 16 | −638 | −548 |
| ±27 | 102,813 | 102,111 | 54.8 | 54.4 | 14.55 → 15 | +377 | +491 |
| ±30 |  98,963 |  98,111 | 60.9 | 60.4 | 13.11 → 13 | −3,473 | −3,509 |
| ±41 | 103,618 | 102,518 | 82.9 | 82.0 |  9.64 → 10 | +1,182 | +898 |
| ±45 | 102,285 | 101,106 | 90.9 | 89.9 |  8.79 →  9 |   −151 |   −514 |
| **±46** | **104,502** | **103,246** | **92.9** | **91.8** |  **8.60 → 9** | **+2,066** | **+1,626** |
| ±47 |  94,918 |  93,675 | 94.9 | 93.7 |  8.42 → **8 (cliff)** | **−7,518** | **−7,945** |
| ±50 |  99,403 |  97,848 | 99.4 | 97.9 |  7.92 →  8 | −3,033 | −3,772 |

The score landscape is **not** concave — it's a sawtooth set by the integer
rounding of `10⁷ / (2h+1) / (N-1)`. Local maxima occur at the largest h that
still rounds to a given per-hit integer: ±20 (pts 20), ±27 (pts 15), ±41
(pts 10), **±46 (pts 9)**. ±46 is the global max because the uniform bulk
saturates the hit rate at ~0.92 there, and the next cliff (pts → 8) eats the
extra coverage.

### File-wins vs each audit opponent at ±46

| opponent          | D4 wins | D5 wins | total |
|-------------------|--------:|--------:|------:|
| current (linear-v2 ±20)  | 5/5 | 5/5 | **10/10** |
| linear-regr       | 4/5 | 5/5 | 9/10 |
| mse               | 5/5 | 4/5 | 9/10 |
| **nic**           | 3/5 | **4/5** | **7/10** ← +1 vs current |
| big-range         | 5/5 | 5/5 | 10/10 |
| correlation-coef  | 5/5 | 5/5 | 10/10 |

**Hard-opponent file-wins (linear-regr + mse + nic): 25/30** — up from 22/30
at the ±20 baseline. The decisive gain is `nic` on Data 5: ±46 scores 103,131
on D5/2 (beats nic's 88,000), 103,266 on D5/3 (loses to 108,000), 103,077 on
D5/4 (beats 102,400 — new win), 103,635 on D5/5 (beats 98,400). The
`nic` audit pass probability climbs from ~65% to ~90%.

### Why the prior analysis missed this

[predictor_analysis.md](predictor_analysis.md) §4 ran an AR(1) correction and
concluded "the score landscape is flat in 95k-103k." That conclusion is
correct *under window-20 centred fits*, because window-20 centre variance
broadens the residual distribution exactly enough to remove the sawtooth.
With the running-OLS-prefix centre (which the linear-v2 PR introduced), the
underlying uniform[-50,+50] structure becomes visible and the sawtooth
returns.

The Theil-Sen analysis in §§ 1-5 is unaffected by this correction *for the
specific widths it tested at* (±15, ±20, ±25, ±30, ±46 with window-30
centre); it correctly showed the centre choice doesn't matter at those
widths under that centre. But the width comparison itself was invalidated:
under the production centre, ±46 wins clearly.

### Recommendation

**Change one constant in [linearstats/predict.go](../linearstats/predict.go):
`fixedRange = 20.0` → `fixedRange = 46.0`.** No other code change. Expected
audit impact:

- D4 mean score: 102,436 → 104,502 (+2.0%)
- D5 mean score: 101,620 → 103,246 (+1.6%)
- Hard-opponent file-wins: 22/30 → 25/30
- Audit pass probability vs `nic`: ~65% → ~90% (the bottleneck moves)

Two further worthwhile cross-checks before shipping:

1. **Validate ±46 on the Dockerized tester** to confirm `server.js` rounding
   matches the simulator (it should — the formula is the same).
2. **Bench the running-OLS Predictor warm-up phase.** At ±46 the wider band
   absorbs the early-sample centre noise that ±20 would miss; expected to
   help on samples < 30, not hurt. Verify on `main_test.go`'s table cases.

---

## 7. Summary

| variant               | D4 mean | D5 mean | hard wins | audit-safe? |
|-----------------------|--------:|--------:|----------:|:------------|
| current (running-OLS ±20)  | 102,432 | 100,720 | 22/30 | yes (baseline) |
| Theil-Sen ±20              | 102,764 | 100,975 | 22/30 | yes (≈ equivalent, +0.3%) |
| Theil-Sen ±15              | 101,790 |  99,678 | 16/30 | **no** (band too narrow) |
| OLS(30) ±15                | 101,738 |  97,770 | 15/30 | **no** (band too narrow) |
| Theil-Sen ±25              | 101,555 | 100,293 | 18/30 | **no** (band slightly too wide) |
| OLS(30) ±25                | 102,163 |  98,668 | 15/30 | **no** (band slightly too wide) |
| Theil-Sen ±30              |  97,848 |  96,701 |  6/30 | **no** (band too wide) |
| Theil-Sen ±46              |  96,084 |  95,173 |  5/30 | **no** (band far too wide) |
| OLS(30) ±46                |  98,141 |  94,251 |  6/30 | **no** (band far too wide) |
| Theil-Sen + adaptive interval (hybrid) | 92,679 | 90,192 | n/a | **no** (loses ~10%) |

**Recommendation:** Keep the OLS centre (Theil-Sen swap buys nothing — see §§
1-5), but **change `fixedRange` from 20 to 46**. See §6 for the rounding-cliff
analysis. Expected lift: D4 +2.0%, D5 +1.6%, hard-opponent file-wins
22/30 → 25/30, audit pass vs `nic` ~65% → ~90%.

§5's claim that ±46 was catastrophic was an artefact of the window-30
centre used in that benchmark; the production centre is full-prefix
running OLS, which preserves the uniform[-50,+50] residual structure that
makes ±46 the global score-landscape maximum.

<a id="footnote-hybrid"></a>
*The hybrid (Theil-Sen centre + empirical-error interval search) was
benchmarked separately via [_sim/theilsen_eval.go](../_sim/theilsen_eval.go);
hit rate collapses to ~2.3% as the interval search overfits white-noise
residuals to width ≈ 1.5. Same failure mode as the OLS-based adaptive
interval evaluated in `_sim/adaptive_eval.go`, ~3% worse.*
