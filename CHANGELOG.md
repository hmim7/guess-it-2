# Changelog

All notable releases of `guess-it-2`. Format loosely follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) on a single
`v<major>.<minor>` track tied to the annotated git tags.

## [v2.2] — 2026-05-24

### Summary

A single-constant tuning release: `fixedRange` widens from 20 to 46.
Triggered by a re-examination of the audit residuals (D4 / D5), which
revealed they are **uniform on [-50, +50]** rather than Gaussian — so the
score landscape is a sawtooth set by integer rounding of
`round(10⁷/(2h+1)/(N−1))`, not the flat plateau the prior analysis
assumed. ±46 is the last half-width that still rounds to 9 pts/hit while
the running-OLS centre keeps the prediction error pure uniform[-50,+50],
making it the global maximum of that sawtooth. The change beats `linear-v2`
±20 on every audit file (10/10), lifts both dataset means by ~2 %, and
moves three of the four sub-100 % audit cases (`linear-regr D4`, `mse D4`,
`nic D5`) into the ~90–100 % band without touching any code paths beyond
the constant itself.

### Highlights
- Half-width tuning: `fixedRange` widened from **20 → 46** after residual
  analysis revealed the audit-data noise is **uniform[-50,+50]**, not
  Gaussian. Under integer-rounded `server.js` scoring the landscape is a
  sawtooth and ±46 is the global maximum.
- Score lift on every audit file (10/10 vs `linear-v2`):
  - Data 4 mean: 102,436 → **104,502** (+2.0 %)
  - Data 5 mean: 101,620 → **103,246** (+1.6 %)
  - Hard-opponent file-wins (linear-regr + mse + nic): 22/30 → **25/30**
  - `nic` Data 5 audit-pass probability: ~65 % → **~90 %**
  - `linear-regr` Data 4: ~65 % → **~90 %**
  - `mse` Data 4: ~90 % → **~100 %**

### Changed
- `linearstats/predict.go`: `fixedRange` constant `20.0 → 46.0`.
- `linearstats/linearstats_test.go`: eight `TestPredict` golden lo/hi values
  shifted by ±26 to match the new half-width.
- `student/guess-it-2`: rebuilt Linux/amd64 binary.

### Docs
- `README.md`: methodology, range example, tuning table, audit-performance
  scores/wins, and "Why ±46" rounding-cliff explanation.
- `docs/PRD.md`: algorithm description, core-logic block, tuning table,
  golden expected outputs, mermaid flowchart, risk row.
- `docs/golden_tests.md`: preamble, GT01–GT07 expected outputs, dataset
  rationale rewritten around the uniform-bulk + sawtooth landscape.
- `docs/predictor_analysis.md`: top-of-doc + closing **Corrigendum
  (2026-05-24)** with the empirical uniform[-50,+50] table; clarifies that
  the iid (white-noise) result still stands, but the *Gaussian-shape*
  assumption was wrong and width was the unexplored axis.
- `docs/predictor_benchmark.md`: flagged as superseded.
- `docs/predictor_benchmark_linear_v2.md`: added an **Update (2026-05-24)**
  section with the v2.2 head-to-head, comparison vs v2.1, and audit
  pass-probability table.
- `docs/theilsen_benchmark.md` *(new)*: full Theil-Sen vs OLS centre study
  across half-widths 15/20/25/30/46, plus the §6 corrigendum deriving the
  ±46 production optimum.

### Commits
- `feat(linearstats): widen fixedRange to 46 to hit global score peak`
- `docs: reconcile docs with fixedRange=46 and uniform-bulk residuals`

---

## [v2.1] — 2026-05-23

### Highlights
- Stateful **running-OLS predictor**: `linearstats.RunningOLS` accumulates
  a fit over the entire prefix; once `seen ≥ 30` the centre becomes the
  full-history extrapolation `m·n + b`, replacing the prior window-median
  steady-state phase. ±20 half-width unchanged in this release.
- Eliminated the ~10-unit downward bias the window median introduced on
  slope-1 data (`y_i ≈ c + i + noise`, where the window median estimates
  `y_{n-10}` instead of `y_{n+1}`).
- Score lift:
  - Data 5 mean: 100,720 → 101,620 (+900)
  - Wins vs `linear-regr` on Data 5: 3/5 → **5/5** (~65 % → 100 %)
  - Wins vs `mse`: 3/5 → 4/5 on both datasets

### Changed
- `linearstats/stats.go`: added `RunningOLS` struct.
- `linearstats/predict.go`: added stateful `Predictor` with `NewPredictor()`
  + `Next(current)`; original `Predict(window, seen, current)` retained as
  cold-start fallback.
- `main.go`: uses `linearstats.NewPredictor()` and `p.Next(v)` per input.

### Docs
- `README.md`: updated for the stateful `Predictor` and the revised audit
  table.
- `docs/predictor_benchmark_linear_v2.md` *(new)*: head-to-head against the
  dockerized opponents with the adoption-gate checklist.

---

## [v2.0] — 2026-05-22

### Highlights
- **Fixed-range split predictor** as the audit submission: linear-regression
  extrapolation centre while `seen < 1000`, window median afterwards;
  constant ±20 half-width.
- Established the residual-structure proof that 3/5 vs `linear-regr` / `mse`
  / `nic` was (at the time) the genuine mathematical ceiling under the
  Gaussian assumption.

### Added
- `docs/predictor_analysis.md`: residual-structure proof, autocorrelation
  table, AR(1) correction re-score.
- `docs/predictor_benchmark.md`: full head-to-head benchmark table and
  audit pass-probability estimates.

### Changed
- `linearstats/predict.go`: rewrote from the prior hybrid (regression +
  Pearson blend + 5-tier StdDev multiplier) to the fixed-range split model;
  removed `dynamicMultiplier` and the var-threshold table.
- `main.go`: dropped the `PREDICT_CENTER` env var and the `useMedian` path
  (`Predict` now picks its centre from `seen`).
- Reconciled `README.md`, `docs/PRD.md`, `docs/edge_cases.md`,
  `docs/golden_tests.md`, `docs/audit_cases.md`, and `.ai/AGENTS.md` with
  the fixed-range split design.

---

## [v1.0] — 2026-05-18

### First release
- Initial `guess-it-2` codebase forked from `guess-it-1`.
- `mathskills` package renamed to `linearstats`; dead `run.go` and
  `predict_test.go` pruned; module path set to `guess-it-2`.
- Hybrid `Predict()` blending regression extrapolation with the robust
  centre via `|r|`, plus a 5-tier StdDev multiplier scaled by
  `(1 − 0.25·|r|)`.
- Streaming `main.go` over `bufio.Scanner`, table-driven `main_test.go`
  with seven I/O scenarios.
- Linux/amd64 student binary at `student/guess-it-2` with
  `student/script.sh` launcher.
- Baseline documentation set: `docs/PRD.md`, `docs/audit_cases.md`,
  `docs/edge_cases.md`, `docs/golden_tests.md`, `.ai/AGENTS.md`,
  `.ai/hmim.ai.log`.

[v2.2]: https://platform.zone01.gr/git/hmim/guess-it-2/compare/v2.1...v2.2
[v2.1]: https://platform.zone01.gr/git/hmim/guess-it-2/compare/v2.0...v2.1
[v2.0]: https://platform.zone01.gr/git/hmim/guess-it-2/compare/v1.0...v2.0
[v1.0]: https://platform.zone01.gr/git/hmim/guess-it-2/src/tag/v1.0
