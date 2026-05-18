package main

import (
	"bytes"
	"strings"
	"testing"
)

// --- I/O Edge Cases ---

func TestRun(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedOut   string
		expectedLines int
	}{
		{
			name:          "Edge 06: Ignores invalid lines",
			input:         "100\ninvalid_string\n200\n",
			expectedLines: 2,
		},
		{
			name:          "Edge 07: Empty input stream",
			input:         "",
			expectedLines: 0,
			expectedOut:   "",
		},
		{
			name:          "Edge 11: Large Magnitude Numbers",
			input:         "9999999999999\n9999999999999\n",
			expectedLines: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			in := strings.NewReader(tc.input)
			var out bytes.Buffer
			err := run(in, &out, "average")
			if err != nil {
				t.Fatalf("run failed unexpectedly: %v", err)
			}

			outStr := strings.TrimSpace(out.String())
			if tc.expectedLines == 0 {
				if outStr != tc.expectedOut {
					t.Errorf("expected no output, got %q", outStr)
				}
			} else {
				lines := strings.Split(outStr, "\n")
				if len(lines) != tc.expectedLines {
					t.Errorf("expected exactly %d output lines, got %d", tc.expectedLines, len(lines))
				}
			}
		})
	}
}
