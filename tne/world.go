// Package tne generates and describes single Traveller: The New Era mainworlds:
// the Universal World Profile plus bases, rolled from a dice.Roller. Rules and
// tables are transcribed from docs/tne/world-generation.md. TNE's Basic Mainworld
// Generation closely mirrors MegaTraveller's; unlike MegaTraveller it defines no
// trade-classification or gas-giant step, so a TNE world is a UWP plus bases.
package tne

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/philoserf/traveller-worldgen/ehex"
)

// World is a single generated TNE mainworld: a name plus the eight UWP
// characteristics and the three possible bases.
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
	MilitaryBase  bool // Non-Imperial Military (base code M)
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

// BaseCode returns the compact library-data base code. The imperial bases give
// "N" (naval), "S" (scout), or "A" (both); a non-imperial military base appends
// "M", so the possibilities are "", "N", "S", "A", "M", "NM", "SM", and "AM".
// Code B (Scout Way Station) is deferred (see § Bases).
func (w World) BaseCode() string {
	var code string
	switch {
	case w.NavalBase && w.ScoutBase:
		code = "A"
	case w.NavalBase:
		code = "N"
	case w.ScoutBase:
		code = "S"
	}
	if w.MilitaryBase {
		code += "M"
	}
	return code
}

// Bases returns a human-readable list of the world's bases: e.g. "Naval, Scout",
// "Military", or "None".
func (w World) Bases() string {
	var names []string
	if w.NavalBase {
		names = append(names, "Naval")
	}
	if w.ScoutBase {
		names = append(names, "Scout")
	}
	if w.MilitaryBase {
		names = append(names, "Military")
	}
	if len(names) == 0 {
		return "None"
	}
	return strings.Join(names, ", ")
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

// LawLevelDesc returns the description for the world's law level plus the
// enforcement throw (the 2D saving throw needed to avoid arrest when encountering
// an enforcement agent, equal to the law level). Law level 0 has no laws and
// therefore no enforcement throw.
func (w World) LawLevelDesc() string {
	desc := lookup(lawDesc, w.LawLevel)
	if w.LawLevel == 0 {
		return desc
	}
	return fmt.Sprintf("%s; %d+ to avoid arrest", desc, w.LawLevel)
}

// TechLevelDesc returns the description for the world's tech level.
func (w World) TechLevelDesc() string { return lookup(techDesc, w.TechLevel) }

// Characteristic is one labeled entry of a world's profile: its extended-hex
// code and human-readable description, both derived from the same field.
type Characteristic struct {
	Label string
	Code  byte
	Desc  string
}

// Profile returns the world's eight characteristics in UWP order. The code and
// description come from the same source, so a renderer never has to re-parse the
// UWP string.
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

// MarshalJSON emits a JSON object with the starport and base code as strings and
// the convenience UWP field alongside the raw values.
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
		MilitaryBase  bool   `json:"militaryBase"`
		BaseCode      string `json:"baseCode"`
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
		MilitaryBase:  w.MilitaryBase,
		BaseCode:      w.BaseCode(),
	})
}

// lookup returns m[k] or beyondRange when k has no entry.
func lookup[K comparable](m map[K]string, k K) string {
	if s, ok := m[k]; ok {
		return s
	}
	return beyondRange
}
