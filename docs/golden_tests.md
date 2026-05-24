# Golden Test Suite: guess-it-2

Mandatory test cases verifying the **running-OLS fixed-range predictor**: a
constant ±46 range centred on the running OLS extrapolation once `seen ≥ 30`,
with a warm-up fallback to the original window-regression / window-median split
below that. Every predicted range spans exactly 92 units.

## 1. Audit Cases

Verification notes:
- Correct per-opponent win rate across both datasets.
- Consistent performance across 3 runs per dataset.
- Robust execution via `script.sh` and the `student/` layout.
- Deterministic stdin/stdout usage in the tester environment.

| ID | Description           | Input Argument                              | Expected Behavior / Output |
|:--:|-----------------------|---------------------------------------------|----------------------------|
| 01 | `big-range` Opponent  | `?guesser=big-range` with Data 4, Data 5    | Win ≥ 2/3 runs per dataset. |
| 02 | `linear-regr`         | `?guesser=linear-regr` with Data 4, Data 5  | Win ≥ 2/3 runs per dataset. |
| 03 | `correlation-coef`    | `?guesser=correlation-coef` with Data 4, 5  | Win ≥ 2/3 runs per dataset. |
| 04 | `mse` (Bonus)         | `?guesser=mse` with Data 4, Data 5          | Win ≥ 2/3 runs per dataset. |
| 05 | `nic` (Bonus)         | `?guesser=nic` with Data 4, Data 5          | Win ≥ 2/3 runs per dataset. |
| 06 | Packaging & Execution | `script.sh` + `student/` at tester root     | Script runs the binary cleanly. |

---

## 2. Predictor Golden Tests

`Predict(window, seen, current)` — every range spans `2 × fixedRange = 92`.

| ID   | Phase       | Window / `seen`                                  | Centre        | Expected Output |
|------|-------------|--------------------------------------------------|---------------|-----------------|
| GT01 | regression  | `[250]`, `seen = 1`                              | `0·1 + 250` = 250 | `204 296` |
| GT02 | regression  | `[100,101,102,103,104]`, `seen = 5`              | `m=1,b=100` → `1·5+100` = 105 | `59 151` |
| GT03 | regression  | `[104,103,102,101,100]`, `seen = 5`              | `m=-1,b=104` → `-1·5+104` = 99 | `53 145` |
| GT04 | regression  | `[50,50,50,50,50]`, `seen = 10`                  | `m=0,b=50` → 50 | `4 96` |
| GT05 | median      | `[10,20,30,40,50]`, `seen ≥ 1000`                | median = 30   | `-16 76` |
| GT06 | median      | `[10,20,30,40]`, `seen ≥ 1000`                   | median = 25   | `-21 71` |
| GT07 | median      | `[100,100,100,100,9000]`, `seen ≥ 1000`          | median = 100  | `54 146` |

---

## 3. Edge-Case Behaviour — see `docs/edge_cases.md`.

| ID | Scenario | Expected Output |
|:--:|---|---|
| 01 | First input | `current ± 46` (warm-up fallback runs the regression on a one-value window since `seen < 30`). |
| 06 | Non-numeric line | Skipped; no output line. |
| 07 | EOF / empty stream | Exit 0, no output. |
| 14 | `seen` reaches 1000 | Centre switches regression → median; width unchanged. |
| 16 | Fractional centre | `roundBounds` rounds both bounds to integers. |

---

## 4. Dataset Analysis (rationale)

Simulated over the dockerized audit datasets Data 4 / Data 5 (5 files × 12,500
numbers each), and the residual distribution characterised by
[_sim/noise_dist.go](../_sim/noise_dist.go):

- Residuals are **uniform on [-50, +50]** (D4 zero outliers; D5 ~0.85%
  outliers up to |r|≈650). Bulk SD = 28.87 = 50/√3 on every file.
- Under integer-rounded `round(10⁷/(2h+1)/(N−1))` per-hit scoring, the score
  landscape is a sawtooth with local maxima at the largest half-width still
  rounding to each per-hit integer: ±20→20pts, ±27→15pts, ±41→10pts,
  **±46→9pts**.
- ±46 is the global maximum: hit rate saturates at ~92% inside [−46, +46],
  the next step ±47 crosses a per-hit cliff (9→8 pts) that costs ~9% of the
  score, and ±45 is below the same per-hit ceiling at lower coverage.

### Output Format
- Two space-separated integers per line: `fmt.Fprintf(w, "%d %d\n", lower, upper)`.
- Every range spans 92 units; `minWidth = 1` is a post-rounding floor.
