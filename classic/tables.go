package classic

import (
	"slices"
	"strings"
)

// All tables in this file are transcribed from docs/classic/world-generation.md, which
// was verified cell-by-cell against Classic Traveller Book 3 (GDW, 1977)
// earlier in this project. Descriptions are condensed from the source prose.

// beyondRange is the description returned for a characteristic value with no
// table entry (a DM-inflated ordinal past the last described row).
const beyondRange = "(beyond described range)"

// starportDesc maps a starport code to its description (Book 3 Starports table).
var starportDesc = map[byte]string{
	'A': "Excellent quality; refined fuel; annual overhaul; shipyard (all ships)",
	'B': "Good quality; refined fuel; annual overhaul; shipyard (non-starships)",
	'C': "Routine quality; unrefined fuel; reasonable repair facilities",
	'D': "Poor quality; unrefined fuel; no repair or shipyard",
	'E': "Frontier installation; bare bedrock; no fuel, facilities, or bases",
	'X': "No starport; no provision for starship landings",
}

// sizeDesc maps a planetary size digit to its description (Book 3 Planetary
// Size table). The 2D-2 roll reaches at most A (10); B and C appear in the
// source table and are retained for completeness.
var sizeDesc = map[int]string{
	0:  "Asteroid/Planetoid Complex",
	1:  "1,000 miles diameter",
	2:  "2,000 miles diameter",
	3:  "3,000 miles diameter",
	4:  "4,000 miles diameter",
	5:  "5,000 miles diameter",
	6:  "6,000 miles diameter",
	7:  "7,000 miles diameter",
	8:  "8,000 miles diameter",
	9:  "9,000 miles diameter",
	10: "10,000 miles diameter",
	11: "11,000 miles diameter",
	12: "12,000 miles diameter",
}

// atmoDesc maps an atmosphere digit to its description (Book 3 Planetary
// Atmosphere table).
var atmoDesc = map[int]string{
	0:  "No atmosphere",
	1:  "Trace",
	2:  "Very thin, tainted",
	3:  "Very thin",
	4:  "Thin, tainted",
	5:  "Thin",
	6:  "Standard",
	7:  "Standard, tainted",
	8:  "Dense",
	9:  "Dense, tainted",
	10: "Exotic",
	11: "Corrosive",
	12: "Insidious",
}

// hydroDesc maps a hydrographic digit to its description (Book 3 Hydrographic
// Percentage table).
var hydroDesc = map[int]string{
	0:  "No free standing water",
	1:  "10%",
	2:  "20%",
	3:  "30%",
	4:  "40%",
	5:  "50%",
	6:  "60%",
	7:  "70%",
	8:  "80%",
	9:  "90%",
	10: "All water; no land masses",
}

// popDesc maps a population digit to its description (Book 3 Population table).
// The digit is an exponent of 10.
var popDesc = map[int]string{
	0:  "No inhabitants",
	1:  "10",
	2:  "100",
	3:  "1,000",
	4:  "10,000",
	5:  "100,000",
	6:  "1,000,000",
	7:  "10,000,000",
	8:  "100,000,000",
	9:  "1,000,000,000",
	10: "10,000,000,000",
}

// govDesc maps a government digit to its description (Book 3 Governmental Type
// table).
var govDesc = map[int]string{
	0:  "No government structure",
	1:  "Company/Corporation",
	2:  "Participating Democracy",
	3:  "Self-Perpetuating Oligarchy",
	4:  "Representative Democracy",
	5:  "Feudal Technocracy",
	6:  "Captive Government",
	7:  "Balkanization",
	8:  "Civil Service Bureaucracy",
	9:  "Impersonal Bureaucracy",
	10: "Charismatic Dictator",
	11: "Non-Charismatic Leader",
	12: "Charismatic Oligarchy",
	13: "Religious Dictatorship",
}

// lawDesc maps a law level to its description (Book 3 Law Levels table). Law
// level is an ordinal scale that can exceed 9; higher values inherit all lower
// prohibitions and are reported via beyondRange.
var lawDesc = map[int]string{
	0: "No weapons laws",
	1: "Body pistols, explosives, and poison gas prohibited",
	2: "Portable energy weapons prohibited",
	3: "Machine guns and automatic rifles prohibited",
	4: "Light assault weapons (SMGs) prohibited",
	5: "Personal concealable firearms prohibited",
	6: "Most firearms (all except shotguns) prohibited",
	7: "Shotguns prohibited",
	8: "Long blades controlled; open public possession prohibited",
	9: "Possession of any weapon outside one's home prohibited",
}

// techRow is one row of the Book 3 Technological Levels tables (weapons/
// computers/communication plus transportation/fuels). Empty fields are the
// source table's blank cells.
type techRow struct {
	Personal      string
	Armor         string
	Special       string
	Computers     string
	Communication string
	Water         string
	Land          string
	Air           string
	Space         string
	Fuels         string
}

// techLevels maps a tech level to its Book 3 table row (levels 0-18). Levels
// with no listed advances (14-16, 18) are intentionally absent.
var techLevels = map[int]techRow{
	0:  {Personal: "Club, cudgel, Spear", Communication: "Runners", Water: "Canoes", Land: "Carts", Fuels: "Muscle"},
	1:  {Personal: "Dagger, pike, Sword", Armor: "Jack", Special: "Catapult", Computers: "Abacus", Communication: "Heliograph", Water: "Galley", Land: "Wagons"},
	2:  {Personal: "Halberd, Broadsword", Special: "Cannon", Fuels: "Wind"},
	3:  {Personal: "Foil, cutlass, Blade, bayonet", Water: "Sailing Ships", Air: "Hot air balloon", Fuels: "Water wheel"},
	4:  {Personal: "Revolver, Shotgun", Armor: "Cloth", Special: "Artillery", Computers: "Adding Machine", Communication: "Telephones", Water: "Steamships", Land: "Trains", Air: "Dirigibles", Fuels: "Coal"},
	5:  {Personal: "Carbine, Rifle, Pistol, SMG", Special: "Sandcasters, Mortars", Computers: "Model/1", Communication: "Radio", Land: "Ground cars", Air: "Fixed wing aircraft", Fuels: "Oil"},
	6:  {Personal: "Auto Rifle", Special: "Missiles, Rocket Launchers", Computers: "Model/1 bis", Communication: "Television", Water: "Submersibles", Land: "ATV, AFV", Air: "Rotary wing aircraft", Fuels: "Fission"},
	7:  {Personal: "Body Pistol", Armor: "Mesh", Special: "Pulse Laser", Computers: "Model/2", Water: "Hovercraft", Land: "Hovercraft", Space: "Non-starships", Fuels: "Solar"},
	8:  {Personal: "Laser Carbine", Special: "Auto-Cannon", Computers: "Model/2 bis", Air: "Air/Raft", Fuels: "Fusion"},
	9:  {Personal: "Laser Rifle", Armor: "Ablat", Special: "Beam Laser", Computers: "Model/3", Space: "Starships"},
	10: {Armor: "Reflec", Computers: "Model/4", Space: "Drives H or less"},
	11: {Computers: "Model/5", Space: "Drives K or less"},
	12: {Computers: "Model/6", Air: "Grav belts", Space: "Drives N or less"},
	13: {Armor: "Battle Dress", Computers: "Model/7", Space: "Drives O or less"},
	14: {Space: "Drives U or less"},
	15: {Space: "All drives"},
	16: {Water: "Matter Transport"},
	17: {Computers: "Artificial Intelligence", Fuels: "Anti-Matter"},
}

// summary joins this row's non-empty advances, in column order, into a compact
// one-line description.
func (r techRow) summary() string {
	fields := []string{
		r.Personal, r.Armor, r.Special, r.Computers, r.Communication,
		r.Water, r.Land, r.Air, r.Space, r.Fuels,
	}
	return strings.Join(slices.DeleteFunc(fields, func(f string) bool { return f == "" }), " · ")
}

// Technological-index DMs (Book 3 Technological Index Matrix). Each is keyed by
// a characteristic's value; a missing key contributes 0. Summed with 1D and
// floored at 0 to yield the tech level.
var (
	starportTechDM = map[byte]int{'A': 6, 'B': 4, 'C': 2, 'X': -4}
	sizeTechDM     = map[int]int{0: 2, 1: 2, 2: 1, 3: 1, 4: 1}
	atmoTechDM     = map[int]int{0: 1, 1: 1, 2: 1, 3: 1, 10: 1, 11: 1, 12: 1, 13: 1, 14: 1}
	hydroTechDM    = map[int]int{9: 1, 10: 2}
	popTechDM      = map[int]int{1: 1, 2: 1, 3: 1, 4: 1, 5: 1, 9: 2, 10: 4}
	govTechDM      = map[int]int{0: 1, 5: 1, 13: -2}
)
