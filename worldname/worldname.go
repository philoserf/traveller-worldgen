// Package worldname generates pronounceable world names from a dice.Roller. Name
// generation is edition-independent — it touches no edition's rules or tables —
// so, like dice and ehex, it lives in its own shared package and every edition
// draws from it. Because it consumes dice from the same Roller, a seed reproduces
// an edition's whole world (UWP and name) as one deterministic stream.
package worldname

import (
	"strings"

	"github.com/philoserf/traveller-worldgen/dice"
)

// Name-fragment tables. Lengths are chosen to divide 36 evenly (the range of two
// independent d6 via choose), so selection is unbiased.
var (
	onsets = []string{ // 18 entries
		"b", "c", "d", "f", "g", "h", "j", "k", "l",
		"m", "n", "p", "r", "s", "t", "v", "br", "dr",
	}
	vowels = []string{ // 9 entries
		"a", "e", "i", "o", "u", "ae", "ia", "ou", "ei",
	}
	codas = []string{ // 12 entries
		"n", "r", "s", "l", "m", "x", "th", "rn", "ss", "ld", "k", "ne",
	}
)

// Generate builds a pronounceable name from r: two or three consonant-vowel
// syllables, with an optional coda on the last one. Deterministic in the roller's
// stream.
func Generate(r dice.Roller) string {
	syllables := 2
	if r.D6(1) >= 4 {
		syllables = 3
	}

	var b strings.Builder
	for i := range syllables {
		b.WriteString(choose(r, onsets))
		b.WriteString(choose(r, vowels))
		if i == syllables-1 && r.D6(1) >= 3 {
			b.WriteString(choose(r, codas))
		}
	}

	name := b.String()
	return strings.ToUpper(name[:1]) + name[1:]
}

// choose selects one element of opts using two independent d6 to form a uniform
// index in [0, 36), reduced modulo len(opts). It always consumes exactly two
// dice.
func choose(r dice.Roller, opts []string) string {
	idx := (r.D6(1)-1)*6 + (r.D6(1) - 1)
	return opts[idx%len(opts)]
}
