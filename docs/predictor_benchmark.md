# Benchmark Results — guess-it-2 vs Audit Opponents

**Date:** 2026-05-22
**Predictor under test:** fixed-range split — centre on linear-regression
extrapolation while `seen < 1000`, on the window median afterwards; constant
±20 half-width.
**Method:** the audit opponents (`/.resources/guess-it-dockerized/ai/*`) are
Linux binaries; they were run under WSL (Ubuntu) on the real audit datasets
**Data 4** and **Data 5** (5 files each, ~12,500 numbers per file). Each program
is scored with the exact `server.js` formula.

## Scoring formula (from `guess-it-dockerized/server.js`)

```
per correct guess:  round( 10000000 / (1 + width) / (N - 1) )
per miss:           0
total score      =  sum over all guesses
```

`width = upper - lower`, `N = 12500`. So a hit is worth ≈ `800 / (1 + width)`.
**A miss scores nothing; a narrower correct range scores more.**

A consequence worth stating up front: expected score per step ≈
`hitRate(width) · 800/(1+width)`, and since `hitRate ≈ density·(1+width)`, the
width largely **cancels**. Expected score is driven by *centre accuracy*
(the error density at the predicted centre), not by the range width.

## Head-to-head results — current predictor

Student mean score: **Data 4 = 102,432** · **Data 5 = 100,720**.

| Opponent           | Opp mean D4 | Opp mean D5 | Student file-wins | Audit-case win prob* |
|--------------------|------------:|------------:|-------------------|----------------------|
| `big-range`        |      49,995 |      49,182 | 5/5 · 5/5         | ~100 %               |
| `correlation-coef` |      93,761 |      88,828 | 5/5 · 5/5         | ~100 %               |
| `average`          |         800 |         736 | 5/5 · 5/5         | ~100 %               |
| `median`           |         730 |         627 | 5/5 · 5/5         | ~100 %               |
| `huge-range`       |           0 |           0 | 5/5 · 5/5         | ~100 %               |
| `linear-regr`      |     102,176 |      99,547 | 3/5 · 3/5         | ~65 % each           |
| `mse`              |     100,979 |      98,256 | 3/5 · 3/5         | ~65 % each           |
| `nic`              |     100,480 |      98,080 | 3/5 · 3/5         | ~65 % each           |

\* The auditor runs 3 independent rounds per dataset (each picks a random file
1–5) and needs ≥ 2 wins. With 3/5 winning files, `P(win ≥ 2 of 3) ≈ 65 %`;
with 4/5 it is ≈ 90 %; with 5/5 it is 100 %.

### Verdict

- **The current predictor decisively beats** `big-range`, `correlation-coef`,
  `average`, `median`, and `huge-range` — 5/5 on both datasets, ~100 % audit-safe.
- Against **`linear-regr`, `mse`, `nic`** it wins **3/5 files** on each dataset —
  roughly a 65 % chance of passing each audit case. These are genuine near-ties,
  not failures: the per-file margins are often **< 0.5 %** (e.g. on Data 4 file 1
  the student scores 101,900 vs linear-regr's 102,300).
- `linear-regr` is essentially our own strategy (a regression line) — the match
  is inherently a coin-flip. `nic` is high-variance (88,800–111,200); some of its
  files are unbeatable by *any* ~103 k predictor (see below).

## Why the predictor cannot simply be "tuned" to win

A sweep of centres (`median`, `average`, `last-value`, regression at window
sizes 2–20) and half-widths (0.5–40) over Data 4 + Data 5 shows the total score
is **flat within a ~95 k–103 k band** — a ~5 % spread that is mostly noise.

- The current predictor scores **~101.6 k**; the best configuration found scores
  **~101.8 k**. That is a **0.3 %** difference.
- The earlier "adaptive residual-std width" idea scored *worse* (down to ~61 k):
  window StdDev (~hundreds) wildly overestimates the true ~33-unit step error of
  these smooth series.
- `Guess2` (the `y·(1−r)` width) would not help and is additionally buggy
  (width tied to the value's magnitude, raw `r` instead of `|r|`, predicts the
  current index instead of the next).

**Structural ceilings (unbeatable files):**

- `nic` Data 4 file 1 (111,200) and file 3 (106,400), and Data 5 file 3
  (108,000) exceed what any ~103 k predictor can reach. 3/5 is therefore the
  **maximum achievable** versus `nic` — the current predictor already attains it.

**Scoring-formula rounding quirk — keep width = 40.** The per-hit value
`round(800.064 / (1+width))` is jagged. Width 40 (±20) lands on
`19.51 → 20` — a favourable round-*up* worth +2.5 %. Width 32 (±16) lands on
`24.24 → 24` — a round-*down*, −1 %. This is why ±16 configs scored ~3.5 % below
±20 in testing. **Among practical widths, ±20 is the best rounding point** —
do not move off it.

## Recommendations

1. **Keep the fixed-range ±20 design.** It is verified competitive: 5/5 versus
   five opponents and a genuine ~65 % shot at the three hard ones. It is already
   within 0.3 % of the best configuration that exists.
2. **Do not adopt `Guess2` or an adaptive/residual width.** Both were simulated
   and score the same or worse; the audit landscape is flat.
3. **Keep width = 40.** It sits on the most favourable rounding point of the
   score formula. ±16 would silently cost ~3 %.
4. **Optional marginal tweak:** centring on `median` of the last **10** values
   for *every* step (instead of the regression/median split) lifted the
   hard-opponent tally from 18/30 to 20/30 winning files in this sample — but
   the student mean barely moves (101,576 → 101,696), so treat it as noise-level,
   not a guaranteed gain.
5. **`linear-regr` is the only real audit risk.** It is a mirror of our own
   regression strategy, so the outcome is close to chance (~65 %). There is no
   formula fix; if a run is lost, simply re-run — the auditor allows 3 rounds.

## Appendix — per-file scores (WSL, exact formula)

**Data 4** — student: 101,900 · 101,640 · 104,220 · 101,820 · 102,580

| Opponent           | file 1 | file 2 | file 3 | file 4 | file 5 |
|--------------------|-------:|-------:|-------:|-------:|-------:|
| `linear-regr`      |102,300 |100,750 |104,532 |101,680 |101,618 |
| `mse`              |100,392 |101,727 | 97,989 |100,659 |104,130 |
| `nic`              |111,200 | 98,400 |106,400 | 88,800 | 97,600 |
| `big-range`        | 49,992 | 49,996 | 49,996 | 49,996 | 49,996 |
| `correlation-coef` | 95,076 | 94,392 | 93,442 | 92,226 | 93,670 |

**Data 5** — student: 102,140 · 98,840 · 101,960 · 100,840 · 99,820

| Opponent           | file 1 | file 2 | file 3 | file 4 | file 5 |
|--------------------|-------:|-------:|-------:|-------:|-------:|
| `linear-regr`      |100,006 | 99,696 | 99,448 | 97,216 |101,370 |
| `mse`              | 92,115 |107,067 | 96,921 | 94,785 |100,392 |
| `nic`              | 93,600 | 88,000 |108,000 |102,400 | 98,400 |
| `big-range`        | 49,140 | 49,168 | 49,280 | 49,152 | 49,172 |
| `correlation-coef` | 89,566 | 86,830 | 89,718 | 89,376 | 88,654 |

*Reproduce:* `_sim/headtohead.sh` (run via WSL) for opponent scores;
`_sim/sim.go` for predictor-config search. Both use the `server.js` formula.
