package main

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/philoserf/traveller-worldgen/dice"
	"github.com/philoserf/traveller-worldgen/t5"
)

// runT5 generates Traveller5 (Core Book 3) mainworlds on the shared edition
// runner. Like classic it has no edition-specific flags, so its generator is
// t5.Generate directly.
func runT5(args []string, stdout, stderr io.Writer) int {
	return runEdition("t5", args, stdout, stderr,
		func(*flag.FlagSet) func() (func(dice.Roller) t5.World, error) {
			return func() (func(dice.Roller) t5.World, error) { return t5.Generate, nil }
		},
		t5UWPLine, t5TextBlock)
}

// t5UWPLine renders one world as a single library-data line: name, UWP, base
// code, the Ix/Ex/Cx extensions, the PBG code, then the trade codes. Empty fields
// render as "-"; trade codes come last as the only variable-length field.
func t5UWPLine(w t5.World) string {
	base := hyphenIfEmpty(w.BaseCode())
	trade := hyphenIfEmpty(strings.Join(w.TradeCodes(), " "))
	return fmt.Sprintf("%-14s %s  %-2s %s %s %s  %s  %s",
		w.Name, w.UWP(), base, w.ImportanceString(), w.Economic.String(), w.Cultural.String(), w.PBG(), trade)
}

// t5TextBlock renders a world as a header line, a per-characteristic breakdown,
// and trailing Trade, extension, and PBG lines.
func t5TextBlock(w t5.World) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s  %s  %s\n", w.Name, w.UWP(), w.Bases())
	for _, c := range w.Profile() {
		fmt.Fprintf(&b, "  %-13s %c  %s\n", c.Label, c.Code, c.Desc)
	}
	fmt.Fprintf(&b, "  %-13s %s\n", "Trade", hyphenIfEmpty(strings.Join(w.TradeCodes(), " ")))
	fmt.Fprintf(&b, "  %-13s %s\n", "Importance", w.ImportanceString())
	fmt.Fprintf(&b, "  %-13s %s  (RU %d)\n", "Economic", w.Economic.String(), w.Economic.RU())
	fmt.Fprintf(&b, "  %-13s %s\n", "Cultural", w.Cultural.String())
	fmt.Fprintf(&b, "  %-13s %d\n", "Gas Giants", w.GasGiants)
	fmt.Fprintf(&b, "  %-13s %d\n", "Belts", w.PlanetoidBelts)
	return b.String()
}
