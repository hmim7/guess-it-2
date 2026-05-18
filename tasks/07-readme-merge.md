# Task 07 — README Merge and AGENTS Context Update

## Goal
Produce a single `README.md` that combines the guess-it-2 content with the structural sections from `README (2).md`, and refresh `docs/AGENTS.md` to match the new project context.

## Steps
1. In `README.md`:
   - Rebrand title and badges to `guess-it-2`.
   - Carry over from `README (2).md`: **Math** section (Linear Regression + Pearson Correlation formulae), **Output Format**, **Project Structure**, **Documentation Reference**, **Code Quality** bullets.
   - Replace `mathskills/` with `linearstats/` in the structure tree.
   - Add a "Methodology" subsection covering the hybrid predictor (regression × |r| blend, multiplier scaling).
   - Drop the "Installation" stanza that references the old `linear-stats` clone URL; keep the guess-it-2 install/usage block.
2. Delete `README (2).md`.
3. Update `docs/AGENTS.md`: rebrand to "Guess It 2", reference `linearstats` package, and the new opponents/datasets.

## Acceptance Criteria
- `README.md` references `guess-it-2`, `linearstats`, hybrid predictor, Data 4/5, and the new opponents.
- `README (2).md` is removed.
- `docs/AGENTS.md` no longer says "Guess It 1" or `math-skills`.

## Validation
```sh
test ! -f "README (2).md" && echo "deleted"
grep -n "guess-it-2\|linearstats\|hybrid" README.md
grep -n "Guess It 1\|math-skills" docs/AGENTS.md && echo "STALE" || echo "ok"
```

## Conventional commit
`docs(readme): merge linear-stats structure into updated guess-it-2 README`
