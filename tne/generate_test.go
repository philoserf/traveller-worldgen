package tne_test

import (
	"testing"

	"github.com/philoserf/traveller-worldgen/dice"
	"github.com/philoserf/traveller-worldgen/tne"
)

// TestWorkedExample feeds a hand-computed dice sequence (Standard subsector) and
// asserts the exact resulting world: starport A, naval+scout+military bases,
// size 8, atmosphere 6, hydrographics 7, population 8, government 4, law 5,
// tech 9 (1D=3 + starport-A +6), name "Varon".
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
		// name: 2 syllables "va"+"ro"+coda "n" -> Varon
		2, 3, 4, 1, 1, 3, 1, 1, 4, 5, 1, 1,
	}
	w := tne.Generate(dice.NewScripted(seq...), tne.Standard)

	want := tne.World{
		Name:          "Varon",
		Starport:      'A',
		Size:          8,
		Atmosphere:    6,
		Hydrographics: 7,
		Population:    8,
		Government:    4,
		LawLevel:      5,
		TechLevel:     9,
		NavalBase:     true,
		ScoutBase:     true,
		MilitaryBase:  true,
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
}

func TestStarportTablePerNature(t *testing.T) {
	t.Parallel()

	// Each nature's 2D->starport column, from § Starport (2D).
	cols := map[tne.Nature][11]byte{
		tne.Backwater: {'A', 'A', 'B', 'B', 'C', 'C', 'C', 'D', 'E', 'E', 'X'},
		tne.Standard:  {'A', 'A', 'A', 'B', 'B', 'C', 'C', 'D', 'E', 'E', 'X'},
		tne.Mature:    {'A', 'A', 'A', 'B', 'B', 'C', 'C', 'D', 'E', 'E', 'E'},
		tne.Cluster:   {'A', 'A', 'A', 'A', 'B', 'B', 'C', 'C', 'D', 'E', 'X'},
	}
	for nature, col := range cols {
		for roll := 2; roll <= 12; roll++ {
			// A/B/C starports trigger extra base rolls, so append spare dice.
			w := tne.Generate(
				dice.NewScripted(roll, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1),
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
	w := tne.Generate(dice.NewScripted(
		12,                                 // starport X
		2,                                  // size -> 0
		2,                                  // population -> 0
		7,                                  // government 2D-7+0 -> 0
		7,                                  // law 2D-7+0 -> 0
		1,                                  // tech 1D=1
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, // name
	), tne.Standard)
	if w.Size != 0 || w.Atmosphere != 0 || w.Hydrographics != 0 {
		t.Fatalf("size 0 world: got size=%d atm=%d hyd=%d, want all 0", w.Size, w.Atmosphere, w.Hydrographics)
	}
}

func TestSizeOneForcesHydroZeroButRollsAtmosphere(t *testing.T) {
	t.Parallel()

	w := tne.Generate(dice.NewScripted(
		12,                                 // starport X
		3,                                  // size -> 1
		12,                                 // atmosphere 2D-7+1 = 6
		2,                                  // population -> 0
		7,                                  // government -> 0
		7,                                  // law -> 0
		1,                                  // tech
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, // name
	), tne.Standard)
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
	w := tne.Generate(dice.NewScripted(
		12,                                 // starport X (no base rolls)
		10,                                 // size -> 8
		12,                                 // atmosphere 2D-7+8 = 13 (>= A)
		12,                                 // hydrographics 2D-7+8 = 13, -4 = 9
		2,                                  // population
		7,                                  // government
		7,                                  // law
		1,                                  // tech
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, // name
	), tne.Standard)
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
	// roll 12 -> 20 (L).
	w := tne.Generate(dice.NewScripted(
		12,                                 // starport X
		12,                                 // size -> 10
		2,                                  // atmosphere 2D-7+10 = 5
		2,                                  // hydrographics 2D-7+10 = 5
		12,                                 // population -> 10
		12,                                 // government 2D-7+10 = 15 (F)
		12,                                 // law 2D-7+15 = 20 (L)
		1,                                  // tech
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, // name
	), tne.Standard)
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

// TestGovernmentETechDM exercises the TNE-specific Government E (-1) tech DM,
// which MegaTraveller lacks. With starport X (-4) + population A (+4) + government
// E (-1) the DMs sum to -1, so a 1D roll of 6 yields tech 5 (not 6).
func TestGovernmentETechDM(t *testing.T) {
	t.Parallel()

	w := tne.Generate(dice.NewScripted(
		12,                                 // starport X (DM -4, no base rolls)
		12,                                 // size -> 10 (DM 0)
		3,                                  // atmosphere 2D-7+10 = 6 (DM 0)
		2,                                  // hydrographics 2D-7+10 = 5 (DM 0)
		12,                                 // population -> 10 (DM +4)
		11,                                 // government 2D-7+10 = 14 (E) (DM -1)
		2,                                  // law 2D-7+14 = 9
		6,                                  // tech 1D=6 + (-4 +4 -1) = 5
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, // name
	), tne.Standard)
	if w.Government != 14 {
		t.Fatalf("government = %d, want 14 (E)", w.Government)
	}
	if w.TechLevel != 5 {
		t.Fatalf("tech = %d, want 5 (1D6=6 with net DM -1; gov-E -1 applied)", w.TechLevel)
	}
}

// TestPropertyRanges generates many seeded worlds across all natures and asserts
// every field stays within its documented range and the UWP re-parses.
func TestPropertyRanges(t *testing.T) {
	t.Parallel()

	natures := []tne.Nature{tne.Backwater, tne.Standard, tne.Mature, tne.Cluster}
	for _, nature := range natures {
		for seed := int64(0); seed < 500; seed++ {
			w := tne.Generate(dice.NewSeeded(seed), nature)

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
		a := tne.Generate(dice.NewSeeded(seed), tne.Standard)
		b := tne.Generate(dice.NewSeeded(seed), tne.Standard)
		if a != b {
			t.Fatalf("seed %d: same seed produced different worlds:\n%+v\n%+v", seed, a, b)
		}
	}
}

func TestParseNature(t *testing.T) {
	t.Parallel()

	cases := map[string]tne.Nature{
		"backwater": tne.Backwater,
		"standard":  tne.Standard,
		"mature":    tne.Mature,
		"cluster":   tne.Cluster,
	}
	for s, want := range cases {
		n, ok := tne.ParseNature(s)
		if !ok || n != want {
			t.Errorf("ParseNature(%q) = %v, %v; want %v, true", s, n, ok, want)
		}
		if n.String() != s {
			t.Errorf("Nature(%q).String() = %q, want %q", s, n.String(), s)
		}
	}
	if _, ok := tne.ParseNature("bogus"); ok {
		t.Error("ParseNature(bogus) = ok, want not ok")
	}
}
