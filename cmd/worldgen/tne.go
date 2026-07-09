package main

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/philoserf/traveller-worldgen/dice"
	"github.com/philoserf/traveller-worldgen/tne"
)

// runTne generates Traveller: The New Era mainworlds, adding the -nature flag to
// the shared edition runner.
func runTne(args []string, stdout, stderr io.Writer) int {
	return runEdition("tne", args, stdout, stderr,
		func(fs *flag.FlagSet) func() (func(dice.Roller) tne.World, error) {
			natureList := strings.Join(tne.NatureNames(), ", ")
			nature := fs.String("nature", "standard", "subsector nature: "+natureList)
			return func() (func(dice.Roller) tne.World, error) {
				natureVal, ok := tne.ParseNature(*nature)
				if !ok {
					return nil, fmt.Errorf("unknown -nature %q (want %s)", *nature, natureList)
				}
				return func(r dice.Roller) tne.World { return tne.Generate(r, natureVal) }, nil
			}
		},
		tneUWPLine, tneTextBlock)
}

// tneUWPLine renders one world as a single line: name, UWP, and the compact base
// code (N/S/A plus a trailing M for a military base, "—" when none).
func tneUWPLine(w tne.World) string {
	code := w.BaseCode()
	if code == "" {
		code = "—"
	}
	return fmt.Sprintf("%-14s %s  %s", w.Name, w.UWP(), code)
}

// tneTextBlock renders a world as a header line plus a per-characteristic
// breakdown.
func tneTextBlock(w tne.World) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s  %s  %s\n", w.Name, w.UWP(), w.Bases())
	for _, c := range w.Profile() {
		fmt.Fprintf(&b, "  %-13s %c  %s\n", c.Label, c.Code, c.Desc)
	}
	return b.String()
}
