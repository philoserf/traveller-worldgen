// Package classic generates and describes single Classic Traveller (Book 3)
// worlds: the eight-characteristic Universal World Profile plus bases, rolled
// from a dice.Roller. Rules and tables are transcribed from
// docs/classic/world-generation.md.
package classic

import (
	"encoding/json"
	"fmt"

	"github.com/philoserf/traveller-worldgen/ehex"
)

// World is a single generated Classic Traveller world: a name plus the eight
// Book 3 characteristics and the two possible bases.
type World struct {
	Name          string
	Starport      byte
	Size          int
	Atmosphere    int
	Hydrographics int
	Population    int
	Government    int
	LawLevel      int
	TechLevel     int
	NavalBase     bool
	ScoutBase     bool
}

// UWP returns the Universal World Profile string, e.g. "A867845-9": starport,
// size, atmosphere, hydrographics, population, government, law level, a dash,
// then the tech level, each rendered in extended hex.
func (w World) UWP() string {
	return fmt.Sprintf(
		"%c%c%c%c%c%c%c-%c",
		w.Starport,
		ehex.Encode(w.Size),
		ehex.Encode(w.Atmosphere),
		ehex.Encode(w.Hydrographics),
		ehex.Encode(w.Population),
		ehex.Encode(w.Government),
		ehex.Encode(w.LawLevel),
		ehex.Encode(w.TechLevel),
	)
}

// Bases returns a human-readable list of the world's bases: "Naval, Scout",
// "Naval", "Scout", or "None".
func (w World) Bases() string {
	switch {
	case w.NavalBase && w.ScoutBase:
		return "Naval, Scout"
	case w.NavalBase:
		return "Naval"
	case w.ScoutBase:
		return "Scout"
	default:
		return "None"
	}
}

// StarportDesc returns the description for the world's starport.
func (w World) StarportDesc() string { return lookup(starportDesc, w.Starport) }

// SizeDesc returns the description for the world's size.
func (w World) SizeDesc() string { return lookup(sizeDesc, w.Size) }

// AtmosphereDesc returns the description for the world's atmosphere.
func (w World) AtmosphereDesc() string { return lookup(atmoDesc, w.Atmosphere) }

// HydrographicsDesc returns the description for the world's hydrographics.
func (w World) HydrographicsDesc() string { return lookup(hydroDesc, w.Hydrographics) }

// PopulationDesc returns the description for the world's population.
func (w World) PopulationDesc() string { return lookup(popDesc, w.Population) }

// GovernmentDesc returns the description for the world's government.
func (w World) GovernmentDesc() string { return lookup(govDesc, w.Government) }

// LawLevelDesc returns the description for the world's law level. Book 3's
// described scale tops out at 9 ("possession of any weapon outside one's home
// prohibited"), but its note makes each level cumulative and ties the raw level
// to an enforcement throw, so a level above 9 keeps the level-9 weapons
// prohibition while enforcement becomes near-certain.
func (w World) LawLevelDesc() string {
	if w.LawLevel > 9 {
		return lawDesc[9] + "; enforcement near-certain (beyond Book 3's described law level 9)"
	}
	return lookup(lawDesc, w.LawLevel)
}

// TechLevelDesc returns a compact summary of the world's tech-level advances.
func (w World) TechLevelDesc() string {
	if r, ok := techLevels[w.TechLevel]; ok {
		return r.summary()
	}
	return beyondRange
}

// Characteristic is one labeled entry of a world's profile: its extended-hex
// code and human-readable description, both derived from the same field.
type Characteristic struct {
	Label string
	Code  byte
	Desc  string
}

// Profile returns the world's eight characteristics in UWP order. The code and
// description come from the same source, so a renderer never has to re-parse
// the UWP string.
func (w World) Profile() []Characteristic {
	return []Characteristic{
		{"Starport", w.Starport, w.StarportDesc()},
		{"Size", ehex.Encode(w.Size), w.SizeDesc()},
		{"Atmosphere", ehex.Encode(w.Atmosphere), w.AtmosphereDesc()},
		{"Hydrographics", ehex.Encode(w.Hydrographics), w.HydrographicsDesc()},
		{"Population", ehex.Encode(w.Population), w.PopulationDesc()},
		{"Government", ehex.Encode(w.Government), w.GovernmentDesc()},
		{"Law Level", ehex.Encode(w.LawLevel), w.LawLevelDesc()},
		{"Tech Level", ehex.Encode(w.TechLevel), w.TechLevelDesc()},
	}
}

// MarshalJSON emits a JSON object with the starport as a string and the
// convenience UWP field alongside the raw characteristic values.
func (w World) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name          string `json:"name"`
		UWP           string `json:"uwp"`
		Starport      string `json:"starport"`
		Size          int    `json:"size"`
		Atmosphere    int    `json:"atmosphere"`
		Hydrographics int    `json:"hydrographics"`
		Population    int    `json:"population"`
		Government    int    `json:"government"`
		LawLevel      int    `json:"lawLevel"`
		TechLevel     int    `json:"techLevel"`
		NavalBase     bool   `json:"navalBase"`
		ScoutBase     bool   `json:"scoutBase"`
	}{
		Name:          w.Name,
		UWP:           w.UWP(),
		Starport:      string(w.Starport),
		Size:          w.Size,
		Atmosphere:    w.Atmosphere,
		Hydrographics: w.Hydrographics,
		Population:    w.Population,
		Government:    w.Government,
		LawLevel:      w.LawLevel,
		TechLevel:     w.TechLevel,
		NavalBase:     w.NavalBase,
		ScoutBase:     w.ScoutBase,
	})
}

// lookup returns m[k] or beyondRange when k has no entry.
func lookup[K comparable](m map[K]string, k K) string {
	if s, ok := m[k]; ok {
		return s
	}
	return beyondRange
}
