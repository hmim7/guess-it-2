# Edge Cases: Guess It 2

How the **fixed-range split predictor** behaves under streaming and statistical
edge conditions. The predictor centres a constant ±20 range on the
linear-regression extrapolation while `seen < 1000`, and on the window median
afterwards — so most edge cases reduce to "where is the centre, and is the
fixed width enough to land."

---

## 1. Input Stream & Data Edge Cases

| ID | Scenario | Description | Expected Behavior |
|----|----------|-------------|-------------------|
| 01 | **First Input** | The very first number; window holds one value. | `LinearRegression([v])` returns `(0, v)` → centre = `v` → output `v-20 v+20`. No zero-width range. |
| 02 | **Insufficient Data** | Fewer points than the window size. | Regression fits whatever points exist; range is still centre ± 20. |
| 03 | **Sudden Trend Change / Spike** | Stream jumps to a new level (`10,11,10,500,...`). | The window slides; within `N` inputs the old values are evicted and the regression centre tracks the new level. |
| 04 | **Consistent Values** | A long run of the same number. | Regression slope ≈ 0, centre = the value → output `value-20 value+20`. |
| 05 | **Oscillating Values** | Alternating values (`10,100,10,100,...`). | Centre ≈ window mean (~55); the fixed ±20 band misses both extremes — an accepted trade (score is maximised by the narrow width, not by hitting). |
| 06 | **Non-Numeric Input** | Garbage line in the stream. | `strconv.ParseFloat` fails → the line is skipped; window and `seen` unchanged. |
| 07 | **Empty or Closed Input Stream** | No data / EOF. | `bufio.Scanner` loop exits; `run` returns nil; exit 0. |
| 12 | **Trend Reversal** | `110,120,130,120,110,...`. | Window eviction lets the regression slope flip sign within `N` inputs. |
| 13 | **Noisy Dataset** | High variance, no trend. | Regression centre ≈ mean; the fixed ±20 lands only on the central core — accepted by design. |

---

## 2. Algorithm, Performance & Output Edge Cases

| ID | Scenario | Description | Expected Behavior |
|----|----------|-------------|-------------------|
| 08 | **Constant Window** | All window values identical. | Regression `(m=0, b=v)` → centre = `v`; median = `v`. Output `v-20 v+20`. |
| 09 | **Narrow-Range Pressure** | A very stable stream tempts a tiny range. | The width is a fixed constant 40; it never collapses. `minWidth = 1` is a further floor after rounding. |
| 10 | **Wide-Range Pressure** | A volatile stream tempts a huge range. | The width is fixed at 40 regardless of volatility — by design, to keep score-per-hit high. |
| 11 | **Large Magnitude Numbers** | Values `> 1e18`. | `float64` arithmetic stays stable; the ±20 band is added without overflow. |
| 14 | **Phase Transition at `seen = 1000`** | The 1000th input arrives. | Centring switches from regression extrapolation to window median; width is unchanged. |
| 15 | **Rapid Input / Performance** | Quick tester mode. | Per-iteration work is `O(N)` over the 20-value window; per-line flush, no lag. |
| 16 | **Non-Integer Prediction** | Regression/median centre is fractional. | `roundBounds` applies `math.Round` to both bounds before printing. |

---

## 3. Detailed Scenarios

#### Case 01: First Input
- **Input:** `189`
- **Behaviour:** `LinearRegression([189]) = (0, 189)`; centre = `0·1 + 189 = 189`; output `169 209`.

#### Case 03: Sudden Trend Change / Spike
- **Input:** `..., 50, 52, 49, 51, 1000, 1002, 998, 1001, ...`
- **Behaviour:** Once the window (`N = 20`) is filled with values near 1000, the regression centre stabilises around the new level; the transient predictions during the slide are accepted misses.

#### Case 05: Oscillating Values
- **Input:** `..., 20, 150, 20, 150, ...`
- **Behaviour:** The regression centre sits near the mean (~85). A fixed ±20 cannot cover both 20 and 150 — the predictor accepts the misses because widening to cover them would cost more score than the extra hits return.

#### Case 14: Phase Transition at `seen = 1000`
- **Behaviour:** For `seen` in `1..999` the centre is `m·len(window) + b`; from `seen = 1000` onward it is `Median(window)`. The dataset analysis showed the median centre scores higher in steady state.

#### Case 16: Non-Integer Prediction
- **Input:** centre computed as e.g. `150.7`
- **Behaviour:** `roundBounds(130.7, 170.7)` → `131 171`.
