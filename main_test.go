package main

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
)

// TestRun is a table-driven exercise of the streaming loop: invalid input is
// skipped, blank lines are ignored, EOF exits cleanly, and every numeric line
// produces one well-formed "lower upper" prediction.
func TestRun(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		expectedLines int
	}{
		{
			name:          "ignores invalid lines",
			input:         "100\ninvalid_string\n200\n",
			expectedLines: 2,
		},
		{
			name:          "empty input stream",
			input:         "",
			expectedLines: 0,
		},
		{
			name:          "large magnitude numbers",
			input:         "9999999999999\n9999999999999\n",
			expectedLines: 2,
		},
		{
			name: "trended stream",
			input: strings.Join([]string{
				"100", "101", "102", "103", "104", "105",
				"106", "107", "108", "109", "110", "111",
			}, "\n") + "\n",
			expectedLines: 12,
		},
		{
			name: "oscillating stream",
			input: strings.Join([]string{
				"10", "100", "10", "100", "10", "100",
				"10", "100", "10", "100", "10", "100",
			}, "\n") + "\n",
			expectedLines: 12,
		},
		{
			name:          "blank lines are skipped",
			input:         "\n\n100\n\n200\n\n",
			expectedLines: 2,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			if err := run(strings.NewReader(tc.input), &out); err != nil {
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
			// Every line must be two integers, lower <= upper.
			for i, line := range lines {
				lo, hi, err := parseRange(line)
				if err != nil {
					t.Errorf("line %d: %v (line=%q)", i, err, line)
					continue
				}
				if lo > hi {
					t.Errorf("line %d: lower %d > upper %d", i, lo, hi)
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
