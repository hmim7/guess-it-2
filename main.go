package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"guess-it-2/linearstats"
)

// main runs the prediction loop and exits with a non-zero status on error.
func main() {
	if err := run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// run reads numbers from r line-by-line and writes the predicted bounds for
// each next value to w.
func run(r io.Reader, w io.Writer) error {
	in := bufio.NewScanner(r)
	in.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	out := bufio.NewWriter(w)
	defer func() { _ = out.Flush() }()

	predictor := linearstats.NewPredictor()

	for in.Scan() {
		line := strings.TrimSpace(in.Text())
		if line == "" {
			continue
		}
		v, err := strconv.ParseFloat(line, 64)
		if err != nil {
			continue
		}

		lower, upper := predictor.Next(v)
		if _, err := fmt.Fprintf(out, "%d %d\n", lower, upper); err != nil {
			return err
		}
		if err := out.Flush(); err != nil {
			return err
		}
	}

	if err := in.Err(); err != nil {
		return err
	}
	return nil
}
