// Package dice provides the six-sided dice abstraction used by the world
// generator.
//
// Classic Traveller Book 3 only ever rolls one or two six-sided dice, so the
// interface is a single D6(n) method rather than a general dice-notation
// parser. Every random draw passes through a Roller, which makes seeded
// reproducibility and scripted-test injection both straightforward.
package dice

import (
	"fmt"
	"math/rand"
)

// Roller is the interface every dice-driven procedure depends on.
type Roller interface {
	// D6 rolls n six-sided dice and returns their sum. n must be >= 1.
	D6(n int) int
}

// Seeded is a production roller backed by a seeded *math/rand.Rand.
type Seeded struct {
	rng *rand.Rand
}

// NewSeeded constructs a Seeded roller with the given seed. The same seed
// always yields the same sequence of rolls.
func NewSeeded(seed int64) *Seeded {
	//nolint:gosec // math/rand is intentional; we are not generating crypto material.
	return &Seeded{rng: rand.New(rand.NewSource(seed))}
}

// D6 implements the Roller interface.
func (s *Seeded) D6(n int) int {
	if n < 1 {
		panic(fmt.Sprintf("dice.Seeded: D6 count must be >= 1, got %d", n))
	}
	total := 0
	for range n {
		total += s.rng.Intn(6) + 1
	}
	return total
}

// Scripted is a test roller that yields preset die-roll sums in order.
//
// Each scripted value is the total of one D6(n) call — the natural roll for
// that call — so a book worked-example that reads "a 2D roll of 9" feeds 9
// directly. Values are consumed in call order regardless of n.
type Scripted struct {
	results []int
	idx     int
}

// NewScripted constructs a Scripted roller that returns the supplied values in
// order, one per D6 call.
func NewScripted(results ...int) *Scripted {
	return &Scripted{results: results}
}

// D6 implements the Roller interface. It ignores n (the scripted value is the
// pre-computed sum) and panics if the sequence is exhausted, which always
// indicates a bug in the test or the procedure under test.
func (s *Scripted) D6(n int) int {
	if s.idx >= len(s.results) {
		panic(fmt.Sprintf("dice.Scripted: exhausted on D6(%d)", n))
	}
	v := s.results[s.idx]
	s.idx++
	return v
}

// Fixed is a roller whose every die shows the same face. D6(n) returns n times
// that value. Useful for pinning one variable in a property test.
type Fixed int

// D6 implements the Roller interface.
func (f Fixed) D6(n int) int { return int(f) * n }
