package dice_test

import (
	"testing"

	"github.com/philoserf/traveller-worldgen/dice"
)

func TestSeededDeterministic(t *testing.T) {
	t.Parallel()

	a := dice.NewSeeded(42)
	b := dice.NewSeeded(42)
	for i := range 1000 {
		ra, rb := a.D6(2), b.D6(2)
		if ra != rb {
			t.Fatalf("roll %d: same seed diverged: %d != %d", i, ra, rb)
		}
	}
}

func TestSeededRange(t *testing.T) {
	t.Parallel()

	r := dice.NewSeeded(7)
	for range 10000 {
		if v := r.D6(1); v < 1 || v > 6 {
			t.Fatalf("D6(1) out of range: %d", v)
		}
		if v := r.D6(2); v < 2 || v > 12 {
			t.Fatalf("D6(2) out of range: %d", v)
		}
		if v := r.D6(3); v < 3 || v > 18 {
			t.Fatalf("D6(3) out of range: %d", v)
		}
	}
}

func TestSeededCoversFullRange(t *testing.T) {
	t.Parallel()

	r := dice.NewSeeded(1)
	seen := map[int]bool{}
	for range 10000 {
		seen[r.D6(1)] = true
	}
	for face := 1; face <= 6; face++ {
		if !seen[face] {
			t.Errorf("D6(1) never produced %d in 10000 rolls", face)
		}
	}
}

func TestSeededPanicsOnBadCount(t *testing.T) {
	t.Parallel()

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic on D6(0)")
		}
	}()
	dice.NewSeeded(1).D6(0)
}

func TestScriptedFIFO(t *testing.T) {
	t.Parallel()

	r := dice.NewScripted(9, 2, 12)
	got := []int{r.D6(2), r.D6(2), r.D6(2)}
	want := []int{9, 2, 12}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("index %d: got %d, want %d", i, got[i], want[i])
		}
	}
}

func TestScriptedPanicsWhenExhausted(t *testing.T) {
	t.Parallel()

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic when scripted roller is exhausted")
		}
	}()
	r := dice.NewScripted(3)
	r.D6(2)
	r.D6(2) // exhausted
}

func TestFixed(t *testing.T) {
	t.Parallel()

	cases := []struct {
		face, n, want int
	}{
		{4, 1, 4},
		{4, 2, 8},
		{6, 3, 18},
		{1, 2, 2},
	}
	for _, c := range cases {
		if got := dice.Fixed(c.face).D6(c.n); got != c.want {
			t.Errorf("Fixed(%d).D6(%d) = %d, want %d", c.face, c.n, got, c.want)
		}
	}
}
