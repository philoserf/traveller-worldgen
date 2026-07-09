package ehex_test

import (
	"testing"

	"github.com/philoserf/traveller-worldgen/ehex"
)

func TestRoundTrip(t *testing.T) {
	t.Parallel()

	for v := range len(ehex.Digits) {
		c := ehex.Encode(v)
		got, ok := ehex.Decode(c)
		if !ok {
			t.Fatalf("Encode(%d)=%q did not round-trip: Decode failed", v, c)
		}
		if got != v {
			t.Fatalf("round-trip mismatch: Encode(%d)=%q, Decode=%d", v, c, got)
		}
	}
}

func TestSkipsIAndO(t *testing.T) {
	t.Parallel()

	for v := range len(ehex.Digits) {
		if c := ehex.Encode(v); c == 'I' || c == 'O' {
			t.Fatalf("Encode(%d) produced forbidden digit %q", v, c)
		}
	}
	for _, c := range []byte{'I', 'O'} {
		if _, ok := ehex.Decode(c); ok {
			t.Errorf("Decode(%q) should be unrecognized", c)
		}
	}
}

func TestKnownValues(t *testing.T) {
	t.Parallel()

	cases := map[int]byte{0: '0', 9: '9', 10: 'A', 17: 'H', 18: 'J', 22: 'N', 23: 'P'}
	for v, want := range cases {
		if got := ehex.Encode(v); got != want {
			t.Errorf("Encode(%d) = %q, want %q", v, got, want)
		}
	}
}

func TestOutOfRange(t *testing.T) {
	t.Parallel()

	if got := ehex.Encode(-1); got != '?' {
		t.Errorf("Encode(-1) = %q, want '?'", got)
	}
	if got := ehex.Encode(len(ehex.Digits)); got != '?' {
		t.Errorf("Encode(overflow) = %q, want '?'", got)
	}
}
