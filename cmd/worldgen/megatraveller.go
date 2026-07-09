package main

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/philoserf/traveller-worldgen/dice"
	"github.com/philoserf/traveller-worldgen/megatraveller"
)

// runMega generates MegaTraveller mainworlds, adding the -nature flag to the
// shared edition runner.
func runMega(args []string, stdout, stderr io.Writer) int {
	return runEdition("mega", args, stdout, stderr,
		func(fs *flag.FlagSet) func() (func(dice.Roller) megatraveller.World, error) {
			natureList := strings.Join(megatraveller.NatureNames(), ", ")
			nature := fs.String("nature", "standard", "subsector nature: "+natureList)
			return func() (func(dice.Roller) megatraveller.World, error) {
				natureVal, ok := megatraveller.ParseNature(*nature)
				if !ok {
					return nil, fmt.Errorf("unknown -nature %q (want %s)", *nature, natureList)
				}
				return func(r dice.Roller) megatraveller.World { return megatraveller.Generate(r, natureVal) }, nil
			}
		},
		megaUWPLine, megaTextBlock)
}

// megaUWPLine renders one world as a single canonical library-data line: name,
// UWP, base code, gas-giant marker ("G" plus the count), planetoid-belt marker
// ("B" plus the count), then the trade codes. Empty fields render as "-". Trade
// codes come last because they are the only variable-length field, so a long
// list never pushes another column out of alignment. Sentinels are ASCII so
// fixed-width padding stays aligned.
func megaUWPLine(w megatraveller.World) string {
	base := hyphenIfEmpty(w.BaseCode())
	gas := "-"
	if w.GasGiants > 0 {
		gas = fmt.Sprintf("G%d", w.GasGiants)
	}
	belts := "-"
	if w.PlanetoidBelts > 0 {
		belts = fmt.Sprintf("B%d", w.PlanetoidBelts)
	}
	trade := hyphenIfEmpty(strings.Join(w.TradeCodes(), " "))
	return fmt.Sprintf("%-14s %s  %-3s %-3s %-3s %s", w.Name, w.UWP(), base, gas, belts, trade)
}

// megaTextBlock renders a world as a header line, a per-characteristic breakdown,
// and trailing Trade, Gas Giants, and Planetoid Belts lines.
func megaTextBlock(w megatraveller.World) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s  %s  %s\n", w.Name, w.UWP(), w.Bases())
	for _, c := range w.Profile() {
		fmt.Fprintf(&b, "  %-13s %c  %s\n", c.Label, c.Code, c.Desc)
	}
	fmt.Fprintf(&b, "  %-13s %s\n", "Trade", hyphenIfEmpty(strings.Join(w.TradeCodes(), " ")))
	fmt.Fprintf(&b, "  %-13s %d\n", "Gas Giants", w.GasGiants)
	fmt.Fprintf(&b, "  %-13s %d\n", "Belts", w.PlanetoidBelts)
	return b.String()
}

// hyphenIfEmpty returns s, or "-" when s is empty.
func hyphenIfEmpty(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
