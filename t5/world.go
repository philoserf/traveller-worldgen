// Package t5 generates and describes single Traveller5 (T5) mainworlds: the
// Universal World Profile (StSAHPGL-T) plus Naval/Scout bases, the PBG counts
// (population digit, planetoid belts, gas giants), the generated trade
// classifications, and the Importance/Economic/Cultural extensions, rolled from a
// dice.Roller. Rules and tables are transcribed from
// docs/t5/world-generation.md.
//
// Unlike the classic, megatraveller, and tne packages, T5's core draw is Flux
// (1D-1D) rather than 2D-7, hydrographics keys off atmosphere rather than size,
// and characteristics are hard-capped rather than kept raw.
package t5

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/philoserf/traveller-worldgen/ehex"
)

// World is a single generated T5 mainworld.
type World struct {
	Name            string
	Starport        byte
	Size            int
	Atmosphere      int
	Hydrographics   int
	Population      int // exponent of 10
	Government      int
	LawLevel        int
	TechLevel       int
	NavalBase       bool
	ScoutBase       bool
	PopulationDigit int // PBG first digit: 1-9 when populated, else 0
	PlanetoidBelts  int
	GasGiants       int
	Importance      int      // Ix
	Economic        Economic // Ex
	Cultural        Cultural // Cx
}

// UWP returns the Universal World Profile string, e.g. "A788899-C": starport,
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

// PBG returns the three-digit Population-Belts-Gas-giants code, e.g. "703".
func (w World) PBG() string {
	return fmt.Sprintf("%c%c%c",
		ehex.Encode(w.PopulationDigit), ehex.Encode(w.PlanetoidBelts), ehex.Encode(w.GasGiants))
}

// BaseCode returns the compact library-data base code: "N" (naval), "S" (scout),
// "A" (both), or "" (none). Depot and Way Station are deferred (see the doc).
func (w World) BaseCode() string {
	switch {
	case w.NavalBase && w.ScoutBase:
		return "A"
	case w.NavalBase:
		return "N"
	case w.ScoutBase:
		return "S"
	}
	return ""
}

// Bases returns a human-readable list of the world's bases: "Naval, Scout",
// "Naval", "Scout", or "None".
func (w World) Bases() string {
	var names []string
	if w.NavalBase {
		names = append(names, "Naval")
	}
	if w.ScoutBase {
		names = append(names, "Scout")
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

// LawLevelDesc returns the description for the world's law level. Unlike the
// earlier editions, T5's law table (Book 3 p. 25) is purely descriptive and
// carries no enforcement saving throw, so none is appended.
func (w World) LawLevelDesc() string { return lookup(lawDesc, w.LawLevel) }

// TechLevelDesc returns the era-band description for the world's tech level.
func (w World) TechLevelDesc() string { return lookup(techDesc, w.TechLevel) }

// TradeCodes returns the world's generated trade classifications (As, Ag, …)
// derived from the UWP, in the source table's order.
func (w World) TradeCodes() []string {
	codes := []string{} // non-nil so MarshalJSON emits [] (not null) for a code-less world
	for _, tc := range tradeClassifications {
		if tc.match(w) {
			codes = append(codes, tc.code)
		}
	}
	return codes
}

// Characteristic is one labeled entry of a world's profile: its extended-hex code
// and human-readable description, both derived from the same field.
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

// MarshalJSON emits a JSON object with the starport and base code as strings, the
// convenience UWP/PBG/extension strings, and the raw characteristic, PBG, trade,
// and extension values.
func (w World) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name            string   `json:"name"`
		UWP             string   `json:"uwp"`
		Starport        string   `json:"starport"`
		Size            int      `json:"size"`
		Atmosphere      int      `json:"atmosphere"`
		Hydrographics   int      `json:"hydrographics"`
		Population      int      `json:"population"`
		Government      int      `json:"government"`
		LawLevel        int      `json:"lawLevel"`
		TechLevel       int      `json:"techLevel"`
		NavalBase       bool     `json:"navalBase"`
		ScoutBase       bool     `json:"scoutBase"`
		BaseCode        string   `json:"baseCode"`
		PopulationDigit int      `json:"populationDigit"`
		PlanetoidBelts  int      `json:"planetoidBelts"`
		GasGiants       int      `json:"gasGiants"`
		PBG             string   `json:"pbg"`
		TradeCodes      []string `json:"tradeCodes"`
		Importance      int      `json:"importance"`
		Ix              string   `json:"ix"`
		Economic        Economic `json:"economic"`
		Ex              string   `json:"ex"`
		RU              int      `json:"ru"`
		Cultural        Cultural `json:"cultural"`
		Cx              string   `json:"cx"`
	}{
		Name:            w.Name,
		UWP:             w.UWP(),
		Starport:        string(w.Starport),
		Size:            w.Size,
		Atmosphere:      w.Atmosphere,
		Hydrographics:   w.Hydrographics,
		Population:      w.Population,
		Government:      w.Government,
		LawLevel:        w.LawLevel,
		TechLevel:       w.TechLevel,
		NavalBase:       w.NavalBase,
		ScoutBase:       w.ScoutBase,
		BaseCode:        w.BaseCode(),
		PopulationDigit: w.PopulationDigit,
		PlanetoidBelts:  w.PlanetoidBelts,
		GasGiants:       w.GasGiants,
		PBG:             w.PBG(),
		TradeCodes:      w.TradeCodes(),
		Importance:      w.Importance,
		Ix:              w.ImportanceString(),
		Economic:        w.Economic,
		Ex:              w.Economic.String(),
		RU:              w.Economic.RU(),
		Cultural:        w.Cultural,
		Cx:              w.Cultural.String(),
	})
}

// lookup returns m[k] or beyondRange when k has no entry.
func lookup[K comparable](m map[K]string, k K) string {
	if s, ok := m[k]; ok {
		return s
	}
	return beyondRange
}
