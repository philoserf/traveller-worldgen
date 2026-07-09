package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/philoserf/traveller-worldgen/dice"
	"github.com/philoserf/traveller-worldgen/megatraveller"
)

// runMega generates MegaTraveller mainworlds. It returns a process exit code:
// 0 on success, 1 on an output error, 2 on a usage error.
func runMega(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("worldgen mega", flag.ContinueOnError)
	fs.SetOutput(io.Discard) // suppress flag's auto-print; we route help/errors ourselves

	natureList := strings.Join(megatraveller.NatureNames(), ", ")
	seed := fs.Int64("seed", 0, "random seed (default: time-based)")
	format := fs.String("format", "text", "output format: text, uwp, or json")
	n := fs.Int("n", 1, "number of independent worlds to generate")
	nature := fs.String("nature", "standard", "subsector nature: "+natureList)

	if err := fs.Parse(args); err != nil {
		// -h/-help is a help request, not a usage error: print to stdout, exit 0
		// (matching the top-level `worldgen --help`).
		if errors.Is(err, flag.ErrHelp) {
			fs.SetOutput(stdout)
			fs.Usage()
			return 0
		}
		errf(stderr, "worldgen mega: %v\n", err)
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
	natureVal, ok := megatraveller.ParseNature(*nature)
	if !ok {
		errf(stderr, "worldgen: unknown -nature %q (want %s)\n", *nature, natureList)
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
	worlds := make([]megatraveller.World, *n)
	for i := range worlds {
		worlds[i] = megatraveller.Generate(roller, natureVal)
	}

	out, err := renderMega(*format, worlds)
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

// renderMega turns the generated worlds into the requested output format.
func renderMega(format string, worlds []megatraveller.World) (string, error) {
	switch format {
	case "json":
		return renderWorldsJSON(worlds)
	case "uwp":
		var b strings.Builder
		for _, w := range worlds {
			b.WriteString(megaUWPLine(w))
			b.WriteByte('\n')
		}
		return b.String(), nil
	default: // text
		var b strings.Builder
		for i, w := range worlds {
			if i > 0 {
				b.WriteByte('\n')
			}
			b.WriteString(megaTextBlock(w))
		}
		return b.String(), nil
	}
}

// megaUWPLine renders one world as a single canonical library-data line: name,
// UWP, base code, gas-giant marker ("G" plus the count), then the trade codes.
// Empty fields render as "-". Trade codes come last because they are the only
// variable-length field, so a long list never pushes another column out of
// alignment. Sentinels are ASCII so fixed-width padding stays aligned.
func megaUWPLine(w megatraveller.World) string {
	base := hyphenIfEmpty(w.BaseCode())
	gas := "-"
	if w.GasGiants > 0 {
		gas = fmt.Sprintf("G%d", w.GasGiants)
	}
	trade := hyphenIfEmpty(strings.Join(w.TradeCodes(), " "))
	return fmt.Sprintf("%-14s %s  %-3s %-3s %s", w.Name, w.UWP(), base, gas, trade)
}

// megaTextBlock renders a world as a header line, a per-characteristic breakdown,
// and trailing Trade and Gas Giants lines.
func megaTextBlock(w megatraveller.World) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s  %s  %s\n", w.Name, w.UWP(), w.Bases())
	for _, c := range w.Profile() {
		fmt.Fprintf(&b, "  %-13s %c  %s\n", c.Label, c.Code, c.Desc)
	}
	fmt.Fprintf(&b, "  %-13s %s\n", "Trade", hyphenIfEmpty(strings.Join(w.TradeCodes(), " ")))
	fmt.Fprintf(&b, "  %-13s %d\n", "Gas Giants", w.GasGiants)
	return b.String()
}

// hyphenIfEmpty returns s, or "-" when s is empty.
func hyphenIfEmpty(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
