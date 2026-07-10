package worldname_test

import (
	"strings"
	"testing"

	"github.com/philoserf/traveller-worldgen/dice"
	"github.com/philoserf/traveller-worldgen/worldname"
)

func TestGenerateDeterministic(t *testing.T) {
	t.Parallel()

	for seed := int64(0); seed < 200; seed++ {
		a := worldname.Generate(dice.NewSeeded(seed))
		b := worldname.Generate(dice.NewSeeded(seed))
		if a != b {
			t.Fatalf("seed %d: name not deterministic: %q != %q", seed, a, b)
		}
	}
}

func TestGenerateWellFormed(t *testing.T) {
	t.Parallel()

	for seed := int64(0); seed < 500; seed++ {
		name := worldname.Generate(dice.NewSeeded(seed))
		if name == "" {
			t.Fatalf("seed %d: empty name", seed)
		}
		if first := name[0]; first < 'A' || first > 'Z' {
			t.Fatalf("seed %d: name %q not capitalized", seed, name)
		}
		if name[1:] != strings.ToLower(name[1:]) {
			t.Fatalf("seed %d: name %q has non-leading uppercase", seed, name)
		}
	}
}
