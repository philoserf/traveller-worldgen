package classic

import (
	"strings"

	"github.com/philoserf/traveller-worldgen/dice"
)

// Name-fragment tables. Lengths are chosen to divide 36 evenly (the range of
// two independent d6 via choose), so selection is unbiased.
var (
	nameOnsets = []string{ // 18 entries
		"b", "c", "d", "f", "g", "h", "j", "k", "l",
		"m", "n", "p", "r", "s", "t", "v", "br", "dr",
	}
	nameVowels = []string{ // 9 entries
		"a", "e", "i", "o", "u", "ae", "ia", "ou", "ei",
	}
	nameCodas = []string{ // 12 entries
		"n", "r", "s", "l", "m", "x", "th", "rn", "ss", "ld", "k", "ne",
	}
)

// generateName builds a pronounceable name from r: two or three consonant-
// vowel syllables, with an optional coda on the last one. Deterministic in the
// roller's stream.
func generateName(r dice.Roller) string {
	syllables := 2
	if r.D6(1) >= 4 {
		syllables = 3
	}

	var b strings.Builder
	for i := range syllables {
		b.WriteString(choose(r, nameOnsets))
		b.WriteString(choose(r, nameVowels))
		if i == syllables-1 && r.D6(1) >= 3 {
			b.WriteString(choose(r, nameCodas))
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
