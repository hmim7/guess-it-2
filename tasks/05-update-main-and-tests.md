# Task 05 — Wire main.go and Expand main_test.go

## Goal
Finalize the stdin-stream entry point against the new `linearstats` API and broaden the IO test coverage.

## Steps
1. `main.go`: keep the `run(io.Reader, io.Writer, string)` shape; update the import to `guess-it-2/linearstats`; rebrand comments away from `guess-it-1`.
2. `main_test.go`: extend the table with:
   - trended stream (e.g. `100,110,120,…`) — assert the produced ranges include each next value.
   - oscillating stream (`10,100,10,100,…`) — assert ranges are wide enough.
   - EOF on empty input — no panic, exit clean.
   - mixed numeric + invalid lines — only numeric inputs produce output.
   - PREDICT_CENTER unset vs `"median"` — both code paths exercised.
3. Run `go vet`, `go fmt`, `go test`.

## Acceptance Criteria
- `main_test.go` is table-driven (one parent `TestRun`).
- Tests cover the five new scenarios above.
- `go test ./... -cover` green; coverage on the `main` package ≥ 80 %.

## Validation
```sh
go vet ./...
go test ./... -v -cover
```

## Conventional commit
`refactor(main): wire linearstats predictor and expand IO tests`
