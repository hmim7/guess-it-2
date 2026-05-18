# Task 08 — Final Verification

## Goal
Confirm the refactor is complete, tests pass, the binary runs, and no legacy references remain.

## Steps
1. `go vet ./...`
2. `go test ./... -cover -coverpkg=./...`
3. Smoke test: `printf '189\n113\n121\n114\n145\n110\n' | go run .`
4. Cross-check: `grep -rIn "mathskills\|guess-it-1" . --exclude-dir=.git --exclude="*.log"`.
5. Verify the student binary: `./student/script.sh < .resources/guess-it-2.md` (or a small sample) produces well-formed `lower upper` lines.

## Acceptance Criteria
- `go vet` clean.
- `go test ./...` green; `linearstats` coverage ≥ 90 %, overall ≥ 80 %.
- The naming sweep returns no matches.
- Smoke test outputs one `lower upper` line per numeric input.

## Validation
```sh
go vet ./...
go test ./... -cover -coverpkg=./...
printf '189\n113\n121\n114\n145\n110\n' | go run .
grep -rIn "mathskills\|guess-it-1" . --exclude-dir=.git --exclude="*.log"
```

## Conventional commit
`test: full verification pass for guess-it-2`
