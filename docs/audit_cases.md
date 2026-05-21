# Audit Cases: Guess It 2

The program is audited by competing against a series of benchmark "guesser" programs on **Data 4** and **Data 5**. The student's program must achieve a higher score than the opponent in at least **2 out of 3 runs** per combination of opponent and dataset.

## Summary

| Case ID | Opponent          | Datasets        | Success Condition              |
|---------|-------------------|-----------------|--------------------------------|
| 01      | `big-range`       | Data 4, Data 5  | Win ≥ 2/3 runs per dataset     |
| 02      | `linear-regr`     | Data 4, Data 5  | Win ≥ 2/3 runs per dataset     |
| 03      | `correlation-coef`| Data 4, Data 5  | Win ≥ 2/3 runs per dataset     |
| 04      | `mse` (Bonus)     | Data 4, Data 5  | Win ≥ 2/3 runs per dataset     |
| 05      | `nic` (Bonus)     | Data 4, Data 5  | Win ≥ 2/3 runs per dataset     |
| 06      | Packaging         | N/A             | Tester executes `script.sh` successfully |

---

## Audit Case 1: `big-range` Opponent
**Description:** A wide, safe opponent. To win, the student must provide a statistically valid but significantly tighter range than the opponent's default safety net.

- **Execution Environment:** Dockerized tester from `guess-it-dockerized.zip`.
- **Opponent:** `?guesser=big-range`
- **Datasets:** `Data 4`, `Data 5`.
- **Procedure:** Run the test 3 times per dataset.
- **Expected Behavior:** Student score > opponent score in ≥ 2 of 3 runs for each dataset.

---

## Audit Case 2: `linear-regr` Opponent
**Description:** Predicts the next value via a linear regression line fitted on the stream. The student must out-perform it by also leveraging regression *and* widening defensively when the trend breaks.

- **Opponent:** `?guesser=linear-regr`
- **Datasets:** `Data 4`, `Data 5`.
- **Procedure:** Run the test 3 times per dataset.
- **Expected Behavior:** Student score > opponent score in ≥ 2 of 3 runs for each dataset.
- **Why the fixed-range predictor wins:** the student centres on the same regression line, but its constant ±20 range is far narrower than a defensively-sized opponent range, so each correct guess scores higher.

---

## Audit Case 3: `correlation-coef` Opponent
**Description:** Modulates its range using the Pearson correlation coefficient. The student wins by holding a constant narrow ±20 range — score-per-hit stays high regardless of how the opponent sizes its band.

- **Opponent:** `?guesser=correlation-coef`
- **Datasets:** `Data 4`, `Data 5`.
- **Procedure:** Run the test 3 times per dataset.
- **Expected Behavior:** Student score > opponent score in ≥ 2 of 3 runs for each dataset.

---

## Audit Case 4: `mse` Opponent (Bonus)
**Description:** Uses mean-squared error to size the range.

- **Opponent:** `?guesser=mse`
- **Datasets:** `Data 4`, `Data 5`.
- **Procedure:** Run the test 3 times per dataset.
- **Expected Behavior:** Student score > opponent score in ≥ 2 of 3 runs for each dataset.

---

## Audit Case 5: `nic` Opponent (Bonus)
**Description:** Reference adaptive opponent. The student relies on the fixed ±20 width keeping score-per-hit high enough to out-score an adaptively-sized opponent range.

- **Opponent:** `?guesser=nic`
- **Datasets:** `Data 4`, `Data 5`.
- **Procedure:** Run the test 3 times per dataset.
- **Expected Behavior:** Student score > opponent score in ≥ 2 of 3 runs for each dataset.

---

## Audit Case 6: Packaging & Execution
**Description:** Verifies the project layout meets the auditor's expectations.

- **Procedure:** The auditor drops `student/` into the tester root and runs `script.sh`.
- **Expected Behavior:**
  - `student/` contains the compiled `guess-it-2` binary and any required assets.
  - `script.sh` is executable and runs the binary without errors.
  - The program reads stdin and writes `lower upper\n` per input line.
- **Failure Condition:** If `script.sh` fails or no output is produced, the audit fails immediately.
