package megatraveller_test

import (
	"testing"

	"github.com/philoserf/traveller-worldgen/dice"
	"github.com/philoserf/traveller-worldgen/megatraveller"
)

// TestWorkedExample feeds a hand-computed dice sequence (Standard subsector) and
// asserts the exact resulting world: starport A, naval+scout+military bases,
// size 8, atmosphere 6, hydrographics 7, population 8, government 4, law 5,
// tech 9 (1D=3 + starport-A +6), gas giants 3, planetoid belts 2, name "Varon".
func TestWorkedExample(t *testing.T) {
	t.Parallel()

	seq := []int{
		3,  // starport 2D (Standard) -> A
		9,  // naval base 2D (>=8) -> yes
		11, // scout base 2D (>=10) -> yes
		10, // military base 2D (>=10) -> yes
		10, // size 2D-2 -> 8
		5,  // atmosphere 2D-7+8 -> 6
		6,  // hydrographics 2D-7+8 -> 7
		10, // population 2D-2 -> 8
		3,  // government 2D-7+8 -> 4
		8,  // law 2D-7+4 -> 5
		3,  // tech 1D=3 (+6 starport A) -> 9
		8,  // gas giants present (2D>=5)
		7,  // gas giant quantity 2D=7 -> 3
		10, // planetoid belts present (2D>=8)
		8,  // planetoid belt quantity 2D=8 -> 2
		// name: 2 syllables "va"+"ro"+coda "n" -> Varon
		2, 3, 4, 1, 1, 3, 1, 1, 4, 5, 1, 1,
	}
	w := megatraveller.Generate(dice.NewScripted(seq...), megatraveller.Standard)

	want := megatraveller.World{
		Name:           "Varon",
		Starport:       'A',
		Size:           8,
		Atmosphere:     6,
		Hydrographics:  7,
		Population:     8,
		Government:     4,
		LawLevel:       5,
		TechLevel:      9,
		NavalBase:      true,
		ScoutBase:      true,
		MilitaryBase:   true,
		GasGiants:      3,
		PlanetoidBelts: 2,
	}
	if w != want {
		t.Fatalf("world mismatch:\n got %+v\nwant %+v", w, want)
	}
	if got := w.UWP(); got != "A867845-9" {
		t.Fatalf("UWP = %q, want %q", got, "A867845-9")
	}
	if got := w.BaseCode(); got != "AM" {
		t.Fatalf("BaseCode = %q, want %q", got, "AM")
	}
	if got := w.Bases(); got != "Naval, Scout, Military" {
		t.Fatalf("Bases = %q, want %q", got, "Naval, Scout, Military")
	}
	if got := w.TradeCodes(); len(got) != 1 || got[0] != "Ri" {
		t.Fatalf("TradeCodes = %v, want [Ri]", got)
	}
}

func TestStarportTablePerNature(t *testing.T) {
	t.Parallel()

	// Each nature's 2D->starport column, from § Starport (2D).
	cols := map[megatraveller.Nature][11]byte{
		megatraveller.Backwater: {'A', 'A', 'B', 'B', 'C', 'C', 'C', 'D', 'E', 'E', 'X'},
		megatraveller.Standard:  {'A', 'A', 'A', 'B', 'B', 'C', 'C', 'D', 'E', 'E', 'X'},
		megatraveller.Mature:    {'A', 'A', 'A', 'B', 'B', 'C', 'C', 'D', 'E', 'E', 'E'},
		megatraveller.Cluster:   {'A', 'A', 'A', 'A', 'B', 'B', 'C', 'C', 'D', 'E', 'X'},
	}
	for nature, col := range cols {
		for roll := 2; roll <= 12; roll++ {
			// A/B/C starports trigger extra base rolls, so append spare dice.
			w := megatraveller.Generate(
				dice.NewScripted(roll, 2, 2, 2, 2, 2, 2, 2, 2, 2, 1, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1),
				nature,
			)
			if want := col[roll-2]; w.Starport != want {
				t.Errorf("%s 2D=%d: starport = %q, want %q", nature, roll, w.Starport, want)
			}
		}
	}
}

func TestSizeZeroForcesAtmosphereAndHydroZero(t *testing.T) {
	t.Parallel()

	// starport X (Standard roll 12) -> no base rolls; size 2D-2 with roll 2 -> 0.
	// No atmosphere/hydro dice are consumed because size 0 forces both to 0.
	w := megatraveller.Generate(dice.NewScripted(
		12,                                 // starport X
		2,                                  // size -> 0
		2,                                  // population -> 0
		7,                                  // government 2D-7+0 -> 0
		7,                                  // law 2D-7+0 -> 0
		1,                                  // tech 1D=1
		2,                                  // gas giants absent (2D<5)
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, // name
	), megatraveller.Standard)
	if w.Size != 0 || w.Atmosphere != 0 || w.Hydrographics != 0 {
		t.Fatalf("size 0 world: got size=%d atm=%d hyd=%d, want all 0", w.Size, w.Atmosphere, w.Hydrographics)
	}
}

func TestSizeOneForcesHydroZeroButRollsAtmosphere(t *testing.T) {
	t.Parallel()

	w := megatraveller.Generate(dice.NewScripted(
		12,                                 // starport X
		3,                                  // size -> 1
		12,                                 // atmosphere 2D-7+1 = 6
		2,                                  // population -> 0
		7,                                  // government -> 0
		7,                                  // law -> 0
		1,                                  // tech
		2,                                  // gas giants absent
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, // name
	), megatraveller.Standard)
	if w.Size != 1 {
		t.Fatalf("size = %d, want 1", w.Size)
	}
	if w.Atmosphere != 6 {
		t.Fatalf("atmosphere = %d, want 6 (should have rolled)", w.Atmosphere)
	}
	if w.Hydrographics != 0 {
		t.Fatalf("hydrographics = %d, want 0 (forced)", w.Hydrographics)
	}
}

func TestHydrographicsMinusFourDM(t *testing.T) {
	t.Parallel()

	// Size 8; atmosphere 2D-7+8 with roll 12 -> 13 (>= A), triggering the -4
	// hydrographics DM. Hydrographics 2D-7+8 = 13, -4 = 9.
	w := megatraveller.Generate(dice.NewScripted(
		12,                                 // starport X (no base rolls)
		10,                                 // size -> 8
		12,                                 // atmosphere 2D-7+8 = 13 (>= A)
		12,                                 // hydrographics 2D-7+8 = 13, -4 = 9
		2,                                  // population
		7,                                  // government
		7,                                  // law
		1,                                  // tech
		2,                                  // gas giants absent
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, // name
	), megatraveller.Standard)
	if w.Atmosphere != 13 {
		t.Fatalf("atmosphere = %d, want 13", w.Atmosphere)
	}
	if w.Hydrographics != 9 {
		t.Fatalf("hydrographics = %d, want 9 (13-4)", w.Hydrographics)
	}
}

func TestGovernmentAndLawReachExtendedRange(t *testing.T) {
	t.Parallel()

	// Population 10; government 2D-7+10 with roll 12 -> 15 (F); law 2D-7+15 with
	// roll 12 -> 20 (L) — MegaTraveller's full extended ranges.
	w := megatraveller.Generate(dice.NewScripted(
		12,                                 // starport X
		12,                                 // size -> 10
		2,                                  // atmosphere 2D-7+10 = 5
		2,                                  // hydrographics 2D-7+10 = 5
		12,                                 // population -> 10
		12,                                 // government 2D-7+10 = 15 (F)
		12,                                 // law 2D-7+15 = 20 (L)
		1,                                  // tech
		2,                                  // gas giants absent
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, // name
	), megatraveller.Standard)
	if w.Government != 15 {
		t.Fatalf("government = %d, want 15 (F)", w.Government)
	}
	if w.LawLevel != 20 {
		t.Fatalf("law = %d, want 20 (L)", w.LawLevel)
	}
	if w.UWP()[5] != 'F' || w.UWP()[6] != 'L' {
		t.Fatalf("UWP = %q, want government F and law L", w.UWP())
	}
}

func TestGasGiantsAndBeltsAbsentConsumeNoQuantityDie(t *testing.T) {
	t.Parallel()

	// A gas-giant presence roll below 5 (and a belt presence roll below 8) means
	// no quantity die is drawn for either, so the following dice feed the name
	// generator. If a quantity die were wrongly consumed, the Scripted roller
	// would exhaust and panic.
	w := megatraveller.Generate(dice.NewScripted(
		12,                                 // starport X
		2,                                  // size -> 0
		2,                                  // population
		7,                                  // government
		7,                                  // law
		1,                                  // tech
		4,                                  // gas giants presence 2D=4 (<5) -> absent
		2,                                  // planetoid belts presence 2D=2 (<8) -> absent
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, // name (no quantity dice consumed)
	), megatraveller.Standard)
	if w.GasGiants != 0 {
		t.Fatalf("gas giants = %d, want 0", w.GasGiants)
	}
	if w.PlanetoidBelts != 0 {
		t.Fatalf("planetoid belts = %d, want 0", w.PlanetoidBelts)
	}
	if w.Name == "" {
		t.Fatal("empty name (a quantity die may have stolen a name die)")
	}
}

// TestPropertyRanges generates many seeded worlds across all natures and asserts
// every field stays within its documented range and the UWP re-parses.
func TestPropertyRanges(t *testing.T) {
	t.Parallel()

	natures := []megatraveller.Nature{
		megatraveller.Backwater, megatraveller.Standard, megatraveller.Mature, megatraveller.Cluster,
	}
	for _, nature := range natures {
		for seed := int64(0); seed < 500; seed++ {
			w := megatraveller.Generate(dice.NewSeeded(seed), nature)

			if w.Name == "" {
				t.Fatalf("%s seed %d: empty name", nature, seed)
			}
			switch w.Starport {
			case 'A', 'B', 'C', 'D', 'E', 'X':
			default:
				t.Fatalf("%s seed %d: invalid starport %q", nature, seed, w.Starport)
			}
			// atmosphere/government = 2D-7+(size|pop), max 15; law = 2D-7+gov, max 20.
			checkRange(t, seed, "size", w.Size, 10)
			checkRange(t, seed, "atmosphere", w.Atmosphere, 15)
			checkRange(t, seed, "hydrographics", w.Hydrographics, 10)
			checkRange(t, seed, "population", w.Population, 10)
			checkRange(t, seed, "government", w.Government, 15)
			checkRange(t, seed, "lawLevel", w.LawLevel, 20)
			checkRange(t, seed, "techLevel", w.TechLevel, 20)
			checkRange(t, seed, "gasGiants", w.GasGiants, 5)
			// A plain 2D belt-quantity roll reaches only 2-12, so at most 2 belts
			// (the printed 13->3 row is unreachable without a DM).
			checkRange(t, seed, "planetoidBelts", w.PlanetoidBelts, 2)

			// Naval only on A/B; scout only on A-D; military only on A-C.
			if w.NavalBase && w.Starport != 'A' && w.Starport != 'B' {
				t.Fatalf("%s seed %d: naval base with starport %q", nature, seed, w.Starport)
			}
			if w.ScoutBase && (w.Starport == 'E' || w.Starport == 'X') {
				t.Fatalf("%s seed %d: scout base with starport %q", nature, seed, w.Starport)
			}
			if w.MilitaryBase && (w.Starport == 'D' || w.Starport == 'E' || w.Starport == 'X') {
				t.Fatalf("%s seed %d: military base with starport %q", nature, seed, w.Starport)
			}

			if len(w.UWP()) != 9 {
				t.Fatalf("%s seed %d: malformed UWP %q", nature, seed, w.UWP())
			}
		}
	}
}

func checkRange(t *testing.T, seed int64, name string, v, hi int) {
	t.Helper()
	if v < 0 || v > hi {
		t.Fatalf("seed %d: %s = %d, out of [0,%d]", seed, name, v, hi)
	}
}

func TestGenerateDeterministic(t *testing.T) {
	t.Parallel()

	for seed := int64(0); seed < 100; seed++ {
		a := megatraveller.Generate(dice.NewSeeded(seed), megatraveller.Standard)
		b := megatraveller.Generate(dice.NewSeeded(seed), megatraveller.Standard)
		if a != b {
			t.Fatalf("seed %d: same seed produced different worlds:\n%+v\n%+v", seed, a, b)
		}
	}
}

func TestParseNature(t *testing.T) {
	t.Parallel()

	cases := map[string]megatraveller.Nature{
		"backwater": megatraveller.Backwater,
		"standard":  megatraveller.Standard,
		"mature":    megatraveller.Mature,
		"cluster":   megatraveller.Cluster,
	}
	for s, want := range cases {
		n, ok := megatraveller.ParseNature(s)
		if !ok || n != want {
			t.Errorf("ParseNature(%q) = %v, %v; want %v, true", s, n, ok, want)
		}
		if n.String() != s {
			t.Errorf("Nature(%q).String() = %q, want %q", s, n.String(), s)
		}
	}
	if _, ok := megatraveller.ParseNature("bogus"); ok {
		t.Error("ParseNature(bogus) = ok, want not ok")
	}
}
