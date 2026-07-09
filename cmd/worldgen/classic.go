package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/philoserf/traveller-worldgen/classic"
	"github.com/philoserf/traveller-worldgen/dice"
)

// runClassic generates Classic Traveller (Book 3) worlds. It returns a process
// exit code: 0 on success, 1 on an output error, 2 on a usage error.
func runClassic(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("worldgen classic", flag.ContinueOnError)
	fs.SetOutput(io.Discard) // suppress flag's auto-print; we route help/errors ourselves

	seed := fs.Int64("seed", 0, "random seed (default: time-based)")
	format := fs.String("format", "text", "output format: text, uwp, or json")
	n := fs.Int("n", 1, "number of independent worlds to generate")

	if err := fs.Parse(args); err != nil {
		// -h/-help is a help request, not a usage error: print to stdout, exit 0
		// (matching the top-level `worldgen --help`).
		if errors.Is(err, flag.ErrHelp) {
			fs.SetOutput(stdout)
			fs.Usage()
			return 0
		}
		errf(stderr, "worldgen classic: %v\n", err)
		return 2
	}
	if *n < 1 {
		errf(stderr, "worldgen: -n must be >= 1, got %d\n", *n)
		return 2
	}
	switch *format {
	case "text", "uwp", "json":
	default:
		errf(stderr, "worldgen: unknown -format %q (want text, uwp, or json)\n", *format)
		return 2
	}

	// An explicit -seed (even -seed 0) is honored; otherwise seed from the clock.
	seeded := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "seed" {
			seeded = true
		}
	})
	if !seeded {
		*seed = time.Now().UnixNano()
	}

	// One roller drives all N worlds, so a seed reproduces the whole batch.
	roller := dice.NewSeeded(*seed)
	worlds := make([]classic.World, *n)
	for i := range worlds {
		worlds[i] = classic.Generate(roller)
	}

	out, err := renderClassic(*format, worlds)
	if err != nil {
		errf(stderr, "worldgen: %v\n", err)
		return 1
	}
	if _, err := io.WriteString(stdout, out); err != nil {
		errf(stderr, "worldgen: writing output: %v\n", err)
		return 1
	}
	return 0
}

// renderClassic turns the generated worlds into the requested output format.
func renderClassic(format string, worlds []classic.World) (string, error) {
	switch format {
	case "json":
		return renderClassicJSON(worlds)
	case "uwp":
		var b strings.Builder
		for _, w := range worlds {
			b.WriteString(uwpLine(w))
			b.WriteByte('\n')
		}
		return b.String(), nil
	default: // text
		var b strings.Builder
		for i, w := range worlds {
			if i > 0 {
				b.WriteByte('\n')
			}
			b.WriteString(textBlock(w))
		}
		return b.String(), nil
	}
}

// renderClassicJSON marshals one world as an object, or several as an array.
func renderClassicJSON(worlds []classic.World) (string, error) {
	var (
		data []byte
		err  error
	)
	if len(worlds) == 1 {
		data, err = json.MarshalIndent(worlds[0], "", "  ")
	} else {
		data, err = json.MarshalIndent(worlds, "", "  ")
	}
	if err != nil {
		return "", fmt.Errorf("encoding json: %w", err)
	}
	return string(data) + "\n", nil
}

// uwpLine renders one world as a single line: name, UWP, and a compact base
// code (N, S, NS, or —).
func uwpLine(w classic.World) string {
	code := ""
	if w.NavalBase {
		code += "N"
	}
	if w.ScoutBase {
		code += "S"
	}
	if code == "" {
		code = "—"
	}
	return fmt.Sprintf("%-14s %s  %s", w.Name, w.UWP(), code)
}

// textBlock renders a world as a header line plus a per-characteristic
// breakdown.
func textBlock(w classic.World) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s  %s  %s\n", w.Name, w.UWP(), w.Bases())
	for _, c := range w.Profile() {
		fmt.Fprintf(&b, "  %-13s %c  %s\n", c.Label, c.Code, c.Desc)
	}
	return b.String()
}
