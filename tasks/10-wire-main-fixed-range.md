# Task 10 — Wire main + script + IO tests

## Goal
Adapt the streaming entry point to the fixed-range split predictor and drop the
now-obsolete centre-mode selector.

## Steps
1. `main.go`: drop the `centerMode` parameter from `run`, drop the
   `useMedian`/`PREDICT_CENTER` logic; `main` calls `run(os.Stdin, os.Stdout)`;
   `Predict` is called as `linearstats.Predict(window, seen, v)`.
2. `student/script.sh`: remove `export PREDICT_CENTER="median"` — the predictor
   chooses its centre from `seen`, no env var needed.
3. `main_test.go`: drop the `centerMode` argument; relax the assertions — a
   fixed ±20 deliberately misses on volatile data, so assert output shape
   (line count, two integers, `lower <= upper`) instead of next-value coverage.

## Acceptance Criteria
- `go build ./...` succeeds; `go vet ./...` clean.
- `go test ./...` green.
- `script.sh` contains no `PREDICT_CENTER` reference.
- Smoke run emits one `lower upper` line (width 40) per numeric input.

## Validation
```sh
go vet ./...
go test ./... -cover
printf '189\n113\n121\n114\n145\n110\n' | go run .
```

## Conventional commit
`refactor(main): drop center-mode selector for fixed-range predictor`
