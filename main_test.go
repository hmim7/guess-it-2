package main

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
)

// TestRun is a table-driven exercise of the streaming loop covering the
// happy path, invalid input, EOF, magnitude, and both centre modes.
func TestRun(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		centerMode    string
		expectedLines int
		coverNext     bool // assert each emitted range contains the next stream value
	}{
		{
			name:          "Edge 06: ignores invalid lines (average centre)",
			input:         "100\ninvalid_string\n200\n",
			centerMode:    "",
			expectedLines: 2,
		},
		{
			name:          "Edge 07: empty input stream",
			input:         "",
			centerMode:    "average",
			expectedLines: 0,
		},
		{
			name:          "Edge 11: large magnitude numbers",
			input:         "9999999999999\n9999999999999\n",
			centerMode:    "median",
			expectedLines: 2,
		},
		{
			name: "Trended stream (median centre): predicted ranges cover the next input",
			input: strings.Join([]string{
				"100", "101", "102", "103", "104", "105",
				"106", "107", "108", "109", "110", "111",
			}, "\n") + "\n",
			centerMode:    "median",
			expectedLines: 12,
			coverNext:     true,
		},
		{
			name: "Oscillating stream (median centre): predicted ranges cover the next input",
			input: strings.Join([]string{
				"10", "100", "10", "100", "10", "100",
				"10", "100", "10", "100", "10", "100",
			}, "\n") + "\n",
			centerMode:    "median",
			expectedLines: 12,
			coverNext:     true,
		},
		{
			name:          "Blank lines are skipped",
			input:         "\n\n100\n\n200\n\n",
			centerMode:    "average",
			expectedLines: 2,
		},
		{
			name:          "PREDICT_CENTER unset falls back to average path",
			input:         "10\n20\n30\n40\n50\n60\n",
			centerMode:    "",
			expectedLines: 6,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			if err := run(strings.NewReader(tc.input), &out, tc.centerMode); err != nil {
				t.Fatalf("run: unexpected error: %v", err)
			}

			outStr := strings.TrimRight(out.String(), "\n")
			if tc.expectedLines == 0 {
				if outStr != "" {
					t.Fatalf("want no output, got %q", outStr)
				}
				return
			}

			lines := strings.Split(outStr, "\n")
			if len(lines) != tc.expectedLines {
				t.Fatalf("expected %d output lines, got %d (out=%q)",
					tc.expectedLines, len(lines), outStr)
			}

			if !tc.coverNext {
				return
			}

			inputs := strings.Split(strings.TrimRight(tc.input, "\n"), "\n")
			for i, line := range lines {
				if i+1 >= len(inputs) {
					break
				}
				// The first 4 predictions use the wide initial range; only
				// assert coverage once the statistical model is active.
				if i < 4 {
					continue
				}
				lo, hi, err := parseRange(line)
				if err != nil {
					t.Errorf("line %d: %v (line=%q)", i, err, line)
					continue
				}
				next, err := strconv.ParseFloat(inputs[i+1], 64)
				if err != nil {
					t.Fatalf("parse next input %q: %v", inputs[i+1], err)
				}
				if next < float64(lo) || next > float64(hi) {
					t.Errorf("range [%d,%d] does not cover next input %v (step %d)",
						lo, hi, next, i)
				}
			}
		})
	}
}

func parseRange(line string) (int64, int64, error) {
	parts := strings.Fields(line)
	if len(parts) != 2 {
		return 0, 0, &parseRangeError{line: line, reason: "want two fields"}
	}
	lo, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	hi, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return lo, hi, nil
}

type parseRangeError struct {
	line   string
	reason string
}

func (e *parseRangeError) Error() string {
	return "parseRange(" + strconv.Quote(e.line) + "): " + e.reason
}
