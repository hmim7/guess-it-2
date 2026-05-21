# Task 09 — Fixed-Range Split Predictor

## Goal
Replace the hybrid predictor (regression/median blend + 5-tier StdDev multiplier
+ jump-guard) with a simple, score-focused fixed-range model, validated against
the datasets in `docs/data-sets/`.

## Dataset analysis (driving the design)
Five dataset groups (1, 2, 3, 5, 9), 5 files each, ~12,500 numbers per file.
Simulated with window = 20:

- Under a `score = 1/width` per-hit model, **fixed ±20 out-scored an adaptive
  `c·meanStep` range on every dataset** — the narrow width's score-per-hit
  outweighs the lower hit rate.
- With a constant ±20 width, score is proportional to hit rate. Centering hit %:

  | file | regression-only | median-only | split (reg n<1000 / median after) |
  |------|----------------:|------------:|----------------------------------:|
  | 1/1  | 39.2 | 40.2 | 40.1 |
  | 3/1  | 32.3 | 40.9 | 40.3 |
  | 5/1  | 38.6 | 40.4 | 40.3 |
  | 9/1  |  8.3 |  9.1 |  9.1 |
  | 9/4  |  7.7 |  7.8 |  7.7 |

  Median centering wins/ties everywhere; the split equals median-only (the first
  1000 points are <8 % of a file) while keeping linear regression in use.

## Algorithm
```
Predict(window, seen, current) -> (low, high):
    if seen < regressionPhaseLimit:        # 1000
        m, b   := LinearRegression(window)
        center := m * float64(len(window)) + b
    else:
        center := Median(window)
    return roundBounds(center - fixedRange, center + fixedRange)   # fixedRange = 20
```

## Steps
1. Rewrite `linearstats/predict.go`: split-model `Predict`; delete
   `dynamicMultiplier`, the 5 var-thresholds, the 5 multipliers, `minStdRange`,
   `initialRange`; add `fixedRange = 20`, `regressionPhaseLimit = 1000`; keep
   `WindowSize`, `minWidth`, `roundBounds`.
2. `stats.go` is untouched — all six functions stay.
3. Rewrite the predictor tests in `linearstats/linearstats_test.go`:
   remove `TestDynamicMultiplier` and the hybrid tests; add table-driven cases
   for the regression phase (`seen < 1000`), the median phase (`seen >= 1000`),
   the ±20 width, and rounding. Keep all stats tests + `TestRoundBounds`.

## Acceptance Criteria
- `Predict` has the signature `Predict(window []float64, seen int, current float64) (int64, int64)`.
- `go vet ./...` clean; `go build ./...` succeeds (after Task 10 wires `main.go`).
- `go test ./linearstats/...` green; coverage ≥ 90 %.
- No `dynamicMultiplier` / tier constants remain in `predict.go`.

## Validation
```sh
go vet ./linearstats/...
go test ./linearstats/... -v -cover
```

## Conventional commit
`refactor(linearstats): replace hybrid model with fixed-range split predictor`
