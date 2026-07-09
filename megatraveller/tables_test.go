package megatraveller_test

import (
	"slices"
	"testing"

	"github.com/philoserf/traveller-worldgen/megatraveller"
)

func TestDescriptions(t *testing.T) {
	t.Parallel()

	w := megatraveller.World{
		Starport: 'A', Size: 10, Atmosphere: 15, Hydrographics: 10,
		Population: 9, Government: 15, LawLevel: 20, TechLevel: 20,
	}
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"starport", w.StarportDesc(), "Excellent; starship shipyard; overhaul; refined fuel"},
		{"size", w.SizeDesc(), "Large (16,000 km)"},
		{"atmosphere", w.AtmosphereDesc(), "Thin, low"},
		{"hydrographics", w.HydrographicsDesc(), "Water World (95–100%)"},
		{"population", w.PopulationDesc(), "High (billions)"},
		{"government", w.GovernmentDesc(), "Totalitarian Oligarchy"},
		{"tech", w.TechLevelDesc(), "Extreme Stellar (the Ancients)"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s = %q, want %q", c.name, c.got, c.want)
		}
	}
}

func TestLawLevelDesc(t *testing.T) {
	t.Parallel()

	// Law 0 has no enforcement throw; any higher level appends it.
	if got := (megatraveller.World{LawLevel: 0}).LawLevelDesc(); got != "No law (no prohibitions)" {
		t.Errorf("law 0 = %q", got)
	}
	want := "Moderate law (all firearms except shotguns prohibited); 6+ to avoid arrest"
	if got := (megatraveller.World{LawLevel: 6}).LawLevelDesc(); got != want {
		t.Errorf("law 6 = %q, want %q", got, want)
	}
	if got := (megatraveller.World{LawLevel: 20}).LawLevelDesc(); got != "Extreme law (totally oppressive and restrictive); 20+ to avoid arrest" {
		t.Errorf("law 20 = %q", got)
	}
}

func TestBaseCodeAndBases(t *testing.T) {
	t.Parallel()

	cases := []struct {
		naval, scout, military bool
		code, human            string
	}{
		{false, false, false, "", "None"},
		{true, false, false, "N", "Naval"},
		{false, true, false, "S", "Scout"},
		{true, true, false, "A", "Naval, Scout"},
		{false, false, true, "M", "Military"},
		{true, true, true, "AM", "Naval, Scout, Military"},
		{false, true, true, "SM", "Scout, Military"},
	}
	for _, c := range cases {
		w := megatraveller.World{NavalBase: c.naval, ScoutBase: c.scout, MilitaryBase: c.military}
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
		w    megatraveller.World
		want []string
	}{
		{
			"agricultural + rich",
			megatraveller.World{Atmosphere: 6, Hydrographics: 6, Population: 7, Government: 5, LawLevel: 3},
			[]string{"Ag", "Ri"}, // pop 7 is outside Ni's 0-6
		},
		{
			"asteroid suppresses vacuum",
			megatraveller.World{Size: 0, Atmosphere: 0, Hydrographics: 0, Population: 4},
			[]string{"As", "Ni"},
		},
		{
			"barren",
			megatraveller.World{Size: 4, Atmosphere: 0, Hydrographics: 0, Population: 0, Government: 0, LawLevel: 0},
			[]string{"Ba", "Lo", "Ni", "Va"}, // atmo 0 excludes De; pop 0 excludes Na
		},
		{
			"industrial excludes atmosphere 5 and 6",
			megatraveller.World{Atmosphere: 5, Hydrographics: 5, Population: 9},
			[]string{"Hi"}, // atmosphere 5 is not in In's {2,3,4,7,9}
		},
		{
			"water world",
			megatraveller.World{Atmosphere: 6, Hydrographics: 10, Population: 5, Government: 4},
			[]string{"Ni", "Wa"},
		},
	}
	for _, c := range cases {
		got := c.w.TradeCodes()
		if !slices.Equal(got, c.want) {
			t.Errorf("%s: TradeCodes = %v, want %v", c.name, got, c.want)
		}
	}
}
