# Golden Test Suite: guess-it-2

Mandatory test cases verifying the **fixed-range split predictor**: a constant
±20 range centred on the linear-regression extrapolation while `seen < 1000`,
and on the window median afterwards. Every predicted range spans exactly
40 units.

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

`Predict(window, seen, current)` — every range spans `2 × fixedRange = 40`.

| ID   | Phase       | Window / `seen`                                  | Centre        | Expected Output |
|------|-------------|--------------------------------------------------|---------------|-----------------|
| GT01 | regression  | `[250]`, `seen = 1`                              | `0·1 + 250` = 250 | `230 270` |
| GT02 | regression  | `[100,101,102,103,104]`, `seen = 5`              | `m=1,b=100` → `1·5+100` = 105 | `85 125` |
| GT03 | regression  | `[104,103,102,101,100]`, `seen = 5`              | `m=-1,b=104` → `-1·5+104` = 99 | `79 119` |
| GT04 | regression  | `[50,50,50,50,50]`, `seen = 10`                  | `m=0,b=50` → 50 | `30 70` |
| GT05 | median      | `[10,20,30,40,50]`, `seen ≥ 1000`                | median = 30   | `10 50` |
| GT06 | median      | `[10,20,30,40]`, `seen ≥ 1000`                   | median = 25   | `5 45` |
| GT07 | median      | `[100,100,100,100,9000]`, `seen ≥ 1000`          | median = 100  | `80 120` |

---

## 3. Edge-Case Behaviour — see `docs/edge_cases.md`.

| ID | Scenario | Expected Output |
|:--:|---|---|
| 01 | First input | `current ± 20` (regression on a one-value window). |
| 06 | Non-numeric line | Skipped; no output line. |
| 07 | EOF / empty stream | Exit 0, no output. |
| 14 | `seen` reaches 1000 | Centre switches regression → median; width unchanged. |
| 16 | Fractional centre | `roundBounds` rounds both bounds to integers. |

---

## 4. Dataset Analysis (rationale)

Simulated over `docs/data-sets/` (5 groups × 5 files, ~12,500 numbers each),
window = 20:

- Under a `score = 1/width` per-hit model, a fixed ±20 range out-scored an
  adaptive `c·meanStep` range on every dataset family.
- With a constant width, score tracks hit rate; the window **median** centred a
  higher-scoring range than regression in steady state (e.g. group 3: 40.9 % vs
  32.3 % hit). The `seen < 1000` regression phase costs <0.5 % since it covers
  <8 % of a file, and keeps the `linear-stats` calculation in use.

### Output Format
- Two space-separated integers per line: `fmt.Fprintf(w, "%d %d\n", lower, upper)`.
- Every range spans 40 units; `minWidth = 1` is a post-rounding floor.
