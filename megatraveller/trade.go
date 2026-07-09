package megatraveller

import "slices"

// Trade classifications, derived from the UWP (see § Trade classifications). Each
// entry pairs a code with a predicate over the world's characteristics; a world
// carries every code whose predicate matches. Order follows the source table.
type tradeClass struct {
	code  string
	match func(World) bool
}

var tradeClassifications = []tradeClass{
	{"Ag", func(w World) bool {
		return between(w.Atmosphere, 4, 9) && between(w.Hydrographics, 4, 8) && between(w.Population, 5, 7)
	}},
	{"As", func(w World) bool {
		return w.Size == 0 && w.Atmosphere == 0 && w.Hydrographics == 0
	}},
	{"Ba", func(w World) bool {
		return w.Population == 0 && w.Government == 0 && w.LawLevel == 0
	}},
	{"De", func(w World) bool {
		return w.Atmosphere >= 2 && w.Hydrographics == 0
	}},
	{"Fl", func(w World) bool {
		return w.Atmosphere >= 10 && w.Hydrographics >= 1
	}},
	{"Hi", func(w World) bool {
		return w.Population >= 9
	}},
	{"Ic", func(w World) bool {
		return between(w.Atmosphere, 0, 1) && w.Hydrographics >= 1
	}},
	{"In", func(w World) bool {
		return inSet(w.Atmosphere, 2, 3, 4, 7, 9) && w.Population >= 9
	}},
	{"Lo", func(w World) bool {
		return between(w.Population, 0, 3)
	}},
	{"Na", func(w World) bool {
		return between(w.Atmosphere, 0, 3) && between(w.Hydrographics, 0, 3) && w.Population >= 6
	}},
	{"Ni", func(w World) bool {
		return between(w.Population, 0, 6)
	}},
	{"Po", func(w World) bool {
		return between(w.Atmosphere, 2, 5) && between(w.Hydrographics, 0, 3)
	}},
	{"Ri", func(w World) bool {
		return inSet(w.Atmosphere, 6, 8) && between(w.Population, 6, 8) && between(w.Government, 4, 9)
	}},
	{"Va", func(w World) bool {
		return w.Atmosphere == 0
	}},
	{"Wa", func(w World) bool {
		return w.Hydrographics == 10
	}},
}

// between reports whether lo <= v <= hi.
func between(v, lo, hi int) bool { return v >= lo && v <= hi }

// inSet reports whether v equals any of the values.
func inSet(v int, values ...int) bool { return slices.Contains(values, v) }
