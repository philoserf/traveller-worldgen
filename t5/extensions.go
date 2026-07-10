package t5

import (
	"fmt"

	"github.com/philoserf/traveller-worldgen/dice"
	"github.com/philoserf/traveller-worldgen/ehex"
)

// Economic is the T5 Economic Extension (Ex), rendered "(RLI±E)". Resources,
// Labor, and Infrastructure are single extended-hex digits; Efficiency is a
// signed value that can be negative (Book 3 p. 27).
type Economic struct {
	Resources      int `json:"resources"`
	Labor          int `json:"labor"`
	Infrastructure int `json:"infrastructure"`
	Efficiency     int `json:"efficiency"`
}

// RU returns the Resource Units: Resources × Labor × Infrastructure × Efficiency.
// A factor of 0 is treated as 1 (per the source, to avoid zeroing the product),
// so RU is negative when Efficiency is.
func (e Economic) RU() int {
	return nonZero(e.Resources) * nonZero(e.Labor) * nonZero(e.Infrastructure) * nonZero(e.Efficiency)
}

// String renders the Economic Extension as "(RLI±E)", e.g. "(D7E+4)".
func (e Economic) String() string {
	return fmt.Sprintf("(%c%c%c%+d)",
		ehex.Encode(e.Resources), ehex.Encode(e.Labor), ehex.Encode(e.Infrastructure), e.Efficiency)
}

// Cultural is the T5 Cultural Extension (Cx), rendered "[HASS]": Heterogeneity,
// Acceptance, Strangeness, and Symbols, each an extended-hex digit (Book 3 p. 27).
type Cultural struct {
	Heterogeneity int `json:"heterogeneity"`
	Acceptance    int `json:"acceptance"`
	Strangeness   int `json:"strangeness"`
	Symbols       int `json:"symbols"`
}

// String renders the Cultural Extension as "[HASS]", e.g. "[9C6D]".
func (c Cultural) String() string {
	return fmt.Sprintf("[%c%c%c%c]",
		ehex.Encode(c.Heterogeneity), ehex.Encode(c.Acceptance),
		ehex.Encode(c.Strangeness), ehex.Encode(c.Symbols))
}

// ImportanceString renders the Importance Extension (Ix) as a signed value in
// braces, e.g. "{+4}" or "{-1}".
func (w World) ImportanceString() string { return fmt.Sprintf("{%+d}", w.Importance) }

// importance computes the Importance Extension (Ix): the signed sum of starport,
// tech-level, trade-code, population, and base modifiers (Book 3 p. 27). The Way
// Station term is always 0 here (Way Stations are deferred).
func importance(w World) int {
	ix := 0
	switch w.Starport {
	case 'A', 'B':
		ix++
	case 'D', 'E', 'X':
		ix--
	}
	if w.TechLevel >= 16 { // TL >= G
		ix++
	}
	if w.TechLevel >= 10 { // TL >= A
		ix++
	}
	if w.TechLevel <= 8 {
		ix--
	}
	for _, code := range w.TradeCodes() {
		switch code {
		case "Ag", "Hi", "In", "Ri":
			ix++
		}
	}
	if w.Population <= 6 {
		ix--
	}
	if w.NavalBase && w.ScoutBase {
		ix++
	}
	return ix
}

// rollEconomic computes the Economic Extension (Ex). It draws Resources (2D),
// Infrastructure (1D for Pop 4-6, 2D for Pop 7+, no die for Pop 0-3), and
// Efficiency (Flux), in that order; Labor is derived. Resources gains the system
// gas-giant and belt counts only at TL 8+.
func rollEconomic(r dice.Roller, w World) Economic {
	resources := r.D6(2)
	if w.TechLevel >= 8 {
		resources += w.GasGiants + w.PlanetoidBelts
	}

	var infrastructure int
	switch {
	case w.Population == 0:
		infrastructure = 0
	case w.Population <= 3:
		infrastructure = w.Importance
	case w.Population <= 6:
		infrastructure = r.D6(1) + w.Importance
	default: // Pop >= 7
		infrastructure = r.D6(2) + w.Importance
	}

	return Economic{
		Resources:      floor0(resources),
		Labor:          floor0(w.Population - 1),
		Infrastructure: floor0(infrastructure),
		Efficiency:     flux(r),
	}
}

// rollCultural computes the Cultural Extension (Cx). An unpopulated world has all
// zeros and draws no dice; otherwise each component is floored at 1. It draws
// Flux for Heterogeneity, Strangeness, and Symbols, in that order; Acceptance is
// derived from the Importance Extension.
func rollCultural(r dice.Roller, w World) Cultural {
	if w.Population == 0 {
		return Cultural{}
	}
	heterogeneity := atLeast1(w.Population + flux(r))
	acceptance := atLeast1(w.Population + w.Importance)
	strangeness := atLeast1(flux(r) + 5)
	symbols := atLeast1(flux(r) + w.TechLevel)
	return Cultural{heterogeneity, acceptance, strangeness, symbols}
}

// nonZero returns v, or 1 when v is 0 (used to keep an RU factor from zeroing the
// product).
func nonZero(v int) int {
	if v == 0 {
		return 1
	}
	return v
}

// atLeast1 clamps a value to a minimum of 1 (the Cx "less than 1 = 1" rule).
func atLeast1(v int) int { return max(1, v) }
