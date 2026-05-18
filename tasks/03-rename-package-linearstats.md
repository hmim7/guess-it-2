# Task 03 — Rename Package to linearstats; Prune Dead Modules

## Goal
Rename the `mathskills` package to `linearstats`, drop dead code carried over from `linear-stats`, and keep the package compiling and tested.

## Steps
1. Rename directory `mathskills/` → `linearstats/`.
2. Update `package mathskills` → `package linearstats` in every file.
3. Update the import in `main.go` to `guess-it-2/linearstats` (alias if necessary).
4. Rename `mathskills_test.go` → `linearstats_test.go`.
5. **Delete** `linearstats/run.go` (file-CLI from linear-stats — not used by the stdin-stream flow).
6. **Merge** the cases from `mathskills/predict_test.go` into `linearstats_test.go` (table-driven), then delete `predict_test.go`.
7. **Remove** parser and run tests from `linearstats_test.go` (`TestParseTrimWhitespace`, `TestParseScientificNotation`, `TestParseInvalid`, `TestParseEmptyLine`, `TestParseOverflow`, `TestParseLineWarn*`, `TestSanitizeForLog`, `TestRun`, `TestValidatePath`, `TestReadData`) — the parser module is intentionally not ported.
8. Ensure the package compiles: `go build ./...`.

## Acceptance Criteria
- No file references `mathskills` (`grep -rIn mathskills .` returns nothing outside `.git`).
- `linearstats/run.go` and `linearstats/predict_test.go` do not exist.
- `go vet ./...` clean.
- `go test ./...` green (Task 04 will add the regression cases).

## Validation
```sh
go vet ./...
go test ./...
grep -rIn "mathskills" . --exclude-dir=.git
ls linearstats
```

## Conventional commit
`refactor(linearstats): rename mathskills package and prune dead modules`
