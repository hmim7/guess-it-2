# guess-it-2

![Go Version](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![DOCKER](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![MIT](https://img.shields.io/badge/MIT-1BB581?style=for-the-badge&logo=opensourceinitiative&logoColor=white)
![Coverage](https://img.shields.io/badge/Coverage-95.1%25-2ECC71?&labelColor=181717&style=for-the-badge&logo=codecov&logoColor=white)
[![Zone01](https://img.shields.io/badge/zone01-Athens-916ADE?&labelColor=181717&style=for-the-badge&logo=data:image/svg%2Bxml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyNCIgaGVpZ2h0PSIyNCIgdmlld0JveD0iMCAwIDI0IDI0IiBmaWxsPSJub25lIiBzdHJva2U9IndoaXRlIiBzdHJva2Utd2lkdGg9IjIiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIgc3Ryb2tlLWxpbmVqb2luPSJyb3VuZCI+PHBhdGggZD0iTTEyIDJMMiA3bDEwIDUgMTAtNS0xMC01eiIvPjxwYXRoIGQ9Ik0yIDE3bDEwIDUgMTAtNU0yIDEybDEwIDUgMTAtNSIvPjwvc3ZnPg==)](https://github.com/01-edu/public/tree/master/subjects/guess-it-2)

A real-time statistical prediction engine written in Go. The program reads a continuous stream of numbers from standard input and predicts the range `[lower, upper]` where the *next* number will fall. `guess-it-2` extends `guess-it-1` by folding `LinearRegression` and `PearsonCorrelation` from the `linear-stats` project into the predictor, producing a hybrid model that tracks clean trends and falls back to the median + 5-tier model on noisy or oscillating data.

## Features

- **Streaming Ingestion** — line-by-line `stdin` processing via `bufio.Scanner`.
- **Sliding Window** — fixed-size (`N = 20`) localized state, adapts to sudden trend changes.
- **Hybrid Predictor** — blends linear-regression extrapolation with the median + 5-tier StdDev model, weighted by `|PearsonCorrelation|`.
- **Dynamic Risk Profile** — five-tier multiplier (**Very Aggressive → Extreme**), additionally scaled by `1 − 0.25·|r|` so a clean trend tightens the band.
- **Adaptive Hit %-Targeting Controller** — closed-loop scaler that tracks realised Hit % over a rolling 30-prediction window and nudges the margin toward `PREDICT_TARGET_HIT` (default `0.30`). Set `PREDICT_ADAPT=off` to disable.
- **Safeguards** — `minStdRange = 1.8`, `minWidth = 1`, wide initial range for `seen < 5`, jump-guard for outliers `> 1.5σ` from centre.

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
120 200  --> range for the next input (in this case for the number 113)
113      --> standard input
160 230  --> range for the next input
...
```

---

## Methodology & Math

This project ports and extends the statistical foundations built in the `linear-stats` project, applying them to a real-time streaming context.

### Statistical Implementation

- **Welford's Algorithm** — variance (and StdDev) computed in a numerically stable single pass to avoid catastrophic cancellation on long streams.
- **Robust Centring (Median)** — `script.sh` exports `PREDICT_CENTER=median`; the median is significantly more resistant to extreme outliers than the mean.
- **Linear Regression Line** — for the current window of `n` values, slope `m` and intercept `b` are fit by the closed-form least-squares solution; the extrapolated next value is `m·n + b`.
- **Pearson Correlation Coefficient (r)** — measures how cleanly the window follows a linear trend. Returns `0` for constant `y` or other degenerate denominators.

### Hybrid Prediction Strategy

For every iteration with a stable window:

```
m, b   := LinearRegression(window)
r      := PearsonCorrelation(window)
regC   := m·len(window) + b
robC   := Median(window)            # or Average if PREDICT_CENTER != "median"
w      := |r|
center := w·regC + (1−w)·robC

std    := StdDev(window)
mul    := dynamicMultiplier(std) · (1 − 0.25·w)
margin := max(mul·std, minStdRange)
if |current − center| > 1.5·std:
    margin = max(margin, |current − center|·1.1)   # jump guard
margin *= multAdjust                                # adaptive Hit %-controller scaler

print roundBounds(center − margin, center + margin)
```

When `|r| ≈ 1` (clean trend), the centre tracks the regression line and the multiplier shrinks ~25 % → tighter, higher-scoring range. When `|r| ≈ 0` (oscillating or noisy data), the centre collapses to the median and the multiplier is untouched → wider, safer range.

### Adaptive Hit %-Targeting Controller

After each prediction, the `AdaptivePredictor` compares the actual `current` against the previous `[lower, upper]`, records the hit/miss into a 30-slot ring buffer, and updates a global `multAdjust ∈ [0.3, 2.0]` via proportional control:

```
err        := recentHitRate − targetHit                 # targetHit = PREDICT_TARGET_HIT (default 0.30)
multAdjust *= 1 − learnRate·err                         # learnRate = 0.15
```

Higher-than-target Hit % shrinks `multAdjust` (tighter bands, higher score per hit); lower-than-target widens it. The controller composes on top of the 5-tier multipliers — tier semantics are preserved, the scaler is a meta-knob. See [TUNING.md § 4](TUNING.md) for derivation and benchmarking guidance.

### Math Reference

**Linear Regression Line.** For `n` data points with `xᵢ = i` and `yᵢ` = the window values:

$$m = (n·Σ(xᵢ·yᵢ) − Σxᵢ·Σyᵢ) / (n·Σxᵢ² − (Σxᵢ)²)$$
$$b = (Σyᵢ − m·Σxᵢ) / n$$


**Pearson Correlation Coefficient.**

$$r = (n·Σ(xᵢ·yᵢ) − Σxᵢ·Σyᵢ) / sqrt((n·Σxᵢ² − (Σxᵢ)²) · (n·Σyᵢ² − (Σyᵢ)²))$$

Returns `0` whenever the denominator is zero (constant `y`) or intermediate values overflow.

---

## Tuning Parameters

Through benchmarking against opponent algorithms (`big-range`, `linear-regr`, `correlation-coef`, plus bonus `mse`, `nic`), these constants are finalized for optimal **Expected Value** (**Score** vs. **Hit %**):

- **Sliding Window Size:** `20`
- **Initial Safety Range:** `±60.0` (for the first 5 inputs; scaled by `|current|` for large-magnitude streams)
- **Minimum Margin Floor:** `1.8` (prevents range collapse when `StdDev ≈ 0`)
- **Minimum Output Width:** `1` (ensures lower and upper bound always differ)
- **Correlation Tightening:** `multiplier *= 1 − 0.25·|r|`
- **Dynamic Multipliers (c):**
  - Very Aggressive (`std < 5.0`): `1.2`
  - Aggressive (`std < 25.0`): `1.4`
  - Balanced (`std < 40.0`): `1.8`
  - Defensive (`std < 80.0`): `2.1`
  - Extreme (`std ≥ 80.0`): `2.4`
- **Adaptive Controller:**
  - `PREDICT_TARGET_HIT` (default `0.30`, clamped to `[0.1, 0.95]`) — target realised Hit %.
  - `PREDICT_ADAPT=off` — disables the controller; predictor reverts to the static 5-tier baseline.
  - Window size: `30` predictions · Learning rate: `0.15` · Scaler bounds: `[0.3, 2.0]`.

---

## Code Quality & Best Practices

```bash
go fmt ./...
go vet ./...
go test ./... -cover -coverpkg=./...
go test -race ./...
go test -bench=. -benchmem ./...
go test -fuzz=Fuzz -fuzztime=10s ./...
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
│   ├── stats.go                  # Average, Median, Variance, StdDev, LinearRegression, PearsonCorrelation
│   ├── predict.go                # Hybrid Predict(), AdaptivePredictor, dynamicMultiplier(), tuning constants
│   └── linearstats_test.go       # Table-driven stats + predictor coverage (≥ 95 %)
├── student/                      # Auditor artefacts
│   ├── guess-it-2                # Linux/amd64 binary (cross-compiled)
│   └── script.sh                 # Launcher: exports PREDICT_CENTER=median, PREDICT_TARGET_HIT=0.30
├── docs/
│   ├── PRD.md                    # Product requirements & architecture
│   ├── audit_cases.md            # Audit procedure & success criteria
│   ├── edge_cases.md             # Streaming / statistical edge cases
│   └── golden_tests.md           # Single source of truth for expected outputs
├── tasks/                        # Implementation task cards (01-08)
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

```
http://localhost:3000/?guesser=big-range
http://localhost:3000/?guesser=correlation-coef
```

Select a **Test Data** set, then click **Quick** to skip the animation and jump straight to the final scores. Click **Clean** to reset the display before the next run.

Recommended opponents for the audit: `big-range`, `linear-regr`, `correlation-coef`, plus bonus `mse` and `nic`.

### Quick-reference checklist

- [ ] Binary rebuilt for Linux: `GOOS=linux GOARCH=amd64 go build -o student/guess-it-2 .`
- [ ] Executable bits set: `chmod +x student/guess-it-2 student/script.sh`
- [ ] `student/` copied into `guess-it-dockerized/student/`
- [ ] Podman socket running: `systemctl --user start podman.socket`
- [ ] Container started: `cd guess-it-dockerized && docker compose up --build`
- [ ] Test on `Data 4` and `Data 5` with `?guesser=correlation-coef` (primary) and others
- [ ] Use **Quick** to fast-forward; **Clean** between runs

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
- [Edge Cases](docs/edge_cases.md) — streaming, statistical, and hybrid-predictor edges
- [Golden Tests](docs/golden_tests.md) — single source of truth for expected behaviour
- [PRD](docs/PRD.md) — product requirements & architecture (with Mermaid flowchart)
- [Tuning Guide](TUNING.md) — directional effect of each `predict.go` constant on Score vs. Hit %
- [Task Cards](tasks/) — implementation breakdown
- [AI Usage Log](.ai/hmim.ai.log) — record of AI-assisted development sessions

---

## License

*This project is licensed under the MIT License.  
Part of the Zone01 School curriculum.*
