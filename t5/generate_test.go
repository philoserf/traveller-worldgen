package t5_test

import (
	"testing"

	"github.com/philoserf/traveller-worldgen/dice"
	"github.com/philoserf/traveller-worldgen/t5"
)

// TestWorkedExample feeds a hand-computed dice sequence and asserts the exact
// resulting world, including the Ix/Ex/Cx extensions. Flux draws consume two dice
// each (minuend then subtrahend), so the scripted sequence lists them
// individually in draw order.
func TestWorkedExample(t *testing.T) {
	t.Parallel()

	seq := []int{
		3,    // starport 2D -> index 1 -> A
		5,    // naval 2D<=6 -> yes
		4,    // scout 2D<=4 -> yes
		10,   // size 2D-2 -> 8
		2, 4, // atmosphere flux 2-4=-2; +size 8 -> 6
		5, 4, // hydrographics flux 5-4=+1; +atm 6 -> 7
		10,   // population 2D-2 -> 8
		1, 5, // government flux 1-5=-4; +pop 8 -> 4
		6, 4, // law flux 6-4=+2; +gov 4 -> 6
		3,    // tech 1D=3; +6 (starport A) -> 9
		2, 1, // population digit: idx=(1)*6+0=6; 6%9+1 -> 7
		5,    // planetoid belts 1D-3 -> 2
		10,   // gas giants 2D/2-2 -> 3
		8,    // Ex resources 2D=8; +GG 3 +belts 2 (TL>=8) -> 13
		6,    // Ex infrastructure 2D=6 (pop>=7); +Ix 3 -> 9
		6, 4, // Ex efficiency flux 6-4 -> +2
		5, 4, // Cx heterogeneity flux 5-4=+1; +pop 8 -> 9
		3, 2, // Cx strangeness flux 3-2=+1; +5 -> 6
		2, 1, // Cx symbols flux 2-1=+1; +TL 9 -> 10
		// name: "va"+"ro"+coda "n" -> Varon
		2, 3, 4, 1, 1, 3, 1, 1, 4, 4, 1, 1,
	}
	w := t5.Generate(dice.NewScripted(seq...))

	want := t5.World{
		Name:            "Varon",
		Starport:        'A',
		Size:            8,
		Atmosphere:      6,
		Hydrographics:   7,
		Population:      8,
		Government:      4,
		LawLevel:        6,
		TechLevel:       9,
		NavalBase:       true,
		ScoutBase:       true,
		PopulationDigit: 7,
		PlanetoidBelts:  2,
		GasGiants:       3,
		Importance:      3,
		Economic:        t5.Economic{Resources: 13, Labor: 7, Infrastructure: 9, Efficiency: 2},
		Cultural:        t5.Cultural{Heterogeneity: 9, Acceptance: 11, Strangeness: 6, Symbols: 10},
	}
	if w != want {
		t.Fatalf("world mismatch:\n got %+v\nwant %+v", w, want)
	}
	if got := w.UWP(); got != "A867846-9" {
		t.Fatalf("UWP = %q, want %q", got, "A867846-9")
	}
	if got := w.BaseCode(); got != "A" {
		t.Fatalf("BaseCode = %q, want %q", got, "A")
	}
	if got := w.PBG(); got != "723" {
		t.Fatalf("PBG = %q, want %q", got, "723")
	}
	if got := w.ImportanceString(); got != "{+3}" {
		t.Fatalf("Ix = %q, want %q", got, "{+3}")
	}
	if got := w.Economic.String(); got != "(D79+2)" {
		t.Fatalf("Ex = %q, want %q", got, "(D79+2)")
	}
	if got := w.Economic.RU(); got != 1638 {
		t.Fatalf("RU = %d, want 1638", got)
	}
	if got := w.Cultural.String(); got != "[9B6A]" {
		t.Fatalf("Cx = %q, want %q", got, "[9B6A]")
	}
	want4 := []string{"Ga", "Ph", "Pa", "Ri"}
	got := w.TradeCodes()
	if len(got) != len(want4) {
		t.Fatalf("TradeCodes = %v, want %v", got, want4)
	}
	for i := range want4 {
		if got[i] != want4[i] {
			t.Fatalf("TradeCodes = %v, want %v", got, want4)
		}
	}
}

func TestStarportTable(t *testing.T) {
	t.Parallel()

	// The single 2D->starport column (identical to Classic Book 3).
	col := [11]byte{'A', 'A', 'A', 'B', 'B', 'C', 'C', 'D', 'E', 'E', 'X'}
	for roll := 2; roll <= 12; roll++ {
		// A/B/C/D starports trigger base rolls; spare dice cover them plus the rest.
		spare := make([]int, 40)
		for i := range spare {
			spare[i] = 1
		}
		w := t5.Generate(dice.NewScripted(append([]int{roll}, spare...)...))
		if want := col[roll-2]; w.Starport != want {
			t.Errorf("2D=%d: starport = %q, want %q", roll, w.Starport, want)
		}
	}
}

func TestSizeZeroForcesAtmosphereAndHydroZero(t *testing.T) {
	t.Parallel()

	// Starport X (roll 12) -> no base rolls; size 2D=2 -> 0. No atmosphere or
	// hydrographics flux is drawn because size 0 forces both to 0; if one were
	// wrongly drawn the sequence would shift and the assertions below would fail.
	spare := filler()
	w := t5.Generate(dice.NewScripted(append([]int{
		12, // starport X
		2,  // size -> 0
		2,  // population 2D-2 -> 0
	}, spare...)...))
	if w.Size != 0 || w.Atmosphere != 0 || w.Hydrographics != 0 {
		t.Fatalf("size 0 world: got size=%d atm=%d hyd=%d, want all 0", w.Size, w.Atmosphere, w.Hydrographics)
	}
}

func TestSizeOneForcesHydroZeroButRollsAtmosphere(t *testing.T) {
	t.Parallel()

	spare := filler()
	w := t5.Generate(dice.NewScripted(append([]int{
		12,   // starport X
		3,    // size 2D-2 -> 1
		6, 1, // atmosphere flux 6-1=+5; +size 1 -> 6
		2, // population -> 0
	}, spare...)...))
	if w.Size != 1 {
		t.Fatalf("size = %d, want 1", w.Size)
	}
	if w.Atmosphere != 6 {
		t.Fatalf("atmosphere = %d, want 6 (should have rolled flux)", w.Atmosphere)
	}
	if w.Hydrographics != 0 {
		t.Fatalf("hydrographics = %d, want 0 (forced by size < 2)", w.Hydrographics)
	}
}

func TestHydrographicsMinusFourDM(t *testing.T) {
	t.Parallel()

	// Atmosphere A (10) triggers the -4 hydrographics DM. Size 8; atm flux +2 ->
	// 10 (A); hydro flux +5 + atm 10 - 4 = 11, capped at A (10).
	spare := filler()
	w := t5.Generate(dice.NewScripted(append([]int{
		12,   // starport X
		10,   // size -> 8
		6, 4, // atmosphere flux +2; +8 -> 10 (A)
		6, 1, // hydrographics flux +5; +10 -4 = 11 -> capped A (10)
		2, // population -> 0
	}, spare...)...))
	if w.Atmosphere != 10 {
		t.Fatalf("atmosphere = %d, want 10", w.Atmosphere)
	}
	if w.Hydrographics != 10 {
		t.Fatalf("hydrographics = %d, want 10 (11 capped at A)", w.Hydrographics)
	}
}

func TestSizeAndPopulationReroll(t *testing.T) {
	t.Parallel()

	// Size 2D-2 = 10 rerolls 1D+9; population 2D-2 = 10 rerolls 2D+3.
	spare := filler()
	w := t5.Generate(dice.NewScripted(append([]int{
		12,   // starport X
		12,   // size 2D-2 = 10 -> reroll
		4,    // size reroll 1D+9 -> 13 (D)
		1, 1, // atmosphere flux -> 0-ish (size 13 != 0 so flux drawn)
		1, 1, // hydrographics flux
		12, // population 2D-2 = 10 -> reroll
		9,  // population reroll 2D+3 -> 12 (C)
	}, spare...)...))
	if w.Size != 13 {
		t.Fatalf("size = %d, want 13 (rerolled 1D+9)", w.Size)
	}
	if w.Population != 12 {
		t.Fatalf("population = %d, want 12 (rerolled 2D+3)", w.Population)
	}
}

// TestPropertyRanges generates many seeded worlds and asserts every field stays
// within its documented range and the UWP/PBG re-parse.
func TestPropertyRanges(t *testing.T) {
	t.Parallel()

	for seed := int64(0); seed < 1000; seed++ {
		w := t5.Generate(dice.NewSeeded(seed))

		if w.Name == "" {
			t.Fatalf("seed %d: empty name", seed)
		}
		switch w.Starport {
		case 'A', 'B', 'C', 'D', 'E', 'X':
		default:
			t.Fatalf("seed %d: invalid starport %q", seed, w.Starport)
		}
		checkRange(t, seed, "size", w.Size, 15)
		checkRange(t, seed, "atmosphere", w.Atmosphere, 15)
		checkRange(t, seed, "hydrographics", w.Hydrographics, 10)
		checkRange(t, seed, "population", w.Population, 15)
		checkRange(t, seed, "government", w.Government, 15)
		checkRange(t, seed, "lawLevel", w.LawLevel, 18)
		checkRange(t, seed, "techLevel", w.TechLevel, 20)
		checkRange(t, seed, "gasGiants", w.GasGiants, 4)
		checkRange(t, seed, "planetoidBelts", w.PlanetoidBelts, 3)

		// PBG population digit is 1-9 when populated, else 0.
		if w.Population == 0 && w.PopulationDigit != 0 {
			t.Fatalf("seed %d: unpopulated world has pop digit %d", seed, w.PopulationDigit)
		}
		if w.Population > 0 && (w.PopulationDigit < 1 || w.PopulationDigit > 9) {
			t.Fatalf("seed %d: pop digit %d out of 1-9", seed, w.PopulationDigit)
		}

		// Bases only on starports that can hold them.
		if w.NavalBase && w.Starport != 'A' && w.Starport != 'B' {
			t.Fatalf("seed %d: naval base with starport %q", seed, w.Starport)
		}
		if w.ScoutBase && (w.Starport == 'E' || w.Starport == 'X') {
			t.Fatalf("seed %d: scout base with starport %q", seed, w.Starport)
		}

		// Extensions: R/L/I and the Cx components are non-negative; an unpopulated
		// world has an all-zero Cultural extension.
		if w.Economic.Resources < 0 || w.Economic.Labor < 0 || w.Economic.Infrastructure < 0 {
			t.Fatalf("seed %d: negative economic component %+v", seed, w.Economic)
		}
		if w.Population == 0 && w.Cultural != (t5.Cultural{}) {
			t.Fatalf("seed %d: unpopulated world has non-zero Cx %+v", seed, w.Cultural)
		}
		if w.Cultural.Heterogeneity < 0 || w.Cultural.Acceptance < 0 ||
			w.Cultural.Strangeness < 0 || w.Cultural.Symbols < 0 {
			t.Fatalf("seed %d: negative cultural component %+v", seed, w.Cultural)
		}

		if len(w.UWP()) != 9 {
			t.Fatalf("seed %d: malformed UWP %q", seed, w.UWP())
		}
		if len(w.PBG()) != 3 {
			t.Fatalf("seed %d: malformed PBG %q", seed, w.PBG())
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
		a := t5.Generate(dice.NewSeeded(seed))
		b := t5.Generate(dice.NewSeeded(seed))
		if a != b {
			t.Fatalf("seed %d: same seed produced different worlds:\n%+v\n%+v", seed, a, b)
		}
	}
}

// filler returns 40 dice all showing 1, used to feed the draws after the
// characteristics under test (bases, extensions, name) so the Scripted roller
// never exhausts. Forty comfortably covers the longest remaining draw chain.
func filler() []int {
	s := make([]int, 40)
	for i := range s {
		s[i] = 1
	}
	return s
}
