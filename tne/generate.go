package tne

import (
	"fmt"
	"slices"

	"github.com/philoserf/traveller-worldgen/dice"
)

// Nature is the referee-chosen character of a subsector, which selects the
// column of the starport table (TNE Basic Mainworld Generation).
type Nature int

// The four subsector natures, in the column order of the starport table.
const (
	Backwater Nature = iota
	Standard
	Mature
	Cluster
)

// natureNames holds each Nature's lowercase name, indexed by the Nature value.
// The order is canonical (matching the starport-table columns), so the slice
// doubles as the source of truth for the CLI's valid-nature list.
var natureNames = [...]string{
	Backwater: "backwater",
	Standard:  "standard",
	Mature:    "mature",
	Cluster:   "cluster",
}

// String returns the lowercase name of the nature.
func (n Nature) String() string {
	if int(n) < 0 || int(n) >= len(natureNames) {
		return fmt.Sprintf("Nature(%d)", int(n))
	}
	return natureNames[n]
}

// ParseNature parses a nature name (case-sensitive lowercase) and reports
// whether it was recognized.
func ParseNature(s string) (Nature, bool) {
	for i, name := range natureNames {
		if name == s {
			return Nature(i), true
		}
	}
	return 0, false
}

// NatureNames returns the valid subsector-nature names in canonical order, so
// CLI help and validation messages draw from the same source of truth.
func NatureNames() []string {
	return slices.Clone(natureNames[:])
}

// Generate rolls a single world from r for the given subsector nature. Draws
// happen in a fixed order — starport, naval base, scout base, military base,
// size, atmosphere, hydrographics, population, government, law level, tech level,
// name — so a given seed always reproduces the same world. Bases and the
// size-0/1 automatic results consume no die when they do not apply.
func Generate(r dice.Roller, nature Nature) World {
	var w World
	w.Starport = rollStarport(r, nature)
	throws := baseThrows[w.Starport]
	w.NavalBase = rollBase(r, throws.naval)
	w.ScoutBase = rollBase(r, throws.scout)
	w.MilitaryBase = rollBase(r, throws.military)
	w.Size = floor0(r.D6(2) - 2)
	w.Atmosphere = rollAtmosphere(r, w.Size)
	w.Hydrographics = rollHydrographics(r, w.Size, w.Atmosphere)
	w.Population = floor0(r.D6(2) - 2)
	w.Government = floor0(r.D6(2) - 7 + w.Population)
	w.LawLevel = floor0(r.D6(2) - 7 + w.Government)
	w.TechLevel = rollTech(r, w)
	w.Name = generateName(r)
	return w
}

// rollStarport rolls 2D on the starport column for the subsector nature.
func rollStarport(r dice.Roller, nature Nature) byte {
	col, ok := starportByNature[nature]
	if !ok {
		col = starportByNature[Standard]
	}
	return col[clamp(r.D6(2)-2, 0, len(col)-1)]
}

// rollBase rolls 2D against a base's per-starport target. A zero target means
// the starport can never have that base, and no die is consumed.
func rollBase(r dice.Roller, target int) bool {
	if target == 0 {
		return false
	}
	return r.D6(2) >= target
}

// rollAtmosphere rolls atmosphere: 2D-7+size, floored at 0. Size 0 forces
// atmosphere 0 with no roll.
func rollAtmosphere(r dice.Roller, size int) int {
	if size == 0 {
		return 0
	}
	return floor0(r.D6(2) - 7 + size)
}

// rollHydrographics rolls hydrographics: 2D-7+size, with a -4 DM when the
// atmosphere is 0, 1, or A+ (exotic and worse), clamped to 0-A (a literal
// percentage). Size 0 or 1 forces 0 with no roll.
func rollHydrographics(r dice.Roller, size, atmosphere int) int {
	if size <= 1 {
		return 0
	}
	v := r.D6(2) - 7 + size
	if atmosphere <= 1 || atmosphere >= 10 {
		v -= 4
	}
	return clamp(v, 0, 10)
}

// rollTech rolls the tech level: 1D plus the summed matrix DMs for the world's
// characteristics, floored at 0.
func rollTech(r dice.Roller, w World) int {
	dm := starportTechDM[w.Starport] +
		sizeTechDM[w.Size] +
		atmoTechDM[w.Atmosphere] +
		hydroTechDM[w.Hydrographics] +
		popTechDM[w.Population] +
		govTechDM[w.Government]
	return floor0(r.D6(1) + dm)
}

// floor0 clamps a value to a minimum of 0 (negative DM results are treated as 0).
func floor0(v int) int { return max(0, v) }

// clamp constrains v to the inclusive range [lo, hi].
func clamp(v, lo, hi int) int { return max(lo, min(v, hi)) }
