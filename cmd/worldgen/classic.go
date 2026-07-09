package main

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/philoserf/traveller-worldgen/classic"
	"github.com/philoserf/traveller-worldgen/dice"
)

// runClassic generates Classic Traveller (Book 3) worlds on the shared edition
// runner. It has no edition-specific flags, so its generator is classic.Generate
// directly.
func runClassic(args []string, stdout, stderr io.Writer) int {
	return runEdition("classic", args, stdout, stderr,
		func(*flag.FlagSet) func() (func(dice.Roller) classic.World, error) {
			return func() (func(dice.Roller) classic.World, error) { return classic.Generate, nil }
		},
		uwpLine, textBlock)
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
