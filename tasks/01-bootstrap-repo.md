# Task 01 — Bootstrap Repo

## Goal
Initialize the version-controlled workspace for guess-it-2 and rebrand the Go module.

## Steps
1. `git init -b main` at project root.
2. Configure local git: `user.name=hmim`, `user.email=helenamim7@gmail.com`.
3. Create `.gitignore` (Go defaults + IDE + OS noise).
4. Update `go.mod` module path to `guess-it-2`.
5. Create the `tasks/` directory and scaffold cards 01–08.
6. Create `.ai/hmim.ai.log` and seed the first (planning) entry.

## Acceptance Criteria
- `git rev-parse --abbrev-ref HEAD` returns `main`.
- `git config user.name` returns `hmim`; `git config user.email` returns `helenamim7@gmail.com`.
- `go.mod` first line is `module guess-it-2`.
- `tasks/` contains cards 01–08; `.ai/hmim.ai.log` exists with the Phase-0 entry.

## Validation
```sh
git status
git config --get user.name
git config --get user.email
head -n 1 go.mod
ls tasks .ai
```

## Conventional commit
`chore: bootstrap guess-it-2 repo and tooling`
