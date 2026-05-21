# Task 11 — Docs Reconcile & Rebuild

## Goal
Bring every document in line with the fixed-range split predictor and ship a
fresh audit binary.

## Steps
1. `README.md`: replace the hybrid/adaptive-controller description with the
   fixed-range split strategy; drop the speculative "Adaptive Hit %-Targeting
   Controller" section and the dead `TUNING.md` link; update tuning constants
   and the project tree.
2. `docs/PRD.md`: rewrite §4.2 (core logic), the §6.1 golden tests, the §7
   architecture/flowchart, milestones, and risks for the fixed-range model.
3. `docs/edge_cases.md`, `docs/golden_tests.md`: reconcile to the fixed-range
   model; remove the hybrid `|r|`-blend cases.
4. `docs/audit_cases.md`: drop the tier/`|r|` rationale from the per-opponent
   notes.
5. `.ai/AGENTS.md`: update the architecture, logic rules, and phase plan.
6. Rebuild `student/guess-it-2` (`GOOS=linux GOARCH=amd64`); keep executable bits.
7. Verify: `go vet`, `go test -cover`, smoke run, dataset score-check.

## Acceptance Criteria
- No doc references the hybrid model, 5-tier multiplier, `dynamicMultiplier`,
  `PREDICT_CENTER`, or the adaptive controller.
- `go vet ./...` clean; `go test ./...` green.
- `student/guess-it-2` is a fresh `ELF 64-bit LSB executable, x86-64`.
- Dataset score-check reproduces the analysis (~40 % hit on groups 1/3/5,
  ~9 % on group 9, constant width 40).

## Validation
```sh
go vet ./...
go test ./... -cover
printf '189\n113\n121\n114\n145\n110\n' | go run .
file student/guess-it-2
```

## Conventional commit
`docs: reconcile docs with the fixed-range split predictor`
