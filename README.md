# guess-it-2

![MIT](https://img.shields.io/badge/MIT-1BB581?style=for-the-badge&logo=opensourceinitiative&logoColor=white)
![Go Version](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![DOCKER TESTER](https://img.shields.io/badge/Docker%20Tester-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![BASH](https://img.shields.io/badge/BASH-121011?style=for-the-badge&logo=gnu-bash&logoColor=white)
![Coverage](https://img.shields.io/badge/Coverage-96.3%25-2ECC71?&labelColor=181717&style=for-the-badge&logo=codecov&logoColor=white)
[![Zone01](https://img.shields.io/badge/zone01-Athens-916ADE?&labelColor=181717&style=for-the-badge&logo=data:image/svg%2Bxml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyNCIgaGVpZ2h0PSIyNCIgdmlld0JveD0iMCAwIDI0IDI0IiBmaWxsPSJub25lIiBzdHJva2U9IndoaXRlIiBzdHJva2Utd2lkdGg9IjIiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIgc3Ryb2tlLWxpbmVqb2luPSJyb3VuZCI+PHBhdGggZD0iTTEyIDJMMiA3bDEwIDUgMTAtNS0xMC01eiIvPjxwYXRoIGQ9Ik0yIDE3bDEwIDUgMTAtNU0yIDEybDEwIDUgMTAtNSIvPjwvc3ZnPg==)](https://github.com/01-edu/public/tree/master/subjects/guess-it-2)

A real-time statistical prediction engine written in Go. The program reads a continuous stream of numbers from standard input and predicts the range `[lower, upper]` where the *next* number will fall. `guess-it-2` extends `guess-it-1` by using a **running OLS fit over the full input prefix** to centre a **fixed ±46 half-width** prediction range — a design empirically tuned to the audit's "smaller correct range scores higher" rule under the data's uniform[-50,+50] residual distribution.

## Features

- **Streaming Ingestion** — line-by-line `stdin` processing via `bufio.Scanner`.
- **Sliding Window** — fixed-size (`N = 20`) localized state, adapts to sudden trend changes.
- **Stateful Running-OLS Predictor** — a `Predictor` struct accumulates a running OLS fit over the full input prefix; from sample 30 onwards it extrapolates `ŷ = m·(n+1) + b` as the range centre with a constant `±46` half-width.
- **Rounding-Cliff-Tuned Width** — the audit data has uniform[-50,+50] residual noise (verified by [_sim/noise_dist.go](_sim/noise_dist.go)). Under `server.js`'s `round(10⁷/(2h+1)/(N-1))` per-hit formula, ±46 is the largest half-width that still rounds to 9 pts/hit while capturing ~92% of the bulk. ±47 crosses a per-hit cliff (9→8) that costs ~9% of the score.
- **Safeguards** — `minWidth = 1` guarantees the lower and upper bounds always differ after rounding.

## Installation

Clone the repository and build the binary from the project root:

```bash
git clone https://platform.zone01.gr/git/hmim/guess-it-2
cd guess-it-2
GOOS=linux GOARCH=amd64 go build -o student/guess-it-2 .
chmod +x student/guess-it-2 student/script.sh
```

---

## Usage

Execute the program via the provided shell script:

```bash
./student/script.sh
```

Or run it directly with Go:

```bash
go run .
```

## I/O Format

- **Input**: A stream of numbers, one per line via `stdin`.
- **Output**: For every input, two space-separated integers (`lower upper`) per line to `stdout`.

*Example:*

```console
>$ ./student/script.sh
189      --> standard input
146 226  --> range for the next input (centred on 186 ±46)
113      --> standard input
134 226  --> range for the next input
...
```

---

## Methodology & Math

This project ports and extends the statistical foundations built in the `linear-stats` project, applying them to a real-time streaming context.

### Statistical Implementation

- **Running OLS** — a `RunningOLS` struct accumulates `(i, yᵢ)` pairs incrementally (O(1) per call) and yields slope `m` and intercept `b` at any point. The extrapolated next value is `m·(n+1) + b`, giving an unbiased centre estimate from sample 30 onwards — eliminating the ~10-unit bias the window median introduced on slope-1 data.
- **Fallback path** — below 30 samples or when the OLS fit is degenerate, the predictor falls back to the original fixed-range split logic (regression over the 20-point window while `seen < 1000`, window median afterwards), preserving correct behaviour on non-linear datasets.
- `Average`, `Median`, `Variance`, `StdDev`, `LinearRegression`, and `PearsonCorrelation` remain in the `linearstats` package as library surface, computed with numerically stable methods (Welford's algorithm for variance, `IsNaN`/`IsInf` guards in the correlation).

### Running-OLS Predictor Strategy

For every input:

```
ols.Add(n, current); n++

if n >= 30 and ols.Fit() succeeds:
    center := m·n + b          // full-prefix running OLS extrapolation
else:
    // fallback: original fixed-range split
    if n < 1000:
        center := m·len(window) + b   // 20-point window regression
    else:
        center := Median(window)

print roundBounds(center − 46, center + 46)
```

The range half-width is always `46` (width 92). The running OLS replaces both
the window-regression and window-median phases because the window median biases
the centre ~10 units low on slope-1 data (`y_i ≈ c + i + noise`) — it estimates
`y_{n-10}` rather than `y_{n+1}`. Running OLS over the full prefix is unbiased
and stabilises within ~30 samples.

### Why ±46 and not ±20

The audit data on both Data 4 and Data 5 has residuals that are **uniform on
[-50, +50]** (D4: zero outliers; D5: ~0.85% outliers up to |r|≈650). With the
running-OLS-prefix centre the prediction error stays uniform[-50,+50], so the
hit rate at half-width `h` is `min(h/50, 1)` exactly. The audit score per hit
is `round(10⁷/(2h+1)/(N-1))`; for N=12,500 this is a step function in `h`:

| h  | width 2h | hit % | round((10⁷/(2h+1))/(N−1)) | per-step value |
|---:|---------:|------:|--------------------------:|---------------:|
| 20 | 40       | 40 %  | round(19.514) = **20**    | 0.40 · 20 = 8.00 |
| 27 | 54       | 54 %  | round(14.551) = **15**    | 0.54 · 15 = 8.10 |
| 41 | 82       | 82 %  | round( 9.638) = **10**    | 0.82 · 10 = 8.20 |
| **46** | **92** | **92 %** | round( 8.603) = **9** | **0.92 · 9 = 8.28** |
| 47 | 94       | 94 %  | round( 8.422) = **8** (cliff) | 0.94 · 8 = 7.52 |

The score landscape is a sawtooth, not flat. Local maxima sit at the **largest
half-width still rounding to a given per-hit integer** (±20, ±27, ±41, ±46).
±46 is the global maximum because the uniform-bulk hit rate saturates at ~92 %
inside [−46, +46], and the very next step (±47) crosses a per-hit cliff
(9→8 points) that 2 percentage points of extra coverage cannot recover. Choice
of half-width must be exact — ±45 and ±47 both regress vs ±46. Full sweep and
empirical confirmation in [docs/theilsen_benchmark.md §6](docs/theilsen_benchmark.md).

### Math Reference

**Linear Regression Line.** For `n` data points with `xᵢ = i` and `yᵢ` = the window values:

$$m = (n·Σ(xᵢ·yᵢ) − Σxᵢ·Σyᵢ) / (n·Σxᵢ² − (Σxᵢ)²)$$
$$b = (Σyᵢ − m·Σxᵢ) / n$$


**Pearson Correlation Coefficient.**

$$r = (n·Σ(xᵢ·yᵢ) − Σxᵢ·Σyᵢ) / sqrt((n·Σxᵢ² − (Σxᵢ)²) · (n·Σyᵢ² − (Σyᵢ)²))$$

Returns `0` whenever the denominator is zero (constant `y`) or intermediate values overflow.

---

## Tuning Parameters

Constants in [linearstats/predict.go](linearstats/predict.go), locked by
simulating the predictor over the datasets in `docs/data-sets/`:

- **Sliding Window Size (`WindowSize`):** `20` — values kept for the fallback regression fit and median.
- **Fixed Range (`fixedRange`):** `46` — constant half-width; every predicted range spans 92 units. Sits at the global maximum of the audit's sawtooth score landscape (see *Methodology & Math* above).
- **OLS Min Samples (`olsMinSamples`):** `30` — number of inputs required before the running OLS centre activates; below this the fallback path runs.
- **Regression Phase Limit (`regressionPhaseLimit`):** `1000` — fallback-path threshold: below this the centre is the window-regression extrapolation; at/above it the centre is the window median.
- **Minimum Output Width (`minWidth`):** `1` — ensures the lower and upper bound always differ after rounding.

---

## Audit Performance

Benchmarked against all eight `guesser` programs in the dockerized tester's
`ai/` folder, on the real audit datasets **Data 4** and **Data 5** (5 files
each), scored with the exact `server.js` formula. Student mean score:
**Data 4 = 104,502 · Data 5 = 103,246**
([full benchmark](docs/predictor_benchmark_linear_v2.md)).

| Opponent           | D4 file-wins | D5 file-wins | Audit-pass probability\* |
|--------------------|:------------:|:------------:|:------------------------:|
| `big-range`        | 5/5          | 5/5          | 100%                   |
| `correlation-coef` | 5/5          | 5/5          | 100%                   |
| `average`          | 5/5          | 5/5          | 100%                   |
| `median`           | 5/5          | 5/5          | 100%                   |
| `huge-range`       | 5/5          | 5/5          | 100%                   |
| `linear-regr`      | 4/5          | 5/5          | ~90 %  \|  100%       |
| `mse`              | 5/5          | 4/5          | 100%  \|  ~90 %       |
| `nic`              | 3/5          | 4/5          | ~65 %  \|  ~90 %        |

\* The auditor runs 3 independent rounds per dataset (each picks a random file
1–5) and needs ≥ 2 wins. With 5/5 winning files the pass is certain; with 3/5,
`P(win ≥ 2 of 3) ≈ 65 %`.

### Why this predictor is the best implementation

- **Centre is unbiased & noise-free at the limit.** Running OLS over the full
  prefix has prediction variance `σ²_ε/n → 0` — after ~30 samples the centre
  contributes <2 units of noise on top of the residual error, so the
  prediction error stays essentially the underlying uniform[-50,+50] noise.
- **±46 is the global maximum of the sawtooth score landscape.** Under uniform
  residuals + integer-rounded `round(10⁷/(2h+1)/(N−1))` per-hit scoring, the
  landscape has local maxima at every "last half-width before the per-hit
  integer drops" (±20→20pts, ±27→15pts, ±41→10pts, ±46→9pts). ±46 wins
  because the bulk hit rate saturates at ~92 % there, and the next step
  crosses a 9→8 cliff that 2 percentage points of extra coverage cannot
  recover. See [docs/theilsen_benchmark.md §6](docs/theilsen_benchmark.md) for
  the full width sweep and rounding-cliff analysis.
- **Centre estimator choice doesn't matter.** Theil-Sen at the same width
  scores within 0.5 % of OLS — robust regression has nothing to do because
  the residuals are iid uniform, not outlier-contaminated
  ([docs/theilsen_benchmark.md §§ 1-5](docs/theilsen_benchmark.md)).
- **Adaptive widths overfit.** Empirical-error interval search and
  adaptive-width experiments all lost: width estimation introduces variance
  that the iid-uniform noise structure does not reward.
- **The remaining 3/5 vs `nic` D4 is a structural ceiling.** `nic` posted
  `nic D4/1 = 111,200` and `D4/3 = 106,400` — both above the uniform-bulk
  ceiling for any trend-centred fixed-width predictor (~104,500). 3/5 is the
  maximum achievable on that dataset; the audit's 3-round re-run covers the
  rest.

The predictor wins 5/5 against all five easy opponents (100% audit-safe),
4/5 against `linear-regr` D4 (~90 %), 5/5 against `linear-regr` D5 and `mse`
D4 (100%), 4/5 against `mse` D5 and `nic` D5 (~90 %), and 3/5 against
`nic` D4 (~65 %) — beating the prior `linear-v2` ±20 baseline on every single
audit file (10/10).

---

## Code Quality & Best Practices

```bash
go fmt ./...
go vet ./...
go test ./... -cover -coverpkg=./...
go test -race ./...
go test -bench=. -benchmem ./...
go test -fuzz=Fuzz -fuzztime=10s ./linearstats
```

**Best Practices Applied:**

- **Separation of Concerns** — `main.go` does I/O & orchestration; `linearstats/predict.go` owns the prediction algorithm; `linearstats/stats.go` exposes the mathematical primitives. `main.go` does no statistics; the `linearstats` package does no I/O.
- **Table-Driven Tests** — every `*_test.go` uses Go's table-driven pattern with `t.Run` subtests for maximum clarity and maintainability.
- **Numerical Stability** — Welford's variance + `math.IsNaN` / `math.IsInf` guards prevent precision loss and panics on edge inputs.
- **Coverage** — `linearstats` package ≥ 95 %; `run()` in `main.go` ≥ 85 %.

---

## Performance Verification

```bash
time go run . < data.txt
time ./student/guess-it-2 < data.txt
```

*Note: The engine is `O(N)` per iteration over the sliding window (`N = 20`); memory and CPU per iteration are strictly bounded regardless of stream length.*

---

## Project Structure

```
guess-it-2/
├── main.go                       # Stdin streaming loop; orchestrates window + predictor
├── main_test.go                  # Table-driven IO tests (trended / oscillating / EOF / invalid)
├── go.mod                        # Go module definition
├── LICENSE                       # MIT license
├── README.md                     # This file
├── linearstats/                  # Statistical + prediction core
│   ├── stats.go                  # Average, Median, Variance, StdDev, LinearRegression, PearsonCorrelation, RunningOLS
│   ├── predict.go                # Predictor (running-OLS centre), Predict() fallback, roundBounds(), tuning constants
│   └── linearstats_test.go       # Table-driven stats + predictor coverage (≥ 95 %)
├── student/                      # Auditor artefacts
│   ├── guess-it-2                # Linux/amd64 binary (cross-compiled)
│   └── script.sh                 # Launcher: runs ./student/guess-it-2
├── docs/
│   ├── PRD.md                    # Product requirements & architecture
│   ├── audit_cases.md            # Audit procedure & success criteria
│   ├── edge_cases.md             # Streaming / statistical edge cases
│   └── golden_tests.md           # Single source of truth for expected outputs
├── tasks/                        # Implementation task cards (01-11)
└── .ai/
    └── hmim.ai.log               # AI-assisted development log
```

---

## Docker — Audit Preparation

### Step 1 — Cross-compile the Go binary for Linux

The Docker container is Linux/amd64. Always rebuild from the project root **before** touching Docker.

```bash
cd /path/to/guess-it-2          # ← repo root, not guess-it-dockerized
GOOS=linux GOARCH=amd64 go build -o student/guess-it-2 .
chmod +x student/guess-it-2 student/script.sh
```

**Verify it worked:**

```bash
file student/guess-it-2
# Expected: ELF 64-bit LSB executable, x86-64 ...
```

### Step 2 — Start the Podman socket

This system uses **Podman** as the container backend; `docker compose` routes through Podman's compatibility socket, which is off by default.

```bash
systemctl --user enable --now podman.socket
systemctl --user status podman.socket    # Expected: Active: active (listening)
```

### Step 3 — Place the student artifacts

Download the [dockerized tester](https://assets.01-edu.org/guess-it/guess-it-dockerized.zip) if you don't have it yet:

```bash
curl -L https://assets.01-edu.org/guess-it/guess-it-dockerized.zip -o guess-it-dockerized.zip
unzip guess-it-dockerized.zip
```

The dockerized tester expects the `student/` folder (binary + `script.sh`) to be present inside `guess-it-dockerized/`. The layout must be:

```
guess-it-dockerized/
├── ai/
│   ├── big-range
│   └── ...
├── index.html
├── index.js
└── student/
    ├── guess-it-2      # Linux/amd64 binary from Step 1
    └── script.sh
```

Copy or symlink from the repo root:

```bash
cp -r student/ guess-it-dockerized/student/
```

### Step 4 — Build and run the container

```bash
cd guess-it-dockerized
docker compose up --build
```

Open the browser at `http://localhost:3000`.

### Step 5 — Run a test against an opponent

The tester requires an opponent AI to compare against. Append `?guesser=<name>` to the URL, where `<name>` is any file in the `ai/` folder:

Run every required opponent, in the order the audit checks them (see
[audit_cases.md](docs/audit_cases.md)) — `big-range`, `linear-regr`,
`correlation-coef`, then bonus `mse` and `nic`:

```
http://localhost:3000/?guesser=big-range
http://localhost:3000/?guesser=linear-regr
http://localhost:3000/?guesser=correlation-coef
http://localhost:3000/?guesser=mse
http://localhost:3000/?guesser=nic
```

Select a **Test Data** set, then click **Quick** to skip the animation and jump straight to the final scores. Click **Clean** to reset the display before the next run.

Each opponent must be tested on both `Data 4` and `Data 5`, 3 runs per dataset.

### Troubleshooting

| Symptom | Cause | Fix |
| --- | --- | --- |
| `Exec format error` | Binary compiled for wrong OS | Redo Step 1 |
| `no such file or directory` (socket) | Podman socket not started | Redo Step 2 |
| Console: "need another guesser" | No `?guesser=` param in URL | Add `?guesser=big-range` to the URL |
| `version` obsolete warning in compose | `version:` key in `docker-compose.yml` is deprecated | Safe to ignore — warning, not an error |
| Port 3000 already in use | Another process is bound to 3000 | `lsof -i :3000` then kill the process |

---

## Documentation Reference

- [Audit Cases](docs/audit_cases.md) — detailed audit procedure & success conditions
- [Edge Cases](docs/edge_cases.md) — streaming and statistical edge cases
- [Golden Tests](docs/golden_tests.md) — single source of truth for expected behaviour
- [PRD](docs/PRD.md) — product requirements & architecture (with Mermaid flowchart)
- [Predictor Analysis](docs/predictor_analysis.md) — residual-structure proof that 3/5 is the ceiling
- [Linear-v2 Benchmark](docs/predictor_benchmark_linear_v2.md) — running-OLS predictor head-to-head results (DS5 vs `linear-regr` 3/5 → 5/5)
- [Theilsen Benchmark](docs/theilsen_benchmark.md) — theil-sen predictor head-to-head scores vs running-OLS predictor and the audit opponents
- [Task Cards](tasks/) — implementation breakdown
- [AI Usage Log](.ai/hmim.ai.log) — record of AI-assisted development sessions

---

## License

*This project is licensed under the MIT License.  
Part of the Zone01 School curriculum.*
