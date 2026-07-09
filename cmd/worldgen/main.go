// Command worldgen generates a single Classic Traveller (Book 3) world.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/philoserf/traveller-worldgen/classic"
	"github.com/philoserf/traveller-worldgen/dice"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run parses args, generates worlds, and writes output. It returns a process
// exit code: 0 on success, 1 on an output error, 2 on a usage error.
func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("worldgen", flag.ContinueOnError)
	fs.SetOutput(stderr)

	seed := fs.Int64("seed", 0, "random seed (default: time-based)")
	format := fs.String("format", "text", "output format: text, uwp, or json")
	n := fs.Int("n", 1, "number of independent worlds to generate")

	if err := fs.Parse(args); err != nil {
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

	out, err := render(*format, worlds)
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

// render turns the generated worlds into the requested output format.
func render(format string, worlds []classic.World) (string, error) {
	switch format {
	case "json":
		return renderJSON(worlds)
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

// renderJSON marshals one world as an object, or several as an array.
func renderJSON(worlds []classic.World) (string, error) {
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

// errf writes a diagnostic to w, deliberately ignoring the write error (there
// is nowhere better to report a failure to write an error message).
func errf(w io.Writer, format string, a ...any) {
	_, _ = fmt.Fprintf(w, format, a...)
}
