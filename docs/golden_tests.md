# Golden Test Suite: guess-it-2

This document records the mandatory test cases used to verify the functional requirements of the `guess-it-2` prediction engine, ensuring robust handling of streaming input, opponent benchmarking, and fault tolerance.

## 1. Audit Cases

Verification Notes:
- Correct per-opponent win rate across both datasets.
- Consistent performance across 3 runs per dataset.
- Robust execution via `script.sh` and `student/` layout.
- Deterministic stdin/stdout usage in the tester environment.

| ID | Description           | Input Argument                              | Expected Behavior / Output                |
|:--:|-----------------------|---------------------------------------------|-------------------------------------------|
| 01 | `big-range` Opponent  | `?guesser=big-range` with Data 4, Data 5    | Win ≥ 2/3 runs per dataset.               |
| 02 | `linear-regr`         | `?guesser=linear-regr` with Data 4, Data 5  | Win ≥ 2/3 runs per dataset.               |
| 03 | `correlation-coef`    | `?guesser=correlation-coef` with Data 4, 5  | Win ≥ 2/3 runs per dataset.               |
| 04 | `mse` (Bonus)         | `?guesser=mse` with Data 4, Data 5          | Win ≥ 2/3 runs per dataset.               |
| 05 | `nic` (Bonus)         | `?guesser=nic` with Data 4, Data 5          | Win ≥ 2/3 runs per dataset.               |
| 06 | Packaging & Execution | `script.sh` + `student/` at tester root     | Script runs the binary cleanly.           |

---

## 2. Edge Cases

| ID | Description | Input Argument | Expected Behavior / Output |
|:--:|---|---|---|
| 01 | First Input | First number in stream | Safe default range (`num ± 60`). |
| 02 | Insufficient Data | Fewer points than initial-phase threshold | Wider, safer prediction. |
| 03 | Sudden Trend Change / Spike | `10,11,10,500,502,...` | Range shifts to the new level after window fills. |
| 04 | Consistent Values | `50,50,50,...` | Range tightens but respects `minStdRange`. |
| 05 | Oscillating Values | `10,100,10,100,...` | Range wide enough; regression weight ≈ 0. |
| 06 | Non-Numeric Input | Garbage line | Ignored, state preserved. |
| 07 | Empty / Closed Stream | EOF | Exit 0 cleanly. |
| 08 | Zero StdDev | All values equal | `minStdRange` floor. |
| 09 | Too-Narrow Range | Stable trend | `minWidth` enforced. |
| 10 | Too-Wide Range | High variance | 5-tier cap. |
| 11 | Large Magnitude Numbers | `> 1e18` | `float64` stable. |
| 12 | Trend Reversal | `110,120,130,120,110` | Window evicts old trend; regression sign flips. |
| 13 | Noisy Dataset | `10,500,20,480` | Defensive multiplier wins. |
| 14 | Aggressive Start | First 5–10 inputs | Wide, then tightens. |
| 15 | Rapid Input | Quick tester mode | O(N) per iteration. |
| 16 | Non-Integer Prediction | Fractional bounds | `math.Round` applied. |
| **17** | **Low `\|r\|` fallback** | Oscillating window | Centre = median; multiplier untouched. |
| **18** | **High `\|r\|` tightening** | Clean trend | Centre tracks `m·n+b`; multiplier shrinks by `0.25·\|r\|`. |
| **19** | **Mid `\|r\|` blend** | `\|r\| ≈ 0.6` | Centre is a true blend; gradual response. |
| **20** | **Constant window, `r = 0`** | All equal | Centre = constant; `minStdRange` floor. |

---

## 3. Detailed Hybrid-Predictor Scenarios

### Case 17: Low `|r|` Fallback
- **Input:** alternating `[10,100,10,100,...]`
- **Mechanism:** `PearsonCorrelation` ≈ 0 → blend weight `w ≈ 0` → `center = median(window)`. Multiplier scale `1 - 0.25·w ≈ 1`. Predictor reduces to the guess-it-1 model.

### Case 18: High `|r|` Tightening
- **Input:** `[100,101,102,103,104,105,...]`
- **Mechanism:** `m ≈ 1`, `r ≈ 1`, `regCenter = m·n + b`. Centre tracks the regression extrapolation; multiplier reduced by ~25 % → tighter band, higher score.

### Case 19: Mid `|r|` Blend
- **Input:** Trending data with noise (e.g. `100, 102, 99, 105, 103, 108, 106, 110, ...`)
- **Mechanism:** `w` between 0.4 and 0.8 — centre is a smooth blend of regression and median; margin shrinks proportionally to `w`.

### Case 20: Constant Window with `r = 0`
- **Input:** `80,80,80,80,80,...`
- **Mechanism:** `PearsonCorrelation` returns 0 (constant-y guard). Blend collapses to median = 80. StdDev = 0 → `margin = minStdRange = 1.8` → output `78 82` after rounding.

---

## 4. Strategic Recommendations

### Hybrid Predictor
```
m, b   = LinearRegression(window)
r      = PearsonCorrelation(window)
regC   = m*len(window) + b
robC   = median(window)            # or average if PREDICT_CENTER != "median"
w      = |r|
center = w*regC + (1-w)*robC
std    = stdDev(window)
mul    = dynamicMultiplier(std) * (1 - 0.25*w)
margin = max(mul*std, minStdRange)
if |current - center| > 1.5*std: margin = max(margin, |current-center|*1.1)
print roundBounds(center-margin, center+margin)
```

### Output Format Reminder
- Smaller range = Higher score.
- Minimum possible range: 1 unit (e.g. `150 151`).
- Format: `fmt.Printf("%d %d\n", lower, upper)`.
