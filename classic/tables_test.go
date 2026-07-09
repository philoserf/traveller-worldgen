package classic_test

import (
	"strings"
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
		{"lawLevel", w.LawLevelDesc(), "Personal concealable firearms prohibited"},
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

func TestLawLevelDescExtendsPastNine(t *testing.T) {
	t.Parallel()

	const level9 = "Possession of any weapon outside one's home prohibited"

	if got := (classic.World{LawLevel: 9}).LawLevelDesc(); got != level9 {
		t.Errorf("law 9 desc = %q, want plain level-9 text", got)
	}

	// Book 3's cumulative rule + enforcement throw give guidance past 9, so law
	// level >9 extends level 9 rather than falling back to the generic marker.
	got := (classic.World{LawLevel: 15}).LawLevelDesc()
	if !strings.Contains(got, level9) {
		t.Errorf("law 15 desc should keep level-9 prohibition: %q", got)
	}
	if !strings.Contains(got, "enforcement near-certain") {
		t.Errorf("law 15 desc should note escalating enforcement: %q", got)
	}
	if strings.Contains(got, "(beyond described range)") {
		t.Errorf("law 15 desc should not use the generic marker: %q", got)
	}
}
