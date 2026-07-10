package classic

import (
	"github.com/philoserf/traveller-worldgen/dice"
	"github.com/philoserf/traveller-worldgen/worldname"
)

// Generate rolls a single world from r. Draws happen in a fixed order —
// starport, naval base, scout base, size, atmosphere, hydrographics,
// population, government, law level, tech level, name — so a given seed always
// reproduces the same world. Bases and the size-0/1 automatic results consume
// no die when they do not apply.
func Generate(r dice.Roller) World {
	var w World
	w.Starport = rollStarport(r)
	w.NavalBase = rollNavalBase(r, w.Starport)
	w.ScoutBase = rollScoutBase(r, w.Starport)
	w.Size = floor0(r.D6(2) - 2)
	w.Atmosphere = rollAtmosphere(r, w.Size)
	w.Hydrographics = rollHydrographics(r, w.Size, w.Atmosphere)
	w.Population = floor0(r.D6(2) - 2)
	w.Government = floor0(r.D6(2) - 7 + w.Population)
	w.LawLevel = floor0(r.D6(2) - 7 + w.Government)
	w.TechLevel = rollTech(r, w)
	w.Name = worldname.Generate(r)
	return w
}

// rollStarport rolls 2D on the Book 3 Starports table.
func rollStarport(r dice.Roller) byte {
	switch r.D6(2) {
	case 2, 3, 4:
		return 'A'
	case 5, 6:
		return 'B'
	case 7, 8:
		return 'C'
	case 9:
		return 'D'
	case 10, 11:
		return 'E'
	default: // 12
		return 'X'
	}
}

// rollNavalBase rolls for a naval base, which only class A and B starports can
// have (2D >= 8).
func rollNavalBase(r dice.Roller, starport byte) bool {
	if starport != 'A' && starport != 'B' {
		return false
	}
	return r.D6(2) >= 8
}

// rollScoutBase rolls for a scout base against the per-starport target
// (A:10+, B:9+, C:8+, D:7+); E and X starports never have one.
func rollScoutBase(r dice.Roller, starport byte) bool {
	var target int
	switch starport {
	case 'A':
		target = 10
	case 'B':
		target = 9
	case 'C':
		target = 8
	case 'D':
		target = 7
	default:
		return false
	}
	return r.D6(2) >= target
}

// rollAtmosphere rolls atmosphere: 2D-7+size, floored at 0. Size 0 forces
// atmosphere 0 with no roll. Large worlds can roll past the last described
// type (C); the value is kept and rendered in extended hex.
func rollAtmosphere(r dice.Roller, size int) int {
	if size == 0 {
		return 0
	}
	return floor0(r.D6(2) - 7 + size)
}

// rollHydrographics rolls hydrographics: 2D-7+size, with a -4 DM when the
// atmosphere is 0, 1, or A+ (exotic/corrosive/insidious), clamped to 0-A (a
// literal percentage). Size 0 or 1 forces 0 with no roll.
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

// rollTech rolls the technological index: 1D plus the summed matrix DMs for the
// world's characteristics, floored at 0.
func rollTech(r dice.Roller, w World) int {
	dm := starportTechDM[w.Starport] +
		sizeTechDM[w.Size] +
		atmoTechDM[w.Atmosphere] +
		hydroTechDM[w.Hydrographics] +
		popTechDM[w.Population] +
		govTechDM[w.Government]
	return floor0(r.D6(1) + dm)
}

// floor0 clamps a value to a minimum of 0 (Book 3 treats negative DM results
// as 0).
func floor0(v int) int { return max(0, v) }

// clamp constrains v to the inclusive range [lo, hi].
func clamp(v, lo, hi int) int { return max(lo, min(v, hi)) }
