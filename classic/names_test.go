package classic_test

import (
	"strings"
	"testing"

	"github.com/philoserf/traveller-worldgen/classic"
	"github.com/philoserf/traveller-worldgen/dice"
)

func TestNamesDeterministic(t *testing.T) {
	t.Parallel()

	for seed := int64(0); seed < 200; seed++ {
		a := classic.Generate(dice.NewSeeded(seed)).Name
		b := classic.Generate(dice.NewSeeded(seed)).Name
		if a != b {
			t.Fatalf("seed %d: name not deterministic: %q != %q", seed, a, b)
		}
	}
}

func TestNamesWellFormed(t *testing.T) {
	t.Parallel()

	for seed := int64(0); seed < 500; seed++ {
		name := classic.Generate(dice.NewSeeded(seed)).Name
		if name == "" {
			t.Fatalf("seed %d: empty name", seed)
		}
		if first := name[0]; first < 'A' || first > 'Z' {
			t.Fatalf("seed %d: name %q not capitalized", seed, name)
		}
		if strings.ToLower(name)[1:] != name[1:] {
			t.Fatalf("seed %d: name %q has non-leading uppercase", seed, name)
		}
	}
}
