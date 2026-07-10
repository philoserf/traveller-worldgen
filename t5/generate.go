package t5

import "github.com/philoserf/traveller-worldgen/dice"

// Generate rolls a single T5 mainworld from r. Draws happen in a fixed,
// documented order so a given seed always reproduces the same world; a
// conditional die is consumed only when its branch applies (see
// docs/t5/world-generation.md and the determinism contract in CLAUDE.md).
//
// The T5 core draw is Flux (1D-1D, -5..+5), not the 2D-7 of the other editions;
// hydrographics keys off atmosphere rather than size; and characteristics are
// hard-capped rather than kept raw.
func Generate(r dice.Roller) World {
	var w World
	w.Starport = rollStarport(r)
	w.NavalBase = rollBase(r, navalTarget[w.Starport])
	w.ScoutBase = rollBase(r, scoutTarget[w.Starport])
	w.Size = rollSize(r)
	w.Atmosphere = rollAtmosphere(r, w.Size)
	w.Hydrographics = rollHydrographics(r, w.Size, w.Atmosphere)
	w.Population = rollPopulation(r)
	w.Government = capAt(flux(r)+w.Population, 15)
	w.LawLevel = capAt(flux(r)+w.Government, 18)
	w.TechLevel = rollTech(r, w)
	w.PopulationDigit = rollPopulationDigit(r, w.Population)
	w.PlanetoidBelts = floor0(r.D6(1) - 3)
	w.GasGiants = floor0(r.D6(2)/2 - 2)
	// Extensions derive from the finished UWP; Ix is computed first because Ex
	// (Infrastructure) and Cx (Acceptance) both read it.
	w.Importance = importance(w)
	w.Economic = rollEconomic(r, w)
	w.Cultural = rollCultural(r, w)
	w.Name = generateName(r)
	return w
}

// flux returns a T5 Flux value: 1D - 1D, an integer in -5..+5. It always consumes
// exactly two dice, the first as the minuend. The two draws are kept in named
// locals because they have side effects (each advances the roller) — the
// subtraction is Flux, not a constant zero.
func flux(r dice.Roller) int {
	high := r.D6(1)
	low := r.D6(1)
	return high - low
}

// rollStarport rolls 2D on the mainworld starport table.
func rollStarport(r dice.Roller) byte {
	return starportType[capAt(r.D6(2)-2, len(starportType)-1)]
}

// rollBase rolls 2D against a base's per-starport target, present on a roll of
// target or less (T5's "N-" notation). A zero target means the starport can never
// hold the base, and no die is consumed.
func rollBase(r dice.Roller, target int) bool {
	if target == 0 {
		return false
	}
	return r.D6(2) <= target
}

// rollSize rolls world size: 2D-2, rerolling a result of 10 as 1D+9 (giving A-F).
func rollSize(r dice.Roller) int {
	s := r.D6(2) - 2
	if s == 10 {
		s = r.D6(1) + 9
	}
	return floor0(s)
}

// rollAtmosphere rolls atmosphere: Flux + Size, floored at 0 and capped at F.
// Size 0 forces atmosphere 0 with no roll.
func rollAtmosphere(r dice.Roller, size int) int {
	if size == 0 {
		return 0
	}
	return capAt(flux(r)+size, 15)
}

// rollHydrographics rolls hydrographics: Flux + Atmosphere, with a -4 DM when the
// atmosphere is below 2 or above 9, clamped to 0-A. Size below 2 forces 0 with no
// roll. Note T5 keys hydrographics off atmosphere, not size.
func rollHydrographics(r dice.Roller, size, atmosphere int) int {
	if size < 2 {
		return 0
	}
	v := flux(r) + atmosphere
	if atmosphere < 2 || atmosphere > 9 {
		v -= 4
	}
	return capAt(v, 10)
}

// rollPopulation rolls the population exponent: 2D-2, rerolling a result of 10 as
// 2D+3 (giving A-F).
func rollPopulation(r dice.Roller) int {
	p := r.D6(2) - 2
	if p == 10 {
		p = r.D6(2) + 3
	}
	return floor0(p)
}

// rollTech rolls the tech level: 1D plus the summed matrix DMs, floored at 0.
func rollTech(r dice.Roller, w World) int {
	dm := starportTechDM[w.Starport] +
		sizeTechDM[w.Size] +
		atmoTechDM[w.Atmosphere] +
		hydroTechDM[w.Hydrographics] +
		popTechDM(w.Population) +
		govTechDM[w.Government]
	return floor0(r.D6(1) + dm)
}

// rollPopulationDigit rolls the PBG population digit: an even 1-9 when the world
// is populated, else 0 with no roll. Two independent d6 form a uniform index in
// [0, 36) reduced modulo 9 (which divides 36 evenly, so the result is unbiased).
func rollPopulationDigit(r dice.Roller, pop int) int {
	if pop == 0 {
		return 0
	}
	idx := (r.D6(1)-1)*6 + (r.D6(1) - 1)
	return idx%9 + 1
}

// floor0 clamps a value to a minimum of 0.
func floor0(v int) int { return max(0, v) }

// capAt floors v at 0 and caps it at hi. Every T5 characteristic clamps to a
// [0, hi] range, so this replaces the general three-argument clamp of the other
// editions.
func capAt(v, hi int) int { return min(floor0(v), hi) }
