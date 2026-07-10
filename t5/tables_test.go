package t5_test

import (
	"slices"
	"testing"

	"github.com/philoserf/traveller-worldgen/t5"
)

func TestDescriptions(t *testing.T) {
	t.Parallel()

	w := t5.World{
		Starport: 'A', Size: 15, Atmosphere: 15, Hydrographics: 10,
		Population: 15, Government: 15, LawLevel: 18, TechLevel: 20,
	}
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"starport", w.StarportDesc(), "Excellent; starship shipyard; overhaul; refined fuel"},
		{"size", w.SizeDesc(), "15,000 miles (24,000 km)"},
		{"atmosphere", w.AtmosphereDesc(), "Unusual"},
		{"hydrographics", w.HydrographicsDesc(), "Water World"},
		{"population", w.PopulationDesc(), "Quadrillions (10^15)"},
		{"government", w.GovernmentDesc(), "Totalitarian Oligarchy"},
		{"law", w.LawLevelDesc(), "Extreme Law. Routine oppression"},
		{"tech", w.TechLevelDesc(), "Extreme Stellar (the Ancients)"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s = %q, want %q", c.name, c.got, c.want)
		}
	}
}

// TestLawLevelNoArrestThrow pins T5's departure from the earlier editions: the
// law description is descriptive only, with no "N+ to avoid arrest" throw.
func TestLawLevelNoArrestThrow(t *testing.T) {
	t.Parallel()

	if got := (t5.World{LawLevel: 0}).LawLevelDesc(); got != "No Law. No prohibitions" {
		t.Errorf("law 0 = %q", got)
	}
	if got := (t5.World{LawLevel: 6}).LawLevelDesc(); got != "Moderate Law. Machine guns prohibited" {
		t.Errorf("law 6 = %q (should carry no arrest throw)", got)
	}
}

func TestBaseCodeAndBases(t *testing.T) {
	t.Parallel()

	cases := []struct {
		naval, scout bool
		code, human  string
	}{
		{false, false, "", "None"},
		{true, false, "N", "Naval"},
		{false, true, "S", "Scout"},
		{true, true, "A", "Naval, Scout"},
	}
	for _, c := range cases {
		w := t5.World{NavalBase: c.naval, ScoutBase: c.scout}
		if got := w.BaseCode(); got != c.code {
			t.Errorf("BaseCode(%v) = %q, want %q", c, got, c.code)
		}
		if got := w.Bases(); got != c.human {
			t.Errorf("Bases(%v) = %q, want %q", c, got, c.human)
		}
	}
}

func TestTradeCodes(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		w    t5.World
		want []string
	}{
		{
			"agricultural + rich",
			t5.World{Atmosphere: 6, Hydrographics: 6, Population: 7},
			[]string{"Ag", "Ri"},
		},
		{
			"asteroid emits both As and Va, plus Barren at TL 0",
			t5.World{Size: 0, Atmosphere: 0, Hydrographics: 0, Population: 0, Government: 0, LawLevel: 0, TechLevel: 0},
			[]string{"As", "Va", "Ba"}, // T5 states no As->Va suppression
		},
		{
			"dieback: pop/gov/law 0 but TL present",
			t5.World{Size: 4, Atmosphere: 0, Hydrographics: 0, Population: 0, Government: 0, LawLevel: 0, TechLevel: 5},
			[]string{"Va", "Di"}, // Di (TL >= 1), not Ba (TL 0)
		},
		{
			"water world",
			t5.World{Size: 5, Atmosphere: 6, Hydrographics: 10, Population: 5, Government: 4},
			[]string{"Wa", "Ni", "Pr"}, // atm 6 + pop 5 also qualifies as Pre-Rich

		},
	}
	for _, c := range cases {
		got := c.w.TradeCodes()
		if !slices.Equal(got, c.want) {
			t.Errorf("%s: TradeCodes = %v, want %v", c.name, got, c.want)
		}
	}
}
