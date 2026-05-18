# Task 02 — Docs Refresh (PRD, Audit, Edge, Golden)

## Goal
Bring the four source-of-truth documents in line with `.resources/guess-it-2*.md` and the hybrid regression + 5-tier predictor.

## Steps
1. Rewrite `docs/PRD.md` for guess-it-2: scope, new opponents (`big-range`, `linear-regr`, `correlation-coef`; bonus `mse`, `nic`), Data 4/5 datasets, hybrid algorithm spec (regression extrapolation × |r| blended with median, multiplier scaled by `1 − 0.25·|r|`), package layout (`linearstats`), and the safeguards (`minStdRange=1.8`, `minWidth=1`, initial-phase, jump-guard).
2. Rewrite `docs/audit_cases.md` per `.resources/guess-it-2_audit.md` (Data 4 + Data 5; 4 opponents + 2 bonus).
3. Update `docs/edge_cases.md` — keep current cases, add `E17 low-|r| fallback` (predictor reverts to median behaviour) and `E18 high-|r| tightening` (multiplier shrinks).
4. Update `docs/golden_tests.md` — refresh the audit-case table and add expected behaviours for the regression/correlation paths.

## Acceptance Criteria
- All four files reference `linear-regr`, `correlation-coef`, Data 4, Data 5 (not `average`, `median`, `Data 1-3`).
- `PRD.md` includes the hybrid predictor formula and the multiplier scaling rule.
- `edge_cases.md` and `golden_tests.md` add the new low/high-|r| cases.
- `grep -rIn "Data 1\|Data 2\|Data 3" docs/` returns no matches *outside* historical comparison context.

## Validation
```sh
grep -rIn "linear-regr\|correlation-coef\|Data 4\|Data 5" docs/
grep -rIn "hybrid\||r|\|PearsonCorrelation" docs/PRD.md docs/golden_tests.md
```

## Conventional commit
`docs: refresh PRD, audit, edge, and golden test docs for guess-it-2`
