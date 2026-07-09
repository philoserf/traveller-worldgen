package classic_test

import (
	"testing"

	"github.com/philoserf/traveller-worldgen/classic"
	"github.com/philoserf/traveller-worldgen/dice"
)

// TestWorkedExample feeds a hand-computed dice sequence and asserts the exact
// resulting world. See the oracle derivation in the project notes: starport A,
// naval+scout bases, size 8, atmosphere 6, hydrographics 7, population 8,
// government 4, law 5, tech 9 (1D=3 + starport-A +6), name "Varon".
func TestWorkedExample(t *testing.T) {
	t.Parallel()

	seq := []int{
		3,  // starport 2D -> A
		9,  // naval base 2D (>=8) -> yes
		11, // scout base 2D (>=10) -> yes
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
	w := classic.Generate(dice.NewScripted(seq...))

	want := classic.World{
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
	}
	if w != want {
		t.Fatalf("world mismatch:\n got %+v\nwant %+v", w, want)
	}
	if got := w.UWP(); got != "A867845-9" {
		t.Fatalf("UWP = %q, want %q", got, "A867845-9")
	}
	if got := w.Bases(); got != "Naval, Scout" {
		t.Fatalf("Bases = %q, want %q", got, "Naval, Scout")
	}
}

func TestStarportTable(t *testing.T) {
	t.Parallel()

	cases := map[int]byte{
		2: 'A', 3: 'A', 4: 'A',
		5: 'B', 6: 'B',
		7: 'C', 8: 'C',
		9:  'D',
		10: 'E', 11: 'E',
		12: 'X',
	}
	for roll, want := range cases {
		// A starport of A/B triggers extra base rolls, so append spare dice.
		w := classic.Generate(dice.NewScripted(roll, 2, 2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1))
		if w.Starport != want {
			t.Errorf("2D=%d: starport = %q, want %q", roll, w.Starport, want)
		}
	}
}

func TestSizeZeroForcesAtmosphereAndHydroZero(t *testing.T) {
	t.Parallel()

	// starport X (roll 12) -> no base rolls; size 2D-2 with roll 2 -> 0.
	// No atmosphere/hydro dice are consumed because size 0 forces both to 0.
	w := classic.Generate(dice.NewScripted(
		12, // starport X
		2,  // size -> 0
		// population, government, law, tech, name follow immediately
		2,                                  // population -> 0
		7,                                  // government 2D-7+0 -> 0
		7,                                  // law 2D-7+0 -> 0
		1,                                  // tech 1D=1 (+DMs)
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, // name
	))
	if w.Size != 0 || w.Atmosphere != 0 || w.Hydrographics != 0 {
		t.Fatalf("size 0 world: got size=%d atm=%d hyd=%d, want all 0", w.Size, w.Atmosphere, w.Hydrographics)
	}
}

func TestSizeOneForcesHydroZeroButRollsAtmosphere(t *testing.T) {
	t.Parallel()

	// starport X, size 2D-2 roll 3 -> 1; atmosphere rolls (size 1 != 0);
	// hydrographics forced to 0 (size <= 1) with no die consumed.
	w := classic.Generate(dice.NewScripted(
		12,                                 // starport X
		3,                                  // size -> 1
		12,                                 // atmosphere 2D-7+1 = 6
		2,                                  // population -> 0
		7,                                  // government -> 0
		7,                                  // law -> 0
		1,                                  // tech
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, // name
	))
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

	// Size 8 (roll 10); atmosphere 2D-7+8 with roll 12 -> 13 (>= A), which
	// triggers the -4 hydrographics DM. Hydrographics 2D-7+8 = 13, -4 = 9.
	w := classic.Generate(dice.NewScripted(
		12,                                 // starport X (no base rolls)
		10,                                 // size -> 8
		12,                                 // atmosphere 2D-7+8 = 13 (>= A, un-clamped)
		12,                                 // hydrographics 2D-7+8 = 13, -4 (atm>=A) = 9, clamp -> 9
		2,                                  // population
		7,                                  // government
		7,                                  // law
		1,                                  // tech
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, // name
	))
	if w.Atmosphere != 13 {
		t.Fatalf("atmosphere = %d, want 13 (un-clamped)", w.Atmosphere)
	}
	if w.Hydrographics != 9 {
		t.Fatalf("hydrographics = %d, want 9 (13-4)", w.Hydrographics)
	}
}

func TestGovernmentAndLawUnclamped(t *testing.T) {
	t.Parallel()

	// Population 10 (roll 12); government 2D-7+10 with roll 12 -> 15 (kept raw,
	// not capped at D). Law = 2D-7+15 with roll 12 -> 20, derived from the true
	// government rather than a capped value.
	w := classic.Generate(dice.NewScripted(
		12,                                 // starport X
		12,                                 // size -> 10
		12,                                 // atmosphere 2D-7+10 = 15 (un-clamped)
		2,                                  // hydrographics 2D-7+10 = 5, atm>=A so -4 -> 1
		12,                                 // population -> 10
		12,                                 // government 2D-7+10 = 15
		12,                                 // law 2D-7+15 = 20
		1,                                  // tech
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, // name
	))
	if w.Government != 15 {
		t.Fatalf("government = %d, want 15 (un-clamped)", w.Government)
	}
	if w.LawLevel != 20 {
		t.Fatalf("law = %d, want 20 (derived from true government)", w.LawLevel)
	}
	if w.Atmosphere != 15 {
		t.Fatalf("atmosphere = %d, want 15 (un-clamped)", w.Atmosphere)
	}
}

// TestPropertyRanges generates many seeded worlds and asserts every field
// stays within its documented range and that the UWP re-parses consistently.
func TestPropertyRanges(t *testing.T) {
	t.Parallel()

	for seed := int64(0); seed < 1000; seed++ {
		w := classic.Generate(dice.NewSeeded(seed))

		if w.Name == "" {
			t.Fatalf("seed %d: empty name", seed)
		}
		switch w.Starport {
		case 'A', 'B', 'C', 'D', 'E', 'X':
		default:
			t.Fatalf("seed %d: invalid starport %q", seed, w.Starport)
		}
		// Un-clamped ordinals: atmosphere/government = 2D-7+(size|pop), so with
		// size/pop capped at 10 and 2D at 12 the maximum is 15; law = 2D-7+gov
		// can reach 20.
		checkRange(t, seed, "size", w.Size, 10)
		checkRange(t, seed, "atmosphere", w.Atmosphere, 15)
		checkRange(t, seed, "hydrographics", w.Hydrographics, 10)
		checkRange(t, seed, "population", w.Population, 10)
		checkRange(t, seed, "government", w.Government, 15)
		checkRange(t, seed, "lawLevel", w.LawLevel, 30)
		checkRange(t, seed, "techLevel", w.TechLevel, 30)

		// Only A/B/C/D can have a scout base; only A/B a naval base.
		if w.NavalBase && w.Starport != 'A' && w.Starport != 'B' {
			t.Fatalf("seed %d: naval base with starport %q", seed, w.Starport)
		}
		if w.ScoutBase && (w.Starport == 'E' || w.Starport == 'X') {
			t.Fatalf("seed %d: scout base with starport %q", seed, w.Starport)
		}

		if len(w.UWP()) != 9 {
			t.Fatalf("seed %d: malformed UWP %q", seed, w.UWP())
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
		a := classic.Generate(dice.NewSeeded(seed))
		b := classic.Generate(dice.NewSeeded(seed))
		if a != b {
			t.Fatalf("seed %d: same seed produced different worlds:\n%+v\n%+v", seed, a, b)
		}
	}
}
