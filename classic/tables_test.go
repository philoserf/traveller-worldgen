package classic_test

import (
	"testing"

	"github.com/philoserf/traveller-worldgen/classic"
)

func TestDescriptionsKnownValues(t *testing.T) {
	t.Parallel()

	w := classic.World{
		Starport:      'A',
		Size:          8,
		Atmosphere:    6,
		Hydrographics: 7,
		Population:    8,
		Government:    4,
		LawLevel:      5,
		TechLevel:     9,
	}
	cases := []struct {
		name, got, want string
	}{
		{"size", w.SizeDesc(), "8,000 miles diameter"},
		{"atmosphere", w.AtmosphereDesc(), "Standard"},
		{"hydrographics", w.HydrographicsDesc(), "70%"},
		{"population", w.PopulationDesc(), "100,000,000"},
		{"government", w.GovernmentDesc(), "Representative Democracy"},
		{"lawLevel", w.LawLevelDesc(), "Personal concealable firearms prohibited; 5+ to avoid arrest"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s desc = %q, want %q", c.name, c.got, c.want)
		}
	}
	if got := w.TechLevelDesc(); got == "" || got == "(beyond described range)" {
		t.Errorf("tech desc for TL9 unexpectedly empty/beyond: %q", got)
	}
}

func TestDescriptionsBeyondRange(t *testing.T) {
	t.Parallel()

	// Atmosphere, government, and tech have no Book 3 guidance past their tables,
	// so they fall back to the generic marker.
	w := classic.World{Atmosphere: 15, Government: 15, TechLevel: 18}
	if got := w.AtmosphereDesc(); got != "(beyond described range)" {
		t.Errorf("atmosphere 15 desc = %q, want beyond-range", got)
	}
	if got := w.GovernmentDesc(); got != "(beyond described range)" {
		t.Errorf("government 15 desc = %q, want beyond-range", got)
	}
	if got := w.TechLevelDesc(); got != "(beyond described range)" {
		t.Errorf("tech 18 desc = %q, want beyond-range", got)
	}
}

func TestLawLevelDescEnforcementThrow(t *testing.T) {
	t.Parallel()

	const level9 = "Possession of any weapon outside one's home prohibited"

	cases := []struct {
		level int
		want  string
	}{
		{0, "No weapons laws"}, // no laws → no enforcement throw
		{5, "Personal concealable firearms prohibited; 5+ to avoid arrest"},
		{9, level9 + "; 9+ to avoid arrest"},
		{15, level9 + "; 15+ to avoid arrest"}, // past 9: keep level-9 prohibition, throw rises
	}
	for _, c := range cases {
		if got := (classic.World{LawLevel: c.level}).LawLevelDesc(); got != c.want {
			t.Errorf("law %d desc = %q, want %q", c.level, got, c.want)
		}
	}
}
