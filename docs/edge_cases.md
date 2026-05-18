# Edge Cases: Guess It 2

This document covers the streaming and statistical edge cases the hybrid predictor must handle. Sections 1–3 mirror the original `guess-it-1` cases; section 4 adds the new correlation-driven edges introduced by the linear-stats integration.

---

## 1. Input Stream & Data Edge Cases

| Case ID | Scenario | Description | Expected Behavior |
|---------|----------|-------------|-------------------|
| 01 | **First Input** | The program receives its very first number. No history exists for StdDev. | Must return a sensible, safe default range (e.g. `num ± 60`) to avoid a zero-width prediction. |
| 02 | **Insufficient Data** | Number of data points < initial-phase threshold. | Provide a wide, safe prediction until the model stabilizes. |
| 03 | **Sudden Trend Change / Spike** | A stream jumps to a different level (`10,11,10,500,...`). | Sliding window forgets the old trend within `N` inputs; range shifts to the new level. |
| 04 | **Consistent Values** | A long stream of the same number. | Range narrows but respects `minStdRange` so it never collapses to a single point. |
| 05 | **Oscillating Values** | Alternating values (`10,100,10,100,...`). | Range stays wide enough to cover both extremes; `|r|` ≈ 0 keeps the hybrid in median-mode. |
| 06 | **Non-Numeric Input** | Garbage line in the stream. | Ignore the line; preserve state; continue. |
| 07 | **Empty or Closed Input Stream** | No data / EOF. | Wait for input or exit 0 on EOF. |
| 12 | **Trend Reversal** | `110,120,130,120,110,...`. | Adjust quickly to the reversal; do not stay biased upward. |
| 13 | **Noisy Dataset** | High variance, no clear trend. | Multiplier escalates (Defensive / Extreme tier); `|r|` low so no regression bias. |

---

## 2. Algorithm, Performance & Prediction Edge Cases

| Case ID | Scenario | Description | Expected Behavior |
|---------|----------|-------------|-------------------|
| 08 | **Zero Standard Deviation** | All window values identical. | `margin = max(multiplier × 0, minStdRange)` floor prevents zero-width range. |
| 09 | **Too-Narrow Range Syndrome** | Very stable trend. | `minWidth = 1` enforced after rounding. |
| 10 | **Too-Wide Range Syndrome** | High variance trend. | 5-tier multiplier caps width at a competitive band. |
| 11 | **Large Magnitude Numbers** | `> 1e18`. | `float64` arithmetic remains stable. |
| 14 | **Aggressive Start Strategy** | First 5–10 inputs. | Wide initial range, tightens after `seen ≥ 5`. |
| 15 | **Rapid Input / Performance** | Quick tester mode. | All computations O(N) over the window; per-line flush, no lag. |
| 16 | **Non-Integer Prediction** | Fractional bounds. | `math.Round` applied before printing. |

---

## 3. Detailed Scenarios

#### Case 01: First Input / Aggressive Start (Case 14)
- **Input:** `189`
- **Expected Behavior:** Use the wide initial range (`current ± 60`, scaled by magnitude) until enough data is collected for the statistical model.

#### Case 03: Sudden Trend Change / Spike
- **Input:** `..., 50, 52, 49, 51, 1000, 1002, 998, 1001, ...`
- **Expected Behavior:** Sliding window forgets the 50s; centre stabilizes around the new mean.

#### Case 05: Oscillating Values
- **Input:** `..., 20, 150, 20, 150, ...`
- **Expected Behavior:** High StdDev forces a wide range that covers both extremes. `PearsonCorrelation ≈ 0` keeps the regression contribution near zero, so the centre stays at the median (~85).

#### Case 08: Zero Standard Deviation
- **Input:** `..., 50, 50, 50, 50, 50`
- **Expected Behavior:** `minStdRange = 1.8` produces a small but non-zero margin.

#### Case 12: Trend Reversal
- **Input:** `..., 110, 120, 130, 140, 130, 120, 110, ...`
- **Expected Behavior:** The window evicts the old peak; `m` flips sign; regression centre adjusts downward.

#### Case 13: Noisy Dataset
- **Input:** `..., 100, 5, 120, 8, 95, 15, ...`
- **Expected Behavior:** High StdDev → Defensive / Extreme tier; low `|r|` → no regression bias.

---

## 4. Hybrid Predictor Edges (new in guess-it-2)

| Case ID | Scenario | Description | Expected Behavior |
|---------|----------|-------------|-------------------|
| **17** | **Low `\|r\|` Fallback** | Window has no clear trend (`\|r\| ≈ 0`). | Centre weight on regression collapses to ≈ 0; predictor behaves like the guess-it-1 median + 5-tier model. |
| **18** | **High `\|r\|` Tightening** | Window has a clean linear trend (`\|r\| ≈ 1`). | Regression centre dominates the blend; multiplier scaled by `1 − 0.25·\|r\|` produces a tighter band. |
| **19** | **Mid `\|r\|` Blend** | Moderate correlation (e.g. `\|r\| ≈ 0.6`). | Centre is a true blend of regression and median; margin shrinks proportionally; no abrupt regime changes. |
| **20** | **Constant Window with `\|r\| = 0`** | All values equal → `PearsonCorrelation` returns `0` (constant `y` guard). | Blend collapses to median = the constant value; `minStdRange` floor applies. |

#### Case 17: Low `|r|` Fallback
- **Input:** `..., 10, 100, 10, 100, 10, 100, ...`
- **Expected Behavior:** `PearsonCorrelation` ≈ 0 → `w` ≈ 0 → centre = median; multiplier untouched by correlation scaling.

#### Case 18: High `|r|` Tightening
- **Input:** `..., 100, 101, 102, 103, 104, 105, ...`
- **Expected Behavior:** `m ≈ 1`, `|r|` ≈ 1, `regCenter ≈ next y`. The blended centre matches the regression extrapolation; the multiplier shrinks by ~25 %.

#### Case 20: Constant Window with `|r| = 0`
- **Input:** `..., 80, 80, 80, 80, 80, ...`
- **Expected Behavior:** `PearsonCorrelation` returns 0 for constant `y` (denominator guard). Centre = 80; margin floored to `minStdRange`.
