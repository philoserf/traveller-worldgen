package tne_test

import (
	"testing"

	"github.com/philoserf/traveller-worldgen/tne"
)

func TestDescriptions(t *testing.T) {
	t.Parallel()

	w := tne.World{
		Starport: 'A', Size: 10, Atmosphere: 15, Hydrographics: 10,
		Population: 9, Government: 15, LawLevel: 20, TechLevel: 13,
	}
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"starport", w.StarportDesc(), "Excellent; starship shipyard; overhaul; refined fuel"},
		{"size", w.SizeDesc(), "Large (16,000 km)"},
		{"atmosphere", w.AtmosphereDesc(), "Exotic (ellipsoid)"},
		{"hydrographics", w.HydrographicsDesc(), "Water World (100%)"},
		{"population", w.PopulationDesc(), "High (billions)"},
		{"government", w.GovernmentDesc(), "Totalitarian Oligarchy"},
		{"tech", w.TechLevelDesc(), "Average Stellar (holographic data storage; the Imperium)"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s = %q, want %q", c.name, c.got, c.want)
		}
	}

	// Smallest populations use TNE's "Inconsequential"; the Ancients anchor TL L.
	if got := (tne.World{Population: 0}).PopulationDesc(); got != "Inconsequential (< ten)" {
		t.Errorf("population 0 = %q", got)
	}
	if got := (tne.World{TechLevel: 20}).TechLevelDesc(); got != "Extreme Stellar (the Ancients)" {
		t.Errorf("tech 20 = %q", got)
	}
}

func TestLawLevelDesc(t *testing.T) {
	t.Parallel()

	if got := (tne.World{LawLevel: 0}).LawLevelDesc(); got != "No law (no prohibitions)" {
		t.Errorf("law 0 = %q", got)
	}
	want := "Moderate law (all firearms except shotguns prohibited); 6+ to avoid arrest"
	if got := (tne.World{LawLevel: 6}).LawLevelDesc(); got != want {
		t.Errorf("law 6 = %q, want %q", got, want)
	}
	if got := (tne.World{LawLevel: 20}).LawLevelDesc(); got != "Extreme law (totally oppressive and restrictive); 20+ to avoid arrest" {
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
		{true, false, true, "NM", "Naval, Military"},
	}
	for _, c := range cases {
		w := tne.World{NavalBase: c.naval, ScoutBase: c.scout, MilitaryBase: c.military}
		if got := w.BaseCode(); got != c.code {
			t.Errorf("BaseCode(%v) = %q, want %q", c, got, c.code)
		}
		if got := w.Bases(); got != c.human {
			t.Errorf("Bases(%v) = %q, want %q", c, got, c.human)
		}
	}
}
