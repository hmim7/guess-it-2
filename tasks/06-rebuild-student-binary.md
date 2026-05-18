# Task 06 — Rebuild Student Binary and Refresh script.sh

## Goal
Produce the Linux/amd64 binary the auditor expects and point `script.sh` at it.

## Steps
1. Remove the stale `student/guess-it-1` binary.
2. Build: `GOOS=linux GOARCH=amd64 go build -o student/guess-it-2 .` from the project root.
3. Update `student/script.sh` to `./student/guess-it-2` (keep `export PREDICT_CENTER="median"`).
4. Ensure executable bits: `chmod +x student/guess-it-2 student/script.sh`.

## Acceptance Criteria
- `file student/guess-it-2` reports `ELF 64-bit LSB executable, x86-64`.
- `student/script.sh` invokes `./student/guess-it-2`.
- No `student/guess-it-1` file remains.

## Validation
```sh
ls -l student
file student/guess-it-2
grep guess-it student/script.sh
```

## Conventional commit
`build(student): rename binary to guess-it-2 and refresh script.sh`
