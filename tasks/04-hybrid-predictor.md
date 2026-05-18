# Task 04 — Hybrid Regression + Tier Predictor

## Goal
Upgrade `linearstats/predict.go` to a hybrid model that blends linear-regression extrapolation with the existing 5-tier StdDev predictor, weighted by `|PearsonCorrelation|`.

## Algorithm
For each prediction (window has at least the initial-phase threshold):

1. `m, b := LinearRegression(window)`; `r := PearsonCorrelation(window)`.
2. `regCenter := m * float64(len(window)) + b`.
3. `robustCenter := median(window)` (or `average` if `PREDICT_CENTER != "median"`).
4. `w := math.Abs(r)`; `center := w*regCenter + (1-w)*robustCenter`.
5. Compute the tier multiplier from `stdDev(window)` (existing 5 tiers).
6. Scale: `multiplier *= 1 - 0.25*w`. Floor by `minStdRange`. Apply the jump-guard if `|current-center| > 1.5*std`.
7. `roundBounds(center - margin, center + margin)` with `minWidth = 1`.

## Steps
1. Implement the new `Predict()`; keep all existing tuning constants and safeguards.
2. Extract `dynamicMultiplier(std float64) float64` for readability and direct testing.
3. Add table-driven cases in `linearstats_test.go`:
   - trended (positive slope) → centre near `m*n+b`, tighter margin.
   - constant values → centre = value, margin = `minStdRange` floor.
   - oscillating values (high std, low |r|) → defensive multiplier wins, no regression bias.
   - sudden spike (jump-guard triggers).
   - high-|r| tightening (compare margin with and without correlation scaling).
5. Run `go fmt ./...` and `go vet ./...`.

## Acceptance Criteria
- `go test ./linearstats/... -cover` ≥ 90 %.
- New tests cover trended, constant, oscillating, spike, low-|r|, high-|r| cases.
- All previously passing predictor tests still pass.
- `go vet ./...` clean.

## Validation
```sh
go vet ./...
go test ./linearstats/... -v -cover
```

## Conventional commit
`feat(linearstats): hybrid regression + tier predictor`
