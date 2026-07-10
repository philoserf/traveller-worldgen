package t5

// Trade classifications derived from the UWP (Book 3 p. 26, Trade Classes). Each
// entry pairs a code with a predicate over the world's characteristics; a world
// carries every code whose predicate matches. Only the Planetary, Population, and
// Economic groups are generated — Climate, Secondary, Political, and Special are
// deferred (see the doc). Order follows the source table.
type tradeClass struct {
	code  string
	match func(World) bool
}

var tradeClassifications = []tradeClass{
	// Planetary
	{"As", func(w World) bool {
		return w.Size == 0 && w.Atmosphere == 0 && w.Hydrographics == 0
	}},
	{"De", func(w World) bool {
		return between(w.Atmosphere, 2, 9) && w.Hydrographics == 0
	}},
	{"Fl", func(w World) bool {
		return between(w.Atmosphere, 10, 12) && between(w.Hydrographics, 1, 10)
	}},
	{"Ga", func(w World) bool {
		return inSet(w.Size, 6, 7, 8) && inSet(w.Atmosphere, 5, 6, 8) && between(w.Hydrographics, 5, 7)
	}},
	{"He", func(w World) bool {
		return between(w.Size, 3, 12) && inSet(w.Atmosphere, 2, 4, 7, 9, 10, 11, 12) && between(w.Hydrographics, 0, 2)
	}},
	{"Ic", func(w World) bool {
		return inSet(w.Atmosphere, 0, 1) && between(w.Hydrographics, 1, 10)
	}},
	{"Oc", func(w World) bool {
		return between(w.Size, 10, 15) && inSet(w.Atmosphere, 3, 4, 5, 6, 7, 8, 9, 13, 14, 15) && w.Hydrographics == 10
	}},
	{"Va", func(w World) bool {
		return w.Atmosphere == 0
	}},
	{"Wa", func(w World) bool {
		return between(w.Size, 3, 9) && inSet(w.Atmosphere, 3, 4, 5, 6, 7, 8, 9, 13, 14, 15) && w.Hydrographics == 10
	}},

	// Population
	{"Di", func(w World) bool {
		return w.Population == 0 && w.Government == 0 && w.LawLevel == 0 && w.TechLevel >= 1
	}},
	{"Ba", func(w World) bool {
		return w.Population == 0 && w.Government == 0 && w.LawLevel == 0 && w.TechLevel == 0
	}},
	{"Lo", func(w World) bool {
		return between(w.Population, 1, 3)
	}},
	{"Ni", func(w World) bool {
		return between(w.Population, 4, 6)
	}},
	{"Ph", func(w World) bool {
		return w.Population == 8
	}},
	{"Hi", func(w World) bool {
		return between(w.Population, 9, 15)
	}},

	// Economic
	{"Pa", func(w World) bool {
		return between(w.Atmosphere, 4, 9) && between(w.Hydrographics, 4, 8) && inSet(w.Population, 4, 8)
	}},
	{"Ag", func(w World) bool {
		return between(w.Atmosphere, 4, 9) && between(w.Hydrographics, 4, 8) && between(w.Population, 5, 7)
	}},
	{"Na", func(w World) bool {
		return between(w.Atmosphere, 0, 3) && between(w.Hydrographics, 0, 3) && between(w.Population, 6, 15)
	}},
	{"Px", func(w World) bool {
		return inSet(w.Atmosphere, 2, 3, 10, 11) && between(w.Hydrographics, 1, 5) &&
			between(w.Population, 3, 6) && between(w.LawLevel, 6, 9)
	}},
	{"Pi", func(w World) bool {
		return inSet(w.Atmosphere, 0, 1, 2, 4, 7, 9) && inSet(w.Population, 7, 8)
	}},
	{"In", func(w World) bool {
		return inSet(w.Atmosphere, 0, 1, 2, 4, 7, 9, 10, 11, 12) && between(w.Population, 9, 15)
	}},
	{"Po", func(w World) bool {
		return between(w.Atmosphere, 2, 5) && between(w.Hydrographics, 0, 3)
	}},
	{"Pr", func(w World) bool {
		return inSet(w.Atmosphere, 6, 8) && inSet(w.Population, 5, 9)
	}},
	{"Ri", func(w World) bool {
		return inSet(w.Atmosphere, 6, 8) && between(w.Population, 6, 8)
	}},
}

// between reports whether lo <= v <= hi.
func between(v, lo, hi int) bool { return v >= lo && v <= hi }

// inSet reports whether v equals any of the values.
func inSet(v int, values ...int) bool {
	for _, x := range values {
		if v == x {
			return true
		}
	}
	return false
}
